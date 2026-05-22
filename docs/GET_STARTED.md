# M3TAL - Get Started

This guide will walk you through the initial setup of the M3TAL system.

## Step 1: Prerequisites

Before proceeding, ensure that you have the following software installed on your system:

*   **Docker Engine**: The containerization platform.
*   **Docker Compose V2**: The tool for defining and running multi-container Docker applications.

You can verify your installation by running the following command in your terminal:

```bash
docker --version && docker compose version
```

If these commands do not output version information, please refer to the official Docker documentation for installation instructions for your operating system.

## Step 2: Install M3TAL via APT

M3TAL is installed using your system's Advanced Packaging Tool (APT). Execute the following three commands in sequence:

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

You will be presented with a series of prompts. Here's what each prompt generally asks for:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a port binding. This is the default and recommended for initial setup or LAN-only access.
    *   `traefik`: Exposes the dashboard through Traefik, allowing access via a domain name (e.g., `dash.yourdomain.com`). Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL will store its persistent state, including the database. The default is typically `./state` relative to the configuration.
*   **`LOG_LEVEL`**: Sets the verbosity of M3TAL's logs. Options include `debug`, `info`, `warn`, `error`. `info` is the default.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session security. **Change this from the default `change_me_immediately`**.
*   **`API_TOKEN`**: A token for API authentication. **Change this from the default `change_me_api_token`**.
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. **Change this from the default `admin_pass`**.
*   **`DOMAIN`**: The domain name used for accessing services if `DASHBOARD_EXPOSE_MODE` is set to `traefik`. The default is `localhost`.
*   **`PUID` / `PGID`**: The User ID and Group ID to run Docker containers under. Often set to your current user's IDs (e.g., `1000`).
*   **`TZ`**: Your local timezone (e.g., `America/Denver`).

Confirm your selections to save the configuration to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

M3TAL uses Docker Compose to manage its services. The `m3tal up` command starts all defined Docker Compose files located in the `/docker/` directory. This includes Traefik, which acts as a reverse proxy for your services.

```bash
m3tal up
```

This command will pull necessary Docker images and start the containers defined in the compose files in `/docker/`.

## Step 5: Start the Dashboard

To launch the M3TAL dashboard, execute the following command:

```bash
m3tal dash up
```

This command will:
1.  Download the necessary Docker image for the M3TAL dashboard.
2.  Consult your `/etc/m3tal/.env` file for the `DASHBOARD_EXPOSE_MODE` setting.
3.  Start the dashboard container, applying the correct Docker Compose override based on the exposure mode.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is `local` (default): `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or `localhost` if running on your local machine).
*   If `DASHBOARD_EXPOSE_MODE` is `traefik`: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured during the wizard, e.g., `http://dash.localhost`).

## Step 7: Log In

You will be presented with a login screen. Use the following default credentials:

*   **Username:** `admin`
*   **Password:** The password you set for `ADMIN_PASSWORD` during the configuration wizard (default: `admin_pass`).

**It is strongly recommended to change the default password immediately after logging in.**

To change the dashboard password outside of the wizard, you can use:

```bash
sudo m3tal dashpass
```

This will prompt you for a new password.

---

## Filesystem Contract

M3TAL relies on a specific filesystem structure for its configuration and state:

| Path                      | Purpose                                                                                   |
| :------------------------ | :---------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | The primary M3TAL configuration file. Managed by the `m3tal config wizard` and `m3tal config set` commands. |
| `/var/lib/m3tal/state.db` | The SQLite database where M3TAL stores its operational state. Automatically created by the API daemon. |
| `/opt/m3tal/stack/`       | The canonical directory containing M3TAL's core Docker Compose files and Traefik configuration. |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for managing all M3TAL stacks and compose files. |
| `/docker/users.json`      | Stores dashboard user credentials, managed by `m3tal dashpass`.                           |

## Port Map

M3TAL utilizes the following ports:

| Port   | Service           | Access Type         | Description                                                                             |
| :----- | :---------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| 80     | Traefik           | Public              | HTTP entry point for external access when Traefik is configured for domain-based routing. |
| 8080   | M3TAL API Daemon  | Host-local          | The core M3TAL API, accessible only from the host machine.                              |
| 8081   | Traefik Dashboard | Host-local only     | The Traefik administrative interface.                                                   |
| 8082   | M3TAL Dashboard   | Direct (local mode) or via Traefik (traefik mode) | The M3TAL web user interface.                                             |

## Firewall Note

If you are exposing Traefik to the public internet (e.g., for domain-based access), you will need to allow traffic on port 80 (and potentially 443 for HTTPS) through your firewall. For example, if you are using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```