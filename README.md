# 🚀 M3TAL Media Server

M3TAL is a modular infrastructure platform designed for precision media server orchestration. It features a DEB-first core for system management and a containerized dashboard for a rich user experience.

---

## 🧠 Final Architecture

- **Core (m3tal)**: Installed via APT, runs on the host. Handles system state, Docker orchestration, and Traefik base configuration.
- **Dash (m3tal-godash)**: Optional UI layer running as a Docker container. Connects to the Core via its HTTP API.
- **UX (/docker)**: A standardized symlink providing a clean entry point for all media stack operations.

---

## 📁 Filesystem Design

| Path | Purpose |
| :--- | :--- |
| `/usr/bin/m3tal` | Primary CLI Binary |
| `/etc/m3tal/.env` | Global Configuration Source of Truth |
| `/var/lib/m3tal/` | Persistent Application Data |
| `/opt/m3tal/stack` | Docker Compose Manifests (Real Path) |
| `/docker` | **User Entry Point** (Symlink to /opt/m3tal/stack) |

---

## 📋 Requirements

- **Linux (Debian/Ubuntu/Mint)**
- **Docker Engine**
- **Docker Compose V2**

---

## 🛠️ Quick Start (APT)

The recommended installation method for production systems.

```bash
# 1. Add the M3TAL repository
curl -sL https://jakej985-rgb.github.io/m3tal-core/public.key | sudo apt-key add -
echo 'deb [arch=amd64] https://jakej985-rgb.github.io/m3tal-core stable main' | sudo tee /etc/apt/sources.list.d/m3tal.list

# 2. Install M3TAL
sudo apt update && sudo apt install -y m3tal

# 3. Initialize environment
sudo m3tal init

# 4. Start the stack
m3tal up
```

### 🌐 Optional Dashboard

```bash
m3tal dash up
```

---

## 🔧 CLI Command Reference

### Core Operations
- `m3tal up`: Start the core infrastructure (Traefik, API, etc.)
- `m3tal down`: Stop all managed services
- `m3tal doctor`: Run system health and configuration diagnostics
- `m3tal init`: Generate default secrets and initialize the filesystem

### Dashboard Operations
- `m3tal dash up`: Start the Docker-based UI
- `m3tal dash down`: Stop the UI
- `m3tal dash status`: Check UI container health

---

## 🔗 Integration Rules

- **Core MUST NOT depend on Dash**: They communicate exclusively via HTTP API.
- **Dash MUST NOT assume host paths**: It operates within the container context and uses API calls for system interaction.
- **Traefik Ownership**: Core defines base entrypoints and networking; Dash defines only its specific router labels.

---

## 🧭 Troubleshooting

- **"Dashboard not installed"**: Run the quick install steps provided by the CLI when running `m3tal dash up`.
- **Permission Errors**: Ensure you have access to the Docker socket or run configuration commands with `sudo`.
- **Path Divergence**: The CLI always reports `/docker` as the stack path for UX consistency, even though files reside in `/opt/m3tal`.

*M3TAL — Modular Infrastructure Platform.*