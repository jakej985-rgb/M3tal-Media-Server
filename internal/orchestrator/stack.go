package orchestrator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// EnvVar represents a single environment variable for a stack.
// FromTemplate true if defined in the .env.template.
// FromCompose true if referenced in the stack's compose file.
// Defaults from the template are stored in Default.
// Required is true if the variable has no default.
//
// Example:
//
//	media.env.template   -> PLEX_CLAIM=XXXX
//	media-compose.yml     -> environment:
//	       - PLEX_CLAIM
type EnvVar struct {
	Name         string
	Default      string
	Required     bool
	FromTemplate bool
	FromCompose  bool
}

// Stack holds the discoverable stack information.
// Name is the stack identifier.
// ComposePath and TemplatePath are absolute paths.
type Stack struct {
	Name         string
	ComposePath  string
	TemplatePath string
}

// StackConfig contains the resolved variables for a stack.
// Vars maps variable name to its EnvVar descriptor.
//
// Additional helper: Merge into a map[string]string of key/value for runtime.
//
// The order of precedence when merging globals and stack env:
//   1. stack env file values override global env values.
//   2. If a variable is only in the template it will have an empty string as value.
//   3. Variables referenced only from compose are included with empty values.
//
// vars may contain a mix of required and optional.
//
// The validator can look for missing required values before launching.
//
// The stack config essentially is the representation that the CLI uses in all commands.
//
// Returning a pointer allows callers to modify in place.
//
// We expose methods for building the stack config.
//
// getRuntimePath returns the runtime env file path for this stack.
//
// BuildStackConfig loads template and compose, merges into a StackConfig.

// StackConfig holds the full stack config.
// NOTE: Vars map holds *EnvVar entries keyed by name.
// The value field is derived from the relevant env file.
type StackConfig struct {
	Stack Stack
	Vars  map[string]*EnvVar
}

// DiscoverStacks scans the given directory for stack compose files.
// It looks for files matching "*-compose.yml".
// For each found, it attempts to locate the corresponding ".env.template".
// If template is missing it will still create a Stack with empty TemplatePath and
// set the Vars map entry but with FromTemplate false.
func DiscoverStacks(dir string) ([]Stack, error) {
	var stacks []Stack
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, "-compose.yml") {
			continue
		}
		stackName := strings.TrimSuffix(name, "-compose.yml")
		composePath := filepath.Join(dir, name)
		templatePath := filepath.Join(dir, stackName+".env.template")
		// If template missing, keep empty path but warn later.
		if _, err := os.Stat(templatePath); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
			templatePath = ""
		}
		stacks = append(stacks, Stack{Name: stackName, ComposePath: composePath, TemplatePath: templatePath})
	}
	return stacks, nil
}

// ParseTemplate reads a *.env.template file and returns a map of EnvVar descriptors.
// Each line is expected to be KEY=value. Empty or comment lines are ignored.
func ParseTemplate(path string) map[string]*EnvVar {
	vars := make(map[string]*EnvVar)
	if path == "" {
		return vars
	}
	f, err := os.Open(path)
	if err != nil {
		return vars
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		vars[key] = &EnvVar{Name: key, Default: value, Required: value == "", FromTemplate: true, FromCompose: false}
	}
	return vars
}

// ParseComposeEnvVars parses a compose YAML file to extract ${VAR} references.
// It returns a map of EnvVar with FromCompose true and Required based on whether a default was found in the template.
func ParseComposeEnvVars(path string, templateVars map[string]*EnvVar) map[string]*EnvVar {
	envs := make(map[string]*EnvVar)
	f, err := os.Open(path)
	if err != nil {
		return envs
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// Regex to find ${VAR} or $VAR
	re := regexp.MustCompile(`\$\{([A-Z0-9_]+)}`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			if len(m) != 2 {
				continue
			}
			key := m[1]
			if _, exists := envs[key]; exists {
				continue
			}
			// Check if exists in template to get default
			tmplVar, found := templateVars[key]
			envs[key] = &EnvVar{
				Name:         key,
				Default:      "",
				Required:     !found || tmplVar.Default == "",
				FromTemplate: found,
				FromCompose:  true,
			}
		}
	}
	return envs
}

