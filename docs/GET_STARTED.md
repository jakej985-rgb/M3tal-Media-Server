# M3TAL - Getting Started Guide

This guide provides step-by-step instructions for setting up M3TAL for the first time.

## Step 1: Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

Verify your installation by running the following commands:

```bash
docker --version
docker compose version
```

If either of these commands fails, please install Docker and Docker Compose according to their official documentation before proceeding.

## Step 2: Install M3TAL via APT

Install the M3TAL CLI binary by adding the official repository and then installing the package.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

Initialize M3TAL's configuration by running the wizard. This process will guide you through setting essential parameters.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's what each prompt means:

*   **`DASHBOARD_PORT`**: The local port where the M3TAL dashboard will be accessible when `DASHBOARD_EXPOSE_MODE` is set to `local`. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is made accessible.
    *   `local`: Exposes the dashboard directly via the `DASHBOARD_PORT`. Recommended for initial setup and LAN-only access.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy, typically via a subdomain (e.g., `dash.yourdomain.com`). Requires Traefik to be running.
*   **`HTTP_PORT`**: The port the M3TAL API daemon listens on. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state, including the database. The default is `./state` relative to the M3TAL data directory.
*   **`LOG_LEVEL`**: Sets the verbosity of M3TAL's logging. Options typically include `debug`, `info`, `warn`, `error`.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management. **Change this from the default immediately.**
*   **`API_TOKEN`**: A token for authenticating API requests. **Change this from the default immediately.**
*   **`ADMIN_PASSWORD`**: The password for the default administrator account in the dashboard. **Change this from the default immediately.**
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL services will use. The default is `m3tal`.
*   **`LOCAL_IP`**: Your local network IP address.
*   **`DOMAIN`**: The domain name you intend to use for M3TAL services. Defaults to `localhost`.
*   **`VPN_USER`**: Username for VPN access (if configured).
*   **`VPN_PASSWORD`**: Password for VPN access (if configured).
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL data storage. Defaults to `./data`.
*   **`MEDIA_PATH`**: Directory for media files. Defaults to `./data/media`.
*   **`CONFIG_PATH`**: Directory for configuration files. Defaults to `./data/config`.
*   **`DOWNLOADS_PATH`**: Directory for downloads. Defaults to `./data/downloads`.
*   **`PUID`**: The user ID to run Docker containers under. Typically your user ID.
*   **`PGID`**: The group ID to run Docker containers under. Typically your user's group ID.
*   **`TZ`**: Your timezone. e.g., `America/Denver`.
*   **`TRAEFIK_WEB_PORT`**: The host port Traefik will use for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The host port Traefik will use for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik uses internally for its own dashboard. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enable or disable debug mode for M3TAL.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection.

## Step 4: Start the Routing Stack (Traefik)

This command starts the core routing infrastructure, primarily Traefik, which acts as a reverse proxy for your M3TAL services. It reads and applies all Docker Compose files found in the `/docker/` directory.

```bash
m3tal up
```

## Step 5: Start the Dashboard

This command specifically pulls the M3TAL dashboard image and starts its container. It will use the `DASHBOARD_EXPOSE_MODE` setting from your configuration.

```bash
m3tal dash up
```

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default), access it at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's local IP address, or `localhost` if running locally).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik` and Traefik is configured, access it via your defined domain: `http://dash.DOMAIN` (replace `DOMAIN` with your configured domain).

## Step 7: Log In

Upon accessing the dashboard, you will be presented with a login screen.

*   **Username:** `admin`
*   **Password:** `admin_pass` (or whatever you set during the wizard)

**It is strongly recommended to change the default password immediately after logging in.**

To change the dashboard password using the CLI:

```bash
sudo m3tal dashpass <new_password>
```

Replace `<new_password>` with your desired strong password.

## Filesystem Contract

M3TAL utilizes a specific filesystem structure for configuration and data. Understanding these paths is crucial for advanced management and troubleshooting.

| Path                     | Purpose                                                          | Managed By                                  |
| :----------------------- | :--------------------------------------------------------------- | :------------------------------------------ |
| `/etc/m3tal/.env`        | Primary environment configuration file.                          | `m3tal config wizard` and `m3tal config set` |
| `/var/lib/m3tal/state.db` | SQLite database for M3TAL's operational state.                 | M3TAL API daemon                            |
| `/opt/m3tal/stack/`      | Canonical directory containing M3TAL's Docker Compose files.     | M3TAL installation                          |
| `/docker`                | Symlink pointing to `/opt/m3tal/stack/`. User-facing path for stack operations. | M3TAL installation                          |
| `/docker/users.json`     | Stores dashboard user credentials (username/password hashes).    | `m3tal dashpass`                            |

## Port Map

The following ports are utilized by M3TAL services:

| Port   | Service         | Access Context                                        |
| :----- | :-------------- | :---------------------------------------------------- |
| 80     | Traefik         | Public ingress for HTTP traffic (when in `traefik` mode). |
| 8080   | M3TAL API daemon | Accessible on the host system.                        |
| 8081   | Traefik         | Traefik's internal dashboard (host-local access only).|
| 8082   | M3TAL Dashboard | Direct port access (when in `local` mode).            |

## Firewall Note

If you intend for Traefik to be accessible from outside your local network (and `DASHBOARD_EXPOSE_MODE` is set to `traefik`), ensure that port 80 is allowed through your firewall:

```bash
sudo ufw allow 80
```

## Service Management (API Daemon)

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```