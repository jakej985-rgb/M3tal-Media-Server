# M3TAL Environment Variables Reference

This document details all environment variables used by the M3TAL ecosystem. These variables are crucial for configuring your M3TAL deployment.

All M3TAL environment variables are read from the `/etc/m3tal/.env` file. This file is managed by the `m3tal config wizard` and can be updated with `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks utilize this file via the `--env-file` option.

## Quick Reference

| Variable Name           | Description                                                                                                                                                              | Default Value      |
|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| **Core**                |                                                                                                                                                                          |                    |
| `DASHBOARD_PORT`        | The port on which the M3TAL dashboard will listen.                                                                                                                       | `8082`             |
| `HTTP_PORT`             | The port on which the M3TAL API daemon will listen.                                                                                                                      | `8080`             |
| `STATE_DIR`             | The directory where M3TAL stores its state database and configuration files.                                                                                             | `./state`          |
| `LOG_LEVEL`             | The logging level for M3TAL services (e.g., `debug`, `info`, `warn`, `error`).                                                                                         | `info`             |
| `DEBUG_MODE`            | Enables or disables debug mode for M3TAL services.                                                                                                                       | `false`            |
| `METRICS_ENABLED`       | Enables or disables metrics collection for M3TAL services.                                                                                                               | `true`             |
| **Authentication**      |                                                                                                                                                                          |                    |
| `DASHBOARD_SECRET`      | A secret key used for securing the M3TAL dashboard session. **Auto-generated on first `m3tal init`. Should only be manually set for rotation.**                         | `change_me_immediately` |
| `API_TOKEN`             | A token for authenticating with the M3TAL API. **Auto-generated on first `m3tal init`. Should only be manually set for rotation.**                                     | `change_me_api_token` |
| `ADMIN_PASSWORD`        | The password for the default administrator user of the M3TAL dashboard.                                                                                                  | `admin_pass`       |
| **Network**             |                                                                                                                                                                          |                    |
| `NETWORK_NAME`          | The name of the Docker network used by M3TAL services.                                                                                                                   | `m3tal`            |
| `LOCAL_IP`              | The local IP address to bind services to.                                                                                                                                | `127.0.0.1`        |
| `DOMAIN`                | The primary domain name for accessing M3TAL services. Used by Traefik for routing rules (e.g., `dash.DOMAIN`, `api.DOMAIN`).                                             | `localhost`        |
| **Storage**             |                                                                                                                                                                          |                    |
| `BASE_STORAGE_PATH`     | The base path for storing M3TAL data, including configuration and media. **Defaults to `/mnt` in production deployments.**                                                | `./data`           |
| `MEDIA_PATH`            | The path within `BASE_STORAGE_PATH` for storing media files.                                                                                                             | `./data/media`     |
| `CONFIG_PATH`           | The path within `BASE_STORAGE_PATH` for storing M3TAL configuration files.                                                                                               | `./data/config`    |
| `DOWNLOADS_PATH`        | The path within `BASE_STORAGE_PATH` for storing downloaded files.                                                                                                        | `./data/downloads` |
| **User/Group IDs**      |                                                                                                                                                                          |                    |
| `PUID`                  | The user ID to run Docker containers as.                                                                                                                                 | `1000`             |
| `PGID`                  | The group ID to run Docker containers as.                                                                                                                                | `1000`             |
| **System**              |                                                                                                                                                                          |                    |
| `TZ`                    | The timezone for M3TAL services.                                                                                                                                         | `America/Denver`  |
| **Traefik**             |                                                                                                                                                                          |                    |
| `DASHBOARD_EXPOSE_MODE` | Controls how the M3TAL dashboard is exposed. Options are `local` (direct port binding) or `traefik` (routed via Traefik).                                               | `local`            |
| `TRAEFIK_WEB_PORT`      | The port Traefik listens on for HTTP traffic.                                                                                                                            | `80`               |
| `TRAEFIK_WEBHTTPS_PORT` | The port Traefik listens on for HTTPS traffic.                                                                                                                           | `443`              |
| `TRAEFIK_DASHBOARD_PORT`| The port on which the Traefik dashboard itself is accessible internally.                                                                                                   | `8080`             |
| **VPN**                 |                                                                                                                                                                          |                    |
| `VPN_USER`              | Username for VPN connection.                                                                                                                                             | `user`             |
| `VPN_PASSWORD`          | Password for VPN connection.                                                                                                                                             | `password`         |

---

## Detailed Environment Variable Reference

### Core

*   **Name**: `DASHBOARD_PORT`
    *   **Description**: The port on which the M3TAL dashboard will listen.
    *   **Default Value**: `8082`
    *   **Example Value**: `8082`
    *   **Used By**: `m3tal-dashboard` container, `CLI binary`

*   **Name**: `HTTP_PORT`
    *   **Description**: The port on which the M3TAL API daemon will listen.
    *   **Default Value**: `8080`
    *   **Example Value**: `8080`
    *   **Used By**: `m3tal-api.service`

*   **Name**: `STATE_DIR`
    *   **Description**: The directory where M3TAL stores its state database and configuration files.
    *   **Default Value**: `./state`
    *   **Example Value**: `/var/lib/m3tal/state`
    *   **Used By**: `m3tal-api.service`, `m3tal-dashboard` container

*   **Name**: `LOG_LEVEL`
    *   **Description**: The logging level for M3TAL services (e.g., `debug`, `info`, `warn`, `error`).
    *   **Default Value**: `info`
    *   **Example Value**: `debug`
    *   **Used By**: `m3tal-api.service`, `m3tal-dashboard` container

*   **Name**: `DEBUG_MODE`
    *   **Description**: Enables or disables debug mode for M3TAL services.
    *   **Default Value**: `false`
    *   **Example Value**: `true`
    *   **Used By**: `m3tal-api.service`, `m3tal-dashboard` container

*   **Name**: `METRICS_ENABLED`
    *   **Description**: Enables or disables metrics collection for M3TAL services.
    *   **Default Value**: `true`
    *   **Example Value**: `false`
    *   **Used By**: `m3tal-api.service`

### Authentication

*   **Name**: `DASHBOARD_SECRET`
    *   **Description**: A secret key used for securing the M3TAL dashboard session. **This variable is auto-generated on the first `m3tal init` command. Users should NOT set this manually unless they intend to rotate the secret.**
    *   **Default Value**: `change_me_immediately`
    *   **Example Value**: `a_very_secure_random_string`
    *   **Used By**: `m3tal-dashboard` container

*   **Name**: `API_TOKEN`
    *   **Description**: A token for authenticating with the M3TAL API. **This variable is auto-generated on the first `m3tal init` command. Users should NOT set this manually unless they intend to rotate the token.**
    *   **Default Value**: `change_me_api_token`
    *   **Example Value**: `another_secure_random_token`
    *   **Used By**: `CLI binary` (for internal API calls)

*   **Name**: `ADMIN_PASSWORD`
    *   **Description**: The password for the default administrator user of the M3TAL dashboard.
    *   **Default Value**: `admin_pass`
    *   **Example Value**: `your_strong_password_here`
    *   **Used By**: `m3tal-dashboard` container

### Network

*   **Name**: `NETWORK_NAME`
    *   **Description**: The name of the Docker network used by M3TAL services.
    *   **Default Value**: `m3tal`
    *   **Example Value**: `m3tal_network`
    *   **Used By**: All M3TAL Docker services

*   **Name**: `LOCAL_IP`
    *   **Description**: The local IP address to bind services to.
    *   **Default Value**: `127.0.0.1`
    *   **Example Value**: `192.168.1.100`
    *   **Used By**: `m3tal-api.service` (if binding to a specific IP)

*   **Name**: `DOMAIN`
    *   **Description**: The primary domain name for accessing M3TAL services. Setting this variable enables Traefik routing rules for `api.DOMAIN` and `dash.DOMAIN`.
    *   **Default Value**: `localhost`
    *   **Example Value**: `m3tal.example.com`
    *   **Used By**: `traefik` container (for routing rules), `m3tal-dashboard` container (for Traefik labels)

### Storage

*   **Name**: `BASE_STORAGE_PATH`
    *   **Description**: The base path for storing M3TAL data, including configuration and media. **In production deployments, this defaults to `/mnt`. In template/development environments, it defaults to `./data`.**
    *   **Default Value**: `./data`
    *   **Example Value**: `/mnt/m3tal_data`
    *   **Used By**: `m3tal-dashboard` container, `m3tal-api.service` (for volumes)

*   **Name**: `MEDIA_PATH`
    *   **Description**: The path within `BASE_STORAGE_PATH` for storing media files.
    *   **Default Value**: `./data/media`
    *   **Example Value**: `/mnt/m3tal_data/media`
    *   **Used By**: `m3tal-dashboard` container (for volumes)

*   **Name**: `CONFIG_PATH`
    *   **Description**: The path within `BASE_STORAGE_PATH` for storing M3TAL configuration files.
    *   **Default Value**: `./data/config`
    *   **Example Value**: `/mnt/m3tal_data/config`
    *   **Used By**: `m3tal-dashboard` container (for volumes)

*   **Name**: `DOWNLOADS_PATH`
    *   **Description**: The path within `BASE_STORAGE_PATH` for storing downloaded files.
    *   **Default Value**: `./data/downloads`
    *   **Example Value**: `/mnt/m3tal_data/downloads`
    *   **Used By**: `m3tal-dashboard` container (for volumes)

### User/Group IDs

*   **Name**: `PUID`
    *   **Description**: The user ID to run Docker containers as. This ensures correct file permissions for mounted volumes.
    *   **Default Value**: `1000`
    *   **Example Value**: `1000`
    *   **Used By**: `m3tal-dashboard` container

*   **Name**: `PGID`
    *   **Description**: The group ID to run Docker containers as. This ensures correct file permissions for mounted volumes.
    *   **Default Value**: `1000`
    *   **Example Value**: `1000`
    *   **Used By**: `m3tal-dashboard` container

### System

*   **Name**: `TZ`
    *   **Description**: The timezone for M3TAL services. This affects log timestamps and any date/time sensitive operations.
    *   **Default Value**: `America/Denver`
    *   **Example Value**: `UTC`
    *   **Used By**: `m3tal-dashboard` container

### Traefik

*   **Name**: `DASHBOARD_EXPOSE_MODE`
    *   **Description**: Controls how the M3TAL dashboard is exposed.
        *   `local`: The dashboard is exposed via a direct port binding on `DASHBOARD_PORT`. Access via `http://HOST_IP:DASHBOARD_PORT`.
        *   `traefik`: The dashboard is routed through Traefik. Access via `http://dash.DOMAIN`.
    *   **Default Value**: `local`
    *   **Example Value**: `traefik`
    *   **Used By**: `CLI binary` (determines which compose override to use), `traefik` container (via labels when set to `traefik`)

