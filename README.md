# M3TAL System Documentation

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2.

## Installation

To install the M3TAL CLI binary and API daemon:

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

## Components

The M3TAL system comprises the following core components:

*   **CLI binary** (`/usr/bin/m3tal`): A unified Go binary installed via APT, serving as the single entrypoint for all system operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on host-local port `8080`. It is responsible for managing Docker interactions, the SQLite state database, and API route handling.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port `8082`. It communicates with the API daemon at `http://host.docker.internal:8080` within the Docker network.
*   **Traefik gateway** (`routing-compose.yml`): A Docker container acting as a reverse proxy, exposing services by domain name on host port `80`. It utilizes a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel Docker container, enabling zero-configuration internet access for exposed services.

## Filesystem Contract

The M3TAL system adheres to the following filesystem contract:

| Path | Purpose |
| :---------------------- | :--------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env` | The primary configuration file, managed by the `m3tal config wizard` command. |
| `/var/lib/m3tal/state.db` | The SQLite state database, automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/` | The canonical directory for all Docker Compose stack files and Traefik dynamic configuration. This is the source of truth. |
| `/docker` | A symbolic link pointing to `/opt/m3tal/stack/`. This path serves as the user-facing alias for all stack operations. |
| `/docker/users.json` | The credential store for the M3TAL Dashboard, managed by the `m3tal dashpass` command. |

## Deployment Lifecycle

M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.

The directory `/opt/m3tal/stack/` is the canonical source of truth where all stack files reside. The `/docker` directory is a user-facing symlink alias for all stack operations, allowing users to interact with their Docker Compose files directly in `/docker`.

### Adding a New Stack

To deploy a new Docker Compose stack:

1.  Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically included by `m3tal up`.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  Run `m3tal up` to orchestrate and start all defined Docker Compose stacks, including your new one.

## Dashboard Access

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`)

This is the default access method.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` is set in `/etc/m3tal/.env`. This utilizes the `m3tal-compose.local.yml` override file.
*   **Mechanism:** A direct port binding is established for the dashboard container, typically mapping host port `8082` to the container's internal port `8082`.
*   **Access:** The dashboard is directly accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **New User Behavior:** A new user performing a default M3TAL installation will access the dashboard directly via port `8082`. This mode does not require Traefik to be running.
*   **Use Cases:** Ideal for LAN-only setups, first-time users, or local development and testing environments.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

This mode routes dashboard traffic through the Traefik Gateway.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik` must be set in `/etc/m3tal/.env`. This utilizes the `m3tal-compose.traefik.yml` override file.
*   **Mechanism:** Traefik labels are added to the dashboard service definition, enabling Traefik to discover and route requests for `dash.${DOMAIN}` to the dashboard container's internal port `8082`. The Traefik Gateway must be running.
*   **Access:** The dashboard is accessible via `http://dash.DOMAIN`, where `DOMAIN` is defined in your `/etc/m3tal/.env` file.
*   **Use Cases:** Suited for domain-based setups and environments where multiple services are exposed behind a single reverse proxy.

## Traefik Gateway

The Traefik Gateway acts as the primary reverse proxy for M3TAL, deployed as a container via `routing-compose.yml`.

Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.

Traefik's configuration includes:
*   Binding host port `80` as the HTTP entry point.
*   Discovering services through the Docker provider, with `exposedByDefault: false` to enforce explicit labeling.
*   Loading dynamic configuration files from `/docker/dynamic/` (via a file provider), which supports hot-reloading.

**Dynamic Configuration Example:**

Traefik uses dynamic configuration files (such as `/docker/dynamic/api.yml`) to route requests to services that may not be Docker containers themselves, or to apply specific routing logic. For instance, `api.DOMAIN` is routed to the M3TAL Go API daemon, which listens on the host-local port `8080`, via `http://host.docker.internal:8080`.

```yaml
# File: /docker/dynamic/api.yml
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

The `dash.DOMAIN` route, when `DASHBOARD_EXPOSE_MODE=traefik`, routes directly to the dashboard container via labels specified in its `m3tal-compose.traefik.yml` override.

### Exposing a Custom User Service via Traefik

To expose your own Docker Compose service (e.g., from `my-app-compose.yml`) through the Traefik Gateway, add appropriate Traefik labels to its service definition:

```yaml
# File: /docker/my-app-compose.yml
services:
  my-app:
    image: nginx:alpine
    container_name: my-custom-nginx
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.myapp.rule=Host(`app.DOMAIN`)"
      - "traefik.http.services.myapp.loadbalancer.server.port=80"
      - "traefik.http.routers.myapp.entrypoints=web"
      - "traefik.docker.network=proxy" # Ensure your service is connected to the 'proxy' network
    networks:
      - proxy

networks:
  proxy:
    external: true # Your service needs to connect to the external 'proxy' network
```

After adding this, run `m3tal up` to deploy the changes. Assuming `DOMAIN` is set in your `.env` (e.g., to `example.com`), `http://app.example.com` will route to your `my-app` service.

## Service Management

The M3TAL API daemon is managed as a systemd service, `m3tal-api.service`. Standard systemctl commands are used for its operation:

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

## Quick Demo

The `m3tal` CLI provides specific commands for managing different aspects of the system:

*   **Start the Dashboard:** To start only the M3TAL Dashboard container (and any necessary overrides based on `DASHBOARD_EXPOSE_MODE`), execute:
    ```bash
    m3tal dash up
    ```
    This command specifically manages the dashboard's lifecycle, downloading the latest compose files and applying the correct expose mode.

*   **Deploy All Stacks:** To orchestrate and deploy all other Docker Compose stacks, including any user-defined `*-compose.yml` files located in the `/docker/` directory, use:
    ```bash
    m3tal up
    ```
    This command will bring up the Traefik gateway, Cloudflared, and any custom services you have added to `/docker/`.

## Port Map

The following table lists the primary network ports utilized by the M3TAL system:

| Port | Service | Access | Description |
| :--- | :-------------------------- | :------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------- |
| 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
| 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
| 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |

Note: These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.