# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install M3TAL, follow these steps using the APT package manager:

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

The M3TAL system comprises several interconnected components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations, installed via APT.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application encapsulated in a Docker container, running internally on port `8082`. It communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container that provides a Cloudflare tunnel for zero-configuration internet access to services.

## Filesystem Contract

The following table details the critical directories and files within the M3TAL filesystem:

| Path                        | Purpose                                                                                                 |
| :-------------------------- | :------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`.                 |
| `/var/lib/m3tal/state.db`   | SQLite state database, used by the API daemon. Automatically created upon API daemon startup.           |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose files and Traefik dynamic configuration.                     |
| `/docker`                   | A symbolic link to `/opt/m3tal/stack/`. This is the user-facing path for all stack management operations. |
| `/docker/users.json`        | Credential store for the M3TAL Dashboard. Managed by `m3tal dashpass`.                                  |

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                    | Access                                        | Description                                                                                             |
| :--- | :------------------------- | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point   | Public                                        | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)      | Host-local                                    | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard          | Host-local only                               | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard            | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all stack files reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/` for all stack operations, simplifying access and management.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy your new stack along with all other existing M3TAL-managed stacks.

## Quick Demo

To quickly get started with the M3TAL Dashboard:

*   To start *only* the dashboard container specifically, use the command:
    ```bash
    m3tal dash up
    ```
    This command downloads the necessary dashboard compose files and starts the dashboard based on your `DASHBOARD_EXPOSE_MODE` setting.
*   To orchestrate and deploy all M3TAL stacks, including any user-defined compose files you have placed in the `/docker/` directory, use:
    ```bash
    m3tal up
    ```

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, configured via the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** This mode utilizes an override file (`m3tal-compose.local.yml`) to add a direct port binding, exposing the dashboard container's internal port `8082` directly on the host. For a new user performing a default installation, this is the out-of-the-box access method.
*   **Access:** `http://HOST_IP:8082` or `http://localhost:8082`
*   **Requirements:** No Traefik gateway is required for this mode. It operates independently.
*   **Use Case:** Ideal for LAN-only setups, first-time users, and local development or testing.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode utilizes an override file (`m3tal-compose.traefik.yml`) to add specific Traefik labels to the dashboard service. Traefik (if running via `m3tal up`) then intercepts requests for `dash.DOMAIN` and routes them internally to the dashboard container on port `8082`.
*   **Access:** `http://dash.DOMAIN`
*   **Requirements:** Traefik must be actively running and configured to serve the specified domain.
*   **Use Case:** Suited for domain-based environments where multiple services are exposed through a single reverse proxy.

## Traefik Gateway

The Traefik gateway automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also supports dynamic configuration files, typically located in `/docker/dynamic/` (which is mapped to `/etc/traefik/dynamic` inside the Traefik container). These files allow routing requests to services that might not be Docker containers themselves, or to specific host-local ports. For example, the `dynamic/api.yml` configuration routes requests for `api.DOMAIN` to the M3TAL Go API daemon listening on host-local port `8080` by directing traffic to `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container when `DASHBOARD_EXPOSE_MODE` is set to `traefik`.

### Exposing a Custom User Service via Traefik

To expose a custom user service (e.g., an Nginx container) via Traefik, you would add appropriate labels to its service definition within your Docker Compose file (e.g., `my-app-compose.yml`) located in `/docker/`:

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
      - proxy # Ensure your service is on the 'proxy' network for Traefik to discover it.
    # Optionally, if your application needs to connect to the M3TAL API:
    # extra_hosts:
    #   - "host.docker.internal:host-gateway"

networks:
  proxy:
    external: true # Assumes the 'proxy' network is external and managed by M3TAL
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically configure routing for `app.DOMAIN` to your `my-app` service.

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle and view its status using standard `systemctl` commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **Restart the API daemon:**
    ```bash
    systemctl restart m3tal-api
    ```
*   **View real-time logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```

## Environment Variables

The following environment variables are utilized by the M3TAL system, primarily defined in `/etc/m3tal/.env`:

| Key                     | Default                   | Description                                                               |
| :---------------------- | :------------------------ | :------------------------------------------------------------------------ |
| `DASHBOARD_PORT`        | `8082`                    | Internal port for the M3TAL Dashboard container.                          |
| `DASHBOARD_EXPOSE_MODE` | `local`                   | `local` for direct port access, `traefik` for Traefik-routed access.      |
| `HTTP_PORT`             | `8080`                    | Port the M3TAL API daemon listens on.                                     |
| `STATE_DIR`             | `./state`                 | Directory for the SQLite state database (`/var/lib/m3tal/state.db`).      |
| `LOG_LEVEL`             | `info`                    | Logging level for the M3TAL API daemon.                                   |
| `DASHBOARD_SECRET`      | `change_me_immediately`   | Secret key for the M3TAL Dashboard, required for session management.      |
| `API_TOKEN`             | `change_me_api_token`     | Authentication token for accessing the M3TAL API.                         |
| `ADMIN_PASSWORD`        | `admin_pass`              | Default password for the dashboard admin user.                            |
| `NETWORK_NAME`          | `m3tal`                   | Default Docker network name for M3TAL services.                           |
| `LOCAL_IP`              | `127.0.0.1`               | Local IP address for internal routing.                                    |
| `DOMAIN`                | `localhost`               | Base domain for Traefik routing (e.g., `api.DOMAIN`, `dash.DOMAIN`).      |
| `VPN_USER`              | `user`                    | Username for VPN configurations (if applicable).                          |
| `VPN_PASSWORD`          | `password`                | Password for VPN configurations (if applicable).                          |
| `BASE_STORAGE_PATH`     | `./data`                  | Base path for M3TAL storage volumes.                                      |
| `MEDIA_PATH`            | `./data/media`            | Path for media files.                                                     |
| `CONFIG_PATH`           | `./data/config`           | Path for configuration files.                                             |
| `DOWNLOADS_PATH`        | `./data/downloads`        | Path for downloaded files.                                                |
| `PUID`                  | `1000`                    | User ID for container processes to ensure correct file permissions.       |
| `PGID`                  | `1000`                    | Group ID for container processes to ensure correct file permissions.      |
| `TZ`                    | `America/Denver`          | Timezone setting for containers.                                          |
| `TRAEFIK_WEB_PORT`      | `80`                      | Traefik's public HTTP entry point.                                        |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                     | Traefik's public HTTPS entry point (if configured).                       |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                    | Internal port for the Traefik dashboard. (Exposed at host-local 8081).    |
| `DEBUG_MODE`            | `false`                   | Enables debug logging and features.                                       |
| `METRICS_ENABLED`       | `true`                    | Enables metrics collection.                                               |