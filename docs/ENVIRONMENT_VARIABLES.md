As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the complete reference for M3TAL environment variables.

---

# M3TAL Environment Variables Reference

All M3TAL configuration is managed through environment variables, primarily defined in the `/etc/m3tal/.env` file. This file serves as the single source of truth for both the `m3tal` CLI and all Docker Compose stacks, which load these variables via the `--env-file` option.

Understanding and managing these variables is key to customizing your M3TAL deployment, from network routing to storage paths and access controls.

---

## Quick Reference Table

| Variable                | Description                                                | Default              | Components                                     |
| :---------------------- | :--------------------------------------------------------- | :------------------- | :--------------------------------------------- |
| `DASHBOARD_PORT`        | Dashboard container port & local host port                 | `8082`               | `m3tal-dashboard`, CLI                         |
| `DASHBOARD_EXPOSE_MODE` | Controls dashboard exposure (`local` or `traefik`)         | `local`              | CLI, `m3tal-dashboard`                         |
| `HTTP_PORT`             | M3TAL API daemon listening port                            | `8080`               | `m3tal-api` daemon, Traefik, `m3tal-dashboard` |
| `STATE_DIR`             | Base directory for component state                         | `./state`            | `m3tal-dashboard`, `m3tal-api` daemon          |
| `LOG_LEVEL`             | Minimum logging verbosity                                  | `info`               | CLI, `m3tal-api` daemon                        |
| `DASHBOARD_SECRET`      | **CRITICAL:** Dashboard session secret                     | `change_me_immediately` | `m3tal-dashboard`                              |
| `API_TOKEN`             | **CRITICAL:** API daemon authentication token              | `change_me_api_token` | CLI, `m3tal-api` daemon                        |
| `ADMIN_PASSWORD`        | Default password for admin user                            | `admin_pass`         | `m3tal-dashboard`, `m3tal-api` daemon          |
| `NETWORK_NAME`          | Default Docker network name for user stacks                | `m3tal`              | Docker Compose (user stacks)                   |
| `LOCAL_IP`              | Host's local IP address                                    | `127.0.0.1`          | Traefik, `m3tal-dashboard`                     |
| `DOMAIN`                | **CRITICAL:** Base domain for Traefik routing              | `localhost`          | Traefik, `m3tal-dashboard`, `m3tal-api` daemon |
| `VPN_USER`              | Username for VPN services                                  | `user`               | VPN containers                                 |
| `VPN_PASSWORD`          | Password for VPN services                                  | `password`           | VPN containers                                 |
| `BASE_STORAGE_PATH`     | **CRITICAL:** Base host path for all data                  | `./data`             | All data-persisting containers                 |
| `MEDIA_PATH`            | Host path for media data                                   | `./data/media`       | Media-related containers                       |
| `CONFIG_PATH`           | Host path for configuration data                           | `./data/config`      | `m3tal-dashboard`, Traefik                     |
| `DOWNLOADS_PATH`        | Host path for downloaded content                           | `./data/downloads`   | Download client containers                     |
| `PUID`                  | User ID for container file permissions                     | `1000`               | Most containers                                |
| `PGID`                  | Group ID for container file permissions                    | `1000`               | Most containers                                |
| `TZ`                    | Timezone for containers                                    | `America/Denver`     | Most containers                                |
| `TRAEFIK_WEB_PORT`      | Traefik's host port for HTTP traffic                       | `80`                 | `traefik` gateway                              |
| `TRAEFIK_WEBHTTPS_PORT` | Traefik's host port for HTTPS traffic                      | `443`                | `traefik` gateway                              |
| `TRAEFIK_DASHBOARD_PORT`| Traefik's internal dashboard container port                | `8080`               | `traefik` gateway                              |
| `DEBUG_MODE`            | Enables debug features and logging                         | `false`              | CLI, `m3tal-api` daemon                        |
| `METRICS_ENABLED`       | Controls exposure of API daemon metrics                    | `true`               | `m3tal-api` daemon                             |

---

## Environment Variable Details

### Core M3TAL Configuration

These variables control fundamental aspects of the M3TAL system, including API communication and logging.

#### `HTTP_PORT`

*   **Description**: The port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming HTTP requests. This port is crucial for internal communication from the dashboard and external routing via Traefik.
*   **Default Value**: `8080`
*   **Example Value**: `5050`
*   **Used By**: `m3tal-api` daemon, Traefik gateway (for routing `api.${DOMAIN}`), `m3tal-dashboard` (to reach the API).

#### `STATE_DIR`

*   **Description**: The base directory within a container where M3TAL components are configured to store their runtime state and temporary data. While the API daemon's primary SQLite database is at `/var/lib/m3tal/state.db` on the host, this variable helps define internal container paths.
*   **Default Value**: `./state`
*   **Example Value**: `/var/lib/m3tal/runtime-state`
*   **Used By**: `m3tal-dashboard` (internal path `/docker/state`), CLI.

