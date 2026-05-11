# 🚀 M3TAL Control Plane (v1.4.0.3)

> **M3tal-Media-Server** is the **Core Orchestrator and Control Plane** for the M3TAL Ecosystem. It provides an autonomous, Go-native observability and orchestration layer for homelabs and small-scale clusters.

![Go](https://img.shields.io/badge/go-1.21%2B-blue?logo=go)
![Python](https://img.shields.io/badge/python-3.10%2B-blue?logo=python)
![Docker](https://img.shields.io/badge/docker-m3tal-blue?logo=docker)
![Status](https://img.shields.io/badge/migration-Go--Native-success)

---

## 🏗️ Architectural Blueprint

The M3TAL ecosystem has transitioned to a **Go-native backend**. The legacy Python orchestrator (`m3tal.py`) has been fully replaced by the high-performance Go-native CLI and API.

### The Component Stack

* **Orchestrator (`./m3tal`)**: The Go-native CLI "Mission Control." It handles host-level lifecycle management and Docker daemon communication.
* **Backend API (`./m3tal-api`)**: The high-performance Go-native engine. It processes metrics and manages container state.
* **Core Library (`pkg/`)**: Importable Go modules containing all container and system management logic.
* **Dashboard**: A containerized Flask-based interface providing visualization and control.

---

## 📦 Prerequisites

* **Docker Engine**: v20.10+
* **Docker Compose**: v2.0+
* **Go**: 1.21+
* **Python**: v3.10+ (For Dashboard development)

---

## 🚀 Step-by-Step Deployment

Follow these steps in order to deploy the M3TAL stack.

### 1. Clone & Initialize Environment

```bash
git clone https://github.com/jakej985-rgb/M3tal-Media-Server.git
cd M3tal-Media-Server

# Create your environment file
cp .env.example .env
```

### 2. Configure Your Settings

Edit the `.env` file. You **must** set secure tokens before building.

```bash
# Generate secure 32-character tokens
openssl rand -hex 32
```

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DASHBOARD_SECRET` | Session security for the web UI. | Required |
| `API_TOKEN` | Auth token for Dashboard -> API comms. | Required |
| `BASE_STORAGE_PATH` | Host path for media/config. | `/mnt` |
| `DOMAIN` | Your local domain (e.g., `home.lan`). | `localhost` |

### 3. Run Automated Host Setup

This script prepares the standardized storage directories and sets correct permissions.

```bash
# Requires sudo to create /mnt directories
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### 4. Build Core Binaries

Compile the Go-native orchestrator and API backend:

```bash
# Build the CLI
go build -o m3tal ./cmd/m3tal

# Build the Backend API
go build -o m3tal-api ./cmd/api
```

### 5. Launch the Stack

The Go CLI manages the entire Docker Compose ecosystem located in `source/m3tal-stack`.

```bash
# Start Dashboard, API, and Media services
./m3tal up
```

---

## 🌐 Service Access

| Service | Port | Default URL | Description |
| :--- | :--- | :--- | :--- |
| **Dashboard** | `8082` | `http://localhost:8082` | Primary Web UI |
| **Backend API** | `5050` | `Internal Only` | REST API (Host Network) |
| **Docker Proxy**| `2375` | `Internal Only` | Secure Socket Access |

> 🔒 **Security Note:** Port `5050` and `2375` should not be exposed to the internet. Access the dashboard via port `8082` or through a reverse proxy like Traefik.

---

## 🛠 Troubleshooting & Verification

### Check System Status
```bash
./m3tal status
```

### Inspect Logs
If a service fails to start, check the container logs:
```bash
# Check Backend API logs
docker logs m3tal-api

# Check Dashboard logs
docker logs m3tal-dashboard
```

### Path Flexibility
M3TAL defaults to `/mnt`. If you use a different path (e.g., `/home/user/data`), update `BASE_STORAGE_PATH` in `.env` and restart the stack.

---

## 📜 License

Licensed under MIT.

*DocCritic Status: Linear tutorial implemented. All blockers resolved.*