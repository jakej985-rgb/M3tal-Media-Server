# Environment Variables Reference

Greetings, fellow engineers and architects! DocSmith here, your guide to the M3TAL Ecosystem. This document details the environment variables that power our system, providing a comprehensive reference for configuration, troubleshooting, and advanced deployments.

All M3TAL environment variables are centrally managed and read from the configuration file located at `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks utilize this file via the `--env-file` option, ensuring a consistent configuration across your entire M3TAL deployment. You can conveniently manage these variables using the `m3tal config wizard` or `m3tal config set KEY value` commands.

---

## Quick Reference

| Variable Name           | Description                                                                                                                                                                                | Default Value          | Example Value            | Component(s)                                   |
| :---------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------- | :----------------------- | :--------------------------------------------- |
| `API_TOKEN`             | Auth token for API daemon. **Auto-generated, do not set manually unless rotating.**                                                                                                        | `change_me_api_token`  | `long_random_token_val`  | `m3tal-api`, CLI, `m3tal-dashboard`            |
| `ADMIN_PASSWORD`        | Default password for the dashboard's initial admin user.                                                                                                                                   | `admin_pass`           | `strong_password`        | `m3tal-dashboard`                              |
| `BASE_STORAGE_PATH`     | Base path for user data/media. Defaults to `/mnt` in production.                                                                                                                           | `./data`               | `/mnt/m3tal`             | `m3tal-dashboard`, user stacks                 |
| `CONFIG_PATH`           | Path for config files and dashboard state.                                                                                                                                                 | `./data/config`        | `/opt/m3tal/config`      | `m3tal-dashboard`                              |
| `DASHBOARD_EXPOSE_MODE` | Controls dashboard exposure: `local` (direct port) or `traefik` (`dash.DOMAIN`).                                                                                                         | `local`                | `traefik`                | `m3tal dash up`                                |
| `DASHBOARD_PORT`        | Port for the `m3tal-dashboard` container.                                                                                                                                                  | `8082`                 | `9000`                   | `m3tal-dashboard`                              |
| `DASHBOARD_SECRET`      | Session secret for `m3tal-dashboard`. **Auto-generated, do not set manually unless rotating.**                                                                                             | `change_me_immediately`| `super_secret_key_123`   | `m3tal-dashboard`                              |
| `DEBUG_MODE`            | Enables debug-level logging and features.                                                                                                                                                  | `false`                | `true`                   | `m3tal-api`, `m3tal-dashboard`                 |
| `DOMAIN`                | Base domain for Traefik routing (`dash.DOMAIN`, `api.DOMAIN`).                                                                                                                             | `localhost`            | `yourdomain.com`         | `traefik`, `routing-compose.yml`               |
| `DOWNLOADS_PATH`        | Path for downloaded files.                                                                                                                                                                 | `./data/downloads`     | `/mnt/downloads`         | User stacks                                    |
| `HTTP_PORT`             | Port for the M3TAL API daemon.                                                                                                                                                             | `8080`                 | `5050`                   | `m3tal-api`, `m3tal-dashboard`                 |
| `LOCAL_IP`              | Host machine's local IP, used for host-gateway.                                                                                                                                            | `127.0.0.1`            | `192.168.1.100`          | `m3tal-dashboard`                              |
| `LOG_LEVEL`             | Logging verbosity (`info`, `debug`, `warn`, `error`).                                                                                                                                    | `info`                 | `debug`                  | `m3tal-api`, `m3tal-dashboard`                 |
| `MEDIA_PATH`            | Path for media files.                                                                                                                                                                      | `./data/media`         | `/mnt/media`             | User stacks                                    |
| `METRICS_ENABLED`       | Enables Prometheus-compatible metrics endpoints.                                                                                                                                           | `true`                 | `false`                  | `m3tal-api`                                    |
| `NETWORK_NAME`          | Name of the internal Docker network for M3TAL services.                                                                                                                                    | `m3tal`                | `m3tal-proxy`            | All containers                                 |
| `PGID`                  | Group ID for containers to ensure host volume permissions.                                                                                                                                 | `1000`                 | `1001`                   | `m3tal-dashboard`, user stacks                 |
| `PUID`                  | User ID for containers to ensure host volume permissions.                                                                                                                                  | `1000`                 | `1001`                   | `m3tal-dashboard`, user stacks                 |
| `STATE_DIR`             | Generic base directory for runtime state. (API uses fixed `/var/lib/m3tal/state.db`).                                                                                                    | `./state`              | `/opt/m3tal/state`       | (Future modules)                               |
| `TRAEFIK_DASHBOARD_PORT`| Internal container port for Traefik's dashboard.                                                                                                                                           | `8080`                 | `8000`                   | `traefik`                                      |
| `TRAEFIK_WEB_PORT`      | Host port for Traefik's HTTP (web) entrypoint.                                                                                                                                             | `80`                   | `8080`                   | `traefik`                                      |
| `TRAEFIK_WEBHTTPS_PORT` | Host port for Traefik's HTTPS (websecure) entrypoint.                                                                                                                                      | `443`                  | `8443`                   | `traefik`                                      |
| `TZ`                    | Timezone for containers.                                                                                                                                                                   | `America/Denver`       | `Europe/London`          | `m3tal-dashboard`, user stacks                 |
| `VPN_PASSWORD`          | Password for VPN client configurations.                                                                                                                                                    | `password`             | `secure_vpn_pass!`       | (Future VPN integrations)                      |
| `VPN_USER`              | Username for VPN client configurations.                                                                                                                                                    | `user`                 | `vpn_user_name`          | (Future VPN integrations)                      |

---

## Detailed Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system's operation.

#### `LOG_LEVEL`
*   **Description**: Sets the logging verbosity for M3TAL components, including the API daemon and dashboard. Useful for debugging or reducing log noise.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Component(s) Used**: `m3tal-api.service`, `m3tal-dashboard` container
*   **Notes**: Common values include `debug`, `info`, `warn`, `error`. Setting to `debug` provides the most detailed output.

#### `STATE_DIR`
*   **Description**: This variable is intended to specify a generic base directory for M3TAL's runtime state. However, the M3TAL API daemon uses a fixed path `/var/lib/m3tal/state.db` for its SQLite database. The `m3tal-dashboard` container sets its *internal* `STATE_DIR` to `/docker/state`, which is mounted from a host path derived from `CONFIG_PATH`. Therefore, this variable's direct impact on current core components is limited, but it might be used by future M3TAL modules or custom stacks.
*   **Default Value**: `./state`
*   **Example Value**: `/opt/m3tal/runtime_state`
*   **Component(s) Used**: (Potentially future M3TAL modules, currently not directly influencing core component host paths)

### Authentication

Variables essential for securing access to M3TAL services.

#### `API_TOKEN`
*   **Description**: An authentication token used to secure access to the M3TAL API daemon. This token is required for the CLI and Dashboard to communicate with the API.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `a_very_long_and_secure_api_token_1a2b3c4d5e6f`
*   **Component(s) Used**: `m3tal-api.service`, `CLI binary`, `m3tal-dashboard` container
*   **Notes**: **This token is auto-generated on first `m3tal init`. Users should NOT set it manually unless rotating the secret for security purposes.**

#### `DASHBOARD_SECRET`
*   **Description**: A secret key used by the M3TAL Dashboard for session management, encryption, and securing internal communications.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `another_super_secret_dashboard_key_abcdefg12345`
*   **Component(s) Used**: `m3tal-dashboard` container
*   **Notes**: **This secret is auto-generated on first `m3tal init`. Users should NOT set it manually unless rotating the secret for security purposes.**

#### `ADMIN_PASSWORD`
*   **Description**: The default password for the initial administrator user account in the M3TAL Dashboard. Users are strongly encouraged to change this immediately upon first login or manage it via the `m3tal dashpass` command.
*   **Default Value**: `admin_pass`
*   **Example Value**: `MyStrongAdminPassword123!`
*   **Component(s) Used**: `m3tal-dashboard` container (used during initial user account setup via `users.json`)

### Network Configuration

These variables define how M3TAL services communicate internally and externally.

#### `DASHBOARD_PORT`
*   **Description**: The internal port on which the M3TAL Dashboard container listens (8082). When `DASHBOARD_EXPOSE_MODE` is `local`, this is the host port it will be directly exposed on.
*   **Default Value**: `8082`
*   **Example Value**: `9000`
*   **Component(s) Used**: `m3tal-dashboard` container, `m3tal dash up` command, `m3tal-compose.local.yml`
*   **Notes**: If you change this, ensure no other service on your host is using the same port.

#### `DASHBOARD_EXPOSE_MODE`
*   **Description**: Controls how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard container's port (`DASHBOARD_PORT`) is directly bound to the host, accessible via `http://HOST_IP:DASHBOARD_PORT`. No Traefik required. Best for LAN-only setups.
    *   `traefik`: Traefik routing rules are enabled, making the dashboard accessible via `http://dash.${DOMAIN}`. Requires Traefik to be running.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Component(s) Used**: `m3tal dash up` command (selects `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`)

