Greetings, M3TAL user. DocSmith, the M3TAL Ecosystem Documentation Architect, here to guide you through the intricate world of M3TAL's environment variables. This document serves as your definitive reference for all configuration options.

---

# Environment Variables Reference

All M3TAL environment variables are centrally managed and read from the `/etc/m3tal/.env` file. This ensures consistent configuration across both the `m3tal` CLI binary and all Docker Compose stacks deployed within your M3TAL ecosystem.

You can manage these variables using the `m3tal config wizard` for interactive setup, or directly via `m3tal config set KEY value` for specific updates.

## Quick Reference Table

| Name                  | Description (Brief)                                                              | Default                   | Components                      |
| :-------------------- | :------------------------------------------------------------------------------- | :------------------------ | :------------------------------ |
| `DASHBOARD_PORT`      | Port for M3TAL Dashboard container.                                              | `8082`                    | Dashboard, Traefik              |
| `DASHBOARD_EXPOSE_MODE` | Controls network exposure of the Dashboard.                                      | `local`                   | Dashboard                       |
| `HTTP_PORT`           | Port for M3TAL API daemon.                                                       | `8080`                    | API Daemon, Traefik             |
| `STATE_DIR`           | Directory for API daemon's state database.                                       | `./state`                 | API Daemon                      |
| `LOG_LEVEL`           | Logging verbosity for CLI and API.                                               | `info`                    | CLI, API Daemon                 |
| `DASHBOARD_SECRET`    | Secret key for Dashboard session management.                                     | `change_me_immediately`   | Dashboard                       |
| `API_TOKEN`           | Authentication token for API daemon access.                                      | `change_me_api_token`     | CLI, API Daemon, Dashboard      |
| `ADMIN_PASSWORD`      | Initial password for Dashboard admin user.                                       | `admin_pass`              | Dashboard                       |
| `NETWORK_NAME`        | Name of the Docker network.                                                      | `m3tal`                   | API Daemon, Compose Stacks, Traefik |
| `LOCAL_IP`            | Host machine IP within Docker network context.                                   | `127.0.0.1`               | Compose Stacks, Traefik         |
| `DOMAIN`              | Primary domain for M3TAL services and Traefik routing.                           | `localhost`               | Traefik, Dashboard              |
| `VPN_USER`            | Username for VPN client configurations.                                          | `user`                    | User VPN Containers             |
| `VPN_PASSWORD`        | Password for VPN client configurations.                                          | `password`                | User VPN Containers             |
| `BASE_STORAGE_PATH`   | Base directory for persistent data.                                              | `./data`                  | Compose Stacks                  |
| `MEDIA_PATH`          | Subdirectory for media files.                                                    | `./data/media`            | Compose Stacks                  |
| `CONFIG_PATH`         | Subdirectory for application configuration.                                      | `./data/config`           | Compose Stacks                  |
| `DOWNLOADS_PATH`      | Subdirectory for downloaded content.                                             | `./data/downloads`        | Compose Stacks                  |
| `PUID`                | User ID (UID) for container file permissions.                                    | `1000`                    | Compose Stacks                  |
| `PGID`                | Group ID (GID) for container file permissions.                                   | `1000`                    | Compose Stacks                  |
| `TZ`                  | Timezone for M3TAL services.                                                     | `America/Denver`          | Compose Stacks, API Daemon      |
| `TRAEFIK_WEB_PORT`    | Host port for Traefik's HTTP entry point.                                        | `80`                      | Traefik Gateway                 |
| `TRAEFIK_WEBHTTPS_PORT` | Host port for Traefik's HTTPS entry point.                                       | `443`                     | Traefik Gateway                 |
| `TRAEFIK_DASHBOARD_PORT` | Internal port Traefik uses for its dashboard.                                    | `8080`                    | Traefik Gateway                 |
| `DEBUG_MODE`          | Enables verbose debugging output.                                                | `false`                   | API Daemon, CLI                 |
| `METRICS_ENABLED`     | Enables Prometheus-compatible metrics from API.                                  | `true`                    | API Daemon                      |

