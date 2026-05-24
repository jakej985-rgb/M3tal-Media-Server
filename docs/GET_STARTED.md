# Getting Started with M3TAL

This guide provides a complete, step-by-step process for first-time users to set up and run the M3TAL Ecosystem.

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your installation with:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the `m3tal` CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The `m3tal config wizard` command will guide you through the initial setup, creating the primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

During the wizard, you will be prompted for several critical configuration values:

*   **`DOMAIN`**: The base domain name for your services (e.g., `yourdomain.com`). This is used by Traefik for routing if you choose the Traefik dashboard mode. If you don't have a domain, `localhost` is a suitable default.
*   **`DASHBOARD_EXPOSE_MODE`**:
    *   `local` (default): The M3TAL Dashboard will be directly accessible via a host port (default 8082). Recommended for initial setup or LAN-only use.
    *   `traefik`: The M3TAL Dashboard will be routed via Traefik at `dash.YOUR_DOMAIN`. Requires Traefik to be running.
*   **`DASHBOARD_PORT`**: (Only prompted if `DASHBOARD_EXPOSE_MODE` is `local`). The host port to expose the M3TAL Dashboard on. Default is 8082.
*   **`PUID`** / **`PGID`**: The User ID and Group ID that Docker containers will run as inside the M3TAL ecosystem. This ensures proper file permissions for mounted volumes. You can typically find these using `id -u` and `id -g` for your current user.
*   **`BASE_STORAGE_PATH`**: The base directory on your host where all M3TAL-managed data will be stored (e.g., `/mnt/m3tal-data`). This path will contain subdirectories for configuration, media, downloads, etc.
*   **`CONFIG_PATH`**: Subdirectory for application configuration.
*   **`MEDIA_PATH`**: Subdirectory for media files.
*   **`DOWNLOADS_PATH`**: Subdirectory for downloads.
*   **`TZ`**: Your timezone (e.g., `America/New_York`). Used by containers for correct time display.
*   **`DASHBOARD_SECRET`**: A strong, random string used to secure dashboard sessions.
*   **`API_TOKEN`**: A strong, random string used to secure API access.

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command initializes and starts all Docker Compose stacks found in the `/docker/` directory. This includes the core routing stack, which deploys Traefik, your reverse proxy.

```bash
m3tal up
```

This command will:
*   Create a Docker network named `proxy`.
*   Deploy Traefik from `routing-compose.yml`, binding host port 80.
*   Configure Traefik to route `api.YOUR_DOMAIN` to the M3TAL API daemon (running on host port 8080).

## 5. Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It retrieves the necessary Docker Compose files and starts the dashboard according to your `DASHBOARD_EXPOSE_MODE` setting.

```bash
m3tal dash up
```

This command will:
*   Download the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` files.
*   Read the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
*   Start the `m3tal-dashboard` container using the appropriate compose override file (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode).

## 6. Open the Browser

Access the M3TAL Dashboard in your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server, or use `http://localhost:8082` if accessing from the server itself).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain configured in your `.env` file). Traefik must be running for this mode to work.

## 7. Log In

The default login credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** `admin_pass` (This is derived from the `ADMIN_PASSWORD` variable in `/etc/m3tal/.env`)

**It is critical to change this default password immediately.** Use the following command:

```bash
sudo m3tal dashpass
```
You will be prompted to enter a new password for the `admin` user.

---

## Filesystem Contract

The following paths are central to the M3TAL ecosystem:

| Path                     | Purpose                                                              |
| :----------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the API daemon.   |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains core compose files and Traefik config. |
| `/docker`                | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.             |

## Port Table

M3TAL utilizes the following network ports:

| Port | Service                     | Access                                                      |
| :--- | :-------------------------- | :---------------------------------------------------------- |
| 80   | Traefik HTTP entry point    | Public (if Traefik is exposed)                              |
| 8080 | M3TAL API daemon (Go)       | Host-local only (internal communication with Dashboard/CLI) |
| 8081 | Traefik Dashboard           | Host-local only (for monitoring Traefik itself)             |
| 8082 | M3TAL Dashboard container   | Direct port (local mode) or via Traefik (traefik mode)      |

## Firewall Note

If you have a firewall enabled (e.g., `ufw`) and are using Traefik to expose services (including the dashboard in `traefik` mode), you must allow inbound traffic on port 80:

```bash
sudo ufw allow 80
```

If you configured Traefik to use HTTPS (port 443), also allow that port:

```bash
sudo ufw allow 443
```

## Service Management

The M3TAL API daemon runs as a systemd service, `m3tal-api.service`. You can manage it using standard `systemctl` commands:

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