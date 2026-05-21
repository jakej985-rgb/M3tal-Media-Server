# M3TAL - Getting Started Guide

This guide will walk you through the initial setup of M3TAL.

## 1. Prerequisites

Before you begin, ensure you have the following software installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify their installation with the following commands:

```bash
docker --version
docker compose version
```

If Docker or Docker Compose is not installed, please refer to the official Docker documentation for installation instructions.

## 2. Install M3TAL via APT

Install M3TAL using our official APT repository with the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up essential parameters. Run the command:

```bash
sudo m3tal config wizard
```

You will be prompted with several questions:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via its port (default). This is recommended for initial setup and LAN access.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy using a domain name. This requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state and configuration files. By default, it's set to `./state` relative to the API's working directory, but it will be managed appropriately by the systemd service.
*   **`LOG_LEVEL`**: Sets the verbosity of M3TAL logs. Options typically include `debug`, `info`, `warn`, `error`. `info` is the default.
*   **`DASHBOARD_SECRET`**: A secret key used for securing the dashboard session. It's crucial to change this from the default.
*   **`API_TOKEN`**: A token used for authenticating API requests. Change this from the default.
*   **`ADMIN_PASSWORD`**: The password for the dashboard administrator account. Change this from the default.
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL will use. The default is `m3tal`.
*   **`LOCAL_IP`**: Your server's local IP address. The default is `127.0.0.1`.
*   **`DOMAIN`**: The domain name you will use to access M3TAL services. The default is `localhost`.
*   **`VPN_USER`**: Username for VPN connection, if configured.
*   **`VPN_PASSWORD`**: Password for VPN connection, if configured.
*   **`BASE_STORAGE_PATH`**: The base directory for storing data. Defaults to `./data`.
*   **`MEDIA_PATH`**: Path for media files. Defaults to `./data/media`.
*   **`CONFIG_PATH`**: Path for configuration files. Defaults to `./data/config`.
*   **`DOWNLOADS_PATH`**: Path for downloads. Defaults to `./data/downloads`.
*   **`PUID`**: The User ID for running containers. Defaults to `1000`.
*   **`PGID`**: The Group ID for running containers. Defaults to `1000`.
*   **`TZ`**: Your timezone. Defaults to `America/Denver`.
*   **`TRAEFIK_WEB_PORT`**: The port Traefik listens on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik listens on for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's own dashboard listens on. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enables or disables debug mode for M3TAL. Defaults to `false`.
*   **`METRICS_ENABLED`**: Enables or disables metrics collection. Defaults to `true`.

## 4. Start the Routing Stack (Traefik)

Traefik acts as the primary reverse proxy for your M3TAL services. To start Traefik and other essential routing components, run:

```bash
m3tal up
```

This command orchestrates the Docker Compose files located within the `/docker/` directory, starting all necessary services, including Traefik.

## 5. Start the Dashboard

The M3TAL dashboard provides a user interface for managing your services. To start the dashboard container:

```bash
m3tal dash up
```

This command will pull the necessary `m3tal-dashboard` Docker image if it's not already present and start the dashboard container, respecting your `DASHBOARD_EXPOSE_MODE` setting from the configuration wizard.

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local`, access the dashboard at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or use `http://localhost:8082` if accessing from the same machine).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik` and Traefik is running, access the dashboard at: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured).

## 7. Log In

When you access the dashboard for the first time, you will be prompted to log in.

*   **Default Credentials**:
    *   Username: `admin`
    *   Password: The password you set during the `m3tal config wizard` for `ADMIN_PASSWORD`.

**Important**: If you forgot or need to change your dashboard password after the initial setup, you can do so using the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the dashboard administrator.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure for configuration and state management.

| Path                     | Purpose                                                                      | Notes                                                                                                                               |
| :----------------------- | :--------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file for M3TAL.                                        | Managed by `m3tal config wizard` and `m3tal config set`. Contains all environment variables.                                      |
| `/var/lib/m3tal/state.db`| SQLite state database.                                                       | Stores operational state, service configurations, and other persistent data for the M3TAL API. Auto-created by the API daemon.        |
| `/opt/m3tal/stack/`      | Canonical location for M3TAL's core Docker Compose files and Traefik config. | Contains `routing-compose.yml`, `m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`, and Traefik config files. |
| `/docker`                | Symlink to `/opt/m3tal/stack/`.                                              | User-facing path for all M3TAL stack operations. Place custom compose files here to have `m3tal up` manage them.                     |
| `/docker/users.json`     | Dashboard credential store.                                                  | Stores encrypted dashboard user credentials. Managed by `m3tal dashpass`.                                                           |

---

## Port Table

The following ports are used by M3TAL and its components:

| Port   | Service                     | Access                                                  | Notes                                                                                                   |
| :----- | :-------------------------- | :------------------------------------------------------ | :------------------------------------------------------------------------------------------------------ |
| 80     | Traefik (HTTP Entrypoint)   | Public (when `DASHBOARD_EXPOSE_MODE=traefik` is used) | If Traefik is exposed to the internet, ensure your firewall allows traffic on this port.                |
| 8080   | M3TAL API Daemon (Go)       | Host-local only                                         | Accessible internally by Docker containers.                                                             |
| 8081   | Traefik Dashboard           | Host-local only (`127.0.0.1:8081`)                      | Traefik's own administrative interface.                                                                 |
| 8082   | M3TAL Dashboard             | Direct port (when `local` mode) or via Traefik (`traefik` mode) | The primary user interface. Accessible as `http://YOUR_IP:8082` or `http://dash.DOMAIN`.               |

---

## Firewall Note

If you intend to expose Traefik to the internet (e.g., for domain-based access), you must allow incoming traffic on port 80. If you are using `ufw` (Uncomplicated Firewall):

```bash
sudo ufw allow 80
```

---

## Service Management (API Daemon)

The M3TAL API daemon (`m3tal-api.service`) is managed by `systemd`. You can monitor its status and view logs using the following commands:

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