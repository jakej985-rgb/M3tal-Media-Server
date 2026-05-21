```markdown
# M3TAL CLI Command Reference

Welcome, Operator. This document serves as your definitive M3TAL CLI cheat-sheet, detailing every command, its purpose, and practical usage examples. M3TAL streamlines the deployment and management of your self-hosted services using Docker Compose.

---

## M3TAL Core Concepts

Before diving into commands, understand these foundational elements of the M3TAL ecosystem:

*   **CLI Binary (`m3tal`)**: Your primary interface, installed at `/usr/bin/m3tal`.
*   **API Daemon (`m3tal-api.service`)**: A Go-based backend service running on the host (port 8080), managed by systemd. It handles Docker interactions, state management, and API routes for the Dashboard.
*   **M3TAL Dashboard (`m3tal-dashboard`)**: A Python/Flask Docker container that provides a user-friendly web interface. It communicates with the `m3tal-api` daemon.
*   **Docker Engine & Compose V2**: M3TAL is built on Docker. All container orchestration uses Docker Compose V2.
*   **Stacks**: Collections of services defined in `*-compose.yml` files located in `/docker/`.
*   **Configuration (`/etc/m3tal/.env`)**: The central configuration file for the entire M3TAL system, managed by the `m3tal config` commands.

---

## Filesystem Contract

The following paths are critical to M3TAL's operation and are explicitly managed or referenced by the system:

| Path                        | Purpose                                                                                                                                                                                                  |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file.** Stores all environment variables that configure the M3TAL system and its services. Managed by `m3tal config wizard`.                                                       |
| `/var/lib/m3tal/state.db`   | SQLite state database. Stores internal M3TAL state, service information, and other dynamic data. Auto-created and managed by the `m3tal-api` daemon.                                                         |
| `/opt/m3tal/stack/`         | Canonical stack directory. This is where M3TAL stores its core compose files (`m3tal-compose.yml`, `routing-compose.yml`) and related configurations (like Traefik dynamic configs).                       |
| `/docker`                   | **User-facing stack directory.** This is a symlink to `/opt/m3tal/stack/`. All user-defined `*-compose.yml` files should be placed here. `m3tal up` processes all compose files in this directory.          |
| `/docker/users.json`        | Dashboard credential store. Contains usernames and hashed passwords for accessing the M3TAL Dashboard. Managed by `m3tal dashpass`.                                                                       |
| `/docker/dynamic/`          | Directory within `/docker` where Traefik's file provider looks for dynamic routing configurations (e.g., `api.yml`). These files are hot-reloaded by Traefik.                                               |

---

## M3TAL Core Commands

### `sudo m3tal`

Opens the interactive Text User Interface (TUI) Control Center. This provides a menu-driven interface for common M3TAL operations, guided by numerical selections.

**Description:** Provides an interactive, menu-driven experience for managing your M3TAL system. Requires `sudo` as it performs system-level operations.

**Usage Example:**

```bash
sudo m3tal
```

### `m3tal init`

Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults.

**Description:** This command should be run immediately after installing M3TAL. It sets up the essential `.env` file with default values, which you can then customize using `m3tal config wizard` or `m3tal config set`. If `/etc/m3tal/.env` already exists, it will prompt for confirmation before overwriting.

**Usage Example:**

```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of the M3TAL environment.

**Description:** Runs a series of diagnostic checks to ensure the system is ready for M3TAL operations. This includes verifying Docker connectivity, checking the validity and existence of `/etc/m3tal/.env`, and confirming required port availability. Essential for troubleshooting initial setup or unexpected behavior.

**Usage Example:**

```bash
m3tal doctor
```

---

## Configuration Management (`m3tal config`)

These commands manage the central configuration file located at `/etc/m3tal/.env`.

### `m3tal config wizard`

Launches an interactive wizard to configure `/etc/m3tal/.env`.

**Description:** A guided, step-by-step process to review and modify the key environment variables in your M3TAL configuration file. Ideal for initial setup or comprehensive adjustments.

**Usage Example:**

```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env`.

**Description:** Allows direct modification of specific configuration keys. After setting a value, it's often necessary to restart relevant services (e.g., `m3tal dash restart`, `m3tal up`, or `systemctl restart m3tal-api`) for changes to take effect.

**Usage Example:**

```bash
m3tal config set DOMAIN "myhomeserver.net"
m3tal config set PUID "1001"
```

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.

**Description:** Quickly retrieve the current setting for a specific configuration key.

**Usage Example:**

```bash
m3tal config get DASHBOARD_EXPOSE_MODE
m3tal config get PUID
```

### `m3tal config scan`

Lists all discoverable environment variables, including their defaults.

**Description:** Provides a comprehensive list of all known environment variables across all defined stacks and M3TAL components. This is useful for understanding all available configuration options, even those not explicitly set in your `.env` file.

**Usage Example:**

```bash
m3tal config scan
```

### `m3tal config list`

Lists the current contents of the `/etc/m3tal/.env` file.

**Description:** Displays the exact key-value pairs as they are currently configured in your primary M3TAL environment file.

**Usage Example:**

```bash
m3tal config list
```

---

## Dashboard Management (`m3tal dash`)

These commands specifically interact with the `m3tal-dashboard` container.

### `m3tal dashpass [username] [password]`

Updates the password for a M3TAL Dashboard user.

**Description:** Manages user credentials stored in `/docker/users.json`. If no `username` and `password` are provided, the command will prompt you interactively for them. Use this to set up or change access to your M3TAL Dashboard.

**Usage Example (Interactive):**

```bash
m3tal dashpass
# Prompts for username and new password
```

**Usage Example (Direct):**

