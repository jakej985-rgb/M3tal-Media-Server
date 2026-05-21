As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the definitive reference for M3TAL environment variables.

---

# Environment Variables Reference (`docs/ENVIRONMENT_VARIABLES.md`)

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables control the behavior of the M3TAL CLI, the API daemon, the Dashboard, Traefik gateway, and all Docker Compose stacks.

All environment variables are read from the primary M3TAL configuration file: `/etc/m3tal/.env`. Both the M3TAL CLI and all Docker Compose stacks (including the M3TAL control plane and user-defined services) are configured to load variables from this file via the `--env-file` option. You can manage these variables conveniently using the `m3tal config wizard` or `m3tal config set <KEY> <value>` commands.

---

## Quick Reference Table

| Variable Name           | Default Value         | Group                     | Description                                                                     |
| :---------------------- | :-------------------- | :------------------------ | :------------------------------------------------------------------------------ |
| `HTTP_PORT`             | `8080`                | Core M3TAL Services       | Port for the M3TAL API daemon.                                                  |
| `LOG_LEVEL`             | `info`                | Core M3TAL Services       | Logging verbosity for the M3TAL API.                                            |
| `DEBUG_MODE`            | `false`               | Core M3TAL Services       | Enables debug logging and features.                                             |
| `METRICS_ENABLED`       | `true`                | Core M3TAL Services       | Enables/disables M3TAL API metrics collection.                                  |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Authentication & Security | Secret key for M3TAL Dashboard session management. **(Auto-generated)**         |
| `API_TOKEN`             | `change_me_api_token` | Authentication & Security | Authentication token for M3TAL API daemon. **(Auto-generated)**                 |
| `ADMIN_PASSWORD`        | `admin_pass`          | Authentication & Security | Initial password for Dashboard admin user.                                      |
| `NETWORK_NAME`          | `m3tal`               | Networking & Docker       | Name of the Docker network for M3TAL containers.                                |
| `LOCAL_IP`              | `127.0.0.1`           | Networking & Docker       | Host machine's IP, used by Docker for `host.docker.internal`.                   |
| `BASE_STORAGE_PATH`     | `./data`              | Storage & Permissions     | Base directory for all M3TAL persistent data.                                   |
| `MEDIA_PATH`            | `./data/media`        | Storage & Permissions     | Subdirectory for media files.                                                   |
| `CONFIG_PATH`           | `./data/config`       | Storage & Permissions     | Subdirectory for configuration files.                                           |
| `DOWNLOADS_PATH`        | `./data/downloads`    | Storage & Permissions     | Subdirectory for downloaded content.                                            |
| `PUID`                  | `1000`                | Storage & Permissions     | User ID for container processes.                                                |
| `PGID`                  | `1000`                | Storage & Permissions     | Group ID for container processes.                                               |
| `DASHBOARD_PORT`        | `8082`                | M3TAL Dashboard           | Host port for the Dashboard in `local` expose mode.                             |
| `DASHBOARD_EXPOSE_MODE` | `local`               | M3TAL Dashboard           | Controls Dashboard access: `local` (direct port) or `traefik` (via domain).     |
| `STATE_DIR`             | `./state`             | M3TAL Dashboard           | Internal path for Dashboard container state files.                              |
| `TZ`                    | `America/Denver`      | M3TAL Dashboard           | Timezone for containers.                                                        |
| `DOMAIN`                | `localhost`           | Traefik Gateway           | Base domain for Traefik routing (e.g., `dash.DOMAIN`, `api.DOMAIN`).            |
| `TRAEFIK_WEB_PORT`      | `80`                  | Traefik Gateway           | Traefik's HTTP entry point port.                                                |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                 | Traefik Gateway           | Traefik's HTTPS entry point port (if configured).                               |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                | Traefik Gateway           | Internal port for Traefik's own dashboard.                                      |
| `VPN_USER`              | `user`                | VPN Integration           | Username for VPN services (e.g., Cloudflared).                                  |
| `VPN_PASSWORD`          | `password`            | VPN Integration           | Password for VPN services (e.g., Cloudflared).                                  |

---

## Detailed Environment Variable Reference

### Core M3TAL Services

Variables controlling the foundational M3TAL API daemon and core system behavior.

