# GETTING STARTED with M3TAL

This guide provides step-by-step instructions for installing and setting up the M3TAL Ecosystem for first-time users.

---

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container orchestration. Ensure both are installed on your system before proceeding.

You can verify your installation by running:

```bash
docker --version && docker compose version
```

If these commands return errors or indicate older versions, refer to the official Docker documentation for installation and upgrade instructions.

## 2. Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.list] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This will install the `/usr/bin/m3tal` CLI binary and the `m3tal-api.service` systemd daemon.

## 3. Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables:

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's an explanation for each:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard will be directly accessible via a host port, typically `http://YOUR_IP:8082`. No Traefik configuration is required for dashboard access. This is suitable for LAN-only setups or initial testing.
    *   `traefik`: The dashboard will be routed through the Traefik reverse proxy, accessible via a domain, e.g., `http://dash.YOUR_DOMAIN`. This requires Traefik to be running and a `DOMAIN` configured.
*   **`DOMAIN`**: Your primary domain for services exposed via Traefik (e.g., `example.com`). If you chose `DASHBOARD_EXPOSE_MODE=traefik`, the dashboard will be available at `dash.DOMAIN`.
*   **`PUID`** / **`PGID`**: The User ID (PUID) and Group ID (PGID) that containers will run as. This ensures proper file permissions for mounted volumes. You can find your current user's PUID and PGID with `id -u` and `id -g`.
*   **`BASE_STORAGE_PATH`**: The base directory for all your M3TAL related data (e.g., `/mnt/data`). Subdirectories like `media`, `config`, and `downloads` will be created within this path.
*   **`CONFIG_PATH`**: Path for M3TAL-specific configuration files (e.g., `/mnt/config`).
*   **`MEDIA_PATH`**: Path for media storage (e.g., `/mnt/data/media`).
*   **`DOWNLOADS_PATH`**: Path for downloads (e.g., `/mnt/data/downloads`).
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management. Generate a strong, random value for this.
*   **`API_TOKEN`**: An authentication token for accessing the M3TAL API. Generate a strong, random value.
*   **`ADMIN_PASSWORD`**: The default password for the M3TAL Dashboard `admin` user. **Change this immediately after setup.**
*   **`TZ`**: Your system's timezone (e.g., `America/New_York`). This ensures containers have correct time settings.
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL and other services will use (default: `m3tal`).

These settings are saved to `/etc/m3tal/.env`.

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory. This includes the `routing-compose.yml` which deploys Traefik, the reverse proxy.

```bash
m3tal up
```

This command will:
*   Ensure the `m3tal-api` daemon is running (it manages Docker resources).
*   Execute `docker compose` operations using all `*-compose.yml` files present in `/docker/`.
*   Start the Traefik gateway, which typically binds to host port 80 to handle incoming HTTP requests and route them to your services.

## 5. Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command will:
1.  Download the latest M3TAL Dashboard compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Read the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container, applying the appropriate override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on your chosen expose mode.
    *   If `DASHBOARD_EXPOSE_MODE=local`, a direct port binding (`8082:8082` by default) is added.
    *   If `DASHBOARD_EXPOSE_MODE=traefik`, Traefik labels are added to route `dash.DOMAIN` to the dashboard.

## 6. Open Browser and Access Dashboard

Once the dashboard container is running, access it via your web browser:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with your server's actual IP address).
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in the wizard). Traefik must be running for this to work.

## 7. Log In to the Dashboard

The default login credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard.

**It is strongly recommended to change the default admin password immediately.**
You can change the dashboard password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

M3TAL uses specific locations for its core components and user-managed stacks. Understanding these paths is crucial for maintenance and deployment.

| Path                     | Purpose                                                                                                                                                                          |
| :----------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Stores environment variables used by M3TAL and its Docker containers. Managed by `m3tal config wizard` and `m3tal config set`.                      |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the M3TAL API daemon. Stores internal application state and metadata.                                                           |
| `/opt/m3tal/stack/`      | Canonical directory for M3TAL's internal Docker Compose files and Traefik configuration.                                                                                         |
| `/docker`                | A symbolic link to `/opt/m3tal/stack/`. This is the user-facing directory for placing all Docker Compose files (`*-compose.yml`) for your services. `m3tal up` scans this directory.|
| `/docker/users.json`     | Dashboard credential store. Contains hashed usernames and passwords for dashboard access. Managed by `m3tal dashpass`.                                                             |

## Port Table

The following ports are used by M3TAL components:

| Port | Service                    | Access                                                         |
| :--- | :------------------------- | :------------------------------------------------------------- |
| 80   | Traefik HTTP entry point   | Publicly exposed if Traefik is running and configured.         |
| 8080 | M3TAL API daemon (Go)      | Host-local only. Accessed by other M3TAL components.           |
| 8081 | Traefik dashboard          | Host-local only. Provides access to Traefik's internal dashboard (often `http://localhost:8081`). |
| 8082 | M3TAL Dashboard (Python)   | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`). |

## Firewall Note

If you expose Traefik on port 80 to the internet or your local network, you may need to open this port in your system's firewall. For `ufw` (Uncomplicated Firewall), you would typically run:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service called `m3tal-api.service`. You can manage its state using standard systemctl commands:

*   **Check Status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View Logs (live tail):**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```

The API daemon is responsible for managing Docker containers, interacting with the state database, and serving the M3TAL API routes.