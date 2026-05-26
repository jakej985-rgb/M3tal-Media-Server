# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI and API daemon via APT:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-apt-key/public.key | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-apt-key stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Core Components

The M3TAL system comprises the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container, listening internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for establishing secure, zero-configuration internet access to services.

## Filesystem Contract

The M3TAL system adheres to the following filesystem structure:

| Path | Purpose |
| :----------------------- | :------------------------------------------------------------------------------- |
| `/etc/m3tal/.env` | The primary configuration file, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | The SQLite state database, automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/` | The canonical directory containing all Docker Compose files and Traefik configuration. |
| `/docker` | A symlink to `/opt/m3tal/stack/`. This path serves as the user-facing alias for all stack operations. |
| `/docker/users.json` | The credential store for the M3TAL Dashboard, managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all stack files reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, providing a convenient user-facing path for all stack management operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables are configured in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to deploy the new stack along with all other managed services.

## Dashboard Access

The M3TAL Dashboard offers two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism**: The `m3tal-dashboard` container's port `8082` is directly bound to the host's `DASHBOARD_PORT` (defaulting to `8082`). No Traefik configuration is involved.
*   **Access**: A new user performing a default installation will access the dashboard directly via port `8082` at `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Use Case**: Ideal for LAN-only setups, first-time users, and local development environments.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: The `m3tal-dashboard` container is configured with Traefik labels, enabling Traefik to route requests for `dash.DOMAIN` to the dashboard's internal port `8082`. Traefik must be running for this mode to function.
*   **Access**: Access is via the domain name configured for M3TAL, typically `http://dash.DOMAIN`.
*   **Use Case**: Suited for domain-based deployments and environments where multiple services are exposed through a single reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik acts as the central ingress point, binding to host port `80` (and `443` if TLS is enabled). It integrates with Docker to dynamically configure routes for services and loads additional dynamic configuration from `/docker/dynamic/` (which is mapped to `/etc/traefik/dynamic` inside the Traefik container).

*   **API Daemon Routing**: The M3TAL API daemon, running on the host machine at `http://127.0.0.1:8080`, is exposed via Traefik using a dynamic configuration file (e.g., `/docker/dynamic/api.yml`). This file routes `api.DOMAIN` to `http://host.docker.internal:8080`, allowing containers within the Docker network to access the host-local API.
*   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE=traefik`, the `m3tal-dashboard` container is configured with Traefik labels to route `dash.DOMAIN` requests to its internal port `8082`.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service (e.g., an Nginx instance) via Traefik, add appropriate labels to its service definition:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true" # Enable Traefik for this service
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)" # Route requests for app.DOMAIN
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Target internal port of the service
      - "traefik.http.routers.myapp.entrypoints=web" # Use the 'web' (HTTP) entry point
    networks:
      - proxy # Connect to the shared Traefik network

networks:
  proxy:
    external: true # Use the existing 'proxy' network
```

After placing this file in `/docker/`, run `m3tal up` to deploy the service. Traefik will automatically detect the labels and expose your `my-app` service at `http://app.DOMAIN`.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Use standard `systemctl` commands for its operation:

*   **Check Status**: `systemctl status m3tal-api`
*   **Restart Service**: `systemctl restart m3tal-api`
*   **View Logs**: `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get the M3TAL dashboard running:

*   Run `m3tal dash up` to specifically start and manage the `m3tal-dashboard` container. This command handles downloading the necessary Compose files and applying the correct override based on your `DASHBOARD_EXPOSE_MODE` setting.
*   To deploy all other M3TAL managed stacks, including Traefik and any user-defined Docker Compose files placed in `/docker/`, execute `m3tal up`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :------------------------------------------- | :------------------------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.