# M3TAL CLI Command Reference

Welcome, M3TAL Operative. This document is your comprehensive cheat-sheet for interacting with the M3TAL Ecosystem via its command-line interface. The `m3tal` CLI binary, installed via APT to `/usr/bin/m3tal`, is your unified entry point for managing configuration, services, and containers.

## M3TAL Ecosystem Overview

M3TAL is a robust, self-hosted platform built upon Docker Engine and Docker Compose V2. It orchestrates a core API daemon, a web dashboard, and user-defined Docker stacks, all managed via a single CLI.

### Filesystem Contract

Understanding the M3TAL filesystem is critical for effective management.

| Path                         | Purpose                                                                                                                                                                                                                                                                                                                                                          |
| :--------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`            | **Primary Configuration File.** This file stores all critical environment variables that configure the M3TAL API, Dashboard, and Docker stacks. It is managed by `m3tal config wizard` and `m3tal config set`.                                                                                                                                                      |
| `/var/lib/m3tal/state.db`    | **SQLite State Database.** Automatically created and managed by the `m3tal-api` daemon. It stores internal state, such as service status and network configurations. Do not modify directly.                                                                                                                                                                        |
| `/opt/m3tal/stack/`          | **Canonical Stack Directory.** This is the internal directory where M3TAL stores its core Docker Compose files (`m3tal-compose.yml`, `routing-compose.yml`) and other stack-related assets.                                                                                                                                                                     |
| `/docker`                    | **User-Facing Stack Directory (Symlink).** This path is a symlink to `/opt/m3tal/stack/`. It is the user-friendly location where you place your custom `*-compose.yml` files for M3TAL to manage. All `m3tal up`/`down` operations target this directory.                                                                                                        |
| `/docker/users.json`         | **Dashboard Credential Store.** This JSON file holds the hashed credentials for accessing the M3TAL Dashboard. It is managed exclusively by the `m3tal dashpass` command.                                                                                                                                                                                          |
| `/docker/dynamic/`           | **Traefik Dynamic Configuration.** This directory (within the `/docker` symlink) is used by Traefik's file provider for dynamic routing rules (e.g., `api.yml`). Traefik hot-reloads changes in this directory.                                                                                                                                                      |
| `/docker/m3tal-compose.yml`  | Base Docker Compose file for the M3TAL Dashboard. This file is updated by `m3tal dash up`.                                                                                                                                                                                                                                                                        |
| `/docker/routing-compose.yml`| Base Docker Compose file for Traefik and Cloudflared. This file is included in `m3tal up` operations.                                                                                                                                                                                                                                                             |

---

## Core M3TAL Commands

These commands provide fundamental system interaction, initialization, and health checks.

### `sudo m3tal`

Opens the interactive Text-User Interface (TUI) Control Center. This allows you to manage M3TAL through a guided, numbered menu system, ideal for initial setup and common operations. Requires `sudo` as it interacts with system-level resources and Docker.

```bash
sudo m3tal
```

### `m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from M3TAL's default values. This command should be run on the first installation or if the `.env` file is missing. It will preserve existing values if the file already exists, only adding new ones.

```bash
m3tal init
# Output: Generating /etc/m3tal/.env from defaults...
# Output: /etc/m3tal/.env created/updated successfully.
```

### `m3tal doctor`

Performs a pre-flight health check of the M3TAL ecosystem. It verifies critical components such as Docker connectivity, the validity and existence of `/etc/m3tal/.env`, and ensures that essential ports (e.g., 8080 for API, 8082 for Dashboard) are available.

```bash
m3tal doctor
# Output: Checking M3TAL ecosystem health...
# Output: Docker Daemon: Connected
# Output: /etc/m3tal/.env: Valid and found
# Output: Port 8080 (M3TAL API): Available
# Output: Port 8082 (M3TAL Dashboard): Available
# Output: All checks passed. M3TAL is ready!
```

---

## Configuration Management (`m3tal config`)

These commands manage the `/etc/m3tal/.env` configuration file.

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for making broad configuration changes, ensuring all dependencies and relationships between variables are handled correctly.

