# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

Follow these steps to install the M3TAL CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The following paths define the M3TAL system's filesystem structure and purpose:

| Path                      | Purpose                                                  |
| ------------------------- | -------------------------------------------------------- |
| `/etc/m3tal/.env`         | Primary configuration file, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database, auto-created by the API daemon. |
| `/opt/m3tal/stack/`       | Canonical stack directory containing compose files and Traefik config. |
| `/docker`                 | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`      | Dashboard credential store, managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker/` directory is a user-facing symlink alias for all stack operations, pointing to the canonical source of truth located at `/opt/m3tal/stack/`.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3. Run `m3tal up` to start all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard offers two access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### Local Mode (Default)

When `DASHBOARD_EXPOSE_MODE=local`, the dashboard is directly accessible via port binding.

- **Access Method:** `http://HOST_IP:8082` (or `http://localhost:8082`).
- **Mechanism:** An override file (`m3tal-compose.local.yml`) adds a direct port mapping of `${DASHBOARD_PORT:-8082}:8082`.
- **Prerequisites:** None beyond M3TAL installation. Traefik is not required for this mode.
- **Use Case:** Ideal for local testing, initial setup, and environments where public domain access is not desired or feasible. A new user performing a default installation will access the dashboard directly via port 8082.

### Traefik Mode

When `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard is exposed via the Traefik reverse proxy.

- **Access Method:** `http://dash.DOMAIN`.
- **Mechanism:** Traefik labels in the dashboard's compose file direct traffic to the dashboard container listening on port 8082. This requires Traefik to be running.
- **Prerequisites:** Traefik must be running and configured to route to the dashboard. The `DOMAIN` environment variable must be set in `/etc/m3tal/.env`.
- **Use Case:** Suitable for production environments, domain-based routing, and when integrating the dashboard with other services managed by Traefik.

## Traefik Gateway

Traefik functions as the primary reverse proxy, automatically discovering and routing traffic to Docker services. This is achieved by interpreting Traefik labels defined within service definitions in Docker Compose files. By default, services are not exposed by Traefik; they require `traefik.enable=true` and other relevant labels to become discoverable and routable.

### Dynamic Configuration

Dynamic configuration files, located in `/docker/dynamic/`, allow for flexible routing rules. For instance, the `dynamic/api.yml` file configures Traefik to route requests for `api.DOMAIN` to the M3TAL API daemon listening on host-local port `8080` via `http://host.docker.internal:8080`.

The dashboard itself is routed via its Traefik labels. When `DASHBOARD_EXPOSE_MODE=traefik`, the `dash.DOMAIN` hostname is directed to the dashboard container.

### Exposing Custom User Services

You can expose custom user services by adding Traefik labels to their respective Docker Compose files. Below is a concrete YAML example for a hypothetical `my-app-compose.yml`:

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

The M3TAL API daemon is managed by `systemd` as the `m3tal-api.service`. The following commands can be used to manage its lifecycle:

- **Status:** `systemctl status m3tal-api`
- **Restart:** `systemctl restart m3tal-api`
- **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

### Starting the Dashboard

To start only the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command specifically manages the dashboard container, downloading the appropriate compose files based on your `DASHBOARD_EXPOSE_MODE` setting and launching it.

### Deploying All Stacks

The `m3tal up` command orchestrates and deploys all stacks defined in the `/docker/` directory. This includes the dashboard, Traefik, and any custom user-defined Docker Compose files you have placed in `/docker/`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                                       | Access                                            | Description                                                                                |
|------|-----------------------------------------------|---------------------------------------------------|--------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point                      | Public (traefik mode)                             | The public-facing HTTP port for services exposed via Traefik.                              |
| 8080 | M3TAL API daemon (Go)                         | Host-local                                        | The internal port the M3TAL API daemon listens on.                                         |
| 8081 | Traefik dashboard                           | Host-local only                                   | The internal Traefik dashboard port, accessible only from the host machine.                |
| 8082 | M3TAL Dashboard                               | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).