# M3TAL Environment Variables Reference

All M3TAL configuration is managed via environment variables, primarily sourced from the `/etc/m3tal/.env` file. This file is automatically managed by the `m3tal config wizard` and can be updated using `m3tal config set KEY value`. Both the M3TAL CLI and all Docker Compose stacks read these variables via the `--env-file` flag.

This document provides a comprehensive reference for all available environment variables, categorized for clarity.

---

### Quick Reference Table

| Variable Name          | Description                                                                                                                                                                                                           | Default Value        | Example Value                                       | Used By                                                                      |
| :--------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------- | :-------------------------------------------------- | :--------------------------------------------------------------------------- |
| **Core**               |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `DASHBOARD_PORT`       | The port on which the M3TAL dashboard listens.                                                                                                                                                                        | `8082`               | `8082`                                              | `m3tal-dashboard` container                                                  |
| `HTTP_PORT`            | The port on which the M3TAL API daemon listens.                                                                                                                                                                       | `8080`               | `8080`                                              | `m3tal-api.service`                                                          |
| `STATE_DIR`            | The directory where M3TAL stores its state, including the SQLite database and user configuration. **Note:** In production, this often maps to a persistent volume like `${CONFIG_PATH}/m3tal/state`.                 | `./state`            | `/mnt/config/m3tal/state`                           | `m3tal-api.service`, `m3tal-dashboard` container                             |
| `LOG_LEVEL`            | The logging verbosity level for M3TAL components.                                                                                                                                                                     | `info`               | `debug`                                             | `m3tal-api.service`, `m3tal-dashboard` container                             |
| `DEBUG_MODE`           | Enables or disables debug mode for M3TAL components.                                                                                                                                                                  | `false`              | `true`                                              | `m3tal-api.service`, `m3tal-dashboard` container                             |
| `METRICS_ENABLED`      | Enables or disables the Prometheus metrics endpoint for the API.                                                                                                                                                      | `true`               | `false`                                             | `m3tal-api.service`                                                          |
| **Auth**               |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `DASHBOARD_SECRET`     | A secret key used for securing the dashboard session. **Auto-generated on first `m3tal init`. Do NOT set manually unless rotating.**                                                                                   | `change_me_immediately` | `a-very-secret-and-long-string-for-rotation`          | `m3tal-dashboard` container                                                  |
| `API_TOKEN`            | A token used for authenticating API requests. **Auto-generated on first `m3tal init`. Do NOT set manually unless rotating.**                                                                                            | `change_me_api_token` | `another-long-and-secure-api-token-for-rotation`    | `m3tal-api.service` (used internally for some operations)                    |
| `ADMIN_PASSWORD`       | The password for the default administrator user.                                                                                                                                                                      | `admin_pass`         | `a-strong-and-unique-password`                      | `m3tal-dashboard` container (for initial login), `m3tal-api.service`         |
| **Network**            |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `NETWORK_NAME`         | The name of the Docker network M3TAL services will use.                                                                                                                                                               | `m3tal`              | `m3tal_net`                                         | Docker Compose stacks                                                        |
| `LOCAL_IP`             | The IP address of the host machine that M3TAL services can bind to or communicate with.                                                                                                                               | `127.0.0.1`          | `192.168.1.100`                                     | `m3tal-api.service` (for `host.docker.internal` resolution in some contexts) |
| `DOMAIN`               | The primary domain name for M3TAL. Setting this enables `dash.DOMAIN` and `api.DOMAIN` routing via Traefik.                                                                                                           | `localhost`          | `m3tal.mydomain.com`                                | Traefik gateway, `m3tal-api.service`                                         |
| **Storage**            |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `BASE_STORAGE_PATH`    | The base directory for storing M3TAL data, including media, configuration, and downloads. **Defaults to `/mnt` in production deployments.**                                                                          | `./data`             | `/mnt`                                              | `m3tal-dashboard` container, all Docker Compose stacks                       |
| `MEDIA_PATH`           | The specific directory within `BASE_STORAGE_PATH` for storing media files.                                                                                                                                            | `./data/media`       | `/mnt/m3tal_storage/media`                          | `m3tal-dashboard` container, all Docker Compose stacks                       |
| `CONFIG_PATH`          | The specific directory within `BASE_STORAGE_PATH` for storing configuration files.                                                                                                                                      | `./data/config`      | `/mnt/m3tal_storage/config`                         | `m3tal-dashboard` container, all Docker Compose stacks                       |
| `DOWNLOADS_PATH`       | The specific directory within `BASE_STORAGE_PATH` for storing downloaded files.                                                                                                                                       | `./data/downloads`   | `/mnt/m3tal_storage/downloads`                      | `m3tal-dashboard` container, all Docker Compose stacks                       |
| **User/Group IDs**     |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `PUID`                 | The User ID to run containers with.                                                                                                                                                                                   | `1000`               | `1000`                                              | `m3tal-dashboard` container                                                  |
| `PGID`                 | The Group ID to run containers with.                                                                                                                                                                                  | `1000`               | `1000`                                              | `m3tal-dashboard` container                                                  |
| **Timezone**           |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `TZ`                   | The timezone for M3TAL components.                                                                                                                                                                                    | `America/Denver`    | `UTC`                                               | `m3tal-dashboard` container                                                  |
| **Traefik**            |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `DASHBOARD_EXPOSE_MODE`| Controls how the dashboard is exposed. `local` (default) exposes it directly via `DASHBOARD_PORT`. `traefik` exposes it via `dash.DOMAIN` through the Traefik gateway.                                            | `local`              | `traefik`                                           | `m3tal-dashboard` container (via compose override), Traefik configuration    |
| `TRAEFIK_WEB_PORT`     | The port Traefik listens on for HTTP traffic.                                                                                                                                                                         | `80`                 | `80`                                                | Traefik container                                                            |
| `TRAEFIK_WEBHTTPS_PORT`| The port Traefik listens on for HTTPS traffic (if configured).                                                                                                                                                        | `443`                | `443`                                               | Traefik container                                                            |
| `TRAEFIK_DASHBOARD_PORT`| The port Traefik's own dashboard is exposed on the host (if enabled).                                                                                                                                                 | `8080`               | `8081`                                              | Traefik container                                                            |
| **VPN**                |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |
| `VPN_USER`             | The username for the VPN connection.                                                                                                                                                                                  | `user`               | `myvpnuser`                                         | Not directly used by core M3TAL services, but can be used by custom stacks   |
| `VPN_PASSWORD`         | The password for the VPN connection.                                                                                                                                                                                  | `password`           | `mysecurevpnpassword`                               | Not directly used by core M3TAL services, but can be used by custom stacks   |
| **System**             |                                                                                                                                                                                                                       |                      |                                                     |                                                                              |

