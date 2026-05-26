# M3TAL CLI Command Reference

Welcome, M3TAL Operative! This document serves as your complete cheat-sheet for interacting with the M3TAL Ecosystem via its command-line interface (CLI). As the M3TAL Ecosystem Documentation Architect, I've compiled this reference to ensure you can efficiently manage your self-hosted infrastructure.

The M3TAL CLI (`/usr/bin/m3tal`) is your single point of entry for managing configurations, deploying services, monitoring health, and interacting with the M3TAL Dashboard.

## 1. Introduction

M3TAL is a unified platform for orchestrating your self-hosted services using Docker and Docker Compose. It streamlines installation, configuration, and management, providing a robust foundation for your digital ecosystem. At its core, M3TAL relies on a powerful Go API daemon, a user-friendly web dashboard, and Traefik for intelligent routing.

## 2. Installation

To get started with M3TAL, follow these steps to install the CLI binary via APT:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## 3. Core Concepts & Architecture

Understanding M3TAL's architecture is key to effective management.

### 3.1. M3TAL Components

*   **CLI binary** (`/usr/bin/m3tal`): The unified Go binary you interact with for all M3TAL operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service, listening on port 8080. It manages Docker, maintains the state DB, and handles API routes for the dashboard and CLI.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container running internally on port 8082, communicating with the API daemon at `http://host.docker.internal:8080`. This provides the web-based user interface.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container exposing services by domain name on host port 80 (HTTP) and 443 (HTTPS). It uses a file provider for dynamic routing and Docker labels for automatic service discovery.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container for secure, zero-config internet access to your services without opening firewall ports.

### 3.2. Filesystem Contract

M3TAL establishes a strict filesystem contract for its operations:

| Path                        | Purpose                                                            |
| :-------------------------- | :----------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary configuration file**. All environment variables are stored here. Managed by `m3tal config wizard`. |
| `/var/lib/m3tal/state.db`   | **SQLite state database**. Auto-created and managed by the M3TAL API daemon. Stores service states, logs, etc. |
| `/opt/m3tal/stack/`         | The canonical directory containing all `*-compose.yml` files and Traefik dynamic configuration. |
| `/docker`                   | **User-facing symlink** to `/opt/m3tal/stack/`. This is where you place your custom `docker compose` files for M3TAL to manage. |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.           |

### 3.3. Dashboard Access Modes

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### 3.3.1. Local Mode (`DASHBOARD_EXPOSE_MODE=local`)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
*   **Mechanism**: Uses the `m3tal-compose.local.yml` override file, which adds a direct port binding to the dashboard container (`${DASHBOARD_PORT:-8082}:8082`).
*   **Access via**: `http://HOST_IP:8082` or `http://localhost:8082`.
*   **Requirements**: No Traefik required. Works out of the box.
*   **Best for**: LAN-only setups, initial installation, local development and testing, or when you don't need domain-based access.

#### 3.3.2. Traefik Mode (`DASHBOARD_EXPOSE_MODE=traefik`)

*   **Configuration**: `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
*   **Mechanism**: Uses the `m3tal-compose.traefik.yml` override, which adds Traefik labels to the dashboard container. Traefik (running via `m3tal up`) then routes requests for `dash.${DOMAIN}` to the dashboard container on its internal port 8082.
*   **Access via**: `http://dash.YOUR_DOMAIN` (e.g., `http://dash.example.com`).
*   **Requirements**: Traefik must be running (started with `m3tal up`), and `DOMAIN` must be correctly configured in `/etc/m3tal/.env`.
*   **Best for**: Domain-based setups, integrating with other services behind Traefik, or when you need internet access (especially with Cloudflared).

### 3.4. Docker & Compose Runtime

M3TAL leverages **Docker Engine** and **Docker Compose V2** for container orchestration. These are hard dependencies for the M3TAL ecosystem.

*   The `m3tal up` command orchestrates all `*-compose.yml` files located within the `/docker/` directory. It effectively runs `docker compose up -d` across all these stacks, applying relevant environment variables from `/etc/m3tal/.env`.
*   The `m3tal dash up` command specifically manages the M3TAL Dashboard container:
    1.  It pulls the latest `m3tal-compose.yml`, `m3tal-compose.local.yml`, and `m3tal-compose.traefik.yml` from GitHub.
    2.  It reads the `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`.
    3.  It then starts the dashboard container using the base compose file and the appropriate override file (local or traefik mode).

