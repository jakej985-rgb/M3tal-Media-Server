# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for the M3TAL command-line interface (CLI). It details each available command and its subcommands with practical usage examples.

---

## Core M3TAL Commands

### `sudo m3tal`

Launches the interactive M3TAL TUI (Text-based User Interface) Control Center. This provides a menu-driven approach to manage various aspects of your M3TAL system.

**Usage Example:**

```bash
sudo m3tal
```

This command will present a numbered menu, allowing you to navigate and select actions like starting/stopping services, managing configurations, and more.

### `m3tal init`

Initializes the M3TAL configuration by generating the `/etc/m3tal/.env` file from default values. This command should be run on the first installation of M3TAL.

**Usage Example:**

```bash
m3tal init
```

This will create or overwrite `/etc/m3tal/.env` with the standard default settings, preparing your system for further configuration.

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, the validity of your `.env` file, and checks for port availability to ensure your system is ready for operation.

**Usage Example:**

```bash
m3tal doctor
```

This command will output a report on the health of your M3TAL environment, highlighting any potential issues.

### `m3tal config wizard`

Launches an interactive wizard to guide you through the configuration of your `/etc/m3tal/.env` file. This is the recommended method for setting up your environment variables.

**Usage Example:**

```bash
m3tal config wizard
```

The wizard will prompt you for values for various configuration parameters, explaining their purpose and saving them to `/etc/m3tal/.env`.

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file. This command allows for direct modification of individual configuration parameters.

**Usage Example:**

```bash
m3tal config set DOMAIN mym3tal.local
```

This command will add or update the `DOMAIN` variable in `/etc/m3tal/.env` to `mym3tal.local`.

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from the `/etc/m3tal/.env` file.

**Usage Example:**

```bash
m3tal config get DASHBOARD_PORT
```

This command will output the current value of the `DASHBOARD_PORT` environment variable.

### `m3tal config scan`

Lists all environment variables across all active M3TAL stacks. This provides an overview of all configurable settings and their current values.

**Usage Example:**

```bash
m3tal config scan
```

This command will display a comprehensive list of environment variables used by M3TAL and its components.

### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file. This is useful for quickly reviewing your active configuration.

**Usage Example:**

```bash
m3tal config list
```

This command will print the entire content of your `/etc/m3tal/.env` file to the console.

---

## Dashboard Management Commands

The following commands specifically manage the M3TAL dashboard container.

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If no username or password is provided, the command will prompt you interactively.

**Usage Examples:**

```bash
m3tal dashpass admin new_secure_password
```

or interactively:

```bash
m3tal dashpass
```

This will prompt for the username and then the new password.

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the dashboard container.

**Usage Example:**

```bash
m3tal dash up
```

This command ensures you have the most recent dashboard definition and starts the service.

### `m3tal dash down`

Stops the M3TAL dashboard container.

**Usage Example:**

```bash
m3tal dash down
```

This command will gracefully shut down the dashboard container.

### `m3tal dash restart`

Restarts the M3TAL dashboard container.

**Usage Example:**

```bash
m3tal dash restart
```

This command stops and then starts the dashboard container, applying any recent configuration changes.

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container in real-time.

**Usage Example:**

```bash
m3tal dash logs
```

This is useful for troubleshooting and monitoring the dashboard's activity. Press `Ctrl+C` to stop streaming.

### `m3tal dash status`

Shows the current status of the M3TAL dashboard container (e.g., running, stopped, exited).

**Usage Example:**

```bash
m3tal dash status
```

This command provides a quick overview of the dashboard container's operational state.

---

## Global M3TAL Stack Management Commands

These commands manage all M3TAL services defined in `*-compose.yml` files.

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This command starts all M3TAL services and any user-added stacks.

**Usage Example:**

```bash
m3tal up
```

This will bring up all defined services, including the API, dashboard, Traefik, and any other stacks you have placed in `/docker/`.

### `m3tal down`

Runs `docker compose down` across all M3TAL stacks. This command stops and removes all containers, networks, and volumes defined in your compose files.

**Usage Example:**

```bash
m3tal down
```

This command will shut down all running M3TAL services.

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks. This command provides a unified view of logs from all your M3TAL services.

**Usage Example:**

```bash
m3tal logs
```

This command continuously displays log output from all active containers. Press `Ctrl+C` to stop streaming.

---

## Systemd Service Management

The M3TAL API daemon is managed by `systemd`.

### `systemctl status m3tal-api`

Displays the current status of the `m3tal-api.service` systemd unit.

**Usage Example:**

```bash
systemctl status m3tal-api
```

This will show if the service is active, loaded, enabled, and provide recent log snippets.

### `journalctl -u m3tal-api -f`

Streams the logs for the `m3tal-api.service` in real-time using `journalctl`.

**Usage Example:**

```bash
journalctl -u m3tal-api -f
```

This command is invaluable for real-time debugging and monitoring of the M3TAL API daemon. Press `Ctrl+C` to stop streaming.

---

## Docker Compose Fallback Commands

In situations where direct control is needed, you can use Docker Compose commands on the M3TAL compose files.

**Note:** It is generally recommended to use the `m3tal` CLI commands for consistency and to ensure proper management by the M3TAL daemon.

### `docker compose up -d` (in `/docker/`)

Starts all services defined in `*-compose.yml` files within the `/docker/` directory in detached mode.

**Usage Example:**

```bash
cd /docker/
docker compose up -d
```

### `docker compose down` (in `/docker/`)

Stops and removes all services defined in `*-compose.yml` files within the `/docker/` directory.

**Usage Example:**

```bash
cd /docker/
docker compose down
```

### `docker compose logs <container_name>` (in `/docker/`)

Fetches logs for a specific container managed by Docker Compose.

**Usage Example:**

```bash
cd /docker/
docker compose logs m3tal-dashboard
```

---

## APT Installation

To install or update M3TAL, use the following APT commands:

**Usage Example:**

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```