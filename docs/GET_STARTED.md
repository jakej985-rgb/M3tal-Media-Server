```markdown
# GET_STARTED.md: First-Time User Setup Guide

This guide provides a step-by-step process for first-time users to set up and get started with the M3TAL Ecosystem.

---

## 1. Prerequisites

Before installing M3TAL, ensure the following are installed on your system:

-   **Docker Engine**
-   **Docker Compose V2**

You can verify their installation and version by running:

```bash
docker --version && docker compose version
```

## 2. Install M3TAL via APT

Use the following commands to add the M3TAL APT repository and install the `m3tal` CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Run the Configuration Wizard

Initialize your M3TAL environment and set essential configuration variables using the interactive wizard.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values:

-   **`DOMAIN`**: The base domain for services exposed via Traefik (e.g., `yourdomain.com`). If you're not using a domain, `localhost` is a suitable default.
-   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed:
    -   `local` (default): Accessible directly via `http://YOUR_IP:8082`. No Traefik configuration required for the dashboard.
    -   `traefik`: Accessible via `http://dash.YOUR_DOMAIN` using Traefik as a reverse proxy.
-   **`TRAEFIK_WEB_PORT`**: The host port Traefik will listen on for HTTP traffic (default: `80`).
-   **`PUID` / `PGID`**: The User ID (PUID) and Group ID (PGID) that containers will run as. Typically, your user's PUID/PGID (default: `1000:1000`). Find yours with `id -u` and `id -g`.
-   **`TZ`**: Your local timezone (e.g., `America/Denver`). This ensures containers have correct time.
-   **`BASE_STORAGE_PATH`**: The base directory for M3TAL application data, including configuration and volumes for other services (default: `/mnt`).
-   **`DASHBOARD_SECRET`**: A unique secret key for the dashboard. **Change this from the default immediately.** This enhances dashboard security.
-   **`API_TOKEN`**: A token for programmatic access to the M3TAL API. **Change this from the default immediately.**

These settings are saved to `/etc/m3tal/.env`.

## 4. Start the Routing Stack

Start the core routing services, including Traefik (the reverse proxy) and optionally Cloudflared (for secure tunnels).

```bash
m3tal up
```

This command initiates Docker Compose operations for all `*-compose.yml` files located in `/docker/` (which is a symlink to `/opt/m3tal/stack/`). This includes `routing-compose.yml` which defines Traefik and Cloudflared.

## 5. Start the M3TAL Dashboard

Launch the M3TAL web-based dashboard for managing your ecosystem.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest M3TAL Dashboard Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the appropriate Compose override file (`m3tal-compose.local.yml` for `local` mode, or `m3tal-compose.traefik.yml` for `traefik` mode). This pulls the dashboard image (`ghcr.io/jakej985-rgb/m3tal-godash:debug`) if not already present, and then starts the container.

## 6. Open M3TAL Dashboard in Browser

Access the M3TAL Dashboard through your web browser based on your configured `DASHBOARD_EXPOSE_MODE`:

-   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server).
-   **If `DASHBOARD_EXPOSE_MODE` is `traefik` and Traefik is properly configured:**
    Open `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured in the wizard).

## 7. Log In to the Dashboard

The first time you access the dashboard, you will be prompted for credentials:

-   **Default Username:** `admin`
-   **Default Password:** `admin_pass` (This is the value of `ADMIN_PASSWORD` in `/etc/m3tal/.env`).

**It is highly recommended to change the default password immediately.** You can do this using the CLI:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL utilizes specific paths for configuration, data, and stack management:

| Path                       | Purpose                                                                |
| :------------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary configuration file, managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`  | SQLite database storing M3TAL's internal state, auto-created by API.   |
| `/opt/m3tal/stack/`        | Canonical directory for M3TAL's core Compose files and Traefik config. |
| `/docker`                  | Symlink that points to `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`       | Dashboard credential store, managed by `m3tal dashpass`.               |

## Port Table

The following ports are used by M3TAL components:

| Port | Service               | Access                                                 |
| :--- | :-------------------- | :----------------------------------------------------- |
| 80   | Traefik HTTP entry point | Public (if `TRAEFIK_WEB_PORT` is 80 and Traefik is running) |
| 8080 | M3TAL API daemon (Go) | Host-local only, used by dashboard and internal services |
| 8081 | Traefik dashboard     | Host-local only, for monitoring Traefik                |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

## Firewall Note

If you plan to expose M3TAL services (like Traefik) to your local network or the internet, you may need to open ports in your firewall. For example, to allow HTTP traffic through Traefik on port 80:

```bash
sudo ufw allow 80
```

If you configure HTTPS (port 443), you would also need to allow that port.

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`:

-   **Check the status of the M3TAL API service:**
    ```bash
    systemctl status m3tal-api
    ```
-   **View real-time logs for the M3TAL API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```

---

## Dashboard Access Modes

The M3TAL Dashboard supports two primary access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. The `m3tal dash up` command dynamically selects the appropriate Docker Compose override file to achieve these modes.

### Mode 1: Local (Default)

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`
-   **Compose Override:** Uses `m3tal-compose.local.yml`. This file adds a direct port binding to the `m3tal-dashboard` container, typically `${DASHBOARD_PORT:-8082}:8082`.
-   **Access Method:** `http://HOST_IP:8082` or `http://localhost:8082`.
-   **Requirements:** No Traefik configuration is needed for the dashboard itself, making it suitable for quick local setups and initial testing.
-   **Best For:** LAN-only deployments, users setting up M3TAL on a home server, or first-time users.

### Mode 2: Traefik

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`
-   **Compose Override:** Uses `m3tal-compose.traefik.yml`. This file applies Traefik labels to the `m3tal-dashboard` container. These labels instruct Traefik to route incoming requests for `dash.YOUR_DOMAIN` to the dashboard container's internal port 8082.
-   **Access Method:** `http://dash.YOUR_DOMAIN` (Traefik must be running via `m3tal up` and configured correctly).
-   **Requirements:** Traefik (from `routing-compose.yml`) must be running and properly configured to handle domain-based routing.
-   **Best For:** Deployments requiring domain-based access, integration with other services behind a reverse proxy, and production environments with proper DNS setup.

## Docker / Compose Runtime

M3TAL is built on **Docker Engine** and **Docker Compose V2**. These are fundamental dependencies for M3TAL's operation.

-   The `m3tal up` command orchestrates all Docker Compose files ending with `*-compose.yml` that are located in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`). This allows you to define and manage multiple services as "stacks."
-   The `m3tal dash up` command is specifically tailored for the M3TAL Dashboard. It automatically:
    1.  Downloads the necessary `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from the M3TAL GitHub repository.
    2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
    3.  Starts the `m3tal-dashboard` container using the base `m3tal-compose.yml` combined with the selected override file, ensuring the dashboard is exposed according to your configuration.
-   **User-defined stacks** should be placed as `my-stack-compose.yml` files within the `/docker/` directory. Running `m3tal up` will then bring up all these services, along with M3TAL's core components.