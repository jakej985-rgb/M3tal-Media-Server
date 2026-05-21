# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

The M3TAL CLI binary is distributed via an APT repository. Follow these steps to install:

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

The M3TAL system is comprised of the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. This daemon manages Docker interactions, the M3TAL state database, and exposes the system's API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application packaged as a Docker container, listening internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Docker container for establishing Cloudflare tunnels, enabling zero-configuration internet access for services.

## Filesystem Contract

The following table details the canonical paths and their purposes within the M3TAL filesystem:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory for Docker Compose stack files and Traefik dynamic configuration. |
| `/docker`                   | Symlink pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. For user convenience and direct interaction, `/docker` is a symlink alias to `/opt/m3tal/stack/`. All user-facing stack operations, such as adding new compose files, should target the `/docker/` directory.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **Local Mode (`DASHBOARD_EXPOSE_MODE=local`)**
    *   This is the default configuration for a new M3TAL installation.
    *   The dashboard container's port `8082` is directly bound to the host machine's port `8082`.
    *   Access is achieved directly via `http://HOST_IP:8082` or `http://localhost:8082`.
    *   No Traefik gateway is required for this mode of access. A new user performing a default installation will access the dashboard directly via port 8082.
    *   This mode is suitable for LAN-only setups, initial testing, or environments without a domain-based reverse proxy.

2.  **Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)**
    *   When this mode is enabled, the dashboard container is configured with Traefik labels.
    *   Traefik routes traffic for the domain `dash.DOMAIN` to the dashboard container's internal port `8082`.
    *   Access is achieved via `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
    *   Requires the Traefik gateway to be running and correctly configured (`m3tal up` will deploy Traefik if not already running).
    *   This mode is ideal for domain-based setups and integrating the dashboard behind a central reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik utilizes dynamic configuration files (such as those found in `/docker/dynamic/`) to define routing for services that are not directly discoverable via Docker labels, or to provide more complex routing rules. For instance, the `dynamic/api.yml` file routes requests for `api.DOMAIN` to the Go API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container when `DASHBOARD_EXPOSE_MODE=traefik` is enabled, leveraging the labels defined in `m3tal-compose.traefik.yml` to target the dashboard container on its internal port `8082`.

### Exposing a Custom Service via Traefik

To expose a custom user service through Traefik, add the necessary labels to its service definition in your Docker Compose file (e.g., `my-app-compose.yml`):

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
This example exposes an Nginx container at `http://app.DOMAIN` through the `web` entrypoint (HTTP port 80) of Traefik, routing to the container's internal port 80. Ensure your service is connected to the `proxy` network.

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. Standard `systemctl` commands are used for its management:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get the M3TAL Dashboard up and running for initial interaction:

1.  Run `m3tal dash up`. This command specifically manages the `m3tal-dashboard` container, ensuring the latest compose files are used and the appropriate override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` based on `DASHBOARD_EXPOSE_MODE`) is applied.
2.  With `DASHBOARD_EXPOSE_MODE=local` (the default), access the dashboard at `http://HOST_IP:8082` after it starts.
3.  To deploy all other M3TAL stacks, including Traefik and any user-defined compose files you've placed in `/docker/`, execute `m3tal up`. This command orchestrates and deploys all active compose files within the `/docker/` directory.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :------ | :----- | :---------- |
| 80   | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.