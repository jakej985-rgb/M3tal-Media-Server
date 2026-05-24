package plugin

import (
	"fmt"
	"strings"
)

// ResolveInstallOrder resolves the install order for a target catalog item.
// It checks the available catalog items, avoids circular dependencies, detects missing dependencies,
// and returns a list of CatalogItems sorted in topological order (dependencies first) that need to be installed.
func ResolveInstallOrder(target CatalogItem, catalog []CatalogItem, installed map[string]bool) ([]CatalogItem, error) {
	// Create lookup map for catalog items by name
	catalogMap := make(map[string]CatalogItem)
	for _, item := range catalog {
		key := strings.ToLower(item.Name)
		catalogMap[key] = item
	}

	var order []CatalogItem
	visited := make(map[string]int) // 0 = unvisited, 1 = visiting, 2 = visited
	var dfs func(item CatalogItem) error

	dfs = func(item CatalogItem) error {
		nameLower := strings.ToLower(item.Name)
		if visited[nameLower] == 2 {
			return nil
		}
		if visited[nameLower] == 1 {
			// Found cycle
			return fmt.Errorf("circular dependency detected: %s", item.Name)
		}

		visited[nameLower] = 1

		for _, dep := range item.Dependencies {
			depLower := strings.ToLower(dep.Name)
			// Skip if already installed
			if installed[depLower] {
				continue
			}

			// Look up dependency in available catalog
			depItem, found := catalogMap[depLower]
			if !found {
				if dep.Required {
					return fmt.Errorf("missing required dependency %q (kind: %s) for plugin %q", dep.Name, dep.Kind, item.Name)
				}
				// Skip non-required missing dependencies
				continue
			}

			if err := dfs(depItem); err != nil {
				if strings.Contains(err.Error(), "circular dependency detected") {
					return fmt.Errorf("%s -> %s", err.Error(), item.Name)
				}
				return err
			}
		}

		visited[nameLower] = 2
		order = append(order, item)
		return nil
	}

	// Start resolution from target
	if err := dfs(target); err != nil {
		return nil, err
	}

	return order, nil
}
