# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables control the configuration, behavior, and network access of your M3TAL installation.

All M3TAL environment variables are read from the `/etc/m3tal/.env` file. This file is managed by the `m3tal config wizard` and can be updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks utilize this `.env` file via the `--env-file` option.

## Quick Reference Table

| Variable Name            | Description                                                                                                            | Default Value     | Example Value                                   | Used By                                     |
|--------------------------|------------------------------------------------------------------------------------------------------------------------|-------------------|-------------------------------------------------|---------------------------------------------|
| **Core**                 |                                                                                                                        |                   |                                                 |                                             |
| `DASHBOARD_PORT`         | The internal port the M3TAL dashboard service listens on.                                                              | `8082`            | `8082`                                          | Dashboard container                         |
| `DASHBOARD_EXPOSE_MODE`  | Controls how the dashboard is exposed: `local` (direct port) or `traefik` (via Traefik).                               | `local`           | `traefik`                                       | Dashboard container, Traefik                |
| `HTTP_PORT`              | The port the M3TAL API daemon listens on.                                                                              | `8080`            | `8080`                                          | API daemon                                  |
| `STATE_DIR`              | The directory where M3TAL stores its state database and configuration files.                                           | `./state`         | `/var/lib/m3tal/state`                          | API daemon, Dashboard container             |
| `LOG_LEVEL`              | The verbosity of logging for M3TAL services.                                                                           | `info`            | `debug`                                         | API daemon                                  |
| `DEBUG_MODE`             | Enables debug mode for M3TAL services.                                                                                 | `false`           | `true`                                          | API daemon                                  |
| `METRICS_ENABLED`        | Enables the collection and exposure of service metrics.                                                                | `true`            | `false`                                         | API daemon                                  |
| **Auth**                 |                                                                                                                        |                   |                                                 |                                             |
| `DASHBOARD_SECRET`       | A secret key used for securing the dashboard session. **Auto-generated on first `m3tal init`.**                       | `change_me_immediately` | `your_super_secret_dashboard_key`               | Dashboard container                         |
| `API_TOKEN`              | A token used for authenticating API requests. **Auto-generated on first `m3tal init`.**                                | `change_me_api_token` | `your_api_authentication_token`                 | API daemon                                  |
| `ADMIN_PASSWORD`         | The password for the default administrator user.                                                                       | `admin_pass`      | `your_strong_admin_password`                    | API daemon, Dashboard container             |
| **Network**              |                                                                                                                        |                   |                                                 |                                             |
| `NETWORK_NAME`           | The name of the Docker network used by M3TAL services.                                                                 | `m3tal`           | `m3tal-net`                                     | Docker Compose, API daemon                  |
| `LOCAL_IP`               | The local IP address of the host machine, used for internal service communication.                                     | `127.0.0.1`       | `192.168.1.100`                                 | API daemon (host.docker.internal mapping) |
| `DOMAIN`                 | The primary domain name for accessing M3TAL services via Traefik. Setting this enables `dash.DOMAIN` and `api.DOMAIN`. | `localhost`       | `example.com`                                   | Traefik, API daemon                         |
| **Storage**              |                                                                                                                        |                   |                                                 |                                             |
| `BASE_STORAGE_PATH`      | The base directory for storing M3TAL data, including configuration, media, and downloads. Defaults to `/mnt` in production. | `./data`          | `/mnt/m3tal_storage`                            | All components, Docker volumes              |
| `MEDIA_PATH`             | The subdirectory within `BASE_STORAGE_PATH` for storing media files.                                                   | `./data/media`    | `/mnt/m3tal_storage/media`                      | All components, Docker volumes              |
| `CONFIG_PATH`            | The subdirectory within `BASE_STORAGE_PATH` for storing configuration files.                                           | `./data/config`   | `/mnt/m3tal_storage/config`                     | All components, Docker volumes              |
| `DOWNLOADS_PATH`         | The subdirectory within `BASE_STORAGE_PATH` for storing downloaded files.                                              | `./data/downloads`| `/mnt/m3tal_storage/downloads`                  | All components, Docker volumes              |
| `PUID`                   | The User ID to run Docker containers with.                                                                             | `1000`            | `1000`                                          | Docker Compose                              |
| `PGID`                   | The Group ID to run Docker containers with.                                                                            | `1000`            | `1000`                                          | Docker Compose                              |
| **Traefik**              |                                                                                                                        |                   |                                                 |                                             |
| `TRAEFIK_WEB_PORT`       | The host port Traefik listens on for HTTP traffic.                                                                     | `80`              | `80`                                            | Traefik container                           |
| `TRAEFIK_WEBHTTPS_PORT`  | The host port Traefik listens on for HTTPS traffic.                                                                    | `443`             | `443`                                           | Traefik container                           |
| `TRAEFIK_DASHBOARD_PORT` | The host port Traefik exposes its own dashboard on (internal access only).                                             | `8080`            | `8080`                                          | Traefik container                           |
| **VPN**                  |                                                                                                                        |                   |                                                 |                                             |
| `VPN_USER`               | The username for the VPN connection.                                                                                   | `user`            | `your_vpn_username`                             | VPN client (if configured)                  |
| `VPN_PASSWORD`           | The password for the VPN connection.                                                                                   | `password`        | `your_vpn_password`                             | VPN client (if configured)                  |
| **System**               |                                                                                                                        |                   |                                                 |                                             |
| `TZ`                     | The timezone to use for M3TAL services.                                                                                | `America/Denver`  | `UTC`                                           | All containers                              |

