# M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the comprehensive command-line interface (CLI) cheat-sheet for the M3TAL ecosystem. This document covers every essential command, its purpose, and real-world usage examples, ensuring you have a complete reference for managing your M3TAL instance.

## M3TAL Ecosystem Overview

M3TAL provides a robust platform for deploying and managing self-hosted services using Docker. It streamlines configuration, service orchestration, and monitoring through a unified CLI and an intuitive web dashboard.

### Core Components
-   **CLI Binary (`/usr/bin/m3tal`)**: Your primary interface for all M3TAL operations, installed via APT.
-   **API Daemon (`m3tal-api.service`)**: A Go binary running as a `systemd` service on port `8080`. It manages Docker interactions, maintains the SQLite state database (`/var/lib/m3tal/state.db`), and exposes an API for the CLI and Dashboard.
-   **Dashboard Container (`m3tal-dashboard`)**: A Python/Flask application providing a web-based control center. It communicates with the API daemon at `http://host.docker.internal:8080`.
-   **Traefik Gateway (`routing-compose.yml`)**: An optional reverse proxy container, binding to port `80` (and `443` for HTTPS if configured). It handles domain-based routing for your services and the M3TAL API/Dashboard.
-   **Cloudflared**: An optional container for secure, zero-trust tunnels, part of the `routing-compose.yml` stack.

### Filesystem Contract

M3TAL adheres to a strict filesystem layout for predictability and ease of management:

| Path                        | Purpose                                                                                                                                                                                           |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | **Primary Configuration File**: Contains all environment variables for M3TAL and its managed services. This file is directly managed by `m3tal config wizard` and `m3tal config set`.                 |
| `/var/lib/m3tal/state.db`   | **SQLite State Database**: Automatically created and managed by the `m3tal-api` daemon to store internal state and service metadata.                                                                  |
| `/opt/m3tal/stack/`         | **Canonical Stack Directory**: The internal location for all Docker Compose files (`*-compose.yml`) and Traefik dynamic configuration files.                                                         |
| `/docker`                   | **User-Facing Stack Directory**: A symlink to `/opt/m3tal/stack/`. This is where you should place your custom `*-compose.yml` files for new services. All `m3tal up` operations scan this directory. |
| `/docker/users.json`        | **Dashboard Credential Store**: Stores hashed user credentials for the M3TAL Dashboard. Managed by `m3tal dashpass`.                                                                                |
| `/docker/dynamic/api.yml`   | **Traefik API Dynamic Config**: Routes `api.${DOMAIN}` to the `m3tal-api` daemon.                                                                                                                   |

### Docker / Compose Runtime

M3TAL is built upon **Docker Engine** and **Docker Compose V2**. These are hard dependencies for the system.
-   The `m3tal up` command orchestrates all `*-compose.yml` files found in the `/docker/` directory.
-   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container, dynamically applying configuration based on `DASHBOARD_EXPOSE_MODE`.
-   User-defined Docker Compose stacks are placed directly in `/docker/`.

### Deployment Lifecycle — Day 2 Operations
1.  **Place your Compose file**: Copy your `my-service-compose.yml` into `/docker/`.
2.  **Configure environment variables**: Use `m3tal config wizard` or `m3tal config set KEY value` to define any necessary environment variables for your new stack (e.g., `PUID`, `PGID`, `MEDIA_PATH`).
3.  **Start all stacks**: Run `m3tal up` to bring up your new service alongside all other M3TAL-managed services.

---

## M3TAL CLI Commands

### 1. Interactive Control Center

#### `sudo m3tal`
Launches the interactive Terminal User Interface (TUI) control center. This provides a guided, menu-driven experience for common M3TAL operations, including managing services, configuration, and monitoring. Requires `sudo` as it interacts with Docker and system-level configurations.

```bash
sudo m3tal
```

*Example Output (truncated):*
```
           _   _ ______ _____  _
          | | | |  ____|  __ \| |
  _ __ ___| |_| | |__  | |__) | |
 | '_ ` _ \ __| |  __| |  ___/| |
 | | | | | | |_| | |____| |    |_|
 |_| |_| |_|\__|_|______|_|    (_)

                 M3TAL: Control Center
