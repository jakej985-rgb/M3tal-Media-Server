As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the definitive CLI cheat-sheet for managing your M3TAL deployment. This document covers every essential command, its function, and practical usage examples, alongside critical architectural insights to empower your system administration.

---

# M3TAL CLI Command Reference

## I. Introduction to the M3TAL Ecosystem

The M3TAL ecosystem provides a robust, opinionated platform for deploying and managing self-hosted services using Docker and Docker Compose. At its core, M3TAL simplifies complex Docker operations through a unified Go CLI binary, a powerful API daemon, and a user-friendly web dashboard.

### M3TAL System Architecture Overview

*   **CLI binary (`/usr/bin/m3tal`)**: Your primary interface, installed via APT. All commands begin here.
*   **API daemon (`m3tal-api.service`)**: A Go binary running as a systemd service (port 8080). It orchestrates Docker, manages the internal state database, and exposes the API routes for the dashboard and CLI.
*   **Dashboard container (`m3tal-dashboard`)**: A Python/Flask web application running inside Docker (internal port 8082). It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway (`routing-compose.yml`)**: A powerful reverse proxy container that exposes your services by domain name on host port 80. It dynamically discovers services via Docker labels and loads static/dynamic configurations.
*   **Cloudflared (`routing-compose.yml`)**: An optional Cloudflare tunnel container, offering secure, zero-config internet access to your services without exposing ports directly.

### Filesystem Contract

M3TAL adheres to a strict filesystem layout for predictability and maintainability.

| Path                   | Purpose                                                                             |
| :--------------------- | :---------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`      | **Primary configuration file**. All M3TAL environment variables reside here. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Automatically created and managed by the API daemon.         |
| `/opt/m3tal/stack/`    | Canonical directory for M3TAL's core compose files and Traefik configuration.       |
| `/docker`              | **Symlink → `/opt/m3tal/stack/`**. This is the user-facing path for all stack operations. All custom `*-compose.yml` files should be placed here. |
| `/docker/users.json`   | Dashboard credential store. Managed by `m3tal dashpass`.                            |
| `/docker/dynamic/`     | Traefik dynamic configuration directory.                                            |

### Dashboard Access — Two Modes

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: M3TAL uses the `m3tal-compose.local.yml` override, which directly binds port `${DASHBOARD_PORT:-8082}` on the host to the dashboard container's internal port 8082.
*   **Access**: Navigate to `http://HOST_IP:8082` or `http://localhost:8082` in your web browser.
*   **Best For**: LAN-only setups, initial installation, local testing, or environments where a reverse proxy isn't immediately desired. No Traefik configuration is required for this mode.

#### Mode 2: `traefik`

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: M3TAL uses the `m3tal-compose.traefik.yml` override. This adds specific Traefik labels to the dashboard container, instructing Traefik to route requests for `dash.${DOMAIN}` to the dashboard on its internal port 8082. **Traefik must be running** via `m3tal up` for this mode to work.
*   **Access**: Navigate to `http://dash.YOUR_DOMAIN` (e.g., `http://dash.m3tal.local`) in your web browser.
*   **Best For**: Domain-based setups, exposing services behind a central reverse proxy, and integrating with other Traefik-managed services.

### Docker / Compose Runtime

M3TAL leverages **Docker Engine** and **Docker Compose V2** as its underlying container orchestration tools. These are hard dependencies for M3TAL's operation.

*   The `m3tal up` command orchestrates the `docker compose up` operation across **all `*-compose.yml` files** discovered within the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`). This includes core M3TAL services, routing infrastructure, and any custom user stacks.
*   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container. It first ensures the latest dashboard compose configurations (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) are pulled from GitHub. It then starts the dashboard, applying the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.

### Deployment Lifecycle — Day 2 Operations

To deploy a new application stack with M3TAL:

1.  **Place your compose file**: Copy your `docker-compose.yml` (e.g., `my-app-compose.yml`) into the `/docker/` directory.
2.  **Configure variables**: Ensure all necessary environment variables for your new stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` for an interactive setup or `m3tal config set KEY value` for specific variables.
3.  **Start all stacks**: Run `m3tal up` to bring up your new stack along with all existing M3TAL services.

