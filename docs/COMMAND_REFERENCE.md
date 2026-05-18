# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for the M3TAL CLI. It details each command, its purpose, and provides real-world usage examples.

---

## Table of Contents

1.  [M3TAL Core Commands](#m3tal-core-commands)
2.  [Dashboard Management](#dashboard-management)
3.  [Configuration Management](#configuration-management)
4.  [System Service Management](#system-service-management)
5.  [Direct Docker Compose Fallback](#direct-docker-compose-fallback)

---

## 1. M3TAL Core Commands

These commands are the primary interface for managing the M3TAL ecosystem.

### `sudo m3tal`

**Description:**
Launches the interactive M3TAL TUI (Text User Interface) Control Center. This provides a numbered menu-driven interface for common operations.

**Usage Example:**
```bash
sudo m3tal
```
*(This will display the TUI menu, allowing you to select actions like starting stacks, managing the dashboard, and more.)*

### `m3tal init`

**Description:**
Generates the default `/etc/m3tal/.env` configuration file. This command should be used during the initial installation of M3TAL to establish a baseline configuration.

**Usage Example:**
```bash
m3tal init
```

### `m3tal doctor`

**Description:**
Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, checks the validity of your `.env` file, and ensures necessary ports are available.

**Usage Example:**
```bash
m3tal doctor
```

### `m3tal config wizard`

**Description:**
Starts an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file. This is the recommended way to set up your environment variables.

**Usage Example:**
```bash
m3tal config wizard
```
*(Follow the on-screen prompts to configure various M3TAL settings.)*

### `m3tal config set KEY VALUE`

**Description:**
Sets a single environment variable in your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config set DOMAIN my-m3tal-domain.com
```

### `m3tal config get KEY`

**Description:**
Reads and displays the value of a specific environment variable from your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config get DASHBOARD_PORT
```
*(This might output `8082` if that's the configured value.)*

### `m3tal config scan`

**Description:**
Lists all environment variables across all configured M3TAL stacks, including their current values and any default values.

**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`

**Description:**
Displays the entire contents of your current `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config list
```

### `m3tal up`

**Description:**
Initiates the startup of all services defined in `*-compose.yml` files within the `/docker/` directory. This brings up your entire M3TAL environment.

**Usage Example:**
```bash
sudo m3tal up
```
*(Ensure you are in the `/docker` directory or that M3TAL can access it.)*

### `m3tal down`

**Description:**
Shuts down all services managed by M3TAL, stopping the containers defined in `*-compose.yml` files in `/docker/`.

**Usage Example:**
```bash
sudo m3tal down
```

### `m3tal logs`

**Description:**
Streams aggregated logs from all currently running M3TAL stacks. This is invaluable for debugging.

**Usage Example:**
```bash
sudo m3tal logs
```
*(Press `Ctrl+C` to stop streaming.)*

---

## 2. Dashboard Management

Commands specifically for managing the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

**Description:**
Updates the password for a dashboard user. If `username` and `password` are omitted, the command becomes interactive, prompting for the username and new password.

**Usage Example (Interactive):**
```bash
m3tal dashpass
```
*(You will be prompted to enter the username and then the new password twice.)*

**Usage Example (Non-Interactive):**
```bash
m3tal dashpass admin newSecurePassword123
```

### `m3tal dash up`

**Description:**
Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container. It respects the `DASHBOARD_EXPOSE_MODE` setting in your `.env` file.

**Usage Example:**
```bash
sudo m3tal dash up
```

### `m3tal dash down`

**Description:**
Stops the `m3tal-dashboard` container.

**Usage Example:**
```bash
sudo m3tal dash down
```

### `m3tal dash restart`

**Description:**
Restarts the `m3tal-dashboard` container.

**Usage Example:**
```bash
sudo m3tal dash restart
```

### `m3tal dash logs`

**Description:**
Streams the logs from the `m3tal-dashboard` container.

**Usage Example:**
```bash
sudo m3tal dash logs
```
*(Press `Ctrl+C` to stop streaming.)*

### `m3tal dash status`

**Description:**
Displays the current status of the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash status
```

---

## 3. Configuration Management

These commands focus on interacting with M3TAL's environment variables and configuration files.

### `m3tal config wizard`

**Description:**
Starts an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file. This is the recommended way to set up your environment variables.

**Usage Example:**
```bash
m3tal config wizard
```
*(Follow the on-screen prompts to configure various M3TAL settings.)*

### `m3tal config set KEY VALUE`

**Description:**
Sets a single environment variable in your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config set DOMAIN my-m3tal-domain.com
```

### `m3tal config get KEY`

**Description:**
Reads and displays the value of a specific environment variable from your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config get DASHBOARD_PORT
```
*(This might output `8082` if that's the configured value.)*

### `m3tal config scan`

**Description:**
Lists all environment variables across all configured M3TAL stacks, including their current values and any default values.

**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`

**Description:**
Displays the entire contents of your current `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config list
```

---

## 4. System Service Management

M3TAL's API daemon is managed via `systemd`.

### `systemctl status m3tal-api`

**Description:**
Checks the current status of the `m3tal-api.service` systemd unit.

**Usage Example:**
```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

**Description:**
Streams the logs from the `m3tal-api.service` in real-time. Use this for detailed debugging of the API daemon.

**Usage Example:**
```bash
journalctl -u m3tal-api -f
```
*(Press `Ctrl+C` to stop streaming.)*

---

## 5. Direct Docker Compose Fallback

In certain scenarios, you might need to interact with Docker Compose directly. M3TAL manages services within the `/docker/` directory.

### `docker compose up`

**Description:**
Manually starts all services defined in `*-compose.yml` files within the current directory (or specified via `-f`). This is what `m3tal up` essentially does.

**Usage Example (from `/docker/` directory):**
```bash
sudo docker compose up -d
```
*(The `-d` flag runs containers in detached mode.)*

### `docker compose down`

**Description:**
Manually stops and removes all services defined in `*-compose.yml` files within the current directory. This is what `m3tal down` essentially does.

**Usage Example (from `/docker/` directory):**
```bash
sudo docker compose down
```

### `docker compose logs [service_name]`

**Description:**
Manually streams logs for a specific service or all services if `service_name` is omitted. This can be an alternative to `m3tal logs`.

**Usage Example:**
```bash
sudo docker compose logs -f m3tal-dashboard
```
*(The `-f` flag streams logs. Press `Ctrl+C` to stop.)*

### `docker compose ps`

**Description:**
Lists the containers for the current Docker Compose project.

**Usage Example:**
```bash
sudo docker compose ps
```

---
