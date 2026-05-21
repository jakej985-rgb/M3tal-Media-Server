# docs/COMMAND_REFERENCE.md

Greetings, M3TAL Operators! DocSmith here, your M3TAL Ecosystem Documentation Architect.

This document serves as your definitive cheat-sheet for the M3TAL Command Line Interface (CLI). From initial setup to day-to-day operations and advanced diagnostics, this guide will equip you with the knowledge to wield the `m3tal` binary with confidence.

## M3TAL System Overview

The M3TAL ecosystem is designed for streamlined deployment and management of self-hosted services using Docker and Docker Compose. It leverages a robust architecture:

*   **CLI Binary (`/usr/bin/m3tal`)**: Your primary interaction point, a Go binary providing a unified interface.
*   **API Daemon (`m3tal-api.service`)**: A Go binary running as a `systemd` service on port 8080, managing Docker interactions, state, and API routes.
*   **M3TAL Dashboard (`m3tal-dashboard`)**: A Python/Flask Docker container providing a web-based control panel, communicating with the API daemon.
*   **Traefik Gateway (`routing-compose.yml`)**: An optional but recommended reverse proxy for exposing services via domain names on port 80.
*   **Cloudflared Tunnel (`routing-compose.yml`)**: An optional component for secure, zero-config external access via Cloudflare.

### Key Filesystem Contract

Understanding the M3TAL filesystem is crucial for effective operation:

| Path                        | Purpose                                                                                                               |
| :-------------------------- | :-------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary Configuration File**. Stores environment variables for the entire M3TAL ecosystem. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the `m3tal-api` daemon.                                            |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains core `docker-compose.yml` files and Traefik dynamic configurations.               |
| `/docker`                   | **User-facing symlink to `/opt/m3tal/stack/`**. This is where you place *all* your `*-compose.yml` files for M3TAL to manage. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                                                              |

### Dashboard Access Modes

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **`local` Mode (Default)**
    *   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
    *   **Mechanism**: Uses `m3tal-compose.local.yml` to directly bind the dashboard container's port to the host.
    *   **Access**: `http://YOUR_HOST_IP:8082` or `http://localhost:8082`
    *   **Usage**: Ideal for LAN-only setups, first-time users, and local development. No Traefik required.

2.  **`traefik` Mode**
    *   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
    *   **Mechanism**: Uses `m3tal-compose.traefik.yml` to apply Traefik labels, routing `dash.YOUR_DOMAIN` to the dashboard container (internal port 8082). Requires Traefik (`routing-compose.yml`) to be running.
    *   **Access**: `http://dash.YOUR_DOMAIN`
    *   **Usage**: Recommended for domain-based setups, exposing services behind a reverse proxy, and production environments.

### Docker and Compose Runtime

M3TAL orchestrates services using **Docker Engine** and **Docker Compose V2**. These are fundamental dependencies.

*   The `m3tal up` command processes all `*-compose.yml` files found in the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container, dynamically applying the correct compose override based on `DASHBOARD_EXPOSE_MODE`.

### Deployment Lifecycle: Adding New Stacks

To integrate a new service or application into your M3TAL ecosystem:

1.  **Place Compose File**: Create or place your `docker-compose.yml` file (e.g., `my-app-compose.yml`) directly within the `/docker/` directory.
2.  **Configure Environment**: Set any required environment variables for your new stack in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set`.
3.  **Start All Stacks**: Run `m3tal up` to deploy your new service along with all existing M3TAL components.

### Traefik Routing Architecture

When enabled, Traefik acts as the intelligent ingress for your M3TAL services:

*   **Entry Point**: Binds host port 80 for HTTP traffic.
*   **Service Discovery**: Automatically discovers Docker containers with Traefik labels.
*   **Dynamic Configuration**: Loads additional rules from `/docker/dynamic/` (supporting hot-reload).
*   **API Routing**: Routes `api.YOUR_DOMAIN` to the `m3tal-api` daemon (running on host port 8080) via `dynamic/api.yml`.
*   **Dashboard Routing**: Routes `dash.YOUR_DOMAIN` to the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE=traefik`.

