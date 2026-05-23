# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables are read from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks via `--env-file`.

The M3TAL CLI manages your environment configuration. For most users, it's recommended to use `m3tal config wizard` or `m3tal config set KEY value` to manage these variables.

## Quick Reference Table

| Variable Name           | Default Value     | Description                                                                                                                               |
|-------------------------|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| **Core**                |                   |                                                                                                                                           |
| `STATE_DIR`             | `./state`         | The directory where M3TAL stores its state database and other internal files.                                                             |
| `LOG_LEVEL`             | `info`            | Controls the verbosity of M3TAL's logs. Options include `debug`, `info`, `warn`, `error`.                                              |
| `DEBUG_MODE`            | `false`           | Enables or disables debug mode for the M3TAL API.                                                                                         |
| `METRICS_ENABLED`       | `true`            | Enables or disables the collection and exposure of system metrics.                                                                        |
| **Auth**                |                   |                                                                                                                                           |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret key for securing the M3TAL Dashboard. **Auto-generated on first `m3tal init`. Rotate if compromised.**                         |
| `API_TOKEN`             | `change_me_api_token` | API authentication token. **Auto-generated on first `m3tal init`. Rotate if compromised.**                                             |
| `ADMIN_PASSWORD`        | `admin_pass`      | Password for the default admin user of the M3TAL Dashboard.                                                                               |
| **Network**             |                   |                                                                                                                                           |
| `HTTP_PORT`             | `8080`            | The port on which the M3TAL API daemon listens.                                                                                           |
| `NETWORK_NAME`          | `m3tal`           | The name of the Docker network used by M3TAL services.                                                                                    |
| `LOCAL_IP`              | `127.0.0.1`       | The local IP address used for internal service communication. Typically `host.docker.internal` in Docker, but `127.0.0.1` for simplicity. |
| **Storage**             |                   |                                                                                                                                           |
| `BASE_STORAGE_PATH`     | `./data`          | The root directory for M3TAL's persistent data. **Defaults to `/mnt` in production deployments.**                                         |
| `MEDIA_PATH`            | `./data/media`    | The directory within `BASE_STORAGE_PATH` for storing media files.                                                                         |
| `CONFIG_PATH`           | `./data/config`   | The directory within `BASE_STORAGE_PATH` for storing configuration files.                                                                 |
| `DOWNLOADS_PATH`        | `./data/downloads`| The directory within `BASE_STORAGE_PATH` for downloaded files.                                                                            |
| `PUID`                  | `1000`            | The user ID for running containers. Matches the host user ID to avoid permission issues.                                                  |
| `PGID`                  | `1000`            | The group ID for running containers. Matches the host group ID to avoid permission issues.                                                |
| **Traefik**             |                   |                                                                                                                                           |
| `DOMAIN`                | `localhost`       | The base domain for Traefik routing. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes.                                          |
| `DASHBOARD_PORT`        | `8082`            | The internal port the M3TAL Dashboard service listens on.                                                                                 |
| `DASHBOARD_EXPOSE_MODE` | `local`           | Controls how the M3TAL Dashboard is exposed. `local` uses direct port binding; `traefik` uses Traefik routing.                           |
| `TRAEFIK_WEB_PORT`      | `80`              | The host port Traefik uses for HTTP traffic.                                                                                              |
| `TRAEFIK_WEBHTTPS_PORT` | `443`             | The host port Traefik uses for HTTPS traffic.                                                                                             |
| `TRAEFIK_DASHBOARD_PORT`| `8080`            | The host port Traefik uses internally to expose its own dashboard.                                                                        |
| **VPN**                 |                   |                                                                                                                                           |
| `VPN_USER`              | `user`            | Username for VPN authentication.                                                                                                          |
| `VPN_PASSWORD`          | `password`        | Password for VPN authentication.                                                                                                          |
| **System**              |                   |                                                                                                                                           |
| `TZ`                    | `America/Denver`  | The timezone to use for M3TAL services.                                                                                                   |

---

## Detailed Variable Reference

### Core

This group of variables controls fundamental aspects of the M3TAL system's operation and logging.

*   **`STATE_DIR`**
    *   **Description:** The directory where M3TAL stores its state database (e.g., `/var/lib/m3tal/state.db`) and other internal files.
    *   **Default Value:** `./state`
    *   **Example Value:** `/mnt/config/m3tal/state`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`LOG_LEVEL`**
    *   **Description:** Controls the verbosity of M3TAL's logs. Higher verbosity provides more detailed information, which can be useful for debugging.
    *   **Options:** `debug`, `info`, `warn`, `error`
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode for the M3TAL API. When enabled, the API may provide more verbose error messages or enable additional debugging endpoints.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api.service`

*   **`METRICS_ENABLED`**
    *   **Description:** Enables or disables the collection and exposure of system metrics. If enabled, metrics can be scraped by monitoring tools.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api.service`

### Auth

Authentication and authorization settings for accessing M3TAL services.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for securing sessions and other sensitive operations within the M3TAL Dashboard. **This variable is auto-generated on the first `m3tal init` and should NOT be set manually unless you intend to rotate the secret.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_very_long_and_random_secret_string_generated_by_m3tal`
    *   **Used By:** `m3tal-dashboard` container

*   **`API_TOKEN`**
    *   **Description:** An API authentication token used by clients to authenticate with the M3TAL API. **This variable is auto-generated on the first `m3tal init` and should NOT be set manually unless you intend to rotate the token.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `a_secure_api_token_generated_by_m3tal`
    *   **Used By:** `m3tal-api.service`

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default `admin` user of the M3TAL Dashboard. It's highly recommended to change this from the default.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `MySecureDashboardPassword123!`
    *   **Used By:** `m3tal-dashboard` container

### Network

Variables related to network configuration and communication ports.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `9000`
    *   **Used By:** `m3tal-api.service`

*   **`NETWORK_NAME`**
    *   **Description:** Specifies the name of the Docker network that M3TAL services will join. This ensures proper communication between containers.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-internal-net`
    *   **Used By:** All M3TAL Docker Compose stacks.

