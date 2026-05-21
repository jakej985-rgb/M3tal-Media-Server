# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for all M3TAL command-line interface (CLI) operations.

## Table of Contents

- [M3TAL System Architecture](#m3tal-system-architecture)
  - [Components](#components)
  - [Filesystem Contract](#filesystem-contract)
  - [Dashboard Access — TWO MODES](#dashboard-access--two-modes)
  - [Docker / Compose Runtime](#docker--compose-runtime)
  - [Deployment Lifecycle — Day 2 Operations](#deployment-lifecycle--day-2-operations)
  - [Traefik Routing Architecture](#traefik-routing-architecture)
  - [Service Management — systemd](#service-management--systemd)
  - [Port Map](#port-map)
- [Core Commands](#core-commands)
  - [`m3tal init`](#m3tal-init)
  - [`m3tal doctor`](#m3tal-doctor)
  - [`m3tal config wizard`](#m3tal-config-wizard)
  - [`m3tal config set KEY VALUE`](#m3tal-config-set-key-value)
  - [`m3tal config get KEY`](#m3tal-config-get-key)
  - [`m3tal config scan`](#m3tal-config-scan)
  - [`m3tal config list`](#m3tal-config-list)
  - [`m3tal up`](#m3tal-up)
  - [`m3tal down`](#m3tal-down)
  - [`m3tal logs`](#m3tal-logs)
- [Dashboard Commands](#dashboard-commands)
  - [`m3tal dashpass [username] [password]`](#m3tal-dashpass-username-password)
  - [`m3tal dash up`](#m3tal-dash-up)
  - [`m3tal dash down`](#m3tal-dash-down)
  - [`m3tal dash restart`](#m3tal-dash-restart)
  - [`m3tal dash logs`](#m3tal-dash-logs)
  - [`m3tal dash status`](#m3tal-dash-status)
- [Interactive TUI](#interactive-tui)
  - [`sudo m3tal`](#sudo-m3tal)
- [Systemd Service Management](#systemd-service-management)
  - [`systemctl status m3tal-api`](#systemctl-status-m3tal-api)
  - [`journalctl -u m3tal-api -f`](#journalctl--u-m3tal-api--f)
- [Direct Docker Compose Commands (Fallback)](#direct-docker-compose-commands-fallback)
- [APT Installation](#apt-installation)

---

## M3TAL System Architecture

### Components

- **CLI binary** (`/usr/bin/m3tal`): Unified Go binary installed via APT. Single entrypoint for all operations.
- **API daemon** (`m3tal-api.service`): Go binary running as a systemd service on port 8080. Manages Docker, state DB, and API routes.
- **Dashboard container** (`m3tal-dashboard`): Python/Flask container running internally on port 8082. Communicates with the API daemon at `http://host.docker.internal:8080`.
- **Traefik gateway** (`routing-compose.yml`): Reverse proxy container exposing services by domain name on port 80. Uses file provider for dynamic routing.
- **Cloudflared** (`routing-compose.yml`): Optional Cloudflare tunnel container for zero-config internet access.

### Filesystem Contract

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

### Dashboard Access — TWO MODES

The dashboard has two access modes, controlled by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`:

**Mode 1: local (default)**
- `DASHBOARD_EXPOSE_MODE=local`
- Uses override: `m3tal-compose.local.yml`
- Adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`
- Access via: `http://HOST_IP:8082` or `http://localhost:8082`
- No Traefik required. Works out of the box on a home server.
- Best for: LAN-only setups, first-time users, local testing.

**Mode 2: traefik**
- `DASHBOARD_EXPOSE_MODE=traefik`
- Uses override: `m3tal-compose.traefik.yml`
- Adds Traefik labels so Traefik routes `dash.${DOMAIN}` → dashboard on port 8082.
- Access via: `http://dash.DOMAIN` (Traefik must be running via `m3tal up`)
- Best for: domain-based setups, multiple services behind a reverse proxy.

### Docker / Compose Runtime

- M3TAL uses **Docker Engine + Docker Compose V2** under the hood. These are hard dependencies.
- The `m3tal up` command runs `docker compose` across all `*-compose.yml` files found in `/docker/`.
- The `m3tal dash up` command specifically manages the dashboard container. It:
  1. Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml` from GitHub.
  2. Reads `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
  3. Starts the dashboard with the appropriate compose override file.
- User stacks live in `/docker/`. Adding a new stack means placing a `*-compose.yml` file there.

### Deployment Lifecycle — Day 2 Operations

Installing a new stack:
1. Place your compose file in `/docker/my-stack-compose.yml`
2. Ensure required variables are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`)
3. Run `m3tal up` to start all stacks.

### Traefik Routing Architecture

Traefik is deployed as a container via `routing-compose.yml`. It:
- Binds port 80 on the host as the HTTP entry point.
- Discovers services automatically via Docker labels.
- Loads dynamic config from `/docker/dynamic/` (file provider, hot-reload).
- Routes `api.DOMAIN` → `http://host.docker.internal:8080` (the Go API daemon) via `dynamic/api.yml`.
- Routes `dash.DOMAIN` → the dashboard container via Traefik labels in `m3tal-compose.traefik.yml` (only when `DASHBOARD_EXPOSE_MODE=traefik`).

### Service Management — systemd

- The API daemon is managed by systemd as `m3tal-api.service`.
- Commands: `systemctl status m3tal-api`, `systemctl restart m3tal-api`, `journalctl -u m3tal-api -f`

### Port Map

| Port | Service | Access |
|------|---------|--------|
| 80 | Traefik HTTP entry point | Public (traefik mode) |
| 8080 | M3TAL API daemon (Go) | Host-local |
| 8081 | Traefik dashboard | Host-local only |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) |

---

## Core Commands

These commands manage the overall M3TAL system, configuration, and running stacks.

### `m3tal init`

Generates `/etc/m3tal/.env` from defaults. Use on first install.

**Usage:**
```bash
sudo m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check. Verifies Docker connectivity, `.env` file validity, and port availability.

**Usage:**
```bash
sudo m3tal doctor
```

### `m3tal config wizard`

Launches an interactive wizard to configure your `/etc/m3tal/.env` file.

**Usage:**
```bash
sudo m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env`.

**Usage:**
```bash
sudo m3tal config set DOMAIN mydomain.com
```

### `m3tal config get KEY`

Reads a single environment variable from `/etc/m3tal/.env`.

**Usage:**
```bash
sudo m3tal config get DOMAIN
```

**Example Output:**
```
mydomain.com
```

### `m3tal config scan`

Lists all environment variables across all known M3TAL stacks, including their current values and defaults.

**Usage:**
```bash
sudo m3tal config scan
```

**Example Output Snippet:**
```json
[
  {
    "key": "DASHBOARD_PORT",
    "default": "8082",
    "value": "8082"
  },
  {
    "key": "DASHBOARD_EXPOSE_MODE",
    "default": "local",
    "value": "local"
  },
  // ... other variables
]
```

### `m3tal config list`

Lists the current contents of the `/etc/m3tal/.env` file.

**Usage:**
```bash
sudo m3tal config list
```

**Example Output:**
```
# M3TAL Configuration
DOMAIN=localhost
DASHBOARD_PORT=8082
# ... other variables
```

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in `/docker/`. This starts all your configured stacks.

**Usage:**
```bash
sudo m3tal up
```

### `m3tal down`

Runs `docker compose down` across all configured stacks. This stops and removes containers, networks, and volumes for all services.

**Usage:**
```bash
sudo m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL stacks.

**Usage:**
```bash
sudo m3tal logs
```

---

## Dashboard Commands

These commands specifically manage the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the dashboard user password. If `username` and `password` are omitted, it will prompt interactively.

**Usage (interactive):**
```bash
sudo m3tal dashpass
```

**Usage (with arguments):**
```bash
sudo m3tal dashpass myuser mysecurepassword123
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container. It respects the `DASHBOARD_EXPOSE_MODE` setting in your `.env` file.

**Usage:**
```bash
sudo m3tal dash up
```

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash down
```

### `m3tal dash restart`

Restarts the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash restart
```

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash status
```

**Example Output:**
```
m3tal-dashboard Up (running)
```

---

## Interactive TUI

### `sudo m3tal`

Opens the interactive TUI Control Center. This presents a numbered menu of common operations, offering a guided experience.

**Usage:**
```bash
sudo m3tal
```

**Example Menu:**
```
Welcome to the M3TAL TUI Control Center!

1. System Health Check (doctor)
2. Configure M3TAL (.env wizard)
3. Start all stacks (up)
4. Stop all stacks (down)
5. View aggregated logs (logs)
6. Manage Dashboard
...
```

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service.

### `systemctl status m3tal-api`

Displays the current status of the `m3tal-api.service`.

**Usage:**
```bash
sudo systemctl status m3tal-api
```

**Example Output:**
```
● m3tal-api.service - M3TAL API Daemon
     Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
     Active: active (running) since Mon 2023-10-27 10:00:00 UTC; 1 day ago
   Main PID: 1234 (m3tal-api)
      Tasks: 10 (limit: 4915)
     Memory: 25.0M
        CPU: 500ms
     CGroup: /system.slice/m3tal-api.service
             └─1234 /usr/bin/m3tal-api

Oct 28 10:00:00 your-server systemd[1]: Started M3TAL API Daemon.
```

### `journalctl -u m3tal-api -f`

Streams the logs for the `m3tal-api.service` in real-time. Use `-f` to follow the logs.

**Usage:**
```bash
sudo journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Commands (Fallback)

In situations where M3TAL CLI commands might not suffice, you can interact with Docker Compose directly. M3TAL orchestrates these commands, but understanding them provides deeper insight.

- **Starting all stacks:**
  ```bash
  cd /docker && sudo docker compose up -d
  ```

- **Stopping all stacks:**
  ```bash
  cd /docker && sudo docker compose down
  ```

- **Viewing logs for a specific stack (e.g., named `my-stack`):**
  ```bash
  cd /docker && sudo docker compose -f my-stack-compose.yml logs -f
  ```

- **Building images for a specific stack:**
  ```bash
  cd /docker && sudo docker compose -f my-stack-compose.yml build
  ```

---

## APT Installation

To install or update M3TAL, use the following APT commands.

**Usage:**
```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update package lists and install
sudo apt update && sudo apt install -y m3tal
```