### Traefik Routing Architecture

Traefik, deployed via `routing-compose.yml`, acts as the central ingress point for your domain-based services.

*   It binds to port 80 (and 443 if configured for HTTPS) on the host machine.
*   It automatically discovers services by reading Docker labels (e.g., those applied to `m3tal-dashboard` in `traefik` mode).
*   It loads dynamic configurations from `/docker/dynamic/` (using the file provider), allowing for hot-reloading of routing rules.
*   **Example**: `api.DOMAIN` is routed to the M3TAL API daemon (`http://host.docker.internal:8080`) via `dynamic/api.yml`.
*   **Example**: `dash.DOMAIN` is routed to the M3TAL dashboard container via Traefik labels, but only when `DASHBOARD_EXPOSE_MODE=traefik`.

### Port Map

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API daemon (Go)       | Host-local only                             |
| 8081 | Traefik dashboard           | Host-local only                             |
| 8082 | M3TAL Dashboard (container) | Direct port (local mode) or via Traefik     |

### APT Installation

If M3TAL is not yet installed on your system, use these commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## II. M3TAL CLI Command Reference

This section details every primary M3TAL command and its subcommands, complete with concrete usage examples.

### 1. `sudo m3tal`

Opens the interactive TUI (Text User Interface) Control Center. This provides a user-friendly, numbered menu for common M3TAL operations, abstracting away individual CLI commands.

*   **Usage:**
    ```bash
    sudo m3tal
    ```
*   **Example Output (Interactive):**
    ```
    M3TAL Control Center
    1. System Health Check (doctor)
    2. Configure M3TAL (.env wizard)
    3. Start All Stacks (up)
    4. Stop All Stacks (down)
    5. Restart M3TAL Dashboard (dash restart)
    6. View Aggregated Logs (logs)
    7. Exit
    Enter your choice: _
    ```

### 2. `m3tal init`

Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults. This command should be run on a first installation or when the `.env` file is missing.

*   **Usage:**
    ```bash
    m3tal init
    ```
*   **Example Output:**
    ```
    INFO: .env file not found. Generating /etc/m3tal/.env from defaults...
    INFO: /etc/m3tal/.env created successfully.
    INFO: Remember to run 'm3tal config wizard' to customize your settings.
    ```

### 3. `m3tal doctor`

Performs a comprehensive pre-flight health check of your M3TAL system. It verifies Docker connectivity, validates the `/etc/m3tal/.env` configuration, checks for critical port availability, and ensures core components are ready.

*   **Usage:**
    ```bash
    m3tal doctor
    ```
*   **Example Output:**
    ```
    INFO: Running M3TAL health checks...
    [OK] Docker Engine is running and accessible.
    [OK] /etc/m3tal/.env file is present and readable.
    [OK] Essential environment variables are set.
    [OK] Port 8080 (M3TAL API) is available.
    [OK] Port 80 (Traefik HTTP) is available.
    [OK] Port 8082 (Dashboard) is available.
    [OK] All core compose files found in /docker/.
    INFO: M3TAL system health is good.
    ```

### 4. `m3tal config wizard`

Launches an interactive command-line wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for managing your M3TAL environment variables.

*   **Usage:**
    ```bash
    m3tal config wizard
    ```
*   **Example Output (Interactive):**
    ```
    M3TAL Configuration Wizard
    Press Enter to keep current value, type new value and press Enter to change.
    
    Current DASHBOARD_EXPOSE_MODE [local]: traefik
    Current DOMAIN [localhost]: m3tal.local
    Current DASHBOARD_SECRET [change_me_immediately]: my_super_secret_key_123
    ... (prompts for other variables) ...
    INFO: /etc/m3tal/.env updated successfully.
    ```

