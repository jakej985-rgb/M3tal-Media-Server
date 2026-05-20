```markdown
# M3TAL Get Started Guide

Welcome to M3TAL! This guide will walk you through the initial setup of your M3TAL ecosystem, ensuring a smooth start for first-time users.

---

### Step 1: Prerequisites

M3TAL relies on Docker for container orchestration. Before proceeding, ensure that **Docker Engine** and **Docker Compose V2** are installed on your system.

You can verify your installation by running:

```bash
docker --version && docker compose version
```

### Step 2: Install M3TAL via APT

M3TAL is distributed as a single Go binary via APT. Follow these commands to add the repository and install M3TAL:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### Step 3: Run the Configuration Wizard

After installation, run the M3TAL configuration wizard to set up your environment variables. This wizard will guide you through essential settings.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here’s what each means:

-   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    -   `local`: The dashboard is directly exposed on a specific port (default 8082). Ideal for local testing or LAN-only access.
    -   `traefik`: The dashboard is routed through Traefik (if running), making it accessible via a subdomain (e.g., `dash.YOUR_DOMAIN`). Best for domain-based setups.
    -   *Default: `local`*
-   **`DOMAIN`**: The base domain used for Traefik routing (e.g., `yourdomain.com`). If using `DASHBOARD_EXPOSE_MODE=local` and not exposing other services via Traefik, `localhost` is sufficient.
    -   *Default: `localhost`*
-   **`PUID` / `PGID`**: The User ID (UID) and Group ID (GID) that containers will run as. This ensures proper file permissions for mounted volumes. For most users, `1000` is the default for the first non-root user.
    -   *Defaults: `1000` / `1000`*
-   **`BASE_STORAGE_PATH`**: The base directory on your host where M3TAL will store application data, configuration, media, and downloads. Subdirectories like `config`, `media`, and `downloads` will be created here.
    -   *Default: `/mnt`*
-   **`DASHBOARD_SECRET`**: A secret key used to secure the M3TAL Dashboard's session management. Generate a strong, unique value.
    -   *Default: `change_me_immediately`*
-   **`ADMIN_PASSWORD`**: The initial password for the `admin` user to log into the M3TAL Dashboard. **Change this immediately after setup.**
    -   *Default: `admin_pass`*
-   **`API_TOKEN`**: An authentication token required for programmatic access to the M3TAL API. Keep this secure.
    -   *Default: `change_me_api_token`*
-   **`TRAEFIK_WEB_PORT`**: The port Traefik will listen on for incoming HTTP requests (typically 80).
    -   *Default: `80`*
-   **`TZ`**: Your system's timezone (e.g., `America/Denver`).
    -   *Default: `America/Denver`*
-   **`HTTP_PORT`**: The port the M3TAL API daemon (Go service) will listen on internally.
    -   *Default: `8080`*

### Step 4: Start the Routing Stack

Start the core routing stack, which includes Traefik. This command starts all Docker Compose files found in the `/docker/` directory.

```bash
m3tal up
```

This will bring up essential services like Traefik (which handles routing for services exposed via domains) and any other user-defined stacks placed in `/docker/`.

### Step 5: Start the Dashboard

Now, start the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command specifically manages the M3TAL Dashboard. It pulls the latest dashboard Docker image (e.g., `ghcr.io/jakej985-rgb/m3tal-godash:debug`) and starts the container using the appropriate configuration based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

### Step 6: Open Browser

Access the M3TAL Dashboard through your web browser. The URL depends on your `DASHBOARD_EXPOSE_MODE` setting:

-   **If `DASHBOARD_EXPOSE_MODE=local`**:
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your M3TAL host, or use `http://localhost:8082` if accessing from the host itself).
-   **If `DASHBOARD_EXPOSE_MODE=traefik`**:
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Traefik must be running for this mode to work.

### Step 7: Log In

The default login credentials for the M3TAL Dashboard are:

-   **Username**: `admin`
-   **Password**: The `ADMIN_PASSWORD` you set during the configuration wizard (default: `admin_pass`).

**It is strongly recommended to change the admin password immediately after your first login.**
You can change the dashboard password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

M3TAL adheres to a defined filesystem structure for its operation and data storage:

| Path                        | Purpose                                                                                                                                                                                                                             |
| :-------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for the M3TAL system. Contains all environment variables managed by `m3tal config wizard` and `m3tal config set`.                                                                                       |
| `/var/lib/m3tal/state.db`   | SQLite state database. This database is automatically created and managed by the M3TAL API daemon to store internal system state.                                                                                                   |
| `/opt/m3tal/stack/`         | Canonical stack directory. This is where M3TAL stores its core Docker Compose files (`m3tal-compose.yml`, `routing-compose.yml`, etc.) and Traefik dynamic configuration.                                                            |
| `/docker`                   | A symbolic link that points to `/opt/m3tal/stack/`. This is the user-facing path where you should place your custom Docker Compose files for M3TAL to manage (e.g., `my-service-compose.yml`).                                     |
| `/docker/users.json`        | Dashboard credential store. This JSON file contains encrypted user credentials for the M3TAL Dashboard. It is managed by the `m3tal dashpass` command.                                                                               |

