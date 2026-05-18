# Environment Variable Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, I've compiled this comprehensive reference for all environment variables used by the M3TAL system. These variables control the behavior of the CLI, API daemon, Dashboard, Traefik, and other M3TAL components.

All M3TAL environment variables are primarily read from the global configuration file: `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (via the `--env-file /etc/m3tal/.env` flag) leverage this central configuration. You can manage these variables using the `m3tal config wizard` or `m3tal config set KEY value` commands.

---

## Quick Reference

| Name                    | Description                                                                  | Default                  | Example                                     |
| :---------------------- | :--------------------------------------------------------------------------- | :----------------------- | :------------------------------------------ |
| `DASHBOARD_PORT`        | Port for the Dashboard container.                                            | `8082`                   | `8082`                                      |
| `DASHBOARD_EXPOSE_MODE` | How the Dashboard container's port is exposed on the host.                   | `local`                  | `all`                                       |
| `HTTP_PORT`             | Port for the M3TAL API daemon.                                               | `8080`                   | `8080`                                      |
| `STATE_DIR`             | Directory for the SQLite state database (`state.db`).                        | `./state`                | `/var/lib/m3tal`                            |
| `LOG_LEVEL`             | Logging verbosity for the API daemon.                                        | `info`                   | `debug`                                     |
| `DASHBOARD_SECRET`      | Secret key for Dashboard session management.                                 | `change_me_immediately`  | `a_very_long_random_string...`              |
| `API_TOKEN`             | API token for programmatic access to the M3TAL API.                          | `change_me_api_token`    | `another_long_random_string...`             |
| `ADMIN_PASSWORD`        | Initial password for the 'admin' user.                                       | `admin_pass`             | `SuperSecurePassword123!`                   |
| `NETWORK_NAME`          | Name of the shared Docker network.                                           | `m3tal`                  | `m3tal-prod-network`                        |
| `LOCAL_IP`              | Host's local IP, used for `host.docker.internal` resolution.                 | `127.0.0.1`              | `192.168.1.100`                             |
| `DOMAIN`                | Primary domain for M3TAL services, enables Traefik routes.                   | `localhost`              | `my-m3tal-server.com`                       |
| `VPN_USER`              | Username for the optional VPN service.                                       | `user`                   | `vpnuser`                                   |
| `VPN_PASSWORD`          | Password for the optional VPN service.                                       | `password`               | `MyStrongVPNPass!`                          |
| `BASE_STORAGE_PATH`     | Base directory for all M3TAL data storage.                                   | `./data`                 | `/mnt/m3tal-data`                           |
| `MEDIA_PATH`            | Directory for media files.                                                   | `./data/media`           | `/mnt/m3tal-data/media`                     |
| `CONFIG_PATH`           | Directory for configuration files.                                           | `./data/config`          | `/mnt/m3tal-data/config`                    |
| `DOWNLOADS_PATH`        | Directory for downloaded files.                                              | `./data/downloads`       | `/mnt/m3tal-data/downloads`                 |
| `PUID`                  | User ID (UID) for containers requiring specific file permissions.            | `1000`                   | `1001`                                      |
| `PGID`                  | Group ID (GID) for containers requiring specific file permissions.           | `1000`                   | `1001`                                      |
| `TZ`                    | Timezone setting for M3TAL services and containers.                          | `America/Denver`         | `Europe/London`                             |
| `TRAEFIK_WEB_PORT`      | Traefik's HTTP entry point port.                                             | `80`                     | `80`                                        |
| `TRAEFIK_WEBHTTPS_PORT` | Traefik's HTTPS entry point port.                                            | `443`                    | `443`                                       |
| `TRAEFIK_DASHBOARD_PORT`| Traefik's internal dashboard container port.                                 | `8080`                   | `8080`                                      |
| `DEBUG_MODE`            | Enables debug logging and development features.                              | `false`                  | `true`                                      |
| `METRICS_ENABLED`       | Enables Prometheus-compatible metrics from the API daemon.                   | `true`                   | `false`                                     |

---

## Detailed Reference

### Core Variables

These variables define fundamental operational parameters for the M3TAL core components.

#### `DASHBOARD_PORT`

*   **Description**: The internal port on which the M3TAL Dashboard container (`m3tal-dashboard`) listens. This port is typically exposed via Traefik or a host mapping.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Components Using It**: Dashboard container (`m3tal-dashboard`), Traefik gateway (for routing to the dashboard).

#### `HTTP_PORT`

*   **Description**: The port on which the M3TAL API daemon (`m3tal-api.service`) listens for HTTP requests. This is the primary interface for programmatic interaction with the M3TAL backend.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Components Using It**: API daemon (`m3tal-api.service`), Traefik gateway (for routing `api.DOMAIN` to the daemon), Dashboard container (to communicate with the API).

#### `LOG_LEVEL`

*   **Description**: Sets the logging verbosity for the M3TAL API daemon. Useful for debugging and monitoring.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Components Using It**: API daemon (`m3tal-api.service`)
*   **Notes**: Common values include `debug`, `info`, `warn`, `error`, `fatal`, `panic`.

### Authentication Variables

These variables are crucial for securing access to the M3TAL Dashboard and API.

#### `DASHBOARD_SECRET`

*   **Description**: A strong, randomly generated secret key used by the Dashboard for session management, token signing, and cryptographic operations.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_very_long_random_string_of_characters_1234567890abcdef`
*   **Components Using It**: Dashboard container (`m3tal-dashboard`)
*   **Notes**: This variable is **auto-generated on the first `m3tal init`** for security. Users should **NOT set it manually** unless they are performing a secure secret rotation. It is critical for the security of your M3TAL Dashboard.

