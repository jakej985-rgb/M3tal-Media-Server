# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system. M3TAL is designed for self-hosted service orchestration, leveraging Docker for containerization and a Go-based API for management.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

M3TAL is distributed via an APT repository. Follow these steps to install the CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Filesystem Contract

The following paths are critical to M3TAL's operation and configuration:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik configuration. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Components

The M3TAL system comprises several interconnected components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on host-local port `8080`. It manages Docker interactions, maintains the SQLite state database, and provides API routes for the dashboard and external tools.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, enabling zero-configuration internet access for exposed services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory is a symlink to `/opt/m3tal/stack/`. This means `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside, and `/docker` is the user-facing symlink alias for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to deploy all defined stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard offers two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   Uses the `m3tal-compose.local.yml` Docker Compose override file.
*   Adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Behavior for new users:** A new user performing a default installation will access the dashboard directly via port `8082`. This direct access via port 8082 is enabled by the default setting `DASHBOARD_EXPOSE_MODE=local`.
*   **Requirements:** No Traefik configuration or external domain is required.
*   **Best for:** LAN-only setups, first-time users, or local testing environments.

### Traefik Mode

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   Uses the `m3tal-compose.traefik.yml` Docker Compose override file.
*   Adds Traefik labels to the dashboard service, enabling Traefik to route `dash.${DOMAIN}` to the dashboard container's internal port `8082`.
*   **Access via:** `http://dash.DOMAIN` (requires Traefik to be running via `m3tal up`).
*   **Requirements:** Traefik must be deployed and configured to handle the specified domain.
*   **Best for:** Domain-based setups and environments where multiple services are managed by a central reverse proxy.

## Traefik Gateway

Traefik acts as the central reverse proxy for M3TAL, deployed as a container via `routing-compose.yml`. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik handles:
*   **Dynamic configuration:** It loads dynamic configuration files (e.g., `dynamic/api.yml`) from `/docker/dynamic/`. These files enable routing requests to services listening on host-local ports, such as routing `api.DOMAIN` to the Go API daemon on host-local port `8080` via `http://host.docker.internal:8080`.
*   **Dashboard routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` routes to the dashboard container on its internal port `8082`.

### Example: Exposing a Custom Service via Traefik

To expose a custom user-defined Docker Compose service through Traefik, add the necessary labels to its service definition:

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
      - "traefik.http.routers.myapp.entrypoints=web"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.docker.network=proxy" # Ensure it's on the same network as Traefik
    networks:
      - proxy

networks:
  proxy:
    external: true # Assumes 'proxy' network is defined as external for Traefik
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically configure routing for `app.DOMAIN` to your `my-app` service.

## Quick Demo

To quickly start the M3TAL Dashboard:

1.  Ensure Docker and Docker Compose V2 are installed.
2.  Install M3TAL using the APT instructions above.
3.  Run `m3tal dash up`. This command specifically downloads the latest dashboard compose files and starts the dashboard container with the appropriate configuration based on `DASHBOARD_EXPOSE_MODE`.

To deploy all other configured stacks, including any user-defined compose files in `/docker/`:

*   Run `m3tal up`. This command orchestrates and deploys all `*-compose.yml` files found in the `/docker/` directory.

## Service Management

The M3TAL API daemon runs as a systemd service. Use standard `systemctl` commands for management:

*   Check service status: `systemctl status m3tal-api.service`
*   Restart the API daemon: `systemctl restart m3tal-api.service`
*   View logs: `journalctl -u m3tal-api.service -f`

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|-----------------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Environment Variables

M3TAL uses environment variables for system configuration, primarily managed through `/etc/m3tal/.env`. The following are key variables and their default values:

| Key | Default Value | Description |
|---|---|---|
| `DASHBOARD_PORT` | `8082` | Internal port for the M3TAL Dashboard. |
| `DASHBOARD_EXPOSE_MODE` | `local` | `local` for direct port access, `traefik` for Traefik routing. |
| `HTTP_PORT` | `8080` | Port for the M3TAL API daemon. |
| `STATE_DIR` | `./state` | Directory for the SQLite state database. |
| `LOG_LEVEL` | `info` | Logging level for M3TAL components. |
| `DASHBOARD_SECRET` | `change_me_immediately` | Secret key for dashboard session management. **Change this immediately.** |
| `API_TOKEN` | `change_me_api_token` | API token for authentication with the M3TAL API. **Change this immediately.** |
| `ADMIN_PASSWORD` | `admin_pass` | Default admin password for the M3TAL Dashboard. **Change this immediately.** |
| `NETWORK_NAME` | `m3tal` | Default Docker network name used by M3TAL. |
| `LOCAL_IP` | `127.0.0.1` | Local IP address for internal routing. |
| `DOMAIN` | `localhost` | Base domain for Traefik routing if not overridden. |
| `VPN_USER` | `user` | Placeholder for VPN username. |
| `VPN_PASSWORD` | `password` | Placeholder for VPN password. |
| `BASE_STORAGE_PATH` | `./data` | Base directory for all M3TAL data storage. |
| `MEDIA_PATH` | `./data/media` | Path for media storage. |
| `CONFIG_PATH` | `./data/config` | Path for configuration files. |
| `DOWNLOADS_PATH` | `./data/downloads` | Path for download storage. |
| `PUID` | `1000` | User ID for container processes. |
| `PGID` | `1000` | Group ID for container processes. |
| `TZ` | `America/Denver` | Timezone setting for containers. |
| `TRAEFIK_WEB_PORT` | `80` | Public HTTP port for Traefik. |
| `TRAEFIK_WEBHTTPS_PORT` | `443` | Public HTTPS port for Traefik. |
| `TRAEFIK_DASHBOARD_PORT` | `8080` | Internal port for the Traefik dashboard. |
| `DEBUG_MODE` | `false` | Enables debug logging and features. |
| `METRICS_ENABLED` | `true` | Enables metrics collection. |