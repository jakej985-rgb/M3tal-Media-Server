# M3TAL CLI Command Reference

Greetings, Operative. This document serves as your essential guide to the M3TAL Command-Line Interface (CLI). As DocSmith, the M3TAL Ecosystem Documentation Architect, my purpose is to equip you with the knowledge to wield M3TAL's full power. The CLI is your primary interface for managing your self-hosted stack, from initial setup to day-to-day operations and diagnostics.

## I. M3TAL System Architecture Overview

The M3TAL ecosystem is designed for robustness and ease of management, built upon a foundation of Docker and systemd.

### Core Components:

*   **CLI Binary (`/usr/bin/m3tal`):** Your single point of interaction, handling all high-level commands.
*   **API Daemon (`m3tal-api.service`):** A Go binary running as a systemd service (port 8080), the control plane for Docker interactions and state management via SQLite.
*   **Dashboard Container (`m3tal-dashboard`):** A Python/Flask application providing a web-based GUI, communicating with the API daemon.
*   **Traefik Gateway (`routing-compose.yml`):** The reverse proxy container, exposing services on port 80 and handling domain-based routing.
*   **Cloudflared (`routing-compose.yml`):** An optional Cloudflare tunnel container for secure, zero-config external access.

### Filesystem Contract:

Understanding the M3TAL filesystem layout is critical for effective management.

| Path                           | Purpose                                                                                                                                                                                                                               |
| :----------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`              | **Primary Configuration File:** Stores all M3TAL environment variables. *Managed by `m3tal config wizard`*.                                                                                                                             |
| `/var/lib/m3tal/state.db`      | **SQLite State Database:** Stores M3TAL's internal state. *Auto-created and managed by the `m3tal-api.service`*.                                                                                                                        |
| `/opt/m3tal/stack/`            | **Canonical Stack Directory:** The actual location where Docker Compose files and Traefik dynamic configurations are stored.                                                                                                            |
| `/docker`                      | **User-Facing Stack Symlink:** A symbolic link to `/opt/m3tal/stack/`. This is the recommended path for users to place their custom `*-compose.yml` files and dynamic Traefik configurations.                                          |
| `/docker/users.json`           | **Dashboard Credential Store:** Encrypted storage for dashboard user accounts. *Managed by `m3tal dashpass`*.                                                                                                                           |
| `/docker/dynamic/`             | **Traefik Dynamic Config:** Directory for Traefik's file provider. Place `.yml` files here for dynamic routing configurations (e.g., `api.yml` for API access).                                                                        |
| `/docker/m3tal-compose.yml`    | Base Docker Compose file for the M3TAL Dashboard and core components.                                                                                                                                                                 |
| `/docker/routing-compose.yml`  | Docker Compose file for Traefik and Cloudflared (if enabled).                                                                                                                                                                         |
| `/docker/traefik.yml`          | Static Traefik configuration.                                                                                                                                                                                                         |

### Docker / Compose Runtime:

M3TAL leverages **Docker Engine** and **Docker Compose V2** for container orchestration.
*   The `m3tal up` command orchestrates all `*-compose.yml` files found within the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the dashboard container, downloading its latest compose configurations and applying overrides based on your `DASHBOARD_EXPOSE_MODE`.
*   To add a new service or "stack," simply place its `*-compose.yml` file into `/docker/` and ensure its required environment variables are configured in `/etc/m3tal/.env`.

### Dashboard Access Modes:

The M3TAL Dashboard can be accessed in two distinct ways, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

1.  **Local Mode (Default):**
    *   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
    *   **Mechanism:** Uses `m3tal-compose.local.yml` to create a direct port binding, exposing the dashboard on the host machine.
    *   **Access:** `http://HOST_IP:8082` or `http://localhost:8082`
    *   **Use Case:** Ideal for LAN-only setups, initial deployment, or when Traefik is not yet configured or desired. No Traefik required.

2.  **Traefik Mode:**
    *   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
    *   **Mechanism:** Uses `m3tal-compose.traefik.yml` to apply Traefik labels, allowing Traefik to route traffic to the dashboard container via a domain. Requires Traefik to be running (`m3tal up`).
    *   **Access:** `http://dash.YOUR_DOMAIN` (e.g., `http://dash.myhomelab.com`)
    *   **Use Case:** Best for domain-based access, integrating into a larger system managed by Traefik, and leveraging Traefik's features (e.g., SSL).

