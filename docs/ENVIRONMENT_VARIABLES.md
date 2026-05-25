# Environment Variables Reference

All M3TAL system and stack-specific environment variables are consolidated and read from the primary configuration file: `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (managed by `m3tal up`) utilize this file via the `--env-file` flag to ensure consistent configuration across the entire ecosystem.

You can manage these variables using the `m3tal config wizard` command for an interactive setup, or directly with `m3tal config set KEY value`.

## Quick Reference

| Variable Name           | Category  | Default Value             | Example Value               | Components                                  |
|-------------------------|-----------|---------------------------|-----------------------------|---------------------------------------------|
| `DASHBOARD_PORT`        | Core      | `8082`                    | `8082`                      | `m3tal-dashboard`, `CLI binary`             |
| `DASHBOARD_EXPOSE_MODE` | Core      | `local`                   | `traefik`                   | `CLI binary`                                |
| `HTTP_PORT`             | Core      | `8080`                    | `8080`                      | `API daemon`                                |
| `STATE_DIR`             | Core      | `./state`                 | `/var/lib/m3tal`            | `API daemon`, `m3tal-dashboard`             |
| `LOG_LEVEL`             | Core      | `info`                    | `debug`                     | `API daemon`, `m3tal-dashboard`             |
| `DASHBOARD_SECRET`      | Auth      | `change_me_immediately`   | `a_very_long_secret_key`    | `m3tal-dashboard`                           |
| `API_TOKEN`             | Auth      | `change_me_api_token`     | `secure_api_token_value`    | `API daemon`                                |
| `ADMIN_PASSWORD`        | Auth      | `admin_pass`              | `myStrongDashboardPassword` | `m3tal-dashboard`                           |
| `NETWORK_NAME`          | Network   | `m3tal`                   | `my_custom_network`         | `CLI binary`, `Docker Compose`              |
| `LOCAL_IP`              | Network   | `127.0.0.1`               | `192.168.1.50`              | `API daemon`, `Traefik gateway`             |
| `DOMAIN`                | Traefik   | `localhost`               | `example.com`               | `Traefik gateway`, `CLI binary`             |
| `TRAEFIK_WEB_PORT`      | Traefik   | `80`                      | `80`                        | `Traefik gateway`                           |
| `TRAEFIK_WEBHTTPS_PORT` | Traefik   | `443`                     | `443`                       | `Traefik gateway`                           |
| `TRAEFIK_DASHBOARD_PORT`| Traefik   | `8080`                    | `8080`                      | `Traefik gateway`                           |
| `VPN_USER`              | VPN       | `user`                    | `myvpnuser`                 | `Cloudflared`                               |
| `VPN_PASSWORD`          | VPN       | `password`                | `myvpnpass`                 | `Cloudflared`                               |
| `BASE_STORAGE_PATH`     | Storage   | `./data`                  | `/mnt/m3tal_data`           | `m3tal-dashboard`, `Docker Compose`         |
| `MEDIA_PATH`            | Storage   | `./data/media`            | `/mnt/media`                | `Docker Compose`                            |
| `CONFIG_PATH`           | Storage   | `./data/config`           | `/etc/m3tal/config`         | `m3tal-dashboard`, `Docker Compose`         |
| `DOWNLOADS_PATH`        | Storage   | `./data/downloads`        | `/mnt/downloads`            | `Docker Compose`                            |
| `PUID`                  | System    | `1000`                    | `1001`                      | `Docker Compose`                            |
| `PGID`                  | System    | `1000`                    | `1001`                      | `Docker Compose`                            |
| `TZ`                    | System    | `America/Denver`          | `Europe/London`             | `Docker Compose`                            |
| `DEBUG_MODE`            | System    | `false`                   | `true`                      | `API daemon`, `m3tal-dashboard`             |
| `METRICS_ENABLED`       | System    | `true`                    | `false`                     | `API daemon`                                |

---

## Detailed Variable Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL ecosystem's operation.

#### `DASHBOARD_PORT`
- **Description**: Specifies the port on which the `m3tal-dashboard` container listens. This port is directly exposed to the host when `DASHBOARD_EXPOSE_MODE` is set to `local`.
- **Default Value**: `8082`
- **Example Value**: `8082`
- **Used by Component(s)**: `m3tal-dashboard`, `CLI binary` (for local mode port binding in `m3tal-compose.local.yml`)

#### `DASHBOARD_EXPOSE_MODE`
- **Description**: Determines how the `m3tal-dashboard` is made accessible.
  - `local`: The dashboard port is directly bound to the host, accessible via `http://HOST_IP:DASHBOARD_PORT`. No Traefik required. Ideal for LAN-only setups.
  - `traefik`: The dashboard is exposed via the Traefik reverse proxy using Docker labels. It becomes accessible via `http://dash.DOMAIN`. Requires Traefik to be running.
- **Default Value**: `local`
- **Example Value**: `traefik`
- **Used by Component(s)**: `CLI binary` (to select the appropriate Docker Compose override file: `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`)

