# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables are read from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks, typically managed by the `m3tal up` command.

## Quick Reference

| Variable Name           | Description                                                                                              | Default Value       |
|-------------------------|----------------------------------------------------------------------------------------------------------|---------------------|
| **Core**                |                                                                                                          |                     |
| `LOG_LEVEL`             | Sets the logging verbosity for M3TAL services.                                                           | `info`              |
| `DEBUG_MODE`            | Enables or disables debug mode for enhanced logging and diagnostics.                                     | `false`             |
| `METRICS_ENABLED`       | Enables or disables the collection and exposure of service metrics.                                      | `true`              |
| **Authentication**      |                                                                                                          |                     |
| `DASHBOARD_SECRET`      | Secret key for securing the M3TAL dashboard sessions. Auto-generated on first `m3tal init`.              | `change_me_immediately` |
| `API_TOKEN`             | API authentication token. Auto-generated on first `m3tal init`.                                          | `change_me_api_token` |
| `ADMIN_PASSWORD`        | Password for the default administrator user of the M3TAL dashboard.                                    | `admin_pass`        |
| **Network**             |                                                                                                          |                     |
| `HTTP_PORT`             | The port the M3TAL API daemon listens on.                                                                | `8080`              |
| `NETWORK_NAME`          | The name of the Docker network M3TAL services will use.                                                  | `m3tal`             |
| `LOCAL_IP`              | The IP address used for internal service communication, particularly for `host.docker.internal`.         | `127.0.0.1`         |
| `DOMAIN`                | The primary domain name for M3TAL services. Used by Traefik for routing rules.                           | `localhost`         |
| `VPN_USER`              | Username for VPN connection.                                                                             | `user`              |
| `VPN_PASSWORD`          | Password for VPN connection.                                                                             | `password`          |
| **Storage**             |                                                                                                          |                     |
| `BASE_STORAGE_PATH`     | The base directory for all M3TAL data storage. Defaults to `/mnt` in production.                         | `./data`            |
| `MEDIA_PATH`            | Path within `BASE_STORAGE_PATH` for storing media files.                                                 | `./data/media`      |
| `CONFIG_PATH`           | Path within `BASE_STORAGE_PATH` for storing M3TAL configuration files.                                   | `./data/config`     |
| `DOWNLOADS_PATH`        | Path within `BASE_STORAGE_PATH` for downloaded files.                                                    | `./data/downloads`  |
| `STATE_DIR`             | Directory for M3TAL state data, including the SQLite database. Note: in production, this is typically under `CONFIG_PATH`. | `./state`           |
| `PUID`                  | User ID for running Docker containers.                                                                   | `1000`              |
| `PGID`                  | Group ID for running Docker containers.                                                                  | `1000`              |
| `TZ`                    | Timezone for M3TAL services.                                                                             | `America/Denver`    |
| **Dashboard**           |                                                                                                          |                     |
| `DASHBOARD_PORT`        | The port the M3TAL dashboard container listens on.                                                       | `8082`              |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed: `local` (direct port) or `traefik` (via reverse proxy).           | `local`             |
| **Traefik**             |                                                                                                          |                     |
| `TRAEFIK_WEB_PORT`      | The host port Traefik listens on for HTTP traffic.                                                       | `80`                |
| `TRAEFIK_WEBHTTPS_PORT` | The host port Traefik listens on for HTTPS traffic.                                                      | `443`               |
| `TRAEFIK_DASHBOARD_PORT`| The port Traefik's own dashboard is exposed on internally.                                               | `8080`              |

---

## Detailed Environment Variable Descriptions

All M3TAL environment variables are managed through the `/etc/m3tal/.env` file. This file is the central configuration hub for the M3TAL ecosystem and is read by the M3TAL CLI and all Docker Compose stacks via `--env-file`.

### Core

*   **`LOG_LEVEL`**
    *   **Description:** Sets the logging verbosity for M3TAL services. Accepted values typically include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Component(s) Using:** `m3tal-api.service`, `m3tal-dashboard`

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode for enhanced logging and diagnostics across M3TAL services.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Component(s) Using:** `m3tal-api.service`, `m3tal-dashboard`

*   **`METRICS_ENABLED`**
    *   **Description:** Enables or disables the collection and exposure of service metrics (e.g., for Prometheus).
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Component(s) Using:** `m3tal-api.service`

### Authentication

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for signing and encrypting session cookies for the M3TAL dashboard. **This variable is auto-generated on the first `m3tal init` and should generally not be set manually unless you are rotating credentials.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_long_and_random_secret_key_for_security`
    *   **Component(s) Using:** `m3tal-dashboard`

*   **`API_TOKEN`**
    *   **Description:** An authentication token used by the dashboard and potentially other services to authenticate with the M3TAL API daemon. **This variable is auto-generated on the first `m3tal init` and should generally not be set manually unless you are rotating credentials.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `a_secure_api_authentication_token_string`
    *   **Component(s) Using:** `m3tal-api.service`, `m3tal-dashboard`

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator user account in the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_strong_admin_password`
    *   **Component(s) Using:** `m3tal-dashboard`

