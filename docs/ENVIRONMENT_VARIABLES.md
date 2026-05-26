# Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL Ecosystem. These variables control the behavior and configuration of the M3TAL CLI, the M3TAL API daemon, the M3TAL Dashboard, Traefik, and other managed Docker services.

All M3TAL environment variables are read from the primary configuration file located at `/etc/m3tal/.env`. Both the `m3tal` CLI and all Docker Compose stacks managed by M3TAL load their configurations from this file via the `--env-file` option, ensuring a single source of truth for your M3TAL deployment. You can manage these variables using `m3tal config wizard` or `m3tal config set KEY value`.

## Quick Reference

| Variable Name             | Default                  | Description                                                                 |
| :------------------------ | :----------------------- | :-------------------------------------------------------------------------- |
| `HTTP_PORT`               | `8080`                   | Port for the M3TAL API daemon.                                              |
| `STATE_DIR`               | `./state`                | Path for the M3TAL API's state database.                                    |
| `LOG_LEVEL`               | `info`                   | Minimum logging level for components.                                       |
| `DEBUG_MODE`              | `false`                  | Enables debug features and logging.                                         |
| `METRICS_ENABLED`         | `true`                   | Controls M3TAL API Prometheus metrics exposure.                             |
| `DASHBOARD_PORT`          | `8082`                   | Port for the M3TAL Dashboard.                                               |
| `DASHBOARD_EXPOSE_MODE`   | `local`                  | How the Dashboard is exposed (`local` or `traefik`).                        |
| `DASHBOARD_SECRET`        | `change_me_immediately`  | Secret key for Dashboard session management. **Auto-generated.**            |
| `ADMIN_PASSWORD`          | `admin_pass`             | Default password for the initial admin user.                                |
| `API_TOKEN`               | `change_me_api_token`    | Token for Dashboard-API communication. **Auto-generated.**                  |
| `NETWORK_NAME`            | `m3tal`                  | Name of the Docker network for M3TAL services.                              |
| `LOCAL_IP`                | `127.0.0.1`              | Host machine's local IP for `host.docker.internal` resolution.              |
| `DOMAIN`                  | `localhost`              | Base domain name for Traefik routing.                                       |
| `TRAEFIK_WEB_PORT`        | `80`                     | Host port for Traefik's HTTP entry point.                                   |
| `TRAEFIK_WEBHTTPS_PORT`   | `443`                    | Host port for Traefik's HTTPS entry point.                                  |
| `TRAEFIK_DASHBOARD_PORT`  | `8080`                   | Internal port for Traefik's management dashboard.                           |
| `BASE_STORAGE_PATH`       | `./data`                 | Base directory for all M3TAL data.                                          |
| `MEDIA_PATH`              | `./data/media`           | Subdirectory for user media files.                                          |
| `CONFIG_PATH`             | `./data/config`          | Subdirectory for persistent configuration.                                  |
| `DOWNLOADS_PATH`          | `./data/downloads`       | Subdirectory for downloaded files.                                          |
| `VPN_USER`                | `user`                   | Username for VPN access.                                                    |
| `VPN_PASSWORD`            | `password`               | Password for VPN access.                                                    |
| `PUID`                    | `1000`                   | User ID for containers.                                                     |
| `PGID`                    | `1000`                   | Group ID for containers.                                                    |
| `TZ`                      | `America/Denver`         | Timezone for containers.                                                    |

---

## Detailed Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL ecosystem's logging and API daemon behavior.

#### `HTTP_PORT`
*   **Description:** The port on which the M3TAL API daemon (Go binary) listens for incoming requests. This is the primary communication endpoint for the M3TAL Dashboard and other internal services.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `m3tal-api` daemon, `m3tal-dashboard` container (for API URL).

