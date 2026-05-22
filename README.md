# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

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

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container running internally on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes Docker services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, facilitating zero-configuration internet access for services.

## Filesystem Contract

The following table details key directories and files within the M3TAL filesystem:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database for the M3TAL API daemon. Auto-created.          |
| `/opt/m3tal/stack/`         | Canonical directory containing all Docker Compose stack files and Traefik configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |

## Configuration

M3TAL uses environment variables for configuration, primarily stored in `/etc/m3tal/.env`. The `m3tal config wizard` command provides an interactive way to manage these settings, or they can be set directly via `m3tal config set KEY value`.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. The `/docker` directory is a symlink alias to `/opt/m3tal/stack/`, serving as the user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (e.g., by using `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard supports two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Setting:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** This mode uses an override file (`m3tal-compose.local.yml`) to add a direct port binding, exposing the dashboard container's internal port `8082` to the host.
*   **Access:** The M3TAL Dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik gateway is required for this mode.
*   **Note:** A new user performing a default M3TAL installation will access the dashboard directly via port `8082` because `DASHBOARD_EXPOSE_MODE` defaults to `local`.

### Traefik Mode

*   **Setting:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** This mode utilizes an override file (`m3tal-compose.traefik.yml`) that adds specific Traefik labels to the dashboard service. Traefik, if running, will interpret these labels to route traffic for `dash.DOMAIN` to the dashboard container on its internal port `8082`.
*   **Access:** The M3TAL Dashboard is accessible via `http://dash.DOMAIN`.
*   **Requirements:** The Traefik gateway must be running (`m3tal up` typically starts it as part of the `routing` stack).

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds to host port `80` (and `443` if HTTPS is configured) as its primary HTTP entry point. It loads dynamic configuration from files located in `/docker/dynamic/` (which is `/opt/m3tal/stack/dynamic` on the canonical path), allowing for hot-reloading of routing rules.

For example, the M3TAL API daemon runs on a host-local port `8080`. Traefik routes requests for `api.DOMAIN` to this host-local service using a dynamic configuration file like `dynamic/api.yml`:

```yaml
# /docker/dynamic/api.yml
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

Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard container (listening internally on port `8082`) is exposed via `dash.DOMAIN` through Traefik labels defined in its compose override file.

### Exposing a Custom Service via Traefik

To expose a custom user service via Traefik, add the necessary Traefik labels to its Docker Compose service definition. Ensure your service is on the `proxy` network, which Traefik monitors.

Example `my-app-compose.yml` snippet:

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

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api.service`
*   **Restart Service:** `systemctl restart m3tal-api.service`
*   **View Logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

The `m3tal` CLI provides convenient commands for managing your Docker stacks:

*   **Start the Dashboard:** To specifically start and manage the `m3tal-dashboard` container, use:
    ```bash
    m3tal dash up
    ```
    This command will download the necessary compose files, apply the correct override based on your `DASHBOARD_EXPOSE_MODE` setting, and bring up the dashboard.
*   **Deploy All Stacks:** To orchestrate and deploy all Docker Compose stacks located in the `/docker/` directory (including any user-defined compose files and M3TAL's core services like Traefik and the Dashboard, if not managed by `m3tal dash up`), use:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :---------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public                                    | The public-facing HTTP port for services exposed via Traefik.                                                                                                                                     |
| 8080 | M3TAL API daemon (Go)       | Host-local                                | The internal port the M3TAL API daemon listens on.                                                                                                                                                |
| 8081 | Traefik dashboard           | Host-local only                           | The internal Traefik dashboard port, accessible only from the host machine.                                                                                                                       |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`.                                                                                    |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.