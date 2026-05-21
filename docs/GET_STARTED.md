```markdown
# docs/GET_STARTED.md — M3TAL First-Time Setup Guide

This guide provides instructions for installing and initially configuring the M3TAL Ecosystem for first-time users.

---

### Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your Docker installation:

```bash
docker --version && docker compose version
```

### Step 2: Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the `m3tal` CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### Step 3: Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables for M3TAL operations. This wizard creates or updates `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt for several values. Below is an explanation of the key prompts:

*   **`DASHBOARD_EXPOSE_MODE`**:
    *   **`local` (default)**: The M3TAL Dashboard will be directly exposed on port `8082` of your host machine. This is suitable for local network access or if you do not plan to use Traefik for domain-based routing. Access will be via `http://YOUR_IP:8082`.
    *   **`traefik`**: The M3TAL Dashboard will be routed through the Traefik reverse proxy. This requires Traefik to be running and a `DOMAIN` configured. Access will be via `http://dash.YOUR_DOMAIN`.
*   **`DOMAIN`**:
    *   Your primary domain name (e.g., `example.com`). This is used by Traefik for routing services like `api.DOMAIN` and `dash.DOMAIN`. If `DASHBOARD_EXPOSE_MODE` is `traefik`, this will determine the dashboard URL. Default is `localhost`.
*   **`PUID` (User ID)** and **`PGID` (Group ID)**:
    *   These specify the user and group IDs that containers will run as. Use `id -u YOUR_USERNAME` and `id -g YOUR_USERNAME` to find your user's IDs. This ensures file permissions are correctly handled between the host and containers. Default is `1000`.
*   **`DASHBOARD_SECRET`**:
    *   A unique secret key used by the M3TAL Dashboard for session management. Generate a strong, random string (e.g., using `openssl rand -hex 16`). **Change this immediately from its default for security.**
*   **`ADMIN_PASSWORD`**:
    *   The password for the default `admin` user to log into the M3TAL Dashboard. **Change this immediately from its default for security.**
*   **`TZ` (Timezone)**:
    *   Your system's timezone (e.g., `America/New_York`). This ensures consistent time reporting across all services. Default is `America/Denver`.
*   **`BASE_STORAGE_PATH`**:
    *   The base directory for all M3TAL data volumes (e.g., `/mnt/m3tal-data`). This path will be mounted into containers. Subdirectories like `media`, `config`, and `downloads` will be created under this path. Default is `./data`.
*   **`HTTP_PORT`**:
    *   The port on which the internal M3TAL API daemon (a Go binary) listens. It is typically accessed only by other M3TAL components or via Traefik. Default is `8080`.

### Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command initializes and starts all Docker Compose stacks found in the `/docker/` directory. This includes the core routing stack (Traefik), which acts as a reverse proxy.

```bash
m3tal up
```

This command will:
*   Create the `proxy` Docker network if it doesn't exist.
*   Start the `traefik` container as defined in `/docker/routing-compose.yml`.
*   Optionally start the `cloudflared` container if configured in `routing-compose.yml`.

### Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command will:
1.  Download the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` definitions.
2.  Read the `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Pull the `ghcr.io/jakej985-rgb/m3tal-godash:debug` Docker image.
4.  Start the `m3tal-dashboard` container, applying the appropriate override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on `DASHBOARD_EXPOSE_MODE`.

### Step 6: Open Browser

Navigate to the M3TAL Dashboard in your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local`:**
    `http://YOUR_IP:8082` (Replace `YOUR_IP` with the actual IP address of your server or `localhost` if accessing from the server itself).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    `http://dash.YOUR_DOMAIN` (Replace `YOUR_DOMAIN` with the domain you configured in Step 3).

### Step 7: Log In

The default credentials for the M3TAL Dashboard are:
*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (default `admin_pass` if not changed).

**It is strongly recommended to change the admin password immediately.**
To change the dashboard administrator password, use the following command:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

M3TAL establishes a clear contract for critical files and directories:

| Path                        | Purpose                                                                                |
| :-------------------------- | :------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file for the M3TAL ecosystem. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Stores M3TAL's internal state. Auto-created by the API daemon.  |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains all Docker Compose files and Traefik configuration.  |
| `/docker`                   | **Symlink** to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations (e.g., placing new compose files). |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                               |

## Port Table

The following ports are used by M3TAL components:

| Port | Service                               | Access                                                                  |
| :--- | :------------------------------------ | :---------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (if `DASHBOARD_EXPOSE_MODE` is `traefik`, or for other services) |
| 8080 | M3TAL API daemon (Go)                 | Host-local (accessed by other M3TAL components or via Traefik)          |
| 8081 | Traefik dashboard                     | Host-local only (e.g., `http://localhost:8081`)                         |
| 8082 | M3TAL Dashboard (Python/Flask)        | Direct port (if `DASHBOARD_EXPOSE_MODE` is `local`) or via Traefik     |

## Firewall Note

If you are exposing Traefik on port 80 to the public internet or your local network, you may need to open this port in your firewall:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api`.
You can manage this service using standard `systemctl` commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
```