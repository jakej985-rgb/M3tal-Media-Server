# M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, I present this comprehensive guide to navigating your M3TAL CLI. This document serves as your cheat-sheet for deploying, managing, and troubleshooting your M3TAL server and its associated stacks.

M3TAL provides a unified command-line interface (`m3tal`) for interacting with its powerful Go API daemon, simplifying the management of Docker containers and configuration.

---

## 🚀 M3TAL Quick Start

### Installation

To get M3TAL up and running on your system, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### First-time Setup

After installation, initialize M3TAL and configure your environment:

1.  **Initialize Configuration**: `m3tal init`
2.  **Configure Environment**: `m3tal config wizard`
3.  **Start Dashboard**: `m3tal dash up`
4.  **Access Dashboard**: See [Dashboard Access Modes](#dashboard-access-modes) below.
5.  **Start All Stacks**: `m3tal up`

---

## 📚 Core M3TAL Concepts

Understanding these core architectural components is key to effectively managing your M3TAL ecosystem.

### M3TAL System Architecture

*   **CLI binary** (`/usr/bin/m3tal`): The single entry point for all operations, installed via APT.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on port `8080`. It orchestrates Docker operations, manages the internal state database, and exposes API routes for the CLI and Dashboard.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container, typically on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): The reverse proxy container, listening on port `80` (and `443` if configured). It dynamically routes traffic to your services based on domain names and Docker labels.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing secure, zero-config internet access to your services without exposing ports directly.

### Filesystem Contract

The following paths are critical to M3TAL's operation and configuration:

| Path                        | Purpose                                                                                                 |
| :-------------------------- | :------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | **Primary Configuration File**. Stores all environment variables for M3TAL and your Docker stacks. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | **SQLite State Database**. Auto-created and managed by the API daemon to persist system state.            |
| `/opt/m3tal/stack/`         | **Canonical Stack Directory**. Contains core M3TAL compose files (`m3tal-compose.yml`, `routing-compose.yml`) and Traefik configuration. **Do not modify directly.** |
| `/docker`                   | **User-facing Stack Directory**. A symlink to `/opt/m3tal/stack/`. This is where you place your custom `*-compose.yml` files for new services. |
| `/docker/users.json`        | **Dashboard Credential Store**. Stores dashboard user hashes. Managed by `m3tal dashpass`.               |

### Docker / Compose Runtime

M3TAL leverages **Docker Engine** and **Docker Compose V2** for all container orchestration. These are hard dependencies.

*   The `m3tal up` command orchestrates all `*-compose.yml` files found in the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container, downloading the latest compose configurations and applying the correct override based on your `DASHBOARD_EXPOSE_MODE`.
*   To add new services, simply place your Docker Compose files (e.g., `my-service-compose.yml`) into the `/docker/` directory and run `m3tal up`.

### Dashboard Access Modes

The M3TAL Dashboard can be accessed in two primary ways, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: Uses `m3tal-compose.local.yml` to add a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access**: Navigate to `http://HOST_IP:8082` or `http://localhost:8082` in your web browser.
*   **Requirements**: No Traefik required. Ideal for LAN-only setups, first-time users, or local development.

#### Mode 2: `traefik`

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: Uses `m3tal-compose.traefik.yml` to apply Traefik labels, routing `dash.${DOMAIN}` to the dashboard container on port `8082`. **Traefik must be running** via `m3tal up`.
*   **Access**: Navigate to `http://dash.yourdomain.com` (replace `yourdomain.com` with your configured `DOMAIN` in `/etc/m3tal/.env`).
*   **Requirements**: Traefik (from `routing-compose.yml`) must be running and configured correctly. Best for domain-based access behind a reverse proxy.

### Traefik Routing Architecture

Traefik, deployed via `routing-compose.yml`, acts as the central ingress point:

