# M3TAL Environment Variables Reference

All M3TAL environment variables are stored in `/etc/m3tal/.env` and are read by both the M3TAL CLI and all Docker Compose stacks. The `m3tal config wizard` and `m3tal config set` commands are used to manage these variables.

## Quick Reference

| Variable Name           | Description                                                                                                                                                                                                                                                                | Default Value       |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| **Core**                |                                                                                                                                                                                                                                                                            |                     |
| `DASHBOARD_PORT`        | The port on which the M3TAL Dashboard service will listen.                                                                                                                                                                                                                 | `8082`              |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed. `local` exposes it directly via a port mapping (accessible via `HOST_IP:DASHBOARD_PORT`), while `traefik` exposes it via Traefik routing rules (accessible via `dash.DOMAIN`).                                                       | `local`             |
| `HTTP_PORT`             | The port on which the M3TAL API daemon will listen.                                                                                                                                                                                                                        | `8080`              |
| `STATE_DIR`             | The directory where M3TAL stores its state, including the SQLite database. This path is relative to the `CONFIG_PATH` if it's a relative path, or absolute otherwise.                                                                                                   | `./state`           |
| `LOG_LEVEL`             | Sets the logging level for M3TAL services.                                                                                                                                                                                                                                 | `info`              |
| `DEBUG_MODE`            | Enables or disables debug mode for M3TAL services.                                                                                                                                                                                                                         | `false`             |
| `METRICS_ENABLED`       | Enables or disables the collection and exposure of metrics for M3TAL services.                                                                                                                                                                                             | `true`              |
| **Auth**                |                                                                                                                                                                                                                                                                            |                     |
| `DASHBOARD_SECRET`      | A secret key used for session management within the M3TAL Dashboard. **Auto-generated on first `m3tal init`. Should NOT be set manually unless rotating.**                                                                                                                  | `change_me_immediately` |
| `API_TOKEN`             | An API token used for authentication with the M3TAL API. **Auto-generated on first `m3tal init`. Should NOT be set manually unless rotating.**                                                                                                                            | `change_me_api_token` |
| `ADMIN_PASSWORD`        | The password for the default administrator account on the M3TAL Dashboard.                                                                                                                                                                                                 | `admin_pass`        |
| **Network**             |                                                                                                                                                                                                                                                                            |                     |
| `NETWORK_NAME`          | The name of the Docker network used by M3TAL services.                                                                                                                                                                                                                     | `m3tal`             |
| `LOCAL_IP`              | The IP address that M3TAL services will bind to locally.                                                                                                                                                                                                                   | `127.0.0.1`         |
| `DOMAIN`                | The primary domain name for your M3TAL instance. Setting this enables Traefik routing rules for `dash.DOMAIN` and `api.DOMAIN`.                                                                                                                                             | `localhost`         |
| **Storage**             |                                                                                                                                                                                                                                                                            |                     |
| `BASE_STORAGE_PATH`     | The base directory for storing all M3TAL-related data, including configuration, media, and downloads. **Defaults to `/mnt` in production deployments.**                                                                                                                     | `./data`            |
| `MEDIA_PATH`            | The path within `BASE_STORAGE_PATH` where media files are stored.                                                                                                                                                                                                          | `./data/media`      |
| `CONFIG_PATH`           | The path within `BASE_STORAGE_PATH` where configuration files are stored.                                                                                                                                                                                                  | `./data/config`     |
| `DOWNLOADS_PATH`        | The path within `BASE_STORAGE_PATH` where downloaded files are stored.                                                                                                                                                                                                     | `./data/downloads`  |
| **System**              |                                                                                                                                                                                                                                                                            |                     |
| `PUID`                  | The User ID (UID) that M3TAL Docker containers will run as.                                                                                                                                                                                                                | `1000`              |
| `PGID`                  | The Group ID (GID) that M3TAL Docker containers will run as.                                                                                                                                                                                                               | `1000`              |
| `TZ`                    | The timezone to be used by M3TAL services.                                                                                                                                                                                                                                 | `America/Denver`   |
| **Traefik**             |                                                                                                                                                                                                                                                                            |                     |
| `TRAEFIK_WEB_PORT`      | The host port that Traefik will use as its HTTP entrypoint.                                                                                                                                                                                                                | `80`                |
| `TRAEFIK_WEBHTTPS_PORT` | The host port that Traefik will use as its HTTPS entrypoint.                                                                                                                                                                                                               | `443`               |
| `TRAEFIK_DASHBOARD_PORT`| The port on which Traefik's own dashboard will be accessible (typically only via `127.0.0.1`).                                                                                                                                                                             | `8080`              |
| **VPN**                 |                                                                                                                                                                                                                                                                            |                     |
| `VPN_USER`              | The username for VPN authentication.                                                                                                                                                                                                                                       | `user`              |
| `VPN_PASSWORD`          | The password for VPN authentication.                                                                                                                                                                                                                                       | `password`          |

---

## Detailed Explanation

All M3TAL environment variables are managed through the `/etc/m3tal/.env` file. This file serves as the central configuration point for the entire M3TAL ecosystem. Both the M3TAL CLI (`/usr/bin/m3tal`) and all Docker Compose stacks (managed via `m3tal up`) read their configurations from this file. It is strongly recommended to use the `m3tal config wizard` or `m3tal config set` commands to modify these variables to ensure consistency and prevent misconfigurations.

### Core

These variables control fundamental aspects of the M3TAL system.

