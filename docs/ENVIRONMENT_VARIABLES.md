# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables are crucial for configuring and managing your M3TAL deployment.

All environment variables are read from the `/etc/m3tal/.env` file by both the M3TAL CLI (`/usr/bin/m3tal`) and all Docker Compose stacks via the `--env-file` option. It is highly recommended to manage these variables using the `m3tal config wizard` or `m3tal config set KEY value` commands.

---

## Quick Reference Table

| Category   | Variable Name            | Default Value       | Description                                                                                                                                                                                                                                                              |
|------------|--------------------------|---------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Core**   | `DASHBOARD_PORT`         | `8082`              | The port on which the M3TAL Dashboard service will run.                                                                                                                                                                                                  |
|            | `HTTP_PORT`              | `8080`              | The port on which the M3TAL API daemon (Go) will listen.                                                                                                                                                                                                 |
|            | `STATE_DIR`              | `./state`           | The directory where M3TAL stores its state, including the SQLite database. This path is relative to the `CONFIG_PATH` if `CONFIG_PATH` is set, otherwise it's relative to the current working directory of the API daemon.                                            |
|            | `LOG_LEVEL`              | `info`              | The logging level for M3TAL services (e.g., `debug`, `info`, `warn`, `error`).                                                                                                                                                                            |
|            | `DEBUG_MODE`             | `false`             | Enables or disables debug mode for M3TAL services.                                                                                                                                                                                                       |
|            | `METRICS_ENABLED`        | `true`              | Enables or disables the collection and reporting of M3TAL metrics.                                                                                                                                                                                         |
|            | `PUID`                   | `1000`              | The user ID to run Docker containers under. This helps in managing file permissions.                                                                                                                                                                       |
|            | `PGID`                   | `1000`              | The group ID to run Docker containers under. This helps in managing file permissions.                                                                                                                                                                      |
|            | `TZ`                     | `America/Denver`    | The timezone for M3TAL services. This ensures correct timestamping in logs and other time-sensitive operations.                                                                                                                                           |
| **Auth**   | `DASHBOARD_SECRET`       | `change_me_immediately` | A secret key used for securing the M3TAL Dashboard session. **Auto-generated on first `m3tal init`. Do not set manually unless rotating.**                                                                                                                   |
|            | `API_TOKEN`              | `change_me_api_token`   | A token used for authenticating API requests. **Auto-generated on first `m3tal init`. Do not set manually unless rotating.**                                                                                                                            |
|            | `ADMIN_PASSWORD`         | `admin_pass`        | The password for the default administrative user of the M3TAL Dashboard.                                                                                                                                                                                 |
| **Network**| `NETWORK_NAME`           | `m3tal`             | The name of the Docker network that M3TAL services will join.                                                                                                                                                                                            |
|            | `LOCAL_IP`               | `127.0.0.1`         | The local IP address that M3TAL services will bind to.                                                                                                                                                                                                   |
|            | `DOMAIN`                 | `localhost`         | The base domain name for M3TAL services. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.                                                                                                                                          |
| **Storage**| `BASE_STORAGE_PATH`      | `./data`            | **Controls where media data is stored.** Defaults to `/mnt` in production deployments. In template/development environments, it defaults to `./data`.                                                                                                       |
|            | `MEDIA_PATH`             | `./data/media`      | The subdirectory within `BASE_STORAGE_PATH` for storing media files.                                                                                                                                                                                     |
|            | `CONFIG_PATH`            | `./data/config`     | The subdirectory within `BASE_STORAGE_PATH` for storing configuration files.                                                                                                                                                                             |
|            | `DOWNLOADS_PATH`         | `./data/downloads`  | The subdirectory within `BASE_STORAGE_PATH` for storing downloaded files.                                                                                                                                                                                |
| **Traefik**| `TRAEFIK_WEB_PORT`       | `80`                | The host port that Traefik will use for HTTP traffic.                                                                                                                                                                                                    |
|            | `TRAEFIK_WEBHTTPS_PORT`  | `443`               | The host port that Traefik will use for HTTPS traffic.                                                                                                                                                                                                   |
|            | `TRAEFIK_DASHBOARD_PORT` | `8080`              | The internal port Traefik listens on for its own dashboard. Accessible at `127.0.0.1:TRAEFIK_DASHBOARD_PORT`.                                                                                                                                              |
| **VPN**    | `VPN_USER`               | `user`              | The username for the VPN connection.                                                                                                                                                                                                                     |
|            | `VPN_PASSWORD`           | `password`          | The password for the VPN connection.                                                                                                                                                                                                                     |
| **System** | `DASHBOARD_EXPOSE_MODE`  | `local`             | Controls how the M3TAL Dashboard is exposed. Options are `local` (direct port access) and `traefik` (via Traefik routing).                                                                                                                             |

---

## Environment Variable Details

All environment variables are read from `/etc/m3tal/.env`.

### Core

These variables control fundamental aspects of the M3TAL system's operation.

