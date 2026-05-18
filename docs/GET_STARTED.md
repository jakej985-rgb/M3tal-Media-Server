# GET_STARTED.md

Welcome to M3TAL! This guide will walk you through the initial setup of your M3TAL ecosystem.

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container management. Ensure these are installed on your system before proceeding.

You can verify their installation and version using the following commands:

```bash
docker --version && docker compose version
```

Expected output will show the installed versions of Docker Engine and Docker Compose V2.

## 2. Install M3TAL via APT

M3TAL is installed via a dedicated APT repository. Execute the following commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

This will install the `m3tal` CLI binary and the `m3tal-api.service` systemd daemon.

## 3. Run the Configuration Wizard

After installation, configure M3TAL using the interactive wizard. This sets up essential environment variables in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration values. Here's an explanation of key prompts:

*   **`DASHBOARD_EXPOSE_MODE`**:
    *   `local` (default): The M3TAL Dashboard will be directly accessible via `http://YOUR_IP:8082`. No Traefik configuration is required for dashboard access in this mode. This is suitable for local network access or initial setup.
    *   `traefik`: The M3TAL Dashboard will be routed through Traefik and accessible via `http://dash.YOUR_DOMAIN`. This mode requires Traefik to be running and a `DOMAIN` configured.
*   **`DOMAIN`**: Your primary domain (e.g., `example.com`). This is crucial if `DASHBOARD_EXPOSE_MODE` is set to `traefik`, or if you plan to use Traefik for other services (e.g., `api.example.com`). If not using Traefik, `localhost` is a suitable default.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL Dashboard will be exposed if `DASHBOARD_EXPOSE_MODE` is `local`. Default is `8082`.
*   **`PUID`** (Process User ID) and **`PGID`** (Process Group ID): These define the user and group IDs that containers will run as, primarily for file system permissions. It's recommended to use the PUID/PGID of a non-root user on your host system (e.g., `1000` for the first created user).
*   **`TZ`** (Timezone): Your desired timezone (e.g., `America/Denver`). This ensures correct time synchronization within containers.
*   **`DASHBOARD_SECRET`**: A secret key used by the dashboard for session management. Generate a strong, random string.
*   **`ADMIN_PASSWORD`**: The initial password for the M3TAL Dashboard administrator user. **Change this immediately after login.**
*   **`API_TOKEN`**: A token for authenticating with the M3TAL API. Generate a strong, random string.
*   **`BASE_STORAGE_PATH`**: The base directory for your container data volumes (e.g., `/mnt/data`).
*   **`CONFIG_PATH`**, **`MEDIA_PATH`**, **`DOWNLOADS_PATH`**: Sub-paths for specific types of data, typically within `BASE_STORAGE_PATH`.

## 4. Start the Routing Stack

The `m3tal up` command starts all Docker Compose stacks defined in the `/docker/` directory. This primarily includes the `routing` stack, which consists of Traefik (the reverse proxy) and optionally Cloudflared (for secure tunnels).

```bash
m3tal up
```

This command will ensure Traefik is running and listening on port 80. If you are using `DASHBOARD_EXPOSE_MODE=traefik`, this step is essential for accessing the dashboard via a domain name.

## 5. Start the M3TAL Dashboard

The M3TAL Dashboard is a containerized application managed specifically by `m3tal dash up`.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Pulls the latest M3TAL Dashboard Docker image.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the dashboard container, applying the correct Docker Compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) to configure its network exposure.

## 6. Access the M3TAL Dashboard

Open your web browser and navigate to the appropriate address based on your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host).
    You can also try `http://localhost:8082` if accessing from the host itself.

*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain configured in step 3).
    This requires DNS records for `dash.YOUR_DOMAIN` to point to your M3TAL host, and Traefik to be running via `m3tal up`.

## 7. Log In

The default username for the M3TAL Dashboard is `admin`. The initial password is the value you set for `ADMIN_PASSWORD` during the `m3tal config wizard` (default is `admin_pass` if not changed).

**It is critical to change this password immediately after your first login.**
You can change the dashboard password at any time using the `m3tal dashpass` command:

```bash
sudo m3tal dashpass
```

## Key Operational Details

### Filesystem Contract

The M3TAL ecosystem maintains a specific filesystem layout for its components and configuration.

| Path                         | Purpose                                                                                                 |
| :--------------------------- | :------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`            | Primary configuration file. Contains environment variables for M3TAL components. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`    | SQLite database storing M3TAL's internal state. Auto-created and managed by the API daemon.             |
| `/opt/m3tal/stack/`          | Canonical directory containing core Docker Compose files (e.g., `routing-compose.yml`).                 |
| `/docker`                    | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for placing and managing your Docker Compose stacks. |
| `/docker/users.json`         | Stores M3TAL Dashboard user credentials. Managed by `m3tal dashpass`.                                   |

### Port Map

M3TAL uses specific ports for its services:

| Port | Service                    | Access Mode                                             |
| :--- | :------------------------- | :------------------------------------------------------ |
| `80` | Traefik HTTP Entry Point   | Public (if Traefik is exposed)                          |
| `8080` | M3TAL API Daemon (Go)      | Host-local only (`http://localhost:8080` or `http://host.docker.internal:8080` for containers) |
| `8081` | Traefik Dashboard          | Host-local only (`http://localhost:8081`)               |
| `8082` | M3TAL Dashboard Container  | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

### Firewall Configuration

If you expose Traefik on port `80` to the internet or your local network, you may need to adjust your firewall rules. For systems using `ufw`, allow HTTP traffic:

```bash
sudo ufw allow 80
```

### M3TAL API Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard systemctl commands:

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

### Docker Compose Runtime

M3TAL leverages Docker Engine and Docker Compose V2.
*   The `m3tal up` command executes `docker compose up -d` across all `*-compose.yml` files present in the `/docker/` directory. This includes the `routing-compose.yml` for Traefik and any additional user-defined stacks.
*   The `m3tal dash up` command specifically manages the `m3tal-dashboard` container. It ensures the necessary compose files (`m3tal-compose.yml` and its overrides) are present and then starts the dashboard, applying the correct expose mode based on your configuration.
*   To add new services, place their `*-compose.yml` definitions into the `/docker/` directory, configure any necessary environment variables via `m3tal config wizard` or `m3tal config set`, and then run `m3tal up` to deploy them.