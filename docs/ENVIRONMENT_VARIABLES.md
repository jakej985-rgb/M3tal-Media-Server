# Environment Variables Reference

As the M3TAL Ecosystem Documentation Architect, it's my duty to ensure you have a clear understanding of your system's configuration. This document serves as a comprehensive reference for all environment variables that govern the behavior of your M3TAL system.

All M3TAL environment variables are read from the primary configuration file: `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (including core M3TAL components and user-defined services) load these variables via `--env-file` for consistent configuration across the entire ecosystem.

---

### Quick Reference

| Variable Name           | Default Value         | Description                                                      |
| :---------------------- | :-------------------- | :--------------------------------------------------------------- |
| `API_TOKEN`             | `change_me_api_token` | Token for CLI authentication with the M3TAL API.                 |
| `ADMIN_PASSWORD`        | `admin_pass`          | Initial password for the dashboard admin user.                   |
| `BASE_STORAGE_PATH`     | `./data`              | Root directory for all persistent M3TAL data.                    |
| `CONFIG_PATH`           | `./data/config`       | Path for application configuration files.                        |
| `DASHBOARD_EXPOSE_MODE` | `local`               | How the dashboard is exposed (direct port or via Traefik).       |
| `DASHBOARD_PORT`        | `8082`                | Direct port for the M3TAL Dashboard in `local` mode.             |
| `DASHBOARD_SECRET`      | `change_me_immediately` | Secret key for dashboard session management.                     |
| `DEBUG_MODE`            | `false`               | Enables debug logging and features.                              |
| `DOMAIN`                | `localhost`           | Base domain for Traefik routing.                                 |
| `DOWNLOADS_PATH`        | `./data/downloads`    | Path for downloaded files.                                       |
| `HTTP_PORT`             | `8080`                | Port for the M3TAL API daemon.                                   |
| `LOCAL_IP`              | `127.0.0.1`           | Local IP address for host-level communication.                   |
| `LOG_LEVEL`             | `info`                | Logging verbosity for the M3TAL API daemon.                      |
| `MEDIA_PATH`            | `./data/media`        | Path for media files.                                            |
| `METRICS_ENABLED`       | `true`                | Enables/disables exposing system and API metrics.                |
| `NETWORK_NAME`          | `m3tal`               | Name of the default Docker network.                              |
| `PGID`                  | `1000`                | Group ID (GID) for containers.                                   |
| `PUID`                  | `1000`                | User ID (UID) for containers.                                    |
| `STATE_DIR`             | `./state`             | Internal dashboard container path for state files.               |
| `TRAEFIK_DASHBOARD_PORT`| `8080`                | Internal Traefik dashboard port.                                 |
| `TRAEFIK_WEB_PORT`      | `80`                  | Port for Traefik's HTTP entry point.                             |
| `TRAEFIK_WEBHTTPS_PORT` | `443`                 | Port for Traefik's HTTPS entry point.                            |
| `TZ`                    | `America/Denver`      | Timezone for the system and containers.                          |
| `VPN_PASSWORD`          | `password`            | Password for potential VPN integration.                          |
| `VPN_USER`              | `user`                | Username for potential VPN integration.                          |

---

### Detailed Variable Reference

#### Core

These variables control fundamental aspects of the M3TAL API daemon and container runtime behavior.

*   **`LOG_LEVEL`**
    *   **Description:** Sets the logging verbosity for the `m3tal-api` daemon.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api` daemon
*   **`DEBUG_MODE`**
    *   **Description:** When set to `true`, enables debug logging and potentially other debug-specific features within M3TAL components.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api` daemon, `m3tal-cli`
*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether system and API metrics are exposed by the `m3tal-api` daemon, typically for Prometheus scraping.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api` daemon
*   **`PUID`**
    *   **Description:** The User ID (UID) that containers will use to run processes and access files, ensuring proper file permissions on mounted volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** All Docker containers (e.g., `m3tal-dashboard`)
*   **`PGID`**
    *   **Description:** The Group ID (GID) that containers will use to run processes and access files, ensuring proper file permissions on mounted volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** All Docker containers (e.g., `m3tal-dashboard`)
*   **`TZ`**
    *   **Description:** Sets the timezone for all M3TAL components and containers, ensuring consistent time reporting.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `America/New_York`
    *   **Used By:** All Docker containers (e.g., `m3tal-dashboard`), `m3tal-api` daemon

#### Auth

These variables control authentication and security for M3TAL's core services.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used by the `m3tal-dashboard` for session management, cookie signing, and other security-sensitive operations.
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `my_long_random_dashboard_secret_key`
    *   **Used By:** `m3tal-dashboard` container
    *   **Note:** This variable is **auto-generated** by `m3tal init` on first setup. Users should **NOT** set it manually unless performing a secret rotation.
*   **`API_TOKEN`**
    *   **Description:** The authentication token used by the `m3tal-cli` to securely communicate with and authorize requests to the `m3tal-api` daemon.
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `a_very_long_and_secure_api_token_string`
    *   **Used By:** `m3tal-cli`, `m3tal-api` daemon
    *   **Note:** This variable is **auto-generated** by `m3tal init` on first setup. Users should **NOT** set it manually unless performing a token rotation.
*   **`ADMIN_PASSWORD`**
    *   **Description:** The initial default password for the administrator user account created within the `m3tal-dashboard`.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `MyStrongAndComplexAdminPassword123!`
    *   **Used By:** `m3tal-dashboard` container

#### Network

These variables configure network settings for M3TAL components and their interaction.

