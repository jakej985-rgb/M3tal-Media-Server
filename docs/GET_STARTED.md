# Getting Started with M3TAL

This guide will walk you through the initial setup and configuration of the M3TAL ecosystem.

---

### Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your Docker installation with:

```bash
docker --version && docker compose version
```

---

### Step 2: Install M3TAL via APT

Execute the following commands to install the M3TAL CLI binary and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

### Step 3: Run the Configuration Wizard

After installation, run the configuration wizard to set up your environment variables. These variables are stored in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's what some key prompts mean:

*   **`DOMAIN`**: This is the base domain name for services exposed via Traefik. For example, if you set it to `example.com`, your dashboard would be accessible at `dash.example.com` (if Traefik mode is enabled). The default is `localhost`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard will be directly bound to a port on your host, typically `8082`. Access via `http://YOUR_IP:8082`. This mode does not require Traefik for dashboard access.
    *   `traefik`: The dashboard will be routed through the Traefik reverse proxy, typically accessible at `http://dash.YOUR_DOMAIN`. This mode requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The internal port for the M3TAL Dashboard. Default is `8082`. Only directly exposed to the host when `DASHBOARD_EXPOSE_MODE` is `local`.
*   **`PUID` / `PGID`**: The User ID (PUID) and Group ID (PGID) used by Docker containers for file permissions. Defaults to `1000` (often the first user on Linux systems). It's recommended to match these to your current user's IDs for correct file permissions.
*   **`TZ`**: Your local timezone (e.g., `America/Denver`). This ensures logs and timestamps are correct.
*   **`BASE_STORAGE_PATH`**: The base directory on your host where M3TAL-related data and configuration will be stored. Defaults to `./data` (relative to the context where M3TAL components expect it, typically `/opt/m3tal/stack`).
*   **`CONFIG_PATH`**: The path for application configuration files, typically nested within `BASE_STORAGE_PATH`.

---

### Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined in `*.yml` files located within the `/docker/` directory. This primarily includes the `routing-compose.yml` which deploys Traefik, the reverse proxy, and optionally Cloudflared.

```bash
m3tal up
```

This command will deploy and start the core routing components.

---

### Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It retrieves the necessary Docker Compose files, reads your `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`, and starts the dashboard with the appropriate configuration.

```bash
m3tal dash up
```

This command will pull the dashboard Docker image (if not present) and start the container.

---

### Step 6: Open the Browser

Once the dashboard is running, you can access it via your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open your browser to `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open your browser to `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in the wizard).

---

### Step 7: Log In to the Dashboard

The default credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** `admin_pass`

You can change the dashboard password using the CLI:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

M3TAL relies on specific filesystem paths for its operation and configuration:

| Path | Purpose |
| :----------------------- | :--------------------------------------------- |
| `/etc/m3tal/.env` | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/` | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker` | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json` | Dashboard credential store. Managed by `m3tal dashpass`. |

---

## Port Map

These are the default ports used by M3TAL components:

| Port | Service | Access |
| :--- | :-------------------------- | :------------------------------------------ |
| 80 | Traefik HTTP entry point | Public (if Traefik is exposed) |
| 8080 | M3TAL API daemon (Go) | Host-local |
| 8081 | Traefik dashboard | Host-local only (accessed via `127.0.0.1:8081` on the host) |
| 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) |

---

## Firewall Note

If you are exposing Traefik on port 80 to the public internet (e.g., for domain-based routing), you might need to allow traffic through your firewall. For `ufw` (Uncomplicated Firewall):

```bash
sudo ufw allow 80
```

---

## Service Management

The core M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands:

*   **Check API daemon status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View API daemon logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```