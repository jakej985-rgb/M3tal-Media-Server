# M3TAL Environment Variables Reference

As the M3TAL Ecosystem Documentation Architect, my goal is to provide a comprehensive reference for all environment variables that configure your M3TAL system. These variables control everything from network ports to storage paths and authentication secrets.

All M3TAL environment variables are read from the primary configuration file: `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (when invoked via `m3tal up` or `m3tal dash up`) utilize this file by passing it via the `--env-file` argument. This ensures a consistent configuration across your entire M3TAL deployment.

---

## Quick Reference Table

| Name                   | Description                                                                                             | Default Value            |
| :--------------------- | :------------------------------------------------------------------------------------------------------ | :----------------------- |
| `DASHBOARD_PORT`       | Port for the M3TAL Dashboard container.                                                                 | `8082`                   |
| `DASHBOARD_EXPOSE_MODE`| Controls how the M3TAL Dashboard is exposed (direct port or Traefik).                                   | `local`                  |
| `HTTP_PORT`            | Port on which the M3TAL API daemon listens.                                                             | `8080`                   |
| `STATE_DIR`            | Directory where the M3TAL API daemon stores its SQLite database and other runtime state.                | `./state`                |
| `DASHBOARD_SECRET`     | Secret key for M3TAL Dashboard session management. **Auto-generated.**                                  | `change_me_immediately`  |
| `API_TOKEN`            | Authentication token for M3TAL API daemon access. **Auto-generated.**                                   | `change_me_api_token`    |
| `ADMIN_PASSWORD`       | Default password for the initial M3TAL Dashboard admin user.                                            | `admin_pass`             |
| `NETWORK_NAME`         | Name of the Docker network used by M3TAL components.                                                    | `m3tal`                  |
| `LOCAL_IP`             | IP address representing the host machine from within Docker containers.                                 | `127.0.0.1`              |
| `DOMAIN`               | Primary domain name for M3TAL services when using Traefik.                                              | `localhost`              |
| `VPN_USER`             | Username for the VPN service.                                                                           | `user`                   |
| `VPN_PASSWORD`         | Password for the VPN service.                                                                           | `password`               |
| `BASE_STORAGE_PATH`    | Base directory for all persistent data storage.                                                         | `./data`                 |
| `MEDIA_PATH`           | Directory where media data is stored.                                                                   | `./data/media`           |
| `CONFIG_PATH`          | Directory for persistent configuration files.                                                           | `./data/config`          |
| `DOWNLOADS_PATH`       | Directory for storing downloaded files.                                                                 | `./data/downloads`       |
| `PUID`                 | User ID (UID) for Docker containers.                                                                    | `1000`                   |
| `PGID`                 | Group ID (GID) for Docker containers.                                                                   | `1000`                   |
| `TZ`                   | Timezone for all M3TAL components.                                                                      | `America/Denver`         |
| `TRAEFIK_WEB_PORT`     | Host port for Traefik's primary HTTP entry point.                                                       | `80`                     |
| `TRAEFIK_WEBHTTPS_PORT`| Host port for Traefik's HTTPS entry point.                                                              | `443`                    |
| `TRAEFIK_DASHBOARD_PORT`| Internal port for the Traefik management dashboard.                                                     | `8080`                   |
| `LOG_LEVEL`            | Minimum logging level for M3TAL components.                                                             | `info`                   |
| `DEBUG_MODE`           | Enables debug-level logging and features.                                                               | `false`                  |
| `METRICS_ENABLED`      | Enables or disables the collection and exposure of system metrics.                                      | `true`                   |

---

## Detailed Environment Variable Reference

### Core M3TAL Configuration

These variables control fundamental aspects of the M3TAL API daemon and Dashboard operation.

#### `DASHBOARD_PORT`

*   **Description**: Specifies the internal port on which the `m3tal-dashboard` container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this port is directly exposed on the host.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Used By**: `m3tal-dashboard` container, `m3tal-compose.local.yml` (for port binding), `m3tal-compose.traefik.yml` (for Traefik service definition).

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Determines how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard is exposed directly on the host's `DASHBOARD_PORT`. Best for LAN-only setups or initial configuration. Access via `http://HOST_IP:8082`.
    *   `traefik`: The dashboard is routed through Traefik using the `dash.${DOMAIN}` hostname. Requires Traefik to be running and `DOMAIN` to be configured. Access via `http://dash.YOUR_DOMAIN`.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Used By**: `m3tal dash up` command (to select the appropriate Docker Compose override file: `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`).

