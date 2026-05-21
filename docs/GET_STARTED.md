```markdown
# M3TAL Get Started Guide

This guide will walk you through the initial setup of M3TAL on your system.

## Step 1: Prerequisites

Before you begin, ensure the following software is installed on your system:

*   **Docker Engine**: The containerization platform.
*   **Docker Compose V2**: The tool for defining and running multi-container Docker applications.

You can verify their installation by running:

```bash
docker --version && docker compose version
```

If these commands do not return version information, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is installed using your system's package manager. Execute the following commands to add the M3TAL repository and install the `m3tal` package:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The `m3tal config wizard` command will guide you through setting up essential M3TAL configuration options.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a breakdown of what each prompt means:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via its port. This is the default and recommended for initial setup.
    *   `traefik`: Exposes the dashboard through Traefik, allowing access via a domain name. Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL will store its state. Defaults to `./state`.
*   **`LOG_LEVEL`**: The verbosity of M3TAL's logs. Options include `debug`, `info`, `warn`, `error`. `info` is the default.
*   **`DASHBOARD_SECRET`**: A secret key for the dashboard. **Change this from the default `change_me_immediately` for security.**
*   **`API_TOKEN`**: An API token for authentication. **Change this from the default `change_me_api_token` for security.**
*   **`ADMIN_PASSWORD`**: The password for the M3TAL dashboard administrator. **Change this from the default `admin_pass` for security.**
*   **`DOMAIN`**: The domain name used for accessing services when `DASHBOARD_EXPOSE_MODE` is set to `traefik`. Defaults to `localhost`.
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL data storage. Defaults to `./data`.
*   **`PUID` / `PGID`**: The User ID and Group ID for running Docker containers. Defaults to `1000`.
*   **`TZ`**: Your system's timezone (e.g., `America/Denver`).

This wizard creates/updates the `/etc/m3tal/.env` file with your selections.

## Step 4: Start the Routing Stack (Traefik)

M3TAL uses Traefik as its reverse proxy to manage incoming traffic to your services. To start Traefik and other core routing components, run:

```bash
m3tal up
```

This command orchestrates all Docker Compose files found within the `/docker/` directory. This includes Traefik and any other core routing infrastructure.

## Step 5: Start the Dashboard

The M3TAL dashboard provides a web interface for managing your services. To start the dashboard container:

```bash
m3tal dash up
```

This command will pull the latest M3TAL dashboard Docker image and start it as a container, respecting the `DASHBOARD_EXPOSE_MODE` setting from your configuration.

## Step 6: Access the M3TAL Dashboard

You can access the M3TAL dashboard in your web browser. The address depends on your `DASHBOARD_EXPOSE_MODE` configuration:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open your browser and navigate to `http://YOUR_IP:8082`. Replace `YOUR_IP` with your server's IP address or `localhost` if accessing from the same machine.

*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open your browser and navigate to `http://dash.DOMAIN`. Replace `DOMAIN` with the domain you configured during the wizard (e.g., `http://dash.localhost`).

## Step 7: Log In

Once the dashboard is loaded, you will be presented with a login screen.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: `admin_pass` (or whatever you set during the wizard for `ADMIN_PASSWORD`)

**IMPORTANT:** It is highly recommended to change the default password immediately after your first login. You can change the administrator password using the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password for the dashboard administrator.

---

## Filesystem Contract

M3TAL adheres to a specific filesystem structure for configuration and data storage. Understanding these locations is crucial for managing and troubleshooting your M3TAL installation.

| Path                      | Purpose                                                              | Notes                                                      |
| :------------------------ | :------------------------------------------------------------------- | :--------------------------------------------------------- |
| `/etc/m3tal/.env`         | Primary M3TAL configuration file.                                    | Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db` | SQLite database storing M3TAL's operational state.                 | Automatically created by the API daemon.                   |
| `/opt/m3tal/stack/`       | Canonical directory for M3TAL's Docker Compose files and configs.    | Contains Traefik configurations and core stack definitions. |
| `/docker`                 | Symlink pointing to `/opt/m3tal/stack/`.                             | This is the user-facing path for all stack operations.     |
| `/docker/users.json`      | Stores dashboard user credentials.                                   | Managed by `m3tal dashpass`.                               |

---

## Port Map

M3TAL utilizes several ports for its core services and for accessing your deployed applications.

| Port       | Service                 | Access Type          | Notes                                                                 |
| :--------- | :---------------------- | :------------------- | :-------------------------------------------------------------------- |
| 80         | Traefik HTTP Entrypoint | Public (Traefik Mode)| The main entry point for HTTP traffic when using Traefik.             |
| 8080       | M3TAL API Daemon (Go)   | Host-Local           | The backend API that manages M3TAL services.                          |
| 8081       | Traefik Dashboard       | Host-Local Only      | The administrative interface for Traefik itself.                      |
| 8082       | M3TAL Dashboard         | Direct/Traefik       | Accessible directly via port (local mode) or through Traefik (traefik mode). |

---

## Firewall Note

If you have a firewall enabled (e.g., `ufw`), you need to allow incoming traffic on port 80 to make your services accessible from the internet via Traefik.

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard systemd commands.

*   **Check the status of the M3TAL API service:**

    ```bash
    systemctl status m3tal-api
    ```

*   **View live logs for the M3TAL API service:**

    ```bash
    journalctl -u m3tal-api -f
    ```

*   **Restart the M3TAL API service:**

    ```bash
    sudo systemctl restart m3tal-api
    ```
```