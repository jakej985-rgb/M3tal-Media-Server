package engine

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// DeployResult captures the outcome of a stack deployment.
type DeployResult struct {
	Stack    string        `json:"stack"`
	Status   string        `json:"status"` // "success", "failed", "timeout"
	Duration time.Duration `json:"duration"`
	Logs     string        `json:"logs,omitempty"`
}

// DefaultTimeout for deployments.
const DefaultTimeout = 120 * time.Second

// DeployStack runs `docker compose up -d` on a compose file with safety.
//   - Context-based timeout (default 120s)
//   - Captured stdout/stderr
//   - Structured result
func DeployStack(composePath string, timeout time.Duration) (*DeployResult, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "up", "-d")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)
	logs := stdout.String() + stderr.String()

	if ctx.Err() == context.DeadlineExceeded {
		return &DeployResult{
			Stack:    composePath,
			Status:   "timeout",
			Duration: duration,
			Logs:     logs,
		}, fmt.Errorf("deployment timed out after %s", timeout)
	}

	if err != nil {
		return &DeployResult{
			Stack:    composePath,
			Status:   "failed",
			Duration: duration,
			Logs:     logs,
		}, fmt.Errorf("deployment failed: %w", err)
	}

	return &DeployResult{
		Stack:    composePath,
		Status:   "success",
		Duration: duration,
		Logs:     logs,
	}, nil
}

// ValidateAndDeploy runs validation before deployment.
// If validation fails, deployment is aborted.
func ValidateAndDeploy(composePath string, routes []RouteInput, timeout time.Duration) (*DeployResult, error) {
	// 1. Parse the compose file
	cf, err := ParseCompose(composePath)
	if err != nil {
		return &DeployResult{
			Stack:  composePath,
			Status: "failed",
			Logs:   err.Error(),
		}, err
	}

	// 2. Validate the compose file structure
	composeErrors := ValidateCompose(cf)
	if len(composeErrors) > 0 {
		var msgs []string
		for _, e := range composeErrors {
			msgs = append(msgs, e.Error())
		}
		return &DeployResult{
			Stack:  composePath,
			Status: "failed",
			Logs:   fmt.Sprintf("validation failed:\n%s", joinErrors(msgs)),
		}, fmt.Errorf("compose validation failed with %d errors", len(composeErrors))
	}

	// 3. Validate each route
	for _, route := range routes {
		routeErrors := ValidateRoute(route)
		if len(routeErrors) > 0 {
			var msgs []string
			for _, e := range routeErrors {
				msgs = append(msgs, e.Error())
			}
			return &DeployResult{
				Stack:  composePath,
				Status: "failed",
				Logs:   fmt.Sprintf("route validation failed for %s:\n%s", route.Service, joinErrors(msgs)),
			}, fmt.Errorf("route validation failed")
		}
	}

	// 4. Deploy
	return DeployStack(composePath, timeout)
}

// PullStack runs `docker compose pull` on a compose file.
func PullStack(composePath string, timeout time.Duration) (*DeployResult, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "pull")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)
	logs := stdout.String() + stderr.String()

	if err != nil {
		return &DeployResult{
			Stack:    composePath,
			Status:   "failed",
			Duration: duration,
			Logs:     logs,
		}, err
	}

	return &DeployResult{
		Stack:    composePath,
		Status:   "success",
		Duration: duration,
		Logs:     logs,
	}, nil
}

// StopStack runs `docker compose down` on a compose file.
func StopStack(composePath string, timeout time.Duration) (*DeployResult, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "down")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)
	logs := stdout.String() + stderr.String()

	if err != nil {
		return &DeployResult{
			Stack:    composePath,
			Status:   "failed",
			Duration: duration,
			Logs:     logs,
		}, err
	}

	return &DeployResult{
		Stack:    composePath,
		Status:   "success",
		Duration: duration,
		Logs:     logs,
	}, nil
}

func joinErrors(msgs []string) string {
	result := ""
	for _, m := range msgs {
		result += "  - " + m + "\n"
	}
	return result
}
