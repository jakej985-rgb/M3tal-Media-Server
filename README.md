# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

M3TAL is distributed as an APT package. Follow these steps to install:

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

## System Components

The M3TAL system is comprised of several interconnected components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the primary command-line interface for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes the M3TAL API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container, internally listening on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container that establishes a Cloudflare tunnel for secure, zero-configuration internet access to services.

## Filesystem Contract

The following paths represent the canonical locations for M3TAL system files and directories:

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created by the API daemon.    |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configurations. |
| `/docker`                   | A symlink pointing to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.           |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all Docker Compose stack files (e.g., `m3tal-compose.yml`, `routing-compose.yml`) and Traefik dynamic configuration files reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, providing a convenient and consistent user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Quick Demo

To quickly start the M3TAL Dashboard:

1.  Execute `m3tal dash up`. This command specifically downloads the latest dashboard compose files and starts the `m3tal-dashboard` container with the appropriate configuration based on your `DASHBOARD_EXPOSE_MODE` setting.
2.  Access the dashboard as described in the Dashboard Access section.

The `m3tal up` command, in contrast, orchestrates and deploys *all* Docker Compose stacks found in the `/docker/` directory, including the `m3tal-dashboard`, `routing-compose.yml`, and any user-defined compose files.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override file, which adds a direct Docker port binding of `${DASHBOARD_PORT:-8082}:8082`. No Traefik configuration is involved.
*   **Access**: A new user performing a default installation will access the dashboard directly via port 8082 at:
    *   `http://HOST_IP:8082`
    *   `http://localhost:8082`
*   **Use Case**: Ideal for LAN-only setups, first-time users, and local testing where domain-based access is not required.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: This mode uses the `m3tal-compose.traefik.yml` override file. This file adds specific Traefik labels to the `m3tal-dashboard` service, allowing Traefik to route incoming requests to the dashboard. Traefik must be running for this mode to function.
*   **Access**: Via a configured domain:
    *   `http://dash.DOMAIN` (e.g., `http://dash.example.com`), where `DOMAIN` is defined in `/etc/m3tal/.env`.
*   **Use Case**: Best for domain-based setups and environments where multiple services are exposed through a single reverse proxy.

## Traefik Gateway

Traefik operates as the central reverse proxy for M3TAL, deployed via `routing-compose.yml`. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds to host port `80` (and `443` if HTTPS is configured) as its primary HTTP entry point. It utilizes a file provider to load dynamic configurations from `/docker/dynamic/`, allowing for hot-reloading of routing rules.

*   **API Routing**: Traefik routes `api.DOMAIN` to the M3TAL API daemon. This is achieved through a dynamic configuration file (e.g., `/docker/dynamic/api.yml`) that defines a router for `api.DOMAIN` and forwards requests to `http://host.docker.internal:8080`, where the Go API daemon listens.
*   **Dashboard Routing**: When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.DOMAIN` to the `m3tal-dashboard` container using labels specified in `m3tal-compose.traefik.yml`.

### Example: Exposing a Custom User Service via Traefik

To expose a hypothetical `my-app` service from a `my-app-compose.yml` file via Traefik:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN}`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy

networks:
  proxy:
    external: true
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically discover the `my-app` service and route requests for `http://app.DOMAIN` (where `DOMAIN` is from your `.env` file) to port `80` of the `my-app` container.

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :-------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------ |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine (e.g., `http://127.0.0.1:8081`). |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.