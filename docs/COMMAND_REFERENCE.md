# M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, my purpose is to equip you with a complete command-line cheat-sheet for managing your M3TAL server. This guide details every command, its function, and provides real-world usage examples to ensure you master your M3TAL environment.

M3TAL provides a unified CLI (`/usr/bin/m3tal`) to interact with its underlying Docker Compose stacks, configuration, and the M3TAL API daemon. For detailed architectural insights, refer to the "M3TAL System Architecture" section at the end of this document.

## 1. Core M3TAL Commands

These commands provide top-level interaction with the M3TAL system.

### `sudo m3tal`

Opens the interactive TUI (Terminal User Interface) Control Center. This provides a numbered menu for navigating common M3TAL operations, ideal for quick management.

**Usage Example:**
```bash
sudo m3tal
```

### `m3tal init`

Generates the primary M3TAL configuration file, `/etc/m3tal/.env`, from system defaults. This command is crucial for the very first installation or if the `.env` file is missing.

**Usage Example:**
```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. This includes verifying Docker connectivity, validating the `/etc/m3tal/.env` file, and checking for essential port availability. It's your first stop for troubleshooting.

**Usage Example:**
```bash
m3tal doctor
```

## 2. Configuration Management

M3TAL's core configuration resides in `/etc/m3tal/.env`. These commands help you manage it.

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating values in `/etc/m3tal/.env`. This is the recommended way to modify your system settings.

**Usage Example:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a specific environment variable within the `/etc/m3tal/.env` file. This command is useful for quick, non-interactive adjustments.

**Usage Example:**
```bash
m3tal config set DOMAIN myhomeserver.tech
```

### `m3tal config get KEY`

Retrieves and displays the value of a specific environment variable from `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal config get DASHBOARD_EXPOSE_MODE
```

### `m3tal config scan`

Lists all known environment variables and their default values across all active Docker Compose stacks managed by M3TAL. This helps identify all possible configuration points.

**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file, showing all active configuration parameters.

**Usage Example:**
```bash
m3tal config list
```

## 3. Dashboard Management

These commands are specifically for controlling the `m3tal-dashboard` container. The dashboard credentials are stored in `/docker/users.json`.

### `m3tal dashpass [username] [password]`

Updates or creates a user password for the M3TAL Dashboard. If `username` and `password` are omitted, it will launch an interactive prompt.

**Usage Examples:**
*   **Interactive:**
    ```bash
    m3tal dashpass
    # Prompt will ask for username and new password
    ```
*   **Direct:**
    ```bash
    m3tal dashpass admin newSecurePa$$w0rd!
    ```

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub, then starts the `m3tal-dashboard` container using the appropriate override (`.local.yml` or `.traefik.yml`) based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal dash up
```

### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container. This is useful after changing configuration settings that affect the dashboard.

**Usage Example:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams the real-time logs from the `m3tal-dashboard` container. Use `Ctrl+C` to stop streaming.

**Usage Example:**
```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` container (e.g., `running`, `exited`).

**Usage Example:**
```bash
m3tal dash status
```

## 4. Stack Management

These commands manage all Docker Compose stacks found in the `/docker/` directory.

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files present in the `/docker/` directory. This starts or recreates all services defined in your M3TAL ecosystem (e.g., Traefik, Cloudflared, and any user-added stacks).

**Usage Example:**
```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This stops and removes all containers, networks, and volumes defined by your M3TAL stacks.

**Usage Example:**
```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running Docker containers managed by M3TAL's stacks. Use `Ctrl+C` to stop streaming.

**Usage Example:**
```bash
m3tal logs
```

## 5. Systemd Service Management

The M3TAL API daemon runs as a systemd service (`m3tal-api.service`). These commands interact directly with systemd to manage the API daemon.

### `systemctl status m3tal-api`

Checks the current status of the M3TAL API daemon, indicating if it's running, exited, or failed.

**Usage Example:**
```bash
systemctl status m3tal-api
```

### `systemctl restart m3tal-api`

Restarts the M3TAL API daemon. This is often necessary after making changes to `/etc/m3tal/.env` that the API depends on.