#### `LOG_LEVEL`

*   **Description**: Sets the minimum log level for output generated by the M3TAL CLI and API daemon. Adjust this for more verbose debugging or to filter down to critical issues.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Used By**: CLI binary, `m3tal-api` daemon.

### Authentication & Access Control

These variables manage authentication tokens and passwords for securing access to M3TAL services.

#### `DASHBOARD_SECRET`

*   **Description**: **CRITICAL**: A strong, randomly generated secret key used by the M3TAL Dashboard for secure session management (e.g., Flask session cookies). **This value is automatically generated on the first `m3tal init` command. Users should NOT set this manually unless performing a security rotation.**
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_very_long_and_complex_random_string_of_characters_12345`
*   **Used By**: `m3tal-dashboard`.

#### `API_TOKEN`

*   **Description**: **CRITICAL**: An authentication token used by the `m3tal` CLI and any other external clients to authenticate with the M3TAL API daemon. This token grants administrative access to the API. **This value is automatically generated on the first `m3tal init` command. Users should NOT set this manually unless performing a security rotation.**
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwic3ViIjoxMjM0NTY3ODkwLCJpYXQiOjE1MTYyMzkwMjJ9.SflKxwRJSMeKKF2QT4fWPW3K`
*   **Used By**: CLI binary, `m3tal-api` daemon.

#### `ADMIN_PASSWORD`

*   **Description**: The default password for the initial "admin" user of the M3TAL Dashboard. You will be prompted to change this during the initial setup process. It may also be used for initial API authentication.
*   **Default Value**: `admin_pass`
*   **Example Value**: `my_very_secure_admin_password`
*   **Used By**: `m3tal-dashboard` (for `/docker/users.json`), `m3tal-api` daemon (for initial authentication checks).

### Network Configuration

Variables for defining network settings, including Docker networks and domain-based routing.

#### `NETWORK_NAME`

*   **Description**: The default name for the internal Docker network created and used by M3TAL for user-deployed Docker Compose stacks to facilitate inter-container communication. Note that core M3TAL services like the dashboard and Traefik typically use an external network named `proxy`.
*   **Default Value**: `m3tal`
*   **Example Value**: `my_custom_network`
*   **Used By**: Docker Compose (for user-defined stacks).

#### `LOCAL_IP`

*   **Description**: The host's local IP address. This can be used by M3TAL components for internal communication or for services that need to bind to a specific host IP.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100`
*   **Used By**: Traefik gateway (for routing to `host.docker.internal`), `m3tal-dashboard` (indirectly via `host.docker.internal`).

#### `DOMAIN`

*   **Description**: **CRITICAL**: The base domain name that M3TAL uses for its services. **Setting this variable enables Traefik to automatically configure routing rules for `dash.${DOMAIN}` (M3TAL Dashboard) and `api.${DOMAIN}` (M3TAL API daemon).** This is essential for exposing your services via a custom domain.
*   **Default Value**: `localhost`
*   **Example Value**: `example.com`
*   **Used By**: Traefik gateway, `m3tal-dashboard` (via Traefik labels), `m3tal-api` daemon (via Traefik dynamic configuration).

### Storage Paths

These variables define the host filesystem paths where M3TAL stores various types of data.

#### `BASE_STORAGE_PATH`

*   **Description**: The fundamental host path where all M3TAL-related media, configuration, and other persistent data are stored. All other storage paths (`MEDIA_PATH`, `CONFIG_PATH`, `DOWNLOADS_PATH`) are typically defined relative to this base path. **In production deployments, this often defaults to `/mnt` to leverage dedicated storage volumes, rather than the `./data` default used in development/templates.**
*   **Default Value**: `./data`
*   **Example Value**: `/srv/m3tal` or `/mnt/m3tal`
*   **Used By**: `m3tal-dashboard`, `ollama`, and any other container that requires persistent storage.

#### `MEDIA_PATH`

*   **Description**: The host path designated for media-related data, such as user-uploaded files, video libraries, or music collections.
*   **Default Value**: `./data/media`
*   **Example Value**: `${BASE_STORAGE_PATH}/media` (e.g., `/mnt/media`)
*   **Used By**: Various user-deployed media containers (e.g., Plex, Jellyfin, photo managers).

#### `CONFIG_PATH`

*   **Description**: The host path where M3TAL component configurations and persistent state files are stored. This includes configurations for the dashboard and Traefik's dynamic rules.
*   **Default Value**: `./data/config`
*   **Example Value**: `${BASE_STORAGE_PATH}/config` (e.g., `/mnt/config`)
*   **Used By**: `m3tal-dashboard`, Traefik gateway (for dynamic configuration files), `m3tal-api` daemon (for `/etc/m3tal` symlinks).

