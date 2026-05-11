# 🚀 M3TAL Control Plane (v1.4.0.3)

> A lightweight, autonomous, and self-healing container orchestration system for homelabs and small-scale clusters.

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
The `install.py` script is the primary entry point for new installations. It performs the following actions:
- **Dependency Audit**: Verifies Git, Docker, and Python availability.
- **Venv Setup**: Creates a Python virtual environment and installs requirements.
- **Scaffolding**: Creates necessary host directories (`/mnt/media`, `/mnt/config`, etc.).
- **Interactive Config**: Launches a wizard to generate your `.env` file and set your administrative password.
- **Network Initialization**: Creates the required `proxy` Docker bridge network.

```bash
# On Linux/macOS
python3 install.py

# On Windows (PowerShell)
python install.py
```

---

### 🎮 System Orchestration
Once installed, use the **M3TAL Orchestrator (`m3tal.py`)** to manage the lifecycle of the stack.

### 2. Launch (Unified Control Plane)

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
> Running `docker compose` inside `docker/media/` or other subdirs is **broken by design** and will fail. Always execute from the repository root to ensure correct volume mounting and networking.

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

## 📜 License & Support

Licensed under MIT.

If you like this project, give it a star ⭐!
