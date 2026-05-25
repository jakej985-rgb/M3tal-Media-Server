# M3TAL Get Started Guide

This guide will walk you through the initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before proceeding, ensure you have the following installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify your installation by running:

```bash
docker --version && docker compose version
```

If Docker or Docker Compose are not installed, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is installed using the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The M3TAL configuration wizard guides you through setting up essential parameters for your M3TAL installation. Run the wizard with the following command:

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a breakdown of each prompt:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. (Default: `8082`)
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via the `DASHBOARD_PORT`. Suitable for LAN-only access and initial setup.
    *   `traefik`: Exposes the dashboard through Traefik, typically via a subdomain. Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. (Default: `8080`)
*   **`STATE_DIR`**: The directory where M3TAL stores its state. (Default: `./state`)
*   **`LOG_LEVEL`**: The verbosity of M3TAL's logs. Common values are `info`, `debug`, `warn`, `error`. (Default: `info`)
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for secure communication. **Change this from the default.**
*   **`API_TOKEN`**: A token for authenticating API requests. **Change this from the default.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. **Change this from the default.**
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL will use. (Default: `m3tal`)
*   **`LOCAL_IP`**: Your local network IP address. (Default: `127.0.0.1`)
*   **`DOMAIN`**: The domain name you intend to use for M3TAL services if using Traefik. (Default: `localhost`)
*   **`VPN_USER`**: Username for VPN connection (if applicable). (Default: `user`)
*   **`VPN_PASSWORD`**: Password for VPN connection (if applicable). (Default: `password`)
*   **`BASE_STORAGE_PATH`**: The base directory for storing M3TAL data. (Default: `./data`)
*   **`MEDIA_PATH`**: Path for media storage within `BASE_STORAGE_PATH`. (Default: `./data/media`)
*   **`CONFIG_PATH`**: Path for configuration files within `BASE_STORAGE_PATH`. (Default: `./data/config`)
*   **`DOWNLOADS_PATH`**: Path for download storage within `BASE_STORAGE_PATH`. (Default: `./data/downloads`)
*   **`PUID`**: The user ID for Docker container processes. (Default: `1000`)
*   **`PGID`**: The group ID for Docker container processes. (Default: `1000`)
*   **`TZ`**: Your local timezone. (Default: `America/Denver`)
*   **`TRAEFIK_WEB_PORT`**: The port Traefik listens on for HTTP traffic. (Default: `80`)
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik listens on for HTTPS traffic. (Default: `443`)
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's dashboard is accessible on (usually for debugging). (Default: `8080`)
*   **`DEBUG_MODE`**: Enable or disable debug logging. (Default: `false`)
*   **`METRICS_ENABLED`**: Enable or disable metrics collection. (Default: `true`)

The wizard will save your choices to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

M3TAL uses Docker Compose to manage its services. The `m3tal up` command will start all Docker Compose files found in the `/docker/` directory. This includes the Traefik reverse proxy, which is essential for routing traffic to your services.

```bash
m3tal up
```

This command will pull the necessary Docker images and start the containers defined in the compose files within `/docker/`.

## Step 5: Start the Dashboard

The M3TAL dashboard provides a user interface for managing your M3TAL environment.

```bash
m3tal dash up
```

This command will pull the `m3tal-dashboard` Docker image if it's not already present and start the dashboard container, respecting the `DASHBOARD_EXPOSE_MODE` setting from your configuration.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local`: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or `http://localhost:8082` if running on your local machine).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured, e.g., `http://dash.localhost`).

## Step 7: Log In

You will be presented with a login screen.

*   **Username:** `admin`
*   **Password:** The password you set during the `m3tal config wizard` for `ADMIN_PASSWORD`.

If you need to change your dashboard password after initial setup, you can use the following command:

```bash
sudo m3tal dashpass
```

This will prompt you to enter a new password and update the credentials.

## Filesystem Contract

M3TAL organizes its configuration and data across the following key locations:

| Path                         | Purpose                                                         |
| :--------------------------- | :-------------------------------------------------------------- |
| `/etc/m3tal/.env`            | Primary M3TAL configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`    | SQLite database storing M3TAL's operational state. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`          | Canonical directory containing M3TAL's Docker Compose files and Traefik configuration. |
| `/docker`                    | Symlink pointing to `/opt/m3tal/stack/`. This is the user-facing path for managing Docker stacks. |
| `/docker/users.json`         | Stores dashboard user credentials. Managed by `m3tal dashpass`. |

## Port Map

The following ports are used by M3TAL and its components:

| Port   | Service           | Access                                 |
| :----- | :---------------- | :------------------------------------- |
| 80     | Traefik (HTTP)    | Public (when Traefik mode is enabled)  |
| 8080   | M3TAL API daemon  | Host-local                             |
| 8081   | Traefik Dashboard | Host-local only                        |
| 8082   | M3TAL Dashboard   | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If Traefik is exposed to the public internet and you are using a firewall such as `ufw`, ensure that port 80 (HTTP) is allowed:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```