### Traefik Routing Architecture:

Traefik (deployed via `routing-compose.yml`) acts as M3TAL's intelligent edge router:

*   Binds to port 80 (and 443 if SSL is configured) on the host.
*   Automatically discovers services via Docker labels (e.g., `m3tal-compose.traefik.yml` for the dashboard).
*   Loads dynamic configurations from `/docker/dynamic/` (e.g., `api.yml` for the M3TAL API).
*   Routes `api.YOUR_DOMAIN` to the M3TAL API daemon (`http://host.docker.internal:8080`).
*   Routes `dash.YOUR_DOMAIN` to the `m3tal-dashboard` container (when `DASHBOARD_EXPOSE_MODE=traefik`).

### Port Map:

| Port | Service                                 | Access                                             |
| :--- | :-------------------------------------- | :------------------------------------------------- |
| 80   | Traefik HTTP entry point                | Public (when `routing-compose.yml` is active)      |
| 8080 | M3TAL API daemon (Go)                   | Host-local only (internal communication)           |
| 8081 | Traefik dashboard (admin interface)     | Host-local only (e.g., `http://127.0.0.1:8081`)    |
| 8082 | M3TAL Dashboard (Python/Flask)          | Direct port (local mode) or via Traefik (traefik mode) |

## II. M3TAL CLI Command Reference

This section details every primary M3TAL command and its subcommands, providing usage, examples, and important notes.

---

### `sudo m3tal`

*   **Purpose:** Opens the interactive M3TAL Control Center (TUI - Text User Interface). This provides a menu-driven interface for common operations, status checks, and configuration adjustments.
*   **Usage:** `sudo m3tal`
*   **Example:**
    ```bash
    sudo m3tal
    # Welcome to the M3TAL Control Center!
    # 1. View System Status
    # 2. Manage Docker Stacks
    # 3. Configure M3TAL
    # 4. Diagnostics & Health Check
    # 5. Exit
    # Enter your choice: _
    ```
*   **Notes:** Requires `sudo` as it interacts with Docker and systemd, which typically require root privileges. Offers a user-friendly alternative to command-line arguments for many operations.

---

### `m3tal init`

*   **Purpose:** Initializes the M3TAL environment by generating the default `/etc/m3tal/.env` configuration file. This is a crucial first step after installation.
*   **Usage:** `m3tal init`
*   **Example:**
    ```bash
    m3tal init
    # /etc/m3tal/.env created successfully with default values.
    # Please review and configure using 'm3tal config wizard'.
    ```
*   **Notes:** Should be run only once on first installation. If `/etc/m3tal/.env` already exists, it will prompt for confirmation before overwriting to prevent data loss.

---

### `m3tal doctor`

*   **Purpose:** Performs a pre-flight health check of the M3TAL ecosystem. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, checks for required network existence, and ensures necessary ports are available.
*   **Usage:** `m3tal doctor`
*   **Example:**
    ```bash
    m3tal doctor
    # [✓] Docker daemon is running and accessible.
    # [✓] /etc/m3tal/.env file exists and is valid.
    # [✓] Docker 'proxy' network exists.
    # [✓] Port 8080 (M3TAL API) is available.
    # [✓] Port 8082 (M3TAL Dashboard) is available (DASHBOARD_EXPOSE_MODE=local).
    # [✓] Port 80 (Traefik) is available.
    # M3TAL system health: GOOD
    ```
*   **Notes:** Essential for troubleshooting and ensuring a clean environment before deploying or updating stacks.

---

### `m3tal config wizard`

*   **Purpose:** Launches an interactive, guided wizard to configure or modify variables in `/etc/m3tal/.env`. This is the recommended way to manage your primary configuration.
*   **Usage:** `m3tal config wizard`
*   **Example:**
    ```bash
    m3tal config wizard
    # M3TAL Configuration Wizard
    # Current value for DASHBOARD_PORT (default: 8082): 8082
    # Enter new value (or press Enter to keep current): 8083
    # ... (prompts for other variables)
    # Configuration saved to /etc/m3tal/.env.
    ```
*   **Notes:** This wizard ensures all necessary environment variables are set correctly. It reads the existing `.env` file and presents default values. Some changes (e.g., `DASHBOARD_EXPOSE_MODE`) may require restarting relevant containers (`m3tal dash restart` or `m3tal up`).

---

### `m3tal config set KEY VALUE`

