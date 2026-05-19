# docs/COMMAND_REFERENCE.md

Welcome to the M3TAL Ecosystem Command Reference. This document provides a comprehensive cheat-sheet for interacting with your M3TAL system via the command-line interface (CLI). M3TAL is designed to streamline the management of your self-hosted services using Docker Compose, providing a unified entry point for configuration, deployment, and monitoring.

All commands are executed via the `m3tal` binary, which acts as the central control plane for the M3TAL API daemon and your Dockerized applications.

---

## M3TAL System Filesystem Contract

Understanding the M3TAL filesystem is crucial for effective management and troubleshooting. These paths are foundational to how M3TAL operates:

| Path                        | Purpose                                                                                                                                                                                                                                                                                             |
| :-------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary Configuration File**. This file stores all environment variables that M3TAL uses to configure itself and the Docker Compose stacks it manages. It is the single source of truth for M3TAL's runtime parameters and is managed primarily by `m3tal config wizard`.                         |
| `/var/lib/m3tal/state.db`   | **SQLite State Database**. An automatically created SQLite database used by the `m3tal-api` daemon to store internal state, such as service statuses, configurations, and other operational data. Do not modify directly.                                                                           |
| `/opt/m3tal/stack/`         | **Canonical Stack Directory**. This is the core directory where M3TAL stores its Docker Compose files (`*-compose.yml`) and related configuration, including Traefik's dynamic configuration files (`dynamic/`). While M3TAL operates on this path internally, users typically interact with its symlink. |
| `/docker`                   | **User-Facing Stack Symlink**. This path is a symbolic link to `/opt/m3tal/stack/`. It is the recommended and user-facing directory for placing all your Docker Compose files (`my-app-compose.yml`, `another-service-compose.yml`, etc.) and any additional Traefik dynamic configuration.        |
| `/docker/users.json`        | **Dashboard Credential Store**. This file stores the hashed credentials for accessing the M3TAL Dashboard. It is managed exclusively by the `m3tal dashpass` command.                                                                                                                                |

---

## M3TAL CLI Command Reference

### Core System Management

#### `sudo m3tal`

Opens the interactive M3TAL TUI (Terminal User Interface) Control Center. This provides a user-friendly, numbered menu for common operations such as managing stacks, configuring M3TAL, and viewing system status. Requires `sudo` as it interacts with core system services and Docker.

```bash
sudo m3tal
# Output:
# M3TAL TUI Control Center
# ------------------------
# 1. Manage Stacks
# 2. Configure M3TAL
# 3. System Status
# 4. Dashboard Management
# 5. Exit
# Enter your choice:
```

#### `m3tal init`

Initializes the M3TAL configuration by generating the primary `/etc/m3tal/.env` file from default values. This command should be run on a first-time installation to ensure a baseline configuration exists.

```bash
m3tal init
# Output:
# [INFO] Generating default /etc/m3tal/.env configuration file...
# [INFO] .env file successfully created at /etc/m3tal/.env.
# [INFO] Run 'm3tal config wizard' to customize your setup.
```

#### `m3tal doctor`

Performs a pre-flight health check of your M3TAL system. It verifies Docker connectivity, validates the `/etc/m3tal/.env` configuration for common issues, and checks for port availability, ensuring your system is ready for operation.

```bash
m3tal doctor
# Output:
# [INFO] M3TAL Doctor: Performing system health check...
# [SUCCESS] Docker daemon is running and accessible.
# [SUCCESS] /etc/m3tal/.env is valid and readable.
# [SUCCESS] Port 80 (Traefik HTTP) is available.
# [SUCCESS] Port 8080 (M3TAL API) is available.
# [SUCCESS] Port 8082 (M3TAL Dashboard) is available.
# [INFO] M3TAL system appears healthy.
```

### Configuration Management (`m3tal config`)

M3TAL's configuration is primarily managed through the `/etc/m3tal/.env` file. These commands provide tools to interact with it.

#### `m3tal config wizard`

