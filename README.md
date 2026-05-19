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

The following table details key file system paths and their purposes within the M3TAL system:

| Path                     | Purpose                                                      |
| :----------------------- | :----------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon.       |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik configuration. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.     |

## Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing a single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on port `8080`. It manages Docker interactions, the state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container, internally on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for establishing zero-configuration internet access to services.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all Docker Compose stack files and associated configurations reside. The `/docker/` path is a symbolic link to `/opt/m3tal/stack/`, serving as the user-facing alias for all stack operations and file placement.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Execute `m3tal up` to deploy all stacks, including your newly added one.

## Quick Demo

*   To start *only* the M3TAL Dashboard container, use:
    ```bash
    m3tal dash up
    ```
    This command specifically manages the dashboard container: it downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub, reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`, and starts the dashboard with the appropriate Compose override file.
*   To orchestrate and deploy all other stacks in the `/docker/` directory, including any user-defined Docker Compose files, use:
    ```bash
    m3tal up
    ```

## Dashboard Access

The M3TAL Dashboard provides two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`). No Traefik configuration is involved.
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Default Behavior**: A new user performing a default M3TAL installation will access the dashboard directly via port `8082`. This behavior is directly linked to the default setting `DASHBOARD_EXPOSE_MODE=local`.
*   **Use Case**: Ideal for LAN-only setups, first-time users, or local testing environments where direct IP-based access is sufficient.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: This mode uses the `m3tal-compose.traefik.yml` override file, which adds Traefik labels to the dashboard service. These labels instruct Traefik to route traffic for `dash.DOMAIN` (where `DOMAIN` is configured in `/etc/m3tal/.env`) to the dashboard container on its internal port `8082`. This requires the Traefik gateway to be running.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Use Case**: Suited for domain-based setups where multiple services are exposed behind a central reverse proxy like Traefik.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Traefik Gateway

The Traefik gateway container is deployed via `routing-compose.yml` and serves as the primary reverse proxy for M3TAL services. It binds host port `80` as its HTTP entry point.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also loads dynamic configuration files from `/docker/dynamic/` (which symlinks to `/opt/m3tal/stack/dynamic/`). These files provide additional routing rules, allowing Traefik to route requests to services not directly managed by Docker Compose labels (e.g., host-local daemons). For example, a dynamic configuration file `dynamic/api.yml` routes `api.DOMAIN` to the M3TAL Go API daemon, which listens on the host-local port `8080`, via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the M3TAL dashboard container (when `DASHBOARD_EXPOSE_MODE=traefik`).

### Exposing a Custom User Service via Traefik

To expose your own Docker Compose service through Traefik, add the necessary Traefik labels to its service definition in your `*-compose.yml` file.

Here's a concrete example for a hypothetical `my-app-compose.yml` file exposing an Nginx container:

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
    external: true # Assumes 'proxy' network is created by M3TAL/Traefik
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will route requests for `app.YOUR_DOMAIN` to the `my-app` Nginx container.

## Service Management

The M3TAL API daemon is managed as a systemd service, `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check Status**: `systemctl status m3tal-api`
*   **Restart Service**: `systemctl restart m3tal-api`
*   **View Logs**: `journalctl -u m3tal-api -f`

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :---------------------------------------- | :------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                    | The public-facing HTTP port for services exposed via Traefik.                                                  |
| 8080 | M3TAL API daemon (Go)       | Host-local                                | The internal port the M3TAL API daemon listens on.                                                             |
| 8081 | Traefik dashboard           | Host-local only                           | The internal Traefik dashboard port, accessible only from the host machine.                                    |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Environment Variables

The M3TAL system relies on various environment variables for configuration, primarily defined in `/etc/m3tal/.env`. The table below lists commonly used variables and their default values:

| Key                       | Default Value         |
| :------------------------ | :-------------------- |
| `DASHBOARD_PORT`          | `8082`                |
| `DASHBOARD_EXPOSE_MODE`   | `local`               |
| `HTTP_PORT`               | `8080`                |
| `STATE_DIR`               | `./state`             |
| `LOG_LEVEL`               | `info`                |
| `DASHBOARD_SECRET`        | `change_me_immediately` |
| `API_TOKEN`               | `change_me_api_token` |
| `ADMIN_PASSWORD`          | `admin_pass`          |
| `NETWORK_NAME`            | `m3tal`               |
| `LOCAL_IP`                | `127.0.0.1`           |
| `DOMAIN`                  | `localhost`           |
| `VPN_USER`                | `user`                |
| `VPN_PASSWORD`            | `password`            |
| `BASE_STORAGE_PATH`       | `./data`              |
| `MEDIA_PATH`              | `./data/media`        |
| `CONFIG_PATH`             | `./data/config`       |
| `DOWNLOADS_PATH`          | `./data/downloads`    |
| `PUID`                    | `1000`                |
| `PGID`                    | `1000`                |
| `TZ`                      | `America/Denver`      |
| `TRAEFIK_WEB_PORT`        | `80`                  |
| `TRAEFIK_WEBHTTPS_PORT`   | `443`                 |
| `TRAEFIK_DASHBOARD_PORT`  | `8080`                |
| `DEBUG_MODE`              | `false`               |
| `METRICS_ENABLED`         | `true`                |