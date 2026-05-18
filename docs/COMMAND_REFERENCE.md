# docs/COMMAND_REFERENCE.md

## M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, my purpose is to provide you with the most precise and actionable command-line interface (CLI) reference for managing your M3TAL deployment. This guide covers every essential command, offering real-world usage examples and detailing the underlying architecture.

M3TAL provides a unified CLI (`/usr/bin/m3tal`) that acts as your central control point for system configuration, service management, and overall health monitoring.

---

### Core CLI Commands

#### `sudo m3tal`
Opens the interactive TUI (Terminal User Interface) Control Center. This provides a numbered menu for common operations, simplifying complex workflows. Requires `sudo` as it often interacts with system-level services and Docker.

```bash
sudo m3tal
```

**Usage Example:**
```bash
sudo m3tal
# You will see a menu like:
# M3TAL Control Center
# 1. Start All Stacks
# 2. Stop All Stacks
# 3. View All Logs
# 4. Configure M3TAL
# 5. Check System Health
# ... (and many more options)
# Enter your choice:
```

#### `m3tal init`
Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults. This command is crucial for the first installation or when you need to regenerate the default `.env` file. It will not overwrite an existing file without confirmation.

```bash
m3tal init
```

**Usage Example:**
```bash
m3tal init
# Output:
# No /etc/m3tal/.env found. Generating from defaults...
# Successfully created /etc/m3tal/.env with default values.
# Please consider running 'm3tal config wizard' to customize your setup.
```

#### `m3tal doctor`
Performs a comprehensive pre-flight health check of your M3TAL system. This includes verifying Docker connectivity, validating the `/etc/m3tal/.env` file, and checking for port availability (e.g., 80, 8080, 8082). Essential for troubleshooting or before deploying new services.

```bash
m3tal doctor
```

**Usage Example:**
```bash
m3tal doctor
# Output:
# M3TAL Doctor: Pre-flight Health Check
#
# Docker Connectivity:
#   - Docker daemon running: [OK]
#   - Docker Compose V2 detected: [OK]
#
# Configuration (/etc/m3tal/.env):
#   - File exists: [OK]
#   - Required variables set: [OK] (DOMAIN, API_TOKEN, DASHBOARD_SECRET)
#   - Valid DOMAIN 'localhost' detected: [WARNING - Consider changing for public access]
#
# Port Availability:
#   - Port 80 (HTTP Traefik): [OK]
#   - Port 8080 (M3TAL API): [OK]
#   - Port 8082 (M3TAL Dashboard): [OK]
#   - Port 8081 (Traefik Dashboard): [OK]
#
# Systemd Service:
#   - m3tal-api.service status: active (running) [OK]
#
# All checks passed. Your M3TAL system is ready for operation.
```

---

### Configuration Management (`m3tal config`)

M3TAL uses `/etc/m3tal/.env` as its primary configuration file. All `m3tal config` commands interact with this file.

#### `m3tal config wizard`
Launches an interactive, step-by-step wizard to guide you through configuring or updating your `/etc/m3tal/.env` file. This is the recommended way to manage your system's configuration.

```bash
m3tal config wizard
```

**Usage Example:**
```bash
m3tal config wizard
# Output:
# M3TAL Configuration Wizard
#
# Welcome! Let's set up your M3TAL environment.
# Current DOMAIN: localhost (default)
# Enter new DOMAIN (e.g., mym3tal.com): example.com
#
# Current DASHBOARD_PORT: 8082
# Enter new DASHBOARD_PORT (leave blank for default):
#
# ... (continues for other critical variables)
# Configuration updated in /etc/m3tal/.env. Please restart services if changes affect running containers.
```

#### `m3tal config set KEY VALUE`
Sets a single environment variable within `/etc/m3tal/.env`. This command directly modifies the configuration file.

```bash
m3tal config set DOMAIN mydomain.com
m3tal config set LOG_LEVEL debug
```

