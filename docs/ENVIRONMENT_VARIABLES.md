# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables are read from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks via the `--env-file` option.

## Quick Reference Table

| Variable Name            | Default Value      | Description                                                                                                   |
| ------------------------ | ------------------ | ------------------------------------------------------------------------------------------------------------- |
| **Core**                 |                    |                                                                                                               |
| `STATE_DIR`              | `./state`          | Directory for M3TAL state data (e.g., database).                                                              |
| `LOG_LEVEL`              | `info`             | Sets the logging verbosity for M3TAL services.                                                                |
| `PUID`                   | `1000`             | User ID for container processes.                                                                              |
| `PGID`                   | `1000`             | Group ID for container processes.                                                                             |
| `TZ`                     | `America/Denver`   | Timezone for M3TAL services.                                                                                  |
| `DEBUG_MODE`             | `false`            | Enables or disables debug mode for M3TAL services.                                                            |
| `METRICS_ENABLED`        | `true`             | Enables or disables metrics collection for M3TAL services.                                                    |
| **Auth**                 |                    |                                                                                                               |
| `DASHBOARD_SECRET`       | `change_me_immediately` | Secret for securing the dashboard. **Auto-generated on first `m3tal init`. Should NOT be set manually.**       |
| `API_TOKEN`              | `change_me_api_token` | API token for authenticating with the M3TAL API. **Auto-generated on first `m3tal init`. Should NOT be set manually.** |
| `ADMIN_PASSWORD`         | `admin_pass`       | Password for the admin user of the dashboard.                                                                 |
| **Network**              |                    |                                                                                                               |
| `HTTP_PORT`              | `8080`             | Port for the M3TAL API daemon.                                                                                |
| `NETWORK_NAME`           | `m3tal`            | Name of the Docker network used by M3TAL services.                                                            |
| `LOCAL_IP`               | `127.0.0.1`        | Local IP address to bind services to.                                                                         |
| `DOMAIN`                 | `localhost`        | Base domain for Traefik routing. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes.                |
| `VPN_USER`               | `user`             | Username for VPN connection.                                                                                  |
| `VPN_PASSWORD`           | `password`         | Password for VPN connection.                                                                                  |
| **Storage**              |                    |                                                                                                               |
| `BASE_STORAGE_PATH`      | `./data`           | Base path for all M3TAL data storage. Defaults to `/mnt` in production deployments.                           |
| `MEDIA_PATH`             | `./data/media`     | Path for storing media files.                                                                                 |
| `CONFIG_PATH`            | `./data/config`    | Path for storing M3TAL configuration files.                                                                 |
| `DOWNLOADS_PATH`         | `./data/downloads` | Path for storing downloaded files.                                                                            |
| **Dashboard**            |                    |                                                                                                               |
| `DASHBOARD_PORT`         | `8082`             | Port the dashboard container listens on.                                                                      |
| `DASHBOARD_EXPOSE_MODE`  | `local`            | Controls how the dashboard is exposed (`local` or `traefik`).                                                 |
| **Traefik**              |                    |                                                                                                               |
| `TRAEFIK_WEB_PORT`       | `80`               | Host port for Traefik's HTTP entry point.                                                                     |
| `TRAEFIK_WEBHTTPS_PORT`  | `443`              | Host port for Traefik's HTTPS entry point (if configured).                                                    |
| `TRAEFIK_DASHBOARD_PORT` | `8080`             | Internal port for the Traefik dashboard itself.                                                               |

---

All M3TAL environment variables are managed and read from the `/etc/m3tal/.env` file. This file is the primary configuration source for the M3TAL ecosystem. The M3TAL CLI (`/usr/bin/m3tal`) and all Docker Compose stacks (including the API daemon and dashboard) read their configurations from this single `.env` file.

The `m3tal config wizard` command can be used to interactively set these variables, and `m3tal config set KEY value` can be used for direct setting.

## Core Variables

These variables control fundamental aspects of the M3TAL system's operation and data management.

### `STATE_DIR`

*   **Description:** Specifies the directory where M3TAL stores its persistent state data, primarily the SQLite database.
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal/state`
*   **Used By:** API daemon, CLI

### `LOG_LEVEL`

*   **Description:** Determines the verbosity of logging output for M3TAL services. Accepted values typically include `debug`, `info`, `warn`, `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** API daemon, CLI, Dashboard

### `PUID`

*   **Description:** The User ID that M3TAL containers will run as. This is important for ensuring correct file permissions when containers interact with host volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** M3TAL Dashboard container

### `PGID`

*   **Description:** The Group ID that M3TAL containers will run as. Similar to `PUID`, this helps manage file permissions.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** M3TAL Dashboard container

### `TZ`

*   **Description:** Sets the timezone for M3TAL services. This ensures that logs and timestamps are recorded correctly.
*   **Default Value:** `America/Denver`
*   **Example Value:** `Europe/London`
*   **Used By:** M3TAL API daemon, M3TAL Dashboard container

### `DEBUG_MODE`

*   **Description:** Enables or disables debug mode for M3TAL services. When enabled, services may provide more detailed logging or enable additional diagnostic features.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** API daemon, CLI

### `METRICS_ENABLED`

