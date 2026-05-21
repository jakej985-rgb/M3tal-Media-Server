# Environment Variables Reference

All M3TAL system and stack configurations are managed via environment variables. These variables are read from the primary configuration file located at `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (managed by `m3tal up` and `m3tal dash up`) source their environment variables from this file using the `--env-file` option.

It is highly recommended to manage these variables using the `m3tal config wizard` or `m3tal config set <KEY> <VALUE>` commands to ensure consistency and prevent errors.

---

## Quick Reference

| Variable | Description | Default Value | Example Value | Component(s) |
| :---------------------- | :----------------------------------------------------------------------------------------------------------------- | :-------------------- | :------------------------------------------------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `DASHBOARD_PORT` | Port on which the M3TAL Dashboard container listens. | `8082` | `8083` | M3TAL Dashboard |
| `DASHBOARD_EXPOSE_MODE` | Controls how the Dashboard is exposed: `local` (direct port) or `traefik` (via reverse proxy). | `local` | `traefik` | CLI, M3TAL Dashboard |
| `HTTP_PORT` | The internal port the M3TAL API daemon listens on. | `8080` | `8081` | M3TAL API Daemon |
| `STATE_DIR` | Base directory for the API daemon's `state.db` if not using system default. | `./state` | `/var/lib/m3tal` | M3TAL API Daemon |
| `LOG_LEVEL` | Logging verbosity for the CLI and API daemon. | `info` | `debug` | CLI, M3TAL API Daemon |
| `DASHBOARD_SECRET` | Secret key used by the Dashboard for session management and encryption. | `change_me_immediately` | `super_secret_dash_key` | M3TAL Dashboard |
| `API_TOKEN` | Token used for authenticating with the M3TAL API. | `change_me_api_token` | `my_strong_api_token` | CLI, M3TAL API Daemon |
| `ADMIN_PASSWORD` | Initial password for the default admin user of the M3TAL Dashboard. | `admin_pass` | `mySecurePass123` | M3TAL Dashboard |
| `NETWORK_NAME` | The name of the Docker network used by M3TAL and all connected stacks. | `m3tal` | `m3tal-proxy` | All Compose Stacks |
| `LOCAL_IP` | Local IP address for host-gateway connections within Docker, or API binding. | `127.0.0.1` | `192.168.1.100` | M3TAL API Daemon, Traefik |
| `DOMAIN` | The base domain for Traefik routing rules (e.g., `api.DOMAIN`, `dash.DOMAIN`). | `localhost` | `m3tal.example.com` | Traefik, CLI, M3TAL API Daemon, M3TAL Dashboard, Cloudflared |
| `VPN_USER` | Username for VPN services (e.g., Cloudflared tunnel credentials). | `user` | `m3taluser` | Cloudflared, User VPN Stacks |
| `VPN_PASSWORD` | Password for VPN services. | `password` | `mySecureVPNPass` | Cloudflared, User VPN Stacks |
| `BASE_STORAGE_PATH` | The root directory on the host where all application data is stored. | `./data` | `/mnt/m3tal-data` | All Compose Stacks, M3TAL Dashboard |
| `MEDIA_PATH` | Subdirectory within `BASE_STORAGE_PATH` for media files. | `./data/media` | `${BASE_STORAGE_PATH}/media` | User Stacks |
| `CONFIG_PATH` | Subdirectory within `BASE_STORAGE_PATH` for configuration files. | `./data/config` | `${BASE_STORAGE_PATH}/config` | User Stacks, M3TAL Dashboard |
| `DOWNLOADS_PATH` | Subdirectory within `BASE_STORAGE_PATH` for downloaded content. | `./data/downloads` | `${BASE_STORAGE_PATH}/downloads` | User Stacks |
| `PUID` | The User ID (UID) used by containers for file permissions. | `1000` | `1001` | All Containers |
| `PGID` | The Group ID (GID) used by containers for file permissions. | `1000` | `1001` | All Containers |
| `TZ` | The timezone for all M3TAL components and containers. | `America/Denver` | `Europe/London` | All Containers, CLI, M3TAL API Daemon |
| `TRAEFIK_WEB_PORT` | The HTTP entrypoint port for Traefik. | `80` | `8080` | Traefik |
| `TRAEFIK_WEBHTTPS_PORT` | The HTTPS entrypoint port for Traefik. | `443` | `8443` | Traefik |
| `TRAEFIK_DASHBOARD_PORT` | The internal port for Traefik's own dashboard (exposed locally on `127.0.0.1:8081`). | `8080` | `8081` | Traefik |
| `DEBUG_MODE` | Enables debug logging and features for relevant components. | `false` | `true` | CLI, M3TAL API Daemon |
| `METRICS_ENABLED` | Enables or disables metrics exposure for the M3TAL API daemon. | `true` | `false` | M3TAL API Daemon |

---

## Detailed Variable Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system, including logging, and user/group IDs for container processes.

#### `LOG_LEVEL`

*   **Description**: Sets the logging verbosity for the M3TAL CLI and API daemon. Valid values typically include `debug`, `info`, `warn`, `error`.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Component(s)**: CLI, M3TAL API Daemon

#### `STATE_DIR`

*   **Description**: Specifies the base directory where the M3TAL API daemon stores its SQLite state database (`state.db`).
    *   **Note**: In production, the API daemon typically defaults to `/var/lib/m3tal/state.db` managed by the system. This variable allows overriding that default for advanced use cases (e.g., local development). The M3TAL Dashboard's internal `/docker/state` path is controlled by the `CONFIG_PATH` variable through a volume mount, not this `STATE_DIR` directly.
*   **Default Value**: `./state`
*   **Example Value**: `/opt/m3tal/state`
*   **Component(s)**: M3TAL API Daemon

#### `PUID`

*   **Description**: The User ID (UID) that all M3TAL containers and user-deployed stacks will use for file permissions within their mounted volumes. This ensures proper file ownership and access rights on the host system.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Component(s)**: All Containers

#### `PGID`

*   **Description**: The Group ID (GID) that all M3TAL containers and user-deployed stacks will use for file permissions within their mounted volumes. This complements `PUID` for managing file access.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Component(s)**: All Containers

#### `TZ`

*   **Description**: Sets the timezone for all M3TAL components and containers. This ensures consistent timestamps in logs and applications.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Component(s)**: All Containers, CLI, M3TAL API Daemon

#### `DEBUG_MODE`

*   **Description**: A boolean flag to enable or disable debug-specific features or logging in relevant components. Setting to `true` may provide more verbose output and enable development utilities.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Component(s)**: CLI, M3TAL API Daemon

#### `METRICS_ENABLED`

*   **Description**: Controls whether the M3TAL API daemon exposes prometheus-compatible metrics for monitoring.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Component(s)**: M3TAL API Daemon

### Authentication

These variables are crucial for securing access to the M3TAL Dashboard and API.

#### `DASHBOARD_SECRET`

*   **Description**: A cryptographic secret key used by the M3TAL Dashboard for session management, cookie signing, and other security-sensitive operations.
    *   **Important**: This variable is **auto-generated** during the first `m3tal init` run. **Users should NOT set this manually** unless performing a security rotation.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `my_super_secure_dashboard_secret_key`
*   **Component(s)**: M3TAL Dashboard

#### `API_TOKEN`

*   **Description**: The authentication token used by the M3TAL CLI and external clients to interact with the M3TAL API daemon. It's also used internally for API daemon security.
    *   **Important**: This variable is **auto-generated** during the first `m3tal init` run. **Users should NOT set this manually** unless performing a security rotation.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `my_very_strong_api_token_123abc`
*   **Component(s)**: CLI, M3TAL API Daemon

#### `ADMIN_PASSWORD`

*   **Description**: Sets the initial password for the default `admin` user of the M3TAL Dashboard. After initial login, you can manage users and passwords via the dashboard.
*   **Default Value**: `admin_pass`
*   **Example Value**: `myComplexDashboardAdminPass`
*   **Component(s)**: M3TAL Dashboard

### Network Configuration

Variables related to how M3TAL components communicate over the network.

#### `HTTP_PORT`

*   **Description**: The internal port on which the M3TAL API daemon (the Go binary) listens for incoming HTTP requests. This port is typically accessed by Traefik or the Dashboard directly via `host.docker.internal`.
*   **Default Value**: `8080`
*   **Example Value**: `8081`
*   **Component(s)**: M3TAL API Daemon

#### `NETWORK_NAME`

*   **Description**: Defines the name of the Docker network that M3TAL's core components and all user-deployed stacks will share. This enables seamless communication between containers.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal-proxy-network`
*   **Component(s)**: All Compose Stacks

