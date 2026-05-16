DocSmith here. The README has been audited and updated to ensure accuracy with the current repository structure and Go-native architectural requirements.

***

# M3TAL Media Server: Core Orchestrator

**M3TAL** is the Core Orchestrator of the M3TAL Ecosystem. This repository houses the Go-native orchestration logic, lifecycle management tools, and the CLI binary required to govern the media infrastructure.

---

## Architecture Overview

The M3TAL ecosystem is decoupled into three primary layers, interacting exclusively via network protocols and standardized file paths:

*   **Orchestrator (`cmd/m3tal`)**: The native Go CLI binary. It manages the Docker lifecycle, validates configurations, and acts as the entry point for all system administration.
*   **Backend API (`m3tal-goback`)**: The data layer. It provides the REST interface used by the Dashboard to query infrastructure status.
*   **Dashboard (`deploy/dashboard`)**: The web-based UI layer (Python/Flask). It communicates with the backend API to visualize system state and process control requests.

---

## Quick Start Guide

### 1. Prerequisites
*   **Docker & Docker Compose** (Plugin v2 recommended).
*   **OS**: Linux (Debian-based recommended for APT compatibility).

### 2. Installation
Building from source is not required. Install the pre-compiled binary via the official M3TAL APT repository:

```bash
# Add the M3TAL public key
curl -sL https://jakej985-rgb.github.io/m3tal-core/public.key | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# Add the repository to your sources
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# Update and install
sudo apt-get update
sudo apt-get install m3tal
```

### 3. Quick Demo
Once installed, initialize your environment:

```bash
# Initialize the configuration directory at /etc/m3tal
m3tal setup

# Start the core infrastructure (defined in deploy/stack)
m3tal up

# Start the dashboard UI (defined in deploy/dashboard)
m3tal dash up
```

---

## Filesystem Standard
To ensure the Orchestrator maintains state across containers, the system enforces these strict path conventions:

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/` | System Configuration |
| `/opt/m3tal/stack` | Docker Compose manifests |
| `/mnt/m3tal-media` | Required media storage mount point |

---

## Deployment: Docker Configuration
The Orchestrator manages the system by interacting with the Docker socket. Ensure all component services maintain path consistency:

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

*   **API-Only Communication**: Components must not share direct file modification access. All interactions between the Dashboard and the Orchestrator must route through the `m3tal-goback` API.
*   **Go-Native Migration**: The ecosystem is transitioning away from legacy wrappers. All core logic in the orchestrator is written in Go.
*   **Path Consistency**: All media assets must be served from the `/mnt` volume tree. This ensures that the Orchestrator, Backend, and Dashboard see the same file structure regardless of the container namespace.

---

## Related Projects
*   [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official web dashboard interface.
*   [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The backend API engine for data persistence and orchestration hooks.

*M3TAL — Modular Infrastructure Platform. Status: Go-Native Migration Active.*