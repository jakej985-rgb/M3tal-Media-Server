# Environment Variables Reference

Welcome, M3TAL operator. This document outlines the environment variables that configure your M3TAL ecosystem.

All variables detailed below are read from the primary configuration file, `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (started via `m3tal up` or `m3tal dash up`) load their configuration from this file using the `--env-file` option.

It is highly recommended to manage these variables using the `m3tal config wizard` command for guided setup, or `m3tal config set KEY value` for direct modifications. After any changes to `/etc/m3tal/.env`, it's generally necessary to restart affected services (e.g., `systemctl restart m3tal-api` for the API daemon, or `m3tal up -d` for Docker containers).

---

## Quick Reference

| Name                      | Description                                                                                                                                                                                                     | Default                | Component(s)                     |
| :------------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------- | :------------------------------- |
| `DASHBOARD_PORT`          | The network port for the M3TAL Dashboard.                                                                                                                                                                       | `8082`                 | `m3tal-dashboard`, `CLI binary`  |
| `DASHBOARD_EXPOSE_MODE`   | Controls how the M3TAL Dashboard is exposed: `local` (direct port) or `traefik` (via `dash.DOMAIN`).                                                                                                            | `local`                | `CLI binary`                     |
| `HTTP_PORT`               | The port on which the M3TAL API daemon listens.                                                                                                                                                                 | `8080`                 | `API daemon`                     |
| `STATE_DIR`               | Internal base directory for state files within certain containers (e.g., Dashboard).                                                                                                                            | `./state`              | `m3tal-dashboard`                |
| `LOG_LEVEL`               | Sets the verbosity of log output.                                                                                                                                                                               | `info`                 | `API daemon`, `CLI binary`, `m3tal-dashboard` |
| `DASHBOARD_SECRET`        | Secret key for M3TAL Dashboard session data. **Auto-generated on `m3tal init`.**                                                                                                                                | `change_me_immediately`| `m3tal-dashboard`                |
| `API_TOKEN`               | Authentication token for M3TAL API access. **Auto-generated on `m3tal init`.**                                                                                                                                  | `change_me_api_token`  | `API daemon`, `CLI binary`       |
| `ADMIN_PASSWORD`          | Initial password for the `admin` user in the M3TAL Dashboard.                                                                                                                                                   | `admin_pass`           | `m3tal-dashboard`                |
| `NETWORK_NAME`            | Name of the shared Docker network for M3TAL components and user stacks.                                                                                                                                         | `m3tal`                | `Docker Compose`                 |
| `LOCAL_IP`                | The host machine's local IP address.                                                                                                                                                                            | `127.0.0.1`            | `CLI binary`, `Traefik gateway`, `API daemon` |
| `DOMAIN`                  | Primary domain name for M3TAL services. Enables `dash.DOMAIN` and `api.DOMAIN` routes.                                                                                                                          | `localhost`            | `Traefik gateway`, `CLI binary`  |
| `VPN_USER`                | Username for VPN/Cloudflare Tunnel authentication.                                                                                                                                                              | `user`                 | `cloudflared`                    |
| `VPN_PASSWORD`            | Password for VPN/Cloudflare Tunnel authentication.                                                                                                                                                              | `password`             | `cloudflared`                    |
| `BASE_STORAGE_PATH`       | Root directory for all M3TAL persistent data on the host. **Defaults to `/mnt` in production.**                                                                                                               | `./data`               | `m3tal-dashboard`, `User stacks` |
| `MEDIA_PATH`              | Sub-directory within `BASE_STORAGE_PATH` for media files.                                                                                                                                                       | `./data/media`         | `User stacks`                    |
| `CONFIG_PATH`             | Sub-directory within `BASE_STORAGE_PATH` for application configurations.                                                                                                                                        | `./data/config`        | `m3tal-dashboard`, `User stacks` |
| `DOWNLOADS_PATH`          | Sub-directory within `BASE_STORAGE_PATH` for downloaded files.                                                                                                                                                  | `./data/downloads`     | `User stacks`                    |
| `PUID`                    | User ID (UID) for Docker containers to ensure correct file permissions.                                                                                                                                         | `1000`                 | `m3tal-dashboard`, `User stacks` |
| `PGID`                    | Group ID (GID) for Docker containers to ensure correct file permissions.                                                                                                                                        | `1000`                 | `m3tal-dashboard`, `User stacks` |
| `TZ`                      | Sets the timezone for M3TAL services and containers.                                                                                                                                                            | `America/Denver`       | `m3tal-dashboard`, `API daemon`, `User stacks` |
| `TRAEFIK_WEB_PORT`        | The port Traefik listens on for HTTP traffic.                                                                                                                                                                   | `80`                   | `Traefik gateway`                |
| `TRAEFIK_WEBHTTPS_PORT`   | The port Traefik listens on for HTTPS traffic.                                                                                                                                                                  | `443`                  | `Traefik gateway`                |
| `TRAEFIK_DASHBOARD_PORT`  | The internal container port for Traefik's management dashboard.                                                                                                                                                 | `8080`                 | `Traefik gateway`                |
| `DEBUG_MODE`              | Enables verbose debugging output and potentially development features.                                                                                                                                          | `false`                | `API daemon`, `CLI binary`       |
| `METRICS_ENABLED`         | Controls whether the M3TAL API daemon exposes Prometheus-compatible metrics.                                                                                                                                    | `true`                 | `API daemon`                     |

---

## Detailed Reference

### Core Variables

Variables fundamental to the M3TAL system's operation.

#### `HTTP_PORT`
*   **Description:** The port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming HTTP connections. This port is internal to the host and typically accessed via Traefik or directly by other M3TAL components.
*   **Default:** `8080`
*   **Example:** `8080`
*   **Component(s):** `API daemon`

#### `STATE_DIR`
*   **Description:** Specifies the internal base directory for state files within certain M3TAL containers. For the M3TAL Dashboard, this is set to `/docker/state` inside the container, which is volume-mounted from `${CONFIG_PATH:-/mnt/config}/m3tal/state` on the host. The M3TAL API daemon's primary state database is fixed at `/var/lib/m3tal/state.db`.
*   **Default:** `./state`
*   **Example:** `/docker/state` (within container)
*   **Component(s):** `m3tal-dashboard` (internal path)

#### `LOG_LEVEL`
*   **Description:** Sets the verbosity of log output for M3TAL components. Available options typically include `debug`, `info`, `warn`, `error`, `fatal`. `debug` provides the most detailed output.
*   **Default:** `info`
*   **Example:** `debug`
*   **Component(s):** `API daemon`, `CLI binary`, `m3tal-dashboard`

#### `DEBUG_MODE`
*   **Description:** When set to `true`, enables verbose debugging output and potentially additional development-focused features or logging in the API daemon and CLI.
*   **Default:** `false`
*   **Example:** `true`
*   **Component(s):** `API daemon`, `CLI binary`

#### `METRICS_ENABLED`
*   **Description:** Controls whether the M3TAL API daemon exposes Prometheus-compatible metrics on its `/metrics` endpoint. Set to `false` to disable metrics collection.
*   **Default:** `true`
*   **Example:** `false`
*   **Component(s):** `API daemon`

### Authentication Variables

Variables related to securing access to M3TAL services.

#### `DASHBOARD_SECRET`
*   **Description:** A unique, cryptographically strong secret key used by the M3TAL Dashboard container to secure session data (e.g., user logins).
    **This variable is auto-generated on the first `m3tal init` and should NOT be set manually by users unless rotating the secret for security reasons.**
*   **Default:** `change_me_immediately`
*   **Example:** `a_long_random_string_of_characters_for_session_security`
*   **Component(s):** `m3tal-dashboard`

#### `API_TOKEN`
*   **Description:** An authentication token required for programmatic access to the M3TAL API daemon. Used by the CLI and other integrated services.
    **This variable is auto-generated on the first `m3tal init` and should NOT be set manually by users unless rotating the token for security reasons.**
*   **Default:** `change_me_api_token`
*   **Example:** `another_long_random_string_of_characters_for_api_auth`
*   **Component(s):** `API daemon`, `CLI binary`

#### `ADMIN_PASSWORD`
*   **Description:** The initial password for the default `admin` user account in the M3TAL Dashboard. This password is used to populate `/docker/users.json` during initial setup. It is strongly recommended to change this immediately after initial setup using the `m3tal dashpass` command.
*   **Default:** `admin_pass`
*   **Example:** `my_strong_admin_password_123!`
*   **Component(s):** `m3tal-dashboard`

### Network Variables

Variables that configure networking for M3TAL services.

#### `DASHBOARD_PORT`
*   **Description:** The network port that the M3TAL Dashboard container exposes.
    *   In `local` exposure mode, this port is directly mapped to the host (e.g., `http://HOST_IP:8082`).
    *   In `traefik` exposure mode, this port is used internally by Traefik to route traffic to the dashboard.
