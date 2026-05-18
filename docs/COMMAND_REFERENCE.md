# M3TAL CLI Command Reference

This document provides a comprehensive cheat sheet for the M3TAL Command Line Interface (CLI).

## Interactive Control Center

The primary entry point for most M3TAL operations is the interactive TUI (Text-based User Interface).

*   **`sudo m3tal`**
    *   Opens the interactive TUI Control Center, presenting a numbered menu of common operations.
    *   **Example:**
        ```bash
        sudo m3tal
        ```

## Initialization and Configuration

Manage M3TAL's core configuration and environment variables.

*   **`m3tal init`**
    *   Generates the default `/etc/m3tal/.env` configuration file. Recommended for first-time installations.
    *   **Example:**
        ```bash
        m3tal init
        ```

*   **`m3tal config wizard`**
    *   Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.
    *   **Example:**
        ```bash
        m3tal config wizard
        ```

*   **`m3tal config set KEY VALUE`**
    *   Sets a single environment variable in the `/etc/m3tal/.env` file.
    *   **Example:**
        ```bash
        m3tal config set DOMAIN mym3tal.local
        ```

*   **`m3tal config get KEY`**
    *   Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.
    *   **Example:**
        ```bash
        m3tal config get DOMAIN
        ```
        ```
        localhost
        ```

*   **`m3tal config scan`**
    *   Lists all environment variables across all configured stacks, showing their current values.
    *   **Example:**
        ```bash
        m3tal config scan
        ```
        ```
        # ... (output showing all env vars and their values) ...
        ```

*   **`m3tal config list`**
    *   Displays the current contents of the `/etc/m3tal/.env` file.
    *   **Example:**
        ```bash
        m3tal config list
        ```
        ```
        # This is the .env file for M3TAL
        # ... (output of the .env file) ...
        ```

## Dashboard Management

Commands specifically for managing the M3TAL dashboard container.

*   **`m3tal dashpass [username] [password]`**
    *   Updates the password for a dashboard user. If username and password are not provided, it will prompt interactively.
    *   **Example (interactive):**
        ```bash
        m3tal dashpass
        ```
    *   **Example (with arguments):**
        ```bash
        m3tal dashpass admin new_secret_password
        ```

*   **`m3tal dash up`**
    *   Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.
    *   **Example:**
        ```bash
        m3tal dash up
        ```

*   **`m3tal dash down`**
    *   Stops and removes the dashboard container.
    *   **Example:**
        ```bash
        m3tal dash down
        ```

*   **`m3tal dash restart`**
    *   Restarts the dashboard container.
    *   **Example:**
        ```bash
        m3tal dash restart
        ```

*   **`m3tal dash logs`**
    *   Streams the logs from the dashboard container in real-time.
    *   **Example:**
        ```bash
        m3tal dash logs
        ```

*   **`m3tal dash status`**
    *   Shows the current status of the dashboard container.
    *   **Example:**
        ```bash
        m3tal dash status
        ```
        ```
        m3tal-dashboard  ghcr.io/jakej985-rgb/m3tal-godash:debug  /usr/bin/python3 server.py  Up 2 hours
        ```

## Stack Management

Commands for managing all deployed M3TAL stacks.

*   **`m3tal up`**
    *   Runs `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This starts all your deployed services.
    *   **Example:**
        ```bash
        m3tal up
        ```

*   **`m3tal down`**
    *   Runs `docker compose down` across all stacks, stopping and removing containers defined in the compose files.
    *   **Example:**
        ```bash
        m3tal down
        ```

*   **`m3tal logs`**
    *   Streams aggregated logs from all running M3TAL stacks.
    *   **Example:**
        ```bash
        m3tal logs
        ```

## Systemd Service Management

The M3TAL API daemon is managed by systemd.

*   **`systemctl status m3tal-api`**
    *   Displays the current status of the `m3tal-api` systemd service.
    *   **Example:**
        ```bash
        sudo systemctl status m3tal-api
        ```
        ```
        ● m3tal-api.service - M3TAL API Daemon
             Loaded: loaded (/etc/systemd/system/m3tal-api.service; enabled; vendor preset: enabled)
             Active: active (running) since Mon 2023-10-27 10:00:00 UTC; 1 day ago
               Docs: https://jakej985-rgb.github.io/m3tal-core/docs/
           Main PID: 12345 (m3tal-api)
              Tasks: 10
              Memory: 25.6M
                 CPU: 1min 30s
              CGroup: /system.slice/m3tal-api.service
                      └─12345 /usr/bin/m3tal-api
        ```

*   **`journalctl -u m3tal-api -f`**
    *   Streams the logs from the `m3tal-api` service in real-time. Use `Ctrl+C` to exit.
    *   **Example:**
        ```bash
        sudo journalctl -u m3tal-api -f
        ```
        ```
        Oct 28 10:00:00 yourhostname m3tal-api[12345]: INFO: Processing request from dashboard.
        Oct 28 10:00:01 yourhostname m3tal-api[12345]: INFO: Stack 'my-app' started successfully.
        ```

## Docker Compose Fallback Commands

In scenarios where direct Docker Compose interaction is needed, you can use the standard Docker Compose commands. M3TAL manages compose files located in `/docker/`.

*   **`docker compose up -d`**
    *   Starts all services defined in `*-compose.yml` files in the current directory (or specified with `-f`).
    *   **Example (from within `/docker/`):**
        ```bash
        cd /docker/
        docker compose up -d
        ```

*   **`docker compose down`**
    *   Stops and removes containers, networks, and volumes defined in compose files.
    *   **Example (from within `/docker/`):**
        ```bash
        cd /docker/
        docker compose down
        ```

*   **`docker compose ps`**
    *   Lists the containers and their status for the services defined in the compose files.
    *   **Example (from within `/docker/`):**
        ```bash
        cd /docker/
        docker compose ps
        ```
        ```
        NAME                 IMAGE                           COMMAND                  SERVICE             CREATED             STATUS              PORTS
        m3tal-dashboard      ghcr.io/jakej985-rgb/m3tal-godash:debug   "python3 server.py"      m3tal-dashboard     2 hours ago         Up 2 hours
        traefik              traefik:latest                  "/entrypoint.sh --..."   traefik             2 hours ago         Up 2 hours
        ```

*   **`docker compose logs [service_name]`**
    *   Displays logs from a specific service.
    *   **Example:**
        ```bash
        cd /docker/
        docker compose logs m3tal-dashboard
        ```

## APT Installation

To install or update M3TAL, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```