**Usage Example:**
```bash
sudo systemctl restart m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the real-time logs from the M3TAL API daemon. This is invaluable for debugging issues related to the API or its interaction with Docker.

**Usage Example:**
```bash
sudo journalctl -u m3tal-api -f
```

## 6. Direct Docker / Compose Fallback

M3TAL uses **Docker Engine + Docker Compose V2** under the hood. While the `m3tal` CLI abstracts most operations, you can always use direct `docker compose` commands as a fallback or for advanced debugging. All M3TAL-managed compose files reside in `/docker/`.

**Listing all compose files:**
```bash
ls /docker/*.yml
```
*Expected output might include:*
```
/docker/m3tal-compose.local.yml
/docker/m3tal-compose.traefik.yml
/docker/m3tal-compose.yml
/docker/routing-compose.yml
```

**Equivalent of `m3tal up` (starting all stacks):**
This involves specifying all relevant compose files.
```bash
docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml up -d
```
*(Note: The actual compose files used by `m3tal up` will depend on your specific setup and any user-added stacks.)*

**Equivalent of `m3tal down` (stopping all stacks):**
```bash
docker compose -f /docker/routing-compose.yml -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml down
```

**Checking a specific container's logs (e.g., Traefik):**
```bash
docker logs -f traefik
```

**Inspecting a container:**
```bash
docker inspect m3tal-dashboard
```

## M3TAL System Architecture (Ground Truth)

Understanding the underlying architecture is key to effective management.

### Components

*   **CLI binary** (`/usr/bin/m3tal`): The unified Go binary, installed via APT, serving as the single entry point for all operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on port `8080`. It manages Docker interactions, the internal state database, and provides API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask Docker container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It uses a file provider for dynamic routing and automatically discovers services via Docker labels.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for secure, zero-config internet access to your services.

### Filesystem Contract

| Path                         | Purpose                                                              |
| :--------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`            | The primary configuration file, managed by `m3tal config wizard`.    |
| `/var/lib/m3tal/state.db`    | The SQLite state database, automatically created by the API daemon.  |
| `/opt/m3tal/stack/`          | The canonical directory for all Docker Compose files and Traefik config. |
| `/docker`                    | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`         | The credential store for the M3TAL Dashboard, managed by `m3tal dashpass`. |
| `/docker/dynamic/`           | Directory for Traefik dynamic configuration files (e.g., `api.yml`). |

### Dashboard Access — TWO MODES

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

**Mode 1: `local` (Default)**

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** Uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's port `8082` to the host's `DASHBOARD_PORT` (defaulting to `8082`).
*   **Access:** Via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik required. Works immediately on a local network.
*   **Best for:** LAN-only setups, initial setup, or local testing.

**Mode 2: `traefik`**

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override, which adds specific Traefik labels to the dashboard container. Traefik (if running via `m3tal up`) then routes `dash.${DOMAIN}` to the dashboard container on its internal port `8082`.
*   **Access:** Via `http://dash.YOUR_DOMAIN` (e.g., `http://dash.myhomeserver.tech`).
*   **Requirements:** Traefik must be running as part of your `m3tal up` stacks. Requires a domain (or `localhost`).
*   **Best for:** Domain-based access, integrating with other services behind a reverse proxy.

### Docker / Compose Runtime

*   M3TAL relies on **Docker Engine** and **Docker Compose V2** as hard dependencies.
*   The `m3tal up` command executes `docker compose` across all `*-compose.yml` files found in the `/docker/` directory.
*   The `m3tal dash up` command is specialized:
    1.  It downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub.
    2.  It reads `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
    3.  It then starts the dashboard container using the base compose file plus the appropriate override (`.local.yml` or `.traefik.yml`).
*   To add new services, simply place your `my-stack-compose.yml` file into `/docker/`. These will be picked up by `m3tal up`.

### Deployment Lifecycle — Day 2 Operations

To deploy a new application stack:
1.  Create your Docker Compose file (e.g., `my-app-compose.yml`) and place it in the `/docker/` directory.
2.  Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to start your new stack alongside all other M3TAL services.

### Traefik Routing Architecture

Traefik is deployed via `routing-compose.yml` and serves as the central reverse proxy:
*   It binds to port `80` on the host, acting as the HTTP entry point (`web` entryPoint).
*   It automatically discovers services via Docker labels on containers within the `proxy` network.
*   It loads dynamic configuration from `/docker/dynamic/` via a file provider, enabling hot-reloading of routing rules.
*   **API Routing:** `api.${DOMAIN}` is routed to the M3TAL API daemon (`http://host.docker.internal:8080`) via a static entry in `/docker/dynamic/api.yml`.
*   **Dashboard Routing:** `dash.${DOMAIN}` is routed to the `m3tal-dashboard` container (on its internal port `8082`) using Traefik labels embedded in `m3tal-compose.traefik.yml` (only active when `DASHBOARD_EXPOSE_MODE=traefik`).

### Port Map

| Port | Service                    | Access                                     |
| :--- | :------------------------- | :----------------------------------------- |
| 80   | Traefik HTTP entry point   | Public (when `DASHBOARD_EXPOSE_MODE=traefik` or other services exposed) |
| 8080 | M3TAL API daemon (Go)      | Host-local only                            |
| 8081 | Traefik dashboard (admin UI) | Host-local only (e.g., `http://localhost:8081`) |
| 8082 | M3TAL Dashboard            | Direct port (when `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (when `DASHBOARD_EXPOSE_MODE=traefik`) |

## APT Installation

To install or update the M3TAL CLI and systemd service:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```