# M3TAL System Documentation

## Overview
This document provides technical details and operational procedures for the M3TAL system. M3TAL consists of a Go CLI binary, a Go API daemon, a Python/Flask dashboard container, a Traefik reverse proxy, and an optional Cloudflared tunnel for secure internet access. These components interact to provide system management and service orchestration capabilities.

## Prerequisites
**Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.**

## APT Installation

To install the M3TAL CLI binary and systemd service:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The following table details key file system paths and their purposes within the M3TAL system:

| Path                        | Purpose                                                      |
|-----------------------------|--------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.|
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.       |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik configuration. |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.     |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all Docker Compose stack files reside, including those for core M3TAL components and routing. For user convenience and direct interaction, the `/docker` directory serves as a symlink alias to `/opt/m3tal/stack/`. All user-facing stack operations, such as adding new services, should target the `/docker/` directory.

### Adding a New Stack
To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy your new stack along with all other configured services.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

**1. Local Mode (Default: `DASHBOARD_EXPOSE_MODE=local`)**
This is the default configuration for new installations.
*   A direct port binding is configured, typically exposing the dashboard at `http://HOST_IP:8082` or `http://localhost:8082`.
*   A new user performing a default installation will access the dashboard directly via port 8082, as this behavior is linked to the `DASHBOARD_EXPOSE_MODE=local` setting.
*   This mode does not require Traefik or domain configuration.
*   **Access:** `http://YOUR_SERVER_IP:8082`

**2. Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)**
This mode leverages Traefik for domain-based routing.
*   The dashboard container is configured with Traefik labels, allowing Traefik to route requests for `dash.DOMAIN` to the dashboard service.
*   Requires Traefik to be running (typically via `m3tal up`).
*   **Access:** `http://dash.DOMAIN` (e.g., `http://dash.example.com` if `DOMAIN=example.com`).

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik handles dynamic configuration by loading YAML files from `/docker/dynamic/` (which symlinks to `/opt/m3tal/stack/dynamic/`). For example, `dynamic/api.yml` is used to route `api.DOMAIN` to the M3TAL Go API daemon, which listens on the host-local port `8080`. This routing is achieved by configuring Traefik to forward requests to `http://host.docker.internal:8080`. Similarly, `dash.DOMAIN` routes to the dashboard container via Traefik labels when `DASHBOARD_EXPOSE_MODE=traefik`.

### Exposing a Custom User Service via Traefik
To expose a custom Docker Compose service (e.g., `my-app`) via Traefik, add the following labels to its service definition in your `my-app-compose.yml` file:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-app
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
    networks:
      - proxy # Ensure your service is on the 'proxy' network if Traefik is
              # configured to use it (default M3TAL setup).

networks:
  proxy:
    external: true # Assumes 'proxy' network is created externally by M3TAL
```
After modifying or adding such a file, run `m3tal up` to deploy the changes.

## Firewall Considerations
If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage it using standard `systemctl` commands:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

*   To start only the M3TAL Dashboard container, run:
    ```bash
    m3tal dash up
    ```
    This command specifically manages the dashboard container, downloading necessary compose files and starting it with the appropriate `DASHBOARD_EXPOSE_MODE` override.

*   To orchestrate and deploy all other stacks configured in the `/docker/` directory, including any user-defined compose files, run:
    ```bash
    m3tal up
    ```

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|-----------------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.