# Getting Started with M3TAL

This guide will walk you through the initial setup of the M3TAL ecosystem.

## 1. Prerequisites

Before you begin, ensure you have the following software installed on your system:

-   **Docker Engine**
-   **Docker Compose V2**

You can verify their installation and versions by running:

```bash
docker --version && docker compose version
```

If these are not installed, please refer to the official Docker documentation for installation instructions.

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

After installation, run the configuration wizard to set up your M3TAL environment:

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's what each prompt means:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a port binding (e.g., `http://YOUR_IP:8082`). This is the default and recommended for initial setup on a local network.
    *   `traefik`: Exposes the dashboard via the Traefik reverse proxy using a domain name (e.g., `http://dash.YOUR_DOMAIN`). Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state, including the `state.db` file. Defaults to `./state` relative to the script, but will be managed by the API service.
*   **`LOG_LEVEL`**: Sets the verbosity of logs (e.g., `info`, `debug`).
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard. It's highly recommended to change this from the default.
*   **`API_TOKEN`**: A token for API authentication. It's highly recommended to change this from the default.
*   **`ADMIN_PASSWORD`**: The password for the dashboard administrator. It's highly recommended to change this from the default.
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL services will use. Defaults to `m3tal`.
*   **`LOCAL_IP`**: Your local network IP address.
*   **`DOMAIN`**: The domain name you will use to access services. Defaults to `localhost`.
*   **`VPN_USER`** and **`VPN_PASSWORD`**: Credentials for a VPN connection if required by your setup.
*   **`BASE_STORAGE_PATH`**: The base directory for storing persistent data for M3TAL services. Defaults to `./data`.
*   **`MEDIA_PATH`**: Path for media storage.
*   **`CONFIG_PATH`**: Path for configuration files.
*   **`DOWNLOADS_PATH`**: Path for downloads.
*   **`PUID`** and **`PGID`**: The User ID and Group ID that Docker containers will run as. Often `1000:1000` for the primary user.
*   **`TZ`**: Your timezone (e.g., `America/Denver`).
*   **`TRAEFIK_WEB_PORT`**: The port Traefik listens on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik listens on for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's own dashboard is accessible on. Defaults to `8080` (often proxied to `127.0.0.1:8080`).
*   **`DEBUG_MODE`**: Enable or disable debug mode.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection.

The wizard will save your choices in `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

M3TAL uses Traefik as its reverse proxy to manage access to various services. To start Traefik and other core routing components, run:

```bash
m3tal up
```

This command reads all Docker Compose files located in the `/docker/` directory and starts the defined services. This includes Traefik, which will be configured to route traffic.

## 5. Start the Dashboard

Next, start the M3TAL dashboard container:

```bash
m3tal dash up
```

This command will:
1.  Download the necessary `m3tal-compose.yml` files and any required override files (e.g., `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.
2.  Pull the `m3tal-dashboard` Docker image if it's not already present.
3.  Start the `m3tal-dashboard` container.

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard's address. This will be:

*   **`http://YOUR_IP:8082`**: If `DASHBOARD_EXPOSE_MODE` is set to `local`. Replace `YOUR_IP` with your server's IP address.
*   **`http://dash.YOUR_DOMAIN`**: If `DASHBOARD_EXPOSE_MODE` is set to `traefik`. Replace `YOUR_DOMAIN` with the domain you configured.

## 7. Log In

You will be presented with the M3TAL login screen.

*   **Default Credentials**:
    *   Username: `admin`
    *   Password: The password you set for `ADMIN_PASSWORD` during the configuration wizard (default: `admin_pass`).

**Changing Your Password**:

If you need to change your dashboard password after the initial setup, use the following command:

```bash
sudo m3tal dashpass
```

This will prompt you for the new password and update the `users.json` file.

## Filesystem Contract

M3TAL relies on a specific file and directory structure for its operation:

| Path                      | Purpose                                                                                                     |
| :------------------------ | :---------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | The primary configuration file for M3TAL. It is managed by the `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db` | The SQLite state database used by the M3TAL API daemon to store persistent information.                      |
| `/opt/m3tal/stack/`       | The canonical directory containing all M3TAL Docker Compose files, Traefik configuration, and related assets. |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory for managing Docker stacks. |
| `/docker/users.json`      | Stores the dashboard credentials. Managed by the `m3tal dashpass` command.                                    |

## Port Map

The following ports are used by M3TAL services:

| Port | Service         | Access Method                               |
| :--- | :-------------- | :------------------------------------------ |
| 80   | Traefik         | Public (HTTP entry point, when exposed)     |
| 8080 | M3TAL API daemon| Host-local only                             |
| 8081 | Traefik Dashboard| Host-local only                             |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If you have a firewall enabled (e.g., `ufw`) and want Traefik to be accessible from outside your local network, you must allow traffic on port 80:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it and view its logs using the following commands:

*   **Check the status of the API service:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live logs of the API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```