**Usage Example:**
```bash
m3tal config set DOMAIN mym3tal.net
# Output:
# Successfully set DOMAIN=mym3tal.net in /etc/m3tal/.env
```

#### `m3tal config get KEY`
Retrieves and displays the value of a specific environment variable from `/etc/m3tal/.env`.

```bash
m3tal config get DOMAIN
m3tal config get API_TOKEN
```

**Usage Example:**
```bash
m3tal config get DOMAIN
# Output:
# mym3tal.net
```

#### `m3tal config scan`
Scans all Docker Compose files in `/docker/` and lists all environment variables referenced across all stacks, along with their default values where applicable. This helps identify all potential configuration points for your entire M3TAL ecosystem.

```bash
m3tal config scan
```

**Usage Example:**
```bash
m3tal config scan
# Output:
# Discovered Environment Variables across all M3TAL Stacks:
#
# KEY                 DEFAULT VALUE       SOURCE
# ------------------- ------------------- ---------------------------------
# DOMAIN              localhost           /docker/routing-compose.yml
# API_TOKEN           change_me_api_token /docker/m3tal-compose.yml
# DASHBOARD_SECRET    change_me_immediately /docker/m3tal-compose.yml
# DASHBOARD_PORT      8082                /docker/m3tal-compose.yml
# TRAEFIK_WEB_PORT    80                  /docker/routing-compose.yml
# PUID                1000                /docker/some-user-stack.yml
# PGID                1000                /docker/some-user-stack.yml
# TZ                  America/Denver      /docker/some-user-stack.yml
# ...
```

#### `m3tal config list`
Displays the full contents of the current `/etc/m3tal/.env` file, showing all configured environment variables and their values.

```bash
m3tal config list
```

**Usage Example:**
```bash
m3tal config list
# Output:
# # M3TAL System Configuration
# DOMAIN=mym3tal.net
# API_TOKEN=your_secure_api_token
# DASHBOARD_SECRET=your_secure_dashboard_secret
# DASHBOARD_PORT=8082
# HTTP_PORT=8080
# STATE_DIR=/var/lib/m3tal/state
# LOG_LEVEL=info
# NETWORK_NAME=m3tal
# ...
```

---

### M3TAL Dashboard Management (`m3tal dash`)

The M3TAL Dashboard (`m3tal-dashboard` container) provides a web-based UI for managing your system. These commands specifically control its lifecycle and access.

#### `m3tal dashpass [username] [password]`
Updates or creates a user password for the M3TAL Dashboard. If `username` and `password` are omitted, the command becomes interactive, prompting you for details. Credentials are stored in `/docker/users.json`.

```bash
m3tal dashpass admin newSecurePassword123
m3tal dashpass
```

**Usage Example (with args):**
```bash
m3tal dashpass admin SuperSecureP@ssw0rd!
# Output:
# Successfully updated password for user 'admin' in /docker/users.json.
```

**Usage Example (interactive):**
```bash
m3tal dashpass
# Output:
# Enter username: myuser
# Enter new password:
# Confirm new password:
# Password for 'myuser' updated.
```

#### `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration from GitHub, then starts or updates the `m3tal-dashboard` container. This specifically targets the dashboard service defined in `/docker/m3tal-compose.yml`. This command also ensures the `m3tal-api.service` is running, as the dashboard relies on it.

```bash
m3tal dash up
```

**Usage Example:**
```bash
m3tal dash up
# Output:
# Pulling latest m3tal-compose.yml...
# Building/starting m3tal-dashboard container...
# [m3tal-dashboard] Pulling image ghcr.io/jakej985-rgb/m3tal-godash:debug
# [m3tal-dashboard] Creating m3tal-dashboard ... done
# M3TAL Dashboard is now running. Access at dash.<YOUR_DOMAIN> or http://localhost:8082
```

#### `m3tal dash down`
Stops and removes the `m3tal-dashboard` container.

```bash
m3tal dash down
```