---

## Core Variables

These variables control fundamental aspects of M3TAL's operation.

### `DASHBOARD_PORT`

*   **Description:** The port on which the M3TAL dashboard listens for incoming connections.
*   **Default Value:** `8082`
*   **Example Value:** `8082`
*   **Used By:** `m3tal-dashboard` container.

### `HTTP_PORT`

*   **Description:** The port on which the M3TAL API daemon (Go binary) listens for incoming connections.
*   **Default Value:** `8080`
*   **Example Value:** `8080`
*   **Used By:** `m3tal-api.service`.

### `STATE_DIR`

*   **Description:** The directory where M3TAL stores its state, including the SQLite database and user configuration. In production deployments, this variable is often mapped to a persistent volume like `${CONFIG_PATH}/m3tal/state` to ensure data is not lost.
*   **Default Value:** `./state`
*   **Example Value:** `/mnt/config/m3tal/state`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `LOG_LEVEL`

*   **Description:** Controls the verbosity of logging output for M3TAL components. Accepted values typically include `debug`, `info`, `warn`, `error`.
*   **Default Value:** `info`
*   **Example Value:** `debug`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `DEBUG_MODE`

*   **Description:** Enables or disables debug mode for M3TAL components, which may provide more detailed logging or enable additional debugging features.
*   **Default Value:** `false`
*   **Example Value:** `true`
*   **Used By:** `m3tal-api.service`, `m3tal-dashboard` container.

