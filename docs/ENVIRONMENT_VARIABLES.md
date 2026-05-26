# M3TAL Environment Variables Reference

All M3TAL configuration is managed through environment variables, primarily read from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks via the `--env-file` directive. This file is automatically managed by the `m3tal config wizard` and `m3tal config set` commands.

## Quick Reference Table

| Variable Name           | Description                                                                                                     | Default Value       |
|-------------------------|-----------------------------------------------------------------------------------------------------------------|---------------------|
| **Core**                |                                                                                                                 |                     |
| `DASHBOARD_PORT`        | The internal port the M3TAL dashboard runs on.                                                                  | `8082`              |
| `HTTP_PORT`             | The internal port the M3TAL API daemon listens on.                                                              | `8080`              |
| `STATE_DIR`             | The directory where M3TAL stores its state, including the SQLite database.                                    | `./state`           |
| `LOG_LEVEL`             | The logging level for M3TAL services (e.g., debug, info, warn, error).                                          | `info`              |
| `DEBUG_MODE`            | Enables debug mode for M3TAL services.                                                                          | `false`             |
| `METRICS_ENABLED`       | Enables Prometheus metrics collection for M3TAL services.                                                       | `true`              |
| **Auth**                |                                                                                                                 |                     |
| `DASHBOARD_SECRET`      | A secret key for securing the M3TAL dashboard session cookies. Auto-generated on first `m3tal init`.           | `change_me_immediately` |
| `API_TOKEN`             | A token for authenticating API requests. Auto-generated on first `m3tal init`.                                  | `change_me_api_token` |
| `ADMIN_PASSWORD`        | The password for the default `admin` user of the M3TAL dashboard.                                               | `admin_pass`        |
| **Network**             |                                                                                                                 |                     |
| `NETWORK_NAME`          | The name of the Docker network M3TAL services will use.                                                         | `m3tal`             |
| `LOCAL_IP`              | The IP address M3TAL services should bind to for local communication.                                           | `127.0.0.1`         |
| `DOMAIN`                | The domain name used for Traefik routing. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes.           | `localhost`         |
| `TRAEFIK_WEB_PORT`      | The host port Traefik listens on for HTTP traffic.                                                              | `80`                |
| `TRAEFIK_WEBHTTPS_PORT` | The host port Traefik listens on for HTTPS traffic.                                                             | `443`               |
| `TRAEFIK_DASHBOARD_PORT`| The internal port Traefik's own dashboard is exposed on (usually for debugging).                                | `8080`              |
| **Storage**             |                                                                                                                 |                     |
| `BASE_STORAGE_PATH`     | The base directory for all M3TAL data storage. Defaults to `/mnt` in production.                                | `./data`            |
| `MEDIA_PATH`            | Path within `BASE_STORAGE_PATH` for media files.                                                                | `./data/media`      |
| `CONFIG_PATH`           | Path within `BASE_STORAGE_PATH` for configuration files (e.g., Docker compose overrides).                     | `./data/config`     |
| `DOWNLOADS_PATH`        | Path within `BASE_STORAGE_PATH` for downloaded files.                                                           | `./data/downloads`  |
| **System**              |                                                                                                                 |                     |
| `PUID`                  | The User ID to run containers with.                                                                             | `1000`              |
| `PGID`                  | The Group ID to run containers with.                                                                            | `1000`              |
| `TZ`                    | The timezone to use for M3TAL services (e.g., `America/Denver`).                                              | `America/Denver`    |
| **VPN**                 |                                                                                                                 |                     |
| `VPN_USER`              | Username for the VPN connection.                                                                                | `user`              |
| `VPN_PASSWORD`          | Password for the VPN connection.                                                                                | `password`          |

---

## Detailed Environment Variable Reference

All environment variables are read from `/etc/m3tal/.env` by both the CLI and all Docker Compose stacks.

### Core

These variables control the fundamental operation of M3TAL services.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port the M3TAL dashboard container listens on. This is used for direct access in `local` expose mode and for Traefik routing in `traefik` mode.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`

*   **`HTTP_PORT`**
    *   **Description:** The internal port the M3TAL API daemon (Go binary) listens on. This is the primary API endpoint for M3TAL.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api.service`

