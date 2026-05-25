# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for the M3TAL command-line interface.

## Introduction

The M3TAL CLI (`m3tal`) is your primary tool for managing the M3TAL ecosystem. It provides a unified interface for initialization, configuration, service management, and more. All commands require root privileges via `sudo` unless otherwise specified.

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

This command launches the interactive M3TAL TUI (Text User Interface) Control Center. It presents a numbered menu of common operations, simplifying management for users.

**Usage Example:**

```bash
sudo m3tal
```

### `m3tal init`

Initializes the M3TAL configuration by generating a default `/etc/m3tal/.env` file. This should be run on the first installation or when you need to reset to default configuration values.

**Usage Example:**

```bash
sudo m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. This includes verifying Docker connectivity, the validity of your `.env` configuration, and checking for port availability conflicts.

**Usage Example:**

```bash
sudo m3tal doctor
```

### `m3tal config wizard`

Launches an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file. This is the recommended way to modify your M3TAL settings.

**Usage Example:**

```bash
sudo m3tal config wizard
```

### `m3tal config set KEY VALUE`

Allows you to set a single environment variable directly in your `/etc/m3tal/.env` file.

**Usage Example:**

```bash
sudo m3tal config set DOMAIN mym3tal.local
```

### `m3tal config get KEY`

Reads and displays the value of a specific environment variable from your `/etc/m3tal/.env` file.

**Usage Example:**

```bash
sudo m3tal config get DOMAIN
```

### `m3tal config scan`

Lists all environment variables across all active M3TAL stacks, providing a comprehensive overview of your configuration.

**Usage Example:**

```bash
sudo m3tal config scan
```

### `m3tal config list`

Displays the current contents of your `/etc/m3tal/.env` file.

**Usage Example:**

```bash
sudo m3tal config list
```

## Dashboard Management

The `m3tal dash` subcommand group is used to manage the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If the username and password are not provided as arguments, the command will prompt you interactively.

**Usage Examples:**

```bash
# Interactive password update for user 'admin'
sudo m3tal dashpass admin

# Set password for user 'admin' to 'secure_password'
sudo m3tal dashpass admin secure_password
```

### `m3tal dash up`

Pulls the latest dashboard Docker compose configuration from GitHub and then starts the dashboard container.

**Usage Example:**

```bash
sudo m3tal dash up
```

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage Example:**

```bash:**
```bash
sudo m3tal dash down
```

### `m3tal dash restart`

Restarts the M3TAL dashboard container.

**Usage Example:**

```bash
sudo m3tal dash restart
```

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container in real-time.

**Usage Example:**

```bash
sudo m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container (e.g., running, exited).

**Usage Example:**

```bash
sudo m3tal dash status
```

## Stack Management

These commands manage all your deployed M3TAL stacks using Docker Compose.

### `sudo m3tal up`

Starts all defined M3TAL services by running `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory.

**Usage Example:**

```bash
sudo m3tal up
```

### `sudo m3tal down`

Stops all defined M3TAL services by running `docker compose down` across all stacks managed by M3TAL.

**Usage Example:**

```bash
sudo m3tal down
```

### `sudo m3tal logs`

Streams aggregated logs from all currently running M3TAL stacks. This provides a unified log view for troubleshooting.

**Usage Example:**

```bash
sudo m3tal logs
```

## Systemd Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl` and view its logs with `journalctl`.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api` systemd service.

**Usage Example:**

```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api` service in real-time. Use `-f` to follow the logs.

**Usage Example:**

```bash
journalctl -u m3tal-api -f
```

## Docker Compose Fallback

While the `m3tal` CLI abstracts Docker Compose operations, you can always interact with Docker Compose directly if needed.

### Direct Docker Compose Usage

All M3TAL stacks are managed via Docker Compose files located in `/docker/`. You can use standard `docker compose` commands in this directory.

**Key Compose Files:**

- `/docker/m3tal-compose.yml`: Base configuration for the dashboard.
- `/docker/m3tal-compose.local.yml`: Override for local dashboard exposure.
- `/docker/m3tal-compose.traefik.yml`: Override for Traefik-based dashboard exposure.
- `/docker/routing-compose.yml`: Configuration for Traefik and Cloudflared.
- Other `*-compose.yml` files in `/docker/` define your custom stacks.

**Usage Examples:**

```bash
# Navigate to the docker directory
cd /docker

# Start all services defined in compose files
sudo docker compose up -d

# Stop all services
sudo docker compose down

# View logs for a specific service (e.g., m3tal-dashboard)
sudo docker compose logs -f m3tal-dashboard

# Rebuild and start a specific service
sudo docker compose up -d --build my-service-name
```