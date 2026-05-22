# Environment Variables Reference

All M3TAL environment variables are centrally managed and read from the `/etc/m3tal/.env` file. This single source of truth is utilized by both the M3TAL CLI and all Docker Compose stacks through the `--env-file` option, ensuring consistent configuration across the entire ecosystem.

You can manage these variables using the `m3tal config wizard` for interactive setup, or `m3tal config set KEY value` for direct updates.

---

## Quick Reference

| Variable Name             | Default Value             | Component(s) Used By                                |
| :------------------------ | :------------------------ | :-------------------------------------------------- |
| `DASHBOARD_PORT`          | `8082`                    | M3TAL Dashboard, CLI (`m3tal dash up`)              |
| `DASHBOARD_EXPOSE_MODE`   | `local`                   | CLI (`m3tal dash up`)                               |
| `HTTP_PORT`               | `8080`                    | M3TAL API Daemon, Traefik                           |
| `STATE_DIR`               | `./state`                 | M3TAL Dashboard                                     |
| `LOG_LEVEL`               | `info`                    | M3TAL API Daemon                                    |
| `DASHBOARD_SECRET`        | `change_me_immediately`   | M3TAL Dashboard                                     |
| `API_TOKEN`               | `change_me_api_token`     | CLI, M3TAL API Daemon, M3TAL Dashboard              |
| `ADMIN_PASSWORD`          | `admin_pass`              | M3TAL Dashboard (initial user)                      |
| `NETWORK_NAME`            | `m3tal`                   | All Docker Compose stacks                           |
| `LOCAL_IP`                | `127.0.0.1`               | M3TAL Dashboard, Traefik                            |
| `DOMAIN`                  | `localhost`               | Traefik, M3TAL Dashboard                            |
| `VPN_USER`                | `user`                    | Custom VPN stacks                                   |
| `VPN_PASSWORD`            | `password`                | Custom VPN stacks                                   |
| `BASE_STORAGE_PATH`       | `./data`                  | M3TAL Dashboard, all Docker Compose stacks          |
| `MEDIA_PATH`              | `./data/media`            | M3TAL Dashboard, all Docker Compose stacks          |
| `CONFIG_PATH`             | `./data/config`           | M3TAL Dashboard, all Docker Compose stacks          |
| `DOWNLOADS_PATH`          | `./data/downloads`        | All Docker Compose stacks                           |
| `PUID`                    | `1000`                    | All Docker Compose containers                       |
| `PGID`                    | `1000`                    | All Docker Compose containers                       |
| `TZ`                      | `America/Denver`          | All Docker Compose containers                       |
| `TRAEFIK_WEB_PORT`        | `80`                      | Traefik Gateway                                     |
| `TRAEFIK_WEBHTTPS_PORT`   | `443`                     | Traefik Gateway                                     |
| `TRAEFIK_DASHBOARD_PORT`  | `8080`                    | Traefik Gateway                                     |
| `DEBUG_MODE`              | `false`                   | M3TAL API Daemon                                    |
| `METRICS_ENABLED`         | `true`                    | M3TAL API Daemon                                    |

---

## Detailed Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system's operation.

#### `DASHBOARD_PORT`

*   **Description**: Specifies the internal port on which the `m3tal-dashboard` container listens. This port is also used for direct host-port mapping when `DASHBOARD_EXPOSE_MODE` is set to `local`.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Used By**: `m3tal-dashboard` container, `m3tal dash up` command (local mode mapping).

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Determines how the M3TAL Dashboard is exposed to the network.
    *   `local`: Direct port binding (`DASHBOARD_PORT`) for LAN access (`http://HOST_IP:8082`). No Traefik required.
    *   `traefik`: Exposes the dashboard via Traefik using the `dash.<DOMAIN>` hostname. Requires Traefik to be running.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Used By**: CLI (`m3tal dash up`) to select the appropriate Docker Compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`).

#### `HTTP_PORT`

*   **Description**: The port on which the `m3tal-api.service` (Go API daemon) listens for incoming connections. This port is primarily for internal host-local communication and is exposed externally via Traefik when `DOMAIN` is set.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: M3TAL API daemon, Traefik (for routing `api.<DOMAIN>` to `http://host.docker.internal:8080`).

