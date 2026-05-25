```markdown
# M3TAL CLI Command Reference

This document provides a comprehensive cheat sheet for the M3TAL CLI, covering all available commands and their usage.

## Table of Contents

- [Getting Started](#getting-started)
- [Interactive Control Center](#interactive-control-center)
- [Configuration Management](#configuration-management)
  - [Initializing Configuration](#initializing-configuration)
  - [Configuration Wizard](#configuration-wizard)
  - [Setting Environment Variables](#setting-environment-variables)
  - [Getting Environment Variables](#getting-environment-variables)
  - [Scanning All Environment Variables](#scanning-all-environment-variables)
  - [Listing Current .env File](#listing-current-env-file)
- [Dashboard Management](#dashboard-management)
  - [Updating Dashboard Password](#updating-dashboard-password)
  - [Starting/Updating Dashboard](#startingupdating-dashboard)
  - [Stopping Dashboard](#stopping-dashboard)
  - [Restarting Dashboard](#restarting-dashboard)
  - [Viewing Dashboard Logs](#viewing-dashboard-logs)
  - [Checking Dashboard Status](#checking-dashboard-status)
- [Stack Management](#stack-management)
  - [Bringing Up All Stacks](#bringing-up-all-stacks)
  - [Bringing Down All Stacks](#bringing-down-all-stacks)
  - [Aggregated Stack Logs](#aggregated-stack-logs)
- [Systemd Service Management](#systemd-service-management)
- [Direct Docker Compose Fallback](#direct-docker-compose-fallback)

## Getting Started

To install M3TAL, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Interactive Control Center

The `sudo m3tal` command launches the TUI Control Center, providing a menu-driven interface for managing M3TAL.

**Usage:**
```bash
sudo m3tal
```

This will present a numbered menu for various operations.

## Configuration Management

M3TAL's configuration is primarily managed through the `/etc/m3tal/.env` file.

### Initializing Configuration

Use this command on a fresh installation to generate the default `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal init
```

### Configuration Wizard

An interactive wizard to guide you through setting up your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config wizard
```

### Setting Environment Variables

Set a single environment variable in your `.env` file.

**Usage:**
```bash
m3tal config set KEY VALUE
```

**Example:**
```bash
m3tal config set DASHBOARD_PORT 8083
```

### Getting Environment Variables

Read the value of a single environment variable from your `.env` file.

**Usage:**
```bash
m3tal config get KEY
```

**Example:**
```bash
m3tal config get DOMAIN
```

### Scanning All Environment Variables

List all environment variables across all managed stacks, including their current values and defaults.

**Usage:**
```bash
m3tal config scan
```

### Listing Current .env File

Display the current contents of your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

## Dashboard Management

Manage the M3TAL dashboard container.

### Updating Dashboard Password

Update the username and password for accessing the M3TAL dashboard. If arguments are omitted, you will be prompted interactively.

**Usage (interactive):**
```bash
m3tal dashpass
```

**Usage (with arguments):**
```bash
m3tal dashpass myuser mysecurepassword
```

### Starting/Updating Dashboard

Pulls the latest dashboard compose configuration from GitHub and starts the dashboard container. This command also handles applying any necessary overrides based on `DASHBOARD_EXPOSE_MODE`.

**Usage:**
```bash
m3tal dash up
```

### Stopping Dashboard

Stops the M3TAL dashboard container.

**Usage:**
```bash
m3tal dash down
```

### Restarting Dashboard

Restarts the M3TAL dashboard container.

**Usage:**
```bash
m3tal dash restart
```

### Viewing Dashboard Logs

Streams the logs from the M3TAL dashboard container in real-time.

**Usage:**
```bash
m3tal dash logs
```

### Checking Dashboard Status

Shows the current status of the M3TAL dashboard container.

**Usage:**
```bash
m3tal dash status
```

## Stack Management

Manage all your deployed M3TAL stacks.

### Bringing Up All Stacks

Starts all services defined in `*-compose.yml` files within the `/docker/` directory using `docker compose up`.

**Usage:**
```bash
m3tal up
```

### Bringing Down All Stacks

Stops and removes all services defined in `*-compose.yml` files within the `/docker/` directory using `docker compose down`.

**Usage:**
```bash
m3tal down
```

### Aggregated Stack Logs

Streams aggregated logs from all currently running M3TAL stacks.

**Usage:**
```bash
m3tal logs
```

## Systemd Service Management

The M3TAL API daemon runs as a systemd service.

**Check the status of the M3TAL API service:**
```bash
systemctl status m3tal-api
```

**Follow the logs of the M3TAL API service in real-time:**
```bash
journalctl -u m3tal-api -f
```

**Restart the M3TAL API service:**
```bash
sudo systemctl restart m3tal-api
```

## Direct Docker Compose Fallback

In cases where direct `docker compose` commands are needed, remember that M3TAL uses Docker Compose V2 and operates on files within `/docker/`.

**To bring up a specific stack (e.g., named `my-stack-compose.yml`):**
```bash
docker compose -f /docker/my-stack-compose.yml up -d
```

**To bring down a specific stack:**
```bash
docker compose -f /docker/my-stack-compose.yml down
```

**To view logs for a specific service within a stack:**
```bash
docker compose -f /docker/my-stack-compose.yml logs -f my-service-name
```
```