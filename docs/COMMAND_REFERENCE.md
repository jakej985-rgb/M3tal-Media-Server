```markdown
# M3TAL CLI Command Reference

Welcome, Operator. This document serves as the complete command-line interface (CLI) cheat-sheet for the M3TAL Ecosystem. As the Documentation Architect, I've compiled this guide to empower you with full control over your M3TAL deployment.

The M3TAL CLI, provided by the `/usr/bin/m3tal` binary, is your unified gateway for managing the M3TAL API daemon, Docker containers, configuration, and overall system health.

## M3TAL Ecosystem Overview

M3TAL orchestrates your self-hosted services using Docker Engine and Docker Compose V2.

### Core Components:
*   **CLI binary** (`/usr/bin/m3tal`): Your primary interface.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a `systemd` service (port 8080), handling Docker interactions, state management, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application providing a web UI, communicating with the API daemon.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container (port 80) for domain-based service exposure.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel for secure, zero-config remote access.

### Critical Filesystem Paths:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file.** Managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.                 |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's internal stack management.             |
| `/docker`                   | **User-facing symlink to `/opt/m3tal/stack/`.** Contains `*-compose.yml` files for all running stacks. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |
| `/docker/dynamic/`          | Traefik dynamic configuration files (e.g., `api.yml`).                 |

## M3TAL CLI Installation

To install the M3TAL CLI and API daemon via APT, execute these commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## M3TAL CLI Commands

### Interactive Control Center

#### `sudo m3tal`
Opens the interactive TUI (Text User Interface) Control Center, providing a numbered menu for common operations. Requires `sudo` as it interacts directly with system resources and Docker.

**Purpose**: Access an interactive, menu-driven interface for managing M3TAL.
**Usage Example**:
```bash
sudo m3tal
```

### System Initialization & Health

#### `m3tal init`
Generates the primary configuration file, `/etc/m3tal/.env`, from built-in defaults. This command is crucial for the first installation or to reset the environment file. It will not overwrite an existing `.env` file unless forced.

**Purpose**: Initialize the M3TAL environment file.
**Usage Example**:
```bash
m3tal init
# Output: Creating /etc/m3tal/.env from defaults... Done.
```

#### `m3tal doctor`
Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability, ensuring your system is ready for operation.

**Purpose**: Diagnose common issues and ensure system readiness.
**Usage Example**:
```bash
m3tal doctor
# Output:
# Docker connectivity: OK
# /etc/m3tal/.env validity: OK
# Port 80 (Traefik HTTP): Available
# Port 8080 (M3TAL API): Available
# Port 8082 (M3TAL Dashboard): Available
# All systems nominal.
```

### Configuration Management

M3TAL uses `/etc/m3tal/.env` as its primary configuration store, which is read by the `m3tal-api.service` and various Docker Compose files.

#### `m3tal config wizard`
Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. It prompts for essential environment variables, offering explanations and default values.

**Purpose**: User-friendly configuration of the M3TAL environment.
**Usage Example**:
```bash
m3tal config wizard
# Output:
# Welcome to the M3TAL Configuration Wizard!
# Current DOMAIN (e.g., example.com) [localhost]: mymetal.net
# Current DASHBOARD_EXPOSE_MODE (local/traefik) [local]: traefik
# ... (continues prompting for other variables)
```

#### `m3tal config set KEY VALUE`
Sets a single environment variable in `/etc/m3tal/.env`. This is useful for making quick, specific configuration changes.

**Purpose**: Modify a specific configuration key.
**Usage Example**:
```bash
m3tal config set DOMAIN mymetal.net
# Output: Set DOMAIN=mymetal.net in /etc/m3tal/.env
```
_Note: After changing critical variables, a restart of `m3tal-api.service` and `m3tal up` might be required for changes to take effect in running containers._

#### `m3tal config get KEY`
Retrieves and displays the value of a single environment variable from `/etc/m3tal/.env`.

**Purpose**: Read a specific configuration key's value.
**Usage Example**:
```bash
m3tal config get DASHBOARD_EXPOSE_MODE
# Output: local
```

#### `m3tal config scan`
Lists all environment variables that M3TAL is aware of, including their default values, and indicates if they are set in `/etc/m3tal/.env` or derived from defaults. This command aggregates variables from all known Docker stacks.

**Purpose**: Discover and audit all available configuration variables.
**Usage Example**:
```bash
m3tal config scan
# Output:
# KEY                    VALUE              SOURCE
# ---                    -----              ------
# PUID                   1000               /etc/m3tal/.env
# PGID                   1000               /etc/m3tal/.env
# DOMAIN                 mymetal.net        /etc/m3tal/.env
# DASHBOARD_PORT         8082               Default
# DASHBOARD_EXPOSE_MODE  traefik            /etc/m3tal/.env
# ...
```

#### `m3tal config list`
Displays the current contents of the `/etc/m3tal/.env` file.

**Purpose**: View the active M3TAL environment file.
**Usage Example**:
```bash
m3tal config list
# Output:
# PUID=1000
# PGID=1000
# DOMAIN=mymetal.net
# DASHBOARD_EXPOSE_MODE=traefik
# ...
```

### Dashboard Management

#### `m3tal dashpass [username] [password]`
Manages user credentials for the M3TAL Dashboard.
*   If `username` and `password` are provided, it sets the password directly.
*   If no arguments are provided, it enters an interactive mode to prompt for a username and new password.

**Purpose**: Securely update M3TAL Dashboard user passwords in `/docker/users.json`.
**Usage Example (Interactive)**:
```bash
m3tal dashpass
# Output:
# Enter username for Dashboard: admin
# Enter new password for admin: (input hidden)
# Confirm new password: (input hidden)
# Password for 'admin' updated successfully.
```
**Usage Example (Direct)**:
```bash
m3tal dashpass admin MySuperSecurePa$$word123
# Output: Password for 'admin' updated successfully.
```

#### `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration from GitHub, then starts or updates the `m3tal-dashboard` container according to your `DASHBOARD_EXPOSE_MODE` setting.

