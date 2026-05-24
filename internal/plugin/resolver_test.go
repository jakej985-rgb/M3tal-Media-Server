package plugin

import (
	"strings"
	"testing"
)

func TestResolveInstallOrder_ValidDAG(t *testing.T) {
	catalog := []CatalogItem{
		{
			Name: "A",
			Kind: "Route",
			Dependencies: []Dependency{
				{Name: "B", Kind: "Traefik", Required: true, AutoInstall: true},
				{Name: "C", Kind: "Middleware", Required: true, AutoInstall: true},
			},
		},
		{
			Name: "B",
			Kind: "Traefik",
			Dependencies: []Dependency{
				{Name: "D", Kind: "Middleware", Required: true, AutoInstall: true},
			},
		},
		{
			Name: "C",
			Kind: "Middleware",
			Dependencies: []Dependency{
				{Name: "D", Kind: "Middleware", Required: true, AutoInstall: true},
			},
		},
		{
			Name:         "D",
			Kind:         "Middleware",
			Dependencies: nil,
		},
	}

	installed := make(map[string]bool)
	target := catalog[0] // A

	order, err := ResolveInstallOrder(target, catalog, installed)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Correct topological order: D must be before B and C, and B & C must be before A.
	if len(order) != 4 {
		t.Fatalf("expected 4 items in resolve order, got: %d", len(order))
	}

	pos := make(map[string]int)
	for i, item := range order {
		pos[strings.ToLower(item.Name)] = i
	}

	if pos["d"] > pos["b"] {
		t.Error("expected D to be installed before B")
	}
	if pos["d"] > pos["c"] {
		t.Error("expected D to be installed before C")
	}
	if pos["b"] > pos["a"] {
		t.Error("expected B to be installed before A")
	}
	if pos["c"] > pos["a"] {
		t.Error("expected C to be installed before A")
	}
}

func TestResolveInstallOrder_Circular(t *testing.T) {
	catalog := []CatalogItem{
		{
			Name: "A",
			Kind: "Route",
			Dependencies: []Dependency{
				{Name: "B", Kind: "Traefik", Required: true},
			},
		},
		{
			Name: "B",
			Kind: "Traefik",
			Dependencies: []Dependency{
				{Name: "C", Kind: "Middleware", Required: true},
			},
		},
		{
			Name: "C",
			Kind: "Middleware",
			Dependencies: []Dependency{
				{Name: "A", Kind: "Route", Required: true},
			},
		},
	}

	installed := make(map[string]bool)
	target := catalog[0] // A

	_, err := ResolveInstallOrder(target, catalog, installed)
	if err == nil {
		t.Fatal("expected circular dependency error, got none")
	}

	if !strings.Contains(err.Error(), "circular dependency detected") {
		t.Errorf("expected circular dependency error message, got: %v", err)
	}
}

func TestResolveInstallOrder_MissingRequired(t *testing.T) {
	catalog := []CatalogItem{
		{
			Name: "A",
			Kind: "Route",
			Dependencies: []Dependency{
				{Name: "B", Kind: "Traefik", Required: true},
			},
		},
	}

	installed := make(map[string]bool)
	target := catalog[0] // A

	_, err := ResolveInstallOrder(target, catalog, installed)
	if err == nil {
		t.Fatal("expected missing required dependency error, got none")
	}

	if !strings.Contains(err.Error(), "missing required dependency") {
		t.Errorf("expected missing required dependency error, got: %v", err)
	}
}

func TestResolveInstallOrder_MissingOptional(t *testing.T) {
	catalog := []CatalogItem{
		{
			Name: "A",
			Kind: "Route",
			Dependencies: []Dependency{
				{Name: "B", Kind: "Traefik", Required: false},
			},
		},
	}

	installed := make(map[string]bool)
	target := catalog[0] // A

	order, err := ResolveInstallOrder(target, catalog, installed)
	if err != nil {
		t.Fatalf("expected no error for optional missing dependency, got: %v", err)
	}

	if len(order) != 1 || order[0].Name != "A" {
		t.Fatalf("expected only A to be resolved, got: %+v", order)
	}
}

func TestResolveInstallOrder_InstalledSkipped(t *testing.T) {
	catalog := []CatalogItem{
		{
			Name: "A",
			Kind: "Route",
			Dependencies: []Dependency{
				{Name: "B", Kind: "Traefik", Required: true},
			},
		},
		{
			Name: "B",
			Kind: "Traefik",
			Dependencies: nil,
		},
	}

	installed := map[string]bool{
		"b": true,
	}
	target := catalog[0] // A

	order, err := ResolveInstallOrder(target, catalog, installed)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(order) != 1 || order[0].Name != "A" {
		t.Fatalf("expected only A in resolution since B is already installed, got: %+v", order)
	}
}