```bash
m3tal config wizard
# Output: Starting M3TAL configuration wizard...
# Output: Current value for DOMAIN (localhost): [Enter new value or press Enter to keep] mydomain.com
# Output: Current value for DASHBOARD_EXPOSE_MODE (local): [Enter 'local' or 'traefik'] traefik
# ... (wizard continues)
# Output: Configuration saved to /etc/m3tal/.env. Restart API and Dashboard for changes to take effect.
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env` to the specified `VALUE`. Use this for quick, targeted adjustments.

```bash
m3tal config set DOMAIN "m3tal.local"
# Output: DOMAIN set to "m3tal.local" in /etc/m3tal/.env.

m3tal config set DASHBOARD_SECRET "mySuperSecretKey123"
# Output: DASHBOARD_SECRET set to "mySuperSecretKey123" in /etc/m3tal/.env.
```

### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

```bash
m3tal config get DOMAIN
# Output: m3tal.local

m3tal config get DASHBOARD_PORT
# Output: 8082
```

### `m3tal config scan`

Lists all known environment variables and their default values across all M3TAL-managed stacks. This provides a comprehensive overview of what can be configured.

```bash
m3tal config scan
# Output: Scanning all known M3TAL environment variables:
# Output: DASHBOARD_PORT (Default: 8082)
# Output: DASHBOARD_EXPOSE_MODE (Default: local)
# Output: DOMAIN (Default: localhost)
# Output: API_TOKEN (Default: change_me_api_token)
# ... (lists all known variables and their defaults)
```

### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file. This is useful for reviewing your active configuration.

```bash
m3tal config list
# Output: # M3TAL Configuration File
# Output: DOMAIN="m3tal.local"
# Output: DASHBOARD_PORT="8082"
# Output: DASHBOARD_EXPOSE_MODE="traefik"
# Output: ... (full contents of your .env file)
```

---

## M3TAL Dashboard Management (`m3tal dash`)

The `m3tal dash` commands specifically manage the M3TAL web dashboard container.

### Dashboard Access Modes (Critical)

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **`DASHBOARD_EXPOSE_MODE=local` (Default)**
    *   **Configuration:** Uses the `m3tal-compose.local.yml` override file.
    *   **Mechanism:** Adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082` to the dashboard container.
    *   **Access:** Access directly via `http://HOST_IP:8082` or `http://localhost:8082`.
    *   **Requirements:** No Traefik required. Works out of the box on a home server.
    *   **Best for:** LAN-only setups, first-time users, local testing.

2.  **`DASHBOARD_EXPOSE_MODE=traefik`**
    *   **Configuration:** Uses the `m3tal-compose.traefik.yml` override file.
    *   **Mechanism:** Adds Traefik labels to the dashboard container, instructing Traefik to route `dash.${DOMAIN}` to the dashboard on its internal port 8082. Traefik must be running (`m3tal up` will start it).
    *   **Access:** Access via `http://dash.YOUR_DOMAIN` (e.g., `http://dash.m3tal.local` if `DOMAIN` is `m3tal.local`).
    *   **Requirements:** Traefik must be deployed and running via `routing-compose.yml`.
    *   **Best for:** Domain-based setups, integration with other services behind a reverse proxy.

The `m3tal dash up` command intelligently applies the correct compose override based on this `.env` variable.

### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, the command becomes interactive, prompting you for the necessary details. This updates the `/docker/users.json` file.

```bash
# Interactive mode (recommended for security)
m3tal dashpass
# Output: Enter username: admin
# Output: Enter new password: ***********
# Output: Confirm new password: ***********
# Output: User 'admin' password updated successfully in /docker/users.json.

# Non-interactive mode (use with caution, avoid in history)
m3tal dashpass admin MySecurePassword123
# Output: User 'admin' password updated successfully in /docker/users.json.
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) directly from GitHub, then starts or updates the `m3tal-dashboard` container. It automatically selects the correct override based on `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.

```bash
m3tal dash up
# Output: Pulling latest dashboard compose configurations from GitHub...
# Output: Downloading m3tal-compose.yml... Done.
# Output: Downloading m3tal-compose.local.yml... Done.
# Output: Downloading m3tal-compose.traefik.yml... Done.
# Output: DASHBOARD_EXPOSE_MODE is 'local'. Starting dashboard with local port binding.
# Output: [+] Running 1/1
# Output:  ✔ Container m3tal-dashboard Started
```

### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

