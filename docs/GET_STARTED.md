# M3TAL Getting Started Guide

This guide provides operational steps for first-time M3TAL users to set up their environment and launch the core services.

## Step 1: Prerequisites

Before installing M3TAL, ensure **Docker Engine** and **Docker Compose V2** are installed on your system.

You can verify their installation and versions with the following commands:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

M3TAL is distributed as a single Go binary via an APT repository. Follow these commands to install it:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, run the M3TAL configuration wizard to set up essential environment variables:

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration options. Here’s what each prompt means:

*   **`DASHBOARD_PORT`**: The host port to expose the M3TAL Dashboard on if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**:
    *   `local` (default): The dashboard will be directly accessible via `http://YOUR_IP:8082` (or the port you specify). No Traefik configuration is needed for dashboard access.
    *   `traefik`: The dashboard will be accessible via a subdomain, typically `http://dash.YOUR_DOMAIN`, routed by Traefik. Requires `DOMAIN` to be set and Traefik to be running.
*   **`DOMAIN`**: Your primary domain name (e.g., `example.com`). This is crucial if `DASHBOARD_EXPOSE_MODE` is `traefik` or if you plan to use Traefik for other services (e.g., `api.YOUR_DOMAIN`).
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management. **Change this immediately from the default for security.**
*   **`API_TOKEN`**: An authentication token for the M3TAL API. **Change this immediately from the default for security.**
*   **`ADMIN_PASSWORD`**: The default password for the `admin` user on the M3TAL Dashboard. **Change this immediately from the default for security.**
*   **`PUID` / `PGID`**: The User ID and Group ID that containers will run as. Typically, this is your current user's UID/GID (e.g., `id -u` and `id -g`). This ensures proper file permissions for mounted volumes.
*   **`TZ`**: Your timezone (e.g., `America/New_York`). Used by containers.
*   **Storage Paths (`BASE_STORAGE_PATH`, `CONFIG_PATH`, `MEDIA_PATH`, `DOWNLOADS_PATH`)**: These define the base directories on your host filesystem where M3TAL will store application data, configurations, media, and downloads respectively. Defaults typically point to subdirectories under `/mnt` or `/data`.
*   Other variables like `HTTP_PORT`, `LOG_LEVEL`, `NETWORK_NAME`, `LOCAL_IP` define internal operational aspects and networking, often safe to leave at defaults for a basic setup.

The wizard writes these settings to `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack

Start the core routing stack, which includes Traefik (the reverse proxy) and potentially Cloudflared for tunnel capabilities:

```bash
m3tal up
```

This command runs `docker compose` across all `*-compose.yml` files located in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`), bringing up services like Traefik. Traefik will typically listen on port 80 to route incoming HTTP traffic.

## Step 5: Start the M3TAL Dashboard

Launch the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest `m3tal-compose.yml` (base configuration) and relevant override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the base compose file and the appropriate override file (either `local` or `traefik`) to expose the dashboard as configured.

## Step 6: Open the M3TAL Dashboard in Your Browser

Access the M3TAL Dashboard using your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host, or `localhost` if accessing from the same machine).

*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Ensure DNS records are pointing correctly if accessing externally.

## Step 7: Log In to the Dashboard

The default login credentials for the M3TAL Dashboard are:

*   **Username**: `admin`
*   **Password**: The value of `ADMIN_PASSWORD` you set in Step 3 (default `admin_pass` if not changed).

**It is highly recommended to change the default password immediately.** You can do this via the dashboard's user management interface or using the CLI:

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL relies on a specific filesystem structure. Understanding these paths is crucial for configuration and troubleshooting:

| Path                        | Purpose                                                                             |
| :-------------------------- | :---------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file containing environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite database storing M3TAL API state and metadata. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical directory containing all Docker Compose files for system and user services. |
| `/docker`                   | **Symlink** to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your custom `*-compose.yml` files here. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                            |

## Port Table

M3TAL uses the following ports for its core services:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if Traefik is exposed)              |
| 8080 | M3TAL API daemon (Go)       | Host-local only (internal communication)    |
| 8081 | Traefik Dashboard (internal)| Host-local only                             |
| 8082 | M3TAL Dashboard (Python)    | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If Traefik is exposed to the internet or local network, you may need to allow traffic on port 80 through your firewall. For `ufw` (Uncomplicated Firewall), you would use:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard `systemctl` commands:

*   **Check API service status:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View API service logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```