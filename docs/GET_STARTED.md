# GET_STARTED.md

This guide provides a comprehensive, step-by-step process for setting up M3TAL for the first time.

---

## Step 1: Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container orchestration. These must be installed on your system before proceeding.

Verify your Docker installation:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

M3TAL is distributed as a single Go binary via an APT repository. Execute the following commands to install it:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, use the configuration wizard to set up essential environment variables. This wizard populates the `/etc/m3tal/.env` file.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several details:

*   **`DASHBOARD_EXPOSE_MODE`**:
    *   **`local` (default)**: The M3TAL Dashboard will be directly exposed on a specific host port (default 8082). This is ideal for quick setups, LAN-only access, or when Traefik is not desired for the dashboard.
    *   **`traefik`**: The M3TAL Dashboard will be exposed via the Traefik reverse proxy, accessible by a domain (e.g., `dash.YOUR_DOMAIN`). This mode requires Traefik to be running.
*   **`DOMAIN`**: The base domain for your services (e.g., `example.com`). If `DASHBOARD_EXPOSE_MODE` is `traefik`, the dashboard will be available at `dash.YOUR_DOMAIN`. If `local` is chosen, this defaults to `localhost`.
*   **`DASHBOARD_PORT`**: The port on the host where the M3TAL Dashboard will be directly accessible if `DASHBOARD_EXPOSE_MODE` is `local`. Defaults to `8082`.
*   **`BASE_STORAGE_PATH`**: The base directory on your host where M3TAL will store user data, configuration files, media, and downloads. This path is volume-mounted into containers. Defaults to `./data` relative to the M3TAL installation.
*   **`PUID` (User ID)** and **`PGID` (Group ID)**: These are the User ID and Group ID that containers will use. Set these to match the user ID and group ID of the user that owns your `BASE_STORAGE_PATH` to prevent permission issues. You can usually find these with `id -u` and `id -g` for your current user. Defaults to `1000`.
*   **`TZ` (Timezone)**: Your preferred timezone (e.g., `America/Denver`). This ensures logs and scheduled tasks reflect your local time.
*   **`ADMIN_PASSWORD`**: The password for the `admin` user to log into the M3TAL Dashboard. **Change this from the default `admin_pass` immediately.**
*   **`DASHBOARD_SECRET`**: A random secret key used for session management in the M3TAL Dashboard. The wizard will generate a strong default.
*   **`API_TOKEN`**: An API authentication token for programmatic access to the M3TAL API. The wizard will generate a strong default.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command initializes and starts all Docker Compose stacks found in the `/docker/` directory. This includes the core routing stack, which deploys Traefik, the reverse proxy gateway.

```bash
m3tal up
```

This command will bring up the Traefik container (and optionally Cloudflared), making port 80 available for routing domain-based services.

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  It downloads the latest M3TAL Dashboard Docker Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from the M3TAL GitHub repository.
2.  It reads the `DASHBOARD_EXPOSE_MODE` variable from your `/etc/m3tal/.env` file.
3.  Based on `DASHBOARD_EXPOSE_MODE`, it starts the `m3tal-dashboard` container, applying the correct Docker Compose override file (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode). This determines how the dashboard is exposed.

## Step 6: Open Browser

Access the M3TAL Dashboard in your web browser based on your chosen `DASHBOARD_EXPOSE_MODE`:

*   **If `DASHBOARD_EXPOSE_MODE=local`:**
    *   Open `http://YOUR_SERVER_IP:8082` (or `http://localhost:8082` if accessing from the server itself). Replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host.
*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    *   Open `http://dash.YOUR_DOMAIN`. Replace `YOUR_DOMAIN` with the domain you configured in Step 3. Ensure your DNS records are pointing `dash.YOUR_DOMAIN` to your server's IP address.

## Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Username:** `admin`
*   **Default Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (default is `admin_pass` if you didn't change it).

**Important:** It is strongly recommended to change the admin password immediately after your first login.
You can change the dashboard password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

The M3TAL ecosystem maintains a clear contract with the host filesystem. Understanding these paths is crucial for configuration and troubleshooting.

| Path                           | Purpose                                                                                                                                                                                                                                 |
| :----------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`              | **Primary configuration file.** Contains all environment variables for the M3TAL API and Docker Compose stacks. This file is managed by `m3tal config wizard` and `m3tal config set`.                                                  |
| `/var/lib/m3tal/state.db`      | **SQLite state database.** This database is automatically created and managed by the M3TAL API daemon. It stores internal state, service information, and other operational data.                                                       |
| `/opt/m3tal/stack/`            | **Canonical stack directory.** This directory contains core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`) and Traefik dynamic configuration.                                                                |
| `/docker`                      | **Symlink to `/opt/m3tal/stack/`.** This is the user-facing path where you will place your own Docker Compose files (e.g., `my-service-compose.yml`) for `m3tal up` to discover and manage.                                           |
| `/docker/users.json`           | **Dashboard credential store.** This file contains the hashed credentials for the M3TAL Dashboard users. It is managed by the `m3tal dashpass` command.                                                                               |

---

## Port Map

The following ports are used by M3TAL components:

| Port | Service               | Access Level                                                                                                               |
| :--- | :-------------------- | :------------------------------------------------------------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed via Traefik)                                    |
| 8080 | M3TAL API daemon (Go) | Host-local only (internal communication, used by Dashboard, exposed externally via Traefik if `api.${DOMAIN}` is configured) |
| 8081 | Traefik dashboard     | Host-local only (accessed via `http://localhost:8081` on the server)                                                       |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`dash.${DOMAIN}`)                                           |

---

## Firewall Configuration

If you have a firewall enabled (e.g., `ufw`), you may need to open ports for external access. If Traefik is exposed on port 80 for public access, allow it:

```bash
sudo ufw allow 80
```

Adjust firewall rules based on your `DASHBOARD_EXPOSE_MODE` and any other services you expose. For `DASHBOARD_EXPOSE_MODE=local`, you might need to open `8082` (e.g., `sudo ufw allow 8082`).

---

## Service Management

The M3TAL API daemon runs as a `systemd` service named `m3tal-api.service`. You can manage its state and view its logs using standard `systemctl` and `journalctl` commands:

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **Stop the API daemon:**
    ```bash
    sudo systemctl stop m3tal-api
    ```