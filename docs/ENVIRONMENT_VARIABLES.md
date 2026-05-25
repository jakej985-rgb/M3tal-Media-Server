# Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. These variables control various aspects of M3TAL's behavior, from network configuration to storage paths and authentication.

All M3TAL environment variables are read from the `/etc/m3tal/.env` file. This file serves as the central configuration hub for both the `m3tal` CLI and all Docker Compose stacks, which load it via the `--env-file` argument. You can manage these variables using the `m3tal config wizard` or `m3tal config set` commands.

## Quick Reference

| Variable Name            | Description                                                        | Default Value          | Used By                                     |
|--------------------------|--------------------------------------------------------------------|------------------------|---------------------------------------------|
| `HTTP_PORT`              | Port for M3TAL API daemon                                          | `8080`                 | `m3tal-api.service`, `m3tal-dashboard`      |
| `STATE_DIR`              | Directory for API daemon's state DB                                | `./state`              | `m3tal-api.service`                         |
| `LOG_LEVEL`              | Logging verbosity for API daemon                                   | `info`                 | `m3tal-api.service`                         |
| `NETWORK_NAME`           | Docker network name for M3TAL components                           | `m3tal`                | All Compose Stacks, `m3tal` CLI             |
| `PUID`                   | User ID for container processes                                    | `1000`                 | All Compose Stacks                          |
| `PGID`                   | Group ID for container processes                                   | `1000`                 | All Compose Stacks                          |
| `TZ`                     | Timezone for containers                                            | `America/Denver`       | All Compose Stacks                          |
| `DEBUG_MODE`             | Enables debug logging and features                                 | `false`                | `m3tal-api.service`, `m3tal-dashboard`      |
| `METRICS_ENABLED`        | Enables/disables metrics collection                                | `true`                 | `m3tal-api.service`, `m3tal-dashboard`      |
| `DASHBOARD_SECRET`       | Secret key for Dashboard sessions                                  | `change_me_immediately`| `m3tal-dashboard`                           |
| `API_TOKEN`              | Auth token for API daemon                                          | `change_me_api_token`  | `m3tal-api.service`, `m3tal` CLI, `m3tal-dashboard` |
| `ADMIN_PASSWORD`         | Default admin password for Dashboard                               | `admin_pass`           | `m3tal-dashboard`                           |
| `LOCAL_IP`               | Host's local IP address                                            | `127.0.0.1`            | Traefik, `m3tal` CLI                        |
| `BASE_STORAGE_PATH`      | Base directory for all M3TAL data                                  | `./data`               | All Compose Stacks                          |
| `MEDIA_PATH`             | Subdirectory for media files                                       | `./data/media`         | User Stacks                                 |
| `CONFIG_PATH`            | Subdirectory for configuration files                               | `./data/config`        | `m3tal-dashboard`, User Stacks              |
| `DOWNLOADS_PATH`         | Subdirectory for downloaded files                                  | `./data/downloads`     | User Stacks                                 |
| `DASHBOARD_PORT`         | Port for M3TAL Dashboard (local mode)                              | `8082`                 | `m3tal-dashboard`, `m3tal` CLI              |
| `DASHBOARD_EXPOSE_MODE`  | Dashboard exposure mode (`local` or `traefik`)                     | `local`                | `m3tal` CLI, `m3tal-dashboard`              |
| `DOMAIN`                 | Primary domain for Traefik routing                                 | `localhost`            | Traefik, `m3tal` CLI                        |
| `TRAEFIK_WEB_PORT`       | Host port for Traefik's HTTP entry point                           | `80`                   | `traefik` container                         |
| `TRAEFIK_WEBHTTPS_PORT`  | Host port for Traefik's HTTPS entry point                          | `443`                  | `traefik` container                         |
| `TRAEFIK_DASHBOARD_PORT` | Internal Traefik dashboard port (host-bound)                       | `8080`                 | `traefik` container                         |
| `VPN_USER`               | Username for VPN services                                          | `user`                 | VPN Stacks                                  |
| `VPN_PASSWORD`           | Password for VPN services                                          | `password`             | VPN Stacks                                  |

## Detailed Reference

### Core Configuration

These variables define the fundamental operational parameters for the M3TAL system, including API communication, logging, and general container execution.

---

Name: `HTTP_PORT`
Description: The port on which the M3TAL API daemon (Go binary) listens for incoming requests. This is the primary interface for the `m3tal` CLI and the Dashboard to communicate with the core system.
Default: `8080`
Example: `9000`
Used by: `m3tal-api.service`, `m3tal-dashboard` container (connects to `host.docker.internal:${HTTP_PORT}`).

