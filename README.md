# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system. It describes the system's architecture, component interactions, deployment mechanisms, and management interfaces.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI and API daemon via APT, execute the following commands:

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

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, facilitating zero-configuration internet access for services.

## Filesystem Contract

The following table details the critical directories and files within the M3TAL filesystem:

| Path                        | Purpose                                                                       |
| :-------------------------- | :---------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | The primary configuration file, managed by `m3tal config wizard`.             |
| `/var/lib/m3tal/state.db`   | The SQLite state database, automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | The canonical directory containing all Docker Compose stack files and Traefik configuration. |
| `/docker`                   | A symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | The credential store for the M3TAL Dashboard, managed by `m3tal dashpass`.    |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack. Note that `/docker` is a symlink to `/opt/m3tal/stack/`, which is the canonical source of truth directory where all stack files reside.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to start all deployed stacks, including your new service.

## Dashboard Access

The M3TAL Dashboard supports two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation.)
*   **Mechanism:** This mode utilizes the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`).
*   **Access:** The M3TAL Dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** Traefik is not required for this mode. A new user performing a default M3TAL installation will access the dashboard directly via port 8082.
*   **Use Case:** Ideal for LAN-only setups, first-time users, and local development or testing environments.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode utilizes the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard service. Traefik then routes traffic from `dash.DOMAIN` to the dashboard container's internal port `8082`.
*   **Access:** The M3TAL Dashboard is accessible via `http://dash.DOMAIN`.
*   **Requirements:** Traefik must be running as part of the M3TAL deployment (`m3tal up`).
*   **Use Case:** Suited for domain-based setups and environments where multiple services are exposed behind a single reverse proxy.

## Traefik Gateway

Traefik is deployed as a Docker container via `routing-compose.yml` and serves as the primary ingress point for external services. It binds host port `80` (and `443` if HTTPS is configured) as its HTTP entry point.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also loads dynamic configuration files (such as those in `/docker/dynamic/`) which allow routing requests to services not directly managed by Docker Compose labels. For instance, the `dynamic/api.yml` file configures Traefik to route `api.DOMAIN` to the Go API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the `m3tal-dashboard` container when `DASHBOARD_EXPOSE_MODE=traefik`.

### Exposing a Custom Service via Traefik

To expose a custom user service, define appropriate Traefik labels in its Docker Compose service definition, ensuring it is connected to the `proxy` network:

```yaml
# In your /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx
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

After placing this file in `/docker/` and ensuring `DOMAIN` is configured in `/etc/m3tal/.env`, run `m3tal up` to deploy and expose `my-app` at `http://app.DOMAIN`.

## Service Management

The M3TAL API daemon is managed as a systemd service, `m3tal-api.service`. You can interact with it using standard `systemctl` commands:

*   **Check status:** `systemctl status m3tal-api.service`
*   **Restart service:** `systemctl restart m3tal-api.service`
*   **View logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

The `m3tal` CLI provides specific commands for managing parts of the system:

*   **Start the M3TAL Dashboard:**
    ```bash
    m3tal dash up
    ```
    This command specifically starts the `m3tal-dashboard` container, ensuring the correct `DASHBOARD_EXPOSE_MODE` override is applied.

*   **Deploy all M3TAL stacks and user-defined services:**
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all Docker Compose stacks found within the `/docker/` directory, including the `routing` stack (Traefik, Cloudflared) and any user-defined compose files.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                    | Access                                        | Description                                                      |
| :--- | :------------------------- | :-------------------------------------------- | :--------------------------------------------------------------- |
| 80   | Traefik HTTP entry point   | Public                                        | The public-facing HTTP port for services exposed via Traefik.    |
| 8080 | M3TAL API daemon (Go)      | Host-local                                    | The internal port the M3TAL API daemon listens on.               |
| 8081 | Traefik dashboard          | Host-local only                               | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard            | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.