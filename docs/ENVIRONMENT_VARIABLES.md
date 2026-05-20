```markdown
# Environment Variables Reference

All M3TAL configuration is managed through environment variables defined in the central configuration file, `/etc/m3tal/.env`. This file is crucial for both the M3TAL CLI and all Docker Compose stacks.

- **CLI Usage:** The `m3tal` binary reads this file directly for its operations.
- **Docker Compose Stacks:** When running `m3tal up` or specific stack commands, Docker Compose automatically loads variables from `/etc/m3tal/.env` via the `--env-file` flag, making them available to all containers.

You can modify these variables using the `m3tal config wizard` for guided setup or `m3tal config set <KEY> <VALUE>` for direct changes.

---

## Quick Reference

| Name | Description | Default Value | Component(s) |
|---|---|---|---|
| `DASHBOARD_PORT` | Port for the M3TAL Dashboard. | `8082` | `m3tal-dashboard` |
| `DASHBOARD_EXPOSE_MODE` | Controls dashboard exposure: `local` (direct port) or `traefik` (via domain). | `local` | `m3tal-dashboard`, CLI |
| `HTTP_PORT` | Port for the M3TAL API daemon. | `8080` | `m3tal-api` |
| `STATE_DIR` | Host path for application state, including the SQLite database and user credentials. | `./state` | `m3tal-api`, `m3tal-dashboard` |
| `LOG_LEVEL` | Minimum logging level for the API daemon and CLI. | `info` | `m3tal-api`, CLI |
| `DASHBOARD_SECRET` | Secret key for M3TAL Dashboard session security. **Auto-generated.** | `change_me_immediately` | `m3tal-dashboard` |
| `API_TOKEN` | Bearer token for CLI and external access to the M3TAL API. **Auto-generated.** | `change_me_api_token` | `m3tal-api`, CLI |
| `ADMIN_PASSWORD` | Initial password for the default `admin` user in the M3TAL Dashboard. | `admin_pass` | `m3tal-dashboard` |
| `NETWORK_NAME` | Name of the Docker network shared by all M3TAL services and user stacks. | `m3tal` | All compose stacks |
| `LOCAL_IP` | Local IP address for host-level services, primarily for CLI internal use. | `127.0.0.1` | CLI |
| `DOMAIN` | Base domain for Traefik routing. Enables `dash.DOMAIN` and `api.DOMAIN` routes. | `localhost` | `m3tal-dashboard`, Traefik gateway |
| `VPN_USER` | Placeholder username for VPN services in user-defined stacks. | `user` | User-defined VPN stacks |
| `VPN_PASSWORD` | Placeholder password for VPN services in user-defined stacks. | `password` | User-defined VPN stacks |
| `BASE_STORAGE_PATH` | Primary directory for all M3TAL media and data storage on the host. | `./data` | All compose stacks |
| `MEDIA_PATH` | Subdirectory for media files, relative to `BASE_STORAGE_PATH`. | `./data/media` | User-defined media stacks |
| `CONFIG_PATH` | Subdirectory for application configurations, relative to `BASE_STORAGE_PATH`. | `./data/config` | User-defined config stacks, `m3tal-dashboard` |
| `DOWNLOADS_PATH` | Subdirectory for downloads, relative to `BASE_STORAGE_PATH`. | `./data/downloads` | User-defined download stacks |
| `PUID` | POSIX User ID for container processes to manage file permissions. | `1000` | All compose stacks |
| `PGID` | POSIX Group ID for container processes to manage file permissions. | `1000` | All compose stacks |
| `TZ` | Timezone setting for all M3TAL containers. | `America/Denver` | All compose stacks |
| `TRAEFIK_WEB_PORT` | External port for Traefik's HTTP entry point. | `80` | Traefik gateway |
| `TRAEFIK_WEBHTTPS_PORT` | External port for Traefik's HTTPS entry point. | `443` | Traefik gateway |
| `TRAEFIK_DASHBOARD_PORT` | Internal port for the Traefik management dashboard. | `8080` | Traefik gateway |
| `DEBUG_MODE` | Enables extended debug logging for the API and CLI. | `false` | `m3tal-api`, CLI |
| `METRICS_ENABLED` | Controls whether system metrics collection is active. | `true` | `m3tal-api`, CLI |

---

## Detailed Variable Reference

### Core Configuration

These variables control fundamental aspects of the M3TAL system's operation and environment.

### `PUID`

**Description:** The POSIX User ID (UID) that container processes will run as. This ensures correct file permissions for volumes mounted into containers, matching the host system's user.
**Default Value:** `1000`
**Example Value:** `1001`
**Used By:** All compose stacks (e.g., `m3tal-dashboard`, user-defined stacks)

### `PGID`

**Description:** The POSIX Group ID (GID) that container processes will run as. Similar to `PUID`, this is critical for file permission management.
**Default Value:** `1000`
**Example Value:** `1001`
**Used By:** All compose stacks (e.g., `m3tal-dashboard`, user-defined stacks)

### `TZ`

**Description:** Sets the timezone for all M3TAL containers, ensuring consistent timestamps and scheduled tasks.
**Default Value:** `America/Denver`
**Example Value:** `Europe/London`
**Used By:** All compose stacks (e.g., `m3tal-dashboard`, user-defined stacks)

### `LOG_LEVEL`

**Description:** Defines the minimum logging level for the M3TAL API daemon and CLI. Options typically include `debug`, `info`, `warn`, `error`, `fatal`.
**Default Value:** `info`
**Example Value:** `debug`
**Used By:** `m3tal-api`, CLI

### `DEBUG_MODE`

**Description:** When set to `true`, enables extended debug logging and potentially other debug-specific behaviors within the M3TAL API and CLI.
**Default Value:** `false`
**Example Value:** `true`
**Used By:** `m3tal-api`, CLI

### `METRICS_ENABLED`

**Description:** Controls whether performance metrics collection is active for the M3TAL API daemon and CLI. Disabling this can reduce resource usage slightly.
**Default Value:** `true`
**Example Value:** `false`
**Used By:** `m3tal-api`, CLI

### Authentication

These variables manage access and security for M3TAL components.

### `DASHBOARD_SECRET`

**Description:** A strong, randomly generated secret key used by the `m3tal-dashboard` container to secure user sessions.
**Default Value:** `change_me_immediately`
**Example Value:** `a_very_long_and_random_string_of_characters_12345`
**Used By:** `m3tal-dashboard`
**Note:** This variable is **auto-generated** on the first `m3tal init` command. Users should **not** set this manually unless performing a security rotation or recovery.

### `API_TOKEN`

**Description:** A bearer token used by the M3TAL CLI and other authenticated clients to interact with the M3TAL API daemon. Provides secure access to API endpoints.
**Default Value:** `change_me_api_token`
**Example Value:** `another_extremely_long_random_string_of_hex_digits_67890`
**Used By:** `m3tal-api`, CLI
**Note:** This variable is **auto-generated** on the first `m3tal init` command. Users should **not** set this manually unless performing a security rotation or recovery.

### `ADMIN_PASSWORD`

**Description:** The initial password for the default `admin` user within the M3TAL Dashboard. It is highly recommended to change this immediately after first login via the dashboard's user management interface.
**Default Value:** `admin_pass`
**Example Value:** `my_strong_admin_password!`
**Used By:** `m3tal-dashboard`

### Network Configuration

Variables governing how M3TAL components communicate and are exposed on the network.

### `DASHBOARD_PORT`

**Description:** The internal and host-exposed port used by the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE` is set to `local`.
**Default Value:** `8082`
**Example Value:** `8083`
**Used By:** `m3tal-dashboard` (when `DASHBOARD_EXPOSE_MODE=local`)