---

Name: `STATE_DIR`
Description: The absolute path on the host filesystem where the M3TAL API daemon stores its SQLite state database (`state.db`) and other runtime data.
Default: `./state`
Example: `/var/lib/m3tal/state`
Used by: `m3tal-api.service`.

---

Name: `LOG_LEVEL`
Description: Controls the verbosity of logging for the M3TAL API daemon. Higher verbosity (e.g., `debug`) provides more detailed output, useful for troubleshooting.
Default: `info`
Example: `debug`
Used by: `m3tal-api.service`.

---

Name: `NETWORK_NAME`
Description: The name of the Docker bridge network that all M3TAL core components and user-defined Docker Compose stacks connect to. This network facilitates inter-container communication.
Default: `m3tal`
Example: `m3tal_internal_net`
Used by: All Docker Compose stacks, `m3tal` CLI.

---

Name: `PUID`
Description: The User ID (UID) under which Docker containers should run processes. Setting this correctly ensures that containers have the appropriate permissions to read and write to host-mounted volumes.
Default: `1000`
Example: `1001`
Used by: All Docker Compose stacks.

---

Name: `PGID`
Description: The Group ID (GID) under which Docker containers should run processes. Similar to `PUID`, this ensures correct file permissions for mounted volumes.
Default: `1000`
Example: `1001`
Used by: All Docker Compose stacks.

---

Name: `TZ`
Description: Sets the timezone for all Docker containers within the M3TAL ecosystem. This ensures that logs, timestamps, and applications inside containers report accurate local times.
Default: `America/Denver`
Example: `Europe/London`
Used by: All Docker Compose stacks.

---

Name: `DEBUG_MODE`
Description: A boolean flag to enable or disable debug logging and potentially other debug-specific features across M3TAL components. Setting this to `true` can provide more diagnostic information.
Default: `false`
Example: `true`
Used by: `m3tal-api.service`, `m3tal-dashboard` container, potentially other services.

---

Name: `METRICS_ENABLED`
Description: A boolean flag to enable or disable internal metrics collection for M3TAL components. When enabled, components may expose metrics endpoints for monitoring tools.
Default: `true`
Example: `false`
Used by: `m3tal-api.service`, `m3tal-dashboard` container, potentially other services.

### Authentication

These variables control authentication and security aspects for the M3TAL Dashboard and API.

---

Name: `DASHBOARD_SECRET`
Description: A cryptographic secret key used by the M3TAL Dashboard for secure session management and encryption of sensitive data.
Default: `change_me_immediately`
Example: `your_super_secret_random_string_here_12345`
Used by: `m3tal-dashboard` container.
Notes: This variable is **auto-generated on the first `m3tal init`**. Users should **NOT set this manually** unless explicitly performing a secret rotation and understanding the implications (e.g., invalidating existing sessions).

---

Name: `API_TOKEN`
Description: An authentication token used by the `m3tal` CLI and the M3TAL Dashboard to authenticate requests with the M3TAL API daemon.
Default: `change_me_api_token`
Example: `another_auto_generated_token_for_api_access_67890`
Used by: `m3tal-api.service`, `m3tal` CLI, `m3tal-dashboard` container.
Notes: This variable is **auto-generated on the first `m3tal init`**. Users should **NOT set this manually** unless explicitly performing a token rotation.

---

Name: `ADMIN_PASSWORD`
Description: The default password for the initial 'admin' user account in the M3TAL Dashboard. This can be changed later via `m3tal dashpass`.
Default: `admin_pass`
Example: `MyNewSecurePassword!23`
Used by: `m3tal-dashboard` container (populates `/docker/users.json` via a volume mount).

### Network Configuration

These variables govern network-related settings, particularly for internal routing and host identification.

---

Name: `LOCAL_IP`
Description: The local IP address of the M3TAL host machine. This is used by services (like Traefik routing to the API daemon) that need to explicitly target the host's IP rather than relying on `host.docker.internal`.
Default: `127.0.0.1`
Example: `192.168.1.100`
Used by: Traefik gateway (`routing-compose.yml` for API routing), `m3tal` CLI.

### Storage Paths

These variables define the base and specific subdirectories for M3TAL's persistent storage, facilitating data management and portability.

---

