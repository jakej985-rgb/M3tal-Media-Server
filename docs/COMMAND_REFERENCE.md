```markdown
# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for all available `m3tal` CLI commands.

## Getting Started

### APT Installation

To install or update M3TAL, use the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### Initializing M3TAL

On first install, generate the default configuration file:

```bash
m3tal init
```

This command creates `/etc/m3tal/.env` with default settings.

## Core Commands

### Interactive Control Center

Launch the full-featured TUI Control Center:

```bash
sudo m3tal
```

This opens a numbered menu allowing you to manage M3TAL services, configuration, and more.

### System Health Check

Perform a pre-flight health check of your M3TAL environment:

```bash
m3tal doctor
```

This checks Docker connectivity, `.env` file validity, and port availability.

## Configuration Management

### Configuration Wizard

An interactive wizard to guide you through configuring `/etc/m3tal/.env`:

```bash
m3tal config wizard
```

### Setting Environment Variables

Set a single environment variable in `/etc/m3tal/.env`:

```bash
m3tal config set DASHBOARD_PORT 8083
```

### Getting Environment Variables

Read a single environment variable from `/etc/m3tal/.env`:

```bash
m3tal config get DOMAIN
```

### Scanning All Environment Variables

List all environment variables across all configured stacks. This command aggregates values from various sources including defaults and potentially overrides:

```bash
m3tal config scan
```

### Listing Current .env Contents

Display the current contents of your `/etc/m3tal/.env` file:

```bash
m3tal config list
```

## Dashboard Management

The M3TAL Dashboard provides a web interface for managing your M3TAL instance.

### Updating Dashboard Password

Update the password for a dashboard user. If username and password are not provided, it will prompt interactively:

```bash
m3tal dashpass admin new_secure_password
```

Or interactively:

```bash
m3tal dashpass
```

### Starting the Dashboard

Pull the latest dashboard compose configuration from GitHub and start the dashboard container:

```bash
m3tal dash up
```

### Stopping the Dashboard

Stop the dashboard container:

```bash
m3tal dash down
```

### Restarting the Dashboard

Restart the dashboard container:

```bash
m3tal dash restart
```

### Viewing Dashboard Logs

Stream logs from the dashboard container:

```bash
m3tal dash logs
```

### Dashboard Container Status

Show the current status of the dashboard container:

```bash
m3tal dash status
```

## Stack Management

### Bringing Up All Stacks

Run `docker compose up` across all `*-compose.yml` files located in `/docker/`:

```bash
m3tal up
```

### Tearing Down All Stacks

Run `docker compose down` across all configured stacks:

```bash
m3tal down
```

### Aggregated Logs

Stream aggregated logs from all currently running M3TAL stacks:

```bash
m3tal logs
```

## Systemd Service Management

The M3TAL API daemon runs as a systemd service.

### Checking Service Status

View the status of the `m3tal-api` service:

```bash
systemctl status m3tal-api
```

### Following Service Logs

Stream logs from the `m3tal-api` service using `journalctl`:

```bash
journalctl -u m3tal-api -f
```

## Direct Docker Compose Commands (Fallback)

While `m3tal` abstracts Docker Compose, you can still use `docker compose` directly for advanced scenarios or debugging.

### Running Specific Compose Files

To run a specific compose file, navigate to its directory (e.g., `/docker/`) and use `docker compose`:

```bash
cd /docker/
docker compose -f my-stack-compose.yml up -d
```

### Stopping Specific Compose Files

```bash
cd /docker/
docker compose -f my-stack-compose.yml down
```

### Listing Docker Services

View all running Docker containers, including M3TAL services:

```bash
docker ps
```

### Viewing Docker Compose Project Names

Identify compose projects managed by M3TAL:

```bash
docker compose ps --all
```

## M3TAL System Architecture Details

### Components

*   **CLI binary** (`/usr/bin/m3tal`): The unified Go binary installed via APT, acting as the single entry point for all operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on port 8080, responsible for managing Docker, the state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port 8082, communicating with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container exposing services by domain name on port 80, utilizing a file provider for dynamic routing.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for zero-config internet access.

### Filesystem Contract

| Path                     | Purpose                                        |
| :----------------------- | :--------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`. |

### Dashboard Access Modes

The dashboard can be accessed in two modes, determined by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`:

*   **`local` (default):**
    *   Uses `m3tal-compose.local.yml` for overrides.
    *   Binds `${DASHBOARD_PORT:-8082}:8082` directly.
    *   Access via `http://HOST_IP:8082` or `http://localhost:8082`.
    *   No Traefik required; suitable for LAN-only or initial setups.

*   **`traefik`:**
    *   Uses `m3tal-compose.traefik.yml` for overrides.
    *   Configures Traefik to route `dash.${DOMAIN}` to the dashboard on port 8082.
    *   Access via `http://dash.DOMAIN` (requires Traefik running via `m3tal up`).
    *   Ideal for domain-based setups behind a reverse proxy.

### Docker / Compose Runtime

M3TAL relies on **Docker Engine** and **Docker Compose V2**.

*   `m3tal up` executes `docker compose` on all `*-compose.yml` files in `/docker/`.
*   `m3tal dash up` specifically manages the dashboard container, downloading its configuration and starting it with the appropriate override based on `DASHBOARD_EXPOSE_MODE`.

### Port Map

| Port | Service                | Access Method                           |
| :--- | :--------------------- | :-------------------------------------- |
| 80   | Traefik HTTP Entrypoint | Public (when `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API daemon (Go)  | Host-local                              |
| 8081 | Traefik dashboard      | Host-local only                         |
| 8082 | M3TAL Dashboard        | Direct port (`local` mode) or via Traefik (`traefik` mode) |
```