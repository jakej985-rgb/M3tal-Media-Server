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

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Overview of Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary serving as the single entrypoint for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It manages Docker interactions, the state database, and exposes the M3TAL API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask Docker container operating internally on port `8082`. It communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services via domain names on host port `80`. It leverages a file provider for dynamic routing configurations.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel Docker container for establishing zero-configuration internet access to services.

## Filesystem Contract

The following table details the key paths and their purposes within the M3TAL filesystem:

| Path                        | Purpose                                                                 |
| :-------------------------- | :---------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.           |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the M3TAL API daemon.            |
| `/opt/m3tal/stack/`         | Canonical directory for Docker Compose stack files and Traefik config.  |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all stack files reside. For user convenience, `/docker` is a symlink alias to `/opt/m3tal/stack/`, making `/docker` the primary user-facing path for all stack operations.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Execute `m3tal up` to start all deployed stacks, including your newly added one.

## Quick Demo

To quickly start the M3TAL Dashboard container specifically, use the command:

```bash
m3tal dash up
```

This command downloads the necessary dashboard compose files and starts the dashboard according to your `DASHBOARD_EXPOSE_MODE` setting.

In contrast, `m3tal up` orchestrates and deploys all other stacks present in the `/docker/` directory, including any user-defined compose files you have placed there.

## Dashboard Access

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local`
*   **Behavior**: This mode utilizes the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`). No Traefik configuration is involved.
*   **Access**: The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Clarification for New Users**: A new user performing a default M3TAL installation will access the dashboard directly via port `8082`. This behavior is a direct result of the default `DASHBOARD_EXPOSE_MODE=local` setting.
*   **Use Case**: Ideal for LAN-only setups, first-time users, and local testing environments.

### Traefik Mode

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik`
*   **Behavior**: This mode uses the `m3tal-compose.traefik.yml` override, which applies Traefik labels to the dashboard container. Traefik then routes incoming requests for `dash.DOMAIN` to the dashboard container's internal port `8082`. This requires Traefik to be running via `m3tal up`.
*   **Access**: The dashboard is accessible via `http://dash.DOMAIN`.
*   **Use Case**: Suited for domain-based setups and environments where multiple services are managed behind a reverse proxy.

## Traefik Gateway

Traefik is deployed as a Docker container via `routing-compose.yml` and acts as the reverse proxy for M3TAL services. It binds host port `80` as its primary HTTP entry point.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable. Traefik also loads dynamic configuration from files placed within `/docker/dynamic/` (symlinked from `/opt/m3tal/stack/dynamic/`), enabling hot-reloading of routing rules.

### Dynamic Routing for Host-Local Services

Traefik's file provider is utilized for routing requests to services listening on host-local ports, such as the M3TAL API daemon. For instance, `dynamic/api.yml` routes `api.DOMAIN` to the Go API daemon on host-local port `8080` by directing traffic to `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the M3TAL Dashboard container (when `DASHBOARD_EXPOSE_MODE=traefik`).

Example dynamic configuration for the M3TAL API (`dynamic/api.yml`):

```yaml
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

### Exposing Custom Services via Traefik

To expose a custom user service via Traefik, you need to add appropriate Traefik labels to its service definition in your Docker Compose file (e.g., `my-app-compose.yml`). The service must also be part of the `proxy` network for Traefik to discover it.

Here is a concrete YAML example for exposing a hypothetical `my-app` service:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.${DOMAIN}`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy
    restart: unless-stopped

networks:
  proxy:
    external: true
```

After placing this file in `/docker/` and running `m3tal up`, your `my-app` service would be accessible via `http://app.DOMAIN`.

## Service Management

The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. You can control and monitor its state using standard `systemctl` commands:

*   **Check status**: `systemctl status m3tal-api`
*   **Restart service**: `systemctl restart m3tal-api`
*   **View logs**: `journalctl -u m3tal-api -f`

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                                                                                                                                                                                        |
| :--- | :-------------------------- | :------------------------------------------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                                                                                                                                                                                      |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                                                                                                                                                                                 |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine (e.g., `http://localhost:8081`).                                                                                                                                                        |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. In `local` mode, it's `http://HOST_IP:8082`; in `traefik` mode, it's routed via `http://dash.DOMAIN`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.