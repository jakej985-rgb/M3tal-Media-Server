# GET_STARTED.md

This guide provides step-by-step instructions for a first-time setup of the M3TAL Ecosystem.

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
You can verify their installation and version with the following commands:

```bash
docker --version && docker compose version
```

Expected output will show versions for Docker Engine and Docker Compose V2.

## 2. Install M3TAL via APT

Execute the following commands in your terminal to install the M3TAL CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

Initialize your M3TAL environment and set essential configuration variables using the interactive wizard:

```bash
sudo m3tal config wizard
```

The wizard will prompt you for the following information:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard will be accessible.
    *   `local` (default): The dashboard container will directly bind to a host port (default 8082), allowing access via `http://YOUR_IP:8082` or `http://localhost:8082`. This mode does not require Traefik.
    *   `traefik`: The dashboard will be routed through the Traefik reverse proxy, accessible via a domain like `http://dash.DOMAIN`. This requires Traefik to be running.
*   **`DOMAIN`**: The base domain for services exposed via Traefik (e.g., `example.com`). If `DASHBOARD_EXPOSE_MODE` is `traefik`, the dashboard will be at `dash.DOMAIN`. Defaults to `localhost`.
*   **`PUID`** and **`PGID`**: The User ID and Group ID that containers will use for file permissions when accessing host volumes. This prevents permission issues with shared storage. Typically, `1000` for both is suitable for the first non-root user.
*   **`BASE_STORAGE_PATH`**: The root directory on your host where all M3TAL-managed persistent data (e.g., configurations, media files, downloads) will be stored.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for security purposes. Generate a strong, random string.
*   **`API_TOKEN`**: An API token required for authentication when interacting with the M3TAL API daemon. Generate a strong, random string.
*   **`ADMIN_PASSWORD`**: The initial password for the `admin` user of the M3TAL Dashboard. Set a strong password.
*   **`TZ`**: Your local timezone (e.g., `America/New_York`). This ensures logs and scheduled tasks reflect your correct time.

Other prompts will set various system-wide environment variables, which you can typically accept the defaults for initially.

## 4. Start the Routing Stack

Start the core routing stack, which includes the Traefik reverse proxy:

```bash
m3tal up
```

This command processes all Docker Compose files (any `*-compose.yml` file) found in the `/docker/` directory. It specifically starts the Traefik reverse proxy and any configured Cloudflared tunnels as defined in `routing-compose.yml`. Traefik will bind to port 80 on your host to handle incoming HTTP requests.

## 5. Start the Dashboard

Start the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command downloads the latest M3TAL Dashboard Docker image and its associated Compose files. It reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env` and starts the dashboard container with the corresponding configuration (either directly binding to a host port or configuring Traefik labels for routing).

## 6. Open the Browser

Access the M3TAL Dashboard through your web browser based on your chosen `DASHBOARD_EXPOSE_MODE`:

*   **If `DASHBOARD_EXPOSE_MODE` is `local`:**
    Open your browser to `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server) or `http://localhost:8082` if accessing from the server itself.

*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open your browser to `http://dash.DOMAIN` (e.g., `http://dash.example.com` if your `DOMAIN` is `example.com`). Ensure DNS records for `dash.DOMAIN` point to your server's IP address if accessing externally.

## 7. Log In

The default login credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (Step 3).

**To change the admin password:**
After logging in, or at any time from the command line, you can change the admin password using:

```bash
sudo m3tal dashpass
```

## Filesystem Contract

The M3TAL ecosystem maintains a strict filesystem contract for its operational files:

| Path                     | Purpose                                                                                                 |
| :----------------------- | :------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`        | Primary configuration file, storing environment variables. Managed by `m3tal config wizard`.            |
| `/var/lib/m3tal/state.db`| SQLite state database, used by the API daemon to manage service states. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`      | Canonical directory for M3TAL's core Docker Compose files and Traefik dynamic configurations.           |
| `/docker`                | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for managing Docker stacks and Traefik configurations. |
| `/docker/users.json`     | Stores M3TAL Dashboard user credentials. Managed by `m3tal dashpass`.                                   |

## Port Table

The following ports are used by M3TAL components on the host:

| Port | Service                                      | Access                                         |
| :--- | :------------------------------------------- | :--------------------------------------------- |
| 80   | Traefik HTTP entry point (routing-compose.yml) | Public (if `DASHBOARD_EXPOSE_MODE=traefik`)    |
| 8080 | M3TAL API daemon (Go)                        | Host-local (accessed internally by dashboard)  |
| 8081 | Traefik Dashboard (admin interface)          | Host-local only (for Traefik's own dashboard)  |
| 8082 | M3TAL Dashboard (Python/Flask)               | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

## Firewall Note

If your server has a firewall (e.g., `ufw`) and you are using Traefik or exposing the dashboard directly, you may need to open the relevant ports. For Traefik, at minimum, allow port 80:

```bash
sudo ufw allow 80
```

If you configured `DASHBOARD_EXPOSE_MODE=local`, you might also need to open port 8082:

```bash
sudo ufw allow 8082
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard `systemctl` and `journalctl` commands:

*   **Check the status of the M3TAL API daemon:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View live logs for the M3TAL API daemon:**
    ```bash
    journalctl -u m3tal-api -f
    ```