```markdown
# M3TAL CLI Command Reference

This document provides a comprehensive reference for all M3TAL command-line interface (CLI) commands. M3TAL is a unified Go binary that serves as the primary entry point for managing your M3TAL ecosystem.

## Table of Contents

- [Core M3TAL Commands](#core-m3tal-commands)
  - [`sudo m3tal`](#sudo-m3tal)
  - [`m3tal init`](#m3tal-init)
  - [`m3tal doctor`](#m3tal-doctor)
- [Configuration Management](#configuration-management)
  - [`m3tal config wizard`](#m3tal-config-wizard)
  - [`m3tal config set KEY VALUE`](#m3tal-config-set-key-value)
  - [`m3tal config get KEY`](#m3tal-config-get-key)
  - [`m3tal config scan`](#m3tal-config-scan)
  - [`m3tal config list`](#m3tal-config-list)
- [Dashboard Management](#dashboard-management)
  - [`m3tal dashpass [username] [password]`](#m3tal-dashpass-username-password)
  - [`m3tal dash up`](#m3tal-dash-up)
  - [`m3tal dash down`](#m3tal-dash-down)
  - [`m3tal dash restart`](#m3tal-dash-restart)
  - [`m3tal dash logs`](#m3tal-dash-logs)
  - [`m3tal dash status`](#m3tal-dash-status)
- [Stack Management](#stack-management)
  - [`m3tal up`](#m3tal-up)
  - [`m3tal down`](#m3tal-down)
  - [`m3tal logs`](#m3tal-logs)
- [Systemd Service Management](#systemd-service-management)
  - [`systemctl status m3tal-api`](#systemctl-status-m3tal-api)
  - [`journalctl -u m3tal-api -f`](#journalctl--u-m3tal-api--f)
- [Direct Docker Compose Fallback](#direct-docker-compose-fallback)

---

## Core M3TAL Commands

### `sudo m3tal`

Opens the interactive TUI (Text User Interface) Control Center, presenting a numbered menu for various operations.

**Usage:**
```bash
sudo m3tal
```

---

### `m3tal init`

Generates the default `/etc/m3tal/.env` configuration file. This command is essential for first-time installations to set up basic configuration.

**Usage:**
```bash
m3tal init
```

---

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL environment. This includes verifying Docker connectivity, checking the validity of your `.env` file, and ensuring necessary ports are available.

**Usage:**
```bash
m3tal doctor
```

---

## Configuration Management

### `m3tal config wizard`

Launches an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file. This is the recommended method for modifying `.env` settings.

**Usage:**
```bash
m3tal config wizard
```

---

### `m3tal config set KEY VALUE`

Sets a single environment variable in your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config set DASHBOARD_PORT 8082
```

---

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config get DOMAIN
```

---

### `m3tal config scan`

Lists all environment variables across all managed stacks, including their current values and defaults.

**Usage:**
```bash
m3tal config scan
```

---

### `m3tal config list`

Displays the current contents of your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

---

## Dashboard Management

The M3TAL dashboard provides a web-based interface for managing your system. Its access mode is controlled by `DASHBOARD_EXPOSE_MODE` in your `.env` file (`local` or `traefik`).

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt for them interactively.

**Usage (interactive):**
```bash
sudo m3tal dashpass
# Prompts for username and password
```

**Usage (with arguments):**
```bash
sudo m3tal dashpass admin new_secure_password
```

---

### `m3tal dash up`

Pulls the latest `m3tal-compose.yml` and associated override files from GitHub and then starts the `m3tal-dashboard` Docker container. It respects the `DASHBOARD_EXPOSE_MODE` setting in your `.env` file.

**Usage:**
```bash
sudo m3tal dash up
```

---

### `m3tal dash down`

Stops the `m3tal-dashboard` Docker container.

**Usage:**
```bash
sudo m3tal dash down
```

---

### `m3tal dash restart`

Restarts the `m3tal-dashboard` Docker container.

**Usage:**
```bash
sudo m3tal dash restart
```

---

### `m3tal dash logs`

Streams the logs from the `m3tal-dashboard` container in real-time.

**Usage:**
```bash
sudo m3tal dash logs
```

---

### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` Docker container.

**Usage:**
```bash
sudo m3tal dash status
```

---

## Stack Management

These commands manage all your deployed M3TAL stacks via Docker Compose.

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This starts all your services defined in the stack files.

**Usage:**
```bash
sudo m3tal up
```

---

### `m3tal down`

Runs `docker compose down` across all managed stacks. This stops and removes all containers, networks, and volumes defined in your compose files.

**Usage:**
```bash
sudo m3tal down
```

---

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks. This provides a unified view of your system's output.

**Usage:**
```bash
sudo m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api.service`.

**Usage:**
```bash
systemctl status m3tal-api
```

---

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time. Use `Ctrl+C` to exit the log stream.

**Usage:**
```bash
journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Fallback

In situations where the `m3tal` CLI might not be sufficient or for direct control, you can use Docker Compose commands directly within the `/docker/` directory.

**Prerequisites:**
- Ensure Docker and Docker Compose V2 are installed and running.
- Navigate to the `/docker/` directory.

**Common Commands:**

- **Start all services:**
  ```bash
  cd /docker
  sudo docker compose up -d
  ```

- **Stop all services:**
  ```bash
  cd /docker
  sudo docker compose down
  ```

- **View logs for a specific service (e.g., `m3tal-dashboard`):**
  ```bash
  cd /docker
  sudo docker compose logs -f m3tal-dashboard
  ```

- **List running services:**
  ```bash
  cd /docker
  sudo docker compose ps
  ```

- **Build or rebuild services (e.g., for development):**
  ```bash
  cd /docker
  sudo docker compose build
  ```
```