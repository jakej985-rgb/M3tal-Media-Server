```markdown
# Getting Started with M3TAL

This guide provides a complete, step-by-step setup for first-time M3TAL users.

---

## Step 1: Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2. Ensure both are installed and operational on your system.

To verify your installation, run:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

M3TAL is distributed as a Go binary through an APT repository. Follow these commands to install it:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, configure M3TAL using the interactive wizard. This sets up your primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several key settings:

*   **`DOMAIN`**: The base domain for your services (e.g., `example.com`). This will be used by Traefik for routing services like `api.example.com` or `dash.example.com`. Default: `localhost`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard is directly accessible on a specific port (default 8082) on your host IP. Best for local network access or initial setup.
    *   `traefik`: The dashboard is routed through Traefik, making it accessible via a domain like `dash.YOUR_DOMAIN`. Requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The port for the M3TAL Dashboard if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default: `8082`.
*   **`ADMIN_PASSWORD`**: The default administrator password for the M3TAL Dashboard. **Change this from the default `admin_pass` immediately.**
*   **`DASHBOARD_SECRET`**: A secret key used by the M3TAL Dashboard for session management. Default: `change_me_immediately`.
*   **`API_TOKEN`**: A token used to authenticate with the M3TAL API. Default: `change_me_api_token`.
*   **`PUID`** (User ID) and **`PGID`** (Group ID): These IDs determine the permissions for containers and volumes, ensuring they match your system's user/group for proper access. Default: `1000`.
*   **`TZ`**: Your timezone (e.g., `America/New_York`). Default: `America/Denver`.
*   **`BASE_STORAGE_PATH`**: The base directory for all M3TAL-related persistent storage (e.g., `/mnt/m3tal-data`). This is where your service data, configurations, and downloads will reside. Default: `./data`.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory. This includes the core routing stack managed by Traefik, which acts as a reverse proxy for your services.

```bash
m3tal up
```

This command will:
*   Start the `traefik` container, exposing port 80 for HTTP traffic.
*   Start the `cloudflared` container (if configured via `routing-compose.yml`), enabling Cloudflare tunnels.
*   Initialize the `proxy` Docker network, which all M3TAL-managed containers use for inter-service communication.

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command will:
1.  Download the latest M3TAL Dashboard Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Read the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container using the appropriate Compose override file (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode) to expose it correctly.
4.  Pull the `ghcr.io/jakej985-rgb/m3tal-godash:debug` Docker image if not already present.

The dashboard container communicates with the M3TAL API daemon (which runs on the host at `http://host.docker.internal:8080`) to manage your M3TAL ecosystem.

## Step 6: Open the M3TAL Dashboard in Your Browser

Access the M3TAL Dashboard based on your `DASHBOARD_EXPOSE_MODE` setting from Step 3:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open your browser to `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your M3TAL host). If accessing from the host itself, use `http://localhost:8082`.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open your browser to `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Traefik must be running (`m3tal up`).

## Step 7: Log In to the Dashboard

When prompted, log in with the following credentials:

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (default: `admin_pass`).

**Important:** You can change the dashboard administrator password at any time using the `m3tal dashpass` command:

```bash
sudo m3tal dashpass
```

This will update the `/docker/users.json` file where dashboard credentials are stored.

---

## Filesystem Contract

M3TAL relies on specific file and directory locations for its operation:

| Path                        | Purpose                                                                                |
| :-------------------------- | :------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the M3TAL API daemon.               |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's core Docker Compose files and Traefik configuration.   |
| `/docker`                   | **Symlink** to `/opt/m3tal/stack/`. This is the user-facing path for placing custom Docker Compose files for your applications. All `*-compose.yml` files here are managed by `m3tal up`. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                               |

## Port Table

These are the default ports used by M3TAL's core components:

| Port | Service                                      | Access                                                    |
| :--- | :------------------------------------------- | :-------------------------------------------------------- |
| 80   | Traefik HTTP entry point                     | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed) |
| 8080 | M3TAL API daemon (Go service)                | Host-local only (accessed by dashboard and Traefik via `host.docker.internal`) |
| 8081 | Traefik dashboard (admin UI)                 | Host-local only (e.g., `http://localhost:8081`)           |
| 8082 | M3TAL Dashboard (Python/Flask container)     | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

## Firewall Note

If Traefik is exposed publicly (e.g., if `DASHBOARD_EXPOSE_MODE=traefik` or you plan to host services on port 80), you must allow traffic on port 80 through your firewall.

For systems using `ufw`:

```bash
sudo ufw allow 80/tcp
sudo ufw enable # if not already enabled
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage its state using standard `systemctl` commands:

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

---
```