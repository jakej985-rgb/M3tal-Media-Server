```markdown
# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for the M3TAL command-line interface (CLI).

## Installation

To install M3TAL, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Core Commands

### Interactive TUI

The `sudo m3tal` command launches the M3TAL interactive Text User Interface (TUI) Control Center, presenting a numbered menu for various operations.

**Usage:**

```bash
sudo m3tal
```

### Initialization

Generates the default `/etc/m3tal/.env` configuration file. This should be run on first install.

**Usage:**

```bash
m3tal init
```

### Health Check

Performs a pre-flight health check of the M3TAL system, verifying Docker connectivity, `.env` file validity, and port availability.

**Usage:**

```bash
m3tal doctor
```

### Configuration Wizard

Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config wizard
```

### Configuration Management

#### Set Environment Variable

Sets a single environment variable in the `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config set API_TOKEN your_super_secret_api_token_here
```

#### Get Environment Variable

Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.

**Usage:**

```bash
m3tal config get DOMAIN
```

#### Scan All Environment Variables

Lists all environment variables across all managed stacks, including their current values.

**Usage:**

```bash
m3tal config scan
```

#### List Current .env File

Displays the entire contents of the current `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config list
```

### Dashboard Management

#### Update Dashboard Password

Updates the password for a dashboard user. If no username or password is provided, it will prompt interactively.

**Usage with username and password:**

```bash
m3tal dashpass admin new_secure_password_123
```

**Usage interactively:**

```bash
m3tal dashpass
```

#### Start Dashboard

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Usage:**

```bash
m3tal dash up
```

#### Stop Dashboard

Stops the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash down
```

#### Restart Dashboard

Restarts the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash restart
```

#### View Dashboard Logs

Streams the logs from the `m3tal-dashboard` container in real-time.

**Usage:**

```bash
m3tal dash logs
```

#### Dashboard Status

Shows the current status of the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash status
```

### Stack Management

#### Bring Up All Stacks

Starts all services defined in `*-compose.yml` files located in `/docker/` using `docker compose`.

**Usage:**

```bash
m3tal up
```

#### Bring Down All Stacks

Stops all services managed by M3TAL across all stacks using `docker compose`.

**Usage:**

```bash
m3metal down
```

#### Stream All Logs

Aggregates and streams logs from all currently running M3TAL-managed Docker containers.

**Usage:**

```bash
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon (`m3tal-api.service`) is managed by systemd.

### Check Service Status

Displays the current status of the `m3tal-api` service.

**Usage:**

```bash
systemctl status m3tal-api
```

### View Service Logs

Streams the journal logs for the `m3tal-api` service in real-time.

**Usage:**

```bash
journalctl -u m3tal-api -f
```

---

## Docker Compose Fallback

In situations where direct CLI control is needed, you can use `docker compose` commands on the compose files located in `/docker/`.

### Bring Up Specific Stack

Starts a specific stack by referencing its compose file.

**Usage (example for a stack named `my-stack`):**

```bash
docker compose -f /docker/my-stack-compose.yml up -d
```

### Bring Down Specific Stack

Stops a specific stack by referencing its compose file.

**Usage (example for a stack named `my-stack`):**

```bash
docker compose -f /docker/my-stack-compose.yml down
```

### View Logs for Specific Stack

Streams logs for services within a specific stack.

**Usage (example for a stack named `my-stack`):**

```bash
docker compose -f /docker/my-stack-compose.yml logs -f
```

---

## M3TAL System Architecture Overview

### Components

*   **CLI binary** (`/usr/bin/m3tal`): Unified Go binary installed via APT. Single entrypoint for all operations.
*   **API daemon** (`m3tal-api.service`): Go binary running as a systemd service on port 8080. Manages Docker, state DB, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): Python/Flask container running internally on port 8082. Communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): Reverse proxy container exposing services by domain name on port 80. Uses file provider for dynamic routing.
*   **Cloudflared** (`routing-compose.yml`): Optional Cloudflare tunnel container for zero-config internet access.

### Filesystem Contract

| Path                       | Purpose                                            |
| :------------------------- | :------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`  | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`        | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                  | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`       | Dashboard credential store. Managed by `m3tal dashpass`. |

### Dashboard Access Modes

The dashboard has two access modes, controlled by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)

*   `DASHBOARD_EXPOSE_MODE=local`
*   Uses override: `m3tal-compose.local.yml`
*   Adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`
*   Access via: `http://HOST_IP:8082` or `http://localhost:8082`
*   No Traefik required. Works out of the box on a home server.
*   **Best for:** LAN-only setups, first-time users, local testing.

#### Mode 2: `traefik`

*   `DASHBOARD_EXPOSE_MODE=traefik`
*   Uses override: `m3tal-compose.traefik.yml`
*   Adds Traefik labels so Traefik routes `dash.${DOMAIN}` → dashboard on port 8082.
*   Access via: `http://dash.DOMAIN` (Traefik must be running via `m3tal up`)
*   **Best for:** Domain-based setups, multiple services behind a reverse proxy.

### Docker / Compose Runtime

M3TAL utilizes **Docker Engine + Docker Compose V2** for container orchestration.

*   The `m3tal up` command runs `docker compose` across all `*-compose.yml` files found in `/docker/`.
*   The `m3tal dash up` command specifically manages the dashboard container:
    1.  Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml` from GitHub.
    2.  Reads `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
    3.  Starts the dashboard with the appropriate compose override file.
*   User stacks reside in `/docker/`. Adding a new stack involves placing a `*-compose.yml` file there.

### Traefik Routing Architecture

Traefik is deployed as a container via `routing-compose.yml`. It acts as a reverse proxy:

*   Binds port 80 on the host as the HTTP entry point.
*   Discovers services automatically via Docker labels.
*   Loads dynamic configuration from `/docker/dynamic/` (file provider with hot-reloading).
*   Routes `api.DOMAIN` → `http://host.docker.internal:8080` (the Go API daemon) via `dynamic/api.yml`.
*   Routes `dash.DOMAIN` → the dashboard container via Traefik labels in `m3tal-compose.traefik.yml` (only when `DASHBOARD_EXPOSE_MODE=traefik`).

### Port Map

| Port | Service                | Access                        |
| :--- | :--------------------- | :---------------------------- |
| 80   | Traefik HTTP entry point | Public (traefik mode)         |
| 8080 | M3TAL API daemon (Go)  | Host-local                    |
| 8081 | Traefik dashboard      | Host-local only               |
| 8082 | M3TAL Dashboard        | Direct port (local mode) or via Traefik (traefik mode) |
```