**Purpose**: Ensure the M3TAL Dashboard is running with the latest configuration.
**Usage Example**:
```bash
m3tal dash up
# Output:
# Pulling latest m3tal-dashboard compose files from GitHub... Done.
# DASHBOARD_EXPOSE_MODE is 'traefik'. Starting dashboard with traefik routing.
# [+] Running 1/0
#  ⠿ Container m3tal-dashboard  Started
```

#### `m3tal dash down`
Stops and removes the `m3tal-dashboard` container.

**Purpose**: Shut down the M3TAL Dashboard.
**Usage Example**:
```bash
m3tal dash down
# Output:
# [!] No compose files specified for dashboard. Attempting to stop default...
# [+] Running 1/0
#  ⠿ Container m3tal-dashboard  Stopped
```

#### `m3tal dash restart`
Restarts the `m3tal-dashboard` container.

**Purpose**: Apply changes or resolve issues by restarting the dashboard.
**Usage Example**:
```bash
m3tal dash restart
# Output:
# [+] Running 1/1
#  ⠿ Container m3tal-dashboard  Restarted
```

#### `m3tal dash logs`
Streams the logs from the `m3tal-dashboard` container, useful for debugging.

**Purpose**: Monitor dashboard container activity.
**Usage Example**:
```bash
m3tal dash logs
# Output:
# m3tal-dashboard  |  * Serving Flask app 'server'
# m3tal-dashboard  |  * Debug mode: on
# ... (real-time logs)
```

#### `m3tal dash status`
Displays the current running status of the `m3tal-dashboard` container.

**Purpose**: Check dashboard operational status.
**Usage Example**:
```bash
m3tal dash status
# Output:
# Container: m3tal-dashboard
# Status:    running (healthy)
# Image:     ghcr.io/jakej985-rgb/m3tal-godash:debug
# Ports:     8082/tcp (if local mode)
```

### Stack Management

M3TAL manages all Docker Compose stacks by iterating through `*-compose.yml` files in the `/docker/` directory.

#### `m3tal up`
Runs `docker compose up -d` for all `*-compose.yml` files found in `/docker/`. This command starts or recreates all defined services in detached mode.

**Purpose**: Deploy or update all M3TAL-managed Docker stacks.
**Usage Example**:
```bash
m3tal up
# Output:
# Running docker compose up -d for all stacks in /docker/...
# [+] Running 4/0
#  ⠿ Container m3tal-dashboard  Started
#  ⠿ Container traefik          Started
#  ⠿ Container cloudflared      Started
#  ⠿ Container ollama           Started
```

#### `m3tal down`
Runs `docker compose down` for all `*-compose.yml` files found in `/docker/`. This command stops and removes all containers, networks, and volumes defined by your stacks.

**Purpose**: Shut down and remove all M3TAL-managed Docker stacks.
**Usage Example**:
```bash
m3tal down
# Output:
# Running docker compose down for all stacks in /docker/...
# [+] Running 4/0
#  ⠿ Container m3tal-dashboard  Stopped
#  ⠿ Container traefik          Removed
#  ⠿ Container cloudflared      Removed
#  ⠿ Container ollama           Removed
```

#### `m3tal logs`
Streams aggregated logs from all currently running M3TAL Docker containers.

**Purpose**: Centralized monitoring of all M3TAL services.
**Usage Example**:
```bash
m3tal logs
# Output:
# traefik          | time="2023-10-27T10:30:00Z" level=info msg="Configuration loaded successfully"
# m3tal-dashboard  |  * Running on http://0.0.0.0:8082 (Press CTRL+C to quit)
# ollama           | INFO [ollama] Starting server on 0.0.0.0:11434 (version 0.1.2)
# ... (aggregated real-time logs)
```

## Dashboard Access Modes

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. This variable dictates which Docker Compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) is used when the dashboard is started with `m3tal dash up`.

### Mode 1: Local Access (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   Uses `m3tal-compose.local.yml`, which adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access via**: `http://YOUR_HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: No Traefik required. Works out of the box.
*   **Best for**: LAN-only setups, initial setup, local testing, or scenarios without a public domain.

