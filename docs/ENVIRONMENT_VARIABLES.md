Alright, fellow M3TAL Engineers! DocSmith here, ready to lay down the law on our critical environment variables. These are the levers and pulleys of your M3TAL ecosystem, and understanding them is key to a smooth operation.

All M3TAL components — from the CLI binary to your Docker Compose stacks — read their configuration from a single, authoritative source: the `.env` file located at `/etc/m3tal/.env`. This file is managed primarily by the `m3tal config wizard` and `m3tal config set` commands, ensuring a consistent and centralized configuration across your entire system. When you execute `docker compose` operations via the `m3tal` CLI, the `--env-file /etc/m3tal/.env` flag is automatically applied, making these variables universally accessible.

---

## Environment Variable Reference

### Quick Reference Table

| Variable Name           | Default Value         | Description                                                      | Component(s)                     |
| :---------------------- | :-------------------- | :--------------------------------------------------------------- | :------------------------------- |
| `DASHBOARD_PORT`        | `8082`                | Internal/External port for the M3TAL Dashboard.                  | `m3tal-dashboard`                |
| `DASHBOARD_EXPOSE_MODE` | `local`               | Mode to expose the Dashboard (direct port or via Traefik).       | `CLI binary`                     |
| `HTTP_PORT`             | `8080`                | Listening port for the M3TAL API daemon.                         | `API daemon`                     |
| `STATE_DIR`             | `./state`             | *Internal* path for dashboard state files within its container.  | `m3tal-dashboard`                |
| `LOG_LEVEL`             | `info`                | Minimum logging level for the M3TAL API daemon.                  | `API daemon`                     |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret key for Dashboard session management.                     | `m3tal-dashboard`                |
| `API_TOKEN`             | `change_me_api_token` | Bearer token for authenticating with M3TAL API.                  | `CLI binary`, `m3tal-dashboard`  |
| `ADMIN_PASSWORD`        | `admin_pass`          | Initial password for the Dashboard's `admin` user.               | `m3tal-dashboard`                |
| `NETWORK_NAME`          | `m3tal`               | Name of the Docker network for inter-container communication.    | `All Compose Stacks`             |
| `LOCAL_IP`              | `127.0.0.1`           | IP address for host-gateway and internal communication.          | `m3tal-dashboard`                |
| `DOMAIN`                | `localhost`           | Base domain name for Traefik routing rules.                      | `Traefik gateway`, `m3tal-dashboard` |
| `VPN_USER`              | `user`                | Username for VPN services.                                       | `Future VPN services`            |
| `VPN_PASSWORD`          | `password`            | Password for VPN services.                                       | `Future VPN services`            |
| `BASE_STORAGE_PATH`     | `./data`              | Root directory for all M3TAL-managed persistent data.            | `All Compose Stacks`             |
| `MEDIA_PATH`            | `./data/media`        | Directory for user-uploaded media files.                         | `m3tal-dashboard` (via mount)    |
| `CONFIG_PATH`           | `./data/config`       | Directory for configuration files and dashboard state on host.   | `m3tal-dashboard` (via mount)    |
| `DOWNLOADS_PATH`        | `./data/downloads`    | Directory for downloaded content.                                | `Future services`                |
| `PUID`                  | `1000`                | User ID (UID) for containers to manage file permissions.         | `All Compose Stacks`             |
| `PGID`                  | `1000`                | Group ID (GID) for containers to manage file permissions.        | `All Compose Stacks`             |
| `TZ`                    | `America/Denver`      | Timezone for containers to ensure correct timestamps.            | `All Compose Stacks`             |
| `TRAEFIK_WEB_PORT`      | `80`                  | The HTTP port Traefik listens on for incoming web requests.      | `Traefik gateway`                |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                 | The HTTPS port Traefik listens on for incoming secure requests.  | `Traefik gateway`                |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                | The internal port Traefik's own dashboard listens on.            | `Traefik gateway`                |
| `DEBUG_MODE`            | `false`               | Enables debug logging and features across components.            | `CLI binary`, `API daemon`, `m3tal-dashboard` |
| `METRICS_ENABLED`       | `true`                | Controls whether system and application metrics are collected.   | `API daemon`                     |