#### `HTTP_PORT`
*   **Description**: The port on which the M3TAL API daemon (`m3tal-api.service`) listens for HTTP requests.
*   **Default Value**: `8080`
*   **Example Value**: `5050`
*   **Component(s) Used**: `m3tal-api.service`, `m3tal-dashboard` container (as `GO_API_URL` uses `http://host.docker.internal:${HTTP_PORT}`), `routing-compose.yml` (for Traefik's API routing)

#### `LOCAL_IP`
*   **Description**: The local IP address of the M3TAL host machine. This variable is primarily used for `host.docker.internal` resolution within containers, allowing them to access host-bound services like the M3TAL API daemon.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100`
*   **Component(s) Used**: `m3tal-dashboard` container (`extra_hosts` for `host.docker.internal`), `traefik` container (for routing to API daemon)

#### `NETWORK_NAME`
*   **Description**: The name of the Docker network that M3TAL's internal services and user-deployed stacks use to communicate with each other. All containers requiring inter-service communication should be attached to this network.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal-internal-net`
*   **Component(s) Used**: All M3TAL Docker containers (`m3tal-dashboard`, `traefik`, `cloudflared`, `ollama`, and any user stacks)

### Storage Configuration

Variables defining where M3TAL stores its persistent data on the host filesystem.

#### `BASE_STORAGE_PATH`
*   **Description**: The primary base path on the host filesystem where M3TAL stores all persistent user data, media, and configuration. Other storage-related paths (`MEDIA_PATH`, `CONFIG_PATH`, `DOWNLOADS_PATH`) are often relative to this.
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal_data`
*   **Component(s) Used**: `m3tal-dashboard` container (for `/mnt` volume mount), any user-deployed stacks requiring persistent storage.
*   **Notes**: **In production deployments (e.g., cloud VMs), this variable typically defaults to `/mnt` to leverage large attached block storage, rather than `./data` (relative to `/docker`) as seen in development templates.**

#### `CONFIG_PATH`
*   **Description**: The path on the host filesystem where M3TAL stores configuration files, including the dashboard's `users.json` and other state data. This path is often mounted into containers.
*   **Default Value**: `./data/config`
*   **Example Value**: `/opt/m3tal/config`
*   **Component(s) Used**: `m3tal-dashboard` container (mounts `${CONFIG_PATH}/m3tal/state` to `/docker/state`)

#### `DOWNLOADS_PATH`
*   **Description**: The designated path on the host filesystem where downloaded files are typically stored by M3TAL-managed download clients or other applications. Often relative to `BASE_STORAGE_PATH`.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/downloads`
*   **Component(s) Used**: Any user-deployed download client stacks

#### `MEDIA_PATH`
*   **Description**: The designated path on the host filesystem where media files (videos, music, photos) are typically stored and served from by M3TAL-managed media server applications. Often relative to `BASE_STORAGE_PATH`.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/media`
*   **Component(s) Used**: Any user-deployed media server stacks

### Traefik Gateway

Variables specific to the Traefik reverse proxy and its routing behavior.

#### `DOMAIN`
*   **Description**: The base domain name for M3TAL services when using Traefik as the reverse proxy. Setting this variable enables Traefik to create routing rules for `dash.${DOMAIN}` (for the Dashboard) and `api.${DOMAIN}` (for the API daemon).
*   **Default Value**: `localhost`
*   **Example Value**: `m3tal.example.com`
*   **Component(s) Used**: `traefik` container, `m3tal-compose.traefik.yml`, `routing-compose.yml`, `dynamic/api.yml`
*   **Notes**: If `DOMAIN` is set to `localhost`, Traefik rules will route to `dash.localhost` and `api.localhost`. Make sure your `/etc/hosts` file or DNS resolves these if not using actual domains.

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description**: The internal container port on which the Traefik dashboard listens. This is exposed on the host via `127.0.0.1:8081` by default for local access.
*   **Default Value**: `8080`
*   **Example Value**: `8000`
*   **Component(s) Used**: `traefik` container (`routing-compose.yml`)

#### `TRAEFIK_WEB_PORT`
*   **Description**: The host port on which Traefik listens for incoming unencrypted HTTP traffic. This serves as the 'web' entrypoint for all services.
*   **Default Value**: `80`
*   **Example Value**: `8080`
*   **Component(s) Used**: `traefik` container (`routing-compose.yml`)

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description**: The host port on which Traefik listens for incoming encrypted HTTPS traffic. This serves as the 'websecure' entrypoint.
*   **Default Value**: `443`
*   **Example Value**: `8443`
*   **Component(s) Used**: `traefik` container (`routing-compose.yml`)

### VPN Integration

Variables reserved for future VPN client integrations.

#### `VPN_USER`
*   **Description**: Username to be used for VPN client configurations (e.g., WireGuard or OpenVPN, if integrated into M3TAL).
*   **Default Value**: `user`
*   **Example Value**: `john_doe_vpn`
*   **Component(s) Used**: (Placeholder for future VPN integration components)

#### `VPN_PASSWORD`
*   **Description**: Password to be used for VPN client configurations (e.g., WireGuard or OpenVPN, if integrated into M3TAL).
*   **Default Value**: `password`
*   **Example Value**: `a_secure_vpn_password!`
*   **Component(s) Used**: (Placeholder for future VPN integration components)

### System Utilities

General system-level variables for container behavior and platform settings.

#### `DEBUG_MODE`
*   **Description**: A boolean flag to enable or disable debug-level logging and potentially other debug features across M3TAL components. Setting to `true` often provides more verbose output useful for development or advanced troubleshooting.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Component(s) Used**: `m3tal-api.service`, `m3tal-dashboard` container

#### `METRICS_ENABLED`
*   **Description**: A boolean flag that controls whether M3TAL components, particularly the API daemon, expose Prometheus-compatible metrics endpoints for monitoring.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Component(s) Used**: `m3tal-api.service`

#### `PGID`
*   **Description**: The Group ID (GID) that containers should run as when accessing host-mounted volumes. Setting this to match the GID of a user on your host ensures proper file permissions and avoids permission denied errors.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Component(s) Used**: `m3tal-dashboard` container, `ollama` container, any user-deployed stacks

#### `PUID`
*   **Description**: The User ID (UID) that containers should run as when accessing host-mounted volumes. Setting this to match the UID of a user on your host ensures proper file permissions and avoids permission denied errors.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Component(s) Used**: `m3tal-dashboard` container, `ollama` container, any user-deployed stacks

#### `TZ`
*   **Description**: Sets the timezone inside M3TAL Docker containers. This affects timestamps in logs, scheduled tasks, and any time-sensitive operations within the containers.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Component(s) Used**: `m3tal-dashboard` container, `ollama` container, any user-deployed stacks
*   **Notes**: Use standard IANA timezone names (e.g., `America/New_York`, `Asia/Tokyo`).