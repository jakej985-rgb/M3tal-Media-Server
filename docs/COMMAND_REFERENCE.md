# M3TAL CLI Command Reference

Greetings, M3TAL Operators. I am DocSmith, your M3TAL Ecosystem Documentation Architect. This document serves as your definitive cheat-sheet for interacting with the M3TAL system via its unified command-line interface. Master these commands, and you will command your digital domain with precision and power.

The M3TAL CLI (`/usr/bin/m3tal`) is the single entry point for all system operations, designed for efficiency and clarity.

---

## 1. M3TAL Installation

For new installations or system recovery, follow these steps to deploy the M3TAL CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 2. M3TAL Filesystem Contract

The M3TAL system relies on a well-defined filesystem layout for its operation. Understanding these paths is crucial for maintenance and troubleshooting.

| Path                     | Purpose                                                                                                     |
| :----------------------- | :---------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | **Primary system configuration file.** Contains all environment variables for M3TAL operations. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`| **SQLite state database.** Stores M3TAL API daemon state, service information, and other operational data. Auto-created. |
| `/opt/m3tal/stack/`      | **Canonical stack directory.** This is the core location where Docker Compose files and Traefik dynamic configurations are stored. |
| `/docker`                | **Symlink to `/opt/m3tal/stack/`.** This is the user-facing path for placing and managing all Docker Compose stack definitions. |
| `/docker/users.json`     | **Dashboard credential store.** Stores hashed passwords for M3TAL Dashboard access. Managed by `m3tal dashpass`. |

---

## 3. M3TAL CLI Commands

This section details every primary command and its subcommands, complete with concrete usage examples.

### 3.1. Interactive TUI Control Center

The M3TAL Control Center offers a guided, interactive menu for common operations.

*   **`sudo m3tal`**
    Opens the interactive Text-User Interface (TUI) Control Center. This provides a numbered menu for common tasks like starting/stopping stacks, checking status, configuring the system, and managing the dashboard. Elevated privileges (`sudo`) are required as many underlying operations interact with Docker.

    ```bash
    # Launch the interactive M3TAL Control Center to manage your system
    sudo m3tal
    ```

### 3.2. System Initialization and Health

Commands for setting up the M3TAL environment and ensuring its operational readiness.

*   **`m3tal init`**
    Generates the primary configuration file, `/etc/m3tal/.env`, from system defaults. This command is essential on the first installation or if `/etc/m3tal/.env` is missing, providing a baseline for your M3TAL environment.

    ```bash
    # Initialize the M3TAL environment by generating the default .env file
    m3tal init
    ```

*   **`m3tal doctor`**
    Performs a pre-flight health check of the M3TAL system. It verifies critical components such as Docker connectivity, the validity and existence of `/etc/m3tal/.env`, and ensures that required ports (e.g., 80, 8080, 8082) are not already in use by other processes.

    ```bash
    # Run a diagnostic check on the M3TAL system to identify potential issues
    m3tal doctor
    ```

### 3.3. Configuration Management (`m3tal config`)

These commands allow you to inspect and modify the core M3TAL configuration file, `/etc/m3tal/.env`.

*   **`m3tal config wizard`**
    Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for making changes, as it provides explanations and validates inputs.

    ```bash
    # Start the interactive configuration wizard to adjust M3TAL settings
    m3tal config wizard
    ```

*   **`m3tal config set KEY VALUE`**
    Sets a specific environment variable (`KEY`) to a new `VALUE` within `/etc/m3tal/.env`. Use this for quick, direct modifications. After changing a value, remember that active Docker containers might need to be restarted (`m3tal up` or `m3tal dash restart`) to pick up the new setting.

    ```bash
    # Set the dashboard exposure mode to 'traefik'
    m3tal config set DASHBOARD_EXPOSE_MODE traefik

    # Update the primary domain used by Traefik and other services
    m3tal config set DOMAIN myhomelab.net
    ```

*   **`m3tal config get KEY`**
    Retrieves and displays the current value of a specified environment variable (`KEY`) from `/etc/m3tal/.env`.

    ```bash
    # Get the currently configured M3TAL Dashboard port
    m3tal config get DASHBOARD_PORT

    # Display the set primary domain
    m3tal config get DOMAIN
    ```

*   **`m3tal config scan`**
    Lists all environment variables known to the M3TAL system, including their defaults and current values, across all defined Docker Compose stacks. This provides a comprehensive overview of your entire M3TAL ecosystem's variable landscape.

    ```bash
    # Scan and list all environment variables recognized by M3TAL
    m3tal config scan
    ```

*   **`m3tal config list`**
    Displays the entire contents of the current M3TAL environment file (`/etc/m3tal/.env`). This is useful for reviewing all active settings at once.

    ```bash
    # Display the full contents of the M3TAL configuration file
    m3tal config list
    ```

### 3.4. Dashboard Management (`m3tal dash`)

Commands specifically designed for managing the M3TAL Dashboard container and its access.

**Dashboard Access Modes: Critical Information**
The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

*   **Mode 1: `local` (Default)**
    *   `DASHBOARD_EXPOSE_MODE=local`
    *   The dashboard container's port `8082` is directly mapped to the host's `8082` (or custom `DASHBOARD_PORT`).
    *   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082`
    *   No Traefik reverse proxy is involved. Ideal for LAN-only setups, initial deployment, or local testing.

