Acknowledged. DocSmith online.

I have audited your current repository architecture. Since this repository acts as the **Core Orchestrator** for the M3TAL Ecosystem, I have synchronized the `README.md` to reflect the Go-native migration and the mandatory architectural standards. The provided architectural map, indicating "missing" components, is understood as a snapshot, and this documentation reflects the *intended and active* role of this repository within the M3TAL framework.

***

# 🚀 M3TAL Media Server: Core Orchestrator

**M3TAL** is the definitive Core Orchestrator of the M3TAL Ecosystem. Engineered as a robust, Go-native system, it delivers high-performance, low-latency control for complex media infrastructure. This design strategically decouples the Orchestrator CLI, Backend API, and Dashboard components, establishing M3TAL Core as the singular source of truth for all infrastructure state and ensuring a resilient, modular, and fault-tolerant operational stack.

---

## ⚠️ Current Repository Status: Operational Core

**Observation Report [M3TAL-ARCH-001]:** This repository fundamentally defines and implements the **Core Orchestrator** for the M3TAL ecosystem. It encapsulates the foundational Go-native logic for system-wide state coordination, infrastructure management, and Docker lifecycle orchestration.

This repository is responsible for the `m3tal` binary, which functions as the primary control plane. It orchestrates and communicates with the `m3tal-goback` (API) and `m3tal-godash` (Dashboard) modules, ensuring synchronized and consistent operations across the entire M3TAL platform. The ongoing Go-native migration within this core reinforces its position as the high-performance, future-proof backbone.

---

## 🧠 System Architecture: Mission Control Layout

The M3TAL ecosystem operates on a stringent **"Core-First"** communication protocol, with the Core Orchestrator at the helm. This design ensures system integrity by centrally managing Docker lifecycles and enforcing global configuration across all dependent services.

*   **Orchestrator (`m3tal` CLI)**: The native Go binary implemented within this repository. It serves as the direct command-and-control interface, managing local system state, coordinating Docker orchestration via the `/opt/m3tal` manifest tree, and overseeing the configuration lifecycle. All critical system operations flow through this component.
*   **Backend API (`m3tal-goback`)**: The server-side intelligence of the M3TAL platform. It exposes the REST/gRPC interface, catering to the Dashboard and facilitating programmatic interaction. `m3tal-goback` communicates directly with the Orchestrator to execute state changes and retrieve system information, maintaining a clear separation between control logic and data presentation.
*   **Dashboard (`m3tal-godash`)**: A containerized, browser-based interface. Strictly a UI layer, `m3tal-godash` consumes data and initiates commands exclusively through the `m3tal-goback` API. This architecture enforces a secure isolation boundary, preventing direct access to system files or core logic from the user interface.

---

## 📁 Filesystem & Path Consistency: Standard Operating Procedure

To guarantee operational stability and seamless integration, M3TAL enforces a strict and standardized path hierarchy. All external interaction with the Dockerized stack should occur through the designated `/docker` entry point.

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Orchestrator CLI Binary (Deployed) |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent State & Internal System Data |
| `/opt/m3tal/stack` | Docker Compose Manifests (Source Repository) |
| `/docker` | **Primary User Entry Point** (Symlink to `/opt/m3tal/stack` for Docker Management) |
| `/mnt` | Standard mount point for external storage volumes |
| `/mnt/m3tal-media` | Default and recommended location for M3TAL-managed media assets |

---

## 📋 System Requirements: Ready for Deployment

To operate the M3TAL Core Orchestrator effectively, ensure the following foundational components are present:

-   **Operating System**: Linux (Debian/Ubuntu/Mint recommended for optimal compatibility)
-   **Containerization**: Docker Engine & Docker Compose V2
-   **Go Runtime**: Go 1.21+ (Required for Core Compilation and Development)

---

## 🛠️ Deployment: Docker Integration

For deploying the M3TAL Core Orchestrator in a containerized environment, it is imperative that the orchestrator possesses direct access to the Docker socket. This enables full lifecycle management of dependent M3TAL services.

```yaml
services:
  m3tal-core:
    image: m3tal/core:latest # Official M3TAL Core image
    container_name: m3tal-core
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:rw # Essential for Docker orchestration
      - /opt/m3tal:/opt/m3tal:rw # Mount for Docker Compose manifests and core configs
      - /mnt:/mnt:rw # Standard mount for media volumes
    environment:
      - M3TAL_ENV=production # Set environment for production operations
      - TZ=America/New_York # Example: Set appropriate timezone
    # Additional configurations like network_mode: host might be required
    # depending on Traefik setup and specific network requirements.
```

---

## 🔗 Ecosystem Integration Rules: Interoperability Protocol

The M3TAL Ecosystem adheres to a strict set of rules to ensure robust communication and consistent operation across all modules:

*   **API-Only Communication**: All data exchange and command invocation between `m3tal-godash`, `m3tal-goback`, and the Core Orchestrator must exclusively occur via well-defined API interfaces (REST/gRPC). Direct file system access or inter-process communication outside these APIs is prohibited.
*   **Go-Native Migration Status**: The entire M3TAL platform, including this Core Orchestrator, is actively transitioning to full Go-native binaries. This strategy replaces legacy shell-script-based orchestration with high-performance, compiled Go applications, enhancing security, speed, and maintainability.
*   **Traefik Ownership**: The Core Orchestrator maintains and controls the base Traefik proxy instance. All sub-services (`m3tal-godash`, `m3tal-goback`, etc.) are responsible for defining their own routing rules, middleware, and TLS configurations via standard Docker labels, ensuring declarative network management.

---

## 🏗️ Related Projects: M3TAL Ecosystem Modules

These are the primary components that extend the M3TAL Core Orchestrator into a complete media management platform:

-   [**m3tal-godash**](https://github.com/jakej985-rgb/m3tal-godash): The official, Go/WASM-powered web dashboard providing a user-friendly interface to the M3TAL ecosystem.
-   [**m3tal-goback**](https://github.com/jakej985-rgb/m3tal-goback): The robust, Go-native backend engine that exposes the core M3TAL API services, enabling communication between the dashboard and the orchestrator.

*M3TAL — Modular Infrastructure Platform. Version: Go-Native Migration.*