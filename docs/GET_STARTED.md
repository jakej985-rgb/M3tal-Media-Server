# GET_STARTED.md

This guide provides step-by-step instructions for installing and setting up M3TAL for the first time.

---

## Step 1: Prerequisites

Before you begin, ensure you have the following software installed on your system:

*   **Docker Engine:** The containerization platform.
*   **Docker Compose V2:** The tool for defining and running multi-container Docker applications.

Verify your installation by running:

```bash
docker --version && docker compose version
```

If either command fails, please refer to the official Docker documentation for installation instructions.

---

## Step 2: Install M3TAL via APT

M3TAL is distributed via an APT repository. Execute the following commands to add the repository and install the M3TAL package:

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

The `m3tal config wizard` command will guide you through setting up essential configuration parameters for M3TAL.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a breakdown of what each prompt means:

*   **`Enter your desired domain name (e.g., localhost, your.domain.com):`** This is the domain name that M3TAL will use to expose its services. For local testing, `localhost` is suitable. If you plan to access M3TAL from outside your local network, you'll need a publicly resolvable domain.
*   **`Enter the timezone for your server (e.g., America/Denver):`** This sets the timezone for your M3TAL services, affecting log timestamps and scheduling.
*   **`Enter the desired port for the M3TAL Dashboard (default: 8082):`** This is the port on which the M3TAL Dashboard will be accessible. The default is `8082`.
*   **`Select the Dashboard expose mode (local/traefik):`**
    *   `local`: Exposes the dashboard directly via a port binding. This is the default and recommended for initial setup and LAN-only access.
    *   `traefik`: Configures the dashboard to be accessible through Traefik (if Traefik is running). This is useful when you have multiple services behind a reverse proxy.
*   **`Set a strong password for the Dashboard admin user:`** This is the password you will use to log into the M3TAL Dashboard. It is highly recommended to choose a strong, unique password.
*   **`Set a strong secret for the Dashboard (used for session security):`** This is a secret key used for securing dashboard sessions. Generate a strong, random string.
*   **`Set a strong API token for programmatic access:`** This token is used by other services or scripts to authenticate with the M3TAL API.
*   **`Enter your user ID (PUID, default: 1000):`** The user ID to run Docker containers as. Typically `1000` for the first user on a Linux system.
*   **`Enter your group ID (PGID, default: 1000):`** The group ID to run Docker containers as. Typically `1000` for the first user on a Linux system.

---

## Step 4: Start the Routing Stack (Traefik)

The M3TAL CLI manages Docker Compose files located in `/opt/m3tal/stack/`. The `m3tal up` command orchestrates all Docker Compose files found within the `/docker/` symlink, which points to `/opt/m3tal/stack/`. This includes the Traefik reverse proxy.

To start Traefik and other core M3TAL services, run:

```bash
m3tal up
```

This command will:
*   Read all `*.compose.yml` files in `/docker/`.
*   Download necessary Docker images if they are not present locally.
*   Start the containers defined in these compose files, including Traefik.
*   Ensure services are connected to the appropriate Docker networks.

---

## Step 5: Start the Dashboard

The M3TAL Dashboard provides a web interface to manage your M3TAL services. To start the dashboard container:

```bash
m3tal dash up
```

This command will:
*   Download the latest `m3tal-compose.yml` and relevant override files (e.g., `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) from the M3TAL repository.
*   Read the `DASHBOARD_EXPOSE_MODE` from your `/etc/m3tal/.env` file.
*   Start the `m3tal-dashboard` container, applying the correct configuration based on the chosen `DASHBOARD_EXPOSE_MODE`.

---

## Step 6: Access the M3TAL Dashboard

You can now access the M3TAL Dashboard through your web browser. The access method depends on the `DASHBOARD_EXPOSE_MODE` configured in Step 3:

*   **If `DASHBOARD_EXPOSE_MODE=local` (default):**
    Open your web browser and navigate to:
    `http://YOUR_IP:8082`
    (Replace `YOUR_IP` with your server's local IP address or `localhost` if accessing from the same machine.)

*   **If `DASHBOARD_EXPOSE_MODE=traefik`:**
    Open your web browser and navigate to:
    `http://dash.DOMAIN`
    (Replace `DOMAIN` with the domain name you configured in Step 3. Traefik must be running for this to work.)

---

## Step 7: Log In

Upon accessing the dashboard, you will be presented with a login screen.

*   **Username:** `admin`
*   **Password:** The password you set during the `sudo m3tal config wizard` step.

### Changing the Dashboard Password

If you need to change the dashboard password after initial setup, use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password and will update the dashboard's credentials securely.

---

## Filesystem Contract

M3TAL relies on a specific directory structure and configuration files for its operation:

| Path                      | Purpose                                                                                                   |
| :------------------------ | :-------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`         | The primary M3TAL configuration file. This file is managed by the `m3tal config wizard` and `m3tal config set` commands. |
| `/var/lib/m3tal/state.db` | The SQLite state database used by the M3TAL API daemon to store its operational state and service information. |
| `/docker`                 | A symbolic link pointing to `/opt/m3tal/stack/`. This is the user-facing directory for all M3TAL stack operations, including Compose files. |
| `/docker/users.json`      | Stores dashboard user credentials. Managed by `m3tal dashpass`.                                           |

---

## Port Map

The following ports are utilized by M3TAL and its components:

| Port | Service / Component         | Access Method                                    |
| :--- | :-------------------------- | :----------------------------------------------- |
| 80   | Traefik (HTTP Entrypoint)   | Public (when Traefik is exposed to the internet) |
| 8080 | M3TAL API Daemon (Go)       | Host-local (accessed by other services)          |
| 8081 | Traefik Dashboard           | Host-local only (for debugging Traefik)          |
| 8082 | M3TAL Dashboard             | Direct port binding (local mode) or via Traefik (traefik mode) |

---

## Firewall Note

If you intend to expose Traefik to the internet (e.g., for accessing services via domain names from outside your local network), ensure that port 80 is open in your firewall.

```bash
sudo ufw allow 80
```

---

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using `systemctl` and `journalctl`.

*   **Check the status of the M3TAL API service:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View live logs of the M3TAL API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```