#### `LOCAL_IP`

*   **Description**: Specifies a local IP address primarily used for Docker's `host.docker.internal` resolution, allowing containers to access services running directly on the host machine. It can also be used for specific API binding if needed.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.10`
*   **Component(s)**: M3TAL API Daemon, Traefik (for API routing)

### Storage Paths

These variables define the base directories on the host filesystem where M3TAL stores various types of data.

#### `BASE_STORAGE_PATH`

*   **Description**: The **root directory** on the host filesystem where all M3TAL application data, media, configuration, and downloads are stored. This path is volume-mounted into containers.
    *   **Important**: While the template defaults to `./data`, in **production deployments**, this should typically be set to `/mnt` or another dedicated high-capacity storage mount point (e.g., `/mnt/m3tal-data`).
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal-data`
*   **Component(s)**: All Compose Stacks, M3TAL Dashboard

#### `MEDIA_PATH`

*   **Description**: Specifies the subdirectory (relative to `BASE_STORAGE_PATH`) where user-deployed media applications (e.g., Plex, Jellyfin) should store their media files.
*   **Default Value**: `./data/media`
*   **Example Value**: `${BASE_STORAGE_PATH}/media` (resolves to `/mnt/m3tal-data/media` if `BASE_STORAGE_PATH=/mnt/m3tal-data`)
*   **Component(s)**: User Stacks

#### `CONFIG_PATH`

*   **Description**: Specifies the subdirectory (relative to `BASE_STORAGE_PATH`) where configuration files and persistent state for applications are stored.
    *   **Note**: For the M3TAL Dashboard, this variable determines where the host directory for `/docker/state` (containing `users.json`) is mounted from. Specifically, the Dashboard uses `${CONFIG_PATH}/m3tal/state` as its host volume source.
*   **Default Value**: `./data/config`
*   **Example Value**: `${BASE_STORAGE_PATH}/config` (resolves to `/mnt/m3tal-data/config` if `BASE_STORAGE_PATH=/mnt/m3tal-data`)
*   **Component(s)**: User Stacks, M3TAL Dashboard

#### `DOWNLOADS_PATH`

*   **Description**: Specifies the subdirectory (relative to `BASE_STORAGE_PATH`) for storing downloaded content from user-deployed download clients (e.g., torrent clients).
*   **Default Value**: `./data/downloads`
*   **Example Value**: `${BASE_STORAGE_PATH}/downloads` (resolves to `/mnt/m3tal-data/downloads` if `BASE_STORAGE_PATH=/mnt/m3tal-data`)
*   **Component(s)**: User Stacks

### Traefik Gateway

Variables controlling the Traefik reverse proxy and its routing behavior.

#### `DOMAIN`

*   **Description**: The base domain name used by Traefik to define routing rules for M3TAL's services. Setting this variable enables routes like `dash.${DOMAIN}` for the Dashboard and `api.${DOMAIN}` for the API daemon.
*   **Default Value**: `localhost`
*   **Example Value**: `m3tal.example.com`
*   **Component(s)**: Traefik, CLI, M3TAL API Daemon, M3TAL Dashboard, Cloudflared

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Determines how the M3TAL Dashboard container is exposed to the network.
    *   `local` (Default): The Dashboard is exposed directly on a host port (default `8082`) via a Docker port binding. Access is `http://HOST_IP:8082`. No Traefik required. Best for LAN-only or initial setup.
    *   `traefik`: The Dashboard is exposed via Traefik using its domain routing. Access is `http://dash.${DOMAIN}`. Requires Traefik to be running and `DOMAIN` to be configured.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Component(s)**: CLI (selects compose override), M3TAL Dashboard (implicit config)

