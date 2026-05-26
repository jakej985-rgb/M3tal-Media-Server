# Environment Variables Reference

All M3TAL environment variables are read from the primary configuration file located at `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (started via `m3tal up` or `m3tal dash up`) utilize this file by passing it via the `--env-file` argument, ensuring a consistent configuration across the entire M3TAL ecosystem.

You can manage these variables using the `m3tal config wizard` for an interactive setup or `m3tal config set KEY value` for specific changes.

---

## Quick Reference

| Variable Name           | Default Value         | Description                                                          |
| :---------------------- | :-------------------- | :------------------------------------------------------------------- |
| `HTTP_PORT`             | `8080`                | Port for the M3TAL API daemon.                                       |
| `STATE_DIR`             | `./state`             | Directory for M3TAL's state database (`state.db`).                   |
| `LOG_LEVEL`             | `info`                | Verbosity of M3TAL API daemon logging.                               |
| `DEBUG_MODE`            | `false`               | Enables debug features and verbose logging.                          |
| `METRICS_ENABLED`       | `true`                | Enables Prometheus-compatible metrics for the API.                   |
| `NETWORK_NAME`          | `m3tal`               | Name of the Docker network for M3TAL components.                     |
| `PUID`                  | `1000`                | User ID (UID) for container processes.                               |
| `PGID`                  | `1000`                | Group ID (GID) for container processes.                              |
| `TZ`                    | `America/Denver`      | Timezone for containers.                                             |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret key for Dashboard session management. **(Auto-generated)**      |
| `API_TOKEN`             | `change_me_api_token` | Bearer token for API authentication. **(Auto-generated)**            |
| `ADMIN_PASSWORD`        | `admin_pass`          | Initial password for the default Dashboard admin user.               |
| `LOCAL_IP`              | `127.0.0.1`           | Local IP address of the host machine.                                |
| `DOMAIN`                | `localhost`           | Primary domain for Traefik routing.                                  |
| `BASE_STORAGE_PATH`     | `./data`              | Base directory for M3TAL-related media and app data.                 |
| `MEDIA_PATH`            | `./data/media`        | Subdirectory for user media files.                                   |
| `CONFIG_PATH`           | `./data/config`       | Subdirectory for application configuration files.                    |
| `DOWNLOADS_PATH`        | `./data/downloads`    | Subdirectory for downloaded content.                                 |
| `DASHBOARD_EXPOSE_MODE` | `local`               | Controls how the M3TAL Dashboard is exposed (`local` or `traefik`). |
| `DASHBOARD_PORT`        | `8082`                | Host port for direct Dashboard access in `local` mode.               |
| `TRAEFIK_WEB_PORT`      | `80`                  | Host port Traefik binds for HTTP traffic.                            |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                 | Host port Traefik binds for HTTPS traffic.                           |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                | Host port Traefik's internal dashboard is accessible on.             |
| `VPN_USER`              | `user`                | Username for VPN authentication (user stacks).                       |
| `VPN_PASSWORD`          | `password`            | Password for VPN authentication (user stacks).                       |

---

## Environment Variables Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL API daemon and general container runtime.

#### `HTTP_PORT`
*   **Description:** The port on which the M3TAL API daemon listens for incoming HTTP connections.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Component:** `m3tal-api.service`

#### `STATE_DIR`
*   **Description:** The host directory where the M3TAL API daemon stores its SQLite state database (`state.db`) and other runtime files. This path is often mounted into containers, for example, the `m3tal-dashboard` mounts it to `/docker/state`.
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal`
*   **Component:** `m3tal-api.service`, `m3tal-dashboard` (via volume mount)

#### `LOG_LEVEL`
*   **Description:** Controls the verbosity of logging output for the M3TAL API daemon. Valid options typically include `debug`, `info`, `warn`, `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Component:** `m3tal-api.service`

#### `DEBUG_MODE`
*   **Description:** When set to `true`, enables additional debug logging and potentially debug-specific features across M3TAL components.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Component:** `m3tal-api.service`, `m3tal-dashboard`

#### `METRICS_ENABLED`
*   **Description:** When set to `true`, the M3TAL API daemon exposes Prometheus-compatible metrics on the API port, typically at `/metrics`.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Component:** `m3tal-api.service`

#### `NETWORK_NAME`
*   **Description:** The name of the Docker network that M3TAL core components (`m3tal-dashboard`, `traefik`) and any user-defined stacks will join for inter-container communication.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal_prod`
*   **Component:** Docker Compose stacks (e.g., `m3tal-dashboard`, `traefik`, `cloudflared`), user-defined stacks

#### `PUID`
*   **Description:** The User ID (UID) that container processes will use when interacting with the host filesystem. This ensures correct file ownership and permissions for mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Component:** `m3tal-dashboard`, user-defined containers

#### `PGID`
*   **Description:** The Group ID (GID) that container processes will use when interacting with the host filesystem. This ensures correct file ownership and permissions for mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Component:** `m3tal-dashboard`, user-defined containers

#### `TZ`
*   **Description:** Specifies the timezone for containers. This ensures logs and scheduled tasks reflect the correct local time.
*   **Default Value:** `America/Denver`
*   **Example Value:** `Europe/London`
*   **Component:** `m3tal-dashboard`, user-defined containers

### Authentication & Security

These variables manage authentication tokens and secrets crucial for the security of your M3TAL instance.

