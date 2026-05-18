# M3TAL Ecosystem: Getting Started Guide

This guide provides the necessary steps for a new user to install and set up the M3TAL ecosystem.

## Step 1: Prerequisites

Before proceeding, ensure you have the following software installed on your system:

*   **Docker Engine**: The core containerization platform.
*   **Docker Compose V2**: The tool for defining and running multi-container Docker applications.

Verify your installation by running the following command in your terminal:

```bash
docker --version && docker compose version
```

If either command returns an error or shows an outdated version, please refer to the official Docker documentation for installation and upgrade instructions.

## Step 2: Install M3TAL via APT

M3TAL is distributed via a dedicated APT repository for straightforward installation and updates. Execute the following three commands sequentially in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This process will install the M3TAL CLI binary (`/usr/bin/m3tal`) and set up necessary systemd services.

## Step 3: Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up essential environment variables. Run the following command:

```bash
sudo m3tal config wizard
```

You will be presented with a series of prompts. Here's a breakdown of each prompt and its purpose:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a port (default, recommended for initial setup).
    *   `traefik`: Exposes the dashboard via Traefik reverse proxy (requires Traefik to be running).
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL will store its state data. Defaults to `./state` relative to the M3TAL installation, but can be set to an absolute path.
*   **`LOG_LEVEL`**: Sets the verbosity of M3TAL logs. Common values include `info`, `debug`, `warn`, `error`.
*   **`DASHBOARD_SECRET`**: A secret key used for securing the dashboard session. It is highly recommended to change this from the default.
*   **`API_TOKEN`**: A token used for authenticating with the M3TAL API. It is highly recommended to change this from the default.
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard with administrator privileges. It is highly recommended to change this from the default.
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL will use. Defaults to `m3tal`.
*   **`LOCAL_IP`**: The local IP address of your host machine. Defaults to `127.0.0.1`.
*   **`DOMAIN`**: The domain name to be used with Traefik for routing. Defaults to `localhost`.
*   **`VPN_USER`**: Username for VPN integration (if used).
*   **`VPN_PASSWORD`**: Password for VPN integration (if used).
*   **`BASE_STORAGE_PATH`**: The base directory for storing M3TAL data. Defaults to `./data`.
*   **`MEDIA_PATH`**: Directory for media files. Defaults to `./data/media`.
*   **`CONFIG_PATH`**: Directory for configuration files. Defaults to `./data/config`.
*   **`DOWNLOADS_PATH`**: Directory for downloads. Defaults to `./data/downloads`.
*   **`PUID`**: The User ID for processes running within Docker containers. Defaults to `1000`.
*   **`PGID`**: The Group ID for processes running within Docker containers. Defaults to `1000`.
*   **`TZ`**: Your local timezone. Defaults to `America/Denver`.
*   **`TRAEFIK_WEB_PORT`**: The host port Traefik will listen on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The host port Traefik will listen on for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The host port Traefik's own dashboard will be accessible on. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enable or disable debug mode for M3TAL.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection.

After completing the wizard, your settings will be saved in `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

The M3TAL ecosystem relies on Traefik as its primary reverse proxy. To start the Traefik service and other essential routing components, run:

```bash
m3tal up
```

This command orchestrates the startup of all Docker Compose files located within the `/docker/` directory. This includes Traefik itself, which will be configured to listen for incoming traffic and route it to the appropriate services based on domain names or direct port mappings.

## Step 5: Start the Dashboard

The M3TAL dashboard provides a web-based interface for managing your M3TAL services. To start the dashboard container:

```bash
m3tal dash up
```

This command will:
1. Download the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` files from the M3TAL repository.
2. Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3. Pull the `m3tal-dashboard` Docker image if it's not already present.
4. Start the dashboard container, applying the appropriate compose override file based on your chosen exposure mode.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's URL. The exact address depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Navigate to `http://YOUR_IP:8082` or `http://localhost:8082`. Replace `YOUR_IP` with your server's actual IP address.
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Navigate to `http://dash.YOUR_DOMAIN` or `http://dash.localhost` (if `DOMAIN` is set to `localhost`). Ensure Traefik is running (`m3tal up`) and DNS is correctly configured if using a custom domain.

## Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: `admin_pass` (This is the default from the `ADMIN_PASSWORD` environment variable. **It is strongly recommended to change this immediately.**)

To change the dashboard password, use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you for a new password and update the `users.json` file accordingly.

---

## Filesystem Contract

The M3TAL ecosystem utilizes specific file paths for configuration and state management. Understanding these locations is crucial for troubleshooting and manual intervention.

| Path                     | Purpose                                                                    |
| :----------------------- | :------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary environment configuration file. Managed by `m3tal config wizard`.  |
| `/var/lib/m3tal/state.db`| SQLite state database used by the M3TAL API daemon to store system state.  |
| `/docker`                | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all M3TAL Docker Compose stack operations. |
| `/opt/m3tal/stack/`      | The canonical directory containing M3TAL's Docker Compose files and Traefik configuration. |
| `/docker/users.json`     | Stores dashboard user credentials (username and hashed passwords). Managed by `m3tal dashpass`. |

## Port Table

The following ports are used by M3TAL and its components:

| Port | Service        | Access Type      | Notes                                                               |
| :--- | :------------- | :--------------- | :------------------------------------------------------------------ |
| 80   | Traefik        | Public           | HTTP entry point for routed services (if exposed externally).       |
| 8080 | M3TAL API      | Host-local       | Internal API daemon, accessible by Docker containers and `host.docker.internal`. |
| 8081 | Traefik Admin  | Host-local       | Traefik's own management dashboard (if enabled).                    |
| 8082 | M3TAL Dashboard| Direct or Traefik| Accessible via direct port binding (`local` mode) or via Traefik (`traefik` mode). |

## Firewall Note

If you are exposing Traefik to the public internet (i.e., your server is accessible from outside your local network), ensure that port 80 (HTTP) is open in your firewall. For example, if using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs (follow in real-time):**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```