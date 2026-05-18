# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, primarily read from `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` and individual `m3tal config set KEY value` commands. Both the M3TAL CLI and all Docker Compose stacks utilize this `.env` file via `--env-file`.

## Quick Reference Table

| Variable Name           | Description                                                                                                                                                                                              | Default Value      |
|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| **Core**                |                                                                                                                                                                                                          |                    |
| `LOG_LEVEL`             | Sets the logging verbosity for M3TAL services.                                                                                                                                                           | `info`             |
| `DEBUG_MODE`            | Enables or disables debug mode for M3TAL services.                                                                                                                                                       | `false`            |
| `METRICS_ENABLED`       | Enables or disables the collection and exposure of service metrics.                                                                                                                                        | `true`             |
| `PUID`                  | The user ID to run containers under.                                                                                                                                                                     | `1000`             |
| `PGID`                  | The group ID to run containers under.                                                                                                                                                                    | `1000`             |
| `TZ`                    | The timezone for M3TAL services and containers.                                                                                                                                                          | `America/Denver`   |
| **Auth**                |                                                                                                                                                                                                          |                    |
| `DASHBOARD_SECRET`      | A secret key used for session management in the M3TAL dashboard. **Auto-generated on first `m3tal init`.**                                                                                                | `change_me_immediately` |
| `API_TOKEN`             | An API token used for authentication with the M3TAL API. **Auto-generated on first `m3tal init`.**                                                                                                      | `change_me_api_token` |
| `ADMIN_PASSWORD`        | The password for the default `admin` user of the M3TAL dashboard.                                                                                                                                        | `admin_pass`       |
| **Network**             |                                                                                                                                                                                                          |                    |
| `HTTP_PORT`             | The internal HTTP port the M3TAL API daemon listens on.                                                                                                                                                  | `8080`             |
| `NETWORK_NAME`          | The name of the Docker network M3TAL services will use.                                                                                                                                                  | `m3tal`            |
| `LOCAL_IP`              | The local IP address that M3TAL services should bind to.                                                                                                                                                 | `127.0.0.1`        |
| **Storage**             |                                                                                                                                                                                                          |                    |
| `BASE_STORAGE_PATH`     | The base directory for storing M3TAL data. **Defaults to `/mnt` in production deployments.**                                                                                                             | `./data`           |
| `MEDIA_PATH`            | The sub-directory within `BASE_STORAGE_PATH` for media files.                                                                                                                                            | `./data/media`     |
| `CONFIG_PATH`           | The sub-directory within `BASE_STORAGE_PATH` for configuration files.                                                                                                                                    | `./data/config`    |
| `DOWNLOADS_PATH`        | The sub-directory within `BASE_STORAGE_PATH` for downloaded files.                                                                                                                                       | `./data/downloads` |
| `STATE_DIR`             | The directory where the M3TAL state database (`state.db`) is stored.                                                                                                                                     | `./state`          |
| **Dashboard**           |                                                                                                                                                                                                          |                    |
| `DASHBOARD_PORT`        | The port the M3TAL dashboard service listens on.                                                                                                                                                         | `8082`             |
| `DASHBOARD_EXPOSE_MODE` | Controls how the M3TAL dashboard is exposed: `local` (direct port binding) or `traefik` (via Traefik reverse proxy).                                                                                    | `local`            |
| **Traefik**             |                                                                                                                                                                                                          |                    |
| `DOMAIN`                | The primary domain name for M3TAL services. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routes via Traefik.                                                                                     | `localhost`        |
| `TRAEFIK_WEB_PORT`      | The host port Traefik listens on for HTTP traffic.                                                                                                                                                       | `80`               |
| `TRAEFIK_WEBHTTPS_PORT` | The host port Traefik listens on for HTTPS traffic (not explicitly configured in provided compose files, but standard).                                                                                 | `443`              |
| `TRAEFIK_DASHBOARD_PORT`| The host port for the Traefik dashboard UI.                                                                                                                                                              | `8080`             |
| **VPN**                 |                                                                                                                                                                                                          |                    |
| `VPN_USER`              | The username for the VPN connection.                                                                                                                                                                     | `user`             |
| `VPN_PASSWORD`          | The password for the VPN connection.                                                                                                                                                                     | `password`         |

---