#### `STATE_DIR`
*   **Description:** Specifies the path to the directory where the M3TAL API daemon stores its SQLite database (`state.db`) and other runtime state information. This directory is mounted into the `m3tal-dashboard` container as `/docker/state`. On the host, this typically resolves to `/var/lib/m3tal/state.db`.
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal/state` (the host path)
*   **Used By:** `m3tal-api` daemon, `m3tal-dashboard` container.

#### `LOG_LEVEL`
*   **Description:** Sets the minimum logging level for M3TAL components, influencing the verbosity of output to logs.
*   **Default Value:** `info`
*   **Example Value:** `debug`, `warn`, `error`, `fatal`
*   **Used By:** `m3tal-api` daemon, `m3tal-dashboard` container.

#### `DEBUG_MODE`
*   **Description:** A boolean flag that, when set to `true`, enables debug-level logging and potentially other debug features within M3TAL components, aiding in troubleshooting.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** `m3tal-api` daemon, `m3tal-dashboard` container.

#### `METRICS_ENABLED`
*   **Description:** Controls whether Prometheus-compatible metrics are exposed by the M3TAL API daemon, allowing for monitoring of the system's performance and health.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** `m3tal-api` daemon.

### Authentication & Security

These variables are crucial for securing access to the M3TAL Dashboard and API.

#### `DASHBOARD_SECRET`
*   **Description:** A secret key used by the M3TAL Dashboard for secure session management and encrypting internal communications. **This value is automatically generated on the first `m3tal init` run.** Users should NOT set this manually unless performing a key rotation. Keep this value confidential.
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a_long_random_string_of_characters_for_security`
*   **Used By:** `m3tal-dashboard` container.

#### `API_TOKEN`
*   **Description:** An authentication token used by the M3TAL Dashboard to communicate securely with the M3TAL API daemon. **This value is automatically generated on the first `m3tal init` run.** Users should NOT set this manually unless performing a key rotation. Keep this value confidential.
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `another_long_random_string_for_api_authentication`
*   **Used By:** `m3tal-api` daemon, `m3tal-dashboard` container.

#### `ADMIN_PASSWORD`
*   **Description:** The default password for the initial `admin` user account created within the M3TAL Dashboard. **It is highly recommended to change this immediately after setup using the `m3tal dashpass` command.**
*   **Default Value:** `admin_pass`
*   **Example Value:** `myStrongAndSecurePassword123`
*   **Used By:** `m3tal-dashboard` container (for initial user setup).

### M3TAL Dashboard

Variables specific to the M3TAL Dashboard container and its exposure.

#### `DASHBOARD_PORT`
*   **Description:** The internal port on which the M3TAL Dashboard container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this is also the host port where the dashboard is exposed.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** `m3tal-dashboard` container, CLI (`m3tal dash up`).

#### `DASHBOARD_EXPOSE_MODE`
*   **Description:** Controls how the M3TAL Dashboard is exposed to the network.
    *   `local` (default): The dashboard is directly exposed on `http://HOST_IP:DASHBOARD_PORT`. Best for LAN-only setups or initial access.
    *   `traefik`: The dashboard is exposed via Traefik, accessible at `http://dash.DOMAIN`. Requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** CLI (`m3tal dash up`), Docker Compose.

### Network & Routing

Configuration for Docker networks and general routing within the M3TAL ecosystem.

