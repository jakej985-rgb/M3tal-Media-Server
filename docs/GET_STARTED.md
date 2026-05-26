# M3TAL - Getting Started Guide

This guide will walk you through the initial setup of the M3TAL system.

## 1. Prerequisites

Before proceeding, ensure you have the following software installed on your system:

*   **Docker Engine**: The containerization platform.
*   **Docker Compose V2**: The tool for defining and running multi-container Docker applications.

You can verify your installation by running the following command:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Install the M3TAL CLI binary using the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

Initialize M3TAL and configure essential settings by running the configuration wizard:

```bash
sudo m3tal config wizard
```

This command will prompt you for various configuration values. Here's a brief explanation of each prompt:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. Defaults to `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed. Options are `local` (default, direct port access) or `traefik` (access via Traefik reverse proxy).
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. Defaults to `8080`.
*   **`STATE_DIR`**: The directory to store M3TAL's state data. Defaults to `./state`.
*   **`LOG_LEVEL`**: The verbosity of M3TAL logs. Options include `debug`, `info`, `warn`, `error`.
*   **`DASHBOARD_SECRET`**: A secret key for the dashboard. **Change this from the default immediately.**
*   **`API_TOKEN`**: A token for API access. **Change this from the default immediately.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. **Change this from the default immediately.**
*   **`NETWORK_NAME`**: The Docker network name M3TAL services will use. Defaults to `m3tal`.
*   **`LOCAL_IP`**: Your local IP address. Defaults to `127.0.0.1`.
*   **`DOMAIN`**: The domain name to use for accessing services via Traefik. Defaults to `localhost`.
*   **`BASE_STORAGE_PATH`**: The base path for storing data for M3TAL managed services. Defaults to `./data`.
*   **`MEDIA_PATH`**: Path for media storage within `BASE_STORAGE_PATH`. Defaults to `./data/media`.
*   **`CONFIG_PATH`**: Path for configuration files within `BASE_STORAGE_PATH`. Defaults to `./data/config`.
*   **`DOWNLOADS_PATH`**: Path for downloads within `BASE_STORAGE_PATH`. Defaults to `./data/downloads`.
*   **`PUID`**: The user ID for Docker container processes. Defaults to `1000`.
*   **`PGID`**: The group ID for Docker container processes. Defaults to `1000`.
*   **`TZ`**: Your timezone. Defaults to `America/Denver`.
*   **`TRAEFIK_WEB_PORT`**: The port Traefik listens on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik listens on for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port for the Traefik dashboard. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enable or disable debug mode. Defaults to `false`.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection. Defaults to `true`.

All configuration options are stored in `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

M3TAL uses Docker Compose to manage its services. The `m3tal up` command starts all Docker Compose files located in the `/docker/` directory. This typically includes Traefik for routing traffic to your services.

```bash
m3tal up
```

This command will pull necessary Docker images and start the routing stack, making your services accessible via domain names if configured.

## 5. Start the Dashboard

To start the M3TAL dashboard container:

```bash
m3tal dash up
```

This command will download the latest `m3tal-compose.yml` and its relevant override file (based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`) and start the dashboard container.

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard's address. The access method depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):** Access the dashboard at `http://YOUR_IP:8082`. Replace `YOUR_IP` with your server's IP address. You can also access it via `http://localhost:8082` if you are on the same machine.
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:** Access the dashboard at `http://dash.DOMAIN`. Replace `DOMAIN` with the domain you configured in the wizard. Traefik must be running for this to work.

## 7. Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: The `ADMIN_PASSWORD` you set during the configuration wizard.

If you need to change your dashboard password after initial setup, use the following command:

```bash
sudo m3tal dashpass
```
This command will prompt you for the new password and update the credentials.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure to manage its configuration and state.

| Path                        | Purpose                                                                                               |
| :-------------------------- | :---------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | The primary configuration file. This file stores all environment variables and is managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`   | The SQLite state database used by the M3TAL API daemon to store system information. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`         | The canonical directory containing M3TAL's core Docker Compose files and Traefik configuration.       |
| `/docker`                   | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for interacting with M3TAL's Docker stacks. |
| `/docker/users.json`        | Stores the M3TAL dashboard credentials. Managed by the `m3tal dashpass` command.                      |

---

## Port Map

M3TAL and its components utilize the following ports:

| Port | Service          | Access                                            |
| :--- | :--------------- | :------------------------------------------------ |
| 80   | Traefik          | Public (when `DASHBOARD_EXPOSE_MODE` is `traefik`) |
| 8080 | M3TAL API daemon | Host-local access                                 |
| 8081 | Traefik dashboard | Host-local access only                           |
| 8082 | M3TAL Dashboard  | Direct port access (when `local` mode) or via Traefik (when `traefik` mode) |

---

## Firewall Note

If you are exposing Traefik to the public internet (e.g., `DASHBOARD_EXPOSE_MODE=traefik`), ensure that port 80 is allowed through your firewall:

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View service logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```