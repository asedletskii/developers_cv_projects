package architecture_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const modulePath = "developers-cv/backend"

var crossModuleImportAllowlist = map[string]map[string]map[string]struct{}{
	"auth": {
		"application": {
			modulePath + "/internal/modules/user/application/queries": {},
		},
		"infrastructure": {
			modulePath + "/internal/modules/user/application/queries": {},
		},
	},
	"resume_link": {
		"application": {
			modulePath + "/internal/modules/resume/application/queries":     {},
			modulePath + "/internal/modules/resume/domain":                  {},
			modulePath + "/internal/modules/analytics/domain":               {},
			modulePath + "/internal/modules/analytics/application/ports":    {},
			modulePath + "/internal/modules/analytics/application/commands": {},
			modulePath + "/internal/modules/analytics/application/queries":  {},
		},
		"infrastructure": {
			modulePath + "/internal/modules/resume/domain":                                   {},
			modulePath + "/internal/modules/analytics/domain":                                {},
			modulePath + "/internal/modules/analytics/application/ports":                     {},
			modulePath + "/internal/modules/analytics/infrastructure/persistence/gormmodel":  {},
			modulePath + "/internal/modules/analytics/infrastructure/persistence/readers":    {},
			modulePath + "/internal/modules/analytics/infrastructure/persistence/repository": {},
		},
		"transport": {
			modulePath + "/internal/modules/resume/application/queries":     {},
			modulePath + "/internal/modules/resume/domain":                  {},
			modulePath + "/internal/modules/analytics/application/ports":    {},
			modulePath + "/internal/modules/analytics/application/commands": {},
			modulePath + "/internal/modules/analytics/transport/http":       {},
		},
	},
	"analytics": {
		"infrastructure": {
			modulePath + "/internal/modules/resume/domain":      {},
			modulePath + "/internal/modules/resume_link/domain": {},
		},
	},
}

var applicationForbiddenImports = []string{
	"go.uber.org/zap",
	"gorm.io/gorm",
	"github.com/gin-gonic/gin",
	modulePath + "/internal/platform/http",
}

var moduleLayers = []string{"domain", "application", "infrastructure", "transport"}

func TestModuleLayerImports(t *testing.T) {
	root := internalRoot(t)
	modules := discoverModules(t, filepath.Join(root, "modules"))
	if len(modules) == 0 {
		t.Fatal("no modules discovered under internal/modules")
	}

	for _, module := range modules {
		moduleRoot := filepath.Join(root, "modules", module)
		t.Run(module, func(t *testing.T) {
			for _, layer := range moduleLayers {
				dir := filepath.Join(moduleRoot, layer)
				t.Run(layer, func(t *testing.T) {
					if !directoryExists(t, dir) {
						t.Skipf("module %s has no %s layer", module, layer)
					}

					forbidden := forbiddenImportsForLayer(module, layer)
					checkImports(t, dir, forbidden)
				})
			}
		})
	}
}

func TestDomainDoesNotDependOnPlatformOrTestUtilities(t *testing.T) {
	root := internalRoot(t)
	modules := discoverModules(t, filepath.Join(root, "modules"))
	for _, module := range modules {
		dir := filepath.Join(root, "modules", module, "domain")
		t.Run(module, func(t *testing.T) {
			if !directoryExists(t, dir) {
				t.Skipf("module %s has no domain layer", module)
			}
			checkImports(t, dir, []string{
				modulePath + "/internal/platform",
				modulePath + "/internal/testutil",
			})
		})
	}
}

func TestResumeLinkAnalyticsDomainDoesNotDependOnResumeDomain(t *testing.T) {
	root := internalRoot(t)
	resumeDomain := moduleImport("resume", "domain")
	checkImports(t, filepath.Join(root, "modules", "resume_link", "domain"), []string{resumeDomain})
	checkImports(t, filepath.Join(root, "modules", "resume_link", "infrastructure", "analytics"), []string{resumeDomain})
}

