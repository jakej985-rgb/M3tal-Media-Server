Hey there, fellow M3TALhead! DocSmith here, your M3TAL Ecosystem Documentation Architect. You're looking for the ultimate command-line cheat sheet to master your M3TAL rig? You've come to the right place. This document is your comprehensive guide to the `m3tal` CLI, covering everything from initial setup to day-to-day operations and advanced configuration. Let's get to it!

---

# `docs/COMMAND_REFERENCE.md` — M3TAL CLI Cheat Sheet

The `m3tal` command-line interface is the primary tool for interacting with and managing your M3TAL ecosystem. It simplifies complex Docker Compose operations, system configuration, and service management into a unified, user-friendly experience.

## M3TAL APT Installation

For a fresh install or to upgrade your M3TAL CLI and API daemon, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## Core CLI Commands

### `sudo m3tal`

**Description:** Opens the interactive TUI (Terminal User Interface) Control Center. This provides a user-friendly, numbered menu for common operations without needing to remember specific commands. It's often the fastest way to perform routine tasks.
**Usage Example:**
```bash
sudo m3tal
```

### `m3tal init`

**Description:** Initializes the M3TAL configuration. This command generates the primary configuration file, `/etc/m3tal/.env`, from system defaults. It is essential to run this on a first installation before configuring your system.
**Usage Example:**
```bash
m3tal init
```

### `m3tal doctor`

**Description:** Performs a pre-flight health check of your M3TAL system. It verifies Docker connectivity, checks the validity of `/etc/m3tal/.env`, and ensures that essential ports are available, helping you diagnose common issues quickly.
**Usage Example:**
```bash
m3tal doctor
```

---

## Configuration Management

M3TAL uses `/etc/m3tal/.env` as its central configuration file. These commands help you manage its contents.

### `m3tal config wizard`

**Description:** Launches an interactive wizard to guide you through configuring or updating the essential variables in `/etc/m3tal/.env`. This is the recommended way to set up your environment after `m3tal init`.
**Usage Example:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

**Description:** Sets or updates a single environment variable in `/etc/m3tal/.env`. Useful for quick, targeted changes. Remember to restart relevant services or stacks if you change critical variables.
**Usage Examples:**
```bash
m3tal config set DOMAIN "my-metal-server.com"
m3tal config set DASHBOARD_EXPOSE_MODE traefik
m3tal config set PUID 1001
```

### `m3tal config get KEY`

**Description:** Retrieves and displays the value of a specific environment variable from `/etc/m3tal/.env`.
**Usage Examples:**
```bash
m3tal config get DOMAIN
m3tal config get DASHBOARD_EXPOSE_MODE
m3tal config get PUID
```

### `m3tal config scan`

**Description:** Lists all recognized environment variables across all active and potential M3TAL stacks, showing their current values from `/etc/m3tal/.env` (or defaults if not set). This provides a comprehensive overview of your system's configuration.
**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`

**Description:** Displays the current contents of the `/etc/m3tal/.env` file. This is useful for reviewing your active configuration.
**Usage Example:**
```bash
m3tal config list
```

---

## Dashboard Management

The M3TAL Dashboard provides a web-based UI for monitoring and controlling your ecosystem. These commands specifically manage the dashboard container.

### `m3tal dashpass [username] [password]`

**Description:** Manages user credentials for the M3TAL dashboard. If `username` and `password` are omitted, it will launch an interactive prompt to set or reset credentials. This data is stored in `/docker/users.json`.
**Usage Examples:**
*   **Interactive mode:**
    ```bash
    m3tal dashpass
    ```
*   **Direct mode:**
    ```bash
    m3tal dashpass admin MyNewSecurePassword123
    ```

### `m3tal dash up`

**Description:** Pulls the latest M3TAL Dashboard Docker Compose configuration from GitHub, then starts the dashboard container. It intelligently applies the correct Docker Compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.
**Usage Example:**
```bash
m3tal dash up
```

### `m3tal dash down`

**Description:** Stops and removes the M3TAL Dashboard container and its associated network resources.
**Usage Example:**
```bash
m3tal dash down
```

### `m3tal dash restart`

**Description:** Restarts the M3TAL Dashboard container. This is useful after making configuration changes or if the dashboard becomes unresponsive.
**Usage Example:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

**Description:** Streams the real-time logs from the M3TAL Dashboard container. Useful for debugging and monitoring dashboard activity.
**Usage Example:**
```bash
m3tal dash logs
```

### `m3tal dash status`

**Description:** Shows the current status of the M3TAL Dashboard container (e.g., running, stopped, restarting).
**Usage Example:**
```bash
m3tal dash status
```

---

## Stack Management

These commands manage all Docker Compose stacks defined in `/docker/`.

### `m3tal up`

**Description:** Orchestrates the startup of all Docker Compose stacks within your M3TAL ecosystem. It runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This brings all your services online.
**Usage Example:**
```bash
m3tal up
```

### `m3tal down`

**Description:** Stops and removes all Docker Compose stacks. It runs `docker compose down` across all `*-compose.yml` files in `/docker/`, gracefully shutting down your entire M3TAL environment.
**Usage Example:**
```bash
m3tal down
```

### `m3tal logs`

**Description:** Streams aggregated logs from all currently running Docker containers managed by M3TAL. This provides a consolidated view of your entire ecosystem's output, essential for system-wide monitoring and debugging.
**Usage Example:**
```bash
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon, `m3tal-api`, runs as a systemd service. Here are the essential commands for managing it.