---

### Detailed Variable Reference

#### Core Configuration

These variables control fundamental aspects of the M3TAL system's operation, logging, and general user/group permissions for containers.

*   **`HTTP_PORT`**
    *   **Description**: The port on which the M3TAL API daemon (the Go binary running as `m3tal-api.service`) listens for incoming HTTP requests. This is the backend API for all M3TAL operations.
    *   **Default Value**: `8080`
    *   **Example Value**: `8000`
    *   **Used by**: `API daemon`, `m3tal-dashboard` (to connect to the API), `Traefik gateway` (for routing `api.${DOMAIN}` to the API).
*   **`STATE_DIR`**
    *   **Description**: Internally, the M3TAL Dashboard container uses this variable to define the path where it stores its state files, such as `users.json`. It is typically mapped from a host volume (defined by `CONFIG_PATH`). *Note: This specific variable from `/etc/m3tal/.env` is not directly consumed by any described component as a host-level environment variable. The Dashboard sets its `STATE_DIR` internally to `/docker/state` within its container.*
    *   **Default Value**: `./state`
    *   **Example Value**: `/app/state`
    *   **Used by**: `m3tal-dashboard` (internally within its container).
*   **`LOG_LEVEL`**
    *   **Description**: Sets the minimum logging severity level for the M3TAL API daemon. Valid values typically include `debug`, `info`, `warn`, `error`, `fatal`.
    *   **Default Value**: `info`
    *   **Example Value**: `debug`
    *   **Used by**: `API daemon`.
*   **`PUID`**
    *   **Description**: The User ID (UID) that Docker containers launched by M3TAL should use. This is crucial for ensuring correct file ownership and permissions for persistent storage volumes, preventing permissions issues when containers write to host directories.
    *   **Default Value**: `1000`
    *   **Example Value**: `1001`
    *   **Used by**: `m3tal-dashboard` container, other compose stacks.
*   **`PGID`**
    *   **Description**: The Group ID (GID) that Docker containers launched by M3TAL should use, complementing `PUID` for file permissions on mounted volumes.
    *   **Default Value**: `1000`
    *   **Example Value**: `1001`
    *   **Used by**: `m3tal-dashboard` container, other compose stacks.
*   **`TZ`**
    *   **Description**: Specifies the timezone for Docker containers. Setting this ensures that logs, timestamps, and scheduled tasks within containers operate according to your local or preferred time zone.
    *   **Default Value**: `America/Denver`
    *   **Example Value**: `Europe/London`
    *   **Used by**: `m3tal-dashboard` container, other compose stacks.
*   **`DEBUG_MODE`**
    *   **Description**: When set to `true`, this variable enables detailed debug logging and potentially activates debug-specific features across various M3TAL components. This is useful for troubleshooting.
    *   **Default Value**: `false`
    *   **Example Value**: `true`
    *   **Used by**: `CLI binary`, `API daemon`, `m3tal-dashboard`.
*   **`METRICS_ENABLED`**
    *   **Description**: Controls whether system and application metrics are collected and exposed by the M3TAL API daemon. Set to `false` to disable metrics collection.
    *   **Default Value**: `true`
    *   **Example Value**: `false`
    *   **Used by**: `API daemon`.

#### Authentication & Security

These variables manage user authentication for the M3TAL Dashboard and API access.

*   **`DASHBOARD_SECRET`**
    *   **Description**: A secret key used by the M3TAL Dashboard for secure session management, cookie signing, and other cryptographic operations. **It is critical to keep this value secure.**
    *   **Default Value**: `change_me_immediately`
    *   **Example Value**: `your_very_long_and_complex_dashboard_secret_string`
    *   **Used by**: `m3tal-dashboard` container.
    *   **Important**: This variable is automatically generated by `m3tal init`. You should NOT set it manually unless you are intentionally rotating your secrets, in which case you must restart the dashboard.
