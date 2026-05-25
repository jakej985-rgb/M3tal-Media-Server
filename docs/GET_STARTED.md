# M3TAL Ecosystem: Getting Started Guide

This guide will walk you through the initial setup and configuration of the M3TAL ecosystem.

## Step 1: Prerequisites

Before proceeding, ensure you have the following software installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify your installations by running the following commands:

```bash
docker --version
docker compose version
```

If either command fails or reports an old version, please refer to the official Docker documentation for installation and upgrade instructions.

## Step 2: Install M3TAL via APT

M3TAL is installed using the Advanced Packaging Tool (APT). Execute the following commands to add the M3TAL repository and install the package:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The `m3tal config wizard` command will guide you through setting up essential M3TAL configurations. You will be prompted for the following settings:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL dashboard is accessed.
    *   `local`: Exposes the dashboard directly via a port. Recommended for initial setup and LAN-only access.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy, allowing access via a domain name.
*   **`DASHBOARD_PORT`**: The port number the dashboard will listen on (default is `8082`). This is only relevant if `DASHBOARD_EXPOSE_MODE` is set to `local`.
*   **`DOMAIN`**: The domain name you intend to use for accessing M3TAL services if `DASHBOARD_EXPOSE_MODE` is `traefik`. Defaults to `localhost`.
*   **`PUID`**: The User ID to run containers as. Defaults to your current user's ID.
*   **`PGID`**: The Group ID to run containers as. Defaults to your current user's group ID.
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL data and configurations. Defaults to `./data` within the directory where `m3tal` commands are run.
*   **`CONFIG_PATH`**: The specific path for configuration files within `BASE_STORAGE_PATH`. Defaults to `./data/config`.
*   **`MEDIA_PATH`**: The directory for media storage. Defaults to `./data/media`.
*   **`DOWNLOADS_PATH`**: The directory for downloads. Defaults to `./data/downloads`.
*   **`TZ`**: Your timezone (e.g., `America/Denver`).

Execute the wizard with:

```bash
sudo m3tal config wizard
```

Follow the on-screen prompts to enter your desired values.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command initiates the core M3TAL services, including the Traefik reverse proxy. This command processes all Docker Compose files located in `/docker/`.

```bash
m3tal up
```

This command will pull necessary Docker images and start the containers defined in your stack.

## Step 5: Start the Dashboard

To start the M3TAL dashboard, use the following command:

```bash
m3tal dash up
```

This command will:
1.  Pull the latest `m3tal-dashboard` Docker image.
2.  Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Start the dashboard container, applying the appropriate configuration based on your chosen expose mode.

## Step 6: Access the M3TAL Dashboard

You can now access the M3TAL dashboard in your web browser. The access method depends on your `DASHBOARD_EXPOSE_MODE` setting during the configuration wizard:

*   **If `DASHBOARD_EXPOSE_MODE=local`:** Open your browser to `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address or `localhost` if accessing from the same machine).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:** Open your browser to `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured, e.g., `http://dash.localhost` if you left the default).

## Step 7: Log In

Upon accessing the dashboard, you will be presented with a login screen.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: `admin_pass` (This is the default value for the `ADMIN_PASSWORD` environment variable. It's strongly recommended to change this immediately.)

**Changing Your Password:**

To change the dashboard password, use the `m3tal dashpass` command followed by the new password:

```bash
sudo m3tal dashpass YOUR_NEW_SECURE_PASSWORD
```

Replace `YOUR_NEW_SECURE_PASSWORD` with a strong, unique password.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure to manage its configuration and state. Understanding these paths is crucial for advanced usage and troubleshooting.

| Path                       | Purpose                                                                                                  |
| :------------------------- | :------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary M3TAL configuration file. Contains all environment variables. Managed by `m3tal config wizard`.   |
| `/var/lib/m3tal/state.db`  | SQLite database storing M3TAL's internal state. Automatically created and managed by the API daemon.     |
| `/opt/m3tal/stack/`        | Canonical directory for M3TAL's core Docker Compose files and Traefik configuration.                     |
| `/docker`                  | **Symlink → `/opt/m3tal/stack/`**. This is the user-facing path for all Docker Compose operations.         |
| `/docker/users.json`       | Stores dashboard credentials. Managed by the `m3tal dashpass` command.                                   |

## Port Map

The following ports are used by M3TAL services:

| Port | Service            | Access                                    |
| :--- | :----------------- | :---------------------------------------- |
| 80   | Traefik (HTTP)     | Public access (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API Daemon   | Host-local access only                      |
| 8081 | Traefik Dashboard  | Host-local access only                      |
| 8082 | M3TAL Dashboard    | Direct port access (if `local`) or via Traefik (if `traefik`) |

## Firewall Note

If your server is protected by a firewall and you intend to expose M3TAL services to the internet (e.g., using `DASHBOARD_EXPOSE_MODE=traefik`), ensure that port 80 is allowed through your firewall. If using UFW (Uncomplicated Firewall):

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard `systemctl` and `journalctl` commands.

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```