// BuildStackConfig combines template and compose parsing into a single config.
func BuildStackConfig(stack Stack) (*StackConfig, error) {
	tmplVars := ParseTemplate(stack.TemplatePath)
	composeVars := ParseComposeEnvVars(stack.ComposePath, tmplVars)
	// Merge both maps: if in both, keep template data but mark FromCompose true too.
	merged := make(map[string]*EnvVar)
	for k, v := range tmplVars {
		merged[k] = v
		if c, ok := composeVars[k]; ok {
			merged[k].FromCompose = true
			if c.Required {
				merged[k].Required = true
			}
		}
	}
	// Add variables only referenced in compose
	for k, v := range composeVars {
		if _, ok := merged[k]; !ok {
			merged[k] = v
		}
	}
	return &StackConfig{Stack: stack, Vars: merged}, nil
}

// LoadEnvFile loads a .env file into a map.
func LoadEnvFile(path string) map[string]string {
	env := make(map[string]string)
	if path == "" {
		return env
	}
	f, err := os.Open(path)
	if err != nil {
		return env
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		env[key] = value
	}
	return env
}

// SaveEnvFile writes a map of env variables to file.
func SaveEnvFile(path string, vars map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	for k, v := range vars {
		// Simple order, could sort for consistency
		_, err := fmt.Fprintf(writer, "%s=%s\n", k, v)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// BuildRuntimeEnv merges global.env and stack env file into a runtime env.
// It creates /run/m3tal/<stack>.env containing the resolved values.
// Global values are overridden by stack values.
// Returns the path to the runtime env file.
func BuildRuntimeEnv(stackName string) (string, error) {
	globalPath := system.ConfigPath
	stackPath := filepath.Join(system.StacksDir, fmt.Sprintf("%s.env", stackName))
	runtimeDir := "/run/m3tal"
	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		return "", err
	}
	runtimePath := filepath.Join(runtimeDir, fmt.Sprintf("%s.env", stackName))

	globals := LoadEnvFile(globalPath)
	stackEnv := LoadEnvFile(stackPath)

	merged := make(map[string]string)
	for k, v := range globals {
		merged[k] = v
	}
	for k, v := range stackEnv {
		merged[k] = v
	}

	// Write out runtime file
	if err := SaveEnvFile(runtimePath, merged); err != nil {
		return "", err
	}
	return runtimePath, nil
}

// Helper function to get the path to the stack's compose file.
func (c *StackConfig) ComposeFile() string {
	return c.Stack.ComposePath
}

// Helper to get the stack project's template path.
func (c *StackConfig) TemplateFile() string {
	return c.Stack.TemplatePath
}

// Helper to get a map of env vars (name -> value) for the stack.
func (c *StackConfig) EnvMap() map[string]string {
	m := make(map[string]string)
	for name, v := range c.Vars {
		val := v.ToString()
		if val != "" {
			m[name] = val
		}
	}
	return m
}

// ToString is an internal helper to get the value for a variable.
func (v *EnvVar) ToString() string {
	if v.FromCompose && v.Default != "" {
		return v.Default
	}
	return v.Default
}

// MergeVars merges the variable map with values from env files for the stack.
// This will be useful for CLI output.
func (c *StackConfig) MergeVars() map[string]string {
	merged := make(map[string]string)
	// read global and stack envs
	globals := LoadEnvFile(system.ConfigPath)
	stackEnv := LoadEnvFile(filepath.Join(system.StacksDir, fmt.Sprintf("%s.env", c.Stack.Name)))
	for k, v := range globals {
		merged[k] = v
	}
	for k, v := range stackEnv {
		merged[k] = v
	}
	// fill defaults for missing variables
	for name, vi := range c.Vars {
		if _, exists := merged[name]; !exists {
			merged[name] = vi.Default
		}
	}
	return merged
}
