# Environment Variable Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables control the behavior of the M3TAL CLI, the API daemon, and all Docker Compose-managed services, including the M3TAL Dashboard, Traefik, and user-defined stacks.

All M3TAL environment variables are read from the primary configuration file: `/etc/m3tal/.env`. Both the `m3tal` CLI and all Docker Compose stacks (via `--env-file /etc/m3tal/.env`) consume these settings. It is highly recommended to manage these variables using the `m3tal config wizard` or `m3tal config set` commands.

---

## Quick Reference Table

| Variable Name           | Description                                                                 | Default Value            | Component(s)                         |
| :---------------------- | :-------------------------------------------------------------------------- | :----------------------- | :----------------------------------- |
| `HTTP_PORT`             | Port for the M3TAL API daemon.                                              | `8080`                   | `m3tal-api.service`                  |
| `STATE_DIR`             | Host path for M3TAL's state database and runtime files.                     | `/var/lib/m3tal/state`   | API, Dashboard                       |
| `LOG_LEVEL`             | Minimum logging level for the M3TAL API daemon.                             | `info`                   | API                                  |
| `DEBUG_MODE`            | Enables verbose logging and debug features.                                 | `false`                  | API, Dashboard, CLI                  |
| `METRICS_ENABLED`       | Enables Prometheus-compatible metrics endpoints.                            | `true`                   | API, Dashboard, Traefik              |
| `NETWORK_NAME`          | Name of the Docker network for inter-container communication.               | `m3tal`                  | All Docker Compose Stacks            |
| `LOCAL_IP`              | Host machine's IP, used for inter-service communication.                    | `127.0.0.1`              | API, Traefik                         |
| `PUID`                  | User ID for container processes.                                            | `1000`                   | All Docker Compose Containers        |
| `PGID`                  | Group ID for container processes.                                           | `1000`                   | All Docker Compose Containers        |
| `TZ`                    | Timezone for containers.                                                    | `America/Denver`         | All Docker Compose Containers        |
| `DASHBOARD_SECRET`      | Secret key for Dashboard session management.                                | `change_me_immediately`  | Dashboard                            |
| `API_TOKEN`             | Bearer token for API authentication.                                        | `change_me_api_token`    | API, CLI, Dashboard                  |
| `ADMIN_PASSWORD`        | Initial password for Dashboard admin user.                                  | `admin_pass`             | Dashboard                            |
| `BASE_STORAGE_PATH`     | Base directory for all M3TAL media, config, and download data.              | `/mnt`                   | All Docker Compose Stacks            |
| `MEDIA_PATH`            | Sub-path for media files.                                                   | `${BASE_STORAGE_PATH}/media` | User Stacks                          |
| `CONFIG_PATH`           | Sub-path for configuration files.                                           | `${BASE_STORAGE_PATH}/config` | Dashboard, User Stacks               |
| `DOWNLOADS_PATH`        | Sub-path for downloaded content.                                            | `${BASE_STORAGE_PATH}/downloads` | User Stacks                          |
| `DASHBOARD_PORT`        | Port where the M3TAL Dashboard container listens internally.                | `8082`                   | Dashboard                            |
| `DASHBOARD_EXPOSE_MODE` | Controls Dashboard exposure: `local` (direct port) or `traefik` (via domain). | `local`                  | CLI (for compose overrides), Dashboard |
| `DOMAIN`                | Primary domain for Traefik routing.                                         | `localhost`              | Traefik, Dashboard                   |
| `TRAEFIK_WEB_PORT`      | HTTP entrypoint port for Traefik.                                           | `80`                     | Traefik                              |
| `TRAEFIK_WEBHTTPS_PORT` | HTTPS entrypoint port for Traefik.                                          | `443`                    | Traefik                              |
| `TRAEFIK_DASHBOARD_PORT`| Internal port for Traefik's own dashboard.                                  | `8080`                   | Traefik                              |
| `VPN_USER`              | Username for VPN services.                                                  | `user`                   | VPN Clients (e.g., WireGuard)        |
| `VPN_PASSWORD`          | Password for VPN services.                                                  | `password`               | VPN Clients (e.g., WireGuard)        |

---

## Detailed Environment Variable Reference

### Core System Configuration

These variables configure the fundamental behavior of the M3TAL API daemon, CLI, and general container runtime.

#### `HTTP_PORT`
-   **Description**: The port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming HTTP requests.
-   **Default Value**: `8080`
-   **Example Value**: `8080`
-   **Component(s) Used By**: `m3tal-api.service`

#### `STATE_DIR`
-   **Description**: The absolute path on the host filesystem where the M3TAL system stores its state database (`state.db`), user credentials (`users.json`), and other runtime configuration files. In production deployments, this path defaults to `/var/lib/m3tal/state`.
-   **Default Value**: `./state` (This is the default in development, but `/var/lib/m3tal/state` is the production default managed by `m3tal init`).
-   **Example Value**: `/var/lib/m3tal/state`
-   **Component(s) Used By**: `m3tal-api.service`, `m3tal-dashboard`

