Greetings, Fellow Automators and System Architects! I am DocSmith, your M3TAL Ecosystem Documentation Architect. This document serves as your complete CLI cheat-sheet, meticulously detailing every command, its purpose, and real-world usage examples within the M3TAL ecosystem.

---

# M3TAL CLI Command Reference

The M3TAL CLI (`/usr/bin/m3tal`) is the unified entry point for managing your self-hosted stack. It orchestrates Docker Compose services, manages configuration, and interacts with the M3TAL API daemon to provide a robust and streamlined experience.

## I. M3TAL Filesystem Contract

Understanding the key locations M3TAL uses is crucial for effective management.

| Path                         | Purpose                                                                                                                                              |
| :--------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`            | **Primary configuration file**. Contains all environment variables for M3TAL and its managed stacks. Managed by `m3tal config wizard` and `m3tal config set`. |
| `/var/lib/m3tal/state.db`    | SQLite state database. Stores internal M3TAL API daemon state, user data, and service information. Auto-created by the API daemon.                         |
| `/opt/m3tal/stack/`          | Canonical stack directory. Contains M3TAL's core Docker Compose files (`routing-compose.yml`, `m3tal-compose.yml`, etc.) and Traefik dynamic configuration. |
| `/docker`                    | **User-facing symlink** → `/opt/m3tal/stack/`. This is the primary directory where users place their custom `*-compose.yml` files for M3TAL to manage.      |
| `/docker/users.json`         | Dashboard credential store. Contains usernames and hashed passwords for the M3TAL Dashboard. Managed by `m3tal dashpass`.                                 |

## II. Core M3TAL Commands

### `sudo m3tal`

Opens the interactive TUI (Text User Interface) Control Center. This provides a user-friendly, menu-driven interface for common operations. Requires `sudo` as it interacts with Docker and system services.

**Usage Example:**

```bash
sudo m3tal
```

### `m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from M3TAL's default values. This command should be run on first installation to set up the basic environment. It will warn if the file already exists.

**Usage Example:**

```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of the M3TAL system. It verifies Docker connectivity, checks the validity and required variables within `/etc/m3tal/.env`, and ensures essential ports (e.g., 80, 8080, 8082) are available.

**Usage Example:**

```bash
m3tal doctor
```

## III. Configuration Management (`m3tal config`)

These commands manage the `/etc/m3tal/.env` file, which is central to your M3TAL deployment.

### `m3tal config wizard`

Launches an interactive, step-by-step wizard to guide you through configuring or updating key environment variables in `/etc/m3tal/.env`. This is the recommended way to manage your M3TAL configuration.

**Usage Example:**

```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a specific environment variable (`KEY`) to a new `VALUE` in `/etc/m3tal/.env`. This allows for direct modification outside the wizard.

**Usage Example:**

```bash
m3tal config set DOMAIN "myawesomehomelab.com"
m3tal config set DASHBOARD_EXPOSE_MODE "traefik"
```

### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable (`KEY`) from `/etc/m3tal/.env`.

**Usage Example:**

```bash
m3tal config get PUID
m3tal config get DASHBOARD_EXPOSE_MODE
```

### `m3tal config scan`

Lists all known environment variables across all managed Docker Compose stacks, showing their current values and whether they are set explicitly or using a default. This is useful for auditing your entire configuration.

**Usage Example:**

```bash
m3tal config scan
```

### `m3tal config list`

Displays the entire content of the current `/etc/m3tal/.env` file, line by line.

**Usage Example:**

```bash
m3tal config list
```

## IV. Dashboard Management (`m3tal dash`)

These commands specifically manage the `m3tal-dashboard` container, M3TAL's web-based control panel.

### Dashboard Access Modes (Critical)

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`:

1.  **`DASHBOARD_EXPOSE_MODE=local` (Default)**
    *   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082`
    *   **Mechanism:** Adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the `m3tal-dashboard` container via the `m3tal-compose.local.yml` override.
    *   **Requires:** No Traefik configuration. Works out-of-the-box on a fresh install.
    *   **Best for:** LAN-only setups, first-time users, local testing, or when Traefik is not desired for the dashboard.

2.  **`DASHBOARD_EXPOSE_MODE=traefik`**
    *   **Access via:** `http://dash.DOMAIN` (e.g., `http://dash.myawesomehomelab.com`)
    *   **Mechanism:** Applies Traefik labels to the `m3tal-dashboard` container via the `m3tal-compose.traefik.yml` override, routing `dash.DOMAIN` to the dashboard's internal port 8082.
    *   **Requires:** Traefik gateway to be running (`m3tal up` will start it if `routing-compose.yml` is present) and a `DOMAIN` configured in `/etc/m3tal/.env`.
    *   **Best for:** Domain-based access, integrating the dashboard behind a reverse proxy with other services.

### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If no username or password is provided, it will launch an interactive prompt to create/update credentials. The credentials are stored in `/docker/users.json`.

**Usage Examples:**

