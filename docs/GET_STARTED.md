# M3TAL Get Started Guide

This guide provides a complete, step-by-step process for first-time users to install, configure, and operate the M3TAL ecosystem.

---

### Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your Docker installation by running:

```bash
docker --version && docker compose version
```

Example output:
```
Docker version 24.0.7, build afdd53b
Docker Compose version v2.23.0-desktop.1
```

If Docker Engine or Docker Compose V2 are not installed, please refer to the official Docker documentation for your operating system.

### Step 2: Install M3TAL via APT

M3TAL is distributed via an APT repository. Execute the following commands to add the repository and install the `m3tal` package:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

The installation process will automatically start the M3TAL API daemon (`m3tal-api.service`).

### Step 3: Run the Configuration Wizard

After installation, run the M3TAL configuration wizard to set up essential environment variables. This wizard creates or updates the primary configuration file at `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values:

-   **`DOMAIN`**: The base domain name for your services (e.g., `example.com`). This is used by Traefik for routing. If you don't have a domain or are setting up for local use, `localhost` is the default and a suitable choice.
-   **`DASHBOARD_EXPOSE_MODE`**:
    -   `local` (default): The dashboard will be accessible directly via a host port binding (e.g., `http://YOUR_IP:8082`). This mode does not require Traefik for dashboard access and is suitable for LAN-only or initial setups.
    -   `traefik`: The dashboard will be exposed via the Traefik reverse proxy, typically at `http://dash.YOUR_DOMAIN`. This mode requires Traefik to be running via `m3tal up`.
-   **`DASHBOARD_PORT`**: The direct port to expose the dashboard on when `DASHBOARD_EXPOSE_MODE` is set to `local`. Defaults to `8082`.
-   **`PUID` / `PGID`**: The User ID and Group ID that containers will run as. These typically default to `1000` for the first non-root user on Linux and ensure proper file permissions for mounted volumes. You can find your current IDs with `id -u` and `id -g`.
-   **`TZ`**: The timezone for your containers (e.g., `America/Denver`).
-   **`BASE_STORAGE_PATH`**: The base directory on your host system where all persistent data for M3TAL-managed services will be stored (e.g., `/mnt`). This directory must exist.
-   **`CONFIG_PATH`**: A sub-path within `BASE_STORAGE_PATH` specifically for configuration files. This is where `users.json` for the dashboard is stored (e.g., `/mnt/config`).
-   **`DASHBOARD_SECRET`**: A secret key used for securing the dashboard's session. **Change this from the default for production environments.**
-   **`API_TOKEN`**: A token used to authenticate with the M3TAL API. **Change this from the default for production environments.**
-   **`ADMIN_PASSWORD`**: The default password for the `admin` user to log into the M3TAL Dashboard. **Change this from the default immediately after setup.**

Review and confirm the settings.

### Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command initializes and starts all Docker Compose stacks found in the `/docker/` directory. This includes the core routing stack (`routing-compose.yml`), which deploys Traefik and optionally Cloudflared.

```bash
m3tal up
```

This command executes `docker compose` operations across all `*-compose.yml` files in `/docker/`, ensuring all defined services (including Traefik, if configured) are brought up.

### Step 5: Start the M3TAL Dashboard

The `m3tal dash up` command specifically manages the M3TAL dashboard container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` files.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using `m3tal-compose.yml` as the base, applying either `m3tal-compose.local.yml` (for `local` mode) or `m3tal-compose.traefik.yml` (for `traefik` mode) as an override. This pulls the necessary Docker image (`ghcr.io/jakej985-rgb/m3tal-godash:debug`) if not already present.

### Step 6: Open the Browser

Access the M3TAL Dashboard using your web browser:

-   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the actual IP address of your host machine, which can be found using `ip a`).
    Alternatively, if accessing from the host machine, use `http://localhost:8082`.

-   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Traefik must be running (`m3tal up`) for this mode to work.

### Step 7: Log In

The default credentials for the M3TAL Dashboard are:
-   **Username:** `admin`
-   **Password:** `admin_pass` (or the `ADMIN_PASSWORD` you set in Step 3)

**Change the default password immediately:**
Use the `m3tal dashpass` command to change the `admin` user's password. This command modifies the `/docker/users.json` file.

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new, strong password.

---

### Filesystem Contract

The M3TAL ecosystem maintains a strict filesystem contract for its operational files:

| Path                        | Purpose                                                                   |
| :-------------------------- | :------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.             |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created by the M3TAL API daemon.     |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains Docker Compose files and Traefik config. |
| `/docker`                   | **Symlink** to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your custom `*-compose.yml` files here. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                  |

### Port Table

M3TAL and its core components utilize the following network ports on the host system:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services use it) |
| 8080 | M3TAL API daemon (Go)       | Host-local (internal communication)         |
| 8081 | Traefik Dashboard (if enabled) | Host-local only                             |
| 8082 | M3TAL Dashboard             | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

### Firewall Note

If Traefik is exposed publicly on port 80 (e.g., for `traefik` dashboard mode or other services), you may need to open this port in your host's firewall. For `ufw` (Uncomplicated Firewall):

```bash
sudo ufw allow 80/tcp
```

### Service Management

The M3TAL API daemon runs as a systemd service. You can manage its state and view its logs using standard `systemctl` and `journalctl` commands:

-   **Check status of the M3TAL API daemon:**
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