As DocSmith, the M3TAL Ecosystem Documentation Architect, here is the complete reference for M3TAL's environment variables.

---

# Environment Variables Reference

All M3TAL environment variables are read from the primary configuration file located at `/etc/m3tal/.env`. Both the `m3tal` CLI binary and all Docker Compose stacks (`m3tal up`, `m3tal dash up`, etc.) are configured to load variables from this file via the `--env-file` option, ensuring a consistent and centralized configuration.

This document provides a comprehensive overview of each variable, its purpose, default value, and the components it affects.

## Quick Reference

| Name                    | Description                                                                  | Default Value         |
| :---------------------- | :--------------------------------------------------------------------------- | :-------------------- |
| `DASHBOARD_PORT`        | Port for the M3TAL Dashboard.                                                | `8082`                |
| `DASHBOARD_EXPOSE_MODE` | Controls how the M3TAL Dashboard is exposed (local port or via Traefik).     | `local`               |
| `HTTP_PORT`             | Port for the M3TAL API daemon.                                               | `8080`                |
| `STATE_DIR`             | Directory for dashboard internal state and configuration.                    | `./state`             |
| `LOG_LEVEL`             | Verbosity of M3TAL's logging output.                                         | `info`                |
| `DASHBOARD_SECRET`      | Secret key for dashboard session management and security.                    | `change_me_immediately` |
| `API_TOKEN`             | Authentication token for the M3TAL API.                                      | `change_me_api_token` |
| `ADMIN_PASSWORD`        | Initial password for the default M3TAL Dashboard administrator.              | `admin_pass`          |
| `NETWORK_NAME`          | Name of the Docker network for M3TAL's core services.                        | `m3tal`               |
| `LOCAL_IP`              | Local IP address of the host machine.                                        | `127.0.0.1`           |
| `DOMAIN`                | Primary domain name for M3TAL services, used by Traefik.                     | `localhost`           |
| `VPN_USER`              | Username for accessing a VPN service.                                        | `user`                |
| `VPN_PASSWORD`          | Password for accessing a VPN service.                                        | `password`            |
| `BASE_STORAGE_PATH`     | Root directory for all M3TAL data storage.                                   | `./data`              |
| `MEDIA_PATH`            | Subdirectory for media files.                                                | `./data/media`        |
| `CONFIG_PATH`           | Subdirectory for configuration files.                                        | `./data/config`       |
| `DOWNLOADS_PATH`        | Subdirectory for downloaded content.                                         | `./data/downloads`    |
| `PUID`                  | User ID (UID) for Docker containers.                                         | `1000`                |
| `PGID`                  | Group ID (GID) for Docker containers.                                        | `1000`                |
| `TZ`                    | Timezone for M3TAL services and containers.                                  | `America/Denver`      |
| `TRAEFIK_WEB_PORT`      | HTTP port Traefik listens on.                                                | `80`                  |
| `TRAEFIK_WEBHTTPS_PORT` | HTTPS port Traefik listens on.                                               | `443`                 |
| `TRAEFIK_DASHBOARD_PORT`| Internal port for Traefik's administrative dashboard.                        | `8080`                |
| `DEBUG_MODE`            | Enables debug-level logging and features.                                    | `false`               |
| `METRICS_ENABLED`       | Controls whether M3TAL components expose Prometheus-compatible metrics.      | `true`                |

---

## Detailed Reference

### Authentication Variables

These variables manage access and security for M3TAL's core components.

*   **`DASHBOARD_SECRET`**
    *   **Description**: A secret key used by the M3TAL Dashboard for session management, cookie signing, and overall application security.
    *   **Default Value**: `change_me_immediately`
    *   **Example Value**: `super_secret_dashboard_key_123abcDEF`
    *   **Used By**: `m3tal-dashboard` container
    *   **Note**: This variable is **auto-generated** on the first `m3tal init` command. Users should not set this manually unless performing a key rotation.
*   **`API_TOKEN`**
    *   **Description**: An authentication token required to interact with the M3TAL API daemon. It secures communication between the dashboard, CLI, and external tools with the API.
    *   **Default Value**: `change_me_api_token`
    *   **Example Value**: `long_random_api_token_XYZ789`
    *   **Used By**: `API daemon`, `m3tal-dashboard` container, `CLI binary`
    *   **Note**: This variable is **auto-generated** on the first `m3tal init` command. Users should not set this manually unless performing a token rotation.