```bash
m3tal dashpass                                   # Interactive prompt for username and password
m3tal dashpass admin newSecurePassword123!     # Set password for 'admin' user directly
```

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub, then starts the `m3tal-dashboard` container using the appropriate compose override based on your `DASHBOARD_EXPOSE_MODE`.

**Usage Example:**

```bash
m3tal dash up
```

### `m3tal dash down`

Stops and removes the `m3tal-dashboard` container.

**Usage Example:**

```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container.

**Usage Example:**

```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams the aggregated logs from the `m3tal-dashboard` container. Use `Ctrl+C` to exit.

**Usage Example:**

```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current operational status of the `m3tal-dashboard` container (e.g., `running`, `exited`).

**Usage Example:**

```bash
m3tal dash status
```

## V. Global Stack Management

These commands manage all Docker Compose stacks defined by `*-compose.yml` files within the `/docker/` directory.

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This command starts or recreates all services defined in your M3TAL deployment, including Traefik, Cloudflared (if configured), and any custom user stacks.

**Usage Example:**

```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in the `/docker/` directory. This command stops and removes all containers, networks, and volumes associated with your M3TAL stacks.

**Usage Example:**

```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL-managed Docker containers. This provides a unified view of all service activity. Use `Ctrl+C` to exit.

**Usage Example:**

```bash
m3tal logs
```

## VI. M3TAL API Daemon Service Management (systemd)

The M3TAL API daemon is a Go binary running as a `systemd` service (`m3tal-api.service`). It's responsible for managing Docker, interacting with the state database (`/var/lib/m3tal/state.db`), and exposing the API routes on port 8080.

### `systemctl status m3tal-api`

Checks the current status of the M3TAL API daemon service, including whether it's active, running, and recent log entries.

**Usage Example:**

```bash
sudo systemctl status m3tal-api
```

### `systemctl restart m3tal-api`

Restarts the M3TAL API daemon service. This is often necessary after making changes to `/etc/m3tal/.env` that the API needs to recognize.

**Usage Example:**

```bash
sudo systemctl restart m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the real-time logs for the M3TAL API daemon service directly from the systemd journal. This is invaluable for debugging issues with the API. Use `Ctrl+C` to exit.

**Usage Example:**

```bash
sudo journalctl -u m3tal-api -f
```

## VII. Direct Docker Compose Commands (Fallback & Advanced Usage)

M3TAL uses **Docker Engine** and **Docker Compose V2** under the hood. While the `m3tal` CLI abstracts most Docker interactions, you can always use direct `docker compose` commands for advanced debugging, specific container management, or fallback scenarios.

All M3TAL-managed compose files (`*-compose.yml`) reside in `/docker/` (which symlinks to `/opt/m3tal/stack/`).

**To manage specific stacks or services directly:**

1.  **List running containers (all Docker):**
    ```bash
    docker ps
    ```

2.  **List all M3TAL services defined in compose files (dry run):**
    Navigate to the `/docker/` directory and list files:
    ```bash
    ls -l /docker/
    ```
    You'll see files like `m3tal-compose.yml`, `routing-compose.yml`, `my-custom-app-compose.yml`, etc.

3.  **Start a specific stack (e.g., just the routing stack):**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml up -d
    ```

4.  **Stop a specific stack:**
    ```bash
    sudo docker compose -f /docker/m3tal-compose.yml down
    ```

5.  **View logs for a specific service within a stack (e.g., Traefik):**
    ```bash
    sudo docker compose -f /docker/routing-compose.yml logs -f traefik
    ```

6.  **Execute a command inside a running container (e.g., bash into the dashboard):**
    ```bash
    sudo docker exec -it m3tal-dashboard bash
    ```

7.  **Pull latest images for all stacks (same as `m3tal pull`, if it existed):**
    ```bash
    sudo docker compose -f /docker/m3tal-compose.yml -f /docker/routing-compose.yml -f /docker/my-app-compose.yml pull
    ```
    *Note: The `m3tal up` command implicitly pulls new images if they are updated.*

## VIII. Network and Port Map

| Port | Service                               | Access                                                          |
| :--- | :------------------------------------ | :-------------------------------------------------------------- |
| 80   | Traefik HTTP entry point              | Public (when `DASHBOARD_EXPOSE_MODE=traefik` and for other services) |
| 8080 | M3TAL API daemon (Go)                 | Host-local (accessed by containers via `host.docker.internal`)  |
| 8081 | Traefik dashboard (admin UI)          | Host-local only (e.g., `http://localhost:8081`)                 |
| 8082 | M3TAL Dashboard                       | Direct port (when `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (when `DASHBOARD_EXPOSE_MODE=traefik`) |

## IX. M3TAL Installation

For a fresh M3TAL installation, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

This comprehensive reference should equip you with the knowledge to efficiently manage your M3TAL ecosystem. Should you encounter any issues or require further clarification, consult the main M3TAL documentation or community channels. Happy automating!