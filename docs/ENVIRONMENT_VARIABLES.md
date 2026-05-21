# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, which are read from `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` and can be updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks utilize this `.env` file via `--env-file` for consistent configuration.

Variables are grouped logically for clarity.

## Quick Reference Table

| Variable Name | Default Value | Description |
|---|---|---|
| **Core** | | |
| `DASHBOARD_PORT` | `8082` | Port for the M3TAL Dashboard. |
| `DASHBOARD_EXPOSE_MODE` | `local` | Controls how the dashboard is exposed (`local` or `traefik`). |
| `HTTP_PORT` | `8080` | Internal HTTP port for the M3TAL API. |
| `STATE_DIR` | `./state` | Directory for storing M3TAL state data (e.g., database). |
| `LOG_LEVEL` | `info` | Logging verbosity for M3TAL services. |
| `DEBUG_MODE` | `false` | Enables debug logging and features. |
| `METRICS_ENABLED` | `true` | Enables Prometheus metrics endpoint. |
| **Auth** | | |
| `DASHBOARD_SECRET` | `change_me_immediately` | Secret key for signing dashboard sessions. **Auto-generated on first `m3tal init`.** |
| `API_TOKEN` | `change_me_api_token` | API token for authenticating external API requests. **Auto-generated on first `m3tal init`.** |
| `ADMIN_PASSWORD` | `admin_pass` | Password for the default admin user in the dashboard. |
| **Network** | | |
| `NETWORK_NAME` | `m3tal` | The Docker network name M3TAL services will join. |
| `LOCAL_IP` | `127.0.0.1` | Local IP address for host services. |
| `DOMAIN` | `localhost` | The primary domain for M3TAL services (used by Traefik). |
| `VPN_USER` | `user` | Username for VPN connection. |
| `VPN_PASSWORD` | `password` | Password for VPN connection. |
| **Storage** | | |
| `BASE_STORAGE_PATH` | `./data` | Base path for all M3TAL data storage. **Defaults to `/mnt` in production.** |
| `MEDIA_PATH` | `./data/media` | Path for media files. |
| `CONFIG_PATH` | `./data/config` | Path for configuration files. |
| `DOWNLOADS_PATH` | `./data/downloads` | Path for downloads. |
| `PUID` | `1000` | User ID for file ownership within containers. |
| `PGID` | `1000` | Group ID for file ownership within containers. |
| **System** | | |
| `TZ` | `America/Denver` | Timezone for M3TAL services. |
| **Traefik** | | |
| `TRAEFIK_WEB_PORT` | `80` | Host port for Traefik's HTTP entrypoint. |
| `TRAEFIK_WEBHTTPS_PORT` | `443` | Host port for Traefik's HTTPS entrypoint. |
| `TRAEFIK_DASHBOARD_PORT` | `8080` | Host port for Traefik's own dashboard. |

---

## Core

These variables configure the fundamental behavior of the M3TAL system.

### `DASHBOARD_PORT`

*   **Description:** The port on which the M3TAL Dashboard service will listen.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** `m3tal-dashboard` container, `m3tal` CLI.

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Controls how the M3TAL Dashboard is exposed to the network.
    *   `local`: Exposes the dashboard directly via a port binding. Access is typically `http://HOST_IP:DASHBOARD_PORT`. This mode does not require Traefik.
    *   `traefik`: Configures Traefik to route traffic to the dashboard via `dash.DOMAIN`. This mode requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** `m3tal-dashboard` container, `m3tal` CLI.

### `HTTP_PORT`

*   **Description:** The internal HTTP port that the M3TAL API daemon (Go binary) listens on. This is generally not exposed directly to the host unless `LOCAL_IP` is `0.0.0.0`.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `m3tal-api.service`.

### `STATE_DIR`

*   **Description:** The directory where M3TAL stores its persistent state, including the SQLite database.
*   **Default Value:** `./state` (This will be relative to the directory where the Compose files are executed, typically `/opt/m3tal/stack/`).
*   **Example Value:** `/var/lib/m3tal/state`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `LOG_LEVEL`

