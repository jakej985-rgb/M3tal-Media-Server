# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

The M3TAL CLI binary is distributed via an APT repository. Follow these steps to install:

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

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the primary entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary operating as a systemd service, listening on host-local port `8080`. It manages Docker interactions, maintains the M3TAL state database, and exposes the M3TAL API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container. It internally listens on port `8082` and communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy. It exposes services by domain name on host port `80` and utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container providing a Cloudflare tunnel for zero-configuration public internet access to services.

## Filesystem Contract

The following paths represent critical locations within the M3TAL filesystem:

| Path                        | Purpose                                                                             |
| :-------------------------- | :---------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for M3TAL. Managed by `m3tal config wizard`.             |
| `/var/lib/m3tal/state.db`   | SQLite state database for the M3TAL API daemon. Auto-created upon API daemon startup. |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's core Docker Compose stack files and Traefik configuration. |
| `/docker`                   | A symbolic link to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by the `m3tal dashpass` command.                |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all core M3TAL stack files reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, providing a user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`)

This is the default mode for M3TAL installations.
When `DASHBOARD_EXPOSE_MODE=local`, the dashboard container's internal port `8082` is directly bound to the host's `8082` port using an override (`m3tal-compose.local.yml`). This provides direct access without requiring Traefik.

*   **Access Method:** `http://HOST_IP:8082` or `http://localhost:8082`
*   **Requirement:** No Traefik configuration is needed.
*   **Behavior:** A new user performing a default installation will access the dashboard directly via port 8082. This behavior is linked to the default setting `DASHBOARD_EXPOSE_MODE=local`.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

When `DASHBOARD_EXPOSE_MODE=traefik`, specific Traefik labels are added to the dashboard service definition via an override (`m3tal-compose.traefik.yml`). Traefik then routes traffic for `dash.DOMAIN` to the dashboard container's internal port `8082`.

*   **Access Method:** `http://dash.DOMAIN`
*   **Requirement:** Traefik must be running and properly configured to route the specified domain.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration from files within `/docker/dynamic/` (which is symlinked from `/opt/m3tal/stack/dynamic/`). This allows for hot-reloading of routing rules. For example, `dynamic/api.yml` defines a router that routes requests for `api.DOMAIN` to the M3TAL Go API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.DOMAIN` to the dashboard container.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service (e.g., an Nginx instance) via Traefik:

1.  Ensure your service is part of the `proxy` network, which Traefik uses for discovery.
2.  Add the necessary Traefik labels to your service definition within your Docker Compose file (e.g., `my-app-compose.yml`) located in `/docker/`:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    labels:
      - "traefik.enable=true" # Enable Traefik for this service
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)" # Route based on domain
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Target port within the container
      - "traefik.http.routers.myapp.entrypoints=web" # Use the 'web' entrypoint (HTTP port 80)
      - "traefik.docker.network=proxy" # Explicitly tell Traefik which network to use
    networks:
      - proxy # Ensure the service is connected to the 'proxy' network

networks:
  proxy:
    external: true # The 'proxy' network is managed by M3TAL's routing stack
```

After adding or modifying the compose file, run `m3tal up` to apply the changes.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

The M3TAL CLI provides specific commands for managing different parts of the system:

*   **Start the M3TAL Dashboard container:**
    ```bash
    m3tal dash up
    ```
    This command specifically fetches the latest dashboard compose files and starts the dashboard container with the appropriate exposure mode (`local` or `traefik`) based on your `/etc/m3tal/.env` configuration.

*   **Deploy all M3TAL stacks and user-defined compose files:**
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all Docker Compose stacks found within the `/docker/` directory, including the M3TAL control plane components (Traefik, Cloudflared) and any user-defined compose files.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                                                              |
| :--- | :-------------------------- | :------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                                                            |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                                                       |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                                                              |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`.                           |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.