### `METRICS_ENABLED`

*   **Description:** Enables or disables the Prometheus metrics endpoint for the M3TAL API service. When enabled, metrics can be scraped by a Prometheus server for monitoring.
*   **Default Value:** `true`
*   **Example Value:** `false`
*   **Used By:** `m3tal-api.service`.

---

## Auth Variables

These variables are crucial for securing your M3TAL instance.

### `DASHBOARD_SECRET`

*   **Description:** A secret key used for securing the dashboard's session cookies and other critical operations. **This variable is auto-generated on the first `m3tal init` command. Users should NOT set this manually unless performing a secret rotation.**
*   **Default Value:** `change_me_immediately`
*   **Example Value:** `a-very-secret-and-long-string-for-rotation`
*   **Used By:** `m3tal-dashboard` container.

### `API_TOKEN`

*   **Description:** A token used for authenticating programmatic access to the M3TAL API. **This variable is auto-generated on the first `m3tal init` command. Users should NOT set this manually unless performing a token rotation.**
*   **Default Value:** `change_me_api_token`
*   **Example Value:** `another-long-and-secure-api-token-for-rotation`
*   **Used By:** `m3tal-api.service` (used internally for some operations).

### `ADMIN_PASSWORD`

*   **Description:** The password for the default administrator user account. This is used for the initial login to the M3TAL dashboard. It's highly recommended to change this immediately after initial setup.
*   **Default Value:** `admin_pass`
*   **Example Value:** `a-strong-and-unique-password`
*   **Used By:** `m3tal-dashboard` container (for initial login), `m3tal-api.service` (for authentication).

---

## Network Variables

These variables configure M3TAL's network interfaces and external access.

### `NETWORK_NAME`

*   **Description:** The name of the Docker network that M3TAL services will be connected to. This ensures proper communication between containers.
*   **Default Value:** `m3tal`
*   **Example Value:** `m3tal_net`
*   **Used By:** Docker Compose stacks.

### `LOCAL_IP`

*   **Description:** The IP address of the host machine that M3TAL services can bind to or communicate with. This is particularly relevant for internal communication where `host.docker.internal` might be used.
*   **Default Value:** `127.0.0.1`
*   **Example Value:** `192.168.1.100`
*   **Used By:** `m3tal-api.service` (for `host.docker.internal` resolution in some contexts).

### `DOMAIN`

*   **Description:** The primary domain name for your M3TAL instance. Setting this variable is essential for enabling Traefik routing rules. When set, M3TAL will configure Traefik to route traffic to `api.DOMAIN` (for the API daemon) and `dash.DOMAIN` (for the dashboard).
*   **Default Value:** `localhost`
*   **Example Value:** `m3tal.mydomain.com`
*   **Used By:** Traefik gateway, `m3tal-api.service`.

---

## Storage Variables

These variables define where M3TAL stores its persistent data.

### `BASE_STORAGE_PATH`

*   **Description:** The base directory on the host machine where M3TAL will store all its persistent data, including media files, configuration, and downloads. **In production deployments, this defaults to `/mnt` to ensure data is stored on a dedicated volume and not within the container's ephemeral filesystem. It does NOT default to `./data` as seen in template configurations.**
*   **Default Value:** `./data`
*   **Example Value:** `/mnt`
*   **Used By:** `m3tal-dashboard` container, all Docker Compose stacks.

### `MEDIA_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` where media files (e.g., uploaded images, videos) are stored.
*   **Default Value:** `./data/media`
*   **Example Value:** `/mnt/m3tal_storage/media`
*   **Used By:** `m3tal-dashboard` container, all Docker Compose stacks.

