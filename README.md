# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system. It describes the core components, their interactions, and administrative practices required for installation, deployment, and ongoing management.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## APT Installation

To install the M3TAL CLI and API daemon via APT:

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

The M3TAL system is comprised of the following key components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the primary entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary operating as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask-based Docker container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services via domain names on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel Docker container, facilitating zero-configuration internet access for services.

## Filesystem Contract

The following table details the critical directories and files within the M3TAL system:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configurations. |
| `/docker`                   | A symlink alias to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth for all stack files and Traefik configurations. The `/docker` directory is a symlink to `/opt/m3tal/stack/`, providing a convenient user-facing alias for all stack management operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Execute `m3tal up` to deploy and start all defined stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for new installations).
*   **Mechanism:** Uses the `m3tal-compose.local.yml` override file, which adds a direct host port binding.
*   **Access:** The dashboard is accessible directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik gateway is required for this mode.
*   **Behavior:** A new user performing a default installation will access the dashboard directly via port `8082`.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override file, which applies Traefik labels to the dashboard container. Traefik then routes `dash.${DOMAIN}` to the dashboard container's internal port `8082`.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   **Requirements:** The Traefik gateway must be running via `m3tal up`.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed as a Docker container via `routing-compose.yml`. It binds host port `80` as its primary HTTP entry point and loads dynamic configuration from files within `/docker/dynamic/` (which is `/opt/m3tal/stack/dynamic/`).

**Dynamic Configuration Example:**

*   **API Daemon Routing:** Traefik routes `api.DOMAIN` to the Go API daemon listening on the host-local port `8080`. This is achieved through a dynamic configuration file, such as `dynamic/api.yml`, which defines a router and service pointing to `http://host.docker.internal:8080`.
*   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik labels on the dashboard container route `dash.DOMAIN` to the dashboard container.

**Example: Exposing a Custom User Service via Traefik**

To expose a hypothetical `my-app` service running Nginx at `app.DOMAIN`:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
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

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically configure routing for `app.DOMAIN` to your Nginx container.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Standard systemctl commands can be used for its administration:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

*   To specifically start the M3TAL Dashboard container, run: `m3tal dash up`. This command fetches the latest dashboard compose files and starts the dashboard with the appropriate expose mode based on your `/etc/m3tal/.env` configuration.
*   To orchestrate and deploy all other Docker Compose stacks present in the `/docker/` directory (including any user-defined compose files), execute: `m3tal up`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :--------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                         | The public-facing HTTP port for services exposed via Traefik.                                                                       |
| 8080 | M3TAL API daemon (Go)       | Host-local                                     | The internal port the M3TAL API daemon listens on.                                                                                  |
| 8081 | Traefik dashboard           | Host-local only                                | The internal Traefik dashboard port, accessible only from the host machine.                                                         |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`.                     |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.