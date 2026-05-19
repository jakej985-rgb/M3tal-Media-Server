# Getting Started with M3TAL

This guide provides a complete, step-by-step process for setting up M3TAL for first-time users.

---

## Step 1: Prerequisites

Before installing M3TAL, ensure that **Docker Engine** and **Docker Compose V2** are installed on your system. M3TAL relies on these tools for container orchestration.

Verify your Docker and Docker Compose versions:

```bash
docker --version && docker compose version
```

---

## Step 2: Install M3TAL via APT

M3TAL is distributed via an APT repository. Execute the following three commands to add the repository and install the `m3tal` CLI:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## Step 3: Run the Configuration Wizard

After installation, configure M3TAL's primary settings using the interactive wizard. This will generate your `/etc/m3tal/.env` file.

```bash
sudo m3tal config wizard
```

During the wizard, you will be prompted for several key settings:

-   **DOMAIN**: Sets the base domain for services exposed via Traefik (e.g., `yourdomain.com`). If you plan to access services like the dashboard or API via subdomains (e.g., `dash.yourdomain.com`, `api.yourdomain.com`), provide your domain here. The default is `localhost`.
-   **DASHBOARD_EXPOSE_MODE**: This critical setting determines how the M3TAL Dashboard is made accessible.
    -   `local` (default): The dashboard will bind directly to `HOST_IP:8082`. It's suitable for local network access and does not require Traefik.
    -   `traefik`: The dashboard will be routed through Traefik, accessible via `http://dash.YOUR_DOMAIN`. This mode requires Traefik to be running and correctly configured with a domain.
-   **DASHBOARD_PORT**: If `DASHBOARD_EXPOSE_MODE` is set to `local`, this is the host port that the M3TAL Dashboard will be accessible on. The default is `8082`.
-   **PUID** / **PGID**: These specify the User ID and Group ID that containers will run as, ensuring correct file permissions for mounted volumes. Default is `1000` (often the first non-root user).
-   **TZ**: Sets the timezone for M3TAL containers (e.g., `America/Denver`).
-   You will also be prompted for or informed about `DASHBOARD_SECRET`, `API_TOKEN`, and `ADMIN_PASSWORD`. It is highly recommended to change these default credentials immediately for security.

---

## Step 4: Start the Routing Stack (Traefik)

M3TAL uses Docker Compose V2 for managing container stacks. The `m3tal up` command starts all Docker Compose files (named `*-compose.yml`) located in the `/docker/` directory.

The default routing stack (`routing-compose.yml`) includes Traefik, which acts as a reverse proxy. It will bind to port 80 on your host, allowing services to be exposed by domain name.

```bash
m3tal up
```

This command starts Traefik and any other services defined in `/docker/*.yml` files.

---

## Step 5: Start the Dashboard

The M3TAL Dashboard is a web interface for managing your ecosystem. Use the `m3tal dash up` command to start it:

```bash
m3tal dash up
```

This command performs the following actions:
1.  Downloads the latest dashboard Compose files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`).
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the appropriate Compose override file based on your chosen exposure mode.

---

## Step 6: Access the Dashboard

The method to access the dashboard depends on the `DASHBOARD_EXPOSE_MODE` you configured in Step 3.

### Mode 1: Local Access (`DASHBOARD_EXPOSE_MODE=local`)

If you selected `local` mode, the dashboard is directly exposed on a specific port.

-   **Access via**: `http://YOUR_IP:8082` (replace `YOUR_IP` with the actual IP address of your server or use `localhost`).
-   **Requirements**: No Traefik required. The `DASHBOARD_PORT` (default 8082) must be open on your server's firewall if accessing from another machine on your local network.

### Mode 2: Traefik Access (`DASHBOARD_EXPOSE_MODE=traefik`)

If you selected `traefik` mode, the dashboard is routed through Traefik using a subdomain.

-   **Access via**: `http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the domain you configured in Step 3).
-   **Requirements**:
    -   Traefik must be running (`m3tal up`).
    -   Port 80 must be open on your server's firewall.
    -   You must have a DNS `A` record configured for `dash.YOUR_DOMAIN` pointing to your server's IP address.

---

## Step 7: Log In to the Dashboard

Upon accessing the dashboard, you will be prompted for credentials.

-   **Default Username**: `admin`
-   **Default Password**: `admin_pass` (This is the default value of the `ADMIN_PASSWORD` variable, set during the `m3tal config wizard`).

**Important**: Change the default password immediately for security.

To change the dashboard password, use the `m3tal dashpass` command:

```bash
sudo m3tal dashpass
```

---

## Filesystem Contract

M3TAL uses a defined filesystem layout for its configuration and data:

| Path                        | Purpose                                                      |
| :-------------------------- | :----------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`   | SQLite state database for the M3TAL API daemon. Auto-created. |
| `/opt/m3tal/stack/`         | Canonical directory for M3TAL's core Docker Compose files and Traefik configuration. |
| `/docker`                   | A symbolic link that points to `/opt/m3tal/stack/`. This is the user-facing path for placing and managing Docker Compose files for your services. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.     |

---

## Port Map

These are the primary ports used by M3TAL components:

| Port | Service                                  | Access                                     |
| :--- | :--------------------------------------- | :----------------------------------------- |
| 80   | Traefik HTTP entry point                 | Public (if Traefik is exposed)             |
| 8080 | M3TAL API daemon (Go service)            | Host-local only                            |
| 8081 | Traefik dashboard (for monitoring Traefik) | Host-local only                            |
| 8082 | M3TAL Dashboard (Python/Flask container) | Direct port (local mode) or via Traefik (traefik mode) |

---

## Firewall Configuration

If you are exposing M3TAL services to your local network or the internet, you may need to open ports in your system's firewall.

For example, using `ufw`:

-   **To allow HTTP traffic (for Traefik on port 80):**
    ```bash
    sudo ufw allow 80/tcp
    ```
-   **To allow direct dashboard access (for `local` expose mode on port 8082):**
    ```bash
    sudo ufw allow 8082/tcp
    ```

---

## Service Management (M3TAL API Daemon)

The core M3TAL API daemon runs as a systemd service (`m3tal-api.service`). You can manage it using standard `systemctl` commands:

-   **Check the status of the API daemon:**
    ```bash
    systemctl status m3tal-api
    ```
-   **Restart the API daemon:**
    ```bash
    systemctl restart m3tal-api
    ```
-   **View real-time logs for the API daemon:**
    ```bash
    journalctl -u m3tal-api -f
    ```

---

## Day 2 Operations: Deploying New Stacks

M3TAL simplifies the deployment of additional Docker Compose-based services.

1.  **Place your Docker Compose file**: Create or copy your `*-compose.yml` file (e.g., `my-app-compose.yml`) into the `/docker/` directory.
2.  **Configure environment variables**: If your stack requires specific environment variables, set them in `/etc/m3tal/.env` using `m3tal config set KEY value` or by running `m3tal config wizard` again.
3.  **Start all stacks**: Run `m3tal up` to detect and start all new and existing Docker Compose stacks in `/docker/`.