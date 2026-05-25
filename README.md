# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

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

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Filesystem Contract

The following table outlines the critical filesystem paths used by M3TAL:

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.             |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.           |

## Core Components

The M3TAL system consists of the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A Go binary providing a unified command-line interface for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container. It operates internally on port `8082` and communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, enabling zero-configuration internet access for services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

It is crucial to note that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside. The `/docker` directory is a user-facing symlink alias, simplifying all stack operations for the user.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Execute `m3tal up` to deploy all defined Docker Compose stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** This mode utilizes the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`). No Traefik configuration is involved.
*   **Access Method:** `http://HOST_IP:8082` or `http://localhost:8082`.
*   **New User Experience:** A new user performing a default M3TAL installation will access the dashboard directly via port `8082` due to `DASHBOARD_EXPOSE_MODE=local` being the default setting.
*   **Use Cases:** Ideal for LAN-only setups, first-time users, and local testing environments.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode utilizes the `m3tal-compose.traefik.yml` override file. This file adds specific Traefik labels to the dashboard service, allowing Traefik to route `dash.${DOMAIN}` to the dashboard container on its internal port `8082`. Requires the Traefik gateway to be running via `m3tal up`.
*   **Access Method:** `http://dash.DOMAIN`
*   **Use Cases:** Suited for domain-based deployments and environments where multiple services are managed behind a single reverse proxy.

## Traefik Gateway

The Traefik gateway is deployed as a Docker container via `routing-compose.yml`. It functions as a reverse proxy, binding to host port `80` (and `443` if configured for HTTPS) as its primary HTTP entry point.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration from the `/docker/dynamic/` directory (a file provider with hot-reload capabilities). This is used for routing requests to services not defined directly as Docker containers, such as the M3TAL API daemon:

*   **API Daemon Routing:** The `dynamic/api.yml` configuration routes `api.DOMAIN` to the Go API daemon, which listens on the host-local port `8080`, via `http://host.docker.internal:8080`.
*   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` is routed to the M3TAL Dashboard container (internal port `8082`) through labels defined in `m3tal-compose.traefik.yml`.

### Exposing a Custom User Service via Traefik

To expose a custom service (e.g., `my-app`) through Traefik, add the necessary Traefik labels to its service definition within your Docker Compose file (e.g., `my-app-compose.yml`):

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
    # ... other service configurations

networks:
  proxy:
    external: true # Ensure your service joins the shared 'proxy' network
```

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Standard systemctl commands apply:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

*   To start only the M3TAL Dashboard container:
    ```bash
    m3tal dash up
    ```
    This command specifically manages the dashboard container, downloading the necessary compose files and applying the correct override based on your `DASHBOARD_EXPOSE_MODE` setting.

*   To orchestrate and deploy all other Docker Compose stacks, including any user-defined `*-compose.yml` files located in the `/docker/` directory:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.