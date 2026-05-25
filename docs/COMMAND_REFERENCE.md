As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the complete CLI cheat-sheet for the M3TAL platform. This document serves as your definitive guide to managing your M3TAL environment, from initial setup to day-to-day operations.

---

# M3TAL CLI Command Reference

M3TAL provides a unified command-line interface (`m3tal`) for managing your entire self-hosted ecosystem. All core operations, from configuration to container orchestration, are accessible via this single entry point.

## Core Concepts & Architecture

Before diving into commands, understand the fundamental building blocks of M3TAL:

*   **CLI binary (`/usr/bin/m3tal`)**: Your primary interaction point. A Go binary handling all user commands.
*   **API daemon (`m3tal-api.service`)**: A Go binary running as a `systemd` service, listening on `http://localhost:8080`. It manages Docker interactions, maintains the SQLite state database (`/var/lib/m3tal/state.db`), and exposes API routes for the CLI and Dashboard.
*   **Dashboard container (`m3tal-dashboard`)**: A Python/Flask application running within Docker. It communicates with the API daemon at `http://host.docker.internal:8080` to provide a user-friendly web interface.
*   **Traefik Gateway (`routing-compose.yml`)**: A Docker container acting as a reverse proxy. It exposes your services (including the M3TAL API and Dashboard in Traefik mode) on port 80, routing requests based on domain names and Docker labels.
*   **Cloudflared**: An optional Cloudflare tunnel container, managed alongside Traefik, for secure, zero-config external access.
*   **Docker Engine + Compose V2**: M3TAL leverages these as its underlying container orchestration technology.

### Filesystem Contract

M3TAL maintains a strict filesystem contract to ensure consistent operation.

| Path                         | Purpose                                                                                                                              |
| :--------------------------- | :----------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`            | **Primary configuration file**. Stores environment variables for the M3TAL API and all Docker Compose stacks. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`    | SQLite state database. Stores persistent M3TAL ecosystem state. Auto-created and managed by the API daemon.                                  |
| `/opt/m3tal/stack/`          | Canonical stack directory. Contains `docker-compose.yml` files (e.g., `routing-compose.yml`, `m3tal-compose.yml`) and Traefik dynamic configuration. |
| `/docker`                    | **User-facing symlink to `/opt/m3tal/stack/`**. This is where you place your custom `*-compose.yml` files for new stacks.               |
| `/docker/users.json`         | Dashboard credential store. Contains usernames and hashed passwords for dashboard access. Managed by `m3tal dashpass`.                   |
| `/docker/dynamic/api.yml`    | Traefik dynamic configuration for routing `api.DOMAIN` to the M3TAL API daemon.                                                         |

### Dashboard Access — TWO MODES

The M3TAL Dashboard can be accessed in two distinct ways, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

**1. Local Mode (Default)**
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: Uses the `m3tal-compose.local.yml` override, which adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the `m3tal-dashboard` container.
*   **Access via**: `http://HOST_IP:8082` or `http://localhost:8082` (where `HOST_IP` is the IP address of your M3TAL server).
*   **Requirements**: No Traefik required. Works immediately after `m3tal dash up`.
*   **Best for**: Local Area Network (LAN) only setups, first-time users, testing, or environments without a domain.

**2. Traefik Mode**
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override. This *removes* the direct port binding and instead adds Traefik labels to the `m3tal-dashboard` container. When Traefik is running, it discovers these labels and routes traffic for `dash.${DOMAIN}` to the dashboard.
*   **Access via**: `http://dash.YOUR_DOMAIN` (e.g., `http://dash.my-metal-server.com`).
*   **Requirements**: `routing-compose.yml` (Traefik) must be running (`m3tal up`). `DOMAIN` must be correctly configured in `/etc/m3tal/.env`. DNS for `dash.YOUR_DOMAIN` must point to your M3TAL server's IP.
*   **Best for**: Domain-based setups, exposing the dashboard alongside other services behind a reverse proxy, and leveraging Traefik's features (like SSL).

### Docker / Compose Runtime

M3TAL acts as a wrapper around `docker compose`.