```bash
m3tal dash down
# Output: Stopping m3tal-dashboard container...
# Output: [+] Stopping 1/1
# Output:  ✔ Container m3tal-dashboard  Stopped
# Output: Container m3tal-dashboard removed.
```

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container. This is useful after making configuration changes that affect the dashboard.

```bash
m3tal dash restart
# Output: Restarting m3tal-dashboard container...
# Output: [+] Restarting 1/1
# Output:  ✔ Container m3tal-dashboard  Restarted
```

### `m3tal dash logs`

Streams the real-time logs from the `m3tal-dashboard` container. Press `Ctrl+C` to exit the log stream.

```bash
m3tal dash logs
# Output: Attaching to m3tal-dashboard
# Output: m3tal-dashboard  |  * Serving Flask app 'server'
# Output: m3tal-dashboard  |  * Debug mode: off
# Output: ... (dashboard logs stream here)
```

### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` container (e.g., `running`, `exited`, `restarting`).

```bash
m3tal dash status
# Output: Container 'm3tal-dashboard': running (Up 5 minutes)
```

---

## M3TAL Stack Management (`m3tal`)

These commands manage all Docker Compose stacks defined in the `/docker/` directory.

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This starts or updates all defined services (including Traefik, if `routing-compose.yml` is present, and any user-defined stacks).

```bash
m3tal up
# Output: Running docker compose up -d for all stacks in /docker/...
# Output: Project "m3tal", build service
# Output: [+] Running 2/2
# Output:  ✔ Container m3tal-dashboard  Running
# Output:  ✔ Container traefik          Running
# Output:  ✔ Container cloudflared      Running
# Output:  ✔ Container my-app           Running
```

### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files found in the `/docker/` directory. This stops and removes all containers, networks, and volumes defined by M3TAL-managed stacks.

```bash
m3tal down
# Output: Running docker compose down for all stacks in /docker/...
# Output: [+] Stopping 4/4
# Output:  ✔ Container m3tal-dashboard  Stopped
# Output:  ✔ Container traefik          Stopped
# Output:  ✔ Container cloudflared      Stopped
# Output:  ✔ Container my-app           Stopped
# Output: Removing 4 containers, 3 networks, 2 volumes...
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL-managed Docker containers. This provides a holistic view of the system's activity. Press `Ctrl+C` to exit the log stream.

```bash
m3tal logs
# Output: Attaching to m3tal-dashboard, traefik, cloudflared, my-app
# Output: m3tal-dashboard  | [2023-10-27 10:30:01,123] INFO: Dashboard started.
# Output: traefik          | time="2023-10-27T10:30:05Z" level=info msg="Configuration loaded from file: /etc/traefik/traefik.yml"
# Output: my-app           | My custom app is doing things...
# ... (aggregated logs stream here)
```

---

## Systemd Service Management

The M3TAL API daemon (`m3tal-api`) runs as a systemd service, managing Docker interactions and providing API routes. These standard systemd commands are used to interact with it.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api.service`. This shows if it's active, running, and recent log entries.

```bash
systemctl status m3tal-api
# Output: ● m3tal-api.service - M3TAL API Daemon
# Output:    Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
# Output:    Active: active (running) since Fri 2023-10-27 10:00:00 UTC; 1h ago
# Output:  Main PID: 12345 (m3tal-api)
# Output:     Tasks: 7 (limit: 4915)
# Output:    Memory: 15.6M
# Output:       CPU: 125ms
# Output:    CGroup: /system.slice/m3tal-api.service
# Output:            └─12345 /usr/bin/m3tal-api --config /etc/m3tal/.env
# Output: Oct 27 10:00:00 m3tal-host systemd[1]: Started M3TAL API Daemon.
# Output: Oct 27 10:00:01 m3tal-host m3tal-api[12345]: INFO: M3TAL API started on :8080
```

### `systemctl restart m3tal-api`

Restarts the `m3tal-api.service`. This is often required after making changes to `/etc/m3tal/.env` that affect the API daemon (e.g., `API_TOKEN`).

```bash
sudo systemctl restart m3tal-api
# Output: (No direct output, but the service will restart)
```

### `journalctl -u m3tal-api -f`

Streams the real-time logs from the `m3tal-api.service`. This is invaluable for debugging issues with the M3TAL core API. Press `Ctrl+C` to exit the log stream.

