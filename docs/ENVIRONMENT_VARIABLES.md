# Environment Variables Reference

All M3TAL environment variables are read from `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` and can also be modified using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks use this file via `--env-file` for consistent configuration.

## Quick Reference

| Variable Name           | Description                                                                         | Default Value     |
|-------------------------|-------------------------------------------------------------------------------------|-------------------|
| **Core**                |                                                                                     |                   |
| `LOG_LEVEL`             | Sets the logging verbosity for M3TAL services.                                      | `info`            |
| `DEBUG_MODE`            | Enables or disables debug mode for M3TAL services.                                  | `false`           |
| `METRICS_ENABLED`       | Enables or disables the collection of operational metrics.                          | `true`            |
| **Authentication**      |                                                                                     |                   |
| `DASHBOARD_SECRET`      | Secret key for securing the dashboard. Auto-generated on `m3tal init`.              | `change_me_immediately` |
| `API_TOKEN`             | API authentication token. Auto-generated on `m3tal init`.                         | `change_me_api_token` |
| `ADMIN_PASSWORD`        | Password for the M3TAL dashboard administrator.                                     | `admin_pass`      |
| **Network**             |                                                                                     |                   |
| `HTTP_PORT`             | The port the M3TAL API daemon listens on.                                           | `8080`            |
| `LOCAL_IP`              | The IP address M3TAL services should bind to locally.                               | `127.0.0.1`       |
| `NETWORK_NAME`          | The name of the Docker network M3TAL services will use.                             | `m3tal`           |
| `VPN_USER`              | Username for VPN connection.                                                        | `user`            |
| `VPN_PASSWORD`          | Password for VPN connection.                                                        | `password`        |
| **Storage**             |                                                                                     |                   |
| `BASE_STORAGE_PATH`     | Base directory for M3TAL data. Defaults to `/mnt` in production.                    | `./data`          |
| `MEDIA_PATH`            | Directory for storing media files.                                                  | `./data/media`    |
| `CONFIG_PATH`           | Directory for M3TAL configuration files.                                            | `./data/config`   |
| `DOWNLOADS_PATH`        | Directory for storing downloaded files.                                             | `./data/downloads`|
| `STATE_DIR`             | Directory for the M3TAL state database.                                             | `./state`         |
| `PUID`                  | User ID for container processes.                                                    | `1000`            |
| `PGID`                  | Group ID for container processes.                                                   | `1000`            |
| **Dashboard**           |                                                                                     |                   |
| `DASHBOARD_PORT`        | The port the M3TAL dashboard listens on.                                            | `8082`            |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed: `local` (direct port) or `traefik`.        | `local`           |
| **Traefik**             |                                                                                     |                   |
| `DOMAIN`                | The domain name used for Traefik routing. Enables `api.DOMAIN` and `dash.DOMAIN`.   | `localhost`       |
| `TRAEFIK_WEB_PORT`      | The host port Traefik listens on for HTTP traffic.                                  | `80`              |
| `TRAEFIK_WEBHTTPS_PORT` | The host port Traefik listens on for HTTPS traffic.                                 | `443`             |
| `TRAEFIK_DASHBOARD_PORT`| The host port Traefik's own dashboard is accessible on.                             | `8080`            |
| **System**              |                                                                                     |                   |
| `TZ`                    | Timezone for M3TAL services.                                                        | `America/Denver`  |

---

## Detailed Environment Variable Reference

All environment variables are read from `/etc/m3tal/.env`.

### Core

These variables control the fundamental behavior and logging of M3TAL services.

*   **`LOG_LEVEL`**
    *   **Description**: Sets the logging verbosity for M3TAL services. Accepted values typically include `debug`, `info`, `warn`, `error`.
    *   **Default Value**: `info`
    *   **Example Value**: `debug`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`

*   **`DEBUG_MODE`**
    *   **Description**: Enables or disables debug mode for M3TAL services. When set to `true`, services may provide more verbose output or enable additional debugging features.
    *   **Default Value**: `false`
    *   **Example Value**: `true`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`

*   **`METRICS_ENABLED`**
    *   **Description**: Enables or disables the collection of operational metrics from M3TAL services. This can be useful for monitoring performance and resource usage.
    *   **Default Value**: `true`
    *   **Example Value**: `false`
    *   **Component(s) Using**: `m3tal-api.service`

### Authentication

These variables are crucial for securing your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description**: A secret key used for securing the M3TAL dashboard's session cookies and other security-sensitive operations.
        **Note**: This variable is auto-generated on the first run of `m3tal init`. Users should **not** set this manually unless they intend to rotate the secret.
    *   **Default Value**: `change_me_immediately`
    *   **Example Value**: `a_very_long_and_random_secret_string`
    *   **Component(s) Using**: `m3tal-dashboard`

*   **`API_TOKEN`**
    *   **Description**: An authentication token used by clients and services to authenticate with the M3TAL API.
        **Note**: This variable is auto-generated on the first run of `m3tal init`. Users should **not** set this manually unless they intend to rotate the token.
    *   **Default Value**: `change_me_api_token`
    *   **Example Value**: `a_unique_api_authentication_token_here`
    *   **Component(s) Using**: `m3tal-api.service` (for generating tokens), `m3tal-dashboard` (for authenticating to API)

*   **`ADMIN_PASSWORD`**
    *   **Description**: The password for the initial administrator user of the M3TAL dashboard.
    *   **Default Value**: `admin_pass`
    *   **Example Value**: `MySecureAdminPassword123`
    *   **Component(s) Using**: `m3tal-dashboard`

### Network

These variables configure network-related aspects of M3TAL.

