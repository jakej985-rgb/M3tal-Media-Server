# M3TAL - Get Started

This guide will walk you through setting up M3TAL on your system.

## Step 1: Prerequisites

Before you begin, ensure you have the following installed:
- **Docker Engine**
- **Docker Compose V2**

You can verify your installation with the following commands:

```bash
docker --version
docker compose version
```

If these commands do not output version information, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

Install the M3TAL CLI by adding our repository and installing the package.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

Initialize M3TAL and configure essential settings using the built-in wizard.

```bash
sudo m3tal config wizard
```

You will be prompted with several questions. Here's a brief explanation of each:

*   **Dashboard Port:** The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **Dashboard Expose Mode:** Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via a port (default). Recommended for initial setup and LAN access.
    *   `traefik`: Exposes the dashboard through Traefik, allowing access via a domain name. Requires Traefik to be running.
*   **HTTP Port:** The port for the M3TAL API daemon. The default is `8080`.
*   **State Directory:** The directory where M3TAL stores its state (e.g., database files). The default is `./state`.
*   **Log Level:** The verbosity of M3TAL's logs. Common options include `info`, `debug`, `warn`, `error`.
*   **Dashboard Secret:** A secret key used for securing dashboard sessions. **Change this from the default for security.**
*   **API Token:** A token for authenticating API requests. **Change this from the default for security.**
*   **Admin Password:** The password for the dashboard's administrator account. **Change this from the default for security.**
*   **Network Name:** The Docker network name M3TAL will use. The default is `m3tal`.
*   **Local IP:** The local IP address of your host machine.
*   **Domain:** The domain name to use for Traefik routing (if `traefik` expose mode is selected). Defaults to `localhost`.
*   **VPN User/Password:** If you intend to use VPN functionality.
*   **Base Storage Path:** The root directory for M3TAL's data storage. Defaults to `./data`.
*   **Media Path:** A sub-directory within the Base Storage Path for media files.
*   **Config Path:** A sub-directory within the Base Storage Path for configuration files.
*   **Downloads Path:** A sub-directory within the Base Storage Path for download files.
*   **PUID/PGID:** The User and Group ID for container processes. Defaults to `1000`.
*   **TZ:** Your timezone.
*   **Traefik Web Port/Web HTTPS Port/Dashboard Port:** Ports for Traefik. Defaults are `80`, `443`, and `8080` respectively.
*   **Debug Mode:** Enable or disable debug mode.
*   **Metrics Enabled:** Enable or disable metrics collection.

## Step 4: Start the Routing Stack (Traefik)

This command starts the core M3TAL routing stack, including Traefik, which acts as a reverse proxy.

```bash
m3tal up
```

This command will:
*   Locate all `.yml` files within the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`).
*   Use Docker Compose V2 to start the services defined in these compose files.

## Step 5: Start the Dashboard

This command specifically pulls the M3TAL dashboard image and starts its container.

```bash
m3tal dash up
```

This command will:
1.  Download the necessary dashboard compose files from GitHub.
2.  Read your `DASHBOARD_EXPOSE_MODE` setting from `/etc/m3tal/.env`.
3.  Start the `m3tal-dashboard` container using the appropriate compose override for your chosen expose mode.

## Step 6: Access the Dashboard

Open your web browser and navigate to the dashboard's address:

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (default): `http://YOUR_IP:8082` (replace `YOUR_IP` with your server's IP address, or `localhost` if on the same machine).
*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik`: `http://dash.DOMAIN` (replace `DOMAIN` with the domain you configured in Step 3, and ensure Traefik is running via `m3tal up`).

## Step 7: Log In

You will be presented with the M3TAL login screen.

*   **Default Credentials:**
    *   Username: `admin`
    *   Password: `admin_pass` (or whatever you set during the wizard).

**It is highly recommended to change the default password immediately after logging in.** You can do this using the following command:

```bash
sudo m3tal dashpass
```

This will prompt you to enter a new password for the dashboard.

---

## Filesystem Contract

M3TAL utilizes a specific directory structure and configuration files:

| Path                       | Purpose                                                                    | Notes                                                                                              |
| :------------------------- | :------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`          | Primary M3TAL configuration file.                                          | Managed by `m3tal config wizard` and manual edits.                                                 |
| `/var/lib/m3tal/state.db`  | SQLite state database for the M3TAL API daemon.                            | Auto-created by the API daemon on first run.                                                       |
| `/opt/m3tal/stack/`        | Canonical directory for M3TAL's Docker Compose files and Traefik configuration. | Contains base compose files and dynamic routing configurations.                                      |
| `/docker`                  | Symlink pointing to `/opt/m3tal/stack/`.                                   | This is the user-facing path for all M3TAL stack operations (`m3tal up`, adding new compose files). |
| `/docker/users.json`       | Stores dashboard user credentials.                                         | Managed by `m3tal dashpass`.                                                                       |

## Port Map

The following ports are used by M3TAL services:

| Port   | Service                 | Access Method                                    |
| :----- | :---------------------- | :----------------------------------------------- |
| 80     | Traefik (HTTP Entrypoint) | Public (if `DASHBOARD_EXPOSE_MODE=traefik`)      |
| 8080   | M3TAL API Daemon (Go)   | Host-local only                                  |
| 8082   | M3TAL Dashboard         | Direct port (if `DASHBOARD_EXPOSE_MODE=local`)   |
|        |                         | Via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

## Firewall Note

If you are exposing M3TAL services to the public internet using Traefik (i.e., `DASHBOARD_EXPOSE_MODE=traefik` and Traefik is configured to listen on port 80), ensure your firewall allows incoming traffic on port 80.

If you are using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl` and view its logs with `journalctl`.

**Check status:**

```bash
systemctl status m3tal-api
```

**View logs in real-time:**

```bash
journalctl -u m3tal-api -f
```