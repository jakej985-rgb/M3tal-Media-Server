Acknowledged. DocSmith online. 

I have audited your current repository architecture. Since this repository acts as the **Core Orchestrator** for the M3TAL Ecosystem, I have synchronized the `README.md` to reflect the Go-native migration and the mandatory architectural standards.

***

# 🚀 M3TAL Media Server: Core Orchestrator

**M3TAL** is the core orchestration engine of the M3TAL Ecosystem. Designed as a Go-native system, it provides high-performance, low-latency management for media infrastructure. By decoupling the CLI, Backend API, and Dashboard, M3TAL ensures a modular, fault-tolerant stack where the Core acts as the definitive source of truth for infrastructure state.

---

## ⚠️ Current Repository Status

**Observation Report [M3TAL-ARCH-001]:** This repository serves as the **Core Orchestrator** for the M3TAL ecosystem. It defines the foundational Go-native logic for infrastructure management, system-wide state coordination, and Docker lifecycle orchestration. 

This repository implements the `m3tal` binary, which acts as the primary control plane for the `m3tal-goback` (API) and `m3tal-godash` (Dashboard) modules.

---

## 🧠 System Architecture

The M3TAL ecosystem operates on a strict **"Core-First"** communication protocol. The Core Orchestrator maintains system integrity by managing Docker lifecycles and enforcing global configuration.

*   **Orchestrator (m3tal CLI)**: The Go-native binary acting as the primary control plane. It manages local system state, Docker orchestration via the `/opt/m3tal` manifest tree, and the configuration lifecycle.
*   **Backend API (`m3tal-goback`)**: The server-side brain. It provides the REST/gRPC interface for the dashboard and system tools. It communicates with the Orchestrator to execute state changes, ensuring a clear separation between control and data planes.
*   **Dashboard (`m3tal-godash`)**: A containerized interface. It is strictly a UI layer that consumes the `m3tal-goback` API. It has no direct access to system files, enforcing a secure isolation boundary.

---

## 📁 Filesystem & Path Consistency

To ensure operational stability, M3TAL enforces a strict path hierarchy. All external interaction should occur through the `/docker` entry point.

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent State & Internal Data |
| `/opt/m3tal/stack` | Docker Compose Manifests (Source) |
| `/docker` | **User Entry Point** (Symlink to `/opt/m3tal/stack`) |
| `/mnt` | Standard mount point for external storage volumes |
| `/mnt/m3tal-media` | Default location for M3TAL-managed media |

---

## 📋 Requirements

-   **Linux (Debian/Ubuntu/Mint)**
-   **Docker Engine & Docker Compose V2**
-   **Go 1.21+** (Required for Core Compilation)

---

## 🛠️ Deployment (Docker)

To run the M3TAL stack in a containerized environment, ensure the orchestrator has access to the Docker socket for lifecycle management:

```yaml
services:
  m3tal-core:
    image: m3tal/core:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /opt/m3tal:/opt/m3tal
      - /mnt:/mnt
    environment:
      - M3TAL_ENV=production
```

---

## 🔗 Ecosystem Integration Rules

*   **API-Only Communication**: All data exchange between `m3tal-godash`, `m3tal-goback`, and the Orchestrator must occur via defined API interfaces.
*   **Go-Native Migration**: The platform is currently transitioning to full Go-native binaries to replace legacy shell-script-based orchestration.
*   **Traefik Ownership**: The Orchestrator maintains the base Traefik proxy. All sub-services must define their own middleware/routers via Docker labels.

---

## 🏗️ Related Projects

- [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official web-based dashboard built with Go/WASM.
- [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The Go-native backend engine providing API services to the ecosystem.

*M3TAL — Modular Infrastructure Platform. Version: Go-Native Migration.*