*   It binds to port `80` (and `443` if HTTPS is enabled) on the host.
*   It automatically discovers services by reading Docker labels on running containers.
*   It loads dynamic routing configurations from `/docker/dynamic/` (e.g., `api.yml` for the M3TAL API daemon), allowing hot-reloads without restarting Traefik.
*   It routes `api.${DOMAIN}` to `http://host.docker.internal:8080` (the M3TAL API daemon).
*   If `DASHBOARD_EXPOSE_MODE=traefik`, it routes `dash.${DOMAIN}` to the `m3tal-dashboard` container.

### Port Map

| Port | Service | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| `80` | Traefik HTTP entry point    | Public (if `DOMAIN` is configured)          |
| `8080` | M3TAL API daemon (Go)       | Host-local only                             |
| `8081` | Traefik dashboard (admin UI) | Host-local only (e.g., `http://127.0.0.1:8081`) |
| `8082` | M3TAL Dashboard             | Direct (`local` mode) or via Traefik (`traefik` mode) |

---

## 🛠️ M3TAL CLI Command Reference

This section details every `m3tal` command, complete with descriptions and real usage examples.

### `sudo m3tal`

Opens the interactive TUI (Terminal User Interface) Control Center. This provides a menu-driven interface to perform common M3TAL operations.

*   **Description**: Launches an interactive, full-screen text-based interface for managing your M3TAL instance. Requires `sudo` as it interacts with Docker and system services.
*   **Usage Example**:
    ```bash
    sudo m3tal
    ```
    *(This will open a numbered menu, e.g., "1) Start All Stacks", "2) Stop All Stacks", etc.)*

---

### `m3tal init`

Generates the initial `/etc/m3tal/.env` file from default values.

*   **Description**: Essential for first-time installations. This command creates the primary configuration file if it doesn't exist, populating it with M3TAL's default environment variables.
*   **Usage Example**:
    ```bash
    m3tal init
    ```
    **Expected Output (if `/etc/m3tal/.env` did not exist):**
    ```
    /etc/m3tal/.env created successfully with default values.
    Please run 'm3tal config wizard' to customize your settings.
    ```
    **Expected Output (if `/etc/m3tal/.env` already exists):**
    ```
    /etc/m3tal/.env already exists. No action taken.
    ```

---

### `m3tal doctor`

Performs a pre-flight health check of the M3TAL ecosystem.

*   **Description**: Diagnoses common issues by checking Docker connectivity, the validity of `/etc/m3tal/.env`, and required port availability. Useful for troubleshooting.
*   **Usage Example**:
    ```bash
    m3tal doctor
    ```
    **Expected Output Example:**
    ```
    [✓] Docker daemon is running and accessible.
    [✓] /etc/m3tal/.env exists and is valid.
    [✓] Port 8080 (M3TAL API) is available.
    [✓] Port 8082 (M3TAL Dashboard) is available.
    [✓] Port 80 (Traefik HTTP) is available.
    [✓] M3TAL API service (m3tal-api.service) is running.
    [✓] All critical checks passed. Your M3TAL system appears healthy.
    ```

---

### `m3tal config wizard`

Launches an interactive wizard to configure `/etc/m3tal/.env`.

*   **Description**: Guides you through configuring key environment variables in `/etc/m3tal/.env`. It prompts for values and provides explanations, making initial setup and ongoing adjustments user-friendly.
*   **Usage Example**:
    ```bash
    m3tal config wizard
    ```
    **Expected Interactive Prompts Example:**
    ```
    Welcome to the M3TAL configuration wizard!

    This wizard will help you set up your /etc/m3tal/.env file.
    Press Enter to keep the current value, or type a new value.

    Current DASHBOARD_PORT (8082):
    Current DASHBOARD_EXPOSE_MODE (local): [local/traefik] traefik
    Current DOMAIN (localhost): example.com
    ... (continues for other variables)
    ```

