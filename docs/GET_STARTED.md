```markdown
# M3TAL Get Started Guide

This guide provides step-by-step instructions for installing and setting up M3TAL for the first time.

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your installation:

```bash
docker --version && docker compose version
```

If these commands do not return version information, please install Docker Engine and Docker Compose V2 according to the official Docker documentation for your operating system before proceeding.

## Step 2: Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the `m3tal` CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

After installation, the M3TAL API daemon (`m3tal-api.service`) will be started automatically.

## Step 3: Run the configuration wizard

The configuration wizard helps set up the primary environment variables in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt for several settings. Here’s an explanation of key prompts:

*   **`DASHBOARD_EXPOSE_MODE` (local/traefik):**
    *   `local`: The M3TAL Dashboard will be directly accessible via `http://YOUR_IP:8082`. This is the default and recommended for initial setup and LAN-only use.
    *   `traefik`: The Dashboard will be routed through the Traefik reverse proxy and accessible via a domain name (e.g., `http://dash.YOUR_DOMAIN`). Requires Traefik to be running and a `DOMAIN` configured.
*   **`DOMAIN` (e.g., `localhost`):** The base domain for services exposed via Traefik. Defaults to `localhost`. Used for `dash.DOMAIN`, `api.DOMAIN`, etc.
*   **`DASHBOARD_PORT` (e.g., `8082`):** The port on which the M3TAL Dashboard container exposes its service. When `DASHBOARD_EXPOSE_MODE=local`, this port will be directly mapped to your host.
*   **`PUID`/`PGID` (e.g., `1000`):** The User ID and Group ID that containers will run as. Use `id -u` and `id -g` to find your current user's IDs. Essential for correct file permissions.
*   **`TZ` (e.g., `America/Denver`):** Your timezone. Used by containers for accurate timestamps.
*   **`BASE_STORAGE_PATH` (e.g., `/mnt/m3tal/data`):** The root directory for M3TAL's data and configuration files. This path is mounted into containers.
*   **`DASHBOARD_SECRET`:** A secret key used by the dashboard for session management. Generated automatically if not provided.
*   **`API_TOKEN`:** An API token for authentication with the M3TAL API. Generated automatically if not provided.

Complete the prompts, providing values appropriate for your environment.

## Step 4: Start the routing stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory. This includes the core routing stack which utilizes Traefik as a reverse proxy.

```bash
m3tal up
```

This command will:
*   Pull required Docker images for the routing stack (e.g., Traefik, Cloudflared).
*   Create and start containers defined in `routing-compose.yml` and any other `*-compose.yml` files present in `/docker/`.
*   Traefik will start on host port 80 (HTTP).

## Step 5: Start the dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest M3TAL Dashboard Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container, applying the appropriate Compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE` setting.

## Step 6: Access the Dashboard

Open your web browser and navigate to the M3TAL Dashboard URL:

*   **If `DASHBOARD_EXPOSE_MODE` is set to `local` (default):**
    *   Navigate to `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your server).
    *   The `DASHBOARD_PORT` value from `/etc/m3tal/.env` determines the port.
*   **If `DASHBOARD_EXPOSE_MODE` is set to `traefik`:**
    *   Navigate to `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the `DOMAIN` configured in Step 3).
    *   Traefik must be running for this mode to work.

## Step 7: Log in

The default credentials for the M3TAL Dashboard are:
*   **Username:** `admin`
*   **Password:** `admin_pass`

It is strongly recommended to change the default password immediately after your first login.
To change the dashboard password, use the following CLI command:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the `admin` user, updating the `/docker/users.json` file.

---

## Filesystem Contract

M3TAL establishes a clear filesystem contract for its operation. Understanding these paths is crucial for configuration and troubleshooting.

| Path                     | Purpose                                                                                                                                                                                                                            |
| :----------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | **Primary configuration file.** Contains environment variables used by the M3TAL API daemon and Docker Compose stacks. Managed by `m3tal config wizard` and `m3tal config set`.                                                        |
| `/var/lib/m3tal/state.db`| **SQLite state database.** Auto-created and managed by the M3TAL API daemon. Stores the state of Docker services, configurations, and other operational data.                                                                          |
| `/opt/m3tal/stack/`      | **Canonical stack directory.** This is where M3TAL stores its core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`) and Traefik dynamic configuration.                                                            |
| `/docker`                | **Symlink → `/opt/m3tal/stack/`**. This is the user-facing path where you should place your custom Docker Compose files (`*-compose.yml`) for M3TAL to manage. All `m3tal up` operations target this directory.                       |
| `/docker/users.json`     | **Dashboard credential store.** Stores the hashed credentials for M3TAL Dashboard users. Managed by the `m3tal dashpass` command.                                                                                                    |

## Port Table

This table outlines the default ports used by M3TAL components.

| Port | Service                               | Access                                                                                                                                                            |
| :--- | :------------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (if `TRAEFIK_WEB_PORT` is 80 and Traefik is running). Used for `http://dash.DOMAIN` and other domain-based services.                                           |
| 8080 | M3TAL API daemon (Go)                 | Host-local. This port is for internal communication between the M3TAL Dashboard container and the host API, or for CLI commands. Not typically exposed publicly. |
| 8081 | Traefik dashboard (internal)          | Host-local only. Accessible via `http://localhost:8081` on the server where Traefik is running.                                                                    |
| 8082 | M3TAL Dashboard (Python/Flask)        | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`).                                                              |

## Firewall Note

If Traefik is exposed on port 80 (the default for domain-based routing) and you have a firewall enabled (e.g., UFW), you will need to allow incoming traffic on that port:

```bash
sudo ufw allow 80/tcp
```

Adjust the port if you have configured Traefik to listen on a different one.

## Service Management

The M3TAL API daemon runs as a systemd service called `m3tal-api.service`. You can manage its state using standard `systemctl` commands:

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **Stop the service:**
    ```bash
    sudo systemctl stop m3tal-api
    ```
*   **Start the service:**
    ```bash
    sudo systemctl start m3tal-api
    ```
```