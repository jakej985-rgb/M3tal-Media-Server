# M3TAL Environment Variables Reference

All M3TAL configuration is managed through environment variables, primarily sourced from `/etc/m3tal/.env`. The `m3tal config wizard` command helps in setting these variables, and changes are applied across all M3TAL services and the CLI by reading from this central file using `--env-file`.

## Quick Reference Table

| Variable Name          | Default Value     | Description                                                                                                                                                                                                                                                                                                                      |
|------------------------|-------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Core**               |                   |                                                                                                                                                                                                                                                                                                                  |
| `HTTP_PORT`            | `8080`            | The port on which the M3TAL API daemon listens for incoming HTTP requests.                                                                                                                                                                                                                                         |
| `STATE_DIR`            | `./state`         | The directory where the M3TAL state database (`state.db`) is stored.                                                                                                                                                                                                                                             |
| `LOG_LEVEL`            | `info`            | The verbosity of logs. Options include `debug`, `info`, `warn`, `error`.                                                                                                                                                                                                                                         |
| `PUID`                 | `1000`            | The User ID for running containers. Should match the host user's UID for proper file permissions.                                                                                                                                                                                                                    |
| `PGID`                 | `1000`            | The Group ID for running containers. Should match the host user's GID for proper file permissions.                                                                                                                                                                                                                  |
| `TZ`                   | `America/Denver`  | The timezone to use for all services, ensuring consistent time logging and operations.                                                                                                                                                                                                                             |
| `DEBUG_MODE`           | `false`           | Enables debug mode for enhanced logging and diagnostics.                                                                                                                                                                                                                                                           |
| `METRICS_ENABLED`      | `true`            | Enables the collection and exposition of Prometheus metrics for monitoring.                                                                                                                                                                                                                                        |
| **Auth**               |                   |                                                                                                                                                                                                                                                                                                                  |
| `DASHBOARD_SECRET`     | `change_me_immediately` | A secret key used for signing dashboard session cookies. **Auto-generated on first `m3tal init`. Do not set manually unless rotating.**                                                                                                                                                                        |
| `API_TOKEN`            | `change_me_api_token` | A token used for authenticating API requests. **Auto-generated on first `m3tal init`. Do not set manually unless rotating.**                                                                                                                                                                                        |
| `ADMIN_PASSWORD`       | `admin_pass`      | The password for accessing the M3TAL dashboard.                                                                                                                                                                                                                                                                  |
| **Network**            |                   |                                                                                                                                                                                                                                                                                                                  |
| `NETWORK_NAME`         | `m3tal`           | The name of the Docker network used by M3TAL services.                                                                                                                                                                                                                                                           |
| `LOCAL_IP`             | `127.0.0.1`       | The IP address used for internal service discovery and communication. Often `host.docker.internal` is used for containers to reach the host.                                                                                                                                                                       |
| **Storage**            |                   |                                                                                                                                                                                                                                                                                                                  |
| `BASE_STORAGE_PATH`    | `./data`          | The root directory for all M3TAL persistent data. **Defaults to `/mnt` in production deployments.**                                                                                                                                                                                                              |
| `MEDIA_PATH`           | `./data/media`    | The sub-directory within `BASE_STORAGE_PATH` where media files are stored.                                                                                                                                                                                                                                       |
| `CONFIG_PATH`          | `./data/config`   | The sub-directory within `BASE_STORAGE_PATH` where configuration files (e.g., dashboard users.json) are stored.                                                                                                                                                                                                   |
| `DOWNLOADS_PATH`       | `./data/downloads` | The sub-directory within `BASE_STORAGE_PATH` where downloaded files are stored.                                                                                                                                                                                                                                  |
| **Dashboard**          |                   |                                                                                                                                                                                                                                                                                                                  |
| `DASHBOARD_PORT`       | `8082`            | The internal port the M3TAL dashboard container listens on.                                                                                                                                                                                                                                                        |
| `DASHBOARD_EXPOSE_MODE`| `local`           | Controls how the dashboard is exposed. `local` (default) uses direct port binding (`HOST_IP:8082`). `traefik` exposes it via Traefik at `dash.DOMAIN`.                                                                                                                                                              |
| **Traefik**            |                   |                                                                                                                                                                                                                                                                                                                  |
| `DOMAIN`               | `localhost`       | The primary domain name for M3TAL services. When set, enables routing for `dash.DOMAIN` and `api.DOMAIN` via Traefik.                                                                                                                                                                                             |
| `TRAEFIK_WEB_PORT`     | `80`              | The host port that Traefik listens on for incoming HTTP traffic.                                                                                                                                                                                                                                                   |
| `TRAEFIK_WEBHTTPS_PORT`| `443`             | The host port that Traefik listens on for incoming HTTPS traffic (if configured).                                                                                                                                                                                                                                  |
| `TRAEFIK_DASHBOARD_PORT`| `8080`            | The internal port Traefik listens on for its own dashboard (accessible via `127.0.0.1:8081` by default).                                                                                                                                                                                                             |
| **VPN**                |                   |                                                                                                                                                                                                                                                                                                                  |
| `VPN_USER`             | `user`            | The username for VPN authentication.                                                                                                                                                                                                                                                                             |
| `VPN_PASSWORD`         | `password`        | The password for VPN authentication.                                                                                                                                                                                                                                                                             |

