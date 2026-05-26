# M3TAL Ecosystem: Getting Started

This guide provides a complete, step-by-step process for first-time users to set up the M3TAL ecosystem.

---

## 1. Prerequisites

Ensure Docker Engine and Docker Compose V2 are installed on your system.
Verify their installation by running:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Use the following commands to install the M3TAL CLI binary and systemd services:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This installs the `/usr/bin/m3tal` CLI binary and sets up the `m3tal-api.service` systemd daemon.

## 3. Run the Configuration Wizard

The M3TAL configuration wizard will guide you through setting up the essential environment variables in `/etc/m3tal/.env`. These variables control the behavior of the M3TAL API daemon and dashboard.

Run the wizard:

```bash
sudo m3tal config wizard
```

You will be prompted for various settings:

-   **`DASHBOARD_EXPOSE_MODE` (local/traefik)**:
    -   `local`: The M3TAL Dashboard will be directly exposed on `http://YOUR_IP:8082`. This is ideal for quick setups, local access, or when not using a domain.
    -   `traefik`: The M3TAL Dashboard will be routed through Traefik, accessible via `http://dash.YOUR_DOMAIN`. This requires Traefik to be running and a `DOMAIN` configured.
-   **`DASHBOARD_PORT`**: The internal port the dashboard container uses, and the external port if `DASHBOARD_EXPOSE_MODE` is `local`. Default is `8082`.
-   **`DOMAIN`**: Your domain name (e.g., `example.com`). This is used by Traefik for routing services (e.g., `dash.example.com`, `api.example.com`). If `DASHBOARD_EXPOSE_MODE` is `traefik`, ensure this is set correctly.
-   **`PUID` (User ID)**: The User ID for container processes, ensuring correct file permissions. Default is `1000`.
-   **`PGID` (Group ID)**: The Group ID for container processes. Default is `1000`.
-   **`BASE_STORAGE_PATH`**: The base directory for all M3TAL-managed persistent data and volumes (e.g., `/mnt/m3tal`).
-   **`CONFIG_PATH`**: The path for configuration files within `BASE_STORAGE_PATH` (e.g., `/mnt/m3tal/config`).
-   **`MEDIA_PATH`**: The path for media files within `BASE_STORAGE_PATH` (e.g., `/mnt/m3tal/media`).
-   **`DOWNLOADS_PATH`**: The path for download files within `BASE_STORAGE_PATH` (e.g., `/mnt/m3tal/downloads`).
-   **`TZ` (Timezone)**: Your server's timezone (e.g., `America/Denver`).
-   **`DASHBOARD_SECRET`**: A strong, random string used to secure dashboard sessions. **Change this from the default immediately.**
-   **`API_TOKEN`**: A secure token for accessing the M3TAL API. **Change this from the default immediately.**
-   **`ADMIN_PASSWORD`**: The initial password for the M3TAL Dashboard `admin` user. **Change this from the default immediately.**
-   **`TRAEFIK_WEB_PORT`**: The HTTP port Traefik will listen on. Default is `80`.
-   **`TRAEFIK_WEBHTTPS_PORT`**: The HTTPS port Traefik will listen on. Default is `443`.

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This includes the `routing-compose.yml` which deploys Traefik, the reverse proxy for M3TAL services.

```bash
m3tal up
```

This command will start Traefik, making it ready to route traffic to your services.

## 5. Start the M3TAL Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It pulls the necessary Docker image and starts the container with the correct configuration based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

```bash
m3tal dash up
```

This command will:
1.  Download the latest `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Read `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container, applying the appropriate compose override for either `local` or `traefik` exposure.

## 6. Open the M3TAL Dashboard in Your Browser

Access the M3TAL Dashboard based on your configured `DASHBOARD_EXPOSE_MODE`:

-   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open your browser to `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host).
-   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open your browser to `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain configured in the wizard, e.g., `http://dash.example.com`). This requires Traefik to be running via `m3tal up` and proper DNS records for your domain.

## 7. Log In to the Dashboard

The default login credentials for the M3TAL Dashboard are:

-   **Username:** `admin`
-   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (default is `admin_pass` if not changed).

**It is highly recommended to change the default password immediately.** You can do this using the M3TAL CLI:

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new password for the `admin` user. This updates the credentials in `/docker/users.json`.

---

## Filesystem Contract

The M3TAL ecosystem interacts with specific files and directories on your system:

| Path                     | Purpose                                                          |
| :----------------------- | :--------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.    |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the API daemon.|
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                | Symlink → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your `*-compose.yml` files here. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.         |

## Port Table

M3TAL uses the following ports on your host machine:

| Port | Service               | Access                                                 |
| :--- | :-------------------- | :----------------------------------------------------- |
| 80   | Traefik HTTP entry point| Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed via Traefik) |
| 8080 | M3TAL API daemon (Go) | Host-local (used by dashboard and other internal components) |
| 8081 | Traefik Dashboard     | Host-local only (for monitoring Traefik itself)        |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`DASHBOARD_EXPOSE_MODE=traefik`) |

## Firewall Note

If you plan to expose Traefik on port 80 (e.g., if `DASHBOARD_EXPOSE_MODE` is `traefik`), you may need to open this port in your firewall. For systems using UFW (Uncomplicated Firewall):

```bash
sudo ufw allow 80
```
Consider allowing 443 as well if using HTTPS with Traefik.

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard systemctl and journalctl commands:

-   **Check the status of the M3TAL API daemon:**
    ```bash
    systemctl status m3tal-api
    ```

-   **View real-time logs for the M3TAL API daemon:**
    ```bash
    journalctl -u m3tal-api -f
    ```

-   **Restart the M3TAL API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```