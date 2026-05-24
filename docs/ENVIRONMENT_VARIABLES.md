# M3TAL Environment Variables Reference

All M3TAL environment variables are read from `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` and can be updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks utilize this file via the `--env-file` flag.

## Quick Reference

| Variable Name          | Description                                                                                                                                                                                                                          | Default Value      |
|------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| **Core**               |                                                                                                                                                                                                                                      |                    |
| `DASHBOARD_PORT`       | The port on which the M3TAL dashboard will listen.                                                                                                                                                                                   | `8082`             |
| `HTTP_PORT`            | The port on which the M3TAL API daemon will listen.                                                                                                                                                                                  | `8080`             |
| `STATE_DIR`            | The directory where M3TAL stores its state database and other configuration files.                                                                                                                                                     | `./state`          |
| `LOG_LEVEL`            | The logging level for M3TAL services.                                                                                                                                                                                                | `info`             |
| `DEBUG_MODE`           | Enables debug mode for M3TAL services.                                                                                                                                                                                               | `false`            |
| `METRICS_ENABLED`      | Enables Prometheus metrics collection for M3TAL services.                                                                                                                                                                            | `true`             |
| **Auth**               |                                                                                                                                                                                                                                      |                    |
| `DASHBOARD_SECRET`     | A secret used for session management within the M3TAL dashboard. **Auto-generated on first `m3tal init`. Rotate manually if compromised.**                                                                                                | `change_me_immediately` |
| `API_TOKEN`            | A token used for authenticating API requests. **Auto-generated on first `m3tal init`. Rotate manually if compromised.**                                                                                                               | `change_me_api_token` |
| `ADMIN_PASSWORD`       | The password for the default administrative user of the M3TAL dashboard.                                                                                                                                                             | `admin_pass`       |
| **Network**            |                                                                                                                                                                                                                                      |                    |
| `NETWORK_NAME`         | The name of the Docker network used by M3TAL services.                                                                                                                                                                               | `m3tal`            |
| `LOCAL_IP`             | The IP address to bind M3TAL services to locally.                                                                                                                                                                                    | `127.0.0.1`        |
| `DOMAIN`               | The primary domain name for M3TAL services. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routing via Traefik.                                                                                                                | `localhost`        |
| `VPN_USER`             | Username for the VPN connection.                                                                                                                                                                                                     | `user`             |
| `VPN_PASSWORD`         | Password for the VPN connection.                                                                                                                                                                                                     | `password`         |
| **Storage**            |                                                                                                                                                                                                                                      |                    |
| `BASE_STORAGE_PATH`    | The base directory for storing all M3TAL data, including media, configuration, and downloads. **Defaults to `/mnt` in production deployments.**                                                                                       | `./data`           |
| `MEDIA_PATH`           | The path within `BASE_STORAGE_PATH` where media files are stored.                                                                                                                                                                    | `./data/media`     |
| `CONFIG_PATH`          | The path within `BASE_STORAGE_PATH` where M3TAL configuration files are stored.                                                                                                                                                      | `./data/config`    |
| `DOWNLOADS_PATH`       | The path within `BASE_STORAGE_PATH` where downloaded files are stored.                                                                                                                                                               | `./data/downloads` |
| `PUID`                 | The User ID (UID) to run Docker containers with.                                                                                                                                                                                     | `1000`             |
| `PGID`                 | The Group ID (GID) to run Docker containers with.                                                                                                                                                                                    | `1000`             |
| **Traefik**            |                                                                                                                                                                                                                                      |                    |
| `TRAEFIK_WEB_PORT`     | The host port that Traefik will use for HTTP traffic.                                                                                                                                                                                | `80`               |
| `TRAEFIK_WEBHTTPS_PORT`| The host port that Traefik will use for HTTPS traffic.                                                                                                                                                                               | `443`              |
| `TRAEFIK_DASHBOARD_PORT`| The host port on which the Traefik dashboard is accessible (usually for local access).                                                                                                                                               | `8080`             |
| **System**             |                                                                                                                                                                                                                                      |                    |
| `DASHBOARD_EXPOSE_MODE`| Controls how the M3TAL dashboard is exposed: `local` (direct port binding) or `traefik` (via Traefik routing).                                                                                                                      | `local`            |
| `TZ`                   | The timezone to use for M3TAL services.                                                                                                                                                                                              | `America/Denver`  |

---

## Detailed Environment Variables Reference

