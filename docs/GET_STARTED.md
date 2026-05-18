```markdown
# GETTING STARTED WITH M3TAL

This guide provides a step-by-step process for setting up M3TAL on your server for the first time.

---

## 1. Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify their installation with:

```bash
docker --version && docker compose version
```

Example expected output:
```
Docker version 24.0.7, build afdd53b
Docker Compose version v2.23.3
```

---

## 2. Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the `m3tal` CLI:

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

The configuration wizard guides you through essential initial setup.

```bash
sudo m3tal config wizard
```

You will be prompted to set several configuration variables. Here's what each important prompt means:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard will be accessible.
    *   `local` (default): The dashboard will be directly accessible via a port on your host machine (e.g., `http://YOUR_IP:8082`). This mode is suitable for local network access or if you don't intend to use Traefik for routing.
    *   `traefik`: The dashboard will be routed via the Traefik reverse proxy (e.g., `http://dash.YOUR_DOMAIN`). This mode requires Traefik to be running and assumes you have a domain configured.
*   **`DOMAIN`**: Your primary domain name (e.g., `example.com`). This is crucial if `DASHBOARD_EXPOSE_MODE` is set to `traefik`, as it forms the basis for subdomains like `dash.example.com`.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be exposed if `DASHBOARD_EXPOSE_MODE` is `local`. The default is `8082`.
*   **`PUID` (User ID)** and **`PGID` (Group ID)**: These specify the user and group IDs that Docker containers will run as, ensuring proper file permissions. You can find your current user's IDs using `id -u` and `id -g`.
*   **`TZ` (Time Zone)**: Your system's time zone (e.g., `America/New_York`).
*   **`BASE_STORAGE_PATH`**, **`CONFIG_PATH`**, **`MEDIA_PATH`**, **`DOWNLOADS_PATH`**: These define the base directories for various types of data managed by M3TAL and its services. By default, they are relative to `/mnt`.
*   **`DASHBOARD_SECRET`**: A secret key used to secure the dashboard. It's automatically generated if left blank, but you should change it.
*   **`API_TOKEN`**: A token used for internal API authentication. It's automatically generated if left blank, but you should change it.
*   **`ADMIN_PASSWORD`**: The default password for the M3TAL Dashboard. You will change this after logging in for the first time.

After completing the wizard, your configuration will be saved to `/etc/m3tal/.env`.

---

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks found in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This includes the core `routing-compose.yml` that deploys Traefik, the reverse proxy.

```bash
m3tal up
```

This command initiates the Traefik gateway and any other core infrastructure services defined in the default compose files. Traefik will expose port 80 (and 443 if HTTPS is configured) on your host.

---

## 5. Start the M3TAL Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It ensures you have the latest compose files and starts the dashboard according to your `DASHBOARD_EXPOSE_MODE` setting.

```bash
m3tal dash up
```

This command:
1.  Downloads the latest `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container, applying the appropriate compose override file to expose it either directly or via Traefik.

---

## 6. Access the M3TAL Dashboard

Open your web browser and navigate to the appropriate address based on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open `http://YOUR_SERVER_IP:8082`
    (Replace `YOUR_SERVER_IP` with the actual IP address of your server).

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open `http://dash.YOUR_DOMAIN`
    (Replace `YOUR_DOMAIN` with the domain you configured in step 3).

---

## 7. Log In to the Dashboard

The first time you access the dashboard, you will be prompted to log in.

*   **Default Username:** `admin`
*   **Default Password:** The `ADMIN_PASSWORD` you set during the configuration wizard (default is `admin_pass` if not changed).

**It is highly recommended to change the default password immediately.** You can do this within the dashboard settings, or via the CLI:

```bash
sudo m3tal dashpass <NEW_PASSWORD>
```

---

## Filesystem Contract

M3TAL uses specific locations for its configuration and state data:

| Path                        | Purpose                                               |
| :-------------------------- | :---------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the API daemon. |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config. |
| `/docker`                   | **Symlink** → `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your Docker Compose files here. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`. |

---

## Port Table

These are the primary ports used by M3TAL components:

| Port | Service                     | Access                       |
| :--- | :-------------------------- | :--------------------------- |
| 80   | Traefik HTTP entry point    | Public (if Traefik is enabled) |
| 8080 | M3TAL API daemon (Go)       | Host-local only              |
| 8082 | M3TAL Dashboard             | Direct (local mode) or via Traefik (traefik mode) |

---

## Firewall Note

If you are running a firewall (e.g., UFW), you will need to allow traffic on the ports M3TAL exposes:

*   If `DASHBOARD_EXPOSE_MODE=traefik` (Traefik is handling routing):
    ```bash
    sudo ufw allow 80/tcp comment "Allow HTTP traffic for Traefik"
    # If using HTTPS: sudo ufw allow 443/tcp comment "Allow HTTPS traffic for Traefik"
    ```
*   If `DASHBOARD_EXPOSE_MODE=local` (direct dashboard access):
    ```bash
    sudo ufw allow 8082/tcp comment "Allow M3TAL Dashboard direct access"
    ```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl`:

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

---

You have successfully set up and accessed your M3TAL Ecosystem!
```