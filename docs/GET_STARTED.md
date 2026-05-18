# GET_STARTED.md

This guide provides step-by-step instructions for setting up M3TAL on your system for the first time.

---

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your installation by running:

```bash
docker --version && docker compose version
```

Example output:
```
Docker version 24.0.7, build afdd53b
Docker Compose version v2.24.1
```

## 2. Install M3TAL via APT

Use the following commands to add the M3TAL APT repository and install the CLI binary and API daemon.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This installs the `m3tal` CLI binary (`/usr/bin/m3tal`) and the `m3tal-api.service` systemd daemon.

## 3. Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables. This creates the `/etc/m3tal/.env` file.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration values. Here's an explanation of key prompts:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard will be directly accessible via a host port (e.g., `http://YOUR_IP:8082`). No Traefik required for dashboard access.
    *   `traefik`: The dashboard will be accessible via a domain (e.g., `http://dash.YOUR_DOMAIN`) through the Traefik reverse proxy. Traefik must be running.
*   **`DASHBOARD_PORT`**: The port on the host machine where the dashboard will be directly exposed if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`DOMAIN`**: Your primary domain for services exposed via Traefik. This is used for `dash.DOMAIN`, `api.DOMAIN`, etc., if Traefik is used. For local setups, `localhost` is a valid default.
*   **`PUID`** (Process User ID) / **`PGID`** (Process Group ID): The user and group IDs that Docker containers will run as. Use `id -u` and `id -g` to find your current user's IDs. This prevents permission issues with mounted volumes.
*   **`TZ`**: Your local timezone (e.g., `America/New_York`).
*   **`BASE_STORAGE_PATH`**: The base directory for all M3TAL data and persistent storage (e.g., `/mnt/m3tal` or `/var/lib/m3tal/data`).
*   **`CONFIG_PATH`**: The base directory for M3TAL configuration files, typically nested under `BASE_STORAGE_PATH` (e.g., `${BASE_STORAGE_PATH}/config`).
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management.
*   **`API_TOKEN`**: An authentication token for accessing the M3TAL API.
*   **`ADMIN_PASSWORD`**: The default password for the dashboard `admin` user.

It is recommended to use the defaults for most values during first setup, except for `PUID`, `PGID`, `TZ`, `DOMAIN`, and the security-related secrets (`DASHBOARD_SECRET`, `API_TOKEN`, `ADMIN_PASSWORD`) which should be changed from their defaults.

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This includes the `routing-compose.yml` which deploys Traefik.

```bash
m3tal up
```

This command will download the necessary Docker images and start the `traefik` and optionally `cloudflared` containers.

## 5. Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest dashboard compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from your `/etc/m3tal/.env` file.
3.  Starts the `m3tal-dashboard` container using the appropriate compose override file (`m3tal-compose.local.yml` for `local` mode, or `m3tal-compose.traefik.yml` for `traefik` mode).

## 6. Access the Dashboard

How you access the dashboard depends on the `DASHBOARD_EXPOSE_MODE` you configured in Step 3.

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    The dashboard is directly exposed on a host port.
    Open your web browser and navigate to:
    `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address)
    or `http://localhost:8082` if accessing from the server itself.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    The dashboard is routed via Traefik. Traefik must be running (`m3tal up`).
    Open your web browser and navigate to:
    `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in the wizard).

## 7. Log In

The default credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the wizard, or `admin_pass` if you used the default.

It is strongly recommended to change the admin password immediately using the `m3tal dashpass` command:

```bash
sudo m3tal dashpass
```

## Filesystem Contract

M3TAL establishes specific locations for its core components and data. Understanding these paths is crucial for maintenance and troubleshooting.

| Path                     | Purpose                                                      |
| :----------------------- | :----------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Contains environment variables. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the API daemon. Stores M3TAL's internal state. |
| `/opt/m3tal/stack/`      | Canonical directory for M3TAL's Docker Compose files and Traefik dynamic configuration. |
| `/docker`                | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for placing custom Docker Compose stacks and interacting with `m3tal up`. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.     |

## Port Table

M3TAL uses specific ports for its components. Ensure these ports are not in use by other services if you intend to expose them.

| Port | Service                    | Access Mode                                                |
| :--- | :------------------------- | :--------------------------------------------------------- |
| 80   | Traefik HTTP entry point   | Public (if Traefik is configured for external access)      |
| 8080 | M3TAL API daemon (Go)      | Host-local. Accessed by dashboard/other containers via `http://host.docker.internal:8080`. |
| 8081 | Traefik dashboard          | Host-local only. Access via `http://localhost:8081`.       |
| 8082 | M3TAL Dashboard (container)| Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`dash.DOMAIN` if `DASHBOARD_EXPOSE_MODE=traefik`). |

## Firewall Note

If you plan to expose Traefik or the M3TAL Dashboard to your network, you may need to open the relevant ports in your system's firewall. For example, to allow HTTP traffic through Traefik (port 80):

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service (`m3tal-api.service`). You can manage it using `systemctl`.

*   **Check API service status:**
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