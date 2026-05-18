# M3TAL Ecosystem Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL relies on Docker Compose V2 for its internal orchestration and uses Docker Engine + Docker Compose V2 internally.

## Installation

M3TAL is distributed via an APT repository. Execute the following commands to install the `m3tal` CLI binary and its associated components.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

**Firewall Note:** If you plan to use the Traefik Gateway for domain-based access, ensure that host port `80` is open in your firewall (e.g., `ufw allow 80/tcp`).

## Filesystem Contract

The M3TAL system adheres to the following filesystem structure:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik configuration. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Core Components

The M3TAL system consists of the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A Go binary providing a unified entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container that operates internally on port `8082`. It communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A containerized reverse proxy that exposes services via domain names on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare Tunnel container for establishing zero-configuration internet access to services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth for all M3TAL stack files, including core and user-defined Docker Compose configurations and Traefik dynamic configuration files. The `/docker` path is a symlink to `/opt/m3tal/stack/`, providing a convenient user-facing alias for all stack management operations.

### Adding a New Stack

To deploy a new Docker Compose stack within the M3TAL ecosystem:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. This can be done using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Execute `m3tal up` from the command line. This command will deploy all Docker Compose stacks present in the `/docker/` directory, including your newly added `my-stack-compose.yml` and any other user-defined compose files.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Use standard `systemctl` commands for its operation:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Dashboard Access

The M3TAL Dashboard provides a web interface for system management. It supports two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### 1. Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting upon a new installation).
*   **Mechanism:** This mode utilizes the `m3tal-compose.local.yml` override file, which adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access:** The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Note:** A new user performing a default M3TAL installation will access the dashboard directly via port `8082`. No Traefik gateway is required for this mode. It is suitable for LAN-only setups or initial local testing.

### 2. Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode uses the `m3tal-compose.traefik.yml` override file, which applies Traefik labels to the dashboard container. Traefik then routes requests for `dash.${DOMAIN}` to the dashboard container on its internal port `8082`.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN`. This requires the Traefik gateway to be running (`m3tal up`).
*   **Note:** This mode is designed for domain-based setups where multiple services are exposed through a reverse proxy.

## Traefik Gateway

The Traefik gateway is deployed as a container via `routing-compose.yml`. It functions as a reverse proxy for all M3TAL services.

*   **Port Binding:** Traefik binds to host port `80` to serve HTTP traffic (and typically `443` for HTTPS, though not explicitly shown in core config).
*   **Service Discovery:** Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions.
*   **Dynamic Configuration:** Traefik utilizes a file provider to load dynamic routing configurations from the `/docker/dynamic/` directory. These configurations support hot-reloading.
*   **API Daemon Routing:** Requests to `api.DOMAIN` are routed to the M3TAL Go API daemon. This is achieved through a dynamic configuration file (e.g., `dynamic/api.yml`) that explicitly routes `api.DOMAIN` to the API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`.
*   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, requests for `dash.DOMAIN` are routed to the `m3tal-dashboard` container based on its Traefik labels.

### Example: Exposing a Custom User Service via Traefik

To expose a custom user-defined service through Traefik, add the appropriate Traefik labels to its service definition in your Docker Compose file (`/docker/my-app-compose.yml`):

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN}`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy

networks:
  proxy:
    external: true
```

After placing this file in `/docker/`, run `m3tal up` to deploy the service and configure Traefik. The service will then be accessible at `http://app.DOMAIN`.

## Quick Demo

To initiate specific M3TAL components or orchestrate the entire system:

*   To start only the M3TAL Dashboard container (and its dependencies) for immediate access, execute:
    ```bash
    m3tal dash up
    ```
    This command will download the necessary dashboard compose files, apply the `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`, and bring up the dashboard container. If using the default `local` mode, the dashboard will be available at `http://HOST_IP:8082`.

*   To orchestrate and deploy all M3TAL core stacks, including the Traefik gateway, API daemon, and any user-defined Docker Compose files placed in the `/docker/` directory, execute:
    ```bash
    m3tal up
    ```
    This command will run `docker compose` across all `*-compose.yml` files in `/docker/`, bringing up the complete M3TAL ecosystem.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access |
|------|---------|--------|
| 80 | Traefik HTTP entry point | Public (when Traefik mode is active) |
| 8080 | M3TAL API daemon (Go) | Host-local |
| 8081 | Traefik dashboard | Host-local only |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) |

**Note:** These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.