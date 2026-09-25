package architecture_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const modulePath = "developers-cv/backend"

func TestResumeLayerImports(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), ".."))

	cases := []struct {
		name      string
		dir       string
		forbidden []string
	}{
		{name: "queries", dir: filepath.Join(root, "modules", "resume", "application", "queries"), forbidden: []string{modulePath + "/internal/modules/resume/application/commands", modulePath + "/internal/modules/resume/infrastructure/persistence/repository", modulePath + "/internal/modules/resume/infrastructure/persistence/readers", modulePath + "/internal/modules/resume/infrastructure/persistence/gormmodel"}},
		{name: "commands", dir: filepath.Join(root, "modules", "resume", "application", "commands"), forbidden: []string{modulePath + "/internal/modules/resume/application/queries", modulePath + "/internal/modules/resume/infrastructure/persistence/readers"}},
		{name: "domain", dir: filepath.Join(root, "modules", "resume", "domain"), forbidden: []string{modulePath + "/internal/modules/resume/application", modulePath + "/internal/modules/resume/infrastructure", modulePath + "/internal/modules/resume/transport"}},
		{name: "transport", dir: filepath.Join(root, "modules", "resume", "transport"), forbidden: []string{modulePath + "/internal/modules/resume/infrastructure"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { checkImports(t, tc.dir, tc.forbidden) })
	}
}

func checkImports(t *testing.T, dir string, forbidden []string) {
	t.Helper()
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || strings.HasSuffix(name, "_gen.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imp := range file.Imports {
			value := strings.Trim(imp.Path.Value, `"`)
			for _, banned := range forbidden {
				if value == banned || strings.HasPrefix(value, banned+"/") {
					t.Errorf("%s imports forbidden package %s", path, value)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