```bash
m3tal dashpass admin SuperSecureP@ssw0rd!
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub, then starts the dashboard container.

**Description:** This is the primary command to deploy or update the M3TAL Dashboard. It automatically fetches the latest `m3tal-compose.yml` and its overrides (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`). It then reads your `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env` to determine whether to use local port binding or Traefik for access, ensuring the dashboard is started correctly.

**Usage Example:**

```bash
m3tal dash up
```

### `m3tal dash down`

Stops and removes the M3TAL Dashboard container.

**Description:** Gracefully shuts down the `m3tal-dashboard` container.

**Usage Example:**

```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the M3TAL Dashboard container.

**Description:** Stops and then immediately starts the `m3tal-dashboard` container. Useful after configuration changes or for general troubleshooting.

**Usage Example:**

```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams aggregated logs from the M3TAL Dashboard container.

**Description:** Displays real-time output from the `m3tal-dashboard` container. Essential for monitoring dashboard operations and debugging issues.

**Usage Example:**

```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the M3TAL Dashboard container.

**Description:** Provides a quick overview of whether the `m3tal-dashboard` container is running, stopped, or in another state.

**Usage Example:**

```bash
m3tal dash status
```

---

## Stack Management (Global)

These commands interact with all Docker Compose stacks defined in `/docker/`.

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files in `/docker/`.

**Description:** This command orchestrates the deployment and update of all services defined in your M3TAL environment. It scans the `/docker/` directory for all files matching `*-compose.yml` (e.g., `routing-compose.yml`, `my-app-compose.yml`) and brings them up. Use this after adding new compose files, modifying existing ones, or making changes to `/etc/m3tal/.env` that affect multiple services.

**Usage Example:**

```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` across all stacks defined in `/docker/`.

**Description:** Stops and removes all containers, networks, and volumes defined by the `*-compose.yml` files in `/docker/`. This effectively brings down your entire M3TAL-managed service infrastructure.

**Usage Example:**

```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all running M3TAL-managed Docker stacks.

**Description:** Provides a unified view of real-time log output from all containers managed by M3TAL. This is invaluable for system-wide monitoring and debugging.

**Usage Example:**

```bash
m3tal logs
```

---

## M3TAL Dashboard Access Modes

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. It's crucial to understand these to access your dashboard correctly.

### Mode 1: `local` (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   **Mechanism**: Uses the `m3tal-compose.local.yml` override file. This adds a direct port binding to the dashboard container, exposing it directly on the host machine.
*   **Access URL**: `http://HOST_IP:8082` or `http://localhost:8082` (if accessing from the host itself). The port can be customized via `DASHBOARD_PORT` in `/etc/m3tal/.env`.
*   **Requirements**: No Traefik instance is required for dashboard access in this mode.
*   **Best For**: LAN-only setups, first-time users, local development, or when you don't need domain-based access.

**To switch to local mode:**

```bash
m3tal config set DASHBOARD_EXPOSE_MODE local
m3tal dash restart
```

### Mode 2: `traefik`

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override file. This adds specific Docker labels to the dashboard container, allowing the Traefik reverse proxy to discover and route traffic to it.
*   **Access URL**: `http://dash.YOUR_DOMAIN` (where `YOUR_DOMAIN` is set via the `DOMAIN` variable in `/etc/m3tal/.env`). Traefik routes requests to the dashboard's internal port (8082).
*   **Requirements**: Traefik must be running as part of your M3TAL stacks (typically via `m3tal up`, which starts `routing-compose.yml`).
*   **Best For**: Domain-based access, integrating the dashboard with other services behind a single reverse proxy, or when you want to utilize Cloudflared tunnels.

**To switch to Traefik mode:**

```bash
m3tal config set DASHBOARD_EXPOSE_MODE traefik
m3tal config set DOMAIN "myhomeserver.net" # Set your actual domain
m3tal dash restart
m3tal up # Ensure Traefik is running
```

---

## M3TAL Systemd Service Management

The `m3tal-api` daemon runs as a systemd service. Use standard `systemctl` commands to manage it.

*   **Check API daemon status:**

    ```bash
    systemctl status m3tal-api
    ```

*   **Restart API daemon:**

    ```bash
    sudo systemctl restart m3tal-api
    ```

*   **View API daemon logs in real-time:**

    ```bash
    journalctl -u m3tal-api -f
    ```

---

## Docker Direct Commands (Fallback)

M3TAL leverages Docker Engine and Docker Compose V2. While the `m3tal` CLI provides a convenient abstraction, you can always interact with Docker directly as a fallback or for advanced debugging.

**Note:** When using direct `docker compose` commands, ensure you are in the `/docker/` directory or specify the compose file paths.

*   **View status of all M3TAL-managed containers:**

    ```bash
    docker compose -f /docker/m3tal-compose.yml -f /docker/routing-compose.yml -f /docker/YOUR_OTHER_STACK.yml ps
    ```
    *(Adjust `-f` flags to include all relevant compose files)*

*   **Start specific stacks (e.g., routing and m3tal-dashboard):**

    ```bash
    docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml up -d
    ```

*   **Stop specific stacks:**

    ```bash
    docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml down
    ```

*   **Stream logs from all running containers:**

    ```bash
    docker compose -f /docker/m3tal-compose.yml -f /docker/routing-compose.yml logs -f
    ```

*   **Restart a specific container (e.g., M3TAL Dashboard):**

    ```bash
    docker restart m3tal-dashboard
    ```

*   **Access Traefik dashboard (if running):**

    ```bash
    http://localhost:8081
    ```
    *(Note: This is typically for internal host-local access, as defined in `routing-compose.yml`)*

---

## APT Installation

To install the M3TAL CLI and core components via APT, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

This concludes the M3TAL CLI Command Reference. For further assistance, consult the M3TAL documentation or community channels.
```