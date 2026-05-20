# M3TAL Getting Started Guide

This guide provides a complete, step-by-step setup for first-time users of the M3TAL Ecosystem.

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your Docker installation by running:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the APT package manager:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

The `m3tal` CLI binary is installed at `/usr/bin/m3tal`, and the `m3tal-api.service` systemd daemon starts automatically.

## Step 3: Run the Configuration Wizard

Initialize your M3TAL environment variables using the configuration wizard:

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration values. Here's an explanation of the key prompts:

*   **`DOMAIN`**: This specifies the base domain for services exposed via Traefik (e.g., `dash.yourdomain.com`, `api.yourdomain.com`). If you don't have a domain or are setting up for local use only, `localhost` is a suitable default.
*   **`DASHBOARD_EXPOSE_MODE`**:
    *   `local` (default): The M3TAL Dashboard will be directly exposed on `http://YOUR_IP:8082`. This is ideal for quick setups or local network access.
    *   `traefik`: The M3TAL Dashboard will be exposed via the Traefik reverse proxy at `http://dash.${DOMAIN}`. This requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL Dashboard will be accessible when `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`PUID`** (User ID) / **`PGID`** (Group ID): These define the user and group IDs that containers will run as. Using your non-root user's PUID/PGID ensures proper file permissions for mounted volumes. You can find these with `id -u` and `id -g`.
*   **`BASE_STORAGE_PATH`**: The primary directory where M3TAL-managed applications will store their data. Default is `./data`, which resolves to `/opt/m3tal/stack/data`.
*   **`CONFIG_PATH`**: The directory for configuration files for M3TAL-managed applications. Default is `./data/config`, resolving to `/opt/m3tal/stack/data/config`.
*   **`DASHBOARD_SECRET`**: A secret key used to secure dashboard user sessions. It is recommended to generate a strong, random value for this.
*   **`ADMIN_PASSWORD`**: The default password for the `admin` user on the M3TAL Dashboard. **Change this immediately after setup.**

## Step 4: Start the Routing Stack (Traefik)

The M3TAL routing stack includes Traefik, the reverse proxy gateway that routes traffic to your services.

```bash
m3tal up
```

This command starts all Docker Compose files located in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This typically includes the `routing-compose.yml` that defines Traefik. Traefik will bind to port 80 on the host system (and 8081 for its own dashboard, accessible only from localhost).

## Step 5: Start the Dashboard

Now, start the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest M3TAL Dashboard Compose files.
2.  Reads your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container, applying the correct Docker Compose override (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode) to configure its network exposure.

## Step 6: Open the Browser

Access the M3TAL Dashboard in your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik` and Traefik is running:**
    Open `http://dash.DOMAIN` (replace `DOMAIN` with the value you configured in Step 3).

## Step 7: Log In to the Dashboard

The default login credentials are:
*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` value you set in the configuration wizard (default is `admin_pass`).

**It is critical to change the default password immediately after your first login.**
You can change the admin password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL uses specific locations for its core configuration and data:

| Path                        | Purpose                                                    |
| :-------------------------- | :--------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file, managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite database storing M3TAL's internal state.           |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for placing Docker Compose files for your applications. |
| `/docker/users.json`        | Dashboard user credential store, managed by `m3tal dashpass`. |

## Port Table

The M3TAL ecosystem uses the following ports:

| Port | Service                                | Access                 |
| :--- | :------------------------------------- | :--------------------- |
| 80   | Traefik HTTP entry point               | Public (when exposed)  |
| 8080 | M3TAL API daemon (Go service)          | Host-local only        |
| 8081 | Traefik Dashboard (internal)           | Host-local only        |
| 8082 | M3TAL Dashboard                        | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If your server has a firewall (e.g., `ufw`) enabled, you may need to explicitly allow traffic on the ports M3TAL uses.

*   To allow Traefik to receive public web traffic (when `DASHBOARD_EXPOSE_MODE=traefik` or for other web services):
    ```bash
    sudo ufw allow 80/tcp
    ```
*   To allow direct access to the M3TAL Dashboard (when `DASHBOARD_EXPOSE_MODE=local`):
    ```bash
    sudo ufw allow 8082/tcp
    ```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands:

*   **Check the status of the M3TAL API daemon:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View real-time logs for the M3TAL API daemon:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the M3TAL API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```