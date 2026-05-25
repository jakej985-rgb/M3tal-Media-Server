# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## APT Installation

To install M3TAL via APT:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The M3TAL system adheres to the following filesystem structure:

| Path                        | Purpose                                                                |
| :-------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.                 |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config.  |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.               |

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all stack files reside. The `/docker` path is a user-facing symlink alias for all stack operations, allowing convenient interaction with Docker Compose files.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to start all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard provides a web-based interface for system management. It supports two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Access Method**: `http://HOST_IP:8082` or `http://localhost:8082`
*   **Details**: In this mode, the dashboard container binds directly to host port `8082`. This is the default behavior for a new installation, allowing immediate direct access via port 8082 without requiring Traefik or domain configuration. This mode is suitable for LAN-only setups, first-time users, and local testing environments.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Access Method**: `http://dash.DOMAIN`
*   **Details**: When configured for Traefik mode, the dashboard container is exposed via Traefik labels. This requires the Traefik Gateway to be running (typically via `m3tal up`). Traefik routes requests for `http://dash.DOMAIN` to the dashboard container internally. This mode is best for domain-based setups and integrating multiple services behind a central reverse proxy.

## Traefik Gateway

Traefik acts as the reverse proxy for M3TAL, managing external access to Docker services. Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik dynamically loads configuration from the `/docker/dynamic/` directory (e.g., `dynamic/api.yml`). This dynamic configuration is used to route requests to services that may not be Docker containers themselves but are accessible from the host, such as the M3TAL Go API daemon. For instance, `api.DOMAIN` is routed to the M3TAL API daemon listening on the host-local port `8080` by directing traffic to `http://host.docker.internal:8080`. Similarly, if `DASHBOARD_EXPOSE_MODE` is set to `traefik`, requests to `dash.DOMAIN` are routed to the M3TAL Dashboard container via its internal port `8082` using Traefik labels defined in its Docker Compose override file.

### Exposing a Custom Service via Traefik

To expose a custom user service through Traefik, add the necessary Traefik labels to its Docker Compose service definition. For example, to expose a hypothetical `my-app` service at `app.DOMAIN`:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-nginx-app
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN:-localhost}`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
      - "traefik.docker.network=proxy" # Ensure your service is on the 'proxy' network
    networks:
      - proxy

networks:
  proxy:
    external: true # Assumes the 'proxy' network is external and managed by M3TAL
```

## Service Management

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service, providing core system functionality and Docker orchestration.

To manage the M3TAL API daemon:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Quick Demo

To quickly get started with the M3TAL dashboard:

*   **Start the M3TAL Dashboard container**: Run `m3tal dash up`. This command specifically downloads and starts the `m3tal-dashboard` container, applying the correct Docker Compose override based on your `DASHBOARD_EXPOSE_MODE` setting.
*   **Start all other M3TAL stacks**: Execute `m3tal up`. This command orchestrates and deploys all `*-compose.yml` files located in the `/docker/` directory, including the Traefik Gateway, Cloudflared, and any user-defined compose files you may have added.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :--------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                   | The public-facing HTTP port for services exposed via Traefik.                                                                                                          |
| 8080 | M3TAL API daemon (Go)       | Host-local                               | The internal port the M3TAL API daemon listens on.                                                                                                                     |
| 8081 | Traefik dashboard           | Host-local only                          | The internal Traefik dashboard port, accessible only from the host machine.                                                                                            |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`.                                                       |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.