*   **`HTTP_PORT`**
    *   **Description**: The port on the host machine that the M3TAL API daemon will listen on. This is typically exposed internally within the Docker network.
    *   **Default Value**: `8080`
    *   **Example Value**: `8080`
    *   **Component(s) Using**: `m3tal-api.service`

*   **`LOCAL_IP`**
    *   **Description**: The IP address that M3TAL services should bind to for local communication. This is often `127.0.0.1` for internal services.
    *   **Default Value**: `127.0.0.1`
    *   **Example Value**: `127.0.0.1`
    *   **Component(s) Using**: `m3tal-api.service` (used in `extra_hosts` to resolve `host.docker.internal`)

*   **`NETWORK_NAME`**
    *   **Description**: Specifies the name of the Docker network that M3TAL services will connect to. This ensures proper inter-service communication.
    *   **Default Value**: `m3tal`
    *   **Example Value**: `m3tal-network`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`, `traefik`

*   **`VPN_USER`**
    *   **Description**: The username to be used for any configured VPN connections by M3TAL services.
    *   **Default Value**: `user`
    *   **Example Value**: `my_vpn_user`
    *   **Component(s) Using**: (Potentially custom stacks or future M3TAL features requiring VPN)

*   **`VPN_PASSWORD`**
    *   **Description**: The password to be used for any configured VPN connections by M3TAL services.
    *   **Default Value**: `password`
    *   **Example Value**: `my_vpn_password`
    *   **Component(s) Using**: (Potentially custom stacks or future M3TAL features requiring VPN)

### Storage

These variables define the locations for M3TAL's persistent data.

*   **`BASE_STORAGE_PATH`**
    *   **Description**: The base directory where M3TAL will store all its persistent data, including configuration, media, and downloads.
        **Note**: In production deployments, this defaults to `/mnt`. In template or development environments, it may default to `./data`.
    *   **Default Value**: `./data`
    *   **Example Value**: `/mnt`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`

*   **`MEDIA_PATH`**
    *   **Description**: The subdirectory within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value**: `./data/media`
    *   **Example Value**: `/mnt/media`
    *   **Component(s) Using**: `m3tal-dashboard`

*   **`CONFIG_PATH`**
    *   **Description**: The subdirectory within `BASE_STORAGE_PATH` where M3TAL configuration files are stored.
    *   **Default Value**: `./data/config`
    *   **Example Value**: `/mnt/config`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`

*   **`DOWNLOADS_PATH`**
    *   **Description**: The subdirectory within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value**: `./data/downloads`
    *   **Example Value**: `/mnt/downloads`
    *   **Component(s) Using**: `m3tal-dashboard`

*   **`STATE_DIR`**
    *   **Description**: The directory where the M3TAL state database (e.g., SQLite) is located.
    *   **Default Value**: `./state`
    *   **Example Value**: `/var/lib/m3tal/state`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`

*   **`PUID`**
    *   **Description**: The User ID (UID) that container processes should run as. This is important for file permissions when using volume mounts.
    *   **Default Value**: `1000`
    *   **Example Value**: `1001`
    *   **Component(s) Using**: `m3tal-dashboard`

*   **`PGID`**
    *   **Description**: The Group ID (GID) that container processes should run as. This is important for file permissions when using volume mounts.
    *   **Default Value**: `1000`
    *   **Example Value**: `1001`
    *   **Component(s) Using**: `m3tal-dashboard`

### Dashboard

These variables control the M3TAL dashboard's networking and access mode.

*   **`DASHBOARD_PORT`**
    *   **Description**: The port on which the M3TAL dashboard service listens internally.
    *   **Default Value**: `8082`
    *   **Example Value**: `8082`
    *   **Component(s) Using**: `m3tal-dashboard`

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description**: Determines how the M3TAL dashboard is made accessible.
        *   `local`: Exposes the dashboard via a direct port binding on the host (default). Access via `http://HOST_IP:DASHBOARD_PORT`.
        *   `traefik`: Configures Traefik to route traffic to the dashboard using a domain name. Access via `http://dash.DOMAIN`.
    *   **Default Value**: `local`
    *   **Example Value**: `traefik`
    *   **Component(s) Using**: `m3tal-dashboard` (influences compose override), `traefik` (when `traefik` mode is enabled)

### Traefik

These variables are used by Traefik for routing external traffic to M3TAL services.

*   **`DOMAIN`**
    *   **Description**: The primary domain name for your M3TAL instance. Setting this variable enables Traefik routing rules for `api.${DOMAIN}` and `dash.${DOMAIN}`.
        **Note**: If this is set, ensure Traefik is running and configured to listen on the appropriate ports.
    *   **Default Value**: `localhost`
    *   **Example Value**: `m3tal.example.com`
    *   **Component(s) Using**: `traefik` (for dynamic routing rules)

*   **`TRAEFIK_WEB_PORT`**
    *   **Description**: The host port that Traefik will listen on for incoming HTTP traffic.
    *   **Default Value**: `80`
    *   **Example Value**: `80`
    *   **Component(s) Using**: `traefik`

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description**: The host port that Traefik will listen on for incoming HTTPS traffic.
    *   **Default Value**: `443`
    *   **Example Value**: `443`
    *   **Component(s) Using**: `traefik`

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description**: The host port on which Traefik's own administrative dashboard will be accessible (typically for debugging Traefik itself).
    *   **Default Value**: `8080`
    *   **Example Value**: `8081`
    *   **Component(s) Using**: `traefik`

### System

These variables configure general system-level settings for M3TAL services.

*   **`TZ`**
    *   **Description**: Sets the timezone for M3TAL services. This ensures consistent logging timestamps and date/time operations.
    *   **Default Value**: `America/Denver`
    *   **Example Value**: `UTC`
    *   **Component(s) Using**: `m3tal-api.service`, `m3tal-dashboard`