# Environment Variables Reference

All M3TAL components, including the CLI (`m3tal` binary) and all Docker Compose stacks, read their configuration from a single source of truth: the `.env` file located at `/etc/m3tal/.env`. This file is paramount to your M3TAL ecosystem's operation.

The `m3tal config wizard` and `m3tal config set` commands are the recommended way to manage these variables. When `m3tal up` is executed, Docker Compose automatically loads these variables via the `--env-file` option, ensuring consistency across all running services.

---

## Quick Reference Table

| Variable Name             | Default Value          | Description                                                                                             | Category | Component(s)       |
| :------------------------ | :--------------------- | :------------------------------------------------------------------------------------------------------ | :------- | :----------------- |
| `HTTP_PORT`               | `8080`                 | Port for the M3TAL API daemon.                                                                          | Core     | API Daemon         |
| `STATE_DIR`               | `./state`              | Directory for API daemon's SQLite state database and runtime files.                                     | Core     | API Daemon         |
| `LOG_LEVEL`               | `info`                 | Verbosity of M3TAL API daemon logging.                                                                  | Core     | API Daemon         |
| `DEBUG_MODE`              | `false`                | Enables debug logging and features for the API daemon.                                                  | Core     | API Daemon         |
| `METRICS_ENABLED`         | `true`                 | Enables Prometheus metrics exposure for the API daemon.                                                 | Core     | API Daemon         |
| `DASHBOARD_SECRET`        | `change_me_immediately`| Secret for Dashboard session management. **Auto-generated.**                                            | Auth     | Dashboard          |
| `API_TOKEN`               | `change_me_api_token`  | Token for M3TAL API authentication. **Auto-generated.**                                                 | Auth     | API Daemon         |
| `ADMIN_PASSWORD`          | `admin_pass`           | Default password for Dashboard's initial admin user.                                                    | Auth     | Dashboard          |
| `DASHBOARD_PORT`          | `8082`                 | Internal and host-exposed port for the M3TAL Dashboard.                                                 | Network  | Dashboard, CLI     |
| `DASHBOARD_EXPOSE_MODE`   | `local`                | How the Dashboard is exposed: `local` (direct port) or `traefik` (via Traefik).                         | Network  | CLI, Dashboard     |
| `NETWORK_NAME`            | `m3tal`                | Name of the Docker network for inter-container communication.                                           | Network  | All Containers     |
| `LOCAL_IP`                | `127.0.0.1`            | Host IP address for internal routing (e.g., Traefik to API daemon).                                     | Network  | Traefik            |
| `BASE_STORAGE_PATH`       | `./data`               | Base directory for all persistent data. **Defaults to `/mnt` in production.**                           | Storage  | All Containers     |
| `MEDIA_PATH`              | `./data/media`         | Subdirectory for media data.                                                                            | Storage  | User Stacks        |
| `CONFIG_PATH`             | `./data/config`        | Subdirectory for M3TAL configuration files (e.g., `users.json`).                                        | Storage  | Dashboard, API     |
| `DOWNLOADS_PATH`          | `./data/downloads`     | Subdirectory for downloads.                                                                             | Storage  | User Stacks        |
| `DOMAIN`                  | `localhost`            | Base domain for Traefik routing (e.g., `dash.DOMAIN`).                                                  | Traefik  | Traefik, CLI       |
| `TRAEFIK_WEB_PORT`        | `80`                   | Host port for Traefik's HTTP entry point.                                                               | Traefik  | Traefik            |
| `TRAEFIK_WEBHTTPS_PORT`   | `443`                  | Host port for Traefik's HTTPS entry point.                                                              | Traefik  | Traefik            |
| `TRAEFIK_DASHBOARD_PORT`  | `8080`                 | Internal port for the Traefik management dashboard (mapped to `127.0.0.1:8081`).                        | Traefik  | Traefik            |
| `VPN_USER`                | `user`                 | Username for optional VPN services.                                                                     | VPN      | User VPN Stacks    |
| `VPN_PASSWORD`            | `password`             | Password for optional VPN services.                                                                     | VPN      | User VPN Stacks    |
| `PUID`                    | `1000`                 | User ID (UID) for containers to ensure correct file permissions.                                        | System   | All Containers     |
| `PGID`                    | `1000`                 | Group ID (GID) for containers to ensure correct file permissions.                                       | System   | All Containers     |
| `TZ`                      | `America/Denver`       | Timezone for Docker containers.                                                                         | System   | All Containers     |

---

## Detailed Variable Reference

### Core Variables

These variables control the fundamental behavior and logging of the M3TAL API daemon.

#### `HTTP_PORT`
-   **Description**: The network port on which the M3TAL API daemon (the Go binary running as `m3tal-api.service`) listens for incoming HTTP requests.
-   **Default Value**: `8080`
-   **Example Value**: `8080`
-   **Component(s) Used By**: API daemon