#### `LOG_LEVEL`
-   **Description**: The minimum logging level for the M3TAL API daemon. Controls the verbosity of logs.
-   **Default Value**: `info`
-   **Example Value**: `debug`, `warn`, `error`
-   **Component(s) Used By**: `m3tal-api.service`

#### `DEBUG_MODE`
-   **Description**: When set to `true`, enables verbose logging and additional debug features across various M3TAL components. This is useful for troubleshooting.
-   **Default Value**: `false`
-   **Example Value**: `true`
-   **Component(s) Used By**: `m3tal-api.service`, `m3tal-dashboard`, `m3tal` CLI

#### `METRICS_ENABLED`
-   **Description**: When set to `true`, enables Prometheus-compatible metrics endpoints for monitoring the M3TAL API daemon and other services.
-   **Default Value**: `true`
-   **Example Value**: `false`
-   **Component(s) Used By**: `m3tal-api.service`, `m3tal-dashboard`, Traefik

#### `NETWORK_NAME`
-   **Description**: The name of the Docker network created and used by M3TAL for inter-container communication. All M3TAL managed containers connect to this network.
-   **Default Value**: `m3tal`
-   **Example Value**: `m3tal_proxy_network`
-   **Component(s) Used By**: All Docker Compose Stacks (as an external network named `proxy`)

#### `LOCAL_IP`
-   **Description**: The IP address of the host machine. This variable is used by containers to refer to services running directly on the host (e.g., the `m3tal-api.service`) via `host.docker.internal`.
-   **Default Value**: `127.0.0.1`
-   **Example Value**: `192.168.1.100`
-   **Component(s) Used By**: `m3tal-api.service` (for binding), Traefik (for routing to host-bound API)

#### `PUID`
-   **Description**: The User ID (UID) that Docker containers will use to run their processes. This ensures proper file permissions when containers interact with host-mounted volumes. It should match the UID of a user on your host system.
-   **Default Value**: `1000`
-   **Example Value**: `1001`
-   **Component(s) Used By**: All Docker Compose Containers

#### `PGID`
-   **Description**: The Group ID (GID) that Docker containers will use to run their processes. This ensures proper file permissions when containers interact with host-mounted volumes. It should match the GID of a group on your host system.
-   **Default Value**: `1000`
-   **Example Value**: `1001`
-   **Component(s) Used By**: All Docker Compose Containers

#### `TZ`
-   **Description**: The timezone for containers, ensuring correct time display and accurate timestamping in logs within the containerized environment.
-   **Default Value**: `America/Denver`
-   **Example Value**: `Europe/London`
-   **Component(s) Used By**: All Docker Compose Containers

### Authentication & Security

These variables manage access and security for the M3TAL Dashboard and API.

#### `DASHBOARD_SECRET`
-   **Description**: A secret key used by the M3TAL Dashboard for session management, encryption, and other security-sensitive operations.
-   **Default Value**: `change_me_immediately`
-   **Note**: This variable is **auto-generated** on first `m3tal init`. Users should **NOT** set it manually unless performing a secret rotation.
-   **Example Value**: `aVeryStrongRandomSecretKey123!_generated_by_m3tal`
-   **Component(s) Used By**: `m3tal-dashboard`

#### `API_TOKEN`
-   **Description**: A bearer token used for authenticating requests to the M3TAL API. This token is used internally by the CLI and Dashboard to communicate with the API.
-   **Default Value**: `change_me_api_token`
-   **Note**: This variable is **auto-generated** on first `m3tal init`. Users should **NOT** set it manually unless performing a token rotation.
-   **Example Value**: `someLongRandomStringForAuth_generated_by_m3tal`
-   **Component(s) Used By**: `m3tal-api.service`, `m3tal` CLI, `m3tal-dashboard`

#### `ADMIN_PASSWORD`
-   **Description**: The initial password for the default administrator user in the M3TAL Dashboard. This is used when the `users.json` file is first created or empty. It can be changed via `m3tal dashpass`.
-   **Default Value**: `admin_pass`
-   **Example Value**: `MySecureAdminPass123`
-   **Component(s) Used By**: `m3tal-dashboard`

### Storage Paths

These variables define the host filesystem paths where M3TAL-managed services store their data.

#### `BASE_STORAGE_PATH`
-   **Description**: The root directory on the host filesystem where M3TAL stores all persistent data, including media, configuration, and downloads.
-   **Default Value**: `./data` (This is the default in the template, but **`/mnt` is the default in production deployments** for mounted storage).
-   **Example Value**: `/mnt/m3tal_data`
-   **Component(s) Used By**: All Docker Compose Stacks (for volume mounts)

