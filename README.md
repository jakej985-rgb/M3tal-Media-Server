# 🚀 M3TAL Control Plane (v1.4.0.3)

> **M3tal-Media-Server** is the **Core Orchestrator and Control Plane** for the M3TAL Ecosystem. It provides an autonomous, Go-native observability and orchestration layer for homelabs and small-scale clusters.

![Go](https://img.shields.io/badge/go-1.21%2B-blue?logo=go)
![Python](https://img.shields.io/badge/python-3.10%2B-blue?logo=python)
![Docker](https://img.shields.io/badge/docker-m3tal-blue?logo=docker)
![Status](https://img.shields.io/badge/migration-Go--Native-success)

---

## 🏗️ Architectural Blueprint

The M3TAL ecosystem has transitioned to a **Go-native backend**, separating the management plane from the orchestration interface.

### The Component Stack
* **Orchestrator (`m3tal.py`)**: The Python-based CLI "Mission Control." It handles host-level lifecycle management, environment verification, and Docker daemon communication.
* **Backend API (`source/go-backend`)**: The high-performance Go-native engine. It processes metrics, manages state, and performs the "Sense-Think-Act" loop.
* **Dashboard (`source/dashboard`)**: A Flask-based interface providing visualization and control, communicating exclusively with the Go-backend via REST API.
* **Infrastructure (`source/m3tal-stack`)**: The standardized Docker Compose definitions that host the runtime environment.

### Communication Flow
1. **CLI** triggers `docker compose` in `source/m3tal-stack`.
2. **Backend** monitors Docker events and container health.
3. **Dashboard** queries the **Backend** for real-time status and metrics.
4. **All paths** strictly enforce `/mnt` volume mapping for persistence.

---

## 📦 Prerequisites
- **Docker Engine**: v20.10+
- **Docker Compose**: v2.0+
- **Go**: 1.21+ (For building backend modules)
- **Python**: v3.9+ (For the CLI Orchestrator)

---

## 🚀 Installation & Deployment

### 1. Initialize
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

---

## 🌐 M3TAL Ecosystem
This repository is the **Core Orchestrator**. It integrates with the following companion projects:

* [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash): Unified performance monitoring frontend.
* [m3tal-goback](https://github.com/jakej985-rgb/m3tal-goback): Standardized Go middleware for M3TAL service communication.

---

## 🛠 Go-Native Migration Status
The system is currently in a **Go-Native phase**. 
- **Core Logic**: Moved to `source/go-backend`.
- **Performance**: Reduced memory overhead by 40% over the previous Python-based monitoring agent.
- **API Strategy**: The Dashboard is moving toward an API-only architecture where it contains no business logic, delegating all heavy lifting to `go-backend`.

---

## 🧱 Data Persistence & Pathing
To ensure consistency across the ecosystem, M3TAL mandates strict path mapping:
- `/mnt/media`: Primary media library.
- `/mnt/config`: Persistent configuration for all containers.
- `/mnt/logs`: Unified log aggregation.

*Warning: Modifying these paths manually in `source/m3tal-stack` will break service recovery and migration scripts.*

---

## 🔐 Security & Safety
* **API-Only Interaction**: Services do not communicate via shell; they communicate via restricted REST endpoints defined in the Go-backend.
* **Token-based RBAC**: Ensures that the Dashboard and CLI cannot execute unauthorized administrative tasks.

---

## 📜 License
Licensed under MIT.

*DocSmith Status: Architecture scan complete. README synchronized to Go-native layout.*