#### `HTTP_PORT`
*   **Description**: The port on which the M3TAL API daemon listens for incoming HTTP requests. This port is accessible host-locally and is the target for Traefik's API routing.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Component Uses It**: `m3tal-api.service`, Traefik gateway (routes to `http://host.docker.internal:8080`)

#### `LOG_LEVEL`
*   **Description**: Controls the verbosity of logging for the M3TAL API daemon. Accepted values typically include `debug`, `info`, `warn`, `error`, `fatal`.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Component Uses It**: `m3tal-api.service`

#### `DEBUG_MODE`
*   **Description**: A boolean flag to enable or disable debug-level logging and potentially other development/debugging features across M3TAL components.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Component Uses It**: `m3tal-api.service`, potentially `m3tal-dashboard`

#### `METRICS_ENABLED`
*   **Description**: A boolean flag to enable or disable the collection and exposure of operational metrics from the M3TAL API daemon.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Component Uses It**: `m3tal-api.service`

### Authentication & Security

Variables critical for M3TAL system security and user authentication.

#### `DASHBOARD_SECRET`
*   **Description**: A crucial secret key used by the M3TAL Dashboard for secure session management, CSRF protection, and other cryptographic operations.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_long_random_string_of_characters_for_security`
*   **Component Uses It**: `m3tal-dashboard`
*   **Note**: This variable is **auto-generated** on first `m3tal init`. Users should **NOT** set it manually unless performing a secret rotation.

#### `API_TOKEN`
*   **Description**: An authentication token required to access the M3TAL API daemon. This token secures communication between the dashboard, CLI, and the API daemon.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `another_long_random_string_for_api_auth`
*   **Component Uses It**: CLI, `m3tal-dashboard`
*   **Note**: This variable is **auto-generated** on first `m3tal init`. Users should **NOT** set it manually unless performing a token rotation.

#### `ADMIN_PASSWORD`
*   **Description**: The initial password for the default administrative user created in the M3TAL Dashboard. It is strongly recommended to change this after initial setup via the dashboard or `m3tal dashpass` command.
*   **Default Value**: `admin_pass`
*   **Example Value**: `MySecureAdminPass123!`
*   **Component Uses It**: `m3tal-dashboard` (user management)

### Networking & Docker

Variables related to Docker networking and host-container communication.

#### `NETWORK_NAME`
*   **Description**: The name of the Docker bridge network used by M3TAL's Docker Compose stacks for internal container communication. All M3TAL control plane services and user-defined stacks should connect to this network for proper functioning.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal-proxy-net`
*   **Component Uses It**: All Docker Compose stacks (including `m3tal-dashboard`, `traefik`, `cloudflared`)

#### `LOCAL_IP`
*   **Description**: Specifies the host machine's IP address. Primarily used by Docker for resolving `host.docker.internal` to allow containers to access services running directly on the host (e.g., the M3TAL API daemon).
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100`
*   **Component Uses It**: Docker Compose (via `extra_hosts` mappings in `m3tal-compose.yml` and `routing-compose.yml`)

### Storage & Permissions

Variables defining where M3TAL stores its persistent data and the user/group IDs for container processes.

#### `BASE_STORAGE_PATH`
*   **Description**: The foundational directory on the host where all M3TAL's persistent data (media, configuration, downloads) is stored. All other storage paths (`MEDIA_PATH`, `CONFIG_PATH`, `DOWNLOADS_PATH`) are typically subdirectories of this.
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal_data`
*   **Component Uses It**: `m3tal-dashboard` (via volume mounts), User-defined Docker Compose stacks
*   **Note**: In production deployments, this variable typically defaults to `/mnt`.

#### `MEDIA_PATH`
*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for storing media files managed by various M3TAL-integrated services.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/m3tal_data/media`
*   **Component Uses It**: User-defined Docker Compose stacks

#### `CONFIG_PATH`
*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for M3TAL's core configuration files, such as the Dashboard's `users.json` and other stack-specific configurations.
*   **Default Value**: `./data/config`
*   **Example Value**: `/mnt/m3tal_data/config`
*   **Component Uses It**: `m3tal-dashboard` (for `users.json` volume mount), CLI (`m3tal config wizard`)

#### `DOWNLOADS_PATH`
*   **Description**: The subdirectory within `BASE_STORAGE_PATH` for storing downloaded content from various services.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/m3tal_data/downloads`
*   **Component Uses It**: User-defined Docker Compose stacks

