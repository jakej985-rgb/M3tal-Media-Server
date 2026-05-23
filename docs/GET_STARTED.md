```markdown
# GET STARTED with M3TAL

This guide provides a complete, step-by-step setup for first-time M3TAL users.

---

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container orchestration. These must be installed on your system before proceeding.

Verify your Docker installation:

```bash
docker --version && docker compose version
```

Expected output will show versions for Docker Engine and Docker Compose V2.

---

## 2. Install M3TAL via APT

M3TAL is distributed as a single Go binary via a Debian APT repository.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## 3. Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables for M3TAL operations. This generates or updates `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values:

*   **`DOMAIN`**: (Default: `localhost`) This is used by Traefik for domain-based routing. If you plan to expose services on the internet, set your domain here (e.g., `example.com`). If only for local access, `localhost` is sufficient.
*   **`DASHBOARD_EXPOSE_MODE`**: (Default: `local`)
    *   **`local`**: The M3TAL dashboard will be directly exposed on a specific port (default 8082) on the host machine. This is ideal for quick local access or when Traefik is not in use. Access via `http://YOUR_IP:8082`.
    *   **`traefik`**: The M3TAL dashboard will be routed through the Traefik reverse proxy. It will be accessible via `http://dash.YOUR_DOMAIN` (e.g., `http://dash.localhost` or `http://dash.example.com`). This requires Traefik to be running.
*   **`DASHBOARD_PORT`**: (Default: `8082`) The port on which the M3TAL dashboard will be exposed directly if `DASHBOARD_EXPOSE_MODE` is set to `local`.
*   **`PUID` / `PGID`**: (Default: `1000` / `1000`) These specify the User ID and Group ID that Docker containers will run as inside the container. It's best practice to match these with a user on your host system to prevent permission issues with mounted volumes.
*   **`TZ`**: (Default: `America/Denver`) Your server's timezone.
*   **`BASE_STORAGE_PATH`**: (Default: `./data`) The base directory on your host that will be mounted into containers, typically for application data. All other path variables are subdirectories of this.
*   **`CONFIG_PATH`**: (Default: `./data/config`) Path for configuration files.
*   **`DOWNLOADS_PATH`**: (Default: `./data/downloads`) Path for downloaded content.
*   **`MEDIA_PATH`**: (Default: `./data/media`) Path for media content.
*   **`DASHBOARD_SECRET`**: (Default: `change_me_immediately`) A secret key used by the dashboard. **Change this immediately for security.**
*   **`API_TOKEN`**: (Default: `change_me_api_token`) A token for M3TAL API access. **Change this immediately for security.**
*   **`ADMIN_PASSWORD`**: (Default: `admin_pass`) The default password for the `admin` user on the M3TAL dashboard. **Change this immediately for security.**

---

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found within the `/docker/` directory. This includes Traefik, which acts as the central routing gateway for your services.

```bash
m3tal up
```

This command will start containers defined in `routing-compose.yml` (Traefik and optionally Cloudflared) and any other compose files you might have placed in `/docker/`.

---

## 5. Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL dashboard container. It downloads the necessary compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) and starts the dashboard based on your `DASHBOARD_EXPOSE_MODE` setting.

```bash
m3tal dash up
```

This pulls the `ghcr.io/jakej985-rgb/m3tal-godash:debug` image (or latest) and starts the dashboard container with the appropriate port bindings or Traefik labels.

---

## 6. Open Browser

Access the M3TAL dashboard in your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your server or `localhost` if accessing from the server itself).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain configured in your `.env` file, e.g., `http://dash.localhost` or `http://dash.example.com`). This requires Traefik to be running via `m3tal up`.

---

## 7. Log In

The default credentials for the M3TAL dashboard are:

*   **Username**: `admin`
*   **Password**: `admin_pass` (as set by the wizard, or the default `ADMIN_PASSWORD` if not changed)

**It is critical to change this password immediately after your first login.**
You can change the dashboard password using the M3TAL CLI:

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL interacts with specific files and directories on your system:

*   `/etc/m3tal/.env`: The primary configuration file, managed by `m3tal config wizard`.
*   `/var/lib/m3tal/state.db`: The SQLite database storing M3TAL's internal state. This is automatically created and managed by the API daemon.
*   `/opt/m3tal/stack/`: The canonical directory where M3TAL stores its core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`).
*   `/docker`: A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for placing your custom Docker Compose files. `m3tal up` will look for all `*-compose.yml` files here.
*   `/docker/users.json`: Stores M3TAL dashboard user credentials. Managed by `m3tal dashpass`.

---

## Port Table

M3TAL uses the following ports:

| Port | Service            | Access                                     |
| :--- | :----------------- | :----------------------------------------- |
| 80   | Traefik HTTP       | Public (if Traefik is exposed)             |
| 8080 | M3TAL API daemon   | Host-local                                 |
| 8081 | Traefik Dashboard  | Host-local only                            |
| 8082 | M3TAL Dashboard    | Direct port (local mode) or via Traefik    |

---

## Firewall Note

If you are exposing Traefik (and thus any services, including the M3TAL Dashboard in `traefik` mode) to the internet, you must open port 80 on your firewall.

For UFW (Uncomplicated Firewall):

```bash
sudo ufw allow 80/tcp
sudo ufw enable
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using standard `systemctl` and `journalctl` commands:

*   **Check API service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live API service logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the API service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
```