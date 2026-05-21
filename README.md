# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Filesystem Contract

The following paths are critical for M3TAL operation and configuration:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## APT Installation

To install M3TAL via APT:

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

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files (e.g., `m3tal-compose.yml`, `routing-compose.yml`, and user-defined compose files) reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/` for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`:

### Local Mode (Default)

-   **Setting**: `DASHBOARD_EXPOSE_MODE=local`
-   **Mechanism**: This mode uses an override (`m3tal-compose.local.yml`) to add a direct port binding, making the dashboard accessible directly on the host's IP address.
-   **Access**: `http://HOST_IP:8082` or `http://localhost:8082`
-   **Requirements**: No Traefik configuration or domain is required. A new user performing a default installation will access the dashboard directly via port `8082`, as this is the behavior linked to the default setting `DASHBOARD_EXPOSE_MODE=local`. This is ideal for LAN-only setups or initial testing.

### Traefik Mode

-   **Setting**: `DASHBOARD_EXPOSE_MODE=traefik`
-   **Mechanism**: This mode uses an override (`m3tal-compose.traefik.yml`) that adds Traefik labels to the dashboard container. Traefik (if running) then routes requests for `dash.DOMAIN` to the dashboard container's internal port 8082.
-   **Access**: `http://dash.DOMAIN` (where `DOMAIN` is configured in `/etc/m3tal/.env`).
-   **Requirements**: Traefik must be running via `m3tal up` for this mode to function. This is suitable for domain-based setups with multiple services behind a reverse proxy.

## Traefik Gateway

Traefik acts as the primary reverse proxy for M3TAL, deployed as a container via `routing-compose.yml`. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files, such as `/docker/dynamic/api.yml`, to route requests. For instance, `api.DOMAIN` is routed to the M3TAL Go API daemon, which listens on the host-local port `8080`. This is achieved by defining a service in the dynamic configuration that points to `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.DOMAIN` to the dashboard container.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service (e.g., defined in `my-app-compose.yml`) via Traefik:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx
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
After placing this file in `/docker/` and running `m3tal up`, your `my-app` service will be accessible via `http://app.YOUR_DOMAIN` (assuming Traefik is running and `YOUR_DOMAIN` is set in `/etc/m3tal/.env`).

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle and inspect its status using standard `systemctl` commands:

-   **Check Status**: `systemctl status m3tal-api`
-   **Restart Service**: `systemctl restart m3tal-api`
-   **View Logs**: `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get started with M3TAL:

-   To start the M3TAL Dashboard container specifically (useful for initial setup or troubleshooting):
    ```bash
    m3tal dash up
    ```
    This command will download the necessary dashboard compose files, read your `DASHBOARD_EXPOSE_MODE` setting, and start the dashboard with the appropriate configuration.
-   To orchestrate and deploy all M3TAL-managed Docker Compose stacks (including the Traefik gateway, Cloudflared, and any user-defined compose files in `/docker/`):
    ```bash
    m3tal up
    ```
    This command ensures all services defined in `*-compose.yml` files within `/docker/` are running, including any custom stacks you have added.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.