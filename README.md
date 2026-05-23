# M3TAL System Documentation

## Overview

This document provides technical details and operational procedures for the M3TAL system. It describes the system's architecture, component interactions, deployment mechanisms, and operational guidelines.

## Prerequisites

**Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.**

## Installation

The M3TAL CLI and API daemon are installed via a dedicated APT repository.

1.  **Add the GPG signing key**
    ```bash
    curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg
    ```

2.  **Add the APT repository**
    ```bash
    echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list
    ```

3.  **Install**
    ```bash
    sudo apt update && sudo apt install -y m3tal
    ```

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Components

The M3TAL system comprises several interconnected components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container designed to expose services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, facilitating zero-configuration internet access for exposed services.

## Filesystem Contract

The following table details critical directories and files within the M3TAL system:

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.             |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.           |

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service               | Access                                        | Description                                                                                                                                                                          |
| :--- | :-------------------- | :-------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point | Public                                        | The public-facing HTTP port for services exposed via Traefik.                                                                                                                        |
| 8080 | M3TAL API daemon (Go)    | Host-local                                    | The internal port the M3TAL API daemon listens on.                                                                                                                                   |
| 8081 | Traefik dashboard     | Host-local only                               | The internal Traefik dashboard port, accessible only from the host machine.                                                                                                          |
| 8082 | M3TAL Dashboard       | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth for all stack files. For user convenience, `/docker` is a symlink alias to `/opt/m3tal/stack/`, making `/docker` the user-facing directory for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Behavior:** This mode uses an override (`m3tal-compose.local.yml`) to add a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the host machine.
*   **Access:** The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **New User Experience:** A new user performing a default installation will access the dashboard directly via port 8082, as this is the behavior linked to the default `DASHBOARD_EXPOSE_MODE=local` setting.
*   **Requirements:** No Traefik or domain configuration is required. Ideal for LAN-only setups, initial configurations, and local testing.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Behavior:** This mode utilizes an override (`m3tal-compose.traefik.yml`) to apply Traefik labels to the dashboard container. Traefik then routes incoming requests for `dash.DOMAIN` to the dashboard container's internal port `8082`.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN`.
*   **Requirements:** Traefik must be running (typically via `m3tal up` which includes `routing-compose.yml`) and configured for your domain. Best for domain-based setups where multiple services are managed by a reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration from files located in `/docker/dynamic/` (which is linked from `/opt/m3tal/stack/dynamic/`). This file provider allows for hot-reloading of routing rules. For instance, the M3TAL API daemon, which listens on host-local port `8080`, is exposed via Traefik through a dynamic configuration file (e.g., `dynamic/api.yml`). This file contains rules to route `api.DOMAIN` to `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` routes to the dashboard container.

### Example: Exposing a Custom User Service via Traefik

To expose a hypothetical `my-app` service running in a custom Docker Compose stack (`my-app-compose.yml`) via Traefik, you would add the following labels to its service definition:

```yaml
# /docker/my-app-compose.yml
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

In this example:
*   `traefik.enable=true` makes the service discoverable by Traefik.
*   `traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)` defines a routing rule that matches requests for `app.YOUR_DOMAIN`.
*   `traefik.http.services.myapp.loadbalancer.server.port=80` specifies that Traefik should forward traffic to port `80` of the `my-app` container.
*   `traefik.http.routers.myapp.entrypoints=web` links this router to Traefik's `web` entrypoint (typically on host port 80).
*   The `proxy` network is necessary for Traefik and your service to communicate.

## Service Management

The M3TAL API daemon operates as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard systemctl commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **Restart the service:**
    ```bash
    systemctl restart m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```

## Quick Demo

To quickly get the M3TAL Dashboard up and running:

1.  **Start the M3TAL API daemon (if not already running):**
    ```bash
    sudo systemctl start m3tal-api
    ```
2.  **Deploy only the M3TAL Dashboard container:**
    ```bash
    m3tal dash up
    ```
    This command specifically downloads the latest dashboard compose files and starts the dashboard container, respecting the `DASHBOARD_EXPOSE_MODE` setting.
3.  **Access the Dashboard:**
    *   If `DASHBOARD_EXPOSE_MODE=local` (the default), navigate to `http://HOST_IP:8082` in your web browser.
    *   If `DASHBOARD_EXPOSE_MODE=traefik` and Traefik is running, navigate to `http://dash.DOMAIN`.

To deploy all other configured M3TAL stacks and any user-defined compose files located in `/docker/`:

*   **Deploy all stacks:**
    ```bash
    m3tal up
    ```
    This command will orchestrate and deploy all `*-compose.yml` files found in the `/docker/` directory, including core M3TAL components and any custom services you have added.