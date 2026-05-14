# 🚀 M3TAL Media Server (v1.7)

| [🚀 Overview](README.md) | [⚙️ Environment](docs/ENVIRONMENT_VARIABLES.md) | [🛠️ Build](docs/BUILD_CONFIGURATION.md) | [🌐 Networking](docs/NETWORKING.md) | [🤖 Architecture](docs/ARCHITECTURE_VISION.md) |
| :---: | :---: | :---: | :---: | :---: |

**DocSmith Status:** *Architecture Scan Complete. Schema Validated. System Optimized.*

M3TAL is a high-performance media server control plane engineered for lifecycle orchestration. Built with **Go 1.21+** (Core Orchestrator) and **Python 3.10** (Legacy Dashboard — runs inside Docker), it provides a unified interface for managing complex media service stacks with absolute path consistency.

---

## 🧠 Architecture: The M3TAL Ecosystem

The `m3tal` binary is the **sole orchestration entry point**. You interact only with this CLI tool; the Python dashboard runs **containerized** inside Docker (via `source/dashboard/Dockerfile`), and the Go binary manages it automatically. **No manual `pip install` is required for end users.**

### System Components

* **Orchestrator (CLI — `./m3tal`)**: The Go-native binary acting as the primary control plane. It interfaces with the Docker socket to manage lifecycle, network configuration, and volume mapping for the stack.
* **Infrastructure (`source/m3tal-stack/`)**: Standardized Docker Compose manifests governing the containerized environment. The `init` command generates these from templates in this directory.
* **Dashboard (`source/dashboard/`)**: Legacy Python/Flask web interface. **Containerized** via its own `Dockerfile`. Being phased out in favor of `m3tal-godash`.

### Dependency Chain

```text
User runs:  ./m3tal up
     ↓
Go binary reads:  .env  →  source/m3tal-stack/*.yml
     ↓
Go binary runs:  docker compose -f network-compose.yml up -d
                 docker compose -f routing-compose.yml up -d
                 docker compose -f m3tal-compose.yml up -d
     ↓
Docker builds/launches:  dashboard, API server, Traefik, media services
```

> **Note on Go-Native Migration**: M3TAL is currently undergoing a structural evolution. While the `dashboard` remains Python-based for legacy compatibility, all core orchestration, infrastructure management, and API-interfacing logic have been fully transitioned to Go-native modules to ensure memory safety and sub-millisecond execution times.

---

## 🔗 Related Projects