*   **`LOCAL_IP`**
    *   **Description:** The local IP address used for internal service communication. For Docker, `host.docker.internal` is often used, but `127.0.0.1` is also supported for simplicity or specific configurations.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `host.docker.internal`
    *   **Used By:** `m3tal-dashboard` container, Traefik routing configuration.

### Storage

These variables define the locations for M3TAL's persistent data and how containerized services should handle user permissions.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory where M3TAL stores all its persistent data, including configuration, media, and downloads. **In production deployments, this defaults to `/mnt` to ensure data is stored on a persistent volume. In development/template setups, it may default to `./data`.**
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt/m3tal-data`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container, all M3TAL Docker Compose stacks.

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files (e.g., uploaded images, video files) are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/m3tal-data/media`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where configuration files and other persistent settings are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/m3tal-data/config`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files (e.g., from external sources) are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/m3tal-data/downloads`
    *   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

*   **`PUID`**
    *   **Description:** The User ID (UID) that containerized M3TAL services will run as. Setting this to your host user's UID (e.g., from `id -u`) is crucial for ensuring correct file permissions on your host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000` (if your host user ID is 1000)
    *   **Used By:** `m3tal-dashboard` container, all M3TAL Docker Compose stacks.

*   **`PGID`**
    *   **Description:** The Group ID (GID) that containerized M3TAL services will run as. Setting this to your host user's GID (e.g., from `id -g`) is crucial for ensuring correct file permissions on your host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000` (if your host group ID is 1000)
    *   **Used By:** `m3tal-dashboard` container, all M3TAL Docker Compose stacks.

### Traefik

Configuration variables specifically for the Traefik reverse proxy. These are essential for exposing M3TAL services via domain names.

*   **`DOMAIN`**
    *   **Description:** The base domain name used for routing traffic to M3TAL services. When set, Traefik will attempt to route requests for `api.DOMAIN` to the M3TAL API and `dash.DOMAIN` to the M3TAL Dashboard. **Setting this enables domain-based access.**
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.mydomain.com`
    *   **Used By:** Traefik routing configuration, `m3tal-dashboard` container (via labels).

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port on which the M3TAL Dashboard service listens. This is the port Traefik will connect to when routing traffic to the dashboard.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container, Traefik routing configuration.

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL Dashboard is made accessible.
        *   `local`: The dashboard is exposed via a direct port binding (`DASHBOARD_PORT`) on the host. Access is typically via `http://HOST_IP:DASHBOARD_PORT`. No Traefik is required for dashboard access in this mode.
        *   `traefik`: The dashboard is exposed via Traefik routing. Access is typically via `http://dash.DOMAIN`. This requires Traefik to be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-dashboard` container (via compose overrides).

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The port on the host machine that Traefik listens on for incoming HTTP traffic. This is the primary entry point for services exposed via Traefik.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik configuration.

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The port on the host machine that Traefik listens on for incoming HTTPS traffic. This is typically used when SSL/TLS is configured.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik configuration.

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The port on the host machine that Traefik uses internally to expose its own administrative dashboard. This is usually only accessible from the host itself.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik configuration.

### VPN

Settings related to VPN client configuration within M3TAL.

*   **`VPN_USER`**
    *   **Description:** The username required for authenticating with a configured VPN service.
    *   **Default Value:** `user`
    *   **Example Value:** `vpnuser123`
    *   **Used By:** VPN client configuration within M3TAL.

*   **`VPN_PASSWORD`**
    *   **Description:** The password required for authenticating with a configured VPN service.
    *   **Default Value:** `password`
    *   **Example Value:** `MyVPNPassword123`
    *   **Used By:** VPN client configuration within M3TAL.

### System

General system-wide settings that affect containerized services.

*   **`TZ`**
    *   **Description:** The timezone setting used by all M3TAL containers. This ensures that logs and timestamps are consistent.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`, `Europe/London`
    *   **Used By:** `m3tal-dashboard` container, `m3tal-api.service`.

---

### How Environment Variables are Managed

M3TAL utilizes a centralized `.env` file for all configuration.

1.  **Primary Configuration File:** All environment variables are read from `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` command or individual `m3tal config set KEY value` commands.
2.  **CLI and Compose Usage:** Both the `m3tal` CLI binary and all Docker Compose stacks (started via `m3tal up`) read their configuration from this `/etc/m3tal/.env` file using the `--env-file` option. This ensures consistency across the entire M3TAL ecosystem.

### Important Notes:

*   **Auto-Generated Secrets:** `DASHBOARD_SECRET` and `API_TOKEN` are automatically generated by M3TAL during the initial setup (`m3tal init`). You should **not** manually set these unless you intend to intentionally rotate them for security reasons.
*   **Storage Path:** The `BASE_STORAGE_PATH` variable is crucial. In production, it defaults to `/mnt` to ensure data is stored on a persistent volume rather than within the Docker overlay filesystem. This is different from the `./data` default seen in development templates.
*   **Domain Routing:** Setting the `DOMAIN` variable is key to enabling Traefik to route traffic to services like the API (`api.DOMAIN`) and Dashboard (`dash.DOMAIN`). If `DOMAIN` is left as `localhost` or not set, these domain-based routes will not function correctly.
*   **Dashboard Exposure:** The `DASHBOARD_EXPOSE_MODE` variable offers two distinct ways to access the M3TAL Dashboard: `local` for direct port access and `traefik` for domain-based access through the reverse proxy.