---------------------------------------------------
1. Service Management (Start/Stop/Restart all stacks)
2. Dashboard Management (Up/Down/Restart/Logs)
3. Configuration Wizard
4. View Logs
5. Health Check
6. Exit
---------------------------------------------------
Enter your choice:
```

### 2. System Initialization & Health

#### `m3tal init`
Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults. This command is crucial for a first-time installation or if the `.env` file is missing. It populates the file with essential variables required for M3TAL's operation.

```bash
m3tal init
```

*Example Usage:*
```bash
m3tal init
```
*Output:*
```
INFO: Initializing /etc/m3tal/.env with default values.
INFO: /etc/m3tal/.env created successfully. Please run 'm3tal config wizard' to customize your environment.
```

#### `m3tal doctor`
Performs a pre-flight health check of the M3TAL ecosystem. This command verifies critical components such as Docker connectivity, the validity of `/etc/m3tal/.env`, and the availability of required network ports. It's essential for troubleshooting and ensuring a healthy M3TAL environment.

```bash
m3tal doctor
```

*Example Usage:*
```bash
m3tal doctor
```
*Output (example with success):*
```
INFO: Running M3TAL pre-flight health checks...
CHECK: Docker daemon connectivity... OK
CHECK: /etc/m3tal/.env file existence... OK
CHECK: /etc/m3tal/.env file syntax... OK
CHECK: Port 8080 (M3TAL API) availability... OK
CHECK: Port 80 (Traefik HTTP) availability... OK
CHECK: Port 8082 (Dashboard Local) availability... OK
CHECK: Required environment variables (DASHBOARD_SECRET, API_TOKEN, DOMAIN)... OK
INFO: All M3TAL health checks passed. Your system is ready.
```
*Output (example with warning/error):*
```bash
m3tal doctor
```
```
INFO: Running M3TAL pre-flight health checks...
CHECK: Docker daemon connectivity... OK
CHECK: /etc/m3tal/.env file existence... OK
CHECK: /etc/m3tal/.env file syntax... OK
CHECK: Port 8080 (M3TAL API) availability... OK
WARNING: Port 80 (Traefik HTTP) is already in use. This may prevent Traefik from starting.
CHECK: Port 8082 (Dashboard Local) availability... OK
CHECK: Required environment variables (DASHBOARD_SECRET, API_TOKEN, DOMAIN)... OK
WARNING: One or more health checks reported issues. Please review the output above.
```

### 3. Configuration Management

M3TAL's configuration is primarily managed through `/etc/m3tal/.env`.

#### `m3tal config wizard`
Launches an interactive, guided wizard to configure or update the `/etc/m3tal/.env` file. This is the recommended way to set up your M3TAL environment, covering common variables like `DOMAIN`, `DASHBOARD_EXPOSE_MODE`, `PUID`, `PGID`, and storage paths.

```bash
m3tal config wizard
```

*Example Usage:*
```bash
m3tal config wizard
```
*Interactive Prompt (example):*
```
Welcome to the M3TAL Configuration Wizard!

This wizard will guide you through setting up your /etc/m3tal/.env file.

Current value for DOMAIN (e.g., mymetal.com): localhost
Enter new value for DOMAIN [default: localhost]: mymetal.home
Value for DOMAIN set to: mymetal.home

Current value for DASHBOARD_EXPOSE_MODE (local/traefik): local
Enter new value for DASHBOARD_EXPOSE_MODE [local]: traefik
Value for DASHBOARD_EXPOSE_MODE set to: traefik

... (continues with other variables)