### Core Port Map

| Port | Service               | Access Context                                        |
| :--- | :-------------------- | :---------------------------------------------------- |
| 80   | Traefik HTTP Entry    | Public (if Traefik is running and exposed)            |
| 8080 | M3TAL API Daemon (Go) | Host-local (accessed by Docker containers via `host.docker.internal`) |
| 8081 | Traefik Dashboard     | Host-local only (`http://localhost:8081`)             |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

### APT Installation

If you haven't installed M3TAL yet, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## M3TAL CLI Command Reference

### Interactive Control Center

#### `sudo m3tal`

Opens the M3TAL Interactive TUI Control Center. This text-based user interface provides a guided, menu-driven way to manage your M3TAL ecosystem. It requires `sudo` because it interacts with system-level services and Docker.

**Usage:**

```bash
sudo m3tal
```

**Notes:** Navigating the TUI allows you to perform many common operations (like starting/stopping services, viewing logs, configuring settings) without memorizing specific CLI commands.

### System Initialization & Health

#### `m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from M3TAL's internal defaults. This command is crucial for the very first installation or if your `.env` file is missing.

**Usage:**

```bash
m3tal init
```

**Notes:** Always run this after a fresh M3TAL installation to ensure your system has a base configuration. Subsequent configuration should be done via `m3tal config wizard` or `m3tal config set`.

#### `m3tal doctor`

Performs a pre-flight health check of your M3TAL system. It verifies Docker connectivity, checks the validity of `/etc/m3tal/.env`, and ensures essential ports (like 8080 for the API daemon) are available.

**Usage:**

```bash
m3tal doctor
```

**Notes:** Use this command to quickly diagnose common setup issues before starting services or if you encounter unexpected behavior.

### Configuration Management

M3TAL's configuration is primarily managed through `/etc/m3tal/.env`.

#### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended way to adjust M3TAL settings.

**Usage:**

```bash
m3tal config wizard
```

**Notes:** The wizard prompts for common settings like domain, expose modes, paths, and user IDs, validating inputs where possible.

#### `m3tal config set KEY VALUE`

Sets a single environment variable `KEY` to `VALUE` in `/etc/m3tal/.env`. This provides direct, non-interactive control over specific settings.

**Usage:**

```bash
m3tal config set DOMAIN mymetalcluster.net
m3tal config set DASHBOARD_EXPOSE_MODE traefik
m3tal config set PUID 1001
```

**Notes:** Changes made with `m3tal config set` often require a restart of the `m3tal-api` daemon or affected Docker containers (e.g., `m3tal dash restart`, `m3tal up`) to take effect.

#### `m3tal config get KEY`

Reads and displays the current value of a single environment variable `KEY` from `/etc/m3tal/.env`.

**Usage:**

```bash
m3tal config get DASHBOARD_EXPOSE_MODE
m3tal config get TRAEFIK_WEB_PORT
```

#### `m3tal config scan`

Lists all known environment variables across all M3TAL stacks, showing their current values and default if applicable. This provides a comprehensive overview of your system's configuration parameters.

**Usage:**

```bash
m3tal config scan
```

**Notes:** This command is useful for auditing and understanding all potential configuration points, even those not explicitly set in your `.env`.

#### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config list
```

**Notes:** This provides a direct view of your active M3TAL configuration.

### Dashboard Management

These commands specifically interact with the `m3tal-dashboard` Docker container.

#### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, it launches an interactive prompt. The credentials are stored in `/docker/users.json`.

**Usage (Interactive):**

```bash
m3tal dashpass
# Follow prompts for username and new password
```

**Usage (Direct):**

```bash
m3tal dashpass admin SuperS3cureP@ssw0rd!
```

**Notes:** Always use strong, unique passwords. Restart the dashboard (`m3tal dash restart`) after changing passwords for changes to take full effect.

#### `m3tal dash up`

Pulls the latest dashboard `docker-compose` configuration files from GitHub, then starts the `m3tal-dashboard` container using the appropriate `DASHBOARD_EXPOSE_MODE` override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`).

