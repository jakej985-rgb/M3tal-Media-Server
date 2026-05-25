# M3TAL Getting Started Guide

This guide provides a complete, step-by-step setup for first-time M3TAL users.

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your installation using the following commands:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

Execute the following commands in your terminal to install the M3TAL CLI binary and related system components.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The M3TAL configuration wizard guides you through setting up essential environment variables for your M3TAL instance. These settings are stored in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

During the wizard, you will be prompted for various settings. Here's an explanation of key prompts:

*   **`DASHBOARD_EXPOSE_MODE`**:
    *   **`local` (default)**: The M3TAL Dashboard will be directly exposed on a specific host port (default 8082). Recommended for initial setup and LAN-only access.
    *   **`traefik`**: The M3TAL Dashboard will be routed through Traefik (if `m3tal up` is run) and accessible via a domain (e.g., `dash.yourdomain.com`). This requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The host port to expose the M3TAL Dashboard on if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`DOMAIN`**: Your primary domain name (e.g., `example.com`). This is used for Traefik routing if `DASHBOARD_EXPOSE_MODE` is `traefik` and for other services. Defaults to `localhost`.
*   **`PUID`** and **`PGID`**: The User ID (PUID) and Group ID (PGID) that Docker containers will run as. Use `id -u` and `id -g` to find your current user's IDs. Defaults to `1000`.
*   **`TZ`**: Your timezone (e.g., `America/New_York`). This ensures containers display correct timestamps. Defaults to `America/Denver`.
*   **`BASE_STORAGE_PATH`**: The base directory on your host where M3TAL-managed services will store their data. Defaults to `/mnt`.
*   **`CONFIG_PATH`**: The directory for configuration files. Defaults to `/mnt/config`.

## Step 4: Start the Routing Stack

This command initializes and starts all Docker Compose stacks found in the `/docker/` directory, including the Traefik reverse proxy (from `routing-compose.yml`).

```bash
m3tal up
```

This ensures that Traefik is running and can route traffic for services, including the M3TAL API daemon and the Dashboard if configured for `traefik` expose mode.

## Step 5: Start the M3TAL Dashboard

This command manages the M3TAL Dashboard container. It performs the following actions:
1.  Pulls the latest M3TAL Dashboard Docker image from `ghcr.io/jakej985-rgb/m3tal-godash`.
2.  Reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Starts the dashboard container using the appropriate Docker Compose override file (`m3tal-compose.local.yml` for `local` mode or `m3tal-compose.traefik.yml` for `traefik` mode).

```bash
m3tal dash up
```

## Step 6: Open the Dashboard in Your Browser

Access the M3TAL Dashboard using your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with the IP address of your M3TAL host, or use `localhost` if accessing from the same machine).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Traefik must be running.

## Step 7: Log In to the Dashboard

The default login credentials for the M3TAL Dashboard are `admin` for the username and `admin` for the password.

**It is strongly recommended to change these default credentials immediately.** You can change the dashboard password using the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL establishes a clear filesystem contract for its operation and data storage:

*   `/etc/m3tal/.env`: The primary configuration file, managed by `sudo m3tal config wizard` and `m3tal config set`.
*   `/var/lib/m3tal/state.db`: The SQLite state database, automatically created and managed by the M3TAL API daemon.
*   `/opt/m3tal/stack/`: The canonical directory where M3TAL stores core Docker Compose files and Traefik configuration.
*   `/docker`: A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for placing and managing your Docker Compose stacks. Any `*-compose.yml` file placed here will be managed by `m3tal up`.
*   `/docker/users.json`: The credential store for the M3TAL Dashboard, managed by `sudo m3tal dashpass`.

## Port Table

M3TAL and its components utilize specific network ports:

| Port | Service                                | Access                                                                 |
| :--- | :------------------------------------- | :--------------------------------------------------------------------- |
| 80   | Traefik HTTP entry point               | Public (when `DASHBOARD_EXPOSE_MODE=traefik` or other services exposed) |
| 8080 | M3TAL API daemon (Go)                  | Host-local only                                                        |
| 8081 | Traefik dashboard                      | Host-local only (for Traefik's own dashboard)                          |
| 8082 | M3TAL Dashboard (Python/Flask)         | Direct port (local mode) or via Traefik (traefik mode)                 |

## Firewall Note

If you intend to expose Traefik on port 80 to the internet or your local network, you must allow traffic through your firewall. For `ufw` (Uncomplicated Firewall) on Ubuntu/Debian:

```bash
sudo ufw allow 80/tcp
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor its status using standard `systemctl` commands:

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the API daemon:**
    ```bash
    sudo systemctl restart m3tal-api
    ```