*   **Mode 2: `traefik`**
    *   `DASHBOARD_EXPOSE_MODE=traefik`
    *   The dashboard is exposed via Traefik, routing `dash.${DOMAIN}` (e.g., `dash.myhomelab.net`) to the container.
    *   **Access via:** `http://dash.YOUR_DOMAIN` (Requires Traefik to be running via `m3tal up`).
    *   Best for domain-based access and integration into a larger Traefik-managed environment.

*   **`m3tal dashpass [username] [password]`**
    Updates the password for a specified dashboard user. If `username` and `password` are omitted, the command will prompt you interactively for the details. User credentials are stored in `/docker/users.json`.

    ```bash
    # Interactively update the password for the 'admin' user
    m3tal dashpass admin

    # Set the password for user 'operator' directly (use with caution in scripts)
    m3tal dashpass operator SecureM3talPass123!
    ```

*   **`m3tal dash up`**
    Ensures the M3TAL Dashboard is running. This command first pulls the latest dashboard compose configuration files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from the official GitHub repository. It then starts the `m3tal-dashboard` container, dynamically applying the correct override based on the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

    ```bash
    # Start or update the M3TAL Dashboard container, pulling the latest configuration
    m3tal dash up
    ```

*   **`m3tal dash down`**
    Stops and removes the `m3tal-dashboard` container. This does not remove its associated data volumes.

    ```bash
    # Stop the M3TAL Dashboard container
    m3tal dash down
    ```

*   **`m3tal dash restart`**
    Restarts the `m3tal-dashboard` container. This is useful after changing dashboard-related environment variables or troubleshooting.

    ```bash
    # Restart the M3TAL Dashboard
    m3tal dash restart
    ```

*   **`m3tal dash logs`**
    Streams the real-time logs from the `m3tal-dashboard` container. This is invaluable for debugging and monitoring dashboard activity.

    ```bash
    # Stream logs from the M3TAL Dashboard
    m3tal dash logs
    ```

*   **`m3tal dash status`**
    Displays the current operational status of the `m3tal-dashboard` container, including its running state, uptime, and port mappings.

    ```bash
    # Check the status of the M3TAL Dashboard container
    m3tal dash status
    ```

### 3.5. Stack Management (`m3tal up`, `m3tal down`, `m3tal logs`)

These commands manage all Docker Compose stacks defined within your `/docker/` directory. M3TAL utilizes **Docker Engine + Docker Compose V2**.