### `DASHBOARD_EXPOSE_MODE`

**Description:** Determines how the M3TAL Dashboard is exposed.
- `local`: The dashboard is directly accessible via `http://HOST_IP:${DASHBOARD_PORT}`. No Traefik required. Best for LAN-only setups or initial configuration.
- `traefik`: The dashboard is exposed via Traefik at `http://dash.${DOMAIN}`. Best for domain-based setups behind a reverse proxy.
**Default Value:** `local`
**Example Value:** `traefik`
**Used By:** `m3tal-dashboard` (determines which compose override to use), CLI (`m3tal dash up` command)

### `HTTP_PORT`

**Description:** The port on which the M3TAL API daemon listens for incoming connections. This is an internal, host-local port, typically accessed by Traefik via `http://host.docker.internal:8080`.
**Default Value:** `8080`
**Example Value:** `8081`
**Used By:** `m3tal-api`, Traefik gateway (for API routing)

### `NETWORK_NAME`

**Description:** The name of the Docker network that M3TAL uses to connect all its core services and user-deployed compose stacks. All containers needing to communicate should be on this network.
**Default Value:** `m3tal`
**Example Value:** `m3tal-internal-net`
**Used By:** All compose stacks

### `LOCAL_IP`

**Description:** Primarily for internal CLI usage to refer to the host machine's IP address. For container-to-host communication, `host.docker.internal` is generally preferred.
**Default Value:** `127.0.0.1`
**Example Value:** `192.168.1.100`
**Used By:** CLI