### 3.5. Deployment Lifecycle: Adding New Stacks

Adding new services to your M3TAL ecosystem is straightforward:

1.  **Place your compose file**: Create or copy your `docker compose` file into the `/docker/` directory (e.g., `/docker/my-stack-compose.yml`).
2.  **Configure environment variables**: Ensure any environment variables required by your new stack (e.g., `PUID`, `PGID`, custom paths) are set in `/etc/m3tal/.env`. Use `m3tal config wizard` for an interactive setup or `m3tal config set KEY value` for specific variables.
3.  **Start all stacks**: Run `m3tal up` to deploy your new service along with all other managed stacks.

### 3.6. Traefik Routing Architecture

Traefik acts as the central reverse proxy for M3TAL when in `traefik` mode. It is deployed as part of `routing-compose.yml` and configured to:

*   Bind to host ports `80` (HTTP) and `443` (HTTPS - if configured) as its primary entry points.
*   Automatically discover services by inspecting Docker container labels (e.g., for the dashboard in `traefik` mode).
*   Load dynamic routing configurations from `/docker/dynamic/` (e.g., `api.yml` for the M3TAL API), which support hot-reloading.
*   Route `api.${DOMAIN}` to `http://host.docker.internal:8080` (the M3TAL Go API daemon).
*   Route `dash.${DOMAIN}` to the `m3tal-dashboard` container (when `DASHBOARD_EXPOSE_MODE=traefik`).

### 3.7. Port Map

