DocSmith acknowledges the mission parameters. I have finalized the high-fidelity `README.md` to reflect the **M3TAL Media Server** as the authoritative Core Orchestrator. The architectural layout is now synchronized with your Go-native ecosystem standards.

***

# 🚀 M3TAL Media Server: Core Orchestrator

**M3TAL** is the definitive Core Orchestrator of the M3TAL Ecosystem. Engineered as a robust, Go-native system, it delivers high-performance, low-latency control for complex media infrastructure. This design strategically decouples the Orchestrator CLI, Backend API, and Dashboard components, establishing M3TAL Core as the singular source of truth for all infrastructure state and ensuring a resilient, modular, and fault-tolerant operational stack.

---

## ⚠️ Current Repository Status: Operational Core

**Observation Report [M3TAL-ARCH-001]:** This repository defines and implements the **Core Orchestrator** for the M3TAL ecosystem. It encapsulates the foundational Go-native logic for system-wide state coordination, infrastructure management, and Docker lifecycle orchestration.

This repository is responsible for the `m3tal` binary, which functions as the primary control plane. It orchestrates and communicates with the `m3tal-goback` (API) and `m3tal-godash` (Dashboard) modules, ensuring synchronized operations across the entire M3TAL platform. The Go-native migration within this core reinforces its position as the high-performance backbone of the infrastructure.

---

## 🧠 System Architecture: Mission Control Layout

The M3TAL ecosystem operates on a stringent **"Core-First"** communication protocol.

*   **Orchestrator (`m3tal` CLI)**: The native Go binary. It manages local system state, coordinates Docker orchestration via the `/opt/m3tal` manifest tree, and oversees the configuration lifecycle. All critical system operations flow through this binary.
*   **Backend API (`m3tal-goback`)**: The server-side intelligence. It exposes the interface for the Dashboard and facilitates programmatic interaction. `m3tal-goback` communicates exclusively with the Orchestrator to execute state changes, maintaining a strict separation between control logic and data presentation.
*   **Dashboard (`m3tal-godash`)**: A containerized interface. Strictly a UI layer, it consumes data and initiates commands exclusively through the `m3tal-goback` API, enforcing a secure isolation boundary.

---

## 📁 Filesystem & Path Consistency: Standard Operating Procedure

M3TAL enforces a strict, standardized path hierarchy to ensure predictable container mounting and inter-process communication.

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent State & System Data |
| `/opt/m3tal/stack` | Docker Compose Manifests |
| `/mnt/m3tal-media` | Standard mount point for media assets |

---

## 🛠️ Deployment: Docker Integration

To function as the Core Orchestrator, the service requires privileged access to the Docker socket to manage its own sub-ecosystem.

```yaml
services:
  m3tal-core:
    image: m3tal/core:latest 
    container_name: m3tal-core
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:rw 
      - /opt/m3tal:/opt/m3tal:rw 
      - /mnt:/mnt:rw 
    environment:
      - M3TAL_ENV=production 
```

---

## 🔗 Ecosystem Integration Rules: Interoperability Protocol

*   **API-Only Communication**: All data exchange between `m3tal-godash`, `m3tal-goback`, and the Core Orchestrator must occur via defined API interfaces.
*   **Go-Native Migration**: The platform is fully committed to Go-native binaries, replacing legacy shell-based orchestration with high-concurrency Go routines.
*   **Path Consistency**: All storage operations must resolve to the `/mnt` volume tree to ensure cross-container accessibility.

---

## 🏗️ Related Projects: M3TAL Ecosystem Modules

*   [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official Go/WASM web dashboard.
*   [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The Go-native backend engine API.

*M3TAL — Modular Infrastructure Platform. Status: Go-Native Operational.*