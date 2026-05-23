I am DocSmith, the M3TAL Ecosystem Documentation Architect. Below is your comprehensive CLI cheat-sheet for managing your M3TAL instance.

---

# M3TAL Command Reference

This document provides a complete reference for the M3TAL Command Line Interface (CLI). The `m3tal` binary is the primary tool for interacting with, configuring, and maintaining your M3TAL ecosystem.

## M3TAL Ecosystem Overview

M3TAL is designed as a modular, containerized application platform. It leverages Docker and Docker Compose for service orchestration and a central Go API daemon for unified control.

### Core Components:
- **CLI binary (`/usr/bin/m3tal`)**: Your single entrypoint for all M3TAL operations.
- **API daemon (`m3tal-api.service`)**: A Go binary running as a `systemd` service (port `8080` internally) that manages Docker interactions, state, and API routes.
- **Dashboard container (`m3tal-dashboard`)**: A Python/Flask container (internal port `8082`) that provides a web-based UI, communicating with the API daemon.
- **Traefik Gateway (`routing-compose.yml`)**: A reverse proxy container exposing services by domain name on port `80`.
- **Cloudflared (`routing-compose.yml`)**: An optional Cloudflare tunnel container for secure, zero-config internet access.

### Filesystem Contract

M3TAL adheres to a strict filesystem layout for configuration and data persistence:

| Path                        | Purpose                                                      |
| :-------------------------- | :----------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file.** Contains all environment variables for M3TAL and its stacks. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains core and user-defined Docker Compose files (`*-compose.yml`) and Traefik dynamic configuration. |
| `/docker`                   | **User-facing symlink to `/opt/m3tal/stack/`.** All stack operations and user-added compose files reside here. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.     |

### Dashboard Access Modes

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **`local` mode (Default)**
    *   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
    *   **Mechanism**: Uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's port `8082` to the host's `DASHBOARD_PORT` (default `8082`).
    *   **Access**: `http://HOST_IP:8082` or `http://localhost:8082`
    *   **Requirements**: No Traefik required. Ideal for LAN-only setups, initial configuration, or testing.

2.  **`traefik` mode**
    *   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
    *   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard container. Traefik then routes `dash.${DOMAIN}` to the dashboard on port `8082`.
    *   **Access**: `http://dash.YOUR_DOMAIN` (Requires Traefik to be running via `m3tal up`).
    *   **Requirements**: Traefik must be deployed and correctly configured. Best for domain-based access behind a reverse proxy.

### Docker / Compose Runtime

M3TAL heavily relies on **Docker Engine** and **Docker Compose V2**. These are hard dependencies for the M3TAL ecosystem.

*   The `m3tal up` command orchestrates all `*-compose.yml` files found in the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the dashboard container:
    1.  It pulls the latest dashboard-specific compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
    2.  It reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
    3.  It then starts the dashboard container using the appropriate compose override file.
*   **User Stacks**: To add a new service, simply place its Docker Compose file (e.g., `my-app-compose.yml`) into `/docker/`. M3TAL will automatically pick it up with `m3tal up`.

### Traefik Routing Architecture

Traefik runs as a container (`routing-compose.yml`) and serves as the primary ingress point on port `80` (and `443` if HTTPS is enabled).

*   It automatically discovers services by reading Docker labels on other containers (e.g., `m3tal-dashboard` in `traefik` mode).
*   It also loads dynamic routing configurations from `/docker/dynamic/` (using the file provider), which enables hot-reloading.
*   Example: `api.DOMAIN` routes to the M3TAL API daemon (`http://host.docker.internal:8080`) via `dynamic/api.yml`.

### Port Map

| Port | Service                               | Access                                                               |
| :--- | :------------------------------------ | :------------------------------------------------------------------- |
| 80   | Traefik HTTP Entry Point              | Public (if Traefik is running and exposed)                           |
| 8080 | M3TAL API Daemon (Go)                 | Host-local only (accessed by dashboard and Traefik via `host.docker.internal`) |
| 8081 | Traefik Dashboard                     | Host-local only (e.g., `http://localhost:8081`)                      |
| 8082 | M3TAL Dashboard (internal container)  | Direct port (local mode) or via Traefik (traefik mode)               |

## Installation

To get M3TAL installed on your system, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## M3TAL CLI Command Reference

### Interactive Control Center

`sudo m3tal`

