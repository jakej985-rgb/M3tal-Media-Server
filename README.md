# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

M3TAL is distributed via an APT repository. Follow these steps to install the CLI binary and system services:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

The following table outlines the key paths and their purposes within the M3TAL filesystem:

| Path                        | Purpose                                                              |
|-----------------------------|----------------------------------------------------------------------|
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.               |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config.|
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. User-facing path for stack operations.|
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |

## Firewall Considerations

If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The `/opt/m3tal/stack/` directory is the canonical source of truth where all stack files reside. For user convenience and interaction, `/docker` serves as a symlink alias to `/opt/m3tal/stack/`. All user-facing stack operations, such as adding or managing Docker Compose files, should target the `/docker/` directory.

### Adding a New Stack

To deploy a new Docker Compose stack:
1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2. Ensure any required environment variables for your new stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3. Run `m3tal up` to deploy all stacks, including your newly added one.

## Dashboard Access

The M3TAL Dashboard provides a web-based interface for system management. It supports two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env`.

### 1. Local Mode (Default)

-   **Setting:** `DASHBOARD_EXPOSE_MODE=local`
-   **Mechanism:** This mode uses an override (`m3tal-compose.local.yml`) to add a direct port binding, exposing the dashboard container's internal port `8082` to the host's `8082` port.
-   **Access:** Directly via `http://HOST_IP:8082` or `http://localhost:8082`.
-   **Requirement:** No Traefik gateway is needed for this mode.
-   **Typical Use:** This is the default behavior for a new installation, allowing immediate access to the dashboard without domain configuration or a reverse proxy. A new user performing a default installation will access the dashboard directly via port 8082.

### 2. Traefik Mode

-   **Setting:** `DASHBOARD_EXPOSE_MODE=traefik`
-   **Mechanism:** This mode uses an override (`m3tal-compose.traefik.yml`) to apply Traefik labels to the dashboard service. Traefik, if running, then routes incoming requests for `dash.DOMAIN` to the dashboard container's internal port `8082`.
-   **Access:** Via `http://dash.DOMAIN`.
-   **Requirement:** Traefik must be deployed and running via `m3tal up` for this mode to function.
-   **Typical Use:** Suited for domain-based setups where multiple services are exposed through a unified reverse proxy.

## Traefik Gateway

The M3TAL system utilizes Traefik as a reverse proxy to manage and route incoming traffic to various Docker services. Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik is deployed via `routing-compose.yml` and listens on host port `80` (HTTP). It is configured to use Docker as a provider for automatic service discovery and a file provider for dynamic configuration.

### Dynamic Configuration

Traefik loads additional routing rules from dynamic configuration files located in `/docker/dynamic/` (symlinked from `/opt/m3tal/stack/dynamic/`). These files allow for routing requests to services that are not necessarily Docker containers managed by labels (e.g., host-local applications).

For instance, the M3TAL API daemon (a Go binary) runs directly on the host machine on port `8080`. Traefik routes `api.DOMAIN` to this host-local service using a dynamic configuration file (e.g., `dynamic/api.yml`) that points to `http://host.docker.internal:8080`.

Example `dynamic/api.yml` content:
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

Similarly, when `DASHBOARD_EXPOSE_MODE=traefik`, Traefik routes `dash.DOMAIN` to the M3TAL Dashboard container on its internal port `8082` through labels defined in `m3tal-compose.traefik.yml`.

### Exposing a Custom User Service via Traefik

To expose a custom Docker Compose service (e.g., defined in `my-app-compose.yml`) via Traefik, you must add specific Traefik labels to its service definition:

```yaml
# /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-app
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)" # Routes app.DOMAIN to this service
      - "traefik.http.services.myapp.loadbalancer.server.port=80" # Internal port of the container
      - "traefik.http.routers.myapp.entrypoints=web" # Uses the 'web' (HTTP) entrypoint
    networks:
      - proxy # Crucial: service must be on the 'proxy' network Traefik listens on

networks:
  proxy:
    external: true # Assumes the 'proxy' network is external and created by Traefik
```
After placing this file in `/docker/` and ensuring `DOMAIN` is set in `/etc/m3tal/.env`, run `m3tal up` to deploy the service and make it discoverable by Traefik.

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard `systemctl` commands:

-   **Check Status:** `systemctl status m3tal-api`
-   **Restart Service:** `systemctl restart m3tal-api`
-   **View Logs:** `journalctl -u m3tal-api -f` (for real-time logs)

## Quick Demo

To quickly get started with M3TAL's primary components:

-   To start *only* the M3TAL Dashboard container with its default (local) port mapping, use:
    ```bash
    m3tal dash up
    ```
    This command specifically manages the dashboard container: it downloads the necessary compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) and starts the dashboard using the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting.

-   To orchestrate and deploy *all* Docker Compose stacks, including the dashboard (if not already running), Traefik, and any user-defined compose files located in `/docker/`, use:
    ```bash
    m3tal up
    ```
    This command ensures all services defined in `*-compose.yml` files within `/docker/` are brought up according to their configurations.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service                     | Access                                      | Description                                                                                                                                                                          |
|------|-----------------------------|---------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 80   | Traefik HTTP entry point    | Public                                      | The public-facing HTTP port for services exposed via Traefik.                                                                                                                        |
| 8080 | M3TAL API daemon (Go)       | Host-local                                  | The internal port the M3TAL API daemon listens on.                                                                                                                                   |
| 8081 | Traefik dashboard           | Host-local only                             | The internal Traefik dashboard port, accessible only from the host machine.                                                                                                          |
| 8082 | M3TAL Dashboard             | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`.                                                                       |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.