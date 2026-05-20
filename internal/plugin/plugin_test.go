package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePlugin_Route(t *testing.T) {
	yaml := `
apiVersion: m3tal/v1
kind: Route
metadata:
  name: test-route
  description: "A test route"
spec:
  service: myapp
  domain: app.local
  port: 8080
  entrypoints: websecure
  network: proxy
`
	p, err := ParsePlugin([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePlugin failed: %v", err)
	}

	if err := p.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if p.Kind != KindRoute {
		t.Errorf("expected kind %q, got %q", KindRoute, p.Kind)
	}

	rp, err := p.AsRoute()
	if err != nil {
		t.Fatalf("AsRoute failed: %v", err)
	}

	if rp.Service != "myapp" {
		t.Errorf("expected service 'myapp', got %q", rp.Service)
	}
	if rp.Domain != "app.local" {
		t.Errorf("expected domain 'app.local', got %q", rp.Domain)
	}
	if rp.Port != 8080 {
		t.Errorf("expected port 8080, got %d", rp.Port)
	}
}

func TestParsePlugin_Stack(t *testing.T) {
	yaml := `
apiVersion: m3tal/v1
kind: Stack
metadata:
  name: network
  description: "Docker network definitions"
spec:
  composePath: network-compose.yml
  priority: 1
  category: infra
`
	p, err := ParsePlugin([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePlugin failed: %v", err)
	}

	if err := p.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	sp, err := p.AsStack()
	if err != nil {
		t.Fatalf("AsStack failed: %v", err)
	}

	if sp.ComposePath != "network-compose.yml" {
		t.Errorf("expected composePath 'network-compose.yml', got %q", sp.ComposePath)
	}
	if sp.Priority != 1 {
		t.Errorf("expected priority 1, got %d", sp.Priority)
	}
	if sp.Category != "infra" {
		t.Errorf("expected category 'infra', got %q", sp.Category)
	}
}

func TestParsePlugin_Middleware(t *testing.T) {
	yaml := `
apiVersion: m3tal/v1
kind: Middleware
metadata:
  name: auth
spec:
  name: m3tal-auth
  type: basicauth
  config:
    usersFile: /etc/traefik/users
`
	p, err := ParsePlugin([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePlugin failed: %v", err)
	}

	if err := p.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	mp, err := p.AsMiddleware()
	if err != nil {
		t.Fatalf("AsMiddleware failed: %v", err)
	}

	if mp.Name != "m3tal-auth" {
		t.Errorf("expected name 'm3tal-auth', got %q", mp.Name)
	}
	if mp.Type != "basicauth" {
		t.Errorf("expected type 'basicauth', got %q", mp.Type)
	}
	if mp.Config["usersFile"] != "/etc/traefik/users" {
		t.Errorf("expected usersFile '/etc/traefik/users', got %q", mp.Config["usersFile"])
	}
}

func TestValidate_InvalidAPIVersion(t *testing.T) {
	p := &Plugin{
		APIVersion: "v2",
		Kind:       KindRoute,
		Metadata:   PluginMetadata{Name: "test"},
		Spec:       map[string]any{"service": "x", "domain": "x", "port": 80},
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for invalid apiVersion")
	}
}

func TestValidate_InvalidKind(t *testing.T) {
	p := &Plugin{
		APIVersion: APIVersion,
		Kind:       "Widget",
		Metadata:   PluginMetadata{Name: "test"},
		Spec:       map[string]any{},
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for invalid kind")
	}
}

func TestValidate_MissingName(t *testing.T) {
	p := &Plugin{
		APIVersion: APIVersion,
		Kind:       KindRoute,
		Metadata:   PluginMetadata{},
		Spec:       map[string]any{"service": "x", "domain": "x", "port": 80},
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for missing name")
	}
}

func TestValidate_MissingRequiredSpec(t *testing.T) {
	p := &Plugin{
		APIVersion: APIVersion,
		Kind:       KindRoute,
		Metadata:   PluginMetadata{Name: "test"},
		Spec:       map[string]any{"service": "x"},
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for missing required spec fields")
	}
}

func TestScanDir(t *testing.T) {
	// Create temp directory with test plugins
	dir := t.TempDir()

	routeDir := filepath.Join(dir, "routes")
	if err := os.MkdirAll(routeDir, 0755); err != nil {
		t.Fatal(err)
	}

	routeYAML := `
apiVersion: m3tal/v1
kind: Route
metadata:
  name: scan-test
spec:
  service: test-svc
  domain: test.local
  port: 9090
`
	if err := os.WriteFile(filepath.Join(routeDir, "test.yml"), []byte(routeYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Also add an invalid file that should be skipped
	if err := os.WriteFile(filepath.Join(routeDir, "bad.yml"), []byte("not: valid: plugin"), 0644); err != nil {
		t.Fatal(err)
	}

	// Add a non-YAML file that should be ignored
	if err := os.WriteFile(filepath.Join(routeDir, "readme.txt"), []byte("ignore me"), 0644); err != nil {
		t.Fatal(err)
	}

	plugins, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("ScanDir failed: %v", err)
	}

	if len(plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(plugins))
	}

	if plugins[0].Metadata.Name != "scan-test" {
		t.Errorf("expected name 'scan-test', got %q", plugins[0].Metadata.Name)
	}
}

func TestScanDir_NonExistent(t *testing.T) {
	plugins, err := ScanDir("/nonexistent/path")
	if err != nil {
		t.Fatalf("ScanDir should not error on nonexistent dir: %v", err)
	}
	if plugins != nil {
		t.Errorf("expected nil plugins for nonexistent dir, got %d", len(plugins))
	}
}

func TestLoadAll_Deduplication(t *testing.T) {
	// System dir with a plugin
	sysDir := t.TempDir()
	sysYAML := `
apiVersion: m3tal/v1
kind: Route
metadata:
  name: shared-route
spec:
  service: sys-svc
  domain: sys.local
  port: 8080
`
	if err := os.WriteFile(filepath.Join(sysDir, "route.yml"), []byte(sysYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// User dir with same name but different spec
	userDir := t.TempDir()
	userYAML := `
apiVersion: m3tal/v1
kind: Route
metadata:
  name: shared-route
spec:
  service: user-svc
  domain: user.local
  port: 9090
`
	if err := os.WriteFile(filepath.Join(userDir, "route.yml"), []byte(userYAML), 0644); err != nil {
		t.Fatal(err)
	}

	reg, err := LoadAll(sysDir, userDir)
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(reg.Routes) != 1 {
		t.Fatalf("expected 1 route after dedup, got %d", len(reg.Routes))
	}

	// User should win
	if reg.Routes[0].Service != "user-svc" {
		t.Errorf("expected user override 'user-svc', got %q", reg.Routes[0].Service)
	}
}

func TestRegistry_StackPrioritySorting(t *testing.T) {
	plugins := []Plugin{
		{
			APIVersion: APIVersion,
			Kind:       KindStack,
			Metadata:   PluginMetadata{Name: "app"},
			Spec:       map[string]any{"composePath": "app-compose.yml", "priority": 10},
		},
		{
			APIVersion: APIVersion,
			Kind:       KindStack,
			Metadata:   PluginMetadata{Name: "network"},
			Spec:       map[string]any{"composePath": "network-compose.yml", "priority": 1},
		},
		{
			APIVersion: APIVersion,
			Kind:       KindStack,
			Metadata:   PluginMetadata{Name: "routing"},
			Spec:       map[string]any{"composePath": "routing-compose.yml", "priority": 2},
		},
	}

	reg, err := BuildRegistry(plugins)
	if err != nil {
		t.Fatalf("BuildRegistry failed: %v", err)
	}

	if len(reg.Stacks) != 3 {
		t.Fatalf("expected 3 stacks, got %d", len(reg.Stacks))
	}

	expected := []string{"network", "routing", "app"}
	for i, name := range expected {
		if reg.Stacks[i].Metadata.Name != name {
			t.Errorf("stack[%d]: expected %q, got %q", i, name, reg.Stacks[i].Metadata.Name)
		}
	}
}

func TestRegistry_GetMethods(t *testing.T) {
	plugins := []Plugin{
		{
			APIVersion: APIVersion,
			Kind:       KindRoute,
			Metadata:   PluginMetadata{Name: "my-route"},
			Spec:       map[string]any{"service": "x", "domain": "x.local", "port": 80},
		},
		{
			APIVersion: APIVersion,
			Kind:       KindStack,
			Metadata:   PluginMetadata{Name: "my-stack"},
			Spec:       map[string]any{"composePath": "x-compose.yml"},
		},
		{
			APIVersion: APIVersion,
			Kind:       KindMiddleware,
			Metadata:   PluginMetadata{Name: "my-mw"},
			Spec:       map[string]any{"name": "test", "type": "basicauth"},
		},
	}

	reg, err := BuildRegistry(plugins)
	if err != nil {
		t.Fatalf("BuildRegistry failed: %v", err)
	}

	if reg.GetRoute("my-route") == nil {
		t.Error("GetRoute('my-route') returned nil")
	}
	if reg.GetRoute("nonexistent") != nil {
		t.Error("GetRoute('nonexistent') should return nil")
	}

	if reg.GetStack("my-stack") == nil {
		t.Error("GetStack('my-stack') returned nil")
	}
	if reg.GetMiddleware("my-mw") == nil {
		t.Error("GetMiddleware('my-mw') returned nil")
	}

	if reg.Count() != 3 {
		t.Errorf("expected count 3, got %d", reg.Count())
	}
}
