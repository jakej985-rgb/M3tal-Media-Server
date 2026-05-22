As DocSmith, the M3TAL Ecosystem Documentation Architect, I present this comprehensive CLI cheat-sheet. It serves as your definitive guide to managing your M3TAL server, covering every command with practical, real-world usage examples.

***

# M3TAL Command Reference

The M3TAL CLI binary (`/usr/bin/m3tal`) is your single entry point for all system operations, from configuration to service management. It orchestrates Docker containers, manages system configurations, and provides crucial insights into your M3TAL ecosystem.

## M3TAL Core Commands

### 1. `sudo m3tal`
Opens the interactive Text-User Interface (TUI) Control Center. This provides a numbered menu for common administrative tasks, offering a guided experience for system management. Requires `sudo` for full access to Docker and system configurations.

**Usage Example:**
```bash
sudo m3tal
```

### 2. `m3tal init`
Initializes the M3TAL system by generating the primary configuration file, `/etc/m3tal/.env`, from system defaults. This command is crucial for your first installation and should be run once.

**Usage Example:**
```bash
m3tal init
```

### 3. `m3tal doctor`
Performs a pre-flight health check of your M3TAL installation. This includes verifying Docker connectivity, validating the `/etc/m3tal/.env` configuration, and checking for port availability to ensure a smooth operation.

**Usage Example:**
```bash
m3tal doctor
```

## Configuration Management (`m3tal config`)

M3TAL's configuration is primarily driven by environment variables stored in `/etc/m3tal/.env`. The `m3tal config` subcommand provides tools to manage this file.

### 1. `m3tal config wizard`
Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for making comprehensive configuration changes.

**Usage Example:**
```bash
m3tal config wizard
```

### 2. `m3tal config set KEY VALUE`
Sets a single environment variable in `/etc/m3tal/.env` to a specified value. This is useful for quick, atomic changes.

**Usage Example:**
```bash
m3tal config set DASHBOARD_PORT 8083
```

### 3. `m3tal config get KEY`
Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal config get DASHBOARD_EXPOSE_MODE
```

### 4. `m3tal config scan`
Lists all environment variables known to M3TAL, including those defined across your Docker Compose stacks. This provides a holistic view of your system's configuration.

**Usage Example:**
```bash
m3tal config scan
```

### 5. `m3tal config list`
Displays the entire contents of the `/etc/m3tal/.env` file, showing all currently configured environment variables and their values.

**Usage Example:**
```bash
m3tal config list
```

## Dashboard Management (`m3tal dash`)

The `m3tal dash` subcommand specifically manages the M3TAL Dashboard container, which runs internally on port 8082 and communicates with the M3TAL API daemon (Go binary) on port 8080.

### 1. `m3tal dashpass [username] [password]`
Updates the password for a specified dashboard user. If `username` and `password` are omitted, it enters an interactive mode to guide you through the process. Dashboard credentials are stored in `/docker/users.json`.

**Usage Examples:**
```bash
# Interactive mode
m3tal dashpass

# Non-interactive mode to set password for 'admin'
m3tal dashpass admin new_secure_password123
```

### 2. `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub, then starts the M3TAL Dashboard container. It automatically selects the correct override based on `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal dash up
```

### 3. `m3tal dash down`
Stops the M3TAL Dashboard container.

**Usage Example:**
```bash
m3tal dash down
```

### 4. `m3tal dash restart`
Restarts the M3TAL Dashboard container.

**Usage Example:**
```bash
m3tal dash restart
```

### 5. `m3tal dash logs`
Streams real-time logs from the M3TAL Dashboard container, useful for debugging and monitoring.

**Usage Example:**
```bash
m3tal dash logs
```

### 6. `m3tal dash status`
Displays the current operational status of the M3TAL Dashboard container.

**Usage Example:**
```bash
m3tal dash status
```

## Stack Management (`m3tal up`, `m3tal down`, `m3tal logs`)

These commands manage all Docker Compose stacks defined in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`).

### 1. `m3tal up`
Executes `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This command starts or recreates all your defined services.

**Usage Example:**
```bash
m3tal up
```

### 2. `m3tal down`
Executes `docker compose down` across all `*-compose.yml` files, stopping and removing all services, networks, and volumes defined in your stacks.

**Usage Example:**
```bash
m3tal down
```

### 3. `m3tal logs`
Streams aggregated, real-time logs from all running Docker containers managed by M3TAL. This provides a centralized view of your entire ecosystem's activity.

**Usage Example:**
```bash
m3tal logs
```

---