*   **Purpose:** Sets or updates a single environment variable in `/etc/m3tal/.env` directly from the command line.
*   **Usage:** `m3tal config set KEY VALUE`
*   **Example:**
    ```bash
    m3tal config set DASHBOARD_SECRET my_new_strong_secret_123
    # DASHBOARD_SECRET set to 'my_new_strong_secret_123' in /etc/m3tal/.env.
    ```
*   **Notes:** Changes made with `m3tal config set` require restarting any affected Docker containers or the `m3tal-api.service` for the changes to take effect. Be cautious with sensitive values on shared systems, as they appear in shell history.

---

### `m3tal config get KEY`

*   **Purpose:** Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.
*   **Usage:** `m3tal config get KEY`
*   **Example:**
    ```bash
    m3tal config get DASHBOARD_PORT
    # 8082

    m3tal config get DOMAIN
    # myhomelab.com
    ```
*   **Notes:** Useful for quickly checking a single configuration setting.

---

### `m3tal config scan`

*   **Purpose:** Lists all environment variables discovered across all Docker Compose files in `/docker/` that are currently defined in `/etc/m3tal/.env`. It helps identify which variables are used by your stacks.
*   **Usage:** `m3tal config scan`
*   **Example:**
    ```bash
    m3tal config scan
    # Detected Environment Variables in use across stacks:
    # ----------------------------------------------------
    # DASHBOARD_PORT (m3tal-compose.yml) = 8082
    # DASHBOARD_EXPOSE_MODE (m3tal-compose.yml) = local
    # DOMAIN (routing-compose.yml, m3tal-compose.traefik.yml) = myhomelab.com
    # PUID (m3tal-compose.yml) = 1000
    # PGID (m3tal-compose.yml) = 1000
    # TZ (m3tal-compose.yml) = America/New_York
    # TRAEFIK_WEB_PORT (routing-compose.yml) = 80
    # ...
    ```
*   **Notes:** This command provides a holistic view of variable usage, aiding in understanding stack dependencies and debugging configuration issues.

---

### `m3tal config list`

*   **Purpose:** Displays the full contents of the `/etc/m3tal/.env` file.
*   **Usage:** `m3tal config list`
*   **Example:**
    ```bash
    m3tal config list
    # # M3TAL Environment Configuration
    # DASHBOARD_PORT=8082
    # DASHBOARD_EXPOSE_MODE=local
    # HTTP_PORT=8080
    # STATE_DIR=./state
    # LOG_LEVEL=info
    # DASHBOARD_SECRET=my_new_strong_secret_123
    # API_TOKEN=change_me_api_token
    # ADMIN_PASSWORD=admin_pass
    # NETWORK_NAME=proxy
    # LOCAL_IP=127.0.0.1
    # DOMAIN=localhost
    # PUID=1000
    # PGID=1000
    # TZ=America/Denver
    # ...
    ```
*   **Notes:** Provides a quick way to review your entire M3TAL configuration. Sensitive information will be displayed, so use with caution in secure environments.

---

### `m3tal dashpass [username] [password]`

*   **Purpose:** Manages user passwords for the M3TAL Dashboard.
*   **Usage:**
    *   Interactive: `m3tal dashpass` (prompts for username and password)
    *   Direct: `m3tal dashpass <username> <password>`
*   **Example (Interactive):**
    ```bash
    m3tal dashpass
    # Enter dashboard username (e.g., admin): newuser
    # Enter new password for newuser: **********
    # Confirm new password: **********
    # Password for 'newuser' updated successfully in /docker/users.json.
    ```
*   **Example (Direct):**
    ```bash
    m3tal dashpass admin MySuperSecurePass123
    # Password for 'admin' updated successfully in /docker/users.json.
    ```
*   **Notes:** This command updates the `/docker/users.json` file. It's crucial for securing your dashboard. Passwords are hashed before storage. Always use strong, unique passwords.

---

### `m3tal dash up`

*   **Purpose:** Pulls the latest M3TAL Dashboard Docker Compose configuration files from GitHub, then starts or updates the `m3tal-dashboard` container based on the `DASHBOARD_EXPOSE_MODE` setting.
*   **Usage:** `m3tal dash up`
*   **Example:**
    ```bash
    m3tal dash up
    # Pulling latest dashboard compose files from GitHub...
    # Using dashboard expose mode: local (from /etc/m3tal/.env)
    # Pulling m3tal-dashboard (ghcr.io/jakej985-rgb/m3tal-godash:debug)...done
    # Creating m3tal-dashboard ... done
    # Dashboard container 'm3tal-dashboard' is up and running.
    # Access at http://localhost:8082
    ```