Launches an interactive, step-by-step wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended way to manage your M3TAL environment variables.

```bash
m3tal config wizard
# Output:
# M3TAL Configuration Wizard
# --------------------------
# Welcome! Let's configure your M3TAL system.
#
# Current value for DOMAIN (e.g., yourdomain.com): localhost
# Enter new value (or press Enter to keep current): myhome.arpa
#
# Current value for DASHBOARD_EXPOSE_MODE (local or traefik): local
# Enter new value (or press Enter to keep current): traefik
# ... (continues for other variables)
# [SUCCESS] Configuration saved to /etc/m3tal/.env.
```

#### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env` to a specified value. This command allows for quick, non-interactive adjustments.

```bash
m3tal config set DOMAIN mydomain.net
# Output:
# [INFO] Setting DOMAIN=mydomain.net in /etc/m3tal/.env
# [SUCCESS] Configuration updated. Remember to restart affected services.

m3tal config set PUID 1001
# Output:
# [INFO] Setting PUID=1001 in /etc/m3tal/.env
# [SUCCESS] Configuration updated. Remember to restart affected services.
```

#### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

```bash
m3tal config get DASHBOARD_EXPOSE_MODE
# Output:
# traefik

m3tal config get BASE_STORAGE_PATH
# Output:
# /mnt/data
```

#### `m3tal config scan`

Scans all active Docker Compose files (`*-compose.yml`) located in `/docker/` and lists all environment variables referenced within them, including their default values if available. This helps identify all potential configuration points across your entire M3TAL ecosystem.

```bash
m3tal config scan
# Output:
# Environment Variables Across All Stacks:
# ----------------------------------------
# DOMAIN (default: localhost)
# DASHBOARD_PORT (default: 8082)
# PUID (default: 1000)
# PGID (default: 1000)
# TZ (default: America/Denver)
# BASE_STORAGE_PATH (default: ./data)
# ...
```

#### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file. Useful for reviewing your complete M3TAL configuration at a glance.

```bash
m3tal config list
# Output:
# # M3TAL System Configuration
# DASHBOARD_PORT=8082
# DASHBOARD_EXPOSE_MODE=traefik
# HTTP_PORT=8080
# STATE_DIR=./state
# LOG_LEVEL=info
# DASHBOARD_SECRET=my_super_secret_key
# API_TOKEN=my_api_token_here
# ADMIN_PASSWORD=admin_pass
# DOMAIN=myhome.arpa
# ...
```

### Dashboard Management (`m3tal dash`)

Commands specifically for managing the M3TAL Dashboard container and its access.

#### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, the command becomes interactive, prompting you for both. User credentials are stored in `/docker/users.json`.

```bash
# Interactive mode
m3tal dashpass
# Output:
# Enter username: admin
# Enter new password for admin:
# Confirm new password:
# [SUCCESS] Password for user 'admin' updated in /docker/users.json.

# Non-interactive mode
m3tal dashpass admin new_strong_password
# Output:
# [SUCCESS] Password for user 'admin' updated in /docker/users.json.
```

#### `m3tal dash up`

