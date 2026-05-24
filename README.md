# M3TAL System Documentation

## Overview
This document provides technical details and operational procedures for the M3TAL system. M3TAL comprises a Go-based CLI binary, an API daemon, a Python/Flask dashboard, and an integrated Docker Compose orchestration layer for managing various services including a Traefik reverse proxy and an optional Cloudflared tunnel.

## Prerequisites
Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## APT Installation
To install the M3TAL CLI binary and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract
The following paths define critical M3TAL system directories and files:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik configuration. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Firewall Considerations
If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Deployment Lifecycle
M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. The `/docker` directory is a symlink to `/opt/m3tal/stack/` and serves as the user-facing alias for all stack operations.

### Adding a New Stack
To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy or update all defined stacks.

## Dashboard Access
The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)
- **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
- **Mechanism:** M3TAL utilizes an override file (`m3tal-compose.local.yml`) that adds a direct Docker port binding (`${DASHBOARD_PORT:-8082}:8082`) to the dashboard container.
- **Access:** `http://HOST_IP:8082` or `http://localhost:8082`.
- **Note:** A new user performing a default installation will access the dashboard directly via port 8082. This mode does not require Traefik to be running. It is suitable for LAN-only setups or initial local testing.

### Traefik Mode
- **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`.
- **Mechanism:** M3TAL utilizes an override file (`m3tal-compose.traefik.yml`) that applies Traefik labels to the dashboard container. These labels instruct Traefik to route requests for `dash.${DOMAIN}` to the dashboard container's internal port `8082`.
- **Access:** `http://dash.DOMAIN` (Traefik must be running and configured, typically via `m3tal up`).
- **Note:** This mode is designed for domain-based access and requires Traefik to be operational.

The M3TAL Dashboard container (`m3tal-dashboard`) communicates with the M3TAL API daemon (Go) at `http://host.docker.internal:8080`.

## Traefik Gateway
Traefik is deployed as a container via `routing-compose.yml` and functions as the M3TAL system's reverse proxy. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds to port 80 on the host machine, serving as the HTTP entry point. It also loads dynamic configuration files from `/docker/dynamic/` (symlinked from `/opt/m3tal/stack/dynamic/`), which enables hot-reloading of routing rules.

### Dynamic Routing Examples:
- **M3TAL API Daemon:** Traefik routes `api.DOMAIN` to the M3TAL API daemon. This is achieved through a dynamic configuration file (e.g., `dynamic/api.yml`) that explicitly routes `api.DOMAIN` to `http://host.docker.internal:8080`, where the Go API daemon listens on the host-local network.
- **M3TAL Dashboard:** When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik labels on the `m3tal-dashboard` service route `dash.DOMAIN` to the dashboard container's internal port `8082`.

### Exposing a Custom User Service via Traefik
To expose a custom Docker Compose service, such as a hypothetical `my-app` service, you must add specific Traefik labels to its service definition:

```yaml
# my-app-compose.yml (located in /docker/)
services:
  my-app:
    image: nginx:alpine
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy # Ensure your service is on the 'proxy' network for Traefik discovery

networks:
  proxy:
    external: true
```
After defining these labels in your Compose file, run `m3tal up` to deploy the service and register it with Traefik.

## Service Management
The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service and listens on host-local port `8080`.

Standard systemctl commands apply for managing the API daemon:
- Check status: `systemctl status m3tal-api`
- Restart service: `systemctl restart m3tal-api`
- View logs: `journalctl -u m3tal-api -f`

## Quick Demo
- To start only the M3TAL Dashboard container, ensuring it's running with the correct exposure mode:
  `m3tal dash up`
- To orchestrate and deploy all other stacks, including any user-defined Docker Compose files placed in `/docker/`, run:
  `m3tal up`

## Port Map
The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.