```bash
sudo journalctl -u m3tal-api -f
# Output: -- Journal begins at Thu 2023-10-26 08:00:00 UTC, ends at Fri 2023-10-27 11:30:00 UTC. --
# Output: Oct 27 11:29:58 m3tal-host m3tal-api[12345]: INFO: Handling API request /status
# Output: Oct 27 11:30:05 m3tal-host m3tal-api[12345]: WARN: Invalid API token provided.
# ... (API daemon logs stream here)
```

---

## Direct Docker Commands (Fallback)

While `m3tal` CLI abstracts Docker Compose operations, it's built on top of Docker Engine + Docker Compose V2. Knowing the direct Docker commands can be useful for advanced debugging or when the `m3tal` CLI itself is having issues.

**Note:** Always prefer the `m3tal` CLI commands for consistency and safety, as they handle file paths, environment variables, and specific overrides correctly.

### Bringing up all M3TAL-managed stacks

The `m3tal up` command internally constructs a `docker compose` command referencing all `*-compose.yml` files in `/docker/`. A common manual equivalent for core services would be:

```bash
# This command needs to include all relevant compose files for your setup.
# The M3TAL API daemon reads /etc/m3tal/.env, but for direct Docker Compose,
# you'd typically either source it or pass variables via -e or --env-file.
# This example uses `env-file` for clarity.

# Ensure /etc/m3tal/.env exists and contains necessary variables (like DOMAIN, DASHBOARD_EXPOSE_MODE, etc.)
docker compose -f /docker/m3tal-compose.yml \
               -f /docker/routing-compose.yml \
               -f /docker/m3tal-compose.$(m3tal config get DASHBOARD_EXPOSE_MODE).yml \
               --env-file /etc/m3tal/.env \
               up -d
```

### Bringing down all M3TAL-managed stacks

Similar to `m3tal up`, `m3tal down` orchestrates the shutdown of all stacks.

```bash
# Again, ensure all compose files are referenced as m3tal does.
docker compose -f /docker/m3tal-compose.yml \
               -f /docker/routing-compose.yml \
               -f /docker/m3tal-compose.$(m3tal config get DASHBOARD_EXPOSE_MODE).yml \
               --env-file /etc/m3tal/.env \
               down
```

### Streaming logs from a specific container

You can stream logs from any individual container directly:

```bash
docker logs -f m3tal-dashboard
docker logs -f traefik
docker logs -f cloudflared
```

### Checking container status

```bash
docker ps -a
docker ps -a --filter "name=m3tal-dashboard"
```

---

## Key Architectural Concepts

*   **Docker Engine + Docker Compose V2:** M3TAL relies heavily on these underlying Docker technologies for container orchestration.
*   **`/docker` Directory:** This symbolic link points to `/opt/m3tal/stack/` and is the canonical location for all Docker Compose files (`*-compose.yml`) managed by `m3tal up`/`down`. Place your custom stack definitions here.
*   **Traefik Routing:** Traefik (`routing-compose.yml`) acts as the reverse proxy.
    *   It binds to host port 80 (and 443 if HTTPS is configured).
    *   It auto-discovers services via Docker labels (e.g., for the dashboard in `traefik` mode).
    *   It also loads dynamic configuration from `/docker/dynamic/` (e.g., `dynamic/api.yml` routes `api.${DOMAIN}` to `http://host.docker.internal:8080` for the M3TAL API).
*   **API Daemon (`m3tal-api.service`):** The core Go API daemon runs directly on the host (not in Docker) on port 8080. It manages the `state.db` and exposes endpoints for the dashboard and CLI.
*   **Dashboard Container (`m3tal-dashboard`):** A Python/Flask application that communicates with the API daemon at `http://host.docker.internal:8080`. Its exposure method is dynamically configured by `DASHBOARD_EXPOSE_MODE`.

## Port Map

| Port | Service               | Access                                   |
| :--- | :-------------------- | :--------------------------------------- |
| 80   | Traefik HTTP entry    | Public (if Traefik is running)           |
| 8080 | M3TAL API daemon (Go) | Host-local only                          |
| 8081 | Traefik dashboard     | Host-local only (`127.0.0.1:8081`)       |
| 8082 | M3TAL Dashboard       | Direct port (local mode) / via Traefik (traefik mode) |

---

## APT Installation (Always include this exact block)

To install the `m3tal` CLI and API daemon on Debian/Ubuntu-based systems:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```