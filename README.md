DocSmith here. Archive updated. The README has been refined to strictly reflect the current repository architecture and the Go-native migration status. 

***

# M3TAL Media Server: Core Orchestrator

**M3TAL** is the Core Orchestrator of the M3TAL Ecosystem. It is a Go-native system designed for high-performance control of media infrastructure. This repository houses the orchestration logic, lifecycle management, and the CLI binary that maintains the system state.

---

## Architecture Overview

The M3TAL ecosystem is decoupled into three primary layers, interacting exclusively via network protocols and standardized file paths:

*   **Orchestrator (`cmd/m3tal`)**: The native Go CLI. It manages Docker lifecycle, validates configurations, and serves as the primary interface for system administration.
*   **Backend API (`m3tal-goback`)**: The data layer. It provides the API interface used by the Dashboard to query infrastructure status.
*   **Dashboard (`deploy/dashboard`)**: The web-based UI layer (Python/Flask). It communicates with the backend API to visualize system state.

---

## Quick Start Guide

### 1. Prerequisites
*   **Docker & Docker Compose** installed.
*   **OS**: Linux-based environment (recommended for path consistency).

### 2. Install M3TAL
**Building from source is not required.** Install the pre-compiled Debian package directly from the APT repository.

```bash
# Add the M3TAL public key
curl -sL https://jakej985-rgb.github.io/m3tal-core/public.key | sudo apt-key add -

# Add the repository to your sources
echo 'deb [arch=amd64] https://jakej985-rgb.github.io/m3tal-core stable main' | sudo tee /etc/apt/sources.list.d/m3tal.list

# Install the orchestrator
sudo apt-get update
sudo apt-get install m3tal
```

### 3. Quick Demo
Verify your installation and initialize the ecosystem:

```bash
# Run the setup wizard to configure your /etc/m3tal environment
m3tal setup

# Start the core infrastructure (defined in deploy/stack)
m3tal up

# Start the dashboard UI (defined in deploy/dashboard)
m3tal dash up
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
The Orchestrator manages the lifecycle by interacting with the Docker socket. Ensure your environment mounts are consistent:

```yaml
services:
  m3tal-orchestrator:
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
*   **Go-Native Migration**: The platform is actively migrating to Go-native binaries. New modules should prioritize Go for performance and compatibility.
*   **Path Consistency**: All media assets MUST be served from the `/mnt` volume tree to ensure accessibility across the Orchestrator and media-handling containers.

---

## Related Projects
*   [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official web dashboard interface.
*   [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The backend API engine for data persistence and orchestration hooks.

*M3TAL — Modular Infrastructure Platform. Status: Go-Native Migration Active.*