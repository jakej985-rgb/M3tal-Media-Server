Hey there, fellow M3TALhead! DocSmith here, your M3TAL Ecosystem Documentation Architect. You've come to the right place for the ultimate cheat sheet to navigating and mastering the M3TAL CLI. This document covers every command, from initial setup to day-to-day operations and advanced configuration, ensuring you have the power to harness your M3TAL stack effectively.

The M3TAL ecosystem is built on a robust foundation of Docker Engine, Docker Compose V2, and a Go API daemon, all orchestrated through a unified `m3tal` CLI binary. Let's dive in!

---

# M3TAL CLI Command Reference

## I. Introduction to M3TAL

M3TAL provides a unified control plane for deploying and managing self-hosted applications using Docker Compose. It streamlines configuration, service management, and monitoring, abstracting away much of the underlying Docker complexity.

**Core Components:**
- **`m3tal` CLI binary**: Your single entry point for all operations.
- **`m3tal-api.service`**: The systemd-managed Go API daemon (port 8080) that handles Docker interactions and state management.
- **`m3tal-dashboard`**: A web-based UI (Python/Flask container) for visual control, accessible locally or via Traefik.
- **Traefik Gateway**: A containerized reverse proxy for exposing services on port 80/443.
- **Docker Engine + Compose V2**: The underlying containerization technology.

## II. Installation

To get M3TAL up and running on your system, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## III. Filesystem Contract

M3TAL adheres to a strict filesystem contract for configuration and data management:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database, auto-created and managed by the API daemon. |
| `/opt/m3tal/stack/` | Canonical directory for Docker Compose files and Traefik config. |
| `/docker` | **Symlink** → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your `*-compose.yml` files here. |
| `/docker/users.json` | Dashboard credential store, managed by `m3tal dashpass`. |

## IV. Core M3TAL Commands

### `sudo m3tal`

Opens the interactive Text User Interface (TUI) Control Center. This provides a user-friendly, numbered menu interface to common M3TAL operations, ideal for quick status checks and actions without remembering specific CLI flags.

**Usage Example:**
```bash
sudo m3tal
```
_This will launch a full-screen terminal application with options like "1. Dashboard Status", "2. Start All Stacks", "3. Stop All Stacks", etc._

### `m3tal init`

Generates the initial `/etc/m3tal/.env` configuration file from M3TAL's defaults. This is crucial for your first installation or if your `.env` file is missing.

**Usage Example:**
```bash
m3tal init
```
_This command should be run once after installation. If `/etc/m3tal/.env` already exists, it will prompt you before overwriting._

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability, helping diagnose common setup issues.

**Usage Example:**
```bash
m3tal doctor
```
**Expected Output Example:**
```
M3TAL Doctor Pre-Flight Health Check:
✓ Docker Daemon: Running (v25.0.3)
✓ Docker Compose: Installed (v2.24.5)
✓ /etc/m3tal/.env: Valid (27 variables loaded)
✓ Port 8080 (M3TAL API): Available
✓ Port 8082 (Dashboard): Available
✓ Port 80 (Traefik HTTP): Available
✓ Essential Docker Networks: 'proxy' found.
✓ API Daemon: Running
M3TAL system health: GOOD.
```

## V. Configuration Management (`m3tal config`)

M3TAL configuration is primarily managed through the `/etc/m3tal/.env` file. These commands help you interact with it.

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating your `/etc/m3tal/.env` file. It presents each environment variable with its current value and default, allowing easy modifications.

**Usage Example:**
```bash
m3tal config wizard
```
_This will prompt you for each configuration variable, e.g., "DASHBOARD_EXPOSE_MODE (current: local, default: local) [local]: "_

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env` to a specified value. Changes are persistent.

**Usage Example:**
```bash
m3tal config set DASHBOARD_EXPOSE_MODE traefik
```
_This will update the `DASHBOARD_EXPOSE_MODE` to `traefik` in your `/etc/m3tal/.env` file. Remember to `m3tal dash restart` for changes to take effect on the dashboard._

### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal config get DOMAIN
```
**Expected Output Example:**
```
yourdomain.com
```

