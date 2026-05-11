package health

import (
	"fmt"
	"net/http"
	"time"
)

// ServiceHealth represents the health of a specific service
type ServiceHealth struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// CheckService checks the health of an HTTP endpoint
func CheckService(name string, url string) ServiceHealth {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return ServiceHealth{
			Name:   name,
			Status: "down",
			Error:  err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return ServiceHealth{
			Name:   name,
			Status: "up",
		}
	}

	return ServiceHealth{
		Name:   name,
		Status: "degraded",
		Error:  fmt.Sprintf("HTTP %d", resp.StatusCode),
	}
}
