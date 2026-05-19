# M3TAL: Get Started Guide

This guide will walk you through the initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify their installation by running the following commands in your terminal:

```bash
docker --version
docker compose version
```

If these commands return version information, you are ready to proceed. If not, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is distributed as a Debian package for easy installation. Execute the following commands to add the M3TAL repository and install the `m3tal` CLI:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The `m3tal config wizard` command will guide you through setting up essential configuration options for M3TAL.

```bash
sudo m3tal config wizard
```

You will be presented with a series of prompts. Here's a breakdown of what each prompt means:

*   **`DASHBOARD_EXPOSE_MODE`**: This determines how the M3TAL dashboard will be accessed.
    *   `local`: The dashboard will be directly accessible via its port (default 8082). This is the recommended mode for initial setup and LAN-only access.
    *   `traefik`: The dashboard will be exposed through the Traefik reverse proxy, accessible via a domain name (e.g., `dash.yourdomain.com`). This is useful if you plan to use Traefik for managing multiple services.
*   **`DOMAIN`**: If you choose `traefik` mode, enter your domain name here. If you are using `local` mode or testing locally, `localhost` is acceptable.
*   **`TRAEFIK_WEB_PORT`**: The port Traefik will listen on for HTTP traffic. The default is 80.
*   **`DASHBOARD_PORT`**: The internal port the M3TAL dashboard container will run on. The default is 8082.
*   **`PUID` / `PGID`**: These are the User ID and Group ID that the Docker containers will run as. It's generally recommended to use your primary user's IDs to avoid permission issues. You can find these by running `id -u` and `id -g` in your terminal.
*   **`TZ`**: Your local timezone. This ensures logs and timestamps are accurate. Select the appropriate timezone from the list (e.g., `America/New_York`).
*   **`BASE_STORAGE_PATH`**: The base directory where M3TAL will store its data and configurations.

The wizard will save your choices to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

M3TAL utilizes Traefik as a reverse proxy to manage incoming traffic for various services. To start the routing stack, run:

```bash
m3tal up
```

This command reads all `*-compose.yml` files found in the `/docker/` directory and starts the services defined within them using Docker Compose. This includes Traefik, which will be configured to listen on port 80.

## Step 5: Start the Dashboard

Now, let's start the M3TAL dashboard. The `m3tal dash up` command handles downloading the necessary image and starting the dashboard container based on your `DASHBOARD_EXPOSE_MODE` setting from the wizard.

```bash
m3tal dash up
```

This command will pull the `ghcr.io/jakej985-rgb/m3tal-godash:debug` image if it's not already present and then start the `m3tal-dashboard` container.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to one of the following addresses:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local`: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or use `http://localhost:8082` if accessing from the same machine).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`: `http://dash.DOMAIN` (replace `yourdomain.com` with the domain you configured in the wizard). Ensure Traefik is running correctly for this to work.

## Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: `admin_pass` (This is the default value for `ADMIN_PASSWORD` in your `.env` file. It is highly recommended to change this.)

*   **Changing Your Password:**
    To change the dashboard password, use the following command:

    ```bash
    sudo m3tal dashpass
    ```
    Follow the prompts to set a new password.

---

## Filesystem Contract

M3TAL relies on specific file locations for its configuration and state management.

| Path                     | Purpose                                                                                                  |
| :----------------------- | :------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary M3TAL configuration file. This file is managed by the `m3tal config wizard` and `m3tal config set` commands. |
| `/var/lib/m3tal/state.db`| SQLite database used to store M3TAL's operational state. This file is automatically created and managed by the M3TAL API daemon. |
| `/docker`                | This is a symbolic link that points to `/opt/m3tal/stack/`. It serves as the user-facing directory for all Docker Compose files that define M3TAL's services. |
| `/docker/users.json`     | Stores the M3TAL dashboard credentials. This file is managed by the `m3tal dashpass` command.           |

## Port Table

The following ports are used by M3TAL and its associated services:

| Port | Service             | Access Method                                     |
| :--- | :------------------ | :------------------------------------------------ |
| 80   | Traefik (HTTP)      | Publicly accessible if Traefik is exposed externally. |
| 8080 | M3TAL API Daemon    | Accessible only from within the host environment. |
| 8081 | Traefik Dashboard   | Accessible from the host machine only (127.0.0.1). |
| 8082 | M3TAL Dashboard     | Accessible via direct port binding (`local` mode) or through Traefik (`traefik` mode). |

## Firewall Note

If you intend to expose Traefik to the public internet (e.g., for accessing services via domains), you must allow traffic on port 80 through your firewall. If you are using `ufw` (Uncomplicated Firewall), run:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check Status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View Logs (real-time):**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```