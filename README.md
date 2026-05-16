DocSmith here. Integration confirmed. The foundational template has been successfully instantiated. As the **Core Orchestrator**, this repository now reflects the precise architectural requirements of the M3TAL ecosystem.

The document is ready for deployment.

***

# 🚀 M3TAL Media Server: Core Orchestrator

**M3TAL** is the definitive Core Orchestrator of the M3TAL Ecosystem. Engineered as a robust, Go-native system, it delivers high-performance, low-latency control for complex media infrastructure. This design strategically decouples the Orchestrator CLI, Backend API, and Dashboard components, establishing M3TAL Core as the singular source of truth for all infrastructure state.

---

## ⚠️ Current Repository Status: Core Infrastructure
**Observation Report [M3TAL-ARCH-001]:** This repository serves as the **Core Orchestrator**. It provides the foundational Go-native logic for system-wide state coordination and container lifecycle management. 

As the primary binary (`m3tal`), it is responsible for the execution of the ecosystem. It does not contain the presentation layer or the secondary API logic; it orchestrates the `m3tal-goback` and `m3tal-godash` modules via API-only communication, ensuring a strictly isolated, fault-tolerant operational stack.

---

## 🧠 System Architecture: Mission Control Layout

The M3TAL ecosystem operates on a stringent **"Core-First"** communication protocol:

*   **Orchestrator (`m3tal` CLI)**: The native Go binary. It manages local system state, coordinates Docker orchestration via the `/opt/m3tal` manifest tree, and oversees the configuration lifecycle.
*   **Backend API (`m3tal-goback`)**: The server-side intelligence. It provides the data-layer for the Dashboard. It communicates with the Orchestrator to execute state changes, maintaining a strict separation between control logic and data presentation.
*   **Dashboard (`m3tal-godash`)**: The containerized UI layer. It consumes data and initiates commands exclusively through the `m3tal-goback` API, enforcing a secure isolation boundary.

---

## 📁 Filesystem & Path Consistency: Standard Operating Procedure

M3TAL enforces a strict, standardized path hierarchy to ensure predictable container mounting and inter-process communication across the ecosystem.

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent State & System Data |
| `/opt/m3tal/stack` | Docker Compose Manifests |
| `/mnt/m3tal-media` | Standard mount point for all media assets |

---

## 🛠️ Deployment: Docker Integration

To function as the Core Orchestrator, the service requires privileged access to the Docker socket to manage its sub-ecosystem. Ensure your `docker-compose.yml` aligns with the M3TAL standard:

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

*   **API-Only Communication**: All data exchange between modules MUST occur via defined API interfaces. No direct file-system manipulation between the Dashboard and Orchestrator is permitted.
*   **Go-Native Migration**: The platform is fully committed to Go-native binaries. All new modules must be written in Go to maintain binary compatibility and performance parity.
*   **Path Consistency**: All storage operations must resolve to the `/mnt` volume tree to ensure cross-container accessibility.

---

## 🏗️ Related Projects: M3TAL Ecosystem Modules

*   [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official Go/WASM web dashboard.
*   [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The Go-native backend engine API.

*M3TAL — Modular Infrastructure Platform. Status: Go-Native Migration in progress.*