*   **Notes:** This command is essential for deploying or updating the dashboard. It intelligently applies `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` as an override based on `DASHBOARD_EXPOSE_MODE` to configure port exposure or Traefik labels, respectively.

---

### `m3tal dash down`

*   **Purpose:** Stops and removes the `m3tal-dashboard` container.
*   **Usage:** `m3tal dash down`
*   **Example:**
    ```bash
    m3tal dash down
    # Stopping m3tal-dashboard ... done
    # Removing m3tal-dashboard ... done
    # Dashboard container 'm3tal-dashboard' has been stopped and removed.
    ```
*   **Notes:** This only affects the dashboard container, not other running stacks.

---

### `m3tal dash restart`

*   **Purpose:** Restarts the `m3tal-dashboard` container. Useful after configuration changes to `/etc/m3tal/.env` that affect the dashboard.
*   **Usage:** `m3tal dash restart`
*   **Example:**
    ```bash
    m3tal dash restart
    # Restarting m3tal-dashboard ... done
    # Dashboard container 'm3tal-dashboard' restarted.
    ```
*   **Notes:** This command essentially performs a `down` followed by an `up` for the dashboard specifically.

---

### `m3tal dash logs`

*   **Purpose:** Streams real-time logs from the `m3tal-dashboard` container.
*   **Usage:** `m3tal dash logs`
*   **Example:**
    ```bash
    m3tal dash logs
    # Attaching to m3tal-dashboard
    # m3tal-dashboard | INFO:werkzeug: * Running on http://0.0.0.0:8082 (Press CTRL+C to quit)
    # m3tal-dashboard | INFO:werkzeug: * Debug mode: off
    # m3tal-dashboard | 172.18.0.1 - - [20/Oct/2023 10:30:45] "GET / HTTP/1.1" 200 -
    # ... (logs continue to stream)
    ```
*   **Notes:** Press `Ctrl+C` to exit the log stream.

---

### `m3tal dash status`

*   **Purpose:** Shows the current status of the `m3tal-dashboard` container.
*   **Usage:** `m3tal dash status`
*   **Example:**
    ```bash
    m3tal dash status
    # Container Name    Status    Ports
    # ------------------------------------
    # m3tal-dashboard   running   0.0.0.0:8082->8082/tcp
    ```
*   **Notes:** Provides a quick overview of whether the dashboard is running and its exposed ports.

---

### `m3tal up`

*   **Purpose:** Orchestrates all Docker Compose stacks defined by `*-compose.yml` files in the `/docker/` directory. This command pulls latest images (if available) and starts/recreates all services.
*   **Usage:** `m3tal up`
*   **Example:**
    ```bash
    m3tal up
    # Pulling routing-traefik (traefik:latest)...done
    # Creating network proxy...
    # Pulling ai-ollama (ollama/ollama:latest)...done
    # Creating container m3tal-dashboard ...
    # Creating container routing-traefik ...
    # Creating container ai-ollama ...
    # All Docker Compose stacks in /docker/ are up and running.
    ```
*   **Notes:** This command is your primary method for deploying and updating your entire M3TAL ecosystem (including Traefik, custom services, etc.). It ensures all services are in their desired state as defined by their compose files and environment variables in `/etc/m3tal/.env`.

---

### `m3tal down`

*   **Purpose:** Stops and removes all Docker Compose services and their associated networks that were started by `m3tal up` across all `*-compose.yml` files in `/docker/`.
*   **Usage:** `m3tal down`
*   **Example:**
    ```bash
    m3tal down
    # Stopping ai-ollama ... done
    # Stopping routing-traefik ... done
    # Stopping m3tal-dashboard ... done
    # Removing ai-ollama ... done
    # Removing routing-traefik ... done
    # Removing m3tal-dashboard ... done
    # All Docker Compose stacks have been stopped and removed.
    ```
*   **Notes:** Use with caution. This command will bring down your entire M3TAL-managed service infrastructure. Data volumes are generally preserved unless explicitly configured otherwise in compose files.

---

### `m3tal logs`

