# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

The M3TAL CLI binary is distributed via APT. Follow these steps to install:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Configuration

M3TAL's primary configuration is managed via an `.env` file. The `m3tal config wizard` command can be used to interactively set crucial variables. Alternatively, individual variables can be set using `m3tal config set KEY value`.

The primary configuration file is located at `/etc/m3tal/.env`. Key variables and their default values include:

| Key                     | Default         | Description                                       |
| ----------------------- | --------------- | ------------------------------------------------- |
| `DASHBOARD_PORT`        | `8082`          | The internal port for the M3TAL Dashboard.        |
| `DASHBOARD_EXPOSE_MODE` | `local`         | Determines dashboard access (`local` or `traefik`). |
| `HTTP_PORT`             | `8080`          | The port the M3TAL API daemon listens on.         |
| `STATE_DIR`             | `./state`       | Directory for the SQLite state database.          |
| `LOG_LEVEL`             | `info`          | Logging level for M3TAL services.                 |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret for dashboard authentication.              |
| `API_TOKEN`             | `change_me_api_token`   | Token for API authentication.                   |
| `ADMIN_PASSWORD`        | `admin_pass`    | Default password for dashboard administrator.     |
| `NETWORK_NAME`          | `m3tal`         | The Docker network name used by M3TAL.            |
| `LOCAL_IP`              | `127.0.0.1`     | The local IP address for host-based services.     |
| `DOMAIN`                | `localhost`     | The domain used for Traefik routing.              |
| `BASE_STORAGE_PATH`     | `./data`        | Base path for persistent data.                    |
| `CONFIG_PATH`           | `./data/config` | Path for persistent configuration files.          |
| `PUID`                  | `1000`          | User ID for container processes.                  |
| `PGID`                  | `1000`          | Group ID for container processes.                 |
| `TZ`                    | `America/Denver`| Timezone for container processes.                 |
| `TRAEFIK_WEB_PORT`      | `80`            | Host port for Traefik HTTP entrypoint.            |
| `TRAEFIK_WEBHTTPS_PORT` | `443`           | Host port for Traefik HTTPS entrypoint (if configured). |

## Filesystem Contract

The following paths are critical to M3TAL's operation and state management:

| Path                    | Purpose                                                     |
| ----------------------- | ----------------------------------------------------------- |
| `/etc/m3tal/.env`       | Primary configuration file, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database, auto-created by the API daemon.      |
| `/opt/m3tal/stack/`     | Canonical stack directory containing compose files and Traefik config. |
| `/docker`               | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`    | Dashboard credential store, managed by `m3tal dashpass`.    |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory serves as the user-facing alias for all stack operations. Internally, it is a symbolic link to `/opt/m3tal/stack/`, which is the canonical source of truth for all stack definitions.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` or set via `m3tal config set`.
3. Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard has two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

When `DASHBOARD_EXPOSE_MODE` is set to `local` (which is the default setting), the dashboard is accessed directly via its exposed port on the host machine. No Traefik is involved in this mode.

A new user performing a default installation will access the dashboard directly via port `8082`.

**Access URL:** `http://HOST_IP:8082` or `http://localhost:8082`

### Traefik Mode

When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, the dashboard is exposed via the Traefik gateway. This mode requires Traefik to be running and configured to route traffic based on domain names.

**Access URL:** `http://dash.DOMAIN` (where `DOMAIN` is defined in `/etc/m3tal/.env`)

## Traefik Gateway

Traefik acts as the reverse proxy and service discovery engine for M3TAL. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml`. Its static configuration (defined in `traefik.yml`) sets up entry points and enables providers.

```yaml
entryPoints:
  web:
    address: ":80"

providers:
  docker:
    exposedByDefault: false
    network: proxy
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

The `dynamic/` directory within `/opt/m3tal/stack/` is used for dynamic configuration files, allowing for hot-reloading of routing rules without restarting Traefik.

For instance, dynamic routing to the M3TAL API daemon is configured via `dynamic/api.yml`:

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
This configuration routes requests to `api.DOMAIN` to the M3TAL API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`. The `dash.DOMAIN` route is handled by Traefik labels defined in the dashboard's compose override file when `DASHBOARD_EXPOSE_MODE=traefik`.

### Exposing Custom User Services

To expose a custom user service via Traefik labels, add the appropriate labels to the service definition within its Docker Compose file. For example, to expose a hypothetical `my-app` service:

```yaml
# Example: my-app-compose.yml (place in /docker/)
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
In this example, `app.DOMAIN` will route to the `my-app` service, which is configured to listen on port `80`. The `proxy` network ensures Traefik can communicate with the service.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon is managed by systemd as `m3tal-api.service`. The following commands can be used to manage its lifecycle:

- **View status:** `systemctl status m3tal-api.service`
- **Restart service:** `systemctl restart m3tal-api.service`
- **View logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

This section outlines common operational commands for a quick demonstration of M3TAL.

### Starting the Dashboard

To start the M3TAL Dashboard container specifically, use the `m3tal dash up` command. This command handles downloading the appropriate compose files and applying overrides based on your `DASHBOARD_EXPOSE_MODE` setting.

```bash
m3tal dash up
```

### Deploying All Stacks

The `m3tal up` command orchestrates and deploys all stacks defined in the `/docker/` directory. This includes the dashboard stack (if not already managed by `m3tal dash up` in isolation) and any user-defined Docker Compose files placed in `/docker/`.

```bash
m3tal up
```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.