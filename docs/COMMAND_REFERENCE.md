# M3TAL CLI Command Reference

This document serves as a comprehensive cheat-sheet for all available M3TAL CLI commands.

## Systemd Service Management

The M3TAL API daemon is managed by systemd.

- **Check service status:**
  ```bash
  systemctl status m3tal-api
  ```

- **Restart the service:**
  ```bash
  systemctl restart m3tal-api
  ```

- **View service logs in real-time:**
  ```bash
  journalctl -u m3tal-api -f
  ```

## Docker Compose Fallback

While the `m3tal` CLI provides a unified interface, you can also use `docker compose` directly for more granular control over your stacks. All user-defined stacks are located in `/docker/`.

- **Run all docker compose services (equivalent to `m3tal up`):**
  ```bash
  docker compose -f /docker/*.yml up -d
  ```

- **Stop all docker compose services (equivalent to `m3tal down`):**
  ```bash
  docker compose -f /docker/*.yml down
  ```

- **View logs for a specific stack (e.g., a stack named `my-app`):**
  ```bash
  docker compose -f /docker/my-app-compose.yml logs -f
  ```

## M3TAL CLI Commands

### `sudo m3tal`

This command launches the interactive TUI Control Center, presenting a numbered menu of common operations.

**Usage Example:**
```bash
sudo m3tal
```

### `m3tal init`

Generates the default `/etc/m3tal/.env` configuration file. This command should be used during the initial installation of M3TAL.

**Usage Example:**
```bash
m3tal init
```

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, checks the validity of your `.env` file, and ensures necessary ports are available.

**Usage Example:**
```bash
m3tal doctor
```

### `m3tal config wizard`

Initiates an interactive wizard to guide you through configuring your `/etc/m3tal/.env` file. This is the recommended method for setting up your environment variables.

**Usage Example:**
```bash
m3tal config wizard
```

### `m3tal config set KEY VALUE`

Sets a single environment variable in your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config set DOMAIN traefik.example.com
```

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config get LOG_LEVEL
```

### `m3tal config scan`

Lists all environment variables across all managed Docker Compose stacks, including their current values.

**Usage Example:**
```bash
m3tal config scan
```

### `m3tal config list`

Displays the current contents of your `/etc/m3tal/.env` file.

**Usage Example:**
```bash
m3tal config list
```

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If no username and password are provided, the command will prompt you interactively.

**Usage Example (interactive):**
```bash
m3tal dashpass
```

**Usage Example (with arguments):**
```bash
m3tal dashpass admin new_secure_password_here
```

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the `m3tal-dashboard` container. This command ensures your dashboard is running with the most up-to-date settings.

**Usage Example:**
```bash
m3tal dash up
```

### `m3tal dash down`

Stops the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash down
```

### `m3tal dash restart`

Restarts the `m3tal-dashboard` container.

**Usage Example:**
```bash
m3tal dash restart
```

### `m3tal dash logs`

Streams the logs from the `m3tal-dashboard` container in real-time.

**Usage Example:**
```bash
m3tal dash logs
```

### `m3tal dash status`

Displays the current status of the `m3tal-dashboard` container (e.g., running, exited, paused).

**Usage Example:**
```bash
m3tal dash status
```

### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files located in the `/docker/` directory. This command brings up all your deployed M3TAL stacks.

**Usage Example:**
```bash
m3tal up
```

### `m3tal down`

Runs `docker compose down` across all stacks managed by M3TAL. This command stops and removes all containers, networks, and volumes associated with your deployed stacks.

**Usage Example:**
```bash
m3tal down
```

### `m3tal logs`

Streams aggregated logs from all currently running M3TAL Docker stacks. This provides a centralized view of your system's activity.

**Usage Example:**
```bash
m3tal logs
```