Configuration saved to /etc/m3tal/.env.
```

#### `m3tal config set KEY VALUE`
Sets a single environment variable in `/etc/m3tal/.env` to a specified value. This command allows for precise, non-interactive modification of your configuration.

```bash
m3tal config set <KEY> <VALUE>
```

*Example Usage:*
```bash
m3tal config set PUID 1000
m3tal config set DASHBOARD_SECRET "supersecretphrase"
m3tal config set DOMAIN "mymetal.com"
```
*Output:*
```
INFO: Set PUID to 1000 in /etc/m3tal/.env
INFO: Set DASHBOARD_SECRET to supersecretphrase in /etc/m3tal/.env
INFO: Set DOMAIN to mymetal.com in /etc/m3tal/.env
```

#### `m3tal config get KEY`
Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

```bash
m3tal config get <KEY>
```

*Example Usage:*
```bash
m3tal config get PUID
m3tal config get DASHBOARD_EXPOSE_MODE
```
*Output:*
```
PUID=1000
DASHBOARD_EXPOSE_MODE=traefik
```

#### `m3tal config scan`
Lists all known environment variables across all configured stacks and their current values, as well as their defaults. This provides a comprehensive overview of your M3TAL environment.

```bash
m3tal config scan
```

*Example Output (truncated):*
```
Scanning all M3TAL environment variables:

| KEY                    | CURRENT VALUE        | DEFAULT VALUE        | STACK            |
|------------------------|----------------------|----------------------|------------------|
| DASHBOARD_PORT         | 8082                 | 8082                 | m3tal-dashboard  |
| DASHBOARD_EXPOSE_MODE  | traefik              | local                | m3tal-dashboard  |
| DOMAIN                 | mymetal.com          | localhost            | routing, general |
| PUID                   | 1000                 | 1000                 | general          |
| PGID                   | 1000                 | 1000                 | general          |
| API_TOKEN              | change_me_api_token  | change_me_api_token  | m3tal-api        |
| BASE_STORAGE_PATH      | /mnt/data            | ./data               | general          |
| ...                    | ...                  | ...                  | ...              |
```

#### `m3tal config list`
Displays the full contents of the current `/etc/m3tal/.env` file. This is useful for reviewing the entire configuration in plain text.

```bash
m3tal config list
```

*Example Output (truncated):*
```
# M3TAL Environment Configuration
# Generated by M3TAL CLI
# --- Core M3TAL Configuration ---
DASHBOARD_PORT=8082
DASHBOARD_EXPOSE_MODE=traefik
HTTP_PORT=8080
STATE_DIR=/var/lib/m3tal
LOG_LEVEL=info
DASHBOARD_SECRET="supersecretphrase"
API_TOKEN="mysecureapitoken"
ADMIN_PASSWORD="superadminpassword"
NETWORK_NAME=m3tal
LOCAL_IP=192.168.1.100
DOMAIN=mymetal.com
TZ=America/Denver
PUID=1000
PGID=1000
# --- Storage Paths ---
BASE_STORAGE_PATH=/mnt/data
MEDIA_PATH=/mnt/data/media
CONFIG_PATH=/mnt/data/config
DOWNLOADS_PATH=/mnt/data/downloads
# --- Traefik Configuration ---
TRAEFIK_WEB_PORT=80
TRAEFIK_WEBHTTPS_PORT=443
TRAEFIK_DASHBOARD_PORT=8080
# --- Debugging & Metrics ---
DEBUG_MODE=false
METRICS_ENABLED=true
```

### 4. Dashboard Management

The `m3tal dash` subcommand group provides specific control over the M3TAL web dashboard container.

#### `m3tal dashpass [username] [password]`
Updates the password for a dashboard user. If `username` and `password` are omitted, the command becomes interactive, prompting for the username and then the new password. Dashboard user credentials are stored in `/docker/users.json`.

```bash
m3tal dashpass [username] [password]
```

*Example Usage (interactive):*
```bash
m3tal dashpass
```
*Interactive Prompt:*
```
Enter dashboard username (e.g., admin): admin
Enter new password for 'admin': <type_password_here>
Confirm new password: <type_password_again>
INFO: Password for user 'admin' updated successfully.
```

*Example Usage (direct):*
```bash
m3tal dashpass admin "MyStrongSecurePassword123!"
```
*Output:*
```
INFO: Password for user 'admin' updated successfully.
```

#### `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration from GitHub, then starts the `m3tal-dashboard` container. This command intelligently applies the correct Docker Compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

