# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, primarily read from `/etc/m3tal/.env` by both the CLI and all Docker Compose stacks via `--env-file`. This file is managed by the `m3tal config wizard` and can be updated with `m3tal config set KEY value`.

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem.

## Quick Reference Table

| Variable Name           | Description                                                                                              | Default Value       | Example Value                         | Used By                                 |
| ----------------------- | -------------------------------------------------------------------------------------------------------- | ------------------- | ------------------------------------- | --------------------------------------- |
| **Core**                |                                                                                                          |                     |                                       |                                         |
| `STATE_DIR`             | Directory to store M3TAL's state database.                                                               | `./state`           | `/var/lib/m3tal/state`                | API Daemon, CLI                         |
| `LOG_LEVEL`             | Sets the logging level for M3TAL services.                                                               | `info`              | `debug`                               | API Daemon, CLI                         |
| `TZ`                    | Timezone for M3TAL services.                                                                             | `America/Denver`    | `UTC`                                 | Dashboard Container, API Daemon         |
| `DEBUG_MODE`            | Enables debug mode for M3TAL services.                                                                   | `false`             | `true`                                | API Daemon, CLI                         |
| `METRICS_ENABLED`       | Enables Prometheus metrics collection.                                                                   | `true`              | `false`                               | API Daemon                              |
| **Authentication**      |                                                                                                          |                     |                                       |                                         |
| `DASHBOARD_SECRET`      | Secret key for session management in the dashboard. Auto-generated on first `m3tal init`.                  | `change_me_immediately` | `s3cr3t_k3y_g3n3r4t3d_by_m3tal`     | Dashboard Container                     |
| `API_TOKEN`             | API token for authentication. Auto-generated on first `m3tal init`.                                      | `change_me_api_token` | `4p1_t0k3n_g3n3r4t3d_by_m3tal`      | API Daemon, CLI                         |
| `ADMIN_PASSWORD`        | Password for the default admin user in the dashboard.                                                    | `admin_pass`        | `my_secure_password`                  | Dashboard Container                     |
| **Network**             |                                                                                                          |                     |                                       |                                         |
| `HTTP_PORT`             | The port the API daemon listens on.                                                                      | `8080`              | `8081`                                | API Daemon                              |
| `NETWORK_NAME`          | The name of the Docker network M3TAL services will use.                                                  | `m3tal`             | `m3tal-net`                           | CLI, Docker Compose                     |
| `LOCAL_IP`              | The local IP address of the host machine.                                                                | `127.0.0.1`         | `192.168.1.100`                       | API Daemon, Dashboard Container         |
| `DOMAIN`                | The primary domain name for M3TAL services. Used for Traefik routing.                                    | `localhost`         | `mydomain.com`                        | Traefik, CLI                            |
| **Storage**             |                                                                                                          |                     |                                       |                                         |
| `BASE_STORAGE_PATH`     | Base path for M3TAL data storage. Defaults to `/mnt` in production deployments.                          | `./data`            | `/mnt`                                | CLI, Docker Compose (Dashboard, API)    |
| `MEDIA_PATH`            | Path for media files within the `BASE_STORAGE_PATH`.                                                     | `./data/media`      | `/mnt/media`                          | CLI, Docker Compose (Dashboard, API)    |
| `CONFIG_PATH`           | Path for configuration files within the `BASE_STORAGE_PATH`.                                             | `./data/config`     | `/mnt/config`                         | CLI, Docker Compose (Dashboard, API)    |
| `DOWNLOADS_PATH`        | Path for downloads within the `BASE_STORAGE_PATH`.                                                       | `./data/downloads`  | `/mnt/downloads`                      | CLI, Docker Compose (Dashboard, API)    |
| **Traefik**             |                                                                                                          |                     |                                       |                                         |
| `DASHBOARD_PORT`        | The port the dashboard container listens on internally.                                                  | `8082`              | `8083`                                | Dashboard Container, Traefik            |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed: `local` (direct port) or `traefik` (via Traefik).               | `local`             | `traefik`                             | CLI, Docker Compose (Dashboard)         |
| `TRAEFIK_WEB_PORT`      | The host port Traefik uses for HTTP traffic.                                                             | `80`                | `8000`                                | Traefik                                 |
| `TRAEFIK_WEBHTTPS_PORT` | The host port Traefik uses for HTTPS traffic.                                                            | `443`               | `8443`                                | Traefik                                 |
| `TRAEFIK_DASHBOARD_PORT`| The host port Traefik listens on for its own dashboard.                                                  | `8080`              | `8081`                                | Traefik                                 |
| **User/System**         |                                                                                                          |                     |                                       |                                         |
| `PUID`                  | User ID for running containers.                                                                          | `1000`              | `1001`                                | Docker Compose (Dashboard)              |
| `PGID`                  | Group ID for running containers.                                                                         | `1000`              | `1001`                                | Docker Compose (Dashboard)              |
| **VPN**                 |                                                                                                          |                     |                                       |                                         |
| `VPN_USER`              | Username for VPN connection.                                                                             | `user`              | `myvpnuser`                           | CLI, Docker Compose                     |
| `VPN_PASSWORD`          | Password for VPN connection.                                                                             | `password`          | `mysecurevpnpassword`                 | CLI, Docker Compose                     |

