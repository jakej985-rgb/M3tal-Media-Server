# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, primarily read from `/etc/m3tal/.env`. This file is automatically managed by the `m3tal config wizard` and can be updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks utilize this `.env` file via `--env-file`.

**Quick Reference Table:**

| Variable Name             | Description                                                                                                                                                               | Default Value        | Example Value                               | Used By                  |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- | ------------------------------------------- | ------------------------ |
| **Core**                  |                                                                                                                                                                           |                      |                                             |                          |
| `DASHBOARD_PORT`          | The internal port the M3TAL dashboard listens on.                                                                                                                         | `8082`               | `8082`                                      | Dashboard Container      |
| `HTTP_PORT`               | The port the M3TAL API daemon listens on.                                                                                                                                 | `8080`               | `5050`                                      | API Daemon               |
| `STATE_DIR`               | Directory for storing M3TAL state.                                                                                                                                        | `./state`            | `/mnt/config/m3tal/state`                   | API Daemon, Dashboard    |
| `LOG_LEVEL`               | Sets the logging verbosity for M3TAL services.                                                                                                                            | `info`               | `debug`                                     | API Daemon, Dashboard    |
| `DEBUG_MODE`              | Enables or disables debug mode.                                                                                                                                           | `false`              | `true`                                      | M3TAL CLI, API Daemon    |
| `METRICS_ENABLED`         | Enables or disables metrics collection.                                                                                                                                   | `true`               | `false`                                     | API Daemon               |
| `PUID`                    | The User ID to run Docker containers as.                                                                                                                                  | `1000`               | `1000`                                      | All Docker Containers    |
| `PGID`                    | The Group ID to run Docker containers as.                                                                                                                                 | `1000`               | `1000`                                      | All Docker Containers    |
| `TZ`                      | Timezone setting for containers.                                                                                                                                          | `America/Denver`     | `UTC`                                       | All Docker Containers    |
| **Auth**                  |                                                                                                                                                                           |                      |                                             |                          |
| `DASHBOARD_SECRET`        | **Auto-generated on first `m3tal init`**. Used for securing dashboard sessions. Rotate only if necessary.                                                                  | `change_me_immediately` | `your_secure_auto_generated_secret`       | Dashboard Container      |
| `API_TOKEN`               | **Auto-generated on first `m3tal init`**. Used for authenticating API requests. Rotate only if necessary.                                                                 | `change_me_api_token`  | `your_secure_auto_generated_api_token`    | API Daemon               |
| `ADMIN_PASSWORD`          | The password for the default admin user of the dashboard.                                                                                                                 | `admin_pass`         | `a_strong_password`                         | Dashboard Container      |
| **Network**               |                                                                                                                                                                           |                      |                                             |                          |
| `NETWORK_NAME`            | The name of the Docker network M3TAL services will use.                                                                                                                   | `m3tal`              | `m3tal_network`                             | Docker Compose           |
| `LOCAL_IP`                | The local IP address of the host machine.                                                                                                                                 | `127.0.0.1`          | `192.168.1.100`                             | API Daemon               |
| `DOMAIN`                  | **Controls Traefik routing rules**. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes. Defaults to `localhost` for local development.                           | `localhost`          | `mydomain.com`                              | Traefik, API Daemon      |
| `TRAEFIK_WEB_PORT`        | The host port Traefik listens on for HTTP traffic.                                                                                                                        | `80`                 | `80`                                        | Traefik Container        |
| `TRAEFIK_WEBHTTPS_PORT`   | The host port Traefik listens on for HTTPS traffic.                                                                                                                       | `443`                | `443`                                       | Traefik Container        |
| `TRAEFIK_DASHBOARD_PORT`  | The host port Traefik listens on for its own dashboard.                                                                                                                   | `8080`               | `8081`                                      | Traefik Container        |
| **Storage**               |                                                                                                                                                                           |                      |                                             |                          |
| `BASE_STORAGE_PATH`       | **Controls where media data is stored**. Defaults to `/mnt` in production deployments.                                                                                      | `./data`             | `/mnt`                                      | All Docker Containers    |
| `MEDIA_PATH`              | Path within `BASE_STORAGE_PATH` for media files.                                                                                                                          | `./data/media`       | `/mnt/media`                                | All Docker Containers    |
| `CONFIG_PATH`             | Path within `BASE_STORAGE_PATH` for configuration files.                                                                                                                  | `./data/config`      | `/mnt/config`                               | All Docker Containers    |
| `DOWNLOADS_PATH`          | Path within `BASE_STORAGE_PATH` for downloads.                                                                                                                            | `./data/downloads`   | `/mnt/downloads`                            | All Docker Containers    |
| **Traefik**               |                                                                                                                                                                           |                      |                                             |                          |
| `DASHBOARD_EXPOSE_MODE`   | Controls how the M3TAL dashboard is exposed: `local` (direct port binding) or `traefik` (via Traefik routing).                                                            | `local`              | `traefik`                                   | M3TAL CLI, Dashboard     |
| **VPN**                   |                                                                                                                                                                           |                      |                                             |                          |
| `VPN_USER`                | Username for VPN connection.                                                                                                                                              | `user`               | `myvpnuser`                                 | VPN Container            |
| `VPN_PASSWORD`            | Password for VPN connection.                                                                                                                                              | `password`           | `mysecretpassword`                          | VPN Container            |

---

## Core

### `DASHBOARD_PORT`

*   **Description:** The internal port that the M3TAL dashboard service listens on within its Docker container.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** `m3tal-dashboard` container.

### `HTTP_PORT`

*   **Description:** The port that the M3TAL API daemon (Go binary) listens on for incoming HTTP requests.
*   **Default Value:** `8080`
*   **Example Value:** `5050`
*   **Used By:** `m3tal-api.service`.

