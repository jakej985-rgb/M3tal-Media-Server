```markdown
# docs/COMMAND_REFERENCE.md - M3TAL CLI Cheat-Sheet

Welcome, M3TAL Engineer! This document serves as your comprehensive guide to the M3TAL Command-Line Interface (CLI). Designed by DocSmith himself, it covers every command, its purpose, and provides real-world usage examples to get you up and running faster than a photon torpedo.

The M3TAL CLI (`/usr/bin/m3tal`) is your single entry point for managing the entire M3TAL ecosystem, from initial setup and configuration to dashboard interaction and container orchestration.

---

## M3TAL System Installation

First, ensure M3TAL is installed on your system. Follow these steps for an APT-based installation:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## Core M3TAL Commands

These commands provide essential system-level interactions and initial setup.

### `sudo m3tal`

Opens the interactive M3TAL Control Center, a Terminal User Interface (TUI) with a numbered menu for common operations. This is your go-to for a guided experience.

**Usage Example:**
```bash
sudo m3tal
```
*Navigates to the interactive TUI, presenting options like "1. Start All Stacks", "2. Stop All Stacks", "3. Configure M3TAL", etc.*

### `m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from M3TAL's default settings. This command is crucial for the very first installation or to reset your configuration to defaults.

**Usage Example:**
```bash
m3tal init
```
*Creates `/etc/m3tal/.env` with default values for variables like `DASHBOARD_PORT`, `DOMAIN`, `PUID`, `PGID`, etc.*

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL environment. It verifies Docker connectivity, checks the validity of `/etc/m3tal/.env`, and ensures that required ports are available. Run this to diagnose common issues.

**Usage Example:**
```bash
m3tal doctor
```
*Outputs a report like:*
```
[INFO] M3TAL Doctor: Pre-flight check initiated...
[SUCCESS] Docker daemon is running and accessible.
[SUCCESS] /etc/m3tal/.env file is valid.
[SUCCESS] Port 80 (Traefik HTTP) is available.
[SUCCESS] Port 8080 (M3TAL API) is available.
[SUCCESS] Port 8082 (M3TAL Dashboard) is available.
[INFO] M3TAL environment is healthy.
```

---

## M3TAL Configuration (`m3tal config`)

M3TAL's primary configuration resides in `/etc/m3tal/.env`. These commands allow you to manage this file effectively.

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating the variables in `/etc/m3tal/.env`. This is the recommended method for making changes, as it provides explanations and validates inputs.

**Usage Example:**
```bash
m3tal config wizard
```
*Prompts the user for values for `DOMAIN`, `DASHBOARD_EXPOSE_MODE`, `PUID`, `PGID`, etc., and saves them to `/etc/m3tal/.env`.*

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env` to a specified value. Use this for quick, direct modifications.

**Usage Example:**
```bash
m3tal config set DOMAIN myhome.tech
```
*Updates the `DOMAIN` variable in `/etc/m3tal/.env` to `myhome.tech`.*

### `m3tal config get KEY`

Retrieves and displays the value of a specific environment variable from `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal config get DASHBOARD_PORT
```
*Outputs:*
```
8082
```

### `m3tal config scan`

Lists all detected environment variables across all active M3TAL stacks. This provides a comprehensive overview of variables being utilized by your services.

**Usage Example:**
```bash
m3tal config scan
```
*Outputs a list of keys and their values, potentially including those inherited from the system or defaults.*

### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file. This is useful for reviewing your explicit M3TAL configuration.

**Usage Example:**
```bash
m3tal config list
```
*Outputs the raw contents of `/etc/m3tal/.env`:*
```
# M3TAL Environment Configuration
DASHBOARD_PORT=8082
DASHBOARD_EXPOSE_MODE=local
DOMAIN=m3tal.local
PUID=1000
PGID=1000
...
```

---

## M3TAL Dashboard Management (`m3tal dash`)

The M3TAL Dashboard is your web-based control panel. These commands manage its lifecycle and access. The dashboard stores user credentials in `/docker/users.json`.

### Understanding Dashboard Access Modes (Critical)

The M3TAL Dashboard can be accessed in two primary modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`. This variable determines which Docker Compose override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) is used when `m3tal dash up` is executed.

