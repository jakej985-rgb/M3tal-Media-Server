# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system, encompassing its core components, deployment mechanisms, and service management.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install M3TAL, execute the following commands:

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

The M3TAL system comprises the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on host port 8080. It manages Docker interactions, the internal state database, and exposes various API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container operating internally on port 8082. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container designed to expose Docker services by domain name on host port 80. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, enabling zero-configuration internet access for exposed services.

## Filesystem Contract

The following table outlines the critical directories and files within the M3TAL filesystem:

| Path                       | Purpose                                                                                             |
| :------------------------- | :-------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary configuration file for environment variables. Managed by `m3tal config wizard`.             |
| `/var/lib/m3tal/state.db`  | SQLite database for M3TAL's internal state. Automatically created and managed by the API daemon.    |
| `/opt/m3tal/stack/`        | The canonical directory containing all Docker Compose stack files and Traefik configuration.          |
| `/docker`                  | A symbolic link to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations.      |
| `/docker/users.json`       | The credential store for the M3TAL Dashboard. Managed by `m3tal dashpass`.                        |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all Docker Compose stack files reside. The `/docker/` directory is a symlink alias to `/opt/m3tal/stack/`, serving as the user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**: This is the default setting for new installations.
*   The dashboard container is exposed via a direct port binding, typically `8082:8082` on the host machine.
*   Access the dashboard directly at `http://HOST_IP:8082` or `http://localhost:8082`.
*   A new user performing a default installation will access the dashboard directly via port 8082.
*   No Traefik instance is required for this mode.

### Traefik Mode

*   **`DASHBOARD_EXPOSE_MODE=traefik`**:
*   The dashboard container is configured with Traefik labels to enable routing via the Traefik gateway.
*   Access the dashboard at `http://dash.DOMAIN` (e.g., `http://dash.example.com`).
*   Requires the Traefik gateway (`routing-compose.yml`) to be running via `m3tal up`. Traefik will route traffic to the dashboard container's internal port 8082.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml` and binds host port 80 as its HTTP entry point. It also processes dynamic configuration files located in `/docker/dynamic/` for additional routing rules. For instance, `dynamic/api.yml` defines a rule to route `api.DOMAIN` to the M3TAL Go API daemon, which listens on the host-local port `8080`, specifically `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes requests for `dash.DOMAIN` to the dashboard container on its internal port 8082 based on its Docker labels.

### Exposing a Custom Service via Traefik

To expose a custom user service through the Traefik gateway, include the necessary Traefik labels in your Docker Compose service definition:

```yaml
# In /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy # Ensure your service is on the 'proxy' network for Traefik discovery

networks:
  proxy:
    external: true
```

After placing this file in `/docker/`, run `m3tal up` to deploy the service and configure Traefik.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Standard `systemctl` commands can be used for its operation:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

The `m3tal` CLI provides specific commands for managing parts of the system:

*   To start only the M3TAL Dashboard container, use `m3tal dash up`. This command will automatically download the latest dashboard compose files and apply the correct override based on your `DASHBOARD_EXPOSE_MODE` setting.
*   To orchestrate and deploy all other stacks, including any user-defined compose files placed in the `/docker/` directory, use `m3tal up`. This command processes all `*-compose.yml` files found.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                             |
| :--- | :-------------------------- | :------------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.