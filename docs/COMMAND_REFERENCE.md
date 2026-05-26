# M3TAL CLI Command Reference

This document provides a comprehensive reference for all available `m3tal` CLI commands.

---

## Interactive TUI Control Center

The `m3tal` command itself launches an interactive Text-based User Interface (TUI) Control Center, presenting a numbered menu of common operations.

**Usage:**
```bash
sudo m3tal
```

**Example:**
```bash
sudo m3tal
```
*This will launch the TUI, allowing you to navigate and select actions like initializing the system, performing health checks, configuring settings, or managing the dashboard and stacks.*

---

## System Initialization and Configuration

### Initialize Environment Variables

Generates the `/etc/m3tal/.env` file with default settings. Use this on your first installation.

**Usage:**
```bash
m3tal init
```

**Example:**
```bash
m3tal init
```
*After running `m3tal init`, a new `/etc/m3tal/.env` file will be created containing all default environment variables for M3TAL.*

### Pre-flight Health Check

Performs a series of checks to ensure M3TAL is running correctly. This includes verifying Docker connectivity, the validity of your `.env` file, and checking for port availability.

**Usage:**
```bash
m3tal doctor
```

**Example:**
```bash
m3tal doctor
```
*This command will output a report on the system's health, highlighting any potential issues that need addressing before proceeding with M3TAL operations.*

### Configuration Wizard

Launches an interactive wizard to guide you through configuring your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config wizard
```

**Example:**
```bash
m3tal config wizard
```
*The wizard will ask you a series of questions about your desired M3TAL configuration and update `/etc/m3tal/.env` accordingly.*

### Set a Single Environment Variable

Allows you to set a specific environment variable in your `.env` file.

**Usage:**
```bash
m3tal config set KEY VALUE
```

**Example:**
```bash
m3tal config set DASHBOARD_PORT 8083
```
*This command will update the `DASHBOARD_PORT` variable in `/etc/m3tal/.env` to `8083`.*

### Get a Single Environment Variable

Reads and displays the value of a specific environment variable from your `.env` file.

**Usage:**
```bash
m3tal config get KEY
```

**Example:**
```bash
m3tal config get LOG_LEVEL
```
*This command will output the current value of the `LOG_LEVEL` environment variable, for instance: `info`.*

### Scan All Environment Variables

Lists all environment variables across all configured M3TAL stacks.

**Usage:**
```bash
m3tal config scan
```

**Example:**
```bash
m3tal config scan
```
*This will display a comprehensive list of all environment variables recognized by M3TAL, including their current values.*

### List Current `.env` File Contents

Displays the current contents of the `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

**Example:**
```bash
m3tal config list
```
*This command will print the entire content of your `/etc/m3tal/.env` file to the console.*

---

## Dashboard Management

The following commands manage the M3TAL dashboard container.

### Update Dashboard Password

Updates the password for a dashboard user. If no username or password is provided, it will prompt interactively.

**Usage:**
```bash
m3tal dashpass [username] [password]
```

**Examples:**
```bash
# Interactive prompt for username and password
sudo m3tal dashpass

# Set password for user 'admin' to 'SuperSecurePassword123'
sudo m3tal dashpass admin SuperSecurePassword123
```
*These commands update the `users.json` file, which stores dashboard credentials.*

### Start and Update Dashboard

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Usage:**
```bash
m3tal dash up
```

**Example:**
```bash
sudo m3tal dash up
```
*This command ensures you have the most recent dashboard configuration and starts its Docker container.*

### Stop Dashboard Container

Stops the dashboard Docker container.

**Usage:**
```bash
m3tal dash down
```

**Example:**
```bash
sudo m3tal dash down
```
*This command gracefully stops the `m3tal-dashboard` container.*

### Restart Dashboard Container

Restarts the dashboard Docker container.

**Usage:**
```bash
m3tal dash restart
```

**Example:**
```bash
sudo m3tal dash restart
```
*This command will stop and then start the `m3tal-dashboard` container.*

### Stream Dashboard Logs

Streams the logs from the dashboard Docker container in real-time.

**Usage:**
```bash
m3tal dash logs
```

**Example:**
```bash
sudo m3tal dash logs
```
*This will display the live output of the dashboard container's logs, useful for debugging.*

### Show Dashboard Container Status

Displays the current status of the dashboard Docker container.

**Usage:**
```bash
m3tal dash status
```

**Example:**
```bash
sudo m3tal dash status
```
*This command will show whether the `m3tal-dashboard` container is running, exited, or in another state.*

---

## Stack Management

These commands manage all running M3TAL stacks based on Docker Compose files found in `/docker/`.

### Start All Stacks

Runs `docker compose up` across all `*-compose.yml` files located in the `/docker/` directory.

**Usage:**
```bash
m3tal up
```

**Example:**
```bash
sudo m3tal up
```
*This command starts all defined M3TAL services and user-added stacks.*

### Stop All Stacks

Runs `docker compose down` across all stacks managed by M3TAL.

**Usage:**
```bash
m3tal down
```

**Example:**
```bash
sudo m3tal down
```
*This command stops and removes containers for all M3TAL services and user-added stacks.*

### Stream Aggregated Logs

Streams aggregated logs from all running M3TAL stacks.

**Usage:**
```bash
m3tal logs
```

**Example:**
```bash
sudo m3tal logs
```
*This command collects and displays logs from all active Docker containers managed by M3TAL in real-time.*

---

## Systemd Service Management

The M3TAL API daemon is managed by systemd.

### Check API Service Status

Displays the current status of the `m3tal-api.service`.

**Usage:**
```bash
systemctl status m3tal-api
```

**Example:**
```bash
systemctl status m3tal-api
```
*This command will show if the API service is active, inactive, failed, and provide recent log entries.*

### View API Service Logs

Streams the logs from the `m3tal-api.service` in real-time.

**Usage:**
```bash
journalctl -u m3tal-api -f
```

**Example:**
```bash
journalctl -u m3tal-api -f
```
*This command is useful for monitoring the API daemon's activity and debugging issues.*

---

## Docker Compose Fallback

In situations where direct `m3tal` commands might not suffice, you can interact with Docker Compose directly. M3TAL uses Docker Compose V2 and orchestrates all `*-compose.yml` files found in `/docker/`.

**Usage:**
```bash
docker compose [options] [command]
```

**Examples:**
```bash
# Navigate to the M3TAL stack directory
cd /docker

# List all running services managed by Docker Compose in the current directory
docker compose ps

# Start a specific stack (e.g., ollama.yml)
docker compose -f ollama.yml up -d

# Stop a specific stack (e.g., ollama.yml)
docker compose -f ollama.yml down

# View logs for a specific service within a stack
docker compose -f ollama.yml logs ollama
```
*Direct Docker Compose commands provide granular control over individual compose files and services within the M3TAL environment.*

---
This document outlines the core functionality of the M3TAL CLI. For more in-depth information on specific commands or M3TAL's architecture, please refer to the official M3TAL documentation.