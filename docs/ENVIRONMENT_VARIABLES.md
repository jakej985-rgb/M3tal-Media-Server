# Environment Variables Reference

All M3TAL environment variables are read from the primary configuration file, `/etc/m3tal/.env`. This file is loaded by the `m3tal` CLI, the `m3tal-api.service`, and all Docker Compose stacks (via the `--env-file` option) to configure the entire M3TAL ecosystem.

You can manage these variables using the `m3tal config wizard` for an interactive setup, or `m3tal config set KEY value` for direct modification.

---

### Quick Reference

| Name | Description | Default Value | Used By |
| :----------------------- | :------------------------------------------------------------------------------------------------------------------------------- | :-------------------- | :-------------------------------------------------------------- |
| `DASHBOARD_PORT` | Port for the M3TAL Dashboard. | `8082` | `m3tal-dashboard`, CLI |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed (direct port binding or via Traefik). | `local` | CLI (`m3tal dash up`) |
| `HTTP_PORT` | Port for the M3TAL API daemon. | `8080` | `m3tal-api.service`, Traefik |
| `DASHBOARD_SECRET` | Secret key for M3TAL Dashboard session management and encryption. | `change_me_immediately` | `m3tal-dashboard` |
| `API_TOKEN` | Authentication token for the M3TAL API. | `change_me_api_token` | `m3tal-api.service`, CLI |
| `ADMIN_PASSWORD` | Default password for the initial M3TAL Dashboard admin user. | `admin_pass` | `m3tal-dashboard` |
| `NETWORK_NAME` | Docker network name for M3TAL components. | `m3tal` | All Docker containers |
| `LOCAL_IP` | IP for `host.docker.internal` within containers. | `127.0.0.1` | All Docker containers |
| `DOMAIN` | Base domain name for Traefik routing rules. | `localhost` | Traefik, `m3tal-dashboard` |
| `BASE_STORAGE_PATH` | Base directory for all M3TAL data. | `./data` | All Docker containers |
| `MEDIA_PATH` | Subdirectory for user media files. | `./data/media` | User stacks |
| `CONFIG_PATH` | Subdirectory for M3TAL configuration files. | `./data/config` | `m3tal-dashboard`, API |
| `DOWNLOADS_PATH` | Subdirectory for downloaded files. | `./data/downloads` | Download clients |
| `STATE_DIR` | Base directory for API daemon's state files. | `./state` | `m3tal-api.service` |
| `TRAEFIK_WEB_PORT` | Traefik's HTTP entry point port. | `80` | Traefik |
| `TRAEFIK_WEBHTTPS_PORT` | Traefik's HTTPS entry point port. | `443` | Traefik |
| `TRAEFIK_DASHBOARD_PORT` | Traefik's internal dashboard port. | `8080` | Traefik |
| `VPN_USER` | Username for VPN client connection. | `user` | VPN client containers |
| `VPN_PASSWORD` | Password for VPN client connection. | `password` | VPN client containers |
| `PUID` | User ID (UID) for containers. | `1000` | All Docker containers |
| `PGID` | Group ID (GID) for containers. | `1000` | All Docker containers |
| `TZ` | Timezone setting for containers. | `America/Denver` | All Docker containers |
| `LOG_LEVEL` | Logging verbosity for the M3TAL API daemon. | `info` | `m3tal-api.service` |
| `DEBUG_MODE` | Enables debug features. | `false` | API, Dashboard |
| `METRICS_ENABLED` | Enables/disables M3TAL API metrics. | `true` | `m3tal-api.service` |

---

### Detailed Variable Reference

#### Core Configuration

These variables control fundamental aspects of the M3TAL system, including default ports and dashboard exposure.

<br>

### `DASHBOARD_PORT`

*   **Description:** The port on which the M3TAL Dashboard container listens. In `local` expose mode, this port is directly bound to the host.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** `m3tal-dashboard` container, M3TAL CLI (when starting the dashboard in local mode).

<br>

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Determines how the M3TAL Dashboard is made accessible.
    *   `local`: The dashboard port (`DASHBOARD_PORT`) is directly bound to the host. Access via `http://HOST_IP:8082`. No Traefik required.
    *   `traefik`: The dashboard is exposed via Traefik, accessible at `http://dash.DOMAIN`. Requires Traefik to be running.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** M3TAL CLI (`m3tal dash up` command) to select the appropriate Docker Compose override file.

<br>

### `HTTP_PORT`