func TestCommandsAndQueriesAreIndependent(t *testing.T) {
	root := internalRoot(t)
	modules := discoverModules(t, filepath.Join(root, "modules"))
	for _, module := range modules {
		applicationRoot := filepath.Join(root, "modules", module, "application")
		t.Run(module, func(t *testing.T) {
			commandsDir := filepath.Join(applicationRoot, "commands")
			queriesDir := filepath.Join(applicationRoot, "queries")
			if !directoryExists(t, commandsDir) || !directoryExists(t, queriesDir) {
				t.Skipf("module %s does not contain both commands and queries", module)
			}
			t.Run("commands-do-not-import-queries", func(t *testing.T) {
				checkImports(t, commandsDir, []string{moduleImport(module, "application/queries")})
			})
			t.Run("queries-do-not-import-commands", func(t *testing.T) {
				checkImports(t, queriesDir, []string{moduleImport(module, "application/commands")})
			})
		})
	}
}

func TestCrossModuleImportsFollowAllowlist(t *testing.T) {
	root := internalRoot(t)
	modules := discoverModules(t, filepath.Join(root, "modules"))
	moduleSet := make(map[string]struct{}, len(modules))
	for _, module := range modules {
		moduleSet[module] = struct{}{}
	}

	for _, sourceModule := range modules {
		moduleRoot := filepath.Join(root, "modules", sourceModule)
		t.Run(sourceModule, func(t *testing.T) {
			err := walkGoSources(moduleRoot, func(path string, imports []string) {
				sourceLayer := sourceLayerForPath(moduleRoot, path)
				for _, importPath := range imports {
					targetModule, targetLayer, ok := importedModuleLayer(importPath, moduleSet)
					if !ok || targetModule == sourceModule {
						continue
					}
					if reason := crossModuleLayerViolation(sourceLayer, targetLayer); reason != "" {
						t.Errorf("module %s: %s imports %s from module %s: %s", sourceModule, path, importPath, targetModule, reason)
						continue
					}
					if _, allowed := crossModuleImportAllowlist[sourceModule][sourceLayer][importPath]; !allowed {
						t.Errorf("module %s layer %s: %s imports cross-module package %s; add this exact package to the layer allowlist only after architectural review", sourceModule, sourceLayer, path, importPath)
					}
				}
			})
			if err != nil {
				t.Fatalf("scan module %s: %v", sourceModule, err)
			}
		})
	}
}

func TestCrossModuleLayerRules(t *testing.T) {
	tests := []struct {
		name          string
		sourceLayer   string
		targetLayer   string
		wantViolation bool
	}{
		{name: "domain may depend on domain", sourceLayer: "domain", targetLayer: "domain"},
		{name: "domain may not depend on application", sourceLayer: "domain", targetLayer: "application", wantViolation: true},
		{name: "domain may not depend on infrastructure", sourceLayer: "domain", targetLayer: "infrastructure", wantViolation: true},
		{name: "application may depend on domain", sourceLayer: "application", targetLayer: "domain"},
		{name: "application may not depend on infrastructure", sourceLayer: "application", targetLayer: "infrastructure", wantViolation: true},
		{name: "application may not depend on transport", sourceLayer: "application", targetLayer: "transport", wantViolation: true},
		{name: "infrastructure may not depend on transport", sourceLayer: "infrastructure", targetLayer: "transport", wantViolation: true},
		{name: "transport may not depend on infrastructure", sourceLayer: "transport", targetLayer: "infrastructure", wantViolation: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotViolation := crossModuleLayerViolation(tc.sourceLayer, tc.targetLayer) != ""
			if gotViolation != tc.wantViolation {
				t.Fatalf("crossModuleLayerViolation(%q, %q) returned violation=%t, want %t", tc.sourceLayer, tc.targetLayer, gotViolation, tc.wantViolation)
			}
		})
	}
}

func TestAuthPersistenceMappersDoNotDependOnRepositories(t *testing.T) {
	root := internalRoot(t)
	checkImports(t, filepath.Join(root, "modules", "auth", "infrastructure", "persistence", "mappers"), []string{
		modulePath + "/internal/modules/auth/infrastructure/persistence/repository",
	})
}

