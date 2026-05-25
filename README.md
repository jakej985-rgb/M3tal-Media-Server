# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system. M3TAL components include a Go CLI binary, a Go API daemon, a Python/Flask dashboard container, and a Traefik-based routing gateway, all orchestrated via Docker Engine and Docker Compose V2.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

M3TAL is distributed via an APT repository. Follow these steps to install the core components:

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

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.             |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.           |

## Core Components

The M3TAL system consists of the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing a single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, enabling zero-configuration internet access for exposed services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory is a symlink to `/opt/m3tal/stack/`. This means `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside, and `/docker` serves as the user-facing alias for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: This mode uses an override (`m3tal-compose.local.yml`) to add a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access**: You can access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirement**: No Traefik is required for this mode.
*   **Note**: A new user performing a default installation will access the dashboard directly via port `8082`, as this is the behavior linked to the default setting `DASHBOARD_EXPOSE_MODE=local`.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: This mode uses an override (`m3tal-compose.traefik.yml`) to add specific Traefik labels to the dashboard container, allowing Traefik to route requests for `dash.${DOMAIN}` to the dashboard on its internal port `8082`.
*   **Access**: You can access the dashboard via `http://dash.DOMAIN` (Traefik must be running via `m3tal up` for this to work).
*   **Requirement**: Traefik must be deployed and configured to handle routing.

## Traefik Gateway

Traefik acts as the central reverse proxy for M3TAL, deployed as a container via `routing-compose.yml`. It binds host port `80` (and `443` if HTTPS is enabled) as the HTTP entry point. Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also loads dynamic configuration files from `/docker/dynamic/` (which is `/opt/m3tal/stack/dynamic/`). These dynamic configurations enable routing for services that are not inherently Docker containers, such as the M3TAL API daemon. For example, `dynamic/api.yml` routes `api.DOMAIN` to the Go API daemon listening on host-local port `8080` by directing requests to `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE=traefik`, utilizing its internal port `8082`.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service (e.g., defined in `my-app-compose.yml`) via Traefik, add the necessary Traefik labels to its service definition:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
      - "traefik.docker.network=proxy" # Ensure this service is on the 'proxy' network
    networks:
      - proxy

networks:
  proxy:
    external: true # Assumes the 'proxy' network is external, created by M3TAL
```
After placing this file in `/docker/`, run `m3tal up` to deploy the service and configure Traefik.

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api.service`
*   **Restart service**: `systemctl restart m3tal-api.service`
*   **View logs**: `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly get started and interact with the M3TAL system:

*   To start only the M3TAL dashboard container, run:
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard compose files and starts the dashboard with the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.
*   To orchestrate and deploy all other stacks, including the Traefik gateway, Cloudflared, and any user-defined compose files located in `/docker/`, run:
    ```bash
    m3tal up
    ```
    This command will bring up all Docker Compose stacks found in the `/docker/` directory.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.