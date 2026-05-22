# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Overview

M3TAL is a system designed to manage and orchestrate services using Docker containers and systemd. It comprises a Go CLI, a Go API daemon, a Python/Flask dashboard, and a Traefik reverse proxy for service routing.

## Installation

To install the M3TAL CLI and API daemon via APT:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The following paths are critical to M3TAL's operation:

| Path                        | Purpose                                                                                                                                                              |
|-----------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`.                                                                              |
| `/var/lib/m3tal/state.db`   | SQLite database storing M3TAL's operational state. Auto-created and managed by the API daemon.                                                                       |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose files and Traefik dynamic configuration. This is the source of truth for stack definitions.                               |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations and where user-defined compose files should be placed.          |
| `/docker/users.json`        | Dashboard credential store for user authentication. Managed by `m3tal dashpass`.                                                                                     |

## Components

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the primary command-line interface for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container, internally listening on port `8082`. It communicates with the API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services on host port `80` by domain name. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container providing Cloudflare tunnel capabilities for secure, zero-config internet access to services.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                                                             |
|------|-----------------------------|---------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                                                           |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                                                      |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                                                             |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`.                        |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all M3TAL-managed stack files reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, intended as the user-facing path for all stack operations and where users should place their custom Docker Compose files.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy your new stack along with all other managed stacks.

## Dashboard Access

The M3TAL Dashboard offers two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`)

This is the default mode.
*   The dashboard container's internal port `8082` is directly bound to the host's port `8082`.
*   Access is via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **A new user performing a default installation will access the dashboard directly via port `8082`**. This mode does not require Traefik to be running.
*   This mode uses the `m3tal-compose.local.yml` override file.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

*   The dashboard is exposed via Traefik labels, routing `dash.DOMAIN` to the dashboard container on its internal port `8082`.
*   Access is via `http://dash.DOMAIN`. This requires Traefik to be running via `m3tal up` and `DOMAIN` to be configured in `/etc/m3tal/.env`.
*   This mode uses the `m3tal-compose.traefik.yml` override file.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml`. It binds host port `80` as its primary HTTP entry point and is configured to load dynamic routing configurations from files placed in `/docker/dynamic/` (which is symlinked from `/opt/m3tal/stack/dynamic/`).

*   **Routing to Containerized Services**:
    For services running within Docker containers, Traefik routes requests based on labels attached to the service in its Docker Compose definition. For example, when `DASHBOARD_EXPOSE_MODE=traefik`, the `dash.DOMAIN` route is handled this way, directing traffic to the `m3tal-dashboard` container.

*   **Routing to Host-Local Services**:
    For services that run directly on the host machine (e.g., the M3TAL API daemon listening on `http://localhost:8080`), Traefik uses dynamic configuration files (such as `/docker/dynamic/api.yml`). This file explicitly instructs Traefik to route `api.DOMAIN` to the Go API daemon by targeting `http://host.docker.internal:8080`. `host.docker.internal` is a special DNS name within Docker networks that resolves to the host machine's IP address.

### Example: Exposing a Custom Service via Traefik

To expose a hypothetical `my-app` service via Traefik, you would add labels to its Docker Compose service definition:

```yaml
# /docker/my-app-compose.yml
version: '3.8'
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
    external: true # Assumes the 'proxy' network is created externally by Traefik.
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically detect `my-app` and route requests for `app.YOUR_DOMAIN` to it.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart the API daemon**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

*   To specifically start only the M3TAL Dashboard container, use the command:
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard Docker Compose files and starts the dashboard container with the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.

*   To orchestrate and deploy all other stacks, including any user-defined Docker Compose files placed in the `/docker/` directory, use:
    ```bash
    m3tal up
    ```
    This command will bring up Traefik, Cloudflared (if configured), and any other `*-compose.yml` files present in `/docker/`.