*   **`m3tal up`**
    Brings up all Docker Compose services defined by `*-compose.yml` files located in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`). This command reads your `/etc/m3tal/.env` variables and applies them to all services, starting or updating them as necessary. This is how you deploy new services or apply updates to existing ones.

    ```bash
    # Start or update all Docker Compose stacks found in /docker/, applying .env settings
    m3tal up
    ```

*   **`m3tal down`**
    Stops and removes all Docker Compose services and their associated networks defined by `*-compose.yml` files in `/docker/`. This effectively brings down your entire M3TAL-managed container infrastructure.

    ```bash
    # Stop and remove all running Docker Compose stacks
    m3tal down
    ```

*   **`m3tal logs`**
    Streams aggregated logs from *all* currently running Docker Compose services managed by M3TAL. This provides a consolidated view of activity across your entire containerized ecosystem.

    ```bash
    # Stream real-time logs from all active M3TAL-managed containers
    m3tal logs
    ```

---

## 4. Systemd Service Management

The M3TAL API daemon (`m3tal-api`) is a critical component running as a Go binary and managed by `systemd`. It runs on port 8080. Understanding how to interact with its systemd service is crucial for maintenance.

*   **Check API Daemon Status:**
    To verify if the M3TAL API daemon is running and healthy:
    ```bash
    systemctl status m3tal-api
    ```

*   **Restart API Daemon:**
    If you've modified system-level M3TAL configurations or are troubleshooting API connectivity:
    ```bash
    sudo systemctl restart m3tal-api
    ```

*   **Stream API Daemon Logs:**
    For real-time debugging of the M3TAL API's operations:
    ```bash
    journalctl -u m3tal-api -f
    ```

---

## 5. Direct Docker Compose Commands (Fallback)

While `m3tal` streamlines Docker Compose operations, understanding the underlying direct `docker compose` commands can be useful for advanced troubleshooting or when M3TAL's CLI is unavailable for some reason. All M3TAL-managed stacks reside in `/docker/` (symlinked from `/opt/m3tal/stack/`).

*   **Start All Stacks Directly:**
    This command operates on all `*-compose.yml` files in `/docker/`.
    ```bash
    # Navigate to the stacks directory
    cd /docker/

    # Start all services defined by compose files in this directory (example listing common files)
    sudo docker compose -f routing-compose.yml -f m3tal-compose.yml -f ollama-compose.yml up -d
    ```
    *Note: You would need to explicitly list all your compose files. `m3tal up` handles this aggregation automatically.*

*   **Stop All Stacks Directly:**
    ```bash
    # Navigate to the stacks directory
    cd /docker/

    # Stop all services defined by compose files
    sudo docker compose -f routing-compose.yml -f m3tal-compose.yml -f ollama-compose.yml down
    ```

*   **Stream Logs from a Specific Stack/Service:**
    ```bash
    # Stream logs from the m3tal-dashboard service within m3tal-compose.yml
    sudo docker compose -f /docker/m3tal-compose.yml logs -f m3tal-dashboard
    ```
    *Note: `m3tal logs` aggregates logs from ALL services, while this targets a specific service within a specific compose file.*

*   **Check Status of a Specific Service:**
    ```bash
    # Check the status of the Traefik container
    sudo docker compose -f /docker/routing-compose.yml ps traefik
    ```

---

## 6. M3TAL System Architecture & Runtime Essentials

The M3TAL ecosystem leverages industry-standard tools to provide a robust and flexible environment.

*   **Docker Engine + Docker Compose V2**: These are fundamental dependencies. M3TAL orchestrates containers using Compose files found in `/docker/`.
*   **`/docker/`**: This directory (a symlink to `/opt/m3tal/stack/`) is your primary workspace for defining and managing your containerized applications. Simply place your `*-compose.yml` files here, and `m3tal up` will integrate them into your ecosystem.
*   **M3TAL API Daemon**: A Go binary (`m3tal-api.service`) running on port 8080, managing Docker interactions, the state database (`/var/lib/m3tal/state.db`), and API routes.
*   **M3TAL Dashboard**: A Python/Flask container (`m3tal-dashboard`) running on internal port 8082, communicating with the API daemon at `http://host.docker.internal:8080`. Its exposure mode (local vs. Traefik) is controlled by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.
*   **Traefik Gateway (`routing-compose.yml`)**: A reverse proxy container exposing services on host port 80. It uses Docker labels and file providers (from `/docker/dynamic/`) for dynamic routing, handling requests for `api.DOMAIN` and, if configured, `dash.DOMAIN`.
*   **Cloudflared (`routing-compose.yml`)**: An optional Cloudflare tunnel for secure, zero-config internet access to your services.

### Traefik Routing Architecture Overview

Traefik acts as the front-door for your domain-based services.
*   It listens on host port 80 (`entryPoints.web`).
*   Automatically discovers services via Docker labels (`providers.docker`).
*   Loads additional dynamic configuration from `/docker/dynamic/` (`providers.file`), enabling custom routing.
*   Example: `api.DOMAIN` is routed to the M3TAL API daemon (Go) at `http://host.docker.internal:8080` via `/docker/dynamic/api.yml`.
*   Example: `dash.DOMAIN` routes to the dashboard container through labels in `m3tal-compose.traefik.yml` when `DASHBOARD_EXPOSE_MODE=traefik`.

---

## 7. M3TAL Port Map

A summary of the key ports used by the M3TAL system:

| Port | Service               | Access Context                                     |
| :--- | :-------------------- | :------------------------------------------------- |
| 80   | Traefik HTTP entry point | Public (when `routing-compose.yml` is active)        |
| 8080 | M3TAL API daemon (Go) | Host-local (internal communication)                |
| 8081 | Traefik dashboard     | Host-local only (internal management interface)    |
| 8082 | M3TAL Dashboard       | Direct port (local mode) or via Traefik (traefik mode) |

---

This comprehensive guide should equip you with the knowledge to effectively manage your M3TAL ecosystem. Remember, precision in command leads to mastery of your machine. Should you require further guidance, consult the M3TAL documentation archives.

DocSmith, signing off.