*   **`m3tal up`**: This command iterates through all `*-compose.yml` files found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`) and runs `docker compose up -d` for them collectively. This includes core M3TAL components like `routing-compose.yml` (Traefik/Cloudflared) and `m3tal-compose.yml` (Dashboard base).
*   **`m3tal dash up`**: This specific command performs additional steps:
    1.  It downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub, ensuring your dashboard setup is current.
    2.  It reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env` to determine which override file (`.local` or `.traefik`) to apply.
    3.  It then starts the `m3tal-dashboard` container using the base and selected override compose files.
*   **Adding New Stacks**: To deploy your own Docker Compose applications, simply place your `my-app-compose.yml` file into `/docker/` and then run `m3tal up`.

### Traefik Routing Architecture

Traefik, deployed by `routing-compose.yml`, acts as the main entry point for domain-based services.
*   It listens on host port 80 (and 443 if SSL is configured).
*   It automatically discovers services by reading Docker labels on running containers (e.g., the dashboard when in `traefik` mode).
*   It loads dynamic routing configuration from `/docker/dynamic/` (which is `/opt/m3tal/stack/dynamic/`).
*   The M3TAL API daemon itself is routed via `api.DOMAIN` using `/docker/dynamic/api.yml`, which points `api.YOUR_DOMAIN` to `http://host.docker.internal:8080`.

### Port Map

| Port | Service                                  | Access                                     |
| :--- | :--------------------------------------- | :----------------------------------------- |
| 80   | Traefik HTTP entry point                 | Public (when Traefik is running)           |
| 8080 | M3TAL API daemon (Go)                    | Host-local only                            |
| 8081 | Traefik dashboard (admin interface)      | Host-local only (`127.0.0.1:8081`)         |
| 8082 | M3TAL Dashboard container (internal)     | Direct port (`HOST_IP:8082`) in local mode |
|      |                                          | Via Traefik (`dash.DOMAIN`) in Traefik mode |

## M3TAL CLI Command Reference

### Installation

To install the M3TAL CLI and API daemon, run the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

### `sudo m3tal`

Opens the interactive Terminal User Interface (TUI) Control Center. This provides a menu-driven way to manage your M3TAL ecosystem.

*   **Description**: Launch the M3TAL TUI for interactive system management.
*   **Usage**: `sudo m3tal`
*   **Example**:
    ```bash
    sudo m3tal
    ```
    (This will open the TUI, presenting a numbered menu of options.)

---

### `m3tal init`

Generates the primary configuration file `/etc/m3tal/.env` from default values. This command should be run immediately after the first M3TAL installation.

*   **Description**: Initialize the M3TAL configuration file.
*   **Usage**: `sudo m3tal init`
*   **Example**:
    ```bash
    sudo m3tal init
    # Output:
    # Initializing /etc/m3tal/.env from defaults...
    # .env file created successfully at /etc/m3tal/.env
    # Please consider running 'm3tal config wizard' to customize your settings.
    ```

---

### `m3tal doctor`

Performs a pre-flight health check on your M3TAL environment, verifying Docker connectivity, `/etc/m3tal/.env` validity, and essential port availability.

*   **Description**: Diagnose common M3TAL system issues.
*   **Usage**: `m3tal doctor`
*   **Example**:
    ```bash
    m3tal doctor
    # Output:
    # M3TAL Doctor - System Health Check
    # ---------------------------------
    # [✓] Docker Daemon: Connected
    # [✓] /etc/m3tal/.env: Valid and readable
    # [✓] Port 8080 (M3TAL API): Available
    # [✓] Port 8082 (Dashboard Local Mode): Available
    # [✓] Port 80 (Traefik): Available
    # [✓] M3TAL API Service: Running (m3tal-api.service)
    #
    # All essential M3TAL components appear healthy.
    ```

---

### `m3tal config` Subcommands

Manage the primary M3TAL configuration file `/etc/m3tal/.env`.

#### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring `/etc/m3tal/.env`. This is the recommended way to set up your environment variables.