## Environment Variables Reference

All M3TAL environment variables are read from the `/etc/m3tal/.env` file. This file is managed by the `m3tal config wizard` command and can be individually updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks (via `--env-file`) load these variables.

### Core

These variables control fundamental aspects of M3TAL's operation and logging.

*   **`LOG_LEVEL`**
    *   **Description:** Sets the logging verbosity for M3TAL services. Higher levels provide more detailed logs.
    *   **Default Value:** `info`
    *   **Example Value:** `debug`
    *   **Used By:** M3TAL API daemon, M3TAL Dashboard.

*   **`DEBUG_MODE`**
    *   **Description:** Enables or disables debug mode for M3TAL services. This may provide more verbose output and enable additional debugging features.
    *   **Default Value:** `false`
    *   **Example Value:** `true`
    *   **Used By:** M3TAL API daemon, M3TAL Dashboard.

*   **`METRICS_ENABLED`**
    *   **Description:** Controls whether M3TAL services expose metrics for monitoring.
    *   **Default Value:** `true`
    *   **Example Value:** `false`
    *   **Used By:** M3TAL API daemon.

*   **`PUID`**
    *   **Description:** The user ID (UID) that Docker containers should run as. This is important for file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** M3TAL Dashboard container.

*   **`PGID`**
    *   **Description:** The group ID (GID) that Docker containers should run as. This is important for file permissions.
    *   **Default Value:** `1000`
    *   **Example Value:** `1001`
    *   **Used By:** M3TAL Dashboard container.

*   **`TZ`**
    *   **Description:** The timezone to be used by M3TAL services and containers for consistent time logging and operations.
    *   **Default Value:** `America/Denver`
    *   **Example Value:** `UTC`, `Europe/London`
    *   **Used By:** M3TAL Dashboard container.

### Auth

These variables are crucial for securing your M3TAL instance.

*   **`DASHBOARD_SECRET`**
    *   **Description:** A secret key used for signing session cookies in the M3TAL dashboard. **This variable is auto-generated on the first `m3tal init` and should not be set manually unless rotating.**
    *   **Default Value:** `change_me_immediately` (This is a placeholder; a real secret is generated on init).
    *   **Example Value:** A long, randomly generated string.
    *   **Used By:** M3TAL Dashboard container.

*   **`API_TOKEN`**
    *   **Description:** An API token used for authenticating requests to the M3TAL API. **This variable is auto-generated on the first `m3tal init` and should not be set manually unless rotating.**
    *   **Default Value:** `change_me_api_token` (This is a placeholder; a real token is generated on init).
    *   **Example Value:** A long, randomly generated string.
    *   **Used By:** M3TAL API daemon.

*   **`ADMIN_PASSWORD`**
    *   **Description:** The password for the default `admin` user account in the M3TAL dashboard.
    *   **Default Value:** `admin_pass`
    *   **Example Value:** `MySecurePassword123!`
    *   **Used By:** M3TAL Dashboard container.

### Network

These variables define network configurations for M3TAL services.

*   **`HTTP_PORT`**
    *   **Description:** The internal HTTP port that the M3TAL API daemon listens on.
    *   **Default Value:** `8080`
    *   **Example Value:** `9090`
    *   **Used By:** M3TAL API daemon.

*   **`NETWORK_NAME`**
    *   **Description:** The name of the Docker network that M3TAL services will join. This allows them to communicate with each other.
    *   **Default Value:** `m3tal`
    *   **Example Value:** `m3tal-network`
    *   **Used By:** All M3TAL Docker containers.

*   **`LOCAL_IP`**
    *   **Description:** The local IP address that M3TAL services should bind to. This is often `127.0.0.1` for internal communication.
    *   **Default Value:** `127.0.0.1`
    *   **Example Value:** `0.0.0.0` (use with caution)
    *   **Used By:** M3TAL API daemon.

### Storage

These variables define the locations for M3TAL's data, configuration, and media.

*   **`BASE_STORAGE_PATH`**
    *   **Description:** The base directory where M3TAL stores all its persistent data, including configuration, media, and downloads. **In production deployments, this defaults to `/mnt`. In template environments, it defaults to `./data`.**
    *   **Default Value:** `./data`
    *   **Example Value:** `/opt/m3tal/data`
    *   **Used By:** M3TAL Dashboard container, M3TAL API daemon (indirectly via other paths).

