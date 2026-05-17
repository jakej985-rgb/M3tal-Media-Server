package main

import (
	"log"
	"os"

	"github.com/jakej985-rgb/m3tal-core/internal/api"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
)

func main() {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		apiToken = "m3tal-secret-token"
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5050"
	}

	// Initialize SQLite store
	dbPath := store.GetStatePath()
	db, err := store.Open(dbPath)
	if err != nil {
		log.Printf("⚠️  Could not open state database at %s: %v", dbPath, err)
		log.Println("⚠️  v2 engine endpoints will be disabled. Starting with v1 only.")
		if err := api.StartServer(port, apiToken); err != nil {
			log.Fatalf("❌ API server failed: %v", err)
		}
		return
	}
	defer db.Close()

	log.Printf("📦 State database: %s\n", dbPath)

	if err := api.StartServerWithStore(port, apiToken, db); err != nil {
		log.Fatalf("❌ API server failed: %v", err)
	}
}