### `m3tal config scan`

Scans all `*-compose.yml` files in `/docker/` to identify and list all environment variables referenced across all your stacks, including their default values if available. Useful for understanding what variables your entire ecosystem depends on.

**Usage Example:**
```bash
m3tal config scan
```
**Expected Output Example:**
```
Key                 Default             Description
-----------------------------------------------------------------
DASHBOARD_PORT      8082                Port for the M3TAL Dashboard
DASHBOARD_EXPOSE_MODE local               Dashboard exposure mode (local/traefik)
HTTP_PORT           8080                M3TAL API daemon port
DOMAIN              localhost           Your public domain name
PUID                1000                User ID for container processes
PGID                1000                Group ID for container processes
TZ                  America/Denver      Timezone for containers
... (and many more)
```

### `m3tal config list`

Displays the entire contents of your current `/etc/m3tal/.env` file, showing all configured environment variables and their values.

**Usage Example:**
```bash
m3tal config list
```
**Expected Output Example:**
```
# M3TAL Environment Configuration
DASHBOARD_PORT=8082
DASHBOARD_EXPOSE_MODE=traefik
HTTP_PORT=8080
STATE_DIR=/var/lib/m3tal
LOG_LEVEL=info
DASHBOARD_SECRET=your_super_secret_key
API_TOKEN=your_api_token
ADMIN_PASSWORD=your_admin_password
NETWORK_NAME=proxy
LOCAL_IP=192.168.1.100
DOMAIN=yourdomain.com
... (and more)
```

## VI. Dashboard Management (`m3tal dash`)

The M3TAL dashboard (`m3tal-dashboard` container) provides a web-based GUI. These commands manage its lifecycle and access.

### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, the command becomes interactive, prompting you for the necessary details. This updates `/docker/users.json`.

**Usage Examples:**
1. **Interactive Mode:**
   ```bash
   m3tal dashpass
   ```
   _Prompts for username and new password._

2. **Direct Mode:**
   ```bash
   m3tal dashpass admin MyS3cur3P@ssw0rd!
   ```
   _Sets the password for the `admin` user to `MyS3cur3P@ssw0rd!`._

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub, then starts the `m3tal-dashboard` container using the appropriate override file based on `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.

**Dashboard Access Modes (Critical):**

The dashboard has two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **`DASHBOARD_EXPOSE_MODE=local` (Default)**
    *   Uses `m3tal-compose.local.yml` for configuration.
    *   Directly binds `DASHBOARD_PORT` (default: `8082`) from the host to the container.
    *   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082`
    *   No Traefik required. Best for LAN-only setups or first-time use.

2.  **`DASHBOARD_EXPOSE_MODE=traefik`**
    *   Uses `m3tal-compose.traefik.yml` for configuration.
    *   Adds Traefik labels, allowing Traefik to route `dash.${DOMAIN}` to the dashboard on port 8082.
    *   **Access via:** `http://dash.YOUR_DOMAIN` (e.g., `http://dash.myhomelab.com`)
    *   Requires Traefik to be running via `m3tal up`. Best for domain-based access behind a reverse proxy.

**Usage Example:**
```bash
m3tal dash up
```
_This will pull the latest dashboard image and start the `m3tal-dashboard` container, making it accessible based on your `DASHBOARD_EXPOSE_MODE`._

### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container. Useful after changing configuration variables or applying updates.

**Usage Example:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams real-time logs from the `m3tal-dashboard` container. Essential for debugging dashboard issues.

**Usage Example:**
```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` container (e.g., `running`, `exited`, `restarting`).

**Usage Example:**
```bash
m3tal dash status
```
**Expected Output Example:**
```
Container: m3tal-dashboard
Image: ghcr.io/jakej985-rgb/m3tal-godash:debug
Status: Up 10 minutes (healthy)
Ports: 0.0.0.0:8082->8082/tcp
```

## VII. Stack Management (`m3tal up/down/logs`)

These commands manage all your Docker Compose stacks defined by `*-compose.yml` files in the `/docker/` directory.