## Detailed Environment Variable Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system.

#### `STATE_DIR`

*   **Description**: The directory where the M3TAL API daemon stores its persistent state, including the SQLite database (`state.db`). While the default is relative (`./state`), in production systemd deployments, it's typically configured to an absolute path like `/var/lib/m3tal`.
*   **Default Value**: `./state`
*   **Example Value**: `/var/lib/m3tal`
*   **Components**: API daemon (`m3tal-api.service`)

#### `LOG_LEVEL`

*   **Description**: Sets the logging verbosity for the M3TAL API daemon and the CLI. Higher verbosity (e.g., `debug`) provides more detailed output, useful for troubleshooting.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Components**: CLI binary, API daemon (`m3tal-api.service`)

### Authentication

These variables are critical for securing access to M3TAL services.

#### `DASHBOARD_SECRET`

*   **Description**: A unique, cryptographically secure secret key used by the M3TAL Dashboard for session management and other security-sensitive operations.
    *   **Note**: This value is **auto-generated** on the first `m3tal init` run. Users should **NOT** set this manually unless performing a security rotation.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_long_random_string_of_characters_for_security_12345`
*   **Components**: Dashboard container (`m3tal-dashboard`)

#### `API_TOKEN`

*   **Description**: An authentication token required for direct access to the M3TAL API daemon. This token secures communication between the dashboard, CLI, and the API.
    *   **Note**: This value is **auto-generated** on the first `m3tal init` run. Users should **NOT** set this manually unless performing a security rotation.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `another_secure_random_token_string_for_api_access_67890`
*   **Components**: CLI binary, API daemon (`m3tal-api.service`), Dashboard container (`m3tal-dashboard`)

#### `ADMIN_PASSWORD`

*   **Description**: The initial password for the default `admin` user of the M3TAL Dashboard. After initial setup, user passwords are managed via the `m3tal dashpass` command.
*   **Default Value**: `admin_pass`
*   **Example Value**: `mySuperSecureAdminPassword!`
*   **Components**: Dashboard container (`m3tal-dashboard`) (for `/docker/users.json` credential store)

### Network

Variables defining how M3TAL services communicate over the network.

#### `DASHBOARD_PORT`

*   **Description**: The internal port on which the M3TAL Dashboard container listens for incoming connections.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Components**: Dashboard container (`m3tal-dashboard`), Traefik gateway

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Controls the network exposure of the M3TAL Dashboard. When set to `local`, the dashboard is primarily intended for access via `localhost` or through Traefik using the `DOMAIN` configuration, preventing direct public exposure from the container.
*   **Default Value**: `local`
*   **Example Value**: `local`
*   **Components**: Dashboard container (`m3tal-dashboard`)

#### `HTTP_PORT`

*   **Description**: The port on which the M3TAL API daemon listens for incoming HTTP requests.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Components**: API daemon (`m3tal-api.service`), Traefik gateway

#### `NETWORK_NAME`

*   **Description**: The name of the Docker network created and used by M3TAL to allow its services and user-defined stacks to communicate with each other.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal-proxy-network`
*   **Components**: API daemon (`m3tal-api.service`), Docker Compose stacks, Traefik gateway

#### `LOCAL_IP`

