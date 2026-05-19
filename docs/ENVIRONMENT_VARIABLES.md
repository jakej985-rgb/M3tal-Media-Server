# M3TAL Environment Variables Reference

All M3TAL configuration is managed through environment variables, typically set in `/etc/m3tal/.env`. This file is read by both the M3TAL CLI and all Docker Compose stacks. Variables can be managed using `m3tal config wizard` or `m3tal config set KEY value`.

## Quick Reference Table

| Variable Name            | Default Value       | Description                                                                                                                                    |
|--------------------------|---------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| **Core**                 |                     |                                                                                                                                                |
| `LOG_LEVEL`              | `info`              | Sets the logging level for M3TAL services.                                                                                                     |
| `DEBUG_MODE`             | `false`             | Enables or disables debug mode for enhanced logging and diagnostics.                                                                             |
| `STATE_DIR`              | `./state`           | Directory for M3TAL's state database and other critical files.                                                                                 |
| `PUID`                   | `1000`              | The User ID to run containers as.                                                                                                              |
| `PGID`                   | `1000`              | The Group ID to run containers as.                                                                                                             |
| `TZ`                     | `America/Denver`    | The timezone for M3TAL services.                                                                                                               |
| `METRICS_ENABLED`        | `true`              | Enables or disables metrics collection for M3TAL services.                                                                                     |
| **Authentication**       |                     |                                                                                                                                                |
| `DASHBOARD_SECRET`       | `change_me_immediately` | Secret key for the dashboard session. Auto-generated on first `m3tal init`. Should NOT be set manually unless rotating.                      |
| `API_TOKEN`              | `change_me_api_token` | API authentication token. Auto-generated on first `m3tal init`. Should NOT be set manually unless rotating.                                    |
| `ADMIN_PASSWORD`         | `admin_pass`        | Password for the M3TAL dashboard administrator account.                                                                                        |
| **Network**              |                     |                                                                                                                                                |
| `HTTP_PORT`              | `8080`              | The port the M3TAL API daemon listens on.                                                                                                      |
| `NETWORK_NAME`           | `m3tal`             | The name of the Docker network M3TAL services will use.                                                                                        |
| `LOCAL_IP`               | `127.0.0.1`         | The local IP address to bind services to.                                                                                                      |
| **Storage**              |                     |                                                                                                                                                |
| `BASE_STORAGE_PATH`      | `./data`            | **Production deployments default to `/mnt`.** Base path for all M3TAL data storage. Controls where media data is stored.                      |
| `MEDIA_PATH`             | `./data/media`      | Path for media files within `BASE_STORAGE_PATH`.                                                                                               |
| `CONFIG_PATH`            | `./data/config`     | Path for configuration files within `BASE_STORAGE_PATH`.                                                                                       |
| `DOWNLOADS_PATH`         | `./data/downloads`  | Path for downloaded files within `BASE_STORAGE_PATH`.                                                                                          |
| **Dashboard**            |                     |                                                                                                                                                |
| `DASHBOARD_PORT`         | `8082`              | The port the M3TAL Dashboard service runs on internally.                                                                                       |
| `DASHBOARD_EXPOSE_MODE`  | `local`             | Controls how the dashboard is exposed: `local` (direct port) or `traefik` (via reverse proxy).                                               |
| **Traefik**              |                     |                                                                                                                                                |
| `DOMAIN`                 | `localhost`         | The primary domain for M3TAL services. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.                                 |
| `TRAEFIK_WEB_PORT`       | `80`                | The host port Traefik uses for HTTP traffic.                                                                                                   |
| `TRAEFIK_WEBHTTPS_PORT`  | `443`               | The host port Traefik uses for HTTPS traffic.                                                                                                  |
| `TRAEFIK_DASHBOARD_PORT` | `8080`              | The internal port Traefik's own dashboard is accessible on.                                                                                    |
| **VPN**                  |                     |                                                                                                                                                |
| `VPN_USER`               | `user`              | Username for VPN connection.                                                                                                                   |
| `VPN_PASSWORD`           | `password`          | Password for VPN connection.                                                                                                                   |

---

## Detailed Environment Variable Reference

All M3TAL configuration is managed through environment variables, typically set in `/etc/m3tal/.env`. This file is read by both the M3TAL CLI and all Docker Compose stacks via `--env-file`.

### Core

These variables control fundamental aspects of the M3TAL system's operation.

*   **`LOG_LEVEL`**
    *   **Description:** Sets the logging level for M3TAL services. This determines the verbosity of log messages.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** CLI, API daemon, Dashboard container

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode. When set to `true`, M3TAL services will provide more detailed logging and diagnostic information.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** CLI, API daemon, Dashboard container

*   **`STATE_DIR`**
    *   **Description:** Specifies the directory where M3TAL stores its state database (e.g., `/var/lib/m3tal/state.db`) and other critical files.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** CLI, API daemon