**Important Note:** The `/docker/` directory is a symlink to `/opt/m3tal/stack/`. When adding new services, place your `my-service-compose.yml` files in `/docker/`.

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This starts or recreates all services defined in your M3TAL ecosystem, including Traefik, Cloudflared, and any user-defined stacks.

**Usage Example:**
```bash
m3tal up
```
_This will start containers for `routing-compose.yml`, `m3tal-compose.yml`, and any other `*-compose.yml` files you have placed in `/docker/`._

### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in the `/docker/` directory. This stops and removes all containers, networks, and volumes associated with your M3TAL stacks.

**Usage Example:**
```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL-managed Docker containers. This provides a unified view of your entire system's activity.

**Usage Example:**
```bash
m3tal logs
```
_This will show logs from `m3tal-dashboard`, `traefik`, `ollama`, and any other running services, prefixed by their container name._

## VIII. Systemd Service Management

The core M3TAL API daemon (`m3tal-api`) runs as a systemd service. These commands are essential for managing and debugging the daemon itself.

### `systemctl status m3tal-api`

Checks the current status of the M3TAL API daemon service.

**Usage Example:**
```bash
systemctl status m3tal-api
```
**Expected Output Example:**
```
● m3tal-api.service - M3TAL API Daemon
     Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
     Active: active (running) since Tue 2024-03-12 10:30:05 UTC; 1min 23s ago
   Main PID: 1234 (m3tal-api)
      Tasks: 7 (limit: 4615)
     Memory: 15.6M
        CPU: 123ms
     CGroup: /system.slice/m3tal-api.service
             └─1234 /usr/bin/m3tal-api

Mar 12 10:30:05 myhost systemd[1]: Started M3TAL API Daemon.
Mar 12 10:30:05 myhost m3tal-api[1234]: [INFO] M3TAL API daemon started on :8080
```

### `journalctl -u m3tal-api -f`

Streams real-time logs from the M3TAL API daemon service. Indispensable for debugging issues related to the API or its interactions with Docker.

**Usage Example:**
```bash
journalctl -u m3tal-api -f
```

### `systemctl restart m3tal-api`

Restarts the M3TAL API daemon. Necessary after manually editing `/etc/m3tal/.env` or if the daemon is misbehaving.

**Usage Example:**
```bash
sudo systemctl restart m3tal-api
```

## IX. Direct Docker / Compose Commands (Fallback & Advanced Use)

M3TAL wraps standard Docker and Docker Compose V2 commands. While M3TAL commands are preferred for consistency and automation, knowing the direct Docker commands is useful for advanced debugging or specific ad-hoc operations.

All M3TAL-managed compose files are located in `/docker/` (which is a symlink to `/opt/m3tal/stack/`).

### Manage a Specific Stack Directly

To manage a single stack, navigate to `/docker/` and use `docker compose` with the specific compose file.

**Example: Start only the routing stack**
```bash
cd /docker/
sudo docker compose -f routing-compose.yml up -d
```

**Example: Stop only the dashboard container**
```bash
cd /docker/
sudo docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml down m3tal-dashboard # or m3tal-compose.traefik.yml
```
_Note: When stopping just one container that uses overrides, you must specify all relevant compose files for Docker Compose to correctly interpret the service definition._

**Example: View logs for a specific container**
```bash
sudo docker logs -f m3tal-dashboard
```

### General Docker Commands

These commands provide insights into your Docker environment.

-   **List all running containers:**
    ```bash
    sudo docker ps
    ```

-   **List all Docker networks:**
    ```bash
    sudo docker network ls
    ```

-   **Inspect a container (e.g., to see its network details or mounted volumes):**
    ```bash
    sudo docker inspect m3tal-dashboard
    ```

-   **Clean up unused Docker resources (images, containers, networks, volumes):**
    ```bash
    sudo docker system prune -a
    ```
    _Use with caution, this removes *all* stopped containers, unused networks, dangling images, and optionally all unused images and volumes._

---

That's the full rundown, M3TAL user! With this cheat sheet, you're equipped to manage your M3TAL ecosystem like a pro. Remember to consult the official M3TAL documentation for deeper dives into specific topics. Happy homelabbing!