*   **Default:** `8082`
*   **Example:** `8082`
*   **Component(s):** `m3tal-dashboard`, `CLI binary`

#### `DASHBOARD_EXPOSE_MODE`
*   **Description:** Controls how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard is directly exposed on `DASHBOARD_PORT` on the host. Best for LAN-only setups or initial testing.
    *   `traefik`: The dashboard is exposed via Traefik under the hostname `dash.${DOMAIN}`. Requires Traefik to be running.
*   **Default:** `local`
*   **Example:** `traefik`
*   **Component(s):** `CLI binary` (`m3tal dash up` command reads this to select the appropriate Docker Compose override file.)

#### `NETWORK_NAME`
*   **Description:** The name of the custom Docker network that M3TAL components and any user-defined Docker Compose stacks will connect to. This ensures inter-container communication across different `docker compose` deployments.
*   **Default:** `m3tal`
*   **Example:** `m3tal_proxy`
*   **Component(s):** `Docker Compose` (all stacks)

#### `LOCAL_IP`
*   **Description:** The local IP address of the host machine. Used by certain components for internal network references, such as resolving `host.docker.internal` for the API daemon or when configuring specific network routes.
*   **Default:** `127.0.0.1`
*   **Example:** `192.168.1.100`
*   **Component(s):** `CLI binary`, `Traefik gateway` (for host-gateway resolution), `API daemon`