*   **`DASHBOARD_PORT`**
    *   **Description:** The host port on which the `m3tal-dashboard` is directly exposed when `DASHBOARD_EXPOSE_MODE` is set to `local`.
    *   **Default Value:** `8082`
    *   **Example Value:** `8083`
    *   **Used By:** `m3tal-cli` (when managing `m3tal-dashboard` in `local` mode)
*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Determines how the `m3tal-dashboard` container is exposed to the network.
        *   `local`: Direct port binding on the host (`DASHBOARD_PORT`). No Traefik required.
        *   `traefik`: Exposed via Traefik using `dash.${DOMAIN}`. Requires Traefik to be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-cli` (specifically `m3tal dash up` command)
*   **`HTTP_PORT`**
    *   **Description:** The port on which the `m3tal-api` daemon listens for incoming HTTP requests. This port is typically only accessible from `host.docker.internal` or `localhost`.
    *   **Default Value:** `8080`
    *   **Example Value:** `9000`
    *   **Used By:** `m3tal-api` daemon
*   **`NETWORK_NAME`**
    *   **Description:** The name of the default Docker network created and managed by M3TAL. All M3TAL core components and user-defined Docker Compose stacks should connect to this network for inter-service communication.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `my_custom_m3tal_network`
    *   **Used By:** `m3tal-cli` (for Docker network creation), all Docker Compose stacks
*   **`LOCAL_IP`**
    *   **Description:** Specifies the local IP address for host-level communication. While `host.docker.internal` typically points to the host's primary IP, this variable can inform other services or configurations that might rely on a specific host IP.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** Potentially various Docker Compose configurations for host binding.

#### Storage

These variables define the filesystem paths where M3TAL stores its persistent data.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory on the host filesystem where all M3TAL-related persistent data, including media, configuration, and downloads, is stored. All other storage paths (e.g., `MEDIA_PATH`, `CONFIG_PATH`) are relative to or contained within this base path.
    *   **Default Value:** `./data` (relative to `/etc/m3tal/`)
    *   **Example Value:** `/mnt/m3tal_data`
    *   **Used By:** All Docker Compose stacks (for volume mounts), `m3tal-cli` (for initial directory setup)
    *   **Note:** In production deployments, this typically defaults to a dedicated mount point like `/mnt` (e.g., `/mnt/data`) rather than a relative path, to ensure data persistence outside the system directory.
*   **`MEDIA_PATH`**
    *   **Description:** The sub-path within `BASE_STORAGE_PATH` designated for storing media files (e.g., movies, TV shows, music) consumed by user-defined media stacks.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Used By:** User-defined media stacks (e.g., Plex, Jellyfin, Radarr, Sonarr)
*   **`CONFIG_PATH`**
    *   **Description:** The sub-path within `BASE_STORAGE_PATH` designated for application configuration files and persistent state that M3TAL components and user services require.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Used By:** `m3tal-dashboard` (for mounting `/docker/state`), all Docker Compose stacks (for application-specific config data)
*   **`DOWNLOADS_PATH`**
    *   **Description:** The sub-path within `BASE_STORAGE_PATH` designated for downloaded files, typically managed by download clients in user-defined stacks.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Used By:** User-defined download stacks (e.g., qBittorrent, Transmission)
*   **`STATE_DIR`**
    *   **Description:** This variable defines the *internal* directory path within the `m3tal-dashboard` container where state files, such as `users.json`, are managed. On the host, this path is mapped from `${CONFIG_PATH}/m3tal/state` to `/docker/state` inside the container.
    *   **Default Value:** `./state`
    *   **Example Value:** `dashboard_state` (This value is used internally by the dashboard container; the host path is controlled by `CONFIG_PATH`).
    *   **Used By:** `m3tal-dashboard` container (as an internal environment variable, usually `STATE_DIR=/docker/state`)

#### Traefik

These variables configure the Traefik reverse proxy for M3TAL's routing.

*   **`DOMAIN`**
    *   **Description:** The base domain name that Traefik will use to route requests to M3TAL services. Setting this value enables convenient access via subdomains like `dash.DOMAIN` for the dashboard and `api.DOMAIN` for the M3TAL API.
    *   **Default Value:** `localhost`
    *   **Example Value:** `example.com`
    *   **Used By:** `m3tal-cli`, `traefik` container, `m3tal-dashboard` (via Traefik labels)
    *   **Note:** Controls Traefik routing rules. Setting it enables `dash.DOMAIN` and `api.DOMAIN` routes.
*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik's HTTP entry point listens on. This is where external HTTP requests are received.
    *   **Default Value:** `80`
    *   **Example Value:** `8080` (if port 80 is already in use)
    *   **Used By:** `traefik` container
*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik's HTTPS entry point listens on. This is where external HTTPS requests are received.
    *   **Default Value:** `443`
    *   **Example Value:** `8443`
    *   **Used By:** `traefik` container
*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which the Traefik management dashboard listens within its container. This is typically mapped to `127.0.0.1:8081` on the host by default for host-local access.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** `traefik` container

#### VPN

These variables are placeholders for potential future VPN integration or custom network configurations.

*   **`VPN_USER`**
    *   **Description:** Placeholder for a username that might be used for authenticating with a VPN service integrated with M3TAL (e.g., Cloudflared authentication, if custom configurations require it).
    *   **Default Value:** `user`
    *   **Example Value:** `my_vpn_user`
    *   **Used By:** `cloudflared` (potential), future VPN services
*   **`VPN_PASSWORD`**
    *   **Description:** Placeholder for a password associated with `VPN_USER`, for authenticating with a VPN service.
    *   **Default Value:** `password`
    *   **Example Value:** `my_secure_vpn_password`
    *   **Used By:** `cloudflared` (potential), future VPN services

---