1.  **`local` mode (Default)**
    *   **Configuration:** `DASHBOARD_EXPOSE_MODE=local` in `/etc/m3tal/.env`.
    *   **Mechanism:** Uses `m3tal-compose.local.yml`, which adds a direct port binding (`${DASHBOARD_PORT:-8082}:8082`) to the dashboard container.
    *   **Access:** Via `http://HOST_IP:8082` or `http://localhost:8082`.
    *   **Requirements:** No Traefik required. Works out-of-the-box.
    *   **Best for:** LAN-only setups, first-time users, local testing.

2.  **`traefik` mode**
    *   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik` in `/etc/m3tal/.env`.
    *   **Mechanism:** Uses `m3tal-compose.traefik.yml`, which adds Traefik labels to the dashboard container. Traefik (running via `m3tal up`) then routes `dash.${DOMAIN}` to the dashboard on port 8082.
    *   **Access:** Via `http://dash.YOUR_DOMAIN` (e.g., `http://dash.m3tal.local`).
    *   **Requirements:** Traefik must be running (`m3tal up` must have been executed for the `routing` stack).
    *   **Best for:** Domain-based setups, integrating with other services behind a reverse proxy.

### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, it will launch an interactive prompt. User credentials are stored in `/docker/users.json`.

**Usage Example (Interactive):**
```bash
m3tal dashpass
```
*Prompts for username and new password, then saves the hashed password to `/docker/users.json`.*

**Usage Example (Direct):**
```bash
m3tal dashpass admin newSecurePass123
```
*Sets the password for the 'admin' user to 'newSecurePass123' in `/docker/users.json`.*

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration files from GitHub, then starts the M3TAL dashboard container using the appropriate expose mode (`local` or `traefik`) defined in `/etc/m3tal/.env`.

**Usage Example:**
```bash
m3tal dash up
```
*Downloads `m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml` to `/docker/`, then starts the `m3tal-dashboard` container based on `DASHBOARD_EXPOSE_MODE`.*

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage Example:**
```bash
m3tal dash down
```
*Stops and removes the `m3tal-dashboard` container.*

### `m3tal dash restart`

Restarts the M3TAL dashboard container. Useful after configuration changes or for troubleshooting.

**Usage Example:**
```bash
m3tal dash restart
```
*Stops, then starts the `m3tal-dashboard` container.*

### `m3tal dash logs`

Streams real-time logs from the M3TAL dashboard container. Essential for debugging.

**Usage Example:**
```bash
m3tal dash logs
```
*Displays continuous log output from the `m3tal-dashboard` container.*

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container (e.g., "running", "exited").

**Usage Example:**
```bash
m3tal dash status
```
*Outputs a line indicating the dashboard container's state, like `m3tal-dashboard Up 5 minutes (healthy)`.*

---

## M3TAL Stack Management (`m3tal`)

These commands manage all Docker Compose-based service stacks within the M3TAL ecosystem. User-defined stacks and M3TAL's core services are managed via Docker Compose files (`*-compose.yml`) located in `/docker/` (which is a symlink to `/opt/m3tal/stack/`).

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files found in `/docker/`. This command brings up all your configured services (including routing, and if configured, the dashboard).

**Usage Example:**
```bash
m3tal up
```
*Starts all Docker containers defined in `/docker/*.yml` files in detached mode, bringing up services like Traefik, Cloudflared, and any user-added stacks (e.g., `ollama`, `jellyfin`).*

### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This command stops and removes all containers, networks, and volumes defined by your M3TAL stacks.

**Usage Example:**
```bash
m3tal down
```
*Stops and removes all Docker containers associated with M3TAL's stacks.*

### `m3tal logs`

Streams aggregated logs from all running M3TAL containers across all stacks. This provides a centralized view for monitoring and debugging.

**Usage Example:**
```bash
m3tal logs
```
*Displays combined, real-time log output from all containers managed by M3TAL (e.g., `traefik`, `m3tal-dashboard`, `ollama`).*

