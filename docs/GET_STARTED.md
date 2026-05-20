# M3TAL Setup Guide

This guide provides a complete, step-by-step setup for first-time users of the M3TAL Ecosystem.

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your installation:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Install the M3TAL CLI and API daemon using the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

Initialize your M3TAL environment using the interactive configuration wizard. This will generate your primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration options:

*   **`DOMAIN`**: The domain name for your M3TAL services (e.g., `example.com`). This is used by Traefik for routing services like `dash.example.com` or `api.example.com`. Default is `localhost`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard is exposed directly on a host port (`DASHBOARD_PORT`). No Traefik configuration is needed for access.
    *   `traefik`: The dashboard is exposed via Traefik at `dash.${DOMAIN}`. This requires Traefik to be running.
*   **`DASHBOARD_PORT`**: If `DASHBOARD_EXPOSE_MODE` is `local`, this is the direct host port where the dashboard will be accessible. Default is `8082`.
*   **`PUID`** (Process User ID) / **`PGID`** (Process Group ID): The user and group IDs that containers will run as. This ensures correct file permissions for mounted volumes. Typically `1000` for the first non-root user.
*   **`TZ`**: Your timezone (e.g., `America/New_York`).
*   **`BASE_STORAGE_PATH`**: The base directory for all your M3TAL data volumes (e.g., `/mnt/data`). This is where configurations, media, and downloads will reside.
*   **`CONFIG_PATH`**: Path for M3TAL-specific configuration files and state. Defaults to a subdirectory within `BASE_STORAGE_PATH`.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management.
*   **`API_TOKEN`**: An authentication token for accessing the M3TAL API.
*   **`ADMIN_PASSWORD`**: The default password for the `admin` user on the M3TAL Dashboard. **It is critical to change this after initial setup.**

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This includes the `routing-compose.yml` which deploys Traefik, the reverse proxy for domain-based access.

```bash
m3tal up
```

## 5. Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container:

1.  It downloads the necessary Docker Compose files (`m3tal-compose.yml` and its overrides) from GitHub.
2.  It reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  It then starts the dashboard container, applying the correct override (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode) to expose it correctly.

```bash
m3tal dash up
```

## 6. Open M3TAL Dashboard in Browser

Access the M3TAL Dashboard based on your `DASHBOARD_EXPOSE_MODE` configuration:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open your browser to `http://YOUR_SERVER_IP:8082` (or `http://localhost:8082` if accessing from the server itself).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open your browser to `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the value configured in `m3tal config wizard`). Traefik must be running for this to work.

## 7. Log In

The default login credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** The value of `ADMIN_PASSWORD` set during the configuration wizard (default `admin_pass`).

**Change the default password immediately** after your first login:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL establishes a clear filesystem contract for its components and data:

| Path                        | Purpose                                                                                                 |
| :-------------------------- | :------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | **Primary Configuration File**. Contains all M3TAL environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the M3TAL API daemon.                                  |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's internal Docker Compose files and Traefik dynamic configurations.           |
| `/docker`                   | **Symlink** to `/opt/m3tal/stack/`. This is the user-facing path for managing Docker Compose stacks.       |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                                                  |

## Port Table

M3TAL uses the following ports for its services:

| Port | Service               | Access Level                                              |
| :--- | :-------------------- | :-------------------------------------------------------- |
| 80   | Traefik HTTP Entry Point | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed) |
| 8080 | M3TAL API Daemon (Go) | Host-local (accessed internally by dashboard/Traefik)     |
| 8081 | Traefik Dashboard      | Host-local only (for Traefik's own management UI)         |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`dash.${DOMAIN}`) |

## Firewall Note

If your server has a firewall (e.g., UFW), you may need to allow traffic on the ports M3TAL uses for external access:

*   To expose Traefik on port 80:
    ```bash
    sudo ufw allow 80
    ```
*   To expose the M3TAL Dashboard directly on port 8082 (if `DASHBOARD_EXPOSE_MODE=local`):
    ```bash
    sudo ufw allow 8082
    ```

## Service Management

The M3TAL API daemon runs as a systemd service. Use `systemctl` to manage it:

*   **Check API service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View API daemon logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```