```markdown
# Environment Variables Reference

As the M3TAL Ecosystem Documentation Architect, I'm providing this comprehensive reference for all environment variables used by the M3TAL system.

All M3TAL environment variables are centrally managed and read from the configuration file at `/etc/m3tal/.env`. This file is used by both the M3TAL CLI (`/usr/bin/m3tal`) and all Docker Compose stacks via the `--env-file` option, ensuring a consistent configuration across the entire ecosystem. You can manage these variables using the `m3tal config wizard` or `m3tal config set KEY value` commands.

---

## Quick Reference

| Name | Description | Default Value | Component | Group |
| :----------------------- | :------------------------------------------------------ | :-------------- | :--------------------------------------------- | :--------------------------------------- |
| `API_TOKEN` | Auth token for M3TAL API. | `change_me_api_token` | CLI, API daemon | Authentication |
| `ADMIN_PASSWORD` | Default password for initial Dashboard admin user. | `admin_pass` | CLI | Authentication |
| `BASE_STORAGE_PATH` | Base directory for all persistent M3TAL data. | `./data` | Dashboard, User Stacks | Storage Paths |
| `CONFIG_PATH` | Subdirectory for M3TAL configuration and state. | `./data/config` | Dashboard | Storage Paths |
| `DASHBOARD_EXPOSE_MODE` | How the Dashboard is exposed (`local` or `traefik`). | `local` | CLI, Dashboard | Network & Routing (Traefik) |
| `DASHBOARD_PORT` | Port for the M3TAL Dashboard container. | `8082` | Dashboard | API & Dashboard |
| `DASHBOARD_SECRET` | Secret key for Dashboard session management. | `change_me_immediately` | Dashboard | Authentication |
| `DEBUG_MODE` | Enables debug logging and features. | `false` | CLI, API daemon, Dashboard | Core Configuration |
| `DOMAIN` | Primary domain for M3TAL services. | `localhost` | Traefik, Dashboard | Network & Routing (Traefik) |
| `DOWNLOADS_PATH` | Subdirectory for downloaded content. | `./data/downloads` | User Stacks | Storage Paths |
| `HTTP_PORT` | Port for the M3TAL API daemon. | `8080` | API daemon | API & Dashboard |
| `LOCAL_IP` | Local IP address of the host machine. | `127.0.0.1` | CLI, API daemon | Core Configuration |
| `LOG_LEVEL` | Logging verbosity. | `info` | CLI, API daemon, Dashboard | Core Configuration |
| `MEDIA_PATH` | Subdirectory for media files. | `./data/media` | User Stacks | Storage Paths |
| `METRICS_ENABLED` | Enables or disables metrics collection. | `true` | API daemon | Core Configuration |
| `NETWORK_NAME` | Name of the Docker network for inter-service communication. | `m3tal` | All Docker containers | Core Configuration |
| `PGID` | Group ID (GID) for containers. | `1000` | All Docker containers | Core Configuration |
| `PUID` | User ID (UID) for containers. | `1000` | All Docker containers | Core Configuration |
| `STATE_DIR` | Internal container path for state database and users.json. | `./state` | Dashboard | API & Dashboard |
| `TRAEFIK_DASHBOARD_PORT` | Internal port for the Traefik dashboard. | `8080` | Traefik | Network & Routing (Traefik) |
| `TRAEFIK_WEB_PORT` | Port Traefik listens on for HTTP traffic. | `80` | Traefik | Network & Routing (Traefik) |
| `TRAEFIK_WEBHTTPS_PORT` | Port Traefik listens on for HTTPS traffic. | `443` | Traefik | Network & Routing (Traefik) |
| `TZ` | Timezone for containers. | `America/Denver` | All Docker containers | Core Configuration |
| `VPN_PASSWORD` | Password for VPN client configurations. | `password` | VPN containers | VPN Integration |
| `VPN_USER` | Username for VPN client configurations. | `user` | VPN containers | VPN Integration |

---

## Detailed Environment Variable Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system, including logging, user/group IDs, and network settings.

#### `LOG_LEVEL`
*   **Description:** Sets the logging verbosity for M3TAL components. Higher verbosity (e.g., `debug`) can be useful for troubleshooting.
*   **Default:** `info`
*   **Example:** `debug`
*   **Component:** CLI binary, API daemon, Dashboard container

#### `DEBUG_MODE`
*   **Description:** Enables debug logging and features across M3TAL components, potentially exposing more detailed information for development or troubleshooting.
*   **Default:** `false`
*   **Example:** `true`
*   **Component:** CLI binary, API daemon, Dashboard container

#### `METRICS_ENABLED`
*   **Description:** Controls whether M3TAL components collect and expose operational metrics.
*   **Default:** `true`
*   **Example:** `false`
*   **Component:** API daemon

#### `NETWORK_NAME`
*   **Description:** The name of the Docker network used for inter-service communication within the M3TAL ecosystem. All M3TAL-managed containers will connect to this network.
*   **Default:** `m3tal`
*   **Example:** `m3tal-net`
*   **Component:** All Docker containers

#### `LOCAL_IP`
*   **Description:** The local IP address of the host machine. This variable can be used for internal routing or specific service configurations that need to reference the host's IP.
*   **Default:** `127.0.0.1`
*   **Example:** `192.168.1.100`
*   **Component:** CLI binary, API daemon (for internal communication resolution)

#### `TZ`
*   **Description:** Sets the timezone for all Docker containers within the M3TAL ecosystem, ensuring logs and timestamps are consistent.
*   **Default:** `America/Denver`
*   **Example:** `America/New_York`, `Europe/London`
*   **Component:** All Docker containers (e.g., `m3tal-dashboard`)

#### `PUID`
*   **Description:** The User ID (UID) used by containers to ensure correct file permissions for mounted volumes. This should typically match the UID of the user owning the M3TAL storage paths on the host.
*   **Default:** `1000`
*   **Example:** `1000`
*   **Component:** All Docker containers (e.g., `m3tal-dashboard`)

#### `PGID`
*   **Description:** The Group ID (GID) used by containers to ensure correct file permissions for mounted volumes. This should typically match the GID of the group owning the M3TAL storage paths on the host.
*   **Default:** `1000`
*   **Example:** `1000`
*   **Component:** All Docker containers (e.g., `m3tal-dashboard`)

### API & Dashboard

These variables are specific to the M3TAL API daemon and the M3TAL Dashboard container.

#### `HTTP_PORT`
*   **Description:** The port on which the M3TAL API daemon (Go binary) listens for incoming requests. This is typically accessed internally or via the Traefik gateway.
*   **Default:** `8080`
*   **Example:** `8080`
*   **Component:** API daemon

#### `DASHBOARD_PORT`
*   **Description:** The internal port on which the M3TAL Dashboard container runs. If `DASHBOARD_EXPOSE_MODE` is set to `local`, this port is directly exposed on the host.
*   **Default:** `8082`
*   **Example:** `8082`
*   **Component:** M3TAL Dashboard container

#### `STATE_DIR`
*   **Description:** This variable defines the *internal* path within the M3TAL Dashboard container where the SQLite state database (`state.db`) and the user credentials file (`users.json`) are stored. The host-level path for this data is controlled by `CONFIG_PATH`.
*   **Default:** `./state`
*   **Example:** `/docker/state` (internal container path)
*   **Component:** M3TAL Dashboard container

### Authentication

These variables control authentication and security for the M3TAL API and Dashboard.

#### `DASHBOARD_SECRET`
*   **Description:** A critical secret key used by the M3TAL Dashboard for session management, encryption, and overall security.
    **Note:** This variable is auto-generated on the first `m3tal init` run. Users should generally *NOT* set it manually unless performing a key rotation process.
*   **Default:** `change_me_immediately`
*   **Example:** `your_super_long_and_random_dashboard_secret_string`
*   **Component:** M3TAL Dashboard container

#### `API_TOKEN`
*   **Description:** The authentication token required to make requests to the M3TAL API daemon. It secures communication between the CLI, Dashboard, and the API.
    **Note:** This variable is auto-generated on the first `m3tal init` run. Users should generally *NOT* set it manually unless performing a token rotation.
*   **Default:** `change_me_api_token`
*   **Example:** `your_very_secure_random_api_token_here`
*   **Component:** CLI binary, API daemon (for request validation)

#### `ADMIN_PASSWORD`
*   **Description:** The default password set for the initial `admin` user on the M3TAL Dashboard. After initial setup, it's recommended to change this via `m3tal dashpass`.
*   **Default:** `admin_pass`
*   **Example:** `MySecureDashboardPass123`
*   **Component:** CLI binary (specifically `m3tal dashpass` for initial user setup)

### Network & Routing (Traefik)

These variables configure how M3TAL services are exposed and routed, primarily through the Traefik gateway.

#### `DOMAIN`
*   **Description:** The primary domain name for your M3TAL deployment. Setting this variable enables Traefik routing rules, creating access points like `dash.DOMAIN` for the dashboard and `api.DOMAIN` for the API.
*   **Default:** `localhost`
*   **Example:** `myhomelab.com`
*   **Component:** Traefik gateway, M3TAL Dashboard container (when `DASHBOARD_EXPOSE_MODE=traefik`)

#### `DASHBOARD_EXPOSE_MODE`
*   **Description:** Controls how the M3TAL Dashboard is made accessible.
    *   `local` (default): The dashboard is directly exposed on `http://HOST_IP:${DASHBOARD_PORT}`. No Traefik required.
    *   `traefik`: The dashboard is routed via Traefik at `http://dash.${DOMAIN}`. Requires Traefik to be running.
