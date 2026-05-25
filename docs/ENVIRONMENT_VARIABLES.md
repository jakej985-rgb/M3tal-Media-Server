Okay, let's get these critical M3TAL environment variables documented. As DocSmith, I understand the importance of clear, precise information for system administrators and developers.

---

# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables control various aspects of the CLI, API daemon, Dashboard, and Docker Compose stacks.

**All environment variables are read from `/etc/m3tal/.env`**. Both the `m3tal` CLI binary and all Docker Compose stacks (via `--env-file /etc/m3tal/.env`) rely on this file for their configuration. The `m3tal config wizard` and `m3tal config set` commands are the recommended way to manage this file.

## Quick Reference Table

| Variable Name            | Default Value      | Description                                                                  |
| :----------------------- | :----------------- | :--------------------------------------------------------------------------- |
| `ADMIN_PASSWORD`         | `admin_pass`       | Default administrator password for the M3TAL Dashboard.                      |
| `API_TOKEN`              | `change_me_api_token` | Authentication token for the M3TAL CLI and API daemon.                       |
| `BASE_STORAGE_PATH`      | `./data`           | Base directory for M3TAL persistent data. Defaults to `/mnt` in production.  |
| `CONFIG_PATH`            | `./data/config`    | Base directory for M3TAL configuration files and internal state.             |
| `DASHBOARD_EXPOSE_MODE`  | `local`            | Controls how the M3TAL Dashboard is exposed: `local` or `traefik`.         |
| `DASHBOARD_PORT`         | `8082`             | Internal port used by the M3TAL Dashboard container.                         |
| `DASHBOARD_SECRET`       | `change_me_immediately` | Secret key for M3TAL Dashboard session management.                           |
| `DEBUG_MODE`             | `false`            | Enables debug mode for M3TAL components.                                     |
| `DOMAIN`                 | `localhost`        | Primary domain for M3TAL services, used by Traefik for routing.              |
| `DOWNLOADS_PATH`         | `./data/downloads` | Base directory for download services.                                        |
| `HTTP_PORT`              | `8080`             | Port used by the M3TAL API daemon.                                           |
| `LOCAL_IP`               | `127.0.0.1`        | Local IP address for internal network configurations.                        |
| `LOG_LEVEL`              | `info`             | Verbosity of M3TAL component logs.                                           |
| `MEDIA_PATH`             | `./data/media`     | Base directory for media storage.                                            |
| `METRICS_ENABLED`        | `true`             | Controls whether M3TAL components expose metrics.                            |
| `NETWORK_NAME`           | `m3tal`            | Name of the Docker network used by M3TAL components.                         |
| `PGID`                   | `1000`             | Group ID for M3TAL container users, for file permissions.                    |
| `PUID`                   | `1000`             | User ID for M3TAL container users, for file permissions.                     |
| `STATE_DIR`              | `./state`          | Directory for general component state (if not explicitly mounted).           |
| `TRAEFIK_DASHBOARD_PORT` | `8080`             | Internal port for Traefik's administration dashboard.                        |
| `TRAEFIK_WEB_PORT`       | `80`               | Host port for Traefik's HTTP entry point.                                    |
| `TRAEFIK_WEBHTTPS_PORT`  | `443`              | Host port for Traefik's HTTPS entry point.                                   |
| `TZ`                     | `America/Denver`   | Timezone for M3TAL containers and services.                                  |
| `VPN_PASSWORD`           | `password`         | Password for VPN service.                                                    |
| `VPN_USER`               | `user`             | Username for VPN service.                                                    |

---

## Environment Variable Reference

### Core Configuration

These variables configure fundamental aspects of the M3TAL ecosystem's operation.

#### `DASHBOARD_PORT`
*   **Description**: The internal port on which the `m3tal-dashboard` container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this port is directly exposed on the host.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Used By**: `m3tal-dashboard` container (via Docker Compose `m3tal-compose.yml`, `m3tal-compose.local.yml`)

