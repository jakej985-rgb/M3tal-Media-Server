# M3TAL CLI Command Reference

Welcome to the M3TAL Ecosystem Command Line Interface (CLI) cheat-sheet. This document provides a comprehensive reference for all `m3tal` commands, designed to help you manage your self-hosted services with ease.

M3TAL orchestrates your home lab or server infrastructure using Docker Engine and Docker Compose V2, with a Go-based API daemon for advanced control and a web dashboard for intuitive management.

## M3TAL Ecosystem Overview

Before diving into the commands, here's a quick overview of the M3TAL architecture:

### Components
*   **CLI binary** (`/usr/bin/m3tal`): Your primary interface for interacting with M3TAL.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service (port 8080), managing Docker, the state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask Docker container (internal port 8082) providing a web UI, communicating with the API daemon.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services on port 80 via domain names.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for secure, zero-config internet access.

### Filesystem Contract
*   `/etc/m3tal/.env`: The primary configuration file, storing all environment variables.
*   `/var/lib/m3tal/state.db`: The SQLite state database, automatically created and managed by the API daemon.
*   `/opt/m3tal/stack/`: The canonical directory for M3TAL's core Docker Compose files and Traefik dynamic configuration.
*   `/docker`: A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path where you should place all your custom `*-compose.yml` files.
*   `/docker/users.json`: The credential store for the M3TAL dashboard.

### Dashboard Access Modes (Critical)

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **Local Mode (Default)**
    *   Set `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
    *   Uses the `m3tal-compose.local.yml` override, which directly maps a host port.
    *   **Access via**: `http://YOUR_HOST_IP:8082` (or `http://localhost:8082`).
    *   **Best for**: LAN-only setups, initial configuration, or environments without a domain/reverse proxy. No Traefik required.

2.  **Traefik Mode**
    *   Set `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
    *   Uses the `m3tal-compose.traefik.yml` override, adding Traefik labels to the dashboard container.
    *   **Access via**: `http://dash.YOUR_DOMAIN` (requires Traefik to be running via `m3tal up`).
    *   **Best for**: Domain-based access, integrating with other services behind Traefik.

### Docker / Compose Runtime Explained
M3TAL is built on Docker Engine and Docker Compose V2.
*   The `m3tal up` command orchestrates all `*-compose.yml` files found in the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container. It ensures the latest compose configurations are pulled from GitHub and starts the dashboard using the appropriate override file based on `DASHBOARD_EXPOSE_MODE`.
*   To install new services, simply place their `*-compose.yml` files in `/docker/` and run `m3tal up`.

### Port Map
| Port | Service | Access |
|------|---------|--------|
| 80 | Traefik HTTP entry point | Public (when Traefik is active) |
| 8080 | M3TAL API daemon (Go) | Host-local |
| 8081 | Traefik dashboard | Host-local only |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) |

---

## M3TAL CLI Command Reference

This section details every `m3tal` command, complete with descriptions and real usage examples.

### `sudo m3tal`

Opens the interactive M3TAL TUI (Text-based User Interface) Control Center. This provides a numbered menu for common operations, configuration, and monitoring. Requires `sudo` as it interacts with Docker.

**Usage Example:**
```bash
sudo m3tal
```

---

### `m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from default values. This command should be run on a first-time installation to set up the necessary environment variables.

**Usage Example:**
```bash
m3tal init
```

---

### `m3tal doctor`

Performs a comprehensive pre-flight health check of your M3TAL system. It verifies Docker connectivity, checks the validity and permissions of `/etc/m3tal/.env`, and ensures critical ports are available on the host system.

**Usage Example:**
```bash
m3tal doctor
```

---

### `m3tal config wizard`

Launches an interactive, step-by-step wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for initial setup and subsequent modifications.

**Usage Example:**
```bash
m3tal config wizard
```

---

### `m3tal config set KEY VALUE`

Sets or updates a single environment variable within the `/etc/m3tal/.env` configuration file. This allows for precise, non-interactive changes.

**Usage Example:**
```bash
m3tal config set DASHBOARD_EXPOSE_MODE traefik
m3tal config set DOMAIN myhomelab.net
```

---

