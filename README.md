# M3TAL System Documentation

## Overview
This document provides technical details and operational procedures for the M3TAL system. It describes the core components, their interactions, deployment mechanisms, and management interfaces.

## Prerequisites
Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation
The M3TAL CLI and API daemon are distributed as a `.deb` package and can be installed via APT.

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
The following table outlines the key directories and files used by the M3TAL system:

| Path                        | Purpose                                                                |
|-----------------------------|------------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file for M3TAL. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.       |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |

## Components
The M3TAL system comprises several key components:

-   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing a single entrypoint for all M3TAL operations, including configuration management, service control, and stack deployment.
-   **API daemon** (`m3tal-api.service`): A Go binary running as a `systemd` service on host-local port `8080`. It manages Docker interactions, the SQLite state database, and exposes API routes for system control.
-   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container running internally on port `8082`. It communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
-   **Traefik Gateway** (`routing-compose.yml`): A reverse proxy container designed to expose services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
-   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing zero-configuration secure internet access to services.

## Deployment Lifecycle
M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/` for all stack operations.

### Adding a New Stack
To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these.
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access
The M3TAL Dashboard provides a web-based interface for system management. Its access method is controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`, which can be either `local` (default) or `traefik`.

### Local Mode (Default: `DASHBOARD_EXPOSE_MODE=local`)
In this mode, the dashboard container binds its internal port `8082` directly to the host's `8082` port.
-   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` is set in `/etc/m3tal/.env`. The `m3tal-compose.local.yml` override file is used, which specifies a direct port binding.
-   **Access:** A new user performing a default installation will access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`. No Traefik configuration is required. This mode is suitable for LAN-only setups or initial local testing.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)
In this mode, the dashboard is routed via the Traefik Gateway using a domain name.
-   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik` is set in `/etc/m3tal/.env`. The `m3tal-compose.traefik.yml` override file is used, which applies Traefik labels to the dashboard service.
-   **Access:** The dashboard is accessible via `http://dash.DOMAIN`, where `DOMAIN` is defined in `/etc/m3tal/.env`. This requires Traefik to be running via `m3tal up` and the `proxy` network to be configured. This mode is suitable for domain-based setups with multiple services behind a reverse proxy.

## Traefik Gateway
Traefik is deployed as a container via `routing-compose.yml`. It functions as a reverse proxy, automatically discovering and routing traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds host port `80` (and `443` if HTTPS is configured) as its primary HTTP entry point. It also loads dynamic configuration files from `/docker/dynamic/` (which is symlinked to `/opt/m3tal/stack/dynamic/`) using its file provider, enabling hot-reloading of routing rules.

-   **API Routing:** Dynamic configuration files, such as `/docker/dynamic/api.yml`, are used to route requests to services listening on host-local ports. For example, `api.DOMAIN` is routed to the Go API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`.
-   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` is routed to the `m3tal-dashboard` container on its internal port `8082` via Traefik labels specified in `m3tal-compose.traefik.yml`.

### Exposing a Custom Service via Traefik
To expose a custom user service through Traefik, add the necessary Traefik labels to its service definition in your Docker Compose file:

```yaml
# In your custom-stack-compose.yml located in /docker/
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN}`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy # Ensure your service is on the 'proxy' network for Traefik discovery

networks:
  proxy:
    external: true # Assumes the 'proxy' network is external and created by M3TAL
```
After placing this file in `/docker/` and running `m3tal up`, your `my-app` service will be accessible via `http://app.DOMAIN`.

## Service Management
The M3TAL API daemon (`m3tal-api.service`) runs as a `systemd` service and can be managed using standard `systemctl` commands:

-   **Check status:** `systemctl status m3tal-api`
-   **Restart service:** `systemctl restart m3tal-api`
-   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo
-   To start the M3TAL Dashboard container specifically, isolating its deployment from other stacks:
    ```bash
    m3tal dash up
    ```
    This command will download the necessary dashboard compose files, apply the correct `DASHBOARD_EXPOSE_MODE` override, and start the dashboard.