*   **Purpose:** Streams aggregated logs from all currently running Docker containers managed by M3TAL.
*   **Usage:** `m3tal logs`
*   **Example:**
    ```bash
    m3tal logs
    # Attaching to m3tal-dashboard, routing-traefik, ai-ollama
    # m3tal-dashboard | INFO:werkzeug: * Running on http://0.0.0.0:8082
    # routing-traefik | traefik.go:123 [INFO] Traefik starting...
    # ai-ollama       | ollama is now listening on 0.0.0.0:11434
    # ... (logs from all services stream concurrently)
    ```
*   **Notes:** An invaluable tool for monitoring the overall health and activity of your ecosystem. Press `Ctrl+C` to exit the log stream.

## III. M3TAL Systemd Service Management

The M3TAL API daemon (`m3tal-api.service`) is a critical component running as a systemd service. It is responsible for handling Docker interactions and maintaining the system's state.

*   **Check API Service Status:**
    ```bash
    sudo systemctl status m3tal-api
    # ● m3tal-api.service - M3TAL API Daemon
    #      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
    #      Active: active (running) since Fri 2023-10-20 09:00:00 UTC; 1h ago
    #    Main PID: 1234 (m3tal-api)
    #       Tasks: 7 (limit: 9277)
    #      Memory: 15.6M
    #         CPU: 123ms
    #      CGroup: /system.slice/m3tal-api.service
    #              └─1234 /usr/bin/m3tal-api
    # Oct 20 09:00:00 hostname systemd[1]: Started M3TAL API Daemon.
    # Oct 20 09:00:01 hostname m3tal-api[1234]: API server listening on :8080
    ```
*   **Stream API Service Logs:**
    ```bash
    sudo journalctl -u m3tal-api -f
    # -- Logs begin at Fri 2023-10-20 08:00:00 UTC, end at Fri 2023-10-20 10:00:00 UTC. --
    # Oct 20 09:00:01 hostname m3tal-api[1234]: API server listening on :8080
    # Oct 20 09:00:05 hostname m3tal-api[1234]: [INFO] Health check successful.
    # ... (logs continue to stream)
    ```
*   **Restart API Service (e.g., after manual `.env` changes affecting the API):**
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **Stop API Service:**
    ```bash
    sudo systemctl stop m3tal-api
    ```
*   **Start API Service:**
    ```bash
    sudo systemctl start m3tal-api
    ```

## IV. Direct Docker Compose Fallback

While the `m3tal` CLI provides a powerful abstraction layer, M3TAL is built directly on Docker Engine and Docker Compose V2. For advanced troubleshooting, specific debugging, or when M3TAL's API daemon might be unavailable, you can always fall back to direct `docker compose` commands.

All M3TAL-managed compose files reside in `/docker/` (which is a symlink to `/opt/m3tal/stack/`).

*   **View Services for a Specific Stack:**
    ```bash
    # Example: Check the services defined in routing-compose.yml
    docker compose -f /docker/routing-compose.yml ps
    ```
*   **Start a Specific Stack:**
    ```bash
    # Example: Start only the routing stack
    docker compose -f /docker/routing-compose.yml up -d
    ```
*   **Stop a Specific Stack:**
    ```bash
    # Example: Stop only the routing stack
    docker compose -f /docker/routing-compose.yml down
    ```
*   **View Logs for a Specific Stack:**
    ```bash
    # Example: Stream logs from the m3tal-dashboard (using its base compose file)
    docker compose -f /docker/m3tal-compose.yml logs -f
    ```
*   **Execute a Command in a Running Container:**
    ```bash
    # Example: Open a shell in the m3tal-dashboard container
    docker compose -f /docker/m3tal-compose.yml exec m3tal-dashboard bash
    ```
*   **Rebuild and Restart a Specific Service:**
    ```bash
    # Example: Rebuild and restart the 'traefik' service within the routing stack
    docker compose -f /docker/routing-compose.yml up -d --build traefik
    ```

**Important Note:** When using direct `docker compose` commands, ensure you specify all relevant compose files if your services span multiple files (e.g., `docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml ps`). However, `m3tal up` automatically handles the discovery and application of these files, including environment-specific overrides.

## V. APT Installation

To install the M3TAL CLI and API daemon via APT, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This concludes your M3TAL CLI cheat-sheet. Embrace the power of M3TAL and build your ultimate self-hosted ecosystem!

DocSmith
M3TAL Ecosystem Documentation Architect