**Configuration Example**:
Set in `/etc/m3tal/.env`:
```ini
DASHBOARD_EXPOSE_MODE=local
DASHBOARD_PORT=8082
```
Then run: `m3tal dash up`
**Access**: Open your web browser to `http://192.168.1.100:8082` (replace with your host's actual IP).

### Mode 2: Traefik-Routed Access

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   Uses `m3tal-compose.traefik.yml`, which adds Traefik labels to the dashboard container, allowing Traefik to route `dash.${DOMAIN}` to port 8082 internally.
*   **Access via**: `http://dash.YOUR_DOMAIN` (e.g., `http://dash.mymetal.net`).
*   **Requirements**: Traefik must be running (`m3tal up` will start it via `routing-compose.yml`), and your `DOMAIN` variable must be correctly configured in `/etc/m3tal/.env`.
*   **Best for**: Domain-based setups, exposing services behind a reverse proxy, and integrating with other Traefik-managed services.

**Configuration Example**:
Set in `/etc/m3tal/.env`:
```ini
DASHBOARD_EXPOSE_MODE=traefik
DOMAIN=mymetal.net
```
Then run: `m3tal up` (this will start Traefik and the dashboard with Traefik labels)
**Access**: Open your web browser to `http://dash.mymetal.net`.

## M3TAL API Daemon Service Management (systemd)

The core M3TAL API daemon runs as a `systemd` service. You can manage its lifecycle using standard `systemctl` commands.

*   **Check API Daemon Status**:
    ```bash
    systemctl status m3tal-api
    # Output:
    # ● m3tal-api.service - M3TAL API Daemon
    #      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
    #      Active: active (running) since ...
    #    Main PID: ...
    #       Tasks: ...
    #      Memory: ...
    #         CPU: ...
    #      CGroup: /system.slice/m3tal-api.service
    #              └─... /usr/bin/m3tal-api
    ```

*   **Restart API Daemon**:
    ```bash
    sudo systemctl restart m3tal-api
    ```

*   **Follow API Daemon Logs**:
    ```bash
    journalctl -u m3tal-api -f
    # Output:
    # Oct 27 10:35:00 myhost m3tal-api[PID]: INFO: API daemon started on :8080
    # Oct 27 10:35:01 myhost m3tal-api[PID]: INFO: Docker client connected
    # ... (real-time logs)
    ```

## Docker Fallback Commands

M3TAL leverages Docker Engine and Docker Compose V2. In advanced scenarios or for debugging, you might need to interact directly with `docker compose`. Remember that M3TAL's stack definition files are located in `/docker/`.

*   **Start all stacks directly with Docker Compose**:
    This is equivalent to `m3tal up`.
    ```bash
    sudo docker compose -f /docker/routing-compose.yml \
                        -f /docker/m3tal-compose.yml \
                        -f /docker/m3tal-compose.local.yml \
                        -f /docker/my-custom-stack-compose.yml \
                        up -d
    ```
    _Note: You'll need to explicitly list all relevant compose files, including override files (`.local.yml`, `.traefik.yml`) based on your `.env` configuration._

*   **Stop a specific stack (e.g., the routing stack)**:
    ```bash
    sudo docker compose -f /docker/routing-compose.yml down
    ```

*   **Stream logs from a specific container (e.g., Traefik)**:
    ```bash
    sudo docker logs -f traefik
    ```

*   **Check status of all containers**:
    ```bash
    sudo docker ps
    ```

## Traefik Routing Architecture

Traefik acts as the reverse proxy for M3TAL, exposing services by domain name on port 80.

*   **Deployment**: Traefik is deployed via `/docker/routing-compose.yml`.
*   **Entry Points**: It binds port 80 on the host for HTTP traffic (`entryPoints.web`).
*   **Service Discovery**: Traefik automatically discovers services within the `proxy` Docker network by inspecting their Docker labels.
*   **Dynamic Configuration**: Additional routing rules can be defined via static YAML files in `/docker/dynamic/`. Traefik monitors this directory for hot-reloads.

### Key Routing Examples:

*   **M3TAL API Daemon**:
    The M3TAL API daemon runs directly on the host (port 8080), not in a Docker container. Traefik routes to it using `host.docker.internal`.
    **Configuration**: Defined in `/docker/dynamic/api.yml`.
    **Access**: `http://api.YOUR_DOMAIN`

    ```yaml
    # /docker/dynamic/api.yml
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

*   **M3TAL Dashboard**:
    When `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard container is configured with Traefik labels in `m3tal-compose.traefik.yml`.
    **Configuration**: Labels on the `m3tal-dashboard` service.
    **Access**: `http://dash.YOUR_DOMAIN`

    ```yaml
    # Excerpt from /docker/m3tal-compose.traefik.yml
    services:
      m3tal-dashboard:
        labels:
          - "traefik.enable=true"
          - "traefik.http.routers.dashboard.rule=Host(`dash.${DOMAIN:-localhost}`)"
          - "traefik.http.routers.dashboard.entrypoints=web"
          - "traefik.http.services.dashboard.loadbalancer.server.port=8082"
          - "traefik.docker.network=proxy"
    ```

This concludes the M3TAL CLI Command Reference. May your M3TAL ecosystem operate with precision and power!
```