*   **`ADMIN_PASSWORD`**
    *   **Description**: The initial password for the default administrator user of the M3TAL Dashboard. This is typically used during initial setup or if no `users.json` file exists.
    *   **Default Value**: `admin_pass`
    *   **Example Value**: `MyStrongAdminPass!2023`
    *   **Used By**: `m3tal-dashboard` container (for initial user creation)

### Network Variables

These variables control how M3TAL components communicate and are exposed over the network.

*   **`DASHBOARD_PORT`**
    *   **Description**: Specifies the internal and, optionally, external port on which the M3TAL Dashboard container is accessible.
    *   **Default Value**: `8082`
    *   **Example Value**: `8082`
    *   **Used By**: `m3tal-dashboard` container (for direct port binding in `local` expose mode)
*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description**: Determines how the M3TAL Dashboard is made available.
        *   `local`: The dashboard is exposed directly on `http://HOST_IP:${DASHBOARD_PORT}` via a direct Docker port binding. This mode does not require Traefik.
        *   `traefik`: The dashboard is exposed via the Traefik reverse proxy at `http://dash.${DOMAIN}`. Requires Traefik to be running.
    *   **Default Value**: `local`
    *   **Example Value**: `traefik`
    *   **Used By**: `CLI binary` (`m3tal dash up` command to select the appropriate Docker Compose override file)
*   **`HTTP_PORT`**
    *   **Description**: The port on which the M3TAL API daemon listens for incoming HTTP requests. This is primarily for internal communication, with Traefik routing external requests to it.
    *   **Default Value**: `8080`
    *   **Example Value**: `8080`
    *   **Used By**: `API daemon`
*   **`NETWORK_NAME`**
    *   **Description**: The name of the Docker network used by M3TAL's core services (e.g., `m3tal-dashboard`, `traefik`) for inter-container communication. This network must be created externally or defined consistently across stacks.
    *   **Default Value**: `m3tal`
    *   **Example Value**: `m3tal_proxy`
    *   **Used By**: `Docker Compose` stacks (implicitly by services connecting to the `proxy` network)
*   **`LOCAL_IP`**
    *   **Description**: The local IP address of the host machine. While not directly used for binding in Docker containers, it can be referenced by the CLI or other services for internal routing and resolving host-specific addresses.
    *   **Default Value**: `127.0.0.1`
    *   **Example Value**: `192.168.1.10`
    *   **Used By**: `CLI binary`, `API daemon` (contextually for internal references)

### Storage Variables

These variables define the filesystem paths where M3TAL stores its data.

*   **`BASE_STORAGE_PATH`**
    *   **Description**: The root directory for all M3TAL data storage, including media, configuration, and downloads. All other storage-related paths are typically subdirectories of this.
    *   **Default Value**: `./data` (relative to the `m3tal` stack directory)
    *   **Example Value**: `/mnt/m3tal_data`
    *   **Used By**: `m3tal-dashboard` container (mounted as `/mnt`), `CLI binary`, other user-defined stacks.
    *   **Note**: In production deployments, this variable often defaults to `/mnt` to leverage dedicated storage mounts.
*   **`MEDIA_PATH`**
    *   **Description**: A subdirectory within `BASE_STORAGE_PATH` designated for storing media files, such as videos, music, or images, managed by M3TAL-enabled applications.
    *   **Default Value**: `./data/media`
    *   **Example Value**: `/mnt/m3tal_data/media`
    *   **Used By**: `m3tal-dashboard` container (for internal volume management), user-defined media stacks
*   **`CONFIG_PATH`**
    *   **Description**: A subdirectory within `BASE_STORAGE_PATH` used to store M3TAL's and other applications' configuration files, ensuring persistent settings.
    *   **Default Value**: `./data/config`
    *   **Example Value**: `/mnt/m3tal_data/config`
    *   **Used By**: `m3tal-dashboard` container (mounted for state and config), `CLI binary`
*   **`DOWNLOADS_PATH`**
    *   **Description**: A subdirectory within `BASE_STORAGE_PATH` intended for storing downloaded content from various M3TAL-managed applications.
    *   **Default Value**: `./data/downloads`
    *   **Example Value**: `/mnt/m3tal_data/downloads`
    *   **Used By**: User-defined download management stacks
*   **`STATE_DIR`**
    *   **Description**: Specifically, the directory where the M3TAL Dashboard stores its internal state, user credentials (`users.json`), and dynamic configuration. This path is often mounted from `CONFIG_PATH`.
    *   **Default Value**: `./state`
    *   **Example Value**: `/var/lib/m3tal/dashboard-state`
    *   **Used By**: `m3tal-dashboard` container (as an internal environment variable `STATE_DIR=/docker/state` which maps to a mounted volume)

