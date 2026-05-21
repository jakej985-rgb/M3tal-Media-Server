```markdown
# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for all available `m3tal` CLI commands.

## Installation

To install M3TAL, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

## Core Commands

### `sudo m3tal`

Opens the interactive TUI (Text User Interface) Control Center. This provides a menu-driven interface for managing M3TAL services.

**Usage:**
```bash
sudo m3tal
```

### `m3tal init`

Generates the `/etc/m3tal/.env` configuration file from default values. This command should be run upon the first installation of M3TAL.

**Usage:**
```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, the validity of your `.env` file, and checks for port availability.

**Usage:**
```bash
m3tal doctor
```

## Configuration Management

### `m3tal config wizard`

Launches an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name and `VALUE` with its desired value.

**Usage:**
```bash
m3tal config set DOMAIN mydomain.com
```

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from the `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name.

**Usage:**
```bash
m3tal config get DASHBOARD_PORT
```

### `m3tal config scan`

Lists all environment variables across all managed stacks, showing their current values as defined in your configuration.

**Usage:**
```bash
m3tal config scan
```

### `m3tal config list`

Displays the current contents of your `/etc/m3tal/.env` file.

**Usage:**
```bash
m3tal config list
```

## Dashboard Management

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt you interactively.

**Usage (interactive):**
```bash
m3tal dashpass
```

**Usage (with arguments):**
```bash
m3tal dashpass admin new_secure_password
```

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the dashboard container.

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

Streams the logs from the M3TAL dashboard container in real-time.

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

## Stack Management

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This starts all your deployed services.

**Usage:**
```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` across all stacks defined in the `/docker/` directory. This stops and removes all containers, networks, and volumes for the managed stacks.

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

The M3TAL API daemon (`m3tal-api.service`) is managed by `systemd`. You can interact with it using the `systemctl` and `journalctl` commands.

### `systemctl status m3tal-api`

Displays the current status of the `m3tal-api.service`.

**Usage:**
```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time. Use `Ctrl+C` to exit the stream.

**Usage:**
```bash
journalctl -u m3tal-api -f
```

## Direct Docker Compose Commands (Fallback)

In situations where the `m3tal` CLI might not be sufficient, you can directly use `docker compose` commands within the relevant directories. M3TAL manages services using Docker Compose V2.

### Running all stacks:
The `m3tal up` and `m3tal down` commands execute `docker compose` within the `/docker` directory. You can replicate this behavior manually:

**Start all stacks:**
```bash
cd /docker
docker compose up -d
```

**Stop all stacks:**
```bash
cd /docker
docker compose down
```

### Managing the dashboard specifically:
The `m3tal dash up/down/restart` commands operate on the dashboard. These typically involve pulling specific compose files and then running `docker compose`.

**Example: Starting the dashboard with local expose mode:**
```bash
cd /docker
# Ensure you have m3tal-compose.yml and m3tal-compose.local.yml
docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml up -d m3tal-dashboard
```

**Example: Starting the dashboard with traefik expose mode:**
```bash
cd /docker
# Ensure you have m3tal-compose.yml and m3tal-compose.traefik.yml
docker compose -f m3tal-compose.yml -f m3tal-compose.traefik.yml up -d m3tal-dashboard
```

**Viewing dashboard logs directly:**
```bash
docker logs -f m3tal-dashboard
```

**Checking dashboard container status:**
```bash
docker ps -a | grep m3tal-dashboard
```
```