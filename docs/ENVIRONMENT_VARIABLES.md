# Environment Variables Reference (`docs/ENVIRONMENT_VARIABLES.md`)

As DocSmith, the M3TAL Ecosystem Documentation Architect, I'm here to provide a definitive guide to the environment variables that power your M3TAL system.

All M3TAL configuration, including these environment variables, is centralized in the `/etc/m3tal/.env` file. Both the `m3tal` CLI binary and all Docker Compose stacks read their configurations from this file via the `--env-file` option, ensuring a single source of truth for your entire deployment.

---

### Quick Reference Table

| Variable Name           | Default Value         | Description                                                      |
| :---------------------- | :-------------------- | :--------------------------------------------------------------- |
| `DASHBOARD_PORT`        | `8082`                | Internal port for the M3TAL Dashboard container.                 |
| `DASHBOARD_EXPOSE_MODE` | `local`               | How the Dashboard is exposed: `local` (port bind) or `traefik` (domain route). |
| `HTTP_PORT`             | `8080`                | Port for the M3TAL API daemon.                                   |
| `STATE_DIR`             | `./state`             | Directory for API database and dashboard state.                  |
| `LOG_LEVEL`             | `info`                | Logging verbosity for M3TAL components.                          |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret for Dashboard session management. **Auto-generated.**     |
| `API_TOKEN`             | `change_me_api_token` | Token for API daemon authentication. **Auto-generated.**         |
| `ADMIN_PASSWORD`        | `admin_pass`          | Initial default password for the Dashboard `admin` user.         |
| `NETWORK_NAME`          | `m3tal`               | Name of the primary Docker network for M3TAL.                    |
| `LOCAL_IP`              | `127.0.0.1`           | Local IP of the host machine.                                    |
| `DOMAIN`                | `localhost`           | Primary domain for Traefik routing (e.g., `dash.DOMAIN`).      |
| `VPN_USER`              | `user`                | Username for user-defined VPN services.                          |
| `VPN_PASSWORD`          | `password`            | Password for user-defined VPN services.                          |
| `BASE_STORAGE_PATH`     | `./data`              | Base directory for all persistent M3TAL data.                    |
| `MEDIA_PATH`            | `./data/media`        | Sub-path within `BASE_STORAGE_PATH` for media files.             |
| `CONFIG_PATH`           | `./data/config`       | Sub-path within `BASE_STORAGE_PATH` for configuration files.     |
| `DOWNLOADS_PATH`        | `./data/downloads`    | Sub-path within `BASE_STORAGE_PATH` for downloads.               |
| `PUID`                  | `1000`                | User ID for containers to ensure correct permissions.            |
| `PGID`                  | `1000`                | Group ID for containers to ensure correct permissions.           |
| `TZ`                    | `America/Denver`      | Timezone for M3TAL containers.                                   |
| `TRAEFIK_WEB_PORT`      | `80`                  | Host port for Traefik's HTTP entry point.                        |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                 | Host port for Traefik's HTTPS entry point.                       |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                | Internal container port for the Traefik dashboard.               |
| `DEBUG_MODE`            | `false`               | Enables or disables debug features and logging.                  |
| `METRICS_ENABLED`       | `true`                | Controls exposure of Prometheus-compatible metrics.              |

---

### Detailed Environment Variable Reference

This section provides an in-depth look at each environment variable, including its purpose, default, example, and the M3TAL component(s) that utilize it.

#### Core Configuration

##### `HTTP_PORT`

*   **Description:** Specifies the port on which the M3TAL API daemon listens for incoming HTTP connections. The Dashboard and Traefik connect to the API via `http://host.docker.internal:${HTTP_PORT}`.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** API daemon, Dashboard, Traefik (for `api.DOMAIN` routing).

##### `STATE_DIR`

*   **Description:** Defines the directory where the M3TAL API daemon stores its SQLite state database (`state.db`) and where the M3TAL Dashboard expects to find container-side configuration (e.g., `users.json`). *Note: On the host filesystem, the API daemon typically uses `/var/lib/m3tal/state.db` and the Dashboard uses `/docker/state` within its container, which maps to `${CONFIG_PATH:-/mnt/config}/m3tal/state` on the host.*
*   **Default Value:** `./state`
*   **Example Value:** `/docker/state` (within a container)
*   **Used By:** API daemon, Dashboard.

