# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for the M3TAL command-line interface (CLI), detailing each command and its usage.

---

## Interactive TUI Control Center

The `sudo m3tal` command launches the M3TAL TUI Control Center, offering an interactive, menu-driven interface for managing your M3TAL system.

**Usage:**
```bash
sudo m3tal
```

---

## M3TAL Configuration Commands

These commands manage the M3TAL environment configuration, primarily through the `/etc/m3tal/.env` file.

### `m3tal init`

Generates the `/etc/m3tal/.env` file with default values. This command should be used on the first installation of M3TAL.

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

Launches an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file.

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

Lists all environment variables across all managed Docker stacks, including their current values.

**Usage:**
```bash
sudo m3tal config scan
```

### `m3tal config list`

Lists the current contents of the `/etc/m3tal/.env` file.

**Usage:**
```bash
sudo m3tal config list
```

---

## M3TAL Dashboard Management

Commands specifically for managing the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt for them interactively.

**Usage (interactive):**
```bash
sudo m3tal dashpass
```

**Usage (with arguments):**
```bash
sudo m3tal dashpass dashboard_user new_secure_password
```

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the dashboard container. This command respects the `DASHBOARD_EXPOSE_MODE` setting in your `.env` file.

**Usage:**
```bash
sudo m3tal dash up
```

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash down
```

### `m3tal dash restart`

Restarts the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash restart
```

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container in real-time.

**Usage:**
```bash
sudo m3tal dash logs
```

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container.

**Usage:**
```bash
sudo m3tal dash status
```

---

## M3TAL Stack Management

These commands manage all Docker Compose stacks defined in your `/docker/` directory.

### `m3tal up`

Starts all Docker Compose services defined in `*-compose.yml` files located in the `/docker/` directory.

**Usage:**
```bash
sudo m3tal up
```

### `m3tal down`

Stops and removes all Docker Compose services defined in `*-compose.yml` files located in the `/docker/` directory.

**Usage:**
```bash
sudo m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL Docker stacks.

**Usage:**
```bash
sudo m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon is managed as a systemd service.

### `systemctl status m3tal-api`

Displays the current status of the `m3tal-api.service`.

**Usage:**
```bash
sudo systemctl status m3tal-api
```

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time using `journalctl`.

**Usage:**
```bash
sudo journalctl -u m3tal-api -f
```

---

## Docker Compose Fallback Commands

In scenarios where direct `m3tal` commands might not be sufficient, you can interact with Docker Compose directly. These commands assume you are in the `/docker/` directory or have it in your PATH.

### Running Docker Compose commands directly:

**Usage (example: starting all services):**
```bash
sudo docker compose up -d
```

**Usage (example: stopping all services):**
```bash
sudo docker compose down
```

**Usage (example: viewing logs for a specific service, e.g., `m3tal-dashboard`):**
```bash
sudo docker compose logs -f m3tal-dashboard
```

---

## APT Installation

To install or update M3TAL, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update package lists and install M3TAL
sudo apt update && sudo apt install -y m3tal
```