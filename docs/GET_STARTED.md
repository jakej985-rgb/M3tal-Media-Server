```markdown
# M3TAL Setup Guide

This guide provides step-by-step instructions for installing and setting up M3TAL for the first time.

---

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify their installation and versions using the following commands:

```bash
docker --version && docker compose version
```

## Step 2: Install M3TAL via APT

Execute the following commands in your terminal to install the M3TAL CLI and API daemon:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, run the configuration wizard to set up essential environment variables for M3TAL operation.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values:

-   **DOMAIN**: The primary domain name for services exposed via Traefik (e.g., `example.com`). Defaults to `localhost`. If you intend to access services like the dashboard via `dash.example.com`, this should be your domain.
-   **LOCAL_IP**: The local IP address of the machine M3TAL is running on. This is used for direct access to services (e.g., the dashboard in local mode) and internal communication.
-   **DASHBOARD_EXPOSE_MODE**: Controls how the M3TAL dashboard is accessed.
    -   `local` (default): The dashboard is exposed directly on a host port (`LOCAL_IP:8082`). No Traefik configuration is needed for dashboard access.
    -   `traefik`: The dashboard is routed via Traefik using a subdomain (e.g., `dash.DOMAIN`). Traefik must be running for dashboard access.
-   **DASHBOARD_PORT**: The host port to expose the M3TAL dashboard on, if `DASHBOARD_EXPOSE_MODE` is set to `local`. Defaults to `8082`.
-   **DASHBOARD_SECRET**: A secret key used for securing dashboard sessions. It is recommended to change the default value.
-   **ADMIN_PASSWORD**: The default password for the `admin` user of the M3TAL dashboard. It is recommended to change the default value.
-   **API_TOKEN**: An API token used for programmatic access to the M3TAL API.
-   **PUID** / **PGID**: The User ID and Group ID that containers will run as inside Docker. Defaults to `1000`, corresponding to the first non-root user on most Linux systems.
-   **TZ**: The timezone for containers (e.g., `America/Denver`).
-   **BASE_STORAGE_PATH**: The base directory on the host where M3TAL will store all persistent data, including stack configurations, media, and downloads (e.g., `/mnt`).

## Step 4: Start the Routing Stack

Start the core M3TAL routing stack, which includes Traefik, the reverse proxy.

```bash
m3tal up
```
This command processes and starts all Docker Compose files located within the `/docker/` directory. By default, this includes the `routing-compose.yml` which deploys Traefik and potentially Cloudflared for network ingress.

## Step 5: Start the Dashboard

Start the M3TAL dashboard container. This command handles pulling the necessary Docker image and configuring the container based on your settings.

```bash
m3tal dash up
```
This command pulls the `m3tal-dashboard` image and starts its container, applying the appropriate configuration based on the `DASHBOARD_EXPOSE_MODE` set in `/etc/m3tal/.env`. For details on access modes, refer to the [Dashboard Access Modes](#dashboard-access-modes) section.

## Step 6: Open Browser and Access Dashboard

Once the dashboard is running, open your web browser to access it:

-   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_LOCAL_IP:8082` (replace `YOUR_LOCAL_IP` with the `LOCAL_IP` configured in Step 3).
-   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the `DOMAIN` configured in Step 3).

## Step 7: Log In to the Dashboard

The default login credentials for the dashboard are:

-   **Username:** `admin`
-   **Password:** The value set for `ADMIN_PASSWORD` during the configuration wizard (default `admin_pass`).

### Change Dashboard Password

To change the `admin` user's password, use the M3TAL CLI:

```bash
sudo m3tal dashpass
```
Follow the prompts to set a new password.

---

## M3TAL Core Services

M3TAL operates through several interconnected components:

-   **CLI binary (`/usr/bin/m3tal`):** The unified command-line tool for all M3TAL operations.
-   **API daemon (`m3tal-api.service`):** A Go binary running as a systemd service on port `8080` (host-local). It manages Docker, the state database, and provides API routes.
-   **Dashboard container (`m3tal-dashboard`):** A Python/Flask container that provides the web interface. It communicates with the API daemon at `http://host.docker.internal:8080`.
-   **Traefik gateway (`routing-compose.yml`):** A reverse proxy container that exposes services by domain name on port `80`. It uses a file provider for dynamic routing and Docker labels for service discovery.
-   **Cloudflared (optional, `routing-compose.yml`):** An optional Cloudflare tunnel container for secure, zero-config internet access to services.

---

## Filesystem Contract

The following paths are critical to M3TAL's operation and store configuration, state, and user data:

| Path                     | Purpose                                                                |
| :----------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`| SQLite state database. Automatically created and managed by the API daemon. |
| `/opt/m3tal/stack/`      | Canonical stack directory containing Docker Compose files and Traefik configuration. |
| `/docker`                | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.               |

---

## Dashboard Access Modes

The M3TAL dashboard supports two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. This setting determines how the dashboard's Docker container is configured.

### Mode 1: Local Access (Default)

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
-   **Mechanism:** M3TAL uses `m3tal-compose.local.yml` which adds a direct port binding to the dashboard container, typically `${DASHBOARD_PORT:-8082}:8082`.
-   **Access Via:** `http://HOST_IP:8082` or `http://localhost:8082` (where `HOST_IP` is your configured `LOCAL_IP`).
-   **Requirements:** No Traefik reverse proxy is required for dashboard access in this mode.
-   **Use Case:** Ideal for initial setup, LAN-only deployments, or local testing environments.

### Mode 2: Traefik Routed Access

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
-   **Mechanism:** M3TAL uses `m3tal-compose.traefik.yml` which adds Traefik labels to the dashboard container. Traefik (if running via `m3tal up`) will discover these labels and route requests from `dash.DOMAIN` to the dashboard on its internal port `8082`.
-   **Access Via:** `http://dash.DOMAIN` (where `DOMAIN` is your configured domain name).
-   **Requirements:** Traefik must be running and properly configured for HTTP access on port 80.
-   **Use Case:** Recommended for domain-based setups, exposing services via a reverse proxy, and managing multiple services under a single domain.

---

## Port Map

The following ports are used by M3TAL core services:

| Port | Service                    | Access                                       |
| :--- | :------------------------- | :------------------------------------------- |
| 80   | Traefik HTTP entry point   | Public (if Traefik is exposed)               |
| 8080 | M3TAL API daemon (Go)      | Host-local only                              |
| 8081 | Traefik dashboard          | Host-local only (for Traefik's own dashboard) |
| 8082 | M3TAL Dashboard (container)| Direct port (local mode) or via Traefik (traefik mode) |

---

## Firewall Configuration

If you intend to expose Traefik publicly to use domain-based routing (e.g., `dash.YOUR_DOMAIN`), you must allow incoming traffic on port 80 through your firewall:

```bash
sudo ufw allow 80/tcp
```
Adjust this command if you are using a different firewall system.

---

## Service Management

The M3TAL API daemon runs as a systemd service named `m3tal-api.service`. You can manage its state using standard systemctl commands:

-   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
-   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```
-   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
```