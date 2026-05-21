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

The M3TAL system is composed of the following key elements:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container exposing services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, facilitating zero-configuration internet access for exposed services.

## Filesystem Contract

The following table outlines the critical filesystem paths used by the M3TAL system:

| Path                        | Purpose                                                                                                                                                                 |
| :-------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`.                                                                                 |
| `/var/lib/m3tal/state.db`   | SQLite state database, used by the API daemon for internal state management. Auto-created by the API daemon upon first run.                                               |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configuration. This is the source of truth for all deployments.                                |
| `/docker`                   | A symbolic link pointing to `/opt/m3tal/stack/`. This path is the user-facing alias for all stack operations and where user-defined Compose files should be placed.      |
| `/docker/users.json`        | Credential store for the M3TAL Dashboard. Managed by the `m3tal dashpass` command.                                                                                      |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth for all M3TAL-managed stack files. `/docker` is a user-facing symlink alias to `/opt/m3tal/stack/` designed for ease of access and management of compose files.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added service.

## Dashboard Access

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (Default)
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` (default setting).
*   **Access Method**: Direct port binding at `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port mapping (`${DASHBOARD_PORT:-8082}:8082`) to the dashboard container.
*   **Requirements**: No Traefik gateway is required for this mode.
*   **Behavior for New Users**: A user performing a default installation will access the dashboard directly via port 8082, as this is the default behavior enabled by `DASHBOARD_EXPOSE_MODE=local`.

### Traefik Mode
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`.
*   **Access Method**: Domain routing at `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Mechanism**: This mode uses the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard container. Traefik interprets these labels to route `dash.${DOMAIN}` to the dashboard service on its internal port `8082`.
*   **Requirements**: The Traefik gateway must be running (`m3tal up` will start it if configured) and configured to handle the specified domain.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik leverages dynamic configuration files, such as `dynamic/api.yml`, to route requests to services not directly exposed as Docker containers (e.g., host-local daemons). For instance, `api.DOMAIN` is routed to the M3TAL Go API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container's internal port `8082` when `DASHBOARD_EXPOSE_MODE=traefik`.

### Example: Exposing a Custom User Service via Traefik

To expose a hypothetical `my-app` service running Nginx through Traefik, you would add the following labels to its service definition in your `my-app-compose.yml` file located in `/docker/`:

```yaml
services:
  my-app:
    image: nginx:alpine
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy
```

This configuration instructs Traefik to enable routing for `my-app`, respond to `app.DOMAIN`, and forward requests to the container's internal port `80`. The `proxy` network must be defined as an external network for Traefik to discover the service.

## Service Management

The M3TAL API daemon is managed as a `systemd` service named `m3tal-api.service`. Standard `systemctl` commands can be used for its lifecycle management:

*   **Check Status**: `systemctl status m3tal-api.service`
*   **Restart Service**: `systemctl restart m3tal-api.service`
*   **View Logs**: `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly get started with the M3TAL Dashboard:

1.  **Start the M3TAL Dashboard**: Execute `m3tal dash up`. This command specifically downloads the latest dashboard Compose files and starts the dashboard container, respecting the `DASHBOARD_EXPOSE_MODE` configured in `/etc/m3tal/.env`. By default (`DASHBOARD_EXPOSE_MODE=local`), you will access the dashboard at `http://HOST_IP:8082`.
2.  **Deploy All M3TAL Stacks**: Use `m3tal up` to orchestrate and deploy all Docker Compose stacks defined by `*-compose.yml` files found in the `/docker/` directory. This includes the `routing-compose.yml` (Traefik, Cloudflared) and any other user-defined compose files.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                    | Access                                        | Description                                                                                             |
| :--- | :------------------------- | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point   | Public                                        | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)      | Host-local                                    | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard          | Host-local only                               | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard            | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.