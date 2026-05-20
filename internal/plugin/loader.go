package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScanDir recursively scans a directory for plugin YAML files (*.yml, *.yaml).
// Returns a slice of parsed Plugins. Invalid files are logged but skipped.
func ScanDir(dir string) ([]Plugin, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Directory doesn't exist yet — not an error
		}
		return nil, fmt.Errorf("cannot stat plugin directory %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	var plugins []Plugin
	var scanErrors []string

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			scanErrors = append(scanErrors, fmt.Sprintf("%s: %v", path, err))
			return nil
		}

		p, err := ParsePlugin(data)
		if err != nil {
			scanErrors = append(scanErrors, fmt.Sprintf("%s: %v", path, err))
			return nil
		}

		p.SourcePath = path
		if err := p.Validate(); err != nil {
			scanErrors = append(scanErrors, fmt.Sprintf("%s: %v", path, err))
			return nil
		}

		plugins = append(plugins, *p)
		return nil
	})

	if err != nil {
		return plugins, fmt.Errorf("error walking plugin directory %s: %w", dir, err)
	}

	// Log scan errors to stderr but don't fail
	for _, e := range scanErrors {
		fmt.Fprintf(os.Stderr, "⚠️  Plugin scan warning: %s\n", e)
	}

	return plugins, nil
}

// LoadAll scans the given directories in order, loading and deduplicating plugins.
// Directories listed later take precedence (user overrides system).
// The standard order is: system dir first, user dir second.
func LoadAll(dirs ...string) (*Registry, error) {
	seen := make(map[string]*Plugin) // key: "kind/name"
	var allPlugins []Plugin

	for _, dir := range dirs {
		plugins, err := ScanDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s: %w", dir, err)
		}

		for i := range plugins {
			p := &plugins[i]
			key := p.Kind + "/" + p.Metadata.Name
			seen[key] = p
		}
	}

	// Build deduplicated list from the seen map
	for _, p := range seen {
		allPlugins = append(allPlugins, *p)
	}

	return BuildRegistry(allPlugins)
}

// LoadFromPaths is a convenience that loads plugins from the system default paths.
// It uses GetPluginDirs from the system package via the provided dirs.
func LoadFromPaths(systemDir, userDir string) (*Registry, error) {
	return LoadAll(systemDir, userDir)
}
