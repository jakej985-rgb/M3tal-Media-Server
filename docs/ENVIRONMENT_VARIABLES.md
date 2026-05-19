# Environment Variable Reference

This document provides a comprehensive reference for all environment variables used within the M3TAL ecosystem. These variables control the behavior and configuration of the M3TAL CLI, API daemon, Dashboard, Traefik gateway, and other integrated services.

All environment variables are read from the primary configuration file: `/etc/m3tal/.env`. This file is sourced by both the `m3tal` CLI binary and all Docker Compose stacks managed by M3TAL via the `--env-file` flag.

Users can manage these variables using the `m3tal config wizard` for interactive setup or `m3tal config set KEY value` for direct modification. After changing variables, relevant services (e.g., `m3tal-api.service` or Docker Compose stacks via `m3tal up`) may need to be restarted for changes to take effect.

---

## Quick Reference

| Variable Name             | Description                                                                                             | Default Value             | Example Value               | Used By                                  |
| :------------------------ | :------------------------------------------------------------------------------------------------------ | :------------------------ | :-------------------------- | :--------------------------------------- |
| `API_TOKEN`               | Security token for API authentication. **Auto-generated.**                                              | `change_me_api_token`     | `a_long_random_string`      | `m3tal-api.service`, `m3tal-dashboard`   |
| `ADMIN_PASSWORD`          | Initial password for the Dashboard's `admin` user.                                                      | `admin_pass`              | `MySecurePassword123!`      | `m3tal-dashboard`                        |
| `BASE_STORAGE_PATH`       | Base directory for all persistent data. **`/mnt` in production.**                                       | `./data`                  | `/mnt/m3tal`                | `m3tal-dashboard`, User Stacks           |
| `CONFIG_PATH`             | Subdirectory for configuration files.                                                                   | `./data/config`           | `/mnt/config`               | `m3tal-dashboard`, `m3tal-api.service`   |
| `DASHBOARD_EXPOSE_MODE`   | How the Dashboard is exposed: `local` (direct port) or `traefik` (via domain).                          | `local`                   | `traefik`                   | `m3tal-dashboard`                        |
| `DASHBOARD_PORT`          | Port for the M3TAL Dashboard.                                                                           | `8082`                    | `8082`                      | `m3tal-dashboard`                        |
| `DASHBOARD_SECRET`        | Secret key for Dashboard session management. **Auto-generated.**                                        | `change_me_immediately`   | `another_long_random_string`| `m3tal-dashboard`                        |
| `DEBUG_MODE`              | Enables debug features/logging.                                                                         | `false`                   | `true`                      | `m3tal-api.service`, `m3tal-dashboard`   |
| `DOMAIN`                  | Primary domain for Traefik routing. Enables `dash.DOMAIN` and `api.DOMAIN`.                             | `localhost`               | `m3tal.example.com`         | `traefik`, `m3tal-dashboard`, `m3tal-api.service` |
| `DOWNLOADS_PATH`          | Subdirectory for downloaded content.                                                                    | `./data/downloads`        | `/mnt/downloads`            | User Stacks                              |
| `HTTP_PORT`               | Port for the M3TAL API daemon.                                                                          | `8080`                    | `8080`                      | `m3tal-api.service`, `traefik`           |
| `LOCAL_IP`                | Host machine's IP, for internal routing/resolution.                                                     | `127.0.0.1`               | `192.168.1.100`             | `m3tal-api.service`                      |
| `LOG_LEVEL`               | Verbosity for API daemon logs (`info`, `debug`, etc.).                                                  | `info`                    | `debug`                     | `m3tal-api.service`                      |
| `MEDIA_PATH`              | Subdirectory for media files.                                                                           | `./data/media`            | `/mnt/media`                | User Stacks                              |
| `METRICS_ENABLED`         | Enables/disables application metrics.                                                                   | `true`                    | `false`                     | `m3tal-api.service`, `m3tal-dashboard`   |
| `NETWORK_NAME`            | Name of the main Docker network for M3TAL services.                                                     | `m3tal`                   | `m3tal_proxy`               | M3TAL Docker Compose stacks              |
| `PGID`                    | Group ID for containers accessing host volumes.                                                         | `1000`                    | `1000`                      | `m3tal-dashboard`, User Stacks           |
| `PUID`                    | User ID for containers accessing host volumes.                                                          | `1000`                    | `1000`                      | `m3tal-dashboard`, User Stacks           |
| `STATE_DIR`               | Directory for API state DB and Dashboard users.json.                                                    | `./state`                 | `/var/lib/m3tal/state`      | `m3tal-api.service`, `m3tal-dashboard`   |
| `TRAEFIK_DASHBOARD_PORT`  | Internal port for Traefik's own API/Dashboard.                                                          | `8080`                    | `8081`                      | `traefik`                                |
| `TRAEFIK_WEB_PORT`        | Traefik's HTTP entry point port.                                                                        | `80`                      | `80`                        | `traefik`                                |
| `TRAEFIK_WEBHTTPS_PORT`   | Traefik's HTTPS entry point port.                                                                       | `443`                     | `443`                       | `traefik`                                |
| `TZ`                      | Timezone for containers (e.g., `America/Denver`).                                                       | `America/Denver`          | `Europe/Berlin`             | `m3tal-dashboard`, User Stacks           |
| `VPN_PASSWORD`            | Password for VPN client connection.                                                                     | `password`                | `mysecurevpnpass`           | VPN client containers (if deployed)      |
| `VPN_USER`                | Username for VPN client connection.                                                                     | `user`                    | `myvpnuser`                 | VPN client containers (if deployed)      |

