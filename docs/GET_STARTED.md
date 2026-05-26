# M3TAL Getting Started Guide

This guide provides a complete, step-by-step process for first-time users to set up the M3TAL Ecosystem.

---

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify their installation with the following commands:

```bash
docker --version
docker compose version
```

## Step 2: Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the APT package manager.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

Initialize your M3TAL environment by running the configuration wizard. This will guide you through setting essential environment variables in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard will be directly accessible on a specific port (default 8082) on your host machine's IP address. Best for simple, LAN-only setups.
    *   `traefik`: The dashboard will be accessible via a subdomain (e.g., `dash.yourdomain.com`) routed by the Traefik reverse proxy. Requires Traefik to be running.
*   **`DOMAIN`**: The base domain name for your services when `DASHBOARD_EXPOSE_MODE` is set to `traefik`. For `local` mode, you can leave this as `localhost`.
*   **`PUID`** and **`PGID`**: The User ID (PUID) and Group ID (PGID) that Docker containers will use. This ensures correct file permissions for volumes. You can find your current user's PUID/PGID with `id -u` and `id -g`.
*   **`TZ`**: Your local timezone (e.g., `America/New_York`). Used by containers for accurate timestamps.
*   **`BASE_STORAGE_PATH`**: The primary path on your host machine where M3TAL-related data and container volumes will be stored (e.g., `/mnt/data`).
*   **`CONFIG_PATH`**: A sub-path within `BASE_STORAGE_PATH` for configuration files (e.g., `/mnt/data/config`).
*   **`MEDIA_PATH`** and **`DOWNLOADS_PATH`**: Sub-paths for media and download storage, typically mounted into containers (e.g., `/mnt/data/media`).
*   **`DASHBOARD_PORT`**: The direct port for the M3TAL Dashboard when `DASHBOARD_EXPOSE_MODE` is `local`. Default is `8082`.
*   **`DASHBOARD_SECRET`**: A random string used to secure the dashboard's session.
*   **`API_TOKEN`**: A token used by the M3TAL Dashboard and other clients to authenticate with the M3TAL API daemon.
*   **`ADMIN_PASSWORD`**: The default password for the `admin` user on the M3TAL Dashboard.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined in the `/docker/` directory. This includes the `routing-compose.yml` which deploys Traefik, the reverse proxy.

```bash
m3tal up
```

This command will pull the necessary Traefik images and start the containers, configuring the core routing layer of your M3TAL ecosystem.

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard. It pulls the latest dashboard image from the registry and starts its container with the appropriate configuration based on your `DASHBOARD_EXPOSE_MODE` setting.

```bash
m3tal dash up
```

This command will download the `m3tal-dashboard` image and initiate the dashboard service.

## Step 6: Open the Browser

Access the M3TAL Dashboard via your web browser. The URL depends on the `DASHBOARD_EXPOSE_MODE` you configured in Step 3:

*   **If `DASHBOARD_EXPOSE_MODE` is `local`**:
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with your host machine's IP address, or use `localhost` if accessing from the same machine).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`**:
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured). Ensure Traefik (started in Step 4) is running and properly configured for your domain.

## Step 7: Log In

The default login credentials for the M3TAL Dashboard are:

*   **Username**: `admin`
*   **Password**: The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard` (default is `admin_pass` if not changed).

### Change Dashboard Password

It is highly recommended to change the default password immediately. Use the following command:

```bash
sudo m3tal dashpass
```

This will prompt you to set a new password for the `admin` user.

---

## M3TAL Key Concepts

### Filesystem Contract

M3TAL establishes a clear filesystem contract for its operation and data storage:

| Path                        | Purpose                                                                 |
| :-------------------------- | :---------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.           |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.      |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's internal Docker Compose files.          |
| `/docker`                   | **Symlink** → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations and where you should place custom compose files. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                |

### M3TAL Ports

The M3TAL ecosystem uses specific ports for its internal and external communications:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if Traefik is exposed)              |
| 8080 | M3TAL API daemon (Go)       | Host-local only (internal communication)    |
| 8081 | Traefik dashboard           | Host-local only (internal Traefik metrics)  |
| 8082 | M3TAL Dashboard container   | Direct port (local mode) or via Traefik (traefik mode) |

### Firewall Note

If Traefik is configured to expose services on port 80 (and/or 443 for HTTPS) to the public internet, you must open these ports in your firewall. For example, using `ufw`:

```bash
sudo ufw allow 80
sudo ufw enable
```

### Service Management

The M3TAL API daemon runs as a systemd service, allowing for standard Linux service management:

*   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```

### Docker and Compose Runtime

M3TAL leverages **Docker Engine** and **Docker Compose V2** as its container orchestration backend.

*   The `m3tal up` command executes `docker compose` operations across all `*-compose.yml` files found within the `/docker/` symlink. This includes core M3TAL services like Traefik (`routing-compose.yml`).
*   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container. It ensures the latest compose files for the dashboard are used and applies the correct compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.
*   To deploy your own services, simply place your Docker Compose files (e.g., `my-app-compose.yml`) into the `/docker/` directory, and M3TAL will manage them with the next `m3tal up` command.

### M3TAL Architecture Overview

*   **CLI binary (`/usr/bin/m3tal`)**: Your primary interface for interacting with the M3TAL ecosystem.
*   **API daemon (`m3tal-api.service`)**: A Go-based backend service running on port 8080. It manages Docker operations, interacts with the `state.db`, and provides an API for the dashboard and other clients.
*   **Dashboard container (`m3tal-dashboard`)**: A Python/Flask application providing a web UI. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway (`routing-compose.yml`)**: A Docker container acting as a reverse proxy, typically on port 80. It routes incoming requests to various services (including the dashboard and API) based on domain names and Docker labels.
*   **Cloudflared**: An optional component within `routing-compose.yml` for creating secure, zero-config tunnels to Cloudflare.

### Dashboard Access Modes

The M3TAL Dashboard offers two distinct access modes, configured by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **Local Mode (`DASHBOARD_EXPOSE_MODE=local`)**:
    *   Directly binds the dashboard container's internal port 8082 to the host machine's port (default `8082`).
    *   No Traefik required for dashboard access.
    *   **Access via**: `http://HOST_IP:8082` or `http://localhost:8082`.
    *   Ideal for initial setup, local testing, or environments where the dashboard only needs to be accessible on the local network.

2.  **Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)**:
    *   Configures the dashboard container with Traefik labels, allowing Traefik to discover and route requests to it.
    *   Assumes Traefik is running and configured with a `DOMAIN`.
    *   **Access via**: `http://dash.YOUR_DOMAIN`.
    *   Best for production environments, domain-based access, and when running multiple services behind Traefik.

### Deployment Lifecycle - Day 2 Operations

To add a new Docker Compose stack to your M3TAL ecosystem:

1.  Create your Docker Compose file (e.g., `my-service-compose.yml`).
2.  Place this file into the `/docker/` directory.
3.  Ensure any required environment variables for your new stack are set in `/etc/m3tal/.env` (using `m3tal config wizard` or `m3tal config set KEY value`).
4.  Run `m3tal up` to start all services, including your new stack. M3TAL will automatically detect and deploy it.