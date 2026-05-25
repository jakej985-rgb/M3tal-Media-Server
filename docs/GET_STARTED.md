# Getting Started with M3TAL

This guide provides step-by-step instructions for installing and setting up the M3TAL ecosystem for the first time.

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 to manage its services. Ensure these are installed and operational on your system before proceeding.

You can verify your Docker installation with the following commands:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Install the M3TAL CLI and API daemon using the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables for M3TAL's operation. These settings are stored in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values:

-   **`DOMAIN`**: The base domain for your services if you plan to use Traefik for domain-based routing (e.g., `example.com`). If you're using M3TAL on a local network without a public domain, `localhost` is a suitable default.
-   **`LOCAL_IP`**: The local IP address of your host machine.
-   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed:
    -   `local` (default): The dashboard will be directly accessible via a host port (default `8082`). No Traefik configuration is needed for dashboard access. Best for local network access or initial setup.
    -   `traefik`: The dashboard will be accessible via Traefik using a domain (e.g., `dash.YOUR_DOMAIN`). Requires Traefik to be running via `m3tal up`.
-   **`DASHBOARD_PORT`**: The port on which the dashboard will be exposed if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
-   **`BASE_STORAGE_PATH`**: The base directory where M3TAL will store application data, volumes, and configurations. Default is `/mnt`.
-   **`CONFIG_PATH`**: The path for M3TAL's core configuration files. Default is `${BASE_STORAGE_PATH}/config`.
-   **`PUID` / `PGID`**: The User ID and Group ID that containers will run as. This helps ensure proper file permissions when containers interact with host volumes. You can typically find these using `id -u` and `id -g` for your current user.
-   **`TZ`**: Your timezone (e.g., `America/Denver`).
-   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management. It's automatically generated if left blank.
-   **`API_TOKEN`**: A token for authentication with the M3TAL API. It's automatically generated if left blank.

## 4. Start the Routing Stack (Traefik)

The M3TAL CLI manages Docker Compose stacks. The `m3tal up` command starts all `*-compose.yml` files found within the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`).

This step will initiate the Traefik reverse proxy and any other default routing components. Traefik typically listens on port 80 for HTTP traffic.

```bash
m3tal up
```

## 5. Start the M3TAL Dashboard

The M3TAL Dashboard is a web-based UI for managing your M3TAL ecosystem. The `m3tal dash up` command specifically pulls the dashboard Docker image and starts its container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the base `m3tal-compose.yml` along with the appropriate override file based on `DASHBOARD_EXPOSE_MODE`.
    -   If `DASHBOARD_EXPOSE_MODE=local`, `m3tal-compose.local.yml` is used, directly binding port `${DASHBOARD_PORT:-8082}:8082` to the host.
    -   If `DASHBOARD_EXPOSE_MODE=traefik`, `m3tal-compose.traefik.yml` is used, adding Traefik labels for routing via `dash.${DOMAIN}`.

## 6. Access the Dashboard

Open your web browser and navigate to the appropriate address based on your `DASHBOARD_EXPOSE_MODE` setting:

-   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    ```
    http://YOUR_IP:8082
    ```
    Replace `YOUR_IP` with the IP address of your M3TAL host.

-   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    ```
    http://dash.DOMAIN
    ```
    Replace `DOMAIN` with the domain you configured during the wizard (e.g., `http://dash.example.com`). Ensure Traefik (started via `m3tal up`) is running and correctly configured for your domain.

## 7. Log In to the Dashboard

The default login credentials for the dashboard are:

-   **Username:** `admin`
-   **Password:** `admin_pass`

You can change the admin password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```
This command updates the password stored in `/docker/users.json`.

---

## Filesystem Contract

M3TAL interacts with specific directories and files on your system:

| Path                        | Purpose                                                                 |
| :-------------------------- | :---------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.           |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.      |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's core compose files and Traefik config.  |
| `/docker`                   | Symlink that points to `/opt/m3tal/stack/`. User-facing path for stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.                |

## Port Table

The following ports are used by M3TAL components:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if Traefik is exposed)              |
| 8080 | M3TAL API daemon (Go)       | Host-local only (accessed by containers)    |
| 8081 | Traefik dashboard           | Host-local only                             |
| 8082 | M3TAL Dashboard container | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If you expose Traefik on port 80 (e.g., for `DASHBOARD_EXPOSE_MODE=traefik` or other services), you may need to open the port in your firewall:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon (`m3tal-api`) runs as a systemd service. You can manage its state and view its logs using standard systemd commands:

-   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
-   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
-   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```