```bash
m3tal dash up
```

*Example Usage:*
```bash
m3tal dash up
```
*Output (example):*
```
INFO: Downloading latest dashboard compose configuration...
INFO: Detected DASHBOARD_EXPOSE_MODE=traefik. Using m3tal-compose.traefik.yml override.
INFO: Starting m3tal-dashboard container...
[+] Running 1/1
 ⠿ Container m3tal-dashboard  Started
INFO: M3TAL Dashboard is now running. Access it via http://dash.mymetal.com
```

**Dashboard Access Modes Explained:**

M3TAL supports two primary ways to access the dashboard, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **`DASHBOARD_EXPOSE_MODE=local` (Default)**
    *   M3TAL uses the `m3tal-compose.local.yml` override.
    *   This adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
    *   **Access via**: `http://YOUR_HOST_IP:8082` or `http://localhost:8082`.
    *   **Best for**: LAN-only setups, initial installation, local development. No Traefik required.

2.  **`DASHBOARD_EXPOSE_MODE=traefik`**
    *   M3TAL uses the `m3tal-compose.traefik.yml` override.
    *   This adds Traefik labels to the dashboard container, allowing Traefik to route requests for `dash.${DOMAIN}` to the dashboard on its internal port `8082`.
    *   **Access via**: `http://dash.YOUR_DOMAIN` (e.g., `http://dash.mymetal.com`).
    *   **Requirements**: Traefik must be running via `m3tal up`.
    *   **Best for**: Domain-based access, environments with a reverse proxy, external access via a domain.

#### `m3tal dash down`
Stops and removes the `m3tal-dashboard` container.

```bash
m3tal dash down
```

*Example Usage:*
```bash
m3tal dash down
```
*Output:*
```
INFO: Stopping m3tal-dashboard container...
[+] Running 1/0
 ⠿ Container m3tal-dashboard  Removed
INFO: M3TAL Dashboard has been stopped.
```

#### `m3tal dash restart`
Restarts the `m3tal-dashboard` container. This is useful for applying configuration changes or resolving minor issues.

```bash
m3tal dash restart
```

*Example Usage:*
```bash
m3tal dash restart
```
*Output:*
```
INFO: Restarting m3tal-dashboard container...
[+] Running 1/1
 ⠿ Container m3tal-dashboard  Restarted
INFO: M3TAL Dashboard has been restarted.
```

#### `m3tal dash logs`
Streams real-time logs from the `m3tal-dashboard` container. This is invaluable for debugging issues related to the dashboard.

```bash
m3tal dash logs
```

*Example Usage:*
```bash
m3tal dash logs
```
*Output (example):*
```
Attaching to m3tal-dashboard
m3tal-dashboard  |  * Serving Flask app 'server' (lazy loading)
m3tal-dashboard  |  * Environment: production
m3tal-dashboard  |    WARNING: This is a development server. Do not use it in a production deployment.
m3tal-dashboard  |    Use a production WSGI server instead.
m3tal-dashboard  |  * Debug mode: off
m3tal-dashboard  |  * Running on http://0.0.0.0:8082 (Press CTRL+C to quit)
m3tal-dashboard  | 172.18.0.1 - - [25/Oct/2023 14:30:05] "GET /api/status HTTP/1.1" 200 -
```

#### `m3tal dash status`
Displays the current running status of the `m3tal-dashboard` container.

```bash
m3tal dash status
```

*Example Usage:*
```bash
m3tal dash status
```
*Output:*
```
INFO: Checking status of m3tal-dashboard...
Name                Command             State         Ports
-------------------------------------------------------------------------
m3tal-dashboard     python3 server.py   Up            8082/tcp
```
*Output (if stopped):*
```
INFO: Checking status of m3tal-dashboard...
Name                Command             State    Ports
--------------------------------------------------------
m3tal-dashboard     python3 server.py   Exited
```

### 5. Global Stack Management

These commands manage all Docker Compose stacks defined by `*-compose.yml` files in the `/docker/` directory.