*   **Description:** Sets the logging verbosity for M3TAL services. Common values include `debug`, `info`, `warn`, `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `DEBUG_MODE`

*   **Description:** Enables additional debugging features and logging within M3TAL services.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `METRICS_ENABLED`

*   **Description:** Toggles the availability of the Prometheus metrics endpoint for M3TAL services.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** `m3tal-api.service`.

---

## Auth

These variables are crucial for securing your M3TAL instance.

### `DASHBOARD_SECRET`

*   **Description:** A secret key used for signing session cookies for the M3TAL Dashboard. This is essential for security.
*   **Default Value:** `change_me_immediately`
*   **Important:** This variable is **auto-generated on the first `m3tal init`** command. Users should **not** set this manually unless rotating the secret.
*   **Example Value:** (A long, randomly generated string)
*   **Used By:** `m3tal-dashboard` container.

### `API_TOKEN`

*   **Description:** A token used for authenticating external API requests to the M3TAL API.
*   **Default Value:** `change_me_api_token`
*   **Important:** This variable is **auto-generated on the first `m3tal init`** command. Users should **not** set this manually unless rotating the token.
*   **Example Value:** (A long, randomly generated string)
*   **Used By:** `m3tal-api.service` (for internal validation).

### `ADMIN_PASSWORD`

*   **Description:** The password for the default administrative user in the M3TAL Dashboard.
*   **Default Value:** `admin_pass`
*   **Example Value:** `mySecurePassword123`
*   **Used By:** `m3tal-dashboard` container.

---

## Network

These variables define network-related configurations for M3TAL.

### `NETWORK_NAME`

*   **Description:** The name of the Docker network that M3TAL services will be attached to. This ensures proper inter-service communication.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal_network`
*   **Used By:** All M3TAL Docker Compose stacks.

### `LOCAL_IP`

*   **Description:** The IP address that M3TAL services will bind to on the host. `127.0.0.1` (localhost) means services are only accessible from the host machine. `0.0.0.0` means services are accessible from any network interface.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `0.0.0.0`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `DOMAIN`

*   **Description:** The primary domain name for your M3TAL instance. This is used by Traefik to configure routing rules.
    *   Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes when `DASHBOARD_EXPOSE_MODE=traefik`.
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.yourdomain.com`
*   **Used By:** Traefik configuration, `m3tal-dashboard` container (for Traefik labels).

### `VPN_USER`

*   **Description:** The username for your VPN connection, if M3TAL is configured to use a VPN.
*   **Default Value:** `user`
*   **Example Value:** `vpnuser123`
*   **Used By:** Potentially used by specific M3TAL services or integrations that require VPN access.

### `VPN_PASSWORD`

*   **Description:** The password for your VPN connection, if M3TAL is configured to use a VPN.
*   **Default Value:** `password`
*   **Example Value:** `VpnPass!@#`
*   **Used By:** Potentially used by specific M3TAL services or integrations that require VPN access.

---

## Storage

These variables define the location and ownership of M3TAL's persistent data.

### `BASE_STORAGE_PATH`

*   **Description:** The base directory where all M3TAL data (configuration, media, downloads, state) will be stored.
    *   **Note:** In production deployments, this defaults to `/mnt`, not `./data` as seen in template configurations.
*   **Default Value:** `./data`
*   **Example Value:** `/mnt`
*   **Used By:** `m3tal-dashboard` container, `m3tal-api.service` (via volume mounts).

### `MEDIA_PATH`

*   **Description:** The specific path within `BASE_STORAGE_PATH` for storing media files.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/media`
*   **Used By:** `m3tal-dashboard` container (via volume mounts).

### `CONFIG_PATH`

*   **Description:** The specific path within `BASE_STORAGE_PATH` for storing configuration files.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/config`
*   **Used By:** `m3tal-dashboard` container (via volume mounts).

### `DOWNLOADS_PATH`

*   **Description:** The specific path within `BASE_STORAGE_PATH` for storing downloaded files.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/downloads`
*   **Used By:** `m3tal-dashboard` container (via volume mounts).

### `PUID`

*   **Description:** The User ID (UID) that Docker containers will run as for file ownership and permissions. This helps ensure containers can read/write to mounted volumes correctly.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** `m3tal-dashboard` container.

### `PGID`

*   **Description:** The Group ID (GID) that Docker containers will run as for file ownership and permissions. This helps ensure containers can read/write to mounted volumes correctly.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** `m3tal-dashboard` container.

---

## System

These variables control system-level configurations.

### `TZ`

*   **Description:** The timezone to be used by M3TAL services. This ensures correct time logging and scheduling.
*   **Default Value:** `America/Denver`
*   **Example Value:** `UTC`
*   **Used By:** `m3tal-dashboard` container.

---

## Traefik

These variables configure the Traefik reverse proxy.

### `TRAEFIK_WEB_PORT`

*   **Description:** The host port that Traefik will listen on for incoming HTTP traffic. This is the main entry point for services exposed via Traefik.
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** Traefik container.

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The host port that Traefik will listen on for incoming HTTPS traffic. This is typically used when TLS is configured.
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** Traefik container.

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The host port that Traefik's own administrative dashboard will be accessible on. This is usually accessed locally via `http://localhost:8081`.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** Traefik container.