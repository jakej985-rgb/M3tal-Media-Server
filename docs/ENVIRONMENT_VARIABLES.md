As DocSmith, the M3TAL Ecosystem Documentation Architect, here is your complete reference for M3TAL environment variables.

---

# Environment Variables Reference

All M3TAL components, including the CLI binary (`/usr/bin/m3tal`) and all Docker Compose stacks, read their configuration from the `/etc/m3tal/.env` file. This file is the primary source of truth for M3TAL's runtime environment variables and is typically managed using the `m3tal config wizard` or `m3tal config set` commands.

## Quick Reference

| Name                        | Default                 | Group    |
|-----------------------------|-------------------------|----------|
| `DASHBOARD_PORT`            | `8082`                  | Core     |
| `DASHBOARD_EXPOSE_MODE`     | `local`                 | Core     |
| `HTTP_PORT`                 | `8080`                  | Core     |
| `STATE_DIR`                 | `./state`               | Core     |
| `LOG_LEVEL`                 | `info`                  | Core     |
| `DASHBOARD_SECRET`          | `change_me_immediately` | Auth     |
| `API_TOKEN`                 | `change_me_api_token`   | Auth     |
| `ADMIN_PASSWORD`            | `admin_pass`            | Auth     |
| `NETWORK_NAME`              | `m3tal`                 | Network  |
| `LOCAL_IP`                  | `127.0.0.1`             | Network  |
| `DOMAIN`                    | `localhost`             | Network  |
| `BASE_STORAGE_PATH`         | `./data`                | Storage  |
| `MEDIA_PATH`                | `./data/media`          | Storage  |
| `CONFIG_PATH`               | `./data/config`         | Storage  |
| `DOWNLOADS_PATH`            | `./data/downloads`      | Storage  |
| `TRAEFIK_WEB_PORT`          | `80`                    | Traefik  |
| `TRAEFIK_WEBHTTPS_PORT`     | `443`                   | Traefik  |
| `TRAEFIK_DASHBOARD_PORT`    | `8080`                  | Traefik  |
| `VPN_USER`                  | `user`                  | VPN      |
| `VPN_PASSWORD`              | `password`              | VPN      |
| `PUID`                      | `1000`                  | System   |
| `PGID`                      | `1000`                  | System   |
| `TZ`                        | `America/Denver`        | System   |
| `DEBUG_MODE`                | `false`                 | System   |
| `METRICS_ENABLED`           | `true`                  | System   |

## Detailed Reference

### Core Variables

These variables control fundamental aspects of the M3TAL ecosystem's operation.

#### `DASHBOARD_PORT`

*   **Description**: The host port on which the `m3tal-dashboard` container is exposed when `DASHBOARD_EXPOSE_MODE` is set to `local`.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Used By**: `m3tal-dashboard` container, CLI (`m3tal dash up` command).

#### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Determines how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard is directly exposed on the host via `DASHBOARD_PORT`. Access via `http://HOST_IP:8082`. No Traefik required. Best for LAN-only setups or first-time users.
    *   `traefik`: The dashboard is routed via Traefik, making it available at `http://dash.${DOMAIN}`. Traefik must be running via `m3tal up`. Best for domain-based setups.
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Used By**: CLI (`m3tal dash up`), `m3tal-dashboard` container (via compose override files: `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`).

#### `HTTP_PORT`

*   **Description**: The port on which the M3TAL API daemon (`m3tal-api.service`) listens for incoming requests locally on the host machine.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: `m3tal-api.service`, Traefik gateway (for `api.${DOMAIN}` routing).

#### `STATE_DIR`

*   **Description**: Specifies the internal path within the `m3tal-dashboard` container where its state files (e.g., `users.json`) are expected. On the host, this typically corresponds to a volume mounted from `${CONFIG_PATH}/m3tal/state`. Note that the `m3tal-api.service` manages its `state.db` separately in `/var/lib/m3tal/state.db`, which is not controlled by this variable.
*   **Default Value**: `./state`
*   **Example Value**: `./state` (internal to container)
*   **Used By**: `m3tal-dashboard` container (internal variable definition).

#### `LOG_LEVEL`

*   **Description**: Sets the minimum logging severity level for M3TAL components. Available levels typically include `debug`, `info`, `warn`, `error`, `fatal`.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Used By**: `m3tal-api.service`, CLI binary.

### Authentication Variables

These variables are crucial for securing access to M3TAL services.

#### `DASHBOARD_SECRET`

*   **Description**: A unique secret key used by the `m3tal-dashboard` for secure session management and cryptographic operations.
    **This variable is auto-generated on first `m3tal init`. Users should NOT set it manually unless rotating existing keys for security reasons.**
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_very_long_random_string_of_characters_for_security`
*   **Used By**: `m3tal-dashboard` container, CLI (`m3tal config wizard`).

#### `API_TOKEN`

*   **Description**: A bearer token used to authenticate requests to the M3TAL API daemon. It provides a secure way for the Dashboard and CLI to interact with the API.
    **This variable is auto-generated on first `m3tal init`. Users should NOT set it manually unless rotating existing tokens for security reasons.**
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `another_long_random_string_for_api_authentication`
*   **Used By**: `m3tal-api.service`, CLI binary, `m3tal-dashboard` container.

#### `ADMIN_PASSWORD`

*   **Description**: The default password for the initial administrator user account in the M3TAL Dashboard. Dashboard user credentials are managed in `/docker/users.json` via the `m3tal dashpass` command.
*   **Default Value**: `admin_pass`
*   **Example Value**: `MySecurePa$$w0rd`
*   **Used By**: `m3tal-dashboard` container (initial setup via `users.json`), CLI (`m3tal dashpass`).

### Network Variables

These variables configure how M3TAL services communicate over the network, both internally and externally.

#### `NETWORK_NAME`

*   **Description**: The name of the Docker bridge network used for inter-container communication within the M3TAL ecosystem. All M3TAL-managed containers will connect to this network.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal`
*   **Used By**: All Docker containers (`m3tal-dashboard`, `traefik`, `cloudflared`, user stacks).

