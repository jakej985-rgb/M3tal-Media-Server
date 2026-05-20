As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the definitive CLI cheat-sheet for the M3TAL ecosystem. This document covers every command, its practical usage, and delves into the underlying architecture and best practices for managing your M3TAL deployment.

---

# M3TAL CLI Command Reference

This document provides a comprehensive guide to the M3TAL Command Line Interface (`m3tal`), your primary tool for interacting with the M3TAL ecosystem.

## Core Commands

### `sudo m3tal`
*   **Description:** Launches the interactive M3TAL Control Center, a text-based user interface (TUI) that provides a guided, menu-driven experience for common M3TAL operations. Requires `sudo` for system-level access.
*   **Usage Example:**
    ```bash
    sudo m3tal
    ```
    *(This command will open the TUI, presenting a numbered menu of system management options.)*

### `m3tal init`
*   **Description:** Initializes the M3TAL ecosystem by generating the primary configuration file, `/etc/m3tal/.env`, from system defaults. This is a critical step during the initial setup of M3TAL. It will not overwrite an existing `.env` file without explicit user confirmation.
*   **Usage Example:**
    ```bash
    sudo m3tal init
    ```
    *(If `/etc/m3tal/.env` is not present, it will be created with default values. If it exists, you will be prompted to confirm an overwrite.)*

### `m3tal doctor`
*   **Description:** Performs a comprehensive pre-flight health check of the M3TAL system. This diagnostic tool verifies Docker connectivity, validates the syntax and content of `/etc/m3tal/.env`, and checks for availability of essential ports (e.g., 80, 8080, 8082). It's indispensable for troubleshooting setup and operational issues.
*   **Usage Example:**
    ```bash
    m3tal doctor
    ```
    *(The output will detail the health of various M3TAL components, highlighting any potential issues or warnings, such as Docker not running or ports being in use.)*

## Configuration Management (`m3tal config`)

M3TAL's primary configuration is stored in `/etc/m3tal/.env`. The `m3tal config` subcommand provides tools for managing this file.

### `m3tal config wizard`
*   **Description:** Initiates an interactive wizard that guides you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for managing your M3TAL settings, offering clear explanations for each environment variable.
*   **Usage Example:**
    ```bash
    sudo m3tal config wizard
    ```
    *(The wizard will prompt you to set variables like `DOMAIN`, `DASHBOARD_EXPOSE_MODE`, `PUID`, `PGID`, and others, providing context for each.)*

### `m3tal config set KEY VALUE`
*   **Description:** Sets a specific environment variable within `/etc/m3tal/.env` to a new value. This command allows for quick, direct modifications of individual settings without launching the full wizard.
*   **Usage Examples:**
    ```bash
    sudo m3tal config set DOMAIN "myhomeserver.net"
    sudo m3tal config set DASHBOARD_EXPOSE_MODE "traefik"
    sudo m3tal config set PUID "1001"
    ```
    *(Remember that changes to `.env` variables often require restarting the affected M3TAL services, e.g., `m3tal dash restart` or `m3tal up`.)*

### `m3tal config get KEY`
*   **Description:** Retrieves and displays the current value of a specified environment variable from `/etc/m3tal/.env`.
*   **Usage Examples:**
    ```bash
    m3tal config get DASHBOARD_EXPOSE_MODE
    # Expected output: traefik

    m3tal config get PUID
    # Expected output: 1000
    ```

### `m3tal config scan`
*   **Description:** Scans all `*-compose.yml` files present in the `/docker/` directory (the symlink to `/opt/m3tal/stack/`) to identify and list all environment variables referenced by M3TAL's Docker stacks. This helps ensure that all necessary variables are defined in your `/etc/m3tal/.env` file.
*   **Usage Example:**
    ```bash
    m3tal config scan
    ```
    *(This command will output a comprehensive list of environment variables used across all discovered Docker Compose files.)*

### `m3tal config list`
*   **Description:** Displays the entire content of the primary M3TAL configuration file, `/etc/m3tal/.env`, directly to your console.
*   **Usage Example:**
    ```bash
    m3tal config list
    ```
    *(This will print all `KEY=VALUE` pairs as they appear in the `.env` file.)*