### `m3tal config get KEY`

Retrieves and displays the current value of a specified environment variable from `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal config get DASHBOARD_PORT
m3tal config get PUID
```

---

### `m3tal config scan`

Scans all Docker Compose files (including core M3TAL stacks and user-defined ones in `/docker/`) to list all recognized environment variables and their default values. This is useful for understanding all available configuration options.

**Usage Example:**
```bash
m3tal config scan
```

---

### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` configuration file.

**Usage Example:**
```bash
m3tal config list
```

---

### `m3tal dashpass [username] [password]`

Manages user passwords for the M3TAL Dashboard, stored in `/docker/users.json`. If `username` and `password` are omitted, the command will prompt for interactive input.

**Usage Examples:**
*   **Interactive password update:**
    ```bash
    m3tal dashpass
    # Prompts for username and new password
    ```
*   **Direct password update for 'admin' user:**
    ```bash
    m3tal dashpass admin SuperSecretP@ssw0rd!
    ```

---

### `m3tal dash up`

Pulls the latest M3TAL Dashboard Docker Compose configuration files from GitHub, then starts the `m3tal-dashboard` container. This command intelligently uses the appropriate compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE` set in `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal dash up
```

---

### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash down
```

---

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash restart
```

---

### `m3tal dash logs`

Streams the logs from the `m3tal-dashboard` container, useful for debugging and monitoring its activity.

**Usage Example:**
```bash
m3tal dash logs
```

---

### `m3tal dash status`

Displays the current status of the `m3tal-dashboard` container (e.g., running, stopped, exited).

**Usage Example:**
```bash
m3tal dash status
```

---

### `m3tal up`

Orchestrates all Docker Compose stacks defined by `*-compose.yml` files found in the `/docker/` directory. This command performs a `docker compose up -d` operation, ensuring all services are started or updated in detached mode.

**Usage Example:**
```bash
m3tal up
```

---

### `m3tal down`

Stops and removes all Docker Compose stacks defined by `*-compose.yml` files in the `/docker/` directory. This performs a `docker compose down` operation.

**Usage Example:**
```bash
m3tal down
```

---

### `m3tal logs`

Aggregates and streams the logs from all currently running Docker containers managed by M3TAL. This provides a unified view for monitoring your entire ecosystem.

**Usage Example:**
```bash
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon runs as a `systemd` service named `m3tal-api.service`. You can manage this service directly using `systemctl` and view its logs with `journalctl`.

*   **Check API daemon status:**
    ```bash
    systemctl status m3tal-api
    ```

*   **Restart API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```

*   **Stream API daemon logs:**
    ```bash
    sudo journalctl -u m3tal-api -f
    ```

---

## Direct Docker Compose Commands (Fallback)

While `m3tal` commands abstract Docker Compose, you can always use `docker compose` directly as a fallback, especially for debugging or advanced scenarios. Remember that `m3tal` commands specifically manage environment variables from `/etc/m3tal/.env` and handle multi-compose setups.

To interact with a specific stack (e.g., the `routing` stack defined by `routing-compose.yml` and `m3tal-compose.yml` for the dashboard), you typically combine the relevant files. For M3TAL's core operations, the current working directory for compose commands is `/docker/`.

**Examples:**

*   **Start all containers directly (similar to `m3tal up`):**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml up -d
    ```
    *(Note: You'd need to manually specify all relevant compose files and override files based on your `DASHBOARD_EXPOSE_MODE` if not using `m3tal up`.)*

*   **Stop specific containers (e.g., Traefik and Cloudflared):**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml down
    ```

*   **Stream logs from a specific container (e.g., m3tal-dashboard):**
    ```bash
    sudo docker logs -f m3tal-dashboard
    ```

*   **Check status of a specific container:**
    ```bash
    sudo docker ps -f name=m3tal-dashboard
    ```

---

## APT Installation

To install the M3TAL CLI and core components, follow these steps:

1.  **Add the GPG signing key**
    ```bash
    curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg
    ```

2.  **Add the APT repository**
    ```bash
    echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list
    ```

3.  **Install**
    ```bash
    sudo apt update && sudo apt install -y m3tal
    ```