## M3TAL System Architecture & Operations Deep Dive

M3TAL leverages Docker Engine and Docker Compose V2 as its foundational runtime. The CLI binary acts as a sophisticated wrapper, simplifying complex Docker operations into user-friendly commands.

### Filesystem Contract

| Path | Purpose |
| :------------------------- | :----------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary configuration file. Managed by `m3tal config wizard`.            |
| `/var/lib/m3tal/state.db`  | SQLite state database for the API daemon. Auto-created.                  |
| `/opt/m3tal/stack/`        | Canonical stack directory. Contains Compose files and Traefik config.    |
| `/docker`                  | **Symlink** to `/opt/m3tal/stack/`. User-facing path for all stack operations. Place your `*-compose.yml` files here. |
| `/docker/users.json`       | Dashboard credential store. Managed by `m3tal dashpass`.                 |

### Dashboard Access Modes

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)
-   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
-   **Mechanism**: Uses the `m3tal-compose.local.yml` override, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`).
-   **Access**: Via `http://HOST_IP:8082` or `http://localhost:8082`.
-   **Requirements**: No Traefik required. Ideal for LAN-only setups, first-time users, or local development.

#### Mode 2: `traefik`
-   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
-   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard container. Traefik (running via `routing-compose.yml`) then routes `dash.${DOMAIN}` to the dashboard on port 8082.
-   **Access**: Via `http://dash.YOUR_DOMAIN` (e.g., `http://dash.example.com`).
-   **Requirements**: Traefik must be running as part of your `m3tal up` stacks. Best for domain-based access and integration into a reverse proxy setup.

### Traefik Routing Architecture

Traefik acts as the central reverse proxy and is deployed via `routing-compose.yml`.
-   It binds to port 80 (HTTP) on the host.
-   Automatically discovers services via Docker labels (e.g., for the dashboard in `traefik` mode).
-   Loads dynamic configuration from `/docker/dynamic/` (e.g., `dynamic/api.yml` for the M3TAL API daemon).
    -   Example: `api.YOUR_DOMAIN` → `http://host.docker.internal:8080` (M3TAL API daemon).

### Port Map

| Port | Service                               | Access                                              |
| :--- | :------------------------------------ | :-------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (if `DASHBOARD_EXPOSE_MODE=traefik`)         |
| 8080 | M3TAL API daemon (Go binary)          | Host-local. Accessible internally by Traefik.       |
| 8081 | Traefik dashboard (admin UI)          | Host-local only.                                    |
| 8082 | M3TAL Dashboard (Python/Flask)        | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

### Systemd Service Management

The M3TAL API daemon, a Go binary, runs as a systemd service named `m3tal-api.service`. This daemon manages Docker interactions, the state database (`/var/lib/m3tal/state.db`), and API routes.

**Check Service Status:**
```bash
systemctl status m3tal-api
```

**Restart the API Daemon:**
```bash
sudo systemctl restart m3tal-api
```

**Stream API Daemon Logs:**
```bash
journalctl -u m3tal-api -f
```

### Direct Docker Compose Fallback

M3TAL's CLI commands (`m3tal up`, `m3tal down`, etc.) wrap standard Docker Compose V2 operations. In scenarios where you need granular control or troubleshooting, you can directly use `docker compose`. All M3TAL-managed compose files reside in `/docker/`.

**To start all services directly with Docker Compose:**
```bash
sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/my-stack-compose.yml up -d
# (Add all your *-compose.yml files as needed)
```

**To stop a specific service (e.g., the dashboard) using its compose file:**
```bash
sudo docker compose -f /docker/m3tal-compose.yml down
```

**To view logs for a specific service using Docker Compose:**
```bash
sudo docker compose -f /docker/m3tal-compose.yml logs -f m3tal-dashboard
```

---

## M3TAL Deployment Lifecycle: Day 2 Operations

**Installing a New Stack:**
1.  **Place your Compose file:** Copy your Docker Compose file (e.g., `my-app-compose.yml`) into the `/docker/` directory.
    ```bash
    sudo cp /path/to/my-app-compose.yml /docker/
    ```
2.  **Configure environment variables:** Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env`. Use the configuration wizard or `m3tal config set`:
    ```bash
    m3tal config wizard
    # OR
    m3tal config set MY_APP_PORT 9000
    ```
3.  **Start all stacks:** Run `m3tal up` to deploy your new stack alongside all existing services.
    ```bash
    m3tal up
    ```

---

## M3TAL APT Installation

To install the M3TAL CLI binary and associated system components, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```