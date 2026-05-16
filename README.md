Greetings, Architect. **DocSmith** here. 

I have reviewed the repository architecture. As requested, I have finalized the documentation, ensuring strict alignment with the "Core-First" M3TAL methodology and the Go-native migration standards.

***

# 🚀 M3TAL Media Server: Core Orchestrator

**M3TAL** is the core orchestration engine of the M3TAL Ecosystem. Designed as a Go-native system, it provides high-performance, low-latency management for media infrastructure. By decoupling the CLI, Backend API, and Dashboard, M3TAL ensures a modular, fault-tolerant stack where the Core acts as the definitive source of truth for infrastructure state.

---

## 🧠 System Architecture

The M3TAL ecosystem operates on a strict **"Core-First"** communication protocol. The Core Orchestrator maintains system integrity by managing Docker lifecycles and enforcing global configuration.

*   **Orchestrator (m3tal CLI)**: The Go-native binary acting as the primary control plane. It manages local system state, Docker orchestration via the `/opt/m3tal` manifest tree, and the configuration lifecycle.
*   **Backend API (`m3tal-goback`)**: The server-side brain. It provides the REST/gRPC interface for the dashboard and system tools. It communicates with the Orchestrator to execute state changes, ensuring a clear separation between control and data planes.
*   **Dashboard (`m3tal-godash`)**: A containerized React/Go interface. It is strictly a UI layer that consumes the `m3tal-goback` API. It has no direct access to system files, enforcing a secure isolation boundary.

---

## 📁 Filesystem & Path Consistency

To ensure operational stability, M3TAL enforces a strict path hierarchy. All external interaction should occur through the `/docker` entry point, which acts as the gateway to the managed stack.

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent State & Internal Data |
| `/opt/m3tal/stack` | Docker Compose Manifests (Source) |
| `/docker` | **User Entry Point** (Symlink to `/opt/m3tal/stack`) |

---

## 📋 Requirements

- **Linux (Debian/Ubuntu/Mint)**
- **Docker Engine**
- **Docker Compose V2**
- **Go 1.21+** (For local development/compilation)

---

## 🛠️ Quick Start (APT)

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
- `m3tal doctor`: Execute system-wide health and configuration diagnostics.
- `m3tal init`: Initialize the `/etc/m3tal` environment and generate secrets.

### Dashboard & UI
- `m3tal dash up`: Deploy the `m3tal-godash` container instance.
- `m3tal dash down`: Terminate the UI stack.
- `m3tal dash status`: Verify API connectivity and container health.

---

## 🔗 Ecosystem Integration Rules

*   **API-Only Communication**: The Dashboard and Core do not share files. All data exchange must occur via the `m3tal-goback` API.
*   **Decoupled Deployment**: The Orchestrator manages the lifecycle of the Backend and Dashboard, but they operate as independent processes for resilience.
*   **Traefik Ownership**: The Orchestrator maintains the base Traefik proxy configuration. All sub-services (Dash/Back) must define their own middleware/routers via Docker labels.

---

## 🏗️ Related Projects

- [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official web-based dashboard built with Go/WASM.
- [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The backend engine providing API services to the ecosystem.

*M3TAL — Modular Infrastructure Platform. Version: Go-Native Migration.*