Pulls the latest M3TAL Dashboard Docker Compose configuration files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub, then starts the dashboard container with the appropriate override based on `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.

```bash
m3tal dash up
# Output:
# [INFO] Pulling latest M3TAL Dashboard compose files from GitHub...
# [SUCCESS] Dashboard compose files updated.
# [INFO] DASHBOARD_EXPOSE_MODE is set to 'traefik'. Using m3tal-compose.traefik.yml.
# [INFO] Starting m3tal-dashboard container...
# [INFO] Docker compose up -d m3tal-dashboard...
# [SUCCESS] M3TAL Dashboard started successfully.
# Access the dashboard at http://dash.myhome.arpa
```

#### `m3tal dash down`

Stops and removes the M3TAL Dashboard container.

```bash
m3tal dash down
# Output:
# [INFO] Stopping m3tal-dashboard container...
# [INFO] Docker compose down m3tal-dashboard...
# [SUCCESS] M3TAL Dashboard stopped and removed.
```

#### `m3tal dash restart`

Restarts the M3TAL Dashboard container. This is useful after changing configuration variables that affect the dashboard.

```bash
m3tal dash restart
# Output:
# [INFO] Restarting m3tal-dashboard container...
# [INFO] Docker compose restart m3tal-dashboard...
# [SUCCESS] M3TAL Dashboard restarted successfully.
```

#### `m3tal dash logs`

Streams the logs from the M3TAL Dashboard container, providing real-time output for debugging or monitoring.

```bash
m3tal dash logs
# Output:
# Attaching to m3tal-dashboard
# m3tal-dashboard    | [2023-10-27 10:00:01] INFO: Dashboard started on 0.0.0.0:8082
# m3tal-dashboard    | [2023-10-27 10:00:05] DEBUG: API request to http://host.docker.internal:8080/api/status
# ... (continues streaming logs)
```

#### `m3tal dash status`

Shows the current status of the M3TAL Dashboard container (e.g., running, stopped, exited).

```bash
m3tal dash status
# Output:
# Container: m3tal-dashboard
# Status: running
# Image: ghcr.io/jakej985-rgb/m3tal-godash:debug
# Ports: 8082/tcp (via Traefik)
```

### Stack Management

M3TAL uses Docker Compose to manage all your services. The `m3tal up` and `m3tal down` commands provide a unified way to control all your deployed stacks. All compose files should be placed in the `/docker/` directory.

#### `m3tal up`

Runs `docker compose up -d` across *all* `*-compose.yml` files found in `/docker/`. This command starts or recreates all services defined in your M3TAL ecosystem, ensuring they are running in detached mode. This includes M3TAL's internal `routing-compose.yml` (Traefik, Cloudflared) and any custom stacks you've added.

```bash
m3tal up
# Output:
# [INFO] Running 'docker compose up -d' for all stacks in /docker/...
# [INFO] Processing /docker/routing-compose.yml...
# [INFO] Processing /docker/m3tal-compose.yml...
# [INFO] Processing /docker/my-app-compose.yml...
# ...
# [SUCCESS] All M3TAL stacks started successfully.
```

#### `m3tal down`

Runs `docker compose down` across *all* `*-compose.yml` files found in `/docker/`. This command stops and removes all containers, networks, and volumes defined by your M3TAL stacks.

```bash
m3tal down
# Output:
# [INFO] Running 'docker compose down' for all stacks in /docker/...
# [INFO] Processing /docker/my-app-compose.yml...
# [INFO] Processing /docker/m3tal-compose.yml...
# [INFO] Processing /docker/routing-compose.yml...
# ...
# [SUCCESS] All M3TAL stacks stopped and removed.
```

#### `m3tal logs`

Streams aggregated logs from all currently running Docker Compose stacks. This provides a consolidated view of all service activity, useful for overall system monitoring and debugging.

```bash
m3tal logs
# Output:
# Attaching to m3tal-dashboard, traefik, my-app-api, my-app-db
# traefik          | time="2023-10-27T10:05:01Z" level=info msg="Configuration loaded from file."
# m3tal-dashboard  | [2023-10-27 10:05:02] INFO: Dashboard refresh triggered.
# my-app-api       | [2023-10-27 10:05:03] INFO: API request received from 172.18.0.3
# traefik          | time="2023-10-27T10:05:04Z" level=debug msg="Provider connection established with docker."
# ... (continues streaming logs from all services)
```

---

## M3TAL Dashboard Access Modes

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. Understanding these modes is critical for correctly accessing your dashboard.

### Mode 1: `local` (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
*   **Mechanism**: Uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's internal port 8082 to the host's `DASHBOARD_PORT` (default 8082).
*   **Access Method**: Access the dashboard directly via your host's IP address or `localhost` on the specified port.
*   **Example**: If your host IP is `192.168.1.100` and `DASHBOARD_PORT=8082`, navigate to `http://192.168.1.100:8082` or `http://localhost:8082` in your web browser.
*   **Best For**: LAN-only setups, first-time installations, local development, or when you don't need Traefik for routing the dashboard.
*   **Traefik Requirement**: No Traefik required for dashboard access in this mode.