*   **`PUID`**
    *   **Description:** The User ID (UID) that Docker containers will run as. This is important for file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** Dashboard container

*   **`PGID`**
    *   **Description:** The Group ID (GID) that Docker containers will run as. This is important for file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** Dashboard container

*   **`TZ`**
    *   **Description:** The timezone to be used by M3TAL services. This ensures correct timestamping of logs and events.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** Dashboard container

*   **`METRICS_ENABLED`**
    *   **Description:** Enables or disables the collection and reporting of system metrics.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** API daemon

### Authentication

These variables are crucial for securing your M3TAL instance. `DASHBOARD_SECRET` and `API_TOKEN` are automatically generated during the initial `m3tal init` process and should generally not be set manually unless you are intentionally rotating them.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for securing dashboard sessions. **Auto-generated on first `m3tal init`. Users should NOT set this manually unless rotating.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_very_long_and_random_secret_string_generated_by_m3tal`
    *   **Used By:** Dashboard container

*   **`API_TOKEN`**
    *   **Description:** An authentication token used by clients to interact with the M3TAL API. **Auto-generated on first `m3tal init`. Users should NOT set this manually unless rotating.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `super_secret_api_key_for_m3tal_operations`
    *   **Used By:** API daemon

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator account of the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `MySecureAdminPassword123!`
    *   **Used By:** Dashboard container

### Network

These variables define network-related configurations for M3TAL services.

*   **`HTTP_PORT`**
    *   **Description:** The port on the host machine that the M3TAL API daemon will listen on.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** API daemon

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will be connected to. This is used for inter-container communication.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal_network`
    *   **Used By:** CLI, API daemon, Dashboard container, Traefik container

*   **`LOCAL_IP`**
    *   **Description:** The local IP address to which services will be bound. This is often `127.0.0.1` for local access or `0.0.0.0` to listen on all interfaces.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `0.0.0.0`
    *   **Used By:** API daemon

### Storage

These variables specify the locations for M3TAL's data and configuration files.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** **In production deployments, this defaults to `/mnt`. In template/development environments, it defaults to `./data`.** This is the root directory for all M3TAL data storage, including media, configuration, and downloads.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Used By:** CLI, API daemon, Dashboard container

*   **`MEDIA_PATH`**
    *   **Description:** The specific path within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/m3tal_storage/media`
    *   **Used By:** API daemon, Dashboard container

*   **`CONFIG_PATH`**
    *   **Description:** The specific path within `BASE_STORAGE_PATH` where M3TAL's configuration files are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/m3tal_storage/config`
    *   **Used By:** API daemon, Dashboard container

*   **`DOWNLOADS_PATH`**
    *   **Description:** The specific path within `BASE_STORAGE_PATH` where downloaded files will be stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/m3tal_storage/downloads`
    *   **Used By:** API daemon, Dashboard container

### Dashboard

These variables control the behavior and accessibility of the M3TAL dashboard.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port on which the M3TAL dashboard service runs. This is distinct from how it's exposed to the user.
    *   **Default Value:** `8082`
    *   **Example Value:** `8083`
    *   **Used By:** Dashboard container

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Determines how the M3TAL dashboard is made accessible.
        *   `local`: The dashboard is exposed directly via a port binding (controlled by `DASHBOARD_PORT`). Access is typically via `http://HOST_IP:DASHBOARD_PORT`. This mode does not require Traefik.
        *   `traefik`: The dashboard is exposed via the Traefik reverse proxy. Access is via `http://dash.DOMAIN`. Traefik must be running and configured correctly.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** CLI, Dashboard container (via compose overrides)

### Traefik

These variables are related to the Traefik reverse proxy, which is used to route external traffic to M3TAL services.

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL instance. Setting this variable enables Traefik to create routing rules for `api.DOMAIN` and `dash.DOMAIN`. If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, Traefik will route to `dash.DOMAIN`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.mydomain.com`
    *   **Used By:** CLI, Traefik container (via dynamic configuration)

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The port on the host machine that Traefik listens on for incoming HTTP traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `8080`
    *   **Used By:** Traefik container

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The port on the host machine that Traefik listens on for incoming HTTPS traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `8443`
    *   **Used By:** Traefik container

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The internal port on which Traefik's own administrative dashboard is accessible. This is typically accessed via `http://localhost:8081` or similar.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** Traefik container

### VPN

These variables are used if you have configured M3TAL to connect to a VPN.

*   **`VPN_USER`**
    *   **Description:** The username for your VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `vpn_client_user`
    *   **Used By:** VPN-related services (if configured)

*   **`VPN_PASSWORD`**
    *   **Description:** The password for your VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `VpnUserP@ssw0rd123`
    *   **Used By:** VPN-related services (if configured)