### Storage Paths

Variables defining the base directories and subdirectories for M3TAL's data storage on the host filesystem.

### `BASE_STORAGE_PATH`

**Description:** The root directory on the host machine where all M3TAL-related data, media, configuration, and state will be stored.
**Default Value:** `./data` (for development/template, relative to `/opt/m3tal/stack/`)
**Example Value:** `/mnt/m3tal_data`
**Used By:** All compose stacks (for volume mounts), `m3tal-dashboard` (as base for its state volume)
**Note:** In production deployments, this variable typically **defaults to `/mnt`** to leverage dedicated storage volumes. The template value of `./data` is for initial setup and local testing.

### `STATE_DIR`

**Description:** The host path where M3TAL stores its critical state files. This includes the SQLite database (`/var/lib/m3tal/state.db`) used by the API daemon and `users.json` for the dashboard credentials.
**Default Value:** `./state` (relative to the compose file location, typically mapped to `/var/lib/m3tal/state` on the host in production setups)
**Example Value:** `/var/lib/m3tal/state`
**Used By:** `m3tal-api` (for `/var/lib/m3tal/state.db`), `m3tal-dashboard` (for `users.json`)

### `MEDIA_PATH`

**Description:** A subdirectory, relative to `BASE_STORAGE_PATH`, designated for media files. User-defined media management stacks typically mount this path.
**Default Value:** `./data/media`
**Example Value:** `/mnt/media`
**Used By:** User-defined media stacks

### `CONFIG_PATH`

**Description:** A subdirectory, relative to `BASE_STORAGE_PATH`, intended for application configuration files.
**Default Value:** `./data/config`
**Example Value:** `/mnt/config`
**Used By:** User-defined config stacks, `m3tal-dashboard` (for its state directory mount)

### `DOWNLOADS_PATH`

**Description:** A subdirectory, relative to `BASE_STORAGE_PATH`, where downloaded content is typically stored.
**Default Value:** `./data/downloads`
**Example Value:** `/mnt/downloads`
**Used By:** User-defined download stacks

### Traefik Gateway

These variables configure the Traefik reverse proxy for exposing services.

### `DOMAIN`

**Description:** The base domain name for Traefik routing. When set, Traefik will automatically configure HTTP routes like `dash.${DOMAIN}` for the M3TAL Dashboard and `api.${DOMAIN}` for the M3TAL API daemon.
**Default Value:** `localhost`
**Example Value:** `mydomain.com`
**Used By:** `m3tal-dashboard` (Traefik labels when `DASHBOARD_EXPOSE_MODE=traefik`), Traefik gateway (dynamic config provider for `api.DOMAIN`)
**Note:** Setting this variable enables domain-based routing for M3TAL services.

### `TRAEFIK_WEB_PORT`

**Description:** The external port on which Traefik listens for unencrypted HTTP traffic (entry point named `web`).
**Default Value:** `80`
**Example Value:** `80`
**Used By:** Traefik gateway (`routing-compose.yml`)

### `TRAEFIK_WEBHTTPS_PORT`

**Description:** The external port on which Traefik listens for encrypted HTTPS traffic (entry point named `websecure`). Requires additional SSL configuration for Traefik and a valid certificate.
**Default Value:** `443`
**Example Value:** `443`
**Used By:** Traefik gateway (`routing-compose.yml`)

### `TRAEFIK_DASHBOARD_PORT`

**Description:** The internal port on which the Traefik management dashboard listens. On the host, this is typically mapped to `127.0.0.1:8081` to prevent public exposure.
**Default Value:** `8080`
**Example Value:** `8080`
**Used By:** Traefik gateway (`routing-compose.yml`)

### VPN Configuration

Placeholder variables for integrating VPN services or other secure tunnels. These are generally used for user-defined VPN stacks rather than M3TAL core components.

### `VPN_USER`

**Description:** A generic placeholder for a username used by VPN services or other secure tunnel configurations in user-defined stacks.
**Default Value:** `user`
**Example Value:** `m3talvpnuser`
**Used By:** User-defined VPN stacks (e.g., WireGuard, OpenVPN, not directly by `cloudflared`)

### `VPN_PASSWORD`

**Description:** A generic placeholder for a password associated with `VPN_USER` for VPN services or secure tunnel configurations in user-defined stacks.
**Default Value:** `password`
**Example Value:** `s3cur3p@ssw0rd`
**Used By:** User-defined VPN stacks (e.g., WireGuard, OpenVPN, not directly by `cloudflared`)
```