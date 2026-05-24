# M3TAL CLI Command Reference

Greetings, M3TAL Operative! This document serves as your definitive guide to the M3TAL Command Line Interface (CLI). The `m3tal` binary is your single point of control for managing your M3TAL ecosystem, from initial setup and configuration to daily operations and diagnostics.

The M3TAL ecosystem is designed for robust, self-hosted applications, leveraging Docker and systemd for stability and ease of management.

**Core Components at a Glance:**
*   **CLI binary (`/usr/bin/m3tal`):** The tool you're documenting here.
*   **API daemon (`m3tal-api.service`):** A Go service listening on port 8080, handling Docker interactions and state.
*   **Dashboard container (`m3tal-dashboard`):** A Python/Flask web interface for M3TAL, running internally on port 8082.
*   **Traefik gateway (`routing-compose.yml`):** The reverse proxy, making services accessible via domain names on port 80.
*   **Cloudflared (`routing-compose.yml`):** Optional secure tunnels for exposing services without port forwarding.

---

## 1. Filesystem Contract

The following paths are critical for the M3TAL ecosystem's operation and configuration:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file.** Managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.     |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains Docker Compose files and Traefik dynamic configuration. |
| `/docker`                   | **Symlink** → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your `*-compose.yml` files here. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |
| `/docker/dynamic/`          | Directory for Traefik's file provider, for dynamic routing configuration. |

---

## 2. APT Installation

Before using the CLI, ensure M3TAL is installed. Execute these commands on your system:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## 3. Core M3TAL Commands

These commands provide fundamental control and diagnostic capabilities for your M3TAL system.

### `sudo m3tal`
Opens the interactive Text-User Interface (TUI) Control Center. This provides a numbered menu for common operations and system status.
```bash
sudo m3tal
```
*Example Output (abbreviated):*
```
======================================
     M3TAL Control Center (TUI)
======================================
1. System Status & Health Check
2. Manage Stacks (Up/Down/Logs)
3. Configure M3TAL (.env wizard)
4. Dashboard Management
5. Exit
Enter your choice (1-5):
```

### `m3tal init`
Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults. This command is crucial for a first-time installation or to reset your configuration. It ensures all necessary environment variables are present.
```bash
m3tal init
```
*Example Usage:*
```bash
# Initialize the .env file if it doesn't exist
m3tal init
```

### `m3tal doctor`
Performs a comprehensive pre-flight health check of your M3TAL system. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability to prevent conflicts.
```bash
m3tal doctor
```
*Example Usage:*
```bash
# Run a system health check before deploying new services
m3tal doctor
```

---

## 4. Configuration Management (`m3tal config`)

The `m3tal config` commands allow you to manage the `/etc/m3tal/.env` configuration file. This file dictates how M3TAL operates, including network settings, storage paths, and dashboard access.

### `m3tal config wizard`
Launches an interactive wizard to guide you through configuring or updating your `/etc/m3tal/.env` file. This is the recommended way to manage your M3TAL settings.
```bash
m3tal config wizard
```
*Example Usage:*
```bash
# Open the interactive wizard to change network settings
m3tal config wizard
```

### `m3tal config set KEY VALUE`
Sets a single environment variable within `/etc/m3tal/.env` to the specified value. Useful for quick, non-interactive changes.
```bash
m3tal config set KEY VALUE
```
*Example Usage:*
```bash
# Set the primary domain for Traefik routing
m3tal config set DOMAIN example.com

# Change the dashboard expose mode to 'traefik'
m3tal config set DASHBOARD_EXPOSE_MODE traefik
```

### `m3tal config get KEY`
Retrieves and displays the value of a specific environment variable from `/etc/m3tal/.env`.
```bash
m3tal config get KEY
```
*Example Usage:*
```bash
# Get the currently configured dashboard port
m3tal config get DASHBOARD_PORT

# Check the configured timezone
m3tal config get TZ
```