---

## Systemd Service Management

The M3TAL API daemon (`m3tal-api.service`) runs as a systemd service, providing the backend for the CLI and Dashboard. These commands interact directly with systemd.

### `systemctl status m3tal-api`

Checks the current status of the M3TAL API daemon.

**Usage Example:**
```bash
systemctl status m3tal-api
```
*Outputs whether the service is active/inactive, running or failed, and recent log entries.*

### `journalctl -u m3tal-api -f`

Streams real-time logs from the M3TAL API daemon. Essential for debugging issues related to the core API or its interactions with Docker.

**Usage Example:**
```bash
journalctl -u m3tal-api -f
```
*Displays continuous log output from the `m3tal-api` service.*

---

## Direct Docker / Compose Fallback

M3TAL leverages Docker Engine and Docker Compose V2. While the `m3tal` CLI abstracts many Docker commands, you can always use direct `docker compose` commands for advanced debugging or specific operations.

The canonical directory for M3TAL's Docker Compose files is `/opt/m3tal/stack/`, which is symlinked to `/docker/` for user convenience.

**To interact with a specific stack directly:**

1.  **Navigate to the stack directory:**
    ```bash
    cd /docker/
    ```

2.  **Use `docker compose` commands (e.g., for the `routing` stack):**

    *   **Start the routing stack in detached mode:**
        ```bash
        docker compose -f routing-compose.yml up -d
        ```

    *   **Stop the routing stack:**
        ```bash
        docker compose -f routing-compose.yml down
        ```

    *   **View logs for a specific service (e.g., `traefik` within `routing-compose.yml`):**
        ```bash
        docker compose -f routing-compose.yml logs -f traefik
        ```

    *   **Inspect all running services:**
        ```bash
        docker ps
        ```

    *   **Inspect all Docker volumes:**
        ```bash
        docker volume ls
        ```

---

## M3TAL Filesystem Contract

Understanding the key file locations is vital for advanced management and troubleshooting:

| Path                        | Purpose                                                                                                                                                                                                                                 |
| :-------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | **Primary Configuration File.** Contains all environment variables for M3TAL's operation. Managed by `m3tal config wizard` and `m3tal config set`.                                                                                       |
| `/var/lib/m3tal/state.db`   | **SQLite State Database.** Auto-created and managed by the M3TAL API daemon. Stores internal state, service information, and other operational data.                                                                                     |
| `/opt/m3tal/stack/`         | **Canonical Stack Directory.** This is where all M3TAL's core Docker Compose files (`m3tal-compose.yml`, `routing-compose.yml`, etc.) and their overrides reside. User-added stack files also belong here.                               |
| `/docker`                   | **User-Facing Stack Symlink.** This is a symbolic link to `/opt/m3tal/stack/`. Use `/docker/` when referring to stack files in commands or documentation. Placing a new `*-compose.yml` file here makes it discoverable by `m3tal up`. |
| `/docker/users.json`        | **Dashboard Credential Store.** Stores hashed usernames and passwords for M3TAL Dashboard access. Managed by `m3tal dashpass`.                                                                                                          |
| `/docker/dynamic/`          | **Traefik Dynamic Configuration.** Traefik loads additional routing rules (e.g., for the API daemon) from `.yml` files in this directory. Supports hot-reloading.                                                                       |

---

## M3TAL Port Map

| Port | Service               | Access                                                 |
| :--- | :-------------------- | :----------------------------------------------------- |
| 80   | Traefik HTTP Entry    | Public (if `routing` stack is running and configured)  |
| 8080 | M3TAL API daemon (Go) | Host-local access only                                 |
| 8081 | Traefik Dashboard     | Host-local only (for Traefik's own UI/metrics)         |
| 8082 | M3TAL Dashboard       | Direct port (if `DASHBOARD_EXPOSE_MODE=local`) or via Traefik (if `DASHBOARD_EXPOSE_MODE=traefik`) |

---

This concludes your M3TAL CLI cheat-sheet. May your stacks be evergreen and your logs be clear!
```