### 5. `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env` to a specified value. This is useful for quick, non-interactive adjustments.

*   **Usage:**
    ```bash
    m3tal config set DASHBOARD_EXPOSE_MODE traefik
    m3tal config set DOMAIN m3tal.local
    ```
*   **Example Output:**
    ```
    INFO: Set DASHBOARD_EXPOSE_MODE to 'traefik' in /etc/m3tal/.env.
    INFO: Set DOMAIN to 'm3tal.local' in /etc/m3tal/.env.
    ```

### 6. `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

*   **Usage:**
    ```bash
    m3tal config get DASHBOARD_EXPOSE_MODE
    ```
*   **Example Output:**
    ```
    traefik
    ```

### 7. `m3tal config scan`

Lists all known environment variables across all detected Docker Compose stacks within `/docker/`, showing their default values and current values from `/etc/m3tal/.env`. This helps identify variables that might be needed by various services.

*   **Usage:**
    ```bash
    m3tal config scan
    ```
*   **Example Output (Partial):**
    ```
    ENV Key                   Default Value           Current Value
    ----------------------------------------------------------------------
    DASHBOARD_PORT            8082                    8082
    DASHBOARD_EXPOSE_MODE     local                   traefik
    HTTP_PORT                 8080                    8080
    DOMAIN                    localhost               m3tal.local
    PUID                      1000                    1000
    PGID                      1000                    1000
    TZ                        America/Denver          America/Denver
    TRAEFIK_WEB_PORT          80                      80
    ...
    ```

### 8. `m3tal config list`

Displays the entire contents of the `/etc/m3tal/.env` configuration file. This is equivalent to `cat /etc/m3tal/.env` but uses the M3TAL CLI for consistency.

*   **Usage:**
    ```bash
    m3tal config list
    ```
*   **Example Output (Partial):**
    ```
    # M3TAL System Environment Variables
    DASHBOARD_PORT=8082
    DASHBOARD_EXPOSE_MODE=traefik
    DOMAIN=m3tal.local
    DASHBOARD_SECRET=my_super_secret_key_123
    API_TOKEN=change_me_api_token
    ...
    ```

### 9. `m3tal dashpass [username] [password]`

Updates a user's password for the M3TAL Dashboard. If `username` and `password` are omitted, the command becomes interactive, prompting for the user and new password. This manages entries in `/docker/users.json`.

*   **Usage (Interactive):**
    ```bash
    m3tal dashpass
    ```
    *   **Example Interactive Flow:**
        ```
        Enter username: admin
        Enter new password: [TYPE_PASSWORD_HERE]
        Confirm new password: [TYPE_PASSWORD_HERE]
        INFO: Password for user 'admin' updated successfully.
        ```
*   **Usage (Direct):**
    ```bash
    m3tal dashpass admin MyStrongP@ssw0rd!
    ```
    *   **Example Output:**
        ```
        INFO: Password for user 'admin' updated successfully.
        ```

### 10. `m3tal dash up`

Pulls the latest dashboard compose configuration files from GitHub, then starts or updates the `m3tal-dashboard` container. It automatically selects the correct override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on `DASHBOARD_EXPOSE_MODE`.

*   **Usage:**
    ```bash
    m3tal dash up
    ```
*   **Example Output:**
    ```
    INFO: Pulling latest M3TAL dashboard compose configuration...
    INFO: Starting/updating m3tal-dashboard container...
    [+] Running 1/1
     ⠿ Container m3tal-dashboard  Started
    INFO: Dashboard is now running.
    ```

### 11. `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

*   **Usage:**
    ```bash
    m3tal dash down
    ```
*   **Example Output:**
    ```
    INFO: Stopping m3tal-dashboard container...
    [+] Running 1/0
     ⠿ Container m3tal-dashboard  Stopped
    INFO: Dashboard container stopped.
    ```

### 12. `m3tal dash restart`

