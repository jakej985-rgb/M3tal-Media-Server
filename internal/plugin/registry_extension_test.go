package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMatchService(t *testing.T) {
	reg := &Registry{
		Routes: []RoutePlugin{
			{
				Metadata: PluginMetadata{Name: "plex-route"},
				Service:  "plex",
				Domain:   "plex.example.com",
				Port:     32400,
				Enabled:  true,
			},
			{
				Metadata: PluginMetadata{Name: "jellyfin-route"},
				Service:  "jellyfin",
				Domain:   "jellyfin.example.com",
				Port:     8096,
				Enabled:  true,
			},
		},
	}

	// 1. Match by service name
	p := reg.MatchService("plex", "", nil)
	if p == nil || p.Metadata.Name != "plex-route" {
		t.Fatalf("expected plex-route, got %v", p)
	}

	// 2. Match by image name
	p = reg.MatchService("my-jellyfin", "jellyfin/jellyfin:latest", nil)
	if p == nil || p.Metadata.Name != "jellyfin-route" {
		t.Fatalf("expected jellyfin-route, got %v", p)
	}

	// 3. Match by labels
	p = reg.MatchService("unknown", "", map[string]string{"traefik.http.routers.plex.rule": "true"})
	if p == nil || p.Metadata.Name != "plex-route" {
		t.Fatalf("expected plex-route, got %v", p)
	}
}

func TestGenerateTraefikConfig(t *testing.T) {
	reg := &Registry{
		Routes: []RoutePlugin{
			{
				Metadata:    PluginMetadata{Name: "plex-route"},
				Service:     "plex",
				Domain:      "plex.example.com",
				Port:        32400,
				Entrypoints: "web",
				Enabled:     true,
			},
		},
		Middlewares: []MiddlewarePlugin{
			{
				Metadata: PluginMetadata{Name: "auth-preset"},
				Name:     "m3tal-auth",
				Type:     "basicauth",
				Config: map[string]string{
					"usersFile": "/etc/traefik/users",
				},
				Enabled: true,
			},
		},
	}

	cfg, err := reg.GenerateTraefikConfig()
	if err != nil {
		t.Fatalf("failed to generate Traefik config: %v", err)
	}

	yamlStr := string(cfg)
	if !strings.Contains(yamlStr, "plex-route") {
		t.Errorf("expected plex-route in config, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "m3tal-auth") {
		t.Errorf("expected m3tal-auth in config, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "/etc/traefik/users") {
		t.Errorf("expected usersFile path in config, got:\n%s", yamlStr)
	}
}

func TestEnableDisablePlugin(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "m3tal-plugin-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test-plugin.yml")
	err = os.WriteFile(filePath, []byte("apiVersion: m3tal/v1\nkind: Route\nmetadata:\n  name: test\nspec:\n  service: test\n  domain: test\n  port: 80"), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// 1. Disable
	newPath, err := DisablePlugin(filePath)
	if err != nil {
		t.Fatalf("failed to disable: %v", err)
	}
	if !strings.HasSuffix(newPath, ".disabled") {
		t.Errorf("expected .disabled suffix, got %s", newPath)
	}

	// 2. Enable
	restoredPath, err := EnablePlugin(newPath)
	if err != nil {
		t.Fatalf("failed to enable: %v", err)
	}
	if strings.HasSuffix(restoredPath, ".disabled") {
		t.Errorf("expected no .disabled suffix, got %s", restoredPath)
	}
}

func TestParameterize(t *testing.T) {
	tpl := "port: ${PORT}\nvolume: ${VOLUME}"
	vars := map[string]string{
		"PORT":   "8080",
		"VOLUME": "/mnt/data",
	}

	res := Parameterize(tpl, vars)
	expected := "port: 8080\nvolume: /mnt/data"
	if res != expected {
		t.Errorf("expected %q, got %q", expected, res)
	}
}