*   **Description:** The port on which the M3TAL API daemon (`m3tal-api.service`) listens on the host. This port is typically not exposed publicly but is accessed internally by containers (via `host.docker.internal`) and Traefik for domain routing.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `m3tal-api.service`, Traefik gateway (for routing to `api.DOMAIN`).

<br>

#### Authentication & Security

These variables manage authentication tokens, secrets, and initial administrative credentials.

<br>

### `DASHBOARD_SECRET`

*   **Description:** A secret key used by the M3TAL Dashboard for secure session management, encryption, and other security-sensitive operations. **Critical for security.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a_long_random_string_of_characters_for_security_12345`
*   **Used By:** `m3tal-dashboard` container.
*   **Note:** This value is **auto-generated** by `m3tal init` on its first run. Users should generally **NOT** set it manually unless performing a secret rotation.

<br>

### `API_TOKEN`

*   **Description:** An authentication token used by the M3TAL CLI and other components to securely communicate with the M3TAL API daemon.
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `another_long_random_string_for_api_auth_67890`
*   **Used By:** `m3tal-api.service` (for verification), M3TAL CLI (for making authenticated API requests).
*   **Note:** This value is **auto-generated** by `m3tal init` on its first run. Users should generally **NOT** set it manually unless performing a token rotation.

<br>

### `ADMIN_PASSWORD`

*   **Description:** The default password for the initial administrator user created in the M3TAL Dashboard. This is used when the `users.json` file is first created. It is highly recommended to change this password immediately after the first login.
*   **Default Value:** `admin_pass`
*   **Example Value:** `secure_admin_password_123`
*   **Used By:** `m3tal-dashboard` (during initial setup of `users.json`).

<br>

#### Networking

Variables related to Docker networking and host communication.

<br>

### `NETWORK_NAME`

*   **Description:** Specifies the name of the Docker network that all M3TAL core components and user-deployed stacks will use for internal communication. This allows containers to resolve each other by name.
*   **Default Value:** `m3tal`
*   **Example Value:** `my_m3tal_network`
*   **Used By:** All M3TAL Docker containers (e.g., `m3tal-dashboard`, `traefik`, `cloudflared`) and user-defined Docker Compose stacks.

<br>

### `LOCAL_IP`

*   **Description:** The IP address that `host.docker.internal` resolves to within all Docker containers. This is crucial for containers to be able to reach services running directly on the host machine, such as the `m3tal-api.service`.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100` (the actual LAN IP address of your host machine)
*   **Used By:** All M3TAL Docker Compose stacks (configured via `extra_hosts`).

<br>

#### Storage & Paths

Defines the filesystem locations for M3TAL's data, configuration, and media.

<br>

### `BASE_STORAGE_PATH`

*   **Description:** The root directory on the host machine where all M3TAL-related media, configuration, and state data will be stored. All other path variables (e.g., `MEDIA_PATH`, `CONFIG_PATH`, `DOWNLOADS_PATH`) are typically subdirectories of this base path.
*   **Default Value:** `./data`
*   **Example Value:** `/var/m3tal/data`
*   **Used By:** All Docker containers (for volume mounts).
*   **Note:** While the template defaults to `./data` for local testing, in **production deployments**, this variable typically defaults to `/mnt` to leverage dedicated storage volumes.

<br>

### `MEDIA_PATH`

*   **Description:** A subdirectory within `BASE_STORAGE_PATH` designated for user media files (e.g., movies, TV shows, music). This path is typically mounted into media server containers.
*   **Default Value:** `./data/media`
*   **Example Value:** `${BASE_STORAGE_PATH}/media` (e.g., `/mnt/media`)
*   **Used By:** User-defined media stacks.

<br>

### `CONFIG_PATH`

*   **Description:** A subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its core configuration and state files. This includes the `m3tal/state` directory which is mounted into the dashboard container, containing `users.json` and other stateful data used by the dashboard.
*   **Default Value:** `./data/config`
*   **Example Value:** `${BASE_STORAGE_PATH}/config` (e.g., `/mnt/config`)
*   **Used By:** `m3tal-dashboard` container (for mounting state directory), other stacks requiring persistent configuration.

<br>

### `DOWNLOADS_PATH`

*   **Description:** A subdirectory within `BASE_STORAGE_PATH` intended for downloaded files. This path is commonly mounted into download client containers (e.g., torrent clients, newsgroup clients).
*   **Default Value:** `./data/downloads`
*   **Example Value:** `${BASE_STORAGE_PATH}/downloads` (e.g., `/mnt/downloads`)
*   **Used By:** Download client containers.

<br>

### `STATE_DIR`

