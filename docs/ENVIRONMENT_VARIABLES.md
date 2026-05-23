# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables control various aspects of the system, from core functionality to network routing, storage, and authentication.

All M3TAL environment variables are read from the `/etc/m3tal/.env` file. This file is the single source of truth for configuration, used by both the `m3tal` CLI binary and all Docker Compose stacks via the `--env-file` option. You can manage these variables easily using the `m3tal config wizard` or `m3tal config set <KEY> <value>` commands.

---

### Quick Reference

| Name                    | Description                                                                 | Default               | Example                    | Component(s)                              |
| :---------------------- | :-------------------------------------------------------------------------- | :-------------------- | :------------------------- | :---------------------------------------- |
| `HTTP_PORT`             | Port for the M3TAL API daemon.                                              | `8080`                | `8081`                     | M3TAL API Daemon                          |
| `STATE_DIR`             | General default directory for M3TAL's internal state files.                 | `./state`             | `/opt/m3tal/state`         | M3TAL CLI, M3TAL Dashboard                |
| `LOG_LEVEL`             | Minimum log level for M3TAL components.                                     | `info`                | `debug`                    | M3TAL API Daemon, M3TAL CLI, Dashboard    |
| `DEBUG_MODE`            | Enables debug logging and features.                                         | `false`               | `true`                     | M3TAL API Daemon, M3TAL CLI, Dashboard    |
| `METRICS_ENABLED`       | Enables /metrics endpoint for Prometheus-style metrics.                     | `true`                | `false`                    | M3TAL API Daemon                          |
| `DASHBOARD_SECRET`      | Secret key for M3TAL Dashboard session management.                          | `change_me_immediately` | `super_secret_key_123`     | M3TAL Dashboard                           |
| `API_TOKEN`             | API token for secure communication with the M3TAL API.                      | `change_me_api_token` | `api_token_abc_xyz`        | M3TAL API Daemon                          |
| `ADMIN_PASSWORD`        | Initial password for the default admin user.                                | `admin_pass`          | `MyStrongPassword123`      | M3TAL API Daemon, M3TAL Dashboard         |
| `NETWORK_NAME`          | Name of the Docker network shared by all M3TAL stacks.                      | `m3tal`               | `my_m3tal_network`         | Docker Compose                            |
| `LOCAL_IP`              | Local IP address of the host machine.                                       | `127.0.0.1`           | `192.168.1.100`            | Traefik, Cloudflared                      |
| `DOMAIN`                | Primary domain for Traefik routing.                                         | `localhost`           | `example.com`              | Traefik, Cloudflared, M3TAL Dashboard, API|
| `DASHBOARD_PORT`        | Host port for the M3TAL Dashboard in `local` expose mode.                   | `8082`                | `9000`                     | M3TAL Dashboard (local mode)              |
| `DASHBOARD_EXPOSE_MODE` | Determines how the Dashboard is exposed.                                    | `local`               | `traefik`                  | M3TAL CLI (for compose management)        |
| `BASE_STORAGE_PATH`     | Base path for all M3TAL persistent data.                                    | `./data`              | `/mnt/m3tal-data`          | Docker Compose (volume mounts), M3TAL CLI |
| `MEDIA_PATH`            | Path for media files, relative to `BASE_STORAGE_PATH`.                      | `./data/media`        | `/mnt/m3tal-data/media`    | Docker Compose (volume mounts)            |
| `CONFIG_PATH`           | Path for M3TAL configuration files, relative to `BASE_STORAGE_PATH`.        | `./data/config`       | `/mnt/m3tal-data/config`   | Docker Compose (volume mounts)            |
| `DOWNLOADS_PATH`        | Path for downloaded files, relative to `BASE_STORAGE_PATH`.                 | `./data/downloads`    | `/mnt/m3tal-data/downloads`| Docker Compose (volume mounts)            |
| `TRAEFIK_WEB_PORT`      | Traefik HTTP entrypoint port on the host.                                   | `80`                  | `8080`                     | Traefik                                   |
| `TRAEFIK_WEBHTTPS_PORT` | Traefik HTTPS entrypoint port on the host.                                  | `443`                 | `8443`                     | Traefik                                   |
| `TRAEFIK_DASHBOARD_PORT`| Traefik's internal dashboard port.                                          | `8080`                | `8081`                     | Traefik                                   |
| `VPN_USER`              | Username for a hypothetical VPN service.                                    | `user`                | `m3talvpn`                 | VPN Component (future/custom)             |
| `VPN_PASSWORD`          | Password for a hypothetical VPN service.                                    | `password`            | `SecureVpnPass123`         | VPN Component (future/custom)             |
| `PUID`                  | User ID for container processes.                                            | `1000`                | `999`                      | Docker Compose (all data containers)      |
| `PGID`                  | Group ID for container processes.                                           | `1000`                | `999`                      | Docker Compose (all data containers)      |
| `TZ`                    | Timezone for containers.                                                    | `America/Denver`      | `Europe/London`            | Docker Compose (all containers)           |

---

### Detailed Variable Reference

#### Core Configuration

