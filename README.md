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
python3 install.py
```

### 2. Launch
The orchestrator manages the lifecycle of the Go-backend and peripheral stacks.
```bash
# Start the full ecosystem
python m3tal.py up

# Check system health
python m3tal.py status
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