```markdown
# M3TAL Ecosystem: Getting Started Guide

This guide provides the necessary steps to set up and start using the M3TAL ecosystem for the first time.

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
You can verify their installation by running:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the following commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This will install the `m3tal` CLI binary to `/usr/bin/m3tal` and the `m3tal-api.service` systemd daemon.

## Step 3: Run the Configuration Wizard

Initialize your M3TAL environment and set essential configuration values using the interactive wizard:

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several details. Here's an explanation of key prompts:

*   **`DOMAIN`**: The base domain name for your services (e.g., `example.com`). If you plan to access services via `dash.example.com` or `api.example.com`, set this accordingly. If you only plan for local IP access, `localhost` is a suitable default.
*   **`DASHBOARD_EXPOSE_MODE`**:
    *   **`local` (default)**: The M3TAL Dashboard container will expose port `8082` directly on your host machine. This is ideal for quick setup or when you only need to access the dashboard via `http://YOUR_IP:8082` on your local network. No Traefik configuration is required for dashboard access in this mode.
    *   **`traefik`**: The M3TAL Dashboard container will be configured with Docker labels, allowing the Traefik reverse proxy to route traffic to it. Access will be via `http://dash.YOUR_DOMAIN`. This mode requires Traefik to be running via `m3tal up` and a configured `DOMAIN`.
*   **`DASHBOARD_PORT`**: The port the M3TAL dashboard will listen on. Default is `8082`. If `DASHBOARD_EXPOSE_MODE` is `local`, this port will be bound directly to your host.
*   **`PUID` (User ID) / `PGID` (Group ID)**: The User ID and Group ID that Docker containers will run as inside the M3TAL stack. It is generally recommended to set this to your current user's PUID/PGID (e.g., `id -u` and `id -g`). This ensures that files created by containers have correct permissions on your host filesystem.
*   **`TZ` (Timezone)**: Your desired timezone (e.g., `America/New_York`). This ensures containers operate with the correct time.

The wizard will also automatically generate secure values for `DASHBOARD_SECRET` and `API_TOKEN`, and set default `ADMIN_PASSWORD` (which you should change immediately). All these settings are stored in `/etc/m3tal/.env`.

## Step 4: Start the Routing Stack (Traefik)

The M3TAL ecosystem uses Docker Compose to manage its services. The `m3tal up` command starts all Docker Compose stacks defined in `/docker/`. This typically includes the `routing` stack (Traefik, Cloudflared) and any other user-defined stacks.

```bash
m3tal up
```

This command will pull the necessary Docker images (e.g., `traefik:latest`) and start the containers, including the Traefik reverse proxy which listens on port 80 (HTTP).

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container:

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest M3TAL Dashboard compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container, applying the appropriate compose override file based on the `DASHBOARD_EXPOSE_MODE` setting.
    *   If `DASHBOARD_EXPOSE_MODE` is `local`, the dashboard will be directly accessible via `http://HOST_IP:8082`.
    *   If `DASHBOARD_EXPOSE_MODE` is `traefik`, the dashboard will be accessible via Traefik at `http://dash.DOMAIN`.

## Step 6: Open Browser

Once the dashboard container is running, open your web browser to access it:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the actual IP address of your M3TAL host).
    You can also use `http://localhost:8082` if accessing from the host machine itself.

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured in the wizard). This requires Traefik to be running (`m3tal up`) and DNS for `dash.DOMAIN` to point to your M3TAL host.

## Step 7: Log In

The default credentials for the M3TAL Dashboard are:

*   **Username:** `admin`
*   **Password:** `admin_pass` (this is configured during the wizard and stored in `/etc/m3tal/.env` as `ADMIN_PASSWORD`)

**It is critical to change this default password immediately.** You can change the dashboard password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password for the `admin` user. The credentials are stored in `/docker/users.json`.

---

## Filesystem Contract

The following paths are critical to the M3TAL ecosystem:

| Path                        | Purpose                                                                                                                                                                                             |
| :-------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for the M3TAL ecosystem. Contains environment variables used by the API daemon and Docker Compose stacks. Managed by `m3tal config wizard` and `m3tal config set`.         |
| `/var/lib/m3tal/state.db`   | SQLite state database. Stores internal M3TAL API daemon state, managed automatically by the `m3tal-api.service`.                                                                                      |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains the base Docker Compose files for M3TAL's internal services (e.g., routing, dashboard).                                                                           |
| `/docker`                   | Symlink that points to `/opt/m3tal/stack/`. This is the user-facing path where all Docker Compose files (`*-compose.yml`) for the M3TAL control plane and user-defined stacks should reside.          |
| `/docker/users.json`        | Dashboard credential store. Contains hashed passwords for M3TAL Dashboard users. Managed by `m3tal dashpass`.                                                                                         |

## Port Table

The following ports are used by core M3TAL services:

| Port | Service                                | Access                                                               |
| :--- | :------------------------------------- | :------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point               | Publicly exposed if Traefik is running and configured for external access (e.g., for `dash.DOMAIN`). |
| 8080 | M3TAL API daemon (Go)                  | Host-local. Accessed by internal services (e.g., M3TAL Dashboard) and Traefik (if routing `api.DOMAIN`). |
| 8081 | Traefik Dashboard (admin interface)    | Host-local only. Access via `http://YOUR_IP:8081` (default bind) when Traefik is running. |
| 8082 | M3TAL Dashboard (Python/Flask)         | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`http://dash.DOMAIN` if `DASHBOARD_EXPOSE_MODE=traefik`). |

## Firewall Note

If you intend to expose Traefik or the M3TAL Dashboard (in `local` mode) to your local network or the internet, you may need to adjust your firewall rules. For example, to allow HTTP traffic on port 80 for Traefik:

```bash
sudo ufw allow 80/tcp
# If using HTTPS/443, you would also allow 443/tcp
# sudo ufw allow 443/tcp
```

## Service Management

The M3TAL API daemon runs as a systemd service called `m3tal-api.service`. You can manage this service using standard `systemctl` commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```