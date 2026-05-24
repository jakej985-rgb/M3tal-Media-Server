# `docs/COMMAND_REFERENCE.md`

# M3TAL CLI Command Reference

Greetings, M3TAL Operators! DocSmith here, your dedicated M3TAL Ecosystem Documentation Architect. This document serves as your definitive guide to the M3TAL Command-Line Interface (CLI). From initial setup to day-to-day operations and advanced troubleshooting, this cheat-sheet provides the commands you need to manage your M3TAL environment effectively.

---

## M3TAL Ecosystem Fundamentals

Before diving into commands, let's establish some foundational knowledge about the M3TAL ecosystem's architecture and key concepts.

### Core Components

*   **CLI binary (`/usr/bin/m3tal`)**: Your primary interface, a unified Go binary for all M3TAL operations.
*   **API daemon (`m3tal-api.service`)**: A Go binary running as a systemd service, listening on port `8080`. It's the brain, managing Docker interactions, the state database, and exposing API routes.
*   **Dashboard container (`m3tal-dashboard`)**: A Python/Flask container (internal port `8082`) that communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik Gateway (`routing-compose.yml`)**: Our intelligent reverse proxy, exposing services by domain on port `80`. It uses a file provider for dynamic routing and Docker labels for discovery.
*   **Cloudflared**: An optional Cloudflare tunnel container (`routing-compose.yml`) for secure, zero-config internet access to your services.

### Filesystem Contract

Understanding the critical locations is paramount for effective management:

| Path                        | Purpose                                                                                                               |
| :-------------------------- | :-------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file.** Contains all environment variables for M3TAL and your Docker stacks. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the `m3tal-api.service` daemon.                                    |
| `/opt/m3tal/stack/`         | The canonical directory for all Docker Compose files and Traefik configuration.                                       |
| `/docker`                   | **User-facing symlink** to `/opt/m3tal/stack/`. This is where you place your `*-compose.yml` files for new stacks.     |
| `/docker/users.json`        | Dashboard credential store. Managed exclusively by `m3tal dashpass`. **Do not edit manually.**                          |

### Dashboard Access Modes (Critical!)

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

1.  **Local Mode (Default)**
    *   `DASHBOARD_EXPOSE_MODE=local`
    *   **How it works**: Uses `m3tal-compose.local.yml` to directly bind the dashboard container's port to your host machine (e.g., `8082:8082`). No Traefik required.
    *   **Access via**: `http://YOUR_HOST_IP:8082` or `http://localhost:8082`
    *   **Best for**: LAN-only setups, first-time users, local development, or if you don't use Traefik for other services.

2.  **Traefik Mode**
    *   `DASHBOARD_EXPOSE_MODE=traefik`
    *   **How it works**: Uses `m3tal-compose.traefik.yml` to apply Traefik labels to the dashboard container, allowing Traefik to route `dash.${DOMAIN}` to the dashboard. Traefik (via `routing-compose.yml`) *must* be running.
    *   **Access via**: `http://dash.yourdomain.com` (replace `yourdomain.com` with your `DOMAIN` env var).
    *   **Best for**: Domain-based access, integrating into a multi-service setup behind Traefik, or when public internet access is desired through Traefik/Cloudflared.

### Docker & Compose Runtime

M3TAL relies on **Docker Engine** and **Docker Compose V2** as hard dependencies.
*   The `m3tal up` command orchestrates all `*-compose.yml` files found in the `/docker/` directory using `docker compose up -d`.
*   The `m3tal dash up` command specifically manages the dashboard container, downloading its latest compose configuration and applying the correct override based on `DASHBOARD_EXPOSE_MODE`.
*   To install a new stack, simply place your `my-stack-compose.yml` file into `/docker/` and then run `m3tal up`.

### Traefik Routing Architecture

Traefik (deployed by `routing-compose.yml`) serves as our intelligent traffic cop:
*   It listens on host port `80` (and `443` if HTTPS is configured).
*   It auto-discovers services via Docker labels (`traefik.enable=true`, etc.).
*   It loads dynamic configurations from `/docker/dynamic/` (using a file provider), allowing for hot-reloads of custom routes.
*   **API Daemon Routing**: `api.${DOMAIN}` routes to the internal M3TAL API daemon (Go) at `http://host.docker.internal:8080` via `/docker/dynamic/api.yml`.
*   **Dashboard Routing**: If `DASHBOARD_EXPOSE_MODE=traefik`, `dash.${DOMAIN}` routes to the `m3tal-dashboard` container.

