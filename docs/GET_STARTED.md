# M3TAL System Setup Guide

This guide will walk you through the initial setup of the M3TAL system.

## Step 1: Prerequisites

Before proceeding, ensure that the following software is installed on your system:

*   **Docker Engine:** The containerization platform.
*   **Docker Compose V2:** The tool for defining and running multi-container Docker applications.

You can verify their installation by running the following commands:

```bash
docker --version
docker compose version
```

If these commands do not return version information, please install Docker Engine and Docker Compose V2 first.

## Step 2: Install M3TAL via APT

M3TAL is installed using the Advanced Packaging Tool (APT). Execute the following three commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up essential parameters for your system. Run the following command:

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's an explanation of each:

*   **`DOMAIN`**: The primary domain name for accessing your M3TAL services. For local testing, `localhost` is often used. If you have a registered domain, enter it here.
*   **`TRAEFIK_WEB_PORT`**: The port Traefik will listen on for incoming HTTP traffic. The default is `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik will listen on for incoming HTTPS traffic. The default is `443`.
*   **`DASHBOARD_PORT`**: The port the M3TAL Dashboard will be accessible on. The default is `8082`.
*   **`HTTP_PORT`**: The port the M3TAL API daemon will listen on. The default is `8080`.
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL Dashboard. **It is highly recommended to change this from the default immediately after setup.**
*   **`PUID` / `PGID`**: The User ID and Group ID for running containers. This is typically your user's ID for file permissions. You can find yours with `id -u` and `id -g`.
*   **`TZ`**: Your timezone. This is important for accurate logging and scheduling. Example: `America/Denver`.
*   **`STATE_DIR`**: The directory where M3TAL will store its state database. The default is `./state`.
*   **`BASE_STORAGE_PATH`**: The base directory for storing data for your stacks. The default is `./data`.
*   **`MEDIA_PATH`**: The directory for media files. Defaults to `./data/media`.
*   **`CONFIG_PATH`**: The directory for configuration files. Defaults to `./data/config`.
*   **`DOWNLOADS_PATH`**: The directory for downloaded files. Defaults to `./data/downloads`.

## Step 4: Start the Routing Stack (Traefik)

M3TAL uses Traefik as a reverse proxy to manage incoming traffic and route it to the appropriate services. This command starts all Docker Compose files located in the `/docker/` directory, including Traefik and any other default routing configurations.

```bash
m3tal up
```

This command will pull necessary Docker images and start the containers defined in the compose files within `/docker/`.

## Step 5: Start the Dashboard

The M3TAL Dashboard provides a web interface for managing your system. This command specifically pulls the dashboard image and starts its container.

```bash
m3tal dash up
```

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard address.

*   If you are using the default `localhost` domain, access it at: `http://localhost:8082`
*   If you have configured a custom `DOMAIN` and Traefik is correctly set up, you can access it via that domain: `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with your configured domain).

## Step 7: Log In to the Dashboard

When you access the dashboard for the first time, you will be presented with a login screen.

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (or `admin_pass` if you used the default and haven't changed it).

**IMPORTANT:** It is crucial to change your dashboard password immediately after logging in. You can do this via the command line:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password for the dashboard.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure for configuration and state management:

| Path                        | Purpose                                                                           | Notes                                                                                      |
| :-------------------------- | :-------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file.                                                       | Managed by `m3tal config wizard` and `m3tal config set`.                                 |
| `/var/lib/m3tal/state.db`   | SQLite state database for M3TAL.                                                  | Auto-created by the API daemon on first run.                                               |
| `/opt/m3tal/stack/`         | Canonical directory containing M3TAL's core Docker Compose files and Traefik config. |                                                                                            |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`.                                                   | This is the user-facing path for all stack operations (`m3tal up`).                        |
| `/docker/users.json`        | Stores dashboard credentials.                                                     | Managed by `m3tal dashpass`.                                                               |
| `/docker/dynamic/`          | Directory for dynamic Traefik routing configuration files.                        | Hot-reloaded by Traefik.                                                                   |
| `/docker/routing-compose.yml` | Docker Compose file for Traefik and other routing services.                     |                                                                                            |
| `/docker/m3tal-compose.yml` | Docker Compose file specifically for the M3TAL dashboard.                         |                                                                                            |

## Port Map

The following ports are used by M3TAL services:

| Port   | Service             | Access Type             |
| :----- | :------------------ | :---------------------- |
| 80     | Traefik HTTP        | Public (if exposed)     |
| 8080   | M3TAL API Daemon    | Host-local (via Traefik) |
| 8081   | Traefik Dashboard   | Host-local only         |
| 8082   | M3TAL Dashboard     | Via Traefik or direct   |

## Firewall Note

If Traefik is exposed to the public internet (i.e., port 80 is accessible externally), you should configure your firewall to allow incoming traffic on that port. For example, using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`.

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View logs (follow in real-time):**
    ```bash
    journalctl -u m3tal-api -f
    ```