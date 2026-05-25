# M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, I've compiled this comprehensive cheat-sheet to guide you through managing your M3TAL instance. The M3TAL CLI (`m3tal`) is your unified entry point for interacting with the M3TAL API daemon, configuring your system, and orchestrating your Dockerized applications.

## M3TAL Ecosystem Overview

M3TAL is a unified platform for self-hosting applications, built on Docker Compose and managed by a central Go API daemon. It provides a robust framework for deploying, configuring, and monitoring your services.

### Core Components:
*   **CLI (`/usr/bin/m3tal`)**: Your primary interface for all M3TAL operations.
*   **API Daemon (`m3tal-api.service`)**: A Go binary running as a systemd service on port `8080`. Manages Docker, state database, and API routes.
*   **M3TAL Dashboard (`m3tal-dashboard`)**: A Python/Flask container running internally on port `8082`. Communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik Gateway (`routing-compose.yml`)**: A reverse proxy container exposing services by domain name on port `80`. Uses a file provider for dynamic routing.
*   **Cloudflared (`routing-compose.yml`)**: An optional Cloudflare tunnel container for zero-config internet access.

### Filesystem Contract

The following paths are critical to the M3TAL ecosystem:

| Path                        | Purpose                                                                                                                                                                                               |
| :-------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Contains environment variables for M3TAL and your Docker stacks. **Managed by `m3tal config wizard`.**                                                                  |
| `/var/lib/m3tal/state.db`   | SQLite state database. Stores M3TAL's internal state, user data, and service information. Auto-created by the API daemon.                                                                             |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains core M3TAL Compose files (`m3tal-compose.yml`, `routing-compose.yml`, etc.) and Traefik dynamic configuration.                                                    |
| `/docker`                   | **User-facing symlink to `/opt/m3tal/stack/`**. This is where you place your custom `*-compose.yml` files for M3TAL to manage.                                                                     |
| `/docker/users.json`        | Dashboard credential store. Contains usernames and hashed passwords for dashboard access. **Managed by `m3tal dashpass`.**                                                                          |
| `/docker/dynamic/`          | Directory for Traefik's dynamic file provider configuration. M3TAL manages `api.yml` here for API daemon routing.                                                                                     |

## Dashboard Access Modes

The M3TAL Dashboard offers two modes for access, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Mode 1: `local` (Default)
*   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
*   **Mechanism**: M3TAL uses `m3tal-compose.local.yml` to add a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`). No Traefik required.
*   **Access**: `http://YOUR_HOST_IP:8082` or `http://localhost:8082` (if accessing from the host machine).
*   **Use Case**: Ideal for LAN-only setups, first-time users, or local development where a domain and Traefik are not yet configured.

### Mode 2: `traefik`
*   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
*   **Mechanism**: M3TAL uses `m3tal-compose.traefik.yml` to add Traefik labels to the dashboard container, allowing Traefik to route `dash.${DOMAIN}` to the dashboard on port `8082`. **Traefik must be running via `m3tal up` for this mode to work.**
*   **Access**: `http://dash.YOUR_DOMAIN` (e.g., `http://dash.mydomain.com`).
*   **Use Case**: Suited for domain-based setups, where the dashboard is one of many services exposed behind Traefik. Requires the `routing-compose.yml` stack to be active.

## M3TAL CLI Command Reference

### `sudo m3tal`

Opens the interactive TUI (Text User Interface) Control Center. This provides a user-friendly, menu-driven way to manage your M3TAL instance, including starting/stopping stacks, viewing logs, and configuring settings.

```bash
sudo m3tal
# Example Interaction:
# ------------------------------------
# |      M3TAL Control Center        |
# ------------------------------------
# | 1. System Status                 |
# | 2. Manage Stacks                 |
# | 3. Configuration Wizard          |
# | 4. View Logs                     |
# | 5. Update Dashboard Password     |
# | 6. Run Doctor Check              |
# | 0. Exit                          |
# ------------------------------------
# Enter your choice: _
```

### `m3tal init`

Generates the `/etc/m3tal/.env` configuration file from default values. This command is essential for a first-time installation to ensure your M3TAL instance has a basic, functional configuration.

