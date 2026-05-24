Greetings, M3TAL Operative.

This document serves as your definitive field guide to the M3TAL Command Line Interface (CLI), a comprehensive cheat-sheet crafted to empower your control over the M3TAL Ecosystem. I am DocSmith, your M3TAL Ecosystem Documentation Architect, and I've meticulously laid out every command, subcommand, and critical operational detail to ensure your mission success.

The M3TAL CLI is your single point of entry to manage, configure, and monitor your M3TAL deployment. From initial setup to day-to-day operations, consider this your essential reference.

---

## M3TAL System Architecture: A Deeper Dive

The M3TAL ecosystem is engineered for robustness and flexibility, built upon several interconnected components:

*   **CLI Binary (`/usr/bin/m3tal`)**: The unified Go binary, installed via APT, acts as the primary user interface. It orchestrates all operations, communicating with the API daemon or directly with Docker where appropriate.
*   **API Daemon (`m3tal-api.service`)**: A Go binary running as a systemd service on host port `8080`. This daemon is the brain of M3TAL, managing Docker interactions, maintaining the SQLite state database (`/var/lib/m3tal/state.db`), and exposing a REST API for the CLI and Dashboard.
*   **Dashboard Container (`m3tal-dashboard`)**: A Python/Flask application running within a Docker container. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080` and provides a user-friendly web interface for managing your stacks.
*   **Traefik Gateway (`routing-compose.yml`)**: Our robust reverse proxy solution, running as a Docker container. It exposes services on host port `80` (and `443` if TLS is configured) by domain name, dynamically discovering services via Docker labels and a file provider.
*   **Cloudflared (`routing-compose.yml`)**: An optional Cloudflare tunnel container, offering secure, zero-config public internet access to your services without exposing ports directly on your router.

### Filesystem Contract

Understanding the M3TAL filesystem layout is crucial for effective management:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file.** Holds all environment variables that M3TAL uses. Managed primarily by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the `m3tal-api.service` daemon. |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains M3TAL's core `*-compose.yml` files and Traefik dynamic configuration. |
| `/docker`                   | **Symlink to `/opt/m3tal/stack/`.** This is the user-facing path for placing all your Docker Compose stack definitions. |
| `/docker/users.json`        | Dashboard credential store. Contains usernames and hashed passwords for dashboard access. Managed exclusively by `m3tal dashpass`. |

### Dashboard Access: Two Modes of Operation

The M3TAL Dashboard offers two distinct exposure modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: Uses the `m3tal-compose.local.yml` override file, which adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access**: Navigate to `http://HOST_IP:8082` or `http://localhost:8082` in your web browser.
*   **Requirements**: No Traefik required. Works immediately on a local network or direct server access.
*   **Best For**: Initial setup, LAN-only deployments, testing, or environments without a public domain.

#### Mode 2: `traefik`
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override file. This adds Traefik labels to the dashboard container, instructing Traefik to route `dash.${DOMAIN}` to the dashboard's internal port `8082`.
*   **Access**: Navigate to `http://dash.YOUR_DOMAIN` (e.g., `http://dash.mymetal.com`) in your web browser.
*   **Requirements**: Traefik must be running (via `m3tal up`), and the `DOMAIN` variable in `/etc/m3tal/.env` must be correctly configured.
*   **Best For**: Domain-based deployments, environments with multiple services behind a Traefik reverse proxy, or when using Cloudflared tunnels.

### Docker / Compose Runtime: The Engine Underneath

M3TAL leverages **Docker Engine** and **Docker Compose V2** as its core container orchestration tools. These are hard dependencies for the system.

*   **`m3tal up`**: This command orchestrates the startup of **all** Docker Compose stacks defined by `*-compose.yml` files located in the `/docker/` directory. It effectively runs `docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml ... up -d`.
*   **`m3tal dash up`**: This command specifically manages the `m3tal-dashboard` container. Its process is more involved:
    1.  Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from the official GitHub repository.
    2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
    3.  Starts the dashboard container using the appropriate base compose file and selected override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`).
*   **User Stacks**: To integrate your own Docker Compose applications, simply place your `my-app-compose.yml` file into the `/docker/` directory. M3TAL will automatically pick it up with the next `m3tal up` command.

### Deployment Lifecycle: Day 2 Operations

Integrating a new application stack into your M3TAL ecosystem is straightforward:

1.  **Place Compose File**: Copy your Docker Compose definition (e.g., `my-app-compose.yml`) into the `/docker/` directory.
2.  **Configure Environment**: Ensure any required environment variables for your new stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` for an interactive setup or `m3tal config set KEY value` for specific variables.
3.  **Deploy**: Run `m3tal up` to start your new stack alongside all other M3TAL services.

### Traefik Routing Architecture