*   **`STATE_DIR`**
    *   **Description:** The directory where M3TAL stores its persistent state, including the SQLite database (`state.db`).
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`LOG_LEVEL`**
    *   **Description:** Sets the logging verbosity for M3TAL services. Accepted values typically include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`DEBUG_MODE`**
    *   **Description:** A boolean flag to enable or disable debug mode for M3TAL services, which may provide more verbose logging or enable additional debugging features.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`METRICS_ENABLED`**
    *   **Description:** A boolean flag to enable or disable the collection and exposure of Prometheus metrics for M3TAL services.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api.service`

### Auth

These variables are critical for securing your M3TAL instance. `DASHBOARD_SECRET` and `API_TOKEN` are auto-generated on first `m3tal init` and should generally not be modified unless you intend to rotate them.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used to encrypt the session cookies for the M3TAL dashboard. **Auto-generated on first `m3tal init`**. Users should NOT set this manually unless rotating the secret.
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_long_and_random_secure_secret_key`
    *   **Used By:** `m3tal-dashboard` container

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating requests to the M3TAL API. **Auto-generated on first `m3tal init`**. Users should NOT set this manually unless rotating the token.
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `a_very_secure_api_authentication_token`
    *   **Used By:** M3TAL CLI, `m3tal-api.service` (to validate incoming requests)

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default `admin` user account in the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_strong_admin_password`
    *   **Used By:** `m3tal-dashboard` container

### Network

These variables control M3TAL's network configuration and how it interacts with other services, especially through Traefik.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will join. This ensures proper inter-container communication.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-network`
    *   **Used By:** All M3TAL Docker Compose stacks.

*   **`LOCAL_IP`**
    *   **Description:** The IP address that M3TAL services will bind to for local network communication. Typically `127.0.0.1` or `host.docker.internal`.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `host.docker.internal`
    *   **Used By:** `m3tal-api.service` (for dashboard communication), `m3tal-dashboard` container (for API communication)

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL instance. Setting this variable enables Traefik to create routing rules for `dash.DOMAIN` and `api.DOMAIN`. If `localhost` is used, Traefik will route to `dash.localhost` and `api.localhost`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `my.m3tal.io`
    *   **Used By:** `m3tal-dashboard` container (via Traefik labels), Traefik configuration (`dynamic/api.yml`, `m3tal-compose.traefik.yml`)

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port on which Traefik will listen for incoming HTTP (port 80) traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** `traefik` container

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port on which Traefik will listen for incoming HTTPS (port 443) traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** `traefik` container

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which Traefik's own administrative dashboard is exposed within the Docker network. This is usually accessed via `http://127.0.0.1:8081`.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `traefik` container

### Storage

These variables define the location and structure of data storage for M3TAL.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory where M3TAL stores all its data, including configuration, media, and downloads. **Defaults to `/mnt` in production deployments, not `./data` as in the template.**
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Used By:** `m3tal-dashboard` container, M3TAL CLI

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files (e.g., images, videos) are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Used By:** `m3tal-dashboard` container

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` used for storing M3TAL-related configuration files, potentially including Docker Compose override files or other application settings.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Used By:** `m3tal-dashboard` container

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Used By:** `m3tal-dashboard` container

### System

These variables configure system-level aspects of the M3TAL deployment, such as user permissions and timezones.

*   **`PUID`**
    *   **Description:** The User ID (UID) that Docker containers should run as. This is important for ensuring correct file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container, and other services that inherit PUID/PGID.

*   **`PGID`**
    *   **Description:** The Group ID (GID) that Docker containers should run as. This is important for ensuring correct file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container, and other services that inherit PUID/PGID.

*   **`TZ`**
    *   **Description:** The timezone to be used by M3TAL services for accurate timestamping of logs and events. This should be a valid IANA timezone name.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`, `Europe/London`
    *   **Used By:** `m3tal-dashboard` container

### VPN

These variables are for configuring an optional VPN client within the M3TAL environment.

*   **`VPN_USER`**
    *   **Description:** The username for authenticating with the VPN service.
    *   **Default Value:** `user`
    *   **Example Value:** `my_vpn_username`
    *   **Used By:** VPN client container (if configured)

*   **`VPN_PASSWORD`**
    *   **Description:** The password for authenticating with the VPN service.
    *   **Default Value:** `password`
    *   **Example Value:** `my_strong_vpn_password`
    *   **Used By:** VPN client container (if configured)