# M3TAL Getting Started Guide

This guide will walk you through the installation and initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before you begin, ensure you have the following software installed on your system:

*   **Docker Engine:** The containerization platform.
*   **Docker Compose V2:** The tool for defining and running multi-container Docker applications.

You can verify your installation by running the following commands:

```bash
docker --version
docker compose version
```

If these commands do not output version information, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is installed using your system's Advanced Package Tool (APT). Execute the following commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, you will run the M3TAL configuration wizard to set up essential parameters.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for various settings. Here's a breakdown of what each prompt typically asks for:

*   **`DASHBOARD_EXPOSE_MODE`**: This determines how the M3TAL dashboard is made accessible.
    *   `local`: Exposes the dashboard directly via a port on your host machine (e.g., `http://YOUR_IP:8082`). This is the default and recommended for initial setup and LAN-only access.
    *   `traefik`: Configures Traefik (a reverse proxy) to route traffic to the dashboard using a domain name (e.g., `http://dash.yourdomain.com`). This requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The port on which the dashboard will be accessible. The default is `8082`.
*   **`DOMAIN`**: Your primary domain name. This is used for Traefik routing if `DASHBOARD_EXPOSE_MODE` is set to `traefik`. Defaults to `localhost`.
*   **`API_TOKEN`**: A secret token for API authentication. **Change this from the default.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. **Change this from the default.**
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard. **Change this from the default.**
*   **`PUID` / `PGID`**: The User ID and Group ID for running Docker containers. Often `1000` on Linux systems.
*   **`TZ`**: Your timezone. Set this to your local timezone (e.g., `America/New_York`).
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL's data and configuration files. Defaults to `./data` relative to where the `m3tal` command is run, but will be absolute after configuration.
*   **`CONFIG_PATH`**: The directory for M3TAL configuration files.
*   **`MEDIA_PATH`**: The directory for media files.
*   **`DOWNLOADS_PATH`**: The directory for downloads.

**Important:** Pay close attention to the prompts and provide appropriate values. It's highly recommended to change default passwords and tokens for security.

## Step 4: Start the Routing Stack (Traefik)

The M3TAL ecosystem utilizes Docker Compose to manage its services. The `m3tal up` command orchestrates all defined compose files located in the `/docker/` directory. This typically includes Traefik, which acts as your reverse proxy.

```bash
m3tal up
```

This command will:
1.  Locate all `*.yml` and `*-compose.yml` files within `/docker/`.
2.  Start the Docker containers defined in these files using Docker Compose. This includes the Traefik reverse proxy, which is essential for routing traffic to various M3TAL services and any other services you add.

## Step 5: Start the Dashboard

The M3TAL dashboard provides a web interface for managing your services. The `m3tal dash up` command is a convenient way to ensure the dashboard container is running and up-to-date.

```bash
m3tal dash up
```

This command will:
1.  Download the latest `m3tal-compose.yml` and any relevant override files (like `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) from the M3TAL repository.
2.  Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` Docker container, applying the appropriate configuration based on your chosen expose mode.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's address. The exact URL depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Access the dashboard via:
    `http://YOUR_IP:8082`
    Replace `YOUR_IP` with your server's local IP address. You can also use `http://localhost:8082` if you are accessing it from the same machine.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Access the dashboard via your configured domain:
    `http://dash.DOMAIN`
    Replace `DOMAIN` with the domain you set during the wizard (e.g., `http://dash.m3tal.local`). This assumes Traefik is running and properly configured to route to the dashboard.

## Step 7: Log In

You will be presented with the M3TAL dashboard login page.

*   **Username:** `admin`
*   **Password:** The password you set during the configuration wizard for `ADMIN_PASSWORD`.

If you need to change your dashboard password after initial setup, you can use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you for the new password and update the credentials securely.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure for its configuration and state management. Understanding these paths is crucial for advanced customization and troubleshooting.

| Path                      | Purpose                                                                                                                  | Notes                                                                    |
| :------------------------ | :----------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | Primary M3TAL environment configuration file. Stores all user-defined settings and defaults.                             | Managed by `m3tal config wizard` and `m3tal config set`.               |
| `/var/lib/m3tal/state.db` | SQLite database storing the persistent state of M3TAL services and configurations.                                       | Auto-created by the M3TAL API daemon.                                    |
| `/opt/m3tal/stack/`       | The canonical directory containing M3TAL's core Docker Compose files, Traefik configurations, and other stack assets.    | This is the source of truth for M3TAL's operational files.               |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all M3TAL stack operations.            | Commands like `m3tal up` operate on files within this directory.         |
| `/docker/users.json`      | Stores the username and hashed passwords for the M3TAL dashboard.                                                        | Managed by `m3tal dashpass` and the dashboard itself.                    |
| `/docker/dynamic/`        | Directory for Traefik dynamic configuration files. This allows for hot-reloading of routing rules without restarting Traefik. | Used when `DASHBOARD_EXPOSE_MODE=traefik`.                               |

---

## Port Table

The following ports are commonly used by M3TAL services.

| Port | Service / Protocol | Description                                     | Access Method                                      |
| :--- | :----------------- | :---------------------------------------------- | :------------------------------------------------- |
| 80   | Traefik (HTTP)     | Inbound HTTP traffic, routed by Traefik.        | Public (if Traefik is exposed to the internet).    |
| 8080 | M3TAL API (Go)     | The primary M3TAL API daemon.                   | Host-local (accessed internally by other services). |
| 8081 | Traefik Dashboard  | Traefik's own administrative interface.         | Host-local access only (e.g., `127.0.0.1:8081`).   |
| 8082 | M3TAL Dashboard    | The M3TAL web dashboard.                        | Directly (if `local` mode) or via Traefik (`traefik` mode). |

---

## Firewall Note

If you intend to expose Traefik (and thus your M3TAL services) to the internet, you must open the necessary ports in your firewall. For example, if you are using `ufw` (Uncomplicated Firewall) and Traefik is listening on port 80:

```bash
sudo ufw allow 80
```

Consider also opening port `443` if you plan to configure HTTPS with Traefik.

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`.

*   **Check the status of the M3TAL API service:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View real-time logs for the M3TAL API service:**
    ```bash
    sudo journalctl -u m3tal-api -f
    ```
    Press `Ctrl+C` to exit the log stream.

*   **Restart the M3TAL API service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```