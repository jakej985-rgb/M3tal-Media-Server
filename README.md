# M3TAL System Documentation

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

### Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## System Components

The M3TAL system comprises several interconnected components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing a single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container, listening internally on port `8082`. It communicates with the API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container that provides a Cloudflare tunnel for secure, zero-configuration internet access to services.

## Filesystem Contract

The following paths are central to the M3TAL system's operation and configuration:

| Path                        | Purpose                                                                   |
| :-------------------------- | :------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configurations. |
| `/docker`                   | A symlink pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                 |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all Docker Compose stack files and associated Traefik dynamic configuration files reside. For user convenience, `/docker` serves as a symlink alias to `/opt/m3tal/stack/` for all stack-related operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy and start all defined stacks, including your new one.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`.

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Configuration

M3TAL's primary configuration is managed through environment variables stored in `/etc/m3tal/.env`. The `m3tal config wizard` command provides an interactive interface for managing these variables. Alternatively, individual variables can be set using `m3tal config set KEY value`.

## Dashboard Access

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Setting**: `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism**: The dashboard container exposes its internal port `8082` directly to the host machine, typically at `HOST_IP:8082`. This is achieved via an override (`m3tal-compose.local.yml`) that adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`).
*   **Access**: `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: No Traefik gateway is required for this mode. A new user performing a default installation will access the dashboard directly via port 8082.
*   **Use Case**: Ideal for LAN-only setups, initial configurations, and local testing environments.

### Traefik Mode

*   **Setting**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: Traefik labels are applied to the dashboard container (`m3tal-compose.traefik.yml`), instructing the Traefik gateway to route requests for `dash.${DOMAIN}` to the dashboard container's internal port `8082`. This requires the Traefik gateway to be running.
*   **Access**: `http://dash.DOMAIN`
*   **Requirements**: The Traefik gateway must be deployed and operational via `m3tal up`.
*   **Use Case**: Suitable for domain-based setups where multiple services are exposed behind Traefik as a central reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration from files located in `/docker/dynamic/` (which symlinks to `/opt/m3tal/stack/dynamic/`). This file provider allows for hot-reloading of routing rules. For instance, the `dynamic/api.yml` file is used to route requests for `api.DOMAIN` to the M3TAL Go API daemon, which listens on the host-local port `8080`. This routing is achieved by directing traffic to `http://host.docker.internal:8080` from within the Traefik container. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard is routed via `dash.DOMAIN` to the dashboard container's internal port `8082`.

### Exposing Custom Services

To expose a custom user service via Traefik labels, include the following within your service's definition in its Docker Compose file (e.g., `my-app-compose.yml`):

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
      - proxy # Ensure your service is connected to the 'proxy' network
```

## Quick Demo

To quickly launch the M3TAL dashboard:

*   Run `m3tal dash up`. This command specifically downloads the latest dashboard Docker Compose files and starts the dashboard container with the appropriate `DASHBOARD_EXPOSE_MODE` override (defaulting to local mode).

To deploy all other M3TAL stacks and any user-defined Docker Compose files placed in the `/docker/` directory:

*   Run `m3tal up`. This command orchestrates and deploys all services defined across all `*-compose.yml` files in the `/docker/` directory.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.