#### `API_TOKEN`

*   **Description**: A strong, randomly generated token used for authenticating requests to the M3TAL API daemon. Both the CLI and the Dashboard use this token for API interactions.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `another_long_random_alphanumeric_string_abcdef1234567890`
*   **Components Using It**: CLI binary (`/usr/bin/m3tal`), API daemon (`m3tal-api.service`), Dashboard container (`m3tal-dashboard`)
*   **Notes**: This variable is **auto-generated on the first `m3tal init`** for security. Users should **NOT set it manually** unless they are performing a secure token rotation. It provides programmatic access to your M3TAL instance.

#### `ADMIN_PASSWORD`

*   **Description**: The initial password for the default 'admin' user of the M3TAL Dashboard. After initial setup, user credentials are managed via `/docker/users.json` (managed by `m3tal dashpass`).
*   **Default Value**: `admin_pass`
*   **Example Value**: `SuperSecurePassword123!`
*   **Components Using It**: Dashboard container (`m3tal-dashboard`)
*   **Notes**: It is highly recommended to change this immediately after setup using `m3tal dashpass`.

### Network Variables

These variables configure how M3TAL components communicate over Docker networks and with the host.

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Controls how the Dashboard container's port (`DASHBOARD_PORT`) is exposed on the host.
    *   `local`: Binds the port only to `127.0.0.1` (localhost), making it accessible only from the host machine.
    *   `all`: Binds the port to `0.0.0.0`, making it accessible from any network interface on the host.
*   **Default Value**: `local`
*   **Example Value**: `all`
*   **Components Using It**: Dashboard container (`m3tal-dashboard`) (via `m3tal-compose.yml`).

#### `NETWORK_NAME`

*   **Description**: The name of the Docker network that connects all M3TAL containers, allowing them to communicate securely and efficiently. User stacks placed in `/docker/` should also join this network.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal-prod-network`
*   **Components Using It**: All Docker containers deployed by M3TAL (Dashboard, Traefik, Cloudflared, user stacks).

#### `LOCAL_IP`

*   **Description**: Specifies the host machine's local IP address. This is essential for services running inside Docker containers to correctly reference services running directly on the host using `host.docker.internal`.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100`
*   **Components Using It**: Traefik gateway (for routing to the API daemon and Dashboard), Dashboard container, any custom user stacks needing to access host services.

### Storage Variables

These variables define the filesystem paths for persistent data used by M3TAL.

#### `STATE_DIR`

*   **Description**: The directory where the M3TAL API daemon stores its SQLite state database (`state.db`). This database holds critical system state.
*   **Default Value**: `./state` (for development/template)
*   **Example Value**: `/var/lib/m3tal`
*   **Components Using It**: API daemon (`m3tal-api.service`)
*   **Notes**: The canonical production path for the SQLite state database is `/var/lib/m3tal/state.db`. `m3tal init` will set up appropriate paths.

#### `BASE_STORAGE_PATH`

*   **Description**: The root directory for all M3TAL-managed data, including media, configuration, and downloads. All other `_PATH` variables are typically subdirectories of this base path.
*   **Default Value**: `./data` (for development/template)
*   **Example Value**: `/mnt/m3tal-data`
*   **Components Using It**: API daemon, Dashboard container, user stacks that require persistent storage.
*   **Notes**: In production deployments, this defaults to `/mnt` to leverage dedicated storage mounts. Ensure this path has appropriate permissions for the Docker user.