#### `STATE_DIR`

*   **Description**: The internal path within the `m3tal-dashboard` container where the SQLite state database (`state.db`) and other configuration files (like `users.json`) are stored. On the host, this typically maps to `${CONFIG_PATH}/m3tal/state`.
*   **Default Value**: `./state`
*   **Example Value**: `/docker/state` (inside container)
*   **Used By**: `m3tal-dashboard` container.

### Authentication

These variables are critical for securing access to M3TAL components.

#### `DASHBOARD_SECRET`

*   **Description**: A unique, cryptographically secure secret key used by the `m3tal-dashboard` for session management, cookie signing, and other security-related operations.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `aSuperStrongRandomStringOfEntropy`
*   **Used By**: `m3tal-dashboard` container.
*   **Note**: This value is **auto-generated** on the first `m3tal init` run. You should **NOT** set this manually unless performing a secret rotation.

#### `API_TOKEN`

*   **Description**: An authentication token used to secure communication with the `m3tal-api.service`. It's used by the CLI and any other components (e.g., the Dashboard) that interact directly with the API.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `anotherVerySecureLongTokenHere`
*   **Used By**: CLI binary, M3TAL API daemon, `m3tal-dashboard` (when communicating with the API).
*   **Note**: This value is **auto-generated** on the first `m3tal init` run. You should **NOT** set this manually unless performing a token rotation.

#### `ADMIN_PASSWORD`

*   **Description**: The default password for the initial `admin` user created in the M3TAL Dashboard's `users.json` file. It's highly recommended to change this immediately after first login.
*   **Default Value**: `admin_pass`
*   **Example Value**: `MyNewSecurePassword123!`
*   **Used By**: M3TAL Dashboard (initial user setup).

### Network Configuration

Variables related to Docker networking and host IP addresses.

#### `NETWORK_NAME`

*   **Description**: The name of the Docker network that all M3TAL-managed containers (including Traefik and the Dashboard) join. This facilitates inter-container communication, especially important for Traefik to discover services. The `proxy` network is the external name used in compose files.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal_network`
*   **Used By**: All Docker Compose stacks (`m3tal-compose.yml`, `routing-compose.yml`, user-defined stacks).

#### `LOCAL_IP`

*   **Description**: The IP address of the host machine, used primarily for internal routing within Docker containers, such as mapping `host.docker.internal` or directly addressing host services like the M3TAL API daemon.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100` (for direct LAN access from containers)
*   **Used By**: M3TAL Dashboard (implicitly for `GO_API_URL` when `host.docker.internal` resolves to this), Traefik (for routing to `http://host.docker.internal:8080`).

### Storage Paths

These variables define the locations on the host filesystem where M3TAL stores its persistent data.

#### `BASE_STORAGE_PATH`

*   **Description**: The root directory on the host machine where all M3TAL-related persistent data (media, configurations, downloads, etc.) is stored. Other path variables are typically relative to this.
*   **Default Value**: `./data` (relative to the `m3tal-api.service` working directory)
*   **Example Value**: `/srv/m3tal`
*   **Used By**: `m3tal-dashboard` container, other Docker Compose stacks mounting volumes.
*   **Note**: In production deployments, this defaults to `/mnt` to leverage dedicated storage mounts, rather than `./data`.

#### `MEDIA_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for media files and content.
*   **Default Value**: `./data/media`
*   **Example Value**: `${BASE_STORAGE_PATH}/media`
*   **Used By**: `m3tal-dashboard` container, other Docker Compose stacks requiring media storage.

#### `CONFIG_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files, including the `state.db` and `users.json` files.
*   **Default Value**: `./data/config`
*   **Example Value**: `${BASE_STORAGE_PATH}/config`
*   **Used By**: `m3tal-dashboard` container (for mounting `/docker/state/config`), other Docker Compose stacks storing configuration.

#### `DOWNLOADS_PATH`

