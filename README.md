# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI binary and related system components, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The M3TAL system relies on a defined filesystem structure for its operation and configuration:

| Path                     | Purpose                                                       |
|--------------------------|---------------------------------------------------------------|
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created by the API daemon.        |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik configuration. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.      |

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/` for all stack operations, simplifying access and management.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy your new stack along with all other configured stacks.

## Dashboard Access

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

When `DASHBOARD_EXPOSE_MODE=local` (which is the default setting for a new installation), the M3TAL Dashboard container exposes its internal port directly to the host.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** Utilizes `m3tal-compose.local.yml` to add a direct port binding, typically `${DASHBOARD_PORT:-8082}:8082`.
*   **Access:** The dashboard is accessible directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirement:** No Traefik gateway is required for this mode. A new user performing a default installation will access the dashboard directly via port 8082.
*   **Use Case:** Ideal for LAN-only setups, first-time users, and local development or testing without domain configuration.

### Traefik Mode

When `DASHBOARD_EXPOSE_MODE=traefik`, the M3TAL Dashboard container is configured to be discoverable and routable by the Traefik Gateway.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** Utilizes `m3tal-compose.traefik.yml` to apply Traefik labels to the dashboard container, allowing Traefik to route traffic based on a defined hostname.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`). This requires the Traefik gateway to be running and correctly configured (via `m3tal up`).
*   **Use Case:** Suitable for domain-based deployments and integrating the dashboard with other services behind a reverse proxy.

## Traefik Gateway

Traefik acts as the M3TAL system's reverse proxy and API gateway. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files, typically located in `/docker/dynamic/` (which is mapped to `/etc/traefik/dynamic` within the Traefik container), to define routes for services that are not managed directly as Docker containers or require specific routing rules. For instance:

*   **M3TAL API Daemon:** The M3TAL API daemon runs on the host-local port `8080`. Traefik routes requests from `api.DOMAIN` to this daemon using a dynamic configuration file like `dynamic/api.yml`. This file defines a router and service that targets `http://host.docker.internal:8080`.
*   **M3TAL Dashboard:** When `DASHBOARD_EXPOSE_MODE=traefik`, requests to `dash.DOMAIN` are routed to the `m3tal-dashboard` container by Traefik, based on labels defined in `m3tal-compose.traefik.yml` (targeting the container's internal port `8082`).

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service through Traefik, you must add specific labels to its service definition. The `proxy` network must also be defined and attached to your service.

Here's an example for a hypothetical `my-app-compose.yml`:

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

After placing this `my-app-compose.yml` file in `/docker/` and running `m3tal up`, Traefik will automatically configure routing for `app.DOMAIN` to your `my-app` service.

## Service Management

The M3TAL API daemon is deployed as a systemd service named `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api.service`
*   **Restart Service:** `systemctl restart m3tal-api.service`
*   **View Logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly get started and interact with M3TAL:

*   To start only the M3TAL dashboard container, run:
    ```bash
    m3tal dash up
    ```
    This command specifically manages the dashboard, downloading the latest compose files and applying the correct override (`.local` or `.traefik`) based on your `/etc/m3tal/.env` configuration.
*   To orchestrate and deploy all other stacks defined in your `/docker/` directory, including any custom user-defined compose files, run:
    ```bash
    m3tal up
    ```
    This command will bring up all services, including Traefik if configured.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                    | Access                                  | Description                                                                                             |
|------|----------------------------|-----------------------------------------|---------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point   | Public                                  | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)      | Host-local                              | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard          | Host-local only                         | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard            | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.