**Usage Example:**
```bash
m3tal dash down
# Output:
# Stopping m3tal-dashboard container...
# [m3tal-dashboard] Stopping m3tal-dashboard ... done
# [m3tal-dashboard] Removing m3tal-dashboard ... done
# M3TAL Dashboard has been stopped and removed.
```

#### `m3tal dash restart`
Restarts the `m3tal-dashboard` container. This is useful after making configuration changes that affect the dashboard.

```bash
m3tal dash restart
```

**Usage Example:**
```bash
m3tal dash restart
# Output:
# Restarting m3tal-dashboard container...
# [m3tal-dashboard] Restarting m3tal-dashboard ... done
# M3TAL Dashboard has been restarted.
```

#### `m3tal dash logs`
Streams the real-time logs from the `m3tal-dashboard` container. Press `Ctrl+C` to exit.

```bash
m3tal dash logs
```

**Usage Example:**
```bash
m3tal dash logs
# Output:
# Attaching to m3tal-dashboard
# m3tal-dashboard  | 2023-10-27 10:30:01 INFO: Starting M3TAL Dashboard on 0.0.0.0:8082
# m3tal-dashboard  | 2023-10-27 10:30:05 INFO: Connected to M3TAL API at http://host.docker.internal:8080
# m3tal-dashboard  | 2023-10-27 10:30:08 INFO: User 'admin' logged in from 172.17.0.1
```

#### `m3tal dash status`
Displays the current status of the `m3tal-dashboard` container (e.g., running, stopped, exited).

```bash
m3tal dash status
```

**Usage Example:**
```bash
m3tal dash status
# Output:
# m3tal-dashboard: Running (Up 5 minutes)
```

---

### M3TAL Stack Management (`m3tal up`, `m3tal down`, `m3tal logs`)

These commands manage all Docker Compose stacks found within the `/docker/` directory. M3TAL automatically discovers and orchestrates all `*-compose.yml` files.

#### `m3tal up`
Runs `docker compose up -d` across *all* `*-compose.yml` files located in the `/docker/` directory. This command orchestrates your entire M3TAL ecosystem, ensuring all services are started and running in detached mode. It also ensures the shared `/etc/m3tal/.env` file is used by all compose files.

```bash
m3tal up
```

**Usage Example:**
```bash
m3tal up
# Output:
# Running docker compose up -d for all stacks in /docker/...
# [routing] Network m3tal_proxy created
# [routing] Creating traefik ... done
# [routing] Creating cloudflared ... done
# [m3tal] Creating m3tal-dashboard ... done
# [my-stack] Creating my-app ... done
# All M3TAL stacks are now running.
```

#### `m3tal down`
Runs `docker compose down` across *all* `*-compose.yml` files in the `/docker/` directory. This stops and removes all containers, networks, and volumes associated with your M3TAL stacks.

```bash
m3tal down
```

**Usage Example:**
```bash
m3tal down
# Output:
# Running docker compose down for all stacks in /docker/...
# [routing] Stopping traefik ... done
# [routing] Stopping cloudflared ... done
# [routing] Removing traefik ... done
# [routing] Removing cloudflared ... done
# [m3tal] Stopping m3tal-dashboard ... done
# [m3tal] Removing m3tal-dashboard ... done
# [my-stack] Stopping my-app ... done
# [my-stack] Removing my-app ... done
# All M3TAL stacks have been stopped and removed.
```

#### `m3tal logs`
Streams aggregated real-time logs from *all* currently running Docker containers managed by M3TAL. Press `Ctrl+C` to exit.

```bash
m3tal logs
```

**Usage Example:**
```bash
m3tal logs
# Output:
# Attaching to traefik, m3tal-dashboard, cloudflared, my-app
# traefik          | 2023-10-27T10:45:01Z INFO message="Starting Traefik..."
# m3tal-dashboard  | 2023-10-27 10:45:02 INFO: Connected to M3TAL API
# my-app           | 2023-10-27 10:45:03 INFO: My application started successfully
# cloudflared      | 2023-10-27T10:45:04Z INFO Tunnel metrics server listening on 127.0.0.1:2006
# ...
```

