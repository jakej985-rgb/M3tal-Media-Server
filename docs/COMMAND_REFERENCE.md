As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the definitive M3TAL Command Line Interface (CLI) cheat-sheet. This document details every essential command, offering crystal-clear usage examples and critical context for managing your M3TAL server.

---

# M3TAL CLI Command Reference

The M3TAL CLI (`/usr/bin/m3tal`) is your single entry point for interacting with the M3TAL ecosystem. It provides a unified interface to manage configuration, deploy services, monitor health, and maintain your server effortlessly.

## M3TAL Ecosystem Overview

M3TAL orchestrates your self-hosted applications using Docker Engine and Docker Compose V2. It consists of:

*   **M3TAL CLI (`/usr/bin/m3tal`):** The primary user interface.
*   **M3TAL API (`m3tal-api.service`):** A Go daemon running on port `8080`, managed by systemd. It handles Docker interactions, state management (in `/var/lib/m3tal/state.db`), and exposes an API.
*   **M3TAL Dashboard (`m3tal-dashboard`):** A Python/Flask Docker container (internal port `8082`) that communicates with the M3TAL API.
*   **Traefik Gateway (`routing-compose.yml`):** A Docker container acting as a reverse proxy, exposing services on port `80` and dynamically routing traffic based on domain names and Docker labels.
*   **Cloudflared (`routing-compose.yml`):** An optional Cloudflare tunnel container for secure, zero-config internet access to your services.

### Core Filesystem Contract

Understanding the filesystem layout is crucial for M3TAL management:

| Path                        | Purpose                                                                                                                                     |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | **Primary Configuration File.** Stores all M3TAL environment variables. Managed by `m3tal config wizard` and `m3tal config set`.              |
| `/var/lib/m3tal/state.db`   | **SQLite State Database.** Automatically created and managed by the `m3tal-api.service` daemon to store internal state.                      |
| `/opt/m3tal/stack/`         | **Canonical Stack Directory.** Contains core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`) and Traefik config.     |
| `/docker`                   | **User-Facing Stack Symlink.** A symlink to `/opt/m3tal/stack/`. This is where you place your custom `*-compose.yml` files for new services. |
| `/docker/users.json`        | **Dashboard Credential Store.** Stores encrypted username/password pairs for the M3TAL Dashboard. Managed by `m3tal dashpass`.               |
| `/docker/dynamic/`          | **Traefik Dynamic Configuration.** Contains `.yml` files for custom Traefik rules (e.g., exposing the M3TAL API). Traefik hot-reloads changes. |

### M3TAL Dashboard Access Modes

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

#### Mode 1: `local` (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   Uses the `m3tal-compose.local.yml` override file.
*   Adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access via:** `http://YOUR_HOST_IP:8082` or `http://localhost:8082`.
*   **Usage:** Ideal for LAN-only setups, initial configuration, and local testing. No Traefik configuration required.

#### Mode 2: `traefik`

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   Uses the `m3tal-compose.traefik.yml` override file.
*   Adds Traefik labels to the dashboard container, allowing Traefik to route `dash.${DOMAIN}` to the dashboard on port `8082`.
*   **Access via:** `http://dash.YOUR_DOMAIN` (e.g., `http://dash.mydomain.com`).
*   **Prerequisite:** Traefik must be running via `m3tal up` (specifically `routing-compose.yml`).
*   **Usage:** Best for domain-based access, exposing services behind a reverse proxy, and integrating with other Traefik-managed applications.

### Docker / Compose Runtime Explained

M3TAL leverages **Docker Engine** and **Docker Compose V2** for all container orchestration.

*   The `m3tal up` command iterates through all `*-compose.yml` files located in `/docker/` and executes `docker compose up -d` for each.
*   The `m3tal dash up` command specifically manages the dashboard:
    1.  It pulls the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub.
    2.  It then starts the dashboard container, applying the correct override based on `DASHBOARD_EXPOSE_MODE`.
*   **To add new applications:** Simply place your `my-app-compose.yml` file into `/docker/` and run `m3tal up`.

### Traefik Routing Architecture

Traefik, deployed via `routing-compose.yml`, serves as the central entry point for all web traffic:

*   It binds host port `80` (HTTP) and optionally `443` (HTTPS) as its entry points.
*   It automatically discovers services by inspecting Docker container labels (e.g., `traefik.enable=true`).
*   It loads additional dynamic configuration from `/docker/dynamic/` (using a file provider), allowing for hot-reloadable custom routing rules.
*   The M3TAL API is exposed via Traefik by a dynamic rule in `/docker/dynamic/api.yml`, routing `api.${DOMAIN}` to `http://host.docker.internal:8080`.