### Traefik Variables

These variables configure the Traefik reverse proxy for routing external traffic to M3TAL services.

*   **`DOMAIN`**
    *   **Description**: The primary domain name under which M3TAL services (e.g., dashboard, API) will be accessible when using Traefik. Setting this enables routing rules like `dash.DOMAIN` and `api.DOMAIN`.
    *   **Default Value**: `localhost`
    *   **Example Value**: `m3tal.example.com`
    *   **Used By**: `Traefik gateway` (for dynamic routing rules), `CLI binary` (for generating Traefik configuration)
    *   **Note**: This variable is crucial for enabling domain-based routing via Traefik.
*   **`TRAEFIK_WEB_PORT`**
    *   **Description**: The HTTP port on the host machine that Traefik listens on for incoming web requests.
    *   **Default Value**: `80`
    *   **Example Value**: `80`
    *   **Used By**: `Traefik gateway`
*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description**: The HTTPS port on the host machine that Traefik listens on for incoming secure web requests.
    *   **Default Value**: `443`
    *   **Example Value**: `443`
    *   **Used By**: `Traefik gateway`
*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description**: The *internal* port of Traefik's own administrative dashboard. On the host, it's typically mapped to `127.0.0.1:8081`.
    *   **Default Value**: `8080`
    *   **Example Value**: `8080`
    *   **Used By**: `Traefik gateway` (internal service configuration)

### VPN Variables

These variables are placeholders for potential VPN configuration, often used with services like Cloudflare Tunnels if they require user authentication or for managing other VPN solutions.

*   **`VPN_USER`**
    *   **Description**: A username that may be used for authenticating with a VPN service, potentially managed or integrated by M3TAL.
    *   **Default Value**: `user`
    *   **Example Value**: `johndoe_vpn`
    *   **Used By**: Hypothetical or future VPN management components (e.g., `cloudflared` if it had a user/pass auth flow, or a WireGuard management stack)
*   **`VPN_PASSWORD`**
    *   **Description**: The corresponding password for the `VPN_USER`, used for VPN service authentication.
    *   **Default Value**: `password`
    *   **Example Value**: `SecureVPN@123`
    *   **Used By**: Hypothetical or future VPN management components

### System Variables

These variables control general system-level behavior, permissions, and logging for M3TAL components.

*   **`PUID`**
    *   **Description**: The User ID (UID) that Docker containers should run as. This is crucial for ensuring that files and directories created or accessed by containers have the correct permissions on the host system.
    *   **Default Value**: `1000`
    *   **Example Value**: `1001`
    *   **Used By**: `Docker Compose` containers (e.g., `m3tal-dashboard`)
*   **`PGID`**
    *   **Description**: The Group ID (GID) that Docker containers should run as, complementing `PUID` for proper file and directory permissions.
    *   **Default Value**: `1000`
    *   **Example Value**: `1001`
    *   **Used By**: `Docker Compose` containers (e.g., `m3tal-dashboard`)
*   **`TZ`**
    *   **Description**: Sets the timezone for all M3TAL services and containers. This ensures consistent timestamps in logs and scheduled tasks.
    *   **Default Value**: `America/Denver`
    *   **Example Value**: `Europe/London`
    *   **Used By**: `Docker Compose` containers (e.g., `m3tal-dashboard`), `API daemon`, `CLI binary`
*   **`LOG_LEVEL`**
    *   **Description**: Controls the verbosity of M3TAL's logging output. Common values include `debug`, `info`, `warn`, `error`. Higher verbosity (e.g., `debug`) provides more detailed insights for troubleshooting.
    *   **Default Value**: `info`
    *   **Example Value**: `debug`
    *   **Used By**: `CLI binary`, `API daemon`, `m3tal-dashboard` container
*   **`DEBUG_MODE`**
    *   **Description**: A boolean flag to enable or disable debug-specific features, logging, or behaviors across M3TAL components. Setting to `true` often provides more verbose output and internal diagnostic information.
    *   **Default Value**: `false`
    *   **Example Value**: `true`
    *   **Used By**: `CLI binary`, `API daemon`, `m3tal-dashboard` container
*   **`METRICS_ENABLED`**
    *   **Description**: Controls whether M3TAL components expose Prometheus-compatible metrics endpoints. Set to `false` to disable metrics collection, for example, in environments where metrics are not needed or to reduce overhead.
    *   **Default Value**: `true`
    *   **Example Value**: `false`
    *   **Used By**: `API daemon`, `m3tal-dashboard` container