| Port | Service                     | Access                                     |
| :--- | :-------------------------- | :----------------------------------------- |
| 80   | Traefik HTTP entry point    | Public (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 443  | Traefik HTTPS entry point   | Public (if configured)                     |
| 8080 | M3TAL API daemon (Go)       | Host-local only                            |
| 8081 | Traefik dashboard           | Host-local only (`127.0.0.1:8081`)         |
| 8082 | M3TAL Dashboard container | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

## 4. Command Reference

### 4.1. `sudo m3tal` - TUI Control Center

Opens an interactive Text User Interface (TUI) control center, providing a numbered menu for common M3TAL operations. This is often the easiest way to interact with M3TAL.

```bash
sudo m3tal
```

### 4.2. `m3tal init` - Initialize Configuration

Generates the primary configuration file (`/etc/m3tal/.env`) from default values. This command should be used on the first installation or if the `.env` file is missing.

```bash
m3tal init
```

### 4.3. `m3tal doctor` - System Health Check

Performs a pre-flight health check of your M3TAL ecosystem. It verifies Docker connectivity, validates the `/etc/m3tal/.env` file, and checks for port availability, providing actionable insights into potential issues.

```bash
m3tal doctor
```

### 4.4. `m3tal config` - Configuration Management

Manages the `/etc/m3tal/.env` configuration file.

#### 4.4.1. `m3tal config wizard` - Interactive Configuration

Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for initial setup and significant changes.

```bash
m3tal config wizard
```

#### 4.4.2. `m3tal config set KEY VALUE` - Set Environment Variable

Sets a single environment variable in `/etc/m3tal/.env` to the specified value.

```bash
# Set the dashboard expose mode to Traefik
m3tal config set DASHBOARD_EXPOSE_MODE traefik

# Set the primary domain for Traefik
m3tal config set DOMAIN mydomain.com

# Update the PUID (Process User ID) for containers
m3tal config set PUID 1000
```

#### 4.4.3. `m3tal config get KEY` - Get Environment Variable

Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.

```bash
# Get the currently configured dashboard expose mode
m3tal config get DASHBOARD_EXPOSE_MODE

# Retrieve the configured API token
m3tal config get API_TOKEN
```

#### 4.4.4. `m3tal config scan` - List All Stack Variables

Scans all `docker compose` files in `/docker/` to identify and list all environment variables used by your stacks, along with their default values. This helps ensure all necessary variables are defined.

```bash
m3tal config scan
```

#### 4.4.5. `m3tal config list` - List Current .env File

Displays the entire contents of the current `/etc/m3tal/.env` configuration file.

```bash
m3tal config list
```

### 4.5. `m3tal dash` - Dashboard Management

Commands specifically for managing the M3TAL Dashboard container.

#### 4.5.1. `m3tal dashpass [username] [password]` - Update Dashboard Password

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt for them interactively. Credentials are stored in `/docker/users.json`.

```bash
# Interactively update the password for the 'admin' user
m3tal dashpass

# Directly set the password for 'admin' to 'SecurePassword123!'
m3tal dashpass admin SecurePassword123!
```

#### 4.5.2. `m3tal dash up` - Start Dashboard Container

Pulls the latest dashboard compose configuration from GitHub (including `m3tal-compose.yml` and its local/traefik overrides) and then starts the `m3tal-dashboard` container. This command ensures the dashboard is running with the correct exposure mode as defined by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.

```bash
m3tal dash up
```

#### 4.5.3. `m3tal dash down` - Stop Dashboard Container

Stops and removes the `m3tal-dashboard` container.

```bash
m3tal dash down
```

#### 4.5.4. `m3tal dash restart` - Restart Dashboard Container

Restarts the `m3tal-dashboard` container. This is useful after changing dashboard-related environment variables or if the dashboard becomes unresponsive.

```bash
m3tal dash restart
```

#### 4.5.5. `m3tal dash logs` - Stream Dashboard Logs

Streams the logs from the `m3tal-dashboard` container, showing real-time output. Useful for debugging dashboard-related issues.

```bash
m3tal dash logs
```

#### 4.5.6. `m3tal dash status` - Show Dashboard Status

Displays the current status of the `m3tal-dashboard` container (e.g., `running`, `exited`).

```bash
m3tal dash status
```

### 4.6. `m3tal up` - Start All Stacks

Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This command starts or recreates all services defined in your M3TAL ecosystem, including Traefik, Cloudflared, and any custom user stacks.

```bash
m3tal up
```

### 4.7. `m3tal down` - Stop All Stacks

Runs `docker compose down` across all `*-compose.yml` files in the `/docker/` directory. This command stops and removes all containers, networks, and volumes associated with your M3TAL-managed stacks.

```bash
m3tal down
```

### 4.8. `m3tal logs` - Stream All Stack Logs

Streams aggregated logs from all currently running M3TAL-managed containers. This provides a consolidated view of your entire ecosystem's activity.

```bash
m3tal logs
```

## 5. M3TAL API Daemon (systemd Service Management)

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service, essential for core M3TAL functionality. You can manage this service directly using `systemctl` and `journalctl`.

*   **Check status**: View the current status of the M3TAL API daemon.
    ```bash
    systemctl status m3tal-api
    ```
*   **Restart the API**: Restart the M3TAL API daemon. This is often necessary after manual changes to `/etc/m3tal/.env` that affect the API itself.
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **Stream API logs**: Follow the logs from the M3TAL API daemon in real-time.
    ```bash
    journalctl -u m3tal-api -f
    ```

## 6. Docker Compose Fallback Commands

While M3TAL provides convenient wrappers, it relies entirely on Docker Engine and Docker Compose. In advanced troubleshooting scenarios or for direct control, you can use raw `docker compose` commands. Remember that these commands won't automatically apply M3TAL-specific logic (like pulling dashboard configs or handling specific `.env` variables in the same way `m3tal` commands do).

**Key Directories**:
*   All M3TAL-managed `docker compose` files are located in `/docker/` (which symlinks to `/opt/m3tal/stack/`).

**Examples:**

*   **Start a specific stack (e.g., `ollama-compose.yml`):**
    ```bash
    cd /docker
    docker compose -f ollama-compose.yml up -d
    ```

*   **Stop a specific stack:**
    ```bash
    cd /docker
    docker compose -f ollama-compose.yml down
    ```

*   **View logs for a specific service within a stack:**
    ```bash
    cd /docker
    docker compose -f routing-compose.yml logs -f traefik
    ```

*   **Inspect a running service:**
    ```bash
    docker inspect m3tal-dashboard
    ```

*   **View the status of all Docker containers:**
    ```bash
    docker ps -a
    ```

*   **Rebuild and restart a specific dashboard container (if `DASHBOARD_EXPOSE_MODE=local` is active):**
    ```bash
    # Note: For dashboard, m3tal dash up is strongly preferred as it handles overrides and pulling latest configs.
    # This is a raw example for general understanding.
    cd /docker
    docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml up -d --build
    ```