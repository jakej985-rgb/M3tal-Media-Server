```markdown
# M3TAL Getting Started Guide

This guide provides a complete, newbie-friendly setup for first-time users of the M3TAL ecosystem.

---

## 1. Prerequisites

M3TAL relies on Docker Engine and Docker Compose V2 for container management. These must be installed on your system before proceeding.

Verify your Docker installation:

```bash
docker --version && docker compose version
```

---

## 2. Install M3TAL via APT

M3TAL is distributed as a single Go binary through an APT repository.

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

After installation, configure M3TAL using the interactive wizard. This sets up the primary configuration file (`/etc/m3tal/.env`).

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several configuration values:

*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard will be accessible.
    *   **`local` (default)**: The dashboard will be directly exposed on port `8082` of your host machine. This is ideal for initial setup and local/LAN access without a reverse proxy.
    *   **`traefik`**: The dashboard will be routed through the Traefik reverse proxy, accessible via a domain (e.g., `dash.yourdomain.com`). This requires Traefik to be running.
*   **`DOMAIN`**: The base domain for services exposed via Traefik (e.g., `example.com`). If you chose `local` for `DASHBOARD_EXPOSE_MODE`, `localhost` is usually sufficient for now.
*   **`DASHBOARD_SECRET`**: A secret key used by the M3TAL Dashboard for session management and security. Generate a strong, random value.
*   **`API_TOKEN`**: An authentication token required for interacting with the M3TAL API daemon. Generate a strong, random value.
*   **`ADMIN_PASSWORD`**: The initial password for the `admin` user of the M3TAL Dashboard. You will use this to log in the first time.
*   **`PUID` / `PGID`**: The User ID (PUID) and Group ID (PGID) that containers will use to run processes. This ensures proper file permissions for mounted volumes. Typically, `1000` for both is the default for the first user on a Linux system.
*   **`TZ`**: Your local timezone (e.g., `America/New_York`).
*   **`BASE_STORAGE_PATH`**: The base directory where M3TAL-managed applications will store their data. Defaults to `/mnt` or similar.
*   **`CONFIG_PATH`**: Directory for application configuration files.
*   **`DOWNLOADS_PATH`**: Directory for application downloads.
*   **`MEDIA_PATH`**: Directory for media files.

---

## 4. Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined by `*-compose.yml` files located in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). This includes the `routing-compose.yml` which deploys Traefik.

```bash
m3tal up
```

This command will bring up the Traefik reverse proxy, making it ready to route traffic to your services.

---

## 5. Start the M3TAL Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It downloads the necessary Compose files, reads your `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`, and starts the dashboard with the appropriate configuration.

```bash
m3tal dash up
```

This command pulls the `m3tal-godash` Docker image (if not already present) and starts the M3TAL Dashboard container.

---

## 6. Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard using the address determined by your `DASHBOARD_EXPOSE_MODE` setting:

*   **If `DASHBOARD_EXPOSE_MODE` is `local`**:
    Access via `http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the actual IP address of your M3TAL host, e.g., `http://192.168.1.100:8082`). You can find your server's IP with `ip a`.
    Alternatively, if accessing from the host itself, use `http://localhost:8082`.

*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`**:
    Access via `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured during the wizard, e.g., `http://dash.example.com`). This assumes Traefik is running and correctly configured with DNS pointing to your M3TAL host.

---

## 7. Log In

Use the following credentials to log in to the M3TAL Dashboard:

*   **Username**: `admin`
*   **Password**: The `ADMIN_PASSWORD` you set during the `m3tal config wizard` (or the value in `/etc/m3tal/.env`).

**To change the default dashboard password after logging in, use the command:**

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

The M3TAL ecosystem establishes a clear contract for file locations:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created by the M3TAL API daemon.         |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains compose files and Traefik config.|
| `/docker`                   | Symlink → `/opt/m3tal/stack/`. User-facing path for all stack operations. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |

---

## M3TAL Core Operations

**Docker Compose Runtime**

M3TAL leverages **Docker Engine + Docker Compose V2**.
The `m3tal up` command orchestrates all Docker Compose applications by executing `docker compose up -d` for every `*-compose.yml` file found within the `/docker/` directory.

The `m3tal dash up` command performs these specific actions for the dashboard:
1.  Downloads the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` files from GitHub into `/docker/`.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container, applying the relevant override (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE`.

**Adding New Stacks (Day 2 Operations)**

To deploy a new Docker Compose stack:
1.  Place your Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
2.  Ensure any required environment variables for your stack are defined in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set KEY value`).
3.  Run `m3tal up` to bring up all defined stacks, including your new one.

---

## Understanding M3TAL Dashboard Access Modes

The `DASHBOARD_EXPOSE_MODE` environment variable in `/etc/m3tal/.env` dictates how the M3TAL Dashboard is exposed:

### Mode 1: `local` (Default)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
*   **Mechanism**: M3TAL uses the `m3tal-compose.local.yml` override file. This file adds a direct port binding to the dashboard container, typically mapping host port `8082` to the container's internal port `8082`.
*   **Access**: You access the dashboard directly via `http://YOUR_HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: No Traefik reverse proxy is required for this mode.
*   **Best for**: LAN-only setups, initial testing, or when you don't need domain-based access.

### Mode 2: `traefik`

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
*   **Mechanism**: M3TAL uses the `m3tal-compose.traefik.yml` override file. This adds specific Traefik labels to the dashboard container, instructing Traefik to route requests for `dash.YOUR_DOMAIN` to the dashboard's internal port `8082`.
*   **Access**: You access the dashboard via `http://dash.YOUR_DOMAIN`.
*   **Requirements**: The Traefik gateway (`routing-compose.yml`) must be running via `m3tal up`, and your DNS must resolve `dash.YOUR_DOMAIN` to your M3TAL host's IP address.
*   **Best for**: Domain-based access, integrating with other services behind Traefik, or when you require TLS/SSL through Traefik.

---

## Port Map

The following ports are relevant to the M3TAL ecosystem:

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| `80` | Traefik HTTP entry point    | Public (if Traefik is exposed)              |
| `8080` | M3TAL API daemon (Go)       | Host-local only (internal communication)    |
| `8081` | Traefik Dashboard (Internal)| Host-local only (internal to Traefik)       |
| `8082` | M3TAL Dashboard             | Direct port (`local` mode) or via Traefik (`traefik` mode) |

---

## Firewall Configuration

If you have a firewall enabled (e.g., `ufw`) and are using Traefik to expose services to the internet (or LAN), you will need to allow traffic on port 80:

```bash
sudo ufw allow 80/tcp
sudo ufw enable # if not already enabled
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands.

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
```