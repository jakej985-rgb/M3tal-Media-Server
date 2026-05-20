# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for the M3TAL command-line interface (CLI), covering all available commands and their usage.

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

---

## Core M3TAL Commands

### Interactive TUI Control Center

This command opens the M3TAL TUI Control Center, providing a menu-driven interface for common operations.

```bash
sudo m3tal
```

---

### `m3tal init`

Generates the default `/etc/m3tal/.env` configuration file. This is typically run on the first install.

**Usage:**

```bash
sudo m3tal init
```

---

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, the validity of your `.env` file, and checks for port availability.

**Usage:**

```bash
sudo m3tal doctor
```

---

### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.

**Usage:**

```bash
sudo m3tal config wizard
```

---

### `m3tal config set KEY VALUE`

Sets a single environment variable in your `/etc/m3tal/.env` file.

**Usage Example:**

```bash
sudo m3tal config set DOMAIN mym3tal.local
```

---

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from your `/etc/m3tal/.env` file.

**Usage Example:**

```bash
sudo m3tal config get DOMAIN
```

---

### `m3tal config scan`

Lists all environment variables across all managed Docker stacks, including those defined in your `.env` file and defaults.

**Usage:**

```bash
sudo m3tal config scan
```

---

### `m3tal config list`

Lists the current contents of your `/etc/m3tal/.env` file.

**Usage:**

```bash
sudo m3tal config list
```

---

## Dashboard Management Commands

The following commands manage the M3TAL Dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If username and password are not provided, the command will become interactive.

**Usage Example (Interactive):**

```bash
sudo m3tal dashpass
```

**Usage Example (with arguments):**

```bash
sudo m3tal dashpass admin new_secure_password
```

---

### `m3tal dash up`

Pulls the latest `m3tal-compose.yml` and its associated override files from GitHub, then starts the dashboard container. It respects the `DASHBOARD_EXPOSE_MODE` setting in your `.env` file.

**Usage:**

```bash
sudo m3tal dash up
```

---

### `m3tal dash down`

Stops and removes the dashboard container.

**Usage:**

```bash
sudo m3tal dash down
```

---

### `m3tal dash restart`

Restarts the dashboard container.

**Usage:**

```bash
sudo m3tal dash restart
```

---

### `m3tal dash logs`

Streams the logs from the dashboard container in real-time.

**Usage:**

```bash
sudo m3tal dash logs
```

---

### `m3tal dash status`

Shows the current status of the dashboard container.

**Usage:**

```bash
sudo m3tal dash status
```

---

## Stack Management Commands

These commands manage your Docker stacks, including the dashboard and any user-defined services.

### `m3tal up`

Starts all Docker services defined in `*-compose.yml` files located in the `/docker/` directory. This includes Traefik, the dashboard (if configured), and any other services you have added.

**Usage:**

```bash
sudo m3tal up
```

---

### `m3tal down`

Stops and removes all Docker containers, networks, and volumes defined by the `*-compose.yml` files in the `/docker/` directory.

**Usage:**

```bash
sudo m3tal down
```

---

### `m3tal logs`

Streams aggregated logs from all running M3TAL-managed Docker stacks in real-time.

**Usage:**

```bash
sudo m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl` and view its logs with `journalctl`.

### Managing the `m3tal-api` service

**Check the status of the M3TAL API service:**

```bash
systemctl status m3tal-api
```

**Restart the M3TAL API service:**

```bash
sudo systemctl restart m3tal-api
```

**View real-time logs of the M3TAL API service:**

```bash
journalctl -u m3tal-api -f
```

---

## Docker Compose Fallback Commands

While the M3TAL CLI provides a convenient abstraction, you can always interact with Docker Compose directly. M3TAL manages compose files located in `/docker/`.

**Start all services in `/docker/`:**

```bash
docker compose --profile default --profile m3tal -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml up -d
```

**Stop all services in `/docker/`:**

```bash
docker compose --profile default --profile m3tal -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml down
```

**Stream logs from all services in `/docker/`:**

```bash
docker compose --profile default --profile m3tal -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml logs -f
```

**Manage only the dashboard using its specific compose files:**

```bash
# Start dashboard with local mode
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml up -d m3tal-dashboard

# Stop dashboard
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml down m3tal-dashboard

# Stream dashboard logs
docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml logs -f m3tal-dashboard
```

---

**Note:** When using direct `docker compose` commands, ensure you include the relevant compose files (e.g., `m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) based on your `DASHBOARD_EXPOSE_MODE` and any other stack configurations. The `--profile` flags might also be necessary depending on how services are organized. The M3TAL CLI handles these complexities for you.