# GETTING STARTED with M3TAL

This guide provides operational steps for first-time users to set up and start the M3TAL Ecosystem.

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your installation:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the `m3tal` CLI:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

The wizard guides you through essential settings for your M3TAL environment.

```bash
sudo m3tal config wizard
```

You will be prompted for several configuration values. Here's what some key prompts mean:

*   **`PUID` (User ID) and `PGID` (Group ID)**: These specify the user and group IDs that Docker containers will run as inside your host system. It's recommended to use the PUID and PGID of your unprivileged user (e.g., `id -u` and `id -g` for your current user) to ensure correct file permissions.
*   **`TZ` (Time Zone)**: Set your local timezone (e.g., `America/Denver`). This ensures containers log and schedule operations with correct time information.
*   **`DOMAIN`**: This is your primary domain name (e.g., `example.com`). If you plan to use Traefik for domain-based routing, set this to your actual domain. If you are only accessing locally by IP, `localhost` is a suitable default.
*   **`DASHBOARD_PORT`**: The direct port on your host machine to access the M3TAL Dashboard if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**:
    *   `local`: The dashboard will be directly accessible via `http://YOUR_IP:DASHBOARD_PORT`. This is suitable for local network access or if you don't use Traefik.
    *   `traefik`: The dashboard will be routed through Traefik, accessible via `http://dash.YOUR_DOMAIN`. This requires Traefik to be running.
*   **`BASE_STORAGE_PATH`**: The base directory on your host for all M3TAL-managed data (e.g., `/mnt/data`).
*   **`DASHBOARD_SECRET` and `API_TOKEN`**: Generate strong, unique values for these. They are used for internal security between components.
*   **`ADMIN_PASSWORD`**: This sets the default password for the `admin` user on the M3TAL Dashboard. **Change this from the default immediately.**

The wizard saves these settings to `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

This command starts all Docker Compose files found in the `/docker/` directory, which include the core routing components like Traefik.

```bash
m3tal up
```

This will pull the necessary Docker images (like `traefik:latest`) and start the containers defined in the `routing-compose.yml` and any other `*-compose.yml` files present in `/docker/`. Traefik typically exposes port 80 for HTTP traffic.

## 5. Start the Dashboard

This command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This process will:
1.  Download the latest M3TAL Dashboard Docker Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Read the `DASHBOARD_EXPOSE_MODE` from your `/etc/m3tal/.env` configuration.
3.  Pull the `ghcr.io/jakej985-rgb/m3tal-godash:debug` image.
4.  Start the `m3tal-dashboard` container, applying the appropriate override (`.local.yml` or `.traefik.yml`) based on your `DASHBOARD_EXPOSE_MODE` setting.

## 6. Open your Browser

Access the M3TAL Dashboard based on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local`**:
    Open your browser to `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`**:
    Open your browser to `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in the wizard). Traefik must be running for this to work.

## 7. Log in

The default login credentials for the M3TAL Dashboard are:

*   **Username**: `admin`
*   **Password**: The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard` (default is `admin_pass`).

**It is strongly recommended to change the default password immediately.** You can do this via the dashboard interface or using the CLI:

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new password for the `admin` user. This updates the `/docker/users.json` file.

---

## Filesystem Contract

M3TAL relies on specific file locations for its operation:

| Path                       | Purpose                                                                            |
| :------------------------- | :--------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary configuration file. Contains environment variables managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`  | SQLite database storing M3TAL's internal state. Auto-created by the API daemon.    |
| `/opt/m3tal/stack/`        | Canonical directory for M3TAL's core Docker Compose files and Traefik configuration. |
| `/docker`                  | Symlink that points to `/opt/m3tal/stack/`. This is the user-facing path for placing custom `*-compose.yml` files. |
| `/docker/users.json`       | Stores dashboard user credentials. Managed by `m3tal dashpass`.                     |

## Port Table

These are the default ports used by M3TAL components:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services exposed) |
| 8080 | M3TAL API daemon (Go)       | Host-local (internal communication)         |
| 8081 | Traefik Dashboard           | Host-local only (internal management interface) |
| 8082 | M3TAL Dashboard (Python/Flask) | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If you plan to expose Traefik or the M3TAL Dashboard to the internet or your local network, you may need to open the relevant ports in your system's firewall. For example, to allow HTTP traffic through Traefik:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service. You can manage it using standard `systemctl` commands:

*   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```
*   **View real-time logs**:
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service**:
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **Stop the service**:
    ```bash
    sudo systemctl stop m3tal-api
    ```
*   **Start the service**:
    ```bash
    sudo systemctl start m3tal-api
    ```