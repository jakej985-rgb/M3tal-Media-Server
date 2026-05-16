Greetings, Architect. **DocSmith** here.

I have reviewed the repository architecture. As requested, I have finalized the documentation, ensuring strict alignment with the "Core-First" M3TAL methodology and the Go-native migration standards. The existing blueprint accurately reflects the current repository state and the M3TAL Ecosystem requirements.

***

# 🚀 M3TAL Media Server: Core Orchestrator

**M3TAL** is the core orchestration engine of the M3TAL Ecosystem. Designed as a Go-native system, it provides high-performance, low-latency management for media infrastructure. By decoupling the CLI, Backend API, and Dashboard, M3TAL ensures a modular, fault-tolerant stack where the Core acts as the definitive source of truth for infrastructure state.

---

## ⚠️ Current Repository Status

**Observation Report [M3TAL-ARCH-001]:** Based on the current architectural scan, this repository **does not contain the direct implementation files** for the M3TAL Orchestrator (CLI), Infrastructure manifests, Backend API (`m3tal-goback`), or Dashboard (`m3tal-godash`). Detected modules: None.

This README serves as the definitive architectural blueprint and functional specification for the **M3TAL Media Server project**. It outlines the intended design, component interactions, and operational procedures, acting as the primary documentation for the ecosystem's core orchestrator. Future development will populate this repository with the relevant Go-native source code and configuration files detailed below.

---

## 🧠 System Architecture

The M3TAL ecosystem operates on a strict **"Core-First"** communication protocol. The Core Orchestrator maintains system integrity by managing Docker lifecycles and enforcing global configuration.

*   **Orchestrator (m3tal CLI)**: The Go-native binary acting as the primary control plane. It manages local system state, Docker orchestration via the `/opt/m3tal` manifest tree, and the configuration lifecycle. The CLI interacts directly with the host system and local Docker daemon.
*   **Backend API (`m3tal-goback`)**: The server-side brain. It provides the REST/gRPC interface for the dashboard and system tools. It communicates with the Orchestrator to execute state changes, ensuring a clear separation between control and data planes. The Backend API is strictly an intermediary, translating external requests into Orchestrator commands.
*   **Dashboard (`m3tal-godash`)**: A containerized React/Go interface. It is strictly a UI layer that consumes the `m3tal-goback` API. It has no direct access to system files, enforcing a secure isolation boundary. The Dashboard provides a visual control panel, making requests to the Backend API, which then relays necessary commands to the Orchestrator.

---

## 📁 Filesystem & Path Consistency

To ensure operational stability, M3TAL enforces a strict path hierarchy. All external interaction should occur through the `/docker` entry point, which acts as the gateway to the managed stack. Furthermore, adherence to the `/mnt` standard is critical for media storage.

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent State & Internal Data |
| `/opt/m3tal/stack` | Docker Compose Manifests (Source) |
| `/docker` | **User Entry Point** (Symlink to `/opt/m3tal/stack`) |
| `/mnt` | Standard mount point for external storage volumes |
| `/mnt/m3tal-media` | Recommended default location for M3TAL-managed media content (e.g., video, audio libraries) |

---

## 📋 Requirements

-   **Linux (Debian/Ubuntu/Mint)**
-   **Docker Engine**
-   **Docker Compose V2**
-   **Go 1.21+** (For local development/compilation of Orchestrator components)

---

## 🛠️ Quick Start (APT)

This section details the deployment process for the *compiled M3TAL Orchestrator binary* via APT. Once installed, the `m3tal` CLI will manage the rest of the ecosystem components (Backend and Dashboard).

```bash
# 1. Add the M3TAL repository
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/public.key | sudo gpg --dearmor -o /etc/apt/keyrings/m3tal.gpg

echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/m3tal.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 2. Install M3TAL
sudo apt update && sudo apt install -y m3tal

# 3. Initialize environment
sudo m3tal init

# 4. Start the stack
m3tal up
```

---

## 🔧 CLI Command Reference

### Core Orchestration
- `m3tal up`: Initialize and start the core infrastructure.
- `m3tal down`: Stop all managed container services.
- `m3tal status`: Check the status of all managed stacks.
- `m3tal doctor`: Execute system-wide health and configuration diagnostics.
- `m3tal init`: Initialize the `/etc/m3tal` environment and generate secrets.

### Dashboard Management
- `m3tal dash start`: Start the M3TAL dashboard.
- `m3tal dash stop`: Stop the M3TAL dashboard.
- `m3tal dash restart`: Restart the dashboard service.
- `m3tal dash logs`: View real-time dashboard logs.
- `m3tal dash status`: Check dashboard container health.

---

## 🔗 Ecosystem Integration Rules

*   **API-Only Communication**: The Dashboard and Core do not share files. All data exchange between `m3tal-godash` and `m3tal-goback`, and between `m3tal-goback` and the Orchestrator, must occur via defined API interfaces.
*   **Decoupled Deployment**: The Orchestrator manages the lifecycle of the Backend and Dashboard, but they operate as independent processes/containers for resilience and scalability.
*   **Traefik Ownership**: The Orchestrator maintains the base Traefik proxy configuration. All sub-services (Dash/Back) must define their own middleware/routers via Docker labels for dynamic routing.

---

## 🏗️ Related Projects

- [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official web-based dashboard built with Go/WASM, consuming the `m3tal-goback` API.
- [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The Go-native backend engine providing API services to the ecosystem, commanded by the M3TAL Orchestrator.

*M3TAL — Modular Infrastructure Platform. Version: Go-Native Migration.*