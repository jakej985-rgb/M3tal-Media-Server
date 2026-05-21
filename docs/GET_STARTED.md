# GET STARTED WITH M3TAL

This guide provides a step-by-step process for installing and setting up M3TAL for the first time.

## 1. Prerequisites

Before proceeding, ensure the following software is installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

Verify your installation by running:

```bash
docker --version && docker compose version
```

If these commands do not return version information, please refer to the official Docker documentation for installation instructions.

## 2. Install M3TAL via APT

M3TAL is distributed as a Debian package. Use the following commands to add the M3TAL repository and install the CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The `m3tal config wizard` command will guide you through the initial configuration of M3TAL.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. Defaults to `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via its port (e.g., `http://YOUR_IP:8082`). Recommended for initial setup.
    *   `traefik`: Exposes the dashboard through Traefik, typically via a subdomain (e.g., `http://dash.DOMAIN`). Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. Defaults to `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state data (like the `state.db` file). Defaults to `./state` within the M3TAL data directory.
*   **`LOG_LEVEL`**: Sets the verbosity of M3TAL logs. Common options include `info`, `debug`, `warn`, `error`.
*   **`DASHBOARD_SECRET`**: A secret key used for securing the dashboard session. **Change this from the default value immediately.**
*   **`API_TOKEN`**: A token for authenticating API requests. **Change this from the default value immediately.**
*   **`ADMIN_PASSWORD`**: The default password for the dashboard administrator. **Change this from the default value immediately.**
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL services will use. Defaults to `m3tal`.
*   **`LOCAL_IP`**: The IP address of your host machine.
*   **`DOMAIN`**: The domain name used for accessing services via Traefik. Defaults to `localhost`.
*   **`VPN_USER`**, **`VPN_PASSWORD`**: Credentials for a VPN service if you intend to use it with M3TAL.
*   **`BASE_STORAGE_PATH`**: The root directory for M3TAL's persistent data. Defaults to `./data`.
*   **`MEDIA_PATH`**: The subdirectory within `BASE_STORAGE_PATH` for media files.
*   **`CONFIG_PATH`**: The subdirectory within `BASE_STORAGE_PATH` for configuration files.
*   **`DOWNLOADS_PATH`**: The subdirectory within `BASE_STORAGE_PATH` for downloaded files.
*   **`PUID`**, **`PGID`**: User and group IDs for running Docker containers. Usually your host user's IDs.
*   **`TZ`**: Your timezone (e.g., `America/Denver`).
*   **`TRAEFIK_WEB_PORT`**: The port Traefik listens on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik listens on for HTTPS traffic. Defaults to `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's own dashboard is accessible on. Defaults to `8080`.
*   **`DEBUG_MODE`**: Enable or disable debug mode.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection.

These settings will be saved in `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

This command starts the core M3TAL services, including Traefik, which acts as the primary reverse proxy for your applications.

```bash
m3tal up
```

This command orchestrates the `docker compose` command across all `.yml` files found in the `/docker/` directory. This typically includes the Traefik routing stack and the M3TAL API daemon.

## 5. Start the Dashboard

The M3TAL dashboard provides a web interface to manage your M3TAL system.

```bash
m3tal dash up
```

This command will:

1.  Download the necessary Docker image for the M3TAL dashboard (`ghcr.io/jakej985-rgb/m3tal-godash:debug`).
2.  Read your `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Start the dashboard container, applying the appropriate configuration (either direct port binding for `local` mode or Traefik labels for `traefik` mode).

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default), access the dashboard at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or use `http://localhost:8082` if on the same machine).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, access the dashboard at: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured, e.g., `http://dash.localhost` if using `localhost` as your domain).

## 7. Log In

Upon first access, you will be presented with the M3TAL login screen.

*   **Username:** `admin`
*   **Password:** The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard`.

If you need to change the administrator password after the initial setup, use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password and update the dashboard's credential store.

## Filesystem Contract

M3TAL organizes its configuration and data across specific directories and files. Understanding this contract is crucial for managing your installation.

| Path                       | Purpose                                                                                                         |
| :------------------------- | :-------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | The primary M3TAL configuration file. All environment variables for M3TAL services are defined here. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`  | SQLite database storing M3TAL's internal state, such as registered services, configurations, and other operational data. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`        | The canonical directory for M3TAL's core Docker Compose files and Traefik configuration.                        |
| `/docker`                  | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all M3TAL stack operations (e.g., `m3tal up`, `m3tal down`). New user-defined services should have their compose files placed here. |
| `/docker/users.json`       | Stores the encrypted credentials for dashboard users. Managed by `m3tal dashpass`.                                  |

## Port Table

The following ports are used by M3TAL and its components:

| Port | Service           | Access                                                    |
| :--- | :---------------- | :-------------------------------------------------------- |
| 80   | Traefik HTTP      | Public (when `DASHBOARD_EXPOSE_MODE=traefik`)             |
| 8080 | M3TAL API daemon  | Host-local (accessible by containers and the host)        |
| 8081 | Traefik dashboard | Host-local only (internal Traefik management interface)   |
| 8082 | M3TAL Dashboard   | Direct port access (`local` mode) or via Traefik (`traefik` mode) |

## Firewall Note

If you are exposing M3TAL services to the public internet and using Traefik on port 80, ensure your firewall allows incoming traffic on this port. For example, if using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon is managed by `systemd`. You can use `systemctl` commands to control and monitor its status.

*   **Check the status of the M3TAL API service:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View real-time logs for the M3TAL API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```