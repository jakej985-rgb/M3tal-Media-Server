# Environment Variables Reference

As the M3TAL Ecosystem Documentation Architect, my goal is to provide you with a comprehensive understanding of the configuration variables that power your M3TAL system.

All environment variables are read by the M3TAL CLI and all Docker Compose stacks from the primary configuration file: `/etc/m3tal/.env`. This file is loaded automatically by `m3tal` commands (like `m3tal up` or `m3tal dash up`) via the `--env-file` argument.

You should primarily manage these variables using the `m3tal config wizard` for interactive setup or `m3tal config set KEY value` for specific changes. Avoid manually editing `/etc/m3tal/.env` directly to prevent syntax errors.

---

## Quick Reference Table

| Name | Description | Default Value | Example Value | Component(s) Used |
| :---------------------- | :------------------------------------------------------------------------------------------------------ | :-------------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| `DASHBOARD_PORT` | Port for the M3TAL Dashboard in `local` expose mode. | `8082` | `8082` | Dashboard (local mode) |
| `DASHBOARD_EXPOSE_MODE` | Controls how the M3TAL Dashboard is exposed: `local` (direct port) or `traefik` (via domain). | `local` | `traefik` | CLI (`m3tal dash up`), Dashboard |
| `HTTP_PORT` | The host port on which the M3TAL API daemon listens. | `8080` | `8080` | API Daemon |
| `STATE_DIR` | Base directory for M3TAL state files (CLI/development). See `CONFIG_PATH` for production persistence. | `./state` | `/var/lib/m3tal` | CLI, Dashboard (internal container path) |
| `LOG_LEVEL` | Minimum logging level for the CLI and API daemon (`debug`, `info`, `warn`, `error`). | `info` | `debug` | CLI, API Daemon |
| `DASHBOARD_SECRET` | Secret key for M3TAL Dashboard session management. **Auto-generated.** | `change_me_immediately` | `super_secret_key_123` | Dashboard |
| `API_TOKEN` | Bearer token for authenticating CLI requests to the M3TAL API. **Auto-generated.** | `change_me_api_token` | `some_long_jwt_token` | CLI, API Daemon |
| `ADMIN_PASSWORD` | Initial password for the default M3TAL Dashboard admin user. | `admin_pass` | `secure_admin_pass` | CLI (`m3tal dashpass`) |
| `NETWORK_NAME` | Base name for user-defined Docker networks created by M3TAL. | `m3tal` | `myproject` | User Stacks (Docker Compose) |
| `LOCAL_IP` | Local IP address for service bindings or internal communication. | `127.0.0.1` | `192.168.1.100` | API Daemon (potential), System |
| `DOMAIN` | Your primary domain for Traefik routing (e.g., `example.com`). Enables `dash.DOMAIN`, `api.DOMAIN` routes. | `localhost` | `m3tal.local` | Traefik, Dashboard (traefik mode), API Daemon (dynamic config) |
| `VPN_USER` | Username for VPN/tunnel services (e.g., Cloudflared Access). | `user` | `admin` | Cloudflared (future/advanced configs) |
| `VPN_PASSWORD` | Password for VPN/tunnel services. | `password` | `myvpnpass` | Cloudflared (future/advanced configs) |
| `BASE_STORAGE_PATH` | Base host path for all M3TAL data (media, config, downloads). | `./data` | `/mnt/data` | CLI, API Daemon, Dashboard, User Stacks |
| `MEDIA_PATH` | Host path for media files. Derived from `BASE_STORAGE_PATH`. | `./data/media` | `/mnt/data/media` | User Stacks |
| `CONFIG_PATH` | Host path for M3TAL configuration files and persistent state. | `./data/config` | `/etc/m3tal/config` | CLI, API Daemon, Dashboard |
| `DOWNLOADS_PATH` | Host path for downloaded content. Derived from `BASE_STORAGE_PATH`. | `./data/downloads` | `/mnt/data/downloads` | User Stacks |
| `PUID` | POSIX User ID for containers to ensure correct file permissions. | `1000` | `1000` | Dashboard, User Stacks |
| `PGID` | POSIX Group ID for containers to ensure correct file permissions. | `1000` | `1000` | Dashboard, User Stacks |
| `TZ` | Timezone for containers to ensure accurate timekeeping. | `America/Denver` | `Europe/London` | Dashboard, User Stacks |
| `TRAEFIK_WEB_PORT` | The host port Traefik uses for HTTP traffic. | `80` | `80` | Traefik |
| `TRAEFIK_WEBHTTPS_PORT` | The host port Traefik uses for HTTPS traffic. | `443` | `443` | Traefik (HTTPS config) |
| `TRAEFIK_DASHBOARD_PORT` | The internal port Traefik's own dashboard listens on within its container. | `8080` | `8080` | Traefik |
| `DEBUG_MODE` | Enable debug logging and potentially other debug features. | `false` | `true` | CLI, API Daemon, Dashboard |
| `METRICS_ENABLED` | Enable or disable Prometheus metrics exposure for the API daemon. | `true` | `false` | API Daemon |

