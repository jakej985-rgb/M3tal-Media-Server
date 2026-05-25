# M3TAL Get Started Guide

This guide provides the necessary steps to install and set up M3TAL for first-time users.

---

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
You can verify their installation and version with the following command:

```bash
docker --version && docker compose version
```

Expected output (versions may vary):
```
Docker version 24.0.7, build afdd53b
Docker Compose version v2.23.3-desktop.1
```

## Step 2: Install M3TAL via APT

Execute the following commands in your terminal to install the M3TAL CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, run the M3TAL configuration wizard to set up essential environment variables. This creates the primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's an explanation of key prompts:

*   **`DASHBOARD_EXPOSE_MODE` (local/traefik):**
    *   `local` (default): The M3TAL Dashboard will be directly accessible via `http://YOUR_IP:8082`. This mode does not require Traefik for dashboard access.
    *   `traefik`: The M3TAL Dashboard will be accessible via `http://dash.YOUR_DOMAIN` (e.g., `http://dash.example.com`). This mode requires Traefik to be running.
*   **`DASHBOARD_PORT` (default: 8082):** The port on which the M3TAL Dashboard will be exposed if `DASHBOARD_EXPOSE_MODE` is set to `local`.
*   **`DOMAIN` (default: localhost):** Your primary domain name. Used by Traefik for routing services (e.g., `api.YOUR_DOMAIN`, `dash.YOUR_DOMAIN`).
*   **`PUID` (default: 1000):** The User ID that M3TAL containers will use for file permissions. Ensure this matches the user ID of the user managing the files on the host.
*   **`PGID` (default: 1000):** The Group ID that M3TAL containers will use for file permissions. Ensure this matches the group ID of the user managing the files on the host.
*   **`TZ` (default: America/Denver):** Your local timezone.
*   **`ADMIN_PASSWORD` (default: admin_pass):** The initial password for the default `admin` user of the M3TAL Dashboard.
*   **`DASHBOARD_SECRET` (default: change_me_immediately):** A secret key used by the M3TAL Dashboard for security purposes. Change this from the default.
*   **`API_TOKEN` (default: change_me_api_token):** An API token for accessing the M3TAL API directly. Change this from the default.
*   **`BASE_STORAGE_PATH`, `CONFIG_PATH`, `MEDIA_PATH`, `DOWNLOADS_PATH`:** Paths on your host system where Docker volumes will be mounted for persistent data storage.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined by `*-compose.yml` files located in the `/docker/` directory. This includes the core routing stack (Traefik).

```bash
m3tal up
```

This command will:
1.  Initialize the Docker `proxy` network if it doesn't exist.
2.  Start the Traefik reverse proxy container, binding to host port 80.
3.  Start the Cloudflared tunnel container (if configured).

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command will:
1.  Download the latest `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Read the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container, applying the appropriate override configuration based on the `DASHBOARD_EXPOSE_MODE` setting.

## Step 6: Open Browser

Access the M3TAL Dashboard in your web browser. The URL depends on your `DASHBOARD_EXPOSE_MODE` setting from Step 3:

*   If `DASHBOARD_EXPOSE_MODE` is `local`:
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server).
    For local access on the server itself, use `http://localhost:8082`.
*   If `DASHBOARD_EXPOSE_MODE` is `traefik`:
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3).
    *Note: Traefik must be running for this mode to work.*

## Step 7: Log In

The default credentials for the M3TAL Dashboard are:
*   **Username:** `admin`
*   **Password:** The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard` (default: `admin_pass`).

To change the dashboard password for the `admin` user:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter and confirm a new password.

---

## Filesystem Contract

M3TAL establishes a clear filesystem structure for its operation and your Docker stacks:

| Path                        | Purpose                                                                                |
| :-------------------------- | :------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for the M3TAL ecosystem. Managed by `m3tal config wizard`.  |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the M3TAL API daemon.               |
| `/opt/m3tal/stack/`         | Canonical directory containing M3TAL's core Docker Compose files and Traefik config.   |
| `/docker`                   | Symlink that points to `/opt/m3tal/stack/`. This is the user-facing path for placing your Docker Compose files. |
| `/docker/users.json`        | Dashboard credential store. Managed by the `m3tal dashpass` command.                   |

---

## Port Table

M3TAL services use the following network ports:

| Port | Service                                  | Access                                                  |
| :--- | :--------------------------------------- | :------------------------------------------------------ |
| 80   | Traefik HTTP entry point                 | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed) |
| 8080 | M3TAL API daemon (Go)                    | Host-local                                              |
| 8081 | Traefik dashboard (admin interface)      | Host-local only (e.g., `http://127.0.0.1:8081`)         |
| 8082 | M3TAL Dashboard (Python/Flask)           | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

---

## Firewall Note

If you expose Traefik (port 80) to the internet or your local network, you must allow incoming connections on your firewall. For `ufw` (Uncomplicated Firewall), use:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon runs as a `systemd` service. You can manage its state and view logs using standard `systemctl` and `journalctl` commands:

*   **Check status:**
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