## Dashboard Management (`m3tal dash`)

The `m3tal dash` subcommand provides dedicated tools for managing the M3TAL dashboard container and its credentials.

### `m3tal dashpass [username] [password]`
*   **Description:** Manages user credentials for accessing the M3TAL dashboard, updating or creating user passwords in `/docker/users.json`. If `username` and `password` arguments are omitted, the command runs interactively, prompting you for input.
*   **Usage Examples:**
    ```bash
    # Interactive mode:
    m3tal dashpass
    # (Prompts for username, password, and confirmation)

    # Non-interactive mode (sets password for 'admin' user):
    m3tal dashpass admin MySecureDashboardPass123
    ```
    *(After changing passwords, it's good practice to restart the dashboard: `m3tal dash restart` for changes to take full effect.)*

### `m3tal dash up`
*   **Description:** Deploys and starts the `m3tal-dashboard` container. This command first fetches the latest dashboard Docker Compose configurations (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub. It then consults `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env` to determine whether to start the dashboard with direct port exposure or via Traefik.
*   **Usage Example:**
    ```bash
    m3tal dash up
    ```
    *(Upon successful execution, the dashboard will be accessible via `http://HOST_IP:8082` (if `DASHBOARD_EXPOSE_MODE=local`) or `http://dash.YOUR_DOMAIN` (if `DASHBOARD_EXPOSE_MODE=traefik` and Traefik is running).)*

### `m3tal dash down`
*   **Description:** Stops and removes the `m3tal-dashboard` container and any associated resources defined in its compose file.
*   **Usage Example:**
    ```bash
    m3tal dash down
    ```

### `m3tal dash restart`
*   **Description:** Restarts the `m3tal-dashboard` container. This is useful after making changes to `/etc/m3tal/.env` (e.g., `DASHBOARD_EXPOSE_MODE`), or updating dashboard credentials with `m3tal dashpass`.
*   **Usage Example:**
    ```bash
    m3tal dash restart
    ```

### `m3tal dash logs`
*   **Description:** Streams real-time log output from the `m3tal-dashboard` container. Essential for diagnosing and monitoring dashboard-specific issues.
*   **Usage Example:**
    ```bash
    m3tal dash logs
    ```
    *(Press `Ctrl+C` to stop streaming logs.)*

### `m3tal dash status`
*   **Description:** Displays the current operational status of the `m3tal-dashboard` container (e.g., `running`, `exited`, `restarting`).
*   **Usage Example:**
    ```bash
    m3tal dash status
    # Expected output: Container m3tal-dashboard is running.
    ```

## Stack Management (`m3tal`)

These commands manage all Docker Compose stacks defined within the M3TAL ecosystem.

### `m3tal up`
*   **Description:** Orchestrates the entire M3TAL Docker environment. This command runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). It starts or recreates all services defined in these files in detached mode, including core components like Traefik and Cloudflared (if configured), along with any user-defined custom stacks.
*   **Usage Example:**
    ```bash
    m3tal up
    ```
    *(This will bring up all your M3TAL-managed Docker containers, ensuring they are running according to their compose definitions.)*

### `m3tal down`
*   **Description:** Stops and removes all Docker containers, networks, and volumes defined by `*-compose.yml` files within the `/docker/` directory. This command effectively brings down your entire M3TAL Docker environment.
*   **Usage Example:**
    ```bash
    m3tal down
    ```
    *(Use this command with caution, as it will halt all M3TAL-managed services.)*

### `m3tal logs`
*   **Description:** Streams aggregated, real-time log output from all currently running Docker containers managed by M3TAL (i.e., all containers launched by `m3tal up`). This provides a consolidated view for comprehensive system monitoring and debugging.
*   **Usage Example:**
    ```bash
    m3tal logs
    ```
    *(Logs from different containers will be prefixed with their respective container names for easy identification. Press `Ctrl+C` to stop streaming.)*

---

## M3TAL System Architecture (Ground Truth)

M3TAL is a unified system designed for efficient management of containerized services.

### Components
-   **CLI binary** (`/usr/bin/m3tal`): The central Go binary for all user interactions, installed via APT.
-   **API daemon** (`m3tal-api.service`): A Go binary running as a `systemd` service on port `8080`. It's the backend for Docker management, state tracking via SQLite, and exposes API routes.
-   **Dashboard container** (`m3tal-dashboard`): A Python/Flask Docker container running internally on port `8082`. It communicates with the host's M3TAL API daemon at `http://host.docker.internal:8080`.
-   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It uses Docker labels and a file provider (`/etc/traefik/dynamic`) for dynamic routing.
-   **Cloudflared** (`routing-compose.yml`): An optional Docker container that establishes secure Cloudflare tunnels for zero-configuration internet access.

### Filesystem Contract

| Path                      | Purpose                                                                                                                                                                                                                                                                 |
| :------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | **Primary M3TAL Configuration File.** Contains all environment variables for the M3TAL system and its services. Managed by `m3tal config wizard` and `m3tal config set`.                                                                                                 |
| `/var/lib/m3tal/state.db` | **SQLite State Database.** Automatically created and managed by the `m3tal-api.service` daemon for internal M3TAL state tracking.                                                                                                                                       |
| `/opt/m3tal/stack/`       | **Canonical Stack Directory.** The base directory containing all core M3TAL Docker Compose files (`*-compose.yml`) and Traefik dynamic configuration.                                                                                                                     |
| `/docker`                 | **User-Facing Stack Directory Symlink.** This is a symlink to `/opt/m3tal/stack/`. Users should place their custom `*-compose.yml` files here to integrate new services into the M3TAL ecosystem via `m3tal up`.                                                           |
| `/docker/users.json`      | **Dashboard Credential Store.** A JSON file holding hashed usernames and passwords for dashboard access. Managed exclusively by the `m3tal dashpass` command.                                                                                                             |

### Dashboard Access — Two Modes (Critical)

The M3TAL dashboard offers two distinct access modes, determined by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)
-   **Configuration:** Set `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
-   **Mechanism:** Uses the `m3tal-compose.local.yml` override, adding a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the dashboard container.
-   **Access Via:** `http://HOST_IP:8082` or `http://localhost:8082`.
-   **Requirements:** No Traefik or domain setup needed.
-   **Best for:** LAN-only setups, initial testing, or users who prefer direct port access.