---

### Key Filesystem Paths

Understanding the M3TAL filesystem contract is crucial for maintenance and deployment:

| Path                        | Purpose                                                                                                 |
|-----------------------------|---------------------------------------------------------------------------------------------------------|
| `/etc/m3tal/.env`           | **Primary Configuration File.** All system-wide environment variables are stored here. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | **SQLite State Database.** Auto-created by the M3TAL API daemon. Stores internal state, user data, etc.  |
| `/opt/m3tal/stack/`         | **Canonical Stack Directory.** Contains core M3TAL compose files and Traefik configuration.               |
| `/docker`                   | **User-Facing Stack Directory.** This is a symlink to `/opt/m3tal/stack/`. All user-defined `*-compose.yml` files for custom applications should be placed here. |
| `/docker/users.json`        | **Dashboard Credential Store.** Stores encrypted dashboard user passwords. Managed by `m3tal dashpass`.   |
| `/docker/dynamic/`          | **Traefik Dynamic Configuration.** Traefik hot-reloads `.yml` files in this directory for dynamic routing rules. |

---

### M3TAL Architecture Overview

The M3TAL ecosystem is built for robustness and extensibility:

#### Core Components
-   **CLI binary** (`/usr/bin/m3tal`): The unified Go binary, installed via APT, serves as the single entrypoint for all operations.
-   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on `http://localhost:8080`. It manages Docker, interacts with `/var/lib/m3tal/state.db`, and exposes API routes for the dashboard and other clients.
-   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running on port 8082. It communicates with the API daemon at `http://host.docker.internal:8080` (from inside the container).
-   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that serves as the entry point for all HTTP traffic on port 80. It uses file providers for dynamic routing rules.
-   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare Tunnel container that provides secure, zero-config internet access to your services without opening firewall ports.

#### Docker / Docker Compose V2 Runtime
M3TAL leverages **Docker Engine** and **Docker Compose V2** as its container orchestration backbone. These are hard dependencies for the system to function.
-   The `m3tal up` command orchestrates *all* Docker Compose files (`*-compose.yml`) found in the `/docker/` directory.
-   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container, defined in `/docker/m3tal-compose.yml`.
-   **User-defined stacks:** To add your own applications, simply place your `*-compose.yml` file into `/docker/`.
-   All Docker Compose commands invoked by M3TAL automatically use the shared configuration file at `/etc/m3tal/.env` via the `--env-file` flag, ensuring a consistent environment.

#### Traefik Routing Architecture (The Gateway)
Traefik is deployed as a container via `routing-compose.yml` and acts as the intelligent ingress for your M3TAL services:
-   It binds to port 80 on the host, serving as the primary HTTP entry point.
-   Services are automatically discovered and routed via Docker labels (e.g., `traefik.http.routers.myapp.rule=Host(\`app.domain.com\`)` applied to your container definitions).
-   **Dynamic Configuration:** Traefik loads additional routing rules from `.yml` files in `/docker/dynamic/` (using a file provider), which are hot-reloaded without restarting Traefik.
-   **Core Routes:**
    -   `api.DOMAIN` routes to the M3TAL API daemon on `http://host.docker.internal:8080` (configured via `/docker/dynamic/api.yml`).
    -   `dash.DOMAIN` routes to the M3TAL dashboard container on port 8082 (configured via labels in `/docker/m3tal-compose.traefik.yml`).
-   The Traefik dashboard itself is accessible locally at `http://localhost:8081` (port 8080 inside the container).

**Traefik Static Configuration (`traefik.yml` - internal to container):**
```yaml
entryPoints:
  web:
    address: ":80" # Host port 80 mapped to this entrypoint

providers:
  docker:
    exposedByDefault: false # Only containers with Traefik labels are exposed
    network: proxy          # Connects to the 'm3tal_proxy' network
  file:
    directory: /etc/traefik/dynamic # This maps to /docker/dynamic on the host
    watch: true # Hot-reload dynamic config changes
```

