```markdown
# M3TAL: Get Started

This guide will walk you through the initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before proceeding, ensure you have the following software installed on your system:

*   **Docker Engine**: The containerization platform.
*   **Docker Compose V2**: The orchestration tool for Docker.

You can verify their installation by running the following command:

```bash
docker --version && docker compose version
```

If either command fails, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

M3TAL is distributed as a Debian package. Install it using the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, you must run the configuration wizard to set up M3TAL's core environment.

```bash
sudo m3tal config wizard
```

This command will guide you through a series of prompts:

*   **`DASHBOARD_EXPOSE_MODE`**: Controls how the M3TAL Dashboard is accessed.
    *   `local` (default): Exposes the dashboard directly via a port on your host machine. Suitable for LAN-only access or initial setup.
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy, allowing access via a domain name (e.g., `dash.yourdomain.com`). Requires Traefik to be running.
*   **`DOMAIN`**: The primary domain name for your M3TAL setup. Used for Traefik routing if `DASHBOARD_EXPOSE_MODE` is set to `traefik`. Defaults to `localhost`.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL Dashboard will be accessible. Defaults to `8082`.
*   **`PUID` / `PGID`**: The User ID and Group ID for running Docker containers. Defaults to `1000`. Ensure these match your primary user's IDs for proper file permissions.
*   **`TZ`**: Your local timezone (e.g., `America/New_York`).
*   **`CONFIG_PATH`**: The directory where M3TAL will store its configuration files and state. Defaults to `./data/config`.
*   **`BASE_STORAGE_PATH`**: The root directory for M3TAL's data storage. Defaults to `./data`.
*   **`MEDIA_PATH`**: Path for media storage.
*   **`DOWNLOADS_PATH`**: Path for download storage.
*   **`API_TOKEN`**: A secret token for API authentication. **Crucially, change this from the default.**
*   **`ADMIN_PASSWORD`**: The password for accessing the M3TAL Dashboard. **Crucially, change this from the default.**
*   **`DASHBOARD_SECRET`**: A secret used internally by the dashboard. **Crucially, change this from the default.**

## Step 4: Start the Routing Stack (Traefik)

The M3TAL ecosystem uses Traefik as its primary reverse proxy and gateway. To start the core routing components, run:

```bash
m3tal up
```

This command orchestrates all Docker Compose files found within the `/docker/` directory. This includes Traefik and any other services defined in your stack.

## Step 5: Start the Dashboard

The M3TAL Dashboard provides a user interface for managing your M3TAL services.

```bash
m3tal dash up
```

This command will:
1. Pull the latest `m3tal-dashboard` Docker image if it's not already present.
2. Read your `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3. Start the dashboard container, applying the correct Docker Compose configuration based on your exposure mode.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's address. The exact URL depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Access the dashboard at `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address) or `http://localhost:8082` if accessing from the same machine.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Access the dashboard at `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured). Traefik must be running (`m3tal up`).

## Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: `admin_pass` (or whatever you set for `ADMIN_PASSWORD` during the wizard)

**IMPORTANT:** It is highly recommended to change the default password immediately after your first login. You can do this from the dashboard's settings or by running:

```bash
sudo m3tal dashpass <new_password>
```

Replace `<new_password>` with your desired secure password.

---

## Filesystem Contract

M3TAL manages its configuration and state across specific filesystem locations:

| Path                       | Purpose                                                                      |
| :------------------------- | :--------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | The primary M3TAL environment configuration file. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`  | SQLite database storing M3TAL's operational state. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`        | The canonical directory containing all M3TAL Docker Compose files and Traefik configurations. |
| `/docker`                  | A symbolic link that points to `/opt/m3tal/stack/`. This is the user-facing path for interacting with M3TAL stacks via the CLI. |
| `/docker/users.json`       | Stores the M3TAL Dashboard user credentials. Managed by `m3tal dashpass`.    |

## Port Map

The following ports are used by M3TAL components:

| Port   | Service               | Access                                        |
| :----- | :-------------------- | :-------------------------------------------- |
| 80     | Traefik (HTTP Entry)  | Public (if `DASHBOARD_EXPOSE_MODE=traefik`)   |
| 8080   | M3TAL API Daemon (Go) | Host-local                                    |
| 8081   | Traefik Dashboard     | Host-local only (accessible via `localhost:8081` or `YOUR_IP:8081`) |
| 8082   | M3TAL Dashboard       | Direct port (if `local` mode) or via Traefik (if `traefik` mode) |

## Firewall Note

If Traefik is exposed to the public internet (e.g., `DASHBOARD_EXPOSE_MODE=traefik` and your server is accessible externally), ensure that port 80 (HTTP) is allowed through your firewall. If using `ufw`:

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

---

This concludes the basic setup guide. You are now ready to start adding and managing services within the M3TAL ecosystem.
```