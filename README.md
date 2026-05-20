# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

The M3TAL CLI binary and API daemon are installed via an APT repository.

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

The M3TAL system consists of the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary that serves as the single entrypoint for all M3TAL operations, installed via APT.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container, listening internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare Tunnel Docker container that provides zero-configuration internet access for services.

## Filesystem Contract

The following table details key file system paths and their purposes within the M3TAL system:

| Path                     | Purpose                                                         |
| :----------------------- | :-------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.   |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created by the API daemon.          |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.        |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all Docker Compose stack files reside. For user convenience, a symlink `/docker` is created, which aliases to `/opt/m3tal/stack/`. All user-facing stack operations, such as adding new compose files, should interact with the `/docker/` directory.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy or update all defined stacks, including your new one.

## Quick Demo

To quickly get the M3TAL Dashboard up and running:

*   **Start the M3TAL Dashboard:**
    ```bash
    m3tal dash up
    ```
    This command specifically manages the M3TAL Dashboard container. It downloads the necessary `m3tal-compose.yml` and its override files from GitHub, reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`, and starts the dashboard using the appropriate Docker Compose configuration.

*   **Deploy all M3TAL stacks and user-defined services:**
    ```bash
    m3tal up
    ```
    This command orchestrates and deploys all `*-compose.yml` files found in the `/docker/` directory. This includes core M3TAL components like Traefik and Cloudflared, as well as any custom Docker Compose files you have placed in `/docker/`.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`) - Default

*   This is the default access mode for a new installation.
*   The dashboard container's internal port `8082` is directly bound to the host machine's port `8082` via the `m3tal-compose.local.yml` override.
*   **Access Method:** You can access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik gateway is required for this mode. It works out-of-the-box for LAN-only setups, first-time users, and local testing environments. A new user performing a default installation will access the dashboard directly via port `8082`.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

*   When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, the `m3tal-compose.traefik.yml` override is applied. This adds Traefik labels to the dashboard service.
*   **Access Method:** The Traefik gateway routes traffic to the dashboard container based on these labels. You access the dashboard via a domain, typically `http://dash.DOMAIN`, where `DOMAIN` is configured in `/etc/m3tal/.env`.
*   **Requirements:** The Traefik gateway (deployed via `m3tal up`) must be running and correctly configured to route `dash.DOMAIN`.

## Traefik Gateway

Traefik is deployed as a Docker container via `routing-compose.yml` and acts as the primary reverse proxy for M3TAL. It binds host port `80` as its HTTP entry point.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also loads dynamic configuration files from `/docker/dynamic/` (symlinked from `/opt/m3tal/stack/dynamic`). This allows for flexible routing of requests, including those to services not running within the Docker network, such as the M3TAL Go API daemon. For instance, the `dynamic/api.yml` configuration routes requests for `api.DOMAIN` to the Go API daemon listening on host-local port `8080` by directing traffic to `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` routes to the M3TAL Dashboard container via Traefik labels.

### Example: Exposing a Custom User Service via Traefik

To expose a hypothetical `my-app` service via Traefik, you would add the necessary labels to its Docker Compose service definition. Place this `my-app-compose.yml` in your `/docker/` directory:

```yaml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    restart: unless-stopped
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
After placing this file and ensuring `DOMAIN` is set in your `.env`, run `m3tal up`. Your `my-app` service will then be accessible via `http://app.DOMAIN`.

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage its lifecycle and inspect its status using standard systemctl commands:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.