#### `MEDIA_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for media files managed by M3TAL or user applications.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/m3tal-data/media`
*   **Components Using It**: API daemon, user stacks interacting with media.

#### `CONFIG_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` for configuration files used by M3TAL components or user applications.
*   **Default Value**: `./data/config`
*   **Example Value**: `/mnt/m3tal-data/config`
*   **Components Using It**: API daemon, user stacks requiring persistent configuration.

#### `DOWNLOADS_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for downloaded content.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/m3tal-data/downloads`
*   **Components Using It**: API daemon, user stacks managing downloads.

### Traefik Variables

These variables configure the Traefik reverse proxy gateway, which handles incoming HTTP traffic and routes it to M3TAL services.

#### `DOMAIN`

*   **Description**: The primary domain name for your M3TAL deployment. Setting this variable enables Traefik to create public routes like `dash.DOMAIN` for the Dashboard and `api.DOMAIN` for the API daemon.
*   **Default Value**: `localhost`
*   **Example Value**: `my-m3tal-server.com`
*   **Components Using It**: Traefik gateway (via `routing-compose.yml` and `/docker/dynamic/api.yml`), CLI (for generating service URLs).
*   **Notes**: Crucial for exposing M3TAL services to the internet or a local network via a friendly domain name. If `localhost`, services are accessible directly on their ports or via Traefik on `http://localhost`.

#### `TRAEFIK_WEB_PORT`

*   **Description**: The port on which the Traefik gateway listens for incoming HTTP traffic. This is typically port 80.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Components Using It**: Traefik gateway (`routing-compose.yml`)

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The port on which the Traefik gateway listens for incoming HTTPS traffic. This is typically port 443.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Components Using It**: Traefik gateway (`routing-compose.yml`)
*   **Notes**: While defined, full HTTPS setup often requires additional configuration (e.g., certificates, Let's Encrypt integration) beyond basic Traefik deployment.

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port on which the Traefik web UI (dashboard) container listens. This port is mapped to `127.0.0.1:8081` on the host by default, making the Traefik dashboard accessible at `http://localhost:8081`.
*   **Default Value**: `8080`
*   **Example Value**: `8080` (container's internal port)
*   **Components Using It**: Traefik gateway (`routing-compose.yml`)
*   **Notes**: The `routing-compose.yml` maps the container's `TRAEFIK_DASHBOARD_PORT` to `127.0.0.1:8081` on the host.

### VPN Variables

These variables are used to configure the optional VPN service, if deployed.

#### `VPN_USER`

*   **Description**: The username for authenticating with the VPN service provided by M3TAL, if enabled.
*   **Default Value**: `user`
*   **Example Value**: `myvpnuser`
*   **Components Using It**: VPN container (if deployed)

#### `VPN_PASSWORD`

*   **Description**: The password for authenticating with the VPN service provided by M3TAL, if enabled.
*   **Default Value**: `password`
*   **Example Value**: `MyStrongVPNPass!`
*   **Components Using It**: VPN container (if deployed)

### System Variables

These variables control system-wide settings, user/group IDs, and debugging.

#### `PUID`

*   **Description**: Specifies the User ID (UID) that containers requiring specific file permissions should run as. This ensures proper file ownership and access when containers interact with host volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Components Using It**: Many Docker containers, especially those mounting persistent volumes from the host (e.g., Dashboard, user stacks).
*   **Notes**: Typically corresponds to the UID of the user who owns the `BASE_STORAGE_PATH` on the host.

#### `PGID`

*   **Description**: Specifies the Group ID (GID) that containers requiring specific file permissions should run as. This ensures proper file group ownership and access when containers interact with host volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Components Using It**: Many Docker containers, especially those mounting persistent volumes from the host (e.g., Dashboard, user stacks).
*   **Notes**: Typically corresponds to the GID of the group that owns the `BASE_STORAGE_PATH` on the host.

#### `TZ`

*   **Description**: Sets the timezone for all M3TAL services and containers. This ensures consistent timestamps and time-based operations across the ecosystem. Values follow the [IANA Time Zone Database](https://www.iana.org/time-zones) format.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`, `Asia/Tokyo`
*   **Components Using It**: API daemon, Dashboard container, all Docker containers.

#### `DEBUG_MODE`

*   **Description**: Enables or disables additional debugging output and potentially development-specific features across M3TAL components.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Components Using It**: API daemon, Dashboard container, CLI.

#### `METRICS_ENABLED`

*   **Description**: Controls whether the M3TAL API daemon exposes Prometheus-compatible metrics on a dedicated endpoint.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Components Using It**: API daemon (`m3tal-api.service`)