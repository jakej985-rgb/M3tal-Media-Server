# M3TAL Ecosystem: Getting Started

Welcome to the M3TAL Ecosystem. This guide provides step-by-step instructions for a first-time setup, from prerequisites to launching your first services.

---

## Step 1: Prerequisites

Before installing M3TAL, ensure **Docker Engine and Docker Compose V2 are installed** on your system.
You can verify their installation and versions with the following commands:

```bash
docker --version && docker compose version
```

Expected output should show Docker Engine and Docker Compose V2 (e.g., `Docker Compose version v2.x.x`).

## Step 2: Install M3TAL via APT

M3TAL is distributed as a single Go binary and installed via a standard APT repository. Execute the following commands in your terminal:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

After installation, configure M3TAL using the interactive wizard. This sets up essential environment variables in `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's an explanation for key prompts:

*   **`DOMAIN`**: The base domain name for your services (e.g., `example.com`). This is used by Traefik for routing services like `dash.example.com` or `api.example.com`. If you don't have a domain, `localhost` is a suitable default.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard is directly accessible on a specific port (`8082` by default) on your host machine's IP address. No Traefik configuration is needed for dashboard access.
    *   `traefik`: The dashboard is exposed via Traefik, accessible at `dash.YOUR_DOMAIN`. This requires Traefik to be running.
*   **`DASHBOARD_PORT`**: The port on which the M3TAL Dashboard will be accessible if `DASHBOARD_EXPOSE_MODE` is set to `local`. Default is `8082`.
*   **`PUID` / `PGID`**: User ID and Group ID used for Docker container processes. This ensures containers have appropriate permissions to access host volumes. You can find your current user's PUID/PGID using `id -u` and `id -g`.
*   **`TZ`**: Timezone for your containers (e.g., `America/Denver`).
*   **`API_TOKEN`**: A secret token used for secure communication with the M3TAL API. Generate a strong, random string.
*   **`DASHBOARD_SECRET`**: A secret key used by the M3TAL Dashboard for session management. Generate a strong, random string.
*   **`ADMIN_PASSWORD`**: The default password for the `admin` user on the M3TAL Dashboard. You will change this later.

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command starts all Docker Compose stacks defined in `/docker/`. This includes the core routing services like Traefik (the reverse proxy) and Cloudflared (if configured).

```bash
m3tal up
```

This command scans `/docker/` for all `*-compose.yml` files and starts them with `docker compose`. Traefik, defined in `routing-compose.yml`, will start and bind to port 80 (HTTP) on your host machine.

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container.

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest base `m3tal-compose.yml` and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the appropriate Compose override file (`m3tal-compose.local.yml` for `local` mode, or `m3tal-compose.traefik.yml` for `traefik` mode). This ensures the dashboard is configured according to your chosen exposure method.

## Step 6: Open Browser

Now you can access the M3TAL Dashboard in your web browser:

*   **If `DASHBOARD_EXPOSE_MODE` is `local` (default):**
    Open `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or `localhost` if accessing from the server itself).
*   **If `DASHBOARD_EXPOSE_MODE` is `traefik`:**
    Open `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3). Ensure your DNS resolves `dash.YOUR_DOMAIN` to your server's IP address.

## Step 7: Log In

Upon accessing the dashboard, you will be prompted to log in.

*   **Username:** `admin`
*   **Password:** The `ADMIN_PASSWORD` you set during the `m3tal config wizard` in Step 3.

**It is strongly recommended to change the default password immediately.** You can do this using the `m3tal dashpass` command:

```bash
sudo m3tal dashpass
```

This command will prompt you to set a new password for the `admin` user, updating the `/docker/users.json` credential store.

---

## Filesystem Contract

M3TAL establishes a clear filesystem contract for its operations and user data:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.   |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains M3TAL's core compose files and Traefik configuration. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations and where user-defined `*-compose.yml` files should be placed. |
| `/docker/users.json`        | M3TAL Dashboard credential store. Managed by `m3tal dashpass`.     |

## Port Map

The M3TAL ecosystem uses the following ports on the host machine:

| Port | Service                                  | Access                                     |
| :--- | :--------------------------------------- | :----------------------------------------- |
| 80   | Traefik HTTP entry point                 | Public (if Traefik is exposed)             |
| 8080 | M3TAL API daemon (Go service)            | Host-local only (internal communication)   |
| 8081 | Traefik dashboard (read-only)            | Host-local only (e.g., `http://127.0.0.1:8081`) |
| 8082 | M3TAL Dashboard (Python/Flask container) | Direct port (local mode) or via Traefik (traefik mode) |

