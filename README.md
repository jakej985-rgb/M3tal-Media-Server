# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Filesystem Contract

The following paths are central to M3TAL's operation and configuration:

| Path                     | Purpose                                                              |
|--------------------------|----------------------------------------------------------------------|
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created by the API daemon.               |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.             |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

It is crucial to understand that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/` for all stack operations, providing a convenient and consistent interface for managing Docker Compose configurations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard provides a web interface for system management. Its access method is determined by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   This is the default mode for a new installation.
*   The dashboard container's internal port `8082` is directly bound to the host's `8082` port (or `DASHBOARD_PORT` if specified).
*   Access the dashboard via `http://HOST_IP:8082` or `http://localhost:8082`.
*   This mode does not require Traefik to be running.
*   A new user performing a default M3TAL installation will access the dashboard directly via port 8082. This behavior is linked to the default setting `DASHBOARD_EXPOSE_MODE=local`.

### Traefik Mode

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   The dashboard container is exposed via Traefik, utilizing Traefik labels defined in `m3tal-compose.traefik.yml`.
*   Access the dashboard via `http://dash.DOMAIN` (e.g., `http://dash.example.com` if `DOMAIN` is `example.com`).
*   This mode requires the Traefik gateway to be running.

## Traefik Gateway

The Traefik gateway automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also uses a file provider to load dynamic configuration from `/docker/dynamic/` (symlinked from `/opt/m3tal/stack/dynamic/`). This allows routing requests to services that are not necessarily Docker containers or are listening on host-local ports. For example, `dynamic/api.yml` routes `api.DOMAIN` to the Go API daemon listening on the host's port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container's internal port 8082 when `DASHBOARD_EXPOSE_MODE=traefik` is active, enabled by the labels in `m3tal-compose.traefik.yml`.

### Exposing a Custom User Service via Traefik

To expose a hypothetical `my-app` service from your `my-app-compose.yml` via Traefik:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN}`)"
      - "traefik.http.routers.myapp.entrypoints=web"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.docker.network=proxy" # Ensure this service is on the 'proxy' network
    networks:
      - proxy

networks:
  proxy:
    external: true # Use the existing 'proxy' network managed by Traefik
```

After placing this file in `/docker/` and running `m3tal up`, your `my-app` service would be accessible at `http://app.DOMAIN` (e.g., `http://app.localhost`).

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. Standard `systemctl` commands apply:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get the M3TAL Dashboard up and running:

*   **Start the dashboard container specifically**:
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard compose files and starts the dashboard container with the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting. By default, you would then access it at `http://HOST_IP:8082`.

*   **Deploy all M3TAL-managed Docker Compose stacks**:
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all `*-compose.yml` files located within the `/docker/` directory. This includes the `routing-compose.yml` for Traefik, `m3tal-compose.yml` for the dashboard (if not already managed by `m3tal dash up`), and any other user-defined compose files you have placed in `/docker/`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                    | Access                                        | Description                                                                                                    |
|------|----------------------------|-----------------------------------------------|----------------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point   | Public                                        | The public-facing HTTP port for services exposed via Traefik.                                                  |
| 8080 | M3TAL API daemon (Go)      | Host-local                                    | The internal port the M3TAL API daemon listens on.                                                             |
| 8081 | Traefik dashboard          | Host-local only                               | The internal Traefik dashboard port, accessible only from the host machine.                                    |
| 8082 | M3TAL Dashboard            | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.