### `m3tal config scan`
Lists all known environment variables across all M3TAL stacks, along with their default values. This helps identify available configuration options.
```bash
m3tal config scan
```
*Example Usage:*
```bash
# See all possible configuration keys and their defaults
m3tal config scan
```

### `m3tal config list`
Displays the current contents of your primary configuration file, `/etc/m3tal/.env`.
```bash
m3tal config list
```
*Example Usage:*
```bash
# Review the active M3TAL configuration
m3tal config list
```

---

## 5. Dashboard Management (`m3tal dash`)

These commands are specifically for managing the M3TAL Dashboard container, which provides a web-based interface for your M3TAL system.

### `m3tal dashpass [username] [password]`
Updates the password for a specified dashboard user. If `username` and `password` are omitted, it will prompt you interactively. Dashboard credentials are stored in `/docker/users.json`.
```bash
m3tal dashpass [username] [password]
```
*Example Usage (interactive):*
```bash
# Interactively set a new password for the 'admin' user
m3tal dashpass admin
# Follow prompts: Enter new password, Confirm new password
```
*Example Usage (direct):*
```bash
# Set a new password for the 'admin' user non-interactively
m3tal dashpass admin MySuperStrongPassword123!
```

### `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the `m3tal-dashboard` container. This command intelligently uses the appropriate compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.
```bash
m3tal dash up
```
*Example Usage:*
```bash
# Start or update the M3TAL Dashboard
m3tal dash up
```

### `m3tal dash down`
Stops and removes the `m3tal-dashboard` container.
```bash
m3tal dash down
```
*Example Usage:*
```bash
# Stop the M3TAL Dashboard container
m3tal dash down
```

### `m3tal dash restart`
Restarts the `m3tal-dashboard` container. This is useful after making configuration changes that affect the dashboard.
```bash
m3tal dash restart
```
*Example Usage:*
```bash
# Restart the dashboard after updating its configuration
m3tal dash restart
```

### `m3tal dash logs`
Streams the logs from the `m3tal-dashboard` container, allowing you to monitor its activity and troubleshoot issues in real-time.
```bash
m3tal dash logs
```
*Example Usage:*
```bash
# View dashboard logs to debug an issue
m3tal dash logs
```

### `m3tal dash status`
Displays the current status of the `m3tal-dashboard` container (e.g., running, stopped, exited).
```bash
m3tal dash status
```
*Example Usage:*
```bash
# Check if the M3TAL Dashboard is currently running
m3tal dash status
```

---

## 6. Dashboard Access Modes (Critical)

The M3TAL Dashboard has two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`. Understanding these is crucial for accessing your dashboard.

### Mode 1: `local` (Default)
*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** Uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's internal port (`8082`) to the host's `DASHBOARD_PORT` (default `8082`).
*   **Access URL:** `http://HOST_IP:8082` or `http://localhost:8082`
*   **Best for:** LAN-only setups, first-time users, local testing, or scenarios where you don't use Traefik for the dashboard.
*   **Example:** If your M3TAL host has IP `192.168.1.100`, you would access the dashboard at `http://192.168.1.100:8082`.

### Mode 2: `traefik`
*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard container. Traefik then routes requests for `dash.${DOMAIN}` to the dashboard on port `8082` (internal to the Docker network). Requires Traefik to be running via `m3tal up`.
*   **Access URL:** `http://dash.DOMAIN` (e.g., `http://dash.example.com`)
*   **Best for:** Domain-based setups, integrating the dashboard behind a reverse proxy, and environments with multiple services.
*   **Example:** If `DOMAIN` is set to `example.com` in `/etc/m3tal/.env`, you would access the dashboard at `http://dash.example.com`.

**Important:** After changing `DASHBOARD_EXPOSE_MODE`, always run `m3tal dash restart` or `m3tal dash up` for the changes to take effect.

---

## 7. Stack Management

