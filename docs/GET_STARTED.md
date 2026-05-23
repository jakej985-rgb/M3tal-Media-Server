# M3TAL Get Started Guide

This guide provides a step-by-step process for setting up M3TAL for the first time.

---

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed and operational on your system.

Verify your Docker installation:

```bash
docker --version && docker compose version
```

---

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

---

## Step 3: Run the Configuration Wizard

Initialize your M3TAL environment by running the configuration wizard. This will create the `/etc/m3tal/.env` file with your system's settings.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for the following information:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard is directly accessible via a port binding (`http://YOUR_IP:8082`). Ideal for LAN-only setups or initial testing. Traefik is not required for dashboard access in this mode.
    *   `traefik`: The dashboard is exposed via Traefik as `dash.DOMAIN`. Requires Traefik to be running and configured.
*   **`DASHBOARD_PORT`**: The host port to expose the M3TAL dashboard on when `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`DOMAIN`**: Your primary domain name, e.g., `example.com`. Used for Traefik routing (e.g., `dash.example.com`, `api.example.com`). Default is `localhost`.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for secure operations. A strong, random value is recommended.
*   **`API_TOKEN`**: An authentication token for the M3TAL API. A strong, random value is recommended.
*   **`ADMIN_PASSWORD`**: The default password for the M3TAL Dashboard `admin` user. This should be changed immediately after initial setup.
*   **`PUID`** and **`PGID`**: The User ID and Group ID that containers will run as, typically `1000`. This ensures proper file permissions on mounted volumes.
*   **`TZ`**: Your timezone, e.g., `America/Denver`.
*   **`BASE_STORAGE_PATH`**: The base directory for all M3TAL-related storage (e.g., `/mnt/m3tal`).
*   Other paths like `MEDIA_PATH`, `CONFIG_PATH`, `DOWNLOADS_PATH` will be relative to `BASE_STORAGE_PATH` by default.

---

## Step 4: Start the Routing Stack

Start the core M3TAL routing stack, including Traefik. This command reads all `*-compose.yml` files located in the `/docker/` symlink and brings up the defined services.

```bash
m3tal up
```

This will deploy the Traefik reverse proxy, which binds to host port `80` (and `443` if configured for HTTPS) to manage routing for services exposed via Traefik.

---

## Step 5: Start the Dashboard

Start the M3TAL Dashboard container. This command will:
1.  Download the necessary M3TAL dashboard Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Read the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Based on the `DASHBOARD_EXPOSE_MODE`, it will start the `m3tal-dashboard` container, applying the appropriate Compose override file (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode).

```bash
m3tal dash up
```

---

## Step 6: Open Browser

Access the M3TAL Dashboard in your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your server).
    *For local testing, you can use `http://localhost:8082`.*

*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3).
    *This requires Traefik to be running via `m3tal up` and proper DNS records for `dash.YOUR_DOMAIN` pointing to your server.*

---

## Step 7: Log in

Log in to the M3TAL Dashboard with the following default credentials:

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (default is `admin_pass`).

**Security Note:** Change the `admin` password immediately after your first login. You can do this via the dashboard's user management interface or by using the CLI:

```bash
sudo m3tal dashpass admin new_strong_password
```

---

## Filesystem Contract

M3TAL establishes a clear filesystem contract for its operation:

| Path                     | Purpose                                                          |
| :----------------------- | :--------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.    |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created by the API daemon.           |
| `/opt/m3tal/stack/`      | Canonical directory for M3TAL's internal stack files (Compose files, Traefik config). |
| `/docker`                | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for placing custom Docker Compose files for new services. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass` or via the Dashboard UI. |

---

## Port Table

M3TAL utilizes the following ports on your host system:

| Port | Service               | Access                                                |
| :--- | :-------------------- | :---------------------------------------------------- |
| `80`   | Traefik HTTP entry point | Public (if Traefik is used for external routing)      |
| `8080` | M3TAL API daemon (Go) | Host-local only (accessed by dashboard and other services) |
| `8081` | Traefik dashboard     | Host-local only (Traefik's internal dashboard)        |
| `8082` | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

---

## Firewall Note

If Traefik is exposed publicly on port `80` (which is the default behavior for `m3tal up` when using Traefik), you must open this port in your firewall. For systems using `ufw`:

```bash
sudo ufw allow 80/tcp
sudo ufw enable # if not already enabled
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage its state using `systemctl` and view its logs with `journalctl`.

Check the status of the M3TAL API daemon:

```bash
systemctl status m3tal-api
```

View real-time logs for the M3TAL API daemon:

```bash
journalctl -u m3tal-api -f
```