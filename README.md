# M3TAL Ecosystem

Welcome to the **M3TAL Core**, a Go-native orchestration layer designed to simplify infrastructure management. M3TAL provides a unified CLI and a robust backend service to deploy, monitor, and manage your containerized stacks.

---

## 🏗️ Architectural Overview

M3TAL operates as a cohesive system designed for stability and security:

*   **Orchestrator (`m3tal` CLI):** The primary entry point for interaction. It manages the lifecycle of the M3TAL daemon and generates configurations.
*   **Systemd Backend (`m3tal-api`):** A Go-native daemon running as a system service. It persists application state and provides the API layer for the frontend.
*   **Dashboard Container:** A Python-based visualization layer. It communicates with the Backend API to fetch stack status and metrics.
*   **Data Isolation:** 
    *   Configuration: `/etc/m3tal/.env`
    *   Persistence: `/var/lib/m3tal/state.db`
    *   Deployment: `/opt/m3tal/stack` (managed via symlinks)

---

## 🚀 Installation

You do not need to build from source. M3TAL provides pre-compiled binaries via our official APT repository for Debian/Ubuntu-based systems.

```bash
# 1. Add GPG Key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add Repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install M3TAL
sudo apt update && sudo apt install -y m3tal
```

---

## 🛠️ Getting Started

### 1. Configuration
Before deploying, initialize your environment settings. The configuration wizard will guide you through setting up your paths and environment variables (stored in `/etc/m3tal/.env`):

```bash
sudo m3tal config wizard
```

### 2. Interactive Menu
The CLI provides an interactive TUI to manage your services and containers:

```bash
sudo m3tal
```

---

## ⚙️ Configuration Reference

M3TAL uses a centralized environment file. Key configurable parameters include:

| Variable | Default | Description |
| :--- | :--- | :--- |
| `DASHBOARD_PORT` | `8082` | Port for the Python dashboard. |
| `HTTP_PORT` | `8080` | Backend API port. |
| `ADMIN_PASSWORD`| `admin_pass` | Dashboard access credentials. |
| `DOMAIN` | `localhost` | Primary domain for stack routing. |
| `BASE_STORAGE_PATH` | `./data` | Root directory for volume mounts. |

*Note: Be sure to update `DASHBOARD_SECRET` and `API_TOKEN` immediately upon initial setup.*

---

## 📁 System Directories

M3TAL adheres to Linux standard practices for file system placement:

*   **/etc/m3tal/**: Contains `.env` and sensitive configuration tokens.
*   **/var/lib/m3tal/**: Contains the `state.db` which tracks the health and status of your ecosystem.
*   **/opt/m3tal/stack**: Holds the actual Compose files and infrastructure definitions utilized by the orchestrator.

---

*Documentation maintained by DocSmith, M3TAL Ecosystem Documentation Architect.*