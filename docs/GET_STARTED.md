# GET_STARTED.md

Welcome to M3TAL! This guide will walk you through the initial setup, from installation to accessing your M3TAL Dashboard.

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 to manage your services. These must be installed on your system before proceeding.

Verify your Docker installation:

```bash
docker --version && docker compose version
```

Expected output will show versions for Docker Engine and Docker Compose.

## 2. Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using our APT repository.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables. This creates `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's a breakdown of common prompts:

*   **`DOMAIN`**: The base domain for your services (e.g., `example.com`). If you plan to access services via domain names like `dash.example.com`, set this. Defaults to `localhost`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): Dashboard is directly exposed on a host port (`8082`). Ideal for LAN-only or local access.
    *   `traefik`: Dashboard is exposed via Traefik on `dash.${DOMAIN}`. Requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The host port for direct access to the dashboard if `DASHBOARD_EXPOSE_MODE` is `local`. Default is `8082`.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard. **Change from default immediately.**
*   **`API_TOKEN`**: An API token for secure communication with the M3TAL API. **Change from default immediately.**
*   **`PUID` / `PGID`**: User and group IDs for containers, ensuring proper file permissions. Typically `1000`.
*   **`TZ`**: Your timezone (e.g., `America/Denver`).
*   **`BASE_STORAGE_PATH`**: The base directory for all M3TAL data (configs, media, downloads). Defaults to `./data` relative to `/opt/m3tal/stack/`.

## 4. Start the Routing Stack (Traefik)

The M3TAL ecosystem uses Traefik as its default reverse proxy for routing services via domain names.

```bash
m3tal up
```

This command scans the `/docker/` directory for all `*-compose.yml` files and starts them using Docker Compose. The `routing-compose.yml` file, which defines Traefik, is included in this process.

## 5. Start the Dashboard

Now, deploy the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command:
1.  Downloads the latest dashboard compose files from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container, applying the correct Docker Compose override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on your chosen `DASHBOARD_EXPOSE_MODE`. It will also pull the dashboard Docker image if not already present.

## 6. Open the Dashboard in Your Browser

Access the M3TAL Dashboard based on your `DASHBOARD_EXPOSE_MODE` configuration:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open your browser to `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server).

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open your browser to `http://dash.DOMAIN` (replace `DOMAIN` with the value you set in step 3). Traefik must be running for this to work.

## 7. Log In to the Dashboard

The first time you access the dashboard, you will be prompted to log in.

*   **Default Username:** `admin`
*   **Default Password:** `admin_pass` (This is set via the `ADMIN_PASSWORD` variable, which you can configure with `m3tal config wizard` or `m3tal config set ADMIN_PASSWORD your_new_password`).

**Important:** Change your dashboard password immediately after logging in.
You can change the dashboard password using the CLI:

```bash
sudo m3tal dashpass
```

## Filesystem Contract

M3TAL establishes a clear filesystem contract for its operational data and configurations:

*   `/etc/m3tal/.env`: The primary configuration file, managed by `m3tal config wizard`.
*   `/var/lib/m3tal/state.db`: The SQLite database storing M3TAL's state, automatically created by the API daemon.
*   `/opt/m3tal/stack/`: The canonical directory containing all Docker Compose files and Traefik configurations.
*   `/docker`: A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations.
*   `/docker/users.json`: Stores dashboard user credentials, managed by `m3tal dashpass`.

## Port Table

M3TAL components use the following ports:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if Traefik is exposed)              |
| 8080 | M3TAL API daemon (Go)       | Host-local only                             |
| 8081 | Traefik dashboard           | Host-local only                             |
| 8082 | M3TAL Dashboard container   | Direct port (local mode) or via Traefik     |

## Firewall Note

If you expose Traefik on port 80 to the public internet, ensure your firewall allows incoming connections on this port. For example, with `ufw`:

```bash
sudo ufw allow 80/tcp
sudo ufw enable
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands:

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

## Docker / Compose Runtime Explained

M3TAL leverages **Docker Engine** and **Docker Compose V2** to manage containerized services.

*   The `m3tal up` command executes `docker compose` operations across all `*-compose.yml` files located in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`). This allows you to manage multiple stacks (e.g., your own services) with a single command.
*   The `m3tal dash up` command is specifically tailored for the M3TAL Dashboard. It automatically downloads the necessary dashboard compose files and starts the dashboard container, dynamically applying the correct configuration (`local` or `traefik` exposure) based on the `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.
*   To add your own Docker Compose stacks, simply place your `my-service-compose.yml` files directly into the `/docker/` directory, and they will be managed by `m3tal up`.

## Dashboard Access Modes

The M3TAL Dashboard offers two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Mode 1: Local (Default)

*   **`DASHBOARD_EXPOSE_MODE=local`**
*   This mode uses the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`).
*   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082`
*   No Traefik reverse proxy is required for this mode, making it ideal for initial setup, local development, or LAN-only environments where direct IP access is preferred.

### Mode 2: Traefik

*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   This mode uses the `m3tal-compose.traefik.yml` override file. Instead of a direct port binding, it adds Traefik labels to the dashboard container. Traefik then discovers these labels and routes traffic for `dash.${DOMAIN}` to the dashboard container's internal port (8082).
*   **Access via:** `http://dash.DOMAIN` (e.g., `http://dash.example.com`)
*   This mode requires Traefik to be running via `m3tal up` and is best suited for domain-based access and integration into a multi-service setup behind a single reverse proxy.