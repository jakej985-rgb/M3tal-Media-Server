# Getting Started with M3TAL

This guide will walk you through the initial setup and configuration of M3TAL for first-time users.

## Step 1: Prerequisites

Before you begin, ensure you have the following software installed on your system:

*   **Docker Engine:** The containerization platform.
*   **Docker Compose V2:** The orchestration tool for Docker.

You can verify your installation by running the following command in your terminal:

```bash
docker --version && docker compose version
```

If either command returns an error or shows an outdated version, please refer to the official Docker documentation for installation and update instructions.

## Step 2: Install M3TAL via APT

M3TAL is distributed via an APT repository for straightforward installation. Execute the following three commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-apt-key/public.key | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-apt-key stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

Once M3TAL is installed, you need to run its configuration wizard to set up essential parameters. This wizard will guide you through each setting.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a brief explanation of each:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL dashboard is accessed.
    *   `local`: (Default) Exposes the dashboard directly via a port on your host machine. Ideal for LAN-only access and initial setup.
    *   `traefik`: Exposes the dashboard through Traefik, allowing access via a domain name (e.g., `dash.yourdomain.com`). Requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. Defaults to `8082`.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. Defaults to `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state, including the database. Defaults to `./state` within the M3TAL configuration path.
*   **`LOG_LEVEL`**: Controls the verbosity of M3TAL logs. Options typically include `debug`, `info`, `warn`, `error`.
*   **`DASHBOARD_SECRET`**: A secret key used for securing the dashboard. **Change this from the default immediately.**
*   **`API_TOKEN`**: An API token for programmatic access to M3TAL. **Change this from the default immediately.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. **Change this from the default immediately.**
*   **`DOMAIN`**: The primary domain name you will use for M3TAL services. Defaults to `localhost`.
*   **`BASE_STORAGE_PATH`**: The base directory for storing data for various M3TAL services. Defaults to `./data`.
*   **`MEDIA_PATH`**: The directory for storing media files.
*   **`CONFIG_PATH`**: The directory where M3TAL will store its configuration files.
*   **`DOWNLOADS_PATH`**: The directory for downloads.
*   **`PUID` / `PGID`**: The User ID (UID) and Group ID (GID) that Docker containers will run as. Typically, these should match your user's UID/GID to avoid permission issues. You can find these by running `id -u` and `id -g`.
*   **`TZ`**: Your local timezone (e.g., `America/Denver`).

The wizard will save your choices to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

The M3TAL routing stack, powered by Traefik, is responsible for managing incoming network traffic to your M3TAL services. To start it, run:

```bash
m3tal up
```

This command orchestrates all Docker Compose files found in the `/docker/` directory, effectively starting all configured services, including Traefik.

## Step 5: Start the Dashboard

Now, start the M3TAL dashboard container. This command will pull the latest dashboard image and start the container based on your `DASHBOARD_EXPOSE_MODE` setting from the configuration wizard.

```bash
m3tal dash up
```

## Step 6: Access the Dashboard

Open your web browser and navigate to the M3TAL dashboard. The access URL depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Access the dashboard at `http://YOUR_IP:8082` or `http://localhost:8082` (replace `YOUR_IP` with your server's IP address).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Access the dashboard at `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured, e.g., `http://dash.localhost`). Traefik must be running for this to work.

## Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Username:** `admin`
*   **Password:** The password you set for `ADMIN_PASSWORD` during the configuration wizard (default is `admin_pass`).

**It is highly recommended to change the default password immediately after logging in.** You can do this using the following command:

```bash
sudo m3tal dashpass <new_password>
```

Replace `<new_password>` with your desired secure password.

## Filesystem Contract

M3TAL utilizes a specific filesystem structure to manage its configuration and data. Understanding these locations is crucial for troubleshooting and advanced configuration:

| Path                         | Purpose                                                                 |
| :--------------------------- | :---------------------------------------------------------------------- |
| `/etc/m3tal/.env`            | Primary environment file containing all M3TAL configuration variables. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`    | SQLite database storing M3TAL's operational state. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`          | Canonical directory containing M3TAL's core Docker Compose files and Traefik configuration. |
| `/docker`                    | A symbolic link that points to `/opt/m3tal/stack/`. This is the user-facing path for all Docker Compose operations. |
| `/docker/users.json`         | Stores dashboard user credentials. Managed by `m3tal dashpass`. |

## Port Map

| Port   | Service           | Access                                          |
| :----- | :---------------- | :---------------------------------------------- |
| 80     | Traefik HTTP      | Public (when `DASHBOARD_EXPOSE_MODE=traefik`)   |
| 8080   | M3TAL API Daemon  | Host-local                                      |
| 8081   | Traefik Dashboard | Host-local only                                 |
| 8082   | M3TAL Dashboard   | Direct port (`local` mode) or via Traefik (`traefik` mode) |

## Firewall Note

If you intend to expose M3TAL services to the public internet (e.g., using Traefik on port 80), you may need to open the relevant ports in your firewall. For example, if Traefik is running on port 80:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard systemd commands:

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