---

## Detailed Environment Variable Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system's operation and logging.

#### `LOG_LEVEL`

*   **Description**: Sets the minimum logging level for the M3TAL CLI and API daemon. Valid values are `debug`, `info`, `warn`, `error`. Setting to `debug` provides the most verbose output, useful for troubleshooting.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Component(s) Used**: CLI binary, API daemon

#### `DEBUG_MODE`

*   **Description**: A general flag to enable verbose debugging across M3TAL components. When `true`, it may activate more detailed logging, enable specific debug endpoints, or alter behavior for diagnostic purposes.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Component(s) Used**: CLI binary, API daemon, Dashboard container

#### `METRICS_ENABLED`

*   **Description**: Controls whether the M3TAL API daemon exposes Prometheus-compatible metrics for monitoring system health and performance. Set to `false` to disable metrics.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Component(s) Used**: API daemon

#### `PUID`

*   **Description**: POSIX User ID. This value is used to set the user ID for processes running inside M3TAL containers (like `m3tal-dashboard`) to ensure consistent file permissions with your host system. It helps prevent permission issues when containers interact with host-mounted volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Component(s) Used**: Dashboard container, User-defined Docker Compose stacks

#### `PGID`

*   **Description**: POSIX Group ID. Similar to `PUID`, this value sets the group ID for processes within M3TAL containers to maintain correct file system permissions.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Component(s) Used**: Dashboard container, User-defined Docker Compose stacks

#### `TZ`

*   **Description**: Sets the timezone for M3TAL containers, ensuring that timestamps and time-based operations within services are accurate and consistent with your local time or preferred zone.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Component(s) Used**: Dashboard container, User-defined Docker Compose stacks

### Authentication

These variables manage access and security credentials for the M3TAL system.

#### `DASHBOARD_SECRET`

*   **Description**: A secret key used by the M3TAL Dashboard for session management and securing internal communications.
    **Note**: This value is automatically generated on the first `m3tal init` command. You should **NOT** set this manually unless you are intentionally rotating the secret for security reasons.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `super_secret_key_1234567890abcdef` (a long, randomly generated string)
*   **Component(s) Used**: Dashboard container

#### `API_TOKEN`

*   **Description**: A bearer token used by the M3TAL CLI to authenticate with the M3TAL API daemon.
    **Note**: This value is automatically generated on the first `m3tal init` command. You should **NOT** set this manually unless you are intentionally rotating the token for security reasons.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIn0.SflKxwRJSMeKKF2QT4fWPfH_GMQ-P1t7Q2_Yq0eP7fQ` (a long JWT)
*   **Component(s) Used**: CLI binary, API daemon

#### `ADMIN_PASSWORD`

*   **Description**: The initial password for the default administrative user in the M3TAL Dashboard. You can change this using `m3tal dashpass`.
*   **Default Value**: `admin_pass`
*   **Example Value**: `secure_admin_pass_!@#`
*   **Component(s) Used**: CLI (`m3tal dashpass` command to create/update dashboard users.json)

### Network Configuration

These variables define how M3TAL services communicate and are exposed.

#### `HTTP_PORT`

*   **Description**: Specifies the host port on which the M3TAL API daemon listens for incoming connections. This is the primary interface for the CLI and Dashboard to interact with the API.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Component(s) Used**: API daemon

#### `DASHBOARD_PORT`

*   **Description**: Defines the host port used to access the M3TAL Dashboard when `DASHBOARD_EXPOSE_MODE` is set to `local`. If set to `traefik`, this port is only used internally by Traefik for routing.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Component(s) Used**: Dashboard container (when `DASHBOARD_EXPOSE_MODE=local`)

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Determines how the M3TAL Dashboard is made accessible:
    *   `local`: The dashboard is exposed directly on `http://HOST_IP:${DASHBOARD_PORT}`. No Traefik required. Best for LAN-only setups or initial configuration.
    *   `traefik`: The dashboard is exposed via Traefik at `http://dash.${DOMAIN}`. Requires Traefik to be running. Best for domain-based setups.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Component(s) Used**: CLI (`m3tal dash up` for selecting compose override), Dashboard container (Traefik labels)

#### `NETWORK_NAME`