#### `m3tal up`
Runs `docker compose up -d` across all `*-compose.yml` files found in `/docker/`. This command starts all your M3TAL-managed services and custom user stacks in detached mode. This implicitly includes `routing-compose.yml` (Traefik, Cloudflared) and the dashboard's compose files (if not already up).

```bash
m3tal up
```

*Example Usage:*
```bash
m3tal up
```
*Output (example):*
```
INFO: Running 'docker compose up -d' for all stacks in /docker/...
[+] Running 3/3
 ⠿ Network proxy      Created
 ⠿ Container traefik  Started
 ⠿ Container cloudflared  Started
 ⠿ Container m3tal-dashboard  Started
INFO: All M3TAL stacks are now running in detached mode.
```

#### `m3tal down`
Runs `docker compose down` across all `*-compose.yml` files found in `/docker/`. This command stops and removes all M3TAL-managed service containers, including their networks and volumes (unless explicitly configured otherwise in compose files).

```bash
m3tal down
```

*Example Usage:*
```bash
m3tal down
```
*Output (example):*
```
INFO: Running 'docker compose down' for all stacks in /docker/...
[+] Running 3/3
 ⠿ Container m3tal-dashboard  Removed
 ⠿ Container traefik  Removed
 ⠿ Container cloudflared  Removed
 ⠿ Network proxy      Removed
INFO: All M3TAL stacks have been stopped and removed.
```

#### `m3tal logs`
Streams aggregated logs from all currently running M3TAL-managed Docker containers. This provides a consolidated view of all service activity, making it easier to monitor and troubleshoot the entire ecosystem.

```bash
m3tal logs
```

*Example Usage:*
```bash
m3tal logs
```
*Output (example of aggregated logs):*
```
Attaching to cloudflared, m3tal-dashboard, traefik
traefik          | time="2023-10-25T14:35:01Z" level=info msg="Starting provider aggregator.ProviderAggregator"
traefik          | time="2023-10-25T14:35:01Z" level=info msg="Starting provider command.Provider"
m3tal-dashboard  | 172.18.0.1 - - [25/Oct/2023 14:35:02] "GET /api/status HTTP/1.1" 200 -
cloudflared      | 2023-10-25T14:35:03Z INF Proxying to metrics service on 127.0.0.1:4445
traefik          | time="2023-10-25T14:35:04Z" level=info msg="Configuration received from provider docker.Provider. The docker.Provider has been enabled."
traefik          | time="2023-10-25T14:35:04Z" level=info msg="Configuration received from provider file.Provider. The file.Provider has been enabled."
```

---

## Systemd Service Management

The core M3TAL API daemon (`m3tal-api`) runs as a `systemd` service. You can interact with it using standard `systemctl` and `journalctl` commands.

### `systemctl status m3tal-api`
Checks the current status of the M3TAL API daemon.

```bash
systemctl status m3tal-api
```

*Example Output:*
```
● m3tal-api.service - M3TAL API Daemon
     Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
     Active: active (running) since Wed 2023-10-25 14:30:00 UTC; 5min ago
   Main PID: 1234 (m3tal-api)
      Tasks: 7 (limit: 9275)
     Memory: 15.2M
        CPU: 18ms
     CGroup: /system.slice/m3tal-api.service
             └─1234 /usr/bin/m3tal-api --config /etc/m3tal/.env

Oct 25 14:30:00 mymetal systemd[1]: Started M3TAL API Daemon.
Oct 25 14:30:01 mymetal m3tal-api[1234]: INFO: M3TAL API daemon started on :8080
```

### `systemctl restart m3tal-api`
Restarts the M3TAL API daemon. This is useful after manually editing `/etc/m3tal/.env` or if the API service becomes unresponsive.

```bash
sudo systemctl restart m3tal-api
```

### `journalctl -u m3tal-api -f`
Streams real-time logs from the M3TAL API daemon via `journalctl`. This is the primary way to monitor the API's internal operations and diagnose issues.

```bash
journalctl -u m3tal-api -f
```