*   **Name**: `TRAEFIK_WEB_PORT`
    *   **Description**: The port Traefik listens on for HTTP traffic.
    *   **Default Value**: `80`
    *   **Example Value**: `80`
    *   **Used By**: `traefik` container

*   **Name**: `TRAEFIK_WEBHTTPS_PORT`
    *   **Description**: The port Traefik listens on for HTTPS traffic.
    *   **Default Value**: `443`
    *   **Example Value**: `443`
    *   **Used By**: `traefik` container

*   **Name**: `TRAEFIK_DASHBOARD_PORT`
    *   **Description**: The port on which the Traefik dashboard itself is accessible internally.
    *   **Default Value**: `8080`
    *   **Example Value**: `8080`
    *   **Used By**: `traefik` container (for Traefik dashboard service)

### VPN

*   **Name**: `VPN_USER`
    *   **Description**: Username for VPN connection.
    *   **Default Value**: `user`
    *   **Example Value**: `myvpnuser`
    *   **Used By**: Potentially used by any service requiring a VPN connection.

*   **Name**: `VPN_PASSWORD`
    *   **Description**: Password for VPN connection.
    *   **Default Value**: `password`
    *   **Example Value**: `mysecurevpnpassword`
    *   **Used By**: Potentially used by any service requiring a VPN connection.