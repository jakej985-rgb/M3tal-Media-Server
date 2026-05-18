# M3TAL System Architecture

This document provides technical details and operational procedures for the M3TAL system.

## Prerequisites

Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation.

## APT Installation

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

| Path | Purpose |
|------|--------|
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

## Deployment Lifecycle

M3TAL utilizes Docker Compose for service orchestration.

### How Stacks Work

Compose files located within the `/docker/` directory define the services for various stacks. The `m3tal up` command orchestrates the bringing up and management of these stacks.

The `/docker` path is a user-facing symlink that points directly to the canonical stack directory: `/opt/m3tal/stack/`. All compose files and associated configurations for M3TAL services and user-added stacks reside here.

### Adding a New Stack

To deploy a new service stack:

1.  Place your Docker Compose file (e.g., `my-app-compose.yml`) into the `/docker/` directory.
2.  Ensure any necessary environment variables for your new stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these.
3.  Execute `m3tal up` to deploy all defined stacks, including your newly added one.

## Dashboard Access

The M3TAL dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (`DASHBOARD_EXPOSE_MODE=local`)

This is the default mode.

*   **Configuration:** Uses the `m3tal-compose.local.yml` override file.
*   **Access:** Provides direct port binding, making the dashboard accessible via `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Dependencies:** Does not require Traefik to be running.
*   **Use Case:** Ideal for home or LAN-only setups, first-time users, and local testing environments where external domain routing is not necessary.

### Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

*   **Configuration:** Uses the `m3tal-compose.traefik.yml` override file.
*   **Access:** Routes traffic to the dashboard via Traefik at `http://dash.DOMAIN`.
*   **Dependencies:** Requires Traefik to be deployed and running via `m3tal up`.
*   **Use Case:** Suitable for domain-based deployments where multiple services are exposed through a single reverse proxy.

A new user performing a default installation will access the dashboard directly via port 8082, not through a domain name.

## Traefik Gateway

Traefik functions as the primary ingress point for services exposed via domain names.

*   Traefik runs as a Docker container, managed by `routing-compose.yml`.
*   It binds to host port 80, serving as the HTTP entrypoint.
*   Services are exposed to Traefik by including specific Traefik labels within their respective Docker Compose service definitions.
*   The Go API daemon is accessible via `http://api.${DOMAIN}`. Traefik dynamically routes this to the API daemon's host-local port (8080) via a dynamic configuration file.
*   The M3TAL dashboard is accessible via `http://dash.${DOMAIN}` when `DASHBOARD_EXPOSE_MODE` is set to `traefik`. This is achieved through Traefik labels defined in the dashboard's compose override.

**Firewall Note:** Ensure that port 80 is open in your firewall (e.g., `ufw`, `iptables`) to allow external access to Traefik.

### Exposing a Custom User Service via Traefik

To expose a custom user-defined service, add the appropriate Traefik labels to its service definition in its compose file. Example for a hypothetical `my-app-compose.yml`:

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

The M3TAL API daemon is managed by `systemd`.

*   **Status:** `systemctl status m3tal-api.service`
*   **Restart:** `systemctl restart m3tal-api.service`
*   **Logs:** `journalctl -u m3tal-api.service -f`

## Quick Demo

### Starting the Dashboard

To specifically start only the dashboard container:

```bash
m3tal dash up
```

This command downloads the necessary compose files for the dashboard and starts it according to the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

### Deploying All Stacks

To orchestrate and deploy all defined stacks, including the dashboard, Traefik, and any user-added services:

```bash
m3tal up
```

This command processes all `*-compose.yml` files found in the `/docker/` directory.

## Port Map

These are the primary M3TAL system ports. User-added stacks may expose additional ports as defined in their respective compose files.

| Port | Service/Purpose | Access Context |
|------|-----------------|----------------|
| 80 | Traefik HTTP Entrypoint | Public (Traefik mode) |
| 8080 | M3TAL API Daemon (Go) | Host-local |
| 8081 | Traefik Dashboard | Host-local only |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (Traefik mode) |