---

## Detailed Reference

### Core Variables

#### `HTTP_PORT`
*   **Description:** The port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming HTTP requests.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `m3tal-api.service`, `traefik` (for routing `api.${DOMAIN}` to `host.docker.internal:8080`).

#### `STATE_DIR`
*   **Description:** The directory where the M3TAL API daemon stores its SQLite state database (`state.db`) and where the Dashboard's user credentials (`users.json`) are managed. The path is typically mounted inside the dashboard container as `/docker/state`.
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal/state` (the canonical path for `state.db` in production).
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

#### `NETWORK_NAME`
*   **Description:** The base name of the Docker network created by M3TAL for inter-container communication. M3TAL components and user stacks connect to this network for service discovery. The main network is typically named `proxy` and is external.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal_proxy`
*   **Used By:** M3TAL Docker Compose stacks (for network creation/attachment).

### Authentication Variables

#### `DASHBOARD_SECRET`
*   **Description:** A critical secret key used by the M3TAL Dashboard for session management, encryption, and securing sensitive user data. **This variable is automatically generated on the first `m3tal init` command and should NOT be set manually unless you are intentionally rotating the secret for security purposes.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a_long_random_string_of_characters_for_security`
*   **Used By:** `m3tal-dashboard` container.

#### `API_TOKEN`
*   **Description:** An authentication token required to securely access the M3TAL API daemon. It's used for securing communication between the Dashboard and the API, and for any external tools interacting with the API. **This variable is automatically generated on the first `m3tal init` command and should NOT be set manually unless you are intentionally rotating the token for security purposes.**
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `another_long_random_string_for_api_auth_purposes`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` (when communicating with the API daemon).

#### `ADMIN_PASSWORD`
*   **Description:** The initial password set for the default `admin` user in the M3TAL Dashboard. Users are strongly advised to change this password immediately after their first login for security reasons.
*   **Default Value:** `admin_pass`
*   **Example Value:** `MySecurePassword123!`
*   **Used By:** `m3tal-dashboard` (specifically, during the creation of `users.json` for initial admin setup).

### Network Variables

#### `DASHBOARD_PORT`
*   **Description:** The internal port on which the M3TAL Dashboard container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this port is directly exposed on the host.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** `m3tal-dashboard` container, `m3tal-compose.local.yml` (for host port binding).

#### `DASHBOARD_EXPOSE_MODE`
*   **Description:** Controls how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard port (`DASHBOARD_PORT`) is directly bound to the host, accessible via `http://HOST_IP:DASHBOARD_PORT`. Best for LAN-only setups or initial configuration.
    *   `traefik`: The dashboard is exposed via the Traefik gateway under `dash.${DOMAIN}`. Requires Traefik to be running. Best for domain-based, reverse-proxied setups.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** `m3tal` CLI (to select the appropriate `m3tal-compose` override file when starting the dashboard).

#### `LOCAL_IP`
*   **Description:** The IP address of the host machine. While Docker's `host-gateway` typically resolves this automatically for `host.docker.internal`, this variable can be used by internal API daemon logic or for specific network configurations.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** `m3tal-api.service` (for internal network binding or resolution).

### Storage Variables

#### `BASE_STORAGE_PATH`
*   **Description:** The fundamental directory on the host filesystem where M3TAL stores all its persistent data, including user-defined media, configuration files, and downloads. **In production deployments (e.g., cloud VMs), this defaults to `/mnt`, not `./data` as in development templates.**
*   **Default Value:** `./data`
*   **Example Value:** `/mnt/m3tal_data`
*   **Used By:** `m3tal-dashboard` container (as a volume mount), and implicitly by other user-defined stacks for volume mappings.