#### `NETWORK_NAME`
*   **Description:** The name of the Docker network created and used by all M3TAL-managed containers for internal communication, allowing them to discover and interact with each other.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal_custom_network`
*   **Used By:** All Docker Compose stacks.

#### `LOCAL_IP`
*   **Description:** The local IP address of the host machine. This is used by containers to resolve `host.docker.internal` to the actual host IP, enabling communication with host-bound services like the M3TAL API daemon.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** Traefik gateway, `m3tal-dashboard` container.

#### `DOMAIN`
*   **Description:** The base domain name that M3TAL uses for services exposed via the Traefik gateway. Setting this variable enables Traefik to create routing rules like `dash.DOMAIN` for the dashboard and `api.DOMAIN` for the API daemon. If not set, Traefik routes will typically default to `localhost`.
*   **Default Value:** `localhost`
*   **Example Value:** `example.com`
*   **Used By:** Traefik gateway, CLI (`m3tal dash up`).

### Traefik Gateway

Variables specifically for configuring the Traefik reverse proxy.

#### `TRAEFIK_WEB_PORT`
*   **Description:** The host port on which Traefik's HTTP (non-HTTPS) entry point listens. This is typically port 80 for standard web traffic.
*   **Default Value:** `80`
*   **Example Value:** `8080` (if port 80 is already in use)
*   **Used By:** `traefik` container.

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description:** The host port on which Traefik's HTTPS entry point listens. This is typically port 443 for secure web traffic.
*   **Default Value:** `443`
*   **Example Value:** `8443` (if port 443 is already in use)
*   **Used By:** `traefik` container.

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description:** The internal port within the Traefik container that hosts Traefik's own management dashboard. This port is *not* exposed directly to the internet by default; instead, it's typically mapped to `127.0.0.1:8081` on the host for local access.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `traefik` container (internal configuration).

### Storage & Paths

Defines the filesystem locations for various M3TAL data.

#### `BASE_STORAGE_PATH`
*   **Description:** The absolute path on the host system that serves as the root directory for all M3TAL-related persistent data, including media, configurations, and downloads. **In production M3TAL deployments, this defaults to `/mnt` to leverage large attached storage, rather than `./data` as seen in template files.**
*   **Default Value:** `./data` (for template/development)
*   **Example Value:** `/mnt/m3tal_data`
*   **Used By:** All containers requiring persistent storage.

#### `MEDIA_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` designated for storing user media files, such as movies, TV shows, and music. This path is often mounted into media-serving containers (e.g., Plex, Jellyfin).
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/media`
*   **Used By:** Media-serving containers (if deployed).

#### `CONFIG_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` used to store persistent configuration files for various M3TAL-managed services and applications.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/config`
*   **Used By:** All containers requiring persistent configuration.

#### `DOWNLOADS_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` where downloaded files from torrent clients or other download services are stored.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/downloads`
*   **Used By:** Download client containers (if deployed).

### VPN Services

Variables for configuring any VPN services deployed via M3TAL.

#### `VPN_USER`
*   **Description:** The username credential for accessing VPN services deployed through M3TAL (e.g., a WireGuard container).
*   **Default Value:** `user`
*   **Example Value:** `m3tal_vpn_user`
*   **Used By:** VPN client configurations (if a VPN service is deployed).

#### `VPN_PASSWORD`
*   **Description:** The password credential for accessing VPN services deployed through M3TAL.
*   **Default Value:** `password`
*   **Example Value:** `vpn_super_secure_password`
*   **Used By:** VPN client configurations (if a VPN service is deployed).

### System & Permissions

General system-level configurations, particularly for user and group IDs within containers.

#### `PUID`
*   **Description:** The User ID (UID) that containers will run as internally. Setting this ensures that containers have appropriate file permissions when interacting with host-mounted volumes, preventing permission issues.
*   **Default Value:** `1000`
*   **Example Value:** `1001` (matching a specific host user)
*   **Used By:** All M3TAL containers (passed as the effective user ID).

#### `PGID`
*   **Description:** The Group ID (GID) that containers will run as internally. Similar to `PUID`, this ensures proper file permissions for containers accessing host volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1001` (matching a specific host group)
*   **Used By:** All M3TAL containers (passed as the effective group ID).

#### `TZ`
*   **Description:** Sets the timezone for all M3TAL containers. This ensures consistent time synchronization across services and accurate timestamping in logs.
*   **Default Value:** `America/Denver`
*   **Example Value:** `Europe/London`, `America/New_York`
*   **Used By:** All M3TAL containers.