Opens the interactive M3TAL TUI (Text User Interface) Control Center. This provides a numbered menu for common operations and system status.
**Example Usage:**
```bash
sudo m3tal
# You will see a menu similar to:
# -------------------------------------
# M3TAL Control Center
# -------------------------------------
# 1. System Status
# 2. Start All Stacks
# 3. Stop All Stacks
# 4. View Aggregated Logs
# 5. Dashboard Management
# 6. Configuration Wizard
# 0. Exit
# Enter your choice:
```

### System Initialization & Health

`m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from default values. This command should be run on the first installation or if the `.env` file is missing.
**Example Usage:**
```bash
m3tal init
# Output:
# Initializing M3TAL environment...
# /etc/m3tal/.env created successfully with default values.
# Please consider running 'm3tal config wizard' to customize your setup.
```

`m3tal doctor`

Performs a pre-flight health check of the M3TAL ecosystem. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability.
**Example Usage:**
```bash
m3tal doctor
# Output:
# M3TAL System Doctor Report:
#   ✓ Docker Daemon: Running and reachable.
#   ✓ /etc/m3tal/.env: Validated successfully.
#   ✓ Required Ports: 80, 8080, 8081, 8082 are available.
#   ✓ Network 'proxy': Exists.
# All essential components are healthy.
```

### Configuration Management (`m3tal config`)

`m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended way to modify your M3TAL environment variables.
**Example Usage:**
```bash
m3tal config wizard
# Output:
# M3TAL Configuration Wizard
#
# Welcome! Let's configure your M3TAL environment.
# Current value for DOMAIN (default: localhost):
# Do you want to set a custom DOMAIN? (y/N): y
# Enter new value for DOMAIN: myhomelab.local
# ... (continues for other variables) ...
# Configuration saved to /etc/m3tal/.env.
```

`m3tal config set KEY VALUE`

Sets a single environment variable directly in `/etc/m3tal/.env`. Use this for quick, non-interactive adjustments.
**Example Usage:**
```bash
m3tal config set DOMAIN mydomain.com
# Output:
# Successfully set DOMAIN=mydomain.com in /etc/m3tal/.env.
```
```bash
m3tal config set DASHBOARD_EXPOSE_MODE traefik
# Output:
# Successfully set DASHBOARD_EXPOSE_MODE=traefik in /etc/m3tal/.env.
```

`m3tal config get KEY`

Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.
**Example Usage:**
```bash
m3tal config get PUID
# Output:
# 1000
```
```bash
m3tal config get DASHBOARD_SECRET
# Output:
# change_me_immediately
```

`m3tal config scan`

Lists all known environment variables and their default values across all recognized M3TAL stacks. This helps identify all possible configuration options.
**Example Usage:**
```bash
m3tal config scan
# Output:
# Listing all known environment variables across stacks:
# - API_TOKEN (default: change_me_api_token)
# - BASE_STORAGE_PATH (default: ./data)
# - CONFIG_PATH (default: ./data/config)
# - DASHBOARD_EXPOSE_MODE (default: local)
# - DASHBOARD_PORT (default: 8082)
# - DASHBOARD_SECRET (default: change_me_immediately)
# - DEBUG_MODE (default: false)
# - DOMAIN (default: localhost)
# - DOWNLOADS_PATH (default: ./data/downloads)
# - HTTP_PORT (default: 8080)
# - LOCAL_IP (default: 127.0.0.1)
# - LOG_LEVEL (default: info)
# - MEDIA_PATH (default: ./data/media)
# - METRICS_ENABLED (default: true)
# - NETWORK_NAME (default: m3tal)
# - PGID (default: 1000)
# - PUID (default: 1000)
# - STATE_DIR (default: ./state)
# - TRAEFIK_DASHBOARD_PORT (default: 8080)
# - TRAEFIK_WEB_PORT (default: 80)
# - TRAEFIK_WEBHTTPS_PORT (default: 443)
# - TZ (default: America/Denver)
# - VPN_PASSWORD (default: password)
# - VPN_USER (default: user)
```

