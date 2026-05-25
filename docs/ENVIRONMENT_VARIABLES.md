# M3TAL Environment Variables Reference

This document details all environment variables used by the M3TAL ecosystem. These variables are crucial for configuring M3TAL's behavior, including its networking, authentication, storage, and overall system operation.

All environment variables are read from the `/etc/m3tal/.env` file by both the M3TAL CLI and all Docker Compose stacks. This centralizes configuration and ensures consistency across the system.

**Quick Reference Table:**

| Variable Name              | Group    | Default Value       |
|----------------------------|----------|---------------------|
| `DASHBOARD_PORT`           | Network  | `8082`              |
| `DASHBOARD_EXPOSE_MODE`    | Network  | `local`             |
| `HTTP_PORT`                | Network  | `8080`              |
| `STATE_DIR`                | System   | `./state`           |
| `LOG_LEVEL`                | System   | `info`              |
| `DASHBOARD_SECRET`         | Auth     | `change_me_immediately` |
| `API_TOKEN`                | Auth     | `change_me_api_token` |
| `ADMIN_PASSWORD`           | Auth     | `admin_pass`        |
| `NETWORK_NAME`             | Network  | `m3tal`             |
| `LOCAL_IP`                 | Network  | `127.0.0.1`         |
| `DOMAIN`                   | Traefik  | `localhost`         |
| `VPN_USER`                 | VPN      | `user`              |
| `VPN_PASSWORD`             | VPN      | `password`          |
| `BASE_STORAGE_PATH`        | Storage  | `./data`            |
| `MEDIA_PATH`               | Storage  | `./data/media`      |
| `CONFIG_PATH`              | Storage  | `./data/config`     |
| `DOWNLOADS_PATH`           | Storage  | `./data/downloads`  |
| `PUID`                     | System   | `1000`              |
| `PGID`                     | System   | `1000`              |
| `TZ`                       | System   | `America/Denver`    |
| `TRAEFIK_WEB_PORT`         | Traefik  | `80`                |
| `TRAEFIK_WEBHTTPS_PORT`    | Traefik  | `443`               |
| `TRAEFIK_DASHBOARD_PORT`   | Traefik  | `8080`              |
| `DEBUG_MODE`               | System   | `false`             |
| `METRICS_ENABLED`          | System   | `true`              |

---

## Core Configuration

### `STATE_DIR`

*   **Description:** Specifies the directory where M3TAL stores its state, including the SQLite database.
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal/state`
*   **Used By:** API daemon, CLI

### `LOG_LEVEL`

*   **Description:** Controls the verbosity of M3TAL's logging output. Options typically include `debug`, `info`, `warn`, `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** API daemon, CLI

### `PUID`

*   **Description:** The User ID (UID) to run Docker containers as. This is important for file permissions.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** M3TAL Dashboard container, CLI

### `PGID`

*   **Description:** The Group ID (GID) to run Docker containers as. This is important for file permissions.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** M3TAL Dashboard container, CLI

### `TZ`

*   **Description:** The timezone to use for logging and any time-sensitive operations within M3TAL services.
*   **Default Value:** `America/Denver`
*   **Example Value:** `UTC`
*   **Used By:** M3TAL Dashboard container

### `DEBUG_MODE`

*   **Description:** Enables or disables debug mode for M3TAL.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** CLI

### `METRICS_ENABLED`

*   **Description:** Enables or disables metrics collection within M3TAL.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** API daemon

---

## Authentication

### `DASHBOARD_SECRET`

*   **Description:** A secret key used for session management and securing the M3TAL Dashboard. **This is auto-generated on the first `m3tal init` and should generally not be set manually unless rotating.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a_very_long_and_random_secret_key`
*   **Used By:** M3TAL Dashboard container

### `API_TOKEN`

*   **Description:** An API token for authenticating requests to the M3TAL API. **This is auto-generated on the first `m3tal init` and should generally not be set manually unless rotating.**
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `a_secure_api_token_string`
*   **Used By:** CLI, M3TAL Dashboard container

### `ADMIN_PASSWORD`

*   **Description:** The password for the default administrator user of the M3TAL Dashboard.
*   **Default Value:** `admin_pass`
*   **Example Value:** `a_strong_and_unique_password`
*   **Used By:** M3TAL Dashboard container

---

## Network Configuration

### `DASHBOARD_PORT`

*   **Description:** The port on which the M3TAL Dashboard will be accessible. In `local` `DASHBOARD_EXPOSE_MODE`, this is the direct port binding.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** M3TAL Dashboard container, CLI

### `HTTP_PORT`

*   **Description:** The port on which the M3TAL API daemon listens for incoming requests.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** API daemon

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Determines how the M3TAL Dashboard is exposed.
    *   `local`: Exposes the dashboard directly via `DASHBOARD_PORT` (e.g., `http://HOST_IP:8082`). No Traefik is required. Ideal for LAN-only setups.
    *   `traefik`: Exposes the dashboard through Traefik, typically at `http://dash.DOMAIN`. Requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** CLI (for configuring the dashboard's compose file)

### `NETWORK_NAME`

*   **Description:** The name of the Docker network M3TAL services will use to communicate.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal-network`
*   **Used By:** CLI, Docker Compose files

### `LOCAL_IP`

*   **Description:** The IP address to bind network services to locally.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `127.0.0.1`
*   **Used By:** CLI

---

## Storage Configuration

### `BASE_STORAGE_PATH`

*   **Description:** The root directory for all M3TAL persistent data, including media, configuration, and downloads. **Note:** In production deployments, this defaults to `/mnt`. In template/development environments, it may default to `./data`.
*   **Default Value:** `./data`
*   **Example Value:** `/mnt`
*   **Used By:** M3TAL Dashboard container, CLI

### `MEDIA_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` for storing media files.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/m3tal/media`
*   **Used By:** M3TAL Dashboard container, CLI

### `CONFIG_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` for storing M3TAL configuration files.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/m3tal/config`
*   **Used By:** M3TAL Dashboard container, CLI

### `DOWNLOADS_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` for storing downloaded files.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/m3tal/downloads`
*   **Used By:** M3TAL Dashboard container, CLI

---

## Traefik Configuration

### `DOMAIN`

*   **Description:** The primary domain name used for routing external traffic to M3TAL services via Traefik. Setting this enables routes like `dash.DOMAIN` and `api.DOMAIN`.
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.example.com`
*   **Used By:** Traefik (dynamic configuration), CLI

### `TRAEFIK_WEB_PORT`

*   **Description:** The host port that Traefik will use for HTTP (port 80) traffic.
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** Traefik container

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The host port that Traefik will use for HTTPS (port 443) traffic.
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** Traefik container

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The internal port on which the Traefik dashboard is accessible (e.g., `127.0.0.1:8081`).
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** Traefik container

---

## VPN Configuration

### `VPN_USER`

*   **Description:** The username for your VPN connection.
*   **Default Value:** `user`
*   **Example Value:** `my_vpn_username`
*   **Used By:** VPN client within M3TAL (if configured)

### `VPN_PASSWORD`

*   **Description:** The password for your VPN connection.
*   **Default Value:** `password`
*   **Example Value:** `my_vpn_password`
*   **Used By:** VPN client within M3TAL (if configured)