# M3TAL Environment Variables Reference

This document details all environment variables used by the M3TAL ecosystem. These variables are primarily configured via the `m3tal config wizard` or by directly editing the `/etc/m3tal/.env` file. All M3TAL components, including the CLI and Docker Compose stacks, read their configuration from this file.

## Quick Reference

| Variable Name         | Default Value    | Description                                                       |
|-----------------------|------------------|-------------------------------------------------------------------|
| **Core**              |                  |                                                                   |
| `HTTP_PORT`           | `8080`           | Port for the M3TAL API daemon.                                    |
| `STATE_DIR`           | `./state`        | Directory for M3TAL's state database.                             |
| `LOG_LEVEL`           | `info`           | Logging level for M3TAL services.                                 |
| `DEBUG_MODE`          | `false`          | Enables debug logging and features.                               |
| `METRICS_ENABLED`     | `true`           | Enables Prometheus metrics collection.                            |
| **Auth**              |                  |                                                                   |
| `DASHBOARD_SECRET`    | `change_me_immediately` | Secret for securing the dashboard. Auto-generated on `m3tal init`. |
| `API_TOKEN`           | `change_me_api_token` | Token for authenticating API requests. Auto-generated on `m3tal init`. |
| `ADMIN_PASSWORD`      | `admin_pass`     | Password for the initial admin user.                              |
| **Network**           |                  |                                                                   |
| `NETWORK_NAME`        | `m3tal`          | Name of the Docker network used by M3TAL services.              |
| `LOCAL_IP`            | `127.0.0.1`      | Local IP address for host-specific networking.                  |
| `DOMAIN`              | `localhost`      | Primary domain for Traefik routing rules.                         |
| `TRAEFIK_WEB_PORT`    | `80`             | Port for Traefik's HTTP entrypoint.                               |
| `TRAEFIK_WEBHTTPS_PORT` | `443`          | Port for Traefik's HTTPS entrypoint.                              |
| `TRAEFIK_DASHBOARD_PORT` | `8080`         | Port Traefik uses internally to access its dashboard.           |
| **Storage**           |                  |                                                                   |
| `BASE_STORAGE_PATH`   | `./data`         | Base path for all M3TAL data storage. Defaults to `/mnt` in production. |
| `MEDIA_PATH`          | `./data/media`   | Path for storing media files.                                     |
| `CONFIG_PATH`         | `./data/config`  | Path for M3TAL configuration files.                               |
| `DOWNLOADS_PATH`      | `./data/downloads` | Path for downloaded files.                                      |
| `PUID`                | `1000`           | User ID for running containers.                                   |
| `PGID`                | `1000`           | Group ID for running containers.                                  |
| **Dashboard**         |                  |                                                                   |
| `DASHBOARD_PORT`      | `8082`           | Port the dashboard service listens on.                            |
| `DASHBOARD_EXPOSE_MODE` | `local`        | Controls dashboard access: `local` (direct port) or `traefik`.  |
| **VPN**               |                  |                                                                   |
| `VPN_USER`            | `user`           | Username for VPN connection.                                      |
| `VPN_PASSWORD`        | `password`       | Password for VPN connection.                                      |
| **System**            |                  |                                                                   |
| `TZ`                  | `America/Denver` | Timezone for M3TAL services.                                      |

---

## Core

### `HTTP_PORT`

*   **Description:** The port on which the M3TAL API daemon (Go binary) listens for incoming HTTP requests.
*   **Default Value:** `8080`
*   **Example Value:** `5050`
*   **Used By:** API daemon, Dashboard container (via `GO_API_URL`)

### `STATE_DIR`

*   **Description:** The directory where M3TAL stores its SQLite state database (`state.db`).
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal/state`
*   **Used By:** API daemon, Dashboard container

### `LOG_LEVEL`

*   **Description:** Controls the verbosity of logging across M3TAL services. Accepted values typically include `debug`, `info`, `warn`, `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** API daemon, Dashboard container

### `DEBUG_MODE`

*   **Description:** Enables debug mode, which may activate more verbose logging, additional debugging endpoints, or developer-focused features.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** API daemon, Dashboard container

### `METRICS_ENABLED`

*   **Description:** When set to `true`, M3TAL services will expose Prometheus-compatible metrics for monitoring.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** API daemon

## Auth

### `DASHBOARD_SECRET`

*   **Description:** A secret key used for securing session cookies and other sensitive operations within the M3TAL Dashboard. **This variable is auto-generated on the first `m3tal init` command.** Users should **not** set this manually unless rotating the secret.
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `aVerySecureRandomStringGeneratedByM3tal`
*   **Used By:** Dashboard container

### `API_TOKEN`

*   **Description:** A token used for authenticating requests to the M3TAL API daemon. **This variable is auto-generated on the first `m3tal init` command.** Users should **not** set this manually unless rotating the token.
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `aRandomAPIAccessToken12345`
*   **Used By:** API daemon (for internal token validation), CLI binary

### `ADMIN_PASSWORD`

