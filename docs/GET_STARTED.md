# M3TAL System - Getting Started

This guide will walk you through the initial setup and configuration of the M3TAL system.

## Step 1: Prerequisites

Before proceeding, ensure you have the following software installed on your system.

Docker Engine and Docker Compose V2 must be installed.

Verify your installation:

```bash
docker --version && docker compose version
```

If these commands do not return version information, please refer to the official Docker documentation for installation instructions.

## Step 2: Install M3TAL via APT

The M3TAL CLI binary is distributed via an APT repository. Execute the following commands to add the repository and install M3TAL:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Step 3: Run the Configuration Wizard

The `m3tal config wizard` command will guide you through setting up essential configuration parameters. You will be prompted for the following:

*   **`DASHBOARD_PORT`**: The port on which the M3TAL dashboard will be accessible. The default is `8082`.
*   **`DASHBOARD_EXPOSE_MODE`**: Determines how the dashboard is exposed.
    *   `local`: Exposes the dashboard directly via its port (default, recommended for initial setup).
    *   `traefik`: Exposes the dashboard through the Traefik reverse proxy (requires Traefik to be running).
*   **`HTTP_PORT`**: The port for the M3TAL API daemon. The default is `8080`.
*   **`STATE_DIR`**: The directory where M3TAL stores its state information, including the database. The default is `./state` relative to the M3TAL configuration directory.
*   **`LOG_LEVEL`**: Sets the verbosity of M3TAL's logs. Common values include `debug`, `info`, `warn`, `error`. The default is `info`.
*   **`DASHBOARD_SECRET`**: A secret key used for securing the dashboard. **It is highly recommended to change this from the default value.**
*   **`API_TOKEN`**: An API token for programmatic access to M3TAL. **It is highly recommended to change this from the default value.**
*   **`ADMIN_PASSWORD`**: The password for the dashboard administrator. **It is highly recommended to change this from the default value.**
*   **`DOMAIN`**: The domain name to be used for accessing services via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`). Defaults to `localhost`.

Run the wizard:

```bash
sudo m3tal config wizard
```

Follow the on-screen prompts to configure your system.

## Step 4: Start the Routing Stack (Traefik)

This command starts the core routing infrastructure, primarily Traefik, which acts as a reverse proxy for your services. It orchestrates all Docker Compose files found in the `/docker/` directory.

```bash
m3tal up
```

This command initiates all defined services managed by Docker Compose in the `/docker/` directory, including Traefik and any other stacks you have added.

## Step 5: Start the Dashboard

This command specifically pulls the latest dashboard image and starts the M3TAL dashboard container.

```bash
m3tal dash up
```

This action ensures you have the most up-to-date dashboard version and launches its container. The specific compose file used depends on your `DASHBOARD_EXPOSE_MODE` setting from the configuration wizard.

## Step 6: Access the M3TAL Dashboard

Open your web browser and navigate to the dashboard's URL.

*   If `DASHBOARD_EXPOSE_MODE` is set to `local` (the default):
    `http://YOUR_IP:8082`
    (Replace `YOUR_IP` with your server's IP address, or `http://localhost:8082` if you are accessing from the same machine).

*   If `DASHBOARD_EXPOSE_MODE` is set to `traefik` and Traefik is configured with a domain:
    `http://dash.DOMAIN`
    (Replace `DOMAIN` with the domain you configured, e.g., `http://dash.localhost` if you used the default).

## Step 7: Log In

You will be presented with the M3TAL dashboard login screen.

*   **Username**: `admin`
*   **Password**: The `ADMIN_PASSWORD` you set during the configuration wizard.

If you need to change the dashboard password after initial setup, use the following command:

```bash
sudo m3tal dashpass
```

This command will prompt you to enter a new password for the `admin` user.

---

## Filesystem Contract

The following file system locations are critical for M3TAL's operation:

| Path                 | Purpose                                                                          |
| :------------------- | :------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`    | The primary configuration file for M3TAL. Modified by `m3tal config wizard`.    |
| `/var/lib/m3tal/`    | Contains runtime data and the state database (`state.db`).                       |
| `/opt/m3tal/stack/`  | The canonical directory for M3TAL's Docker Compose files and Traefik configuration. |
| `/docker`            | A symbolic link pointing to `/opt/m3tal/stack/`. User-facing path for stacks.    |
| `/docker/users.json` | Stores dashboard credentials. Managed by `m3tal dashpass`.                       |

## Port Usage

The following ports are used by M3TAL and its components:

| Port | Service      | Access                                    |
| :--- | :----------- | :---------------------------------------- |
| 80   | Traefik      | Public (if Traefik is exposed)            |
| 8080 | M3TAL API    | Host-local (accessible by other services) |
| 8081 | Traefik API  | Host-local only                           |
| 8082 | M3TAL Dashboard | Host-local or Public (via Traefik)        |

## Firewall Configuration

If you are exposing Traefik to the public internet (e.g., on port 80), ensure your firewall allows incoming traffic on that port. For systems using `ufw`:

```bash
sudo ufw allow 80
```

## Service Management

The M3TAL API daemon runs as a systemd service. You can manage and monitor it using the following commands:

*   **Check status**:
    ```bash
    systemctl status m3tal-api
    ```

*   **View logs (follow in real-time)**:
    ```bash
    journalctl -u m3tal-api -f
    ```