#### `STATE_DIR`
-   **Description**: The base directory where the M3TAL API daemon stores its SQLite state database (`state.db`) and other runtime data. While the API daemon typically uses `/var/lib/m3tal/state.db` in a systemd deployment, this variable may influence its internal referencing or be used for CLI operations.
-   **Default Value**: `./state`
-   **Example Value**: `/var/lib/m3tal`, `./state`
-   **Component(s) Used By**: API daemon (for state.db management)

#### `LOG_LEVEL`
-   **Description**: Controls the verbosity of logging output for the M3TAL API daemon. Useful for debugging or reducing log noise.
-   **Default Value**: `info`
-   **Example Value**: `debug`, `info`, `warn`, `error`
-   **Component(s) Used By**: API daemon

#### `DEBUG_MODE`
-   **Description**: A boolean flag that, when set to `true`, enables detailed debug logging and potentially other debug-specific features within the M3TAL API daemon.
-   **Default Value**: `false`
-   **Example Value**: `true`, `false`
-   **Component(s) Used By**: API daemon

#### `METRICS_ENABLED`
-   **Description**: A boolean flag that, when set to `true`, enables the exposure of Prometheus-compatible metrics from the M3TAL API daemon. This allows for monitoring of the M3TAL core services.
-   **Default Value**: `true`
-   **Example Value**: `false`, `true`
-   **Component(s) Used By**: API daemon

### Authentication Variables

These variables are critical for securing access to the M3TAL Dashboard and API.

#### `DASHBOARD_SECRET`
-   **Description**: A unique secret key used by the M3TAL Dashboard for secure session management (e.g., cookies) and encrypting internal data.
    **Note: This value is automatically generated on the first `m3tal init` command. Users should generally NOT set this manually unless performing a rotation or specific recovery procedure.**
-   **Default Value**: `change_me_immediately`
-   **Example Value**: `a_very_long_and_random_string_of_characters_generated_by_m3tal`
-   **Component(s) Used By**: Dashboard

#### `API_TOKEN`
-   **Description**: A token used to authenticate requests made to the M3TAL API daemon, ensuring only authorized clients can interact with the core system.
    **Note: This value is automatically generated on the first `m3tal init` command. Users should generally NOT set this manually unless performing a rotation or specific recovery procedure.**
-   **Default Value**: `change_me_api_token`
-   **Example Value**: `your_secure_api_token_here_generated_by_m3tal`
-   **Component(s) Used By**: API daemon (for validating inbound requests), CLI (for authenticating API calls)

#### `ADMIN_PASSWORD`
-   **Description**: The default password for the initial administrator user of the M3TAL Dashboard. It is highly recommended to change this immediately after the first login via the Dashboard UI or using `m3tal dashpass`.
-   **Default Value**: `admin_pass`
-   **Example Value**: `MyStrongAndSecurePass123!`
-   **Component(s) Used By**: Dashboard (for initial user setup)

### Network Variables

These variables define how M3TAL components communicate with each other and are exposed to the network.

#### `DASHBOARD_PORT`
-   **Description**: The internal port on which the `m3tal-dashboard` container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this port is directly mapped to the host's `DASHBOARD_PORT`.
-   **Default Value**: `8082`
-   **Example Value**: `8082`, `9000`
-   **Component(s) Used By**: Dashboard, CLI (`m3tal dash up` in local mode)

#### `DASHBOARD_EXPOSE_MODE`
-   **Description**: Controls how the M3TAL Dashboard is exposed to the host network.
    -   `local`: The dashboard is exposed directly on the host's `DASHBOARD_PORT` (e.g., `http://HOST_IP:8082`). No Traefik is required. Best for LAN-only setups or initial configuration.
    -   `traefik`: The dashboard is exposed via the Traefik reverse proxy at `http://dash.DOMAIN`. Requires Traefik to be running. Best for domain-based access.
-   **Default Value**: `local`
-   **Example Value**: `traefik`
-   **Component(s) Used By**: CLI (`m3tal dash up`), Docker Compose files for the Dashboard

#### `NETWORK_NAME`
-   **Description**: The name of the custom Docker network used by M3TAL for all its Docker containers (e.g., Dashboard, Traefik, Cloudflared, and user stacks). This provides isolated and consistent inter-container communication.
-   **Default Value**: `m3tal`
-   **Example Value**: `m3tal-proxy`, `my-custom-net`
-   **Component(s) Used By**: All Docker containers

#### `LOCAL_IP`
-   **Description**: The IP address of the host machine. This is primarily used for internal routing by services like Traefik to communicate with applications running directly on the host, such as the M3TAL API daemon (which listens on `http://host.docker.internal:8080`).
-   **Default Value**: `127.0.0.1`
-   **Example Value**: `192.168.1.100`, `10.0.0.5`
-   **Component(s) Used By**: Traefik (for routing to host services)

### Storage Variables

These variables define the filesystem locations for M3TAL's persistent data.

#### `BASE_STORAGE_PATH`
-   **Description**: The absolute path on the host filesystem that serves as the root directory for all M3TAL persistent data. This includes media, configuration, and downloads, which are typically mounted as volumes into Docker containers.
    **Note: In production deployments, this variable typically defaults to `/mnt` to align with common Linux server data partitioning strategies, rather than the `./data` default used in development/template setups.**
