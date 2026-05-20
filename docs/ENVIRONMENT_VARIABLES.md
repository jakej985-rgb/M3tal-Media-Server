# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, which are read from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks via the `--env-file` option. The `m3tal config wizard` and `m3tal config set` commands are the recommended ways to manage these variables.

This document provides a comprehensive reference for all available environment variables, grouped by their functional area.

## Quick Reference Table

| Variable Name | Description | Default Value | Example Value | Used By |
|---|---|---|---|---|
| **Core** | | | | |
| `DASHBOARD_PORT` | Port for the M3TAL Dashboard service. | `8082` | `8082` | `m3tal-dashboard` |
| `HTTP_PORT` | Port for the M3TAL API daemon. | `8080` | `8080` | `m3tal-api` |
| `STATE_DIR` | Directory for the M3TAL state database. | `./state` | `/var/lib/m3tal/state` | `m3tal-api`, `m3tal-dashboard` |
| `LOG_LEVEL` | Logging level for M3TAL services. | `info` | `debug` | `m3tal-api`, `m3tal-dashboard` |
| `DEBUG_MODE` | Enables or disables debug mode. | `false` | `true` | `m3tal-api`, `m3tal-dashboard` |
| `METRICS_ENABLED` | Enables or disables metrics collection. | `true` | `false` | `m3tal-api` |
| **Auth** | | | | |
| `DASHBOARD_SECRET` | Secret key for the M3TAL Dashboard. **Auto-generated on first `m3tal init`.** | `change_me_immediately` | `a_secure_random_string` | `m3tal-dashboard` |
| `API_TOKEN` | API token for accessing M3TAL services. **Auto-generated on first `m3tal init`.** | `change_me_api_token` | `another_secure_random_string` | `m3tal-api` |
| `ADMIN_PASSWORD` | Password for the M3TAL dashboard administrator. | `admin_pass` | `a_strong_password` | `m3tal-dashboard` |
| **Network** | | | | |
| `NETWORK_NAME` | The name of the Docker network M3TAL services will join. | `m3tal` | `m3tal_network` | All M3TAL containers |
| `LOCAL_IP` | The IP address of the host machine. Used for internal communication. | `127.0.0.1` | `192.168.1.100` | `m3tal-api` (via `host.docker.internal`) |
| `DOMAIN` | The primary domain name for M3TAL services. Used for Traefik routing. | `localhost` | `m3tal.example.com` | `traefik` |
| `TRAEFIK_WEB_PORT` | The port Traefik listens on for HTTP traffic. | `80` | `80` | `traefik` |
| `TRAEFIK_WEBHTTPS_PORT` | The port Traefik listens on for HTTPS traffic. | `443` | `443` | `traefik` |
| `TRAEFIK_DASHBOARD_PORT` | The port Traefik exposes its own dashboard on. | `8080` | `8081` | `traefik` |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed (`local` or `traefik`). See Dashboard Access section for details. | `local` | `traefik` | `m3tal-dashboard` |
| **Storage** | | | | |
| `BASE_STORAGE_PATH` | Base directory for all M3TAL data. Defaults to `/mnt` in production. | `./data` | `/mnt/m3tal_storage` | All M3TAL containers |
| `MEDIA_PATH` | Directory for media storage. | `./data/media` | `/mnt/m3tal_storage/media` | `m3tal-dashboard` |
| `CONFIG_PATH` | Directory for M3TAL configuration files. | `./data/config` | `/mnt/m3tal_storage/config` | `m3tal-api`, `m3tal-dashboard` |
| `DOWNLOADS_PATH` | Directory for downloaded files. | `./data/downloads` | `/mnt/m3tal_storage/downloads` | `m3tal-dashboard` |
| **System** | | | | |
| `PUID` | User ID for running containers. | `1000` | `1001` | All M3TAL containers |
| `PGID` | Group ID for running containers. | `1000` | `1001` | All M3TAL containers |
| `TZ` | Timezone for M3TAL services. | `America/Denver` | `Etc/UTC` | `m3tal-api`, `m3tal-dashboard` |
| **VPN** | | | | |
| `VPN_USER` | Username for the VPN connection. | `user` | `myvpnuser` | `m3tal-dashboard` |
| `VPN_PASSWORD` | Password for the VPN connection. | `password` | `myvpnpassword` | `m3tal-dashboard` |

---

## Detailed Environment Variable Reference

All environment variables are read from `/etc/m3tal/.env`.

### Core

These variables control fundamental aspects of M3TAL's operation.

*   **`DASHBOARD_PORT`**
    *   **Description:** The port on which the M3TAL Dashboard service listens.
    *   **Default Value:** `8082`
    *   **Example Value:** `8082`
    *   **Used By:** `m3tal-dashboard` container.

*   **`HTTP_PORT`**
    *   **Description:** The port on which the M3TAL API daemon listens for incoming HTTP requests.
    *   **Default Value:** `8080`
    *   **Example Value:** `8080`
    *   **Used By:** `m3tal-api` service.