---

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env`.

*   **Description**: Allows direct modification of a specific configuration key. If the key exists, its value is updated; otherwise, it's added.
*   **Usage Example**:
    *   Set the `DASHBOARD_EXPOSE_MODE` to `traefik`:
        ```bash
        m3tal config set DASHBOARD_EXPOSE_MODE traefik
        ```
    *   Set the primary `DOMAIN`:
        ```bash
        m3tal config set DOMAIN myhome.tech
        ```
    **Expected Output:**
    ```
    KEY 'DASHBOARD_EXPOSE_MODE' set to 'traefik' in /etc/m3tal/.env
    ```

---

### `m3tal config get KEY`

Reads the value of a single environment variable from `/etc/m3tal/.env`.

*   **Description**: Retrieves and displays the current value of a specified environment variable.
*   **Usage Example**:
    *   Get the current `DASHBOARD_PORT`:
        ```bash
        m3tal config get DASHBOARD_PORT
        ```
    **Expected Output:**
    ```
    8082
    ```
    *   Get the current `DOMAIN`:
        ```bash
        m3tal config get DOMAIN
        ```
    **Expected Output:**
    ```
    myhome.tech
    ```

---

### `m3tal config scan`

Lists all environment variables and their sources across all stacks.

*   **Description**: Provides a comprehensive overview of all environment variables M3TAL recognizes, including their default values and which Docker Compose stacks (or M3TAL itself) might use them.
*   **Usage Example**:
    ```bash
    m3tal config scan
    ```
    **Expected Output Example:**
    ```
    Environment Variables Scanned:
    - API_TOKEN (default: change_me_api_token)
    - ADMIN_PASSWORD (default: admin_pass)
    - BASE_STORAGE_PATH (default: ./data)
    - CONFIG_PATH (default: ./data/config)
    - DASHBOARD_EXPOSE_MODE (default: local)
    - DASHBOARD_PORT (default: 8082)
    - DASHBOARD_SECRET (default: change_me_immediately)
    - DEBUG_MODE (default: false)
    - DOMAIN (default: localhost)
    - DOWNLOADS_PATH (default: ./data/downloads)
    - HTTP_PORT (default: 8080)
    - LOCAL_IP (default: 127.0.0.1)
    - LOG_LEVEL (default: info)
    - MEDIA_PATH (default: ./data/media)
    - METRICS_ENABLED (default: true)
    - NETWORK_NAME (default: m3tal)
    - PGID (default: 1000)
    - PUID (default: 1000)
    - STATE_DIR (default: ./state)
    - TZ (default: America/Denver)
    - TRAEFIK_DASHBOARD_PORT (default: 8080)
    - TRAEFIK_WEBHTTPS_PORT (default: 443)
    - TRAEFIK_WEB_PORT (default: 80)
    - VPN_PASSWORD (default: password)
    - VPN_USER (default: user)
    ```

---

### `m3tal config list`

Lists the current contents of the `/etc/m3tal/.env` file.

*   **Description**: Displays the active key-value pairs configured in your primary environment file.
*   **Usage Example**:
    ```bash
    m3tal config list
    ```
    **Expected Output Example:**
    ```ini
    # M3TAL Environment Configuration
    DASHBOARD_PORT=8082
    DASHBOARD_EXPOSE_MODE=traefik
    HTTP_PORT=8080
    STATE_DIR=/var/lib/m3tal/state
    LOG_LEVEL=info
    DASHBOARD_SECRET=super_secret_dashboard_key
    API_TOKEN=super_secret_api_token
    ADMIN_PASSWORD=my_secure_admin_pass
    NETWORK_NAME=m3tal
    LOCAL_IP=192.168.1.100
    DOMAIN=myhome.tech
    TZ=America/New_York
    PUID=1000
    PGID=1000
    # ... more variables ...
    ```

---

### `m3tal dashpass [username] [password]`

Updates the M3TAL Dashboard user password.

*   **Description**: Manages user credentials for the M3TAL Dashboard, stored in `/docker/users.json`. If no `username` or `password` is provided, it launches an interactive prompt.
*   **Usage Examples**:
    *   **Interactive Mode (recommended)**:
        ```bash
        m3tal dashpass
        ```
        **Expected Interactive Prompts:**
        ```
        Enter username (default: admin): newuser
        Enter new password:
        Confirm new password:
        Password for 'newuser' updated successfully.
        ```
    *   **Direct Mode**:
        ```bash
        m3tal dashpass admin newSecurePassword123!
        ```
        **Expected Output:**
        ```
        Password for 'admin' updated successfully.
        ```

---

### `m3tal dash up`

Pulls the latest dashboard compose configuration and starts the dashboard container.

*   **Description**: Ensures your `m3tal-dashboard` is running with the most current configuration files and image. It downloads `m3tal-compose.yml` and its overrides, then starts the container using the appropriate `DASHBOARD_EXPOSE_MODE` setting.
*   **Usage Example**:
    ```bash
    m3tal dash up
    ```
    **Expected Output Example:**
    ```
    Fetching latest dashboard compose files...
    Files downloaded successfully to /opt/m3tal/stack/.
    Detected DASHBOARD_EXPOSE_MODE=traefik. Using m3tal-compose.traefik.yml.
    [+] Running 1/0
    ✔ Container m3tal-dashboard Started
    Dashboard container started successfully. Access at http://dash.myhome.tech (if Traefik is running) or http://HOST_IP:8082 (if local mode).
    ```

---

### `m3tal dash down`

Stops the M3TAL Dashboard container.

*   **Description**: Shuts down the `m3tal-dashboard` container.
*   **Usage Example**:
    ```bash
    m3tal dash down
    ```
    **Expected Output Example:**
    ```
    [+] Stopping 1/0
    ✔ Container m3tal-dashboard Stopped
    Dashboard container stopped.
    ```

---

### `m3tal dash restart`

Restarts the M3TAL Dashboard container.

*   **Description**: Stops and then starts the `m3tal-dashboard` container. Useful after configuration changes or for basic troubleshooting.
*   **Usage Example**:
    ```bash
    m3tal dash restart
    ```
    **Expected Output Example:**
    ```
    [+] Restarting 1/0
    ✔ Container m3tal-dashboard Restarted
    Dashboard container restarted.
    ```

---

### `m3tal dash logs`

Streams logs from the M3TAL Dashboard container.

*   **Description**: Displays real-time log output from the `m3tal-dashboard` container. Useful for debugging dashboard-specific issues.
*   **Usage Example**:
    ```bash
    m3tal dash logs
    ```
    **Expected Output Example:**
    ```
    Attaching to m3tal-dashboard
    m3tal-dashboard  |  * Serving Flask app 'server'
    m3tal-dashboard  |  * Debug mode: off
    m3tal-dashboard  | WARNING: This is a development server. Do not use it in a production deployment. Use a production WSGI server instead.
    m3tal-dashboard  |  * Running on all addresses (0.0.0.0)
    m3tal-dashboard  |  * Port: 8082
    m3tal-dashboard  | Press CTRL+C to quit
    m3tal-dashboard  | 192.168.1.1 - - [01/Jan/2023 12:34:56] "GET / HTTP/1.1" 200 -
    ```

---

### `m3tal dash status`

Shows the current status of the M3TAL Dashboard container.

*   **Description**: Reports whether the `m3tal-dashboard` container is running, stopped, or in another state.
*   **Usage Example**:
    ```bash
    m3tal dash status
    ```
    **Expected Output Example (Running):**
    ```
    Container m3tal-dashboard is Running (healthy)
    ```
    **Expected Output Example (Stopped):**
    ```
    Container m3tal-dashboard is Stopped
    ```

---

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files in `/docker/`.

*   **Description**: Starts or recreates all services defined in every Docker Compose file found in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`). This includes M3TAL's core components (Traefik, Cloudflared) and any user-defined stacks.
*   **Usage Example**:
    ```bash
    m3tal up
    ```
    **Expected Output Example:**
    ```
    [+] Running 4/4
    ✔ Container m3tal-dashboard Started
    ✔ Container traefik Started
    ✔ Container cloudflared Started
    ✔ Container ollama Started
    All M3TAL stacks and user services are now up.
    ```