-   **Default Value**: `./data`
-   **Example Value**: `/mnt`, `/data/m3tal`
-   **Component(s) Used By**: All Docker containers (as bind mounts), API daemon

#### `MEDIA_PATH`
-   **Description**: A subdirectory within `BASE_STORAGE_PATH` designated for storing media files (e.g., photos, videos, music).
-   **Default Value**: `./data/media`
-   **Example Value**: `/mnt/media`, `/data/m3tal/media`
-   **Component(s) Used By**: User-defined media services (e.g., Plex, Jellyfin), M3TAL Dashboard (for media management features)

#### `CONFIG_PATH`
-   **Description**: A subdirectory within `BASE_STORAGE_PATH` used to store M3TAL's configuration files, such as `users.json` for dashboard credentials.
-   **Default Value**: `./data/config`
-   **Example Value**: `/mnt/config`, `/data/m3tal/config`
-   **Component(s) Used By**: Dashboard, API daemon

#### `DOWNLOADS_PATH`
-   **Description**: A subdirectory within `BASE_STORAGE_PATH` dedicated to storing downloaded content.
-   **Default Value**: `./data/downloads`
-   **Example Value**: `/mnt/downloads`, `/data/m3tal/downloads`
-   **Component(s) Used By**: User-defined download services (e.g., torrent clients, newsreaders)

### Traefik Variables

These variables configure the Traefik reverse proxy for exposing M3TAL and user services.

#### `DOMAIN`
-   **Description**: The base domain name that Traefik uses for routing rules. Setting this value (e.g., to `example.com`) enables domain-based access such as `dash.example.com` and `api.example.com`. If left at `localhost`, routes will be `dash.localhost` and `api.localhost`.
    **Note: This variable is crucial for enabling Traefik's dynamic routing based on hostnames.**
-   **Default Value**: `localhost`
-   **Example Value**: `example.com`, `my-home-server.local`
-   **Component(s) Used By**: Traefik, CLI (for generating Traefik dynamic configuration)

#### `TRAEFIK_WEB_PORT`
-   **Description**: The host port on which the Traefik container listens for incoming HTTP (non-secure) traffic. This is Traefik's primary entry point for web services.
-   **Default Value**: `80`
-   **Example Value**: `80`, `8080`
-   **Component(s) Used By**: Traefik

#### `TRAEFIK_WEBHTTPS_PORT`
-   **Description**: The host port on which the Traefik container listens for incoming HTTPS (secure) traffic. This is typically used for services configured with SSL/TLS.
-   **Default Value**: `443`
-   **Example Value**: `443`, `8443`
-   **Component(s) Used By**: Traefik

#### `TRAEFIK_DASHBOARD_PORT`
-   **Description**: The internal port on which the Traefik management dashboard listens within its container. By default, Traefik's dashboard is mapped to `127.0.0.1:8081` on the host, making it accessible only locally.
-   **Default Value**: `8080`
-   **Example Value**: `8080`
-   **Component(s) Used By**: Traefik

### VPN Variables

These variables are placeholders for configuring optional VPN services within the M3TAL ecosystem.

#### `VPN_USER`
-   **Description**: The username to be used by optional VPN services deployed as Docker containers within M3TAL.
-   **Default Value**: `user`
-   **Example Value**: `m3taluser`, `myvpnclient`
-   **Component(s) Used By**: User-defined VPN service containers

#### `VPN_PASSWORD`
-   **Description**: The password corresponding to the `VPN_USER` for optional VPN services.
-   **Default Value**: `password`
-   **Example Value**: `MySecureVPNPass!`
-   **Component(s) Used By**: User-defined VPN service containers

### System Variables

These variables control system-level settings for Docker containers and the host environment.

#### `PUID`
-   **Description**: The numeric User ID (UID) that Docker containers should run as. Setting this ensures that files created or modified by containers, especially when using bind mounts, have the correct permissions on the host filesystem. This should match the UID of your unprivileged user.
-   **Default Value**: `1000`
-   **Example Value**: `1000`, `33` (for `www-data` user)
-   **Component(s) Used By**: All Docker containers

#### `PGID`
-   **Description**: The numeric Group ID (GID) that Docker containers should run as. Similar to `PUID`, this ensures correct group permissions on host volumes mounted into containers. This should match the GID of your unprivileged user's primary group.
-   **Default Value**: `1000`
-   **Example Value**: `1000`, `33` (for `www-data` group)
-   **Component(s) Used By**: All Docker containers

#### `TZ`
-   **Description**: Specifies the timezone for Docker containers. Setting this ensures that logs, timestamps, and scheduled tasks within containers reflect the correct local time.
-   **Default Value**: `America/Denver`
-   **Example Value**: `Europe/London`, `Asia/Tokyo`, `America/New_York`
-   **Component(s) Used By**: All Docker containers