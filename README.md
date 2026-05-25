# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

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

The M3TAL system is composed of several interacting components:

*   **CLI binary** (`/usr/bin/m3tal`): A Go binary installed via APT, serving as the unified entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the internal state database, and exposes RESTful API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running within a Docker container, listening internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It dynamically discovers Docker services via labels and loads additional routing configurations from a file provider.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing zero-configuration internet access for services.

## Filesystem Contract

The M3TAL system relies on a specific filesystem layout:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.   |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose stack files and Traefik dynamic configuration. |
| `/docker`                   | Symlink pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. For user convenience, `/docker` is a symlink alias to `/opt/m3tal/stack/`, making it the user-facing directory for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to deploy your new stack along with all other M3TAL-managed services.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

**1. Local Mode (Default)**
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism**: A direct port binding (`${DASHBOARD_PORT:-8082}:8082`) is added to the `m3tal-dashboard` container. No Traefik configuration is involved.
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Note for new users**: A new user performing a default M3TAL installation will access the dashboard directly via host port `8082`.

**2. Traefik Mode**
*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism**: Traefik labels are applied to the `m3tal-dashboard` container, allowing Traefik to route traffic from a specified domain to the dashboard's internal port `8082`. This requires the Traefik gateway to be running.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN`, where `DOMAIN` is configured in `/etc/m3tal/.env`.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik also utilizes dynamic configuration files (e.g., `dynamic/api.yml` located in `/docker/dynamic/`) to route requests to services that are not running as Docker containers or require specific routing rules. For instance, `api.DOMAIN` is routed to the M3TAL Go API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`. Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, `dash.DOMAIN` is routed to the dashboard container based on its Traefik labels.

### Exposing a Custom User Service via Traefik

To expose a new service, add the following labels to your service definition in its Docker Compose file (e.g., `my-app-compose.yml` in `/docker/`):

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
      - proxy # Crucial: Your service must be on the 'proxy' network for Traefik to discover it.
```

After updating your Docker Compose file, run `m3tal up` to apply the changes.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`.

*   **Check service status**: `systemctl status m3tal-api`
*   **Restart the API daemon**: `systemctl restart m3tal-api`
*   **View service logs**: `journalctl -u m3tal-api -f`

## Quick Demo

*   To start only the M3TAL Dashboard container specifically, without deploying other general-purpose stacks:
    ```bash
    m3tal dash up
    ```
    This command downloads the latest dashboard compose files and starts the dashboard with the appropriate override based on `DASHBOARD_EXPOSE_MODE`.
*   To orchestrate and deploy all other stacks in the `/docker/` directory, including any user-defined Docker Compose files:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :------------------------------------------- | :-------------------------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.