*   **Description**: This variable can be used as a prefix or base name for user-defined Docker networks created by M3TAL or user stacks. While the core control plane uses a fixed `proxy` network, this offers flexibility for custom deployments.
*   **Default Value**: `m3tal`
*   **Example Value**: `myproject`
*   **Component(s) Used**: User-defined Docker Compose stacks

#### `LOCAL_IP`

*   **Description**: Specifies a local IP address that M3TAL components might use for binding to network interfaces or for internal host-level communication.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100`
*   **Component(s) Used**: API daemon (potential binding), System

### Storage Paths

These variables define where M3TAL stores its data on the host filesystem.

#### `BASE_STORAGE_PATH`

*   **Description**: The primary host directory for all M3TAL-related data, including configuration, media, and downloads.
    **Note**: In production deployments, this defaults to `/mnt` (e.g., if `/mnt` is a dedicated data volume), not `./data` as in the template. The `m3tal init` process typically sets this for you. Containers will often mount this path internally as `/mnt`.
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal_data`
*   **Component(s) Used**: CLI binary, API daemon, Dashboard container (volume mounts), User-defined Docker Compose stacks

#### `CONFIG_PATH`

*   **Description**: The host path dedicated to M3TAL's core configuration files and persistent state, including database files and user credentials. This path is crucial for persisting system state across reboots and container updates.
*   **Default Value**: `./data/config`
*   **Example Value**: `/etc/m3tal/config`
*   **Component(s) Used**: CLI binary, API daemon, Dashboard container (volume mounts)

#### `MEDIA_PATH`

*   **Description**: The host path designated for media files managed by M3TAL or user-deployed applications. This path is often a subdirectory of `BASE_STORAGE_PATH`.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/m3tal_data/media`
*   **Component(s) Used**: User-defined Docker Compose stacks (for shared media libraries)

#### `DOWNLOADS_PATH`

*   **Description**: The host path intended for downloaded content managed by M3TAL or other services. Typically a subdirectory of `BASE_STORAGE_PATH`.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/m3tal_data/downloads`
*   **Component(s) Used**: User-defined Docker Compose stacks (for download clients)

#### `STATE_DIR`

*   **Description**: A base directory for various M3TAL state files. While the API daemon primarily uses `/var/lib/m3tal/state.db` and the Dashboard maps its internal `/docker/state` via `CONFIG_PATH`, this variable can serve as a fallback or for local development environments to specify where temporary or local state data is stored.
*   **Default Value**: `./state`
*   **Example Value**: `/tmp/m3tal_state`
*   **Component(s) Used**: CLI binary (local operations), Dashboard container (internal path)

### Traefik Gateway

These variables configure the Traefik reverse proxy for routing external traffic.

#### `DOMAIN`

*   **Description**: Your primary domain name. Setting this enables Traefik to create routing rules like `dash.YOUR_DOMAIN` for the M3TAL Dashboard and `api.YOUR_DOMAIN` for the M3TAL API. If not set, services may only be accessible via IP address or `localhost`.
*   **Default Value**: `localhost`
*   **Example Value**: `m3tal.local`
*   **Component(s) Used**: Traefik gateway, Dashboard container (when `DASHBOARD_EXPOSE_MODE=traefik`), API daemon (Traefik dynamic configuration)

#### `TRAEFIK_WEB_PORT`

*   **Description**: The host port Traefik binds to for handling incoming HTTP (non-HTTPS) requests.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Component(s) Used**: Traefik gateway

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port Traefik binds to for handling incoming HTTPS requests. While not always directly used in basic HTTP setups, it's essential for configuring secure communication.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Component(s) Used**: Traefik gateway (for HTTPS configuration)

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port within the Traefik container where its own management dashboard is exposed. On the host, this is typically mapped to `127.0.0.1:8081` for local access.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Component(s) Used**: Traefik gateway

### VPN Integration (Cloudflared)

These variables are for configuring VPN or tunnel services, such as Cloudflared, for secure remote access.

#### `VPN_USER`

*   **Description**: Username credential that may be used by VPN or tunnel services (e.g., Cloudflared Access) for authentication. Specific usage depends on the VPN solution configured.
*   **Default Value**: `user`
*   **Example Value**: `admin`
*   **Component(s) Used**: Cloudflared container (for advanced configurations like Cloudflare Access)

#### `VPN_PASSWORD`

*   **Description**: Password credential that may be used by VPN or tunnel services (e.g., Cloudflared Access) for authentication. Specific usage depends on the VPN solution configured.
*   **Default Value**: `password`
*   **Example Value**: `mysecurevpnpass`
*   **Component(s) Used**: Cloudflared container (for advanced configurations like Cloudflare Access)