*   **`API_TOKEN`**
    *   **Description**: The Bearer token used for authenticating requests to the M3TAL API daemon. The `m3tal` CLI and the M3TAL Dashboard use this token to interact with the API. **Keep this token highly confidential.**
    *   **Default Value**: `change_me_api_token`
    *   **Example Value**: `a_super_secure_api_token_generated_by_m3tal_init`
    *   **Used by**: `CLI binary`, `API daemon` (for validation), `m3tal-dashboard` (for API calls).
    *   **Important**: This variable is automatically generated by `m3tal init`. You should NOT set it manually unless you are intentionally rotating your secrets, in which case you must restart any services that use it.
*   **`ADMIN_PASSWORD`**
    *   **Description**: The initial default password for the `admin` user account within the M3TAL Dashboard. After initial setup or using `m3tal dashpass`, this value may become irrelevant as user credentials are stored in `/docker/users.json`.
    *   **Default Value**: `admin_pass`
    *   **Example Value**: `MySecureAdminPassword123!`
    *   **Used by**: `m3tal-dashboard` container (for initial `users.json` generation).

#### Network Configuration

Variables pertaining to Docker networking and internal IP addressing.

*   **`NETWORK_NAME`**
    *   **Description**: The name of the Docker network that M3TAL uses for communication between its core services and any user-deployed applications. This provides a segregated and manageable network environment.
    *   **Default Value**: `m3tal`
    *   **Example Value**: `my_m3tal_network`
    *   **Used by**: All M3TAL Docker Compose stacks (e.g., `m3tal-compose.yml`, `routing-compose.yml`).
*   **`LOCAL_IP`**
    *   **Description**: Specifies an IP address that Docker containers can use to reach the host system. For containers that need to communicate with services running directly on the host (like the M3TAL API daemon), this IP is aliased as `host.docker.internal`.
    *   **Default Value**: `127.0.0.1`
    *   **Example Value**: `192.168.1.100` (if `host.docker.internal` needs to point to a specific LAN IP).
    *   **Used by**: `m3tal-dashboard` container (`extra_hosts` directive).

#### Storage Paths

These variables define the base paths where M3TAL stores various types of persistent data on your host system.

*   **`BASE_STORAGE_PATH`**
    *   **Description**: The primary root directory on the host system where M3TAL and its deployed services will store all persistent data, including media, configs, and downloads. All other storage paths (`MEDIA_PATH`, `CONFIG_PATH`, `DOWNLOADS_PATH`) are typically subdirectories of this base path.
    *   **Default Value**: `./data`
    *   **Example Value**: `/srv/m3tal`
    *   **Used by**: All M3TAL Docker Compose stacks (for volume mounts).
    *   **Important**: In production deployments, this variable defaults to `/mnt` (e.g., `/mnt/data`), reflecting a common practice for mounting large storage devices. The default `./data` in the template is for local development or quick starts.
*   **`MEDIA_PATH`**
    *   **Description**: The dedicated host directory for storing user-uploaded media files and other multimedia content managed by M3TAL applications. This path is expected to be a subdirectory of `BASE_STORAGE_PATH`.
    *   **Default Value**: `./data/media`
    *   **Example Value**: `/mnt/m3tal_data/media`
    *   **Used by**: `m3tal-dashboard` container (via its `/mnt` mount), other media-related compose stacks.
*   **`CONFIG_PATH`**
    *   **Description**: The host directory designated for configuration files and application-specific state that needs to persist outside containers. This is where the M3TAL Dashboard's `users.json` is stored on the host (within `/etc/m3tal/state/config`).
    *   **Default Value**: `./data/config`
    *   **Example Value**: `/etc/m3tal/config`
    *   **Used by**: `m3tal-dashboard` container (for mounting its `/docker/state` directory from `${CONFIG_PATH}/m3tal/state`).
