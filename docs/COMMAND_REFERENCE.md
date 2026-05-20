Greetings, M3TAL Operatives! DocSmith here, your M3TAL Ecosystem Documentation Architect.

This document, `docs/COMMAND_REFERENCE.md`, is your comprehensive cheat-sheet for interacting with the M3TAL ecosystem via its command-line interface (CLI). It details every core command, provides real-world usage examples, and outlines critical system architecture for effective Day 2 operations.

Master these commands, understand the underlying architecture, and you'll be orchestrating your digital infrastructure with the precision of a master craftsman.

---

# M3TAL CLI Command Reference

The M3TAL CLI (`/usr/bin/m3tal`) is your primary interface for managing the M3TAL ecosystem. It's a unified Go binary designed for simplicity and power.

## Core Commands

### `sudo m3tal`
Launch the interactive M3TAL TUI (Text-based User Interface) Control Center. This provides a user-friendly, numbered menu for common administrative tasks, guided configuration, and system oversight.
**Purpose**: Interactive system management and oversight.
**Usage Example**:
```bash
sudo m3tal
```
_This command will open a full-screen interactive menu in your terminal._

### `m3tal init`
Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults. This command is crucial for your *first installation* to establish the core environment variables needed by the M3TAL API daemon and services.
**Purpose**: Initialize the system configuration.
**Usage Example**:
```bash
m3tal init
```
_After running, review and customize `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set`._

### `m3tal doctor`
Performs a pre-flight health check of your M3TAL installation. This includes verifying Docker connectivity, validating the `/etc/m3tal/.env` file, and checking for port availability required by M3TAL services.
**Purpose**: Diagnose common setup issues.
**Usage Example**:
```bash
m3tal doctor
```
_This will output a report indicating any detected issues, such as Docker daemon status or port conflicts._

## Configuration Management (`m3tal config`)

M3TAL's configuration is primarily driven by environment variables stored in `/etc/m3tal/.env`. These commands allow you to manage this file.

### `m3tal config wizard`
Launches an interactive, step-by-step wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for new users or for significant configuration changes.
**Purpose**: User-friendly configuration of `/etc/m3tal/.env`.
**Usage Example**:
```bash
m3tal config wizard
```
_The wizard will prompt you for values for various M3TAL environment variables._

### `m3tal config set KEY VALUE`
Sets a specific environment variable (`KEY`) to a designated `VALUE` in `/etc/m3tal/.env`. After setting, you typically need to restart the M3TAL API daemon or relevant Docker containers for changes to take effect.
**Purpose**: Direct modification of a single configuration variable.
**Usage Example**:
```bash
m3tal config set DOMAIN myhome.tech
m3tal config set DASHBOARD_EXPOSE_MODE traefik
```
_Remember to restart affected services (e.g., `sudo systemctl restart m3tal-api` or `m3tal dash restart`) for changes to apply._

### `m3tal config get KEY`
Retrieves and displays the current value of a specific environment variable (`KEY`) from `/etc/m3tal/.env`.
**Purpose**: Inspect a single configuration variable.
**Usage Example**:
```bash
m3tal config get PUID
m3tal config get DOMAIN
```

### `m3tal config scan`
Lists all known environment variables across all M3TAL stacks, including their default values and descriptions where available. This provides a comprehensive overview of configurable parameters.
**Purpose**: Discover all available configuration options.
**Usage Example**:
```bash
m3tal config scan
```

### `m3tal config list`
Displays the current contents of the `/etc/m3tal/.env` file, showing all defined environment variables and their values.
**Purpose**: Review the active configuration file.
**Usage Example**:
```bash
m3tal config list
```

## Dashboard Management (`m3tal dash`)

The M3TAL Dashboard is a central control panel. These commands help you manage its credentials and lifecycle.

### `m3tal dashpass [username] [password]`
Updates the password for a specified dashboard user. If `username` and `password` are omitted, an interactive prompt will guide you. This command modifies the `/docker/users.json` file.
**Purpose**: Manage dashboard user credentials.
**Usage Examples**:
```bash
# Interactive mode
m3tal dashpass

# Set password for 'admin' user directly
m3tal dashpass admin MySuperSecurePassword!
```
_After updating passwords, restarting the dashboard (`m3tal dash restart`) is recommended._

