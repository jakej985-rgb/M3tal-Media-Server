# M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, my mission is to equip you with the knowledge to wield the M3TAL CLI with absolute precision. This document serves as your definitive cheat sheet, detailing every core command, its real-world usage, and the underlying system mechanics that make the M3TAL ecosystem tick.

The `m3tal` CLI binary (`/usr/bin/m3tal`) is your single entry point for all operations, orchestrating the M3TAL API daemon and your Dockerized applications. It acts as the bridge between your commands and the robust M3TAL architecture.

---

## I. M3TAL Core CLI Commands

### `sudo m3tal`
Opens the interactive TUI (Terminal User Interface) Control Center. This interface provides a numbered menu for common operations, configuration, and status checks, offering a user-friendly alternative to direct CLI commands. Requires `sudo` as it may perform privileged operations like Docker management.

**Usage Example:**
```bash
sudo m3tal
```

### `m3tal init`
Generates the primary configuration file, `/etc/m3tal/.env`, from system defaults. This command should be used on first installation or if your `.env` file is missing. It ensures all critical environment variables are present before M3TAL services are started.

**Usage Example:**
```bash
m3tal init
```

### `m3tal doctor`
Performs a crucial pre-flight health check of the M3TAL ecosystem. This includes verifying Docker connectivity, checking the validity and integrity of `/etc/m3tal/.env`, and ensuring that required ports (e.g., 80, 8080, 8082) are not already in use by other processes.

**Usage Example:**
```bash
m3tal doctor
```

### `m3tal config wizard`
Launches an interactive, guided wizard to configure the `/etc/m3tal/.env` file. This is the recommended method for making changes to your M3TAL environment variables, as it provides explanations and validates inputs.

**Usage Example:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`
Sets a single environment variable in `/etc/m3tal/.env` to a specified value. This command allows for quick, non-interactive updates to your configuration. Note that changes may require restarting the `m3tal-api.service` or M3TAL Docker containers to take effect.

**Usage Example:**
```bash
m3tal config set DASHBOARD_EXPOSE_MODE traefik
```

### `m3tal config get KEY`
Retrieves and displays the current value of a single environment variable from `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal config get DOMAIN
```
**Output Example:**
```
example.com
```

### `m3tal config scan`
Lists all known environment variables and their values across all M3TAL stacks. This provides a comprehensive overview of your entire M3TAL configuration landscape, including those defined within specific compose files.

**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`
Displays the current contents of the primary M3TAL configuration file, `/etc/m3tal/.env`. This provides a direct view of the variables managed by the `m3tal config` commands.

**Usage Example:**
```bash
m3tal config list
```

### `m3tal dashpass [username] [password]`
Updates the password for a user accessing the M3TAL Dashboard. If `username` and `password` are omitted, the command becomes interactive, prompting you for the necessary details. Credentials are stored in `/docker/users.json`.

**Usage Example (Interactive):**
```bash
m3tal dashpass
# Enter username: admin
# Enter new password: myStrongPassword!
# Confirm new password: myStrongPassword!
```

**Usage Example (Direct):**
```bash
m3tal dashpass admin newSecurePassword
```

