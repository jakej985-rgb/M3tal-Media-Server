# M3TAL Get Started Guide

This guide will walk you through the initial setup of the M3TAL system.

## Step 1: Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Docker Engine**: The containerization platform.
*   **Docker Compose V2**: The orchestration tool for Docker.

You can verify their installation by running the following commands in your terminal:

```bash
docker --version && docker compose version
```

If these commands output version information, you are ready to proceed.

## Step 2: Install M3TAL via APT

M3TAL is distributed via an APT repository for easy installation. Follow these commands to add the repository and install the M3TAL CLI:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The `m3tal config wizard` will guide you through setting up essential configuration variables.

Run the wizard by executing:

```bash
sudo m3tal config wizard
```

You will be prompted for the following:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local` (default): Exposes the dashboard directly via `HOST_IP:DASHBOARD_PORT`. Recommended for initial setup and LAN-only access.
    *   `traefik`: Exposes the dashboard via Traefik using a domain name (e.g., `dash.yourdomain.com`). Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL will store its state. The default is `./state` relative to the current directory.
*   **`LOG_LEVEL`**: The verbosity of M3TAL's logs. Options include `debug`, `info`, `warn`, `error`. `info` is the default.
*   **`DASHBOARD_SECRET`**: A secret key for the dashboard. **Change this from the default `change_me_immediately` for security.**
*   **`API_TOKEN`**: An API token for authentication. **Change this from the default `change_me_api_token` for security.**
*   **`ADMIN_PASSWORD`**: The password for the dashboard's admin user. **Change this from the default `admin_pass` for security.**
*   **`NETWORK_NAME`**: The Docker network name M3TAL will use. The default is `m3tal`.
*   **`LOCAL_IP`**: Your local network IP address. The default is `127.0.0.1`.
*   **`DOMAIN`**: The primary domain name for your M3TAL setup. `localhost` is the default.
*   **`VPN_USER`**: Username for VPN integration (if applicable).
*   **`VPN_PASSWORD`**: Password for VPN integration (if applicable).
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL's data storage. Defaults to `./data`.
*   **`MEDIA_PATH`**: Path for media files. Defaults to `./data/media`.
*   **`CONFIG_PATH`**: Path for configuration files. Defaults to `./data/config`.
*   **`DOWNLOADS_PATH`**: Path for download files. Defaults to `./data/downloads`.
*   **`PUID`**: The User ID to run Docker containers as. Defaults to `1000`.
*   **`PGID`**: The Group ID to run Docker containers as. Defaults to `1000`.
*   **`TZ`**: Your timezone. Defaults to `America/Denver`.
*   **`TRAEFIK_WEB_PORT`**: The host port Traefik listens on for HTTP. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The host port Traefik listens on for HTTPS. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The host port Traefik's own dashboard listens on. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enable or disable debug mode. Defaults to `false`.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection. Defaults to `true`.

The wizard will save these settings to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

The routing stack, managed by Traefik, acts as the entry point for your M3TAL services. If you chose `traefik` mode for dashboard exposure in the wizard, this step is crucial.

To start the routing stack and any other defined services, run:

```bash
m3tal up
```

This command orchestrates all `*-compose.yml` files found in the `/docker/` directory, starting services like Traefik and the API daemon.

## Step 5: Start the Dashboard

Next, you need to start the M3TAL dashboard container.

Execute the following command:

```bash
m3tal dash up
```

This command will pull the latest M3TAL dashboard Docker image and start the container. It respects the `DASHBOARD_EXPOSE_MODE` setting from your `.env` file.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` in `/etc/m3tal/.env`, access it at:
    `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address or `localhost` if accessing locally).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, and Traefik is running, access it at:
    `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured in `/etc/m3tal/.env`).

## Step 7: Log In

You will be presented with the M3TAL dashboard login screen.

*   **Username**: `admin`
*   **Password**: The `ADMIN_PASSWORD` you set during the configuration wizard.

**To change the admin password after initial setup, use the following command:**

```bash
sudo m3tal dashpass admin NEW_PASSWORD
```

Replace `NEW_PASSWORD` with your desired secure password.

---

## Filesystem Contract

M3TAL utilizes specific locations in the filesystem for its configuration and data.

| Path                 | Purpose                                                                | Managed By             |
| :------------------- | :--------------------------------------------------------------------- | :--------------------- |
| `/etc/m3tal/.env`    | Primary M3TAL configuration file.                                      | `m3tal config wizard`  |
| `/var/lib/m3tal/`    | Directory for M3TAL's internal state and databases.                    | M3TAL API daemon       |
| `/var/lib/m3tal/state.db` | SQLite state database for the M3TAL API.                             | M3TAL API daemon       |
| `/opt/m3tal/stack/`  | Canonical directory containing M3TAL's Docker Compose files and configs. | M3TAL installation     |
| `/docker`            | Symlink to `/opt/m3tal/stack/`. User-facing path for stack operations. | M3TAL installation     |
| `/docker/users.json` | Stores dashboard credentials.                                          | `m3tal dashpass`       |

## Port Map

| Port | Service           | Access                                       |
| :--- | :---------------- | :------------------------------------------- |
| 80   | Traefik (HTTP)    | Public (if exposed and `traefik` mode used)  |
| 8080 | M3TAL API daemon  | Host-local                                   |
| 8081 | Traefik Dashboard | Host-local only                              |
| 8082 | M3TAL Dashboard   | Direct port (if `local` mode) or via Traefik |

## Firewall Note

If Traefik is exposed to the internet (e.g., if your server is publicly accessible and you intend to use domain-based access), ensure port 80 is open in your firewall. For example, using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it as follows:

*   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs in real-time**:
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```