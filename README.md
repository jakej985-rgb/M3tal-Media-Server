# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

### APT Installation

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Configuration

M3TAL's primary configuration is managed through the `/etc/m3tal/.env` file. This file is automatically created and managed by the `m3tal config wizard` command or can be modified using `m3tal config set KEY value`.

The following environment variables are configurable:

| Key | Default | Description |
|---|---|---|
| `DASHBOARD_PORT` | `8082` | The internal port the M3TAL Dashboard container listens on. |
| `DASHBOARD_EXPOSE_MODE` | `local` | Controls dashboard access mode: `local` (direct port) or `traefik` (domain). |
| `HTTP_PORT` | `8080` | The internal port the M3TAL API daemon listens on. |
| `STATE_DIR` | `./state` | The directory for the SQLite state database. |
| `LOG_LEVEL` | `info` | The logging level for M3TAL services. |
| `DASHBOARD_SECRET` | `change_me_immediately` | Secret used for dashboard session management. **Must be changed.** |
| `API_TOKEN` | `change_me_api_token` | Token for API authentication. **Must be changed.** |
| `ADMIN_PASSWORD` | `admin_pass` | Password for the administrator user. **Must be changed.** |
| `NETWORK_NAME` | `m3tal` | The Docker network name used by M3TAL services. |
| `LOCAL_IP` | `127.0.0.1` | The IP address to bind services to. |
| `DOMAIN` | `localhost` | The domain name for Traefik routing (when `DASHBOARD_EXPOSE_MODE=traefik`). |
| `VPN_USER` | `user` | VPN username (if applicable). |
| `VPN_PASSWORD` | `password` | VPN password (if applicable). |
| `BASE_STORAGE_PATH` | `./data` | Base path for persistent data storage. |
| `MEDIA_PATH` | `./data/media` | Path for media storage. |
| `CONFIG_PATH` | `./data/config` | Path for configuration storage. |
| `DOWNLOADS_PATH` | `./data/downloads` | Path for download storage. |
| `PUID` | `1000` | User ID for container processes. |
| `PGID` | `1000` | Group ID for container processes. |
| `TZ` | `America/Denver` | Timezone for containers. |
| `TRAEFIK_WEB_PORT` | `80` | The host port Traefik listens on for HTTP. |
| `TRAEFIK_WEBHTTPS_PORT` | `443` | The host port Traefik listens on for HTTPS (if configured). |
| `TRAEFIK_DASHBOARD_PORT` | `8080` | The host port for the Traefik dashboard (host-local access only). |
| `DEBUG_MODE` | `false` | Enables debug logging and features. |
| `METRICS_ENABLED` | `true` | Enables Prometheus metrics collection. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack. The `/docker/` directory is a user-facing symlink alias for all stack operations, pointing to `/opt/m3tal/stack/`, which is the canonical source of truth directory where all stack files reside.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env`.
3. Run `m3tal up` to deploy your new stack along with all other M3TAL-managed stacks.

## Dashboard Access

The M3TAL Dashboard has two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`:

### Local Mode (Default)

- **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
- **Access:** Direct port binding at `http://HOST_IP:8082`.
- **Description:** In this default mode, the dashboard container is directly exposed on the host's IP address using the port specified by `DASHBOARD_PORT` (defaulting to 8082). No Traefik is required for dashboard access in this mode, making it ideal for simple local setups or initial installations. A new user performing a default installation will access the dashboard directly via port 8082.

### Traefik Mode

- **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
- **Access:** Domain routing at `http://dash.DOMAIN` via Traefik.
- **Description:** When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, the dashboard is exposed through the Traefik reverse proxy. Traefik is configured to route traffic for `dash.DOMAIN` to the dashboard container. This mode requires Traefik to be running and properly configured.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

The Traefik gateway is deployed via `routing-compose.yml` and exposes services on host port `80` (HTTP) and potentially `443` (HTTPS). Traefik uses a file provider to load dynamic configuration files from the `/docker/dynamic/` directory, allowing for hot-reloading of routing rules without restarting Traefik.

For example, the dynamic configuration file `dynamic/api.yml` routes requests for `api.DOMAIN` to the M3TAL API daemon, which is listening on host-local port `8080`:

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

This configuration allows external access to the API via `api.DOMAIN`, which is internally directed to the Go API daemon running on the host.

### Exposing a Custom User Service

To expose a custom user service via Traefik labels, add the appropriate labels to your service definition within its Docker Compose file located in the `/docker/` directory. For instance, to expose a hypothetical `my-app` service:

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
This configuration will make your `my-app` service accessible via `http://app.DOMAIN` if Traefik is configured to use the `proxy` network.

## Service Management — systemd

The M3TAL API daemon runs as a systemd service. You can manage it using the following commands:

- **Check Status:** `systemctl status m3tal-api`
- **Restart Service:** `systemctl restart m3tal-api`
- **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

1.  **Start the Dashboard:** To specifically start the dashboard container, use the `m3tal dash up` command. This command downloads the necessary compose files and applies the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.

2.  **Deploy All Stacks:** The `m3tal up` command orchestrates and deploys all other stacks defined in the `/docker/` directory, including any user-defined compose files you have added. This command ensures all your M3TAL services and any custom stacks are running.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|---|---|---|---|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Filesystem Contract

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Docker Services State

The following Docker services are managed by M3TAL:

```json
[
  {
    "name": "m3tal-dashboard",
    "image": "ghcr.io/jakej985-rgb/m3tal-godash:debug",
    "ports": [],
    "stack": "m3tal"
  },
  {
    "name": "traefik",
    "image": "traefik:latest",
    "ports": [
      "127.0.0.1:8081:8080"
    ],
    "stack": "routing"
  },
  {
    "name": "cloudflared",
    "image": "cloudflare/cloudflared:latest",
    "ports": [],
    "stack": "routing"
  }
]
```