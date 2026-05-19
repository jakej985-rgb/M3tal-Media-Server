# M3TAL - Getting Started Guide

This guide will walk you through the initial setup of M3TAL on your system.

## Step 1: Prerequisites

Before proceeding, ensure you have the following software installed on your system:

*   **Docker Engine**: The containerization platform.
*   **Docker Compose V2**: The tool for defining and running multi-container Docker applications.

You can verify your installation by running the following command:

```bash
docker --version && docker compose version
```

If either command fails, please install Docker Engine and Docker Compose V2 according to your operating system's documentation.

## Step 2: Install M3TAL via APT

M3TAL is installed using your system's package manager. Execute the following three commands to add the M3TAL repository and install the CLI:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, you'll need to configure M3TAL. Run the interactive configuration wizard:

```bash
sudo m3tal config wizard
```

This command will guide you through several prompts to set up your M3TAL environment. Here's an explanation of each prompt:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL dashboard is accessed.
    *   `local`: Exposes the dashboard directly via a port on your host machine (default, recommended for initial setup).
    *   `traefik`: Exposes the dashboard through Traefik, requiring Traefik to be running and accessible via a domain name.
*   **`DASHBOARD_PORT`**: The port on your host machine that the dashboard will listen on (if `DASHBOARD_EXPOSE_MODE=local`). Defaults to `8082`.
*   **`DOMAIN`**: The domain name you intend to use for M3TAL services when `DASHBOARD_EXPOSE_MODE=traefik`. Defaults to `localhost`.
*   **`DASHBOARD_SECRET`**: A secret key used for securing dashboard sessions. It's highly recommended to change this from the default.
*   **`API_TOKEN`**: An API token for authenticating with the M3TAL API. Change this from the default.
*   **`ADMIN_PASSWORD`**: The password for the default administrator user of the M3TAL dashboard. Change this from the default.
*   **`PUID` / `PGID`**: The User ID (UID) and Group ID (GID) to run Docker containers with. Using your host user's UID/GID is common. You can find yours by running `id -u` and `id -g`.
*   **`BASE_STORAGE_PATH`**: The base directory for storing M3TAL-related data, including configurations and media.
*   **`MEDIA_PATH`**: The specific directory for media files within the `BASE_STORAGE_PATH`.
*   **`CONFIG_PATH`**: The specific directory for configuration files within the `BASE_STORAGE_PATH`.
*   **`DOWNLOADS_PATH`**: The directory where downloaded files will be stored.
*   **`TZ`**: Your timezone. This is important for correct log timestamps.

## Step 4: Start the Routing Stack (Traefik)

M3TAL uses Traefik as a reverse proxy to manage incoming traffic and route it to various services. To start the Traefik routing stack, run:

```bash
m3tal up
```

This command orchestrates the startup of all Docker Compose files located within the `/docker/` directory, including Traefik itself.

## Step 5: Start the Dashboard

Next, you need to start the M3TAL dashboard. The command will pull the necessary Docker image if it's not already present and then start the dashboard container according to your `DASHBOARD_EXPOSE_MODE` setting from the wizard.

```bash
m3tal dash up
```

## Step 6: Access the Dashboard

Once the dashboard is running, you can access it through your web browser. The access method depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open your browser and navigate to `http://YOUR_IP:8082`. Replace `YOUR_IP` with your server's IP address or `localhost` if you are accessing it from the same machine.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open your browser and navigate to `http://dash.DOMAIN`. Replace `DOMAIN` with the domain you configured during the wizard (e.g., `http://dash.example.com`). Traefik must be running for this to work.

## Step 7: Log In

You will be presented with the M3TAL dashboard login screen.

*   **Username:** `admin`
*   **Password:** The password you set for `ADMIN_PASSWORD` during the configuration wizard.

If you need to change the dashboard password after initial setup, use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you for the new password and update the `users.json` file accordingly.

---

## Filesystem Contract

M3TAL organizes its configuration and data in specific locations. Understanding this structure is key to managing your installation:

| Path                     | Purpose                                                                                    |
| :----------------------- | :----------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | The primary environment file for M3TAL configuration. Modified by `m3tal config wizard`.   |
| `/var/lib/m3tal/state.db`| SQLite database used by the M3TAL API daemon to store operational state. Auto-created.     |
| `/opt/m3tal/stack/`      | The canonical directory containing M3TAL's core Docker Compose files and Traefik configuration. |
| `/docker`                | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations and user-added compose files. |
| `/docker/users.json`     | Stores dashboard user credentials. Managed by `m3tal dashpass`.                              |

## Port Table

The following ports are utilized by M3TAL and its components:

| Port   | Service           | Access Method                      | Notes                                      |
| :----- | :---------------- | :--------------------------------- | :----------------------------------------- |
| 80     | Traefik (HTTP)    | Public (via firewall)              | Entrypoint for routed services.            |
| 8080   | M3TAL API Daemon  | Host-local                         | Internal API for the dashboard and CLI.    |
| 8081   | Traefik Dashboard | Host-local                         | Only accessible directly on the host machine. |
| 8082   | M3TAL Dashboard   | Direct port (`local` mode) or Traefik (`traefik` mode) | Access point for the M3TAL web UI.   |

## Firewall Note

If Traefik is exposed to the internet, you may need to allow incoming traffic on port 80 through your firewall. For example, if you are using `ufw`:

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