# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install M3TAL, follow these steps:

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

The M3TAL system is comprised of several interacting components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare Tunnel container for establishing secure, zero-config internet access to services.

## Filesystem Contract

The following paths represent the canonical filesystem structure and their respective purposes within the M3TAL system:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for M3TAL, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database, automatically created by the API daemon.         |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store, managed by `m3tal dashpass`.               |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack. It is important to note that `/opt/m3tal/stack/` is the canonical source of truth directory where all stack files reside, and `/docker` is the user-facing symlink alias for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are set in `/etc/m3tal/.env` (e.g., using `m3tal config wizard` or `m3tal config set KEY value`).
3.  Execute `m3tal up` to start all deployed stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   This is the default configuration for a new M3TAL installation.
*   The dashboard container's port is directly bound to the host machine's network.
*   Access the dashboard via `http://HOST_IP:8082` or `http://localhost:8082`.
*   This mode does **not** require Traefik to be running for dashboard access. A new user performing a default installation will access the dashboard directly via port 8082.
*   **Best for:** LAN-only setups, first-time users, and local development or testing.

### Traefik Mode

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   In this mode, the dashboard container is configured with Traefik labels to enable routing via the Traefik gateway.
*   Access the dashboard via `http://dash.DOMAIN` (where `DOMAIN` is configured in `/etc/m3tal/.env`, defaulting to `localhost`).
*   This mode **requires** Traefik to be running via `m3tal up` to correctly route traffic to the dashboard.
*   **Best for:** Domain-based setups and environments where multiple services are exposed behind a central reverse proxy.

## Traefik Gateway

Traefik serves as the central reverse proxy for M3TAL, automatically discovering and routing traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml` and binds to host port `80` (HTTP entry point). It also uses a file provider to load dynamic configurations from `/docker/dynamic/`, allowing for hot-reloading of routing rules.

### Dynamic Routing Examples

*   **M3TAL API Daemon:** Traefik routes `api.DOMAIN` to the M3TAL Go API daemon, which listens on host-local port `8080`. This is configured via a dynamic configuration file, such as `/docker/dynamic/api.yml`, which translates `http://api.DOMAIN` requests to `http://host.docker.internal:8080`.
*   **M3TAL Dashboard:** When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.DOMAIN` to the `m3tal-dashboard` container, targeting its internal port `8082`. This routing is configured via Traefik labels within `m3tal-compose.traefik.yml`.

### Exposing a Custom Service via Traefik

To expose your own Docker Compose service through Traefik, add the necessary labels to its service definition. Ensure your service is on the `proxy` network, which Traefik also utilizes.

Example `my-app-compose.yml`:

```yaml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app-container
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

After placing this file in `/docker/` and running `m3tal up`, `http://app.DOMAIN` will route to your `my-app` container.

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

The `m3tal` CLI provides specific commands for managing different parts of the system:

*   **Start the Dashboard:** To specifically bring up only the M3TAL Dashboard container (useful for initial setup or troubleshooting), use:
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard Docker Compose files and starts the dashboard with the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.

*   **Deploy All Stacks:** To orchestrate and deploy all other Docker Compose stacks present in the `/docker/` directory, including any user-defined compose files, use:
    ```bash
    m3tal up
    ```
    This command will bring up the Traefik gateway, Cloudflared (if configured), and any other services defined in `*-compose.yml` files.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.