### Core

*   **`DASHBOARD_PORT`**
    *   **Description:** The port on which the M3TAL dashboard will listen.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container, Traefik routing rules (when `DASHBOARD_EXPOSE_MODE=traefik`).

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon will listen.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api.service`, Traefik routing rules (`api.DOMAIN`).

*   **`STATE_DIR`**
    *   **Description:** The directory where M3TAL stores its state database (`state.db`) and other configuration files.
    *   **Default Value:** `./state`
    *   **Example Value:** `/mnt/config/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`LOG_LEVEL`**
    *   **Description:** The logging level for M3TAL services. Common values include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`DEBUG_MODE`**
    *   **Description:** Enables debug mode for M3TAL services, which can provide more verbose logging and enable additional debugging features.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`METRICS_ENABLED`**
    *   **Description:** Enables Prometheus metrics collection for M3TAL services. If enabled, metrics will be exposed on a dedicated port for scraping.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api.service`.

### Auth

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret used for session management within the M3TAL dashboard. **This is auto-generated on the first `m3tal init` and should not be set manually unless you are rotating it due to a security concern.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `somerandomsupersecretstring`
    *   **Used By:** `m3tal-dashboard` container.

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating API requests. **This is auto-generated on the first `m3tal init` and should not be set manually unless you are rotating it due to a security concern.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `anotherrandomsecuretoken`
    *   **Used By:** `m3tal-api.service` (internally for authentication).

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrative user of the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `mysecurepassword123`
    *   **Used By:** `m3tal-dashboard` container.

### Network

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network used by M3TAL services. This ensures services can communicate with each other.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal_network`
    *   **Used By:** All M3TAL Docker Compose stacks.

*   **`LOCAL_IP`**
    *   **Description:** The IP address to bind M3TAL services to locally. This is typically `127.0.0.1` for local access.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `127.0.0.1`
    *   **Used By:** Traefik configuration, `m3tal-api.service`.

*   **`DOMAIN`**
    *   **Description:** The primary domain name for M3TAL services. Setting this variable enables routing rules in Traefik for `dash.DOMAIN` and `api.DOMAIN`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.example.com`
    *   **Used By:** Traefik routing rules, `m3tal-dashboard` container (for constructing URLs).

*   **`VPN_USER`**
    *   **Description:** Username for the VPN connection, if M3TAL is configured to use a VPN.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** VPN client configuration within Docker images.

*   **`VPN_PASSWORD`**
    *   **Description:** Password for the VPN connection, if M3TAL is configured to use a VPN.
    *   **Default Value:** `password`
    *   **Example Value:** `myvpnpassword`
    *   **Used By:** VPN client configuration within Docker images.

### Storage

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The base directory for storing all M3TAL data, including media, configuration, and downloads. **This variable controls where media data is stored. In production deployments, it defaults to `/mnt`, not `./data` as in the template.**
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container, and generally all services for persistent storage.

*   **`MEDIA_PATH`**
    *   **Description:** The path within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`CONFIG_PATH`**
    *   **Description:** The path within `BASE_STORAGE_PATH` where M3TAL configuration files are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The path within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`PUID`**
    *   **Description:** The User ID (UID) to run Docker containers with. This is important for file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container, and potentially other services to ensure consistent file ownership.

*   **`PGID`**
    *   **Description:** The Group ID (GID) to run Docker containers with. This is important for file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container, and potentially other services to ensure consistent file ownership.

### Traefik

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik will use for incoming HTTP traffic. This is the primary entry point for services exposed via Traefik.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik service configuration.

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik will use for incoming HTTPS traffic. This is typically used for SSL/TLS termination.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik service configuration.

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The host port on which the Traefik dashboard is accessible. This is usually bound to localhost for local access only.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik service configuration.

### System

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL dashboard is exposed to the network.
        *   `local`: The dashboard is exposed via a direct port binding (`${DASHBOARD_PORT}:8082`). Access via `http://HOST_IP:${DASHBOARD_PORT}`. No Traefik is required for this mode.
        *   `traefik`: The dashboard is exposed via Traefik routing rules. Access via `http://dash.${DOMAIN}`. Traefik must be running for this mode to work.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-dashboard` Docker Compose configuration, Traefik routing rules (when set to `traefik`).

*   **`TZ`**
    *   **Description:** The timezone to use for M3TAL services. This ensures logs and timestamps are accurate.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.