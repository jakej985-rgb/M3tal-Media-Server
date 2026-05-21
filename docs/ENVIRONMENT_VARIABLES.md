# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables are typically managed via the `m3tal config wizard` command and stored in `/etc/m3tal/.env`. Both the M3TAL CLI and all Docker Compose stacks read their configuration from this file using the `--env-file` option.

## Quick Reference

| Variable Name              | Default Value      | Description                                                                                                   |
| -------------------------- | ------------------ | ------------------------------------------------------------------------------------------------------------- |
| **Core**                   |                    |                                                                                                               |
| `LOG_LEVEL`                | `info`             | The logging level for M3TAL components.                                                                       |
| `STATE_DIR`                | `./state`          | The directory to store M3TAL's state database.                                                                |
| `CONFIG_PATH`              | `./data/config`    | The base directory for M3TAL configuration files.                                                             |
| `DEBUG_MODE`               | `false`            | Enables or disables debug mode for enhanced logging and diagnostics.                                          |
| `METRICS_ENABLED`          | `true`             | Enables or disables the collection and exposition of M3TAL metrics.                                           |
| **Auth**                   |                    |                                                                                                               |
| `DASHBOARD_SECRET`         | `change_me_immediately` | Secret key for securing the M3TAL dashboard session. Auto-generated on first `m3tal init`.                  |
| `API_TOKEN`                | `change_me_api_token`   | API authentication token for programmatic access to M3TAL. Auto-generated on first `m3tal init`.             |
| `ADMIN_PASSWORD`           | `admin_pass`       | Password for the default M3TAL administrator account.                                                         |
| `VPN_USER`                 | `user`             | Username for VPN authentication.                                                                              |
| `VPN_PASSWORD`             | `password`         | Password for VPN authentication.                                                                              |
| **Network**                |                    |                                                                                                               |
| `HTTP_PORT`                | `8080`             | The port on which the M3TAL API daemon listens.                                                               |
| `DASHBOARD_PORT`           | `8082`             | The port on which the M3TAL dashboard service listens internally.                                             |
| `NETWORK_NAME`             | `m3tal`            | The name of the Docker network used by M3TAL services.                                                        |
| `LOCAL_IP`                 | `127.0.0.1`        | The local IP address that M3TAL services should bind to.                                                      |
| `DOMAIN`                   | `localhost`        | The primary domain name for M3TAL services. Affects Traefik routing rules (e.g., `dash.${DOMAIN}`).        |
| `TRAEFIK_WEB_PORT`         | `80`               | The host port exposed by Traefik for HTTP traffic.                                                            |
| `TRAEFIK_WEBHTTPS_PORT`    | `443`              | The host port exposed by Traefik for HTTPS traffic (if configured).                                           |
| `TRAEFIK_DASHBOARD_PORT`   | `8080`             | The internal port Traefik uses to access its own dashboard.                                                   |
| **Storage**                |                    |                                                                                                               |
| `BASE_STORAGE_PATH`        | `./data`           | The base directory for M3TAL media and data storage. Defaults to `/mnt` in production deployments.      |
| `MEDIA_PATH`               | `./data/media`     | The directory within `BASE_STORAGE_PATH` where media files are stored.                                        |
| `DOWNLOADS_PATH`           | `./data/downloads` | The directory within `BASE_STORAGE_PATH` for downloaded files.                                                |
| **Traefik**                |                    |                                                                                                               |
| `DASHBOARD_EXPOSE_MODE`    | `local`            | Controls how the M3TAL dashboard is exposed: `local` (direct port) or `traefik` (via Traefik routing).    |
| **System**                 |                    |                                                                                                               |
| `PUID`                     | `1000`             | The user ID (UID) to run Docker containers as.                                                                |
| `PGID`                     | `1000`             | The group ID (GID) to run Docker containers as.                                                               |
| `TZ`                       | `America/Denver`   | The timezone to use for M3TAL services.                                                                       |

---

## Core

These variables control the fundamental operation and logging of M3TAL services.

### `LOG_LEVEL`

*   **Description**: The logging level for M3TAL components. Supported values include `debug`, `info`, `warn`, `error`.
*   **Default Value**: `info`
*   **Example Value**: `debug`
*   **Used By**: CLI, API daemon, Dashboard container

### `STATE_DIR`

*   **Description**: The directory where M3TAL's state database (`state.db`) is stored.
*   **Default Value**: `./state`
*   **Example Value**: `/var/lib/m3tal/state`
*   **Used By**: CLI, API daemon

### `CONFIG_PATH`

*   **Description**: The base directory for M3TAL configuration files, including user-specific configurations and compose files.
*   **Default Value**: `./data/config`
*   **Example Value**: `/mnt/config/m3tal`
*   **Used By**: CLI, API daemon, Dashboard container

### `DEBUG_MODE`

*   **Description**: Enables or disables debug mode for enhanced logging and diagnostics across M3TAL components.
*   **Default Value**: `false`
*   **Example Value**: `true`
*   **Used By**: CLI, API daemon, Dashboard container

### `METRICS_ENABLED`

*   **Description**: Enables or disables the collection and exposition of M3TAL metrics, typically for Prometheus scraping.
*   **Default Value**: `true`
*   **Example Value**: `false`
*   **Used By**: API daemon

## Auth

These variables are critical for securing access to M3TAL and its services.

### `DASHBOARD_SECRET`

*   **Description**: A secret key used to secure M3TAL dashboard sessions. This is auto-generated on the first `m3tal init` command. **Users should not set this manually unless rotating it.**
*   **Default Value**: `change_me_immediately`
*   **Example Value**: `a_very_long_and_random_secret_string`
*   **Used By**: Dashboard container

### `API_TOKEN`

