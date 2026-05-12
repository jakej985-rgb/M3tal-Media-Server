# M3TAL Platform Makefile

.PHONY: build clean up down status help

help:
	@echo "M3TAL Management Utility"
	@echo "Usage:"
	@echo "  make build    - Compile CLI and API binaries"
	@echo "  make clean    - Remove compiled binaries"
	@echo "  make up       - Start the M3TAL stack using the Go orchestrator"
	@echo "  make down     - Stop the M3TAL stack"
	@echo "  make status   - List container status"

build:
	@echo "🚀 Building M3TAL CLI..."
	go build -o m3tal ./cmd/m3tal
	@echo "🚀 Building M3TAL API..."
	go build -o m3tal-api ./cmd/api
	@echo "✅ Build complete."

clean:
	@echo "🧹 Cleaning binaries..."
	rm -f m3tal m3tal-api
	@echo "✅ Cleanup complete."

up: build
	./m3tal up

down:
	./m3tal down

status:
	./m3tal list