### Port Map

| Port | Service                    | Access                                        |
| :--- | :------------------------- | :-------------------------------------------- |
| 80   | Traefik HTTP entry point   | Public (Traefik mode)                         |
| 8080 | M3TAL API daemon (Go)      | Host-local only                               |
| 8081 | Traefik dashboard          | Host-local only (`127.0.0.1:8081`)            |
| 8082 | M3TAL Dashboard container | Direct port (local mode) or via Traefik (Traefik mode) |

---

## M3TAL CLI Commands

This section details every M3TAL command and its practical usage.

### `sudo m3tal`

The ultimate entry point to the M3TAL TUI (Text-User Interface) Control Center. This interactive menu allows you to manage the entire ecosystem without memorizing every CLI command.

*   **Description**: Opens the interactive TUI Control Center with a numbered menu for common operations like starting/stopping stacks, viewing logs, managing configuration, and more. Requires `sudo` as it interacts with Docker and system-level configurations.
*   **Usage Example**:
    ```bash
    sudo m3tal
    ```

### `m3tal init`

Initializes the M3TAL environment.

*   **Description**: Generates the primary configuration file, `/etc/m3tal/.env`, from default values. This command is crucial for your first installation or if your `.env` file is missing. It sets up the basic structure for your M3TAL deployment.
*   **Usage Example**:
    ```bash
    m3tal init
    ```

### `m3tal doctor`

Your pre-flight health check.

*   **Description**: Performs a comprehensive pre-flight health check of your M3TAL system. It verifies Docker connectivity, checks the validity and integrity of `/etc/m3tal/.env`, and ensures that essential ports (like 8080, 8082, 80) are available and not in use by other processes. Indispensable for troubleshooting startup issues.
*   **Usage Example**:
    ```bash
    m3tal doctor
    ```

---

### M3TAL Configuration Management (`m3tal config`)

These commands are essential for managing `/etc/m3tal/.env`.

#### `m3tal config wizard`

The interactive guide to M3TAL configuration.

*   **Description**: Launches an interactive wizard that guides you through configuring or updating your `/etc/m3tal/.env` file. It presents each environment variable, explains its purpose, and prompts you for a value, making setup straightforward and error-proof.
*   **Usage Example**:
    ```bash
    m3tal config wizard
    ```

#### `m3tal config set KEY VALUE`

Set a specific environment variable.

*   **Description**: Sets a single key-value pair in your `/etc/m3tal/.env` file. This allows for quick, non-interactive adjustments to your configuration.
*   **Usage Example**:
    ```bash
    m3tal config set DOMAIN myhome.lan
    ```

#### `m3tal config get KEY`

Retrieve a specific environment variable.

*   **Description**: Reads and displays the current value of a specified environment variable from `/etc/m3tal/.env`. Useful for quickly checking a setting.
*   **Usage Example**:
    ```bash
    m3tal config get DASHBOARD_EXPOSE_MODE
    ```
    *Expected output*: `local`

#### `m3tal config scan`

List all known configuration variables.

*   **Description**: Scans and lists all known environment variables across all M3TAL and common Docker stack configurations, including their defaults. This provides a comprehensive overview of configurable parameters available in the ecosystem.
*   **Usage Example**:
    ```bash
    m3tal config scan
    ```

#### `m3tal config list`

Display the current `.env` file.

*   **Description**: Displays the entire content of the current `/etc/m3tal/.env` file. This is useful for reviewing your active configuration in its entirety.
*   **Usage Example**:
    ```bash
    m3tal config list
    ```

---

### M3TAL Dashboard Management (`m3tal dash`)

Commands specifically for controlling the M3TAL Dashboard container. Remember to configure `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env` for proper access.

#### `m3tal dashpass [username] [password]`

Manage Dashboard user credentials.

*   **Description**: Updates user credentials for the M3TAL Dashboard, stored in `/docker/users.json`. If `username` and `password` are omitted, it launches an interactive prompt. This is the *only* way to safely manage dashboard user accounts.
*   **Usage Example (Interactive)**:
    ```bash
    m3tal dashpass
    ```
    *Prompts for username and password.*
*   **Usage Example (Direct)**:
    ```bash
    m3tal dashpass admin SuperSecureP@ssw0rd!
    ```

#### `m3tal dash up`

Start the M3TAL Dashboard.

*   **Description**: Pulls the latest dashboard compose configuration (`m3tal-compose.yml` and its overrides) from GitHub, then starts the `m3tal-dashboard` container with the appropriate `DASHBOARD_EXPOSE_MODE` override.
*   **Usage Example**:
    ```bash
    m3tal dash up
    ```