*   **`DASHBOARD_PORT`**
    *   **Description:** Specifies the internal port on which the M3TAL Dashboard service listens.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container, Traefik (when `DASHBOARD_EXPOSE_MODE=traefik`).

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Determines how the M3TAL Dashboard is made accessible.
        *   `local`: The dashboard is directly exposed via a port mapping (`${DASHBOARD_PORT}:8082`). Access is typically via `http://HOST_IP:DASHBOARD_PORT` or `http://localhost:DASHBOARD_PORT`. This mode does not require Traefik and is suitable for LAN-only access or initial setup.
        *   `traefik`: The dashboard is exposed through Traefik routing rules, accessible via `http://dash.${DOMAIN}`. This mode requires Traefik to be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-dashboard` container (via compose overrides), Traefik.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon (the Go binary) listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api.service`, Traefik (for routing to the API).

*   **`STATE_DIR`**
    *   **Description:** The directory where M3TAL stores its persistent state. This includes the SQLite database file (`state.db`). This path can be relative or absolute. If relative, it's resolved relative to the `CONFIG_PATH`.
    *   **Default Value:** `./state`
    *   **Example Value:** `/mnt/config/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`LOG_LEVEL`**
    *   **Description:** Controls the verbosity of logging output for M3TAL services. Common values include `debug`, `info`, `warn`, and `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug-specific features and logging within M3TAL services.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether M3TAL services expose Prometheus-compatible metrics for monitoring.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api.service`.

### Auth

These variables are crucial for securing access to your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A strong, randomly generated secret used for signing session cookies in the M3TAL Dashboard. This helps prevent session hijacking. **This variable is auto-generated by the `m3tal init` command. Manual modification is generally not required unless you need to rotate the secret for security reasons.**
    *   **Default Value:** `change_me_immediately` (This indicates it's a placeholder and MUST be changed).
    *   **Example Value:** `s0m3_r4nd0m_s3cr3t_k3y_f0r_s3ss10ns`
    *   **Used By:** `m3tal-dashboard` container.

*   **`API_TOKEN`**
    *   **Description:** A secret token used for authenticating programmatic access to the M3TAL API. This token should be kept confidential. **This variable is auto-generated by the `m3tal init` command. Manual modification is generally not required unless you need to rotate the token for security reasons.**
    *   **Default Value:** `change_me_api_token` (This indicates it's a placeholder and MUST be changed).
    *   **Example Value:** `42f8a1b9c0d7e6f5a3b1c0d7e6f5a3b1c0d7e6f5a3b1c0d7e6f5a3b1`
    *   **Used By:** M3TAL CLI, `m3tal-api.service`.

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator account on the M3TAL Dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_secure_admin_password_123`
    *   **Used By:** `m3tal-dashboard` container.

### Network

These variables define network-related settings for M3TAL.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will be connected to. This is often the same network Traefik is also connected to, allowing services to communicate with each other and with Traefik.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-net`
    *   **Used By:** M3TAL CLI, Docker Compose stacks.

*   **`LOCAL_IP`**
    *   **Description:** The IP address that M3TAL services will bind to on the host machine. This is typically `127.0.0.1` for local access or a specific IP address if you need to bind to a particular interface.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** `m3tal-api.service` (internal binding), `m3tal-dashboard` container (for `host.docker.internal` resolution).

*   **`DOMAIN`**
    *   **Description:** The primary domain name that M3TAL will use for routing. Setting this variable is crucial for enabling Traefik-based access to services. When set, Traefik will be configured to route requests to `api.${DOMAIN}` and `dash.${DOMAIN}` (if `DASHBOARD_EXPOSE_MODE=traefik`).
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.mydomain.com`
    *   **Used By:** Traefik (via compose labels), M3TAL CLI (for generating URLs).

### Storage

These variables configure the locations where M3TAL stores its data.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory on the host machine where M3TAL will store all its persistent data, including configuration files, media, and downloads. **In production deployments, this defaults to `/mnt` to ensure data is stored on a persistent volume outside of the Docker overlay filesystem.** In development or template environments, it might default to `./data`.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt/m3tal_data`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container, M3TAL CLI.

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files (e.g., uploaded images, videos) are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/m3tal_data/media`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files, including the state database.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/m3tal_data/config`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files from M3TAL services are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/m3tal_data/downloads`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### System

These variables relate to system-level configurations, including user permissions and timezone.

*   **`PUID`**
    *   **Description:** The User ID (UID) that the Docker containers will run as. This is important for ensuring that the containers have the correct file permissions on the host system, especially when volumes are mounted.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container (via compose file user directive).

*   **`PGID`**
    *   **Description:** The Group ID (GID) that the Docker containers will run as. Similar to `PUID`, this ensures correct file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container (via compose file user directive).

*   **`TZ`**
    *   **Description:** Specifies the timezone to be used by M3TAL services. This ensures that timestamps in logs and other time-sensitive data are accurate.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `Europe/London`
    *   **Used By:** `m3tal-dashboard` container.

### Traefik

These variables configure Traefik, the reverse proxy used by M3TAL.

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The port on the host machine that Traefik will listen on for incoming HTTP traffic. This is the primary entry point for services exposed via Traefik.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik container.

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The port on the host machine that Traefik will listen on for incoming HTTPS traffic. This is typically used when SSL/TLS is configured.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik container.

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The port on which Traefik's administrative dashboard is exposed internally. This is usually only accessible from `localhost` and is not intended for external access.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik container.

### VPN

These variables are used for configuring VPN connections if your M3TAL setup requires it.

*   **`VPN_USER`**
    *   **Description:** The username for authenticating with a VPN service.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** M3TAL CLI (potentially), specific VPN client containers.

*   **`VPN_PASSWORD`**
    *   **Description:** The password for authenticating with a VPN service.
    *   **Default Value:** `password`
    *   **Example Value:** `my_super_secret_vpn_password`
    *   **Used By:** M3TAL CLI (potentially), specific VPN client containers.