# 🚀 M3TAL Media Server (v1.7)

| [🚀 Overview](README.md) | [⚙️ Environment](docs/ENVIRONMENT_VARIABLES.md) | [🛠️ Build](docs/BUILD_CONFIGURATION.md) | [🌐 Networking](docs/NETWORKING.md) | [🤖 Architecture](docs/ARCHITECTURE_VISION.md) |
| :---: | :---: | :---: | :---: | :---: |

**DocSmith Status:** *Architecture Scan Complete. Schema Validated. System Optimized.*

M3TAL is a high-performance media server control plane engineered for lifecycle orchestration. Built with **Go 1.21+** (Core Orchestrator) and **Python 3.10** (Dashboard Service), it provides a unified interface for managing complex media service stacks with absolute path consistency.

---

## 🧠 Architecture: The M3TAL Ecosystem

The `m3tal` binary acts as the **Source of Truth** for the entire ecosystem. It abstracts the container lifecycle, ensuring that the Go-native orchestrator manages the state of the infrastructure defined in `source/m3tal-stack/`.

### System Components

* **Orchestrator (CLI)**: The Go-native binary acting as the primary control plane. It interfaces with the Docker socket to manage lifecycle, network configuration, and volume mapping for the stack.
* **Infrastructure (`source/m3tal-stack`)**: The standardized Docker Compose manifests governing the containerized environment.
* **Dashboard (`source/dashboard`)**: The legacy-compatible Python/Flask web interface. Note: This service is currently being phased out in favor of `m3tal-godash`.

> **Note on Go-Native Migration**: M3TAL is currently undergoing a structural evolution. While the `dashboard` remains Python-based for legacy compatibility, all core orchestration, infrastructure management, and API-interfacing logic have been fully transitioned to Go-native modules to ensure memory safety and sub-millisecond execution times.

---

## 🔗 Related Projects

* [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash): High-performance, Go-native dashboard rewrite.
* [m3tal-goback](https://github.com/jakej985-rgb/m3tal-goback): The evolving Go-native backend API implementation.

---

## 🛠️ Prerequisites

1. **Docker Engine**: v20.10+ (Ensure `docker` group membership).
1. **Go Environment**: v1.21+ (Required for binary compilation).
1. **DNS Mapping**: Required for local service discovery. Add to `/etc/hosts`:

```text
127.0.0.1 m3tal.localhost api.localhost traefik.localhost
```

1. **Storage Logic**: M3TAL enforces a strict `/mnt` mapping. Your host data at `BASE_STORAGE_PATH` is always mounted to `/mnt` inside containers, ensuring absolute path consistency across all microservices.

---

## 🚀 Quick Start (10 Minutes)

```bash
# 1. Clone & Setup
git clone https://github.com/jakej985-rgb/m3tal-core.git
cd m3tal-core
cp template.env .env

# 2. Build & Compile Orchestrator
chmod +x build.sh
./build.sh

# 3. Initialize Environment
./m3tal init
./m3tal up
```

---

## ⚙️ Service Routing & Communication

M3TAL utilizes an **API-Only Communication** model. The Frontend (Dashboard) communicates with the backend via internal Docker networks, while all traffic is ingress-filtered through the **Traefik Gateway**.

| Service | Host Header | Internal Port |
| :--- | :--- | :--- |
| **Dashboard** | `m3tal.localhost` | 8082 |
| **Backend API** | `api.localhost` | 5050 |
| **Traefik Admin** | `traefik.localhost` | 8080 |

---

## 🛠️ CLI Command Reference

The `m3tal` binary provides a "Mission Control" interface for ecosystem management:

* `./m3tal up`: Boots the defined stack via Go-orchestrated Docker Compose.
* `./m3tal down`: Graceful shutdown of all services.
* `./m3tal config set <key> <val>`: Update environment variables safely.
* `./m3tal list`: Displays status of active containers managed by the ecosystem.
* `./m3tal dashpass <user> <pass>`: Securely rotates dashboard credentials.

---

## 🧭 Troubleshooting

* **Orchestrator Desync**: If manual changes occur to the stack files, run `./m3tal init` to re-sync the environment.
* **Log Access**: Use `./m3tal list` to identify service names, then `docker logs <service>` for deep inspection.
* **Pathing**: If data is unreachable, verify that `BASE_STORAGE_PATH` in your `.env` is an absolute path on your host machine.

---

*M3TAL Core - Precision Media Infrastructure.*