### `systemctl status m3tal-api`

**Description:** Checks the current status of the M3TAL API daemon, including whether it's active, running, and its most recent log entries.
**Usage Example:**
```bash
systemctl status m3tal-api
```

### `systemctl restart m3tal-api`

**Description:** Restarts the M3TAL API daemon. This is often necessary after making changes to `/etc/m3tal/.env` that affect the API's behavior (e.g., `API_TOKEN`, `LOG_LEVEL`).
**Usage Example:**
```bash
systemctl restart m3tal-api
```

### `journalctl -u m3tal-api -f`

**Description:** Streams real-time logs from the M3TAL API daemon. This is invaluable for detailed troubleshooting and monitoring of the API's operations.
**Usage Example:**
```bash
journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Fallback

M3TAL abstracts away direct `docker compose` commands, but understanding how to use them can be useful for advanced debugging or specific operations. M3TAL relies on **Docker Engine** and **Docker Compose V2**.

All M3TAL Docker Compose files, and any user-added stack files, reside in `/docker/` (which is a symlink to `/opt/m3tal/stack/`).

### Common Docker Compose Commands:

*   **Start a specific stack:**
    ```bash
    docker compose -f /docker/routing-compose.yml up -d
    ```
*   **Stop a specific stack:**
    ```bash
    docker compose -f /docker/ollama-compose.yml down
    ```
*   **View logs for a specific stack:**
    ```bash
    docker compose -f /docker/m3tal-compose.yml logs -f
    ```
*   **Rebuild and restart a specific service (e.g., after custom image changes):**
    ```bash
    docker compose -f /docker/m3tal-compose.yml up -d --build m3tal-dashboard
    ```
*   **View all running containers:**
    ```bash
    docker ps
    ```
*   **View all Docker Compose projects/stacks (v2):**
    ```bash
    docker compose ls
    ```

---

## M3TAL System Architecture Overview

The M3TAL ecosystem is designed for robustness and ease of management, built upon several interconnected components:

*   **CLI binary (`/usr/bin/m3tal`):** The unified Go binary that serves as the single entry point for all user-facing operations, installed via APT.
*   **API daemon (`m3tal-api.service`):** A Go binary running as a `systemd` service, listening on port 8080. It manages Docker interactions, maintains the SQLite state database, and exposes RESTful API routes for the dashboard and other clients.
*   **Dashboard container (`m3tal-dashboard`):** A Python/Flask application running in a Docker container (internally on port 8082). It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik Gateway (`routing-compose.yml`):** A robust reverse proxy container, deployed as a core stack component. It exposes services by domain name on host port 80 (and 443 with HTTPS) and uses a file provider for dynamic routing configuration.
*   **Cloudflared (`routing-compose.yml`):** An optional Cloudflare tunnel container, also part of the routing stack, enabling secure, zero-config internet access for your services without opening firewall ports.

### Filesystem Contract

The following paths are critical to the M3TAL ecosystem:

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for the entire M3TAL system. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite database storing M3TAL's internal state. Auto-created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose files and Traefik dynamic configuration. |
| `/docker`                   | A symbolic link to `/opt/m3tal/stack/`. This is the user-facing path for placing custom `*-compose.yml` files. |
| `/docker/users.json`        | Credential store for M3TAL Dashboard users. Managed by `m3tal dashpass`. |
| `/docker/dynamic/`          | Directory for Traefik dynamic configuration files (e.g., `api.yml`). These files are hot-reloaded by Traefik. |

### Dashboard Access — TWO MODES

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   Uses the Docker Compose override file: `m3tal-compose.local.yml`.
*   Directly binds the dashboard container's internal port (`8082`) to a host port (default `8082`). This means the dashboard is accessible directly via your server's IP address.
*   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requires no Traefik:** Ideal for LAN-only setups, first-time users, or local development environments where a domain name is not yet configured.

    ```yaml
    # m3tal-compose.local.yml
    services:
      m3tal-dashboard:
        ports:
          - "${DASHBOARD_PORT:-8082}:8082"
    ```

#### Mode 2: `traefik`

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   Uses the Docker Compose override file: `m3tal-compose.traefik.yml`.
*   Adds Traefik labels to the dashboard container, instructing Traefik to route incoming requests for `dash.${DOMAIN}` (e.g., `dash.my-metal-server.com`) to the dashboard container on its internal port `8082`.
*   **Access via:** `http://dash.DOMAIN` (e.g., `http://dash.my-metal-server.com`).
*   **Requires Traefik:** The Traefik gateway must be running (`m3tal up` will start it as part of `routing-compose.yml`). Best for domain-based setups and integrating the dashboard behind a reverse proxy alongside other services.

    ```yaml
    # m3tal-compose.traefik.yml
    services:
      m3tal-dashboard:
        labels:
          - "traefik.enable=true"
          - "traefik.http.routers.dashboard.rule=Host(`dash.${DOMAIN:-localhost}`)"
          - "traefik.http.routers.dashboard.entrypoints=web"
          - "traefik.http.services.dashboard.loadbalancer.server.port=8082"
          - "traefik.docker.network=proxy"
    ```

