# M3TAL - Get Started

This guide will walk you through the initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before you begin, ensure you have the following software installed on your system:

*   **Docker Engine**: The core containerization platform.
*   **Docker Compose V2**: The tool for defining and running multi-container Docker applications.

To verify your installation, run the following command in your terminal:

```bash
docker --version && docker compose version
```

If either command fails, please install Docker Engine and Docker Compose V2 according to your operating system's documentation.

## Step 2: Install M3TAL via APT

M3TAL provides an official APT repository for easy installation. Follow these steps precisely:

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

You will be prompted with several questions. Here's a brief explanation of each:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL dashboard is accessible.
    *   `local`: Exposes the dashboard directly via a port on your host machine (e.g., `http://YOUR_IP:8082`). This is the default and recommended for initial setup.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy, typically via a domain name (e.g., `http://dash.DOMAIN`). Requires Traefik to be running.
*   **`DOMAIN`**: The primary domain name for your M3TAL setup. If you are using `traefik` mode for dashboard access or other services, this is the domain that will be used. For `local` mode, it can often be left as `localhost`.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`PUID` / `PGID`**: The User ID and Group ID for running containers. Typically, this should match your current user's UID/GID to ensure proper file permissions. You can find these by running `id -u` and `id -g`.
*   **`BASE_STORAGE_PATH`**: The root directory for M3TAL data, configurations, and media.
*   **`CONFIG_PATH`**: The directory within `BASE_STORAGE_PATH` where M3TAL stores its configuration files.
*   **`MEDIA_PATH`**: The directory within `BASE_STORAGE_PATH` for media files.
*   **`DOWNLOADS_PATH`**: The directory within `BASE_STORAGE_PATH` for download-related data.
*   **`TZ`**: Your timezone, used for consistent logging and scheduling.

The wizard will save your choices to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

This command starts the core routing infrastructure, primarily Traefik, which acts as a reverse proxy for your services.

```bash
m3tal up
```

This command orchestrates the Docker Compose files located in `/docker/`. It will pull necessary images and start containers for the routing stack (Traefik).

## Step 5: Start the Dashboard

Next, start the M3TAL dashboard container.

```bash
m3tal dash up
```

This command will:
1.  Download the necessary Docker Compose configuration files for the dashboard.
2.  Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Pull the `m3tal-dashboard` Docker image if it's not already present.
4.  Start the `m3tal-dashboard` container based on your chosen exposure mode.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default):
    `http://YOUR_IP:8082`
    (Replace `YOUR_IP` with your server's IP address or `localhost` if accessing from the same machine.)

*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`:
    `http://dash.DOMAIN`
    (Replace `DOMAIN` with the domain you configured during the wizard. Ensure Traefik is running and your DNS is correctly configured.)

## Step 7: Log In

You will be presented with the M3TAL dashboard login screen.

*   **Default Credentials**:
    *   Username: `admin`
    *   Password: `admin_pass` (This is the default value for `ADMIN_PASSWORD` in `/etc/m3tal/.env`)

**It is crucial to change the default password immediately after your first login.** You can do this via the command line:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the dashboard.

---

## Filesystem Contract

M3TAL organizes its configuration and state in specific locations:

| Path                      | Purpose                                                                                                  |
| :------------------------ | :------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | Primary environment configuration file. Managed by the `m3tal config wizard` and `m3tal config set`.     |
| `/var/lib/m3tal/state.db` | SQLite database used by the M3TAL API daemon to store system state. Automatically created on first run. |
| `/opt/m3tal/stack/`       | Canonical directory containing M3TAL's core Docker Compose files and Traefik configuration.              |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations.   |
| `/docker/users.json`      | Stores dashboard user credentials. Managed by `m3tal dashpass`.                                          |

---

## Port Map

The following ports are relevant to M3TAL operations:

| Port   | Service           | Access                                                        |
| :----- | :---------------- | :------------------------------------------------------------ |
| 80     | Traefik           | Public access for HTTP services (when Traefik is enabled).    |
| 8080   | M3TAL API Daemon  | Internal access from other containers or `host.docker.internal`. |
| 8081   | Traefik Dashboard | Host-local access only.                                       |
| 8082   | M3TAL Dashboard   | Directly via host IP/port (local mode) or via Traefik (traefik mode). |

---

## Firewall Note

If Traefik is exposed to the internet and you are using a firewall like `ufw`, ensure that port 80 (HTTP) is open:

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```

*   **View logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```
    (Press `Ctrl+C` to exit the log stream.)