*   **Description:** Controls whether M3TAL services expose metrics endpoints for monitoring.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** API daemon, CLI

## Auth Variables

These variables are crucial for securing access to your M3TAL instance.

### `DASHBOARD_SECRET`

*   **Description:** A secret key used for signing session cookies and other security-sensitive operations within the dashboard. **This variable is auto-generated on the first `m3tal init` and should NOT be set manually unless you intend to rotate it.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a_very_strong_and_unique_secret_key`
*   **Used By:** M3TAL Dashboard container

### `API_TOKEN`

*   **Description:** A token used for authenticating API requests to the M3TAL API daemon. **This variable is auto-generated on the first `m3tal init` and should NOT be set manually unless you intend to rotate it.**
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `another_super_secret_api_token_for_clients`
*   **Used By:** API daemon, CLI

### `ADMIN_PASSWORD`

*   **Description:** The password for the default administrative user of the M3TAL dashboard.
*   **Default Value:** `admin_pass`
*   **Example Value:** `your_secure_admin_password`
*   **Used By:** M3TAL Dashboard container, CLI (for initial setup)

## Network Variables

These variables configure network-related aspects of M3TAL, including ports and domain routing.

### `HTTP_PORT`

*   **Description:** The port on which the M3TAL API daemon (Go application) listens for incoming HTTP requests.
*   **Default Value:** `8080`
*   **Example Value:** `9090`
*   **Used By:** API daemon, CLI

### `NETWORK_NAME`

*   **Description:** The name of the Docker network that M3TAL services will join. This network is typically created by M3TAL to facilitate communication between containers.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal-network`
*   **Used By:** API daemon, M3TAL Dashboard container, Traefik container

### `LOCAL_IP`

*   **Description:** The IP address to which services will bind locally. This is often `127.0.0.1` for local access or a specific host IP if needed.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** API daemon

### `DOMAIN`

*   **Description:** The base domain name for your M3TAL instance. Setting this variable enables Traefik to configure routing rules for `api.${DOMAIN}` and `dash.${DOMAIN}`. If set to `localhost`, Traefik will use `localhost` in its routing rules.
*   **Default Value:** `localhost`
*   **Example Value:** `my-m3tal.example.com`
*   **Used By:** API daemon (for internal communication), Traefik container

### `VPN_USER`

*   **Description:** The username used for establishing a VPN connection, if M3TAL is configured to use one for external access or secure communication.
*   **Default Value:** `user`
*   **Example Value:** `vpn_client_user`
*   **Used By:** API daemon (potentially for VPN clients)

### `VPN_PASSWORD`

*   **Description:** The password used for establishing a VPN connection.
*   **Default Value:** `password`
*   **Example Value:** `your_vpn_password`
*   **Used By:** API daemon (potentially for VPN clients)

## Storage Variables

These variables define the locations for M3TAL's data, including media, configuration, and downloads.

### `BASE_STORAGE_PATH`

*   **Description:** The root directory for all M3TAL data storage. This variable controls where media files, configuration, and downloads are physically located on the host system. **In production deployments, this defaults to `/mnt`, not `./data` as seen in template examples.**
*   **Default Value:** `./data`
*   **Example Value:** `/mnt/m3tal_storage`
*   **Used By:** API daemon, M3TAL Dashboard container

### `MEDIA_PATH`

*   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files (images, videos, etc.) are stored.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/m3tal_storage/media`
*   **Used By:** API daemon, M3TAL Dashboard container

### `CONFIG_PATH`

*   **Description:** The subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/m3tal_storage/config`
*   **Used By:** API daemon, M3TAL Dashboard container

### `DOWNLOADS_PATH`

*   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files are stored.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/m3tal_storage/downloads`
*   **Used By:** API daemon, M3TAL Dashboard container

## Dashboard Variables

These variables specifically control how the M3TAL dashboard is exposed and accessed.

### `DASHBOARD_PORT`

*   **Description:** The internal port on which the M3TAL dashboard container listens for HTTP requests.
*   **Default Value:** `8082`
*   **Example Value:** `8083`
*   **Used By:** M3TAL Dashboard container

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Determines how the M3TAL dashboard is made accessible.
    *   `local`: Exposes the dashboard directly via a port binding (`${DASHBOARD_PORT}:8082`). Access via `http://HOST_IP:${DASHBOARD_PORT}`. This mode does not require Traefik.
    *   `traefik`: Configures Traefik to route traffic for `dash.${DOMAIN}` to the dashboard. Requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** M3TAL Dashboard container (via compose overrides), CLI (for setup)

## Traefik Variables

These variables configure the Traefik reverse proxy for routing external traffic to M3TAL services.

### `TRAEFIK_WEB_PORT`

*   **Description:** The port on the host machine that Traefik will use as its primary HTTP entry point.
*   **Default Value:** `80`
*   **Example Value:** `8000`
*   **Used By:** Traefik container

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The port on the host machine that Traefik will use as its HTTPS entry point. This is relevant if you are configuring TLS/SSL with Traefik.
*   **Default Value:** `443`
*   **Example Value:** `8443`
*   **Used By:** Traefik container

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The internal port on which the Traefik dashboard service listens. This is typically accessed locally for monitoring Traefik's configuration and traffic.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** Traefik container (for its own dashboard)