### `m3tal dash up`
This command is specifically designed to manage the M3TAL Dashboard container (`m3tal-dashboard`). It performs the following critical steps:
1.  **Downloads latest compose configurations**: Fetches `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub.
2.  **Reads `DASHBOARD_EXPOSE_MODE`**: Determines the appropriate dashboard access mode (see "M3TAL Dashboard Access Modes" below) from `/etc/m3tal/.env`.
3.  **Starts the dashboard container**: Brings up the `m3tal-dashboard` container using the base compose file and the relevant override (`local` or `traefik`).
**Purpose**: Ensure the dashboard is running with the latest configuration and correct exposure mode.
**Usage Example**:
```bash
m3tal dash up
```
_Ensure `DASHBOARD_EXPOSE_MODE` is correctly set in `/etc/m3tal/.env` before running._

### `m3tal dash down`
Stops and removes the M3TAL Dashboard container.
**Purpose**: Shut down the dashboard.
**Usage Example**:
```bash
m3tal dash down
```

### `m3tal dash restart`
Restarts the M3TAL Dashboard container. This is useful after changing dashboard-related environment variables or credentials.
**Purpose**: Apply new configurations or troubleshoot the dashboard.
**Usage Example**:
```bash
m3tal dash restart
```

### `m3tal dash logs`
Streams the real-time logs from the M3TAL Dashboard container. Useful for debugging and monitoring dashboard activity.
**Purpose**: Monitor dashboard operations.
**Usage Example**:
```bash
m3tal dash logs
```

### `m3tal dash status`
Displays the current operational status of the M3TAL Dashboard container.
**Purpose**: Check if the dashboard is running and healthy.
**Usage Example**:
```bash
m3tal dash status
```

## Stack Management

These commands interact with all Docker Compose stacks deployed in the `/docker/` directory.

### `m3tal up`
Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This command brings up or updates all your configured Docker services in detached mode.
**Purpose**: Start or update all M3TAL-managed Docker stacks.
**Usage Example**:
```bash
m3tal up
```
_This will start services like Traefik (`routing-compose.yml`) and any other custom stacks you've placed in `/docker/`._

### `m3tal down`
Runs `docker compose down` across all `*-compose.yml` files found in `/docker/`. This command stops and removes all M3TAL-managed Docker services and their associated networks.
**Purpose**: Stop and remove all M3TAL-managed Docker stacks.
**Usage Example**:
```bash
m3tal down
```

### `m3tal logs`
Streams aggregated real-time logs from all running Docker containers managed by M3TAL. This provides a consolidated view for overall system monitoring and debugging.
**Purpose**: Centralized logging for all Docker services.
**Usage Example**:
```bash
m3tal logs
```

---

# Core System Architecture & Concepts

Understanding M3TAL's architecture is key to effective management.

## The M3TAL Filesystem Contract

The following paths are critical to M3TAL's operation and configuration:

| Path                     | Purpose                                                                |
| :----------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file for the M3TAL ecosystem. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the M3TAL API daemon. |
| `/opt/m3tal/stack/`      | Canonical stack directory containing core compose files and Traefik config. |
| `/docker`                | **User-facing symlink** to `/opt/m3tal/stack/`. Place all custom `*-compose.yml` files here. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.               |

## M3TAL Dashboard Access Modes

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

To switch modes, set the variable using `m3tal config set DASHBOARD_EXPOSE_MODE <mode>` and then run `m3tal dash up` or `m3tal dash restart` to apply the change.

### Mode 1: `local` (Default)
*   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
*   **Mechanism**: Uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's port 8082 to the host's `DASHBOARD_PORT` (defaulting to 8082).
*   **Access**: `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: No Traefik required.
*   **Best for**: LAN-only setups, first-time users, local development, or environments without a domain.

### Mode 2: `traefik`
*   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
*   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override. This adds Traefik labels to the dashboard container, instructing Traefik (running from `routing-compose.yml`) to route requests for `dash.${DOMAIN}` to the dashboard on its internal port 8082.
*   **Access**: `http://dash.YOUR_DOMAIN` (e.g., `http://dash.myhome.tech`).
*   **Requirements**: Traefik must be running (via `m3tal up`), and the `DOMAIN` variable must be correctly set in `/etc/m3tal/.env`.
*   **Best for**: Domain-based setups, exposing services via a reverse proxy, and integration into a larger web service architecture.

