# M3TAL Ecosystem Documentation

As DocSmith, the M3TAL Ecosystem Documentation Architect, this README provides a technical overview of the M3TAL system architecture, deployment, and management.

---

## Table of Contents

1.  [Overview](#overview)
2.  [Installation](#installation)
3.  [Core Components](#core-components)
4.  [Filesystem Contract](#filesystem-contract)
5.  [Runtime Environment](#runtime-environment)
6.  [Deployment Lifecycle](#deployment-lifecycle)
7.  [Dashboard Access](#dashboard-access)
8.  [Traefik Gateway](#traefik-gateway)
9.  [Service Management](#service-management)
10. [Port Map](#port-map)
11. [Firewall Configuration](#firewall-configuration)
12. [Quick Demo](#quick-demo)

---

## 1. Overview

M3TAL is an opinionated Docker Compose orchestration system. It provides a CLI, an API daemon, and a web dashboard for managing Docker container stacks. The system relies on Docker Engine and Docker Compose V2 for container lifecycle management and Traefik for application routing.

## 2. Installation

M3TAL is installed via an APT repository. Docker Engine and Docker Compose V2 are prerequisites for the M3TAL system.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Core Components

The M3TAL ecosystem consists of the following components:

*   **CLI binary (`/usr/bin/m3tal`)**: A Go binary providing a single entry point for system operations.
*   **API daemon (`m3tal-api.service`)**: A Go binary running as a systemd service, listening on port 8080. It manages Docker interactions, maintains the state database, and exposes API routes.
*   **Dashboard container (`m3tal-dashboard`)**: A Python/Flask application running in a Docker container. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway (`routing-compose.yml`)**: A Docker container acting as a reverse proxy. It exposes services via domain names on port 80 and uses a file provider for dynamic routing configuration.
*   **Cloudflared (`routing-compose.yml`)**: An optional Docker container for establishing Cloudflare tunnels, enabling zero-configuration internet access to services.

## 4. Filesystem Contract

The following paths are critical for M3TAL operations:

| Path                        | Purpose                                                |
| :-------------------------- | :----------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary system configuration file, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database, auto-created by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL-managed compose files and Traefik dynamic configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for placing and managing Docker Compose stacks. |
| `/docker/users.json`        | Dashboard credential store, managed by `m3tal dashpass`. |

## 5. Runtime Environment

M3TAL uses **Docker Engine + Docker Compose V2** as its underlying container orchestration and management layer. These are fundamental dependencies for the system to function.

The `m3tal` CLI interacts directly with the Docker daemon and orchestrates services using `docker compose` commands.

## 6. Deployment Lifecycle

### Stack Management

The `m3tal up` command executes `docker compose` across all `*-compose.yml` files present in the `/docker/` directory. This brings up or updates all defined container stacks.

The `m3tal dash up` command specifically manages the `m3tal-dashboard` container. It performs the following actions:
1.  Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` files from the official GitHub repository.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the dashboard container using the appropriate compose override file based on the configured exposure mode.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Create or place your Docker Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these variables.
3.  Execute `m3tal up` to start all stacks, including your newly added stack.

## 7. Dashboard Access

The M3TAL dashboard supports two distinct access modes, configured via the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Mode 1: Local (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: Traefik is **not** required for this mode. It provides out-of-the-box access for first-time users or setups on local networks.
*   **Use Case**: Recommended for initial setup, LAN-only environments, and local testing.
*   **Note for new users**: By default, after installation and `m3tal dash up`, the M3TAL Dashboard will be available at `http://HOST_IP:8082`, not via a domain name.

### Mode 2: Traefik

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: This mode uses the `m3tal-compose.traefik.yml` override file. This file adds specific Traefik labels to the dashboard service, enabling Traefik to route `dash.${DOMAIN}` to the dashboard container on its internal port 8082.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN`.
*   **Requirements**: The Traefik gateway must be running (e.g., via `m3tal up`) for this mode to function.
*   **Use Case**: Suitable for domain-based deployments where multiple services are exposed through Traefik.

## 8. Traefik Gateway

Traefik operates as a Docker container, defined in `routing-compose.yml`. Its primary functions include:

*   **Port Binding**: Traefik binds to port 80 on the host system, serving as the HTTP entry point for all routed services.
*   **Service Discovery**: It automatically discovers services by inspecting Docker container labels within the `proxy` network.
*   **Dynamic Configuration**: Traefik loads additional routing rules from `/docker/dynamic/` using a file provider. This configuration directory supports hot-reloading.
*   **API Daemon Routing**: The `api.DOMAIN` route is configured dynamically via `/docker/dynamic/api.yml` to point to the M3TAL API daemon running on the host at `http://host.docker.internal:8080`.
*   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, `dash.DOMAIN` is routed to the `m3tal-dashboard` container via Traefik labels defined in `m3tal-compose.traefik.yml`.

## 9. Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Standard systemctl commands apply:

*   **Check Status**: `sudo systemctl status m3tal-api`
*   **Restart Service**: `sudo systemctl restart m3tal-api`
*   **View Logs**: `sudo journalctl -u m3tal-api -f` (for following logs in real-time)

## 10. Port Map

The following network ports are used by M3TAL and its components:

| Port | Service               | Access Context                                           |
| :--- | :-------------------- | :------------------------------------------------------- |
| 80   | Traefik HTTP entry    | Publicly accessible (when Traefik is used for routing)   |
| 8080 | M3TAL API daemon (Go) | Host-local access only (internal communication)          |
| 8081 | Traefik dashboard     | Host-local only (internal Traefik metrics/dashboard)     |
| 8082 | M3TAL Dashboard       | Direct port binding (local mode) or via Traefik (traefik mode) |

## 11. Firewall Configuration

If you are using a firewall (e.g., UFW or iptables) on your host system, ensure that port 80 is open to allow external access to services routed by Traefik.

Example for UFW:
`sudo ufw allow 80/tcp`

## 12. Quick Demo

Follow these steps for a rapid demonstration of M3TAL:

1.  **Install M3TAL**:
    ```bash
    # 1. Add the GPG signing key
    curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

    # 2. Add the APT repository
    echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

    # 3. Install
    sudo apt update && sudo apt install -y m3tal
    ```

2.  **Run Configuration Wizard**:
    ```bash
    m3tal config wizard
    ```
    Follow the prompts. Accept defaults for `DASHBOARD_EXPOSE_MODE` (which is `local`).

3.  **Start M3TAL Dashboard**:
    ```bash
    m3tal dash up
    ```
    This will download and start the dashboard container in local mode.

4.  **Access Dashboard**:
    Open your web browser and navigate to `http://HOST_IP:8082` (replace `HOST_IP` with the IP address of your M3TAL server). You will be presented with the M3TAL dashboard login screen.