##### `LOG_LEVEL`

*   **Description:** Sets the logging verbosity for M3TAL components, controlling the amount of detail output to logs.
*   **Default Value:** `info`
*   **Example Value:** `debug`, `warn`, `error`
*   **Used By:** API daemon, Dashboard.

##### `DEBUG_MODE`

*   **Description:** Enables or disables debug logging and features across M3TAL components for troubleshooting.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** API daemon, Dashboard.

##### `METRICS_ENABLED`

*   **Description:** Controls whether M3TAL components, primarily the API daemon, expose Prometheus-compatible metrics endpoints.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** API daemon.

#### Authentication

##### `DASHBOARD_SECRET`

*   **Description:** A secret key used by the M3TAL Dashboard for session management, cookie signing, and securing sensitive operations. This value helps protect the integrity and confidentiality of user sessions. **This variable is automatically generated on the first `m3tal init` run. Users should generally NOT set this manually unless performing a secret rotation.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `my_strong_dashboard_secret_1234567890abcdef`
*   **Used By:** Dashboard, CLI (`m3tal init`).

##### `API_TOKEN`

*   **Description:** An authentication token used to secure access to the M3TAL API daemon. The CLI and Dashboard use this token to authenticate with the API. **This variable is automatically generated on the first `m3tal init` run. Users should generally NOT set this manually unless performing a token rotation.**
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `my_strong_api_token_0987654321fedcba`
*   **Used By:** API daemon, CLI (`m3tal init`), Dashboard.

##### `ADMIN_PASSWORD`

*   **Description:** The initial default password for the `admin` user account in the M3TAL Dashboard. It is highly recommended to change this immediately after your first login via the dashboard's user management interface or by using the `m3tal dashpass` CLI command.
*   **Default Value:** `admin_pass`
*   **Example Value:** `MySecureAdminPassword123!`
*   **Used By:** Dashboard (initial user setup), CLI (`m3tal dashpass`).

#### Network Configuration

##### `NETWORK_NAME`

*   **Description:** Defines the name of the primary Docker network that M3TAL core components and user-deployed stacks utilize for internal inter-container communication.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal-proxy-network`
*   **Used By:** CLI (for `docker network create`), all compose stacks.

##### `LOCAL_IP`

*   **Description:** Specifies the local IP address of the host machine. While most inter-service communication within Docker containers uses `host.docker.internal` to reach the host API daemon, this variable can be used for specific host bindings or network configurations if required by certain services.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** API daemon (for binding), CLI.

#### Storage Paths

##### `BASE_STORAGE_PATH`

*   **Description:** The root directory on the host filesystem where M3TAL stores all persistent data, including user-defined media, configurations, and downloads. **In production M3TAL deployments, this variable defaults to `/mnt`, not `./data` as seen in development templates.**
*   **Default Value:** `./data` (template default) / `/mnt` (production default)
*   **Example Value:** `/mnt/m3tal-data`, `/var/lib/m3tal-volumes`
*   **Used By:** Dashboard, CLI, User Stacks.

##### `MEDIA_PATH`

*   **Description:** A sub-path, typically relative to `BASE_STORAGE_PATH`, where media files are expected to be stored. This is primarily used by user-deployed stacks (e.g., media servers like Plex or Jellyfin) that require access to media libraries.
*   **Default Value:** `./data/media`
*   **Example Value:** `${BASE_STORAGE_PATH}/media`, `/mnt/media`
*   **Used By:** User Stacks.

##### `CONFIG_PATH`

*   **Description:** A sub-path, typically relative to `BASE_STORAGE_PATH`, where M3TAL's internal configuration files (such as `state.db` and `users.json`) are stored on the host filesystem. This is distinct from the `/etc/m3tal/.env` file itself.
*   **Default Value:** `./data/config`
*   **Example Value:** `${BASE_STORAGE_PATH}/config`, `/mnt/config`
*   **Used By:** Dashboard, API daemon (indirectly via volume mounts).

##### `DOWNLOADS_PATH`

*   **Description:** A sub-path, typically relative to `BASE_STORAGE_PATH`, where downloaded files are expected to be stored. This is commonly used by user-deployed download clients.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `${BASE_STORAGE_PATH}/downloads`, `/mnt/downloads`
*   **Used By:** User Stacks.

#### Traefik Gateway Configuration

##### `DOMAIN`

*   **Description:** The primary domain name used for Traefik routing rules. Setting this variable enables Traefik to expose core M3TAL services at subdomains like `dash.DOMAIN` (for the Dashboard) and `api.DOMAIN` (for the API daemon).
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.example.com`
*   **Used By:** Traefik, CLI (for generating Traefik dynamic configs), Dashboard (indirectly via Traefik labels in `traefik` expose mode).

