package architecture_test

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestAuthLayerImports(t *testing.T) {
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
		{
			name: "domain",
			dir:  filepath.Join(root, "modules", "auth", "domain"),
			forbidden: []string{
				modulePath + "/internal/modules/auth/application",
				modulePath + "/internal/modules/auth/infrastructure",
				modulePath + "/internal/modules/auth/transport",
			},
		},
		{
			name: "application",
			dir:  filepath.Join(root, "modules", "auth", "application"),
			forbidden: []string{
				modulePath + "/internal/modules/auth/infrastructure",
				modulePath + "/internal/modules/auth/transport",
			},
		},
		{
			name: "transport",
			dir:  filepath.Join(root, "modules", "auth", "transport"),
			forbidden: []string{
				modulePath + "/internal/modules/auth/infrastructure",
			},
		},
		{
			name: "infrastructure",
			dir:  filepath.Join(root, "modules", "auth", "infrastructure"),
			forbidden: []string{
				modulePath + "/internal/modules/auth/transport",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { checkImports(t, tc.dir, tc.forbidden) })
	}
}

func TestAuthPersistenceMappersDoNotDependOnRepositories(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
	checkImports(t, filepath.Join(root, "modules", "auth", "infrastructure", "persistence", "mappers"), []string{
		modulePath + "/internal/modules/auth/infrastructure/persistence/repository",
	})
}
