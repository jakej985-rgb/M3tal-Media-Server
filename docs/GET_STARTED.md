# M3TAL - Get Started Guide

This guide provides a complete, step-by-step setup for first-time M3TAL users.

---

### Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify their installation by running:

```bash
docker --version && docker compose version
```

Example expected output:
```
Docker version 24.0.5, build 24.0.5-0ubuntu1~22.04.1
Docker Compose version v2.20.2
```

### Step 2: Install M3TAL via APT

Execute the following commands in your terminal to install the M3TAL CLI binary and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### Step 3: Run the Configuration Wizard

After installation, configure M3TAL using the interactive wizard. This sets up the primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration values. Explanation for common prompts:

*   **`DASHBOARD_EXPOSE_MODE`**:
    *   **`local` (default)**: The M3TAL dashboard will be directly accessible via `http://YOUR_IP:8082`. No Traefik configuration is needed for dashboard access. Ideal for local network use or initial setup.
    *   **`traefik`**: The M3TAL dashboard will be routed through the Traefik reverse proxy, typically accessible at `http://dash.YOUR_DOMAIN`. This requires Traefik to be running and a domain configured.
*   **`DOMAIN`**: The base domain for services exposed via Traefik (e.g., `example.com`). If `DASHBOARD_EXPOSE_MODE` is `traefik`, the dashboard will be at `dash.DOMAIN`. Default is `localhost`.
*   **`DASHBOARD_PORT`**: The internal port for the M3TAL dashboard. Default is `8082`. If `DASHBOARD_EXPOSE_MODE` is `local`, this port will be directly exposed on your host.
*   **`PUID`** (User ID) / **`PGID`** (Group ID): These define the user and group IDs that containers will run as, crucial for file permissions on mounted volumes. Use `id -u YOUR_USERNAME` and `id -g YOUR_USERNAME` to find your user's PUID/PGID. Default is `1000`.
*   **`TZ`** (Timezone): Your system's timezone (e.g., `America/Denver`).
*   **`DASHBOARD_SECRET`**: A secret key used to secure dashboard sessions. It is highly recommended to change the default value.
*   **`API_TOKEN`**: A token used to authenticate requests to the M3TAL API. It is highly recommended to change the default value.
*   **`ADMIN_PASSWORD`**: The initial password for the default `admin` user in the M3TAL dashboard. This should be changed to a strong password.
*   **`BASE_STORAGE_PATH`**, **`MEDIA_PATH`**, **`CONFIG_PATH`**, **`DOWNLOADS_PATH`**: These define the base directories for persistent storage used by various containers. The wizard will suggest sensible defaults; adjust if you have specific storage requirements (e.g., dedicated `/mnt/data`).

### Step 4: Start the Routing Stack (Traefik)

Initiate the core routing and network infrastructure. This command processes all Docker Compose files (`*-compose.yml`) located in the `/docker/` directory, including Traefik.

```bash
m3tal up
```

This will pull necessary images and start containers defined in `/docker/routing-compose.yml` (which includes Traefik and optionally Cloudflared) and any other stack compose files present.

### Step 5: Start the M3TAL Dashboard

Start the M3TAL web dashboard container. This command handles downloading the latest dashboard compose files and launching the dashboard service according to your `DASHBOARD_EXPOSE_MODE` setting from Step 3.

```bash
m3tal dash up
```

This command will:
1.  Download the latest `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Read the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container, applying the correct override for local or Traefik exposure.

### Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the M3TAL Dashboard using the appropriate URL based on your `DASHBOARD_EXPOSE_MODE` configuration:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host). If accessing from the same machine, `http://localhost:8082` also works.
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain configured in `m3tal config wizard`). This requires Traefik to be running via `m3tal up` and proper DNS resolution for `dash.YOUR_DOMAIN`.

### Step 7: Log In to the Dashboard

The default login credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard`, or `admin_pass` if not changed.

**Change the Dashboard Password:**
It is highly recommended to change the default password immediately. Use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password for the `admin` user.

---

### Filesystem Contract

The following paths are critical to the M3TAL ecosystem:

| Path                        | Purpose                                                                                |
| :-------------------------- | :------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for the M3TAL system. Managed by `m3tal config wizard`.     |
| `/var/lib/m3tal/state.db`   | SQLite database storing M3TAL's internal state. Auto-created by the API daemon.        |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's internal Docker Compose files and Traefik configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations (e.g., adding custom `*-compose.yml` files). |
| `/docker/users.json`        | Credential store for the M3TAL Dashboard. Managed by `m3tal dashpass`.                 |

---

### Port Overview

M3TAL utilizes the following network ports:

| Port | Service                               | Access                                                                  |
| :--- | :------------------------------------ | :---------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services exposed) |
| 8080 | M3TAL API daemon (Go)                 | Host-local (accessed by dashboard container and Traefik)                |
| 8081 | Traefik dashboard (for Traefik itself)| Host-local only (not for M3TAL dashboard)                               |
| 8082 | M3TAL Dashboard                       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

---

### Firewall Configuration

If you plan to expose Traefik on port 80 to the public internet or your local network, ensure your firewall allows incoming connections on this port.

For systems using `ufw`:

```bash
sudo ufw allow 80/tcp
```

### M3TAL API Service Management

The M3TAL API daemon runs as a systemd service. Use standard `systemctl` commands to manage it:

*   **Check the API service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View API service logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the API service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```