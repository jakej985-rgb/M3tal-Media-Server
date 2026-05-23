# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

M3TAL is distributed via an APT repository for easy installation and updates.

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

The M3TAL system consists of several integrated components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary that serves as the single entrypoint for all M3TAL operations, installed via APT.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service. It listens on host-local port `8080` and is responsible for managing Docker interactions, maintaining the SQLite state database, and serving API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container. It listens internally on port `8082` and communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy. It exposes services by domain name on host port `80` and uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container that provides a Cloudflare tunnel for zero-configuration internet access to services.

## Filesystem Contract

The following paths define the core filesystem structure and their purposes within the M3TAL system:

| Path                        | Purpose                                                                                                                                                             |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | The primary configuration file for M3TAL, containing environment variables. It is managed by the `m3tal config wizard` command.                                       |
| `/var/lib/m3tal/state.db`   | The SQLite state database used by the M3TAL API daemon. This file is automatically created and managed by the API daemon.                                           |
| `/opt/m3tal/stack/`         | The canonical directory for all Docker Compose stack files and Traefik dynamic configuration files. This is the source of truth for M3TAL's orchestrated services. |
| `/docker`                   | A symbolic link pointing to `/opt/m3tal/stack/`. This path serves as the user-facing alias for all stack-related operations.                                       |
| `/docker/users.json`        | The credential store for the M3TAL Dashboard, managed by the `m3tal dashpass` command.                                                                              |

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                                                                     |
| :--- | :-------------------------- | :------------------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                                                                   |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                                                              |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                                                                     |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory is a user-facing symlink to `/opt/m3tal/stack/`. This means that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside, and `/docker` is the convenient alias for all user stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**: This is the default setting for M3TAL installations.
*   **Access Method**: Direct port binding at `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Behavior**: A new user performing a default M3TAL installation will access the dashboard directly via port `8082`. This mode explicitly adds a port binding (`${DASHBOARD_PORT:-8082}:8082`) to the `m3tal-dashboard` container, bypassing Traefik. No Traefik configuration is required for this mode.
*   **Use Cases**: Ideal for LAN-only setups, initial configurations, and local testing environments.

### Traefik Mode

*   **`DASHBOARD_EXPOSE_MODE=traefik`**:
*   **Access Method**: Domain routing at `http://dash.DOMAIN`.
*   **Behavior**: In this mode, Traefik labels are added to the `m3tal-dashboard` container definition. Traefik, if running, then routes requests for `dash.DOMAIN` to the dashboard container's internal port `8082`. Requires the Traefik gateway to be deployed and operational.
*   **Use Cases**: Suitable for domain-based setups where multiple services are exposed through a single reverse proxy.

## Traefik Gateway

Traefik acts as the central reverse proxy for M3TAL, automatically discovering and routing traffic to Docker services. This is achieved by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files, such as `dynamic/api.yml` located within `/docker/dynamic/`, to route requests to services that are not directly Docker containers but are listening on host-local ports. For example, `api.DOMAIN` is routed to the M3TAL Go API daemon running on the host machine at `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` is routed to the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE=traefik` is configured.

### Exposing a Custom User Service via Traefik

To expose your own Docker Compose service via Traefik, add the appropriate labels to its service definition:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy

networks:
  proxy:
    external: true
```

After placing this file in `/docker/` and running `m3tal up`, your `my-app` service would be accessible at `http://app.DOMAIN` (assuming `DOMAIN` is set in `/etc/m3tal/.env`).

## Service Management

The M3TAL API daemon operates as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check Status**: `systemctl status m3tal-api`
*   **Restart Service**: `systemctl restart m3tal-api`
*   **View Logs**: `journalctl -u m3tal-api -f` (for real-time logs)

## Quick Demo

To quickly get started with the M3TAL Dashboard:

1.  **Start the Dashboard**: Execute `m3tal dash up`. This command specifically downloads the latest dashboard Docker Compose files and starts the `m3tal-dashboard` container using the appropriate configuration (defaulting to local mode).
2.  **Access Dashboard**: If `DASHBOARD_EXPOSE_MODE` is set to `local` (the default), access your dashboard directly via `http://HOST_IP:8082`.
3.  **Deploy All Stacks**: To deploy all other M3TAL core stacks (including Traefik) and any user-defined Docker Compose files placed in `/docker/`, run `m3tal up`. This command orchestrates and deploys everything else in the directory.