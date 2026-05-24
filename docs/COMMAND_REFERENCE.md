# M3TAL CLI Command Reference

This document provides a comprehensive reference for all M3TAL CLI commands. It serves as a cheat-sheet for managing your M3TAL ecosystem.

## Table of Contents

1.  [Core Commands](#core-commands)
2.  [Configuration Management](#configuration-management)
3.  [Dashboard Management](#dashboard-management)
4.  [Stack Management](#stack-management)
5.  [Systemd Service Management](#systemd-service-management)
6.  [Direct Docker Compose Commands](#direct-docker-compose-commands)

---

## Core Commands

### `sudo m3tal`

Opens the interactive M3TAL TUI Control Center. This is your primary interface for navigating and managing the M3TAL ecosystem.

**Usage:**

```bash
sudo m3tal
```

---

### `m3tal init`

Generates the `/etc/m3tal/.env` file from default configurations. This command should be used on the initial installation of M3TAL to set up the basic environment variables.

**Usage:**

```bash
m3tal init
```

---

### `m3tal doctor`

Performs a pre-flight health check of your M3TAL installation. It verifies Docker connectivity, the validity of your `.env` file configuration, and checks for port availability to ensure a smooth operation.

**Usage:**

```bash
m3tal doctor
```

---

## Configuration Management

### `m3tal config wizard`

Launches an interactive wizard to configure the `/etc/m3tal/.env` file. This is the recommended method for setting up and modifying your M3TAL environment variables.

**Usage:**

```bash
m3tal config wizard
```

---

### `m3tal config set KEY VALUE`

Sets a single environment variable in the `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name and `VALUE` with its desired value.

**Usage:**

```bash
m3tal config set DOMAIN mydomain.local
```

---

### `m3tal config get KEY`

Reads and displays the value of a single environment variable from the `/etc/m3tal/.env` file. Replace `KEY` with the environment variable name you wish to query.

**Usage:**

```bash
m3tal config get DASHBOARD_PORT
```

---

### `m3tal config scan`

Lists all environment variables across all managed stacks within your M3TAL configuration. This provides a comprehensive overview of your system's settings.

**Usage:**

```bash
m3tal config scan
```

---

### `m3tal config list`

Displays the current contents of the `/etc/m3tal/.env` file. This is a quick way to review your active configuration.

**Usage:**

```bash
m3tal config list
```

---

## Dashboard Management

### `m3tal dashpass [username] [password]`

Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt you interactively for this information.

**Usage (interactive):**

```bash
m3tal dashpass
```

**Usage (with arguments):**

```bash
m3tal dashpass myuser newsecurepassword123
```

---

### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub and then starts the M3TAL dashboard container. This ensures your dashboard is running with the most up-to-date settings.

**Usage:**

```bash
m3tal dash up
```

---

### `m3tal dash down`

Stops and removes the M3TAL dashboard Docker container.

**Usage:**

```bash
m3tal dash down
```

---

### `m3tal dash restart`

Restarts the M3TAL dashboard Docker container. This is useful for applying configuration changes or recovering from unexpected issues.

**Usage:**

```bash
m3tal dash restart
```

---

### `m3tal dash logs`

Streams the logs from the M3TAL dashboard container in real-time. This is invaluable for debugging and monitoring the dashboard's activity.

**Usage:**

```bash
m3tal dash logs
```

---

### `m3tal dash status`

Shows the current status of the M3TAL dashboard Docker container (e.g., running, stopped, exited).

**Usage:**

```bash
m3tal dash status
```

---

## Stack Management

### `m3tal up`

Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory. This command brings all your defined M3TAL stacks online.

**Usage:**

```bash
m3tal up
```

---

### `m3tal down`

Runs `docker compose down` across all managed stacks. This command stops and removes containers, networks, and volumes defined in your Compose files.

**Usage:**

```bash
m3tal down
```

---

### `m3tal logs`

Streams aggregated logs from all running M3TAL stacks. This provides a centralized view of your system's logs for easier troubleshooting.

**Usage:**

```bash
m3tal logs
```

---

## Systemd Service Management

M3TAL's core API daemon runs as a systemd service. You can manage it using `systemctl` and view its logs with `journalctl`.

### `systemctl status m3tal-api`

Displays the current status of the `m3tal-api.service`. This will show if the service is active, loaded, and running.

**Usage:**

```bash
systemctl status m3tal-api
```

---

### `journalctl -u m3tal-api -f`

Streams the logs from the `m3tal-api.service` in real-time. The `-f` flag allows you to follow the logs as they are generated.

**Usage:**

```bash
journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Commands

As a fallback or for advanced users, you can interact with Docker Compose directly. M3TAL uses Docker Compose V2. The primary location for stack Compose files is `/docker/`.

**Note:** These commands are executed from within the `/docker/` directory or by specifying the Compose file path.

### `docker compose up -d`

Starts all services defined in `docker-compose.yml` (or other specified Compose files) in detached mode. This is equivalent to `m3tal up` if your stacks are structured accordingly.

**Usage (from `/docker/`):**

```bash
cd /docker/
docker compose up -d
```

**Usage (specifying files):**

```bash
docker compose -f /docker/my-stack-compose.yml -f /docker/another-stack-compose.yml up -d
```

---

### `docker compose down`

Stops and removes containers, networks, and volumes defined in the Compose file. This is equivalent to `m3tal down`.

**Usage (from `/docker/`):**

```bash
cd /docker/
docker compose down
```

---

### `docker compose logs -f [service_name]`

Streams logs for a specific service or all services if `[service_name]` is omitted. This is a more targeted way to view logs than `m3tal logs`.

**Usage (streaming all logs):**

```bash
cd /docker/
docker compose logs -f
```

**Usage (streaming a specific service's logs):**

```bash
cd /docker/
docker compose logs -f m3tal-dashboard
```

---

### `docker compose ps`

Lists the running containers for the services defined in the Compose file.

**Usage (from `/docker/`):**

```bash
cd /docker/
docker compose ps
```

---

### `docker compose build`

Builds, (re)creates, starts, and attaches to containers for a service.

**Usage (from `/docker/`):**

```bash
cd /docker/
docker compose build
```

---

### `docker compose pull`

Pulls service images.

**Usage (from `/docker/`):**

```bash
cd /docker/
docker compose pull
```

---