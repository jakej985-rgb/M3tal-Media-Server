# M3TAL System Documentation

## Overview
This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites
Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## APT Installation
To install the M3TAL CLI binary and API daemon:

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

## Filesystem Contract
The following table outlines the key directories and files used by the M3TAL system:

| Path                        | Purpose                                                                                                                                                                  |
|-----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`.                                                                                  |
| `/var/lib/m3tal/state.db`   | SQLite state database storing system configuration and operational data. Auto-created and managed by the API daemon.                                                     |
| `/opt/m3tal/stack/`         | Canonical stack directory. This location contains all core Docker Compose files and Traefik dynamic configuration files. This is the source of truth for all stack files. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose stack operations and where user-defined Compose files should be placed.                 |
| `/docker/users.json`        | Dashboard credential store. Contains user hashes for dashboard access. Managed by `m3tal dashpass`.                                                                        |

## Components
The M3TAL system consists of several integrated components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary providing a single entrypoint for all M3TAL operations, including configuration, service management, and Docker orchestration.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker containers, interacts with the SQLite state database, and exposes RESTful API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container. It runs internally on port `8082` and communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It uses a file provider for dynamic routing and automatically discovers Docker services via labels.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing secure, zero-config internet access to services without requiring open inbound firewall ports.

## Deployment Lifecycle
M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all core Docker Compose stack files and Traefik configuration reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, providing a convenient, user-facing path for all stack operations.

### Adding a New Stack
To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy your new stack along with all other M3TAL-managed services.

## Dashboard Access
The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### 1. Local Mode (Default)
*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting upon initial installation).
*   **Mechanism:** Uses the `m3tal-compose.local.yml` override file, which adds a direct port binding to the host machine.
*   **Access:** The M3TAL Dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirement:** No Traefik gateway is required for this mode. A new user performing a default installation will access the dashboard directly via port 8082.
*   **Use Case:** Ideal for LAN-only setups, first-time users, and local development or testing environments where domain routing is not necessary.

### 2. Traefik Mode
*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`.
*   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override file. This file applies Traefik labels to the dashboard container, allowing Traefik to route traffic.
*   **Access:** The M3TAL Dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Requirement:** The Traefik gateway must be running (typically started via `m3tal up`) and configured to handle the specified domain.
*   **Use Case:** Best for domain-based deployments where multiple services are exposed behind a single reverse proxy.

## Traefik Gateway
Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files (such as `dynamic/api.yml` located in `/docker/dynamic/`) for routing requests to services that are not managed directly by Docker Compose labels, such as host-local daemons. For example, `api.DOMAIN` is routed to the M3TAL Go API daemon, which listens on host-local port `8080`, by routing to `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE=traefik`, leveraging labels on the dashboard service.

### Exposing a Custom Service via Traefik
To expose a custom user service named `my-app` via Traefik, you would add labels to its Docker Compose service definition as shown in this example for a hypothetical `my-app-compose.yml`:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx-app
    labels:
      - "traefik.enable=true" # Enable Traefik for this service
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)" # Route based on domain, e.g., http://app.example.com
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Target port inside the container
      - "traefik.http.routers.myapp.entrypoints=web" # Use the 'web' entrypoint (HTTP port 80)
    networks:
      - proxy # Ensure the service is connected to the 'proxy' network for Traefik discovery

networks:
  proxy:
    external: true # Use the existing 'proxy' network managed by M3TAL
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically configure routing for `app.DOMAIN` to your `my-app` service.

## Service Management
The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api.service`
*   **Restart Service:** `systemctl restart m3tal-api.service`
*   **View Logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo
This section provides a quick overview of how to manage M3TAL services.

*   To start only the M3TAL Dashboard container, use the specific command:
    ```bash
    m3tal dash up
    ```
    This command will download the necessary dashboard Compose files and start the dashboard according to your `DASHBOARD_EXPOSE_MODE` setting.

*   To orchestrate and deploy all other stacks, including any user-defined Docker Compose files placed in the `/docker/` directory, use the general deployment command:
    ```bash
    m3tal up
    ```
    This command processes all `*-compose.yml` files in `/docker/`, ensuring all services are brought up or updated as defined.

## Port Map
The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|-----------------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine (e.g., `http://localhost:8081`). |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.