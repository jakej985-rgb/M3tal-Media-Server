# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, primarily read from `/etc/m3tal/.env`. This file is automatically managed by the `m3tal config wizard` and can be individually updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks (via `--env-file`) source this file for configuration.

## Quick Reference Table

| Variable Name          | Description                                                                                                                                             | Default Value     | Example Value                     | Used By                                                                           |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------|-----------------------------------|-----------------------------------------------------------------------------------|
| **Core**               |                                                                                                                                                         |                   |                                   |                                                                                   |
| `STATE_DIR`            | Directory for M3TAL's state database.                                                                                                                   | `./state`         | `/var/lib/m3tal/state`            | API Daemon, CLI                                                                   |
| `LOG_LEVEL`            | Sets the logging verbosity for M3TAL services.                                                                                                          | `info`            | `debug`                           | API Daemon, Dashboard                                                             |
| `TZ`                   | Sets the timezone for M3TAL services.                                                                                                                   | `America/Denver`  | `UTC`                             | Dashboard, API Daemon                                                             |
| `DEBUG_MODE`           | Enables debug logging and potentially other debugging features.                                                                                         | `false`           | `true`                            | API Daemon, Dashboard                                                             |
| `METRICS_ENABLED`      | Enables Prometheus metrics collection for M3TAL services.                                                                                               | `true`            | `false`                           | API Daemon                                                                        |
| **Authentication**     |                                                                                                                                                         |                   |                                   |                                                                                   |
| `DASHBOARD_SECRET`     | Secret key for securing the dashboard session. **Auto-generated on first `m3tal init`. Do NOT set manually unless rotating.**                           | `change_me_immediately` | `a_very_secret_key_generated_by_m3tal` | Dashboard                                                                         |
| `API_TOKEN`            | API token for authentication with the M3TAL API. **Auto-generated on first `m3tal init`. Do NOT set manually unless rotating.**                          | `change_me_api_token` | `an_api_token_generated_by_m3tal`     | API Daemon, CLI                                                                   |
| `ADMIN_PASSWORD`       | The password for the default `admin` user of the dashboard.                                                                                             | `admin_pass`      | `my_strong_admin_password`        | Dashboard                                                                         |
| **Networking**         |                                                                                                                                                         |                   |                                   |                                                                                   |
| `HTTP_PORT`            | The port the M3TAL API daemon listens on.                                                                                                               | `8080`            | `8081`                            | API Daemon                                                                        |
| `LOCAL_IP`             | The IP address M3TAL services should bind to or be accessible on.                                                                                       | `127.0.0.1`       | `192.168.1.100`                   | API Daemon, Dashboard (host.docker.internal mapping)                              |
| `NETWORK_NAME`         | The name of the Docker network M3TAL services will use.                                                                                                 | `m3tal`           | `m3tal_net`                       | API Daemon, Dashboard, Traefik                                                    |
| `DASHBOARD_PORT`       | The internal port the dashboard container listens on.                                                                                                   | `8082`            | `8083`                            | Dashboard                                                                         |
| `DASHBOARD_EXPOSE_MODE`| Controls how the dashboard is exposed: `local` (direct port) or `traefik` (via Traefik reverse proxy).                                                 | `local`           | `traefik`                         | Dashboard, Traefik (configuration via compose overrides)                          |
| **Storage**            |                                                                                                                                                         |                   |                                   |                                                                                   |
| `BASE_STORAGE_PATH`    | **Controls where media data is stored.** Defaults to `/mnt` in production deployments.                                                                  | `./data`          | `/mnt`                            | API Daemon, Dashboard                                                             |
| `MEDIA_PATH`           | Subdirectory within `BASE_STORAGE_PATH` for media files.                                                                                                | `./data/media`    | `/mnt/media`                      | API Daemon, Dashboard                                                             |
| `CONFIG_PATH`          | Subdirectory within `BASE_STORAGE_PATH` for configuration files.                                                                                        | `./data/config`   | `/mnt/config`                     | API Daemon, Dashboard                                                             |
| `DOWNLOADS_PATH`       | Subdirectory within `BASE_STORAGE_PATH` for downloads.                                                                                                  | `./data/downloads`| `/mnt/downloads`                  | API Daemon, Dashboard                                                             |
| **User/Group IDs**     |                                                                                                                                                         |                   |                                   |                                                                                   |
| `PUID`                 | The User ID for running containers.                                                                                                                     | `1000`            | `1001`                            | Dashboard                                                                         |
| `PGID`                 | The Group ID for running containers.                                                                                                                    | `1000`            | `1001`                            | Dashboard                                                                         |
| **Traefik**            |                                                                                                                                                         |                   |                                   |                                                                                   |
| `DOMAIN`               | The domain name used for Traefik routing. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes.                                                  | `localhost`       | `mydomain.com`                    | Traefik (dynamic config), Dashboard (compose override)                            |
| `TRAEFIK_WEB_PORT`     | The host port Traefik uses for HTTP traffic.                                                                                                            | `80`              | `80`                              | Traefik                                                                           |
| `TRAEFIK_WEBHTTPS_PORT`| The host port Traefik uses for HTTPS traffic.                                                                                                           | `443`             | `443`                             | Traefik                                                                           |
| `TRAEFIK_DASHBOARD_PORT`| The host port Traefik listens on for its own dashboard. (Note: This is the port Traefik itself exposes, not the M3TAL dashboard's port).             | `8080`            | `8081`                            | Traefik                                                                           |
| **VPN**                |                                                                                                                                                         |                   |                                   |                                                                                   |
| `VPN_USER`             | Username for the VPN connection.                                                                                                                        | `user`            | `myvpnuser`                       | VPN Client Container (if configured)                                              |
| `VPN_PASSWORD`         | Password for the VPN connection.                                                                                                                        | `password`        | `myvpnpassword`                   | VPN Client Container (if configured)                                              |

---

## Variable Groups

All M3TAL configuration is managed via environment variables, primarily read from `/etc/m3tal/.env`. This file is automatically managed by the `m3tal config wizard` and can be individually updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks (via `--env-file`) source this file for configuration.

### Core

These variables control fundamental aspects of the M3TAL system's operation.

*   **`STATE_DIR`**
    *   **Description:** Specifies the directory where M3TAL's SQLite state database (`state.db`) will be stored.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** API Daemon, CLI

*   **`LOG_LEVEL`**
    *   **Description:** Sets the logging verbosity for M3TAL services. Options typically include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** API Daemon, Dashboard

*   **`TZ`**
    *   **Description:** Sets the timezone for M3TAL services. This ensures consistent timestamping in logs and for any time-sensitive operations.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** Dashboard, API Daemon

*   **`DEBUG_MODE`**
    *   **Description:** Enables debug logging and potentially other debugging features across M3TAL services.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** API Daemon, Dashboard

*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether M3TAL services expose Prometheus metrics for monitoring.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** API Daemon

### Authentication

These variables are crucial for securing your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for securing the dashboard session cookies. **This variable is auto-generated on the first `m3tal init` command and should NOT be set manually unless you intend to rotate the secret.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_very_secret_key_generated_by_m3tal`
    *   **Used By:** Dashboard

*   **`API_TOKEN`**
    *   **Description:** An API token used for authentication with the M3TAL API daemon. **This variable is auto-generated on the first `m3tal init` command and should NOT be set manually unless you intend to rotate the token.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `an_api_token_generated_by_m3tal`
    *   **Used By:** API Daemon, CLI

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default `admin` user of the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_strong_admin_password`
    *   **Used By:** Dashboard

### Networking

These variables configure network access and communication for M3TAL services.

*   **`HTTP_PORT`**
    *   **Description:** The port on the host that the M3TAL API daemon (Go binary) listens on.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** API Daemon

*   **`LOCAL_IP`**
    *   **Description:** The IP address that M3TAL services should bind to or be made accessible on. This is often `127.0.0.1` for local access or a specific host IP for broader network availability. It's also used for `host.docker.internal` mapping for the dashboard.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** API Daemon, Dashboard (host.docker.internal mapping)

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will be connected to. This network is also used by Traefik for service discovery.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal_net`
    *   **Used By:** API Daemon, Dashboard, Traefik

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port on which the M3TAL dashboard container listens.
    *   **Default Value:** `8082`
    *   **Example Value:** `8083`
    *   **Used By:** Dashboard

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL dashboard is made accessible.
        *   `local`: The dashboard is exposed directly via its `DASHBOARD_PORT` using a Docker port mapping. Access is typically `http://<HOST_IP>:<DASHBOARD_PORT>`.
        *   `traefik`: The dashboard is exposed via the Traefik reverse proxy. Access is typically `http://dash.<DOMAIN>`.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** Dashboard, Traefik (configuration via compose overrides)

### Storage

These variables define the locations for persistent data used by M3TAL.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** **This is a critical variable that controls where M3TAL stores its primary data, including media files, configuration, and downloads.** In production deployments, this defaults to `/mnt` for better management of persistent storage. In template/development environments, it might default to `./data`.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Used By:** API Daemon, Dashboard

*   **`MEDIA_PATH`**
    *   **Description:** A subdirectory within `BASE_STORAGE_PATH` designated for storing media files managed by M3TAL.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Used By:** API Daemon, Dashboard

*   **`CONFIG_PATH`**
    *   **Description:** A subdirectory within `BASE_STORAGE_PATH` designated for storing configuration files. This is often where dashboard credentials or other settings reside.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Used By:** API Daemon, Dashboard

*   **`DOWNLOADS_PATH`**
    *   **Description:** A subdirectory within `BASE_STORAGE_PATH` where downloaded files will be stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Used By:** API Daemon, Dashboard

### User/Group IDs

These variables are used to set the User ID (UID) and Group ID (GID) for processes running inside Docker containers. This helps ensure correct file permissions when containers access host volumes.

*   **`PUID`**
    *   **Description:** The User ID to be used when running Docker containers. This should typically match the owner of the persistent storage directories on the host.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** Dashboard

*   **`PGID`**
    *   **Description:** The Group ID to be used when running Docker containers. This should typically match the group owner of the persistent storage directories on the host.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** Dashboard

### Traefik

These variables configure the Traefik reverse proxy, which handles routing external traffic to M3TAL services.

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL instance. Setting this variable enables Traefik to create routing rules for services like `api.DOMAIN` and `dash.DOMAIN`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `mydomain.com`
    *   **Used By:** Traefik (dynamic config), Dashboard (compose override)

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The port on the host that Traefik listens on for incoming HTTP traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** Traefik

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The port on the host that Traefik listens on for incoming HTTPS traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** Traefik

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The port on the host that Traefik listens on for its own administrative dashboard (which is typically accessed via `http://localhost:8081` by default). This is distinct from the M3TAL dashboard's port.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** Traefik

### VPN

These variables are used if you configure M3TAL to use a VPN client.

*   **`VPN_USER`**
    *   **Description:** The username required to connect to your VPN service.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** VPN Client Container (if configured)

*   **`VPN_PASSWORD`**
    *   **Description:** The password required to connect to your VPN service.
    *   **Default Value:** `password`
    *   **Example Value:** `myvpnpassword`
    *   **Used By:** VPN Client Container (if configured)