*   **Description**: Interactive setup for M3TAL environment variables.
*   **Usage**: `sudo m3tal config wizard`
*   **Example**:
    ```bash
    sudo m3tal config wizard
    # Output:
    # M3TAL Configuration Wizard
    #
    # This wizard will help you set up your /etc/m3tal/.env file.
    # Press Enter to accept the default value or type a new one.
    #
    # DASHBOARD_PORT (Current: 8082, Default: 8082):
    # DASHBOARD_EXPOSE_MODE (Current: local, Default: local) [local/traefik]: traefik
    # DOMAIN (Current: localhost, Default: localhost): my-metal-server.com
    # ... (continues for other variables)
    #
    # Configuration saved to /etc/m3tal/.env.
    # Restart the M3TAL API service and dashboard for changes to take effect.
    ```

#### `m3tal config set KEY VALUE`

Set a single environment variable in `/etc/m3tal/.env`.

*   **Description**: Modify a specific environment variable.
*   **Usage**: `sudo m3tal config set <KEY> <VALUE>`
*   **Example**:
    ```bash
    sudo m3tal config set DASHBOARD_EXPOSE_MODE traefik
    # Output:
    # Set DASHBOARD_EXPOSE_MODE=traefik in /etc/m3tal/.env
    # Restart the M3TAL API service and dashboard for changes to take effect.
    ```

#### `m3tal config get KEY`

Read the value of a single environment variable from `/etc/m3tal/.env`.

*   **Description**: Retrieve the value of an environment variable.
*   **Usage**: `m3tal config get <KEY>`
*   **Example**:
    ```bash
    m3tal config get DOMAIN
    # Output: my-metal-server.com
    ```

#### `m3tal config scan`

Lists all known environment variables and their default values, indicating which stack components utilize them. This helps understand the full range of configurable options.

*   **Description**: Display all environment variables recognized by M3TAL and its stacks.
*   **Usage**: `m3tal config scan`
*   **Example**:
    ```bash
    m3tal config scan
    # Output (truncated):
    # KEY                    DEFAULT             STACKS
    # -----------------------------------------------------------
    # DASHBOARD_PORT         8082                m3tal
    # DASHBOARD_EXPOSE_MODE  local               m3tal
    # DOMAIN                 localhost           m3tal, routing
    # PUID                   1000                m3tal, user_stacks
    # PGID                   1000                m3tal, user_stacks
    # TRAEFIK_WEB_PORT       80                  routing
    # ...
    ```

#### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file.

*   **Description**: Show the active M3TAL configuration.
*   **Usage**: `m3tal config list`
*   **Example**:
    ```bash
    m3tal config list
    # Output (truncated):
    # DASHBOARD_PORT=8082
    # DASHBOARD_EXPOSE_MODE=traefik
    # HTTP_PORT=8080
    # STATE_DIR=/var/lib/m3tal/state
    # ...
    # DOMAIN=my-metal-server.com
    ```

---

### `m3tal dashpass [username] [password]`

Manages user credentials for the M3TAL Dashboard, stored in `/docker/users.json`. If no arguments are provided, it runs interactively.

*   **Description**: Update or set a dashboard user's password.
*   **Usage**:
    *   `m3tal dashpass` (interactive)
    *   `m3tal dashpass <username> <password>` (non-interactive)
*   **Examples**:
    1.  **Interactive**:
        ```bash
        m3tal dashpass
        # Output:
        # Enter username (default: admin): admin
        # Enter new password:
        # Confirm new password:
        # Password for user 'admin' updated successfully.
        ```
    2.  **Non-interactive**:
        ```bash
        m3tal dashpass docsmith "MySecureDashPass#42"
        # Output:
        # Password for user 'docsmith' updated successfully.
        ```

---

### `m3tal dash` Subcommands

Manage the `m3tal-dashboard` container.

#### `m3tal dash up`

Pulls the latest dashboard compose configuration files from GitHub, then starts the `m3tal-dashboard` container with the appropriate `DASHBOARD_EXPOSE_MODE` override.

*   **Description**: Start or update the M3TAL Dashboard container.
*   **Usage**: `m3tal dash up`
*   **Example**:
    ```bash
    m3tal dash up
    # Output:
    # Pulling latest dashboard compose files from GitHub...
    # Reading DASHBOARD_EXPOSE_MODE from /etc/m3tal/.env (Current: traefik)...
    # Starting m3tal-dashboard container...
    # [+] Running 1/0
    #  ⠿ Container m3tal-dashboard  Started
    # Dashboard started. Access via: http://dash.my-metal-server.com
    ```

#### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

*   **Description**: Stop the M3TAL Dashboard container.
*   **Usage**: `m3tal dash down`
*   **Example**:
    ```bash
    m3tal dash down
    # Output:
    # Stopping m3tal-dashboard container...
    # [!] No stopped containers
    # [+] Running 1/0
    #  ⠿ Container m3tal-dashboard  Removed
    # Dashboard stopped.
    ```

#### `m3tal dash restart`

Restarts the `m3tal-dashboard` container.

*   **Description**: Restart the M3TAL Dashboard container.
*   **Usage**: `m3tal dash restart`
*   **Example**:
    ```bash
    m3tal dash restart
    # Output:
    # Restarting m3tal-dashboard container...
    # [+] Running 1/0
    #  ⠿ Container m3tal-dashboard  Restarted
    # Dashboard restarted.
    ```

#### `m3tal dash logs`

Streams logs from the `m3tal-dashboard` container to your terminal. Useful for debugging.

*   **Description**: View real-time logs for the M3TAL Dashboard.
*   **Usage**: `m3tal dash logs`
*   **Example**:
    ```bash
    m3tal dash logs
    # Output (truncated):
    # Attaching to m3tal-dashboard
    # m3tal-dashboard  |  * Serving Flask app 'server'
    # m3tal-dashboard  |  * Debug mode: on
    # m3tal-dashboard  | WARNING: This is a development server. Do not use it in a production deployment. Use a production WSGI server instead.
    # m3tal-dashboard  |  * Running on http://0.0.0.0:8082 (Press CTRL+C to quit)
    # m3tal-dashboard  |  * Restarting with stat
    # m3tal-dashboard  |  * Debugger is active!
    ```

#### `m3tal dash status`

Shows the current status (running, stopped, etc.) of the `m3tal-dashboard` container.

*   **Description**: Check the operational status of the M3TAL Dashboard.
*   **Usage**: `m3tal dash status`
*   **Example**:
    ```bash
    m3tal dash status
    # Output:
    # m3tal-dashboard: Running (healthy)
    ```

---

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in `/docker/` (e.g., `routing-compose.yml`, `m3tal-compose.yml`, `ollama-compose.yml`). This starts or updates all your configured Docker stacks.

*   **Description**: Start all M3TAL-managed Docker Compose stacks in detached mode.
*   **Usage**: `m3tal up`
*   **Example**:
    ```bash
    m3tal up
    # Output:
    # Starting all M3TAL Docker Compose stacks...
    # [+] Running 5/5
    #  ⠿ Network proxy             Created
    #  ⠿ Container m3tal-dashboard Started
    #  ⠿ Container traefik         Started
    #  ⠿ Container cloudflared     Started
    #  ⠿ Container ollama          Started
    # All stacks are up and running.
    ```

---

### `m3tal down`

Runs `docker compose down` across all configured Docker stacks, stopping and removing their containers, networks, and volumes (unless explicitly preserved).

*   **Description**: Stop and remove all M3TAL-managed Docker Compose stacks.
*   **Usage**: `m3tal down`
*   **Example**:
    ```bash
    m3tal down
    # Output:
    # Stopping and removing all M3TAL Docker Compose stacks...
    # [+] Running 5/5
    #  ⠿ Container ollama          Removed
    #  ⠿ Container cloudflared     Removed
    #  ⠿ Container traefik         Removed
    #  ⠿ Container m3tal-dashboard Removed
    #  ⠿ Network m3tal_proxy       Removed
    # All stacks are down.
    ```

---

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL-managed Docker containers.

*   **Description**: View real-time aggregated logs from all active M3TAL containers.
*   **Usage**: `m3tal logs`
*   **Example**:
    ```bash
    m3tal logs
    # Output (truncated):
    # Attaching to m3tal-dashboard, traefik, ollama
    # m3tal-dashboard  |  * Running on http://0.0.0.0:8082 (Press CTRL+C to quit)
    # traefik          | time="2023-10-27T10:30:05Z" level=info msg="Starting provider aggregator.ProviderAggregator"
    # ollama           | time="2023-10-27T10:30:06.123Z" level=info msg="api server listening on [::]:11434"
    # m3tal-dashboard  | 172.18.0.1 - - [27/Oct/2023 10:30:07] "GET /api/v1/status HTTP/1.1" 200 -
    ```