## Docker & Compose Runtime

M3TAL leverages **Docker Engine** and **Docker Compose V2** for container orchestration. These are hard dependencies for the M3TAL ecosystem.

*   The `m3tal up` command executes `docker compose` operations across *all* `*-compose.yml` files found within the `/docker/` directory.
*   The `m3tal dash up` command is a specialized wrapper that ensures the M3TAL Dashboard is deployed with the correct configuration by dynamically selecting the appropriate compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on `DASHBOARD_EXPOSE_MODE`.
*   Your custom Docker Compose stacks should be placed directly in the `/docker/` directory.

## Traefik Routing Architecture

Traefik, the M3TAL reverse proxy, is deployed as a Docker container via `routing-compose.yml`.

*   **Host Port Binding**: Traefik binds to port `80` (and `443` if configured for HTTPS) on the host system, serving as the primary HTTP/S entry point.
*   **Service Discovery**: It automatically discovers and routes to Docker services based on their labels (e.g., for `m3tal-dashboard` in `traefik` mode).
*   **File Provider**: Traefik loads additional dynamic routing configurations from `/etc/traefik/dynamic/` (which in the M3TAL context, maps to `/docker/dynamic/`). This allows for hot-reloading of routing rules.
*   **M3TAL API Routing**: Traefik routes requests for `api.${DOMAIN}` to `http://host.docker.internal:8080`, directly accessing the M3TAL API daemon running on the host system via `dynamic/api.yml`.
*   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.${DOMAIN}` to the `m3tal-dashboard` container.

## Deployment Lifecycle — Day 2 Operations: Installing a New Stack

To add a new application stack to your M3TAL ecosystem:

1.  **Place Compose File**: Create your Docker Compose file (e.g., `my-app-compose.yml`) and place it into the `/docker/` directory.
2.  **Configure Environment Variables**: Ensure all environment variables required by your new stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to add them.
3.  **Start All Stacks**: Run `m3tal up` to deploy your new stack alongside existing M3TAL services.

---

# M3TAL API Daemon & Systemd Service Management

The M3TAL API daemon is a critical Go binary that manages Docker interactions, the state database, and API routes. It runs as a systemd service.

### `systemctl status m3tal-api`
Checks the current status of the `m3tal-api` systemd service, including whether it's active, its uptime, and recent log entries.
**Purpose**: Verify the API daemon's operational status.
**Usage Example**:
```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`
Streams real-time logs from the `m3tal-api` systemd service. This is invaluable for debugging issues with the M3TAL API or its interactions with Docker.
**Purpose**: Live monitoring of the API daemon's logs.
**Usage Example**:
```bash
journalctl -u m3tal-api -f
```

---

# Direct Docker Compose Fallback

While the `m3tal` CLI provides convenient wrappers, you can always interact directly with Docker Compose for troubleshooting or advanced scenarios. Remember that M3TAL operates on all `*-compose.yml` files in `/docker/`.

**Example: Check status of all M3TAL-managed services**:
```bash
docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/my-app-compose.yml ps
```
_You would include `-f` for every `*-compose.yml` file present in `/docker/`._

**Example: Stream logs from a specific M3TAL-managed service (e.g., Traefik)**:
```bash
docker logs -f traefik
```

**Example: Restart a specific service from a stack**:
```bash
docker compose -f /docker/routing-compose.yml restart traefik
```

---

# M3TAL System Ports

| Port | Service                          | Access                                      |
| :--- | :------------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point         | Public (via `routing-compose.yml`)          |
| 8080 | M3TAL API daemon (Go binary)     | Host-local (accessible to dashboard/Traefik)|
| 8081 | Traefik dashboard (admin UI)     | Host-local only (e.g., `http://localhost:8081`) |
| 8082 | M3TAL Dashboard (Python/Flask)   | Direct port (local mode) or via Traefik (traefik mode) |

---

# M3TAL APT Installation

For first-time setup or updating your `m3tal` CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

This concludes your M3TAL CLI cheat-sheet. May your commands be swift and your systems stable!

_DocSmith, M3TAL Ecosystem Documentation Architect_