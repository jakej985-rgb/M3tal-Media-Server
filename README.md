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

The following paths are integral to the M3TAL system:

| Path                        | Purpose                                                        |
| :-------------------------- | :------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.  |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.       |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. User-facing path for stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.       |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth for all M3TAL stack definition files. The `/docker` directory is a user-facing symlink alias for all stack operations, pointing directly to `/opt/m3tal/stack/`.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Execute `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`:

**Local Mode (Default)**
When `DASHBOARD_EXPOSE_MODE=local`, the dashboard is directly accessible on your host machine's IP address. A new user performing a default installation will access the dashboard directly via port 8082. This mode bypasses Traefik for dashboard access.
Access via: `http://HOST_IP:8082` or `http://localhost:8082`.

**Traefik Mode**
When `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard is exposed through the Traefik gateway. Traefik will route traffic for `http://dash.DOMAIN` to the dashboard container. This mode requires Traefik to be running.
Access via: `http://dash.DOMAIN`.

## Traefik Gateway

Traefik acts as the reverse proxy for M3TAL, routing traffic to various services based on domain names. It is deployed via `routing-compose.yml`.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik binds to host port `80` (and optionally `443` for HTTPS) as its public entry point.

### Dynamic Configuration

Traefik utilizes a file provider to load dynamic routing configurations from the `/etc/traefik/dynamic/` directory within the Traefik container. These files allow for fine-grained control over routing and can be hot-reloaded without restarting Traefik.

For example, the `dynamic/api.yml` file configures Traefik to route requests for `api.DOMAIN` to the M3TAL API daemon, which listens on host-local port `8080`.

**Dynamic routing example (`dynamic/api.yml`):**
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

### Exposing Custom Services

To expose a custom user service via Traefik, define the appropriate Traefik labels within your service's Docker Compose definition.

**Example YAML for a hypothetical `my-app-compose.yml`:**
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

The M3TAL API daemon is managed by systemd. The following commands can be used to interact with it:

*   `systemctl status m3tal-api`: Check the status of the M3TAL API service.
*   `systemctl restart m3tal-api`: Restart the M3TAL API service.
*   `journalctl -u m3tal-api -f`: View the logs for the M3TAL API service in real-time.

## Quick Demo

This section outlines basic commands for interacting with M3TAL.

### Starting the Dashboard

To specifically start only the dashboard container:

```bash
m3tal dash up
```

This command ensures the dashboard is running, respecting the `DASHBOARD_EXPOSE_MODE` setting.

### Deploying All Stacks

To orchestrate and deploy all stacks defined in the `/docker/` directory, including any user-added compose files:

```bash
m3tal up
```

This command is a wrapper around `docker compose` and manages the lifecycle of all M3TAL-managed Docker stacks.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                                        | Access                                            | Description                                                                            |
| :--- | :--------------------------------------------- | :------------------------------------------------ | :------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point                       | Public                                            | The public-facing HTTP port for services exposed via Traefik.                          |
| 8080 | M3TAL API daemon (Go)                          | Host-local                                        | The internal port the M3TAL API daemon listens on.                                     |
| 8081 | Traefik dashboard                              | Host-local only                                   | The internal Traefik dashboard port, accessible only from the host machine.            |
| 8082 | M3TAL Dashboard                                | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).