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

## Configuration

The primary configuration file for M3TAL is located at `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` command or can be updated using `m3tal config set KEY value`.

## Filesystem Contract

The following table outlines the critical file paths and their purposes within the M3TAL system:

| Path                     | Purpose                                                              |
|--------------------------|----------------------------------------------------------------------|
| `/etc/m3tal/.env`        | Primary environment configuration file.                              |
| `/var/lib/m3tal/state.db` | SQLite state database for the API daemon.                            |
| `/opt/m3tal/stack/`      | Canonical directory for Docker Compose stack files and Traefik configs. |
| `/docker`                | User-facing symlink alias for `/opt/m3tal/stack/` for stack operations. |
| `/docker/users.json`     | Stores dashboard credentials, managed via `m3tal dashpass`.          |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/docker` directory serves as the user-facing alias for all stack operations, and it is a symbolic link to `/opt/m3tal/stack/`. This means that `/opt/m3tal/stack/` is the canonical source of truth for all stack files.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any necessary environment variables are configured in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value`.
3. Run `m3tal up` to initiate the deployment of all stacks, including your newly added one.

## Dashboard Access

The M3TAL dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`:

### Local Mode (Default)

- **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
- **Access Method:** Direct port binding. A new user performing a default installation will access the dashboard directly via port 8082.
- **URL:** `http://HOST_IP:8082` (replace `HOST_IP` with your server's IP address) or `http://localhost:8082`.
- **Requirements:** No Traefik is needed for access in this mode.
- **Use Case:** Ideal for local testing, home server setups, and initial deployments where domain-based access is not yet configured.

### Traefik Mode

- **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
- **Access Method:** Domain routing through Traefik.
- **URL:** `http://dash.DOMAIN` (replace `DOMAIN` with your configured domain).
- **Requirements:** Traefik must be running and configured to route traffic.
- **Use Case:** Suitable for production environments, multi-service deployments, and when utilizing a reverse proxy for centralized access.

## Traefik Gateway

Traefik acts as the reverse proxy for M3TAL, automatically discovering and routing traffic to Docker services. This is achieved by interpreting Traefik labels defined within service definitions in Docker Compose files.

Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik loads dynamic configuration from files within the `/opt/m3tal/stack/dynamic/` directory. This allows for routing requests to services listening on host-local ports. For example, the `api.DOMAIN` route is configured to direct traffic to the Go API daemon running on host-local port `8080` via the `dynamic/api.yml` file:

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

Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard container is configured via Traefik labels in `m3tal-compose.traefik.yml` to be accessible at `http://dash.DOMAIN`.

### Exposing Custom User Services via Traefik

To expose a custom user service through Traefik, define the appropriate Traefik labels within its Docker Compose service definition. Here is an example for a hypothetical `my-app-compose.yml`:

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
This example configures Traefik to route requests for `app.DOMAIN` to the `my-app` service.

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon is managed by `systemd` as `m3tal-api.service`. The following commands can be used to manage its lifecycle:

- **Check status:** `systemctl status m3tal-api`
- **Restart service:** `systemctl restart m3tal-api`
- **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

This section provides a brief overview of how to start and manage M3TAL services.

### Starting the Dashboard

To start only the M3TAL dashboard container, use the following command:

```bash
m3tal dash up
```

This command specifically targets the dashboard service and its relevant configuration, allowing for quick access and testing of the dashboard interface.

### Deploying All Stacks

The `m3tal up` command orchestrates and deploys all other stacks defined in the `/docker/` directory. This includes any user-defined Docker Compose files placed in `/docker/`, ensuring that your entire M3TAL setup is running.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                   | Access                                         | Description                                                                                                |
|------|---------------------------|------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point  | Public (when Traefik is active)                | The public-facing HTTP port for services exposed via Traefik.                                              |
| 8080 | M3TAL API daemon (Go)     | Host-local                                     | The internal port the M3TAL API daemon listens on. Accessible by services running on the host or via Traefik. |
| 8081 | Traefik dashboard         | Host-local only                                | The internal Traefik dashboard port, accessible only from the host machine.                                 |
| 8082 | M3TAL Dashboard           | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.