### `m3tal dash up`
Manages the M3TAL Dashboard container (`m3tal-dashboard`). This command performs several critical actions:
1.  Pulls the latest dashboard compose configurations (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub.
2.  Reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
3.  Starts the `m3tal-dashboard` container using the appropriate compose override file based on `DASHBOARD_EXPOSE_MODE`.

**Usage Example:**
```bash
m3tal dash up
```

### `m3tal dash down`
Stops and removes the M3TAL Dashboard container.

**Usage Example:**
```bash
m3tal dash down
```

### `m3tal dash restart`
Restarts the M3TAL Dashboard container. This is useful after making configuration changes that affect the dashboard or if it becomes unresponsive.

**Usage Example:**
```bash
m3tal dash restart
```

### `m3tal dash logs`
Streams the real-time logs from the M3TAL Dashboard container. Essential for troubleshooting and monitoring.

**Usage Example:**
```bash
m3tal dash logs
```

### `m3tal dash status`
Shows the current status of the M3TAL Dashboard container (e.g., running, stopped, exited).

**Usage Example:**
```bash
m3tal dash status
```

### `m3tal up`
Initiates `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory (which symlinks to `/opt/m3tal/stack/`). This command starts or recreates all defined M3TAL services and user-defined stacks. This is your primary command for deploying or updating your entire M3TAL ecosystem.

**Usage Example:**
```bash
m3tal up
```

### `m3tal down`
Executes `docker compose down` across all `*-compose.yml` files in `/docker/`. This stops and removes all M3TAL-managed containers, networks, and volumes for all stacks.

**Usage Example:**
```bash
m3tal down
```

### `m3tal logs`
Streams aggregated logs from all currently running M3TAL-managed Docker containers. This provides a unified view of your entire system's operational output.

**Usage Example:**
```bash
m3tal logs
```

---

## II. M3TAL System Architecture & Filesystem Contract

Understanding the M3TAL filesystem contract is paramount for effective management and troubleshooting.

| Path                     | Purpose                                                                |
| :----------------------- | :--------------------------------------------------------------------- |
| `/etc/m3tal/.env`        | Primary configuration file. Managed by `m3tal config wizard`.          |
| `/var/lib/m3tal/state.db`| SQLite state database. Auto-created and managed by the API daemon.     |
| `/opt/m3tal/stack/`      | Canonical directory for M3TAL's internal stack files (compose, Traefik config). |
| `/docker`                | Symlink to `/opt/m3tal/stack/`. This is the user-facing path for all stack operations. Place your `*-compose.yml` files here. |
| `/docker/users.json`     | Dashboard credential store. Managed by `m3tal dashpass`.               |
| `/docker/dynamic/`       | Directory for Traefik dynamic configuration files (e.g., `api.yml`). These hot-reload. |

---

## III. M3TAL Dashboard Access Modes

The M3TAL Dashboard offers two distinct access modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. The `m3tal dash up` command dynamically selects the appropriate Docker Compose override file based on this setting.

### Mode 1: `local` (Default)
-   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
-   **Mechanism:** Uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's internal port (`8082`) to the host's `DASHBOARD_PORT` (defaulting to `8082`).
-   **Access Via:** `http://HOST_IP:8082` or `http://localhost:8082`.
-   **Requirements:** No Traefik required.
-   **Best For:** LAN-only setups, first-time users, local development, or environments where a reverse proxy is not desired.

**Usage Example:**
After setting `DASHBOARD_EXPOSE_MODE=local` with `m3tal config set DASHBOARD_EXPOSE_MODE local`, then run:
```bash
m3tal dash up
# Access the dashboard at http://localhost:8082 or http://YOUR_SERVER_IP:8082
```

### Mode 2: `traefik`
-   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
-   **Mechanism:** Uses the `m3tal-compose.traefik.yml` override, which adds specific Traefik labels to the dashboard container. Traefik (running via `routing-compose.yml`) then detects these labels and routes traffic for `dash.${DOMAIN}` to the dashboard.
-   **Access Via:** `http://dash.YOUR_DOMAIN` (e.g., `http://dash.example.com`).
-   **Requirements:** Traefik must be running as part of your `m3tal up` services.
-   **Best For:** Domain-based setups, exposing services behind a central reverse proxy, and integrating with other services managed by Traefik.

**Usage Example:**
After setting `DASHBOARD_EXPOSE_MODE=traefik` with `m3tal config set DASHBOARD_EXPOSE_MODE traefik`, and ensuring `DOMAIN` is set (e.g., `m3tal config set DOMAIN example.com`), then run:
```bash
m3tal dash up
m3tal up # Ensure Traefik is running
# Access the dashboard at http://dash.example.com
```

---

## IV. Docker / Compose Runtime Explained

M3TAL leverages **Docker Engine** and **Docker Compose V2** as its underlying containerization and orchestration platform. These are hard dependencies for the M3TAL ecosystem.

-   The `m3tal up` command executes `docker compose up -d` (in detached mode) across *all* `*-compose.yml` files it finds within the `/docker/` directory. This unified approach simplifies the deployment of complex, multi-service applications.
-   The `m3tal dash up` command is specifically tailored for the M3TAL dashboard. It not only starts the dashboard container but also ensures you're running the latest configuration by downloading updated compose files and applying the correct `local` or `traefik` override based on your `/etc/m3tal/.env` settings.
-   **User Stacks:** To add your own Docker Compose applications to the M3TAL ecosystem, simply place your `my-application-compose.yml` file into the `/docker/` directory. M3TAL will discover and manage it alongside its core services.

---

## V. Deployment Lifecycle — Day 2 Operations

**Installing a New Stack:**

1.  **Place your Compose file:** Copy your Docker Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
    ```bash
    sudo cp ~/my-stack-compose.yml /docker/
    ```
2.  **Configure Environment Variables:** Ensure all required environment variables for your new stack are defined. Use the `m3tal config wizard` for an interactive experience, or `m3tal config set` for direct updates to `/etc/m3tal/.env`.
    ```bash
    m3tal config set MY_STACK_DATA_PATH /mnt/my-stack-data
    ```
3.  **Deploy All Stacks:** Run `m3tal up` to start your new stack along with all other M3TAL-managed services.
    ```bash
    m3tal up
    ```

---

## VI. Traefik Routing Architecture

Traefik, deployed as a container via `routing-compose.yml`, is the central reverse proxy for the M3TAL ecosystem.

-   It binds to port `80` on the host, serving as the HTTP entry point for all incoming web traffic.
-   **Automatic Service Discovery:** Traefik automatically discovers and configures routes for services that carry specific Docker labels (e.g., the M3TAL Dashboard when in `traefik` mode).
-   **Dynamic Configuration:** Traefik loads additional dynamic configurations from the `/docker/dynamic/` directory (e.g., `api.yml`). These files are hot-reloaded, meaning changes take effect instantly without restarting Traefik.
-   **M3TAL API Routing:** A dedicated dynamic configuration (`/docker/dynamic/api.yml`) routes requests for `api.${DOMAIN}` directly to the M3TAL API daemon running on the host at `http://host.docker.internal:8080`.
-   **Dashboard Routing:** When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik labels on the dashboard container route `dash.${DOMAIN}` to the M3TAL Dashboard on its internal port `8082`.

---

## VII. Systemd Service Management

The M3TAL API daemon is a critical Go binary that runs as a systemd service named `m3tal-api.service`. This daemon is responsible for managing Docker interactions, the state database (`/var/lib/m3tal/state.db`), and exposing the API routes that the `m3tal` CLI and Dashboard interact with.

**Check Service Status:**
```bash
systemctl status m3tal-api
```

**Restart the API Daemon:**
```bash
sudo systemctl restart m3tal-api
```

**View Real-time API Daemon Logs:**
```bash
journalctl -u m3tal-api -f
```

---

## VIII. Port Map

| Port | Service               | Access                                  |
| :--- | :-------------------- | :-------------------------------------- |
| 80   | Traefik HTTP entry point | Public (when Traefik is active)         |
| 8080 | M3TAL API daemon (Go)   | Host-local only                         |
| 8081 | Traefik dashboard     | Host-local only (e.g., `http://localhost:8081`) |
| 8082 | M3TAL Dashboard       | Direct port (local mode) or via Traefik (traefik mode) |

---

## IX. Direct Docker Compose Fallback

While the `m3tal` CLI provides a convenient, aggregated interface for managing your Docker stacks, you can always resort to direct Docker Compose V2 commands as a fallback or for specific, granular operations.

Remember that `m3tal` aggregates all `*-compose.yml` files in `/docker/`. When using direct `docker compose` commands, you'll need to explicitly specify all relevant compose files using multiple `-f` flags.

**Example: Starting M3TAL core services (dashboard & routing) directly:**
```bash
sudo docker compose \
  -f /docker/m3tal-compose.yml \
  -f /docker/m3tal-compose.local.yml \
  -f /docker/routing-compose.yml \
  up -d
```
*(Note: Replace `m3tal-compose.local.yml` with `m3tal-compose.traefik.yml` if `DASHBOARD_EXPOSE_MODE=traefik`.)*

**Example: Stopping M3TAL core services directly:**
```bash
sudo docker compose \
  -f /docker/m3tal-compose.yml \
  -f /docker/m3tal-compose.local.yml \
  -f /docker/routing-compose.yml \
  down
```

**Example: Viewing logs for a specific service (e.g., dashboard):**
```bash
sudo docker compose -f /docker/m3tal-compose.yml logs -f m3tal-dashboard
```

---

## X. Installation (APT)

If you haven't installed M3TAL yet, follow these steps to add the APT repository and install the `m3tal` package:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```