---

## Detailed Environment Variable Descriptions

All environment variables are read from `/etc/m3tal/.env`.

### Core

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port the M3TAL dashboard service listens on.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container.

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the dashboard is exposed to the network.
        *   `local`: Exposes the dashboard directly via a port binding. Access via `http://HOST_IP:DASHBOARD_PORT`. This mode does not require Traefik.
        *   `traefik`: Exposes the dashboard through Traefik, allowing access via `http://dash.DOMAIN`. This mode requires Traefik to be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-dashboard` container (for Traefik labels), Traefik gateway (for routing rules).

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon (Go binary) listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api.service`.

*   **`STATE_DIR`**
    *   **Description:** Specifies the directory where M3TAL stores its persistent state, including the SQLite database (`state.db`) and configuration files.
    *   **Default Value:** `./state` (relative to M3TAL's working directory, typically `/opt/m3tal/stack/`)
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`LOG_LEVEL`**
    *   **Description:** Determines the verbosity of logs generated by M3TAL services. Common values include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`.

*   **`DEBUG_MODE`**
    *   **Description:** Enables additional debugging output and features within M3TAL services.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api.service`.

*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether M3TAL services expose metrics for monitoring (e.g., Prometheus).
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api.service`.

### Auth

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used to sign and verify session cookies for the M3TAL dashboard. **This variable is auto-generated on the first `m3tal init` and should not be set manually unless you intend to rotate it.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `your_super_secret_dashboard_key`
    *   **Used By:** `m3tal-dashboard` container.

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating requests to the M3TAL API daemon. **This variable is auto-generated on the first `m3tal init` and should not be set manually unless you intend to rotate it.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `your_api_authentication_token`
    *   **Used By:** `m3tal-api.service` (for internal use or potential external API access).

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator user account.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `your_strong_admin_password`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### Network

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services are connected to. This is crucial for inter-service communication within Docker.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-net`
    *   **Used By:** Docker Compose files, `m3tal-api.service`.

*   **`LOCAL_IP`**
    *   **Description:** The IP address of the host machine that Docker containers can use to reach services running directly on the host. This is often used to map `host.docker.internal`.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** `m3tal-api.service` (e.g., `GO_API_URL=${GO_API_URL:-http://${LOCAL_IP:-host.docker.internal}:${HTTP_PORT:-8080}}`).

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL installation. Setting this variable is essential for Traefik to correctly route traffic to services like the dashboard and API. When set, Traefik will configure routes for `dash.DOMAIN` and `api.DOMAIN`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.example.com`
    *   **Used By:** Traefik gateway (routing rules), `m3tal-dashboard` container (Traefik labels).

### Storage

*   **`BASE_STORAGE_PATH`**
    *   **Description:** This environment variable controls the root directory where M3TAL stores all its persistent data, including configuration files, media, and downloads. **In production deployments, this defaults to `/mnt` to ensure data is stored on a dedicated or persistent volume, not a local subdirectory like `./data` which is typical in template or development environments.**
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt/m3tal_storage`
    *   **Used By:** All components that require persistent storage (API daemon, Dashboard container), Docker volumes configuration.

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` dedicated to storing media files.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/m3tal_storage/media`
    *   **Used By:** All components that handle or store media data, Docker volumes configuration.

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/m3tal_storage/config`
    *   **Used By:** All components that require access to configuration, Docker volumes configuration.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` used for storing any files downloaded by M3TAL services.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/m3tal_storage/downloads`
    *   **Used By:** All components that perform downloads, Docker volumes configuration.

*   **`PUID`**
    *   **Description:** The User ID (UID) that Docker containers should run as. This helps in managing file permissions for mounted volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** Docker Compose files for container execution context.

*   **`PGID`**
    *   **Description:** The Group ID (GID) that Docker containers should run as. This complements `PUID` for managing file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** Docker Compose files for container execution context.

### Traefik

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik binds to for incoming HTTP (non-encrypted) traffic. This is the primary entry point for services exposed via Traefik.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik container (static configuration).

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik binds to for incoming HTTPS (encrypted) traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik container (static configuration).

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which the Traefik dashboard itself is exposed within the Docker network. This is typically accessed via `http://host.docker.internal:TRAEFIK_DASHBOARD_PORT`.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik container (static configuration).

### VPN

*   **`VPN_USER`**
    *   **Description:** The username required for establishing a VPN connection, if a VPN client is configured and used by M3TAL services.
    *   **Default Value:** `user`
    *   **Example Value:** `your_vpn_username`
    *   **Used By:** VPN client configuration (if applicable).

*   **`VPN_PASSWORD`**
    *   **Description:** The password required for establishing a VPN connection, if a VPN client is configured and used by M3TAL services.
    *   **Default Value:** `password`
    *   **Example Value:** `your_vpn_password`
    *   **Used By:** VPN client configuration (if applicable).

### System

*   **`TZ`**
    *   **Description:** Specifies the timezone for all M3TAL containers. This ensures consistent timestamping in logs and for time-sensitive operations.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** All containers (set as environment variable in Docker Compose).