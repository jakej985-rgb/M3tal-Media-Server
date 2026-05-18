# M3TAL Ecosystem Documentation

This document describes the M3TAL system architecture and operational procedures.

## Overview

M3TAL is a system for managing and orchestrating containerized applications using Docker Engine and Docker Compose V2. It provides a unified CLI for deployment, configuration, and service management.

## Runtime Environment

M3TAL is designed to run on systems with **Docker Engine** and **Docker Compose V2** installed. All application stacks and management services are deployed and managed as Docker containers.

## Installation

To install M3TAL, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Configuration

The primary configuration file is located at `/etc/m3tal/.env`. This file stores environment variables that control M3TAL's behavior. Use the `m3tal config wizard` command to interactively set up your environment. Individual configuration variables can be set using `m3tal config set KEY value`.

## Filesystem Contract

M3TAL adheres to the following filesystem paths:

| Path                           | Purpose                                       |
| :----------------------------- | :-------------------------------------------- |
| `/etc/m3tal/.env`              | Primary configuration file.                   |
| `/var/lib/m3tal/state.db`      | SQLite state database for the API daemon.     |
| `/opt/m3tal/stack/`            | Canonical stack directory.                    |
| `/docker`                      | Symlink to `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`           | Dashboard credential store.                   |

## Deployment Lifecycle

M3TAL manages application deployments as "stacks."

### How Stacks Work

Stacks are defined by Docker Compose files. M3TAL looks for all files ending in `*-compose.yml` within the `/docker` directory. The `m3tal up` command orchestrates the startup of all services defined across these Compose files using `docker compose`.

### Adding a New Stack

To add a new application stack:

1.  Place your Docker Compose file (e.g., `my-app-compose.yml`) into the `/docker/` directory.
2.  Ensure any necessary environment variables for your new stack are configured in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set`).
3.  Run `m3tal up` to deploy and start all stacks, including your newly added one.

## Dashboard Access

The M3TAL dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Local Mode (Default)

When `DASHBOARD_EXPOSE_MODE` is set to `local` (this is the default setting upon fresh installation), the dashboard is directly exposed on the host machine.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Access:** `http://HOST_IP:8082` or `http://localhost:8082` (where `HOST_IP` is your server's IP address).
*   **Mechanism:** An override Compose file (`m3tal-compose.local.yml`) adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`).
*   **Use Case:** Ideal for home lab or local network setups where direct IP access is sufficient and Traefik is not yet deployed or configured. A new user will access the dashboard at `http://HOST_IP:8082` immediately after a default installation.

### Traefik Mode

When `DASHBOARD_EXPOSE_MODE` is set to `traefik`, the dashboard is exposed via the Traefik reverse proxy.

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Access:** `http://dash.DOMAIN` (where `DOMAIN` is configured in `/etc/m3tal/.env`).
*   **Mechanism:** An override Compose file (`m3tal-compose.traefik.yml`) adds Traefik labels to the dashboard service. Traefik then routes incoming requests for `dash.DOMAIN` to the dashboard container.
*   **Prerequisite:** Traefik must be running and configured to accept traffic on port 80.
*   **Use Case:** Suitable for domain-based access, especially when running multiple services behind a single entry point.

## Traefik Gateway

Traefik acts as the primary reverse proxy for M3TAL services.

*   **Deployment:** Traefik runs as a container, typically defined in `routing-compose.yml`, and binds to host port 80 for HTTP traffic.
*   **Service Exposure:** Services are made accessible through Traefik by adding specific `traefik.*` labels to their respective Docker Compose service definitions.
*   **Dynamic Routing:** Traefik utilizes a file provider to load dynamic routing configurations from `/etc/traefik/dynamic/`.
*   **API Routing:** The Go API daemon is exposed via `api.DOMAIN` (e.g., `api.localhost` or `api.yourdomain.com`). This is configured in `dynamic/api.yml`.
*   **Dashboard Routing:** The M3TAL dashboard is exposed via `dash.DOMAIN` when `DASHBOARD_EXPOSE_MODE` is set to `traefik`. This is handled by Traefik labels in `m3tal-compose.traefik.yml`.

**Traefik Static Configuration (`traefik.yml`):**
```yaml
entryPoints:
  web:
    address: ":80"

providers:
  docker:
    exposedByDefault: false
    network: proxy
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

**Dynamic Routing Example (`dynamic/api.yml`):**
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

## Service Management

The M3TAL API daemon (`m3tal-api.service`) is managed by `systemd`. Use the following commands for service control:

*   **Check status:** `sudo systemctl status m3tal-api`
*   **Restart service:** `sudo systemctl restart m3tal-api`
*   **View logs:** `journalctl -u m3tal-api -f`

## Quick Demo

1.  **Install M3TAL:**
    ```bash
    # Follow the APT installation steps above.
    ```
2.  **Run Configuration Wizard:**
    ```bash
    m3tal config wizard
    ```
    Follow the prompts to set essential parameters.
3.  **Start the Dashboard:**
    ```bash
    m3tal dash up
    ```
    This command will download necessary Compose files and start the dashboard container. By default, it will be accessible via your host IP address.
4.  **Access Dashboard:**
    Open your web browser and navigate to `http://<your_host_ip>:8082`.

## Port Map

| Port | Service(s)                       | Access Method                                  |
| :--- | :------------------------------- | :--------------------------------------------- |
| 80   | Traefik HTTP Entry Point         | Public (requires `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API Daemon (Go)            | Host-local only (accessed by containers)       |
| 8081 | Traefik Dashboard                | Host-local only (`127.0.0.1:8081`)             |
| 8082 | M3TAL Dashboard                  | Direct port (`http://HOST_IP:8082`) in local mode, or via Traefik (`http://dash.DOMAIN`) in traefik mode. |

## Firewall Configuration

If you are using Traefik (and therefore exposing services on port 80), ensure that port 80 is allowed through your host's firewall (e.g., `ufw` or `iptables`).

```bash
# Example for ufw:
sudo ufw allow 80/tcp
sudo ufw reload
```