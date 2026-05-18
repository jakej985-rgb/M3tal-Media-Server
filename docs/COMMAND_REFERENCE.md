# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for all available M3TAL CLI commands. It provides a quick reference for common operations, configuration, and system management tasks.

## Table of Contents

1.  [Interactive Control Center](#interactive-control-center)
2.  [Initialization & Configuration](#initialization--configuration)
    *   [m3tal init](#m3tal-init)
    *   [m3tal config wizard](#m3tal-config-wizard)
    *   [m3tal config set KEY VALUE](#m3tal-config-set-key-value)
    *   [m3tal config get KEY](#m3tal-config-get-key)
    *   [m3tal config scan](#m3tal-config-scan)
    *   [m3tal config list](#m3tal-config-list)
3.  [System Health & Diagnostics](#system-health--diagnostics)
    *   [m3tal doctor](#m3tal-doctor)
4.  [Dashboard Management](#dashboard-management)
    *   [m3tal dashpass [username] [password]](#m3tal-dashpass-username-password)
    *   [m3tal dash up](#m3tal-dash-up)
    *   [m3tal dash down](#m3tal-dash-down)
    *   [m3tal dash restart](#m3tal-dash-restart)
    *   [m3tal dash logs](#m3tal-dash-logs)
    *   [m3tal dash status](#m3tal-dash-status)
5.  [Stack Management](#stack-management)
    *   [m3tal up](#m3tal-up)
    *   [m3tal down](#m3tal-down)
    *   [m3tal logs](#m3tal-logs)
6.  [Systemd Service Management](#systemd-service-management)
    *   [systemctl status m3tal-api](#systemctl-status-m3tal-api)
    *   [journalctl -u m3tal-api -f](#journalctl--u-m3tal-api--f)
7.  [Direct Docker Compose Commands (Fallback)](#direct-docker-compose-commands-fallback)

---

## Interactive Control Center

The `sudo m3tal` command launches the interactive TUI (Text User Interface) Control Center, offering a menu-driven approach to managing your M3TAL environment.

```bash
sudo m3tal
```

**Example Usage:**
Running `sudo m3tal` will present a numbered menu allowing you to select operations like initializing, configuring, starting/stopping stacks, and managing the dashboard.

---

## Initialization & Configuration

These commands are used for setting up and managing the M3TAL configuration, primarily through the `/etc/m3tal/.env` file.

### `m3tal init`

Generates the `/etc/m3tal/.env` file with default values. Use this on your first installation to establish the base configuration.

```bash
m3tal init
```

**Example Usage:**
After a fresh installation, run `m3tal init` to create the initial configuration file.

### `m3tal config wizard`

Starts an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file. This is the recommended method for comprehensive configuration.

```bash
m3tal config wizard
```

**Example Usage:**
When you need to adjust multiple settings, run `m3tal config wizard` and follow the prompts to update your environment variables.

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file.

```bash
m3tal config set DOMAIN mym3tal.local
```

**Example Usage:**
To change your primary domain name to `mym3tal.local`, you would execute: `m3tal config set DOMAIN mym3tal.local`.

### `m3tal config get KEY`

Reads and displays the current value of a single environment variable from the `/etc/m3tal/.env` file.

```bash
m3tal config get DASHBOARD_PORT
```

**Example Usage:**
To check the current port for the M3TAL dashboard, run: `m3tal config get DASHBOARD_PORT`. This would likely output `8082`.

### `m3tal config scan`

Lists all environment variables across all defined M3TAL stacks, showing their current values.

```bash
m3tal config scan
```

**Example Usage:**
To get a comprehensive overview of all configuration variables and their current settings across your M3TAL setup, run `m3tal config scan`.

### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file.

```bash
m3tal config list
```

**Example Usage:**
To view the complete configuration file, execute: `m3tal config list`.

---

## System Health & Diagnostics

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL environment. This includes verifying Docker connectivity, the validity of your `.env` file, and checking for port availability conflicts.

```bash
m3tal doctor
```

**Example Usage:**
Before deploying new stacks or troubleshooting issues, run `m3tal doctor` to ensure your system is healthy.

---

## Dashboard Management

Commands specifically for managing the M3TAL Dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt you interactively.

```bash
m3tal dashpass m3taladmin newsecurepassword123
```

**Example Usage:**
To change the password for the user `m3taladmin` to `newsecurepassword123`, use: `m3tal dashpass m3taladmin newsecurepassword123`. If run without arguments, it will prompt for the username and new password.

### `m3tal dash up`

Pulls the latest dashboard Compose configuration from GitHub and then starts the `m3tal-dashboard` container.

```bash
m3tal dash up
```

**Example Usage:**
To ensure you have the latest dashboard version and to start it, run `m3tal dash up`.

### `m3tal dash down`

Stops the `m3tal-dashboard` container.

```bash
m3tal dash down
```

**Example Usage:**
To stop the dashboard service, execute: `m3tal dash down`.

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container.

```bash
m3tal dash restart
```

**Example Usage:**
If the dashboard is unresponsive, try restarting it with: `m3tal dash restart`.

### `m3tal dash logs`

Streams the logs from the `m3tal-dashboard` container in real-time.

```bash
m3tal dash logs
```

**Example Usage:**
To view live logs from the dashboard and troubleshoot issues, run: `m3tal dash logs`. Press `Ctrl+C` to stop streaming.

### `m3tal dash status`

Shows the current status of the `m3tal-dashboard` container (e.g., running, exited).

```bash
m3tal dash status
```

**Example Usage:**
To check if the dashboard is running and its current state, execute: `m3tal dash status`.

---

## Stack Management

These commands manage all your deployed Docker stacks using Docker Compose.

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This starts all your defined services.

```bash
m3tal up
```

**Example Usage:**
After adding a new service or making changes to your stack configurations, run `m3tal up` to bring all services online.

### `m3tal down`

Runs `docker compose down` across all stacks managed by M3TAL, effectively stopping all running services.

```bash
m3tal down
```

**Example Usage:**
To stop all M3TAL managed services, execute: `m3tal down`.

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL stacks.

```bash
m3tal logs
```

**Example Usage:**
To monitor the logs of all your running services simultaneously, use: `m3tal logs`. Press `Ctrl+C` to stop streaming.

---

## Systemd Service Management

M3TAL's core API daemon is managed by `systemd`.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api.service` unit.

```bash
systemctl status m3tal-api
```

**Example Usage:**
To verify if the M3TAL API daemon is running correctly, execute: `systemctl status m3tal-api`.

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time. The `-f` flag follows the log output.

```bash
journalctl -u m3tal-api -f
```

**Example Usage:**
To view live logs from the M3TAL API daemon and diagnose issues, run: `journalctl -u m3tal-api -f`. Press `Ctrl+C` to stop streaming.

---

## Direct Docker Compose Commands (Fallback)

In cases where direct control is needed or for advanced troubleshooting, you can use standard Docker Compose commands. M3TAL orchestrates these commands based on the files in `/docker/`.

**Example Usage:**
To directly start all services defined in `/docker/`:
```bash
docker compose -f /docker/*.yml up -d
```

**Example Usage:**
To directly stop all services defined in `/docker/`:
```bash
docker compose -f /docker/*.yml down
```

**Example Usage:**
To view the logs of a specific container, e.g., `m3tal-dashboard`:
```bash
docker compose -f /docker/m3tal-compose.yml logs m3tal-dashboard
```