*   **Description**: An authentication token used for programmatic access to the M3TAL API. This is auto-generated on the first `m3tal init` command. **Users should not set this manually unless rotating it.**
*   **Default Value**: `change_me_api_token`
*   **Example Value**: `another_long_and_random_api_token`
*   **Used By**: CLI

### `ADMIN_PASSWORD`

*   **Description**: The password for the default M3TAL administrator account.
*   **Default Value**: `admin_pass`
*   **Example Value**: `a_secure_password_here`
*   **Used By**: CLI, Dashboard container

### `VPN_USER`

*   **Description**: The username used for VPN authentication, if a VPN is configured.
*   **Default Value**: `user`
*   **Example Value**: `myvpnuser`
*   **Used By**: CLI, VPN client (if applicable)

### `VPN_PASSWORD`

*   **Description**: The password used for VPN authentication, if a VPN is configured.
*   **Default Value**: `password`
*   **Example Value**: `mysecurevpnpassword`
*   **Used By**: CLI, VPN client (if applicable)

## Network

These variables define network-related configurations for M3TAL services.

### `HTTP_PORT`

*   **Description**: The port on which the M3TAL API daemon listens for incoming HTTP requests.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: API daemon

### `DASHBOARD_PORT`

*   **Description**: The internal port on which the M3TAL dashboard service listens. This is used for direct access in `local` mode or by Traefik in `traefik` mode.
*   **Default Value**: `8082`
*   **Example Value**: `8082`
*   **Used By**: Dashboard container

### `NETWORK_NAME`

*   **Description**: The name of the Docker network that M3TAL services will join. This ensures they can communicate with each other.
*   **Default Value**: `m3tal`
*   **Example Value**: `m3tal_network`
*   **Used By**: CLI, API daemon, Docker Compose stacks

### `LOCAL_IP`

*   **Description**: The local IP address that M3TAL services should bind to. Typically `127.0.0.1` for localhost access.
*   **Default Value**: `127.0.0.1`
*   **Example Value**: `192.168.1.100`
*   **Used By**: API daemon

### `DOMAIN`

*   **Description**: The primary domain name for your M3TAL installation. This variable is crucial for Traefik routing rules, enabling access to services like `dash.${DOMAIN}` and `api.${DOMAIN}`.
*   **Default Value**: `localhost`
*   **Example Value**: `m3tal.mydomain.com`
*   **Used By**: CLI, API daemon, Traefik configuration

### `TRAEFIK_WEB_PORT`

*   **Description**: The host port that Traefik will expose for incoming HTTP traffic. This is the primary entry point for services routed by Traefik.
*   **Default Value**: `80`
*   **Example Value**: `80`
*   **Used By**: Traefik container

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description**: The host port that Traefik will expose for incoming HTTPS traffic. This is used if you have SSL/TLS termination configured with Traefik.
*   **Default Value**: `443`
*   **Example Value**: `443`
*   **Used By**: Traefik container

### `TRAEFIK_DASHBOARD_PORT`

*   **Description**: The internal port on which Traefik's own dashboard is accessible. This is typically only exposed via `127.0.0.1` for local access to the Traefik UI.
*   **Default Value**: `8080`
*   **Example Value**: `8080`
*   **Used By**: Traefik container

## Storage

These variables define the locations for persistent data and media files used by M3TAL.

### `BASE_STORAGE_PATH`

*   **Description**: The base directory where M3TAL stores its persistent data, including media files, configuration backups, and other essential data. **In production deployments, this defaults to `/mnt` to utilize dedicated storage.**
*   **Default Value**: `./data`
*   **Example Value**: `/mnt`
*   **Used By**: CLI, API daemon, Dashboard container

### `MEDIA_PATH`

*   **Description**: The sub-directory within `BASE_STORAGE_PATH` where media files are stored.
*   **Default Value**: `./data/media`
*   **Example Value**: `/mnt/media`
*   **Used By**: CLI, API daemon, Dashboard container

### `DOWNLOADS_PATH`

*   **Description**: The sub-directory within `BASE_STORAGE_PATH` where downloaded files are temporarily or permanently stored.
*   **Default Value**: `./data/downloads`
*   **Example Value**: `/mnt/downloads`
*   **Used By**: CLI, API daemon, Dashboard container

## Traefik

These variables specifically control how the M3TAL dashboard is exposed and accessed.

### `DASHBOARD_EXPOSE_MODE`

*   **Description**: Controls how the M3TAL dashboard is exposed to the network.
    *   `local`: The dashboard is exposed via a direct port binding (`DASHBOARD_PORT`). Access is typically `http://HOST_IP:8082` or `http://localhost:8082`. This mode does not require Traefik.
    *   `traefik`: The dashboard is exposed through Traefik. Access is via `http://dash.${DOMAIN}` (requires Traefik to be running and `DOMAIN` to be set correctly).
*   **Default Value**: `local`
*   **Example Value**: `traefik`
*   **Used By**: CLI, Dashboard container (via compose overrides)

## System

These variables influence the operating environment and user context for Docker containers.

### `PUID`

*   **Description**: The User ID (UID) that Docker containers will run as. This is important for file permissions and ownership on the host system.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Used By**: Docker Compose stacks (e.g., `m3tal-dashboard`)

### `PGID`

*   **Description**: The Group ID (GID) that Docker containers will run as. This is important for file permissions and ownership on the host system.
*   **Default Value**: `1000`
*   **Example Value**: `1000`
*   **Used By**: Docker Compose stacks (e.g., `m3tal-dashboard`)

### `TZ`

*   **Description**: The timezone to be used by M3TAL services. This ensures accurate timestamping in logs and for time-sensitive operations.
*   **Default Value**: `America/Denver`
*   **Example Value**: `UTC`
*   **Used By**: API daemon, Dashboard container