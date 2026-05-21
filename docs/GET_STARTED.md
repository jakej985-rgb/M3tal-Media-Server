# Getting Started with M3TAL

This guide will walk you through the installation and initial setup of the M3TAL ecosystem.

## Table of Contents

1.  [Prerequisites](#step-1-prerequisites)
2.  [Install M3TAL via APT](#step-2-install-m3tal-via-apt)
3.  [Run the Configuration Wizard](#step-3-run-the-configuration-wizard)
4.  [Start the Routing Stack (Traefik)](#step-4-start-the-routing-stack-traefik)
5.  [Start the Dashboard](#step-5-start-the-dashboard)
6.  [Access the Dashboard](#step-6-access-the-dashboard)
7.  [Log In](#step-7-log-in)
8.  [Filesystem Contract](#filesystem-contract)
9.  [Port Table](#port-table)
10. [Firewall Note](#firewall-note)
11. [Service Management](#service-management)

---

## Step 1: Prerequisites

Before proceeding, ensure you have the following installed on your system:

*   **Docker Engine**
*   **Docker Compose V2**

You can verify their installation by running:

```bash
docker --version && docker compose version
```

---

## Step 2: Install M3TAL via APT

Execute the following commands to add the M3TAL repository and install the CLI binary:

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

This wizard will guide you through essential initial configuration.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a breakdown of what each prompt means:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via its assigned port (`DASHBOARD_PORT`). This is ideal for LAN-only access and initial setup.
    *   `traefik`: Exposes the dashboard through Traefik. This requires Traefik to be running and configured with a domain.
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state database (`state.db`). The default is `./state` (relative to where `m3tal` commands are run if not using absolute paths, but typically managed by systemd which will use a default location like `/var/lib/m3tal`).
*   **`LOG_LEVEL`**: Sets the verbosity of logs. Options include `debug`, `info`, `warn`, `error`. `info` is generally recommended for normal operation.
*   **`DASHBOARD_SECRET`**: A secret key for securing the dashboard. **Change this from the default immediately.**
*   **`API_TOKEN`**: A token for authenticating with the M3TAL API. **Change this from the default immediately.**
*   **`ADMIN_PASSWORD`**: The password for the default dashboard administrator user. **Change this from the default immediately.**
*   **`NETWORK_NAME`**: The name of the Docker network M3TAL will use. The default is `m3tal`.
*   **`LOCAL_IP`**: Your server's local IP address. This is often detected automatically.
*   **`DOMAIN`**: The domain name you intend to use for accessing M3TAL services. Defaults to `localhost`.
*   **`PUID` / `PGID`**: The User ID (UID) and Group ID (GID) to run Docker containers as. Often defaults to `1000` for the first user.
*   **`TZ`**: Your server's timezone.
*   **`TRAEFIK_WEB_PORT`**: The port Traefik will use for HTTP traffic. Default is `80`.
*   **`TRAEFIK_WEBHTTPS_PORT`**: The port Traefik will use for HTTPS traffic. Default is `443`.

The wizard will save your choices to `/etc/m3tal/.env`.

---

## Step 4: Start the Routing Stack (Traefik)

This command initializes and starts all Docker Compose files located in the `/docker/` directory. This includes Traefik, which acts as your primary ingress point for services.

```bash
m3tal up
```

This command will:
*   Read all `*.yml` files within `/docker/`.
*   Construct a unified Docker Compose command.
*   Start the containers defined in these files.

---

## Step 5: Start the Dashboard

This command specifically pulls the M3TAL dashboard image and starts its container, respecting the `DASHBOARD_EXPOSE_MODE` setting from your configuration.

```bash
m3tal dash up
```

This command will:
1.  Download the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from the M3TAL repository.
2.  Read the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container using the appropriate compose file override for your chosen mode.

---

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard URL:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default), access it at: `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or `localhost` if accessing from the server itself).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik` and Traefik is running, access it at: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured).

---

## Step 7: Log In

Upon first access, you will be presented with a login screen.

*   **Default Username:** `admin`
*   **Default Password:** `admin_pass` (This is the default from the `ADMIN_PASSWORD` environment variable. **You should change this immediately.**)

To change the administrator password, use the following command:

```bash
sudo m3tal dashpass admin new_password
```

Replace `new_password` with your desired strong password.

---

## Filesystem Contract

M3TAL utilizes a specific filesystem structure for configuration and state management.

| Path                  | Purpose                                                                                                 |
| :-------------------- | :------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`     | The primary configuration file. Managed by the `m3tal config wizard` and `m3tal config set` commands.    |
| `/var/lib/m3tal/state.db` | SQLite state database. Auto-created by the API daemon. Stores operational state and configurations. |
| `/docker`             | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. |
| `/docker/users.json`  | Stores dashboard user credentials. Managed by `m3tal dashpass`.                                         |
| `/opt/m3tal/stack/`   | The canonical directory containing M3TAL's core Docker Compose files, Traefik configurations, etc.    |

---

## Port Table

The following ports are used by M3TAL and its components.

| Port | Service           | Access                                     |
| :--- | :---------------- | :----------------------------------------- |
| 80   | Traefik (HTTP)    | Public (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API daemon  | Host-local only                            |
| 8081 | Traefik Dashboard | Host-local only (`127.0.0.1:8081`)         |
| 8082 | M3TAL Dashboard   | Direct port (local mode) or via Traefik     |

---

## Firewall Note

If you are exposing Traefik to the public internet (e.g., `DASHBOARD_EXPOSE_MODE=traefik` and your server has a public IP), you must allow traffic on port 80 through your firewall. If you are using `ufw`:

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status:**
    ```bash
    systemctl status m3tal-api
    ```
*   **View logs in real-time:**
    ```bash
    journalctl -u m3tal-api -f
    ```
*   **Restart the service:**
    ```bash
    sudo systemctl restart m3tal-api
    ```