## Port Table

The following ports are used by M3TAL components:

| Port | Service                                   | Access                                                                  |
| :--- | :---------------------------------------- | :---------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point                  | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed) |
| 8080 | M3TAL API daemon (Go application)         | Host-local only (accessed by containers via `host.docker.internal`)     |
| 8081 | Traefik dashboard                         | Host-local only (accessible at `http://localhost:8081` on the host)     |
| 8082 | M3TAL Dashboard (Python/Flask container) | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`dash.DOMAIN` if `traefik` mode) |

## Firewall Note

If Traefik is exposed on port 80 (or 443 for HTTPS) and you have a firewall enabled (e.g., UFW), you will need to allow traffic on these ports. For example, to allow HTTP traffic on port 80:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon runs as a systemd service called `m3tal-api`. You can manage its state using standard `systemctl` commands:

-   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```
-   **View live logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```
-   **Restart the service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```

---

## M3TAL Dashboard Access Modes

The M3TAL Dashboard has two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Mode 1: Local Access (`DASHBOARD_EXPOSE_MODE=local`)

This is the default and recommended mode for first-time users, LAN-only setups, or local testing.

-   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
-   **Mechanism**: M3TAL uses the `m3tal-compose.local.yml` Docker Compose override file, which adds a direct port binding to the dashboard container (e.g., `${DASHBOARD_PORT:-8082}:8082`). This exposes the dashboard directly on your host's IP address.
-   **Access**: Open your browser to `http://YOUR_HOST_IP:8082` (or `http://localhost:8082`).
-   **Requirements**: No Traefik configuration is required for dashboard access in this mode.

### Mode 2: Traefik Access (`DASHBOARD_EXPOSE_MODE=traefik`)

This mode integrates the dashboard behind the Traefik reverse proxy, making it accessible via a domain or subdomain. This is suitable for domain-based setups where you manage multiple services behind Traefik.

-   **Configuration**: Set `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
-   **Mechanism**: M3TAL uses the `m3tal-compose.traefik.yml` Docker Compose override file. This file adds Traefik labels to the dashboard container, instructing Traefik to route requests for `dash.${DOMAIN}` to the dashboard container on its internal port 8082.
-   **Access**: Open your browser to `http://dash.YOUR_DOMAIN` (ensure your `DOMAIN` variable is set correctly in `/etc/m3tal/.env` and points to your M3TAL host).
-   **Requirements**: Traefik must be running via `m3tal up` for this mode to function.

## Docker / Compose Runtime

M3TAL orchestrates services using **Docker Engine** and **Docker Compose V2**. These are hard dependencies.

-   The `m3tal up` command acts as a wrapper for `docker compose`. It discovers and starts all Docker Compose projects (`*-compose.yml` files) located in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`).
-   The `m3tal dash up` command specifically manages the M3TAL Dashboard. It performs the following steps:
    1.  Downloads the latest base `m3tal-compose.yml` and relevant override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from the official GitHub repository.
    2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
    3.  Starts the dashboard container using `docker compose` with the appropriate override file applied based on the configured exposure mode.
-   Users can deploy their own Docker Compose stacks by simply placing their `*-compose.yml` files into the `/docker/` directory. M3TAL will automatically pick them up on the next `m3tal up` command.

## Deployment Lifecycle — Day 2 Operations

To add a new Docker Compose stack to your M3TAL ecosystem:

1.  **Place your compose file**: Create or copy your Docker Compose file (e.g., `my-application-compose.yml`) into the `/docker/` directory.
2.  **Set required variables**: If your new stack uses environment variables, ensure they are defined in `/etc/m3tal/.env` using `m3tal config wizard` or `m3tal config set KEY value`.
3.  **Start all stacks**: Run `m3tal up` to deploy your new stack alongside existing ones.

## Traefik Routing Architecture

Traefik, deployed as a container via `routing-compose.yml`, functions as the primary reverse proxy for M3TAL.

-   It binds to port 80 (and optionally 443 for HTTPS) on the host system, serving as the HTTP entry point for all incoming web traffic.
-   Traefik automatically discovers services by inspecting Docker container labels.
-   It loads dynamic routing configurations from files placed in `/docker/dynamic/` (which points to `/etc/traefik/dynamic` inside the Traefik container). These configurations are hot-reloaded without restarting Traefik.
-   **API Routing**: Traefik is configured to route `api.YOUR_DOMAIN` to the M3TAL API daemon, which runs on `http://host.docker.internal:8080` on the host. This is defined in `routing-compose.yml` and a dynamic configuration file like `dynamic/api.yml`.
-   **Dashboard Routing**: If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, Traefik routes `dash.YOUR_DOMAIN` to the M3TAL Dashboard container, utilizing the Traefik labels defined in `m3tal-compose.traefik.yml`.
```