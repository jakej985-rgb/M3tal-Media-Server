# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

---

## Prerequisites

**Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.**

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Installation

To install the M3TAL CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-apt-key/public.key | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-apt-key stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## System Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and exposes the M3TAL API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It utilizes a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, enabling zero-configuration internet access for exposed services.

## Filesystem Contract

The following table details the primary filesystem paths utilized by the M3TAL system:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file, managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`   | SQLite state database, automatically created by the API daemon.        |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik config. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store, managed by `m3tal dashpass`.               |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth where all Docker Compose stack files and associated configurations reside. For user convenience, a symlink `/docker` is created, which aliases `/opt/m3tal/stack/`. All user-facing stack operations, such as adding new stacks, should target the `/docker` directory.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard supports two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism:** M3TAL utilizes an override (`m3tal-compose.local.yml`) which adds a direct port binding, typically `${DASHBOARD_PORT:-8082}:8082`.
*   **Access:** A new user performing a default installation will access the dashboard directly via port 8082. This means it is accessible at `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements:** No Traefik configuration or external domain is required. This mode functions out-of-the-box for LAN-only setups or initial testing.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** M3TAL employs an override (`m3tal-compose.traefik.yml`) that injects Traefik labels into the dashboard container's definition. These labels instruct Traefik to route requests for `dash.${DOMAIN}` to the dashboard container's internal port `8082`.
*   **Access:** Accessible at `http://dash.DOMAIN`. This mode requires Traefik to be running (typically via `m3tal up`) and correctly configured to handle the specified domain.
*   **Requirements:** Suitable for domain-based deployments and environments where multiple services are managed behind a single reverse proxy.

## Traefik Gateway

Traefik serves as the M3TAL system's reverse proxy, automatically discovering and routing traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed as a container via `routing-compose.yml` and binds host port `80` as its primary HTTP entry point. It loads dynamic configuration files from `/docker/dynamic/` (which supports hot-reloading changes).

*   **API Daemon Routing:** Traefik routes `api.DOMAIN` to the M3TAL Go API daemon, which listens on the host-local port `8080`. This is achieved through a dynamic configuration file (e.g., `dynamic/api.yml`) that directs requests to `http://host.docker.docker.internal:8080`.
*   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.DOMAIN` to the M3TAL Dashboard container, which listens internally on port `8082`.

### Example: Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service, `my-app`, via Traefik:

Place your `my-app-compose.yml` in the `/docker/` directory with the following Traefik labels:

```yaml
# /docker/my-app-compose.yml
version: '3.8'

services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true" # Enable Traefik for this service
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN:-localhost}`)" # Route based on host header, e.g., app.example.com
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Target port within the container
      - "traefik.http.routers.myapp.entrypoints=web" # Use the 'web' entrypoint (HTTP port 80)
      - "traefik.docker.network=proxy" # Ensure service is on the Traefik-managed 'proxy' network
    networks:
      - proxy

networks:
  proxy:
    external: true # Use the existing 'proxy' network created by Traefik
```

After placing this file, run `m3tal up` to deploy the service and enable Traefik routing.

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage it using standard `systemctl` commands:

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get started with the M3TAL Dashboard:

*   To start *only* the dashboard container, including fetching the latest compose files and applying the correct `DASHBOARD_EXPOSE_MODE` override, execute:
    ```bash
    m3tal dash up
    ```
    If `DASHBOARD_EXPOSE_MODE` is set to `local` (the default), the dashboard will be accessible at `http://HOST_IP:8082`. If set to `traefik`, it will be available via `http://dash.DOMAIN` (assuming Traefik is also running).

*   To orchestrate and deploy *all* M3TAL-managed Docker Compose stacks, including any user-defined compose files located in `/docker/`, execute:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                             |
| :--- | :-------------------------- | :------------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.