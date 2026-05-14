# M3TAL Platform Makefile

.PHONY: build clean deb up down status help

VERSION := $(shell cat VERSION 2>/dev/null | sed 's/^v//' || echo "0.0.0")

help:
	@echo "M3TAL Management Utility"
	@echo ""
	@echo "Build:"
	@echo "  make build    - Compile CLI and API binaries"
	@echo "  make deb      - Build binaries & package into .deb (requires nfpm)"
	@echo "  make clean    - Remove compiled binaries and .deb artifacts"
	@echo ""
	@echo "Runtime:"
	@echo "  make up       - Start the M3TAL stack using the Go orchestrator"
	@echo "  make down     - Stop the M3TAL stack"
	@echo "  make status   - List container status"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-dash - Build the optional GUI dashboard Docker image"

build:
	@echo "🚀 Building M3TAL CLI..."
	go build -o m3tal ./cmd/m3tal
	@echo "🚀 Building M3TAL API..."
	go build -o m3tal-api ./cmd/api
	@echo "✅ Build complete."

deb: build
	@echo "📦 Packaging M3TAL Core .deb..."
	VERSION=$(VERSION) nfpm package --config packaging/nfpm.yaml --packager deb
	@echo "✅ .deb package created."

clean:
	@echo "🧹 Cleaning binaries..."
	rm -f m3tal m3tal-api
	rm -f *.deb
	@echo "✅ Cleanup complete."

up: build
	./m3tal up

down:
	./m3tal down

status:
	./m3tal list

docker-dash:
	@echo "🐳 Building M3TAL Dashboard Docker image..."
	docker build -t m3tal-dashboard:latest ./source/dashboard
	@echo "✅ Dashboard Docker image built."
	@echo "   Run with: docker compose -f /usr/share/m3tal/stack/m3tal-compose.yml up dash"