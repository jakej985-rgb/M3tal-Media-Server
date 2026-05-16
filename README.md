DocSmith here. Archive updated. The repository documentation has been synchronized with the **Core Orchestrator** architectural profile. 

The structure is now optimized for the M3TAL Go-native standard, ensuring the `m3tal` binary is clearly defined as the nexus of the ecosystem.

***

# M3TAL Media Server: Core Orchestrator

**M3TAL** is the Core Orchestrator of the M3TAL Ecosystem. It is a Go-native system designed for high-performance control of media infrastructure. This repository houses the orchestration logic, lifecycle management, and the CLI binary that maintains the system state.

---

## Architecture Overview

The M3TAL ecosystem is decoupled into three primary layers, interacting exclusively via network protocols and standardized file paths.

*   **Orchestrator (`cmd/m3tal`)**: The native Go CLI. It manages Docker lifecycle, validates configurations, and serves as the primary interface for system administration.
*   **Backend API (`m3tal-goback`)**: The data layer. It provides the REST/gRPC interface used by the Dashboard to query infrastructure status.
*   **Dashboard (`deploy/dashboard`)**: The web-based UI layer (Python/Flask). It communicates with the backend API to visualize system state.

---

## Quick Start Guide

### 1. Prerequisites
*   **Go 1.21+** installed on the host.
*   **Docker & Docker Compose** installed.
*   **OS**: Linux-based environment (recommended for path consistency).

### 2. Build the Orchestrator
Navigate to the root directory to build the Go-native binary:
```bash
go build -o m3tal ./cmd/m3tal/main.go
sudo cp m3tal /usr/local/bin/
```

### 3. Deploy the Ecosystem
The repository includes a standardized deployment stack. Deploy the secondary modules via the `deploy/stack` directory:
```bash
cd deploy/stack
docker-compose up -d
```

---

## Filesystem Standard
To ensure the Orchestrator can manage assets and configuration across containers, the system follows these mandatory path conventions:

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/` | Configuration directory |
| `/opt/m3tal/stack` | Docker Compose manifests |
| `/mnt/m3tal-media` | Primary storage mount for media assets |

---

## Deployment: Docker Configuration
The Orchestrator requires access to the Docker socket to manage the ecosystem. Ensure your environment matches this configuration:

```yaml
services:
  m3tal-orchestrator:
    build: .
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /opt/m3tal:/opt/m3tal
      - /mnt:/mnt
    environment:
      - M3TAL_ROOT=/opt/m3tal
```

---

## Ecosystem Integration Rules

*   **API-Only Communication**: Modules must not communicate via shared databases or direct file modification. All interaction between the Dashboard and Orchestrator occurs via the `m3tal-goback` API.
*   **Go-Native Migration**: The platform is migrating to Go-native binaries. New modules should prioritize Go for performance and compatibility.
*   **Path Consistency**: All media assets MUST be served from the `/mnt` volume tree to ensure accessibility across the Orchestrator and the media-handling containers.

---

## Related Projects
*   [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official web dashboard interface.
*   [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The backend API engine for data persistence and orchestration hooks.

*M3TAL — Modular Infrastructure Platform. Status: Go-Native Migration Active.*