Restarts the `m3tal-dashboard` container. This is useful after making configuration changes that affect the dashboard.

*   **Usage:**
    ```bash
    m3tal dash restart
    ```
*   **Example Output:**
    ```
    INFO: Restarting m3tal-dashboard container...
    [+] Running 1/1
     ⠿ Container m3tal-dashboard  Restarted
    INFO: Dashboard container restarted successfully.
    ```

### 13. `m3tal dash logs`

Streams the real-time logs from the `m3tal-dashboard` container. Press `Ctrl+C` to exit the log stream.

*   **Usage:**
    ```bash
    m3tal dash logs
    ```
*   **Example Output (Streaming):**
    ```
    Attaching to m3tal-dashboard
    m3tal-dashboard |  * Running on http://0.0.0.0:8082/ (Press CTRL+C to quit)
    m3tal-dashboard | 172.18.0.1 - - [21/Jul/2024 10:30:05] "GET / HTTP/1.1" 200 -
    m3tal-dashboard | 172.18.0.1 - - [21/Jul/2024 10:30:06] "GET /static/css/style.css HTTP/1.1" 200 -
    ... (logs continue until Ctrl+C) ...
    ```

### 14. `m3tal dash status`

Displays the current operational status of the `m3tal-dashboard` container.

*   **Usage:**
    ```bash
    m3tal dash status
    ```
*   **Example Output:**
    ```
    INFO: Status of m3tal-dashboard:
    [+] Running 1/1
     ⠿ Container m3tal-dashboard  Running 15 minutes ago
    ```

### 15. `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This command starts or recreates all defined services in the background (`-d`). Use this to deploy new stacks or update existing ones.

*   **Usage:**
    ```bash
    m3tal up
    ```
*   **Example Output:**
    ```
    INFO: Running docker compose up -d for all stacks in /docker/...
    [+] Running 4/4
     ⠿ Network proxy         Created
     ⠿ Container traefik     Started
     ⠿ Container ollama      Started
     ⠿ Container m3tal-dashboard  Started
     ⠿ Container cloudflared Started
    INFO: All M3TAL stacks are now up and running.
    ```

### 16. `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This command gracefully stops and removes all containers, networks, and volumes defined by the compose files.

*   **Usage:**
    ```bash
    m3tal down
    ```
*   **Example Output:**
    ```
    INFO: Running docker compose down for all stacks in /docker/...
    [+] Running 4/4
     ⠿ Container traefik     Stopped
     ⠿ Container ollama      Stopped
     ⠿ Container m3tal-dashboard  Stopped
     ⠿ Container cloudflared Stopped
     ⠿ Network proxy         Removed
    INFO: All M3TAL stacks have been taken down.
    ```

### 17. `m3tal logs`

Streams aggregated, real-time logs from all currently running M3TAL Docker containers. This provides a consolidated view of system activity. Press `Ctrl+C` to exit the log stream.

*   **Usage:**
    ```bash
    m3tal logs
    ```
*   **Example Output (Streaming):**
    ```
    Attaching to cloudflared, m3tal-dashboard, ollama, traefik
    m3tal-dashboard |  * Running on http://0.0.0.0:8082/ (Press CTRL+C to quit)
    traefik          | time="2024-07-21T10:35:10Z" level=info msg="Starting provider aggregator.ProviderAggregator"
    ollama           | time="2024-07-21 10:35:11" level=info msg="Listening on 0.0.0.0:11434 (version 0.1.34)"
    cloudflared      | 2024-07-21T10:35:12Z INF Starting tunnel tunnelID=your-tunnel-id...
    m3tal-dashboard | 172.18.0.1 - - [21/Jul/2024 10:35:15] "GET /api/v1/status HTTP/1.1" 200 -
    ... (logs continue from all services until Ctrl+C) ...
    ```

---

## III. M3TAL Systemd Service Management

The M3TAL API daemon is a critical component running as a systemd service. Understanding how to manage it directly is essential for advanced troubleshooting or service restarts.

