package architecture_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var moduleSchema = map[string]string{
	"analytics":   "analytics",
	"auth":        "auth",
	"resume":      "resume",
	"resume_link": "resume_link",
	"user":        "users",
}

type foreignSchemaViolation struct {
	Line    int
	Schemas []string
	SQL     string
}

func TestNoForeignSchemaSQL(t *testing.T) {
	root := internalRoot(t)
	modulesRoot := filepath.Join(root, "modules")
	modules := discoverModules(t, modulesRoot)
	schemas := make([]string, 0, len(moduleSchema)+1)
	for _, schema := range moduleSchema {
		schemas = append(schemas, schema)
	}
	schemas = append(schemas, "public")
	sort.Strings(schemas)

	for _, module := range modules {
		ownSchema, ok := moduleSchema[module]
		if !ok {
			t.Fatalf("module %s has no schema mapping; add it to moduleSchema", module)
		}
		foreignSchemas := make([]string, 0, len(schemas)-1)
		for _, schema := range schemas {
			if schema != ownSchema {
				foreignSchemas = append(foreignSchemas, schema)
			}
		}

		moduleRoot := filepath.Join(modulesRoot, module)
		err := filepath.WalkDir(moduleRoot, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				if path != moduleRoot && (entry.Name() == "vendor" || strings.HasPrefix(entry.Name(), ".")) {
					return filepath.SkipDir
				}
				return nil
			}
			name := entry.Name()
			if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || strings.HasSuffix(name, "_gen.go") || name == "wire_gen.go" || strings.HasPrefix(name, "zz_generated.") {
				return nil
			}

			source, err := os.ReadFile(path) //nolint:gosec
			if err != nil {
				return err
			}
			violations, err := findForeignSchemaSQL(path, source, foreignSchemas)
			if err != nil {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			for _, violation := range violations {
				for _, schema := range violation.Schemas {
					t.Errorf("%s:%d: module %s must not reference foreign schema %q: %s", relative, violation.Line, module, schema, compactSQL(violation.SQL))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("scan SQL literals in module %s: %v", module, err)
		}
	}
}

func TestForeignSchemaScannerDetectsViolation(t *testing.T) {
	const source = `package sample

const query = "SELECT * FROM resume.resumes WHERE id = $1"
`
	violations, err := findForeignSchemaSQL("sample.go", []byte(source), []string{"resume", "resume_link"})
	if err != nil {
		t.Fatalf("scan synthetic source: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected one violation, got %d: %+v", len(violations), violations)
	}
	if violations[0].Line != 3 {
		t.Fatalf("expected violation on line 3, got %d", violations[0].Line)
	}
	if len(violations[0].Schemas) != 1 || violations[0].Schemas[0] != "resume" {
		t.Fatalf("expected resume schema violation, got %+v", violations[0].Schemas)
	}
	if !strings.Contains(violations[0].SQL, "resume.resumes") {
		t.Fatalf("expected original SQL in violation, got %q", violations[0].SQL)
	}
}

func TestForeignSchemaScannerAllowsOwnSchema(t *testing.T) {
	const source = `package sample

const query = "SELECT * FROM analytics.analytics_views"
`
	violations, err := findForeignSchemaSQL("sample.go", []byte(source), []string{"resume", "resume_link"})
	if err != nil {
		t.Fatalf("scan synthetic source: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %+v", violations)
	}
}

func findForeignSchemaSQL(filename string, source []byte, foreignSchemas []string) ([]foreignSchemaViolation, error) {
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filename, source, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	var violations []foreignSchemaViolation
	ast.Inspect(parsed, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}
		value, err := strconv.Unquote(literal.Value)
		if err != nil {
			return true
		}
		foreign := make([]string, 0, len(foreignSchemas))
		for _, schema := range foreignSchemas {
			if strings.Contains(value, schema+".") {
				foreign = append(foreign, schema)
			}
		}
		if len(foreign) == 0 {
			return true
		}
		violations = append(violations, foreignSchemaViolation{
			Line:    fileSet.Position(literal.Pos()).Line,
			Schemas: foreign,
			SQL:     value,
		})
		return true
	})
	return violations, nil
}

func compactSQL(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) > 240 {
		return value[:240] + "..."
	}
	return value
}