*   **Description**: The subdirectory within `BASE_STORAGE_PATH` intended for downloaded content from various services.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `${BASE_STORAGE_PATH}/downloads`
*   **Used By**: Any Docker Compose stack that performs downloads.

### Traefik Gateway

Variables specific to the Traefik reverse proxy and its routing behavior.

#### `DOMAIN`

*   **Description**: The primary domain name for your M3TAL services. Setting this enables Traefik to route `dash.<DOMAIN>` to the Dashboard and `api.<DOMAIN>` to the M3TAL API daemon. Traefik rules are dynamically configured based on this.
*   **Default Value**: `localhost`
*   **Example Value**: `mym3talserver.com`
*   **Used By**: Traefik Gateway (`routing-compose.yml`), M3TAL Dashboard (when `DASHBOARD_EXPOSE_MODE=traefik`).

#### `TRAEFIK_WEB_PORT`

*   **Description**: The port Traefik listens on for incoming HTTP (`web`) traffic. This is typically mapped to port 80 on the host.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Used By**: `traefik` container in `routing-compose.yml`.

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The port Traefik listens on for incoming HTTPS (`websecure`) traffic. This is typically mapped to port 443 on the host, often used with TLS configurations (e.g., Let's Encrypt).
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Used By**: `traefik` container in `routing-compose.yml` (if HTTPS entrypoint is enabled).

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port for Traefik's own management dashboard. Mapped to `127.0.0.1:8081` on the host by default, making it accessible only locally via `http://localhost:8081`.
*   **Default Value**: `8080` (internal container port)
*   **Example Value**: `8080`
*   **Used By**: `traefik` container in `routing-compose.yml`.

### VPN Configuration

Placeholder variables for potential VPN integration.

#### `VPN_USER`

*   **Description**: Username for a generic VPN service. This variable is a placeholder and would be used by custom user-defined VPN stacks.
*   **Default Value**: `user`
*   **Example Value**: `m3talvpnuser`
*   **Used By**: Custom VPN stacks (e.g., OpenVPN, WireGuard client containers).

#### `VPN_PASSWORD`

*   **Description**: Password for a generic VPN service. This variable is a placeholder and would be used by custom user-defined VPN stacks.
*   **Default Value**: `password`
*   **Example Value**: `m3talvpnpass123`
*   **Used By**: Custom VPN stacks (e.g., OpenVPN, WireGuard client containers).

### System Settings

General system-wide configuration options for containers.

#### `PUID`

*   **Description**: The User ID (UID) that Docker containers will use to run processes and access mounted volumes. Setting this ensures proper file permissions and ownership on your host system.
*   **Default Value**: `1000` (common for the first non-root user)
*   **Example Value**: `1001`
*   **Used By**: `m3tal-dashboard` container, all other Docker Compose containers that specify `user: "${PUID:-1000}:${PGID:-1000}"`.

#### `PGID`

*   **Description**: The Group ID (GID) that Docker containers will use to run processes and access mounted volumes. Setting this ensures proper file permissions and ownership on your host system.
*   **Default Value**: `1000` (common for the first non-root user's primary group)
*   **Example Value**: `1001`
*   **Used By**: `m3tal-dashboard` container, all other Docker Compose containers that specify `user: "${PUID:-1000}:${PGID:-1000}"`.

#### `TZ`

*   **Description**: Sets the timezone for containers, ensuring that logs and timestamps within the containers reflect the correct local time.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Used By**: `m3tal-dashboard` container, other Docker Compose containers.

#### `LOG_LEVEL`

*   **Description**: Controls the verbosity of logging for the M3TAL API daemon. Higher levels provide more detailed output.
*   **Default Value**: `info`
*   **Example Value**: `debug`, `warn`, `error`
*   **Used By**: M3TAL API daemon.

#### `DEBUG_MODE`

*   **Description**: Enables additional debugging features and logging for the M3TAL API daemon. Useful for troubleshooting.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Used By**: M3TAL API daemon.

#### `METRICS_ENABLED`

*   **Description**: Controls whether the M3TAL API daemon exposes Prometheus-compatible metrics at a `/metrics` endpoint (typically on port `HTTP_PORT`).
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Used By**: M3TAL API daemon.