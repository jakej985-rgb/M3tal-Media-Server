# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for the M3TAL CLI, covering all available commands and their usage.

## Table of Contents

*   [Interactive Control Center](#interactive-control-center)
*   [Initialization and Configuration](#initialization-and-configuration)
*   [Dashboard Management](#dashboard-management)
*   [Stack Management](#stack-management)
*   [System Service Management](#system-service-management)
*   [Direct Docker Compose Commands](#direct-docker-compose-commands)

---

## Interactive Control Center

The `sudo m3tal` command launches the M3TAL TUI Control Center, offering an interactive, numbered menu to navigate and execute various M3TAL functions.

```bash
sudo m3tal
```

This command will present a menu like the following (actual options may vary):

```
Welcome to M3TAL Control Center!

1. Initialize M3TAL configuration
2. Run M3TAL configuration wizard
3. Manage M3TAL .env configuration
4. Manage M3TAL Dashboard
5. Deploy/Undeploy M3TAL stacks
6. View system logs
7. Run pre-flight health check
8. Exit

Enter your choice:
```

---

## Initialization and Configuration

These commands manage the core M3TAL configuration, primarily the `/etc/m3tal/.env` file.

### `m3tal init`

Generates the `/etc/m3tal/.env` file with default values. This is typically run once after installing M3TAL.

**Usage:**

```bash
m3tal init
```

**Example:**

```bash
sudo m3tal init
# Output: .env file generated at /etc/m3tal/.env with default values.
```

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file. This is the recommended method for initial or significant configuration changes.

**Usage:**

```bash
m3tal config wizard
```

**Example:**

```bash
sudo m3tal config wizard
# This will prompt you for various configuration values like DOMAIN, PUID, PGID, etc.
# Example interaction:
# ? Enter your desired domain (e.g., your.domain.com or localhost): localhost
# ? Enter the User ID (PUID) for running containers (default: 1000): 1000
# ? Enter the Group ID (PGID) for running containers (default: 1000): 1000
# Configuration updated.
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config set <KEY> <VALUE>
```

**Example:**

```bash
sudo m3tal config set DOMAIN mylocal.m3tal
# Output: Setting DOMAIN to mylocal.m3tal in /etc/m3tal/.env
```

### `m3tal config get KEY`

Reads a single environment variable from the `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config get <KEY>
```

**Example:**

```bash
sudo m3tal config get DOMAIN
# Output: localhost
```

### `m3tal config scan`

Lists all environment variables across all managed stacks. This command provides a consolidated view of configurable settings.

**Usage:**

```bash
m3tal config scan
```

**Example:**

```bash
sudo m3tal config scan
# Output will be a JSON array like:
# [
#   {
#     "key": "DASHBOARD_PORT",
#     "default": "8082",
#     "value": "8082"
#   },
#   {
#     "key": "DASHBOARD_EXPOSE_MODE",
#     "default": "local",
#     "value": "local"
#   },
#   ...
# ]
```

### `m3tal config list`

Lists the current contents of the `/etc/m3tal/.env` file.

**Usage:**

```bash
m3tal config list
```

**Example:**

```bash
sudo m3tal config list
# Output:
# DASHBOARD_PORT=8082
# DASHBOARD_EXPOSE_MODE=local
# DOMAIN=localhost
# ...
```

---

## Dashboard Management

Commands for managing the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, it will prompt interactively.

**Usage:**

```bash
m3tal dashpass [username] [password]
```

**Example (Interactive):**

```bash
sudo m3tal dashpass
# Output:
# Enter dashboard username: admin
# Enter new password:
# Confirm new password:
# Password for user 'admin' updated successfully.
```

**Example (Directly):**

```bash
sudo m3tal dashpass admin MySuperSecurePassword123
# Output: Password for user 'admin' updated successfully.
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Usage:**

```bash
m3tal dash up
```

**Example:**

```bash
sudo m3tal dash up
# Output:
# Pulling latest dashboard compose configurations from GitHub...
# Starting m3tal-dashboard container...
# Container m3tal-dashboard is running.
```

### `m3tal dash down`

Stops the dashboard container.

**Usage:**

```bash
m3tal dash down
```

**Example:**

```bash
sudo m3tal dash down
# Output: Stopping m3tal-dashboard container...
# Container m3tal-dashboard stopped.
```

### `m3tal dash restart`

Restarts the dashboard container.

**Usage:**

```bash
m3tal dash restart
```

**Example:**

```bash
sudo m3tal dash restart
# Output: Restarting m3tal-dashboard container...
# Container m3tal-dashboard restarted.
```

### `m3tal dash logs`

Streams the logs from the dashboard container in real-time.

**Usage:**

```bash
m3tal dash logs
```

**Example:**

```bash
sudo m3tal dash logs
# Output will show live logs from the m3tal-dashboard container:
# 2023-10-27 10:00:00 INFO: Starting Flask server on port 8082
# 2023-10-27 10:00:05 INFO: Dashboard loaded successfully.
# ...
```

Press `Ctrl+C` to stop streaming logs.

### `m3tal dash status`

Shows the current status of the dashboard container.

**Usage:**

```bash
m3tal dash status
```

**Example:**

```bash
sudo m3tal dash status
# Output:
# m3tal-dashboard is running (created 2 days ago, restarted 5 times)
```
or
```bash
sudo m3tal dash status
# Output:
# m3tal-dashboard is stopped.
```

---

## Stack Management

Commands for deploying and managing all M3TAL stacks defined by compose files in `/docker/`.

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in `/docker/`. This deploys all defined M3TAL stacks.

**Usage:**

```bash
m3tal up
```

**Example:**

```bash
sudo m3tal up
# Output will show Docker Compose starting all services defined in files within /docker/:
# [+] Running 2/2
# ✔ Network m3tal_proxy  Created                                                     0.2s
# ✔ Container traefik      Started                                                     1.5s
# ✔ Container m3tal-dashboard Started                                                 2.0s
# ...
```

### `m3tal down`

Runs `docker compose down` across all stacks defined in `/docker/`. This stops and removes all containers, networks, and volumes associated with the deployed stacks.

**Usage:**

```bash
m3tal down
```

**Example:**

```bash
sudo m3tal down
# Output will show Docker Compose stopping and removing all services:
# [+] Running 2/2
# ✔ Container m3tal-dashboard Removed                                               0.5s
# ✔ Container traefik         Removed                                               1.0s
# ✔ Network m3tal_proxy       Removed                                               1.2s
# ...
```

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks. This provides a unified log view.

**Usage:**

```bash
m3tal logs
```

**Example:**

```bash
sudo m3tal logs
# Output will show live logs from all running containers:
# traefik      | time="2023-10-27T10:05:00Z" level=info msg="...traefik logs..."
# m3tal-dashboard | 2023-10-27 10:05:01 INFO: Processing request...
# ...
```

Press `Ctrl+C` to stop streaming logs.

---

## System Service Management

M3TAL utilizes systemd to manage its core API daemon.

### `systemctl status m3tal-api`

Checks the status of the M3TAL API service.

**Usage:**

```bash
systemctl status m3tal-api
```

**Example:**

```bash
sudo systemctl status m3tal-api
# Output will show if the service is active, loaded, and running:
# ● m3tal-api.service - M3TAL API Daemon
#      Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
#      Active: active (running) since Fri 2023-10-27 10:00:00 UTC; 2 days ago
#        Docs: https://jakej985-rgb.github.io/m3tal-core/docs/
#      Main PID: 1234 (m3tal-api)
#         Tasks: 10 (limit: 4915)
#        Memory: 50.0M
#           CPU: 1min 30s
#      CGroup: /system.slice/m3tal-api.service
#              └─1234 /usr/bin/m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time.

**Usage:**

```bash
journalctl -u m3tal-api -f
```

**Example:**

```bash
sudo journalctl -u m3tal-api -f
# Output will show live logs from the M3TAL API daemon:
# Oct 27 10:00:00 server m3tal-api[1234]: INFO[0000] Starting M3TAL API daemon on port 8080
# Oct 27 10:00:01 server m3tal-api[1234]: INFO[0001] Connected to Docker daemon
# ...
```

Press `Ctrl+C` to stop streaming logs.

---

## Direct Docker Compose Commands

As a fallback, you can directly use `docker compose` commands within the `/docker/` directory.

**Note:** The `m3tal up` and `m3tal down` commands abstract this for you, but direct usage is available.

### `docker compose up`

Starts all services defined in `*-compose.yml` files in the current directory.

**Usage:**

```bash
cd /docker
sudo docker compose up
```

**Example:**

```bash
cd /docker
sudo docker compose up -d # Use -d for detached mode
# Output:
# [+] Running 2/2
# ✔ Network m3tal_proxy  Created                                                     0.2s
# ✔ Container traefik      Started                                                     1.5s
# ✔ Container m3tal-dashboard Started                                                 2.0s
```

### `docker compose down`

Stops and removes containers, networks, and volumes for services defined in `*-compose.yml` files in the current directory.

**Usage:**

```bash
cd /docker
sudo docker compose down
```

**Example:**

```bash
cd /docker
sudo docker compose down
# Output:
# [+] Running 2/2
# ✔ Container m3tal-dashboard Removed                                               0.5s
# ✔ Container traefik         Removed                                               1.0s
# ✔ Network m3tal_proxy       Removed                                               1.2s
```

---

This document should serve as your primary reference for interacting with the M3TAL system via its command-line interface.