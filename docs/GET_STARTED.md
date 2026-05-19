# M3TAL Ecosystem - Get Started Guide

This guide will walk you through the initial setup of the M3TAL ecosystem.

## 1. Prerequisites

Before proceeding, ensure you have the following installed on your system:

*   **Docker Engine:** The core containerization platform.
*   **Docker Compose V2:** The tool for defining and running multi-container Docker applications.

Verify your installation by running the following command:

```bash
docker --version && docker compose version
```

If these commands do not return version information, please refer to the official Docker documentation for installation instructions.

## 2. Install M3TAL via APT

M3TAL is distributed via an APT repository for easy installation and updates.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up essential parameters.

```bash
sudo m3tal config wizard
```

You will be presented with several prompts:

*   **`DASHBOARD_PORT`**: The port the M3TAL dashboard will listen on. (Default: `8082`)
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is accessed.
    *   `local`: Exposes the dashboard directly on the specified `DASHBOARD_PORT`. Ideal for LAN-only access.
    *   `traefik`: Exposes the dashboard through Traefik, typically via a subdomain. Requires Traefik to be running.
*   **`HTTP_PORT`**: The port the M3TAL API daemon will listen on. (Default: `8080`)
*   **`STATE_DIR`**: The directory where M3TAL stores its state. (Default: `./state`)
*   **`LOG_LEVEL`**: The verbosity of M3TAL logs (e.g., `info`, `debug`). (Default: `info`)
*   **`DASHBOARD_SECRET`**: A secret used for session management within the dashboard. **Change this from the default.** (Default: `change_me_immediately`)
*   **`API_TOKEN`**: A token used for API authentication. **Change this from the default.** (Default: `change_me_api_token`)
*   **`ADMIN_PASSWORD`**: The initial password for the dashboard administrator. **Change this from the default.** (Default: `admin_pass`)
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL uses. (Default: `m3tal`)
*   **`LOCAL_IP`**: The local IP address of your host. (Default: `127.0.0.1`)
*   **`DOMAIN`**: The domain name you intend to use for accessing services. (Default: `localhost`)
*   **`VPN_USER`**: Username for VPN connections (if applicable).
*   **`VPN_PASSWORD`**: Password for VPN connections (if applicable).
*   **`BASE_STORAGE_PATH`**: The base directory for storing M3TAL data. (Default: `./data`)
*   **`MEDIA_PATH`**: Path for media storage.
*   **`CONFIG_PATH`**: Path for configuration files.
*   **`DOWNLOADS_PATH`**: Path for downloads.
*   **`PUID` / `PGID`**: The User ID and Group ID for running containers. Usually matches your user's ID. (Default: `1000`)
*   **`TZ`**: Your timezone (e.g., `America/Denver`). (Default: `America/Denver`)
*   **`TRAEFIK_WEB_PORT`**: The port Traefik uses for HTTP traffic. (Default: `80`)
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik uses for HTTPS traffic. (Default: `443`)
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's own dashboard listens on. (Default: `8080`)
*   **`DEBUG_MODE`**: Enable debug logging. (Default: `false`)
*   **`METRICS_ENABLED`**: Enable metrics collection. (Default: `true`)

The wizard will save these settings to `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

This command initiates the core routing infrastructure, primarily Traefik, which acts as your API gateway.

```bash
m3tal up
```

This command orchestrates all Docker Compose files found in `/docker/` that are not explicitly excluded. This typically includes `routing-compose.yml` which deploys Traefik and other essential routing components.

## 5. Start the Dashboard

The M3TAL dashboard provides a web interface for managing your system.

```bash
m3tal dash up
```

This command will:
1.  Fetch the latest dashboard Docker image (`ghcr.io/jakej985-rgb/m3tal-godash:debug`).
2.  Read the `DASHBOARD_EXPOSE_MODE` setting from your `/etc/m3tal/.env` file.
3.  Start the dashboard container, applying the appropriate configuration (direct port binding for `local` mode or Traefik labels for `traefik` mode).

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local`, access it at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or use `http://localhost:8082` if on the same machine).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik` and Traefik is running, access it at: `http://dash.DOMAIN` (replace `dash.DOMAIN` with your configured domain for the dashboard).

## 7. Log In

Upon first access, you will be prompted to log in.

*   **Default Username:** `admin`
*   **Default Password:** `admin_pass` (This is the value of `ADMIN_PASSWORD` from the configuration wizard. **It is strongly recommended to change this immediately.**)

To change the administrator password after logging in:

```bash
sudo m3tal dashpass
```

This command will prompt you for the new password.

---

## Filesystem Contract

Understanding the M3TAL filesystem layout is crucial for managing your installation.

| Path                      | Purpose                                                                                                     |
| :------------------------ | :---------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | The primary configuration file. Contains all environment variables for M3TAL and its components. Managed by the `m3tal config wizard` and `m3tal config set` commands. |
| `/var/lib/m3tal/state.db` | SQLite database used by the M3TAL API daemon to store system state, service configurations, and other persistent data. Auto-created by the API daemon on first run. |
| `/opt/m3tal/stack/`       | The canonical directory containing all Docker Compose files and Traefik configuration for the M3TAL system. |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all M3TAL stack operations (e.g., `m3tal up`). |
| `/docker/users.json`      | Stores the dashboard credentials. Managed by the `m3tal dashpass` command and used by the dashboard container. |

---

## Port Table

The following ports are utilized by the M3TAL system.

| Port | Service                     | Access                                            |
| :--- | :-------------------------- | :------------------------------------------------ |
| 80   | Traefik (HTTP Entrypoint)   | Public (when `DASHBOARD_EXPOSE_MODE=traefik`)     |
| 8080 | M3TAL API Daemon (Go)       | Host-local (accessible by containers and host)    |
| 8081 | Traefik Dashboard           | Host-local only                                   |
| 8082 | M3TAL Dashboard             | Direct port (if `local` mode) or via Traefik      |

---

## Firewall Configuration

If Traefik is exposed to the internet (typically via `traefik` mode or if `TRAEFIK_WEB_PORT` is not localhost), ensure port 80 is open in your firewall.

If you are using `ufw`:

```bash
sudo ufw allow 80
```

---

## Service Management (API Daemon)

The M3TAL API daemon runs as a systemd service. You can manage and monitor it as follows:

**Check status:**

```bash
systemctl status m3tal-api
```

**Restart the service:**

```bash
sudo systemctl restart m3tal-api
```

**View logs in real-time:**

```bash
sudo journalctl -u m3tal-api -f
```