# Getting Started with M3TAL

This guide provides a step-by-step process for first-time installation and basic setup of the M3TAL Ecosystem.

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify their installation and version:

```bash
docker --version
docker compose version
```

## 2. Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The M3TAL CLI provides an interactive wizard to set up essential environment variables in `/etc/m3tal/.env`. These variables control the behavior of the M3TAL API daemon, Dashboard, and Docker Compose stacks.

Run the wizard:

```bash
sudo m3tal config wizard
```

You will be prompted to provide values for several settings. Key prompts include:

*   **`DOMAIN`**: The public domain name you intend to use for M3TAL services (e.g., `example.com`). This is crucial for Traefik routing to expose services like the M3TAL API (`api.DOMAIN`) and Dashboard (`dash.DOMAIN`). If you do not have a domain configured, `localhost` is a valid default for local access.
*   **`LOCAL_IP`**: The local IP address of your host machine. This is used by services to communicate internally, especially between Docker containers and the host-bound M3TAL API daemon.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL Dashboard container will listen (default: `8082`).
*   **`HTTP_PORT`**: The port on which the M3TAL API daemon will listen (default: `8080`).
*   **`API_TOKEN`**: A secure token for the API. It is recommended to change the default `change_me_api_token` for production environments.
*   **`ADMIN_PASSWORD`**: The initial password for the `admin` user of the M3TAL Dashboard. It is recommended to change the default `admin_pass` immediately.
*   **`PUID`/`PGID`**: User and Group IDs for container permissions, typically `1000` for a standard user.
*   **`TZ`**: Your local timezone (e.g., `America/Denver`).

All these settings are stored in `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

The M3TAL ecosystem uses Traefik as its reverse proxy to manage incoming network requests and route them to the correct services, often using domain names.

This command starts all Docker Compose stacks found in the `/docker/` directory, including the Traefik routing stack:

```bash
m3tal up
```

The `m3tal up` command automatically iterates through all files matching `*-compose.yml` within the `/docker/` symlink (which points to `/opt/m3tal/stack/`) and executes `docker compose` for them using the environment variables from `/etc/m3tal/.env`. This will typically start Traefik, which listens on port 80 (or 443 for HTTPS if configured).

## 5. Start the M3TAL Dashboard

The M3TAL Dashboard provides a web interface for managing your ecosystem.

This command specifically pulls the `m3tal-dashboard` Docker image (if not already present) and starts its container. It also ensures the `m3tal-api.service` daemon is running, which the dashboard relies on for backend communication. The dashboard container is defined in `/docker/m3tal-compose.yml`.

```bash
m3tal dash up
```

## 6. Access the Dashboard

Open your web browser and navigate to the M3TAL Dashboard:

*   If you configured `DOMAIN` as `localhost` or wish to access directly via IP:
    `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address or `localhost`)

*   If you configured `DOMAIN` with a public domain and Traefik is properly set up with corresponding DNS records:
    `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with your configured domain, e.g., `dash.example.com`)

## 7. Log In

Use the following credentials to log in:

*   **Username**: `admin`
*   **Password**: The value you set for `ADMIN_PASSWORD` during the `config wizard` (default: `admin_pass`).

It is strongly recommended to change this password immediately. You can do this using the CLI:

```bash
sudo m3tal dashpass admin new_secure_password
```

## Filesystem Contract

The following paths are critical to the M3TAL Ecosystem:

| Path                        | Purpose                                                                                                                                                                                                                               |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | Primary configuration file. This file contains all environment variables for the M3TAL API daemon and Docker Compose stacks. It is managed by `m3tal config wizard` and `m3tal config set`.                                          |
| `/var/lib/m3tal/state.db`   | SQLite state database for the M3TAL API daemon. This database stores M3TAL's internal state, configuration, and registered services. It is automatically created and managed by the API daemon.                                     |
| `/opt/m3tal/stack/`         | Canonical stack directory. This directory contains M3TAL's core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`) and Traefik's static and dynamic configurations.                                           |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose files and dynamic Traefik configurations. Users should place their custom `*-compose.yml` files here for M3TAL to manage them with `m3tal up`. |
| `/docker/users.json`        | M3TAL Dashboard credential store. This JSON file contains the hashed passwords for Dashboard users. It is managed by the `m3tal dashpass` command.                                                                                 |

## Port Table

These are the default ports used by core M3TAL components:

| Port | Service                            | Access                                                                  |
| :--- | :--------------------------------- | :---------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point           | Public (if exposed). This port handles all incoming HTTP requests for services routed via Traefik. |
| 8080 | M3TAL API daemon (Go)              | Host-local. The M3TAL Dashboard and other internal services communicate with the API daemon on this port. Typically accessed via Traefik (e.g., `api.DOMAIN`) or `http://host.docker.internal:8080` from containers. |
| 8081 | Traefik dashboard                  | Host-local only. Provides a web interface to monitor Traefik's configuration and status. |
| 8082 | M3TAL Dashboard (Python/Flask) | Via Traefik (e.g., `dash.DOMAIN`) or direct IP access (e.g., `http://YOUR_IP:8082`). |

## Firewall Note

If you have a firewall enabled (e.g., UFW) and intend to expose Traefik publicly to handle incoming HTTP requests (on port 80), you must allow the port:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon is a critical backend component that runs as a systemd service called `m3tal-api.service`. You can manage its state and view logs using standard systemd commands:

*   **Check service status**:
    ```bash
    systemctl status m3tal-api
    ```
*   **View live service logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```

Note that `m3tal` CLI commands, such as `m3tal dash up`, often interact with and manage the `m3tal-api.service` automatically to ensure it is running when needed.