### Mode 2: `traefik`

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
*   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override, which adds specific Traefik labels to the dashboard container. Traefik (running via `routing-compose.yml`) then discovers these labels and routes incoming traffic from `dash.${DOMAIN}` to the dashboard container's internal port 8082.
*   **Access Method**: Access the dashboard via its designated subdomain, typically `dash.${DOMAIN}` (where `DOMAIN` is set in `/etc/m3tal/.env`).
*   **Example**: If `DOMAIN=myhome.arpa`, navigate to `http://dash.myhome.arpa` in your web browser.
*   **Best For**: Domain-based access, centralizing all services behind Traefik, or when you need secure (HTTPS) access via Traefik.
*   **Traefik Requirement**: Traefik *must* be running (`m3tal up` will ensure this as `routing-compose.yml` is always included).

To switch between modes, use `m3tal config set DASHBOARD_EXPOSE_MODE [local|traefik]` and then run `m3tal dash restart` (or `m3tal dash up` if not running).

---

## Systemd Service Management

The M3TAL API daemon, a Go binary, runs as a systemd service (`m3tal-api.service`) on your host machine, listening on port 8080. This daemon is crucial for the M3TAL CLI and Dashboard to interact with Docker and the internal state database.

### Check Service Status

To verify if the M3TAL API daemon is running and healthy:

```bash
systemctl status m3tal-api
# Output:
# ● m3tal-api.service - M3TAL API Daemon
#      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
#      Active: active (running) since Fri 2023-10-27 09:30:00 UTC; 1h ago
#    Main PID: 12345 (m3tal-api)
#       Tasks: 8 (limit: 9386)
#      Memory: 15.2M
#         CPU: 100ms
#      CGroup: /system.slice/m3tal-api.service
#              └─12345 /usr/bin/m3tal-api --config /etc/m3tal
# ...
```

### Stream Service Logs

To view real-time logs from the M3TAL API daemon for debugging:

```bash
journalctl -u m3tal-api -f
# Output:
# Oct 27 09:30:00 hostname systemd[1]: Started M3TAL API Daemon.
# Oct 27 09:30:00 hostname m3tal-api[12345]: [INFO] M3TAL API Daemon starting on port 8080.
# Oct 27 09:30:01 hostname m3tal-api[12345]: [INFO] Connected to Docker socket.
# Oct 27 09:30:05 hostname m3tal-api[12345]: [DEBUG] API endpoint /api/status accessed by 127.0.0.1
# ... (continues streaming logs)
```

### Restart the Service

If you make changes to M3TAL's core configuration (e.g., in `/etc/m3tal/.env`) that affect the API daemon (like `HTTP_PORT` or `LOG_LEVEL`), you might need to restart it:

```bash
sudo systemctl restart m3tal-api
# Output:
# [SUCCESS] M3TAL API Daemon restarted.
```

---

## Docker Direct Commands (Fallback)

M3TAL provides a high-level abstraction over Docker Engine and Docker Compose V2. However, for advanced users or troubleshooting, you can always interact with your Docker containers and stacks directly using `docker` and `docker compose` commands.

M3TAL's `up`, `down`, `logs` commands are essentially wrappers around `docker compose` calls to the `/docker/` directory.

**Important Note:** When using direct Docker Compose commands, navigate to the `/docker/` directory first, as compose files are typically relative to the execution directory.

### Listing all running containers

```bash
docker ps
# Output:
# CONTAINER ID   IMAGE                                COMMAND                  CREATED         STATUS         PORTS                                      NAMES
# a1b2c3d4e5f0   ghcr.io/jakej985-rgb/m3tal-godash    "python3 server.py"      2 hours ago     Up 2 hours     8082/tcp                                   m3tal-dashboard
# f1e2d3c4b5a0   traefik:latest                       "/entrypoint.sh --pr…"   2 hours ago     Up 2 hours     0.0.0.0:80->80/tcp, 127.0.0.1:8081->8080/tcp   traefik
# 0123456789ab   myorg/my-app:latest                  "python app.py"          1 hour ago      Up 1 hour      8000/tcp                                   my-app-api
```

