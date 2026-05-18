# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables. These variables are loaded from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks via the `--env-file` argument. It is strongly recommended to manage these variables using the `m3tal config wizard` or `m3tal config set KEY VALUE` commands.

## Quick Reference Table

| Variable Name          | Default Value    | Description                                                                                                                                                              |
|------------------------|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Core**               |                  |                                                                                                                                                                          |
| `STATE_DIR`            | `./state`        | Directory for storing M3TAL's internal state database.                                                                                                                   |
| `LOG_LEVEL`            | `info`           | Sets the verbosity of M3TAL's logs.                                                                                                                                      |
| `DEBUG_MODE`           | `false`          | Enables or disables debug mode for enhanced logging and diagnostics.                                                                                                       |
| `METRICS_ENABLED`      | `true`           | Enables or disables the collection and exposure of system metrics.                                                                                                         |
| **Authentication**     |                  |                                                                                                                                                                          |
| `DASHBOARD_SECRET`     | `change_me_immediately` | Secret key for signing dashboard session cookies. **Auto-generated on first `m3tal init`**. Manually set only for rotation.                                    |
| `API_TOKEN`            | `change_me_api_token`   | API authentication token for programmatic access. **Auto-generated on first `m3tal init`**. Manually set only for rotation.                                    |
| `ADMIN_PASSWORD`       | `admin_pass`     | Password for the default administrator user of the dashboard.                                                                                                            |
| **Network**            |                  |                                                                                                                                                                          |
| `HTTP_PORT`            | `8080`           | The port the M3TAL API daemon listens on.                                                                                                                                |
| `NETWORK_NAME`         | `m3tal`          | The name of the Docker network M3TAL services will use.                                                                                                                  |
| `LOCAL_IP`             | `127.0.0.1`      | The IP address M3TAL should bind to locally.                                                                                                                             |
| `DOMAIN`               | `localhost`      | The base domain for accessing M3TAL services. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.                                                   |
| **Storage**            |                  |                                                                                                                                                                          |
| `BASE_STORAGE_PATH`    | `./data`         | The base path for all M3TAL data storage. Defaults to `/mnt` in production deployments.                                                                                  |
| `MEDIA_PATH`           | `./data/media`   | Path within `BASE_STORAGE_PATH` for storing media files.                                                                                                                 |
| `CONFIG_PATH`          | `./data/config`  | Path within `BASE_STORAGE_PATH` for storing configuration files.                                                                                                         |
| `DOWNLOADS_PATH`       | `./data/downloads` | Path within `BASE_STORAGE_PATH` for downloaded files.                                                                                                                    |
| **Dashboard**          |                  |                                                                                                                                                                          |
| `DASHBOARD_PORT`       | `8082`           | The port the M3TAL Dashboard application listens on internally.                                                                                                          |
| `DASHBOARD_EXPOSE_MODE`| `local`          | Controls how the dashboard is exposed: `local` (direct port binding) or `traefik` (via Traefik reverse proxy).                                                          |
| **System Users**       |                  |                                                                                                                                                                          |
| `PUID`                 | `1000`           | The user ID for running Docker containers.                                                                                                                               |
| `PGID`                 | `1000`           | The group ID for running Docker containers.                                                                                                                              |
| `TZ`                   | `America/Denver` | The timezone to use for logging and other time-sensitive operations.                                                                                                     |
| **Traefik**            |                  |                                                                                                                                                                          |
| `TRAEFIK_WEB_PORT`     | `80`             | The host port Traefik listens on for HTTP traffic.                                                                                                                       |
| `TRAEFIK_WEBHTTPS_PORT`| `443`            | The host port Traefik listens on for HTTPS traffic.                                                                                                                      |
| `TRAEFIK_DASHBOARD_PORT`| `8080`           | The port Traefik's own dashboard is accessible on the host.                                                                                                              |
| **VPN (Placeholder)**  |                  |                                                                                                                                                                          |
| `VPN_USER`             | `user`           | Username for VPN authentication.                                                                                                                                         |
| `VPN_PASSWORD`         | `password`       | Password for VPN authentication.                                                                                                                                         |

---

## Detailed Environment Variable Reference

All environment variables are read from `/etc/m3tal/.env`.

### Core

These variables control fundamental aspects of the M3TAL system's operation and logging.

*   **`STATE_DIR`**
    *   **Description:** Specifies the directory where M3TAL stores its internal SQLite state database.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** API daemon (`m3tal-api.service`), CLI binary

*   **`LOG_LEVEL`**
    *   **Description:** Determines the verbosity of the logs produced by M3TAL components. Accepted values typically include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** API daemon, CLI binary

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode, which can provide more detailed logging and diagnostic information for troubleshooting.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** API daemon, CLI binary