```bash
m3tal init
# Example Output:
# INFO: .env file created at /etc/m3tal/.env with default values.
# INFO: Please run 'm3tal config wizard' to customize your setup.
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL system. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability to help diagnose common issues before starting services.

```bash
m3tal doctor
# Example Output:
# [SUCCESS] Docker daemon is running and accessible.
# [SUCCESS] /etc/m3tal/.env is present and valid.
# [SUCCESS] Port 8080 (M3TAL API) is available.
# [INFO] Port 80 (Traefik Web) is in use by another process. This might be expected if Traefik is already running.
# [WARNING] DASHBOARD_SECRET is set to default. Please change it with 'm3tal config wizard' or 'm3tal config set DASHBOARD_SECRET <new_secret>'.
# M3TAL system health check completed.
```

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring your `/etc/m3tal/.env` file. This is the recommended way to set up and modify your environment variables, ensuring proper validation and helpful prompts.

```bash
m3tal config wizard
# Example Interaction:
# Welcome to the M3TAL Configuration Wizard!
# Current DOMAIN: localhost. Enter new DOMAIN (e.g., mydomain.com): example.com
# Current DASHBOARD_EXPOSE_MODE: local. Select new DASHBOARD_EXPOSE_MODE (local/traefik): traefik
# Current API_TOKEN: change_me_api_token. Enter new API_TOKEN: my_strong_api_token_123
# ... (continues for other relevant variables)
# Configuration saved to /etc/m3tal/.env.
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env`. Useful for quick, specific changes without going through the full wizard.

```bash
m3tal config set DOMAIN mycoolserver.net
# Example Output:
# INFO: Successfully set DOMAIN=mycoolserver.net in /etc/m3tal/.env.
```

### `m3tal config get KEY`

Reads and displays the value of a specific environment variable from `/etc/m3tal/.env`.

```bash
m3tal config get PUID
# Example Output:
# 1000
```

### `m3tal config scan`

Lists all detected environment variables that are relevant to M3TAL across all configured stacks. This includes variables used by the API daemon, dashboard, and common compose files in `/docker/`.

```bash
m3tal config scan
# Example Output:
# Key                      Default Value      Description
# ---                      -------------      -----------
# API_TOKEN                change_me_api_token  API authentication token.
# DOMAIN                   localhost          Base domain for Traefik routing.
# DASHBOARD_PORT           8082               Port for the M3TAL Dashboard.
# PUID                     1000               User ID for container permissions.
# ... (more variables listed)
```

### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file.

```bash
m3tal config list
# Example Output:
# PUID=1000
# PGID=1000
# TZ=America/Denver
# DASHBOARD_PORT=8082
# DASHBOARD_EXPOSE_MODE=traefik
# DOMAIN=mycoolserver.net
# API_TOKEN=my_secure_api_token_123
# ...
```

### `m3tal dashpass [username] [password]`

Updates a user's password for the M3TAL Dashboard. If `username` and `password` are omitted, it will prompt for interactive input. This command manages the `/docker/users.json` file.

```bash
# Interactive mode (prompts for username and password)
m3tal dashpass
# Example Interaction:
# Enter username to update: admin
# Enter new password for admin: **********
# Confirm new password: **********
# INFO: Password for user 'admin' updated successfully.

# Non-interactive mode
m3tal dashpass admin SuperSecurePassword123!
# Example Output:
# INFO: Password for user 'admin' updated successfully.
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the M3TAL Dashboard container. This command intelligently selects the correct override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.

```bash
m3tal dash up
# Example Output:
# INFO: Pulling latest dashboard compose files from GitHub...
# INFO: DASHBOARD_EXPOSE_MODE is 'traefik'. Starting dashboard with Traefik override.
# [+] Running 1/0
#  ⠿ Container m3tal-dashboard  Started
# M3TAL Dashboard is now running.
# Access at: http://dash.mycoolserver.net (or http://HOST_IP:8082 if in local mode)
```

### `m3tal dash down`

Stops and removes the M3TAL Dashboard container.

