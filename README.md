# 🚀 M3TAL Media Server (v1.7)

| [🚀 Overview](README.md) | [⚙️ Environment](docs/ENVIRONMENT_VARIABLES.md) | [🛠️ Build](docs/BUILD_CONFIGURATION.md) | [🌐 Networking](docs/NETWORKING.md) | [🤖 Architecture](docs/ARCHITECTURE_VISION.md) |
| :---: | :---: | :---: | :---: | :---: |

**DocSmith Status:** *Architecture Scan Complete. Schema Validated. System Optimized.*

M3TAL is a high-performance media server control plane engineered for lifecycle orchestration. This repository delivers the core orchestrator and its integrated dashboard, leveraging **Go 1.21+** for robust orchestration and core logic. It provides a unified interface for managing complex media service stacks with absolute path consistency.

---

## 🧠 Architecture: The M3TAL Ecosystem

The `m3tal` binary acts as the **primary Orchestrator/Core** for this system. It manages the lifecycle and configuration of the infrastructure defined by `m3tal-stack`, ensuring precise control over the deployed services, including the `dashboard`.

### System Components

*   **Orchestrator (`m3tal`)**: The Go-native binary compiled from this repository's source. It serves as the control plane, interfacing directly with the Docker socket to manage container lifecycle, network configurations, and volume mappings for the `m3tal-stack`.
*   **Infrastructure (`source/m3tal-stack/`)**: This directory contains the standardized Docker Compose manifests. These definitions orchestrate the deployment of critical services, including the Traefik proxy and the `dashboard` container, all managed by the `m3tal` orchestrator to ensure strict networking and storage compliance.
*   **Dashboard (`source/dashboard/`)**: The current web interface for the M3TAL system. It is a Python/Flask application, containerized via its dedicated `Dockerfile`, and deployed as part of the `m3tal-stack` managed by the Go orchestrator.

### Relationship Mapping

```mermaid
graph TD
    A[M3TAL CLI (./m3tal)] -->|Manages Lifecycle & Configuration| B[m3tal-stack (Docker Compose Definitions)]
    B -->|Deploys & Orchestrates| C[Traefik Gateway Container]
    B -->|Deploys & Orchestrates| D[Dashboard Container (source/dashboard)]
    C -->|Routes HTTP/S Traffic To| D
    D -.->|API Requests To (External)| E[m3tal-goback (Remote M3TAL Backend API)]
```

**Operational Flow:**
The `M3TAL CLI` (`./m3tal`) is the primary interface for managing the system. It interacts with the `m3tal-stack` (Docker Compose definitions) to deploy and orchestrate the necessary containers, including the `Traefik Gateway` and the `Dashboard`. The `Traefik Gateway` then routes external HTTP/S requests to the `Dashboard` service. The `Dashboard` itself, while part of this repository's deployment, communicates with the `m3tal-goback` service, which is an external, dedicated M3TAL Backend API instance, typically deployed as a separate entity within the M3TAL ecosystem.

> **Note on Go-Native Migration**: This repository's `m3tal` orchestrator is fully Go-native, providing memory safety and sub-millisecond execution for all infrastructure management and control plane operations. While the `dashboard` (`source/dashboard/`) remains Python-based for legacy compatibility within this stack, the wider M3TAL ecosystem is undergoing a structural evolution, with the core backend API logic transitioning to a dedicated Go-native module (`m3tal-goback`) to ensure optimal performance and reliability across the platform.

---

## 🔗 Related Projects

As a core component of the M3TAL ecosystem, this repository integrates with and relies on external, purpose-built modules:

*   [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash): The next-generation, high-performance, Go-native dashboard rewrite, slated to replace the `source/dashboard` in future releases.
*   [m3tal-goback](https://github.com/jakej985-rgb/m3tal-goback): The evolving Go-native backend API implementation, providing the core data and business logic for the M3TAL ecosystem. The `source/dashboard` within this repository communicates with this external service via API calls.

---

## 📄 Environment Configuration

All critical configuration parameters for the `m3tal` orchestrator and the `m3tal-stack` reside in the root `.env` file. For a comprehensive list and detailed descriptions, refer to the [Environment Variables](docs/ENVIRONMENT_VARIABLES.md) documentation.

| Variable            | Required | Purpose                                                      |
| :------------------ | :------- | :----------------------------------------------------------- |
| `BASE_STORAGE_PATH` | **YES**  | The absolute host path where all media data is stored. This path is consistently mounted to `/mnt` inside all relevant containers. |
| `API_TOKEN`         | **YES**  | Authentication token used for secure communication between the `Dashboard` and the external `m3tal-goback` API. |
| `DASHBOARD_SECRET`  | **YES**  | A cryptographic key used by the `Dashboard` for session signing and other security-sensitive operations. |

---

## 🛠️ Quick Start

Initiate and manage your M3TAL Media Server deployment with these commands:

```bash
# 1. Compile the M3TAL orchestrator binary
go build -o m3tal main.go 

# 2. Initialize the infrastructure and validate paths
./m3tal init

# 3. Launch the containerized M3TAL stack
./m3tal up
```

### Path Consistency Rule
A fundamental principle of the M3TAL ecosystem is strict path consistency. Every containerized service within the `m3tal-stack` utilizes a standardized volume mapping: the host's `BASE_STORAGE_PATH` is always mounted to `/mnt` inside the container. **It is imperative not to alter these internal mount points** within the compose files, as deviations will disrupt the orchestration layer's visibility and management of your media data.

---

## 🛠️ CLI Command Reference

The `m3tal` orchestrator provides a suite of commands for comprehensive system management:

| Command               | Description                                                  |
| :-------------------- | :----------------------------------------------------------- |
| `./m3tal init`        | Initializes the M3TAL infrastructure, generates necessary configuration files, and performs initial path validations. |
| `./m3tal up`          | Deploys and starts the containerized M3TAL stack as defined by the Go-orchestrated Compose manifests. |
| `./m3tal doctor`      | Executes a diagnostic scan, validating host health, Docker socket accessibility, and required port availability. |
| `./m3tal config`      | Provides an interface for managing global `.env` settings and other M3TAL configuration parameters. |
| `./m3tal down`        | Stops and removes all containers, networks, and volumes associated with the M3TAL stack. |

---

## 🧭 Troubleshooting

*   **Orchestrator Desynchronization**: If infrastructure configurations appear inconsistent or outdated, execute `./m3tal init` to refresh Compose file templates and re-synchronize the orchestrator with the intended state.
*   **Storage Path Issues**: Ensure that `BASE_STORAGE_PATH` in your `.env` file specifies an absolute path on the host system. The containerized stack rigidly assumes `/mnt` inside containers is the root of your media drive. Incorrect host path configuration will lead to inaccessible media.
*   **Network Unreachability**: If deployed services are not accessible, verify that Traefik is running correctly within the `m3tal-stack` and that the `traefik_public` network is properly attached to all target containers that require external access.
*   **Dashboard API Errors**: If the dashboard reports issues communicating with the backend, confirm that the external `m3tal-goback` service is operational and accessible from the dashboard container, and that `API_TOKEN` is correctly configured in the `.env` file.

*M3TAL Core — Precision Media Infrastructure.*