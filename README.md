# M3TAL System Documentation

## Overview
This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites
Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract
The M3TAL system adheres to the following filesystem structure for critical components and user interaction:

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains Docker Compose files and Traefik dynamic configuration. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Components
The M3TAL system is composed of the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all system operations and management.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It is responsible for managing Docker interactions, maintaining the SQLite state database, and handling API routes for system control.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask Docker container operating internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It uses a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel Docker container, facilitating zero-configuration internet access for exposed services.

## Port Map
The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|-----------------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations
If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management
The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Standard `systemctl` commands can be used for its operation:

*   **Check status:** `systemctl status m3tal-api.service`
*   **Restart service:** `systemctl restart m3tal-api.service`
*   **View logs:** `journalctl -u m3tal-api.service -f`

## Deployment Lifecycle
M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack. The `/opt/m3tal/stack/` directory serves as the canonical source of truth where all stack files reside, while `/docker` is a user-facing symlink alias for all stack operations.

### Adding a New Stack
To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access
The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (Default)
*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism:** M3TAL applies an override (`m3tal-compose.local.yml`) that adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`).
*   **Access:** The M3TAL Dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik gateway is required for this mode.
*   **Use Case:** Ideal for LAN-only setups, initial configurations, and first-time users who will access the dashboard directly via port 8082.

### Traefik Mode
*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`.
*   **Mechanism:** M3TAL applies an override (`m3tal-compose.traefik.yml`) that adds specific Traefik labels to the dashboard service. Traefik, when running, interprets these labels to route traffic for `dash.DOMAIN` to the dashboard container on its internal port `8082`.
*   **Access:** The M3TAL Dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.localhost`).
*   **Requirements:** The Traefik gateway must be deployed and running via `m3tal up` for this mode to function.
*   **Use Case:** Suitable for domain-based setups and environments where multiple services are managed behind a single reverse proxy.

## Traefik Gateway
Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed as a Docker container via `routing-compose.yml`. It binds host port `80` as its primary HTTP entry point and dynamically loads additional routing configurations from `/docker/dynamic/` (using a file provider with hot-reloading).

### Dynamic Configuration
For services that are not Docker containers or need specialized routing (e.g., to a host-local daemon), Traefik uses dynamic configuration files. For example, `dynamic/api.yml` routes `api.DOMAIN` to the M3TAL Go API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container when `DASHBOARD_EXPOSE_MODE=traefik` via labels.

Example `dynamic/api.yml`:
```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)"
      service: api
      entryPoints:
        - web

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080"
```

### Exposing a Custom Service via Traefik
To expose a custom user service (e.g., `my-app`) via Traefik, you need to add specific labels to its service definition in your Docker Compose file (e.g., `my-app-compose.yml`). The service must also be part of the `proxy` network, which Traefik monitors.

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true" # Enable Traefik for this service
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN}`)" # Define the host rule (e.g., app.localhost)
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Target port inside the container
      - "traefik.http.routers.myapp.entrypoints=web" # Use the 'web' entrypoint (HTTP port 80)
      - "traefik.docker.network=proxy" # Explicitly tell Traefik which network to look for the service on
    networks:
      - proxy # Connect the service to the external 'proxy' network

networks:
  proxy:
    external: true # Refer to the existing external Traefik proxy network
```
After placing this file in `/docker/`, run `m3tal up` to deploy the service and make it discoverable by Traefik.

## Quick Demo
This section outlines quick commands for common deployment scenarios.

*   **Start the M3TAL Dashboard specifically:**
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard Docker Compose files and starts the dashboard container with the appropriate `DASHBOARD_EXPOSE_MODE` override (defaulting to local mode, accessible at `http://HOST_IP:8082`).

*   **Deploy all M3TAL stacks and user-defined Docker Compose files:**
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all `*-compose.yml` files found within the `/docker/` directory, including the core M3TAL components (like Traefik and any services defined in `routing-compose.yml`) and any user-defined Docker Compose files you have placed there.