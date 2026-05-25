# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

| Path                      | Purpose                                       |
| :------------------------ | :-------------------------------------------- |
| `/etc/m3tal/.env`         | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`       | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                 | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`      | Dashboard credential store. Managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth for all stack definition files and Traefik configuration. The `/docker/` directory is a user-facing symlink alias for all stack operations, pointing directly to `/opt/m3tal/stack/`. This ensures that all user-initiated stack commands operate on the correct, centralized set of configuration files.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.

### Initial Deployment and Updates

The `m3tal up` command is used to deploy and update all defined Docker Compose stacks.

## Dashboard Access

The dashboard has two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable within `/etc/m3tal/.env`.

**Local mode (default)**
- Setting: `DASHBOARD_EXPOSE_MODE=local`
- This mode utilizes an override file that adds a direct port binding: `${DASHBOARD_PORT:-8082}:8082` to the `m3tal-dashboard` service.
- A new user performing a default installation will access the dashboard directly via `http://HOST_IP:8082` or `http://localhost:8082`.
- Traefik is not required for access in this mode. This is ideal for LAN-only setups, initial testing, or home server environments.

**Traefik mode**
- Setting: `DASHBOARD_EXPOSE_MODE=traefik`
- This mode uses an override file that configures Traefik to route traffic.
- Access is provided via a domain name at `http://dash.DOMAIN`.
- Traefik must be running and configured to manage this domain for access. This mode is suited for domain-based setups and managing multiple services behind a reverse proxy.

## Traefik Gateway

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml` and configured to bind port `80` (and optionally `443` for HTTPS) on the host as its HTTP/HTTPS entry points. It utilizes the Docker provider for automatic service discovery and the file provider to load dynamic configuration from the `/etc/traefik/dynamic/` directory, allowing for hot-reloading of routing rules.

For instance, Traefik can route requests to the M3TAL API daemon. A dynamic configuration file (such as `dynamic/api.yml`) can be created within the `/docker/dynamic/` directory to define routing for `api.DOMAIN`. This configuration will direct traffic to the Go API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`.

The `dash.DOMAIN` route is configured to direct traffic to the dashboard container when `DASHBOARD_EXPOSE_MODE=traefik`.

Here is a concrete YAML example of how to expose a custom user service via Traefik labels within a hypothetical `my-app-compose.yml` file:

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

## Service Management — systemd

The M3TAL API daemon is managed by systemd as `m3tal-api.service`. The following commands can be used to interact with this service:

- `systemctl status m3tal-api`: View the current status of the API service.
- `systemctl restart m3tal-api`: Restart the API service.
- `journalctl -u m3tal-api -f`: Tail the logs of the API service in real-time.

## Quick Demo

To start the M3TAL dashboard container specifically, run:

```bash
m3tal dash up
```

The `m3tal up` command orchestrates and deploys all other stacks defined in the `/docker/` directory, including any user-defined compose files that have been added.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                               | Access                                    | Description                                                                                             |
|------|---------------------------------------|-------------------------------------------|---------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point              | Public (when Traefik is active)           | The public-facing HTTP port for services exposed via Traefik.                                           |
| 8080 | M3TAL API daemon (Go)                 | Host-local                                | The internal port the M3TAL API daemon listens on.                                                      |
| 8081 | Traefik dashboard                     | Host-local only                           | The internal Traefik dashboard port, accessible only from the host machine.                             |
| 8082 | M3TAL Dashboard                       | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).