### `STATE_DIR`

*   **Description:** Specifies the directory where M3TAL stores its state database and other critical configuration files.
*   **Default Value:** `./state`
*   **Example Value:** `/mnt/config/m3tal/state`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `LOG_LEVEL`

*   **Description:** Controls the verbosity of logging for M3TAL services. Accepted values typically include `debug`, `info`, `warn`, and `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `DEBUG_MODE`

*   **Description:** A boolean flag to enable or disable debug mode for M3TAL. This can provide more verbose output and enable additional debugging features.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** M3TAL CLI, `m3tal-api.service`.

### `METRICS_ENABLED`

*   **Description:** A boolean flag to enable or disable the collection and exposure of Prometheus metrics from the M3TAL API.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** `m3tal-api.service`.

### `PUID`

*   **Description:** The User ID (UID) that Docker containers will run as. This is important for file permissions within mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** All Docker Containers.

### `PGID`

*   **Description:** The Group ID (GID) that Docker containers will run as. This is important for file permissions within mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** All Docker Containers.

### `TZ`

*   **Description:** Specifies the timezone to be used by all Docker containers. This ensures consistent time logging and operations.
*   **Default Value:** `America/Denver`
*   **Example Value:** `UTC`
*   **Used By:** All Docker Containers.

## Auth

### `DASHBOARD_SECRET`

*   **Description:** **This variable is auto-generated on the first `m3tal init` command.** It is used as a secret key for securing session cookies and other sensitive operations within the M3TAL dashboard. Users should **not** set this manually unless performing a secret rotation.
*   **Default Value:** `change_me_immediately` (This is a placeholder; it will be auto-generated with a strong random value on first initialization).
*   **Example Value:** `your_secure_auto_generated_secret_string_here`
*   **Used By:** `m3tal-dashboard` container.

### `API_TOKEN`

*   **Description:** **This variable is auto-generated on the first `m3tal init` command.** It serves as an API token for authenticating requests to the M3TAL API daemon. Users should **not** set this manually unless performing a token rotation.
*   **Default Value:** `change_me_api_token` (This is a placeholder; it will be auto-generated with a strong random value on first initialization).
*   **Example Value:** `your_secure_auto_generated_api_token_string_here`
*   **Used By:** `m3tal-api.service`.

### `ADMIN_PASSWORD`

*   **Description:** The password for the default administrative user of the M3TAL dashboard.
*   **Default Value:** `admin_pass`
*   **Example Value:** `a_strong_password`
*   **Used By:** `m3tal-dashboard` container.

## Network

### `NETWORK_NAME`

*   **Description:** The name of the Docker network that M3TAL services will be connected to.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal_network`
*   **Used By:** Docker Compose for defining network configurations.

### `LOCAL_IP`

*   **Description:** The IP address of the host machine that M3TAL services can bind to or be accessed on.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** `m3tal-api.service` (for internal communication with `host.docker.internal`).

### `DOMAIN`

*   **Description:** **This variable is critical for Traefik routing.** When set to a valid domain name (e.g., `mydomain.com`), M3TAL will configure Traefik to route requests to `api.DOMAIN` and `dash.DOMAIN`. If left as `localhost`, Traefik rules will reflect that, suitable for local development.
*   **Default Value:** `localhost`
*   **Example Value:** `mydomain.com`
*   **Used By:** Traefik container (for routing rules), `m3tal-api.service` (potentially for API host configuration).

### `TRAEFIK_WEB_PORT`

*   **Description:** The host port that Traefik will expose for incoming HTTP (port 80) traffic.
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** `traefik` container.

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The host port that Traefik will expose for incoming HTTPS (port 443) traffic.
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** `traefik` container.

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The host port that Traefik will expose for accessing its own administrative dashboard. This is typically bound to `127.0.0.1` for local access only.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** `traefik` container.

## Storage

### `BASE_STORAGE_PATH`

*   **Description:** **This variable controls the base directory where all persistent M3TAL data, including media, configuration, and downloads, is stored.** In production deployments, this defaults to `/mnt` to ensure data is stored on a potentially persistent volume. In template/development environments, it might default to `./data`.
*   **Default Value:** `./data`
*   **Example Value:** `/mnt`
*   **Used By:** All Docker Containers for volume mounts.

### `MEDIA_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` where media files (e.g., uploaded images, videos) are stored.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/media`
*   **Used By:** All Docker Containers for media volume mounts.

### `CONFIG_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` where M3TAL configuration files are stored. This often includes state database files.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/config`
*   **Used By:** All Docker Containers for configuration volume mounts.

### `DOWNLOADS_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` where downloaded files or artifacts are stored.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/downloads`
*   **Used By:** All Docker Containers for downloads volume mounts.

## Traefik

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Determines how the M3TAL dashboard is made accessible.
    *   `local`: The dashboard is exposed directly via a host port mapping (controlled by `DASHBOARD_PORT`). Access is typically via `http://<HOST_IP>:<DASHBOARD_PORT>`. This mode does not require Traefik to be running.
    *   `traefik`: The dashboard is exposed through the Traefik reverse proxy. Access is via `http://dash.<DOMAIN>`. This mode requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** M3TAL CLI (during `m3tal dash up`), `m3tal-dashboard` container (for Traefik labels when in `traefik` mode).

## VPN

### `VPN_USER`

*   **Description:** The username used to authenticate with a VPN service.
*   **Default Value:** `user`
*   **Example Value:** `myvpnuser`
*   **Used By:** Cloudflared container (if configured for VPN access).

### `VPN_PASSWORD`

*   **Description:** The password used to authenticate with a VPN service.
*   **Default Value:** `password`
*   **Example Value:** `mysecretpassword`
*   **Used By:** Cloudflared container (if configured for VPN access).