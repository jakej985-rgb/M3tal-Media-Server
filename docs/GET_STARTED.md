# M3TAL Getting Started Guide

Welcome to M3TAL! This guide provides a step-by-step process for first-time users to set up and start using the M3TAL ecosystem.

---

### Step 1: Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container management. These must be installed on your system before proceeding.

Verify your Docker installation:

```bash
docker --version && docker compose version
```

Expected output will show the installed versions for Docker Engine and Docker Compose V2.

### Step 2: Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the APT package manager.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### Step 3: Run the Configuration Wizard

Initialize your M3TAL environment by running the configuration wizard. This will generate the primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration values:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard will be directly accessible via a host port (e.g., `http://YOUR_IP:8082`). This is suitable for LAN-only setups or when Traefik is not desired.
    *   `traefik`: The dashboard will be routed through the Traefik reverse proxy (e.g., `http://dash.DOMAIN`). This requires Traefik to be running.
*   **`DOMAIN`**: The base domain for services exposed via Traefik (e.g., `example.com`). If you set `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard will be available at `dash.YOUR_DOMAIN`.
*   **`PUID` / `PGID`**: The User ID (PUID) and Group ID (PGID) that containers will run as. It is recommended to use the PUID/PGID of your non-root user (e.g., `id -u` and `id -g`). This ensures proper file permissions for mounted volumes.
*   **`TZ`**: Your local timezone (e.g., `America/Denver`). This sets the timezone for containers.
*   **`BASE_STORAGE_PATH`**: The base directory for all M3TAL data volumes (e.g., `/mnt/m3tal-data`). This path will contain subdirectories for config, media, and downloads.
*   **`CONFIG_PATH`**: The directory for application configuration files. Defaults to `${BASE_STORAGE_PATH}/config`.
*   **`MEDIA_PATH`**: The directory for media files. Defaults to `${BASE_STORAGE_PATH}/media`.
*   **`DOWNLOADS_PATH`**: The directory for downloads. Defaults to `${BASE_STORAGE_PATH}/downloads`.
*   **`DASHBOARD_SECRET`**: A secret key for the dashboard. **Change this from the default for security.**
*   **`API_TOKEN`**: An authentication token for the M3TAL API. **Change this from the default for security.**
*   **`ADMIN_PASSWORD`**: The default password for the dashboard admin user. **Change this from the default for security.**

### Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined in `/docker/`. This includes the core routing stack (Traefik).

```bash
m3tal up
```

This command will start containers defined in files like `routing-compose.yml` (for Traefik) and any other `*-compose.yml` files present in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`).

### Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It will pull the necessary Docker image and start the container according to your `DASHBOARD_EXPOSE_MODE` setting.

```bash
m3tal dash up
```

This command downloads the latest `m3tal-compose.yml` and applies either `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` as an override, depending on the `DASHBOARD_EXPOSE_MODE` configured in `/etc/m3tal/.env`.

### Step 6: Access the Dashboard

Open your web browser and navigate to the M3TAL Dashboard. The address depends on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local`**:
    `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`**:
    `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain configured in Step 3). Traefik must be running for this mode to work.

### Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Default Username**: `admin`
*   **Default Password**: The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard` (or `admin_pass` if you used the default).

**To change the dashboard password:**
After logging in, you can change the password for the `admin` user using the CLI:

```bash
sudo m3tal dashpass admin new_secure_password
```

---

### Filesystem Contract

M3TAL maintains several important files and directories on your system:

*   `/etc/m3tal/.env`: The primary configuration file for M3TAL. This file stores environment variables used by the API daemon and Docker Compose stacks. It is managed by `m3tal config wizard` and `m3tal config set`.
*   `/var/lib/m3tal/state.db`: The SQLite state database used by the M3TAL API daemon to manage service states, user settings, and other operational data. It is automatically created and managed by the API.
*   `/opt/m3tal/stack/`: The canonical directory where M3TAL stores its core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`).
*   `/docker`: A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for placing your own Docker Compose files (e.g., `my-app-compose.yml`) to integrate them into the M3TAL ecosystem. All `*-compose.yml` files placed here will be managed by `m3tal up`.
*   `/docker/users.json`: The credential store for the M3TAL Dashboard. This file is managed by `m3tal dashpass`.

### Port Table

The following ports are used by M3TAL components:

| Port | Service | Access | Description |
| :--- | :-------------- | :------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| 80   | Traefik HTTP    | Public   | Main HTTP entry point for services routed through Traefik. Exposed if Traefik is running and configured (e.g., `m3tal up`). |
| 8080 | M3TAL API Daemon | Host-local | The Go API daemon listens on this port. It's typically accessed by `host.docker.internal` from inside containers. |
| 8081 | Traefik Dashboard | Host-local | Traefik's own administrative dashboard, available only on the host machine. |
| 8082 | M3TAL Dashboard | Direct / Traefik | The M3TAL Dashboard container. Accessible directly via `http://YOUR_IP:8082` (local mode) or via Traefik at `http://dash.YOUR_DOMAIN` (traefik mode). |

### Firewall Note

If you intend to expose Traefik or the M3TAL Dashboard to your local network or the internet, you may need to adjust your firewall rules. For example, to allow HTTP traffic on port 80 for Traefik:

```bash
sudo ufw allow 80/tcp
```

### Service Management

The M3TAL API daemon runs as a systemd service. You can manage its lifecycle using `systemctl`.

**Check the status of the M3TAL API daemon:**

```bash
systemctl status m3tal-api
```

**View live logs for the M3TAL API daemon:**

```bash
journalctl -u m3tal-api -f
```