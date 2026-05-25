# M3TAL Ecosystem: Getting Started Guide

This guide provides a step-by-step process for installing and setting up the M3TAL ecosystem for the first time.

## Step 1: Prerequisites

Before proceeding, ensure you have the following software installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify their installation by running the following commands:

```bash
docker --version
docker compose version
```

## Step 2: Install M3TAL via APT

Install the M3TAL CLI binary and its associated systemd service using the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

Initialize M3TAL's configuration by running the wizard. This will guide you through setting essential environment variables for your setup.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for the following settings. Press Enter to accept the default value shown in parentheses, or type your desired value and press Enter.

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. Defaults to `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a host port (default, recommended for beginners).
    *   `traefik`: Exposes the dashboard via the Traefik reverse proxy (requires Traefik to be running).
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. Defaults to `8080`.
*   **`STATE_DIR`**: The directory for storing M3TAL's state data. Defaults to `./state` within the base storage path.
*   **`LOG_LEVEL`**: The verbosity of M3TAL's logs (`info`, `debug`, `warn`, `error`). Defaults to `info`.
*   **`DASHBOARD_SECRET`**: A secret key for the dashboard. **Change this from the default value.**
*   **`API_TOKEN`**: A token for API authentication. **Change this from the default value.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. **Change this from the default value.**
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL will use. Defaults to `m3tal`.
*   **`LOCAL_IP`**: Your local network IP address. Defaults to `127.0.0.1`.
*   **`DOMAIN`**: The domain name to use for services exposed via Traefik. Defaults to `localhost`.
*   **`VPN_USER`**: Username for VPN connection (if applicable).
*   **`VPN_PASSWORD`**: Password for VPN connection (if applicable).
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL's data and configuration files. Defaults to `./data`.
*   **`MEDIA_PATH`**: Path for media storage, relative to `BASE_STORAGE_PATH`.
*   **`CONFIG_PATH`**: Path for configuration storage, relative to `BASE_STORAGE_PATH`.
*   **`DOWNLOADS_PATH`**: Path for download storage, relative to `BASE_STORAGE_PATH`.
*   **`PUID`**: The user ID for running Docker containers. Defaults to `1000`.
*   **`PGID`**: The group ID for running Docker containers. Defaults to `1000`.
*   **`TZ`**: Your timezone (e.g., `America/Denver`).
*   **`TRAEFIK_WEB_PORT`**: The port Traefik will listen on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik will listen on for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's internal dashboard listens on. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enable debug mode for M3TAL. Defaults to `false`.
*   **`METRICS_ENABLED`**: Enable metrics collection. Defaults to `true`.

## Step 4: Start the Routing Stack (Traefik)

This command initiates the core M3TAL services, including the Traefik reverse proxy. It reads and applies all `*-compose.yml` files found in the `/docker/` directory.

```bash
m3tal up
```

This command orchestrates the startup of services defined in your Docker Compose files, essential for routing traffic to your applications.

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL dashboard container. It will pull the latest dashboard image and start the container based on your `DASHBOARD_EXPOSE_MODE` configuration.

```bash
m3tal dash up
```

This command ensures the M3TAL dashboard is running and accessible.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default), access it at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address or `localhost`).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, access it at: `http://dash.DOMAIN` (replace `DOMAIN` with your configured domain, e.g., `http://dash.localhost`).

## Step 7: Log In

Use the following default credentials to log in to the M3TAL dashboard:

*   **Username**: `admin`
*   **Password**: The password you set during the `m3tal config wizard` (or `admin_pass` if not changed).

To change your dashboard password after the initial login, use the following command, replacing `new_password_here` with your desired password:

```bash
sudo m3tal dashpass new_password_here
```

---

## Filesystem Contract

M3TAL adheres to a specific filesystem structure for configuration and data.

| Path                       | Purpose                                                                                              |
| :------------------------- | :--------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | The primary M3TAL environment configuration file. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`  | SQLite database storing M3TAL's operational state. Automatically created by the API daemon.          |
| `/opt/m3tal/stack/`        | The canonical directory containing M3TAL's core stack files, including compose configurations and Traefik setup. |
| `/docker`                  | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory for managing stacks and compose files. |
| `/docker/users.json`       | Stores dashboard user credentials. Managed via the `m3tal dashpass` command.                           |

---

## Port Table

| Port   | Service            | Access                                                               |
| :----- | :----------------- | :------------------------------------------------------------------- |
| 80     | Traefik (HTTP)     | Publicly accessible if Traefik is configured for external access.    |
| 8080   | M3TAL API (Go)     | Accessible from within the host machine, used by other M3TAL components. |
| 8081   | Traefik (Dashboard)| Accessible locally on the host for Traefik management.               |
| 8082   | M3TAL Dashboard    | Accessible directly (local mode) or via Traefik (traefik mode).      |

---

## Firewall Note

If you intend to expose Traefik to the public internet on port 80, you will need to allow this port through your firewall. If you are using `ufw`:

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

**Check the status of the M3TAL API service:**

```bash
systemctl status m3tal-api
```

**View live logs for the M3TAL API service:**

```bash
journalctl -u m3tal-api -f
```