## Firewall Note

If Traefik is used for public access and your server has a firewall (e.g., `ufw`) enabled, you must open port 80 (and 443 for HTTPS if configured).

```bash
sudo ufw allow 80/tcp
# sudo ufw allow 443/tcp # If you configure HTTPS with Traefik
sudo ufw reload
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using standard `systemctl` commands:

*   **Check API service status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **Restart the API service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **View API service logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```

---

## Dashboard Access Modes

The M3TAL Dashboard offers two distinct access modes, configured via the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

### Mode 1: `local` (Default)

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
*   **Mechanism:** M3TAL uses `m3tal-compose.local.yml` as an override, which adds a direct port binding to the dashboard container (e.g., `${DASHBOARD_PORT:-8082}:8082`). This exposes the dashboard container's internal port 8082 directly on your host machine's specified port (default 8082).
*   **Access via:** `http://YOUR_HOST_IP:8082` or `http://localhost:8082`
*   **Requirements:** No Traefik gateway is required for dashboard access in this mode.
*   **Best for:** LAN-only setups, first-time users, or local testing where a domain is not yet configured.

### Mode 2: `traefik`

*   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
*   **Mechanism:** M3TAL uses `m3tal-compose.traefik.yml` as an override. This file adds specific Traefik labels to the dashboard container. Traefik (running via `m3tal up`) discovers these labels and automatically configures a route (e.g., `dash.${DOMAIN}`) to forward requests to the dashboard container's internal port 8082.
*   **Access via:** `http://dash.YOUR_DOMAIN` (Requires Traefik to be running via `m3tal up` and proper DNS resolution for `dash.YOUR_DOMAIN`).
*   **Best for:** Domain-based setups where you want all services, including the dashboard, to be accessible via a reverse proxy and domain names.

---

## Docker / Compose Runtime

M3TAL leverages **Docker Engine and Docker Compose V2** as its core runtime. These are hard dependencies for the system.

*   The `m3tal up` command orchestrates all Docker Compose applications (stacks) found in the `/docker/` directory. When you run `m3tal up`, M3TAL executes `docker compose --project-directory /docker/ up -d`, effectively starting all `*-compose.yml` files in that location.
*   The `m3tal dash up` command is a specialized operation for the M3TAL Dashboard. It ensures the dashboard is always deployed correctly by:
    1.  Downloading the latest `m3tal-compose.yml` and its mode-specific overrides (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
    2.  Reading the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
    3.  Starting the dashboard container using the base compose file combined with the appropriate override for the configured mode.
*   **User Stacks:** To deploy your own Docker Compose applications, simply place your `*-compose.yml` files into the `/docker/` directory.

---

## Deployment Lifecycle — Day 2 Operations

Once M3TAL is set up, deploying new services or "stacks" is straightforward:

1.  **Place your compose file:** Create or place your Docker Compose file (e.g., `my-app-compose.yml`) into the `/docker/` directory.
2.  **Set environment variables:** If your stack requires specific environment variables (e.g., ports, paths), ensure they are defined in `/etc/m3tal/.env`. You can use `m3tal config wizard` or `m3tal config set KEY value` to manage these.
3.  **Start all stacks:** Run `m3tal up` to start your new stack along with any existing ones. M3TAL will automatically detect the new compose file.

---

## Traefik Routing Architecture

Traefik acts as the central reverse proxy for your M3TAL ecosystem. It is deployed as a Docker container via `routing-compose.yml` (managed by `m3tal up`).

*   **Entry Point:** Traefik binds to port 80 on the host machine, serving as the primary HTTP entry point for all web traffic.
*   **Service Discovery:** It automatically discovers services by inspecting Docker containers for specific labels.
*   **Dynamic Configuration:** Traefik loads additional dynamic configurations from the `/etc/traefik/dynamic` directory (which symlinks to `/docker/dynamic`). This allows for hot-reloading routing rules without restarting Traefik.
*   **API Routing:** Traefik routes requests for `api.YOUR_DOMAIN` directly to the M3TAL API daemon running internally on `http://host.docker.internal:8080`. This is configured via `dynamic/api.yml`.
*   **Dashboard Routing:** If `DASHBOARD_EXPOSE_MODE` is set to `traefik`, it routes `dash.YOUR_DOMAIN` requests to the M3TAL Dashboard container (internal port 8082) using Traefik labels defined in `m3tal-compose.traefik.yml`.