---

### `m3tal down`

Runs `docker compose down` across all stacks.

*   **Description**: Stops and removes all containers, networks, and volumes defined in all `*-compose.yml` files within `/docker/`. This effectively shuts down your entire M3TAL ecosystem and all hosted services.
*   **Usage Example**:
    ```bash
    m3tal down
    ```
    **Expected Output Example:**
    ```
    [+] Stopping 4/4
    ✔ Container m3tal-dashboard Stopped
    ✔ Container traefik Stopped
    ✔ Container cloudflared Stopped
    ✔ Container ollama Stopped
    All M3TAL stacks and user services are now down.
    ```

---

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks.

*   **Description**: Collects and displays real-time log output from all Docker containers managed by M3TAL. Useful for monitoring overall system health and debugging issues across multiple services.
*   **Usage Example**:
    ```bash
    m3tal logs
    ```
    **Expected Output Example:**
    ```
    Attaching to m3tal-dashboard, traefik, cloudflared, ollama
    m3tal-dashboard  |  * Running on http://0.0.0.0:8082 (Press CTRL+C to quit)
    traefik          | time="2023-01-01T12:00:01Z" level=info msg="Configuration loaded from file."
    ollama           | time="2023-01-01T12:00:02Z" level=info msg="Ollama server started on 0.0.0.0:11434"
    cloudflared      | INFO Tunnel started with ID 123abc456
    m3tal-dashboard  | 172.18.0.1 - - [01/Jan/2023 12:00:05] "GET /api/status HTTP/1.1" 200 -
    ```

