# M3TAL Environment Variables Reference

All M3TAL environment variables are read from `/etc/m3tal/.env` by both the CLI and all compose stacks via `--env-file`. The `m3tal config wizard` command is the recommended way to manage these variables.

## Quick Reference Table

| Variable Name          | Group    | Default Value        |
|------------------------|----------|----------------------|
| `DASHBOARD_PORT`       | Network  | `8082`               |
| `DASHBOARD_EXPOSE_MODE`| Network  | `local`              |
| `HTTP_PORT`            | Network  | `8080`               |
| `STATE_DIR`            | System   | `./state`            |
| `LOG_LEVEL`            | System   | `info`               |
| `DASHBOARD_SECRET`     | Auth     | `change_me_immediately` |
| `API_TOKEN`            | Auth     | `change_me_api_token` |
| `ADMIN_PASSWORD`       | Auth     | `admin_pass`         |
| `NETWORK_NAME`         | Network  | `m3tal`              |
| `LOCAL_IP`             | Network  | `127.0.0.1`          |
| `DOMAIN`               | Traefik  | `localhost`          |
| `VPN_USER`             | VPN      | `user`               |
| `VPN_PASSWORD`         | VPN      | `password`           |
| `BASE_STORAGE_PATH`    | Storage  | `./data`             |
| `MEDIA_PATH`           | Storage  | `./data/media`       |
| `CONFIG_PATH`          | Storage  | `./data/config`      |
| `DOWNLOADS_PATH`       | Storage  | `./data/downloads`   |
| `PUID`                 | System   | `1000`               |
| `PGID`                 | System   | `1000`               |
| `TZ`                   | System   | `America/Denver`     |
| `TRAEFIK_WEB_PORT`     | Traefik  | `80`                 |
| `TRAEFIK_WEBHTTPS_PORT`| Traefik  | `443`                |
| `TRAEFIK_DASHBOARD_PORT`| Traefik  | `8080`               |
| `DEBUG_MODE`           | System   | `false`              |
| `METRICS_ENABLED`      | System   | `true`               |

---

## Core

These variables control fundamental aspects of the M3TAL system.

### `STATE_DIR`

*   **Description:** Specifies the directory where M3TAL stores its state, including the SQLite database.
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal/state`
*   **Used By:** CLI, API daemon

### `LOG_LEVEL`

*   **Description:** Sets the logging verbosity for M3TAL services.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** CLI, API daemon

### `PUID`

*   **Description:** The user ID for running Docker containers. Defaults to the user ID of the user running the `m3tal` command.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** CLI, Dashboard container

### `PGID`

*   **Description:** The group ID for running Docker containers. Defaults to the group ID of the user running the `m3tal` command.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** CLI, Dashboard container

### `TZ`

*   **Description:** Sets the timezone for containers, ensuring consistent time-based logging and operations.
*   **Default Value:** `America/Denver`
*   **Example Value:** `UTC`
*   **Used By:** Dashboard container

### `DEBUG_MODE`

*   **Description:** Enables or disables debug mode for M3TAL. When `true`, more verbose logging and debugging features may be activated.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** CLI, API daemon

### `METRICS_ENABLED`

*   **Description:** Controls whether Prometheus metrics are exposed by the M3TAL API daemon.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** API daemon

---

## Authentication

These variables manage authentication and security for M3TAL.

### `DASHBOARD_SECRET`

*   **Description:** A secret key used for session management and signing cookies within the M3TAL dashboard. **This is auto-generated on first `m3tal init` and should not be set manually unless rotating.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `s3cr3tK3yF0rD4shb0@rd`
*   **Used By:** Dashboard container

### `API_TOKEN`

*   **Description:** A token used for authenticating API requests. **This is auto-generated on first `m3tal init` and should not be set manually unless rotating.**
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `a1b2c3d4e5f67890g1h2i3j4k5l6m7n8`
*   **Used By:** CLI

### `ADMIN_PASSWORD`

*   **Description:** The password for the default administrator user of the M3TAL dashboard.
*   **Default Value:** `admin_pass`
*   **Example Value:** `new_secure_password`
*   **Used By:** Dashboard container

---

## Network

These variables configure network-related settings for M3TAL services.

### `DASHBOARD_PORT`

*   **Description:** The port on which the M3TAL dashboard is exposed locally. This is used for direct access when `DASHBOARD_EXPOSE_MODE` is `local`.
*   **Default Value:** `8082`
*   **Example Value:** `9090`
*   **Used By:** CLI, Dashboard container

### `HTTP_PORT`

*   **Description:** The port on which the M3TAL API daemon listens for incoming HTTP requests.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** API daemon

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Controls how the M3TAL dashboard is exposed:
    *   `local`: Exposes the dashboard directly via `DASHBOARD_PORT` (default). Access via `http://HOST_IP:8082`.
    *   `traefik`: Exposes the dashboard via Traefik. Access via `http://dash.DOMAIN`.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** CLI

### `NETWORK_NAME`

*   **Description:** The name of the Docker network that M3TAL services will join.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal-network`
*   **Used By:** CLI

### `LOCAL_IP`

*   **Description:** The IP address that M3TAL services will bind to on the host machine.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** CLI

---

## Storage

These variables define the paths for M3TAL's data storage.

### `BASE_STORAGE_PATH`

*   **Description:** The base directory for all M3TAL data storage. In production, this defaults to `/mnt` to allow for external drive mounting. For template/development environments, it defaults to `./data`. This controls where media data is stored.
*   **Default Value:** `./data`
*   **Example Value:** `/mnt`
*   **Used By:** CLI, Dashboard container

### `MEDIA_PATH`

*   **Description:** The specific path within `BASE_STORAGE_PATH` for storing media files.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/m3tal/media`
*   **Used By:** CLI, Dashboard container

### `CONFIG_PATH`

*   **Description:** The specific path within `BASE_STORAGE_PATH` for storing configuration files.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/m3tal/config`
*   **Used By:** CLI, Dashboard container

### `DOWNLOADS_PATH`

*   **Description:** The specific path within `BASE_STORAGE_PATH` for storing downloaded files.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/m3tal/downloads`
*   **Used By:** CLI

---

## Traefik

These variables configure Traefik, the reverse proxy.

### `DOMAIN`

*   **Description:** The primary domain name for your M3TAL instance. Setting this enables Traefik routing rules for `dash.DOMAIN` and `api.DOMAIN`.
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.example.com`
*   **Used By:** CLI, Traefik (via labels)

### `TRAEFIK_WEB_PORT`

*   **Description:** The host port that Traefik uses for HTTP traffic (usually 80).
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** Traefik container

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The host port that Traefik uses for HTTPS traffic (usually 443).
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** Traefik container

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The port on which Traefik's own dashboard is accessible on the host.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** Traefik container

---

## VPN

These variables are used for configuring a VPN connection.

### `VPN_USER`

*   **Description:** The username for connecting to the VPN.
*   **Default Value:** `user`
*   **Example Value:** `myvpnuser`
*   **Used By:** CLI

### `VPN_PASSWORD`

*   **Description:** The password for connecting to the VPN.
*   **Default Value:** `password`
*   **Example Value:** `mysecurevpnpassword`
*   **Used By:** CLI