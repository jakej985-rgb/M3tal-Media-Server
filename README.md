# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Filesystem Contract

The following table outlines the key file paths and their purposes within the M3TAL system:

| Path                        | Purpose                                                                 |
|-----------------------------|-------------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.           |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.                  |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config.   |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                |

## Installation

To install the M3TAL system, follow these steps:

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

The directory `/opt/m3tal/stack/` serves as the canonical source of truth where all Docker Compose stack files reside. The `/docker` directory is a symlink to `/opt/m3tal/stack/`, providing a convenient, user-facing alias for all stack-related operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to start all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

-   **Setting:** `DASHBOARD_EXPOSE_MODE=local`
-   **Access:** When `DASHBOARD_EXPOSE_MODE` is set to `local` (the default for a new installation), the M3TAL Dashboard container exposes port `8082` directly on the host machine. Access is achieved via `http://HOST_IP:8082` or `http://localhost:8082`. This mode bypasses Traefik and is suitable for local network access or initial setup.
-   A new user performing a default M3TAL installation will access the dashboard directly via port `8082`, not a domain, due to the default `DASHBOARD_EXPOSE_MODE=local` setting.

### Traefik Mode

-   **Setting:** `DASHBOARD_EXPOSE_MODE=traefik`
-   **Access:** When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, the dashboard is exposed via the Traefik Gateway using domain routing. Access is achieved via `http://dash.DOMAIN`, where `DOMAIN` is configured in `/etc/m3tal/.env`. This mode requires Traefik to be running and correctly configured to route traffic based on the Traefik labels defined for the dashboard service.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik utilizes dynamic configuration files, such as `/docker/dynamic/api.yml`, to define routing rules for services not managed directly as Docker containers or for specific host-local services. For example, `api.DOMAIN` is routed to the M3TAL Go API daemon listening on host-local port `8080` via `http://host.docker.internal:8080` through this mechanism. Similarly, `dash.DOMAIN` is routed to the M3TAL Dashboard container when `DASHBOARD_EXPOSE_MODE=traefik` by interpreting specific Traefik labels on the dashboard service definition.

To expose a custom user service via Traefik, add appropriate labels to its Docker Compose service definition:

```yaml
# Example: my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy # Ensure your service is on the 'proxy' network
```

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`, which can be managed using standard systemctl commands:

-   Check service status: `systemctl status m3tal-api`
-   Restart the API daemon: `systemctl restart m3tal-api`
-   View logs in real-time: `journalctl -u m3tal-api -f`

## Quick Demo

-   To specifically start and manage the M3TAL Dashboard container and its dependencies (like the `proxy` network), use `m3tal dash up`. This command ensures the dashboard's `m3tal-compose.yml` and relevant overrides are up-to-date and deployed.
-   The `m3tal up` command orchestrates and deploys *all* Docker Compose stacks found in the `/docker/` directory, including the dashboard (if not already managed by `m3tal dash up` or if explicitly configured) and any user-defined compose files.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.