*Example Output:*
```
Oct 25 14:38:00 mymetal m3tal-api[1234]: INFO: API service restarting...
Oct 25 14:38:01 mymetal m3tal-api[1234]: INFO: M3TAL API daemon started on :8080
Oct 25 14:38:05 mymetal m3tal-api[1234]: INFO: Request to /api/containers received from 127.0.0.1
Oct 25 14:38:05 mymetal m3tal-api[1234]: DEBUG: Executing docker client command: containers list
```

---

## Direct Docker Compose Commands (Fallback)

While the `m3tal` CLI abstracts many Docker Compose operations, it's useful to know the underlying commands for advanced debugging or direct interaction. M3TAL utilizes **Docker Compose V2**.

All M3TAL's Docker Compose files are located in `/docker/` (which symlinks to `/opt/m3tal/stack/`).

### Bring up all services
This command mimics `m3tal up`.
```bash
docker compose -f /docker/m3tal-compose.yml \
               -f /docker/routing-compose.yml \
               -f /docker/m3tal-compose.local.yml \
               -f /docker/m3tal-compose.traefik.yml \
               -f /docker/my-custom-stack-compose.yml \
               up -d
```
*(Note: You would only include the relevant dashboard override file based on `DASHBOARD_EXPOSE_MODE`.)*

### Bring down all services
This command mimics `m3tal down`.
```bash
docker compose -f /docker/m3tal-compose.yml \
               -f /docker/routing-compose.yml \
               -f /docker/m3tal-compose.local.yml \
               -f /docker/m3tal-compose.traefik.yml \
               -f /docker/my-custom-stack-compose.yml \
               down
```

### Stream logs from all services
This command mimics `m3tal logs`.
```bash
docker compose -f /docker/m3tal-compose.yml \
               -f /docker/routing-compose.yml \
               -f /docker/m3tal-compose.local.yml \
               -f /docker/m3tal-compose.traefik.yml \
               -f /docker/my-custom-stack-compose.yml \
               logs -f
```

### Specific dashboard management
To specifically manage the dashboard container using direct Docker Compose:
```bash
# To start the dashboard in local mode (assuming DASHBOARD_EXPOSE_MODE=local)
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml up -d m3tal-dashboard

# To start the dashboard in traefik mode (assuming DASHBOARD_EXPOSE_MODE=traefik)
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.traefik.yml up -d m3tal-dashboard

# To stop the dashboard
docker compose -f /docker/m3tal-compose.yml down m3tal-dashboard

# To restart the dashboard
docker compose -f /docker/m3tal-compose.yml restart m3tal-dashboard

# To stream dashboard logs
docker compose -f /docker/m3tal-compose.yml logs -f m3tal-dashboard
```
*(Note: For `up` and `restart`, you should include the correct override `-f` file based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.)*

---

## Traefik Routing Architecture

Traefik acts as the central reverse proxy for M3TAL, deployed via `routing-compose.yml`.
-   It listens on port `80` (and `443` for HTTPS if configured).
-   It automatically discovers Docker services based on their labels (e.g., for `m3tal-dashboard` when `DASHBOARD_EXPOSE_MODE=traefik`).
-   It uses a file provider to load dynamic configuration from `/docker/dynamic/` for services without Docker labels, like the `m3tal-api` daemon.
-   **API Routing**: `/docker/dynamic/api.yml` configures Traefik to route `api.${DOMAIN}` (e.g., `api.mymetal.com`) to the `m3tal-api` daemon running on the host at `http://host.docker.internal:8080`.
-   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE=traefik`, labels are added to the `m3tal-dashboard` container to route `dash.${DOMAIN}` (e.g., `dash.mymetal.com`) to the dashboard.

**Port Map**

| Port | Service                    | Access                                        |
| :--- | :------------------------- | :-------------------------------------------- |
| 80   | Traefik HTTP entry point   | Public (used by `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API daemon (Go)      | Host-local (accessed by Dashboard/CLI/Traefik) |
| 8081 | Traefik dashboard (admin UI) | Host-local only (for Traefik's own dashboard) |
| 8082 | M3TAL Dashboard (Python)   | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik |

---

## APT Installation

To install M3TAL and ensure you receive future updates, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```