```markdown
# GET_STARTED.md

This guide provides a complete setup for first-time M3TAL users.

---

### Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify their installation and version:

```bash
docker --version && docker compose version
```

Example output:
```
Docker version 24.0.7, build afdd53b
Docker Compose version v2.23.0-desktop.1
```

### Step 2: Install M3TAL via APT

Execute the following commands to add the M3TAL APT repository and install the `m3tal` CLI binary:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This installs the `m3tal` CLI binary to `/usr/bin/m3tal` and the M3TAL API daemon (`m3tal-api.service`) which runs as a systemd service.

### Step 3: Run the Configuration Wizard

Initialize your M3TAL environment by running the configuration wizard. This wizard guides you through setting up essential environment variables stored in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for the following information:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard is directly accessible via a host port binding (e.g., `http://YOUR_IP:8082`). This mode uses the `m3tal-compose.local.yml` Docker Compose override. It's suitable for LAN-only setups or initial testing without a domain.
    *   `traefik`: The dashboard is exposed via Traefik (e.g., `http://dash.YOUR_DOMAIN`). This mode applies Traefik labels using the `m3tal-compose.traefik.yml` Docker Compose override, allowing Traefik to route traffic to the dashboard container. This requires Traefik to be running.
*   **`DOMAIN`**: Your primary domain (e.g., `example.com`). This is used for Traefik routing rules if `DASHBOARD_EXPOSE_MODE` is set to `traefik`, and also for other services you deploy. If using `local` mode, `localhost` is a common choice.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management. Generate a strong, random value.
*   **`API_TOKEN`**: A token used to authenticate requests to the M3TAL API daemon. Generate a strong, random value.
*   **`ADMIN_PASSWORD`**: The default password for the M3TAL Dashboard `admin` user. Choose a strong password. You can change this later.
*   **`PUID`** and **`PGID`**: The User ID and Group ID that containers will run as. Typically, these are `1000` for the default user on most Linux distributions. This ensures proper file permissions for mounted volumes.
*   **`TZ`**: Your timezone (e.g., `America/Denver`).
*   **`BASE_STORAGE_PATH`**: The base directory for all your M3TAL data and persistent storage (e.g., `/mnt/m3tal-data`).
*   **`CONFIG_PATH`**: A subdirectory within `BASE_STORAGE_PATH` for configuration files (e.g., `${BASE_STORAGE_PATH}/config`).
*   **`MEDIA_PATH`**: A subdirectory within `BASE_STORAGE_PATH` for media files (e.g., `${BASE_STORAGE_PATH}/media`).
*   **`DOWNLOADS_PATH`**: A subdirectory within `BASE_STORAGE_PATH` for downloads (e.g., `${BASE_STORAGE_PATH}/downloads`).

### Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose files located in the `/docker/` directory. This includes the `routing-compose.yml` file, which deploys Traefik, the reverse proxy gateway.

```bash
m3tal up
```

This command initiates the Traefik container, which binds to host port 80 and serves as the HTTP entry point for all services. Traefik automatically discovers services via Docker labels and loads dynamic configurations from `/etc/traefik/dynamic`. It routes `api.DOMAIN` to the M3TAL API daemon.

### Step 5: Start the M3TAL Dashboard

The `m3tal dash up` command is specifically used to deploy and manage the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` files.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the base `m3tal-compose.yml` and applies the appropriate override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE`.

The dashboard container (`m3tal-dashboard`) communicates with the M3TAL API daemon at `http://host.docker.internal:8080`.

### Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the M3TAL Dashboard based on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE=local`**:
    Access via `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your host machine) or `http://localhost:8082` if accessing from the host itself.
*   **If `DASHBOARD_EXPOSE_MODE=traefik`**:
    Access via `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Traefik must be running for this mode to work.

### Step 7: Log In to the Dashboard

The default login credentials for the M3TAL Dashboard are:

*   **Username**: `admin`
*   **Password**: The `ADMIN_PASSWORD` you set during the `m3tal config wizard` in Step 3.

To change the dashboard password for the `admin` user, use the following command:

```bash
sudo m3tal dashpass admin new_strong_password
```

---

### Filesystem Contract

M3TAL relies on specific file and directory locations for its operation. Understanding this contract is essential for configuration and troubleshooting.

| Path                         | Purpose                                                     |
| :--------------------------- | :---------------------------------------------------------- |
| `/etc/m3tal/.env`            | Primary configuration file. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`    | SQLite state database. Auto-created and managed by the API daemon. |
| `/opt/m3tal/stack/`          | Canonical directory for M3TAL's internal Docker Compose files and Traefik configuration. |
| `/docker`                    | Symlink that points to `/opt/m3tal/stack/`. This is the user-facing path for placing custom Docker Compose stacks and performing stack operations. |
| `/docker/users.json`         | Dashboard credential store. Managed by `m3tal dashpass`. |

---

### Port Map

The following ports are used by M3TAL components:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP entry point    | Public (if Traefik is exposed)              |
| 8080 | M3TAL API daemon (Go)       | Host-local only                             |
| 8081 | Traefik dashboard           | Host-local only                             |
| 8082 | M3TAL Dashboard container | Direct port (local mode) or via Traefik (traefik mode) |

### Firewall Configuration

If you are exposing Traefik (e.g., `DASHBOARD_EXPOSE_MODE=traefik` or serving other internet-facing services), you might need to allow incoming traffic on port 80 (and 443 for HTTPS if configured) through your firewall.

Example for `ufw`:

```bash
sudo ufw allow 80/tcp
# sudo ufw allow 443/tcp # If you configure HTTPS
sudo ufw enable
```

### M3TAL API Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage its state using `systemctl` and view its logs using `journalctl`.

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View real-time logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
```