#### `m3tal dash down`

Stop the M3TAL Dashboard.

*   **Description**: Stops and removes the `m3tal-dashboard` container.
*   **Usage Example**:
    ```bash
    m3tal dash down
    ```

#### `m3tal dash restart`

Restart the M3TAL Dashboard.

*   **Description**: Restarts the `m3tal-dashboard` container. This is useful after making configuration changes that affect the dashboard, such as changing `DASHBOARD_EXPOSE_MODE`.
*   **Usage Example**:
    ```bash
    m3tal dash restart
    ```

#### `m3tal dash logs`

Stream M3TAL Dashboard logs.

*   **Description**: Streams real-time logs from the `m3tal-dashboard` container. Essential for debugging dashboard-specific issues.
*   **Usage Example**:
    ```bash
    m3tal dash logs
    ```

#### `m3tal dash status`

Show M3TAL Dashboard status.

*   **Description**: Displays the current status of the `m3tal-dashboard` container (e.g., running, stopped, exited).
*   **Usage Example**:
    ```bash
    m3tal dash status
    ```

---

### M3TAL Global Stack Management

These commands operate on all Docker Compose stacks defined in `/docker/`.

#### `m3tal up`

Bring all M3TAL stacks online.

*   **Description**: Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This command starts all your defined services (Traefik, your custom applications, etc.) in detached mode.
*   **Usage Example**:
    ```bash
    m3tal up
    ```

#### `m3tal down`

Take all M3TAL stacks offline.

*   **Description**: Runs `docker compose down` across all `*-compose.yml` files in `/docker/`, gracefully stopping and removing all containers, networks, and volumes defined by your stacks.
*   **Usage Example**:
    ```bash
    m3tal down
    ```

#### `m3tal logs`

Stream aggregated logs from all running stacks.

*   **Description**: Streams aggregated, real-time logs from all containers managed by M3TAL's Docker Compose stacks. This provides a consolidated view for overall system monitoring and debugging.
*   **Usage Example**:
    ```bash
    m3tal logs
    ```

---

## M3TAL Systemd Service Management

The core M3TAL API daemon runs as a `systemd` service. Use these commands to manage its lifecycle and view its logs.

*   **Check API Service Status**:
    ```bash
    systemctl status m3tal-api
    ```
    *   *Expected output*: Shows active/inactive status, uptime, and recent logs.
*   **Restart API Service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```
    *   *Usage*: Use after modifying `/etc/m3tal/.env` or if the API daemon is misbehaving.
*   **Stream API Service Logs**:
    ```bash
    sudo journalctl -u m3tal-api -f
    ```
    *   *Usage*: Provides a live stream of logs from the M3TAL API daemon, crucial for debugging.

---

## Docker Fallback Commands

While the `m3tal` CLI provides a convenient abstraction, you can always interact directly with Docker Compose V2 if needed. These commands assume you are in the `/docker/` directory or specify the compose files directly.

*   **Start all stacks directly**:
    ```bash
    cd /docker/
    docker compose -f routing-compose.yml -f m3tal-compose.yml -f your-stack-compose.yml up -d
    # Or, if you want to find all automatically:
    docker compose -f $(find . -maxdepth 1 -name "*-compose.yml" -print0 | xargs -0 echo "-f") up -d
    ```
    *   *Note*: The `m3tal up` command handles the discovery of all `*-compose.yml` files automatically.

*   **Stop all stacks directly**:
    ```bash
    cd /docker/
    docker compose -f routing-compose.yml -f m3tal-compose.yml -f your-stack-compose.yml down
    # Or, to find all automatically:
    docker compose -f $(find . -maxdepth 1 -name "*-compose.yml" -print0 | xargs -0 echo "-f") down
    ```

*   **Stream logs from a specific service**:
    ```bash
    cd /docker/
    docker compose -f m3tal-compose.yml logs -f m3tal-dashboard
    ```

*   **Check status of all services**:
    ```bash
    cd /docker/
    docker compose -f routing-compose.yml -f m3tal-compose.yml -f your-stack-compose.yml ps
    ```

---

## M3TAL Installation via APT

For new installations or updates, follow these steps to install the `m3tal` CLI and its associated components.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

That concludes your comprehensive M3TAL CLI cheat-sheet. Remember, DocSmith is always here to help you navigate the M3TALverse. Keep those systems running smoothly!