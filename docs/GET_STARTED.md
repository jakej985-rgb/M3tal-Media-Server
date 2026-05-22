# M3TAL Get Started Guide

This guide will walk you through the initial setup of M3TAL on your system.

## Step 1: Prerequisites

Before proceeding, ensure you have the following installed:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify your installation by running:

```bash
docker --version
docker compose version
```

## Step 2: Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The M3TAL configuration wizard will guide you through essential settings. Run the following command:

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a breakdown of what each prompt means:

*   **`DASHBOARD_PORT`**: The port on your host machine where the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is made available.
    *   `local` (default): Exposes the dashboard directly via a host port. Ideal for LAN-only access or initial setup.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy (requires `m3tal up` to be run first).
*   **`HTTP_PORT`**: The port the M3TAL API daemon listens on. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state data.
*   **`LOG_LEVEL`**: The verbosity of M3TAL's logs (e.g., `info`, `debug`, `warn`).
*   **`DASHBOARD_SECRET`**: A secret key for securing your dashboard. **Change this from the default.**
*   **`API_TOKEN`**: An API token for programmatic access to M3TAL. **Change this from the default.**
*   **`ADMIN_PASSWORD`**: The password for the dashboard's admin user. **Change this from the default.**
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL services will use.
*   **`LOCAL_IP`**: Your machine's local IP address.
*   **`DOMAIN`**: The domain name you intend to use with M3TAL services (e.g., `localhost`, `yourdomain.com`).
*   **`VPN_USER`**, **`VPN_PASSWORD`**: Credentials if you intend to use a VPN with M3TAL.
*   **`BASE_STORAGE_PATH`**: The base directory for all persistent data storage for M3TAL services.
*   **`MEDIA_PATH`**: A sub-directory within `BASE_STORAGE_PATH` for media files.
*   **`CONFIG_PATH`**: A sub-directory within `BASE_STORAGE_PATH` for configuration files.
*   **`DOWNLOADS_PATH`**: A sub-directory within `BASE_STORAGE_PATH` for downloads.
*   **`PUID`**, **`PGID`**: User and group IDs for running Docker containers. Typically your host user's IDs.
*   **`TZ`**: Your local timezone.
*   **`TRAEFIK_WEB_PORT`**: The port Traefik will listen on for HTTP traffic.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik will listen on for HTTPS traffic.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik's own dashboard will be accessible on.
*   **`DEBUG_MODE`**: Enable or disable debug mode.
*   **`METRICS_ENABLED`**: Enable or disable metrics collection.

## Step 4: Start the Routing Stack (Traefik)

This command starts the core M3TAL routing components, including Traefik, which acts as a reverse proxy for your services.

```bash
m3tal up
```

This command orchestrates all Docker Compose files found in the `/docker/` directory, starting services like Traefik.

## Step 5: Start the Dashboard

This command specifically pulls the M3TAL dashboard image and starts its container. It respects the `DASHBOARD_EXPOSE_MODE` setting from your configuration.

```bash
m3tal dash up
```

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard URL.

*   If `DASHBOARD_EXPOSE_MODE` is set to `local`, access the dashboard at:
    `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address or `localhost`)

*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik` and `m3tal up` has been run, access the dashboard at:
    `http://dash.DOMAIN` (replace `DOMAIN` with your configured domain, e.g., `dash.localhost`)

## Step 7: Log In

You will be presented with a login screen.

*   **Default Credentials**:
    *   Username: `admin`
    *   Password: `admin_pass` (This is the default value for `ADMIN_PASSWORD` if you did not change it during the wizard. It is highly recommended to change this.)

To change the admin password after logging in, use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the dashboard.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure to manage its configuration and state.

| Path                      | Purpose                                                                  |
| :------------------------ | :----------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | The primary configuration file for M3TAL. Modified by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite database storing M3TAL's operational state. Auto-created by the API daemon. |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory for all M3TAL Docker Compose files and related configurations. |
| `/docker/users.json`      | Stores dashboard user credentials. Managed by `m3tal dashpass`.          |

## Port Map

| Port | Service             | Access Method                                           |
| :--- | :------------------ | :------------------------------------------------------ |
| 80   | Traefik (HTTP)      | Accessible publicly if Traefik is exposed and configured. |
| 8080 | M3TAL API Daemon    | Accessible on the host machine.                         |
| 8081 | Traefik Dashboard   | Accessible on the host machine only.                    |
| 8082 | M3TAL Dashboard     | Accessible directly via host port (local mode) or via Traefik (traefik mode). |

## Firewall Note

If you are exposing Traefik to the public internet and wish to allow HTTP traffic, ensure your firewall permits connections on port 80:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`:

*   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```