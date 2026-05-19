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

## Configuration

The primary configuration file for M3TAL is located at `/etc/m3tal/.env`. This file is managed by the `m3tal config wizard` command. Environment variables such as `DASHBOARD_EXPOSE_MODE`, `DOMAIN`, and various path variables are defined here.

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack definition files reside. The `/docker/` directory serves as a user-facing symlink alias for all stack operations, pointing directly to `/opt/m3tal/stack/`. When you place a Docker Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory, it is automatically recognized and included by `m3tal up`.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any necessary environment variables required by your stack are configured in `/etc/m3tal/.env`. The `m3tal config wizard` or `m3tal config set KEY value` commands can be used for this purpose.
3.  Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard has two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

*   **Local Mode (default):**
    *   `DASHBOARD_EXPOSE_MODE=local`
    *   When `DASHBOARD_EXPOSE_MODE` is set to `local` (or is not explicitly defined, as `local` is the default), the dashboard container is directly exposed on the host machine's IP address at port 8082.
    *   Access is provided via `http://HOST_IP:8082` or `http://localhost:8082`.
    *   This mode does not require Traefik to be running.
    *   A new user performing a default installation will access the dashboard directly via port 8082.

*   **Traefik Mode:**
    *   `DASHBOARD_EXPOSE_MODE=traefik`
    *   When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, Traefik is configured to route traffic to the dashboard.
    *   Access is provided via a domain name, typically `http://dash.DOMAIN` (where `DOMAIN` is configured in `/etc/m3tal/.env`).
    *   This mode requires Traefik to be actively running as part of a deployed stack.

## Traefik Gateway

Traefik acts as the primary reverse proxy for M3TAL, managing external access to various services. It automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik dynamically loads configuration files from the `/docker/dynamic/` directory. These files define rules for routing requests. For instance, a file like `dynamic/api.yml` is used to route requests for `api.DOMAIN` to the M3TAL API daemon, which is listening on host-local port `8080` via `http://host.docker.internal:8080`.

### Dynamic Routing Examples

*   **`dash.DOMAIN` routes to the dashboard container:** This is configured via Traefik labels within the dashboard's compose file (specifically `m3tal-compose.traefik.yml`).
*   **`api.DOMAIN` routes to the M3TAL API daemon:** This is configured in `dynamic/api.yml` as follows:

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

### Exposing a Custom User Service via Traefik

To expose a custom user service (e.g., a hypothetical `my-app`) via Traefik labels, you would include these labels within your custom Docker Compose file (e.g., `my-app-compose.yml`) placed in the `/docker/` directory:

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

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Service Management

The M3TAL API daemon is managed by `systemd`. The service name is `m3tal-api.service`.

*   **Check status:** `systemctl status m3tal-api`
*   **Restart service:** `systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

To quickly start only the dashboard container:

```bash
m3tal dash up
```

This command specifically targets the dashboard container. The `m3tal up` command orchestrates and deploys all other stacks present in the `/docker/` directory, including any custom Docker Compose files you have added.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
|------|---------|--------|-------------|
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.