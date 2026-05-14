# 🚀 M3TAL Media Server (v1.7)

**DocSmith Status:** *Architecture Scan Complete. Schema Validated.*

M3TAL is a high-performance media server control plane engineered for lifecycle orchestration. Built with **Go 1.26** (Core Orchestrator) and **Python 3.10** (Dashboard Service), it provides a unified interface for managing complex media service stacks.

---

## 🧠 Architecture: The M3TAL Ecosystem

The `./m3tal` binary acts as the **Source of Truth** for the entire ecosystem. It abstracts the container lifecycle, ensuring that the Go-native orchestrator manages the state of the infrastructure defined in `source/m3tal-stack/`.

### System Components
*   **Orchestrator (`./m3tal`)**: The Go-native binary acting as the primary control plane for lifecycle management, network configuration, and volume mapping.
*   **Infrastructure (`source/m3tal-stack`)**: The standardized Docker Compose manifests governing the containerized environment.
*   **Dashboard (`source/dashboard`)**: The legacy-compatible Python/Flask web interface. Note: This service is currently being phased out in favor of `m3tal-godash`.

> **Note on Migration**: We are currently in a Go-native migration phase. While the Dashboard remains Python-based for flexibility, all system-level orchestration, networking, and API interactions have been transitioned to Go to ensure memory safety and sub-millisecond execution times.

---

## 🔗 Related Projects
*   [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash): High-performance, Go-native dashboard rewrite.
*   [m3tal-goback](https://github.com/jakej985-rgb/m3tal-goback): The evolving Go-native backend API implementation.

---

## 🛠️ Prerequisites

1.  **Docker Engine**: v20.10+ (Ensure `docker` group membership).
2.  **Go Environment**: v1.26+ (Required for binary compilation).
3.  **DNS Mapping**: Required for local service discovery. Add to `/etc/hosts`:
    ```text
    127.0.0.1 m3tal.localhost api.localhost traefik.localhost
    ```
4.  **Storage Logic**: M3TAL enforces a strict `/mnt` mapping. Your host data at `BASE_STORAGE_PATH` is always mounted to `/mnt` inside containers, ensuring absolute path consistency across all microservices.

---

## 🚀 Quick Start (10 Minutes)

```bash
# 1. Clone & Setup
git clone https://github.com/jakej985-rgb/m3tal-core.git
cd m3tal-core
cp template.env .env
# Edit .env with your configuration (see Environment Variables below)

# 2. Build & Compile Orchestrator
chmod +x build.sh
./build.sh

# 3. Initialize Environment
./m3tal init
./m3tal up
```

### Pre-flight Checklist

| Requirement | Status |
|-------------|--------|
| Docker Engine 20.10+ | ✅ Required |
| Go 1.21+ | ✅ Required for build.sh |
| Ports 80, 443, 8080, 8082 Free | ✅ Required |
| `/mnt` Directory Writable | ✅ Required |

### Platform-Specific Setup

**Linux (Ubuntu/Debian)**:
```bash
# Add user to docker group (avoid sudo for docker)
sudo usermod -aG docker $USER

# Update system
sudo apt update && sudo apt upgrade -y

# Verify Docker access (log out and back in after adding to group)
docker ps
```

**macOS**:
```bash
# Install Docker Desktop
brew install --cask docker

# Start Docker Desktop and wait for it to be ready
# M3TAL will use Docker Desktop networking
```

**Windows**:
```bash
# Install Docker Desktop for Windows
# Ensure WSL2 backend is enabled
# M3TAL supports Linux containers on WSL2
```

---

## 📋 Environment Variables

The `.env` file requires the following variables:

| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `BASE_STORAGE_PATH` | Absolute path to media directory | `/mnt/media` or `/home/user/media` | Yes |
| `DOMAIN` | Your domain for Traefik routing | `m3tal.local` | Yes |
| `TZ` | Timezone for container clocks | `America/Denver` | Yes |
| `PUID` | User ID for file permissions | `1000` | Yes |
| `PGID` | Group ID for file permissions | `1000` | Yes |
| `DASHBOARD_PORT` | Dashboard web port | `8082` | No (default: 8082) |

### Quick .env Setup

```bash
# Create .env from template
cp template.env .env

# Edit with your values
nano .env  # or vi, vim, code, etc.

# Example minimal .env:
BASE_STORAGE_PATH=/mnt/media
DOMAIN=m3tal.local
TZ=America/Denver
PUID=1000
PGID=1000
```

---

## ⚙️ Service Routing & Communication
M3TAL uses an **API-Only Communication** model. The Frontend (Dashboard) communicates with the Backend via internal Docker networks, while all traffic is ingress-filtered through the **Traefik Gateway**.

| Service | Host Header | Internal Port |
| :--- | :--- | :--- |
| **Dashboard** | `m3tal.localhost` | 8082 |
| **Backend API** | `api.localhost` | 5050 |
| **Traefik Admin** | `traefik.localhost` | 8080 |

---

## 🛠️ CLI Command Reference

The `m3tal` binary provides a "Mission Control" interface:

*   `./m3tal up` : Boots the defined stack via Go-orchestrated Docker Compose.
*   `./m3tal down` : Graceful shutdown of all services.
*   `./m3tal config set <key> <val>` : Update environment variables safely.
*   `./m3tal list` : Displays status of active containers managed by the ecosystem.
*   `./m3tal dashpass <user> <pass>` : Securely rotates dashboard credentials.

---

## 🧭 Troubleshooting

*   **Orchestrator Desync**: If manual changes occur to the stack files, run `./m3tal init` to re-sync the environment.
*   **Log Access**: The Orchestrator captures all standard output from the stack. Use `./m3tal list` to identify service names, then `docker logs <service>` for deep inspection.
*   **Pathing**: If data is unreachable, verify that `BASE_STORAGE_PATH` in your `.env` is an absolute path on your host machine.

---

*M3TAL Core - Precision Media Infrastructure.*