#### `LOCAL_IP`

*   **Description**: The local IP address of the host machine. This is used by containers to resolve `host.docker.internal` and access services running directly on the host (e.g., the M3TAL API daemon on port 8080).
*   **Default Value**: `127.0.0.1` (fallback, but should be set to actual host IP if containers need to reach host services)
*   **Example Value**: `192.168.1.100`
*   **Used By**: `m3tal-api.service` (for binding), Docker containers (for `host-gateway` resolution).

#### `DOMAIN`

*   **Description**: The primary domain name for M3TAL services. **Setting this variable controls Traefik routing rules, enabling access to the Dashboard at `http://dash.${DOMAIN}` and the API at `http://api.${DOMAIN}`.** If unset or set to `localhost`, Traefik will use `localhost` for routing rules.
*   **Default Value**: `localhost`
*   **Example Value**: `myhomelab.com`
*   **Used By**: `routing-compose.yml` (Traefik, Cloudflared), `m3tal-dashboard` container (Traefik labels for routing).

### Storage Variables

These variables define the base paths for persistent data storage on the host system.

#### `BASE_STORAGE_PATH`

*   **Description**: The root directory on the host machine where all M3TAL data storage (media, config, downloads) is organized.
    **In production deployments, this variable typically defaults to `/mnt` (e.g., `/mnt/m3tal-data`) instead of `./data` as specified in development templates.**
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal-data`
*   **Used By**: `m3tal-dashboard` container (volume mounts), user stacks, CLI.

#### `MEDIA_PATH`

*   **Description**: The specific path on the host for storing media files, often a subdirectory of `BASE_STORAGE_PATH`.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/m3tal-data/media`
*   **Used By**: `m3tal-dashboard` container (volume mounts), user stacks.

#### `CONFIG_PATH`

*   **Description**: The base path on the host for M3TAL's configuration and persistent state files. This directory contains `m3tal/state`, which is mounted into containers and houses the `state.db` (for API data) and `users.json` (for Dashboard user credentials).
*   **Default Value**: `./data/config`
*   **Example Value**: `/mnt/m3tal-data/config`
*   **Used By**: `m3tal-dashboard` container (volume mounts), CLI.

#### `DOWNLOADS_PATH`

*   **Description**: The path on the host for storing downloaded content.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/m3tal-data/downloads`
*   **Used By**: User stacks (for applications managing downloads).

### Traefik Variables

Variables specifically for configuring the Traefik reverse proxy gateway.

#### `TRAEFIK_WEB_PORT`

*   **Description**: The host port Traefik listens on for incoming HTTP (`web`) traffic. This is typically port 80.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Used By**: `traefik` container.

#### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port Traefik listens on for incoming HTTPS (`websecure`) traffic. This is typically port 443. Requires proper TLS configuration for Traefik to be effective.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Used By**: `traefik` container.

#### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal container port for Traefik's own dashboard. This port is typically mapped to `127.0.0.1:8081` on the host, making the Traefik dashboard accessible locally.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: `traefik` container.

### VPN Variables

Variables for configuring VPN client connections, primarily used by components like `cloudflared` if connecting to a VPN service.

#### `VPN_USER`

*   **Description**: Username for VPN client authentication, if a VPN service is configured for components like `cloudflared`.
*   **Default Value**: `user`
*   **Example Value**: `johndoe`
*   **Used By**: `cloudflared` container (if VPN functionality is enabled).

#### `VPN_PASSWORD`

*   **Description**: Password for VPN client authentication, if a VPN service is configured for components like `cloudflared`.
*   **Default Value**: `password`
*   **Example Value**: `MyStrongVPNPass`
*   **Used By**: `cloudflared` container (if VPN functionality is enabled).

### System Variables

General system-level settings for containers and core components.

#### `PUID`

*   **Description**: Specifies the User ID (UID) that containers should run as. This is crucial for ensuring correct file permissions when containers interact with host volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Used By**: All Docker containers.

#### `PGID`

*   **Description**: Specifies the Group ID (GID) that containers should run as. This is crucial for ensuring correct file permissions when containers interact with host volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Used By**: All Docker containers.

#### `TZ`

*   **Description**: Sets the timezone for containers and CLI components, ensuring consistent timestamping and localized time display.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Used By**: All Docker containers, CLI binary.

#### `DEBUG_MODE`

*   **Description**: A boolean flag to enable or disable debug logging and potentially other debug-specific features across M3TAL components. Set to `true` for verbose output and troubleshooting.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Used By**: CLI, `m3tal-api.service`, `m3tal-dashboard` container (if integrated).

#### `METRICS_ENABLED`

*   **Description**: A boolean flag to enable or disable the collection and exposure of application metrics (e.g., Prometheus endpoints) for monitoring purposes.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Used By**: `m3tal-api.service`.