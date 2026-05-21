# Getting Started with M3TAL

Welcome to the M3TAL Ecosystem! This guide provides a complete, step-by-step process to get your M3TAL instance up and running for the first time.

---

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container management. Ensure these are installed and operational on your system before proceeding.

You can verify your Docker installation with the following commands:

```bash
docker --version && docker compose version
```

Expected output will show the installed versions for Docker Engine and Docker Compose V2.

## 2. Install M3TAL via APT

Install the M3TAL CLI binary and API daemon using the APT package manager.

```bash
# 1. Add the GPG signing key for M3TAL's repository
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository to your system's sources
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update your package lists and install the m3tal package
sudo apt update && sudo apt install -y m3tal
```

This will install the `/usr/bin/m3tal` CLI binary and set up the `m3tal-api.service` systemd daemon.

## 3. Run the Configuration Wizard

Initialize your M3TAL configuration using the interactive wizard. This will create or update `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's an explanation of key prompts:

-   **`DOMAIN`**: The base domain for your services (e.g., `example.com`). This is used by Traefik for routing. If you don't have a domain or are testing locally, `localhost` is a suitable default.
-   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    -   `local` (default): The dashboard will be directly accessible via a port binding (`http://YOUR_IP:8082`). Ideal for local network access or initial setup without a domain.
    -   `traefik`: The dashboard will be routed through the Traefik reverse proxy, accessible via a subdomain (e.g., `http://dash.YOUR_DOMAIN`). Requires Traefik to be running.
-   **`DASHBOARD_PORT`**: The direct port for the dashboard if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
-   **`PUID`** and **`PGID`**: The User ID and Group ID that containers will run as. Use `id -u` and `id -g` to find your current user's IDs. This is important for file permissions with mounted volumes.
-   **`TZ`**: Your local timezone (e.g., `America/Denver`). Used by containers for correct timestamps.
-   **`BASE_STORAGE_PATH`**, **`CONFIG_PATH`**, **`MEDIA_PATH`**, **`DOWNLOADS_PATH`**: These define the base directories for various data types. The wizard will suggest sensible defaults.
-   **`DASHBOARD_SECRET`**, **`API_TOKEN`**, **`ADMIN_PASSWORD`**: **These are critical security parameters.** Change the default values to strong, unique secrets immediately. These secure your dashboard login and API access.

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined in `/docker/`. This includes the core routing stack (Traefik).

```bash
m3tal up
```

This command will:
-   Pull necessary Docker images (e.g., `traefik:latest`).
-   Start the Traefik reverse proxy container, which listens on port 80 (HTTP) by default.
-   Start any other Docker Compose files found in `/docker/`.
-   Ensure the `proxy` Docker network is created for inter-container communication.

## 5. Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command will:
-   Download the latest dashboard compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
-   Read your `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
-   Pull the `ghcr.io/jakej985-rgb/m3tal-godash:debug` Docker image.
-   Start the `m3tal-dashboard` container with the appropriate configuration for direct port exposure or Traefik routing.

## 6. Open the M3TAL Dashboard in Your Browser

Access the M3TAL Dashboard using your web browser. The URL depends on your `DASHBOARD_EXPOSE_MODE` setting from the configuration wizard:

-   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server, or use `localhost` if accessing from the same machine).
-   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the `DOMAIN` you configured in the wizard). Traefik must be running for this mode to work.

## 7. Log In to the Dashboard

The dashboard will present a login screen.

-   **Default credentials:**
    -   **Username:** `admin`
    -   **Password:** The `ADMIN_PASSWORD` you set during the `m3tal config wizard` (default is `admin_pass` if you didn't change it).

**It is highly recommended to change your dashboard password immediately after your first login.**
You can do this using the CLI:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the `admin` user. The credentials are stored in `/docker/users.json`.

---

## Filesystem Contract

M3TAL adheres to a specific filesystem layout for configuration and data:

| Path                        | Purpose                                                                                |
| :-------------------------- | :------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file, managed by `m3tal config wizard`.                          |
| `/var/lib/m3tal/state.db`   | SQLite state database, auto-created and managed by the `m3tal-api.service` daemon.     |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's internal Docker Compose files and Traefik configuration. |
| `/docker`                   | A symbolic link to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations and where you place custom Docker Compose files. |
| `/docker/users.json`        | Dashboard credential store, managed by `m3tal dashpass`.                               |

## Port Table

These are the default ports used by M3TAL components:

| Port | Service                               | Access                                                          |
| :--- | :------------------------------------ | :-------------------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (if `DASHBOARD_EXPOSE_MODE=traefik` or other services are exposed via Traefik) |
| 8080 | M3TAL API daemon (Go)                 | Host-local access only (accessed by Dashboard and Traefik via `host.docker.internal`) |
| 8081 | Traefik dashboard (admin UI)          | Host-local only (e.g., `http://localhost:8081`)                 |
| 8082 | M3TAL Dashboard (Python/Flask)        | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (`dash.DOMAIN`) |

### Firewall Note

If you expose Traefik on port 80 to the public internet, you will need to open this port in your firewall. For example, using `ufw`:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon runs as a systemd service called `m3tal-api.service`. You can manage it using standard `systemctl` commands:

-   **Check the status:**
    ```bash
    systemctl status m3tal-api
    ```
-   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
-   **View real-time logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```

---

You have now successfully set up and configured your M3TAL Ecosystem. Explore the dashboard to manage your services!