`m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file.
**Example Usage:**
```bash
m3tal config list
# Output:
# # M3TAL Environment Variables
# PUID=1000
# PGID=1000
# TZ=America/Denver
# DASHBOARD_PORT=8082
# DASHBOARD_EXPOSE_MODE=local
# DOMAIN=mydomain.com
# DASHBOARD_SECRET=supersecretpassword
# API_TOKEN=mysecureapitoken
# ... (rest of .env file) ...
```

### Dashboard Management (`m3tal dash`)

`m3tal dashpass [username] [password]`

Updates a user's password for the M3TAL Dashboard, stored in `/docker/users.json`. If `username` and `password` are omitted, the command becomes interactive.
**Example Usage (Interactive):**
```bash
m3tal dashpass
# Output:
# Enter username to update: admin
# Enter new password for admin:
# Confirm new password:
# Password for 'admin' updated successfully in /docker/users.json.
```
**Example Usage (Direct):**
```bash
m3tal dashpass admin newSecurePassword123!
# Output:
# Password for 'admin' updated successfully in /docker/users.json.
```

`m3tal dash up`

Pulls the latest dashboard compose configuration files from GitHub, then starts the `m3tal-dashboard` container. This command intelligently uses either `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` based on the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.
**Example Usage:**
```bash
m3tal dash up
# Output:
# Pulling latest dashboard compose files...
# Dashboard files updated successfully.
# Starting m3tal-dashboard in 'local' mode...
# [Container startup logs from Docker Compose]
# m3tal-dashboard started successfully.
```

`m3tal dash down`

Stops and removes the `m3tal-dashboard` container.
**Example Usage:**
```bash
m3tal dash down
# Output:
# Stopping m3tal-dashboard...
# [Container shutdown logs from Docker Compose]
# m3tal-dashboard stopped and removed.
```

`m3tal dash restart`

Restarts the `m3tal-dashboard` container. This is useful after making configuration changes that affect the dashboard.
**Example Usage:**
```bash
m3tal dash restart
# Output:
# Restarting m3tal-dashboard...
# [Container restart logs from Docker Compose]
# m3tal-dashboard restarted successfully.
```

`m3tal dash logs`

Streams the logs from the `m3tal-dashboard` container. This is invaluable for debugging dashboard issues.
**Example Usage:**
```bash
m3tal dash logs -f
# Output:
# Attaching to m3tal-dashboard
# m3tal-dashboard |  * Running on http://0.0.0.0:8082
# m3tal-dashboard |  * Debug mode: off
# m3tal-dashboard | 172.18.0.1 - - [21/Jul/2024 14:35:01] "GET /api/status HTTP/1.1" 200 -
# ... (continues streaming logs) ...
```

`m3tal dash status`

Shows the current status of the `m3tal-dashboard` container (e.g., running, exited, paused).
**Example Usage:**
```bash
m3tal dash status
# Output:
# Container 'm3tal-dashboard': running (Up 5 minutes)
```

### General Stack Management

`m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This command starts all M3TAL core services (routing, dashboard) and any user-defined stacks. It's often run with `-d` for detached mode.
**Example Usage:**
```bash
m3tal up -d
# Output:
# Creating network "proxy" with the default driver
# Creating m3tal-dashboard ... done
# Creating traefik        ... done
# Creating cloudflared    ... done
# Creating my-app         ... done
# All M3TAL stacks started in detached mode.
```

