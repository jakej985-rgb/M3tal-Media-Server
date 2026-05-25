```markdown
# M3TAL: Getting Started

This guide will walk you through the installation and initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before you begin, ensure you have the following software installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify your installation by running the following commands in your terminal:

```bash
docker --version
docker compose version
```

If these commands do not return version information, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is installed using the Advanced Packaging Tool (APT) from a trusted repository. Execute the following three commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up essential parameters. Run the following command:

```bash
sudo m3tal config wizard
```

You will be prompted for the following information:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL dashboard is accessed.
    *   `local`: Access the dashboard directly via its IP address and port (e.g., `http://YOUR_IP:8082`). This is the default and recommended for initial setup.
    *   `traefik`: Expose the dashboard through the Traefik reverse proxy, typically via a subdomain (e.g., `http://dash.DOMAIN`). Requires Traefik to be running.
*   **`DOMAIN`**: The primary domain name for your M3TAL setup if you are using `traefik` mode for exposing services. Defaults to `localhost`.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. Defaults to `8082`.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. Defaults to `8080`.
*   **`PUID` / `PGID`**: The User ID and Group ID to run Docker containers. Defaults to `1000`. These should match your primary user's IDs to ensure proper file permissions. You can find your PUID/PGID by running `id -u` and `id -g` respectively.
*   **`TZ`**: Your local timezone (e.g., `America/Denver`).
*   **`BASE_STORAGE_PATH`**: The base directory for M3TAL to store its data and configuration. Defaults to `./data`.
*   **`CONFIG_PATH`**: The subdirectory within `BASE_STORAGE_PATH` for configuration files. Defaults to `./data/config`.
*   **`MEDIA_PATH`**: The subdirectory within `BASE_STORAGE_PATH` for media files. Defaults to `./data/media`.
*   **`DOWNLOADS_PATH`**: The subdirectory within `BASE_STORAGE_PATH` for download files. Defaults to `./data/downloads`.
*   **`DASHBOARD_SECRET`**: A secret key for the M3TAL dashboard. **Change this from the default `change_me_immediately` for security.**
*   **`API_TOKEN`**: A token for API authentication. **Change this from the default `change_me_api_token` for security.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL dashboard. Defaults to `admin_pass`. **Change this for security.**

## Step 4: Start the Routing Stack (Traefik)

M3TAL utilizes Traefik as its reverse proxy and routing gateway. To start the Traefik stack, run:

```bash
m3tal up
```

This command orchestrates all Docker Compose files found in the `/docker/` directory, including the Traefik configuration, starting the necessary containers for your M3TAL environment.

## Step 5: Start the Dashboard

The M3TAL dashboard provides a user interface to manage your M3TAL services. Start the dashboard container with the following command:

```bash
m3tal dash up
```

This command will pull the latest M3TAL dashboard Docker image if it's not already present and then start the dashboard container based on your `DASHBOARD_EXPOSE_MODE` configuration.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard's URL:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default), access it at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address) or `http://localhost:8082` if you are accessing it from the same machine.
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, and Traefik is running, access it at: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured).

## Step 7: Log In

You will be presented with the M3TAL dashboard login page.

*   **Default Username**: `admin`
*   **Default Password**: The password you set or the default `admin_pass` if you haven't changed it during the wizard.

**Important:** For security, it is highly recommended to change the default password immediately. You can do this via the command line:

```bash
sudo m3tal dashpass <new_password>
```

Replace `<new_password>` with your desired secure password.

## Filesystem Contract

M3TAL adheres to a specific filesystem structure for configuration and data persistence. Understanding these paths is crucial for managing your M3TAL instance:

| Path                     | Purpose                                                                 | Managed By                    |
| :----------------------- | :---------------------------------------------------------------------- | :---------------------------- |
| `/etc/m3tal/.env`        | Primary M3TAL configuration file. Stores environment variables.         | `m3tal config wizard`, `m3tal config set` |
| `/var/lib/m3tal/state.db`| SQLite database storing M3TAL's internal state and service information. | M3TAL API Daemon (`m3tal-api`) |
| `/docker`                | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all M3TAL Docker Compose stacks. | M3TAL CLI, User                   |
| `/docker/users.json`     | Stores M3TAL dashboard user credentials.                                | `m3tal dashpass`              |

## Port Map

The following ports are used by M3TAL services:

| Port | Service         | Access Method                             |
| :--- | :-------------- | :---------------------------------------- |
| 80   | Traefik HTTP    | Public (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API       | Host-local                                |
| 8082 | M3TAL Dashboard | Direct port (`local` mode) or Traefik (`traefik` mode) |

## Firewall Configuration

If Traefik is exposed to the public internet (e.g., `DASHBOARD_EXPOSE_MODE=traefik` and your server is publicly accessible), ensure that port 80 is open in your firewall. If you are using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
    (Press `Ctrl+C` to exit log view.)
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```

You have now successfully set up M3TAL. You can begin adding your desired applications and services through the dashboard or by placing their Docker Compose files in the `/docker/` directory and running `m3tal up`.
```