*   **Description:** Specifies the base directory on the host where the M3TAL API daemon (`m3tal-api.service`) stores its state files, including the primary SQLite database (`state.db`).
*   **Default Value:** `./state`
*   **Example Value:** `/var/lib/m3tal`
*   **Used By:** `m3tal-api.service`.
*   **Note:** In standard systemd deployments, the API daemon's primary SQLite database (`state.db`) is typically configured to use the absolute path `/var/lib/m3tal/state.db`, which may override or take precedence over this variable's effect for the main database file itself. This variable might be more relevant for other auxiliary state files or during development.

<br>

#### Traefik Gateway

Variables specific to the Traefik reverse proxy for controlling routing and dashboard access.

<br>

### `DOMAIN`

*   **Description:** The base domain name used by Traefik to define routing rules for M3TAL services. Setting this variable enables Traefik to route requests for `dash.DOMAIN` to the M3TAL Dashboard and `api.DOMAIN` to the M3TAL API.
*   **Default Value:** `localhost`
*   **Example Value:** `example.com`
*   **Used By:** Traefik gateway, `m3tal-dashboard` (when `DASHBOARD_EXPOSE_MODE=traefik`).

<br>

### `TRAEFIK_WEB_PORT`

*   **Description:** The port on the host machine that Traefik listens on for incoming HTTP traffic. This is Traefik's primary entry point for web services.
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** `traefik` container.

<br>

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The port on the host machine that Traefik listens on for incoming HTTPS traffic. This is typically used when SSL/TLS certificates are configured for secure web access.
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** `traefik` container.

<br>

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The internal port used by Traefik for its own administrative dashboard within the container.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `traefik` container.
*   **Note:** In standard M3TAL deployments, Traefik's internal dashboard port `8080` is mapped to `127.0.0.1:8081` on the host, ensuring it's only accessible locally for security.

<br>

#### VPN Integration

Variables for configuring VPN client containers that may be part of user-deployed stacks.

<br>

### `VPN_USER`

*   **Description:** The username required for authentication when connecting a VPN client container (e.g., OpenVPN, WireGuard) to a VPN service.
*   **Default Value:** `user`
*   **Example Value:** `my_vpn_username`
*   **Used By:** VPN client containers in user-defined stacks.

<br>

### `VPN_PASSWORD`

*   **Description:** The password associated with the `VPN_USER` for authenticating a VPN client connection.
*   **Default Value:** `password`
*   **Example Value:** `my_secure_vpn_password`
*   **Used By:** VPN client containers in user-defined stacks.

<br>

#### System & Runtime

General system-level settings, including user/group IDs, timezone, and logging.

<br>

### `PUID`

*   **Description:** The User ID (UID) that Docker containers will use to run their processes. Setting this to match a non-root user on your host machine helps ensure correct file ownership and permissions for mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1000` (common for the first non-root user on Linux)
*   **Used By:** All M3TAL Docker containers (e.g., `m3tal-dashboard`) and user-defined Docker Compose stacks.

<br>

### `PGID`

*   **Description:** The Group ID (GID) that Docker containers will use. Similar to `PUID`, setting this to match a group on your host machine helps ensure correct file ownership and permissions for mounted volumes.
*   **Default Value:** `1000`
*   **Example Value:** `1000` (common for the first non-root user's primary group on Linux)
*   **Used By:** All M3TAL Docker containers and user-defined Docker Compose stacks.

<br>

### `TZ`

*   **Description:** Sets the timezone for all M3TAL Docker containers. This ensures logs, timestamps, and scheduled tasks operate within the correct local time.
*   **Default Value:** `America/Denver`
*   **Example Value:** `America/New_York`, `Europe/London`, `Asia/Tokyo`
*   **Used By:** All M3TAL Docker containers (e.g., `m3tal-dashboard`) and user-defined Docker Compose stacks.

<br>

### `LOG_LEVEL`

*   **Description:** Controls the verbosity of logging output for the M3TAL API daemon. Useful for debugging or reducing log noise.
*   **Default Value:** `info`
*   **Example Value:** `debug`, `warn`, `error`, `fatal`, `panic`
*   **Used By:** `m3tal-api.service`.

<br>

### `DEBUG_MODE`

*   **Description:** When set to `true`, enables extended debug logging and potentially other debug-specific features within M3TAL components.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` (and potentially other future components).

<br>

### `METRICS_ENABLED`

*   **Description:** When set to `true`, the M3TAL API daemon collects and exposes internal operational metrics (e.g., for Prometheus).
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** `m3tal-api.service`.