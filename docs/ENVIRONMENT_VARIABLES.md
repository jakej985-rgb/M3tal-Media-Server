# M3TAL Environment Variables Reference

This document provides a comprehensive reference for all environment variables used by the M3TAL ecosystem. All M3TAL components, including the CLI and Dockerized services, read their configuration from `/etc/m3tal/.env`. This file is automatically managed by the `m3tal config wizard` and can be manually updated using `m3tal config set KEY value`.

## Quick Reference Table

| Variable Name             | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| **Core**                  |                                                                    |                     |                                             |                                                                             |
| `LOG_LEVEL`               | Sets the logging verbosity.                                        | `info`              | `debug`                                     | CLI, API daemon, Dashboard                                                  |
| `DEBUG_MODE`              | Enables debug mode for enhanced logging and diagnostics.           | `false`             | `true`                                      | API daemon, Dashboard                                                       |
| `METRICS_ENABLED`         | Enables Prometheus metrics collection for monitoring.              | `true`              | `false`                                     | API daemon                                                                  |
| `STATE_DIR`               | Directory for the M3TAL state database and configuration.          | `./state`           | `/var/lib/m3tal/state`                      | API daemon, Dashboard                                                       |
| `CONFIG_PATH`             | Base path for M3TAL configuration files.                           | `./data/config`     | `/mnt/config/m3tal`                         | Dashboard                                                                   |
| `PUID`                    | User ID for running Docker containers.                             | `1000`              | `1001`                                      | Dashboard                                                                   |
| `PGID`                    | Group ID for running Docker containers.                            | `1000`              | `1001`                                      | Dashboard                                                                   |
| `TZ`                      | Timezone for all components.                                       | `America/Denver`    | `UTC`                                       | API daemon, Dashboard                                                       |
| **Authentication**        |                                                                    |                     |                                             |                                                                             |
| `DASHBOARD_SECRET`        | Secret key for session security in the dashboard.                  | `change_me_immediately` | `s3cr3tK3yF0rD4shB04rd`                  | Dashboard                                                                   |
| `API_TOKEN`               | API token for authentication with the M3TAL API.                   | `change_me_api_token` | `p@$$w0rdF0rAP1`                            | CLI, Dashboard                                                              |
| `ADMIN_PASSWORD`          | Password for the default administrator user.                       | `admin_pass`        | `MyS3cur3P@ssw0rd!`                         | Dashboard                                                                   |
| **Network**               |                                                                    |                     |                                             |                                                                             |
| `HTTP_PORT`               | Port the M3TAL API daemon listens on.                              | `8080`              | `5050`                                      | API daemon                                                                  |
| `NETWORK_NAME`            | Name of the Docker network M3TAL services will join.               | `m3tal`             | `m3tal-network`                             | All services                                                                |
| `LOCAL_IP`                | The local IP address of the host machine.                          | `127.0.0.1`         | `192.168.1.100`                             | API daemon (for `host.docker.internal`)                                     |
| `DOMAIN`                  | The primary domain used for Traefik routing.                       | `localhost`         | `example.com`                               | Traefik, API daemon (for routing)                                           |
| **Storage**               |                                                                    |                     |                                             |                                                                             |
| `BASE_STORAGE_PATH`       | Base directory for M3TAL data storage.                             | `./data`            | `/mnt` (Production)                         | All services (for volumes)                                                  |
| `MEDIA_PATH`              | Path for storing media files within `BASE_STORAGE_PATH`.           | `./data/media`      | `/mnt/storage/media`                        | All services                                                                |
| `DOWNLOADS_PATH`          | Path for storing downloaded files within `BASE_STORAGE_PATH`.      | `./data/downloads`  | `/mnt/storage/downloads`                    | All services                                                                |
| **Traefik**               |                                                                    |                     |                                             |                                                                             |
| `DASHBOARD_PORT`          | Port the M3TAL dashboard container listens on internally.          | `8082`              | `8082`                                      | Dashboard, Traefik                                                          |
| `DASHBOARD_EXPOSE_MODE`   | Controls how the dashboard is exposed (`local` or `traefik`).      | `local`             | `traefik`                                   | Dashboard compose override logic                                            |
| `TRAEFIK_WEB_PORT`        | The host port Traefik listens on for HTTP traffic.                 | `80`                | `80`                                        | Traefik                                                                     |
| `TRAEFIK_WEBHTTPS_PORT`   | The host port Traefik listens on for HTTPS traffic.                | `443`               | `443`                                       | Traefik                                                                     |
| `TRAEFIK_DASHBOARD_PORT`  | The host port Traefik exposes its own dashboard on (internal).     | `8080`              | `8081`                                      | Traefik                                                                     |
| **VPN**                   |                                                                    |                     |                                             |                                                                             |
| `VPN_USER`                | Username for VPN authentication.                                   | `user`              | `myvpnuser`                                 | API daemon (if VPN is configured)                                           |
| `VPN_PASSWORD`            | Password for VPN authentication.                                   | `password`          | `mysecretvpnpassword`                       | API daemon (if VPN is configured)                                           |

---

## Detailed Environment Variable Reference

All environment variables are read from `/etc/m3tal/.env` by both the M3TAL CLI and all Docker Compose stacks. The `m3tal config wizard` is the recommended way to manage these variables.

### Core