#### `PUID`
*   **Description**: The User ID (UID) that containers will run as. Used to ensure correct file permissions when containers interact with host-mounted volumes, preventing permission issues.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Component Uses It**: `m3tal-dashboard`, User-defined Docker Compose stacks

#### `PGID`
*   **Description**: The Group ID (GID) that containers will run as. Used to ensure correct file permissions when containers interact with host-mounted volumes, preventing permission issues.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Component Uses It**: `m3tal-dashboard`, User-defined Docker Compose stacks

### M3TAL Dashboard

Variables specifically related to the M3TAL Dashboard container and its exposure.

#### `DASHBOARD_PORT`
*   **Description**: The host port on which the M3TAL Dashboard container is exposed when `DASHBOARD_EXPOSE_MODE` is set to `local`. This is also the internal port the Dashboard container listens on.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Component Uses It**: `m3tal-dashboard` (container), CLI (`m3tal dash up`)

#### `DASHBOARD_EXPOSE_MODE`
*   **Description**: Determines how the M3TAL Dashboard is made accessible on the network.
    *   `local`: Uses a direct port binding (`HOST_IP:DASHBOARD_PORT`). No Traefik required. Best for LAN-only setups or initial testing.
    *   `traefik`: Exposed via the Traefik gateway at `http://dash.DOMAIN`. Requires the Traefik gateway to be running. Best for domain-based setups.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Component Uses It**: CLI (`m3tal dash up`)

#### `STATE_DIR`
*   **Description**: This variable's primary use is as an internal environment variable within the `m3tal-dashboard` container, where it's typically set to `/docker/state`. On the host, the actual path for dashboard state files (like `users.json`) is derived from `CONFIG_PATH` via a volume mount (`${CONFIG_PATH}/m3tal/state:/docker/state`). The API daemon's main state database is located at `/var/lib/m3tal/state.db` and is not directly controlled by this variable.
*   **Default Value**: `./state`
*   **Example Value**: `/docker/state` (internal container path)
*   **Component Uses It**: `m3tal-dashboard` (internal container env)

#### `TZ`
*   **Description**: Sets the timezone for containers, ensuring that logs and timestamps are synchronized with the host or a desired timezone. Essential for accurate timekeeping across services.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Component Uses It**: `m3tal-dashboard`, User-defined Docker Compose stacks

### Traefik Gateway

Variables configuring the Traefik reverse proxy for external service exposure.

#### `DOMAIN`
*   **Description**: The base domain name used by Traefik for defining dynamic routing rules. Setting this enables accessible routes like `http://dash.DOMAIN` for the dashboard and `http://api.DOMAIN` for the API daemon.
*   **Default Value**: `localhost`
*   **Example Value**: `mym3tal.com`
*   **Component Uses It**: Traefik gateway (`routing-compose.yml`), `m3tal-dashboard` (when `DASHBOARD_EXPOSE_MODE=traefik`), Traefik dynamic configuration (`dynamic/api.yml`)
*   **Note**: Setting this variable enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.

#### `TRAEFIK_WEB_PORT`
*   **Description**: The port Traefik listens on for unencrypted HTTP traffic (defined as the `web` entrypoint in Traefik's static configuration). This is typically port 80 on the host.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Component Uses It**: Traefik gateway

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description**: The port Traefik listens on for encrypted HTTPS traffic (defined as the `websecure` entrypoint in Traefik's static configuration, if HTTPS is configured). This is typically port 443 on the host.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Component Uses It**: Traefik gateway

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description**: The internal port on which Traefik's own administrative dashboard (for monitoring Traefik itself) is exposed within its container. This is mapped to `127.0.0.1:8081` on the host for local access.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Component Uses It**: Traefik gateway

### VPN Integration

Variables for configuring Virtual Private Network (VPN) services, such as Cloudflared tunnels.

#### `VPN_USER`
*   **Description**: The username for authenticating with a VPN service, typically used by containers like `cloudflared` for establishing secure tunnels to external networks.
*   **Default Value**: `user`
*   **Example Value**: `myvpnuser`
*   **Component Uses It**: `cloudflared` (and other VPN-related containers)

#### `VPN_PASSWORD`
*   **Description**: The password for authenticating with a VPN service, typically used by containers like `cloudflared`.
*   **Default Value**: `password`
*   **Example Value**: `MySecureVPNPass`
*   **Component Uses It**: `cloudflared` (and other VPN-related containers)