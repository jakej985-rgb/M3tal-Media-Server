# M3TAL CLI Command Reference

This document provides a comprehensive cheat sheet for the M3TAL CLI, covering all available commands and their usage.

---

## Table of Contents

1.  [M3TAL TUI Control Center](#m3tal-tui-control-center)
2.  [Configuration Management](#configuration-management)
    *   [Initialize Configuration](#initialize-configuration)
    *   [Configuration Wizard](#configuration-wizard)
    *   [Set Environment Variable](#set-environment-variable)
    *   [Get Environment Variable](#get-environment-variable)
    *   [Scan All Environment Variables](#scan-all-environment-variables)
    *   [List Current .env File](#list-current-env-file)
3.  [Dashboard Management](#dashboard-management)
    *   [Update Dashboard Password](#update-dashboard-password)
    *   [Start/Update Dashboard](#startupdate-dashboard)
    *   [Stop Dashboard](#stop-dashboard)
    *   [Restart Dashboard](#restart-dashboard)
    *   [View Dashboard Logs](#view-dashboard-logs)
    *   [Dashboard Status](#dashboard-status)
4.  [Stack Management](#stack-management)
    *   [Start All Stacks](#start-all-stacks)
    *   [Stop All Stacks](#stop-all-stacks)
    *   [Aggregate Logs](#aggregate-logs)
5.  [System Health](#system-health)
    *   [Run Pre-flight Checks](#run-pre-flight-checks)
6.  [Systemd Service Management](#systemd-service-management)
7.  [Direct Docker Compose Commands (Fallback)](#direct-docker-compose-commands-fallback)

---

## M3TAL TUI Control Center

The interactive TUI provides a user-friendly menu-driven interface for managing M3TAL.

**Command:**
```bash
sudo m3tal
```

**Example Usage:**
```bash
sudo m3tal
```
This command will launch the M3TAL TUI, presenting a numbered menu for various actions like starting/stopping services, managing configurations, and more.

---

## Configuration Management

### Initialize Configuration

Generates the `/etc/m3tal/.env` file with default values. Use this on the first installation to set up your environment.

**Command:**
```bash
m3tal init
```

**Example Usage:**
```bash
m3tal init
```
This command creates `/etc/m3tal/.env` if it doesn't exist, populating it with the default configuration settings.

### Configuration Wizard

An interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.

**Command:**
```bash
m3tal config wizard
```

**Example Usage:**
```bash
m3tal config wizard
```
This launches an interactive process where you will be prompted to enter values for various configuration parameters, which will then be saved to `/etc/m3tal/.env`.

### Set Environment Variable

Sets a single environment variable within the `/etc/m3tal/.env` file.

**Command:**
```bash
m3tal config set KEY VALUE
```

**Example Usage:**
```bash
m3tal config set DASHBOARD_PORT 8083
```
This command will set the `DASHBOARD_PORT` variable in `/etc/m3tal/.env` to `8083`.

### Get Environment Variable

Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.

**Command:**
```bash
m3tal config get KEY
```

**Example Usage:**
```bash
m3tal config get LOG_LEVEL
```
This command will output the current value of the `LOG_LEVEL` variable from your `.env` file.

### Scan All Environment Variables

Lists all environment variables across all managed stacks, showing their current values and defaults.

**Command:**
```bash
m3tal config scan
```

**Example Usage:**
```bash
m3tal config scan
```
This command will display a comprehensive list of all configuration variables, their current settings, and their default values if not explicitly set.

### List Current .env File

Displays the entire contents of the current `/etc/m3tal/.env` file.

**Command:**
```bash
m3tal config list
```

**Example Usage:**
```bash
m3tal config list
```
This command will print the complete content of your `/etc/m3tal/.env` file to the console.

---

## Dashboard Management

The following commands manage the M3TAL Dashboard container.

### Update Dashboard Password

Updates the password for the dashboard user. If no username or password is provided, it will prompt interactively.

**Command:**
```bash
m3tal dashpass [username] [password]
```

**Example Usage (Interactive):**
```bash
sudo m3tal dashpass
```
This will prompt you to enter a username and then a new password for the dashboard.

**Example Usage (with arguments):**
```bash
sudo m3tal dashpass admin new_secure_password
```
This command sets the password for the `admin` user to `new_secure_password` directly.

### Start/Update Dashboard

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Command:**
```bash
m3tal dash up
```

**Example Usage:**
```bash
sudo m3tal dash up
```
This command ensures you have the latest `m3tal-compose.yml` and related overrides, then starts the `m3tal-dashboard` container.

### Stop Dashboard

Stops the M3TAL dashboard container.

**Command:**
```bash
m3tal dash down
```

**Example Usage:**
```bash
sudo m3tal dash down
```
This command will gracefully stop the `m3tal-dashboard` container.

### Restart Dashboard

Restarts the M3TAL dashboard container.

**Command:**
```bash
m3tal dash restart
```

**Example Usage:**
```bash
sudo m3tal dash restart
```
This command will stop and then start the `m3tal-dashboard` container.

### View Dashboard Logs

Streams the logs from the M3TAL dashboard container in real-time.

**Command:**
```bash
m3tal dash logs
```

**Example Usage:**
```bash
sudo m3tal dash logs
```
This command will continuously display the output from the `m3tal-dashboard` container. Press `Ctrl+C` to exit.

### Dashboard Status

Shows the current status of the M3TAL dashboard container.

**Command:**
```bash
m3tal dash status
```

**Example Usage:**
```bash
sudo m3tal dash status
```
This command will provide information about whether the `m3tal-dashboard` container is running, stopped, or in an error state.

---

## Stack Management

These commands manage all your deployed M3TAL stacks.

### Start All Stacks

Runs `docker compose up` across all `*-compose.yml` files found in `/docker/`.

**Command:**
```bash
m3tal up
```

**Example Usage:**
```bash
sudo m3tal up
```
This command starts all defined services and stacks within your M3TAL environment by processing all compose files in the `/docker/` directory.

### Stop All Stacks

Runs `docker compose down` across all stacks managed by M3TAL.

**Command:**
```bash
m3tal down
```

**Example Usage:**
```bash
sudo m3tal down
```
This command stops and removes all containers, networks, and volumes defined by your M3TAL stacks.

### Aggregate Logs

Streams aggregated logs from all running M3TAL stacks.

**Command:**
```bash
m3tal logs
```

**Example Usage:**
```bash
sudo m3tal logs
```
This command will display a unified stream of logs from all active containers managed by M3TAL. Press `Ctrl+C` to exit.

---

## System Health

### Run Pre-flight Checks

Performs a comprehensive pre-flight health check of your M3TAL installation. This includes verifying Docker connectivity, `.env` file validity, and port availability.

**Command:**
```bash
m3tal doctor
```

**Example Usage:**
```bash
sudo m3tal doctor
```
This command will execute a series of checks and report any issues found with your M3TAL setup.

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl`.

**Command:** `systemctl status m3tal-api`
**Example Usage:**
```bash
sudo systemctl status m3tal-api
```
This command shows the current status of the `m3tal-api.service`.

**Command:** `systemctl restart m3tal-api`
**Example Usage:**
```bash
sudo systemctl restart m3tal-api
```
This command restarts the M3TAL API daemon.

**Command:** `journalctl -u m3tal-api -f`
**Example Usage:**
```bash
sudo journalctl -u m3tal-api -f
```
This command streams the logs from the `m3tal-api.service` in real-time. Press `Ctrl+C` to exit.

---

## Direct Docker Compose Commands (Fallback)

In cases where the M3TAL CLI might not cover a specific Docker Compose operation, you can directly interact with Docker Compose. M3TAL orchestrates services using compose files located in `/docker/`.

**Example: Starting a specific stack**
If you have a stack defined in `/docker/my-app-compose.yml`, you can start it directly.

```bash
cd /docker/
sudo docker compose -f my-app-compose.yml up -d
```

**Example: Stopping a specific stack**
```bash
cd /docker/
sudo docker compose -f my-app-compose.yml down
```

**Example: Viewing logs for a specific service**
Assuming your stack has a service named `my-service`.

```bash
sudo docker compose logs -f my-service
```

**Note:** When using direct `docker compose` commands, ensure you are in the `/docker/` directory or specify the correct compose file path. Remember to use `sudo` if Docker requires root privileges on your system.

---

## APT Installation

To install M3TAL on your system:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```