*   **`STATE_DIR`**
    *   **Description:** The directory where the M3TAL state database (`state.db`) is stored.
    *   **Default Value:** `./state` (relative to `CONFIG_PATH` when running in a container)
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** `m3tal-api` service and `m3tal-dashboard` container.

*   **`LOG_LEVEL`**
    *   **Description:** Sets the verbosity of logging for M3TAL services. Can be `debug`, `info`, `warn`, `error`.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** `m3tal-api` service and `m3tal-dashboard` container.

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode for M3TAL services. When `true`, additional debugging information may be logged or features enabled.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** `m3tal-api` service and `m3tal-dashboard` container.

*   **`METRICS_ENABLED`**
    *   **Description:** Enables or disables the collection and exposure of performance metrics.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** `m3tal-api` service.

### Auth

These variables are crucial for securing your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for signing session cookies and other security-sensitive operations within the M3TAL Dashboard. **This variable is auto-generated on the first `m3tal init` and should not be manually set unless you intend to rotate it.**
    *   **Default Value:** `change_me_immediately`
    *   **Example Value:** `a_very_secure_and_random_string_generated_by_m3tal_init`
    *   **Used By:** `m3tal-dashboard` container.

*   **`API_TOKEN`**
    *   **Description:** A token used for authenticating API requests to the M3TAL API daemon. **This variable is auto-generated on the first `m3tal init` and should not be manually set unless you intend to rotate it.**
    *   **Default Value:** `change_me_api_token`
    *   **Example Value:** `another_secure_token_for_api_access`
    *   **Used By:** `m3tal-api` service.

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default administrator account on the M3TAL Dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `a_strong_and_unique_password`
    *   **Used By:** `m3tal-dashboard` container.

### Network

These variables configure M3TAL's network connectivity and integration with external services like Traefik.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will join. This ensures proper communication between containers.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal_network`
    *   **Used By:** All M3TAL containers.

*   **`LOCAL_IP`**
    *   **Description:** The IP address of the host machine. This is primarily used internally by Docker to expose services via `host.docker.internal` for containers to communicate with the host.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `192.168.1.100`
    *   **Used By:** `m3tal-api` service (internally via `host.docker.internal`).

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL deployment. This is crucial for Traefik routing rules, enabling access to services like `dash.${DOMAIN}` and `api.${DOMAIN}`.
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.example.com`
    *   **Used By:** `traefik` container.

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTP traffic. This is the public-facing HTTP port.
    *   **Default Value:** `80`
    *   **Example Value:** `80`
    *   **Used By:** `traefik` container.

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTPS traffic. This is the public-facing HTTPS port.
    *   **Default Value:** `443`
    *   **Example Value:** `443`
    *   **Used By:** `traefik` container.

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The port Traefik itself exposes its management dashboard on. This is typically accessed internally via `localhost:8081`.
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** `traefik` container.

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL Dashboard is exposed to the network.
        *   `local`: Exposes the dashboard directly via a port binding (default `8082`). Access via `http://HOST_IP:8082` or `http://localhost:8082`. No Traefik is required. Best for LAN-only or initial setups.
        *   `traefik`: Configures Traefik to route traffic to the dashboard via `http://dash.${DOMAIN}`. Requires Traefik to be running. Best for domain-based setups behind a reverse proxy.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** `m3tal-dashboard` container (via Docker Compose overrides).

### Storage

These variables define the locations for M3TAL's persistent data.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The base directory where all M3TAL persistent data is stored. In production deployments, this defaults to `/mnt` to utilize mounted drives. In template or development setups, it may default to `./data`.
    *   **Default Value:** `./data`
    *   **Example Value:** `/mnt/m3tal_storage`
    *   **Used By:** All M3TAL containers for their respective volume mounts.

*   **`MEDIA_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `/mnt/m3tal_storage/media`
    *   **Used By:** `m3tal-dashboard` container.

*   **`CONFIG_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where configuration files and the state database are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `/mnt/m3tal_storage/config`
    *   **Used By:** `m3tal-api` service and `m3tal-dashboard` container.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The subdirectory within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `/mnt/m3tal_storage/downloads`
    *   **Used By:** `m3tal-dashboard` container.

### System

These variables control system-level settings for containers.

*   **`PUID`**
    *   **Description:** The User ID (UID) that the containers will run as. This is important for ensuring correct file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** All M3TAL containers.

*   **`PGID`**
    *   **Description:** The Group ID (GID) that the containers will run as. This is important for ensuring correct file permissions on the host system.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** All M3TAL containers.

*   **`TZ`**
    *   **Description:** The timezone to be used by the containers. This ensures correct timekeeping for logs and application functionality.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `Etc/UTC`
    *   **Used By:** `m3tal-api` service and `m3tal-dashboard` container.

### VPN

These variables are used if you configure M3TAL to use a VPN.

*   **`VPN_USER`**
    *   **Description:** The username for your VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** `m3tal-dashboard` container (if VPN configuration is active).

*   **`VPN_PASSWORD`**
    *   **Description:** The password for your VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `myvpnpassword`
    *   **Used By:** `m3tal-dashboard` container (if VPN configuration is active).