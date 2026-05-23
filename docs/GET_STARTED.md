# M3TAL Getting Started Guide

This guide provides instructions for setting up the M3TAL ecosystem for first-time users.

---

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your Docker installation:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

Execute the configuration wizard to set up essential environment variables.

```bash
sudo m3tal config wizard
```

You will be prompted for several values:

-   **DOMAIN**: Your primary domain for Traefik routing (e.g., `example.com`). If you plan to use `DASHBOARD_EXPOSE_MODE=local`, this can be `localhost`.
-   **DASHBOARD_EXPOSE_MODE**:
    -   `local` (default): The dashboard will be directly exposed on `http://YOUR_IP:8082`. No Traefik configuration is required for dashboard access. Ideal for local or LAN-only setups.
    -   `traefik`: The dashboard will be routed through Traefik and accessible at `http://dash.DOMAIN`. Requires Traefik to be running via `m3tal up`. Ideal for domain-based setups.
-   **DASHBOARD_SECRET**: A secret key used for dashboard session security.
-   **ADMIN_PASSWORD**: The default password for the `admin` user to log into the M3TAL Dashboard.
-   **API_TOKEN**: An authentication token for the M3TAL API.
-   **PUID** / **PGID**: The User ID and Group ID that M3TAL's containers will run as. Typically `1000` for the first non-root user. This ensures correct file permissions for mounted volumes.
-   **TZ**: Your desired timezone (e.g., `America/New_York`).
-   **BASE_STORAGE_PATH**, **CONFIG_PATH**, **MEDIA_PATH**, **DOWNLOADS_PATH**: These define the base directories for M3TAL's data, configuration files, and common media/download storage paths for your containers. They are used as default volume mounts.

These settings are saved to `/etc/m3tal/.env`.

## 4. Start the Routing Stack

This command starts the core routing services, including Traefik.

```bash
m3tal up
```

This command processes all Docker Compose files (`*-compose.yml`) found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). It starts and manages all declared services. The `routing-compose.yml` file, specifically, deploys the Traefik reverse proxy, which binds to host port 80 for HTTP traffic.

## 5. Start the Dashboard

This command initiates the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command specifically manages the `m3tal-dashboard` container. It downloads the necessary `m3tal-compose.yml` files, reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`, and starts the dashboard using the appropriate Docker Compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`). This pulls the required Docker image and starts the container.

## 6. Access the Dashboard

Open your web browser and navigate to the appropriate address:

-   If `DASHBOARD_EXPOSE_MODE` is `local`:
    `http://YOUR_IP:8082` (Replace `YOUR_IP` with the IP address of your server).
-   If `DASHBOARD_EXPOSE_MODE` is `traefik`:
    `http://dash.DOMAIN` (Replace `DOMAIN` with the value configured in Step 3).
    *Note: For `traefik` mode, ensure `m3tal up` has successfully started Traefik and that your DNS resolves `dash.DOMAIN` to your server's IP.*

## 7. Log In

Use the following default credentials to log in:

-   **Username**: `admin`
-   **Password**: The `ADMIN_PASSWORD` you set during the configuration wizard.

To change the dashboard password for the `admin` user, use the following command:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

The M3TAL ecosystem adheres to a specific filesystem layout for its components:

| Path                     | Purpose                                                                |
| :----------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the API daemon.     |
| `/opt/m3tal/stack/`      | Canonical stack directory. Contains Docker Compose files and Traefik config. |
| `/docker`                | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store for users. Managed by `m3tal dashpass`.     |

## Port Map

The following ports are used by M3TAL components on the host system:

| Port | Service                                | Access                                             |
| :--- | :------------------------------------- | :------------------------------------------------- |
| 80   | Traefik HTTP entry point               | Public (when Traefik is running and exposed)       |
| 8080 | M3TAL API daemon (Go)                  | Host-local only                                    |
| 8081 | Traefik dashboard                      | Host-local only (accessible via `http://localhost:8081` on the host where Traefik runs) |
| 8082 | M3TAL Dashboard container              | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

## Firewall Note

If you plan to expose Traefik on port 80 to the internet or your local network, you may need to adjust your firewall rules. For systems using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands:

-   Check the status of the API daemon:
    ```bash
    systemctl status m3tal-api
    ```
-   View real-time logs for the API daemon:
    ```bash
    journalctl -u m3tal-api -f
    ```