#### `HTTP_PORT`
- **Description**: Defines the port on which the `m3tal-api.service` (Go API daemon) listens for incoming requests. This port is typically only accessible from the host system or other containers via `host.docker.internal`.
- **Default Value**: `8080`
- **Example Value**: `8080`
- **Used by Component(s)**: `API daemon`

#### `STATE_DIR`
- **Description**: Specifies the directory used by the `API daemon` to store its SQLite state database (`state.db`) and by the `m3tal-dashboard` for its internal configuration and user data (`users.json`).
- **Default Value**: `./state` (This is typically relative to the current working directory during development or simple setups.)
- **Example Value**: `/var/lib/m3tal` (In production deployments, this defaults to `/var/lib/m3tal` for system-managed persistence.)
- **Used by Component(s)**: `API daemon`, `m3tal-dashboard`

#### `LOG_LEVEL`
- **Description**: Controls the verbosity of log output for M3TAL core components.
  - `debug`: Most verbose, useful for troubleshooting.
  - `info`: General informational messages (default).
  - `warn`: Warnings and non-critical errors.
  - `error`: Only critical error messages.
- **Default Value**: `info`
- **Example Value**: `debug`
- **Used by Component(s)**: `API daemon`, `m3tal-dashboard`

### Authentication & Security

These variables manage access and security for M3TAL components.

#### `DASHBOARD_SECRET`
- **Description**: A long, random string used by the `m3tal-dashboard` for session management, cookie signing, and other cryptographic operations. It is crucial for dashboard security.
- **Default Value**: `change_me_immediately`
- **Example Value**: `a_very_long_and_random_string_of_characters_for_security`
- **Used by Component(s)**: `m3tal-dashboard`
- **Important Note**: This variable is **auto-generated** on the first `m3tal init` run. Users should generally **NOT** set it manually unless performing a secret rotation.

#### `API_TOKEN`
- **Description**: The bearer token required to authenticate requests to the `m3tal-api.service`. This token protects your M3TAL API from unauthorized access.
- **Default Value**: `change_me_api_token`
- **Example Value**: `another_super_strong_random_token_for_api_access_xyz`
- **Used by Component(s)**: `API daemon`
- **Important Note**: This variable is **auto-generated** on the first `m3tal init` run. Users should generally **NOT** set it manually unless performing a token rotation.

#### `ADMIN_PASSWORD`
- **Description**: Sets the initial password for the default `admin` user within the `m3tal-dashboard`'s `users.json` file. It's used when the `m3tal dashpass` command creates or resets the admin user.
- **Default Value**: `admin_pass`
- **Example Value**: `myStrongDashboardPassword`
- **Used by Component(s)**: `m3tal-dashboard` (specifically when user management operations are performed through `m3tal dashpass`)

### Network Configuration

Variables related to Docker networking and host IP addresses.

#### `NETWORK_NAME`
- **Description**: Defines the name of the Docker network that M3TAL and all its managed Docker Compose stacks will utilize. This allows containers from different stacks to communicate with each other.
- **Default Value**: `m3tal`
- **Example Value**: `my_custom_network`
- **Used by Component(s)**: `CLI binary`, `Docker Compose`

#### `LOCAL_IP`
- **Description**: The IP address of the host machine where M3TAL is running. This is used by some services, particularly Traefik, to route traffic to the `API daemon` which runs directly on the host, via `http://host.docker.internal:PORT`.
- **Default Value**: `127.0.0.1`
- **Example Value**: `192.168.1.50`
- **Used by Component(s)**: `API daemon` (implicitly via `host.docker.internal`), `Traefik gateway` (for routing to host-based services)

### Traefik Gateway

Configuration options for the Traefik reverse proxy.

