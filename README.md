# M3TAL System Documentation

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

## System Components

The M3TAL system comprises the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations, installed via APT.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application container running internally on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for establishing zero-configuration internet access to services.

## Filesystem Contract

The M3TAL system adheres to the following filesystem contract:

| Path                      | Purpose                                                          |
|---------------------------|------------------------------------------------------------------|
| `/etc/m3tal/.env`         | Primary configuration file. Managed by `m3tal config wizard`.    |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon.           |
| `/opt/m3tal/stack/`       | Canonical stack directory. Contains Docker Compose files and Traefik dynamic configuration. |
| `/docker`                 | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`      | Dashboard credential store. Managed by `m3tal dashpass`.         |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory serves as the canonical source of truth where all M3TAL-managed Docker Compose files and Traefik dynamic configuration files reside. The `/docker` directory is a user-facing symlink alias to `/opt/m3tal/stack/`, providing a convenient, consistent path for all stack-related operations.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (e.g., by using `m3tal config wizard` or `m3tal config set KEY value`).
3. Execute `m3tal up` to deploy all Docker Compose stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### 1. Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` (This is the default setting for a new installation).
*   **Mechanism:** This mode utilizes a direct port binding, exposing the dashboard container's internal port `8082` directly on the host machine. It does not require Traefik.
*   **Access:** Access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Note:** A new user performing a default M3TAL installation will access the dashboard directly via port `8082` on their host machine, as this behavior is linked to the default `DASHBOARD_EXPOSE_MODE=local` setting.

### 2. Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** In this mode, Traefik labels are applied to the dashboard container, allowing Traefik to route traffic to it. This requires the Traefik gateway to be running.
*   **Access:** Access the dashboard via a configured domain, typically `http://dash.DOMAIN`.
*   **Prerequisite:** Traefik must be deployed and operational via `m3tal up`.

## Traefik Gateway

Traefik serves as the M3TAL system's edge router, deployed as a Docker container via `routing-compose.yml`. Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds to host port `80` (and `443` if HTTPS is configured) as its primary HTTP (and HTTPS) entry point. It also loads dynamic configuration files from `/docker/dynamic/` (which symlinks to `/opt/m3tal/stack/dynamic/`), enabling hot-reloading of routing rules.

### Dynamic Configuration for Host-Local Services

Dynamic configuration files are used to route requests to services not managed directly by Traefik through Docker labels (e.g., host-local daemons). For instance, the `m3tal-api` daemon runs on the host machine at port `8080`. Traefik routes `api.DOMAIN` to this daemon using a file provider configuration like `dynamic/api.yml`:

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

Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, the `dash.DOMAIN` route is handled by Traefik labels on the `m3tal-dashboard` container, directing traffic to its internal port `8082`.

### Exposing a Custom User Service via Traefik

To expose your own Docker Compose service via Traefik, add the necessary Traefik labels to its service definition. For example, to expose a hypothetical `my-app` service on `app.DOMAIN`:

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
      - proxy # Crucial: Your service must be on the 'proxy' network
networks:
  proxy:
    external: true # Define the external 'proxy' network
```

After placing this file in `/docker/` and running `m3tal up`, Traefik will automatically discover and route requests for `app.DOMAIN` to your `my-app` container.

## Service Management

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service. You can manage its lifecycle using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api.service`
*   **Restart Service:** `systemctl restart m3tal-api.service`
*   **View Logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

To quickly start and access the M3TAL Dashboard:

1.  **Start the M3TAL Dashboard specifically:**
    ```bash
    m3tal dash up
    ```
    This command specifically downloads the latest dashboard Docker Compose files and starts the `m3tal-dashboard` container with the appropriate configuration based on your `DASHBOARD_EXPOSE_MODE` setting.

2.  **Access the Dashboard:**
    *   **If `DASHBOARD_EXPOSE_MODE=local` (default):** Open your web browser and navigate to `http://HOST_IP:8082` (replace `HOST_IP` with your machine's IP address or `localhost`).
    *   **If `DASHBOARD_EXPOSE_MODE=traefik`:** Open your web browser and navigate to `http://dash.DOMAIN` (replace `DOMAIN` with your configured domain).

To deploy all other M3TAL stacks and any user-defined Docker Compose files:

```bash
m3tal up
```
This command orchestrates and deploys *all* Docker Compose stacks found within the `/docker/` directory, including the Traefik gateway, Cloudflared (if configured), and any custom `*-compose.yml` files you have added.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.