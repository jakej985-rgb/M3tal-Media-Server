# M3TAL CLI Command Reference

This document provides a comprehensive reference for all M3TAL command-line interface (CLI) commands. The M3TAL CLI is your primary interface for managing and interacting with the M3TAL ecosystem.

## Systemd Service Management

The M3TAL API daemon is managed by systemd. You can interact with the service using `systemctl` and view logs with `journalctl`.

*   **Check the status of the M3TAL API service:**
    ```bash
    systemctl status m3tal-api
    ```

*   **View real-time logs for the M3TAL API service:**
    ```bash
    journalctl -u m3tal-api -f
    ```

## Docker Compose Fallback Commands

While the M3TAL CLI provides a unified interface, you can directly use Docker Compose commands as a fallback if needed. All stack compose files are located in `/docker/`.

*   **Start all services defined in compose files within `/docker/`:**
    ```bash
    docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml -f /docker/routing-compose.yml up -d
    ```
    *Note: The actual files used will depend on your configuration and `DASHBOARD_EXPOSE_MODE`.*

*   **Stop all services defined in compose files within `/docker/`:**
    ```bash
    docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml -f /docker/routing-compose.yml down
    ```

*   **Stream aggregated logs from all running Docker containers managed by M3TAL:**
    ```bash
    docker compose -f /docker/m3tal-compose.yml -f /docker/m3tal-compose.local.yml -f /docker/routing-compose.yml logs -f
    ```

## M3TAL CLI Commands

### Interactive Control Center

*   **`sudo m3tal`**
    Opens the interactive TUI Control Center, presenting a numbered menu of common M3TAL operations.
    ```bash
    sudo m3tal
    ```

### Initialization and Configuration

*   **`m3tal init`**
    Generates the default `/etc/m3tal/.env` file. Use this on first installation to establish base configuration.
    ```bash
    m3tal init
    ```

*   **`m3tal config wizard`**
    Launches an interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.
    ```bash
    m3tal config wizard
    ```

*   **`m3tal config set KEY VALUE`**
    Sets a single environment variable in `/etc/m3tal/.env`.
    ```bash
    m3tal config set DOMAIN mydomain.com
    ```

*   **`m3tal config get KEY`**
    Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.
    ```bash
    m3tal config get DOMAIN
    ```

*   **`m3tal config scan`**
    Lists all environment variables across all active M3TAL stacks and their current values.
    ```bash
    m3tal config scan
    ```

*   **`m3tal config list`**
    Displays the entire contents of the current `/etc/m3tal/.env` file.
    ```bash
    m3tal config list
    ```

### Dashboard Management

*   **`m3tal dashpass [username] [password]`**
    Updates the password for a dashboard user. If username and password are not provided, the command will prompt interactively.
    ```bash
    m3tal dashpass admin new_secure_password
    ```
    Or, for interactive mode:
    ```bash
    m3tal dashpass
    ```

*   **`m3tal dash up`**
    Pulls the latest dashboard Docker Compose configuration from GitHub and starts the dashboard container.
    ```bash
    m3tal dash up
    ```

*   **`m3tal dash down`**
    Stops the dashboard Docker container.
    ```bash
    m3tal dash down
    ```

*   **`m3tal dash restart`**
    Restarts the dashboard Docker container.
    ```bash
    m3tal dash restart
    ```

*   **`m3tal dash logs`**
    Streams the logs from the dashboard Docker container in real-time.
    ```bash
    m3tal dash logs
    ```

*   **`m3tal dash status`**
    Shows the current status of the dashboard Docker container.
    ```bash
    m3tal dash status
    ```

### System-wide Stack Management

*   **`m3tal up`**
    Runs `docker compose up` for all `*-compose.yml` files found in the `/docker/` directory, starting all your M3TAL services.
    ```bash
    m3tal up
    ```

*   **`m3tal down`**
    Runs `docker compose down` for all stacks defined in the `/docker/` directory, stopping all your M3TAL services.
    ```bash
    m3tal down
    ```

*   **`m3tal logs`**
    Streams aggregated logs from all running M3TAL stacks. This command provides a consolidated view of the output from all managed containers.
    ```bash
    m3tal logs
    ```

## APT Installation

To install M3TAL, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```