##### `TRAEFIK_WEB_PORT`

*   **Description:** The host port on which the Traefik gateway listens for standard HTTP (unencrypted) traffic.
*   **Default Value:** `80`
*   **Example Value:** `8080` (if port 80 is occupied by another service)
*   **Used By:** Traefik.

##### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The host port on which the Traefik gateway listens for HTTPS (encrypted) traffic. Utilizing this port requires additional Traefik configuration for SSL certificate management.
*   **Default Value:** `443`
*   **Example Value:** `8443`
*   **Used By:** Traefik (when HTTPS is configured).

##### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The *internal* container port where the Traefik dashboard itself is exposed. The external host port for the Traefik dashboard is typically mapped to `127.0.0.1:8081` by default.
*   **Default Value:** `8080`
*   **Example Value:** `8082`
*   **Used By:** Traefik.

#### Dashboard Specific

##### `DASHBOARD_PORT`

*   **Description:** Controls the internal port on which the M3TAL Dashboard container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this port also determines the direct host port binding (`HOST_IP:${DASHBOARD_PORT}`).
*   **Default Value:** `8082`
*   **Example Value:** `8083`
*   **Used By:** Dashboard, CLI (for `local` mode port binding).

##### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Determines how the M3TAL Dashboard service is made accessible.
    *   `local`: The dashboard container's port is directly bound to a host port, typically `HOST_IP:8082`. This mode requires no Traefik configuration and is ideal for LAN-only setups or first-time users.
    *   `traefik`: The dashboard container is configured with Traefik labels, allowing Traefik to route traffic from `dash.DOMAIN` (if `DOMAIN` is set) to the dashboard on port 8082. This mode requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** CLI (during `m3tal dash up` to select the appropriate compose override), Dashboard (indirectly via compose labels).

#### VPN Configuration

##### `VPN_USER`

*   **Description:** Specifies the username for a user-defined VPN service that you might integrate with your M3TAL setup. *Note: This variable is not utilized by core M3TAL components (like Cloudflared) but can be leveraged by custom user stacks requiring VPN credentials.*
*   **Default Value:** `user`
*   **Example Value:** `my_vpn_username`
*   **Used By:** User Stacks (e.g., a VPN client container).

##### `VPN_PASSWORD`

*   **Description:** Specifies the password for a user-defined VPN service. *Note: Similar to `VPN_USER`, this variable is intended for use by custom user stacks requiring VPN credentials, not core M3TAL components.*
*   **Default Value:** `password`
*   **Example Value:** `my_strong_vpn_password`
*   **Used By:** User Stacks (e.g., a VPN client container).

#### System Configuration

##### `PUID`

*   **Description:** The User ID (UID) that M3TAL containers (e.g., Dashboard) should run as inside the container. Setting this ensures that containers have correct file ownership and permissions when accessing mounted volumes on the host.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** Dashboard, User Stacks.

##### `PGID`

*   **Description:** The Group ID (GID) that M3TAL containers (e.g., Dashboard) should run as inside the container. This works in conjunction with `PUID` to manage file system permissions on mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** Dashboard, User Stacks.

##### `TZ`

*   **Description:** Sets the timezone for M3TAL containers, ensuring that timestamps in logs, scheduled tasks, and internal operations are correctly localized.
*   **Default Value:** `America/Denver`
*   **Example Value:** `America/New_York`, `Europe/London`, `Asia/Tokyo`
*   **Used By:** Dashboard, User Stacks.