---

## Detailed Environment Variable Reference

All M3TAL environment variables are read from `/etc/m3tal/.env` by both the CLI (`/usr/bin/m3tal`) and all Docker Compose stacks via the `--env-file` flag.

### Core Variables

These variables control fundamental aspects of the M3TAL system.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api.service`

*   **`STATE_DIR`**
    *   **Description:** The directory where the M3TAL state database (`state.db`) is stored. This database is managed by the API daemon.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`LOG_LEVEL`**
    *   **Description:** The verbosity of logs emitted by M3TAL services. Options include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`PUID`**
    *   **Description:** The User ID (UID) for running containers. It is highly recommended to set this to your host user's UID to ensure proper file ownership and permissions for mounted volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container, user-defined stacks

*   **`PGID`**
    *   **Description:** The Group ID (GID) for running containers. It is highly recommended to set this to your host user's GID to ensure proper file ownership and permissions for mounted volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Used By:** `m3tal-dashboard` container, user-defined stacks

*   **`TZ`**
    *   **Description:** The timezone to use for all services, ensuring consistent time logging and operations across the system.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** `m3tal-dashboard` container

*   **`DEBUG_MODE`**
    *   **Description:** Enables debug mode for enhanced logging and diagnostics across M3TAL services.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`METRICS_ENABLED`**
    *   **Description:** Enables the collection and exposition of Prometheus metrics for monitoring purposes.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api.service`

### Authentication Variables

These variables are critical for securing access to M3TAL services.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for signing dashboard session cookies. **This variable is auto-generated on the first `m3tal init` command. Users should NOT set it manually unless they are intentionally rotating it for security reasons.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_very_secure_random_string_generated_by_m3tal_init`
    *   **Used By:** `m3tal-dashboard` container

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating API requests. **This variable is auto-generated on the first `m3tal init` command. Users should NOT set it manually unless they are intentionally rotating it for security reasons.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `another_secure_random_token_from_m3tal_init`
    *   **Used By:** `m3tal-api.service` (for internal API authentication)

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for accessing the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `MySuperSecurePassword123!`
    *   **Used By:** `m3tal-dashboard` container

### Network Variables

These variables configure network-related settings for M3TAL.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network used by M3TAL services. This ensures services can communicate with each other.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-internal-net`
    *   **Used By:** All M3TAL Docker Compose stacks.

*   **`LOCAL_IP`**
    *   **Description:** The IP address used for internal service discovery and communication. In Docker environments, `host.docker.internal` is often used to allow containers to reach services running directly on the host machine.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `host.docker.internal`
    *   **Used By:** `m3tal-dashboard` container (for `GO_API_URL`), Traefik configuration.

### Storage Variables

These variables define the locations for persistent data storage.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory for all M3TAL persistent data. **This variable defaults to `/mnt` in production deployments. In template or development environments, it may default to `./data`.** This is where `MEDIA_PATH`, `CONFIG_PATH`, and `DOWNLOADS_PATH` are relative to.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Used By:** `m3tal-dashboard` container, user-defined stacks

*   **`MEDIA_PATH`**
    *   **Description:** The sub-directory within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Used By:** `m3tal-dashboard` container, user-defined stacks

*   **`CONFIG_PATH`**
    *   **Description:** The sub-directory within `BASE_STORAGE_PATH` where configuration files (e.g., dashboard `users.json`) are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Used By:** `m3tal-dashboard` container, user-defined stacks

*   **`DOWNLOADS_PATH`**
    *   **Description:** The sub-directory within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Used By:** User-defined stacks

### Dashboard Variables

These variables control the M3TAL dashboard's behavior and accessibility.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port that the M3TAL dashboard container listens on.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL dashboard is exposed to users.
        *   `local` (default): Uses a direct port binding, making the dashboard accessible via `http://HOST_IP:8082` or `http://localhost:8082`. This mode does not require Traefik to be running.
        *   `traefik`: Exposes the dashboard through Traefik, making it accessible via `http://dash.DOMAIN` (if `DOMAIN` is set). This mode requires Traefik to be running and configured.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-dashboard` container (via compose overrides)

### Traefik Variables

These variables are used to configure Traefik, M3TAL's reverse proxy.

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL installation. Setting this variable enables Traefik routing rules for subdomains like `dash.DOMAIN` and `api.DOMAIN`. If not set, Traefik will default to using `localhost` in its rules.
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.yourdomain.com`
    *   **Used By:** Traefik configuration (dynamic/api.yml, m3tal-compose.traefik.yml)

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTP traffic. This is typically port 80.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik configuration

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTPS traffic. This is typically port 443.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik configuration

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port that Traefik listens on for its own dashboard. This dashboard is usually accessed via `127.0.0.1:8081` and not exposed publicly.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik container configuration

### VPN Variables

These variables are used if you configure M3TAL to use a VPN for network access.

*   **`VPN_USER`**
    *   **Description:** The username for VPN authentication.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** VPN client configuration (if applicable)

*   **`VPN_PASSWORD`**
    *   **Description:** The password for VPN authentication.
    *   **Default Value:** `password`
    *   **Example Value:** `MyVpnP@ssword123`
    *   **Used By:** VPN client configuration (if applicable)