-   To orchestrate and deploy all other stacks, including any user-defined compose files placed in the `/docker/` directory, along with core M3TAL services like Traefik:
    ```bash
    m3tal up
    ```
    This command will process all `*-compose.yml` files in `/docker/` and bring up their respective services.

## Port Map
The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                    | Access                                        | Description                                                          |
| :--- | :------------------------- | :-------------------------------------------- | :------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point   | Public                                        | The public-facing HTTP port for services exposed via Traefik.        |
| 8080 | M3TAL API daemon (Go)      | Host-local                                    | The internal port the M3TAL API daemon listens on.                   |
| 8081 | Traefik dashboard          | Host-local only                               | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard            | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Configuration
M3TAL uses `/etc/m3tal/.env` as its primary configuration file. This file contains environment variables that influence the behavior of the API daemon and Docker Compose stacks.

You can manage this configuration using the M3TAL CLI:
-   **Interactive wizard:** `m3tal config wizard`
-   **Set a specific key:** `m3tal config set KEY value`

### Key Environment Variables
The following environment variables are commonly used within the M3TAL ecosystem:

| Key                         | Default Value             | Description                                                   |
| :-------------------------- | :------------------------ | :------------------------------------------------------------ |
| `DASHBOARD_PORT`            | `8082`                    | Internal port for the M3TAL Dashboard container.              |
| `DASHBOARD_EXPOSE_MODE`     | `local`                   | Defines how the dashboard is exposed (`local` or `traefik`).  |
| `HTTP_PORT`                 | `8080`                    | Port the M3TAL API daemon listens on.                         |
| `STATE_DIR`                 | `./state`                 | Directory for the SQLite state database.                      |
| `LOG_LEVEL`                 | `info`                    | Logging verbosity for M3TAL components.                       |
| `DASHBOARD_SECRET`          | `change_me_immediately`   | Secret key for dashboard session management. **Critical to change.** |
| `API_TOKEN`                 | `change_me_api_token`     | Authentication token for the M3TAL API. **Critical to change.** |
| `ADMIN_PASSWORD`            | `admin_pass`              | Initial administrator password for the dashboard. **Critical to change.** |
| `NETWORK_NAME`              | `m3tal`                   | Name of the default Docker network created by M3TAL.          |
| `LOCAL_IP`                  | `127.0.0.1`               | Local IP address used in certain configurations.              |
| `DOMAIN`                    | `localhost`               | Base domain for Traefik routing (e.g., `api.DOMAIN`).         |
| `VPN_USER`                  | `user`                    | Username for VPN services (if used).                          |
| `VPN_PASSWORD`              | `password`                | Password for VPN services (if used).                          |
| `BASE_STORAGE_PATH`         | `./data`                  | Base directory for all M3TAL data volumes.                    |
| `MEDIA_PATH`                | `./data/media`            | Subdirectory for media files.                                 |
| `CONFIG_PATH`               | `./data/config`           | Subdirectory for configuration files.                         |
| `DOWNLOADS_PATH`            | `./data/downloads`        | Subdirectory for downloads.                                   |
| `PUID`                      | `1000`                    | User ID for container processes.                              |
| `PGID`                      | `1000`                    | Group ID for container processes.                             |
| `TZ`                        | `America/Denver`          | Timezone setting for containers.                              |
| `TRAEFIK_WEB_PORT`          | `80`                      | Traefik HTTP entry point port.                                |
| `TRAEFIK_WEBHTTPS_PORT`     | `443`                     | Traefik HTTPS entry point port.                               |
| `TRAEFIK_DASHBOARD_PORT`    | `8080`                    | Internal Traefik dashboard port (host-local `8081` is usually proxied to this). |
| `DEBUG_MODE`                | `false`                   | Enables debug logging and features.                           |
| `METRICS_ENABLED`           | `true`                    | Enables metrics collection.                                   |