---

## Systemd Service Management

The M3TAL API daemon runs as a `systemd` service named `m3tal-api.service`. You can manage its lifecycle and view its logs directly using `systemctl` and `journalctl`.

*   **Check API daemon status**:
    ```bash
    systemctl status m3tal-api
    # Output (truncated):
    # ● m3tal-api.service - M3TAL API Daemon
    #      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
    #      Active: active (running) since Fri 2023-10-27 10:28:30 UTC; 1min 30s ago
    #    Main PID: 1234 (m3tal-api)
    #       Tasks: 7 (limit: 4627)
    #      Memory: 9.5M
    #         CPU: 18ms
    #      CGroup: /system.slice/m3tal-api.service
    #              └─1234 /usr/bin/m3tal-api
    ```

*   **Restart the API daemon**:
    ```bash
    sudo systemctl restart m3tal-api
    # Output:
    # (No direct output, but service will restart)
    ```

*   **Stream API daemon logs**:
    ```bash
    sudo journalctl -u m3tal-api -f
    # Output (truncated):
    # Oct 27 10:28:30 host systemd[1]: Starting M3TAL API Daemon...
    # Oct 27 10:28:30 host m3tal-api[1234]: [INFO] M3TAL API starting on port 8080...
    # Oct 27 10:28:30 host m3tal-api[1234]: [INFO] Connected to SQLite database at /var/lib/m3tal/state.db
    # Oct 27 10:28:30 host systemd[1]: Started M3TAL API Daemon.
    ```

---

## Direct Docker Commands (Fallback)

While the `m3tal` CLI abstracts away direct Docker Compose interactions, it's useful to know the underlying commands for advanced debugging or manual intervention. All M3TAL Docker Compose files reside in `/docker/` (symlinked to `/opt/m3tal/stack/`).

*   **List all active M3TAL containers**:
    ```bash
    docker ps --filter label=m3tal.stack
    # Output (truncated):
    # CONTAINER ID   IMAGE                                     COMMAND                  CREATED          STATUS          PORTS                                                                       NAMES
    # a1b2c3d4e5f6   ghcr.io/jakej985-rgb/m3tal-godash:debug   "python3 server.py"      5 minutes ago    Up 5 minutes    8082/tcp                                                                    m3tal-dashboard
    # f1e2d3c4b5a6   traefik:latest                            "/entrypoint.sh traefk"  5 minutes ago    Up 5 minutes    0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp, 127.0.0.1:8081->8080/tcp         traefik
    ```

*   **Start a specific M3TAL stack (e.g., routing)**:
    ```bash
    sudo docker compose -f /docker/routing-compose.yml up -d
    # Output:
    # [+] Running 2/2
    #  ⠿ Container traefik         Started
    #  ⠿ Container cloudflared     Started
    ```

*   **Stop a specific M3TAL stack (e.g., dashboard)**:
    ```bash
    sudo docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml down
    # Output:
    # [!] No stopped containers
    # [+] Running 1/0
    #  ⠿ Container m3tal-dashboard  Removed
    ```
    *Note: When `DASHBOARD_EXPOSE_MODE` is `traefik`, replace `m3tal-compose.local.yml` with `m3tal-compose.traefik.yml`.*

*   **View logs for a specific container (e.g., Traefik)**:
    ```bash
    docker logs -f traefik
    # Output (truncated):
    # time="2023-10-27T10:35:10Z" level=info msg="Traefik started and listening on :80"
    # time="2023-10-27T10:35:11Z" level=debug msg="Configuration received from provider docker. Refreshing...
    ```

*   **Execute a command inside a running container (e.g., shell into dashboard)**:
    ```bash
    docker exec -it m3tal-dashboard sh
    # Output:
    # # ls
    # server.py  state  users.json
    # # exit
    ```