# 🚀 M3TAL Media Server (v1.4)

M3TAL is a high-performance media server control plane built with **Go 1.26** and **Python 3.10**. It features a native Go orchestrator that manages a multi-stack Docker environment for media services, routing, and system monitoring.

## 🧠 Architecture: The M3TAL Orchestrator

The `./m3tal` binary is the **Source of Truth** for the M3TAL lifecycle.
- **Important**: Do **not** run `docker compose` commands manually unless debugging. Interacting directly with the compose files may cause state desync with the Go Orchestrator.
- The orchestrator programmatically manages:
   - `network-compose.yml`: Virtual networking and DNS.
   - `routing-compose.yml`: Traefik gateway and SSL.
   - `m3tal-compose.yml`: Core API, Dashboard, and Agents.

---

## 🛠️ Prerequisites

1. **Docker Engine**: v20.10+ (Ensure your user is in the `docker` group).
2. **Go 1.26+**: Required to build the core binaries.
3. **Python 3.10+**: Required for the dashboard runtime.
4. **Storage**: A mount point for media.
   - Default: `/mnt`. 
   - Override: Set `BASE_STORAGE_PATH` in your `.env`.
   - *Pre-flight*: `mkdir -p /your/path && sudo chown $USER:$USER /your/path`.

---

## 🚀 Quick Start (Linux/WSL)

### 1. Build the Platform
Use the provided `Makefile` to automate compilation:
```bash
make build
```

### 2. Initialize Environment
Prepare dependencies and configuration:
```bash
# Setup Dashboard dependencies
pip install -r source/dashboard/requirements.txt

# Initialize .env
cp .env.example .env
```

### 3. Start the Stack
Initialize the orchestrator and launch services:
```bash
./m3tal pull
make up
```

---

## ⚙️ Configuration (.env)

| Variable | Description | Default |
| :--- | :--- | :--- |
| `BASE_STORAGE_PATH` | Host path for all media data | `/mnt` |
| `DOMAIN` | Domain name for external access | `localhost` |
| `DASHBOARD_PORT` | Port for the web interface | `8082` |
| `API_PORT` | Port for the Go-native API | `5050` |
| `API_TOKEN` | Secure token for CLI-to-API auth | (Generated) |

---

## ✅ Verification & Access

### Service Endpoints & Firewall
| Service | Local URL | Port | Firewall Required |
| :--- | :--- | :--- | :--- |
| **Dashboard** | `http://localhost:8082` | `8082` | **Yes** (External Access) |
| **Backend API** | `http://localhost:5050` | `5050` | **No** (Internal Only) |
| **Traefik Web** | `http://localhost:80` | `80` | **Yes** (HTTP) |
| **Traefik SSL** | `http://localhost:443` | `443` | **Yes** (HTTPS) |

### Initial Login
Run the following to set your admin password:
```bash
./m3tal dashpass admin yourpassword
```

---

## 🧭 Troubleshooting

- **Logs**: Always check the orchestrator logs first: `./m3tal list`
- **Manual Debug**: If necessary, check container logs directly: `docker logs m3tal-dashboard`
- **Reverse Proxy**: Traefik is included by default. For custom domains, ensure your `DOMAIN` variable in `.env` matches your DNS records.

---

*M3TAL Core - Built for Performance and Stability.*