**Usage:**

```bash
m3tal dash up
```

**Notes:** This command ensures you're running the most up-to-date dashboard version and automatically configures it based on your `DASHBOARD_EXPOSE_MODE` setting.

#### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash down
```

#### `m3tal dash restart`

Restarts the `m3tal-dashboard` container. This is useful after making configuration changes or updating user passwords.

**Usage:**

```bash
m3tal dash restart
```

#### `m3tal dash logs`

Streams the aggregated logs from the `m3tal-dashboard` container to your terminal. Press `Ctrl+C` to stop streaming.

**Usage:**

```bash
m3tal dash logs
```

**Notes:** Essential for troubleshooting dashboard-related issues.

#### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` container (e.g., running, stopped, unhealthy).

**Usage:**

```bash
m3tal dash status
```

### Core Stack Management

These commands manage all `*-compose.yml` files located in `/docker/`, representing your entire M3TAL ecosystem's services.

#### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in `/docker/`. This command starts all your M3TAL services, including the Traefik gateway (if configured) and any user-defined stacks.

**Usage:**

```bash
m3tal up
```

**Notes:** This command will also pull new images if available and create/recreate containers as needed. Run this after adding new `compose` files to `/docker/` or making significant configuration changes.

#### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This command stops and removes all containers, networks, and volumes defined by your M3TAL stacks.

**Usage:**

```bash
m3tal down
```

**Notes:** Use with caution as it removes volumes by default. If you want to keep volumes, you'll need to use direct `docker compose` commands (see fallback section).

#### `m3tal logs`

Streams aggregated logs from all currently running M3TAL Docker containers (i.e., all containers started by `m3tal up`). Press `Ctrl+C` to stop streaming.

**Usage:**

```bash
m3tal logs
```

**Notes:** Invaluable for monitoring the health and activity of your entire M3TAL deployment.

---

## Systemd Service Management

The `m3tal-api` daemon is a critical component, running as a `systemd` service. You can manage it directly using `systemctl` commands.

*   **Check API daemon status:**

    ```bash
    sudo systemctl status m3tal-api
    ```

*   **Restart API daemon (after `.env` changes, for example):**

    ```bash
    sudo systemctl restart m3tal-api
    ```

*   **Stream API daemon logs:**

    ```bash
    sudo journalctl -u m3tal-api -f
    ```

---

## Direct Docker Compose Fallback

M3TAL leverages Docker Engine and Docker Compose V2. While the `m3tal` CLI provides a convenient abstraction, you can always use direct `docker compose` commands as a fallback or for more granular control. All M3TAL stack files reside in `/docker/`.

To manage individual stacks or use specific `docker compose` flags not exposed by the `m3tal` CLI:

1.  **Navigate to the stack directory:**

    ```bash
    cd /docker
    ```

2.  **Run `docker compose` commands:**

    *   **Start all services without building images (if already built):**
        ```bash
        docker compose up -d --no-build
        ```
    *   **Start specific services from a particular compose file:**
        ```bash
        docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml up -d m3tal-dashboard
        ```
    *   **Stop and remove services (keeping named volumes):**
        ```bash
        docker compose down --remove-orphans --volumes
        # Note: 'm3tal down' removes volumes by default, this example keeps them.
        ```
    *   **View logs for all services (similar to `m3tal logs`):**
        ```bash
        docker compose logs -f
        ```
    *   **Check status of all services:**
        ```bash
        docker compose ps
        ```
    *   **Pull latest images for a specific compose file:**
        ```bash
        docker compose -f routing-compose.yml pull
        ```

**Important:** When using direct `docker compose` commands, ensure you are in the `/docker/` directory or explicitly specify the compose files (`-f`). M3TAL's `up`/`down` commands handle orchestrating *all* compose files in that directory automatically.