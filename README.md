# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Overview

M3TAL is a system composed of a Go-based CLI, a Go API daemon, a Python/Flask dashboard, and Dockerized services including Traefik for routing and optional Cloudflared for tunneling. It uses Docker Compose V2 for container orchestration and systemd for API daemon management. System configuration is managed via an `.env` file, and state is stored in an SQLite database.

## Installation

To install the M3TAL CLI binary and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The following table details the primary paths and their purposes within the M3TAL system:

| Path                        | Purpose                                                                             |
|-----------------------------|-------------------------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the API daemon.         |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose files and Traefik dynamic configurations. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                            |

## Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the primary command-line interface for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, system state within `/var/lib/m3tal/state.db`, and exposes M3TAL's API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container. It runs internally on port `8082` and communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy. It exposes services via domain names on host port `80` and uses a file provider for dynamic configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for establishing secure, zero-config internet access to services.

## Service Management

The M3TAL API daemon is managed as a systemd service, `m3tal-api.service`. Standard systemctl commands apply:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all M3TAL-managed and user-defined stack files reside. For user convenience and direct interaction, `/docker` is provided as a symlink alias to `/opt/m3tal/stack/`. All user-facing stack operations should target the `/docker/` directory.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env`. This can be done using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Execute `m3tal up` to deploy your new stack along with all other M3TAL-managed services.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, configured by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (this is the default setting after a new installation).
*   **Mechanism:** This mode utilizes the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`8082:8082`).
*   **Access:** The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirement:** No Traefik gateway is required for this mode.
*   **User Experience:** A new user performing a default M3TAL installation will access the dashboard directly via port `8082`.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`.
*   **Mechanism:** This mode utilizes the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard service. These labels instruct Traefik to route requests for `dash.DOMAIN` to the dashboard container's internal port `8082`.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Requirement:** The Traefik gateway must be running (`m3tal up` typically deploys Traefik).

## Traefik Gateway

Traefik acts as the central reverse proxy for M3TAL, automatically discovering and routing traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also integrates with dynamic configuration files, typically located in `/docker/dynamic/` (which is symlinked from `/opt/m3tal/stack/dynamic/`). For example, the `dynamic/api.yml` file is used to route requests for `api.DOMAIN` to the Go API daemon, which listens on the host-local port `8080`. This is achieved by routing the request to `http://host.docker.internal:8080` from within the Traefik container. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` routes to the M3TAL Dashboard container via its defined Traefik labels.

### Exposing a Custom Service via Traefik

To expose a custom user service (e.g., `my-app`) through Traefik, add the necessary Traefik labels to its service definition in your Docker Compose file (`my-app-compose.yml`):

```yaml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-app
    restart: unless-stopped
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

After placing this `my-app-compose.yml` in `/docker/` and ensuring `DOMAIN` is set in `/etc/m3tal/.env`, run `m3tal up` to deploy the service and configure Traefik.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Quick Demo

To quickly get the M3TAL Dashboard up and running for initial setup or testing:

*   **Start the Dashboard:** Run `m3tal dash up`. This command specifically downloads the latest dashboard Compose files and starts only the M3TAL Dashboard container, respecting your `DASHBOARD_EXPOSE_MODE` setting.
*   **Start All Stacks:** To deploy all other M3TAL-managed services (like Traefik, Cloudflared) and any user-defined Docker Compose files placed in `/docker/`, execute `m3tal up`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Environment Variables

Key environment variables that configure the M3TAL system are managed in `/etc/m3tal/.env`:

| Key                     | Default Value          | Description                                                    |
|-------------------------|------------------------|----------------------------------------------------------------|
| `DASHBOARD_PORT`        | `8082`                 | Internal port for the M3TAL Dashboard container.               |
| `DASHBOARD_EXPOSE_MODE` | `local`                | Controls how the dashboard is exposed (`local` or `traefik`).  |
| `HTTP_PORT`             | `8080`                 | Port for the M3TAL API daemon (Go).                            |
| `STATE_DIR`             | `./state`              | Directory for the SQLite state database.                       |
| `LOG_LEVEL`             | `info`                 | Logging verbosity for M3TAL components.                        |
| `DASHBOARD_SECRET`      | `change_me_immediately`| Secret key for dashboard session management.                   |
| `API_TOKEN`             | `change_me_api_token`  | Token for M3TAL API authentication.                            |
| `ADMIN_PASSWORD`        | `admin_pass`           | Default password for dashboard admin user.                     |
| `NETWORK_NAME`          | `m3tal`                | Default Docker network name for M3TAL services.                |
| `LOCAL_IP`              | `127.0.0.1`            | Local IP address for internal routing.                         |
| `DOMAIN`                | `localhost`            | Base domain for Traefik routing (e.g., `example.com`).         |
| `VPN_USER`              | `user`                 | Username for VPN configurations.                               |
| `VPN_PASSWORD`          | `password`             | Password for VPN configurations.                               |
| `BASE_STORAGE_PATH`     | `./data`               | Base directory for all data storage.                           |
| `MEDIA_PATH`            | `./data/media`         | Path for media storage.                                        |
| `CONFIG_PATH`           | `./data/config`        | Path for configuration files.                                  |
| `DOWNLOADS_PATH`        | `./data/downloads`     | Path for downloaded content.                                   |
| `PUID`                  | `1000`                 | User ID for container processes.                               |
| `PGID`                  | `1000`                 | Group ID for container processes.                              |
| `TZ`                    | `America/Denver`       | Timezone for containers.                                       |
| `TRAEFIK_WEB_PORT`      | `80`                   | Traefik's HTTP entry point port.                               |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                  | Traefik's HTTPS entry point port.                              |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                 | Internal Traefik dashboard port (host-local access only).      |
| `DEBUG_MODE`            | `false`                | Enables debug logging and features.                            |
| `METRICS_ENABLED`       | `true`                 | Enables collection and exposure of system metrics.             |