---

## ⚙️ Systemd Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage and monitor this service using standard `systemctl` and `journalctl` commands.

*   **Check API Service Status**:
    ```bash
    sudo systemctl status m3tal-api
    ```
    **Expected Output Example:**
    ```
    ● m3tal-api.service - M3TAL API Daemon
         Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
         Active: active (running) since Mon 2023-01-01 10:00:00 UTC; 1 day ago
       Main PID: 1234 (m3tal-api)
          Tasks: 7 (limit: 9272)
         Memory: 15.6M
            CPU: 1min 2.345s
         CGroup: /system.slice/m3tal-api.service
                 └─1234 /usr/bin/m3tal-api
    ```

*   **Restart API Service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```

*   **View Real-time API Service Logs**:
    ```bash
    sudo journalctl -u m3tal-api -f
    ```
    **Expected Output Example:**
    ```
    Jan 01 10:00:00 hostname m3tal-api[1234]: [INFO] M3TAL API started on :8080
    Jan 01 10:00:05 hostname m3tal-api[1234]: [INFO] Handling Docker request: /containers/json
    Jan 01 10:00:10 hostname m3tal-api[1234]: [INFO] Database updated.
    ```

---

## 🐳 Direct Docker Compose Fallback

For advanced users or troubleshooting scenarios, you can interact directly with Docker Compose. M3TAL's CLI commands largely abstract these operations, but knowing the direct commands can be helpful. Remember that M3TAL operates on the `/docker/` directory, which is a symlink to `/opt/m3tal/stack/`.

To manage all M3TAL-related stacks:

```bash
# Navigate to the M3TAL stack directory
cd /docker

# Start all defined services
sudo docker compose up -d

# Stop all defined services
sudo docker compose down

# View logs from all services
sudo docker compose logs -f

# List all services and their status
sudo docker compose ps
```

To manage a specific stack (e.g., only the `routing` stack with Traefik/Cloudflared):

```bash
# Assuming routing-compose.yml is in /docker/
cd /docker

# Start only the routing stack
sudo docker compose -f routing-compose.yml up -d

# Stop only the routing stack
sudo docker compose -f routing-compose.yml down

# View logs from only the routing stack
sudo docker compose -f routing-compose.yml logs -f
```

---