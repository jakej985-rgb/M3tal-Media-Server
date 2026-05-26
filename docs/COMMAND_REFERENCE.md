# M3TAL CLI Command Reference

This document provides a comprehensive reference for all M3TAL CLI commands, serving as a cheat sheet for managing your M3TAL ecosystem.

---

## Core M3TAL Commands

These commands form the primary interface for interacting with and configuring your M3TAL installation.

### `sudo m3tal`

Opens the interactive TUI Control Center. Navigate the numbered menu to perform various operations.

**Usage Example:**
```bash
sudo m3tal
```

### `m3tal init`

Generates the `/etc/m3tal/.env` file from default values. This command should be used upon first installation to establish the base configuration.

**Usage Example:**
```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL environment. It verifies Docker connectivity, the validity of your `.env` file, and checks for port availability.

**Usage Example:**
```bash
m3tal doctor
```

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file. This is the recommended method for comprehensive configuration.

**Usage Example:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name and `VALUE` with its desired setting.

**Usage Example:**
```bash
m3tal config set DASHBOARD_PORT 8083
```

### `m3tal config get KEY`

Retrieves the value of a single environment variable from the `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name.

**Usage Example:**
```bash
m3tal config get LOG_LEVEL
```

### `m3tal config scan`

Lists all environment variables currently recognized across all M3TAL stacks, including their current values and default values if applicable.

**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config list
```

---

## Dashboard Management Commands

These commands are specifically for managing the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt you interactively.

**Usage Example (interactive):**
```bash
m3tal dashpass
```

**Usage Example (with arguments):**
```bash
m3tal dashpass admin new_secure_password
```

### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Usage Example:**
```bash
m3tal dash up
```

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage Example:**
```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the M3TAL dashboard container.

**Usage Example:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container in real-time.

**Usage Example:**
```bash
m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container.

**Usage Example:**
```bash
m3tal dash status
```

---

## Stack Management Commands

These commands manage the Docker Compose stacks within your M3TAL environment.

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This starts all your defined services.

**Usage Example:**
```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` across all defined stacks, stopping and removing containers, networks, and volumes associated with them.

**Usage Example:**
```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL stacks.

**Usage Example:**
```bash
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service. These commands allow you to manage it.

### `systemctl status m3tal-api`

Checks the current status of the `m3tal-api.service`.

**Usage Example:**
```bash
systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time. Use `Ctrl+C` to exit.

**Usage Example:**
```bash
journalctl -u m3tal-api -f
```

---

## Docker Compose Fallback Commands

In scenarios where direct Docker Compose interaction is needed, you can use the following commands. Note that M3TAL manages Compose V2.

### `docker compose up -f /docker/your-stack.yml`

Starts a specific stack defined by a compose file. Replace `/docker/your-stack.yml` with the actual path to your compose file.

**Usage Example:**
```bash
docker compose up -f /docker/ollama-compose.yml
```

### `docker compose down -f /docker/your-stack.yml`

Stops a specific stack.

**Usage Example:**
```bash
docker compose down -f /docker/ollama-compose.yml
```

### `docker compose ps`

Lists all containers managed by Docker Compose.

**Usage Example:**
```bash
docker compose ps
```

### `docker compose logs [container_name]`

Streams logs for a specific container.

**Usage Example:**
```bash
docker compose logs m3tal-dashboard
```

---

## APT Installation

To install or update M3TAL via APT, use the following steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-apt-key/public.key | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-apt-key stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install or update
sudo apt update && sudo apt install -y m3tal
```