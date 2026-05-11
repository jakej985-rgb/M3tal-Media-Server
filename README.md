# 🚀 M3TAL Control Plane (v1.4.0.3)

> This repository, named **M3tal-Media-Server**, serves as the **Core Orchestrator and Control Plane** for the broader M3TAL Ecosystem. It is a lightweight, autonomous, and self-healing container orchestration system for homelabs and small-scale clusters.

![Python](https://img.shields.io/badge/python-3.10%2B-blue?logo=python) 
![Docker](https://img.shields.io/badge/docker-m3tal-blue?logo=docker)
![Version](https://img.shields.io/badge/version-v1.4.0.3-green)

M3TAL (Modern Media Management & Management) is an "Autonomous Local Cloud" that manages your Docker containers so you don't have to. It detects failures, scales services, and ensures your media stack is always online.

---

## 🔥 Key Features

* **🧠 Autonomous Self-Healing**: Detects crashed containers and restarts them automatically within seconds.
* **🌏 Distributed Leadership**: High Availability (HA) support — switch between nodes automatically if the master fails.
* **📈 Deep Metrics**: Real-time monitoring of CPU, Memory, and I/O for both the system and every individual container.
* **🔄 Auto-Scaling**: Automatically adjusts service replicas based on load (Upscale on high CPU, Downscale on idle).
* **🛡️ Hardened Security**: Token-based RBAC, BCrypt password hashing, and shell-injection protection.
* **🖥️ Web Dashboard**: Simple UI to manage your cluster, view metrics history, and approve AI actions.
* **🚑 Disaster Recovery**: Built-in `backup.sh` and `restore.sh` scripts for one-click stack recovery.

---

### 📦 Prerequisites

Before deploying the M3TAL Control Plane, ensure your host system meets the following requirements:

- **Docker Engine**: v20.10+ (Check with `docker info`)
- **Docker Compose**: v2.0+ (Check with `docker compose version`)
- **Python**: v3.9+ (For the Orchestrator CLI)
- **Git**: For repository cloning and updates

> [!NOTE]
> The M3TAL Go Backend is built and deployed within the Docker infrastructure. **You do not need a local Go toolchain** installed on your host system.

---

### 🚀 Installation Phase

The M3TAL platform uses an interactive, cross-platform installer that handles dependency verification, environment scaffolding, and initial configuration.

#### 1. Clone and Initialize
```bash
git clone https://github.com/jakej985-rgb/M3tal-Media-Server.git
cd M3tal-Media-Server
```

#### 2. Run the Interactive Installer
The `install.py` script is the primary entry point for new installations. It handles the venv creation, scaffolding, and initial configuration.

```bash
# On Linux/macOS
python3 install.py

# On Windows (PowerShell)
python install.py
```

---

### ⚙️ Detailed Configuration

If you prefer manual configuration or need to audit the system requirements:

#### 1. Infrastructure Requirements (Linux)
M3TAL assumes a standardized storage layout. If the installer didn't create these, run:

```bash
# Create standardized storage paths
sudo mkdir -p /mnt/media /mnt/config /mnt/downloads

# Ensure correct permissions (Standard UID/GID 1000)
sudo chown -R 1000:1000 /mnt/media /mnt/config /mnt/downloads
sudo chmod -R 775 /mnt/media /mnt/config /mnt/downloads
```

#### 2. Environment Variables (`.env`)
The system requires a `.env` file at the root. Key variables include:

```ini
# --- Core Config ---
DASHBOARD_PORT=8080
STATE_DIR=./state
LOG_LEVEL=info

# --- Auth ---
DASHBOARD_SECRET=your_super_secret_token_here
ADMIN_PASSWORD=your_secure_password

# --- Network ---
NETWORK_NAME=proxy
LOCAL_IP=192.168.1.100

# --- Storage ---
BASE_STORAGE_PATH=/mnt
```

---

### 🎮 System Orchestration

Once installed, follow these steps to manage your cluster.

#### 1. Activate Environment
Always ensure your Python virtual environment is active before running the orchestrator:
```bash
# On Linux/macOS
source venv/bin/activate

# On Windows (PowerShell)
.\venv\Scripts\activate
```

#### 2. Launch (Unified Control Plane)
M3TAL provides a unified CLI for all system orchestration.

```bash
# 1. Start the entire environment (Registry, Backend, and Stacks)
python m3tal.py up

# 2. Check the status of all services
python m3tal.py status

# 3. Stop everything safely
python m3tal.py stop
```

#### CLI Reference
| Command | Aliases | Description |
| :--- | :--- | :--- |
| `up` | `start` | Starts all Docker stacks in priority order. |
| `status` | `ps`, `ls` | Shows the health and status of all containers. |
| `stop` | `down` | Safely shuts down the entire system. |
| `restart` | — | Full system recycle (stop -> start). |

> [!NOTE]
> For a full list of all commands and advanced options, see the [CLI Guide](CLI_GUIDE.md).

> [!TIP]
> Add an alias for easier access: `alias m3tal="python m3tal.py"`

> [!WARNING]
> Running `docker compose` directly inside `source/m3tal-stack/` is **intentionally unsupported** for core services. To add your own services, see the [Custom Stacks](#-adding-your-own-media-services) section below.

---

### 🌐 External Access & Reverse Proxy

M3TAL uses a dedicated `proxy` Docker network to isolate services. To expose your media services (Plex, Sonarr, etc.) to the internet:
- **Integrated Gateway**: The system includes a Traefik-based reverse proxy pre-configured for the `proxy` network.
- **Automatic Discovery**: Services added with the correct Docker labels (e.g., `traefik.http.routers...`) are automatically picked up and assigned SSL certificates via Let's Encrypt.
- **Dashboard Access**: By default, the Dashboard is exposed on port `8080` for initial setup, but should be proxied behind a domain for production use.

---

### 🛠 Adding Your Own Media Services

To orchestrate your own services (Plex, Radarr, Nextcloud, etc.) alongside the M3TAL core:

1. **Decentralized Discovery**: M3TAL automatically scans the repository for any directory named `docker/`. It will discover and manage any files matching `*-compose.yml` or `docker-compose.yml` found within those folders.
2. **Standard Setup**: Simply create a `docker/` folder in your project directory and add your compose file (e.g., `source/my-plex/docker/plex-compose.yml`).
3. **Network Integration**: Ensure your services join the `proxy` network to communicate with the M3TAL backend and reverse proxy.
4. **Deployment**: Run `python m3tal.py up` — the orchestrator automatically discovers and deploys all matched stacks in the correct order.

### 3. Login

Open your browser to `http://YOUR_SERVER_IP:8080`.

* **Username**: `admin`
* **Password**: the admin password you chose during the interactive setup
* *⚠️ If you need to rotate or recover the admin password later, run `python scripts/config/manage_users.py --reset-admin`.*

---

## 🧱 Architecture

The system uses a high-performance Go-native observability backend that continuously executes a "Sense-Think-Act" loop:

1. **Registry** → Discovers your Docker stacks.
2. **Monitor** → Senses container health.
3. **Metrics** → Gathers performance data.
4. **Anomaly** → Identifies issues (crashes, leaks).
5. **Decision** → Plans recovery or scaling actions.
6. **Reconcile** → Executes actions (restart/scale).

### 📂 Components
- **Infrastructure (Docker Compose)**: The `source/m3tal-stack` directory contains the standardized Docker Compose configurations for the entire platform. This includes the Go Backend, Dashboard, and core services.
- **Core Orchestrator**: The `m3tal.py` CLI acts as the system's "Mission Control," wrapping Docker commands into high-level lifecycle actions (`up`, `down`, `status`, `logs`).
- **Control Plane**: The `agents/` directory houses the Python-based decision and healing agents.

---

## 🔐 Security & Safety

M3TAL is designed to be **safe**:

* **No Direct Calls**: Agents communicate only via atomic JSON state files.
* **Cooldowns**: Prevents "flapping" or restart loops by enforcing wait times between actions.
* **Allowlisting**: Only approved images and container names are permitted via the API.

---

## 🗺 Roadmap

* [ ] React-based "Admin Center" UI
* [ ] Predictive AI Scaling (predicting load spikes)
* [ ] Gossip protocol node discovery
* [ ] Plugin system for custom agents

---

## 🧹 Uninstallation & Cleanup

To completely remove M3TAL and its associated data:

1. **Stop Services**: `python m3tal.py down`
2. **Remove Infrastructure**:
   ```bash
   docker network rm proxy
   # Optional: docker system prune -a (WARNING: removes all unused images/networks)
   ```
3. **Delete Files**:
   Remove the repository directory and any host-mounted storage paths (e.g., the directory you chose for `/mnt/media` and `/mnt/config` during installation).

---

## 📜 License & Support

Licensed under MIT.

If you like this project, give it a star ⭐!
