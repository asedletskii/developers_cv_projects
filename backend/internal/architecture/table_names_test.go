package architecture_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestGORMModelsUseOwnedSchemas(t *testing.T) {
	root := internalRoot(t)
	modules := discoverModules(t, filepath.Join(root, "modules"))
	total := 0

	for _, module := range modules {
		expectedSchema, ok := moduleSchema[module]
		if !ok {
			t.Fatalf("module %s has no schema mapping", module)
		}
		dir := filepath.Join(root, "modules", module, "infrastructure", "persistence", "gormmodel")
		if !directoryExists(t, dir) {
			continue
		}

		schema, err := readGORMModelSchemaConstant(dir)
		if err != nil {
			t.Fatalf("module %s: %v", module, err)
		}
		if schema != expectedSchema {
			t.Errorf("module %s: gormmodel schema constant is %q, want %q", module, schema, expectedSchema)
		}

		count := 0
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("module %s: read gormmodel directory: %v", module, err)
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			path := filepath.Join(dir, name)
			fileSet := token.NewFileSet()
			parsed, err := parser.ParseFile(fileSet, path, nil, 0)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			for _, declaration := range parsed.Decls {
				function, ok := declaration.(*ast.FuncDecl)
				if !ok || function.Name.Name != "TableName" || function.Recv == nil {
					continue
				}
				table, err := qualifiedGORMTableName(function, schema)
				position := fileSet.Position(function.Pos())
				if err != nil {
					t.Errorf("%s: %v", position, err)
					continue
				}
				prefix, tableName, ok := strings.Cut(table, ".")
				if !ok || prefix == "" || tableName == "" {
					t.Errorf("%s: TableName() must return schema.table, got %q", position, table)
					continue
				}
				if prefix != expectedSchema {
					t.Errorf("%s: TableName() returns schema %q, want %q", position, prefix, expectedSchema)
					continue
				}
				count++
			}
		}
		if count == 0 {
			t.Errorf("module %s has gormmodel package but no TableName() methods", module)
		}
		total += count
	}
	if total == 0 {
		t.Fatal("no GORM TableName() methods were checked")
	}
}

func readGORMModelSchemaConstant(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read gormmodel directory: %w", err)
	}
	value := ""
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", path, err)
		}
		for _, declaration := range parsed.Decls {
			generic, ok := declaration.(*ast.GenDecl)
			if !ok || generic.Tok != token.CONST {
				continue
			}
			for _, specification := range generic.Specs {
				valueSpec, ok := specification.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, name := range valueSpec.Names {
					if name.Name != "schema" || i >= len(valueSpec.Values) {
						continue
					}
					literal, ok := valueSpec.Values[i].(*ast.BasicLit)
					if !ok || literal.Kind != token.STRING {
						return "", fmt.Errorf("%s: schema constant must be a string literal", path)
					}
					schema, err := strconv.Unquote(literal.Value)
					if err != nil {
						return "", fmt.Errorf("%s: unquote schema constant: %w", path, err)
					}
					if value != "" && value != schema {
						return "", fmt.Errorf("%s: conflicting schema constants %q and %q", dir, value, schema)
					}
					value = schema
				}
			}
		}
	}
	if value == "" {
		return "", fmt.Errorf("%s: schema constant not found", dir)
	}
	return value, nil
}

func qualifiedGORMTableName(function *ast.FuncDecl, schema string) (string, error) {
	if function.Body == nil || len(function.Body.List) != 1 {
		return "", fmt.Errorf("TableName() must contain exactly one return statement")
	}
	returnStatement, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returnStatement.Results) != 1 {
		return "", fmt.Errorf("TableName() must contain exactly one return value")
	}

	switch expression := returnStatement.Results[0].(type) {
	case *ast.BasicLit:
		if expression.Kind != token.STRING {
			return "", fmt.Errorf("TableName() must return a string literal")
		}
		return strconv.Unquote(expression.Value)
	case *ast.BinaryExpr:
		if expression.Op != token.ADD {
			return "", fmt.Errorf("TableName() must concatenate schema with a table name")
		}
		schemaIdentifier, ok := expression.X.(*ast.Ident)
		if !ok || schemaIdentifier.Name != "schema" {
			return "", fmt.Errorf("TableName() must start with the package schema constant")
		}
		suffixLiteral, ok := expression.Y.(*ast.BasicLit)
		if !ok || suffixLiteral.Kind != token.STRING {
			return "", fmt.Errorf("TableName() must concatenate schema with a string literal")
		}
		suffix, err := strconv.Unquote(suffixLiteral.Value)
		if err != nil {
			return "", fmt.Errorf("unquote table suffix: %w", err)
		}
		if !strings.HasPrefix(suffix, ".") || len(suffix) == 1 {
			return "", fmt.Errorf("TableName() suffix must start with a dot and a table name, got %q", suffix)
		}
		return schema + suffix, nil
	default:
		return "", fmt.Errorf("TableName() must return a string literal or schema + literal")
	}
}