---

## M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, primarily read from `/etc/m3tal/.env` by both the CLI and all Docker Compose stacks via `--env-file`. This file is managed by the `m3tal config wizard` and can be updated with `m3tal config set KEY value`.

### Core

These variables control fundamental aspects of M3TAL's operation.

*   **`STATE_DIR`**
    *   **Description:** Specifies the directory where M3TAL's state database (`state.db`) will be stored. This database holds critical information about your M3TAL installation.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** API Daemon, CLI

*   **`LOG_LEVEL`**
    *   **Description:** Determines the verbosity of logging output for M3TAL services. Options typically include `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** API Daemon, CLI

*   **`TZ`**
    *   **Description:** Sets the timezone for M3TAL services, ensuring consistent timestamping of logs and events.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`
    *   **Used By:** Dashboard Container, API Daemon

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode for M3TAL services. When enabled, more verbose debugging information may be logged.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** API Daemon, CLI

*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether M3TAL exposes Prometheus metrics for monitoring.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** API Daemon

### Authentication

These variables are crucial for securing your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for signing session cookies in the M3TAL dashboard. This is **auto-generated** on the first `m3tal init` and should **not be set manually** unless you intend to rotate it.
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `s3cr3t_k3y_g3n3r4t3d_by_m3tal`
    *   **Used By:** Dashboard Container

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating API requests. This is **auto-generated** on the first `m3tal init` and should **not be set manually** unless you intend to rotate it.
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `4p1_t0k3n_g3n3r4t3d_by_m3tal`
    *   **Used By:** API Daemon, CLI

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator account in the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `my_secure_password`
    *   **Used By:** Dashboard Container

### Network

These variables configure network-related settings for M3TAL.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** API Daemon

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will connect to. This allows containers to communicate with each other.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-net`
    *   **Used By:** CLI, Docker Compose

*   **`LOCAL_IP`**
    *   **Description:** The IP address of the host machine that M3TAL services can reach. This is particularly important for inter-container communication.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** API Daemon, Dashboard Container

*   **`DOMAIN`**
    *   **Description:** The primary domain name configured for your M3TAL instance. This variable is crucial for Traefik routing rules. When set, Traefik will configure routes for `dash.${DOMAIN}` and `api.${DOMAIN}`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `mydomain.com`
    *   **Used By:** Traefik, CLI

### Storage

These variables define where M3TAL stores its data.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The root directory for all M3TAL-related data storage. In production deployments, this defaults to `/mnt`. In template/development environments, it might default to `./data`. This controls where media data is stored.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt`
    *   **Used By:** CLI, Docker Compose (Dashboard, API)

*   **`MEDIA_PATH`**
    *   **Description:** The specific path within `BASE_STORAGE_PATH` where media files (e.g., uploaded images, videos) are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/media`
    *   **Used By:** CLI, Docker Compose (Dashboard, API)

*   **`CONFIG_PATH`**
    *   **Description:** The specific path within `BASE_STORAGE_PATH` where M3TAL configuration files are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/config`
    *   **Used By:** CLI, Docker Compose (Dashboard, API)

*   **`DOWNLOADS_PATH`**
    *   **Description:** The specific path within `BASE_STORAGE_PATH` used for downloaded content or temporary storage of downloaded files.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/downloads`
    *   **Used By:** CLI, Docker Compose (Dashboard, API)

### Traefik

These variables configure the Traefik reverse proxy integration.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port the M3TAL dashboard container listens on. This is distinct from the host port if `DASHBOARD_EXPOSE_MODE` is `local`.
    *   **Default Value:** `8082`
    *   **Example Value:** `8083`
    *   **Used By:** Dashboard Container, Traefik

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Determines how the M3TAL dashboard is exposed.
        *   `local`: The dashboard is exposed directly via a port binding (controlled by `DASHBOARD_PORT`). Access via `http://HOST_IP:DASHBOARD_PORT`. No Traefik is required for dashboard access.
        *   `traefik`: The dashboard is exposed via Traefik. Access via `http://dash.${DOMAIN}`. Traefik must be running.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** CLI, Docker Compose (Dashboard)

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTP traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `8000`
    *   **Used By:** Traefik

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTPS traffic.
    *   **Default Value:** `443`
    *   **Example Value:** `8443`
    *   **Used By:** Traefik

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The host port on which Traefik's own dashboard is accessible.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** Traefik

### User/System

These variables relate to user and system-level configurations for container execution.

*   **`PUID`**
    *   **Description:** The User ID (UID) that Docker containers should run as. This is important for file permissions when containers access host volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** Docker Compose (Dashboard)

*   **`PGID`**
    *   **Description:** The Group ID (GID) that Docker containers should run as. This is important for file permissions when containers access host volumes.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** Docker Compose (Dashboard)

### VPN

These variables are used if M3TAL is configured to use a VPN for network access.

*   **`VPN_USER`**
    *   **Description:** The username for establishing a VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** CLI, Docker Compose

*   **`VPN_PASSWORD`**
    *   **Description:** The password for establishing a VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `mysecurevpnpassword`
    *   **Used By:** CLI, Docker Compose