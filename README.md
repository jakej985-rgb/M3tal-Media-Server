# M3TAL Ecosystem Documentation

M3TAL is a Docker Compose management system designed for self-hosting environments. It provides a unified CLI, an API daemon for Docker orchestration, and a web dashboard for service management. All services are deployed using Docker Engine and Docker Compose V2.

## Core Components

*   **CLI binary** (`/usr/bin/m3tal`): A single Go binary installed via APT, serving as the entry point for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on host port 8080. It manages Docker interactions, maintains the M3TAL state database, and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container, internally on port 8082. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy. It exposes services by domain name on host port 80 and uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing zero-configuration internet access for exposed services.

## Filesystem Contract

The following paths are central to the M3TAL system:

| Path                        | Purpose                                                      |
| :-------------------------- | :----------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.       |
| `/opt/m3tal/stack/`         | Canonical directory for Docker Compose stack files and Traefik dynamic configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.     |

## Runtime Environment

M3TAL relies on **Docker Engine** and **Docker Compose V2** for container orchestration. These are hard dependencies for the system to function.

*   The `m3tal up` command executes `docker compose` across all `*-compose.yml` files found in the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the dashboard container. It downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub. It then reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env` and starts the dashboard container with the appropriate compose override file.

## Installation

Ensure Docker Engine and Docker Compose V2 are installed on your system before proceeding with M3TAL installation.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Configuration

M3TAL uses `/etc/m3tal/.env` for all system-wide configuration. This file contains environment variables that control the behavior of the API daemon, dashboard, and other core components. The `.env` file is managed interactively via the `m3tal config wizard` command. Individual settings can be updated with `m3tal config set KEY value`.

## Deployment Lifecycle

M3TAL manages services as Docker Compose "stacks". All user-defined compose files are placed in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`).

### Bringing Up Stacks

The `m3tal up` command brings up all defined stacks by invoking `docker compose` on all `*-compose.yml` files present in `/docker/`.

### Adding a New Stack

To deploy a new Docker Compose stack with M3TAL:

1.  Place your compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to deploy the new stack alongside existing ones.

## Traefik Gateway

Traefik is deployed as a Docker container via `routing-compose.yml`. Its primary function is to act as a reverse proxy for all domain-based services.

*   Traefik binds host port 80 as its HTTP entry point.
*   It automatically discovers services by reading Docker labels applied to service definitions within compose files.
*   Dynamic routing configurations, such as for the M3TAL API, are loaded from YAML files placed in `/docker/dynamic/` (which maps to `/etc/traefik/dynamic` inside the Traefik container). Traefik hot-reloads these configurations on changes.
*   The M3TAL API daemon is exposed at `api.DOMAIN`. Traefik routes traffic for this domain to `http://host.docker.internal:8080` (the Go API daemon) via a dynamic configuration file in `/docker/dynamic/api.yml`.
*   When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, the M3TAL dashboard is exposed via Traefik at `dash.DOMAIN`, routing to the dashboard container based on its Traefik labels.

**Firewall Note:** Ensure host port 80 is allowed through your firewall (e.g., UFW, iptables) for Traefik to function as a public entry point.

## Dashboard Access

The M3TAL dashboard supports two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### 1. Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port binding, exposing the dashboard container's internal port 8082 directly on the host at `${DASHBOARD_PORT:-8082}`.
*   **Access:** Navigate to `http://HOST_IP:8082` or `http://localhost:8082` in your web browser.
*   **Requirements:** No Traefik instance is required for this mode. It works out-of-the-box on a local network or home server.
*   **Recommendation:** This is the default setup for initial installation, first-time users, and LAN-only deployments.

### 2. Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode uses the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard service. Traefik then routes requests for `dash.DOMAIN` to the dashboard container's internal port 8082.
*   **Access:** Navigate to `http://dash.DOMAIN` in your web browser (e.g., `http://dash.example.com`).
*   **Requirements:** Traefik must be running via `m3tal up` for this mode to function.
*   **Recommendation:** This mode is suitable for setups requiring domain-based access and routing multiple services behind Traefik.

A new user performing a default installation will access the dashboard directly at `http://HOST_IP:8082`, not via a domain name.

## Service Management

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service. Its status and logs can be monitored using standard systemctl commands:

*   **Check Status:** `sudo systemctl status m3tal-api`
*   **Restart Service:** `sudo systemctl restart m3tal-api`
*   **View Live Logs:** `journalctl -u m3tal-api -f`

## Port Map

| Port | Service                                      | Access                                               |
| :--- | :------------------------------------------- | :--------------------------------------------------- |
| 80   | Traefik HTTP entry point                     | Public (if Traefik is enabled and configured)        |
| 8080 | M3TAL API daemon (Go)                        | Host-local only (internal communication)             |
| 8081 | Traefik dashboard (Traefik's own UI)         | Host-local only (internal management)                |
| 8082 | M3TAL Dashboard (Python/Flask)               | Direct port (local mode) or via Traefik (Traefik mode) |

## Quick Demo

Follow these steps to quickly get M3TAL running and access the dashboard:

1.  **Install M3TAL:**
    ```bash
    # 1. Add the GPG signing key
    curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

    # 2. Add the APT repository
    echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list | sudo tee /etc/apt/sources.list.d/m3tal.list

    # 3. Install
    sudo apt update && sudo apt install -y m3tal
    ```
2.  **Run the configuration wizard:**
    ```bash
    m3tal config wizard
    ```
    Follow the prompts. At minimum, set a strong `DASHBOARD_SECRET` and `ADMIN_PASSWORD`.
3.  **Bring up the M3TAL dashboard:**
    ```bash
    m3tal dash up
    ```
4.  **Open your browser:**
    Navigate to `http://HOST_IP:8082` (replace `HOST_IP` with your server's actual IP address).
    Log in with the username `admin` and the `ADMIN_PASSWORD` you configured in step 2.