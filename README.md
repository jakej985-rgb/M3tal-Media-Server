# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

**Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.**

## APT Installation

To install the M3TAL CLI binary and API daemon, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The M3TAL system relies on a specific filesystem structure for its operation and configuration:

| Path                        | Purpose                                                              |
|-----------------------------|----------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file, managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database, auto-created by the M3TAL API daemon.         |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configuration. |
| `/docker`                   | Symlink pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | M3TAL Dashboard credential store, managed by `m3tal dashpass`.       |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

It is crucial to understand that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside. The `/docker` directory is a user-facing symlink alias provided for convenience in all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. This can be done using `m3tal config wizard` or `m3tal config set KEY value`.
3. Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`) - Default

This is the default mode, ideal for initial setup and local network access.
*   **Mechanism**: A direct port binding is added to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`), bypassing Traefik.
*   **Access**: You can access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Clarification**: A new user performing a default installation will access the dashboard directly via port 8082. This direct access behavior is a result of the default `DASHBOARD_EXPOSE_MODE=local` setting.
*   **Requirements**: No Traefik gateway is required to be running for dashboard access in this mode.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

This mode integrates the dashboard with the Traefik reverse proxy for domain-based access.
*   **Mechanism**: Traefik labels are applied to the dashboard container, allowing Traefik to route traffic to it.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Requirements**: The Traefik gateway must be running (`m3tal up` typically starts it) and configured with the appropriate domain.

## Traefik Gateway

Traefik functions as the central reverse proxy within the M3TAL system, automatically discovering and routing traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik dynamically loads its configuration from `/docker/dynamic/` (which symlinks to `/opt/m3tal/stack/dynamic/`). This allows for hot-reloading of routing rules. For instance, the M3TAL API daemon, which runs on the host-local port `8080`, is exposed via Traefik using a dynamic configuration file:

*   **API Daemon Routing**: Requests to `api.DOMAIN` (e.g., `api.example.com`) are routed by Traefik through the dynamic configuration file `/docker/dynamic/api.yml` to the M3TAL Go API daemon listening on `http://host.docker.internal:8080`.

*   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE=traefik`, requests to `dash.DOMAIN` are routed to the M3TAL Dashboard container, which internally listens on port 8082. This routing is managed via Traefik labels applied to the dashboard service within its Docker Compose override file.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service through Traefik, you need to add specific labels to its service definition. For example, to expose a hypothetical `my-app` service via `app.DOMAIN`:

```yaml
# my-app-compose.yml (placed in /docker/)
services:
  my-app:
    image: nginx:alpine
    container_name: my-app-nginx
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
      - "traefik.docker.network=proxy" # Ensure it's on the proxy network
    networks:
      - proxy

networks:
  proxy:
    external: true # Assumes 'proxy' network is created by Traefik stack
```

After adding this file to `/docker/` and running `m3tal up`, Traefik will automatically detect and route `http://app.DOMAIN` to your `my-app` service.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage its lifecycle using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get started with the M3TAL Dashboard:

1.  **Start the M3TAL Dashboard specifically**:
    ```bash
    m3tal dash up
    ```
    This command downloads the necessary Docker Compose files for the dashboard and starts it using the configuration specified by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`. By default (`DASHBOARD_EXPOSE_MODE=local`), you will be able to access it via `http://HOST_IP:8082`.

2.  **Deploy all M3TAL-managed and user-defined stacks**:
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all Docker Compose stacks defined by `*-compose.yml` files located within the `/docker/` directory. This includes core M3TAL components like the Traefik gateway, as well as any user-defined compose files you have placed there.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|-----------------|-------------|
| 80   | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.