| Variable Name | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| `LOG_LEVEL`   | Sets the logging verbosity for M3TAL components.                   | `info`              | `debug`                                     | CLI, API daemon, Dashboard                                                  |
| `DEBUG_MODE`  | Enables debug mode for enhanced logging and diagnostics.           | `false`             | `true`                                      | API daemon, Dashboard                                                       |
| `METRICS_ENABLED` | Enables Prometheus metrics collection for monitoring purposes.   | `true`              | `false`                                     | API daemon                                                                  |
| `STATE_DIR`   | Specifies the directory where the M3TAL state database (`state.db`) and other core configuration files are stored. | `./state`           | `/var/lib/m3tal/state`                      | API daemon, Dashboard                                                       |
| `CONFIG_PATH` | Defines the base path for M3TAL-related configuration files, often used by services to mount their configs. | `./data/config`     | `/mnt/config/m3tal`                         | Dashboard                                                                   |
| `PUID`        | The User ID (UID) to use when running Docker containers. Essential for file permissions. | `1000`              | `1001`                                      | Dashboard                                                                   |
| `PGID`        | The Group ID (GID) to use when running Docker containers. Essential for file permissions. | `1000`              | `1001`                                      | Dashboard                                                                   |
| `TZ`          | The timezone to be used by all M3TAL components.                 | `America/Denver`    | `UTC`                                       | API daemon, Dashboard                                                       |

### Authentication

**Note:** `DASHBOARD_SECRET` and `API_TOKEN` are auto-generated on the first `m3tal init`. Users should generally **not** set these manually unless they need to rotate them.

| Variable Name | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| `DASHBOARD_SECRET` | A secret key used for securing sessions and sensitive data within the M3TAL Dashboard. | `change_me_immediately` | `s3cr3tK3yF0rD4shB04rd`                  | Dashboard                                                                   |
| `API_TOKEN`   | An API token used for authenticating requests to the M3TAL API.    | `change_me_api_token` | `p@$$w0rdF0rAP1`                            | CLI, Dashboard                                                              |
| `ADMIN_PASSWORD` | The password for the default administrator user of the M3TAL Dashboard. | `admin_pass`        | `MyS3cur3P@ssw0rd!`                         | Dashboard                                                                   |

### Network

| Variable Name | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| `HTTP_PORT`   | The port on the host machine that the M3TAL API daemon listens on. | `8080`              | `5050`                                      | API daemon                                                                  |
| `NETWORK_NAME` | The name of the Docker network that M3TAL services will join for inter-service communication. | `m3tal`             | `m3tal-network`                             | All services                                                                |
| `LOCAL_IP`    | The local IP address of the host machine. This is used by containers to refer to the host via `host.docker.internal`. | `127.0.0.1`         | `192.168.1.100`                             | API daemon (for `host.docker.internal`)                                     |
| `DOMAIN`      | The primary domain name used for Traefik routing. Setting this enables routes like `dash.DOMAIN` and `api.DOMAIN`. | `localhost`         | `example.com`                               | Traefik, API daemon (for routing)                                           |

### Storage

**Note:** `BASE_STORAGE_PATH` controls where media data is stored. In production deployments, this defaults to `/mnt` for better separation from the OS filesystem, rather than `./data` as seen in some template examples.

| Variable Name | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| `BASE_STORAGE_PATH` | The base directory on the host system where M3TAL will store its persistent data, including media, configuration, and downloads. | `./data`            | `/mnt` (Production)                         | All services (for volumes)                                                  |
| `MEDIA_PATH`  | The specific path within `BASE_STORAGE_PATH` used for storing media files. | `./data/media`      | `/mnt/storage/media`                        | All services                                                                |
| `DOWNLOADS_PATH` | The specific path within `BASE_STORAGE_PATH` used for storing downloaded files. | `./data/downloads`  | `/mnt/storage/downloads`                    | All services                                                                |

### Traefik

| Variable Name | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| `DASHBOARD_PORT` | The internal port the M3TAL Dashboard container listens on.      | `8082`              | `8082`                                      | Dashboard, Traefik                                                          |
| `DASHBOARD_EXPOSE_MODE` | Controls how the dashboard is exposed. Use `local` for direct access (e.g., `http://localhost:8082`) and `traefik` for access via Traefik (e.g., `http://dash.DOMAIN`). | `local`             | `traefik`                                   | Dashboard compose override logic                                            |
| `TRAEFIK_WEB_PORT` | The host port that Traefik will use as its HTTP entry point.     | `80`                | `80`                                        | Traefik                                                                     |
| `TRAEFIK_WEBHTTPS_PORT` | The host port that Traefik will use as its HTTPS entry point.    | `443`               | `443`                                       | Traefik                                                                     |
| `TRAEFIK_DASHBOARD_PORT` | The internal host port Traefik uses to expose its own dashboard. Access is typically `127.0.0.1:TRAEFIK_DASHBOARD_PORT`. | `8080`              | `8081`                                      | Traefik                                                                     |

### VPN

| Variable Name | Description                                                        | Default Value       | Example Value                               | Used By                                                                     |
| :------------ | :----------------------------------------------------------------- | :------------------ | :------------------------------------------ | :-------------------------------------------------------------------------- |
| `VPN_USER`    | The username for establishing a VPN connection.                    | `user`              | `myvpnuser`                                 | API daemon (if VPN is configured and utilized)                              |
| `VPN_PASSWORD` | The password for establishing a VPN connection.                    | `password`          | `mysecretvpnpassword`                       | API daemon (if VPN is configured and utilized)                              |