### Direct `docker compose` for a specific stack

If you only want to manage a single stack, say `my-app-compose.yml`, without affecting others:

```bash
cd /docker/

# Start only 'my-app' stack
docker compose -f my-app-compose.yml up -d

# Stop only 'my-app' stack
docker compose -f my-app-compose.yml down

# View logs for 'my-app' stack
docker compose -f my-app-compose.yml logs -f
```

### Rebuilding a specific service

To rebuild and restart a single service within a stack (e.g., `my-app-api` in `my-app-compose.yml`):

```bash
cd /docker/
docker compose -f my-app-compose.yml up -d --build my-app-api
```

---

## Deployment Lifecycle — Day 2 Operations

M3TAL simplifies the process of adding and managing new Docker Compose-based applications.

### Installing a New Stack

1.  **Place your Compose File**: Copy your Docker Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory. This directory is symlinked from `/opt/m3tal/stack/`.
    ```bash
    sudo cp ~/my-stack-compose.yml /docker/
    ```
2.  **Configure Environment Variables**: If your new stack requires specific environment variables (e.g., `APP_PORT`, `DB_PASSWORD`), ensure they are defined in `/etc/m3tal/.env`. You can use `m3tal config wizard` for an interactive setup or `m3tal config set KEY value` for direct modification.
    ```bash
    m3tal config set APP_PORT 8000
    m3tal config wizard # Or run the wizard to review and set all necessary variables
    ```
3.  **Start All Stacks**: Run `m3tal up` to start all Docker Compose stacks, including your newly added one. M3TAL will automatically detect and include any `*-compose.yml` file in `/docker/`.
    ```bash
    m3tal up
    # Output will show your new stack being brought up.
    ```

---

## Traefik Routing Architecture

Traefik acts as the central reverse proxy for M3TAL, deployed via `routing-compose.yml` in `/docker/`. It's responsible for routing external requests to the correct internal Docker containers or host-local services.

*   **Entry Points**: Traefik binds to host port 80 (HTTP) as its primary entry point.
*   **Service Discovery**: It automatically discovers services running within the Docker `proxy` network by reading Docker labels (e.g., for the M3TAL Dashboard when `DASHBOARD_EXPOSE_MODE=traefik`).
*   **File Provider**: Traefik also loads dynamic configuration from `/docker/dynamic/` (which is `/opt/m3tal/stack/dynamic/`). This allows M3TAL to configure routing for host-local services like the M3TAL API daemon.
    *   **M3TAL API Routing**: A file like `/docker/dynamic/api.yml` routes `api.${DOMAIN}` to `http://host.docker.internal:8080` (the M3TAL API daemon running on the host). This `host.docker.internal` alias is critical for containers to reach the host machine.
*   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik labels on the `m3tal-dashboard` container (`m3tal-compose.traefik.yml`) route `dash.${DOMAIN}` to the dashboard.

This architecture allows for flexible, domain-based access to both containerized and host-local M3TAL components.

---

## Port Map

The following ports are used by M3TAL components and services:

| Port | Service                            | Access                                                    |
| :--- | :--------------------------------- | :-------------------------------------------------------- |
| 80   | Traefik HTTP entry point           | Public (when Traefik is running and configured)           |
| 8080 | M3TAL API daemon (Go binary)       | Host-local only (accessed by dashboard/Traefik via `host.docker.internal`) |
| 8081 | Traefik dashboard (admin UI)       | Host-local only (default, can be configured)              |
| 8082 | M3TAL Dashboard (Python/Flask) | Direct port (local mode) or via Traefik (traefik mode)    |

---

## APT Installation

To install the M3TAL CLI binary and systemd service on your Debian-based system, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```