#### `HTTP_PORT`

*   **Description**: The network port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming requests. This port is generally only accessible locally on the host or via Traefik.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: `m3tal-api.service`, Traefik (`dynamic/api.yml` for routing to the API).

#### `STATE_DIR`

*   **Description**: Defines the directory where the M3TAL API daemon stores its persistent data, including the core SQLite state database (`state.db`).
*   **Default Value**: `./state`
*   **Example Value**: `/var/lib/m3tal/state`
*   **Used By**: `m3tal-api.service`, `m3tal-dashboard` (to mount the state directory).

### Authentication & Security

These variables manage authentication credentials and secrets for M3TAL components.

#### `DASHBOARD_SECRET`

*   **Description**: A cryptographic secret key used by the M3TAL Dashboard for secure session management, cookie signing, and other security-sensitive operations. **This variable is automatically generated on the first `m3tal init` run.**
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_very_long_and_random_string_of_characters_123abc`
*   **Used By**: `m3tal-dashboard` container.
*   **IMPORTANT**: Users should **NOT** set this manually unless performing a secret rotation and understand the implications (e.g., invalidating existing user sessions). Always use `m3tal config wizard` or `m3tal init` for initial setup.

#### `API_TOKEN`

*   **Description**: The bearer token required to authenticate requests to the M3TAL API daemon. The dashboard uses this token to communicate with the API. **This variable is automatically generated on the first `m3tal init` run.**
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `another_very_long_and_random_string_of_characters_xyz789`
*   **Used By**: `m3tal-api.service` (for validation), `m3tal-dashboard` (for API calls).
*   **IMPORTANT**: Users should **NOT** set this manually unless performing a token rotation. Always use `m3tal config wizard` or `m3tal init` for initial setup.

#### `ADMIN_PASSWORD`

*   **Description**: The default password for the initial administrator user (`admin`) of the M3TAL Dashboard. It is highly recommended to change this immediately after first access via the dashboard's user management or the `m3tal dashpass` CLI command.
*   **Default Value**: `admin_pass`
*   **Example Value**: `MySuperSecurePassword123!`
*   **Used By**: `m3tal-dashboard` (for initial user setup), `m3tal dashpass` (CLI for managing dashboard user passwords).

### Network Configuration

These variables control Docker networking and local IP routing.

#### `NETWORK_NAME`

*   **Description**: Defines the name of the custom Docker network used by all M3TAL-managed containers. This network facilitates secure and efficient inter-container communication.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal_internal_net`
*   **Used By**: All M3TAL Docker Compose stacks (e.g., `routing-compose.yml`, `m3tal-compose.yml`) via their `networks` section.

#### `LOCAL_IP`

*   **Description**: Specifies the IP address that Docker containers use to reach the host machine. This is crucial for services running directly on the host (like the M3TAL API daemon) to be accessible from within containers. Typically configured as `host-gateway` or a specific LAN IP.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100` (if your host has a static LAN IP)
*   **Used By**: `m3tal-dashboard` (`extra_hosts` for `host.docker.internal`), Traefik (`dynamic/api.yml` for routing to `http://host.docker.internal:8080`).

### Storage Paths

These variables define where M3TAL stores various types of persistent data on the host filesystem.

#### `BASE_STORAGE_PATH`

*   **Description**: The fundamental root directory on the host filesystem for all persistent M3TAL data. Other `_PATH` variables are often relative to this. **In production deployments, this defaults to `/mnt` (e.g., for disk mounts), whereas the template defaults to `./data` for local testing.**
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal_data`
*   **Used By**: All Docker Compose stacks as the base for volume mounts.

#### `MEDIA_PATH`

*   **Description**: The subdirectory, relative to `BASE_STORAGE_PATH` (or an absolute path if specified), where large media files (e.g., videos, music, photos) should be stored.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/media`
*   **Used By**: Docker Compose stacks for media-related services (e.g., Plex, Jellyfin).

#### `CONFIG_PATH`

