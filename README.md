# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL system, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Core Components

The M3TAL system is comprised of several interconnected components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing the primary command-line interface for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on port `8080`. It is responsible for managing Docker interactions, maintaining the system's state database, and exposing API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container. It runs internally on port `8082` and communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes Docker services by domain name, typically on host port `80`. It dynamically configures routing based on Docker service labels and external configuration files.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container that facilitates secure, zero-config internet access to services without requiring open inbound firewall ports.

## Filesystem Contract

The following table details the critical directories and files within the M3TAL filesystem:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.   |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains Docker Compose files and Traefik configuration. |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth where all M3TAL-managed Docker Compose stack files and Traefik dynamic configurations reside. For user convenience and interaction, `/docker` is established as a symlink alias to `/opt/m3tal/stack/`. All user-facing stack operations, such as adding new Compose files, should be directed to the `/docker/` path.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all deployed stacks, including your new service.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`)

This is the default access method upon a new installation.
*   **Configuration**: `DASHBOARD_EXPOSE_MODE` is set to `local` in `/etc/m3tal/.env`. This enables the `m3tal-compose.local.yml` override, which adds a direct port binding to the dashboard container.
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: No Traefik gateway is required for this mode.
*   **Use Case**: Ideal for LAN-only setups, initial configurations, first-time users performing a default installation, and local testing environments. A new user performing a default installation will access the dashboard directly via port `8082`.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

This mode routes dashboard access through the Traefik gateway.
*   **Configuration**: `DASHBOARD_EXPOSE_MODE` is set to `traefik` in `/etc/m3tal/.env`. This enables the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard service.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN`, where `DOMAIN` is configured in `/etc/m3tal/.env`.
*   **Requirements**: The Traefik gateway must be running via `m3tal up` for this mode to function.
*   **Use Case**: Suitable for domain-based environments where multiple services are exposed through a centralized reverse proxy.

## Traefik Gateway

Traefik operates as the central ingress point for HTTP traffic within M3TAL. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik's configuration is managed through a combination of its static configuration (e.g., `traefik.yml`) and dynamic configuration files loaded via its file provider. For instance, dynamic configuration files such as `/docker/dynamic/api.yml` are used to define routes for services that are not Docker containers themselves but listen on host-local ports, such as the M3TAL API daemon. This file routes requests from `api.DOMAIN` to the Go API daemon listening on `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE=traefik` is enabled, leveraging Traefik labels within the dashboard's Docker Compose definition.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service through Traefik, you need to add specific labels to its service definition. Here is an example for a hypothetical `my-app-compose.yml` that exposes an Nginx container as `app.DOMAIN`:

```yaml
# /docker/my-app-compose.yml
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

networks:
  proxy:
    external: true
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically detect and configure routing for `app.DOMAIN` to your `my-app` service.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon is managed by `systemd` as `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check Status**: `systemctl status m3tal-api.service`
*   **Restart Service**: `systemctl restart m3tal-api.service`
*   **View Logs**: `journalctl -u m3tal-api.service -f`

## Quick Demo

The M3TAL CLI provides specific commands to manage individual components or the entire system:

*   **Start the Dashboard Only**: To start just the M3TAL Dashboard container, use `m3tal dash up`. This command specifically downloads the latest dashboard Docker Compose files and starts the dashboard with the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.
*   **Deploy All Stacks**: To orchestrate and deploy all Docker Compose stacks defined in the `/docker/` directory, including any user-defined compose files, run `m3tal up`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :--------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                         | The public-facing HTTP port for services exposed via Traefik.                                                                          |
| 8080 | M3TAL API daemon (Go)       | Host-local                                     | The internal port the M3TAL API daemon listens on.                                                                                     |
| 8081 | Traefik dashboard           | Host-local only                                | The internal Traefik dashboard port, accessible only from the host machine.                                                            |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.