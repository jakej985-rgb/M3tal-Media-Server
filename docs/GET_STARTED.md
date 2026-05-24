```markdown
# M3TAL Getting Started Guide

This guide provides a complete, step-by-step setup for first-time M3TAL users. It covers installation, initial configuration, and how to access the M3TAL Dashboard.

---

## Step 1: Prerequisites

Docker Engine and Docker Compose V2 must be installed on your system.
Verify your Docker installation by running:

```bash
docker --version && docker compose version
```

Example expected output:
```
Docker version 24.0.6, build d680c7dd1c
Docker Compose version v2.21.0
```

If Docker or Docker Compose V2 are not installed, please refer to the official Docker documentation for your operating system.

---

## Step 2: Install M3TAL via APT

M3TAL is installed as a single CLI binary via the APT package manager. Follow these commands to add the repository and install M3TAL:

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

After installation, run the M3TAL configuration wizard to set up essential environment variables. This wizard will guide you through setting up `/etc/m3tal/.env`.

```bash
sudo m3tal config wizard
```

The wizard will prompt you for several values. Here's an explanation for the common prompts:

*   **`DOMAIN`**: The base domain for your services (e.g., `example.com`). This is used for Traefik routing. Defaults to `localhost`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the M3TAL Dashboard is exposed.
    *   `local` (default): The dashboard will be directly accessible via `http://YOUR_IP:8082`. No Traefik configuration is needed for the dashboard itself. Ideal for local networks or first-time setup.
    *   `traefik`: The dashboard will be accessible via `http://dash.YOUR_DOMAIN` through the Traefik reverse proxy. Requires Traefik to be running.
*   **`PUID` / `PGID`**: User ID and Group ID for containers to ensure proper file permissions (e.g., `1000` for the first non-root user). You can find your current user's IDs with `id -u` and `id -g`.
*   **`TZ`**: Your timezone (e.g., `America/New_York`).
*   **`DASHBOARD_SECRET`**: A secret key used by the M3TAL Dashboard for session management. Generate a strong, random string.
*   **`API_TOKEN`**: A token for secure communication with the M3TAL API. Generate a strong, random string.
*   **`ADMIN_PASSWORD`**: The default password for the `admin` user on the M3TAL Dashboard. **Change this immediately after setup.**

---

## Step 4: Start the Routing Stack (Traefik)

The `m3tal up` command brings up all Docker Compose stacks found in the `/docker/` directory. This includes the Traefik reverse proxy, which is part of the `routing-compose.yml` stack.

```bash
m3tal up
```

This command will initialize and start all services defined in your M3TAL ecosystem, including the core routing infrastructure.

---

## Step 5: Start the Dashboard

The `m3tal dash up` command specifically manages the M3TAL Dashboard container. It pulls the latest dashboard image from the registry and starts the container using the appropriate configuration based on your `DASHBOARD_EXPOSE_MODE` setting from Step 3.

```bash
m3tal dash up
```

This command automatically selects either `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml` to apply the correct exposure settings for the dashboard.

---

## Step 6: Access the Dashboard

How you access the M3TAL Dashboard depends on the `DASHBOARD_EXPOSE_MODE` configured in Step 3.

### Mode 1: Local Access (`DASHBOARD_EXPOSE_MODE=local`)

If you set `DASHBOARD_EXPOSE_MODE` to `local`, the dashboard is exposed directly on port 8082 of your host machine.

Open your web browser and navigate to:
`http://YOUR_SERVER_IP:8082` (replace `YOUR_SERVER_IP` with the IP address of your M3TAL host)
or `http://localhost:8082` if accessing from the host itself.

### Mode 2: Traefik Access (`DASHBOARD_EXPOSE_MODE=traefik`)

If you set `DASHBOARD_EXPOSE_MODE` to `traefik`, the dashboard is routed through Traefik using a subdomain.

Open your web browser and navigate to:
`http://dash.YOUR_DOMAIN` (replace `YOUR_DOMAIN` with the `DOMAIN` you configured in Step 3).

---

## Step 7: Log in to the Dashboard

The default login credentials for the M3TAL Dashboard are:

*   **Username**: `admin`
*   **Password**: The value you set for `ADMIN_PASSWORD` during the `m3tal config wizard` (default is `admin_pass` if not changed).

**It is highly recommended to change the default password immediately.** You can do this using the CLI command:

```bash
sudo m3tal dashpass
```

Follow the prompts to set a new password for the `admin` user.

---

## Filesystem Contract

M3TAL relies on specific file and directory locations for its operation and configuration.

| Path                     | Purpose                                                            |
| :----------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file, managed by `m3tal config wizard`.      |
| `/var/lib/m3tal/state.db`| SQLite state database, auto-created and managed by the API daemon. |
| `/opt/m3tal/stack/`      | Canonical directory containing M3TAL's core compose files.         |
| `/docker`                | **Symlink** to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your custom `*-compose.yml` files here. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.           |

---

## Port Table

These are the default ports used by M3TAL components. Ensure these ports are not in use by other services on your host.

| Port | Service               | Access                                                 |
| :--- | :-------------------- | :----------------------------------------------------- |
| 80   | Traefik HTTP entry point | Public (if `DASHBOARD_EXPOSE_MODE=traefik` and Traefik is running) |
| 8080 | M3TAL API daemon (Go) | Host-local only (internal communication)               |
| 8081 | Traefik dashboard     | Host-local only (for Traefik's own dashboard)          |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

---

## Firewall Note

If you have a firewall enabled (e.g., `ufw`) and have configured Traefik to be exposed publicly (e.g., via `DOMAIN` in the wizard), you will need to allow external access to port 80 (and 443 for HTTPS if configured):

```bash
sudo ufw allow 80/tcp
# sudo ufw allow 443/tcp # If you configure HTTPS for Traefik
sudo ufw enable # if ufw is not already enabled
```

---

## Service Management

The M3TAL API daemon (`m3tal-api`) runs as a systemd service. You can manage it using standard `systemctl` commands:

*   **Check service status:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View live logs:**
    ```bash
    journalctl -u m3tal-api -f
    ```

*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```

---
```