#### `MEDIA_PATH`
-   **Description**: The path within `BASE_STORAGE_PATH` designated for storing media files (e.g., movies, TV shows, music). This path is typically mounted into media server containers.
-   **Default Value**: `${BASE_STORAGE_PATH}/media` (e.g., `/mnt/media`)
-   **Example Value**: `/mnt/storage/media`
-   **Component(s) Used By**: User-defined media server stacks

#### `CONFIG_PATH`
-   **Description**: The path within `BASE_STORAGE_PATH` designated for storing configuration files for various services (e.g., `users.json` for the M3TAL Dashboard, or configuration for other user-defined stacks).
-   **Default Value**: `${BASE_STORAGE_PATH}/config` (e.g., `/mnt/config`)
-   **Example Value**: `/mnt/m3tal/configs`
-   **Component(s) Used By**: `m3tal-dashboard`, User-defined stacks

#### `DOWNLOADS_PATH`
-   **Description**: The path within `BASE_STORAGE_PATH` designated for downloaded content. This path is typically mounted into download client containers.
-   **Default Value**: `${BASE_STORAGE_PATH}/downloads` (e.g., `/mnt/downloads`)
-   **Example Value**: `/mnt/storage/downloads`
-   **Component(s) Used By**: User-defined download client stacks

### Dashboard Specifics

These variables specifically control the M3TAL Dashboard container's behavior and exposure.

#### `DASHBOARD_PORT`
-   **Description**: The internal port on which the M3TAL Dashboard container listens. When `DASHBOARD_EXPOSE_MODE` is `local`, this port is directly mapped to the host.
-   **Default Value**: `8082`
-   **Example Value**: `8082`
-   **Component(s) Used By**: `m3tal-dashboard`

#### `DASHBOARD_EXPOSE_MODE`
-   **Description**: Controls how the M3TAL Dashboard is made accessible.
    -   `local`: The dashboard is exposed directly via a port binding (`${DASHBOARD_PORT}:8082`) on the host. Access via `http://HOST_IP:${DASHBOARD_PORT}`. Best for LAN-only setups.
    -   `traefik`: The dashboard is exposed via Traefik, using a domain-based route (`http://dash.DOMAIN`). Traefik must be running. Best for domain-based setups.
-   **Default Value**: `local`
-   **Example Value**: `traefik`
-   **Component(s) Used By**: `m3tal` CLI (to select `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` overrides), `m3tal-dashboard`

### Network & Routing (Traefik)

These variables configure the Traefik reverse proxy and domain-based routing.

#### `DOMAIN`
-   **Description**: The primary domain name for M3TAL services. **Setting this variable enables Traefik to automatically route traffic for `dash.DOMAIN` (to the dashboard) and `api.DOMAIN` (to the API daemon).** This requires Traefik to be running.
-   **Default Value**: `localhost`
-   **Example Value**: `example.com`
-   **Component(s) Used By**: Traefik gateway, `m3tal-dashboard` (when `DASHBOARD_EXPOSE_MODE=traefik`)

#### `TRAEFIK_WEB_PORT`
-   **Description**: The port on the host that Traefik binds for HTTP traffic. This is the primary entrypoint for web services.
-   **Default Value**: `80`
-   **Example Value**: `80`
-   **Component(s) Used By**: Traefik gateway

#### `TRAEFIK_WEBHTTPS_PORT`
-   **Description**: The port on the host that Traefik binds for HTTPS traffic. Enabling HTTPS requires additional Traefik configuration (e.g., Let's Encrypt).
-   **Default Value**: `443`
-   **Example Value**: `443`
-   **Component(s) Used By**: Traefik gateway

#### `TRAEFIK_DASHBOARD_PORT`
-   **Description**: The internal port where Traefik's own management dashboard is exposed within its container. This is typically mapped to `127.0.0.1:8081` on the host.
-   **Default Value**: `8080`
-   **Example Value**: `8080`
-   **Component(s) Used By**: Traefik gateway

### VPN Integration

These variables are used for integrating VPN services (e.g., WireGuard or OpenVPN clients) with M3TAL-managed stacks. They are not directly used by the M3TAL core components but are available for user stacks.

#### `VPN_USER`
-   **Description**: The username for VPN services integrated with M3TAL. This is typically used by containers that act as VPN clients.
-   **Default Value**: `user`
-   **Example Value**: `m3tal_vpn_user`
-   **Component(s) Used By**: VPN client containers (e.g., WireGuard, OpenVPN)

#### `VPN_PASSWORD`
-   **Description**: The password for VPN services integrated with M3TAL. This is typically used by containers that act as VPN clients.
-   **Default Value**: `password`
-   **Example Value**: `MyVpnSecret!`
-   **Component(s) Used By**: VPN client containers (e.g., WireGuard, OpenVPN)