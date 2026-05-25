# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for the M3TAL CLI, covering all available commands and their usage.

## Table of Contents

- [M3TAL System Architecture](#m3tal-system-architecture)
- [Filesystem Contract](#filesystem-contract)
- [Dashboard Access Modes](#dashboard-access-modes)
- [Docker / Compose Runtime](#docker--compose-runtime)
- [Deployment Lifecycle](#deployment-lifecycle)
- [Traefik Routing Architecture](#traefik-routing-architecture)
- [Service Management — systemd](#service-management--systemd)
- [Port Map](#port-map)
- [APT Installation](#apt-installation)
- [Command Reference](#command-reference)
  - [`sudo m3tal`](#sudo-m3tal)
  - [`m3tal init`](#m3tal-init)
  - [`m3tal doctor`](#m3tal-doctor)
  - [`m3tal config wizard`](#m3tal-config-wizard)
  - [`m3tal config set`](#m3tal-config-set)
  - [`m3tal config get`](#m3tal-config-get)
  - [`m3tal config scan`](#m3tal-config-scan)
  - [`m3tal config list`](#m3tal-config-list)
  - [`m3tal dashpass`](#m3tal-dashpass)
  - [`m3tal dash up`](#m3tal-dash-up)
  - [`m3tal dash down`](#m3tal-dash-down)
  - [`m3tal dash restart`](#m3tal-dash-restart)
  - [`m3tal dash logs`](#m3tal-dash-logs)
  - [`m3tal dash status`](#m3tal-dash-status)
  - [`m3tal up`](#m3tal-up)
  - [`m3tal down`](#m3tal-down)
  - [`m3tal logs`](#m3tal-logs)
- [Direct Docker Compose Commands](#direct-docker-compose-commands)

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

### Dashboard Access Modes

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

## APT Installation

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## Command Reference

### `sudo m3tal`

Opens the interactive TUI Control Center with a numbered menu.

**Usage:**
```bash
sudo m3tal
```

**Example:**
```bash
sudo m3tal
```
*(This command will launch an interactive TUI.)*

---

### `m3tal init`

Generates `/etc/m3tal/.env` from defaults. Use on first install.

**Usage:**
```bash
m3tal init
```

**Example:**
```bash
m3tal init
```
*(This command creates or overwrites `/etc/m3tal/.env` with default configuration values.)*

---

### `m3tal doctor`

Performs a pre-flight health check: Docker connectivity, `.env` validity, and port availability.

**Usage:**
```bash
m3tal doctor
```

**Example:**
```bash
m3tal doctor
```
*(This command will output a report on the system's health status.)*

---

### `m3tal config wizard`

Starts an interactive wizard to configure `/etc/m3tal/.env`.

**Usage:**
```bash
m3tal config wizard
```

**Example:**
```bash
m3tal config wizard
```
*(This command will guide you through setting up your M3TAL configuration interactively.)*

---

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env`.

**Usage:**
```bash
m3tal config set <KEY> <VALUE>
```

**Example:**
```bash
m3tal config set DASHBOARD_PORT 8083
```
*(This command sets the `DASHBOARD_PORT` variable to `8083` in your `.env` file.)*

---

### `m3tal config get KEY`

Reads a single environment variable from `/etc/m3tal/.env`.

**Usage:**
```bash
m3tal config get <KEY>
```

**Example:**
```bash
m3tal config get DOMAIN
```
*(This command will output the current value of the `DOMAIN` variable.)*

---

### `m3tal config scan`

Lists all environment variables across all stacks managed by M3TAL, including their current values.

**Usage:**
```bash
m3tal config scan
```

**Example:**
```bash
m3tal config scan
```
*(This command will display a comprehensive list of all environment variables and their settings.)*

---

### `m3tal config list`

Lists the current contents of the `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

**Example:**
```bash
m3tal config list
```
*(This command will print the entire content of your `.env` file to the console.)*

---

### `m3tal dashpass [username] [password]`

Updates the dashboard user password. If arguments are omitted, it will prompt interactively.

**Usage:**
```bash
m3tal dashpass [username] [password]
```

**Example (interactive):**
```bash
sudo m3tal dashpass
```
*(This command will prompt you for the username and new password.)*

**Example (with arguments):**
```bash
sudo m3tal dashpass myuser newsecurepassword123
```
*(This command sets the password for `myuser` to `newsecurepassword123`.)*

---

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Usage:**
```bash
m3tal dash up
```

**Example:**
```bash
sudo m3tal dash up
```
*(This command ensures you have the latest dashboard configuration and then starts the `m3tal-dashboard` container.)*

---

### `m3tal dash down`

Stops the dashboard container.

**Usage:**
```bash
m3tal dash down
```

**Example:**
```bash
sudo m3tal dash down
```
*(This command stops and removes the `m3tal-dashboard` container.)*

---

### `m3tal dash restart`

Restarts the dashboard container.

**Usage:**
```bash
m3tal dash restart
```

**Example:**
```bash
sudo m3tal dash restart
```
*(This command stops and then starts the `m3tal-dashboard` container.)*

---

### `m3tal dash logs`

Streams the logs from the dashboard container.

**Usage:**
```bash
m3tal dash logs
```

**Example:**
```bash
sudo m3tal dash logs
```
*(This command will display the real-time logs from the `m3tal-dashboard` container.)*

---

### `m3tal dash status`

Shows the current status of the dashboard container.

**Usage:**
```bash
m3tal dash status
```

**Example:**
```bash
sudo m3tal dash status
```
*(This command will output information about the `m3tal-dashboard` container's state.)*

---

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in `/docker/`. This starts all your deployed stacks.

**Usage:**
```bash
m3tal up
```

**Example:**
```bash
sudo m3tal up
```
*(This command will start all services defined in the compose files within the `/docker` directory.)*

---

### `m3tal down`

Runs `docker compose down` across all stacks managed by M3TAL.

**Usage:**
```bash
m3tal down
```

**Example:**
```bash
sudo m3tal down
```
*(This command will stop and remove all containers, networks, and volumes for all managed stacks.)*

---

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks.

**Usage:**
```bash
m3tal logs
```

**Example:**
```bash
sudo m3tal logs
```
*(This command will display a unified stream of logs from all active Docker containers managed by M3TAL.)*

---

## Direct Docker Compose Commands

As a fallback and for finer control, you can also directly use `docker compose` commands on the individual stack files located in `/docker/`.

**Example:** To start only the Traefik gateway:
```bash
sudo docker compose -f /docker/routing-compose.yml up -d
```

**Example:** To stop all services defined in `my-app-compose.yml`:
```bash
sudo docker compose -f /docker/my-app-compose.yml down
```

**Example:** To view logs for a specific service within a stack (e.g., `web` service in `my-app-compose.yml`):
```bash
sudo docker compose -f /docker/my-app-compose.yml logs web
```