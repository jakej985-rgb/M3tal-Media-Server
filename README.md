# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

---

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Overview

M3TAL is a system comprising a Go-based CLI, an API daemon, a Python/Flask dashboard, and a Traefik-powered routing gateway. It leverages Docker Engine and Docker Compose V2 for container orchestration, providing a structured environment for deploying and managing services. All configuration is managed via `.env` files and a SQLite state database.

## Filesystem Contract

The following table details the canonical paths and their purposes within the M3TAL filesystem:

| Path                       | Purpose                                                                                                 |
|----------------------------|---------------------------------------------------------------------------------------------------------|
| `/etc/m3tal/.env`          | Primary configuration file. Managed by `m3tal config wizard` or `m3tal config set`.                     |
| `/var/lib/m3tal/state.db`  | SQLite state database. Auto-created and managed by the M3TAL API daemon.                                |
| `/opt/m3tal/stack/`        | Canonical directory for M3TAL-managed Docker Compose files and Traefik dynamic configuration.           |
| `/docker`                  | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations.   |
| `/docker/users.json`       | Dashboard credential store. Managed by the `m3tal dashpass` command.                                    |

## Installation

To install the M3TAL CLI binary and API daemon, follow these steps:

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

The M3TAL system consists of the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing a single entrypoint for all M3TAL operations, including configuration management, service control, and stack deployment.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal SQLite state database, and exposes the M3TAL API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container running internally on port `8082`. It communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy. It exposes services via domain names on host port `80` and uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container providing Cloudflare tunnel capabilities for secure, zero-config internet access to services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The canonical source of truth for all stack files is `/opt/m3tal/stack/`. The `/docker` directory is a symlink to `/opt/m3tal/stack/`, serving as the user-facing alias for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these variables.
3.  Execute `m3tal up` to deploy and start all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the dashboard container.
*   **Access:** The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Prerequisites:** No Traefik configuration or domain setup is required.
*   **Behavior:** A new user performing a default installation will access the dashboard directly via port 8082, as this is the behavior linked to the default setting `DASHBOARD_EXPOSE_MODE=local`.
*   **Use Case:** Ideal for LAN-only setups, first-time users, and local development or testing environments.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode utilizes the `m3tal-compose.traefik.yml` override file. This file applies Traefik labels to the dashboard service, enabling Traefik to route traffic from a specified domain to the dashboard container's internal port 8082.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Prerequisites:** The Traefik gateway must be running (typically started via `m3tal up`) and configured to handle the `dash.DOMAIN` route.
*   **Use Case:** Suited for domain-based setups where multiple services are exposed behind a single reverse proxy.

## Traefik Gateway

The Traefik gateway is deployed as a Docker container, typically via `routing-compose.yml`. It acts as a reverse proxy, handling incoming HTTP requests and routing them to the appropriate backend services.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

### Dynamic Configuration

Traefik loads dynamic configuration from files located in `/docker/dynamic/` (which is symlinked from `/opt/m3tal/stack/dynamic/`). This allows for hot-reloading of routing rules without restarting Traefik.

A key example is routing requests for `api.DOMAIN` to the M3TAL Go API daemon, which listens on the host-local port `8080`. This is achieved through a dynamic configuration file like `dynamic/api.yml`:

```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)"
      service: api
      entryPoints:
        - web

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080"
```
This configuration routes `api.DOMAIN` to the host-local port `8080`, specifically targeting the M3TAL API daemon. Similarly, `dash.DOMAIN` routes to the dashboard container when `DASHBOARD_EXPOSE_MODE=traefik`, leveraging labels within the dashboard's Docker Compose definition.

### Exposing Custom Services

To expose a custom user service via Traefik, add the necessary Traefik labels to its service definition in your Docker Compose file (e.g., `my-app-compose.yml`):

```yaml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-app
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Internal port of the service
      - "traefik.http.routers.myapp.entrypoints=web" # Assuming 'web' is your HTTP entrypoint
    networks:
      - proxy # Ensure the service is on the same network as Traefik
    restart: unless-stopped

networks:
  proxy:
    external: true # Traefik typically uses an external 'proxy' network
```

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get started with M3TAL:

*   To specifically start only the M3TAL dashboard container (along with any necessary base configuration), use the command:
    `m3tal dash up`
*   To orchestrate and deploy all other stacks defined by `*-compose.yml` files in the `/docker/` directory, including any user-defined compose files and the Traefik gateway:
    `m3tal up`

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                             |
|------|-----------------------------|---------------------------------------------|---------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Configuration (Environment Variables)

M3TAL uses environment variables for system configuration, primarily defined in `/etc/m3tal/.env`. The following are key variables and their default values:

| Key                       | Default Value         | Description                                                          |
|---------------------------|-----------------------|----------------------------------------------------------------------|
| `DASHBOARD_PORT`          | `8082`                | Internal port for the M3TAL Dashboard container.                     |
| `DASHBOARD_EXPOSE_MODE`   | `local`               | `local` for direct port access, `traefik` for domain-based access. |
| `HTTP_PORT`               | `8080`                | Port for the M3TAL API daemon.                                       |
| `STATE_DIR`               | `./state`             | Directory for the SQLite state database.                             |
| `LOG_LEVEL`               | `info`                | Logging verbosity for the API daemon.                                |
| `DASHBOARD_SECRET`        | `change_me_immediately` | Secret key for dashboard session management. **Critical to change.** |
| `API_TOKEN`               | `change_me_api_token` | API token for authentication with the M3TAL API. **Critical to change.** |
| `ADMIN_PASSWORD`          | `admin_pass`          | Default password for the dashboard admin user. **Critical to change.** |
| `NETWORK_NAME`            | `m3tal`               | Name of the Docker network M3TAL services connect to.                |
| `LOCAL_IP`                | `127.0.0.1`           | Local IP address of the host.                                        |
| `DOMAIN`                  | `localhost`           | Base domain for Traefik-routed services.                             |
| `VPN_USER`                | `user`                | Username for VPN services (if integrated).                           |
| `VPN_PASSWORD`            | `password`            | Password for VPN services (if integrated).                           |
| `BASE_STORAGE_PATH`       | `./data`              | Base directory for all data storage.                                 |
| `MEDIA_PATH`              | `./data/media`        | Path for media files.                                                |
| `CONFIG_PATH`             | `./data/config`       | Path for configuration files.                                        |
| `DOWNLOADS_PATH`          | `./data/downloads`    | Path for downloads.                                                  |
| `PUID`                    | `1000`                | User ID for container processes.                                     |
| `PGID`                    | `1000`                | Group ID for container processes.                                    |
| `TZ`                      | `America/Denver`      | Timezone setting for containers.                                     |
| `TRAEFIK_WEB_PORT`        | `80`                  | Traefik's public HTTP entrypoint port.                               |
| `TRAEFIK_WEBHTTPS_PORT`   | `443`                 | Traefik's public HTTPS entrypoint port.                              |
| `TRAEFIK_DASHBOARD_PORT`  | `8080`                | Traefik's internal dashboard port.                                   |
| `DEBUG_MODE`              | `false`               | Enable debug logging and features.                                   |
| `METRICS_ENABLED`         | `true`                | Enable system metrics collection.                                    |