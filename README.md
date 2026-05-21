# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

### APT Installation

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The M3TAL system utilizes the following critical filesystem paths:

| Path                       | Purpose                                                     |
|----------------------------|-------------------------------------------------------------|
| `/etc/m3tal/.env`          | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`  | SQLite state database. Auto-created by the API daemon.    |
| `/opt/m3tal/stack/`        | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                  | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`       | Dashboard credential store. Managed by `m3tal dashpass`.    |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker/` directory is a user-facing symlink alias for all stack operations, pointing to the canonical source of truth for stack files located at `/opt/m3tal/stack/`. When `m3tal up` is executed, it targets all `*-compose.yml` files within `/docker/`, ensuring that all defined stacks, including any custom user-added ones, are deployed.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any necessary environment variables required by your custom stack are defined in `/etc/m3tal/.env` or managed through `m3tal config set KEY value`.
3.  Execute `m3tal up` to deploy your new stack alongside the M3TAL core services.

## Dashboard Access

The M3TAL dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable within `/etc/m3tal/.env`.

### Local Mode (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Access:** Direct port binding at `http://HOST_IP:8082`. No Traefik is required for this mode.
*   **Description:** When a new user performs a default installation, they will access the dashboard directly via port 8082. This is the default behavior due to the `DASHBOARD_EXPOSE_MODE=local` setting, which maps the internal dashboard port to the host's IP address on port 8082. This mode is ideal for LAN-only setups, local testing, or initial deployments where a public domain and reverse proxy are not yet configured.

### Traefik Mode

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Access:** Domain routing at `http://dash.DOMAIN` via Traefik. Traefik must be running and configured to manage this domain.
*   **Description:** In this mode, Traefik is responsible for routing traffic to the dashboard. Access is established via the domain name specified in `DOMAIN` environment variable, prefixed with `dash.`. This mode is suitable for deployments that utilize a domain-based routing strategy and have Traefik configured as the reverse proxy.

## Traefik Gateway

Traefik functions as the primary reverse proxy for the M3TAL system, automatically discovering and routing traffic to Docker services based on Traefik labels defined within their respective Docker Compose service definitions. Services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik leverages a file provider mechanism to load dynamic configuration files from the `/docker/dynamic/` directory. This allows for hot-reloading of routing configurations without restarting Traefik. For instance, a dynamic configuration file like `dynamic/api.yml` is used to route requests from `api.DOMAIN` to the Go API daemon listening on the host-local port `8080` via `http://host.docker.internal:8080`.

The dashboard itself is routed via Traefik when `DASHBOARD_EXPOSE_MODE=traefik` is set. This is configured via Traefik labels within the dashboard's Compose file, directing traffic for `dash.DOMAIN` to the dashboard container.

To expose a custom user service via Traefik labels, you would typically include them in your custom Docker Compose file:

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
```

## Service Management

The M3TAL API daemon is managed by systemd. The following `systemctl` commands can be used to control its lifecycle:

*   **Check status:** `systemctl status m3tal-api.service`
*   **Restart service:** `systemctl restart m3tal-api.service`
*   **View logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

To get started quickly, you can launch the dashboard container specifically:

*   **Start the dashboard:** `m3tal dash up`
    This command downloads the latest dashboard Compose definitions, checks your `/etc/m3tal/.env` for `DASHBOARD_EXPOSE_MODE`, and starts the dashboard container accordingly.

To deploy all other stacks, including any custom Docker Compose files you have added to the `/docker/` directory:

*   **Start all stacks:** `m3tal up`
    This command orchestrates and deploys all services defined in `*-compose.yml` files within the `/docker/` directory.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                               | Description                                                                         |
|------|-----------------------------|--------------------------------------|-------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point    | Public (traefik mode)                | The public-facing HTTP port for services exposed via Traefik.                       |
| 8080 | M3TAL API daemon (Go)       | Host-local                           | The internal port the M3TAL API daemon listens on.                                  |
| 8081 | Traefik dashboard           | Host-local only                      | The internal Traefik dashboard port, accessible only from the host machine.         |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.