*   **Default:** `local`
*   **Example:** `traefik`
*   **Component:** CLI binary (controls which compose override is used for the dashboard), M3TAL Dashboard container

#### `TRAEFIK_WEB_PORT`
*   **Description:** The external port on the host machine where Traefik listens for incoming HTTP (non-HTTPS) traffic.
*   **Default:** `80`
*   **Example:** `80`
*   **Component:** Traefik gateway

#### `TRAEFIK_WEBHTTPS_PORT`
*   **Description:** The external port on the host machine where Traefik listens for incoming HTTPS traffic.
*   **Default:** `443`
*   **Example:** `443`
*   **Component:** Traefik gateway (implicitly used for HTTPS entrypoints, even if not explicitly defined in provided compose)

#### `TRAEFIK_DASHBOARD_PORT`
*   **Description:** The internal port within the Traefik container where its own dashboard (management UI) runs. This is exposed on `127.0.0.1:8081` by default.
*   **Default:** `8080`
*   **Example:** `8080`
*   **Component:** Traefik gateway

### Storage Paths

These variables define the base and subdirectories for all persistent data managed by M3TAL and user-deployed stacks.

#### `BASE_STORAGE_PATH`
*   **Description:** The top-level directory on the host filesystem where M3TAL stores all its persistent data, including configuration, media, and downloads.
    **Note:** While the template default is `./data`, in production M3TAL deployments, this variable typically defaults to `/mnt` to align with common server storage conventions.
