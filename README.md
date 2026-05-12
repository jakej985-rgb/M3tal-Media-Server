# 🚀 M3TAL Media Server (v1.5)

M3TAL is a high-performance media server control plane built with **Go 1.26** and **Python 3.10**. It features a native Go orchestrator that manages a multi-stack Docker environment for media services, routing, and system monitoring.

## 🧠 Architecture: The M3TAL Orchestrator

The `./m3tal` binary is the **Source of Truth** for the M3TAL lifecycle.
- **Boot Sequence**: Running `./m3tal up` starts all containers, including the Backend API and the Dashboard. No separate Python startup is required for the production stack.
- **Important**: Do **not** run `docker compose` commands manually unless debugging. Interacting directly with the compose files may cause desync with the Go Orchestrator.

---

## 🛠️ Prerequisites

1. **Docker Engine**: v20.10+ (Ensure your user is in the `docker` group).
2. **Go 1.26+**: Required to build the core binaries (`m3tal`, `m3tal-api`). 
   - *Linux Tip*: If you see `go: No such file or directory`, run: `export PATH=$PATH:/usr/local/go/bin` (or your Go install path).
3. **Storage**: A mount point for media.
   - Default: `./data` (Portable, user-local).
   - Override: Set `BASE_STORAGE_PATH` in your `.env`.

### 🚦 Pre-flight Check
Before starting, ensure the following:
*   **Port Availability**: Ports `80`, `443`, `8080`, and `5050` must not be in use by other services (like Nginx or Apache).
*   **Permissions**: If you are not in the `docker` group, prefix commands with `sudo` or run `sudo usermod -aG docker $USER`.
*   **Source of Truth**: The orchestrator manages files in `source/m3tal-stack/`. For system-wide installs, these reside in `/usr/share/m3tal/stack/`.

---

## 🚀 Quick Start (Linux/WSL)

### 0. First Time Setup (Bootstrap)
Follow these steps exactly to prepare your environment:
```bash
# 1. Clone the repository
git clone https://github.com/jakej985-rgb/m3tal-core.git
cd m3tal-core

# 2. Prepare environment template
cp template.env .env

# 3. Build and Download Dependencies
chmod +x build.sh
./build.sh

# 4. Initialize Secrets and Paths
./m3tal init
```

### 1. Build the Platform
The `./build.sh` script (or `build.ps1` for Windows) handles dependency downloading (`go mod download`), verification, and compilation. It is the recommended way to build.

---

## ⚙️ Configuration & Storage

### Environment Variables (.env)
| Variable | Description | Default |
| :--- | :--- | :--- |
| `BASE_STORAGE_PATH` | **Host Path** for all media data | `./data` |
| `HTTP_PORT` | Traefik web entrypoint | `8080` |
| `DASHBOARD_PORT` | Port for the web interface | `8082` |
| `API_PORT` | Port for the Go-native API | `5050` |

### 📂 Storage Mapping Logic
M3TAL uses a **Consistency Model** for storage paths. 
*   **Host Path**: Defined by `BASE_STORAGE_PATH` (e.g., `/home/user/media`).
*   **Container Path**: The Orchestrator automatically maps your host path to **`/mnt`** inside all media containers.
*   **Example**: If you set `BASE_STORAGE_PATH=/data/movies`, the movie service will see its data at `/mnt/movies`.

---

## 🌐 Networking & SSL

### Docker Networks
The Orchestrator creates and manages a dedicated Docker network named **`m3tal`**.
*   **External Access**: To connect external tools (like Portainer or custom monitors) to the stack, ensure they are attached to the `m3tal` network.
*   **Gateway**: Traefik acts as the entrypoint for all traffic on port `8080`.

### Accessing the Services
Once the stack is up, use these local URLs (via Traefik):
*   **Dashboard**: `http://m3tal.localhost:8080`
*   **Backend API**: `http://api.localhost:8080`
*   **Traefik Admin**: `http://traefik.localhost:8080` (Internal only)

---

## ✅ Verification & Access
## ✅ Verification & Access

### Service Endpoints & Firewall
| Service | Local URL | Port | Firewall Required |
| :--- | :--- | :--- | :--- |
| **Dashboard** | `http://m3tal.localhost:8080` | `8080` | **Yes** (External Access) |
| **Backend API** | `http://api.localhost:8080` | `8080` | **No** (Internal Only) |
| **Traefik Web** | `http://localhost:8080` | `8080` | **Yes** (HTTP) |
| **Traefik SSL** | `http://localhost:443` | `443` | **Yes** (HTTPS) |

### SSL & External Access
M3TAL uses **Cloudflare Tunnel (`cloudflared`)** for secure external access.
*   **SSL**: Terminated at the Cloudflare edge. No local certificate management required.
*   **Custom Domains**: To expose a service, update the `DOMAIN` variable:
    ```bash
    ./m3tal config set DOMAIN yourdomain.com
    ```

### Initial Login
Run the following to set your admin password:
```bash
./m3tal dashpass admin yourpassword
```

### Configuration Management
Manage your `.env` variables directly via the CLI:
```bash
# List all variables
./m3tal config list

# Set a variable (e.g., your domain)
./m3tal config set DOMAIN mydomain.com

# Get a specific value
./m3tal config get API_PORT
```

---

## 🧭 Troubleshooting

- **Logs**: Always check the orchestrator logs first: `./m3tal list`
- **Manual Debug**: If necessary, check container logs directly: `docker logs m3tal-dashboard`
- **Clean-up**: To stop the entire stack, run `make down` or `./m3tal down`.
- **Docker on Windows (WSL)**: If you see `failed to connect to the docker API at npipe:////./pipe/dockerDesktopLinuxEngine`:
  1. Ensure **Docker Desktop** is running.
  2. Check "Settings > General > Use the WSL 2 based engine".
  3. Verify "Settings > Resources > WSL Integration" is enabled for your distribution.
  4. Try restarting Docker Desktop.

---

*M3TAL Core - Built for Performance and Stability.*