*   **Description**: The subdirectory, relative to `BASE_STORAGE_PATH` (or an absolute path if specified), where persistent configuration files (like the Dashboard's `users.json`) are stored.
*   **Default Value**: `./data/config`
*   **Example Value**: `/mnt/m3tal_config`
*   **Used By**: `m3tal-dashboard` (for `users.json`), potentially other services requiring persistent config.

#### `DOWNLOADS_PATH`

*   **Description**: The subdirectory, relative to `BASE_STORAGE_PATH` (or an absolute path if specified), designated for storing downloaded files.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/downloads`
*   **Used By**: Docker Compose stacks for download management services (e.g., qBittorrent, Transmission).

### Traefik Gateway

These variables configure the Traefik reverse proxy and its interaction with M3TAL services.

#### `DOMAIN`

*   **Description**: The primary domain name under which your M3TAL services (Dashboard, API, etc.) will be accessible when using Traefik. **Setting this enables Traefik routing rules for `dash.DOMAIN` and `api.DOMAIN` routes.**
*   **Default Value**: `localhost`
*   **Example Value**: `my-server.example.com`
*   **Used By**: Traefik (`routing-compose.yml`, `dynamic/api.yml`, `m3tal-compose.traefik.yml`), Cloudflared.

#### `TRAEFIK_WEB_PORT`

*   **Description**: The host port Traefik binds to for handling incoming HTTP (non-HTTPS) traffic. This is Traefik's primary entry point for web services.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Used By**: `traefik` container (`routing-compose.yml`).

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port Traefik binds to for handling incoming HTTPS traffic. Requires additional Traefik configuration for SSL/TLS certificates.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Used By**: `traefik` container (`routing-compose.yml`).

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port of the Traefik container where its own management dashboard is exposed. By default, Traefik's dashboard is mapped to `127.0.0.1:8081` on the host, making it accessible only locally.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: `traefik` container (`routing-compose.yml`).

### VPN Configuration

These variables are for configuring an integrated VPN service, if used.

#### `VPN_USER`

*   **Description**: The username for authenticating with an integrated VPN service (e.g., OpenVPN, WireGuard).
*   **Default Value**: `user`
*   **Example Value**: `m3tal_vpn_user`
*   **Used By**: VPN Stack (if deployed).

#### `VPN_PASSWORD`

*   **Description**: The password for authenticating with an integrated VPN service.
*   **Default Value**: `password`
*   **Example Value**: `SecureVPNPass#42`
*   **Used By**: VPN Stack (if deployed).

### System & Runtime

These variables control system-wide settings, permissions, and debugging.

#### `PUID`

*   **Description**: The numeric User ID (UID) that Docker containers should use to run their processes. This is crucial for ensuring correct file ownership and permissions on volumes mounted from the host filesystem. It should typically match the UID of your unprivileged user on the host.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Used By**: All Docker Compose stacks (via the `user` directive for services).

#### `PGID`

*   **Description**: The numeric Group ID (GID) that Docker containers should use to run their processes. Similar to `PUID`, this ensures correct group ownership and permissions on mounted volumes. It should typically match the GID of your unprivileged user on the host.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Used By**: All Docker Compose stacks (via the `user` directive for services).

#### `TZ`

*   **Description**: Specifies the timezone for all M3TAL components and containers. This affects timestamps in logs, scheduled tasks, and any time-sensitive operations.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Used By**: All Docker Compose stacks (via the `environment` directive for services), M3TAL API daemon.

#### `LOG_LEVEL`

*   **Description**: Sets the minimum severity level for logs generated by M3TAL components. Available levels typically include `debug`, `info`, `warn`, `error`, `fatal`. Setting to `debug` provides the most verbose output.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Used By**: `m3tal-api.service`, `m3tal-dashboard` container, potentially other services.

#### `DEBUG_MODE`

*   **Description**: A boolean flag that enables or disables debug-specific features, verbose logging, or development-related behaviors across M3TAL components.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Used By**: M3TAL API daemon, M3TAL Dashboard, potentially other services.

#### `METRICS_ENABLED`

*   **Description**: A boolean flag to enable or disable the collection and exposure of system metrics by M3TAL components. When enabled, metrics endpoints (e.g., `/metrics`) might become available for monitoring tools like Prometheus.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Used By**: M3TAL API daemon, potentially other services.