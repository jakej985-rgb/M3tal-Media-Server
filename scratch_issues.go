package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	data, err := os.ReadFile("/home/m3tal/.gemini/antigravity-ide/brain/086ab5ea-7ffe-4e61-ab10-11406eb74b77/.system_generated/steps/888/content.md")
	if err != nil {
		panic(err)
	}

	// Pattern to match /issues/<number> along with some title/context
	re := regexp.MustCompile(`"/jakej985-rgb/m3tal-core/issues/(\d+)"[^>]*>.*?<span>(.*?)</span>`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	fmt.Println("Found Issues:")
	seen := make(map[string]bool)
	for _, m := range matches {
		id := m[1]
		if !seen[id] {
			seen[id] = true
			fmt.Printf(" - Issue #%s: %s\n", id, m[2])
		}
	}

	// Fallback general regex
	re2 := regexp.MustCompile(`/issues/(\d+)`)
	matches2 := re2.FindAllStringSubmatch(string(data), -1)
	fmt.Println("\nRaw Issue Numbers Referenced:")
	seen2 := make(map[string]bool)
	for _, m := range matches2 {
		id := m[1]
		if !seen2[id] && !seen[id] {
			seen2[id] = true
			fmt.Printf(" - #%s\n", id)
		}
	}
}