**Dynamic API Routing Example (`/docker/dynamic/api.yml`):**
```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)" # Routes traffic for api.yourdomain.com
      service: api
      entryPoints:
        - web

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080" # Points to the M3TAL API daemon on the host
```

---

### Deployment Lifecycle: Day 2 Operations

Installing a new Docker Compose stack into your M3TAL ecosystem is straightforward:

1.  **Place your compose file:** Put your Docker Compose definition (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
2.  **Configure environment variables:** Ensure any required variables for your new stack are set in `/etc/m3tal/.env`. Use `m3tal config wizard` for an interactive setup, or `m3tal config set KEY value` for specific variables.
3.  **Start your stack(s):**
    -   Run `m3tal up` to start all stacks, including your newly added one.
    -   Alternatively, for a single stack, use direct Docker Compose: `docker compose -f /docker/my-stack-compose.yml --env-file /etc/m3tal/.env up -d`.

---

### Service Management (systemd)

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. While CLI commands like `m3tal dash up` might interact with it, direct systemd commands provide granular control and monitoring.

-   **Check API daemon status:**
    ```bash
    systemctl status m3tal-api
    ```
    **Usage Example:**
    ```bash
    systemctl status m3tal-api
    # Output:
    # ● m3tal-api.service - M3TAL API Daemon
    #      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
    #      Active: active (running) since Fri 2023-10-27 10:00:00 UTC; 1h 20min ago
    #    Main PID: 1234 (m3tal-api)
    #       Tasks: 8 (limit: 9278)
    #      Memory: 25.6M
    #         CPU: 1.234s
    #      CGroup: /system.slice/m3tal-api.service
    #              └─1234 /usr/bin/m3tal-api
    ```

-   **Restart the API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```

-   **View real-time logs for the API daemon:**
    ```bash
    journalctl -u m3tal-api -f
    ```
    **Usage Example:**
    ```bash
    journalctl -u m3tal-api -f
    # Output:
    # Oct 27 10:00:00 hostname systemd[1]: Started M3TAL API Daemon.
    # Oct 27 10:00:01 hostname m3tal-api[1234]: INFO: M3TAL API starting on port 8080
    # Oct 27 10:00:05 hostname m3tal-api[1234]: INFO: SQLite database opened at /var/lib/m3tal/state.db
    ```

---

### Direct Docker / Docker Compose Fallbacks

While the `m3tal` CLI provides a convenient abstraction, you can always interact directly with Docker and Docker Compose. This is useful for advanced debugging or specific container operations. Remember to include `--env-file /etc/m3tal/.env` when using `docker compose` directly to ensure your configurations are applied.

-   **List all running containers:**
    ```bash
    docker ps
    ```

-   **Stop and remove a specific stack (e.g., `routing`):**
    ```bash
    docker compose -f /docker/routing-compose.yml --env-file /etc/m3tal/.env down
    ```

-   **Start a specific stack in detached mode (e.g., `my-app`):**
    ```bash
    docker compose -f /docker/my-app-compose.yml --env-file /etc/m3tal/.env up -d
    ```

-   **View logs for a specific container (e.g., `traefik`):**
    ```bash
    docker logs -f traefik
    ```

-   **Execute a command inside a running container:**
    ```bash
    docker exec -it m3tal-dashboard bash
    ```

---

### Port Map

| Port | Service                            | Access                                      |
|------|------------------------------------|---------------------------------------------|
| 80   | Traefik HTTP entry point           | Public (via `routing-compose.yml`)          |
| 8080 | M3TAL API daemon (Go)              | Host-local (via Traefik or direct host IP)  |
| 8081 | Traefik dashboard                  | Host-local only                             |
| 8082 | M3TAL Dashboard (Python/Flask)     | Via Traefik (dash.DOMAIN) or direct host IP |

---

### APT Installation

To install or update the `m3tal` CLI and its associated components:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```