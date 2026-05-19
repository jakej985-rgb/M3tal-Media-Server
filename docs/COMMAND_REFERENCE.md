# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for all M3TAL CLI commands, designed for quick reference and efficient system management.

---

## Table of Contents

1.  [Core M3TAL Commands](#core-m3tal-commands)
2.  [Dashboard Management Commands](#dashboard-management-commands)
3.  [Configuration Management Commands](#configuration-management-commands)
4.  [Systemd Service Management](#systemd-service-management)
5.  [Direct Docker Compose Fallback](#direct-docker-compose-fallback)
6.  [APT Installation](#apt-installation)

---

## Core M3TAL Commands

These commands are the primary interface for managing your M3TAL environment.

### `sudo m3tal`

Opens the interactive TUI Control Center, presenting a numbered menu for various operations.

**Usage:**

```bash
sudo m3tal
```

### `m3tal init`

Generates the `/etc/m3tal/.env` file from default settings. This command is crucial for first-time installations or when resetting configuration to defaults.

**Usage:**

```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, the validity of your `/etc/m3tal/.env` file, and checks for port availability.

**Usage:**

```bash
m3tal doctor
```

### `m3tal up`

Starts all services defined in `*-compose.yml` files within the `/docker/` directory using `docker compose`. This command is used to bring up your entire M3TAL stack, including any custom services you've added.

**Usage:**

```bash
m3tal up
```

### `m3tal down`

Stops and removes all services defined in `*-compose.yml` files within the `/docker/` directory using `docker compose`. This command gracefully shuts down your entire M3TAL stack.

**Usage:**

```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks. This is invaluable for monitoring the health and activity of your services.

**Usage:**

```bash
m3tal logs
```

---

## Dashboard Management Commands

These commands specifically manage the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt you interactively.

**Usage (interactive):**

```bash
m3tal dashpass
```

**Usage (with arguments):**

```bash
m3tal dashpass myuser newsecurepassword123
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the `m3tal-dashboard` container. This ensures you're running the most up-to-date dashboard service.

**Usage:**

```bash
m3tal dash up
```

### `m3tal dash down`

Stops the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams logs specifically from the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` container.

**Usage:**

```bash
m3tal dash status
```

---

## Configuration Management Commands

These commands allow you to manage the M3TAL environment variables stored in `/etc/m3tal/.env`.

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file. This is the recommended way to set up your environment variables initially.

**Usage:**

```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env`.

**Usage:**

```bash
m3tal config set DASHBOARD_PORT 8085
```

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.

**Usage:**

```bash
m3tal config get DOMAIN
```

### `m3tal config scan`

Lists all environment variables across all M3TAL stacks, including their current values and default values if applicable.

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

---

## Systemd Service Management

M3TAL's API daemon runs as a systemd service. You can manage it and view its logs using standard `systemctl` and `journalctl` commands.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api` systemd service.

**Usage:**

```bash
systemctl status m3tal-api
```

### `systemctl restart m3tal-api`

Restarts the `m3tal-api` systemd service.

**Usage:**

```bash
systemctl restart m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs for the `m3tal-api` service in real-time. Press `Ctrl+C` to exit.

**Usage:**

```bash
journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Fallback

In some advanced scenarios, or for debugging, you may need to interact directly with Docker Compose. M3TAL utilizes Docker Compose V2.

Ensure you are in the `/docker/` directory or specify the compose file path.

### `docker compose up`

Starts services defined in `docker-compose.yml` files.

**Usage (in `/docker/`):**

```bash
docker compose up -d
```

**Usage (specific file):**

```bash
docker compose -f /docker/my-custom-stack.yml up -d
```

### `docker compose down`

Stops and removes containers, networks, and volumes.

**Usage (in `/docker/`):**

```bash
docker compose down
```

**Usage (specific file):**

```bash
docker compose -f /docker/my-custom-stack.yml down
```

### `docker compose ps`

Lists the containers for the current Compose project.

**Usage (in `/docker/`):**

```bash
docker compose ps
```

---

## APT Installation

This section details the commands required to install M3TAL via APT.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update package list and install
sudo apt update && sudo apt install -y m3tal
```