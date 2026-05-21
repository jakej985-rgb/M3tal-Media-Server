# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

The M3TAL CLI binary and API daemon are distributed via an APT repository.

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

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and exposes API routes for M3TAL operations.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container, internally listening on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy Docker container that exposes services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel Docker container for establishing secure, zero-config internet access to services.

## Filesystem Contract

| Path                         | Purpose                                                     |
| :--------------------------- | :---------------------------------------------------------- |
| `/etc/m3tal/.env`            | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`    | SQLite state database. Auto-created by the API daemon.      |
| `/opt/m3tal/stack/`          | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                    | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`         | Dashboard credential store. Managed by `m3tal dashpass`.    |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth where all M3TAL-managed stack files and Traefik configuration reside. The `/docker` directory is a symlink to `/opt/m3tal/stack/`, making `/docker` the user-facing alias for all stack operations, allowing users to easily interact with compose files.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard has two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`) - Default

This is the default configuration for a new installation. The M3TAL Dashboard container exposes its internal port `8082` directly onto the host machine.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (uses `m3tal-compose.local.yml` override).
*   **Access:** Direct port binding at `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik required.
*   **Use Case:** Ideal for LAN-only setups, first-time users, and local development/testing. A new user performing a default installation will access the dashboard directly via port `8082`.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

This mode integrates the M3TAL Dashboard with the Traefik gateway for domain-based access.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik` (uses `m3tal-compose.traefik.yml` override).
*   **Access:** Domain routing at `http://dash.DOMAIN` via Traefik labels.
*   **Requirements:** Traefik must be running (`m3tal up` deploys Traefik).
*   **Use Case:** Suited for domain-based deployments and environments where multiple services are exposed through a single reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds to host port `80` (HTTP) and is configured to load dynamic configuration files from `/docker/dynamic/` (which is mapped to `/etc/traefik/dynamic` inside the container). This allows for hot-reloading of routing rules.

An example of dynamic configuration is routing `api.DOMAIN` to the M3TAL API daemon. Since the API daemon runs directly on the host on port `8080` (not within a Docker network accessible by name), a dynamic configuration file (e.g., `/docker/dynamic/api.yml`) is used to route `api.DOMAIN` to `http://host.docker.internal:8080`, directing requests to the host-local API daemon. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` is routed to the dashboard container via its Traefik labels.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service (e.g., `my-app`) via Traefik, add the necessary labels to its service definition:

```yaml
# In /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app-nginx
    labels:
      - "traefik.enable=true" # Enable Traefik for this service
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)" # Route requests for app.DOMAIN
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Target the internal port of the service
      - "traefik.http.routers.myapp.entrypoints=web" # Use the 'web' entrypoint (port 80)
    networks:
      - proxy # Connect to the Traefik proxy network

networks:
  proxy:
    external: true # Ensure this network is external and named 'proxy'
```
After placing this file in `/docker/` and running `m3tal up`, `app.DOMAIN` will route to your `my-app` service.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`.

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get started and see the M3TAL Dashboard:

1.  **Start the Dashboard:**
    ```bash
    m3tal dash up
    ```
    This command specifically downloads and starts the M3TAL Dashboard container using the appropriate `docker compose` override based on your `DASHBOARD_EXPOSE_MODE` setting.
    If using the default `DASHBOARD_EXPOSE_MODE=local`, you can then access the dashboard directly at `http://HOST_IP:8082`.

2.  **Deploy all M3TAL Stacks:**
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all `*-compose.yml` files found in the `/docker/` directory. This includes core M3TAL components like Traefik and Cloudflared (if configured), as well as any user-defined compose files you have placed in `/docker/`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                             |
| :--- | :-------------------------- | :------------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Environment Variables

The M3TAL system relies on a set of environment variables for configuration, primarily managed through `/etc/m3tal/.env`.

| Key                     | Default                 | Description                                                                 |
| :---------------------- | :---------------------- | :-------------------------------------------------------------------------- |
| `DASHBOARD_PORT`        | `8082`                  | Port for the M3TAL Dashboard container.                                     |
| `DASHBOARD_EXPOSE_MODE` | `local`                 | Determines dashboard exposure: `local` (direct port) or `traefik` (domain). |
| `HTTP_PORT`             | `8080`                  | Port for the M3TAL API daemon.                                              |
| `STATE_DIR`             | `./state`               | Directory for the SQLite state database.                                    |
| `LOG_LEVEL`             | `info`                  | Logging level for the API daemon.                                           |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret key for dashboard session management.                                |
| `API_TOKEN`             | `change_me_api_token`   | Authentication token for API access.                                        |
| `ADMIN_PASSWORD`        | `admin_pass`            | Default administrator password for the dashboard.                           |
| `NETWORK_NAME`          | `m3tal`                 | Name of the Docker network used by M3TAL services.                          |
| `LOCAL_IP`              | `127.0.0.1`             | Local IP address (used for some internal bindings).                         |
| `DOMAIN`                | `localhost`             | Primary domain for Traefik routing (e.g., `dash.DOMAIN`).                   |
| `VPN_USER`              | `user`                  | VPN username (if applicable).                                               |
| `VPN_PASSWORD`          | `password`              | VPN password (if applicable).                                               |
| `BASE_STORAGE_PATH`     | `./data`                | Base directory for all M3TAL data storage.                                  |
| `MEDIA_PATH`            | `./data/media`          | Path for media storage.                                                     |
| `CONFIG_PATH`           | `./data/config`         | Path for configuration files.                                               |
| `DOWNLOADS_PATH`        | `./data/downloads`      | Path for downloads.                                                         |
| `PUID`                  | `1000`                  | User ID for container processes.                                            |
| `PGID`                  | `1000`                  | Group ID for container processes.                                           |
| `TZ`                    | `America/Denver`        | Timezone for containers.                                                    |
| `TRAEFIK_WEB_PORT`      | `80`                    | Traefik's public HTTP entry point port.                                     |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                   | Traefik's public HTTPS entry point port.                                    |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                  | Internal port for the Traefik dashboard.                                    |
| `DEBUG_MODE`            | `false`                 | Enable debug logging.                                                       |
| `METRICS_ENABLED`       | `true`                  | Enable system metrics collection.                                           |