```bash
m3tal dash down
# Example Output:
# [+] Stopping 1/0
#  ⠿ Container m3tal-dashboard  Stopped
# INFO: M3TAL Dashboard container stopped.
```

### `m3tal dash restart`

Restarts the M3TAL Dashboard container.

```bash
m3tal dash restart
# Example Output:
# [+] Restarting 1/0
#  ⠿ Container m3tal-dashboard  Restarted
# INFO: M3TAL Dashboard container restarted.
```

### `m3tal dash logs`

Streams logs from the M3TAL Dashboard container, useful for debugging and monitoring its activity.

```bash
m3tal dash logs
# Example Output (continuous stream):
# m3tal-dashboard | INFO: Starting M3TAL Dashboard server...
# m3tal-dashboard | INFO: Listening on http://0.0.0.0:8082
# m3tal-dashboard | INFO: Connected to M3TAL API at http://host.docker.internal:8080
# ...
```

### `m3tal dash status`

Shows the current status of the M3TAL Dashboard container (e.g., running, stopped, exited).

```bash
m3tal dash status
# Example Output:
# Container Name: m3tal-dashboard
# Status:         running (healthy)
# Image:          ghcr.io/jakej985-rgb/m3tal-godash:debug
# Ports:          8082/tcp (via Traefik, or 0.0.0.0:8082->8082/tcp if in local mode)
# Labels:         m3tal.stack=control-plane
```

### `m3tal up`

Starts all Docker Compose stacks defined by `*-compose.yml` files found in the `/docker/` directory. This includes core M3TAL services like Traefik (`routing-compose.yml`) and any custom user stacks.

```bash
m3tal up
# Example Output:
# INFO: Running docker compose up across all stacks in /docker/...
# [+] Running 3/0
#  ⠿ Network proxy  Created
#  ⠿ Container traefik          Started
#  ⠿ Container ollama           Started
#  ⠿ Container m3tal-dashboard  Started
# All M3TAL managed services are now running.
```

### `m3tal down`

Stops and removes all Docker Compose stacks defined by `*-compose.yml` files in the `/docker/` directory.

```bash
m3tal down
# Example Output:
# INFO: Running docker compose down across all stacks in /docker/...
# [+] Stopping 3/0
#  ⠿ Container ollama           Stopped
#  ⠿ Container m3tal-dashboard  Stopped
#  ⠿ Container traefik          Stopped
# All M3TAL managed services have been stopped.
```

### `m3tal logs`

Streams aggregated logs from all currently running Docker containers managed by M3TAL. This provides a consolidated view of your entire M3TAL ecosystem's activity.

```bash
m3tal logs
# Example Output (continuous stream):
# traefik         | time="2023-10-27T10:30:05Z" level=info msg="Configuration loaded from flags."
# m3tal-dashboard | INFO: Dashboard started successfully.
# ollama          | INFO: server listening on 0.0.0.0:11434
# ...
```

---

## Systemd Service Management

The M3TAL API daemon is managed as a systemd service, `m3tal-api.service`. This service is crucial for all M3TAL CLI operations, as the CLI communicates directly with it.

### Checking API Service Status

To check if the M3TAL API daemon is running and healthy:

```bash
systemctl status m3tal-api
# Example Output:
# ● m3tal-api.service - M3TAL API Daemon
#      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
#      Active: active (running) since Fri 2023-10-27 10:25:01 UTC; 5min ago
#    Main PID: 1234 (m3tal-api)
#       Tasks: 8 (limit: 4627)
#      Memory: 12.5M
#         CPU: 123ms
#      CGroup: /system.slice/m3tal-api.service
#              └─1234 /usr/bin/m3tal-api
# Oct 27 10:25:01 myhost systemd[1]: Started M3TAL API Daemon.
# Oct 27 10:25:01 myhost m3tal-api[1234]: INFO: M3TAL API started on :8080
```

### Restarting the API Service

If you modify `/etc/m3tal/.env` directly or encounter issues with the API, a restart might be necessary:

```bash
sudo systemctl restart m3tal-api
```

### Viewing API Service Logs

To stream real-time logs from the M3TAL API daemon:

