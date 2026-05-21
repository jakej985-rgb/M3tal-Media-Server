# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI and API daemon:

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

## System Components

The M3TAL system is composed of the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary that serves as the single entrypoint for all M3TAL operations, installed via APT.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, enabling zero-configuration internet access for services.

## Filesystem Contract

The M3TAL system relies on a specific filesystem structure for its operation and configuration:

| Path                        | Purpose                                                                                |
| :-------------------------- | :------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | The primary system-wide configuration file, managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`   | The SQLite state database, automatically created and managed by the API daemon.        |
| `/opt/m3tal/stack/`         | The canonical directory containing all Docker Compose stack files and Traefik config.  |
| `/docker`                   | A symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | The Dashboard credential store, managed by `m3tal dashpass`.                           |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

It is crucial to understand that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside. The `/docker` directory is a user-facing symlink alias for `/opt/m3tal/stack/`, intended for all user interactions with stack files.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all deployed stacks, including your newly added service.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** Uses the `m3tal-compose.local.yml` override, which adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082`.
*   **Access:** Directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik configuration is required.
*   **Usage:** A new user performing a default M3TAL installation will access the dashboard directly via port 8082. This mode is ideal for LAN-only setups, first-time users, and local testing.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard container definition. These labels instruct Traefik to route `dash.${DOMAIN}` to the dashboard container on its internal port `8082`.
*   **Access:** Via a domain at `http://dash.DOMAIN`.
*   **Requirements:** Traefik must be running via `m3tal up` for this routing to function.
*   **Usage:** Best suited for domain-based setups and environments where multiple services are exposed behind a single reverse proxy.

## Traefik Gateway

Traefik acts as the M3TAL system's reverse proxy, automatically discovering and routing traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files (such as `dynamic/api.yml` located in `/docker/dynamic/`). For example, `dynamic/api.yml` defines how requests to `api.DOMAIN` are routed to the M3TAL API daemon, which listens on the host-local port `8080`. This routing is achieved by directing traffic to `http://host.docker.internal:8080` from within the Traefik container. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` routes directly to the dashboard container via its internal port `8082`.

### Exposing a Custom User Service via Traefik

To expose your own Docker Compose service through Traefik, you need to add appropriate Traefik labels to its service definition. Ensure your service is part of the `proxy` network, which Traefik uses for discovery.

Here is a concrete YAML example for a hypothetical `my-app-compose.yml` that deploys an Nginx container and exposes it via Traefik at `app.DOMAIN`:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-nginx-app
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy
    restart: unless-stopped

networks:
  proxy:
    external: true
```

After placing this file in `/docker/` and running `m3tal up`, your Nginx app would be accessible at `http://app.DOMAIN` (assuming `DOMAIN` is correctly set in `/etc/m3tal/.env` and Traefik is running).

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

To quickly start the M3TAL Dashboard:

1.  First, ensure Docker is running and `m3tal-api.service` is active:
    ```bash
    sudo systemctl start m3tal-api
    ```
2.  Run the specific command to start *only* the dashboard container:
    ```bash
    m3tal dash up
    ```
    This command will download the necessary dashboard compose files and start the dashboard container using the appropriate expose mode (defaulting to local mode).
3.  Access the dashboard at `http://localhost:8082` (or `http://YOUR_HOST_IP:8082`).

The `m3tal up` command, in contrast, orchestrates and deploys all other Docker Compose stacks found in the `/docker/` directory, including the `routing-compose.yml` (Traefik, Cloudflared) and any user-defined compose files you've added.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                             |
| :--- | :-------------------------- | :------------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Key Configuration Variables

The primary M3TAL configuration is managed via environment variables stored in `/etc/m3tal/.env`. These variables dictate various aspects of system behavior. You can use `m3tal config wizard` to set these or `m3tal config set KEY value`.

| Key                         | Default Value           | Description                                                                    |
| :-------------------------- | :---------------------- | :----------------------------------------------------------------------------- |
| `DASHBOARD_PORT`            | `8082`                  | The internal port the M3TAL Dashboard container listens on.                    |
| `DASHBOARD_EXPOSE_MODE`     | `local`                 | Controls how the Dashboard is exposed (local port binding or via Traefik).     |
| `HTTP_PORT`                 | `8080`                  | The host-local port the M3TAL API daemon listens on.                           |
| `STATE_DIR`                 | `./state`               | Directory for the SQLite state database.                                       |
| `LOG_LEVEL`                 | `info`                  | Sets the logging verbosity for the API daemon.                                 |
| `DASHBOARD_SECRET`          | `change_me_immediately` | Secret key for Dashboard session management. **MUST BE CHANGED.**              |
| `API_TOKEN`                 | `change_me_api_token`   | API token for authentication with the M3TAL API. **MUST BE CHANGED.**          |
| `ADMIN_PASSWORD`            | `admin_pass`            | Default administrator password for initial Dashboard access. **MUST BE CHANGED.** |
| `NETWORK_NAME`              | `m3tal`                 | Name of the Docker network used by M3TAL services.                             |
| `LOCAL_IP`                  | `127.0.0.1`             | Local IP address of the host.                                                  |
| `DOMAIN`                    | `localhost`             | The base domain name used for Traefik routing (e.g., `dash.DOMAIN`).           |
| `VPN_USER`                  | `user`                  | Placeholder for VPN username (if applicable).                                  |
| `VPN_PASSWORD`              | `password`              | Placeholder for VPN password (if applicable).                                  |
| `BASE_STORAGE_PATH`         | `./data`                | Base directory for all data volumes.                                           |
| `MEDIA_PATH`                | `./data/media`          | Path for media storage.                                                        |
| `CONFIG_PATH`               | `./data/config`         | Path for configuration files.                                                  |
| `DOWNLOADS_PATH`            | `./data/downloads`      | Path for downloads storage.                                                    |
| `PUID`                      | `1000`                  | User ID for container processes.                                               |
| `PGID`                      | `1000`                  | Group ID for container processes.                                              |
| `TZ`                        | `America/Denver`        | Timezone for containers.                                                       |
| `TRAEFIK_WEB_PORT`          | `80`                    | Public HTTP port for Traefik.                                                  |
| `TRAEFIK_WEBHTTPS_PORT`     | `443`                   | Public HTTPS port for Traefik.                                                 |
| `TRAEFIK_DASHBOARD_PORT`    | `8080`                  | Internal port for Traefik's own dashboard.                                     |
| `DEBUG_MODE`                | `false`                 | Enables debug logging for M3TAL components.                                    |
| `METRICS_ENABLED`           | `true`                  | Enables metric collection for M3TAL components.                                |