### Network

*   **`HTTP_PORT`**
    *   **Description:** The TCP port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Component(s) Using:** `m3tal-api.service`

*   **`NETWORK_NAME`**
    *   **Description:** Specifies the name of the Docker network that M3TAL services will join. This is crucial for inter-container communication.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-network`
    *   **Component(s) Using:** All M3TAL Docker containers

*   **`LOCAL_IP`**
    *   **Description:** The IP address that M3TAL services will use for internal communication. This is particularly important for the dashboard to communicate with the API daemon via `host.docker.internal`.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `127.0.0.1`
    *   **Component(s) Using:** `m3tal-dashboard`

*   **`DOMAIN`**
    *   **Description:** The primary domain name configured for your M3TAL installation. This variable is essential for Traefik to set up correct routing rules. For example, if `DOMAIN` is `m3tal.example.com`, Traefik will attempt to route `api.m3tal.example.com` to the API and `dash.m3tal.example.com` to the dashboard (when in `traefik` mode).
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.example.com`
    *   **Component(s) Using:** `traefik`, `m3tal-dashboard` (via Traefik labels)

*   **`VPN_USER`**
    *   **Description:** The username to be used for a configured VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Component(s) Using:** Potentially used by specific M3TAL services that rely on VPN connectivity for external access.

*   **`VPN_PASSWORD`**
    *   **Description:** The password to be used for a configured VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `my_very_secure_vpn_password`
    *   **Component(s) Using:** Potentially used by specific M3TAL services that rely on VPN connectivity for external access.

### Storage

*   **`BASE_STORAGE_PATH`**
    *   **Description:** This variable controls the root directory where all M3TAL persistent data is stored. In production deployments, this defaults to `/mnt`. In development or template setups, it may default to `./data`.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Component(s) Using:** All M3TAL Docker containers (for persistent volumes)

*   **`MEDIA_PATH`**
    *   **Description:** The sub-path within `BASE_STORAGE_PATH` where media files (e.g., uploaded assets, generated content) are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `./data/media`
    *   **Component(s) Using:** All M3TAL Docker containers

*   **`CONFIG_PATH`**
    *   **Description:** The sub-path within `BASE_STORAGE_PATH` where M3TAL configuration files and potentially state databases are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Component(s) Using:** All M3TAL Docker containers

*   **`DOWNLOADS_PATH`**
    *   **Description:** The sub-path within `BASE_STORAGE_PATH` designated for storing downloaded files.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `./data/downloads`
    *   **Component(s) Using:** All M3TAL Docker containers

*   **`STATE_DIR`**
    *   **Description:** Specifies the directory where M3TAL stores its state information, most notably the SQLite database (`state.db`). In production environments, this path is typically a subdirectory within `CONFIG_PATH`.
    *   **Default Value:** `./state`
    *   **Example Value:** `/mnt/config/m3tal/state`
    *   **Component(s) Using:** `m3tal-api.service`, `m3tal-dashboard`

*   **`PUID`**
    *   **Description:** The User ID (UID) that Docker containers will run as. This is important for ensuring correct file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Component(s) Using:** All M3TAL Docker containers

*   **`PGID`**
    *   **Description:** The Group ID (GID) that Docker containers will run as. This is important for ensuring correct file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1000`
    *   **Component(s) Using:** All M3TAL Docker containers

*   **`TZ`**
    *   **Description:** Sets the timezone for all M3TAL services, ensuring consistent time logging and scheduling.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `Europe/London`
    *   **Component(s) Using:** `m3tal-api.service`, `m3tal-dashboard`

### Dashboard

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal TCP port that the M3TAL dashboard container listens on.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Component(s) Using:** `m3tal-dashboard`

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Determines how the M3TAL dashboard is made accessible.
        *   `local`: The dashboard is exposed directly via a host port (defined by `DASHBOARD_PORT`). Access is typically via `http://HOST_IP:DASHBOARD_PORT`. This mode does not require Traefik.
        *   `traefik`: The dashboard is exposed via the Traefik reverse proxy. Access is via `http://dash.DOMAIN`. This mode requires Traefik to be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Component(s) Using:** `m3tal-dashboard` (influences compose configuration), `traefik` (via labels when set to `traefik`)

### Traefik

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik uses as its entry point for incoming HTTP (unencrypted) traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Component(s) Using:** `traefik`

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik uses as its entry point for incoming HTTPS (encrypted) traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Component(s) Using:** `traefik`

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which the Traefik dashboard service itself listens. This is typically exposed on `127.0.0.1:8081` by default for local access to the Traefik dashboard.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Component(s) Using:** `traefik`