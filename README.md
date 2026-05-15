**DocSmith Status:** *Architectural Scan Complete. Repository README updated to reflect current Go-native state and module relationships.*

# 🚀 M3TAL Media Server (v1.7)

| [🚀 Overview](README.md) | [⚙️ Environment](docs/ENVIRONMENT_VARIABLES.md) | [🛠️ Build](docs/BUILD_CONFIGURATION.md) | [🌐 Networking](docs/NETWORKING.md) | [🤖 Architecture](docs/ARCHITECTURE_VISION.md) |
| :---: | :---: | :---: | :---: | :---: |

M3TAL is a high-performance media server control plane engineered for lifecycle orchestration. This repository serves as the **M3TAL Media Server**, the core Orchestrator for the integrated M3TAL ecosystem, leveraging **Go 1.21+** for robust, sub-millisecond control plane operations. It provides a unified interface for managing complex media service stacks with absolute path consistency.

---

## 🧠 Architecture: The M3TAL Ecosystem

The `m3tal` binary acts as the **primary Orchestrator/Core**. It manages the lifecycle and configuration of the infrastructure defined by the `m3tal-stack`, ensuring absolute control over the deployed services.

### System Components

*   **Orchestrator (`m3tal` CLI)**: The Go-native binary compiled from the repository root. It serves as the control plane, interfacing directly with the Docker socket to manage container lifecycles, network configurations, and volume mappings for the `m3tal-stack`.
*   **Infrastructure (`deploy/stack/`)**: Standardized Docker Compose manifests. These define the deployment of critical services, including the Traefik proxy and the `dashboard` container, all managed by the `m3tal` orchestrator to ensure strict networking and storage compliance.
*   **Dashboard (`deploy/dashboard/`)**: The web interface for the M3TAL system. Built as a Python/Flask application, it is containerized and deployed as a core component of the `m3tal-stack`.

### Relationship Mapping

```mermaid
graph TD
    A[M3TAL CLI ./m3tal] -->|Manages Lifecycle & Configuration| B[m3tal-stack Docker Compose]
    B -->|Deploys| C[Traefik Gateway]
    B -->|Deploys| D[Dashboard Container]
    C -->|Routes Traffic To| D
    D -.->|API Requests To| E[m3tal-goback Remote]
```

**Operational Flow:**
The `M3TAL CLI` is the primary interface for managing the system. It interacts with `m3tal-stack` (Docker Compose definitions) to deploy and orchestrate containers. The `Traefik Gateway` routes external traffic to the `Dashboard`. The `Dashboard` communicates exclusively via API with the `m3tal-goback` service—an external, dedicated M3TAL Backend API—to execute business logic.

> **Go-Native Migration Status**: The Orchestrator is now fully Go-native, providing memory safety and peak performance for all control plane operations. While the `dashboard` remains Python-based for legacy support, the ecosystem is rapidly migrating to Go-native modules.

---

## 🔗 Related Projects

This repository is a core component of the M3TAL ecosystem:

*   [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash): The next-generation, high-performance, Go-native dashboard replacement.
*   [m3tal-goback](https://github.com/jakej985-rgb/m3tal-goback): The evolving Go-native backend API implementation providing the core data/logic layer.

---

## 🛠️ Quick Start (Recommended)

M3TAL is now distributed as a native Debian package. This is the recommended way to install it on Linux.

```bash
# 1. Add the M3TAL repository
curl -sL https://jakej985-rgb.github.io/m3tal-core/public.key | sudo apt-key add -
echo 'deb [arch=amd64] https://jakej985-rgb.github.io/m3tal-core stable main' | sudo tee /etc/apt/sources.list.d/m3tal.list

# 2. Install M3TAL Core
sudo apt-get update && sudo apt-get install -y m3tal-core

# 3. Initialize and Start
sudo m3tal init
m3tal up
```

### Development / Local Build
If you prefer to build from source:
```bash
go build -o m3tal ./cmd/m3tal
./m3tal init
./m3tal up
```

### Path Consistency Rule
The M3TAL ecosystem mandates that the host `BASE_STORAGE_PATH` is always mounted to `/mnt` inside every container. **Do not modify these mount points**, as the orchestrator relies on this structure for deterministic lifecycle management.

---

## 🛠️ CLI Command Reference

| Command           | Description                                                  |
| :---------------- | :----------------------------------------------------------- |
| `m3tal init`      | Syncs configuration and validates host path requirements.    |
| `m3tal up`        | Deploys the containerized stack via Docker Compose.          |
| `m3tal doctor`    | Diagnostic scan (Docker, ports, and path health).            |
| `m3tal config`    | Interface for updating global configuration parameters.      |
| `m3tal down`      | Stops and purges the M3TAL stack.                            |

---

## 🧭 Troubleshooting

*   **Desynchronization**: If infrastructure states drift, run `m3tal init` to refresh templates.
*   **Storage Path Issues**: Ensure `BASE_STORAGE_PATH` is an absolute path. The stack assumes `/mnt` is the internal media root.
*   **Dashboard API Errors**: Ensure the `m3tal-api` is reachable and the `API_TOKEN` matches.

*M3TAL Core — Precision Media Infrastructure.*