# 🚀 M3TAL Media Server (v1.4)

M3TAL is a high-performance media server control plane built with **Go 1.26** and **Python 3.10**. It features a native Go orchestrator that manages a multi-stack Docker environment for media services, routing, and system monitoring.

## 🧠 Architecture: The M3TAL Orchestrator

The `./m3tal` binary is more than a CLI; it is a **native Go orchestrator**. When you run `./m3tal up`, the following happens:
1. It imports `pkg/orchestrator`.
2. It programmatically executes `docker compose` across three specialized stacks:
   - `network-compose.yml`: Virtual networking and DNS.
   - `routing-compose.yml`: Traefik gateway and SSL.
   - `m3tal-compose.yml`: Core API, Dashboard, and Agents.

---

## 🛠️ Prerequisites

1. **Docker Engine**: v20.10+
2. **Docker Compose**: v2.0+
3. **Go 1.26+**: Required to build the core binaries.
4. **Storage**: A mount point for media (Default: `/mnt`). 
   - *Note*: Ensure this directory exists: `sudo mkdir -p /mnt && sudo chown $USER:$USER /mnt`.

---

## 🚀 Quick Start (Linux/WSL)

### 1. Build the Platform
Standardize your binaries by building the CLI and the API:
```bash
go build -o m3tal ./cmd/m3tal
go build -o m3tal-api ./cmd/api
```

### 2. Initialize Environment
Prepare the directory structure and generate your secret keys:
```bash
cp .env.example .env
# Edit .env with your domain and storage paths
```

### 3. Start the Stack
Initialize the orchestrator and pull images:
```bash
./m3tal pull
./m3tal up
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

Once the orchestrator reports success, verify your deployment:

1. **Dashboard**: Access at [http://localhost:8082](http://localhost:8082)
2. **API Health**: `curl http://localhost:5050/api/metrics`
3. **Container Status**: `./m3tal list` or `docker compose ps`

### Initial Login
Run the following to set your admin password:
```bash
./m3tal dashpass admin yourpassword
```

---

## 🧭 Troubleshooting

- **Permissions**: If the stack fails to start, verify that your user is in the `docker` group.
- **Port Conflicts**: Ensure ports `80`, `443`, `8082`, and `5050` are not in use.
- **Logs**: View real-time logs with `docker compose -f source/m3tal-stack/m3tal-compose.yml logs -f`.

---

*M3TAL Core - Built for Performance and Stability.*