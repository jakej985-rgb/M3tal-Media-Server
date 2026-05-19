# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

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

## Core Components

The M3TAL system comprises the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary operating as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes the M3TAL API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container that runs internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container providing a Cloudflare tunnel for secure, zero-configuration internet access to services.

## Filesystem Contract

The following table details the critical directories and files within the M3TAL system:

| Path | Purpose |
| :----------------------- | :---------------------------------------------------------------------- |
| `/etc/m3tal/.env` | Primary configuration file for environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database, automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/` | The canonical directory for all Docker Compose stack files and Traefik dynamic configurations. |
| `/docker` | A symlink to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations. |
| `/docker/users.json` | The credential store for the M3TAL Dashboard. Managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

It is important to note that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside. The `/docker/` directory is a symlink alias to `/opt/m3tal/stack/`, providing a convenient user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all deployed stacks, including your newly added stack.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**: This is the default setting upon installation.
*   A direct port binding is established for the dashboard container (e.g., `8082:8082`).
*   **Access Method**: `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Behavior**: For a new user performing a default installation, the dashboard is directly accessible via port `8082`. No Traefik configuration is required. This mode is suitable for LAN-only setups, first-time users, and local testing environments.

### Traefik Mode

*   **`DASHBOARD_EXPOSE_MODE=traefik`**:
*   The dashboard container is configured with Traefik labels, allowing Traefik to route incoming requests for `dash.DOMAIN` to the dashboard's internal port `8082`.
*   **Access Method**: `http://dash.DOMAIN` (requires Traefik to be running via `m3tal up`).
*   **Behavior**: This mode leverages Traefik for domain-based routing, ideal for environments with multiple services behind a reverse proxy and publicly accessible domains.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Traefik Gateway

Traefik is deployed as a Docker container via `routing-compose.yml` and serves as the M3TAL system's reverse proxy. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files, typically located in `/docker/dynamic/` (symlinked from `/opt/m3tal/stack/dynamic/`). These files allow for routing requests to services that may not be directly Docker containers, such as the M3TAL Go API daemon. For example, the `dynamic/api.yml` file configures Traefik to route `api.DOMAIN` to the Go API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container when `DASHBOARD_EXPOSE_MODE=traefik` is set, leveraging the Traefik labels defined in `m3tal-compose.traefik.yml`.

### Example: Exposing a Custom User Service via Traefik

To expose a hypothetical user service (e.g., an Nginx container) via Traefik, you would add appropriate labels to its Docker Compose service definition:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
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
In this example, `app.DOMAIN` would route to the `my-app` service. The `proxy` network must be created externally for Traefik to connect to services using it.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

*   To start only the M3TAL Dashboard container, use:
    ```bash
    m3tal dash up
    ```
    (This command automatically fetches the necessary compose files and applies the correct `DASHBOARD_EXPOSE_MODE` override.)

*   To orchestrate and deploy all other Docker Compose stacks, including any user-defined compose files placed in the `/docker/` directory, use:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :--------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Configuration Variables

M3TAL uses environment variables for its configuration, primarily managed through the `/etc/m3tal/.env` file. These can be adjusted using `m3tal config wizard` or `m3tal config set KEY value`. Below is a list of common configuration variables and their default values:

| Key | Default Value | Description |
| :---------------------- | :-------------------- | :------------------------------------------------------------------------------------------ |
| `DASHBOARD_PORT` | `8082` | The internal port on which the M3TAL Dashboard container listens. |
| `DASHBOARD_EXPOSE_MODE` | `local` | Controls how the dashboard is exposed (`local` for direct port, `traefik` for Traefik routing). |
| `HTTP_PORT` | `8080` | The port for the M3TAL API daemon. |
| `STATE_DIR` | `./state` | Directory for the SQLite state database. |
| `LOG_LEVEL` | `info` | Logging verbosity for the API daemon (e.g., `info`, `debug`, `warn`). |
| `DASHBOARD_SECRET` | `change_me_immediately` | Secret key for the dashboard for session management. **MUST be changed.** |
| `API_TOKEN` | `change_me_api_token` | API token for authentication with the M3TAL API. **MUST be changed.** |
| `ADMIN_PASSWORD` | `admin_pass` | Default admin password for the dashboard. **MUST be changed.** |
| `NETWORK_NAME` | `m3tal` | The Docker network name used by M3TAL services. |
| `LOCAL_IP` | `127.0.0.1` | Local IP address for host-internal communication. |
| `DOMAIN` | `localhost` | The base domain used for Traefik routing rules. |
| `VPN_USER` | `user` | Default VPN user, if VPN components are enabled. |
| `VPN_PASSWORD` | `password` | Default VPN password, if VPN components are enabled. |
| `BASE_STORAGE_PATH` | `./data` | Base path for all M3TAL data storage. |
| `MEDIA_PATH` | `./data/media` | Path for media files. |
| `CONFIG_PATH` | `./data/config` | Path for configuration files. |
| `DOWNLOADS_PATH` | `./data/downloads` | Path for downloads. |
| `PUID` | `1000` | User ID for container processes to ensure correct file permissions. |
| `PGID` | `1000` | Group ID for container processes to ensure correct file permissions. |
| `TZ` | `America/Denver` | Timezone setting for all M3TAL components. |
| `TRAEFIK_WEB_PORT` | `80` | Public HTTP port for Traefik. |
| `TRAEFIK_WEBHTTPS_PORT` | `443` | Public HTTPS port for Traefik. |
| `TRAEFIK_DASHBOARD_PORT` | `8080` | Internal port for the Traefik dashboard. |
| `DEBUG_MODE` | `false` | Enables debug mode for various components. |
| `METRICS_ENABLED` | `true` | Enables metrics collection. |