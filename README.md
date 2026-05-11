# 🚀 M3TAL Control Plane (v1.4.0.3)

> **M3tal-Media-Server** is the **Core Orchestrator and Control Plane** for the M3TAL Ecosystem. It provides an autonomous, Go-native observability and orchestration layer for homelabs and small-scale clusters.

![Go](https://img.shields.io/badge/go-1.21%2B-blue?logo=go)
![Python](https://img.shields.io/badge/python-3.10%2B-blue?logo=python)
![Docker](https://img.shields.io/badge/docker-m3tal-blue?logo=docker)
![Status](https://img.shields.io/badge/migration-Go--Native-success)

---

## 🏗️ Architectural Blueprint

The M3TAL ecosystem has transitioned to a **Go-native backend**, separating the management plane from the orchestration interface. Legacy Python logic (`m3tal.py`) has been fully decommissioned in favor of high-performance Go binaries.

### The Component Stack

* **Orchestrator (`./m3tal`)**: The Go-native CLI "Mission Control." It handles host-level lifecycle management and Docker daemon communication.
* **Backend API (`./m3tal-api`)**: The high-performance Go-native engine. It processes metrics and manages container state.
* **Core Library (`pkg/`)**: Importable Go modules containing all container and system management logic.
* **Dashboard**: A Flask-based interface (containerized) providing visualization and control.

---

## 📦 Prerequisites

* **Docker Engine**: v20.10+
* **Docker Compose**: v2.0+ (or `docker-compose-plugin`)
* **Go**: 1.21+ (For building backend modules)
* **Python**: v3.10+ (For local dashboard development/testing)

---

## 🚀 Installation & Deployment

### 1. Initialize & Automate Setup

The provided setup script is mandatory for first-time installations. It prepares the host environment and **compiles the core binaries**.

```bash
git clone https://github.com/jakej985-rgb/M3tal-Media-Server.git
cd M3tal-Media-Server

# Run the automated setup (Requires sudo for /mnt creation)
chmod +x scripts/setup.sh
./scripts/setup.sh
```

**The setup script performs the following:**
- Verifies system dependencies (Docker, Go, Python).
- Creates standardized storage paths: `/mnt/media`, `/mnt/config`, `/mnt/downloads`.
- Initializes your `.env` file from the template.
- **Compiles** the `./m3tal` CLI and `./m3tal-api` backend.

### 2. Configure Environment

**You must edit the `.env` file** before launching. Use the following commands to generate secure tokens:

```bash
# Generate a strong secret for DASHBOARD_SECRET and API_TOKEN
openssl rand -hex 32
```

Edit the file:
```bash
nano .env
```

| Variable | Description | Recommended Value |
| :--- | :--- | :--- |
| `DASHBOARD_SECRET` | Session security for the web UI. | Result of `openssl rand -hex 32` |
| `API_TOKEN` | Auth token for Dashboard -> API comms. | Result of `openssl rand -hex 32` |
| `LOCAL_IP` | Host machine IP for internal routing. | `192.168.x.x` or `127.0.0.1` |
| `DOMAIN` | Root domain for Traefik/Service Discovery. | `localhost` or your domain. |

---

## 🎮 System Orchestration

M3TAL provides a unified Go CLI for all system orchestration.

### Launching the Ecosystem

Running `./m3tal up` automatically invokes the Docker Compose stack in `source/m3tal-stack`, launching the API, Dashboard, and core services.

```bash
# Start the entire environment (API + Dashboard + Stack)
./m3tal up

# Check status of all containers
./m3tal status

# Stop everything safely
./m3tal down
```

---

## 🌐 Service Access & Networking

Once the stack is up, services are accessible via the following default ports. 

| Service | Internal Port | Default URL | Description |
| :--- | :--- | :--- | :--- |
| **Dashboard** | `8082` | `http://localhost:8082` | Primary Web Interface |
| **Backend API** | `5050` | `http://localhost:5050/api` | REST API (Host Network) |
| **Docker Proxy** | `2375` | `http://localhost:2375` | Secure Socket Proxy |

### Ingress & Traefik
The `m3tal-compose.yml` includes Traefik labels by default. If you have a Traefik instance running on the `m3tal` network, the dashboard will be automatically discovered at `http://m3tal.<YOUR_DOMAIN>`.

---

## 🧱 Data Persistence & Pathing

M3TAL mandates strict path mapping to ensure service portability and recovery:

* `/mnt/media`: Primary media library.
* `/mnt/config`: Persistent configuration for all containers.
* `/mnt/downloads`: Unified download directory.

*Note: If you need to change these, update the `BASE_STORAGE_PATH` in your `.env` and re-run setup.*

---

## 🔐 Security & Safety

* **Docker Socket Proxy**: The orchestrator does not touch the raw `/var/run/docker.sock`. It communicates via a restricted proxy (`docker-proxy`) that only allows safe operations.
* **Token Authentication**: All communication between the Dashboard and the Backend API requires the `X-API-Token` header, matching your `.env` configuration.

---

## 📜 License

Licensed under MIT.

*DocCritic Status: All blockers resolved. Build instructions and port maps finalized.*