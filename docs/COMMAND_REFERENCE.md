# M3TAL CLI Command Reference

This document provides a comprehensive reference for all available M3TAL CLI commands, serving as a cheat-sheet for managing your M3TAL ecosystem.

## Getting Started

### APT Installation

To install or update M3TAL, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Interactive Control Center

### `sudo m3tal`

This command launches the M3TAL Interactive TUI Control Center, providing a numbered menu for common operations.

**Usage:**
```bash
sudo m3tal
```

## Initialization and Configuration

### `m3tal init`

Generates the default `/etc/m3tal/.env` configuration file. This should be run upon first installation.

**Usage:**
```bash
m3tal init
```

### `m3tal config wizard`

Initiates an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config set DASHBOARD_PORT 8083
```

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from the `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config get DOMAIN
```

### `m3tal config scan`

Lists all environment variables across all configured stacks.

**Usage:**
```bash
m3tal config scan
```

### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

## Dashboard Management

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If username and password are not provided, the command will prompt interactively.

**Usage (interactive):**
```bash
m3tal dashpass
```

**Usage (with arguments):**
```bash
m3tal dashpass admin new_secure_password
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and starts the dashboard container.

**Usage:**
```bash
m3tal dash up
```

### `m3tal dash down`

Stops the dashboard container.

**Usage:**
```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the dashboard container.

**Usage:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams logs from the dashboard container in real-time.

**Usage:**
```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the dashboard container.

**Usage:**
```bash
m3tal dash status
```

## Stack Management

### `m3tal up`

Starts all services defined in `*-compose.yml` files located in `/docker/`.

**Usage:**
```bash
m3tal up
```

### `m3tal down`

Stops all services managed by M3TAL's docker compose configurations.

**Usage:**
```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL stacks.

**Usage:**
```bash
m3tal logs
```

## Systemd Service Management

The M3TAL API daemon is managed by systemd. You can monitor and control its state using `systemctl` and `journalctl`.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api` service.

**Usage:**
```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs for the `m3tal-api` service in real-time.

**Usage:**
```bash
journalctl -u m3tal-api -f
```

## Direct Docker Compose Commands (Fallback)

In cases where direct control is needed, you can use Docker Compose commands on the stack compose files. M3TAL primarily uses files within the `/docker/` directory.

### `docker compose up`

Starts services defined in `*-compose.yml` files within the current directory or specified context.

**Usage (example for all stacks):**
```bash
docker compose -f /docker/routing-compose.yml -f /docker/my-app-compose.yml up -d
```

### `docker compose down`

Stops and removes containers, networks, and volumes defined in compose files.

**Usage (example for all stacks):**
```bash
docker compose -f /docker/routing-compose.yml -f /docker/my-app-compose.yml down
```

### `docker compose logs`

Displays logs from services.

**Usage (example for a specific service):**
```bash
docker compose -f /docker/my-app-compose.yml logs my-service-name
```

### `docker compose ps`

Lists containers for the services defined in the compose file.

**Usage:**
```bash
docker compose -f /docker/routing-compose.yml ps
```

---
**Note:** Always ensure you are in the correct directory or provide the full path to your `docker-compose.yml` files when using direct `docker compose` commands. The `m3tal` CLI abstracts these interactions for convenience.