#### `DOMAIN`
*   **Description:** The primary domain name associated with your M3TAL deployment.
    **Setting this variable is crucial if you intend to use Traefik for routing, as it enables domain-based routes such as `dash.YOUR_DOMAIN` for the Dashboard and `api.YOUR_DOMAIN` for the API.**
*   **Default:** `localhost`
*   **Example:** `example.com`
*   **Component(s):** `Traefik gateway`, `CLI binary`

### Storage Variables

Variables defining the host filesystem paths for persistent data.

#### `BASE_STORAGE_PATH`
*   **Description:** The root directory on the host filesystem where all M3TAL-related persistent data is stored. This includes media, configurations, downloads, and other application data.
    **In production M3TAL deployments, this variable defaults to `/mnt`, not `./data` as seen in development templates.**
*   **Default:** `./data`
*   **Example:** `/mnt/m3tal_data`
*   **Component(s):** `m3tal-dashboard`, `User stacks`

#### `MEDIA_PATH`
*   **Description:** The specific sub-directory, typically within `BASE_STORAGE_PATH`, where media files (videos, music, photos) are expected to be stored. This path is commonly mounted into user-defined media server containers (e.g., Plex, Jellyfin).
*   **Default:** `./data/media`
*   **Example:** `/mnt/m3tal_data/media`
*   **Component(s):** `User stacks` (e.g., media servers)

#### `CONFIG_PATH`
*   **Description:** The specific sub-directory, typically within `BASE_STORAGE_PATH`, where application configuration files for M3TAL components and user stacks are stored. This path is used to persist the M3TAL Dashboard's internal state (e.g., `users.json`).
*   **Default:** `./data/config`
*   **Example:** `/mnt/m3tal_data/config`
*   **Component(s):** `m3tal-dashboard`, `User stacks`

#### `DOWNLOADS_PATH`
*   **Description:** The specific sub-directory, typically within `BASE_STORAGE_PATH`, where downloaded files are stored. This path is commonly mounted into download client containers.
*   **Default:** `./data/downloads`
*   **Example:** `/mnt/m3tal_data/downloads`
*   **Component(s):** `User stacks` (e.g., download clients)

### Traefik Variables

Variables configuring the Traefik reverse proxy gateway.

#### `TRAEFIK_WEB_PORT`
*   **Description:** The port on the host machine that Traefik listens on for incoming standard HTTP traffic.
*   **Default:** `80`
*   **Example:** `80`
*   **Component(s):** `Traefik gateway`

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description:** The port on the host machine that Traefik listens on for incoming secure HTTPS traffic.
*   **Default:** `443`
*   **Example:** `443`
*   **Component(s):** `Traefik gateway`

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description:** The internal port on which the Traefik management dashboard service listens within its container. This port is typically mapped to `127.0.0.1:8081` on the host for local, private access (as configured in `routing-compose.yml`).
*   **Default:** `8080`
*   **Example:** `8080` (internal container port)
*   **Component(s):** `Traefik gateway`

### VPN Variables

Variables used for VPN or secure tunnel configurations, primarily with Cloudflared.

#### `VPN_USER`
*   **Description:** The username required for authentication with a configured VPN service or Cloudflare Tunnel, if such a service is used for remote access or secure networking.
*   **Default:** `user`
*   **Example:** `mycloudflaretunneluser`
*   **Component(s):** `cloudflared`

#### `VPN_PASSWORD`
*   **Description:** The password required for authentication with a configured VPN service or Cloudflare Tunnel.
*   **Default:** `password`
*   **Example:** `my_tunnel_secret_password`
*   **Component(s):** `cloudflared`

### System Variables

General system-level configuration variables.

#### `PUID`
*   **Description:** The User ID (UID) that Docker containers should run processes as. This is crucial for ensuring proper file permissions when containers interact with host volumes (e.g., `BASE_STORAGE_PATH`, `MEDIA_PATH`). It should match the UID of a user on your host system with appropriate permissions.
*   **Default:** `1000`
*   **Example:** `998`
*   **Component(s):** `m3tal-dashboard`, `User stacks`

#### `PGID`
*   **Description:** The Group ID (GID) that Docker containers should run processes as. Similar to `PUID`, this ensures correct group permissions on host volumes. It should match the GID of a group on your host system.
*   **Default:** `1000`
*   **Example:** `998`
*   **Component(s):** `m3tal-dashboard`, `User stacks`

#### `TZ`
*   **Description:** Sets the timezone for M3TAL services and containers. This is important for accurate logging timestamps and time-based operations. Uses standard TZ database names (e.g., `America/New_York`, `Europe/London`).
*   **Default:** `America/Denver`
*   **Example:** `America/Los_Angeles`
*   **Component(s):** `m3tal-dashboard`, `API daemon`, `User stacks`