*   **`HTTP_PORT`**
    *   **Description**: Specifies the port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming connections. Traefik will route `api.${DOMAIN}` to this port internally.
    *   **Default Value**: `8080`
    *   **Example Value**: `8081`
    *   **Component(s) using it**: M3TAL API Daemon, Traefik (for routing)

*   **`STATE_DIR`**
    *   **Description**: This variable specifies a *general purpose default* directory for M3TAL's internal state files and configuration. While the M3TAL API daemon's primary database is *always* located at `/var/lib/m3tal/state.db`, and the dashboard's user configuration (`users.json`) is managed under the `CONFIG_PATH/m3tal/state` host directory, `STATE_DIR` can serve as a fallback or base path for other internal components or legacy configurations that might reference a generic state location within containers.
    *   **Default Value**: `./state`
    *   **Example Value**: `/opt/m3tal/state`
    *   **Component(s) using it**: M3TAL CLI, M3TAL Dashboard (as an internal environment variable mapping to the host's `CONFIG_PATH/m3tal/state` for `users.json`).

*   **`LOG_LEVEL`**
    *   **Description**: Sets the minimum severity level for logs generated by M3TAL components. Available levels typically include `debug`, `info`, `warn`, `error`, `fatal`.
    *   **Default Value**: `info`
    *   **Example Value**: `debug`
    *   **Component(s) using it**: M3TAL API Daemon, M3TAL CLI, M3TAL Dashboard, other M3TAL-managed containers.

*   **`DEBUG_MODE`**
    *   **Description**: When set to `true`, enables verbose debug logging and potentially additional development features across M3TAL components.
    *   **Default Value**: `false`
    *   **Example Value**: `true`
    *   **Component(s) using it**: M3TAL API Daemon, M3TAL CLI, M3TAL Dashboard.

*   **`METRICS_ENABLED`**
    *   **Description**: Controls whether the M3TAL API daemon exposes a `/metrics` endpoint for Prometheus-style metrics collection.
    *   **Default Value**: `true`
    *   **Example Value**: `false`
    *   **Component(s) using it**: M3TAL API Daemon.

#### Authentication & Access

*   **`DASHBOARD_SECRET`**
    *   **Description**: A strong secret key used by the M3TAL Dashboard for session management, cookie signing, and other security-sensitive operations. **This variable is auto-generated on first `m3tal init` and should NOT be set manually by users unless rotating the secret.**
    *   **Default Value**: `change_me_immediately` (placeholder)
    *   **Example Value**: `my_super_secure_dashboard_secret_12345`
    *   **Component(s) using it**: M3TAL Dashboard.

*   **`API_TOKEN`**
    *   **Description**: The API token used for authenticating requests to the M3TAL API daemon. This token secures communication between the dashboard, CLI, and potentially other services with the core API. **This variable is auto-generated on first `m3tal init` and should NOT be set manually by users unless rotating the token.**
    *   **Default Value**: `change_me_api_token` (placeholder)
    *   **Example Value**: `m3tal_api_token_xyz_789_abc`
    *   **Component(s) using it**: M3TAL API Daemon (for validation), M3TAL Dashboard (for making requests), M3TAL CLI.

*   **`ADMIN_PASSWORD`**
    *   **Description**: Sets the initial password for the default administrative user of the M3TAL system. This is used during the first setup or for password resets via the CLI.
    *   **Default Value**: `admin_pass`
    *   **Example Value**: `MyStrongAndSecurePass!2024`
    *   **Component(s) using it**: M3TAL API Daemon (for user management), M3TAL Dashboard (for initial login).

#### Network & Routing

*   **`NETWORK_NAME`**
    *   **Description**: Defines the name of the external Docker network that all M3TAL-managed compose stacks will share. This allows services in different `*-compose.yml` files to communicate with each other.
    *   **Default Value**: `m3tal`
    *   **Example Value**: `my_m3tal_network`
    *   **Component(s) using it**: Docker Compose (all stacks).

*   **`LOCAL_IP`**
    *   **Description**: Specifies the local IP address of the M3TAL host machine. This can be used by services that need to bind to a specific IP or reference the host.
    *   **Default Value**: `127.0.0.1`
    *   **Example Value**: `192.168.1.100`
    *   **Component(s) using it**: Traefik, Cloudflared (if configured), M3TAL API Daemon (if binding to specific IP).

*   **`DOMAIN`**
    *   **Description**: Sets the primary domain for your M3TAL deployment. When Traefik is enabled, this variable controls the routing rules, enabling access to services like `dash.DOMAIN` for the dashboard and `api.DOMAIN` for the M3TAL API.
    *   **Default Value**: `localhost`
    *   **Example Value**: `myhomelab.net`
    *   **Component(s) using it**: Traefik, Cloudflared, M3TAL Dashboard (for Traefik labels), M3TAL API Daemon (for URL generation).

*   **`DASHBOARD_PORT`**
    *   **Description**: Specifies the host port used to expose the M3TAL Dashboard directly when `DASHBOARD_EXPOSE_MODE` is set to `local`.
    *   **Default Value**: `8082`
    *   **Example Value**: `9000`
    *   **Component(s) using it**: M3TAL Dashboard (in `local` expose mode).

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description**: Controls how the M3TAL Dashboard container is made accessible.
        *   **`local` (default)**: The dashboard is exposed directly on the host's `DASHBOARD_PORT` (e.g., `http://HOST_IP:8082`). This mode uses the `m3tal-compose.local.yml` override, adding a direct port binding. Ideal for LAN-only setups or initial configuration without Traefik.
        *   **`traefik`**: The dashboard is exposed via the Traefik gateway at `http://dash.${DOMAIN}`. This mode uses the `m3tal-compose.traefik.yml` override, adding Traefik labels for routing. Requires Traefik to be running via `m3tal up`.
    *   **Default Value**: `local`
    *   **Example Value**: `traefik`
    *   **Component(s) using it**: M3TAL CLI (determines which compose override to use), M3TAL Dashboard (via Traefik labels).

#### Storage & Paths

*   **`BASE_STORAGE_PATH`**
    *   **Description**: The root directory on the host machine where all M3TAL persistent data (media, configurations, downloads) will be stored.
        *   **Important**: While the template defaults to `./data`, in production M3TAL deployments, this variable typically defaults to `/mnt` to ensure data persistence on dedicated storage.
    *   **Default Value**: `./data`
    *   **Example Value**: `/mnt/m3tal-data`
    *   **Component(s) using it**: Docker Compose (all volume mounts), M3TAL CLI.

*   **`MEDIA_PATH`**
    *   **Description**: Specifies the path for storing media files, relative to `BASE_STORAGE_PATH`.
    *   **Default Value**: `./data/media`
    *   **Example Value**: `/mnt/m3tal-data/media`
    *   **Component(s) using it**: Docker Compose (volume mounts for media-consuming containers).

*   **`CONFIG_PATH`**
    *   **Description**: Specifies the path for M3TAL-specific configuration files (e.g., `users.json` for the dashboard), relative to `BASE_STORAGE_PATH`.
    *   **Default Value**: `./data/config`
    *   **Example Value**: `/mnt/m3tal-data/config`
    *   **Component(s) using it**: Docker Compose (volume mounts for config-storing containers), M3TAL Dashboard (for `users.json`).

*   **`DOWNLOADS_PATH`**
    *   **Description**: Specifies the path for downloaded files, relative to `BASE_STORAGE_PATH`.
    *   **Default Value**: `./data/downloads`
    *   **Example Value**: `/mnt/m3tal-data/downloads`
    *   **Component(s) using it**: Docker Compose (volume mounts for download clients).

#### Traefik Gateway

*   **`TRAEFIK_WEB_PORT`**
    *   **Description**: The host port on which Traefik listens for incoming HTTP (non-HTTPS) traffic.
    *   **Default Value**: `80`
    *   **Example Value**: `8080`
    *   **Component(s) using it**: Traefik.

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description**: The host port on which Traefik listens for incoming HTTPS traffic.
    *   **Default Value**: `443`
    *   **Example Value**: `8443`
    *   **Component(s) using it**: Traefik.

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description**: The internal port used by Traefik for its own management dashboard. This port is typically exposed only on localhost (e.g., `127.0.0.1:8081`) for administrative access.
    *   **Default Value**: `8080`
    *   **Example Value**: `8081`
    *   **Component(s) using it**: Traefik.

#### VPN Configuration

*   **`VPN_USER`**
    *   **Description**: Username for a hypothetical VPN service that could be integrated with the M3TAL ecosystem. (Note: No VPN component is currently included in the default M3TAL services list, but these variables are reserved for future or custom integrations).
    *   **Default Value**: `user`
    *   **Example Value**: `m3talvpn`
    *   **Component(s) using it**: VPN Component (future/custom).

*   **`VPN_PASSWORD`**
    *   **Description**: Password for a hypothetical VPN service that could be integrated with the M3TAL ecosystem. (Note: No VPN component is currently included in the default M3TAL services list, but these variables are reserved for future or custom integrations).
    *   **Default Value**: `password`
    *   **Example Value**: `SecureVpnPass123`
    *   **Component(s) using it**: VPN Component (future/custom).

#### System & Container Settings

*   **`PUID`**
    *   **Description**: Specifies the User ID (UID) that container processes will run as. This ensures correct file permissions for volumes mounted from the host system.
    *   **Default Value**: `1000`
    *   **Example Value**: `999`
    *   **Component(s) using it**: Docker Compose (all data-persisting containers).

*   **`PGID`**
    *   **Description**: Specifies the Group ID (GID) that container processes will run as. This ensures correct file permissions for volumes mounted from the host system.
    *   **Default Value**: `1000`
    *   **Example Value**: `999`
    *   **Component(s) using it**: Docker Compose (all data-persisting containers).

*   **`TZ`**
    *   **Description**: Sets the timezone inside M3TAL containers, ensuring consistent time reporting and scheduling.
    *   **Default Value**: `America/Denver`
    *   **Example Value**: `Europe/London`
    *   **Component(s) using it**: Docker Compose (all containers).