These commands manage your entire collection of Docker Compose stacks defined in `/docker/`. M3TAL uses **Docker Engine + Docker Compose V2**.

**Deployment Lifecycle — Day 2 Operations:**
1.  **Place your compose file:** Put your Docker Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
2.  **Configure environment variables:** Ensure any required variables for your new stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  **Start all stacks:** Run `m3tal up` to deploy your new stack along with existing ones.

### `m3tal up`
Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This brings up all your configured Docker services in detached mode.
```bash
m3tal up
```
*Example Usage:*
```bash
# Deploy all services including Traefik, Cloudflared, and any custom stacks
m3tal up
```

### `m3tal down`
Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This stops and removes all containers, networks, and volumes defined in these compose files.
```bash
m3tal down
```
*Example Usage:*
```bash
# Stop and remove all M3TAL-managed Docker services
m3tal down
```

### `m3tal logs`
Streams aggregated logs from all running Docker containers managed by M3TAL. This provides a consolidated view for troubleshooting.
```bash
m3tal logs
```
*Example Usage:*
```bash
# View real-time logs from all active Docker containers
m3tal logs
```

---

## 8. Systemd Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage it directly using standard `systemctl` commands.

### `systemctl status m3tal-api`
Checks the current status of the M3TAL API daemon service.
```bash
systemctl status m3tal-api
```
*Example Output (abbreviated):*
```
● m3tal-api.service - M3TAL API Daemon
     Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
     Active: active (running) since ...
```

### `journalctl -u m3tal-api -f`
Streams the logs specifically from the M3TAL API daemon service, following new log entries in real-time.
```bash
journalctl -u m3tal-api -f
```
*Example Usage:*
```bash
# Monitor the M3TAL API daemon for errors or activity
journalctl -u m3tal-api -f
```

### Other useful systemctl commands:
```bash
# Restart the M3TAL API daemon
sudo systemctl restart m3tal-api

# Stop the M3TAL API daemon
sudo systemctl stop m3tal-api

# Start the M3TAL API daemon
sudo systemctl start m3tal-api
```

---

## 9. Direct Docker / Compose Fallback

M3TAL automates `docker compose` commands. However, you can always execute `docker compose` commands directly from within the `/docker/` directory as a fallback or for advanced debugging.

**M3TAL's Docker Compose Runtime:**
M3TAL relies on **Docker Engine** and **Docker Compose V2**. The `m3tal up` command translates to `docker compose -f routing-compose.yml -f m3tal-compose.yml -f <other-stacks-compose.yml> up -d` (simplified). It automatically includes overrides like `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` as needed.

Navigate to `/docker/` and use standard `docker compose` commands:

```bash
# Navigate to the M3TAL stack directory
cd /docker/

# Bring up a specific stack (e.g., Traefik & Cloudflared)
# Note: For multiple files, you need to specify them. M3TAL does this automatically.
docker compose -f routing-compose.yml up -d

# Stop a specific stack
docker compose -f routing-compose.yml down

# View logs for a specific service (e.g., traefik within routing-compose.yml)
docker compose -f routing-compose.yml logs -f traefik

# Check the status of all services defined in a compose file
docker compose -f m3tal-compose.yml ps

# Build images for a specific compose file (if needed)
docker compose -f my-custom-stack-compose.yml build

# Stop and remove all services managed by M3TAL's routing stack
docker compose -f routing-compose.yml down
```

---

## 10. Port Map

| Port | Service | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if `routing-compose.yml` is up)     |
| 8080 | M3TAL API daemon (Go)       | Host-local only                             |
| 8081 | Traefik Dashboard           | Host-local only (`127.0.0.1:8081` by default for Traefik's own dashboard) |
| 8082 | M3TAL Dashboard             | Direct port (`local` mode) or via Traefik (`traefik` mode) |

---

This concludes the M3TAL CLI Command Reference. With these commands and the foundational understanding of the M3TAL architecture, you are well-equipped to manage your M3TAL ecosystem effectively. Happy orchestrating!