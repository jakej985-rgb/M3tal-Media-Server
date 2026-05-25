# M3TAL: Getting Started Guide

This guide will walk you through the initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

Verify your installation by running the following commands:

```bash
docker --version
docker compose version
```

If these commands do not output version information, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is installed using the Advanced Package Tool (APT) from our official repository. Execute the following commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, you must run the configuration wizard to set up M3TAL's core settings.

```bash
sudo m3tal config wizard
```

You will be prompted with a series of questions. Here's a brief explanation of each:

*   **`DASHBOARD_PORT (default: 8082)`**: The local port on which the M3TAL dashboard will be accessible.
*   **`DASHBOARD_EXPOSE_MODE (default: local)`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a host port. Ideal for LAN-only access.
    *   `traefik`: Exposes the dashboard through Traefik (if configured). Requires a domain name.
*   **`HTTP_PORT (default: 8080)`**: The port on which the M3TAL API daemon will listen. This is typically accessed by other M3TAL services or internal tooling.
*   **`STATE_DIR (default: ./state)`**: The directory where M3TAL's state database (`state.db`) will be stored. This is relative to the `CONFIG_PATH`.
*   **`LOG_LEVEL (default: info)`**: Sets the verbosity of M3TAL's logs. Options typically include `debug`, `info`, `warn`, `error`.
*   **`DASHBOARD_SECRET (default: change_me_immediately)`**: A secret key used for securing the dashboard. **It is highly recommended to change this from the default.**
*   **`API_TOKEN (default: change_me_api_token)`**: A token used for authenticating API requests. **It is highly recommended to change this from the default.**
*   **`ADMIN_PASSWORD (default: admin_pass)`**: The default password for the dashboard administrator. **It is highly recommended to change this from the default.**
*   **`NETWORK_NAME (default: m3tal)`**: The name of the Docker network M3TAL services will use.
*   **`LOCAL_IP (default: 127.0.0.1)`**: The local IP address of your host machine.
*   **`DOMAIN (default: localhost)`**: The domain name you will use to access services. Use `localhost` for local-only access or your actual domain name if using Traefik.
*   **`VPN_USER (default: user)`**: Username for VPN connection (if applicable).
*   **`VPN_PASSWORD (default: password)`**: Password for VPN connection (if applicable).
*   **`BASE_STORAGE_PATH (default: ./data)`**: The base directory for M3TAL data storage.
*   **`MEDIA_PATH (default: ./data/media)`**: Path for media files.
*   **`CONFIG_PATH (default: ./data/config)`**: Path for configuration files, including `.env` and state data.
*   **`DOWNLOADS_PATH (default: ./data/downloads)`**: Path for downloaded files.
*   **`PUID (default: 1000)`**: The User ID to run Docker containers as.
*   **`PGID (default: 1000)`**: The Group ID to run Docker containers as.
*   **`TZ (default: America/Denver)`**: Your local timezone.
*   **`TRAEFIK_WEB_PORT (default: 80)`**: The host port Traefik will listen on for HTTP traffic.
*   **`TRAEFIK_WEBHTTPS_PORT (default: 443)`**: The host port Traefik will listen on for HTTPS traffic (if configured).
*   **`TRAEFIK_DASHBOARD_PORT (default: 8080)`**: The host port Traefik's own dashboard will be accessible on (usually for internal use or troubleshooting).
*   **`DEBUG_MODE (default: false)`**: Enables or disables debug logging for M3TAL.
*   **`METRICS_ENABLED (default: true)`**: Enables or disables the collection of performance metrics.

## Step 4: Start the Routing Stack (Traefik)

The routing stack, primarily Traefik, acts as your reverse proxy and service gateway. It is essential for managing access to your deployed services.

```bash
m3tal up
```

This command orchestrates the startup of all Docker Compose files located within the `/docker/` directory. This includes Traefik and any other services you may have added.

## Step 5: Start the Dashboard

The M3TAL dashboard provides a user-friendly interface for managing your M3TAL ecosystem.

```bash
m3tal dash up
```

This command will:
1.  Pull the latest `m3tal-dashboard` Docker image from the container registry.
2.  Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Start the dashboard container, applying the correct configuration for your chosen exposure mode.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard's address. The URL depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Access the dashboard at `http://YOUR_IP:8082` or `http://localhost:8082`. Replace `YOUR_IP` with your server's IP address.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Access the dashboard at `http://dash.DOMAIN`. Ensure Traefik is running (`m3tal up`) and that your `DOMAIN` is correctly configured in `/etc/m3tal/.env`.

## Step 7: Log In to the Dashboard

When prompted for credentials:

*   **Username:** `admin`
*   **Password:** The password you set during the configuration wizard (default is `admin_pass`).

**To change your dashboard password after initial setup, use the following command:**

```bash
sudo m3tal dashpass
```

This will prompt you for the new password and update the credential store.

---

## Filesystem Contract

M3TAL relies on a specific filesystem structure for configuration and state. Understanding these paths is crucial for advanced management and troubleshooting.

| Path                     | Purpose                                                                      |
| :----------------------- | :--------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | The primary M3TAL configuration file. All user-defined settings are stored here and managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db` | SQLite database storing M3TAL's operational state, service configurations, and other critical data. Automatically created and managed by the `m3tal-api` service. |
| `/opt/m3tal/stack/`      | The canonical directory containing M3TAL's core Docker Compose files and Traefik configuration. |
| `/docker`                | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory for managing your M3TAL service stacks. New stack compose files should be placed here. |
| `/docker/users.json`     | Stores dashboard user credentials. Managed by `m3tal dashpass`. |

---

## Port Map

The following ports are used by M3TAL and its components:

| Port | Service        | Access                                               |
| :--- | :------------- | :--------------------------------------------------- |
| 80   | Traefik        | Public (if `DASHBOARD_EXPOSE_MODE=traefik`)          |
| 8080 | M3TAL API      | Host-local (accessible by internal M3TAL services) |
| 8081 | Traefik Dashboard | Host-local only (for Traefik management)             |
| 8082 | M3TAL Dashboard | Direct port (if `local` mode) or via Traefik         |

---

## Firewall Note

If you are exposing Traefik to the public internet (e.g., `DASHBOARD_EXPOSE_MODE=traefik`), ensure that port 80 is open in your firewall. For example, if using `ufw`:

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard systemd commands.

*   **Check the status of the API service:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View live logs from the API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```