*   **`DOWNLOADS_PATH`**
    *   **Description**: The host directory intended for storing downloaded content from various M3TAL services. This path is typically a subdirectory of `BASE_STORAGE_PATH`.
    *   **Default Value**: `./data/downloads`
    *   **Example Value**: `/mnt/m3tal_data/downloads`
    *   **Used by**: Future download-management services, applications that interact with a downloads directory.

#### Traefik Gateway Configuration

Variables specifically for configuring the Traefik reverse proxy and how services are exposed.

*   **`DOMAIN`**
    *   **Description**: The base domain name that Traefik will use to define routing rules for exposed services. For example, setting `DOMAIN=m3tal.com` will enable access to `dash.m3tal.com` for the dashboard and `api.m3tal.com` for the API.
    *   **Default Value**: `localhost`
    *   **Example Value**: `yourdomain.com`
    *   **Used by**: `Traefik gateway` (for dynamic routing), `m3tal-compose.traefik.yml` (for dashboard routing).
    *   **Important**: Setting this variable enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.
*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description**: Determines how the M3TAL Dashboard container is exposed to the network.
        *   `local` (default): The dashboard is exposed directly on a host port (`${DASHBOARD_PORT:-8082}`). Access via `http://HOST_IP:8082`. No Traefik required. Best for LAN-only or local access.
        *   `traefik`: The dashboard is exposed via the Traefik reverse proxy, accessible at `http://dash.${DOMAIN}`. Requires Traefik to be running. Best for domain-based setups.
    *   **Default Value**: `local`
    *   **Example Value**: `traefik`
    *   **Used by**: `CLI binary` (`m3tal dash up` command) to select the appropriate Docker Compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`).
*   **`DASHBOARD_PORT`**
    *   **Description**: The internal port on which the M3TAL Dashboard container listens. This port is either directly mapped to the host (in `local` expose mode) or routed by Traefik (in `traefik` expose mode).
    *   **Default Value**: `8082`
    *   **Example Value**: `8088`
    *   **Used by**: `m3tal-dashboard` container (internal port), `m3tal-compose.local.yml` (host port mapping), `m3tal-compose.traefik.yml` (Traefik service definition).
*   **`TRAEFIK_WEB_PORT`**
    *   **Description**: The host port on which Traefik listens for incoming HTTP (non-secure) requests. This serves as the primary entry point for web traffic routed by Traefik.
    *   **Default Value**: `80`
    *   **Example Value**: `8080` (if Traefik is behind another reverse proxy).
    *   **Used by**: `Traefik gateway` container.
*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description**: The host port on which Traefik listens for incoming HTTPS (secure) requests. This is used for secure web traffic, often accompanied by SSL/TLS certificate management by Traefik.
    *   **Default Value**: `443`
    *   **Example Value**: `8443`
    *   **Used by**: `Traefik gateway` container.
*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description**: The internal port on which Traefik's own management dashboard listens. This port is typically mapped to a host-local address (e.g., `127.0.0.1:8081`) to restrict access.
    *   **Default Value**: `8080`
    *   **Example Value**: `9000`
    *   **Used by**: `Traefik gateway` container.

#### VPN Configuration

Placeholder variables for future VPN integration. Currently, Cloudflared is used, which does not consume these variables directly.

*   **`VPN_USER`**
    *   **Description**: A username credential intended for use with integrated VPN services to establish secure tunnels.
    *   **Default Value**: `user`
    *   **Example Value**: `mysecurevpnuser`
    *   **Used by**: Placeholder for future VPN integrations.
*   **`VPN_PASSWORD`**
    *   **Description**: A password credential intended for use with integrated VPN services to establish secure tunnels.
    *   **Default Value**: `password`
    *   **Example Value**: `MyStrongVPNPass!`
    *   **Used by**: Placeholder for future VPN integrations.

---

This comprehensive guide should provide you with all the necessary details to configure your M3TAL ecosystem effectively. Remember to always consult `m3tal config wizard` and `m3tal config set` for managing these variables securely.

DocSmith, signing off.