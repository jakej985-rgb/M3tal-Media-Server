# M3TAL Ecosystem

This document details the architecture and operational procedures for the M3TAL system. M3TAL provides a unified management interface for Docker-based services, leveraging Docker Engine and Docker Compose V2 for container orchestration and Traefik for ingress routing.

## Core Components

The M3TAL ecosystem is comprised of the following components:

*   **CLI binary (`/usr/bin/m3tal`)**: A Go binary, installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API daemon (`m3tal-api.service`)**: A Go binary running as a systemd service, listening on port `8080`. This daemon is responsible for managing Docker interactions, persistent state in its SQLite database, and exposing the M3TAL API routes.
*   **Dashboard container (`m3tal-dashboard`)**: A Python/Flask application running in a Docker container on port `8082`. It communicates with the M3TAL API daemon via `http://host.docker.internal:8080`.
*   **Traefik gateway (`routing-compose.yml`)**: A Docker container acting as a reverse proxy. It exposes services on port `80` by domain name, using a file provider for dynamic configuration.
*   **Cloudflared (`routing-compose.yml`)**: An optional Cloudflare Tunnel container that provides secure, zero-configuration internet access for exposed services.

## Filesystem Contract

The following paths represent the canonical locations for M3TAL's operational files:

| Path                        | Purpose                                                 |
| :-------------------------- | :------------------------------------------------------ |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon.  |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's core compose files and Traefik dynamic configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`. |

## Runtime Environment

M3TAL operates on **Docker Engine** and requires **Docker Compose V2** to be installed and available on the host system. All container orchestration is performed using `docker compose` commands executed by the M3TAL CLI and API daemon.

The `m3tal up` command orchestrates all `*-compose.yml` files located within the `/docker/` directory. The `m3tal dash up` command specifically manages the dashboard container and its dependencies, defined in `/docker/m3tal-compose.yml`. All compose files are configured to use the shared environment variables from `/etc/m3tal/.env` via the `--env-file` flag.

## APT Installation

To install M3TAL, execute the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Deployment Lifecycle

### Stacks

M3TAL manages application stacks using Docker Compose files. These files are located within the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`).

To bring up all configured stacks, execute:
```bash
m3tal up
```
This command iterates through all `*-compose.yml` files in `/docker/` and executes `docker compose up -d` for each.

### Adding a New Stack

To add a new application stack to M3TAL:

1.  Create your Docker Compose file (e.g., `my-app-compose.yml`) and place it in the `/docker/` directory.
2.  Ensure any environment variables required by your new stack are configured in `/etc/m3tal/.env`. You can manage these variables using the CLI:
    *   `m3tal config wizard` for interactive setup.
    *   `m3tal config set KEY value` to set individual variables.
3.  Run `m3tal up` to start your new stack along with all other configured stacks.
    Alternatively, to manage a single stack, you can use standard `docker compose` commands directly:
    `docker compose -f /docker/my-app-compose.yml up -d`

## Traefik Gateway

Traefik runs as a Docker container, deployed via `routing-compose.yml`, and serves as the primary HTTP entry point for all M3TAL-managed services.

*   **Port Binding**: Traefik binds to port `80` on the host system.
*   **Service Exposure**: Services are exposed by adding Traefik labels to their respective service definitions within their Docker Compose files. These labels define routing rules (e.g., hostnames) and other configurations.
*   **Dynamic Configuration**: Traefik uses a file provider to load dynamic routing configurations from `/docker/dynamic/`. These configurations are hot-reloaded when changes are detected.
*   **M3TAL API Routing**: The `api.DOMAIN` hostname is routed to the M3TAL API daemon. This is achieved via a dynamic configuration file (e.g., `dynamic/api.yml`) which directs traffic to `http://host.docker.internal:8080`.
*   **M3TAL Dashboard Routing**: The `dash.DOMAIN` hostname routes traffic to the M3TAL dashboard container, typically listening on port `8082` within its container. This routing is configured via Traefik labels within `m3tal-compose.traefik.yml`.
*   **Traefik Dashboard**: The Traefik dashboard itself is accessible on `http://localhost:8081` from the host system only.

### Traefik Static Configuration (`traefik.yml` example)

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

### Dynamic Routing Example (`dynamic/api.yml`)

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

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`.

*   **Status**: `systemctl status m3tal-api`
*   **Restart**: `systemctl restart m3tal-api`
*   **Logs**: `journalctl -u m3tal-api -f`

The M3TAL CLI also interacts with this systemd service; for instance, `m3tal dash up` will implicitly trigger `systemctl start m3tal-api` if the daemon is not running.

## Port Map

| Port | Service                            | Access                               |
| :--- | :--------------------------------- | :----------------------------------- |
| 80   | Traefik HTTP entry point           | Public (via configured `DOMAIN`)     |
| 8080 | M3TAL API daemon (Go)              | Host-local (accessed via Traefik or direct) |
| 8081 | Traefik dashboard                  | Host-local only                      |
| 8082 | M3TAL Dashboard (Python/Flask)     | Via Traefik or direct (via configured `DOMAIN`) |

## Firewall Configuration

Ensure that port `80` is allowed through your host's firewall (e.g., `ufw` or `iptables`) to permit external access to services exposed via the Traefik gateway.

Example for UFW: `sudo ufw allow 80/tcp`

## Quick Demo

Follow these steps to quickly set up and access the M3TAL Dashboard:

1.  **Install M3TAL**:
    ```bash
    # 1. Add the GPG signing key
    curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

    # 2. Add the APT repository
    echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

    # 3. Install
    sudo apt update && sudo apt install -y m3tal
    ```

2.  **Run Configuration Wizard**:
    ```bash
    m3tal config wizard
    ```
    Follow the prompts. For a local demo, accept the default `DOMAIN=localhost`.

3.  **Deploy Dashboard and Core Services**:
    ```bash
    m3tal dash up
    ```
    This command starts the M3TAL API daemon, the Traefik gateway, and the M3TAL Dashboard container.

4.  **Access Dashboard**:
    Open your web browser and navigate to `http://dash.localhost`.
    Use the `admin_pass` for the default password, or the password you set during the `config wizard` for the `ADMIN_PASSWORD` variable.