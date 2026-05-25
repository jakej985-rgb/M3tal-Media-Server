# M3TAL Get Started Guide

This guide will walk you through the initial setup of the M3TAL ecosystem.

## 1. Prerequisites

Before proceeding, ensure the following are installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify their installation by running:

```bash
docker --version
docker compose version
```

## 2. Install M3TAL via APT

Execute the following commands to add the M3TAL repository and install the CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The `m3tal config wizard` command will guide you through setting up essential configuration parameters.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a brief explanation of each:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a port. Ideal for LAN-only setups and initial testing.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy using a domain name. Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state, including the database.
*   **`LOG_LEVEL`**: Controls the verbosity of M3TAL's logging (e.g., `info`, `debug`).
*   **`DASHBOARD_SECRET`**: A secret key for securing the dashboard. **It is highly recommended to change this from the default.**
*   **`API_TOKEN`**: A token for authenticating with the M3TAL API. **It is highly recommended to change this from the default.**
*   **`ADMIN_PASSWORD`**: The password for the dashboard's administrative user. **It is highly recommended to change this from the default.**
*   **`DOMAIN`**: The domain name used for accessing services when `DASHBOARD_EXPOSE_MODE` is set to `traefik`. Defaults to `localhost`.
*   **`PUID` / `PGID`**: The User ID and Group ID for running Docker containers. Often defaults to `1000`.
*   **`TZ`**: Your local timezone (e.g., `America/Denver`).
*   **`TRAEFIK_WEB_PORT`**: The port Traefik listens on for HTTP traffic. Defaults to `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik listens on for HTTPS traffic. Defaults to `443`.

## 4. Start the Routing Stack (Traefik)

This command starts the core routing infrastructure, primarily Traefik, which acts as a reverse proxy for your services.

```bash
m3tal up
```

This command orchestrates all Docker Compose files located in the `/docker/` directory. This typically includes Traefik and any other core M3TAL services.

## 5. Start the Dashboard

This command pulls the M3TAL dashboard image and starts its container, respecting your `DASHBOARD_EXPOSE_MODE` setting from the configuration wizard.

```bash
m3tal dash up
```

## 6. Access the Dashboard

Open your web browser and navigate to the dashboard URL. The exact URL depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local`**: Navigate to `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address or `localhost`).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`**: Navigate to `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured, or `http://dash.localhost` if you are using the default).

## 7. Log In

You will be presented with the M3TAL dashboard login screen.

*   **Default Credentials**:
    *   Username: `admin`
    *   Password: The password you set during the `m3tal config wizard` (or `admin_pass` if you did not change it).

**To change your dashboard password after logging in, use the following command in your terminal:**

```bash
sudo m3tal dashpass
```

Follow the prompts to set a new password.

---

## Filesystem Contract

The M3TAL ecosystem relies on a specific filesystem layout for configuration and state.

| Path                     | Purpose                                                                                                           |
| :----------------------- | :---------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | The primary environment configuration file. Managed by the `m3tal config wizard` and other `m3tal config` commands. |
| `/var/lib/m3tal/state.db`| SQLite database storing the operational state of M3TAL. Auto-created and managed by the API daemon.             |
| `/docker`                | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory where Docker Compose files reside. |
| `/docker/users.json`     | Stores dashboard credentials. Managed by the `m3tal dashpass` command.                                            |

## Port Map

| Port   | Service           | Access Method                                               |
| :----- | :---------------- | :---------------------------------------------------------- |
| 80     | Traefik           | Public (when `DASHBOARD_EXPOSE_MODE=traefik`)               |
| 8080   | M3TAL API daemon  | Host-local (accessed by other M3TAL components internally) |
| 8081   | Traefik dashboard | Host-local only                                             |
| 8082   | M3TAL Dashboard   | Direct port access (`local` mode) or via Traefik (`traefik` mode) |

## Firewall Note

If Traefik is exposed to the public internet (e.g., `DASHBOARD_EXPOSE_MODE=traefik`), ensure that port 80 is open in your firewall:

```bash
sudo ufw allow 80
```

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
    (Press `Ctrl+C` to exit the log stream)