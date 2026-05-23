# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

**Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.**

## APT Installation

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

The following table details the critical directories and files within the M3TAL system:

| Path                        | Purpose                                                              |
|-----------------------------|----------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.               |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config.|
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. User-facing path for stack operations.|
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |

## Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container, internally listening on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container that provides a Cloudflare tunnel for zero-configuration internet access to services.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|-----------------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory is a user-facing symlink alias for stack operations. The canonical source of truth for all stack files, including `m3tal-compose.yml` and routing configurations, is `/opt/m3tal/stack/`. All user-defined Docker Compose files should be placed in the `/docker/` directory.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. These can be managed using `m3tal config wizard` or `m3tal config set KEY value`.
3. Run `m3tal up` to deploy your new stack along with all other existing M3TAL-managed stacks.

## Dashboard Access

The M3TAL Dashboard offers two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`).
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirement**: No Traefik gateway is required for this mode. A new user performing a default M3TAL installation will access the dashboard directly via port 8082, due to the default `DASHBOARD_EXPOSE_MODE=local` setting.
*   **Use Case**: Ideal for LAN-only setups, initial configurations, and local testing environments.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: This mode utilizes the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard service. Traefik then routes traffic based on these labels.
*   **Access**: The dashboard is accessible via a domain, typically `http://dash.DOMAIN` (where `DOMAIN` is configured in `/etc/m3tal/.env`).
*   **Requirement**: The Traefik gateway must be running (`m3tal up` will deploy it).
*   **Use Case**: Suited for environments requiring domain-based access and integration with a reverse proxy for multiple services.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml` and binds host port `80` as its HTTP entry point. It also loads dynamic configuration files from `/docker/dynamic/` (which is `/opt/m3tal/stack/dynamic/`), allowing for hot-reloading of routing rules.

An example of dynamic configuration is the routing of `api.DOMAIN`. The file `/docker/dynamic/api.yml` defines a router that directs requests for `api.DOMAIN` to the M3TAL API daemon, which listens on the host-local port `8080`. This is achieved by routing to `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` routes to the dashboard container via Traefik labels defined in its compose configuration.

To expose a custom user service via Traefik, add appropriate labels to its Docker Compose service definition. For instance, in a `my-app-compose.yml`:

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
      - proxy # Ensure your service is on the 'proxy' network for Traefik discovery
```

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. Standard `systemctl` commands are used for its management:

*   **Check Status**: `systemctl status m3tal-api`
*   **Restart Service**: `systemctl restart m3tal-api`
*   **View Logs**: `journalctl -u m3tal-api -f`

## Quick Demo

*   To start the M3TAL Dashboard container specifically, including fetching the latest compose files and applying the correct override based on your `DASHBOARD_EXPOSE_MODE` setting, execute:
    ```bash
    m3tal dash up
    ```
*   To orchestrate and deploy all M3TAL-managed Docker Compose stacks, including the Traefik gateway, Cloudflared, and any user-defined compose files located in `/docker/`, run:
    ```bash
    m3tal up
    ```