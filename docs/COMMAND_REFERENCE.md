# M3TAL CLI Command Reference

This document provides a comprehensive cheat sheet for all available `m3tal` CLI commands.

## Table of Contents

*   [Interactive TUI](#interactive-tui)
*   [M3TAL Initialization](#m3tal-initialization)
*   [System Health and Diagnostics](#system-health-and-diagnostics)
*   [Configuration Management](#configuration-management)
    *   [Wizard Configuration](#wizard-configuration)
    *   [Setting a Single Environment Variable](#setting-a-single-environment-variable)
    *   [Getting a Single Environment Variable](#getting-a-single-environment-variable)
    *   [Scanning All Environment Variables](#scanning-all-environment-variables)
    *   [Listing Current .env File Contents](#listing-current-env-file-contents)
*   [Dashboard Management](#dashboard-management)
    *   [Updating Dashboard Password](#updating-dashboard-password)
    *   [Starting/Updating Dashboard](#startingupdating-dashboard)
    *   [Stopping Dashboard](#stopping-dashboard)
    *   [Restarting Dashboard](#restarting-dashboard)
    *   [Viewing Dashboard Logs](#viewing-dashboard-logs)
    *   [Checking Dashboard Status](#checking-dashboard-status)
*   [Stack Management](#stack-management)
    *   [Bringing Up All Stacks](#bringing-up-all-stacks)
    *   [Bringing Down All Stacks](#bringing-down-all-stacks)
    *   [Streaming Aggregated Logs](#streaming-aggregated-logs)
*   [Systemd Service Management](#systemd-service-management)
*   [Direct Docker Compose Commands (Fallback)](#direct-docker-compose-commands-fallback)

---

## Interactive TUI

This command launches the M3TAL TUI Control Center, providing a menu-driven interface for common operations.

```bash
sudo m3tal
```

---

## M3TAL Initialization

Generates the default `/etc/m3tal/.env` configuration file. This is typically run once during the initial installation.

```bash
m3tal init
```

---

## System Health and Diagnostics

Performs a pre-flight health check of your M3TAL environment. This includes verifying Docker connectivity, the validity of your `.env` file, and checking for port availability.

```bash
m3tal doctor
```

---

## Configuration Management

### Wizard Configuration

Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.

```bash
m3tal config wizard
```

### Setting a Single Environment Variable

Sets a specific environment variable within your `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name and `VALUE` with its desired value.

```bash
m3tal config set DOMAIN localhost
```

### Getting a Single Environment Variable

Reads and displays the value of a specific environment variable from your `/etc/m3tal/.env` file.

```bash
m3tal config get LOG_LEVEL
```

### Scanning All Environment Variables

Lists all environment variables across all managed stacks, showing their current values.

```bash
m3tal config scan
```

### Listing Current .env File Contents

Displays the complete contents of the current `/etc/m3tal/.env` file.

```bash
m3tal config list
```

---

## Dashboard Management

### Updating Dashboard Password

Updates the password for a dashboard user. If no username and password are provided, the command will prompt for them interactively.

```bash
# Interactive mode (prompts for username and password)
m3tal dashpass

# With username and password arguments
m3tal dashpass admin new_secure_password
```

### Starting/Updating Dashboard

Pulls the latest `m3tal-compose.yml` and its relevant override (e.g., `m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) from GitHub and then starts or restarts the dashboard container.

```bash
m3tal dash up
```

### Stopping Dashboard

Stops the `m3tal-dashboard` container.

```bash
m3tal dash down
```

### Restarting Dashboard

Restarts the `m3tal-dashboard` container.

```bash
m3tal dash restart
```

### Viewing Dashboard Logs

Streams the logs from the `m3tal-dashboard` container in real-time.

```bash
m3tal dash logs
```

### Checking Dashboard Status

Shows the current status of the `m3tal-dashboard` container.

```bash
m3tal dash status
```

---

## Stack Management

### Bringing Up All Stacks

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This starts all your defined services.

```bash
m3tal up
```

### Bringing Down All Stacks

Runs `docker compose down` for all stacks managed by M3TAL, stopping all running services.

```bash
m3tal down
```

### Streaming Aggregated Logs

Streams aggregated logs from all currently running M3TAL stacks.

```bash
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl` and view its logs with `journalctl`.

**Check the status of the M3TAL API service:**

```bash
systemctl status m3tal-api
```

**Restart the M3TAL API service:**

```bash
sudo systemctl restart m3tal-api
```

**View real-time logs for the M3TAL API service:**

```bash
journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Commands (Fallback)

In situations where you need more granular control or are debugging, you can directly use Docker Compose commands within your stack directories. M3TAL orchestrates these commands for you, but direct usage is possible.

**Navigate to a specific stack directory (e.g., `/docker/my-stack`)**

```bash
cd /docker/my-stack
```

**Start a specific stack:**

```bash
docker compose -f my-stack-compose.yml up -d
```

**Stop a specific stack:**

```bash
docker compose -f my-stack-compose.yml down
```

**View logs for a specific stack:**

```bash
docker compose -f my-stack-compose.yml logs -f
```

**Example for the M3TAL dashboard in local mode (using its default compose files):**

```bash
# Start dashboard in local mode
cd /opt/m3tal/stack/
docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml up -d m3tal-dashboard

# Stop dashboard
cd /opt/m3tal/stack/
docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml down m3tal-dashboard
```