# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI and API daemon via APT:

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

## Components

The M3TAL system is composed of several integrated components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container, running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name, typically on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing secure, zero-config internet access to internal services.

## Filesystem Contract

The following table details the critical filesystem paths and their purposes within the M3TAL system:

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for Docker Compose stack files and Traefik configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Credential store for the M3TAL Dashboard. Managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, serving as the user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all deployed stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Access Method**: Direct port binding at `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override, which adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the dashboard container.
*   **Traefik Requirement**: No Traefik configuration or setup is required for this mode.
*   **New User Experience**: A new user performing a default installation will access the dashboard directly via port 8082. This behavior is linked to the default setting `DASHBOARD_EXPOSE_MODE=local`.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Access Method**: Domain routing at `http://dash.DOMAIN` via Traefik labels.
*   **Mechanism**: This mode uses the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard container. Traefik then routes requests for `dash.DOMAIN` to the dashboard container's internal port `8082`.
*   **Traefik Requirement**: Traefik must be running (typically via `m3tal up`) and properly configured to route traffic.

## Traefik Gateway

Traefik is deployed as a Docker container via `routing-compose.yml` and serves as the primary ingress point for services exposed by M3TAL. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik listens on host port `80` (for HTTP) and uses a file provider to load dynamic configuration from `/docker/dynamic/` (which supports hot-reloading).

### API Daemon Routing

For services listening on host-local ports, such as the M3TAL Go API daemon on `8080`, Traefik utilizes dynamic configuration files to route external requests. For example, `dynamic/api.yml` routes `api.DOMAIN` to `http://host.docker.internal:8080` (the Go API daemon) through Traefik's load balancer.

Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, requests for `dash.DOMAIN` are routed to the dashboard container on port `8082` based on labels in `m3tal-compose.traefik.yml`.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service through Traefik, you need to add specific labels to its service definition. The `proxy` network must also be defined as external for Traefik to discover the service.

Here's a concrete YAML example for a hypothetical `my-app-compose.yml`:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy
    restart: unless-stopped

networks:
  proxy:
    external: true # Assumes 'proxy' network is created by routing-compose.yml
```

After placing this file in `/docker/` and running `m3tal up`, your `my-app` service would be accessible via `http://app.YOUR_DOMAIN`.

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api.service`
*   **Restart service**: `systemctl restart m3tal-api.service`
*   **View logs**: `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly get started with the M3TAL Dashboard:

*   To start *only* the dashboard container, which includes downloading the latest compose files and applying the correct expose mode based on your `.env` configuration:
    ```bash
    m3tal dash up
    ```
    Once started, access the dashboard at `http://HOST_IP:8082` (if `DASHBOARD_EXPOSE_MODE=local`) or `http://dash.DOMAIN` (if `DASHBOARD_EXPOSE_MODE=traefik` and Traefik is running).

*   To orchestrate and deploy all other stacks in the `/docker/` directory, including any user-defined compose files you have placed there:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :--------------------------------------------- | :--------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.