*   **`DASHBOARD_PORT`**
    *   **Description:** The port on which the M3TAL Dashboard service will run.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon (Go) will listen.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api.service`.

*   **`STATE_DIR`**
    *   **Description:** The directory where M3TAL stores its state, including the SQLite database. This path is relative to the `CONFIG_PATH` if `CONFIG_PATH` is set, otherwise it's relative to the current working directory of the API daemon.
    *   **Default Value:** `./state`
    *   **Example Value:** `./state` or `/mnt/config/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`LOG_LEVEL`**
    *   **Description:** The logging level for M3TAL services (e.g., `debug`, `info`, `warn`, `error`).
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode for M3TAL services.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** M3TAL CLI, `m3tal-api.service`, `m3tal-dashboard` container.

*   **`METRICS_ENABLED`**
    *   **Description:** Enables or disables the collection and reporting of M3TAL metrics.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** M3TAL CLI, `m3tal-api.service`, `m3tal-dashboard` container.

*   **`PUID`**
    *   **Description:** The user ID to run Docker containers under. This helps in managing file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container (in its compose file).

*   **`PGID`**
    *   **Description:** The group ID to run Docker containers under. This helps in managing file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container (in its compose file).

*   **`TZ`**
    *   **Description:** The timezone for M3TAL services. This ensures correct timestamping in logs and other time-sensitive operations.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC` or `Europe/London`
    *   **Used By:** `m3tal-dashboard` container.

### Auth

These variables are related to authentication and security within the M3TAL ecosystem.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for securing the M3TAL Dashboard session. **Auto-generated on first `m3tal init`. Do not set manually unless rotating.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** (Auto-generated)
    *   **Used By:** `m3tal-dashboard` container.

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating API requests. **Auto-generated on first `m3tal init`. Do not set manually unless rotating.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** (Auto-generated)
    *   **Used By:** M3TAL CLI, `m3tal-api.service`.

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrative user of the M3TAL Dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_secure_password`
    *   **Used By:** `m3tal-dashboard` container (via `/docker/users.json` which is populated from this variable on first run).

### Network

These variables configure network-related settings for M3TAL services.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will join.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal`
    *   **Used By:** M3TAL CLI, Docker Compose files.

*   **`LOCAL_IP`**
    *   **Description:** The local IP address that M3TAL services will bind to.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `127.0.0.1`
    *   **Used By:** M3TAL CLI, `m3tal-api.service`.

*   **`DOMAIN`**
    *   **Description:** The base domain name for M3TAL services. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.
    *   **Default Value:** `localhost`
    *   **Example Value:** `yourdomain.com`
    *   **Used By:** M3TAL CLI, `m3tal-dashboard` container (via Traefik labels), Traefik configuration.

### Storage

These variables define the locations for M3TAL's data storage.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** **Controls where media data is stored.** Defaults to `/mnt` in production deployments. In template/development environments, it defaults to `./data`.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt` (production) or `./data` (development)
    *   **Used By:** `m3tal-dashboard` container (volume mounts), M3TAL CLI.

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` for storing media files.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `./data/media` or `/mnt/media`
    *   **Used By:** M3TAL CLI.

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` for storing configuration files.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `./data/config` or `/mnt/config`
    *   **Used By:** `m3tal-dashboard` container (volume mounts), M3TAL CLI.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` for storing downloaded files.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `./data/downloads` or `/mnt/downloads`
    *   **Used By:** M3TAL CLI.

### Traefik

These variables are specific to the Traefik reverse proxy configuration.

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik will use for HTTP traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** `traefik` container (port mapping).

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik will use for HTTPS traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** `traefik` container (port mapping).

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port Traefik listens on for its own dashboard. Accessible at `127.0.0.1:TRAEFIK_DASHBOARD_PORT`.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `traefik` container (port mapping).

### VPN

These variables are used if you intend to configure a VPN connection for M3TAL.

*   **`VPN_USER`**
    *   **Description:** The username for the VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** M3TAL CLI (potentially for VPN configuration within services).

*   **`VPN_PASSWORD`**
    *   **Description:** The password for the VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `mysecretvpnpassword`
    *   **Used By:** M3TAL CLI (potentially for VPN configuration within services).

### System

These variables control system-level behaviors and deployment modes.

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL Dashboard is exposed. Options are `local` (direct port access) and `traefik` (via Traefik routing).
        *   **`local`**: Uses override `m3tal-compose.local.yml`. Adds a direct port binding `${DASHBOARD_PORT:-8082}:8082`. Access via `http://HOST_IP:8082` or `http://localhost:8082`. No Traefik required.
        *   **`traefik`**: Uses override `m3tal-compose.traefik.yml`. Adds Traefik labels so Traefik routes `dash.${DOMAIN}` → dashboard on port 8082. Access via `http://dash.DOMAIN` (Traefik must be running via `m3tal up`).
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** M3TAL CLI (when managing the dashboard service), `m3tal-dashboard` container (via compose overrides).