#### `MEDIA_PATH`
*   **Description:** The subdirectory, relative to `BASE_STORAGE_PATH`, designated for storing media files and libraries.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/media`
*   **Used By:** User-defined media management stacks that require access to media storage.

#### `CONFIG_PATH`
*   **Description:** The subdirectory, relative to `BASE_STORAGE_PATH`, designated for configuration files. The API daemon's `state.db` and Dashboard's `users.json` are typically stored within `${CONFIG_PATH}/m3tal/state`.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/config`
*   **Used By:** `m3tal-dashboard` container (for volume mounts like `/docker/state`), `m3tal-api.service`.

#### `DOWNLOADS_PATH`
*   **Description:** The subdirectory, relative to `BASE_STORAGE_PATH`, designated for downloaded content.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/downloads`
*   **Used By:** User-defined download client stacks (e.g., torrent clients, newsreaders).

### Traefik Variables

#### `DOMAIN`
*   **Description:** The primary domain name for M3TAL services. When set, Traefik automatically configures routing rules to expose the M3TAL API at `api.DOMAIN` and the Dashboard at `dash.DOMAIN` (if `DASHBOARD_EXPOSE_MODE` is `traefik`). If left as `localhost`, services are typically accessed via direct IP/port.
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.example.com`
*   **Used By:** `traefik` container (for dynamic routing rules), `m3tal-dashboard` container (for Traefik labels in `m3tal-compose.traefik.yml`), `m3tal-api.service` (for dynamic Traefik configuration to expose the API).

#### `TRAEFIK_WEB_PORT`
*   **Description:** The port on which the Traefik gateway listens for incoming unencrypted HTTP traffic. This is the main entry point for web services.
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** `traefik` container.

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description:** The port on which the Traefik gateway listens for incoming encrypted HTTPS traffic. This is the secure entry point for web services (requires additional TLS configuration not detailed here).
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** `traefik` container (for HTTPS entry point configuration).

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description:** The internal port on which Traefik's own API and Dashboard listen. By default, Traefik's dashboard is often exposed on `127.0.0.1:8081` on the host for local access.
*   **Default Value:** `8080`
*   **Example Value:** `8081` (as exposed by default for Traefik's dashboard on the host)
*   **Used By:** `traefik` container.

### VPN Variables

#### `VPN_USER`
*   **Description:** The username used for authenticating with a VPN service, if a VPN client container is deployed within your M3TAL stacks.
*   **Default Value:** `user`
*   **Example Value:** `myvpnuser`
*   **Used By:** VPN client containers (e.g., OpenVPN, WireGuard clients).

#### `VPN_PASSWORD`
*   **Description:** The password used for authenticating with a VPN service, complementing `VPN_USER`, if a VPN client container is deployed.
*   **Default Value:** `password`
*   **Example Value:** `mysecurevpnpass`
*   **Used By:** VPN client containers.

### System Variables

#### `PUID`
*   **Description:** The User ID (UID) that Docker containers should run as when accessing host volumes. Setting this ensures correct file ownership and permissions for data stored on the host, preventing permission issues.
*   **Default Value:** `1000`
*   **Example Value:** `1000` (common for the first user on Linux)
*   **Used By:** `m3tal-dashboard` container, and widely used by other user-defined containers for volume permissions.

#### `PGID`
*   **Description:** The Group ID (GID) that Docker containers should run as when accessing host volumes. Similar to `PUID`, this ensures correct group ownership and permissions for host-mounted data.
*   **Default Value:** `1000`
*   **Example Value:** `1000` (common for the primary group of the first user on Linux)
*   **Used By:** `m3tal-dashboard` container, and widely used by other user-defined containers for volume permissions.

#### `TZ`
*   **Description:** Sets the timezone for Docker containers, ensuring accurate timestamps in logs and for time-sensitive applications. Uses standard TZ database format (e.g., `America/New_York`, `Europe/London`, `Asia/Tokyo`).
*   **Default Value:** `America/Denver`
*   **Example Value:** `Europe/Berlin`
*   **Used By:** `m3tal-dashboard` container, and generally recommended for all user-defined containers.

#### `LOG_LEVEL`
*   **Description:** Controls the verbosity of logging output for the M3TAL API daemon. Higher verbosity levels (e.g., `debug`) provide more detailed information, useful for troubleshooting. Accepted values typically include `debug`, `info`, `warn`, `error`, `fatal`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** `m3tal-api.service`.

#### `DEBUG_MODE`
*   **Description:** A boolean flag that enables debug-specific features or increases logging detail in various M3TAL components. Setting this to `true` can provide additional insights for development or troubleshooting.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` (if it implements a debug mode).

#### `METRICS_ENABLED`
*   **Description:** A boolean flag that controls whether application metrics are collected and exposed by M3TAL components. Disabling this can slightly reduce resource usage if metrics are not required.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` (if it exposes metrics).