#### `DOWNLOADS_PATH`

*   **Description**: The host path specifically for downloaded content, such as files from torrent clients or direct downloads.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `${BASE_STORAGE_PATH}/downloads` (e.g., `/mnt/downloads`)
*   **Used By**: Download client containers (e.g., qBittorrent, Transmission).

### Traefik Gateway Configuration

Variables specific to the Traefik reverse proxy, controlling its exposure and dashboard.

#### `DASHBOARD_PORT`

*   **Description**: This variable defines the internal port that the M3TAL Dashboard container listens on (always `8082`). When `DASHBOARD_EXPOSE_MODE` is set to `local`, this variable also determines the host port the dashboard is directly exposed on.
*   **Default Value**: `8082`
*   **Example Value**: `9000` (for direct host exposure)
*   **Used By**: `m3tal-dashboard` container, `m3tal dash up` (CLI).

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Controls how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard is exposed directly on the host's `DASHBOARD_PORT` (e.g., `http://HOST_IP:8082`). Ideal for LAN-only setups or initial configuration without Traefik.
    *   `traefik`: The dashboard is exposed via Traefik, routing `dash.${DOMAIN}` to the container. Requires Traefik to be running.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Used By**: `m3tal dash up` (CLI), `m3tal-dashboard` (via Docker Compose overrides).

#### `TRAEFIK_WEB_PORT`

*   **Description**: The host port on which the Traefik gateway listens for standard HTTP (web) traffic. This is the primary entry point for domain-based services.
*   **Default Value**: `80`
*   **Example Value**: `8080` (if port 80 is in use by another service on the host)
*   **Used By**: `traefik` gateway.

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port on which the Traefik gateway listens for HTTPS (websecure) traffic.
*   **Default Value**: `443`
*   **Example Value**: `8443` (if port 443 is in use)
*   **Used By**: `traefik` gateway.

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The *internal container port* for Traefik's administrative dashboard. By default, this internal port (`8080`) is mapped to host port `8081` and accessible only locally (`127.0.0.1:8081`).
*   **Default Value**: `8080`
*   **Example Value**: `9090` (typically not changed as the host mapping is usually fixed).
*   **Used By**: `traefik` gateway.

### VPN Configuration

Variables used for configuring optional VPN services, such as Cloudflared tunnels or dedicated VPN client containers.

#### `VPN_USER`

*   **Description**: Username credential for authenticating with VPN services or secure tunnels, such as those provided by Cloudflared or a dedicated VPN container.
*   **Default Value**: `user`
*   **Example Value**: `vpnadmin`
*   **Used By**: VPN containers (e.g., `cloudflared` if configured for authenticated tunnels, or other VPN stacks).

#### `VPN_PASSWORD`

*   **Description**: Password credential for authenticating with VPN services or secure tunnels.
*   **Default Value**: `password`
*   **Example Value**: `secure_vpn_pass_123`
*   **Used By**: VPN containers (e.g., `cloudflared` if configured for authenticated tunnels, or other VPN stacks).

### System & Utility Variables

General system-wide variables for user/group IDs, timezone, and debugging.

#### `PUID`

*   **Description**: The User ID (UID) that Docker containers should use when running applications. This is crucial for ensuring proper file ownership and permissions on host volumes mounted into containers. It should generally match the UID of the user who owns your M3TAL data directories on the host system.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Used By**: `m3tal-dashboard`, `ollama`, and most user-deployed containers.

#### `PGID`

*   **Description**: The Group ID (GID) that Docker containers should use. Similar to `PUID`, this ensures correct group ownership and permissions for files on mounted host volumes. It should typically match the GID of the group that owns your M3TAL data directories.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Used By**: `m3tal-dashboard`, `ollama`, and most user-deployed containers.

#### `TZ`

*   **Description**: Sets the timezone for M3TAL containers, ensuring accurate logging and time-sensitive operations. This should be set to a valid IANA timezone identifier (e.g., `Europe/London`, `Asia/Tokyo`).
*   **Default Value**: `America/Denver`
*   **Example Value**: `America/New_York`
*   **Used By**: `m3tal-dashboard`, `ollama`, `traefik`, `cloudflared`, and most user-deployed containers.

#### `DEBUG_MODE`

*   **Description**: Enables or disables debug-specific features and highly verbose logging within the M3TAL CLI and API daemon. Set to `true` for detailed troubleshooting information.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Used By**: CLI binary, `m3tal-api` daemon.

#### `METRICS_ENABLED`

*   **Description**: Controls whether the M3TAL API daemon exposes Prometheus-compatible metrics endpoints, allowing for monitoring of M3TAL's internal performance.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Used By**: `m3tal-api` daemon.