*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether M3TAL exposes system and service metrics for monitoring.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** API daemon

### Authentication

These variables are crucial for securing your M3TAL installation.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for signing session cookies for the M3TAL Dashboard. **This variable is auto-generated on the first `m3tal init`**. Users should NOT set this manually unless they intend to rotate the secret.
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_very_secret_and_random_string_generated_by_m3tal`
    *   **Used By:** Dashboard container (`m3tal-dashboard`), API daemon

*   **`API_TOKEN`**
    *   **Description:** An authentication token used for programmatic access to the M3TAL API. **This variable is auto-generated on the first `m3tal init`**. Users should NOT set this manually unless they intend to rotate the token.
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `m3tal-api-token-abcdef123456`
    *   **Used By:** API daemon (for generating token), CLI binary (for authentication)

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator user of the M3TAL Dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_secure_dashboard_password`
    *   **Used By:** Dashboard container

### Network

These variables define network-related configurations for M3TAL services.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon (Go binary) listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** API daemon (`m3tal-api.service`)

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will use to communicate with each other.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal_internal_net`
    *   **Used By:** All M3TAL Docker Compose stacks

*   **`LOCAL_IP`**
    *   **Description:** The IP address to which M3TAL services should bind locally. This is often used for internal service-to-service communication.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** API daemon, Traefik configuration

*   **`DOMAIN`**
    *   **Description:** The base domain name for your M3TAL installation. Setting this variable is essential for enabling domain-based routing through Traefik. When set, M3TAL will configure Traefik to route traffic to `dash.${DOMAIN}` and `api.${DOMAIN}`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.mydomain.com`
    *   **Used By:** Traefik configuration, Dashboard container (via `GO_API_URL` environment variable which is often constructed using `DOMAIN`)

### Storage

These variables control where M3TAL stores its persistent data.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory for all M3TAL persistent data, including configuration, media, and downloads. In production deployments, this defaults to `/mnt` to ensure data is stored on a dedicated volume. In development or template environments, it may default to `./data`.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt/m3tal_data`
    *   **Used By:** All M3TAL Docker Compose stacks (volumes)

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `./data/media`
    *   **Used By:** All M3TAL Docker Compose stacks (volumes)

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where configuration files are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `./data/config`
    *   **Used By:** All M3TAL Docker Compose stacks (volumes)

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `./data/downloads`
    *   **Used By:** All M3TAL Docker Compose stacks (volumes)

### Dashboard

These variables control the M3TAL Dashboard's behavior and accessibility.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port on which the M3TAL Dashboard application listens.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** Dashboard container (`m3tal-dashboard`)

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Determines how the M3TAL Dashboard is made accessible.
        *   `local`: The dashboard is exposed via a direct port binding on the host. Access via `http://HOST_IP:DASHBOARD_PORT`. This mode does not require Traefik.
        *   `traefik`: The dashboard is exposed as a service routed by Traefik. Access via `http://dash.DOMAIN`. This mode requires Traefik to be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** Dashboard container (`m3tal-dashboard`) (influences compose override)

### System Users

These variables are used to configure the user and group IDs for running Docker containers, ensuring proper file permissions.

*   **`PUID`**
    *   **Description:** The User ID (UID) to run Docker containers as. This is important for file ownership and permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** All M3TAL Docker Compose stacks

*   **`PGID`**
    *   **Description:** The Group ID (GID) to run Docker containers as. This is important for file ownership and permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** All M3TAL Docker Compose stacks

*   **`TZ`**
    *   **Description:** The timezone to be used by M3TAL components for logging and any time-related operations.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** API daemon, Dashboard container

### Traefik

These variables configure the Traefik reverse proxy used for routing external traffic to M3TAL services.

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTP traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik container (`traefik`)

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTPS traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik container (`traefik`)

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which Traefik's own dashboard is accessible on the host. Note that Traefik's dashboard is typically exposed on `127.0.0.1:8081` by default.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik container (`traefik`)

### VPN (Placeholder)

These variables are placeholders and indicate future or optional VPN integration. Their current functionality and usage may vary.

*   **`VPN_USER`**
    *   **Description:** Username for VPN authentication.
    *   **Default Value:** `user`
    *   **Example Value:** `vpn_user_123`
    *   **Used By:** Potentially VPN client configuration

*   **`VPN_PASSWORD`**
    *   **Description:** Password for VPN authentication.
    *   **Default Value:** `password`
    *   **Example Value:** `my_vpn_password`
    *   **Used By:** Potentially VPN client configuration