# M3TAL Get Started Guide

This guide provides the essential steps for setting up M3TAL on your system.

## 1. Prerequisites

Before you begin, ensure you have the following installed:

- **Docker Engine**
- **Docker Compose V2**

You can verify their installation with the following commands:

```bash
docker --version
docker compose version
```

## 2. Install M3TAL

M3TAL is installed via APT using the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up essential environment variables.

Run the wizard:

```bash
sudo m3tal config wizard
```

You will be prompted for the following values:

- **`DASHBOARD_PORT`**: The port the M3TAL dashboard will listen on. Defaults to `8082`.
- **`DASHBOARD_EXPOSE_MODE`**: How the dashboard is exposed. Choose `local` for direct access or `traefik` for access via the Traefik reverse proxy. Defaults to `local`.
- **`HTTP_PORT`**: The port for the M3TAL API daemon. Defaults to `8080`.
- **`STATE_DIR`**: The directory for M3TAL's state. Defaults to `./state` (relative to M3TAL's current working directory when running commands).
- **`LOG_LEVEL`**: The logging verbosity. Options include `debug`, `info`, `warn`, `error`. Defaults to `info`.
- **`DASHBOARD_SECRET`**: A secret key used for the dashboard session. **Change this from the default `change_me_immediately`**.
- **`API_TOKEN`**: A token for API authentication. **Change this from the default `change_me_api_token`**.
- **`ADMIN_PASSWORD`**: The password for the dashboard admin user. **Change this from the default `admin_pass`**.
- **`NETWORK_NAME`**: The Docker network name M3TAL will use. Defaults to `m3tal`.
- **`LOCAL_IP`**: The local IP address of your host machine. Defaults to `127.0.0.1`.
- **`DOMAIN`**: The domain name to use for accessing services via Traefik. Defaults to `localhost`.
- **`VPN_USER`**: Username for VPN connections (if applicable). Defaults to `user`.
- **`VPN_PASSWORD`**: Password for VPN connections (if applicable). Defaults to `password`.
- **`BASE_STORAGE_PATH`**: The base directory for M3TAL data. Defaults to `./data`.
- **`MEDIA_PATH`**: Path for media files. Defaults to `./data/media`.
- **`CONFIG_PATH`**: Path for configuration files. Defaults to `./data/config`.
- **`DOWNLOADS_PATH`**: Path for downloads. Defaults to `./data/downloads`.
- **`PUID`**: The User ID for Docker container processes. Defaults to `1000`.
- **`PGID`**: The Group ID for Docker container processes. Defaults to `1000`.
- **`TZ`**: Your timezone. Defaults to `America/Denver`.
- **`TRAEFIK_WEB_PORT`**: The host port for Traefik's web entrypoint. Defaults to `80`.
- **`TRAEFIK_WEBHTTPS_PORT`**: The host port for Traefik's HTTPS entrypoint. Defaults to `443`.
- **`TRAEFIK_DASHBOARD_PORT`**: The host port for Traefik's dashboard. Defaults to `8080`.
- **`DEBUG_MODE`**: Enable or disable debug mode. Defaults to `false`.
- **`METRICS_ENABLED`**: Enable or disable metrics collection. Defaults to `true`.

These settings are saved in `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

This command starts the Traefik reverse proxy and other core routing components. It processes all Docker Compose files located in the `/docker/` directory.

```bash
m3tal up
```

## 5. Start the Dashboard

This command specifically pulls the M3TAL dashboard image and starts the dashboard container based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

```bash
m3tal dash up
```

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard's address:

- If `DASHBOARD_EXPOSE_MODE` is set to `local` (default), access the dashboard at:
  `http://YOUR_IP:8082`
  (Replace `YOUR_IP` with your server's IP address or `localhost` if running locally.)

- If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, and Traefik is running, you can access the dashboard at:
  `http://dash.DOMAIN`
  (Replace `DOMAIN` with the domain you configured in `/etc/m3tal/.env`.)

## 7. Log In

Use the following default credentials to log into the M3TAL dashboard:

- **Username**: `admin`
- **Password**: `admin_pass` (This is the default value for `ADMIN_PASSWORD` from the wizard. **You should have changed this during the wizard.**)

To change the admin password after initial setup, run:

```bash
sudo m3tal dashpass
```

## Filesystem Contract

M3TAL utilizes a specific directory structure and configuration files:

| Path                      | Purpose                                                      | Managed By      |
|---------------------------|--------------------------------------------------------------|-----------------|
| `/etc/m3tal/.env`         | Primary M3TAL environment configuration file.                | `m3tal config`  |
| `/var/lib/m3tal/state.db` | SQLite database storing M3TAL's operational state.         | M3TAL API       |
| `/opt/m3tal/stack/`       | Canonical directory containing M3TAL's Docker Compose files. | M3TAL Core      |
| `/docker`                 | Symlink to `/opt/m3tal/stack/`. User-facing path for stacks. | M3TAL Core      |
| `/docker/users.json`      | Stores dashboard user credentials (username and hashed password). | `m3tal dashpass`|

## Port Map

| Port | Service          | Access Method                                          |
|------|------------------|--------------------------------------------------------|
| 80   | Traefik (HTTP)   | Public (when `DASHBOARD_EXPOSE_MODE=traefik`)          |
| 8080 | M3TAL API Daemon | Host-local (internal communication for dashboard)      |
| 8081 | Traefik Dashboard| Host-local only                                        |
| 8082 | M3TAL Dashboard  | Direct port (`local` mode) or via Traefik (`traefik` mode) |

## Firewall Note

If Traefik is exposed to the internet and you are using `ufw` (Uncomplicated Firewall), ensure port 80 is allowed:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`:

- **Check the service status**:
  ```bash
  systemctl status m3tal-api
  ```

- **View logs in real-time**:
  ```bash
  journalctl -u m3tal-api -f
  ```