*   **Description**: The IP address that represents the host machine within the Docker network context. This is crucial for Traefik and other containers to communicate with host-bound services like the M3TAL API daemon via `host.docker.internal`.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `172.17.0.1` (common Docker bridge IP)
*   **Components**: Docker Compose stacks (e.g., Traefik's dynamic configuration)

### Storage

These variables define the paths for persistent data storage. All paths are typically relative to `BASE_STORAGE_PATH`.

#### `BASE_STORAGE_PATH`

*   **Description**: The base directory on the host machine where M3TAL services store all persistent data, including media, application configurations, and downloads.
    *   **Note**: While the template defaults to `./data`, in production M3TAL deployments, this variable **defaults to `/mnt`** to align with common server storage practices.
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal-data`
*   **Components**: Docker Compose stacks (for volume mounts)

#### `MEDIA_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for media files (e.g., movies, TV shows, music). This path is typically mounted into media management containers.
*   **Default Value**: `./data/media`
*   **Example Value**: `${BASE_STORAGE_PATH}/media` (e.g., `/mnt/m3tal-data/media`)
*   **Components**: Docker Compose stacks (for specific service volume mounts)

#### `CONFIG_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` used to store persistent application configuration files for M3TAL services and user-deployed applications.
*   **Default Value**: `./data/config`
*   **Example Value**: `${BASE_STORAGE_PATH}/config` (e.g., `/mnt/m3tal-data/config`)
*   **Components**: Docker Compose stacks (for specific service volume mounts)

#### `DOWNLOADS_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` where downloaded content from torrent clients, Usenet downloaders, etc., is stored.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `${BASE_STORAGE_PATH}/downloads` (e.g., `/mnt/m3tal-data/downloads`)
*   **Components**: Docker Compose stacks (for specific service volume mounts)

### Traefik Gateway

Variables specific to the Traefik reverse proxy and its routing capabilities.

#### `DOMAIN`

*   **Description**: The primary domain for your M3TAL services. Setting this variable is essential for Traefik to correctly establish routing rules, enabling access to services like the Dashboard (`dash.DOMAIN`) and the API (`api.DOMAIN`) via friendly URLs.
*   **Default Value**: `localhost`
*   **Example Value**: `m3tal.local` or `my-awesome-server.com`
*   **Components**: Traefik gateway, Dashboard container (via labels)

#### `TRAEFIK_WEB_PORT`

*   **Description**: The host port that Traefik's HTTP entry point binds to, making your services accessible via standard HTTP.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Components**: Traefik gateway (`routing-compose.yml`)

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port that Traefik's HTTPS entry point binds to, enabling secure, encrypted access to your services. Requires additional configuration for SSL certificates.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Components**: Traefik gateway (`routing-compose.yml`)

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port that Traefik is configured to use for its own web dashboard. While the external mapping often defaults to `8081` for local access, this variable defines the port within Traefik's configuration.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Components**: Traefik gateway (`routing-compose.yml`)

### VPN Configuration

Variables for integrating VPN services with M3TAL containers.

#### `VPN_USER`

*   **Description**: The username used for authenticating VPN client containers within the M3TAL ecosystem. This is typically passed as an environment variable to user-defined VPN-enabled service stacks.
*   **Default Value**: `user`
*   **Example Value**: `myvpnuser`
*   **Components**: User-defined VPN client containers

#### `VPN_PASSWORD`

*   **Description**: The password corresponding to the `VPN_USER` for VPN client container authentication.
*   **Default Value**: `password`
*   **Example Value**: `mysecurevpnpass123`
*   **Components**: User-defined VPN client containers

### System Configuration

General system-level variables affecting various M3TAL components.

#### `PUID`

*   **Description**: The User ID (UID) that containers should use when accessing mounted volumes. Setting this to the UID of the user running Docker ensures proper file permissions and avoids permission denied errors for services writing to host paths.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Components**: Docker Compose stacks

#### `PGID`

*   **Description**: The Group ID (GID) that containers should use when accessing mounted volumes. Similar to `PUID`, this ensures containers have the correct group permissions for files and directories on the host.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Components**: Docker Compose stacks

#### `TZ`

*   **Description**: Sets the timezone for all M3TAL containers and the API daemon. This ensures consistent time reporting and scheduling across all services. Use standard IANA timezone database names (e.g., `America/New_York`).
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Components**: Docker Compose stacks, API daemon (`m3tal-api.service`)

#### `DEBUG_MODE`

*   **Description**: A boolean flag that, when set to `true`, enables verbose debugging output for the API daemon and potentially other services. Useful for diagnostics and development.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Components**: API daemon (`m3tal-api.service`), CLI binary

#### `METRICS_ENABLED`

*   **Description**: A boolean flag that controls whether the M3TAL API daemon exposes Prometheus-compatible metrics. When enabled, metrics can be scraped by monitoring systems.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Components**: API daemon (`m3tal-api.service`)