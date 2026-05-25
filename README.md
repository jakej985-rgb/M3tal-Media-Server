# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

The M3TAL CLI binary is installed via APT.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Configuration

M3TAL's primary configuration is managed via the `/etc/m3tal/.env` file, which is populated and maintained using the `m3tal config wizard` and `m3tal config set KEY value` commands.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth for all stack definition files. The `/docker/` directory is a user-facing symlink alias to `/opt/m3tal/stack/`, simplifying stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these variables.
3.  Execute `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

*   **Local Mode (default)**: When `DASHBOARD_EXPOSE_MODE=local`, the dashboard is directly exposed on the host's IP address and a configurable port (defaulting to `8082`). Access is granted via `http://HOST_IP:8082` or `http://localhost:8082`. This mode does not require Traefik to be running and is ideal for local development or simple home server setups. A new user performing a default installation will access the dashboard directly via port 8082.

*   **Traefik Mode**: When `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard is routed through the Traefik reverse proxy. Access is granted via a domain name, typically `http://dash.DOMAIN`, provided Traefik is running and configured to route to the dashboard container. This mode requires Traefik to be operational.

## Traefik Gateway

Traefik acts as the primary reverse proxy for M3TAL, managing external access to containerized services. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik dynamically configures its routing rules based on files located in the `/opt/m3tal/stack/dynamic/` directory. This allows for granular control over inbound requests. For example, the `dynamic/api.yml` file is used to route requests from `api.DOMAIN` to the M3TAL API daemon, which listens on host-local port `8080`, via `http://host.docker.internal:8080`. The dashboard is routed via `dash.DOMAIN` when `DASHBOARD_EXPOSE_MODE=traefik`.

Here is a concrete YAML example of how to expose a custom user service via Traefik labels in a hypothetical `my-app-compose.yml` file:

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
```

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon is managed by systemd as `m3tal-api.service`. You can control and monitor its status using the following `systemctl` commands:

*   View service status: `systemctl status m3tal-api.service`
*   Restart the service: `systemctl restart m3tal-api.service`
*   View service logs: `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly start the M3TAL dashboard container without deploying other stacks:

```bash
m3tal dash up
```

The `m3tal up` command orchestrates and deploys all other stacks defined in the `/docker/` directory, including any user-defined compose files.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port   | Service                                   | Access                                                    | Description                                                                                                 |
| :----- | :---------------------------------------- | :-------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------- |
| 80     | Traefik HTTP entry point                  | Public (when using Traefik mode)                          | The public-facing HTTP port for services exposed via Traefik.                                               |
| 8080   | M3TAL API daemon (Go)                     | Host-local                                                | The internal port the M3TAL API daemon listens on.                                                          |
| 8081   | Traefik dashboard                         | Host-local only                                           | The internal Traefik dashboard port, accessible only from the host machine.                                 |
| 8082   | M3TAL Dashboard                           | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Filesystem Contract

| Path                           | Purpose                                                                                             |
| :----------------------------- | :-------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`              | Primary configuration file. Managed by `m3tal config wizard`.                                       |
| `/var/lib/m3tal/state.db`      | SQLite state database. Auto-created by the API daemon.                                            |
| `/opt/m3tal/stack/`            | Canonical stack directory. Contains compose files and Traefik configuration.                        |
| `/docker`                      | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations.               |
| `/docker/users.json`           | Dashboard credential store. Managed by `m3tal dashpass`.                                            |
| `/opt/m3tal/stack/dynamic/`    | Directory for Traefik dynamic configuration files.                                                  |
| `/opt/m3tal/stack/routing-compose.yml` | Docker Compose file for deploying Traefik and Cloudflared.                                      |
| `/opt/m3tal/stack/m3tal-compose.yml` | Base Docker Compose file for the M3TAL dashboard.                                               |
| `/opt/m3tal/stack/m3tal-compose.local.yml` | Docker Compose override for local dashboard access mode.                                        |
| `/opt/m3tal/stack/m3tal-compose.traefik.yml` | Docker Compose override for Traefik dashboard access mode.                                      |

## Environment Variables

The following environment variables are managed by M3TAL and can be configured via `/etc/m3tal/.env`:

| Key                    | Default        | Description                                                             |
| :--------------------- | :------------- | :---------------------------------------------------------------------- |
| `DASHBOARD_PORT`       | `8082`         | The internal port the M3TAL Dashboard container listens on.             |
| `DASHBOARD_EXPOSE_MODE`| `local`        | Determines how the dashboard is exposed (`local` or `traefik`).         |
| `HTTP_PORT`            | `8080`         | The port the M3TAL API daemon listens on.                               |
| `STATE_DIR`            | `./state`      | Directory for the M3TAL state database.                                 |
| `LOG_LEVEL`            | `info`         | The logging level for M3TAL services.                                   |
| `DASHBOARD_SECRET`     | `change_me_immediately` | Secret key for dashboard session management.                            |
| `API_TOKEN`            | `change_me_api_token`  | API token for programmatic access.                                      |
| `ADMIN_PASSWORD`       | `admin_pass`   | Default password for the M3TAL dashboard administrator.                 |
| `NETWORK_NAME`         | `m3tal`        | The name of the Docker network M3TAL services will use.                 |
| `LOCAL_IP`             | `127.0.0.1`    | The local IP address to bind to.                                        |
| `DOMAIN`               | `localhost`    | The domain name used for routing when `DASHBOARD_EXPOSE_MODE=traefik`.  |
| `VPN_USER`             | `user`         | Username for VPN configuration.                                         |
| `VPN_PASSWORD`         | `password`     | Password for VPN configuration.                                         |
| `BASE_STORAGE_PATH`    | `./data`       | Base path for persistent data storage.                                  |
| `MEDIA_PATH`           | `./data/media` | Path for media storage.                                                 |
| `CONFIG_PATH`          | `./data/config`| Path for configuration files.                                           |
| `DOWNLOADS_PATH`       | `./data/downloads` | Path for download storage.                                              |
| `PUID`                 | `1000`         | User ID for container processes.                                        |
| `PGID`                 | `1000`         | Group ID for container processes.                                       |
| `TZ`                   | `America/Denver`| Timezone for container processes.                                       |
| `TRAEFIK_WEB_PORT`     | `80`           | The host port Traefik uses for HTTP entry points.                       |
| `TRAEFIK_WEBHTTPS_PORT`| `443`          | The host port Traefik uses for HTTPS entry points (if configured).      |
| `TRAEFIK_DASHBOARD_PORT`| `8080`         | The internal port Traefik listens on for its dashboard.                 |
| `DEBUG_MODE`           | `false`        | Enables or disables debug mode for M3TAL services.                      |
| `METRICS_ENABLED`      | `true`         | Enables or disables metrics collection for M3TAL services.              |