* [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash): High-performance, Go-native dashboard rewrite.
* [m3tal-goback](https://github.com/jakej985-rgb/m3tal-goback): The evolving Go-native backend API implementation.

---

## 📄 Environment Configuration

All configuration lives in a single `.env` file copied from `template.env`. See the [full reference](docs/ENVIRONMENT_VARIABLES.md) for details.

| Variable | Required | Default | Purpose |
| :--- | :--- | :--- | :--- |
| `BASE_STORAGE_PATH` | **YES** | `./data` | Host path for media + config (mounted to `/mnt` inside containers) |
| `DASHBOARD_SECRET` | **YES** | `change_me_immediately` | Session signing key (generate with `openssl rand -hex 32`) |
| `API_TOKEN` | **YES** | `change_me_api_token` | Dashboard-to-backend auth (generate with `openssl rand -hex 16`) |
| `ADMIN_PASSWORD` | **YES** | `admin_pass` | Initial dashboard login password |
| `DOMAIN` | No | `localhost` | Root domain for Traefik routing |
| `HTTP_PORT` | No | `8080` | Host port for HTTP traffic (Traefik entrypoint) |
| `DASHBOARD_PORT` | No | `8082` | Host port for dashboard UI |

> **Tip**: After first boot, change the admin password with: `./m3tal dashpass admin <new-password>`

---

### 🛠️ Prerequisites

1. **Docker Engine**: v20.10+ with `docker` group membership. Verify: `docker ps`
2. **Go Environment**: v1.21+ (Only required for building from source. See [Build Guide](docs/BUILD_CONFIGURATION.md)).
3. **DNS Mapping**: Add to `/etc/hosts` for local service discovery:

```text
127.0.0.1 m3tal.localhost api.localhost traefik.localhost
```

### 🖥️ Host Preparation

#### 1. Verify Docker Permissions

```bash
groups | grep docker || echo "WARNING: User not in docker group."
# Fix with: sudo usermod -aG docker $USER  (then log out and back in)
```

#### 2. Prepare Storage Directory

Create a **real, writable directory** on your host where M3TAL will store media and configs:

```bash
# Example: create a directory in your home folder
mkdir -p /home/yourusername/m3tal-media
```

> ⚠️ The path you create here will become `BASE_STORAGE_PATH` in your `.env` file. It is **mounted inside all containers at `/mnt`** for path consistency. If the directory does not exist or is not writable by Docker, containers will start with an **empty `/mnt`** — your data will not be accessible.

#### 3. Verify Required Ports Are Free

M3TAL needs the following ports **open and available** on your host:

| Port | Service | Purpose |
| :--- | :--- | :--- |
| `80`  | Traefik HTTP | Main web traffic entrypoint |
| `443` | Traefik HTTPS | Encrypted web traffic |
| `8080` | Traefik Dashboard | Admin panel (localhost only) |
| `8082` | M3TAL Dashboard | Web UI |

Check for conflicts:

```bash
# Linux
ss -tlnp | grep -E ':(80|443|8080|8082) '

# macOS
lsof -i :80 -i :443 -i :8080 -i :8082

# Windows PowerShell
netstat -ano | findstr ":80 :443 :8080 :8082"
```

If any are in use (e.g., by Apache, Nginx, or another service), either stop that service or change the port in your `.env` file.

---

## 🚀 Quick Start (10 Minutes)

```bash
# 1. Clone & Setup
git clone https://github.com/jakej985-rgb/m3tal-core.git
cd m3tal-core

# 2. Configure Environment
cp template.env .env
# ✏️  Edit .env — at minimum set BASE_STORAGE_PATH and DASHBOARD_SECRET

# 3. Build & Compile Orchestrator
chmod +x build.sh
./build.sh
# ✅ Produces two binaries: ./m3tal (CLI) and ./m3tal-api (backend API)

# 4. Initialize Environment
./m3tal init
# 🔍 Validates .env, generates secrets if missing, checks storage path

# 5. (Optional) Run Pre-Flight Check
./m3tal doctor
# ✅ Checks: Docker daemon, .env validity, storage path, port availability

# 6. Start the Stack
./m3tal up
```

### What `./m3tal init` Does

When you run `init`, the following happens:

1. **Reads** `template.env` (or existing `.env` on re-init)
2. **Auto-generates** secure values for `DASHBOARD_SECRET` and `API_TOKEN` if missing
3. **Validates** that `BASE_STORAGE_PATH` exists and is writable (with error guidance if not)
4. **Creates** the `.env` file with your settings
5. **Prepares** the Docker Compose files in `source/m3tal-stack/` for execution

### What `./build.sh` Produces

```
📥 Downloading dependencies...
🚀 Building M3TAL CLI...
     → Writes: ./m3tal (the main CLI binary)
🚀 Building M3TAL API...
     → Writes: ./m3tal-api (the backend API binary)
✅ Build complete. Binaries: ./m3tal, ./m3tal-api
```

After `init` and `up`:

| Service | URL | Purpose |
| :--- | :--- | :--- |
| **M3TAL Dashboard** | http://m3tal.localhost:8082 | Web UI for managing media services |
| **Backend API** | http://api.localhost:5050 | Internal REST API |
| **Traefik Admin** | http://traefik.localhost:8080 | Reverse proxy dashboard (localhost only) |

---

## ⚙️ Service Routing & Communication

M3TAL utilizes an **API-Only Communication** model. The Frontend (Dashboard) communicates with the backend via internal Docker networks, while all traffic is ingress-filtered through the **Traefik Gateway**.

| Service | Host Header | Internal Port | Status |
| :--- | :--- | :--- | :--- |
| **M3TAL Dashboard** | `m3tal.localhost` | Port 8082 | Proxy/Internal |
| **Backend API** | `api.localhost` | Port 5050 | Proxy/Internal |
| **Traefik Admin** | `traefik.localhost` | Port 8080 | Localhost Only |

> **Note on Dashboards**: The legacy Dashboard (`m3tal.localhost`) is included for v1.7 stability. If you are deploying for the first time, we recommend checking the [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash) repository for the modern replacement.

---

## 🛠️ CLI Command Reference

The `m3tal` binary provides a "Mission Control" interface for ecosystem management:

| Command | Description |
| :--- | :--- |
| `./m3tal init` | Initialize environment — creates `.env`, validates storage, generates secrets |
| `./m3tal doctor` | Run comprehensive pre-flight health check (Docker, env, ports, paths) |
| `./m3tal up` | Boot the defined stack via Go-orchestrated Docker Compose |
| `./m3tal down` | Graceful shutdown of all services |
| `./m3tal config set <key> <val>` | Update environment variables safely |
| `./m3tal config list` | Display current configuration |
| `./m3tal list` (or `ps`) | Status of active containers managed by the ecosystem |
| `./m3tal dashpass <user> <pass>` | Securely rotates dashboard credentials |
| `./m3tal logs` | View stack logs |

---

## 🧭 Troubleshooting

* **Orchestrator Desync**: If manual changes occur to the stack files, run `./m3tal init` to re-sync the environment.
* **Log Access**: Use `./m3tal list` to identify service names, then `docker logs <service>` for deep inspection.
* **Pathing**: If data is unreachable, verify that `BASE_STORAGE_PATH` in your `.env` is an absolute path on your host machine, and that it exists: `ls -la "$BASE_STORAGE_PATH"`.
* **Port Conflicts**: If services fail to start, check for port conflicts: `ss -tlnp | grep -E ':(80|443|8080|8082)'`. Stop any conflicting services (Apache, Nginx, etc.) or adjust port mappings in `.env`.
* **Traefik 404/503**: See the [Networking Guide](docs/NETWORKING.md) for debugging reverse proxy issues.

---

*M3TAL Core — Precision Media Infrastructure.*