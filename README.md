# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system. M3TAL is designed to orchestrate and manage Docker containers, providing a unified CLI, an API daemon, and a web dashboard for system administration.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

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

## Filesystem Contract

The following paths define the M3TAL filesystem contract:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory is a user-facing symlink alias for `/opt/m3tal/stack/`, which is the canonical source of truth directory where all Docker Compose stack files reside.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy your new stack along with all other M3TAL-managed services.

## Dashboard Access

The M3TAL Dashboard provides a web-based interface for system management. Its access method is controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation.)
*   **Mechanism**: A direct port binding is added to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`), exposing it directly on the host network. Traefik is not involved in routing the dashboard in this mode.
*   **Access**: `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Note**: A new user performing a default installation will access the dashboard directly via port 8082. This behavior is linked to the default `DASHBOARD_EXPOSE_MODE=local` setting.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: Traefik labels are applied to the dashboard container, allowing Traefik to route incoming requests for `dash.DOMAIN` to the dashboard's internal port (8082). This mode requires Traefik to be running via `m3tal up`.
*   **Access**: `http://dash.DOMAIN`

## Traefik Gateway

Traefik acts as the reverse proxy gateway for the M3TAL system. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik uses a file provider to load dynamic configuration from `/docker/dynamic/` (specifically `/etc/traefik/dynamic` within the container). This allows for routing requests to services that are not necessarily Docker containers managed by labels. For instance, `api.DOMAIN` is routed to the M3TAL Go API daemon, which listens on the host-local port `8080`, by configuring a service in `dynamic/api.yml` to target `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the M3TAL dashboard container when `DASHBOARD_EXPOSE_MODE=traefik`.

### Example: Exposing a Custom Service via Traefik

To expose a hypothetical `my-app` service via Traefik, you would add specific labels to its Docker Compose service definition in `my-app-compose.yml` (located in `/docker/`):

```yaml
services:
  my-app:
    image: nginx:alpine
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy

networks:
  proxy:
    external: true
```

In this example:
*   `traefik.enable=true` makes the service discoverable by Traefik.
*   `traefik.http.routers.myapp.rule=Host(\`app.DOMAIN\`)` defines a router that matches requests for `app.DOMAIN`.
*   `traefik.http.services.myapp.loadbalancer.server.port=80` tells Traefik the internal port of the `my-app` container to which traffic should be directed.
*   `traefik.http.routers.myapp.entrypoints=web` specifies that this router listens on the 'web' entry point (which typically maps to host port 80).
*   The `proxy` network is required for Traefik to communicate with the service.

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage its lifecycle using standard systemctl commands:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

### Starting the Dashboard

To specifically start or restart the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command ensures the dashboard is running with the correct configuration based on your `DASHBOARD_EXPOSE_MODE` setting.

### Deploying All Stacks

To orchestrate and deploy all M3TAL-managed Docker Compose stacks, including any user-defined compose files placed in the `/docker/` directory:

```bash
m3tal up
```

This command processes all `*-compose.yml` files in `/docker/` and brings up their respective services.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.