#### `DASHBOARD_SECRET`
*   **Description:** A secret key used by the M3TAL Dashboard for secure session management, CSRF protection, and other cryptographic operations.
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a_long_random_string_of_characters`
*   **Component:** `m3tal-dashboard`
*   **Note:** This value is **auto-generated** on the first `m3tal init` and stored in `/etc/m3tal/.env`. Users should **NOT** set it manually unless performing a secret rotation.

#### `API_TOKEN`
*   **Description:** A bearer token used for authenticating requests to the M3TAL API daemon. The dashboard uses this token to communicate with the API.
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `another_long_random_api_key_string`
*   **Component:** `m3tal-api.service`, `m3tal-dashboard`
*   **Note:** This value is **auto-generated** on the first `m3tal init` and stored in `/etc/m3tal/.env`. Users should **NOT** set it manually unless performing a token rotation.

#### `ADMIN_PASSWORD`
*   **Description:** The initial password set for the default admin user of the M3TAL Dashboard. This is used to populate `/docker/users.json` on first setup. You can manage dashboard users and passwords with `m3tal dashpass`.
*   **Default Value:** `admin_pass`
*   **Example Value:** `my_secure_admin_password`
*   **Component:** `m3tal-dashboard` (initial user setup)

### Network & Routing

These variables define networking aspects and domain routing for your M3TAL setup.

#### `LOCAL_IP`
*   **Description:** The local IP address of the host machine. Used by components for host-local communication or for services that need to bind to a specific interface.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Component:** `m3tal-api.service` (binding), potentially user-defined stacks

#### `DOMAIN`
*   **Description:** The primary domain name for your M3TAL services.
*   **Default Value:** `localhost`
*   **Example Value:** `myhomelab.com`
*   **Component:** `traefik`, `m3tal-dashboard` (in Traefik mode), `cloudflared`
*   **Note:** **Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik**, allowing access to the dashboard and API via domain names.

### Storage Paths

These variables define the base and subdirectories for persistent storage on the host machine.

#### `BASE_STORAGE_PATH`
*   **Description:** The root directory on the host where all M3TAL-related media and application data are stored. Other path variables are typically subdirectories of this.
*   **Default Value:** `./data`
*   **Example Value:** `/mnt/m3tal_data`
*   **Component:** `m3tal-dashboard`, user-defined stacks
*   **Note:** In production deployments, this defaults to `/mnt` to leverage dedicated storage mounts, rather than `./data` (which is common for development or local testing).

#### `MEDIA_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` designated for user media files (e.g., videos, music, photos).
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/m3tal_data/media`
*   **Component:** User-defined media-serving stacks

#### `CONFIG_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` for application configuration files (e.g., `users.json` for the dashboard).
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/m3tal_data/config`
*   **Component:** `m3tal-dashboard` (for `users.json`), user-defined stacks

#### `DOWNLOADS_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` for downloaded content (e.g., by download managers, torrent clients).
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/m3tal_data/downloads`
*   **Component:** User-defined download managers

### Dashboard Exposure

These variables specifically control how the M3TAL Dashboard container is made accessible.

#### `DASHBOARD_EXPOSE_MODE`
*   **Description:** Determines the method of exposing the M3TAL Dashboard.
    *   `local`: The dashboard is directly exposed via a host port binding (`DASHBOARD_PORT`). Best for LAN-only or local testing.
    *   `traefik`: The dashboard is exposed via Traefik as `dash.DOMAIN`. Requires Traefik to be running. Best for domain-based setups.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Component:** `m3tal dash up` command, `m3tal-dashboard` compose configuration (uses `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` override).

#### `DASHBOARD_PORT`
*   **Description:** The host port that the M3TAL Dashboard container will bind to when `DASHBOARD_EXPOSE_MODE` is set to `local`.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Component:** `m3tal-dashboard` (in `local` expose mode)

### Traefik Gateway

These variables configure the Traefik reverse proxy for handling incoming web traffic.

#### `TRAEFIK_WEB_PORT`
*   **Description:** The host port Traefik binds to listen for incoming HTTP traffic (entrypoint `web`).
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Component:** `traefik`

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description:** The host port Traefik binds to listen for incoming HTTPS traffic (entrypoint `websecure`, though not explicitly configured in ground truth, this is a standard Traefik expectation).
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Component:** `traefik`

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description:** The host port on which Traefik's own internal dashboard (not the M3TAL Dashboard) is accessible.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Component:** `traefik`
*   **Note:** While the variable defaults to `8080` as per the JSON, the `traefik` service typically maps `127.0.0.1:8081:8080`, meaning the host port is `8081` by default to avoid conflicts. This variable would control the *host* port if configured in Traefik's compose file.

### VPN Services

These variables are placeholders for configuring user-defined VPN service containers.

#### `VPN_USER`
*   **Description:** Username for authentication within a user-defined VPN stack (e.g., for WireGuard, OpenVPN containers).
*   **Default Value:** `user`
*   **Example Value:** `vpnadmin`
*   **Component:** User-defined VPN stack (e.g., `cloudflared` if it uses credentials, or other VPN containers)

#### `VPN_PASSWORD`
*   **Description:** Password for authentication within a user-defined VPN stack.
*   **Default Value:** `password`
*   **Example Value:** `strong_vpn_pass`
*   **Component:** User-defined VPN stack (e.g., `cloudflared` if it uses credentials, or other VPN containers)