*   **Description:** The password for the initial administrative user created when M3TAL is first set up.
*   **Default Value:** `admin_pass`
*   **Example Value:** `MySecretAdminPassword123`
*   **Used By:** API daemon, Dashboard container

## Network

### `NETWORK_NAME`

*   **Description:** The name of the Docker network that M3TAL services will join. This is crucial for inter-service communication.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal-net`
*   **Used By:** All M3TAL Docker Compose stacks

### `LOCAL_IP`

*   **Description:** The IP address on the host machine that M3TAL services will bind to for local network access. This is often used for internal communication (e.g., `host.docker.internal`).
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** API daemon, Dashboard container

### `DOMAIN`

*   **Description:** The primary domain name configured for your M3TAL instance. This variable is critical for Traefik routing. Setting this variable enables services like `dash.DOMAIN` and `api.DOMAIN` via Traefik.
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.mydomain.com`
*   **Used By:** Traefik gateway, API daemon (potentially for generating URLs)

### `TRAEFIK_WEB_PORT`

*   **Description:** The host port that Traefik will bind to for incoming HTTP traffic.
*   **Default Value:** `80`
*   **Example Value:** `8080`
*   **Used By:** Traefik container

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The host port that Traefik will bind to for incoming HTTPS traffic.
*   **Default Value:** `443`
*   **Example Value:** `8443`
*   **Used By:** Traefik container

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The internal port Traefik uses to expose its own management dashboard. This is typically bound to `127.0.0.1` for secure local access.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** Traefik container

## Storage

### `BASE_STORAGE_PATH`

*   **Description:** The root directory on the host where M3TAL stores its persistent data, including configuration, media, and downloads. **In production deployments, this defaults to `/mnt` to ensure data is stored on a more permanent or intended volume, rather than `./data` as seen in template or development setups.**
*   **Default Value:** `./data`
*   **Example Value:** `/mnt/m3tal_data`
*   **Used By:** API daemon, Dashboard container, all other M3TAL stacks

### `MEDIA_PATH`

*   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files (e.g., uploaded images, videos) are stored.
*   **Default Value:** `./data/media`
*   **Example Value:** `./data/media` (if `BASE_STORAGE_PATH` is `./data`) or `/mnt/m3tal_data/media` (if `BASE_STORAGE_PATH` is `/mnt/m3tal_data`)
*   **Used By:** API daemon, Dashboard container

### `CONFIG_PATH`

*   **Description:** The subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files, including the `state.db` if `STATE_DIR` is relative to `BASE_STORAGE_PATH`.
*   **Default Value:** `./data/config`
*   **Example Value:** `./data/config` (if `BASE_STORAGE_PATH` is `./data`) or `/mnt/m3tal_data/config` (if `BASE_STORAGE_PATH` is `/mnt/m3tal_data`)
*   **Used By:** API daemon, Dashboard container

### `DOWNLOADS_PATH`

*   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files or artifacts are stored.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `./data/downloads` (if `BASE_STORAGE_PATH` is `./data`) or `/mnt/m3tal_data/downloads` (if `BASE_STORAGE_PATH` is `/mnt/m3tal_data`)
*   **Used By:** API daemon, Dashboard container

### `PUID`

*   **Description:** The User ID (UID) that containers will run as. This is important for file permissions on the host system.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** All M3TAL Docker Compose stacks

### `PGID`

*   **Description:** The Group ID (GID) that containers will run as. This is important for file permissions on the host system.
*   **Default Value:** `1000`
*   **Example Value:** `1001`
*   **Used By:** All M3TAL Docker Compose stacks

## Dashboard

### `DASHBOARD_PORT`

*   **Description:** The port on which the M3TAL Dashboard service (Python/Flask container) listens internally.
*   **Default Value:** `8082`
*   **Example Value:** `8083`
*   **Used By:** Dashboard container

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Controls how the M3TAL Dashboard is exposed to the network.
    *   `local`: The dashboard is exposed directly via a port mapping (controlled by `DASHBOARD_PORT`). Access via `http://HOST_IP:DASHBOARD_PORT`. No Traefik is required for this mode.
    *   `traefik`: The dashboard is exposed via the Traefik reverse proxy, using the `dash.${DOMAIN}` hostname. Traefik must be running for this mode to work.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** Dashboard container (via compose override logic), Traefik configuration

## VPN

### `VPN_USER`

*   **Description:** The username used for establishing a VPN connection, if M3TAL is configured to use a VPN for specific outbound traffic.
*   **Default Value:** `user`
*   **Example Value:** `myvpnuser`
*   **Used By:** Components requiring VPN access (specific implementations may vary)

### `VPN_PASSWORD`

*   **Description:** The password for the VPN connection.
*   **Default Value:** `password`
*   **Example Value:** `MyVpnSecurePassword`
*   **Used By:** Components requiring VPN access (specific implementations may vary)

## System

### `TZ`

*   **Description:** The timezone setting for all M3TAL containers and services. This ensures consistent timestamp logging and date/time operations.
*   **Default Value:** `America/Denver`
*   **Example Value:** `UTC` or `Europe/London`
*   **Used By:** API daemon, Dashboard container