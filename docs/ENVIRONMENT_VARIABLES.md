# M3TAL Environment Variables Reference

All M3TAL configuration is managed through environment variables, primarily defined in `/etc/m3tal/.env`. This file is read by both the M3TAL CLI and all Docker Compose stacks via `--env-file`. The `m3tal config wizard` or `m3tal config set` commands are the recommended ways to manage these variables.

## Quick Reference Table

| Variable Name        | Default Value     | Description                                                                                                  | Components Using It                                  |
|----------------------|-------------------|--------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| **Core**             |                   |                                                                                                              |                                                      |
| `HTTP_PORT`          | `8080`            | The port the M3TAL API daemon listens on.                                                                    | `m3tal-api.service`                                  |
| `STATE_DIR`          | `./state`         | The directory where M3TAL stores its state database and other essential files.                               | `m3tal-api.service`, `m3tal-dashboard`               |
| `LOG_LEVEL`          | `info`            | The logging level for M3TAL services (e.g., `debug`, `info`, `warn`, `error`).                               | `m3tal-api.service`, `m3tal-dashboard`               |
| `DEBUG_MODE`         | `false`           | Enables debug mode for M3TAL services.                                                                       | `m3tal-api.service`, `m3tal-dashboard`               |
| `METRICS_ENABLED`    | `true`            | Enables or disables the collection and exposure of M3TAL metrics.                                            | `m3tal-api.service`                                  |
| **Auth**             |                   |                                                                                                              |                                                      |
| `DASHBOARD_SECRET`   | `change_me_immediately` | A secret key used for securing the dashboard session. **Auto-generated on first `m3tal init`.**               | `m3tal-dashboard`                                    |
| `API_TOKEN`          | `change_me_api_token`   | An API token for authenticating with the M3TAL API. **Auto-generated on first `m3tal init`.**                | `m3tal-api.service`                                  |
| `ADMIN_PASSWORD`     | `admin_pass`      | The password for the default admin user of the M3TAL dashboard.                                              | `m3tal-dashboard`                                    |
| **Network**          |                   |                                                                                                              |                                                      |
| `NETWORK_NAME`       | `m3tal`           | The name of the Docker network M3TAL services will use.                                                      | `m3tal-api.service`, `m3tal-dashboard`, `traefik`    |
| `LOCAL_IP`           | `127.0.0.1`       | The local IP address to bind services to.                                                                    | `m3tal-api.service`                                  |
| **Storage**          |                   |                                                                                                              |                                                      |
| `BASE_STORAGE_PATH`  | `./data`          | The base directory for M3TAL data storage. **Defaults to `/mnt` in production deployments.**                   | `m3tal-dashboard`, `m3tal-api.service`               |
| `MEDIA_PATH`         | `./data/media`    | The directory within `BASE_STORAGE_PATH` for media files.                                                    | `m3tal-dashboard`                                    |
| `CONFIG_PATH`        | `./data/config`   | The directory within `BASE_STORAGE_PATH` for configuration files.                                            | `m3tal-dashboard`                                    |
| `DOWNLOADS_PATH`     | `./data/downloads`| The directory within `BASE_STORAGE_PATH` for downloads.                                                      | `m3tal-dashboard`                                    |
| `PUID`               | `1000`            | The user ID to run containers with.                                                                          | `m3tal-dashboard`                                    |
| `PGID`               | `1000`            | The group ID to run containers with.                                                                         | `m3tal-dashboard`                                    |
| **Traefik**          |                   |                                                                                                              |                                                      |
| `DASHBOARD_PORT`     | `8082`            | The internal port the M3TAL dashboard container listens on.                                                  | `m3tal-dashboard`, `traefik`                         |
| `DASHBOARD_EXPOSE_MODE` | `local`       | Controls how the dashboard is exposed. `local` (default) uses direct port binding, `traefik` uses Traefik routing. | `m3tal-dashboard`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml` |
| `DOMAIN`             | `localhost`       | The domain name used for Traefik routing rules (e.g., `api.${DOMAIN}`, `dash.${DOMAIN}`).                  | `traefik`, `m3tal-compose.traefik.yml`               |
| `TRAEFIK_WEB_PORT`   | `80`              | The host port Traefik uses for HTTP traffic.                                                                 | `traefik`                                            |
| `TRAEFIK_WEBHTTPS_PORT` | `443`          | The host port Traefik uses for HTTPS traffic.                                                                | `traefik`                                            |
| `TRAEFIK_DASHBOARD_PORT` | `8080`        | The port Traefik exposes its own dashboard on (internal).                                                    | `traefik`                                            |
| **VPN**              |                   |                                                                                                              |                                                      |
| `VPN_USER`           | `user`            | The username for VPN connection.                                                                             | Not explicitly used by documented components.        |
| `VPN_PASSWORD`       | `password`        | The password for VPN connection.                                                                             | Not explicitly used by documented components.        |
| **System**           |                   |                                                                                                              |                                                      |
| `TZ`                 | `America/Denver`  | The timezone to be used by M3TAL services.                                                                   | `m3tal-dashboard`                                    |

---

## Detailed Environment Variable Descriptions

All M3TAL configuration is managed through environment variables, primarily defined in `/etc/m3tal/.env`. This file is read by both the M3TAL CLI and all Docker Compose stacks via `--env-file`. The `m3tal config wizard` or `m3tal config set` commands are the recommended ways to manage these variables.

### Core

These variables control fundamental aspects of M3TAL's operation.

*   **`HTTP_PORT`**
    *   **Description:** The port the M3TAL API daemon listens on.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Components Using It:** `m3tal-api.service`

*   **`STATE_DIR`**
    *   **Description:** The directory where M3TAL stores its state database and other essential files.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Components Using It:** `m3tal-api.service`, `m3tal-dashboard`

*   **`LOG_LEVEL`**
    *   **Description:** The logging level for M3TAL services. Accepted values include `debug`, `info`, `warn`, and `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Components Using It:** `m3tal-api.service`, `m3tal-dashboard`

