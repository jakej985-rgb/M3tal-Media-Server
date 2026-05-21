# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## APT Installation

To install the M3TAL CLI binary and systemd service:

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

## Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary operating as a systemd service, listening on host-local port `8080`. This daemon manages Docker interactions, maintains the SQLite state database, and exposes M3TAL's API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application packaged as a Docker container, running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container for establishing Cloudflare tunnels, enabling secure, zero-configuration internet access to services.

## Filesystem Contract

The following table outlines the critical file system paths and their purposes within the M3TAL system:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for M3TAL. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by the `m3tal dashpass` command.     |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all Docker Compose stack files (e.g., `m3tal-compose.yml`, `routing-compose.yml`) and Traefik dynamic configuration files reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/`, intended for all stack-related operations by the user.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy all currently defined stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`) - Default

*   This is the default mode for a new installation.
*   The `m3tal-dashboard` container is configured with a direct port binding, typically `${DASHBOARD_PORT:-8082}:8082`. This is achieved via the `m3tal-compose.local.yml` override file.
*   Access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   This mode does not require Traefik to be running for dashboard access.
*   A new user performing a default installation will access the dashboard directly via port 8082.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

*   When this mode is enabled, Traefik labels are applied to the `m3tal-dashboard` service via the `m3tal-compose.traefik.yml` override file.
*   These labels instruct Traefik to route traffic for `dash.DOMAIN` to the dashboard container, which internally listens on port 8082.
*   Access the dashboard via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   This mode requires Traefik to be running and properly configured for domain routing.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration from files located in `/docker/dynamic/` (which is symlinked from `/opt/m3tal/stack/dynamic/`). This file provider allows for hot-reloading of routing rules. For instance, the `dynamic/api.yml` file defines a router for `api.DOMAIN` that routes requests to the M3TAL Go API daemon, which listens on host-local port `8080`. This routing is achieved by targeting `http://host.docker.internal:8080` from within the Traefik container. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, the `dash.DOMAIN` route is handled by Traefik labels on the dashboard container, directing traffic to its internal port `8082`.

### Example: Exposing a Custom Service via Traefik

To expose a hypothetical `my-app` service through Traefik, you would add labels to its service definition in your Docker Compose file (e.g., `my-app-compose.yml`):

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

The M3TAL API daemon is managed by systemd as the `m3tal-api.service`. Standard `systemctl` commands can be used for its management:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

*   To specifically start only the M3TAL Dashboard container (useful for initial setup or troubleshooting):
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard compose files, reads your `DASHBOARD_EXPOSE_MODE` setting, and starts the dashboard with the appropriate overrides.

*   To orchestrate and deploy all other Docker Compose stacks, including any user-defined compose files you have placed in the `/docker/` directory:
    ```bash
    m3tal up
    ```
    This command will bring up Traefik, Cloudflared, and any other services defined by `*-compose.yml` files in the `/docker/` directory.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.