### Docker / Compose Runtime

M3TAL is built on **Docker Engine** and leverages **Docker Compose V2** for container orchestration.

*   The `m3tal up` command orchestrates all Docker Compose services by iterating through all `*-compose.yml` files present in the `/docker/` directory and running `docker compose up -d` for each.
*   The `m3tal dash up` command is specialized for the dashboard:
    1.  It ensures the latest `m3tal-compose.yml` (and its local/traefik overrides) are downloaded from GitHub.
    2.  It dynamically selects the appropriate override file based on `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
    3.  It then starts the dashboard container.
*   **Deployment Lifecycle - Day 2 Operations:** To deploy a new service or stack, simply place your custom `my-new-stack-compose.yml` file into `/docker/`. Ensure all required environment variables for that stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set`). Then, run `m3tal up` to bring your new stack online.

### Traefik Routing Architecture

Traefik acts as the central reverse proxy for your M3TAL ecosystem, deployed via `routing-compose.yml`.

*   **Entry Points:** It binds to port 80 on the host (and optionally 443 for HTTPS) as its primary HTTP entry point (`web` entrypoint).
*   **Service Discovery:** Traefik automatically discovers and configures routes for Docker containers based on their labels.
*   **File Provider:** It also loads dynamic configuration from `/docker/dynamic/` (specifically mapped to `/etc/traefik/dynamic` inside the container), allowing for hot-reloadable custom routing rules.
*   **API Daemon Routing:** Traefik routes requests for `api.DOMAIN` to the M3TAL API daemon (running on the host at `http://host.docker.internal:8080`) via a dynamic configuration file like `dynamic/api.yml`.
*   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, requests for `dash.DOMAIN` are routed to the dashboard container using its Docker labels.

**Traefik Static Configuration (from `traefik.yml`):**

```yaml
entryPoints:
  web:
    address: ":80" # Binds to host port 80

providers:
  docker:
    exposedByDefault: false # Only expose containers with specific Traefik labels
    network: proxy # Connects to the 'proxy' network for container discovery
  file:
    directory: /etc/traefik/dynamic # Watches this directory for dynamic config (mapped from /docker/dynamic)
    watch: true # Hot-reloads config changes
```

**Dynamic API Routing Example (`/docker/dynamic/api.yml`):**

```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)" # Routes requests for api.yourdomain.com
      service: api
      entryPoints:
        - web

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080" # Targets the host's M3TAL API daemon
```

### Port Map

| Port | Service                                  | Access                                     |
| :--- | :--------------------------------------- | :----------------------------------------- |
| 80   | Traefik HTTP entry point                 | Public (when Traefik is running)           |
| 8080 | M3TAL API daemon (Go binary)             | Host-local only                            |
| 8081 | Traefik dashboard                        | Host-local only (e.g., `http://localhost:8081`) |
| 8082 | M3TAL Dashboard container (internal)     | Direct port (local mode) or via Traefik (traefik mode) |
| 443  | Traefik HTTPS entry point                | Public (requires HTTPS setup in Traefik)   |
| 11434| Ollama (example AI stack service)        | Typically exposed directly or via Traefik  |

---

This comprehensive reference should equip you with the knowledge to navigate and master your M3TAL ecosystem. Keep building, and remember: with M3TAL, you're always in control!