### M3TAL Port Map

| Port | Service                               | Access                                                          |
| :--- | :------------------------------------ | :-------------------------------------------------------------- |
| 80   | Traefik HTTP Entry Point              | Public (if `routing-compose.yml` is active)                     |
| 443  | Traefik HTTPS Entry Point             | Public (if `routing-compose.yml` is active and configured for SSL) |
| 8080 | M3TAL API Daemon (Go)                 | Host-local only (accessed by containers via `host.docker.internal`) |
| 8081 | Traefik Dashboard (Internal)          | Host-local only (for Traefik's own UI, bound via `127.0.0.1:8081:8080`) |
| 8082 | M3TAL Dashboard Container (Internal) | Direct port (local mode) or via Traefik (traefik mode)        |

---

## M3TAL CLI Command Reference

### Interactive TUI Control Center

#### `sudo m3tal`

Opens the M3TAL interactive Terminal User Interface (TUI), providing a numbered menu for common operations like managing services, checking logs, and updating configurations. Requires `sudo` for full system access.

**Description:** The primary interactive interface for managing your M3TAL system.
**Usage Example:**
```bash
sudo m3tal
```

### Initial Setup & Health Checks

#### `m3tal init`

Generates the `/etc/m3tal/.env` configuration file from M3TAL's default values. This is typically run only once during the first installation.

**Description:** Initializes the M3TAL environment configuration file.
**Usage Example:**
```bash
m3tal init
```

#### `m3tal doctor`

Performs a pre-flight health check on your M3TAL system. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability, helping diagnose common setup issues.

**Description:** Diagnoses common M3TAL system health issues.
**Usage Example:**
```bash
m3tal doctor
```

### Configuration Management (`m3tal config`)

The `m3tal config` subcommand group manages the `/etc/m3tal/.env` file, which is central to your M3TAL deployment.

#### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring essential variables in `/etc/m3tal/.env`. This is the recommended way to set up your initial M3TAL environment.

**Description:** Interactive wizard to configure the primary `/etc/m3tal/.env` file.
**Usage Example:**
```bash
m3tal config wizard
```

#### `m3tal config set KEY VALUE`

Sets a specific environment variable (`KEY`) to a new `VALUE` in `/etc/m3tal/.env`. Changes take effect after restarting affected services.

**Description:** Sets a single environment variable in the M3TAL configuration.
**Usage Example:**
```bash
# Set the dashboard to expose directly on a specific port
m3tal config set DASHBOARD_PORT 8083

# Switch the dashboard to Traefik mode
m3tal config set DASHBOARD_EXPOSE_MODE traefik
```

#### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable (`KEY`) from `/etc/m3tal/.env`.

**Description:** Retrieves the value of a specific environment variable.
**Usage Example:**
```bash
# Get the current dashboard exposure mode
m3tal config get DASHBOARD_EXPOSE_MODE

# Read the domain configured for Traefik
m3tal config get DOMAIN
```

#### `m3tal config scan`

Lists all known environment variables across all configured Docker Compose stacks, showing their current values (if set) and default values. Useful for understanding what variables are available for configuration.

**Description:** Lists all environment variables relevant to M3TAL, including defaults and current values.
**Usage Example:**
```bash
m3tal config scan
```

#### `m3tal config list`

Displays the entire content of the current `/etc/m3tal/.env` file.

**Description:** Shows the complete contents of the active M3TAL environment file.
**Usage Example:**
```bash
m3tal config list
```

### Dashboard Management (`m3tal dash`)

The `m3tal dash` subcommand group is dedicated to managing the M3TAL Dashboard container (`m3tal-dashboard`).

#### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, it will launch an interactive prompt to set the `admin` password. This command modifies `/docker/users.json`.

**Description:** Manages M3TAL Dashboard user passwords.
**Usage Example:**
```bash
# Set the 'admin' user password interactively
m3tal dashpass

# Set the 'admin' user password directly (replace with strong password)
m3tal dashpass admin MySuperSecureDashboardPass!23

# Set a password for another user (if multiple users are supported)
m3tal dashpass jake P@$$w0rd4Jake
```

#### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub, then starts or updates the `m3tal-dashboard` container based on the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

**Description:** Starts or updates the M3TAL Dashboard container.
**Usage Example:**
```bash
m3tal dash up
```

#### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

**Description:** Stops and removes the M3TAL Dashboard container.
**Usage Example:**
```bash
m3tal dash down
```

#### `m3tal dash restart`

Restarts the `m3tal-dashboard` container. Useful after changing configuration variables affecting the dashboard (e.g., `DASHBOARD_PORT`).

**Description:** Restarts the M3TAL Dashboard container.
**Usage Example:**
```bash
m3tal dash restart
```

#### `m3tal dash logs`

Streams the real-time logs from the `m3tal-dashboard` container. Useful for debugging dashboard-related issues.

**Description:** Streams logs from the M3TAL Dashboard container.
**Usage Example:**
```bash
m3tal dash logs
```

#### `m3tal dash status`

Displays the current status of the `m3tal-dashboard` container (e.g., running, stopped, unhealthy).

**Description:** Shows the operational status of the M3TAL Dashboard container.
**Usage Example:**
```bash
m3tal dash status
```

### Stack Management

These commands interact with all Docker Compose files (`*-compose.yml`) located in the `/docker/` directory.

#### `m3tal up`

Runs `docker compose up -d` for all Docker Compose files found in `/docker/`. This command brings up all your configured services (including Traefik, if `routing-compose.yml` is present, and the M3TAL Dashboard, if its compose file is present).

**Description:** Starts or updates all Docker Compose stacks defined in `/docker/`.
**Usage Example:**
```bash
# Start all services including Traefik, M3TAL Dashboard, and any custom apps
m3tal up
```

#### `m3tal down`

Runs `docker compose down` for all Docker Compose stacks across all `*-compose.yml` files in `/docker/`. This stops and removes all containers, networks, and volumes defined in those files.

**Description:** Stops and removes all Docker Compose stacks.
**Usage Example:**
```bash
# Stop and remove all running services
m3tal down
```

#### `m3tal logs`

Streams aggregated real-time logs from all running Docker containers managed by M3TAL. This provides a unified view of your entire ecosystem's activity.

**Description:** Streams aggregated logs from all running M3TAL-managed containers.
**Usage Example:**
```bash
# View live logs from all your services
m3tal logs
```

---

## M3TAL API Daemon Service Management (systemd)

The M3TAL API daemon (`m3tal-api`) runs as a systemd service, managing internal state and Docker interactions.

#### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api.service`. This shows if the service is active, any recent errors, and its process ID.

**Description:** Displays the status of the M3TAL API systemd service.
**Usage Example:**
```bash
systemctl status m3tal-api
```

#### `journalctl -u m3tal-api -f`

Streams the real-time logs from the `m3tal-api.service`. Essential for debugging issues related to the API daemon itself.

**Description:** Streams logs from the M3TAL API systemd service.
**Usage Example:**
```bash
journalctl -u m3tal-api -f
```

#### Other `systemctl` commands (general knowledge)

```bash
# Restart the M3TAL API service
sudo systemctl restart m3tal-api

# Stop the M3TAL API service
sudo systemctl stop m3tal-api

# Start the M3TAL API service
sudo systemctl start m3tal-api
```

---

## Direct Docker Compose Commands (Fallback)

While M3TAL's CLI simplifies operations, you can always interact directly with Docker Compose as a fallback or for advanced debugging. The M3TAL CLI essentially wraps these commands.

**Important Note:** The `/docker/` directory is a symlink to `/opt/m3tal/stack/`. When using `docker compose` directly, you should `cd` into `/docker/` first or specify the full path to your compose files.

#### Running specific stacks directly

To bring up individual Docker Compose stacks:

```bash
# Change to the stacks directory first
cd /docker/

# Start the core routing and dashboard components (assuming local mode for dashboard)
docker compose -f routing-compose.yml -f m3tal-compose.yml -f m3tal-compose.local.yml up -d

# Start a custom application stack (e.g., 'ollama-compose.yml')
docker compose -f ollama-compose.yml up -d

# Combine multiple specific compose files for a custom setup
docker compose -f routing-compose.yml -f my-app-compose.yml up -d
```

#### Stopping specific stacks directly

```bash
# Stop the core routing components
cd /docker/
docker compose -f routing-compose.yml down

# Stop the dashboard components (assuming local mode for dashboard)
cd /docker/
docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml down

# Stop a custom application stack
cd /docker/
docker compose -f ollama-compose.yml down
```

#### Viewing logs for specific containers

```bash
# Stream logs for the M3TAL Dashboard container
docker logs -f m3tal-dashboard

# Stream logs for the Traefik container
docker logs -f traefik

# Stream logs for a custom container (e.g., 'ollama')
docker logs -f ollama
```

---

## M3TAL Installation

For new installations, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```