`m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This stops and removes all M3TAL core services and user-defined stacks.
**Example Usage:**
```bash
m3tal down
# Output:
# Stopping m3tal-dashboard ... done
# Stopping traefik        ... done
# Stopping cloudflared    ... done
# Stopping my-app         ... done
# Removing m3tal-dashboard ... done
# Removing traefik        ... done
# Removing cloudflared    ... done
# Removing my-app         ... done
# All M3TAL stacks stopped and removed.
```

`m3tal logs`

Streams aggregated logs from all currently running Docker containers managed by M3TAL. This provides a consolidated view for overall system monitoring. Use `-f` to follow logs in real-time.
**Example Usage:**
```bash
m3tal logs -f
# Output:
# Attaching to cloudflared, m3tal-dashboard, traefik
# traefik         | time="2024-07-21T14:40:05Z" level=info msg="Configuration loaded from file: /etc/traefik/traefik.yml"
# m3tal-dashboard | 172.18.0.1 - - [21/Jul/2024 14:40:06] "GET /api/status HTTP/1.1" 200 -
# cloudflared     | 2024-07-21T14:40:07Z INF Initializing Cloudflare Tunnel
# ... (continues streaming logs from all services) ...
```

---

## Advanced Management

### Systemd Service Management

The M3TAL API daemon, `m3tal-api`, runs as a `systemd` service. You can manage it using standard `systemctl` commands. This daemon is crucial as it powers the `m3tal` CLI and the Dashboard's interaction with Docker.

**Check Service Status:**
```bash
systemctl status m3tal-api
# Output:
# ● m3tal-api.service - M3TAL API Daemon
#      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
#      Active: active (running) since Sun 2024-07-21 14:30:00 UTC; 10min ago
#    Main PID: 1234 (m3tal-api)
#       Tasks: 8 (limit: 4915)
#      Memory: 15.6M
#         CPU: 125ms
#      CGroup: /system.slice/m3tal-api.service
#              └─1234 /usr/bin/m3tal-api
```

**Restart the API Daemon:**
```bash
sudo systemctl restart m3tal-api
# Output: (no direct output, but the service will restart)
```

**Stream API Daemon Logs:**
```bash
sudo journalctl -u m3tal-api -f
# Output:
# -- Journal begins at Sat 2024-07-20 10:00:00 UTC, ends at Sun 2024-07-21 14:45:00 UTC. --
# Jul 21 14:30:00 hostname systemd[1]: Started M3TAL API Daemon.
# Jul 21 14:30:00 hostname m3tal-api[1234]: [INFO] M3TAL API daemon starting...
# Jul 21 14:30:01 hostname m3tal-api[1234]: [INFO] Database initialized at /var/lib/m3tal/state.db
# Jul 21 14:30:01 hostname m3tal-api[1234]: [INFO] Listening on :8080
# ... (continues streaming API logs) ...
```

### Direct Docker Compose Fallback

The `m3tal` CLI commands abstract `docker compose` operations for convenience. However, you can always interact directly with Docker Compose, especially for troubleshooting or advanced scenarios. All M3TAL compose files, including core and user-defined stacks, reside in `/docker/` (which is a symlink to `/opt/m3tal/stack/`).

**Prerequisites:**
- Ensure you have `docker` and `docker compose` (V2) installed and configured.
- Navigate to the `/docker/` directory before running `docker compose` commands.

**1. Listing Docker Compose Files:**
To see which compose files M3TAL manages, look in `/docker/`:
```bash
ls -l /docker/
# Output:
# total 36
# -rw-r--r-- 1 root root  234 Jul 20 10:00 m3tal-compose.local.yml
# -rw-r--r-- 1 root root  987 Jul 20 10:00 m3tal-compose.yml
# -rw-r--r-- 1 root root  567 Jul 20 10:00 m3tal-compose.traefik.yml
# -rw-r--r-- 1 root root  890 Jul 20 10:00 routing-compose.yml
# -rw-r--r-- 1 root root 1234 Jul 21 12:00 my-custom-stack-compose.yml
# drwxr-xr-x 2 root root 4096 Jul 20 10:00 dynamic
# -rw-r--r-- 1 root root  234 Jul 20 10:00 users.json
```

**2. Starting Specific Stacks (similar to `m3tal up`):**
To start the core routing and dashboard stacks, explicitly specify their compose files:
```bash
cd /docker/
docker compose -f routing-compose.yml -f m3tal-compose.yml -f m3tal-compose.local.yml up -d
# Note: You'd use either m3tal-compose.local.yml OR m3tal-compose.traefik.yml based on your .env setting.
```
To start a user-defined stack:
```bash
cd /docker/
docker compose -f my-custom-stack-compose.yml up -d
```
To start ALL stacks, you would chain all `*-compose.yml` files.

**3. Stopping Specific Stacks (similar to `m3tal down`):**
```bash
cd /docker/
docker compose -f routing-compose.yml -f m3tal-compose.yml down
```

**4. Viewing Logs for a Specific Service:**
```bash
cd /docker/
docker compose -f routing-compose.yml logs -f traefik
# Output: (logs specifically for the Traefik container)
```

**5. Checking Status of Services:**
```bash
cd /docker/
docker compose -f routing-compose.yml -f m3tal-compose.yml ps
# Output:
# NAME                COMMAND                  SERVICE             STATUS              PORTS
# m3tal-dashboard     "python3 server.py"      m3tal-dashboard     running             0.0.0.0:8082->8082/tcp
# traefik             "/entrypoint.sh --pr…"   traefik             running             0.0.0.0:80->80/tcp, 127.0.0.1:8081->8080/tcp
# cloudflared         "/usr/local/bin/clou…"   cloudflared         running
```

---