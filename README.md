# M3TAL System Documentation

## Overview
This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites
Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation
The M3TAL CLI binary and API daemon are installed via APT.

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
The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the primary command-line interface for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary operating as a systemd service, listening on host-local port `8080`. It manages Docker interactions, maintains the M3TAL state database, and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy Docker container that exposes services on host port `80` by domain name. It utilizes a file provider for dynamic configuration and Docker labels for service discovery.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel Docker container, facilitating secure zero-configuration public internet access to services.

## Filesystem Contract
The following table outlines key file system paths and their purposes within the M3TAL system:

| Path | Purpose |
| :------------------------- | :--------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | The primary configuration file for M3TAL, managed by the `m3tal config wizard` command.            |
| `/var/lib/m3tal/state.db`  | The SQLite state database for the M3TAL API daemon, automatically created upon first run.          |
| `/opt/m3tal/stack/`        | The canonical directory for all Docker Compose stack files and associated Traefik dynamic configurations. |
| `/docker`                  | A symbolic link to `/opt/m3tal/stack/`, serving as the user-facing path for all stack operations. |
| `/docker/users.json`       | The credential store for the M3TAL Dashboard, managed by the `m3tal dashpass` command.             |

## Deployment Lifecycle
M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack. The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside, and `/docker` is a user-facing symlink alias for all stack operations.

### Adding a New Stack
To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`.
3.  Execute `m3tal up` to deploy all stacks, including your newly added service.

## Dashboard Access
The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism**: This mode uses the `m3tal-compose.local.yml` override file to add a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the host.
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirement**: No Traefik instance is required for this mode. A new user performing a default installation will access the dashboard directly via port 8082, as this is the behavior linked to the default `DASHBOARD_EXPOSE_MODE=local` setting.
*   **Use Case**: Ideal for LAN-only setups, initial configurations, and local development.

### Traefik Mode
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: This mode utilizes the `m3tal-compose.traefik.yml` override file to add Traefik labels to the dashboard container. Traefik interprets these labels to route `dash.${DOMAIN}` to the dashboard container's internal port `8082`.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN`.
*   **Requirement**: The Traefik gateway must be running (`m3tal up` will start it if configured).
*   **Use Case**: Suitable for domain-based deployments where multiple services are exposed through a reverse proxy.

## Traefik Gateway
Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed as a Docker container, typically binding host port `80` as its HTTP entry point. It loads dynamic configurations from `/docker/dynamic/` (using a file provider with hot-reload).

For example, Traefik routes requests for `api.DOMAIN` to the M3TAL Go API daemon. This is achieved through a dynamic configuration file (e.g., `dynamic/api.yml`) that maps `api.DOMAIN` to `http://host.docker.internal:8080`, where the API daemon is listening on a host-local port. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, Traefik labels on the dashboard container enable routing `dash.DOMAIN` to the dashboard container.

### Exposing a Custom User Service via Traefik
To expose a custom Docker Compose service (e.g., `my-app`) via Traefik, add the necessary Traefik labels to its service definition. Ensure your service is part of the `proxy` network, which Traefik is configured to monitor.

Example `my-app-compose.yml`:
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

networks:
  proxy:
    external: true
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically configure routing for `app.DOMAIN` to your `my-app` service.

## Service Management
The M3TAL API daemon (`m3tal-api.service`) is managed as a systemd service. Standard systemctl commands apply:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo
To quickly start and access the M3TAL Dashboard:

*   Run `m3tal dash up` to specifically start the dashboard container. This command handles downloading the necessary compose files and applying the correct override based on your `DASHBOARD_EXPOSE_MODE` setting (defaulting to local access via port 8082).
*   Run `m3tal up` to orchestrate and deploy all other stacks found in the `/docker/` directory, including any user-defined compose files and the Traefik gateway if configured.

## Port Map
The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :--------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                   | The public-facing HTTP port for services exposed via Traefik.                                                                      |
| 8080 | M3TAL API daemon (Go)       | Host-local                               | The internal port the M3TAL API daemon listens on.                                                                                 |
| 8081 | Traefik dashboard           | Host-local only                          | The internal Traefik dashboard port, accessible only from the host machine.                                                        |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.