func forbiddenImportsForLayer(module, layer string) []string {
	switch layer {
	case "domain":
		return []string{
			moduleImport(module, "application"),
			moduleImport(module, "infrastructure"),
			moduleImport(module, "transport"),
			"github.com/gin-gonic/gin",
			"gorm.io/gorm",
		}
	case "application":
		return append([]string{
			moduleImport(module, "infrastructure"),
			moduleImport(module, "transport"),
		}, applicationForbiddenImports...)
	case "infrastructure":
		return []string{moduleImport(module, "transport")}
	case "transport":
		return []string{moduleImport(module, "infrastructure")}
	default:
		return nil
	}
}

func discoverModules(t *testing.T, modulesRoot string) []string {
	t.Helper()
	entries, err := os.ReadDir(modulesRoot)
	if err != nil {
		t.Fatalf("read modules directory %q: %v", modulesRoot, err)
	}

	var modules []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		module := entry.Name()
		moduleRoot := filepath.Join(modulesRoot, module)
		hasLayer := false
		for _, layer := range moduleLayers {
			if directoryExists(t, filepath.Join(moduleRoot, layer)) {
				hasLayer = true
				break
			}
		}
		if !hasLayer {
			t.Errorf("directory %q is under internal/modules but contains no recognized layer (%s)", moduleRoot, strings.Join(moduleLayers, ", "))
			continue
		}
		modules = append(modules, module)
	}
	sort.Strings(modules)
	return modules
}

func internalRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func moduleImport(module, suffix string) string {
	return modulePath + "/internal/modules/" + module + "/" + suffix
}

func directoryExists(t *testing.T, dir string) bool {
	t.Helper()
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		t.Fatalf("stat directory %q: %v", dir, err)
	}
	if !stat.IsDir() {
		t.Fatalf("expected layer path %q to be a directory", dir)
	}
	return true
}

func checkImports(t *testing.T, dir string, forbidden []string) {
	t.Helper()
	err := walkGoSources(dir, func(path string, imports []string) {
		for _, importPath := range imports {
			for _, banned := range forbidden {
				if importPath == banned || strings.HasPrefix(importPath, banned+"/") {
					t.Errorf("file %s imports forbidden package %s (rule: %s is not allowed in %s)", path, importPath, banned, dir)
				}
			}
		}
	})
	if err != nil {
		t.Fatalf("scan imports in %s: %v", dir, err)
	}
}

func walkGoSources(dir string, visit func(path string, imports []string)) error {
	return filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != dir && (entry.Name() == "vendor" || strings.HasPrefix(entry.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || strings.HasSuffix(name, "_gen.go") || name == "wire_gen.go" || strings.HasPrefix(name, "zz_generated.") {
			return nil
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		imports := make([]string, 0, len(parsed.Imports))
		for _, spec := range parsed.Imports {
			value, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			imports = append(imports, value)
		}
		visit(path, imports)
		return nil
	})
}

func sourceLayerForPath(moduleRoot, sourcePath string) string {
	relative, err := filepath.Rel(moduleRoot, sourcePath)
	if err != nil {
		return ""
	}
	layer, _, _ := strings.Cut(filepath.ToSlash(relative), "/")
	return layer
}

func crossModuleLayerViolation(sourceLayer, targetLayer string) string {
	switch sourceLayer {
	case "domain":
		if targetLayer != "domain" {
			return "domain code may depend on another module's domain only, never its application or adapters"
		}
	case "application":
		if targetLayer == "infrastructure" || targetLayer == "transport" {
			return "application code must not depend on infrastructure or transport layers"
		}
	case "infrastructure":
		if targetLayer == "transport" {
			return "infrastructure code must not depend on transport layers"
		}
	case "transport":
		if targetLayer == "infrastructure" {
			return "transport code must not depend on infrastructure layers"
		}
	}
	return ""
}

func importedModuleLayer(importPath string, modules map[string]struct{}) (string, string, bool) {
	const prefix = modulePath + "/internal/modules/"
	if !strings.HasPrefix(importPath, prefix) {
		return "", "", false
	}
	remainder := strings.TrimPrefix(importPath, prefix)
	module, packagePath, _ := strings.Cut(remainder, "/")
	_, exists := modules[module]
	if !exists || packagePath == "" {
		return module, "", exists
	}
	layer, _, _ := strings.Cut(packagePath, "/")
	for _, knownLayer := range moduleLayers {
		if layer == knownLayer {
			return module, layer, true
		}
	}
	return module, "", true
}
