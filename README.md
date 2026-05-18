# M3TAL Ecosystem Documentation

This document describes the M3TAL system architecture, deployment, and operational procedures.

## Runtime Environment

M3TAL operates on a host system running **Docker Engine** and **Docker Compose V2**. These are hard dependencies for container orchestration. The `m3tal` CLI manages Docker Compose operations, primarily by executing `docker compose` commands against stack definitions.

## Installation

The M3TAL CLI binary (`/usr/bin/m3tal`) is distributed via an APT repository.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Filesystem Contract

M3TAL maintains a specific directory structure for configuration and operational data.

| Path                     | Purpose                                          |
|--------------------------|--------------------------------------------------|
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`. |

## Components

The M3TAL ecosystem consists of the following primary components:

*   **CLI binary** (`/usr/bin/m3tal`): A Go binary providing a single entrypoint for system operations, configuration management, and orchestrating Docker Compose stacks.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on port 8080. It manages Docker interactions, the state database (`/var/lib/m3tal/state.db`), and exposes API routes for M3TAL operations.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask application running in a Docker container on port 8082. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy running in a Docker container, binding host port 80. It exposes services via domain names and utilizes file providers for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, providing zero-configuration internet access for exposed services.

## Configuration

System-wide configuration is stored in `/etc/m3tal/.env`. This file is managed interactively using the `m3tal config wizard` command.
Docker Compose stacks reference this file via the `--env-file` flag during their lifecycle.

Key environment variables include:

*   `DOMAIN`: The base domain used for Traefik routing (e.g., `localhost`).
*   `ADMIN_PASSWORD`: Password for dashboard administration.
*   `API_TOKEN`: Token for API authentication.
*   `DASHBOARD_SECRET`: Secret for dashboard session management.
*   `NETWORK_NAME`: Name of the Docker network used by M3TAL services (default: `m3tal`).

## Deployment Lifecycle

M3TAL manages Docker Compose stacks. The `m3tal up` command iterates through all `*-compose.yml` files located in the `/docker/` directory and brings them up using `docker compose`.

### Adding a New Stack

To deploy a new application stack:

1.  Create your Docker Compose definition file (e.g., `my-stack-compose.yml`).
2.  Place this file in the `/docker/` directory.
3.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these.
4.  Run `m3tal up` to start all defined stacks, including your new one.
    *   Alternatively, to manage a single stack, use `docker compose -f /docker/my-stack-compose.yml up -d`.

The `m3tal dash up` command specifically manages the `m3tal-dashboard` container, defined in `/docker/m3tal-compose.yml`.

## Traefik Gateway

Traefik runs as a Docker container, deployed via `routing-compose.yml`. It functions as the primary HTTP entry point for M3TAL services.

*   **Port Binding:** Traefik binds host port 80 as its HTTP entry point.
*   **Service Discovery:** Services are exposed by configuring Traefik labels directly within their Docker Compose service definitions.
*   **Dynamic Configuration:** Traefik loads additional routing rules from `/etc/traefik/dynamic` (symlinked from `/docker/dynamic`), which supports hot-reloading.
*   **API Daemon Routing:** `api.DOMAIN` routes to the M3TAL API daemon (running on the host at `http://host.docker.internal:8080`) via a dynamic configuration file, `/docker/dynamic/api.yml`.
*   **Dashboard Routing:** `dash.DOMAIN` routes traffic to the `m3tal-dashboard` container via Traefik labels defined within `m3tal-compose.traefik.yml`.
*   **Traefik Dashboard:** The Traefik dashboard itself is accessible locally at `http://localhost:8081`.

### Traefik Static Configuration (`traefik.yml`)

```yaml
entryPoints:
  web:
    address: ":80"

providers:
  docker:
    exposedByDefault: false
    network: proxy # Assumes a 'proxy' network, which should be shared by services
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

### Dynamic Routing Example (`/docker/dynamic/api.yml`)

This configuration routes requests for `api.DOMAIN` to the M3TAL API daemon.

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

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service.

*   **Check Status:** `systemctl status m3tal-api`
*   **Restart Service:** `systemctl restart m3tal-api`
*   **View Logs:** `journalctl -u m3tal-api -f`

The `m3tal` CLI also interacts with the API daemon. For example, `m3tal dash up` initiates a `systemctl start m3tal-api` command if the API is not running.

## Port Map

| Port | Service                         | Access                                    |
|------|---------------------------------|-------------------------------------------|
| 80   | Traefik HTTP entry point        | Public                                    |
| 8080 | M3TAL API daemon (Go)           | Host-local (accessible via Traefik or direct) |
| 8081 | Traefik dashboard               | Host-local only                           |
| 8082 | M3TAL Dashboard (Python/Flask)  | Via Traefik or direct                     |

## Firewall Considerations

For Traefik to function as the public entry point, ensure that host port 80 is allowed through your system's firewall (e.g., UFW or iptables).

Example for UFW:
`sudo ufw allow 80/tcp`

## Quick Demo

Follow these steps to get a basic M3TAL setup operational:

1.  **Install M3TAL:**
    ```bash
    # 1. Add the GPG signing key
    curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

    # 2. Add the APT repository
    echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

    # 3. Install
    sudo apt update && sudo apt install -y m3tal
    ```

2.  **Run Configuration Wizard:**
    ```bash
    m3tal config wizard
    ```
    Follow the prompts to set your `DOMAIN`, `ADMIN_PASSWORD`, and other necessary variables.

3.  **Deploy Dashboard and Routing:**
    ```bash
    m3tal dash up
    ```
    This command starts the `m3tal-api.service`, the Traefik gateway, and the M3TAL dashboard container.

4.  **Open Browser:**
    Navigate to `http://dash.<YOUR_DOMAIN>` (replace `<YOUR_DOMAIN>` with the value you set during the config wizard, e.g., `http://dash.localhost`). You should see the M3TAL Dashboard login page.