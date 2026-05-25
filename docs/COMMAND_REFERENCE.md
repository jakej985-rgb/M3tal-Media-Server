# M3TAL CLI Command Reference

This document provides a comprehensive cheat sheet for the M3TAL Command Line Interface (CLI). It covers all available commands and their typical usage.

## Table of Contents

1.  [Getting Started & Initialization](#getting-started--initialization)
2.  [Configuration Management](#configuration-management)
3.  [Dashboard Management](#dashboard-management)
4.  [System & Stack Management](#system--stack-management)
5.  [Systemd Service Management](#systemd-service-management)
6.  [Direct Docker Compose Commands (Fallback)](#direct-docker-compose-commands-fallback)

---

## Getting Started & Initialization

### `sudo m3tal`

Opens the interactive TUI Control Center, presenting a numbered menu for navigating M3TAL's features.

**Usage:**
```bash
sudo m3tal
```
*(This command is interactive and will present a menu.)*

### `m3tal init`

Generates the default `/etc/m3tal/.env` configuration file. This command should be used on the first install or when you need to reset to default configurations.

**Usage:**
```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL system. It verifies:
*   Docker connectivity.
*   The validity of your `/etc/m3tal/.env` file.
*   The availability of necessary ports.

**Usage:**
```bash
m3tal doctor
```

---

## Configuration Management

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring your `/etc/m3tal/.env` file. This is the recommended way to set up M3TAL's environment variables.

**Usage:**
```bash
m3tal config wizard
```
*(This command is interactive and will prompt for input.)*

### `m3tal config set KEY VALUE`

Sets a single environment variable in your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config set DOMAIN mym3tal.local
```

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config get DASHBOARD_PORT
```

### `m3tal config scan`

Lists all environment variables across all defined M3TAL stacks. This command provides a consolidated view of your configuration.

**Usage:**
```bash
m3tal config scan
```

### `m3tal config list`

Lists the current contents of your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

---

## Dashboard Management

The M3TAL dashboard provides a web-based interface for managing your M3TAL system. Its access mode (local or Traefik) is controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt for them interactively.

**Usage Examples:**
```bash
# Interactive mode
m3tal dashpass

# Set password for a specific user
m3tal dashpass admin new_secure_password
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container. This command ensures you are running the most up-to-date dashboard setup.

**Usage:**
```bash
m3tal dash up
```

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage:**
```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the M3TAL dashboard container.

**Usage:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container in real-time. Press `Ctrl+C` to stop streaming.

**Usage:**
```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container (e.g., running, exited).

**Usage:**
```bash
m3tal dash status
```

---

## System & Stack Management

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This command starts all your defined M3TAL stacks and services.

**Usage:**
```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` for all stacks managed by M3TAL. This command stops all running services.

**Usage:**
```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL stacks. This provides a centralized view of your system's activity. Press `Ctrl+C` to stop streaming.

**Usage:**
```bash
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon (`m3tal-api.service`) is managed by `systemd`.

### `systemctl status m3tal-api`

Checks the current status of the M3TAL API service.

**Usage:**
```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs from the M3TAL API service in real-time. This is useful for debugging issues with the backend API. Press `Ctrl+C` to stop streaming.

**Usage:**
```bash
journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Commands (Fallback)

In situations where the M3TAL CLI might not cover a specific need, you can directly interact with Docker Compose. M3TAL uses Docker Compose V2.

### Running all stacks:

**Start all services:**
```bash
docker compose -f /docker/*.yml up -d
```

**Stop all services:**
```bash
docker compose -f /docker/*.yml down
```

**View logs for all services:**
```bash
docker compose -f /docker/*.yml logs -f
```

### Managing specific stacks:

Navigate to the directory containing your custom compose file (e.g., `/docker/my-custom-stack/`) and use `docker compose`:

**Start a specific stack (e.g., `my-stack.yml`):**
```bash
docker compose -f /docker/my-stack.yml up -d
```

**Stop a specific stack:**
```bash
docker compose -f /docker/my-stack.yml down
```

**View logs for a specific stack:**
```bash
docker compose -f /docker/my-stack.yml logs -f
```

---
This concludes the M3TAL CLI command reference. Always refer to the official M3TAL documentation for the most up-to-date information and advanced usage.