*   **Default:** `./data`
*   **Example:** `/mnt/m3tal`, `/srv/m3tal`
*   **Component:** M3TAL Dashboard container (volume mounts), User-defined stacks

#### `MEDIA_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` designated for media files (e.g., movies, TV shows, music). This path is commonly mounted into media server containers.
*   **Default:** `./data/media`
*   **Example:** `/mnt/m3tal/media`
*   **Component:** User-defined media server containers

#### `CONFIG_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` used for M3TAL's internal configuration files and state data, including the path where the dashboard maps its `STATE_DIR`.
*   **Default:** `./data/config`
*   **Example:** `/mnt/m3tal/config`
*   **Component:** M3TAL Dashboard container (volume mounts for `/docker/state`)

#### `DOWNLOADS_PATH`
*   **Description:** A subdirectory within `BASE_STORAGE_PATH` intended for downloaded content. This path is commonly mounted into download client containers.
*   **Default:** `./data/downloads`
*   **Example:** `/mnt/m3tal/downloads`
*   **Component:** User-defined download client containers

### VPN Integration

These variables are placeholders for optional VPN client integrations within the M3TAL ecosystem, such as WireGuard or OpenVPN clients.

#### `VPN_USER`
*   **Description:** Username for authenticating with an optional VPN service.
*   **Default:** `user`
*   **Example:** `m3tal_vpn_user`
*   **Component:** External VPN client containers (e.g., `cloudflared` or other custom VPN stacks)

#### `VPN_PASSWORD`
*   **Description:** Password for authenticating with an optional VPN service.
*   **Default:** `password`
*   **Example:** `MySecureVPNPassw0rd`
*   **Component:** External VPN client containers (e.g., `cloudflared` or other custom VPN stacks)
```