```bash
sudo journalctl -u m3tal-api -f
# Example Output (continuous stream):
# Oct 27 10:35:10 myhost m3tal-api[1234]: INFO: M3TAL API started on :8080
# Oct 27 10:35:15 myhost m3tal-api[1234]: INFO: Dashboard password updated by CLI.
# ...
```

---

## Docker Direct Commands (Fallback)

M3TAL leverages Docker Engine and Docker Compose V2. While the `m3tal` CLI is the preferred interface, you can always use direct `docker compose` commands as a fallback or for advanced debugging.

**Important Note**:
*   `m3tal up` processes *all* `*-compose.yml` files in `/docker/`.
*   `m3tal dash up` specifically manages the dashboard, downloading its compose files and applying the correct override based on `DASHBOARD_EXPOSE_MODE`. The relevant compose files are located in `/opt/m3tal/stack/` (or `/docker/` via symlink).

### Start All Stacks (Equivalent to `m3tal up`)

To start all services defined by `*-compose.yml` files in `/docker/`:

```bash
# This command uses shell globbing to include all *-compose.yml files.
# Ensure you are in a directory where this glob resolves correctly or specify full paths.
docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/ollama-compose.yml up -d
# If you have many files, globbing can be convenient (shell dependent):
# docker compose -f /docker/*.yml up -d
```

### Stop All Stacks (Equivalent to `m3tal down`)

```bash
docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/ollama-compose.yml down
# docker compose -f /docker/*.yml down
```

### Start M3TAL Dashboard (Equivalent to `m3tal dash up`)

**Note:** You must specify the base `m3tal-compose.yml` and the correct override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

**For `DASHBOARD_EXPOSE_MODE=local`:**
```bash
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml up -d m3tal-dashboard
```

**For `DASHBOARD_EXPOSE_MODE=traefik`:**
```bash
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.traefik.yml up -d m3tal-dashboard
```

### Stop M3TAL Dashboard (Equivalent to `m3tal dash down`)

```bash
# The specific override file isn't strictly necessary for 'down', but harmless.
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml down m3tal-dashboard
```

### Restart M3TAL Dashboard (Equivalent to `m3tal dash restart`)

```bash
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml restart m3tal-dashboard
```

### View All Logs (Equivalent to `m3tal logs`)

```bash
# Views logs from all services defined in the specified compose files.
docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/ollama-compose.yml logs -f
# docker compose -f /docker/*.yml logs -f
```

---

## M3TAL System Architecture

### Deployment Lifecycle — Day 2 Operations

**Installing a new Docker Compose stack:**

1.  Place your new compose file (e.g., `my-service-compose.yml`) into the `/docker/` directory.
2.  Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to start your new stack alongside all existing M3TAL-managed services.

### Traefik Routing Architecture

Traefik (deployed as a container via `routing-compose.yml`) acts as the central reverse proxy:
*   **HTTP Entry Point**: Binds port 80 (and 443 if HTTPS is configured) on the host.
*   **Service Discovery**: Automatically detects services via Docker labels (e.g., for the dashboard when `DASHBOARD_EXPOSE_MODE=traefik`).
*   **Dynamic Configuration**: Loads additional routing rules from `/docker/dynamic/` (Traefik's file provider). For example, `dynamic/api.yml` routes `api.DOMAIN` to the M3TAL API daemon. This directory allows hot-reloading for changes without restarting Traefik.

**API Daemon Routing Example (`/docker/dynamic/api.yml`):**
```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)"
      service: api
      entryPoints:
        - web

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080" # Routes to the M3TAL API daemon
```

---

## Port Map

| Port | Service               | Access                                                                  |
| :--- | :-------------------- | :---------------------------------------------------------------------- |
| 80   | Traefik HTTP Entry    | Public (if `routing-compose.yml` is active)                             |
| 8080 | M3TAL API Daemon (Go) | Host-local (communicates with dashboard/CLI)                            |
| 8081 | Traefik Dashboard     | Host-local only (typically `127.0.0.1:8081` for Traefik's own dashboard) |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

---

## APT Installation

To install the M3TAL CLI binary and API daemon on Debian/Ubuntu-based systems:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```