*   **`DEBUG_MODE`**
    *   **Description:** Enables debug mode for M3TAL services, providing more verbose logging and potentially other debugging features.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Components Using It:** `m3tal-api.service`, `m3tal-dashboard`

*   **`METRICS_ENABLED`**
    *   **Description:** Enables or disables the collection and exposure of M3TAL metrics, which can be useful for monitoring.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Components Using It:** `m3tal-api.service`

### Auth

These variables are crucial for securing access to your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for securing the dashboard session. **This variable is auto-generated on first `m3tal init`. Users should NOT set it manually unless performing a secret rotation.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `s3cr3tK3yF0rD4shB04rd`
    *   **Components Using It:** `m3tal-dashboard`

*   **`API_TOKEN`**
    *   **Description:** An API token for authenticating with the M3TAL API. **This variable is auto-generated on first `m3tal init`. Users should NOT set it manually unless performing a secret rotation.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `a_v3ry_s3cure_t0k3n`
    *   **Components Using It:** `m3tal-api.service`

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default admin user of the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `myS3cureP@sswOrd!`
    *   **Components Using It:** `m3tal-dashboard`

### Network

These variables define network-related settings for M3TAL services.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network M3TAL services will use for inter-service communication.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal_network`
    *   **Components Using It:** `m3tal-api.service`, `m3tal-dashboard`, `traefik`

*   **`LOCAL_IP`**
    *   **Description:** The local IP address to bind services to. This is typically `127.0.0.1` for local access.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `127.0.0.1`
    *   **Components Using It:** `m3tal-api.service`

### Storage

These variables configure where M3TAL stores its data and media.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The base directory for M3TAL data storage. This controls where media data is stored. **Defaults to `/mnt` in production deployments, not `./data` as seen in template configurations.**
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt` (for production)
    *   **Components Using It:** `m3tal-dashboard`, `m3tal-api.service`

*   **`MEDIA_PATH`**
    *   **Description:** The directory within `BASE_STORAGE_PATH` specifically designated for storing media files.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Components Using It:** `m3tal-dashboard`

*   **`CONFIG_PATH`**
    *   **Description:** The directory within `BASE_STORAGE_PATH` for configuration files, including user-specific settings.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Components Using It:** `m3tal-dashboard`

*   **`DOWNLOADS_PATH`**
    *   **Description:** The directory within `BASE_STORAGE_PATH` used for downloaded files or temporary storage.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Components Using It:** `m3tal-dashboard`

*   **`PUID`**
    *   **Description:** The User ID (UID) to assign to processes running within M3TAL containers. This is important for file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Components Using It:** `m3tal-dashboard`

*   **`PGID`**
    *   **Description:** The Group ID (GID) to assign to processes running within M3TAL containers. This is important for file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Components Using It:** `m3tal-dashboard`

### Traefik

These variables are specifically for configuring the Traefik reverse proxy and how M3TAL services are exposed.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port the M3TAL dashboard container listens on. This is distinct from the port users access it through.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Components Using It:** `m3tal-dashboard`, `traefik`

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL dashboard is exposed to users.
        *   `local` (default): Uses a direct port binding (`${DASHBOARD_PORT}:8082`), accessible via `http://HOST_IP:${DASHBOARD_PORT}`. No Traefik is required for this mode.
        *   `traefik`: Configures Traefik to route traffic for `dash.${DOMAIN}` to the dashboard container. Traefik must be running and accessible for this mode to work.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Components Using It:** `m3tal-dashboard`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`

*   **`DOMAIN`**
    *   **Description:** The base domain name used for Traefik routing rules. Setting this variable enables routes like `api.${DOMAIN}` and `dash.${DOMAIN}`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `yourdomain.com`
    *   **Components Using It:** `traefik`, `m3tal-compose.traefik.yml`

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik will use for incoming HTTP traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Components Using It:** `traefik`

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik will use for incoming HTTPS traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Components Using It:** `traefik`

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which Traefik exposes its own management dashboard. This is typically accessed via `http://127.0.0.1:8081`.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Components Using It:** `traefik`

### VPN

These variables are for configuring VPN client settings, though their direct usage within the documented core M3TAL components is not explicitly detailed in the provided context.

*   **`VPN_USER`**
    *   **Description:** The username for establishing a VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Components Using It:** Not explicitly used by documented components.

*   **`VPN_PASSWORD`**
    *   **Description:** The password for establishing a VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `supersecretpassword`
    *   **Components Using It:** Not explicitly used by documented components.

### System

These variables configure system-level settings for M3TAL services.

*   **`TZ`**
    *   **Description:** The timezone to be used by M3TAL services. This ensures accurate logging and time-based operations.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Components Using It:** `m3tal-dashboard`