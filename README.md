# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

M3TAL is distributed as an APT package. Follow these steps to install:

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

The following paths define the M3TAL system's filesystem structure:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Components

The M3TAL system comprises the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on host port `8080`. It manages Docker interactions, the SQLite state database, and exposes REST API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare Tunnel Docker container for secure, zero-configuration internet access to services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. The `/docker` path is a symlink alias to `/opt/m3tal/stack/`, serving as the user-facing directory for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Execute `m3tal up` to deploy all stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard provides two distinct access modes, configured via the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** Uses the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (e.g., `${DASHBOARD_PORT:-8082}:8082`).
*   **Access:** Direct access via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik instance is required for this mode. A new user performing a default installation will access the dashboard directly via port 8082. This direct access via port 8082 is due to the default setting `DASHBOARD_EXPOSE_MODE=local`.
*   **Use Case:** Ideal for LAN-only deployments, initial setups, and local testing environments.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override file, which applies Traefik labels to the dashboard service. Traefik then routes `dash.${DOMAIN}` to the dashboard container on its internal port `8082`.
*   **Access:** Via a domain at `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Requirements:** Traefik must be running (`m3tal up` will start it if configured).
*   **Use Case:** Suited for domain-based setups and environments where multiple services are managed behind Traefik as a reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration files from `/docker/dynamic/` (which is a symlink to `/opt/m3tal/stack/dynamic/`). This allows for hot-reloading of routing rules. For instance, `dynamic/api.yml` is used to route requests for `api.DOMAIN` to the M3TAL API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container when `DASHBOARD_EXPOSE_MODE=traefik` is set.

### Exposing a Custom User Service via Traefik

To expose your own Docker Compose service through Traefik, add appropriate labels to its service definition:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-web-app
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy # Essential for Traefik to discover the service

networks:
  proxy:
    external: true # Assumes 'proxy' network is created by M3TAL's routing stack
```

After placing this file in `/docker/` and running `m3tal up`, your service will be accessible via `http://app.DOMAIN`.

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard systemctl commands:

*   **Check Status:** `systemctl status m3tal-api.service`
*   **Restart Service:** `systemctl restart m3tal-api.service`
*   **View Logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly get started and see the M3TAL Dashboard:

1.  Ensure Docker Engine and Docker Compose V2 are installed.
2.  Install M3TAL using the APT instructions above.
3.  Start the M3TAL Dashboard container specifically:
    ```bash
    m3tal dash up
    ```
    This command will download the necessary compose files and start the dashboard container using the default `local` expose mode. You can then access it at `http://HOST_IP:8082`.
4.  To deploy the full M3TAL system, including Traefik and any user-defined Docker Compose stacks placed in `/docker/`:
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all `*-compose.yml` files found in the `/docker/` directory.

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

The M3TAL system relies on various environment variables, typically managed within `/etc/m3tal/.env`. Below are the primary variables and their default values:

| Key | Default Value | Description |
|---|---|---|
| `DASHBOARD_PORT` | `8082` | The internal port for the M3TAL Dashboard container. |
| `DASHBOARD_EXPOSE_MODE` | `local` | Controls how the dashboard is exposed: `local` (direct port) or `traefik` (via domain). |
| `HTTP_PORT` | `8080` | The host-local port on which the M3TAL API daemon listens. |
| `STATE_DIR` | `./state` | Directory for the SQLite state database and other runtime data. |
| `LOG_LEVEL` | `info` | Logging verbosity for the API daemon (`debug`, `info`, `warn`, `error`). |
| `DASHBOARD_SECRET` | `change_me_immediately` | Secret key for dashboard session management. **Change this immediately post-install.** |
| `API_TOKEN` | `change_me_api_token` | API token for authentication with the M3TAL API. **Change this immediately post-install.** |
| `ADMIN_PASSWORD` | `admin_pass` | Default administrator password for the M3TAL Dashboard. **Change this immediately post-install.** |
| `NETWORK_NAME` | `m3tal` | Default Docker network name for M3TAL services. |
| `LOCAL_IP` | `127.0.0.1` | Local IP address for specific configurations. |
| `DOMAIN` | `localhost` | Base domain for Traefik routing when `DASHBOARD_EXPOSE_MODE=traefik`. |
| `VPN_USER` | `user` | Username for VPN services, if integrated. |
| `VPN_PASSWORD` | `password` | Password for VPN services, if integrated. |
| `BASE_STORAGE_PATH` | `./data` | Base directory for all persistent storage volumes. |
| `MEDIA_PATH` | `./data/media` | Subdirectory for media storage. |
| `CONFIG_PATH` | `./data/config` | Subdirectory for configuration files. |
| `DOWNLOADS_PATH` | `./data/downloads` | Subdirectory for downloads. |
| `PUID` | `1000` | User ID for container processes to ensure correct permissions. |
| `PGID` | `1000` | Group ID for container processes to ensure correct permissions. |
| `TZ` | `America/Denver` | Timezone setting for containers. |
| `TRAEFIK_WEB_PORT` | `80` | Host port for Traefik's HTTP entry point. |
| `TRAEFIK_WEBHTTPS_PORT` | `443` | Host port for Traefik's HTTPS entry point. |
| `TRAEFIK_DASHBOARD_PORT` | `8080` | Internal port for the Traefik dashboard. (Note: exposed as `127.0.0.1:8081` on host). |
| `DEBUG_MODE` | `false` | Enables debug features or verbose logging. |
| `METRICS_ENABLED` | `true` | Enables or disables metrics collection. |