Name: `BASE_STORAGE_PATH`
Description: The primary directory on the host filesystem where M3TAL stores all persistent data, including media, configurations, and downloads. All other `*_PATH` variables are typically subdirectories of this base path.
Default: `./data`
Example: `/mnt/m3tal_data`
Used by: All Docker Compose stacks (as volume mounts).
Notes: **In production deployments, this variable typically defaults to `/mnt`**, ensuring data persistence on a dedicated mount point, not the current working directory.

---

Name: `MEDIA_PATH`
Description: A subdirectory within `BASE_STORAGE_PATH` designated for media files (e.g., videos, music, photos) managed by M3TAL-integrated applications.
Default: `./data/media`
Example: `/mnt/media`
Used by: User-defined Docker Compose stacks (as volume mounts).

---

Name: `CONFIG_PATH`
Description: A subdirectory within `BASE_STORAGE_PATH` for storing configuration files and application-specific settings. This path is used to volume mount the host's M3TAL configuration data into containers, for example, for the Dashboard's `/docker/state` directory.
Default: `./data/config`
Example: `/mnt/config`
Used by: `m3tal-dashboard` container, user-defined Docker Compose stacks (as volume mounts).

---

Name: `DOWNLOADS_PATH`
Description: A subdirectory within `BASE_STORAGE_PATH` for storing downloaded files managed by M3TAL-integrated applications.
Default: `./data/downloads`
Example: `/mnt/downloads`
Used by: User-defined Docker Compose stacks (as volume mounts).

### Traefik and Routing

These variables configure how M3TAL services are exposed, particularly when using the Traefik reverse proxy.

---

Name: `DASHBOARD_PORT`
Description: The internal port on which the M3TAL Dashboard container listens. When `DASHBOARD_EXPOSE_MODE` is set to `local`, this port is directly exposed on the host.
Default: `8082`
Example: `8083`
Used by: `m3tal-dashboard` container, `m3tal` CLI (for `local` mode port binding).

---

Name: `DASHBOARD_EXPOSE_MODE`
Description: Controls how the M3TAL Dashboard is made accessible.
  *   `local`: The dashboard is exposed directly on the host's `DASHBOARD_PORT` (e.g., `http://HOST_IP:8082`). No Traefik required. Ideal for LAN-only or initial setups.
  *   `traefik`: The dashboard is routed through Traefik, accessible via `http://dash.${DOMAIN}`. Requires Traefik to be running via `m3tal up`.
Default: `local`
Example: `traefik`
Used by: `m3tal` CLI (specifically `m3tal dash up` to select compose override), `m3tal-dashboard` container (indirectly via Traefik labels in `m3tal-compose.traefik.yml`).

---

Name: `DOMAIN`
Description: The primary domain name used for M3TAL services when Traefik is enabled. Setting this enables routing rules like `dash.${DOMAIN}` for the Dashboard and `api.${DOMAIN}` for the M3TAL API daemon.
Default: `localhost`
Example: `myhomelab.com`
Used by: Traefik gateway (`routing-compose.yml`), `m3tal` CLI (for Traefik configuration generation).

---

Name: `TRAEFIK_WEB_PORT`
Description: The host port that Traefik's HTTP entry point binds to. This is where external HTTP traffic enters the Traefik proxy.
Default: `80`
Example: `8080`
Used by: `traefik` container.

---

Name: `TRAEFIK_WEBHTTPS_PORT`
Description: The host port that Traefik's HTTPS entry point binds to. This is where external HTTPS traffic enters the Traefik proxy.
Default: `443`
Example: `8443`
Used by: `traefik` container.

---

Name: `TRAEFIK_DASHBOARD_PORT`
Description: The internal port on which the Traefik management dashboard listens. This is explicitly bound to `127.0.0.1` on the host by default, making it accessible only from the host machine (e.g., `http://127.0.0.1:8081` mapping to container port 8080).
Default: `8080`
Example: `8080` (fixed internal port mapping on host `127.0.0.1:8081`)
Used by: `traefik` container.

### VPN Configuration

These variables are specifically for configuring optional VPN services managed by M3TAL.

---

Name: `VPN_USER`
Description: The username to be used for VPN services deployed and managed by M3TAL (e.g., WireGuard configurations).
Default: `user`
Example: `m3talvpnuser`
Used by: VPN-related Docker Compose stacks (if deployed).

---

Name: `VPN_PASSWORD`
Description: The password associated with the `VPN_USER` for VPN services managed by M3TAL.
Default: `password`
Example: `MyStrongVPNPass`
Used by: VPN-related Docker Compose stacks (if deployed).