Traefik, deployed as a container via `routing-compose.yml`, acts as the intelligent ingress for your M3TAL services:

*   **Entry Points**: Binds to host port `80` (HTTP) and optionally `443` (HTTPS) to receive incoming traffic.
*   **Service Discovery**: Automatically discovers new services by inspecting Docker container labels within the `proxy` network.
*   **Dynamic Configuration**: Loads additional routing rules from `.yml` files in `/docker/dynamic/` (which is mapped to Traefik's internal `/etc/traefik/dynamic` directory). This enables hot-reloading of configuration without restarting Traefik.
*   **M3TAL API Routing**: Routes `api.${DOMAIN}` to the M3TAL API daemon (`http://host.docker.internal:8080`) using a dynamic configuration file (`/docker/dynamic/api.yml`).
*   **Dashboard Routing**: Routes `dash.${DOMAIN}` to the `m3tal-dashboard` container (on its internal port `8082`) when `DASHBOARD_EXPOSE_MODE=traefik`, configured via labels in `m3tal-compose.traefik.yml`.

### Port Map

Key ports utilized by the M3TAL ecosystem:

| Port | Service                    | Access Mode                                                                      |
| :--- | :------------------------- | :------------------------------------------------------------------------------- |
| `80` | Traefik HTTP Entry Point   | Public (if Traefik is running and exposed)                                       |
| `443`| Traefik HTTPS Entry Point  | Public (if Traefik is running and exposed with HTTPS configuration)              |
| `8080`| M3TAL API Daemon (Go)      | Host-local only. Accessed by M3TAL Dashboard and CLI.                            |
| `8081`| Traefik Dashboard          | Host-local only (e.g., `http://127.0.0.1:8081`). For Traefik's own management UI. |
| `8082`| M3TAL Dashboard Container  | Direct port binding (`http://HOST_IP:8082` in `local` mode) or via Traefik (`http://dash.DOMAIN` in `traefik` mode). |

---

## M3TAL CLI Command Reference

This section details every M3TAL CLI command, providing its purpose and a practical usage example.

### 1. `sudo m3tal`

The primary interactive entry point into the M3TAL Control Center.

*   **Description**: Launches the interactive Text-based User Interface (TUI) Control Center. This interface provides a menu-driven way to manage your M3TAL deployment, offering guided operations for common tasks.
*   **Usage Example**:
    ```bash
    sudo m3tal
    ```
    This command will open the TUI, presenting a numbered menu of options, such as starting/stopping services, configuring variables, or checking status. Root privileges are required as it interacts with system-level services and Docker.

### 2. `m3tal init`

Initializes the M3TAL environment.

*   **Description**: Generates the `/etc/m3tal/.env` configuration file from default values. This command is crucial for the very first installation of M3TAL or to reset configuration to defaults. It ensures all necessary environment variables are present before any services are started.
*   **Usage Example**:
    ```bash
    m3tal init
    ```
    If `/etc/m3tal/.env` already exists, this command will prompt for confirmation before overwriting it.

### 3. `m3tal doctor`

Performs a pre-flight health check.

*   **Description**: Runs a diagnostic scan to verify the health and readiness of your M3TAL system. It checks for critical prerequisites like Docker connectivity, the validity of the `/etc/m3tal/.env` file, and ensures essential ports are not already in use by other services.
*   **Usage Example**:
    ```bash
    m3tal doctor
    ```
    This will output a report indicating any issues found, such as Docker not running, syntax errors in `.env`, or port `80` being occupied.

### 4. `m3tal config wizard`

Interactive configuration manager.

*   **Description**: Launches an interactive wizard that guides you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended way to manage your M3TAL environment variables, ensuring proper format and providing helpful descriptions for each setting.
*   **Usage Example**:
    ```bash
    m3tal config wizard
    ```
    The wizard will present each configuration option, allowing you to enter new values or keep existing ones.

### 5. `m3tal config set KEY VALUE`

Sets a specific environment variable.

*   **Description**: Directly sets a single environment variable in the `/etc/m3tal/.env` file to the specified value. This is useful for quick adjustments or scripting.
*   **Usage Example**:
    ```bash
    m3tal config set DOMAIN mymetal.com
    m3tal config set DASHBOARD_EXPOSE_MODE traefik
    ```
    After setting a configuration, you typically need to restart relevant services (e.g., `m3tal dash restart` or `m3tal up`) for changes to take effect.

### 6. `m3tal config get KEY`

Retrieves a specific environment variable's value.

*   **Description**: Displays the current value of a specified environment variable from the `/etc/m3tal/.env` file.
*   **Usage Example**:
    ```bash
    m3tal config get PUID
    ```
    This will output the value associated with the `PUID` key, for example: `1000`.

### 7. `m3tal config scan`

Lists all discoverable environment variables.

*   **Description**: Scans all M3TAL-managed Docker Compose files (`/docker/*.yml`) and the M3TAL API daemon to identify all potential environment variables used across the entire ecosystem. It displays their keys and default values (if available) or current values from `.env`. This provides a comprehensive overview of all configurable parameters.
*   **Usage Example**:
    ```bash
    m3tal config scan
    ```
    This command helps identify what variables a newly added stack might require.

### 8. `m3tal config list`

Lists the current `.env` file contents.

*   **Description**: Prints the entire content of the `/etc/m3tal/.env` file, showing all currently defined environment variables and their values.
*   **Usage Example**:
    ```bash
    m3tal config list
    ```
    This provides a direct look at your active M3TAL configuration.

### 9. `m3tal dashpass [username] [password]`

Manages dashboard user passwords.

*   **Description**: Updates the password for a specified dashboard user. If `username` and `password` are omitted, it launches an interactive prompt to create or update credentials. Dashboard user credentials are stored in `/docker/users.json`.
*   **Usage Examples**:
    *   **Interactive Mode**:
        ```bash
        m3tal dashpass
        ```
        This will prompt you for a username and a new password.
    *   **Direct Mode**:
        ```bash
        m3tal dashpass admin SuperSecureP@ssw0rd!
        ```
        This directly sets the password for the `admin` user. Ensure strong passwords are used.

### 10. `m3tal dash up`

Starts the M3TAL Dashboard container.

*   **Description**: Pulls the latest M3TAL Dashboard Docker Compose configuration from GitHub (`m3tal-compose.yml` and its overrides) and starts the `m3tal-dashboard` container using the configuration specified by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.
*   **Usage Example**:
    ```bash
    m3tal dash up
    ```
    This command ensures your dashboard is running with the most up-to-date image and correct exposure settings.

### 11. `m3tal dash down`

Stops the M3TAL Dashboard container.

*   **Description**: Gracefully stops and removes the `m3tal-dashboard` container.
*   **Usage Example**:
    ```bash
    m3tal dash down
    ```

### 12. `m3tal dash restart`

Restarts the M3TAL Dashboard container.

*   **Description**: Stops, then starts the `m3tal-dashboard` container. This is useful after making changes to the `.env` file that affect the dashboard, such as `DASHBOARD_EXPOSE_MODE`.
*   **Usage Example**:
    ```bash
    m3tal dash restart
    ```

### 13. `m3tal dash logs`

Streams M3TAL Dashboard container logs.

*   **Description**: Displays a real-time stream of logs from the `m3tal-dashboard` container. Essential for debugging issues with the dashboard.
*   **Usage Example**:
    ```bash
    m3tal dash logs
    ```
    Press `Ctrl+C` to stop streaming logs.

### 14. `m3tal dash status`

Shows M3TAL Dashboard container status.

*   **Description**: Reports the current status of the `m3tal-dashboard` container (e.g., `running`, `exited`, `restarting`).
*   **Usage Example**:
    ```bash
    m3tal dash status
    ```
    This will quickly show if your dashboard is active and healthy.

### 15. `m3tal up`

Starts all M3TAL-managed Docker stacks.

*   **Description**: Executes `docker compose up -d` for all `*-compose.yml` files found in the `/docker/` directory. This command brings up your entire M3TAL ecosystem, including Traefik, Cloudflared (if configured), and any custom user stacks.
*   **Usage Example**:
    ```bash
    m3tal up
    ```
    Use this after initial setup or to restart all services after a system reboot or significant configuration changes.

### 16. `m3tal down`

Stops all M3TAL-managed Docker stacks.

*   **Description**: Executes `docker compose down` for all `*-compose.yml` files in `/docker/`, gracefully stopping and removing all associated containers, networks, and volumes (if not explicitly marked as external or persistent).
*   **Usage Example**:
    ```bash
    m3tal down
    ```
    Use this to shut down your entire M3TAL deployment.

### 17. `m3tal logs`

Streams aggregated logs from all running stacks.

*   **Description**: Provides a unified, real-time stream of logs from all Docker containers managed by M3TAL. This is invaluable for monitoring overall system health and debugging interactions between different services.
*   **Usage Example**:
    ```bash
    m3tal logs
    ```
    Press `Ctrl+C` to stop streaming logs.

---

## Systemd Service Management

The M3TAL API daemon, `m3tal-api.service`, runs as a critical background service managed by `systemd`. You can interact with it using standard `systemctl` commands.

*   **Check API Status**:
    ```bash
    systemctl status m3tal-api
    ```
    This command shows if the API daemon is active, its uptime, and recent log entries.

*   **Restart API Daemon**:
    ```bash
    sudo systemctl restart m3tal-api
    ```
    Use this after making any changes that might affect the API daemon's operation, though most configuration is handled via `.env` and `m3tal up` for Docker.

*   **Stream API Daemon Logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```
    This command streams real-time logs from the M3TAL API daemon. Essential for diagnosing issues related to the core API functionality. Press `Ctrl+C` to exit the stream.

---

## Direct Docker Compose Commands (Fallback & Advanced Use)

While the `m3tal` CLI provides a convenient abstraction layer, M3TAL is fundamentally built on Docker Engine and Docker Compose. For advanced users or in troubleshooting scenarios where the `m3tal` CLI might be compromised, you can directly interact with Docker Compose.

Remember that `/docker` is a symlink to `/opt/m3tal/stack/`. The `m3tal` CLI typically operates on all compose files in this directory.

*   **Bring up all stacks directly**:
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/my-app-compose.yml up -d
    ```
    You would need to explicitly list *all* your `*-compose.yml` files.

*   **Bring down all stacks directly**:
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/my-app-compose.yml down
    ```

*   **View logs for a specific service**:
    ```bash
    sudo docker compose -f /docker/m3tal-compose.yml logs -f m3tal-dashboard
    ```
    This example specifically targets the dashboard service's logs.

*   **Check status of all services in a compose file**:
    ```bash
    sudo docker compose -f /docker/routing-compose.yml ps
    ```
    This shows the status of containers defined within `routing-compose.yml`.

**Important**: When using direct `docker compose` commands, ensure you are in the correct directory (e.g., `/docker/`) or explicitly specify all compose files using `-f`. The `m3tal` CLI automates this discovery for you.

---

## M3TAL APT Installation (Always use this exact block)

To install or update the M3TAL CLI and API daemon on Debian/Ubuntu-based systems:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## Appendices

### Docker Services State JSON

```json
[
  {
    "name": "ollama",
    "image": "ollama/ollama:latest",
    "ports": [
      "11434:11434"
    ],
    "stack": "ai"
  },
  {
    "name": "m3tal-dashboard",
    "image": "ghcr.io/jakej985-rgb/m3tal-godash:debug",
    "ports": [],
    "stack": "m3tal"
  },
  {
    "name": "traefik",
    "image": "traefik:latest",
    "ports": [
      "${TRAEFIK_WEB_PORT:-80}:80",
      "${TRAEFIK_WEBHTTPS_PORT:-443}:443",
      "127.0.0.1:8081:8080"
    ],
    "stack": "routing"
  },
  {
    "name": "cloudflared",
    "image": "cloudflare/cloudflared:latest",
    "ports": [],
    "stack": "routing"
  }
]
```

### Environment Variables State JSON

```json
[
  {
    "key": "DASHBOARD_PORT",
    "default": "8082"
  },
  {
    "key": "DASHBOARD_EXPOSE_MODE",
    "default": "local"
  },
  {
    "key": "HTTP_PORT",
    "default": "8080"
  },
  {
    "key": "STATE_DIR",
    "default": "./state"
  },
  {
    "key": "LOG_LEVEL",
    "default": "info"
  },
  {
    "key": "DASHBOARD_SECRET",
    "default": "change_me_immediately"
  },
  {
    "key": "API_TOKEN",
    "default": "change_me_api_token"
  },
  {
    "key": "ADMIN_PASSWORD",
    "default": "admin_pass"
  },
  {
    "key": "NETWORK_NAME",
    "default": "m3tal"
  },
  {
    "key": "LOCAL_IP",
    "default": "127.0.0.1"
  },
  {
    "key": "DOMAIN",
    "default": "localhost"
  },
  {
    "key": "VPN_USER",
    "default": "user"
  },
  {
    "key": "VPN_PASSWORD",
    "default": "password"
  },
  {
    "key": "BASE_STORAGE_PATH",
    "default": "./data"
  },
  {
    "key": "MEDIA_PATH",
    "default": "./data/media"
  },
  {
    "key": "CONFIG_PATH",
    "default": "./data/config"
  },
  {
    "key": "DOWNLOADS_PATH",
    "default": "./data/downloads"
  },
  {
    "key": "PUID",
    "default": "1000"
  },
  {
    "key": "PGID",
    "default": "1000"
  },
  {
    "key": "TZ",
    "default": "America/Denver"
  },
  {
    "key": "TRAEFIK_WEB_PORT",
    "default": "80"
  },
  {
    "key": "TRAEFIK_WEBHTTPS_PORT",
    "default": "443"
  },
  {
    "key": "TRAEFIK_DASHBOARD_PORT",
    "default": "8080"
  },
  {
    "key": "DEBUG_MODE",
    "default": "false"
  },
  {
    "key": "METRICS_ENABLED",
    "default": "true"
  }
]
```