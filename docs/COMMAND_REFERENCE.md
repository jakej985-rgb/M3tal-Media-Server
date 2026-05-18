# M3TAL CLI Command Reference

This document provides a comprehensive reference for all commands available through the M3TAL CLI.

## Table of Contents

* [General Commands](#general-commands)
* [Configuration Management](#configuration-management)
* [Dashboard Management](#dashboard-management)
* [System Management](#system-management)
* [Docker and Compose Management](#docker-and-compose-management)
* [Systemd Service Management](#systemd-service-management)
* [Direct Docker Compose Fallback](#direct-docker-compose-fallback)

---

## General Commands

### `sudo m3tal`
Opens the interactive TUI Control Center, providing a numbered menu for various operations.

**Usage:**
```bash
sudo m3tal
```

---

## Configuration Management

### `m3tal init`
Generates the `/etc/m3tal/.env` configuration file from default values. This command should be used on the first install or when a fresh configuration is desired.

**Usage:**
```bash
sudo m3tal init
```

### `m3tal doctor`
Performs a pre-flight health check of your M3TAL system. It verifies Docker connectivity, the validity of your `.env` file, and checks for port availability.

**Usage:**
```bash
sudo m3tal doctor
```

### `m3tal config wizard`
Launches an interactive wizard to guide you through configuring your `/etc/m3tal/.env` file. This is the recommended method for detailed configuration.

**Usage:**
```bash
sudo m3tal config wizard
```

### `m3tal config set KEY VALUE`
Sets a single environment variable in the `/etc/m3tal/.env` file.

**Usage:**
```bash
sudo m3tal config set DOMAIN mydomain.com
```

### `m3tal config get KEY`
Reads and displays the value of a single environment variable from the `/etc/m3tal/.env` file.

**Usage:**
```bash
sudo m3tal config get DOMAIN
```

### `m3tal config scan`
Lists all environment variables known to M3TAL across all managed stacks, including their current values and defaults.

**Usage:**
```bash
sudo m3tal config scan
```

### `m3tal config list`
Displays the current contents of the `/etc/m3tal/.env` file.

**Usage:**
```bash
sudo m3tal config list
```

---

## Dashboard Management

### `m3tal dashpass [username] [password]`
Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt for them interactively.

**Usage (interactive):**
```bash
sudo m3tal dashpass
```
*Follow the prompts to enter the username and new password.*

**Usage (non-interactive):**
```bash
sudo m3tal dashpass admin new_secure_password_here
```

### `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the dashboard container. This command ensures you are running the most up-to-date version of the dashboard.

**Usage:**
```bash
sudo m3tal dash up
```

### `m3tal dash down`
Stops and removes the dashboard Docker container.

**Usage:**
```bash
sudo m3tal dash down
```

### `m3tal dash restart`
Restarts the dashboard Docker container.

**Usage:**
```bash
sudo m3tal dash restart
```

### `m3tal dash logs`
Streams the real-time logs from the dashboard Docker container. This is useful for debugging.

**Usage:**
```bash
sudo m3tal dash logs
```

### `m3tal dash status`
Shows the current status of the dashboard Docker container (e.g., running, exited).

**Usage:**
```bash
sudo m3tal dash status
```

---

## System Management

### `m3tal up`
Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This command starts all your deployed M3TAL stacks.

**Usage:**
```bash
sudo m3tal up
```

### `m3tal down`
Runs `docker compose down` across all managed Docker Compose stacks. This command stops and removes all containers, networks, and volumes associated with your M3TAL services.

**Usage:**
```bash
sudo m3tal down
```

### `m3tal logs`
Streams aggregated logs from all currently running M3TAL stacks. This provides a consolidated view of your system's output.

**Usage:**
```bash
sudo m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon is managed by systemd. You can interact with it using `systemctl` and `journalctl`.

### `systemctl status m3tal-api`
Displays the current status of the `m3tal-api.service`.

**Usage:**
```bash
sudo systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`
Streams the logs from the `m3tal-api.service` in real-time. Use `-f` to follow the log output.

**Usage:**
```bash
sudo journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Fallback

In situations where the `m3tal` CLI commands may not be sufficient, you can interact directly with Docker Compose. M3TAL orchestrates your stacks using Docker Compose V2.

**Note:** All direct Docker Compose commands should be executed from the `/docker/` directory, as this is where M3TAL expects to find your `*-compose.yml` files.

### Listing Docker Compose Files
M3TAL uses all `*-compose.yml` files in `/docker/`.

### Starting All Stacks
```bash
cd /docker
sudo docker compose up -d
```

### Stopping All Stacks
```bash
cd /docker
sudo docker compose down
```

### Streaming All Logs
```bash
cd /docker
sudo docker compose logs -f
```

### Starting a Specific Stack
If you have a file named `my-app-compose.yml` in `/docker/`:
```bash
cd /docker
sudo docker compose -f my-app-compose.yml up -d
```

### Stopping a Specific Stack
```bash
cd /docker
sudo docker compose -f my-app-compose.yml down
```

### Viewing a Specific Stack's Logs
```bash
cd /docker
sudo docker compose -f my-app-compose.yml logs -f
```
---