#### Mode 2: `traefik`
-   **Configuration:** Set `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
-   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard container. These labels configure Traefik to route `dash.${DOMAIN}` to the dashboard's internal port `8082`.
-   **Access Via:** `http://dash.YOUR_DOMAIN` (where `YOUR_DOMAIN` is set in `/etc/m3tal/.env`).
-   **Requirements:** Traefik must be running (via `m3tal up`), and `DOMAIN` must be configured.
-   **Best for:** Domain-based access, integrating with other services behind a reverse proxy, and leveraging Traefik's advanced features.

### Docker / Compose Runtime
-   M3TAL relies on **Docker Engine** and **Docker Compose V2**. These are hard dependencies.
-   The `m3tal up` command orchestrates all `*-compose.yml` files in `/docker/` using `docker compose up -d`.
-   `m3tal dash up` specifically manages the dashboard: it downloads the latest compose files, reads `DASHBOARD_EXPOSE_MODE`, and starts the dashboard with the appropriate override.
-   User-defined Docker Compose stacks should be placed in `/docker/`.

### Deployment Lifecycle — Day 2 Operations

**Installing a new stack:**
1.  Place your `my-stack-compose.yml` file into `/docker/`.
2.  Ensure all necessary environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start your new stack alongside existing services.

### Traefik Routing Architecture
Traefik is deployed via `routing-compose.yml`.
-   Binds host port `80` as the HTTP entry point.
-   Discovers services via Docker labels.
-   Loads dynamic configuration from `/docker/dynamic/` (file provider, hot-reload).
-   Routes `api.DOMAIN` to `http://host.docker.internal:8080` (M3TAL API daemon) via `dynamic/api.yml`.
-   Routes `dash.DOMAIN` to the dashboard container via Traefik labels in `m3tal-compose.traefik.yml` (only when `DASHBOARD_EXPOSE_MODE=traefik`).

