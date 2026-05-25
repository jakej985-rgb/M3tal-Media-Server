# M3TAL - Get Started

This guide will walk you through the initial setup of the M3TAL ecosystem.

## Step 1: Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify your installations by running:

```bash
docker --version && docker compose version
```

If Docker or Docker Compose is not installed, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

Use the following commands to add the M3TAL APT repository and install the M3TAL CLI binary.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, run the configuration wizard to set up essential M3TAL parameters.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a breakdown of each prompt:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via its port (e.g., `http://YOUR_IP:8082`). This is the default and recommended for initial setup or LAN-only access.
    *   `traefik`: Exposes the dashboard through Traefik, allowing access via a domain name (e.g., `http://dash.yourdomain.com`). Requires Traefik to be running.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state database. The default is `./state` (relative to where commands are run, but typically managed by the systemd service).
*   **`LOG_LEVEL`**: The verbosity of M3TAL logs. Options include `debug`, `info`, `warn`, `error`. `info` is the default.
*   **`DASHBOARD_SECRET`**: A secret key for the dashboard. Change this from the default for security.
*   **`API_TOKEN`**: A token for API authentication. Change this from the default for security.
*   **`ADMIN_PASSWORD`**: The password for the dashboard's admin user. Change this from the default for security.
*   **`NETWORK_NAME`**: The Docker network name that M3TAL services will join. The default is `m3tal`.
*   **`LOCAL_IP`**: The local IP address of your host machine. `127.0.0.1` is the default.
*   **`DOMAIN`**: The domain name you intend to use for M3TAL services if `DASHBOARD_EXPOSE_MODE` is set to `traefik`. `localhost` is the default.
*   **`BASE_STORAGE_PATH`**: The base directory for persistent storage of M3TAL data. The default is `./data`.
*   **`MEDIA_PATH`**: A sub-directory within `BASE_STORAGE_PATH` for media files. The default is `./data/media`.
*   **`CONFIG_PATH`**: A sub-directory within `BASE_STORAGE_PATH` for configuration files. The default is `./data/config`.
*   **`DOWNLOADS_PATH`**: A sub-directory within `BASE_STORAGE_PATH` for downloaded files. The default is `./data/downloads`.
*   **`PUID`**: The User ID for running Docker containers. Defaults to `1000`.
*   **`PGID`**: The Group ID for running Docker containers. Defaults to `1000`.
*   **`TZ`**: Your system's timezone. `America/Denver` is the default.
*   **`TRAEFIK_WEB_PORT`**: The host port for Traefik's HTTP entrypoint. The default is `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The host port for Traefik's HTTPS entrypoint. The default is `443`.
*   **`TRAEFIK_DASHBOARD_PORT`**: The port Traefik itself listens on for its dashboard. The default is `8080`.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command orchestrates all M3TAL-related Docker Compose files located within the `/docker/` directory. This typically includes the Traefik reverse proxy.

```bash
m3tal up
```

This command will:
*   Locate all `*.yml` files in `/docker/`.
*   Use Docker Compose to bring up the defined services for each file.
*   This will start Traefik, which acts as your primary ingress point for M3TAL services.

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL dashboard container.

```bash
m3tal dash up
```

This command will:
*   Download the necessary M3TAL dashboard Docker image if it's not already present.
*   Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
*   Start the `m3tal-dashboard` container, applying the appropriate configuration override for either local port exposure or Traefik routing.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (the default):
    `http://YOUR_IP:8082`
    (Replace `YOUR_IP` with your server's IP address or `localhost` if running locally).

*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`:
    `http://dash.DOMAIN`
    (Replace `DOMAIN` with the domain you configured, e.g., `http://dash.localhost`).

## Step 7: Log In

When you first access the dashboard, you will be presented with a login screen.

*   **Default Credentials:**
    *   **Username:** `admin`
    *   **Password:** `admin_pass` (This is the default value for `ADMIN_PASSWORD` from the configuration wizard. **It is strongly recommended to change this immediately.**)

To change the administrator password, use the following command:

```bash
sudo m3tal dashpass <new_password>
```

Replace `<new_password>` with your desired secure password.

---

## Filesystem Contract

M3TAL relies on a specific filesystem layout for its configuration and data.

| Path                 | Purpose                                                                                                     |
| :------------------- | :---------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`    | The primary M3TAL environment configuration file. This file is managed by the `m3tal config wizard` and CLI. |
| `/var/lib/m3tal/state.db` | The SQLite database used by the M3TAL API daemon to store its operational state. Auto-created by the API. |
| `/docker`            | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory for all Docker Compose files. |
| `/opt/m3tal/stack/`  | The canonical directory containing all M3TAL Docker Compose files, Traefik configurations, and related assets. |
| `/docker/users.json` | Stores dashboard credentials. Managed via the `m3tal dashpass` command.                                     |

---

## Port Table

| Port | Service        | Access Type                                         |
| :--- | :------------- | :-------------------------------------------------- |
| 80   | Traefik        | Public ingress for HTTP (when `traefik` mode is used). |
| 8080 | M3TAL API      | Host-local. Accessed by internal M3TAL services.    |
| 8081 | Traefik Dashboard | Host-local only (e.g., `http://127.0.0.1:8081`).     |
| 8082 | M3TAL Dashboard | Accessible directly (local mode) or via Traefik.    |

---

## Firewall Note

If you intend to expose Traefik (and thus your M3TAL services) to the internet, ensure that port 80 is allowed through your firewall. If you are using `ufw`:

```bash
sudo ufw allow 80
```

---

## Service Management (API Daemon)

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`.

*   **Check the status of the API service:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View live logs of the API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```