#### `DASHBOARD_EXPOSE_MODE`
*   **Description**: Controls how the M3TAL Dashboard is exposed to the network.
    *   `local` (default): The dashboard is accessible via a direct port binding, typically `http://HOST_IP:8082`. No Traefik required. Best for LAN-only or initial setups.
    *   `traefik`: The dashboard is configured with Traefik labels, making it accessible via `http://dash.DOMAIN` (requires Traefik to be running).
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Used By**: `m3tal` CLI (specifically `m3tal dash up`), Docker Compose overrides (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`)

#### `HTTP_PORT`
*   **Description**: The port on which the M3TAL API daemon (`m3tal-api.service`) listens. This API is typically accessed locally by the Dashboard or other services, or via Traefik.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: M3TAL API daemon, `m3tal-dashboard` container (via `GO_API_URL`), Traefik gateway (for routing `api.DOMAIN`)

#### `STATE_DIR`
*   **Description**: A general directory for component-specific state files. While the primary API daemon uses `/var/lib/m3tal/state.db`, and the dashboard's internal state is mapped from `${CONFIG_PATH}/m3tal/state`, this variable could define the location for other components.
*   **Default Value**: `./state`
*   **Example Value**: `/var/lib/m3tal/other_state`
*   **Used By**: Potentially future M3TAL components or helper scripts. (Note: Current core components use explicitly defined paths or volume mounts).

### Authentication & Security

These variables manage user authentication and cryptographic secrets within M3TAL.

#### `DASHBOARD_SECRET`
*   **Description**: A unique, strong secret key used by the M3TAL Dashboard for secure session management, cryptographic operations, and protecting sensitive data.
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `aVeryLongAndComplexRandomStringForDashboardSecurity`
*   **Used By**: `m3tal-dashboard` container
*   **Important**: This variable is **auto-generated on first `m3tal init`**. Users should **NOT set it manually** unless performing a secret rotation and fully understanding the implications.

#### `API_TOKEN`
*   **Description**: An authentication token required by the M3TAL CLI and other clients to interact with the M3TAL API daemon. It secures access to core system functionalities.
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `aDifferentLongAndComplexRandomTokenForM3talApiAccess`
*   **Used By**: M3TAL CLI, M3TAL API daemon
*   **Important**: This variable is **auto-generated on first `m3tal init`**. Users should **NOT set it manually** unless performing a token rotation and fully understanding the implications.

#### `ADMIN_PASSWORD`
*   **Description**: The initial or default administrator password for accessing the M3TAL Dashboard. It is highly recommended to change this immediately after initial setup using `m3tal dashpass`.
*   **Default Value**: `admin_pass`
*   **Example Value**: `MySecureAdminPass123!`
*   **Used By**: `m3tal dashpass` command, M3TAL Dashboard (initial user setup)

### Network Configuration

Variables related to M3TAL's internal and external network setup.

#### `NETWORK_NAME`
*   **Description**: The name of the Docker network that M3TAL uses for inter-container communication. All M3TAL-managed containers are connected to this network.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal_prod_network`
*   **Used By**: Docker Compose stacks (`m3tal-compose.yml`, `routing-compose.yml`), Traefik Docker provider.

#### `LOCAL_IP`
*   **Description**: Specifies a local IP address that M3TAL components might use for internal network bindings or host-bound services.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.10`
*   **Used By**: Potentially internal M3TAL services for binding to a specific host interface. (Note: Traefik and Dashboard often use `host.docker.internal` for API access).

### Storage Paths

These variables define the base directories for M3TAL's persistent data storage.

#### `BASE_STORAGE_PATH`
*   **Description**: The root directory on the host filesystem where M3TAL stores all its media, configuration, and other persistent data.
*   **Default Value**: `./data`
*   **Example Value**: `/mnt/m3tal_storage`
*   **Used By**: `m3tal-dashboard` container (volume mounts), `m3tal` CLI for initializing stack directories.
*   **Important**: In production deployments, this variable typically defaults to `/mnt` (e.g., `/mnt/m3tal_data`) to align with common server storage conventions, rather than the template's `./data`.

#### `MEDIA_PATH`
*   **Description**: The subdirectory within `BASE_STORAGE_PATH` designated for user-managed media files (e.g., photos, videos, music).
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/m3tal_storage/media`
*   **Used By**: Hypothetical media management services or containers requiring access to user media.

#### `CONFIG_PATH`
*   **Description**: The subdirectory within `BASE_STORAGE_PATH` used for storing M3TAL-specific configuration files and internal state that needs to persist across container restarts. This includes dashboard user credentials (`users.json`).
*   **Default Value**: `./data/config`
*   **Example Value**: `/mnt/m3tal_storage/config`
*   **Used By**: `m3tal-dashboard` container (for mounting `/docker/state/config`), `m3tal` CLI.

#### `DOWNLOADS_PATH`
*   **Description**: The subdirectory within `BASE_STORAGE_PATH` where download-related services would store downloaded content.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/m3tal_storage/downloads`
*   **Used By**: Hypothetical download management services or containers.

### Traefik Gateway

Variables controlling the Traefik reverse proxy and its routing rules.

#### `DOMAIN`
*   **Description**: The primary domain name for your M3TAL services. Setting this variable enables Traefik to automatically route traffic for `dash.DOMAIN` to the M3TAL Dashboard and `api.DOMAIN` to the M3TAL API daemon.
*   **Default Value**: `localhost`
*   **Example Value**: `myhomelab.com`
*   **Used By**: Traefik gateway (`routing-compose.yml`), `m3tal-dashboard` (via `m3tal-compose.traefik.yml` labels), Traefik dynamic configuration files (`dynamic/api.yml`).

#### `TRAEFIK_WEB_PORT`
*   **Description**: The host port on which the Traefik gateway listens for incoming unencrypted HTTP traffic.
*   **Default Value**: `80`
*   **Example Value**: `8080` (if port 80 is in use by another service on the host)
*   **Used By**: `traefik` container (`routing-compose.yml`)

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description**: The host port on which the Traefik gateway listens for incoming encrypted HTTPS traffic.
*   **Default Value**: `443`
*   **Example Value**: `8443`
*   **Used By**: `traefik` container (`routing-compose.yml`)

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description**: The internal port used by Traefik for its own administration dashboard. By default, this internal port is mapped to `127.0.0.1:8081` on the host, making it accessible only from the local machine.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: `traefik` container (implicitly, as Traefik's internal dashboard runs on this port; the host mapping `127.0.0.1:8081:8080` is explicit in `routing-compose.yml`).

### VPN Services

Variables for configuring an optional VPN service.

#### `VPN_USER`
*   **Description**: The username for authenticating with a VPN service integrated into the M3TAL ecosystem.
*   **Default Value**: `user`
*   **Example Value**: `vpnclient`
*   **Used By**: Hypothetical VPN client services (e.g., OpenVPN, WireGuard).

#### `VPN_PASSWORD`
*   **Description**: The password for authenticating with a VPN service integrated into the M3TAL ecosystem.
*   **Default Value**: `password`
*   **Example Value**: `MyStrongVpnPass!@#`
*   **Used By**: Hypothetical VPN client services.

### System Configuration

General system-wide variables affecting various M3TAL components.

#### `PUID`
*   **Description**: The User ID (UID) that M3TAL Docker containers should use when running. This is crucial for ensuring correct file permissions on mounted volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Used By**: Docker Compose stacks (e.g., `m3tal-dashboard` container `user: "${PUID:-1000}:${PGID:-1000}"`).

#### `PGID`
*   **Description**: The Group ID (GID) that M3TAL Docker containers should use when running. Similar to `PUID`, this ensures correct file permissions on mounted volumes.
*   **Default Value**: `1000`
*   **Example Value**: `1001`
*   **Used By**: Docker Compose stacks (e.g., `m3tal-dashboard` container `user: "${PUID:-1000}:${PGID:-1000}"`).

#### `TZ`
*   **Description**: Specifies the timezone for M3TAL containers and services, ensuring consistent time reporting across the ecosystem.
*   **Default Value**: `America/Denver`
*   **Example Value**: `Europe/London`
*   **Used By**: Various Docker containers (e.g., `m3tal-dashboard` container).

#### `DEBUG_MODE`
*   **Description**: A flag to enable debug mode for M3TAL components. When `true`, it typically results in more verbose logging, additional diagnostic information, and sometimes enables development-specific features.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Used By**: M3TAL CLI, M3TAL API daemon, `m3tal-dashboard` container.

#### `METRICS_ENABLED`
*   **Description**: Controls whether M3TAL components expose operational metrics (e.g., Prometheus format) for monitoring and observability.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Used By**: M3TAL API daemon, `m3tal-dashboard` container.

#### `LOG_LEVEL`
*   **Description**: Sets the verbosity level for logs generated by M3TAL components. Common levels include `debug`, `info`, `warn`, `error`, and `fatal`.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Used By**: M3TAL CLI, M3TAL API daemon, `m3tal-dashboard` container.

---