*   **`MEDIA_PATH`**
    *   **Description:** The sub-directory within `BASE_STORAGE_PATH` where media files are stored.
    *   **Default Value:** `./data/media`
    *   **Example Value:** `./data/assets`
    *   **Used By:** M3TAL Dashboard container.

*   **`CONFIG_PATH`**
    *   **Description:** The sub-directory within `BASE_STORAGE_PATH` where M3TAL configuration files are stored.
    *   **Default Value:** `./data/config`
    *   **Example Value:** `./data/settings`
    *   **Used By:** M3TAL Dashboard container.

*   **`DOWNLOADS_PATH`**
    *   **Description:** The sub-directory within `BASE_STORAGE_PATH` where downloaded files are stored.
    *   **Default Value:** `./data/downloads`
    *   **Example Value:** `./data/files`
    *   **Used By:** M3TAL Dashboard container.

*   **`STATE_DIR`**
    *   **Description:** The directory where the M3TAL state database (`state.db`) is stored.
    *   **Default Value:** `./state`
    *   **Example Value:** `/var/lib/m3tal/state`
    *   **Used By:** M3TAL API daemon, M3TAL Dashboard container.

### Dashboard

These variables control the M3TAL dashboard's behavior and accessibility.

*   **`DASHBOARD_PORT`**
    *   **Description:** The internal port that the M3TAL dashboard service listens on.
    *   **Default Value:** `8082`
    *   **Example Value:** `9092`
    *   **Used By:** M3TAL Dashboard container.

*   **`DASHBOARD_EXPOSE_MODE`**
    *   **Description:** Controls how the M3TAL dashboard is exposed to users.
        *   `local`: The dashboard is exposed directly via a host port mapping (defaulting to `DASHBOARD_PORT`). Access is via `http://HOST_IP:DASHBOARD_PORT`.
        *   `traefik`: The dashboard is exposed via the Traefik reverse proxy. Access is via `http://dash.DOMAIN`.
    *   **Default Value:** `local`
    *   **Example Value:** `traefik`
    *   **Used By:** M3TAL Dashboard container (influences compose override).

### Traefik

These variables configure the Traefik reverse proxy for M3TAL.

*   **`DOMAIN`**
    *   **Description:** The primary domain name for your M3TAL instance. Setting this variable enables Traefik to correctly route traffic to services like the dashboard (`dash.DOMAIN`) and API (`api.DOMAIN`).
    *   **Default Value:** `localhost`
    *   **Example Value:** `m3tal.mydomain.com`
    *   **Used By:** Traefik (routing rules), M3TAL Dashboard container (when `DASHBOARD_EXPOSE_MODE=traefik`).

*   **`TRAEFIK_WEB_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTP (non-encrypted) traffic.
    *   **Default Value:** `80`
    *   **Example Value:** `8000`
    *   **Used By:** Traefik container.

*   **`TRAEFIK_WEBHTTPS_PORT`**
    *   **Description:** The host port that Traefik listens on for incoming HTTPS (encrypted) traffic. Note: SSL/TLS configuration for Traefik is managed separately.
    *   **Default Value:** `443`
    *   **Example Value:** `8443`
    *   **Used By:** Traefik container.

*   **`TRAEFIK_DASHBOARD_PORT`**
    *   **Description:** The host port on which Traefik's own dashboard UI is accessible. **This is typically `127.0.0.1:8080` (port 8080 on localhost) for security.**
    *   **Default Value:** `8080`
    *   **Example Value:** `8081`
    *   **Used By:** Traefik container.

### VPN

These variables are used if you configure M3TAL to use a VPN for its network traffic.

*   **`VPN_USER`**
    *   **Description:** The username for your VPN connection.
    *   **Default Value:** `user`
    *   **Example Value:** `myvpnuser`
    *   **Used By:** M3TAL VPN client configuration (if applicable).

*   **`VPN_PASSWORD`**
    *   **Description:** The password for your VPN connection.
    *   **Default Value:** `password`
    *   **Example Value:** `MyVpnPassword123!`
    *   **Used By:** M3TAL VPN client configuration (if applicable).

---

## APT Installation

To install or update M3TAL, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```