*   **Check API daemon status:**
    ```bash
    systemctl status m3tal-api
    ```
    *   **Example Output:**
        ```
        ● m3tal-api.service - M3TAL API Daemon
             Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
             Active: active (running) since Sun 2024-07-21 10:20:00 UTC; 20min ago
           Main PID: 1234 (m3tal-api)
              Tasks: 7 (limit: 9277)
             Memory: 15.6M
                CPU: 234ms
             CGroup: /system.slice/m3tal-api.service
                     └─1234 /usr/bin/m3tal-api
        
        Jul 21 10:20:00 hostname systemd[1]: Started M3TAL API Daemon.
        Jul 21 10:20:01 hostname m3tal-api[1234]: INFO: M3TAL API daemon started on :8080
        ```

*   **Restart API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
    *   **Example Output:** (No direct output, but `systemctl status` would show a new uptime.)

*   **Stream API daemon logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
    *   **Example Output (Streaming):**
        ```
        Jul 21 10:45:00 hostname m3tal-api[1234]: INFO: Handling request /api/v1/health
        Jul 21 10:45:05 hostname m3tal-api[1234]: INFO: Processing Docker event: container start m3tal-dashboard
        ... (logs continue until Ctrl+C) ...
        ```

---

## IV. Docker Direct Commands (Fallback)

While M3TAL's CLI provides a streamlined interface, directly interacting with Docker Compose can be useful for advanced debugging, inspecting individual services, or when the M3TAL CLI itself encounters an issue. Remember that M3TAL explicitly manages compose files within the `/docker/` directory (symlinked from `/opt/m3tal/stack/`).

*   **General Principle**: To replicate `m3tal up` or `m3tal down` using direct Docker Compose, you typically need to specify all relevant `*-compose.yml` files.

*   **Start specific M3TAL core stacks (e.g., routing and dashboard):**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml up -d
    ```
    *   **Explanation**: This command tells Docker Compose to combine `routing-compose.yml` and `m3tal-compose.yml` (and their respective overrides if present) and start the services in detached mode.

*   **Start a single user-defined stack (e.g., `ollama`):**
    Assume you have `ollama-compose.yml` in `/docker/`.
    ```bash
    sudo docker compose -f /docker/ollama-compose.yml up -d
    ```

*   **View status of all running containers from specific stacks:**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml ps
    ```
    *   **Example Output (Partial):**
        ```
        NAME                IMAGE                           COMMAND                  SERVICE             CREATED        STATUS                  PORTS
        cloudflared         cloudflare/cloudflared:latest   "cloudflared --no-aut…"  cloudflared         2 minutes ago  Up 2 minutes
        m3tal-dashboard     ghcr.io/jakej985-rgb/m3tal-go…  "python3 server.py"      m3tal-dashboard     2 minutes ago  Up 2 minutes            8082/tcp
        traefik             traefik:latest                  "/entrypoint.sh --api…"  traefik             2 minutes ago  Up 2 minutes (healthy)  0.0.0.0:80->80/tcp, 127.0.0.1:8081->8080/tcp
        ```

*   **Stream logs from all services in specific stacks:**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml logs -f
    ```

*   **Stop and remove specific M3TAL core stacks:**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml down
    ```

*   **Inspect a specific container (e.g., `m3tal-dashboard`):**
    ```bash
    sudo docker inspect m3tal-dashboard
    ```

*   **Execute a command inside a running container (e.g., check Python version in dashboard):**
    ```bash
    sudo docker exec -it m3tal-dashboard python3 --version
    ```
    *   **Example Output:**
        ```
        Python 3.9.18
        ```

---

This comprehensive guide should equip you with the knowledge to effectively manage your M3TAL ecosystem. For further assistance or to explore advanced configurations, consult the official M3TAL documentation and community resources.

**DocSmith, M3TAL Ecosystem Documentation Architect**