### `CONFIG_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files. This can include database files, user credentials (e.g., `users.json`), and other service-specific configurations.
*   **Default Value:** `./data/config`
*   **Example Value:** `/mnt/m3tal_storage/config`
*   **Used By:** `m3tal-dashboard` container, all Docker Compose stacks.

### `DOWNLOADS_PATH`

*   **Description:** The specific subdirectory within `BASE_STORAGE_PATH` where downloaded files are stored.
*   **Default Value:** `./data/downloads`
*   **Example Value:** `/mnt/m3tal_storage/downloads`
*   **Used By:** `m3tal-dashboard` container, all Docker Compose stacks.

---

## User/Group ID Variables

These variables control the User ID (UID) and Group ID (GID) that containers will run as. This is important for file permissions on the host system.

### `PUID`

*   **Description:** The User ID (UID) that Docker containers will use to run processes. This should typically match the UID of the user who owns the data directories on the host system.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** `m3tal-dashboard` container.

### `PGID`

*   **Description:** The Group ID (GID) that Docker containers will use to run processes. This should typically match the GID of the group that owns the data directories on the host system.
*   **Default Value:** `1000`
*   **Example Value:** `1000`
*   **Used By:** `m3tal-dashboard` container.

---

## Timezone Variable

### `TZ`

*   **Description:** The timezone setting for M3TAL components. This ensures that timestamps in logs and other time-sensitive operations are accurate.
*   **Default Value:** `America/Denver`
*   **Example Value:** `UTC`
*   **Used By:** `m3tal-dashboard` container.

---

## Traefik Variables

These variables are specific to the Traefik reverse proxy configuration.

### `DASHBOARD_EXPOSE_MODE`

*   **Description:** Controls how the M3TAL dashboard is exposed.
    *   `local` (default): The dashboard is exposed directly via the port defined by `DASHBOARD_PORT` (e.g., `http://HOST_IP:8082`). No Traefik configuration is needed for dashboard access. This is ideal for LAN-only setups or initial testing.
    *   `traefik`: The dashboard is exposed through the Traefik reverse proxy, accessible via `http://dash.DOMAIN`. This requires Traefik to be running and configured with your `DOMAIN`.
*   **Default Value:** `local`
*   **Example Value:** `traefik`
*   **Used By:** `m3tal-dashboard` container (via compose override), Traefik configuration.

### `TRAEFIK_WEB_PORT`

*   **Description:** The port on the host machine that Traefik listens on for incoming HTTP traffic. This is the main entry point for web requests when using Traefik.
*   **Default Value:** `80`
*   **Example Value:** `80`
*   **Used By:** Traefik container.

### `TRAEFIK_WEBHTTPS_PORT`

*   **Description:** The port on the host machine that Traefik listens on for incoming HTTPS traffic. This is used when you have SSL/TLS certificates configured for Traefik.
*   **Default Value:** `443`
*   **Example Value:** `443`
*   **Used By:** Traefik container.

### `TRAEFIK_DASHBOARD_PORT`

*   **Description:** The port on the host machine that Traefik's own administrative dashboard is exposed on. This is typically used for monitoring and managing Traefik itself.
*   **Default Value:** `8080`
*   **Example Value:** `8081`
*   **Used By:** Traefik container.

---

## VPN Variables

These variables are intended for use with custom Docker stacks or configurations that require a VPN connection.

### `VPN_USER`

*   **Description:** The username for establishing a VPN connection. This is not directly used by core M3TAL services but can be leveraged by external Docker stacks or custom configurations.
*   **Default Value:** `user`
*   **Example Value:** `myvpnuser`
*   **Used By:** Not directly used by core M3TAL services, but can be used by custom stacks.

### `VPN_PASSWORD`

*   **Description:** The password for establishing a VPN connection. Similar to `VPN_USER`, this is not directly used by core M3TAL services but is available for use in custom configurations.
*   **Default Value:** `password`
*   **Example Value:** `mysecurevpnpassword`
*   **Used By:** Not directly used by core M3TAL services, but can be used by custom stacks.