#### `DASHBOARD_PORT`

*   **Description**: The internal port on which the M3TAL Dashboard container listens.
    *   If `DASHBOARD_EXPOSE_MODE` is `local`, this port is directly mapped to the host (e.g., `8082:8082`).
    *   If `DASHBOARD_EXPOSE_MODE` is `traefik`, Traefik routes traffic for `dash.${DOMAIN}` to this port internally.
*   **Default Value**: `8082`
*   **Example Value**: `8083`
*   **Component(s)**: M3TAL Dashboard

#### `TRAEFIK_WEB_PORT`

*   **Description**: The host port Traefik uses as its primary HTTP entry point. This is typically port `80` for standard web traffic.
*   **Default Value**: `80`
*   **Example Value**: `8080`
*   **Component(s)**: Traefik

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port Traefik uses as its primary HTTPS entry point. This is typically port `443` for secure web traffic.
*   **Default Value**: `443`
*   **Example Value**: `8443`
*   **Component(s)**: Traefik

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port used by Traefik for its own management dashboard. In M3TAL, this is typically exposed on `127.0.0.1:8081` on the host, not publicly.
*   **Default Value**: `8080`
*   **Example Value**: `8081`
*   **Component(s)**: Traefik

### VPN Configuration

Variables related to Virtual Private Network services, particularly Cloudflared.

#### `VPN_USER`

*   **Description**: Username credential for VPN services (e.g., a Cloudflare Tunnel credential).
*   **Default Value**: `user`
*   **Example Value**: `m3tal_cloudflared_user`
*   **Component(s)**: Cloudflared, User VPN Stacks

#### `VPN_PASSWORD`

*   **Description**: Password credential for VPN services (e.g., a Cloudflare Tunnel credential).
*   **Default Value**: `password`
*   **Example Value**: `mySecureCloudflareTunnelPass`
*   **Component(s)**: Cloudflared, User VPN Stacks