**Traefik static config (traefik.yml extract):**
```yaml
entryPoints:
  web:
    address: ":80"
providers:
  docker:
    exposedByDefault: false
    network: proxy
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

**Dynamic routing example (dynamic/api.yml):**
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
          - url: "http://host.docker.internal:8080"
```

### Port Map

| Port | Service                               | Access                                                                      |
| :--- | :------------------------------------ | :-------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (via Traefik, if `TRAEFIK_WEB_PORT=80`)                              |
| 8080 | M3TAL API daemon (Go)                 | Host-local only. Accessed by the Dashboard (`host.docker.internal:8080`) and Traefik (`api.DOMAIN`). |
| 8081 | Traefik dashboard (internal)          | Host-local only. Access via `http://localhost:8081`.                        |
| 8082 | M3TAL Dashboard (internal container)  | Direct (`http://HOST_IP:8082`) if `DASHBOARD_EXPOSE_MODE=local`. Via Traefik (`http://dash.DOMAIN`) if `DASHBOARD_EXPOSE_MODE=traefik`. |

---

## Systemd Service Management

The M3TAL API daemon (`m3tal-api.service`) is a critical component managed by `systemd`.

### `systemctl status m3tal-api`
*   **Description:** Checks the current status of the M3TAL API daemon, showing if it's active, running, or failed, along with recent log entries.
*   **Usage Example:**
    ```bash
    systemctl status m3tal-api
    ```
    *(Expected output will indicate `Active: active (running)` if the service is operational.)*

### `journalctl -u m3tal-api -f`
*   **Description:** Streams real-time log output from the M3TAL API daemon. This is invaluable for live debugging and monitoring of the API's operations.
*   **Usage Example:**
    ```bash
    journalctl -u m3tal-api -f
    ```
    *(Press `Ctrl+C` to exit the log stream.)*

### Other Useful Systemd Commands:
```bash
# Restart the M3TAL API daemon
sudo systemctl restart m3tal-api

# Stop the M3TAL API daemon
sudo systemctl stop m3tal-api

# Start the M3TAL API daemon
sudo systemctl start m3tal-api
```

---

## Docker Fallback Commands

M3TAL is built on **Docker Engine** and **Docker Compose V2**. While the `m3tal` CLI abstracts many Docker interactions, you can always use direct `docker compose` commands for advanced control or debugging.

All M3TAL's Docker Compose files are located in `/docker/` (a symlink to `/opt/m3tal/stack/`). The `m3tal up` command essentially constructs and executes a `docker compose` command referencing all `*-compose.yml` files in this directory.

To interact directly with M3TAL's Docker environment:

*   **To start all M3TAL stacks manually (equivalent to `m3tal up`):**
    ```bash
    # Navigate to the stacks directory
    cd /docker/
    # Dynamically find all compose files and build the command
    COMPOSE_FILES=$(find . -maxdepth 1 -name "*-compose.yml" | sort | sed 's/^/-f /' | tr '\n' ' ')
    # Execute docker compose up
    docker compose ${COMPOSE_FILES} up -d
    ```

*   **To stop all M3TAL stacks manually (equivalent to `m3tal down`):**
    ```bash
    # From the /docker/ directory
    docker compose ${COMPOSE_FILES} down
    ```

*   **To view real-time logs for a specific container (e.g., Traefik):**
    ```bash
    docker logs -f traefik
    ```

*   **To restart a specific M3TAL container (e.g., `m3tal-dashboard`):**
    ```bash
    docker restart m3tal-dashboard
    ```

*   **To inspect detailed information about a container:**
    ```bash
    docker inspect m3tal-dashboard
    ```

*   **To list all running containers managed by M3TAL:**
    ```bash
    docker ps --filter "label=m3tal.stack"
    ```

---

## APT Installation

To install the M3TAL CLI binary and systemd service on your Debian/Ubuntu based system, follow these steps:

```bash
# 1. Add the GPG signing key for package verification
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository to your sources list
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update your package lists and install the m3tal package
sudo apt update && sudo apt install -y m3tal
```