#### `DOMAIN`
- **Description**: The base domain name used by Traefik to define routing rules for M3TAL's core services. When set, Traefik will expose the Dashboard at `dash.DOMAIN` and the M3TAL API at `api.DOMAIN`.
- **Default Value**: `localhost`
- **Example Value**: `example.com`
- **Used by Component(s)**: `Traefik gateway`, `CLI binary` (for generating Traefik dynamic configuration and dashboard labels)
- **Important Note**: Setting this variable enables the `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.

#### `TRAEFIK_WEB_PORT`
- **Description**: The port on the host machine where the `Traefik gateway` listens for incoming HTTP traffic.
- **Default Value**: `80`
- **Example Value**: `80`
- **Used by Component(s)**: `Traefik gateway`

#### `TRAEFIK_WEBHTTPS_PORT`
- **Description**: The port on the host machine where the `Traefik gateway` listens for incoming HTTPS traffic.
- **Default Value**: `443`
- **Example Value**: `443`
- **Used by Component(s)**: `Traefik gateway`

#### `TRAEFIK_DASHBOARD_PORT`
- **Description**: The internal port within the `Traefik gateway` container where Traefik's own management dashboard is exposed. By default, this is mapped to `127.0.0.1:8081` on the host, making it only locally accessible.
- **Default Value**: `8080`
- **Example Value**: `8080`
- **Used by Component(s)**: `Traefik gateway`

### VPN / Cloudflare Tunnel

Variables for configuring optional VPN-like access, such as Cloudflare Tunnels.

#### `VPN_USER`
- **Description**: The username credential used by the `Cloudflared` container for establishing a secure tunnel, if applicable.
- **Default Value**: `user`
- **Example Value**: `cloudflare_tunnel_user`
- **Used by Component(s)**: `Cloudflared` (if configured to use credentials for tunnel setup)

#### `VPN_PASSWORD`
- **Description**: The password credential used by the `Cloudflared` container for establishing a secure tunnel, if applicable.
- **Default Value**: `password`
- **Example Value**: `cloudflare_tunnel_pass`
- **Used by Component(s)**: `Cloudflared` (if configured to use credentials for tunnel setup)

### Storage Paths

Root and subdirectory paths for persistent data storage.

#### `BASE_STORAGE_PATH`
- **Description**: The root directory on the host filesystem where M3TAL-managed applications store their persistent data, including media, configurations, and downloads.
- **Default Value**: `./data` (This is a common default for local development or simple setups.)
- **Example Value**: `/mnt/m3tal_data` (In production deployments, this variable defaults to `/mnt` to align with typical Linux mounting points for data volumes.)
- **Used by Component(s)**: `m3tal-dashboard`, `Docker Compose` (for mounting volumes into various user-managed and core containers)

#### `MEDIA_PATH`
- **Description**: A subdirectory within `BASE_STORAGE_PATH` designated for storing media files managed by M3TAL-enabled applications.
- **Default Value**: `./data/media`
- **Example Value**: `${BASE_STORAGE_PATH}/media` (e.g., `/mnt/media`)
- **Used by Component(s)**: `Docker Compose` (for mounting media volumes into user-managed containers)

#### `CONFIG_PATH`
- **Description**: A subdirectory within `BASE_STORAGE_PATH` (or an absolute path) where various application-specific configuration files (beyond the main `.env` file) can be stored and mounted into containers. The dashboard uses this to mount internal configs like `users.json`.
- **Default Value**: `./data/config` (This is a common default for local development or simple setups.)
- **Example Value**: `/etc/m3tal/config` (In production deployments, this may point to a more canonical system configuration path like `/etc/m3tal/config` or `/mnt/config` as seen in compose files.)
- **Used by Component(s)**: `m3tal-dashboard`, `Docker Compose` (for mounting configuration volumes into user-managed containers)

#### `DOWNLOADS_PATH`
- **Description**: A subdirectory within `BASE_STORAGE_PATH` dedicated to storing files downloaded by M3TAL-enabled applications.
- **Default Value**: `./data/downloads`
- **Example Value**: `${BASE_STORAGE_PATH}/downloads` (e.g., `/mnt/downloads`)
- **Used by Component(s)**: `Docker Compose` (for mounting download volumes into user-managed containers)

### System & Runtime

General system-level settings for containers.

#### `PUID`
- **Description**: The User ID (UID) that containers should use when running processes inside, primarily for managing file ownership and permissions on mounted volumes. This ensures consistency between container processes and host file system permissions.
- **Default Value**: `1000` (Typically the first non-root user ID on Linux systems.)
- **Example Value**: `1001`
- **Used by Component(s)**: `Docker Compose` (for setting `user` within containers, including `m3tal-dashboard` and user-managed stacks)

#### `PGID`
- **Description**: The Group ID (GID) that containers should use when running processes inside, primarily for managing file group ownership and permissions on mounted volumes. This ensures consistency between container processes and host file system permissions.
- **Default Value**: `1000` (Typically the primary group ID for the user with UID 1000.)
- **Example Value**: `1001`
- **Used by Component(s)**: `Docker Compose` (for setting `user` within containers, including `m3tal-dashboard` and user-managed stacks)

#### `TZ`
- **Description**: Sets the timezone for all M3TAL-managed Docker containers. This ensures that logs, timestamps, and scheduled tasks within containers reflect the correct local time.
- **Default Value**: `America/Denver`
- **Example Value**: `Europe/London`
- **Used by Component(s)**: `Docker Compose` (as an environment variable for `m3tal-dashboard` and user-managed stacks)

#### `DEBUG_MODE`
- **Description**: A boolean flag to enable or disable extensive debugging output and potentially expose development-specific endpoints or features within M3TAL core components.
- **Default Value**: `false`
- **Example Value**: `true`
- **Used by Component(s)**: `API daemon`, `m3tal-dashboard`

#### `METRICS_ENABLED`
- **Description**: A boolean flag to control whether the `API daemon` collects and exposes application performance metrics. Disabling this can reduce overhead if metrics collection is not desired.
- **Default Value**: `true`
- **Example Value**: `false`
- **Used by Component(s)**: `API daemon`