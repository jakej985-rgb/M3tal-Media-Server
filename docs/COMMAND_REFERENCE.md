# M3TAL CLI Command Reference

This document provides a comprehensive cheat sheet for all available `m3tal` CLI commands.

## Table of Contents

* [M3TAL CLI Commands](#m3tal-cli-commands)
    * [Interactive TUI](#interactive-tui)
    * [Initialization and Configuration](#initialization-and-configuration)
    * [Dashboard Management](#dashboard-management)
    * [Stack Management](#stack-management)
    * [Systemd Service Management](#systemd-service-management)
* [Direct Docker Compose Commands (Fallback)](#direct-docker-compose-commands-fallback)
* [APT Installation](#apt-installation)

---

## M3TAL CLI Commands

### Interactive TUI

The `sudo m3tal` command launches the M3TAL TUI Control Center, providing a menu-driven interface for common operations.

```bash
sudo m3tal
```

### Initialization and Configuration

These commands manage M3TAL's configuration, primarily through the `/etc/m3tal/.env` file.

*   **`m3tal init`**: Generates `/etc/m3tal/.env` from default values. Use this command upon initial installation.

    ```bash
    m3tal init
    ```

*   **`m3tal config wizard`**: Launches an interactive wizard to guide you through configuring `/etc/m3tal/.env`.

    ```bash
    m3tal config wizard
    ```

*   **`m3tal config set KEY VALUE`**: Sets a single environment variable in `/etc/m3tal/.env`.

    ```bash
    m3tal config set DASHBOARD_PORT 8083
    ```

*   **`m3tal config get KEY`**: Reads and displays the value of a single environment variable from `/etc/m3tal/.env`.

    ```bash
    m3tal config get DOMAIN
    ```

*   **`m3tal config scan`**: Lists all environment variables across all managed stacks, including their current values.

    ```bash
    m3tal config scan
    ```

*   **`m3tal config list`**: Displays the current contents of the `/etc/m3tal/.env` file.

    ```bash
    m3tal config list
    ```

### Dashboard Management

Commands specifically for managing the M3TAL dashboard container.

*   **`m3tal dashpass [username] [password]`**: Updates the password for a dashboard user. If `username` and `password` are omitted, the command will prompt for them interactively.

    ```bash
    m3tal dashpass admin your_new_secure_password
    ```
    or for interactive input:
    ```bash
    m3tal dashpass
    ```

*   **`m3tal dash up`**: Pulls the latest dashboard Compose configuration from GitHub and then starts the dashboard container.

    ```bash
    m3tal dash up
    ```

*   **`m3tal dash down`**: Stops the dashboard container.

    ```bash
    m3tal dash down
    ```

*   **`m3tal dash restart`**: Restarts the dashboard container.

    ```bash
    m3tal dash restart
    ```

*   **`m3tal dash logs`**: Streams the logs from the dashboard container in real-time.

    ```bash
    m3tal dash logs
    ```

*   **`m3tal dash status`**: Shows the current status of the dashboard container.

    ```bash
    m3tal dash status
    ```

### Stack Management

Commands for managing all deployed stacks via Docker Compose.

*   **`m3tal up`**: Runs `docker compose up` across all `*-compose.yml` files found in `/docker/`, starting all managed services.

    ```bash
    m3tal up
    ```

*   **`m3tal down`**: Runs `docker compose down` across all stacks, stopping all managed services.

    ```bash
    m3tal down
    ```

*   **`m3tal logs`**: Streams aggregated logs from all running M3TAL stacks.

    ```bash
    m3tal logs
    ```

---

## Systemd Service Management

The M3TAL API daemon is managed by systemd. You can interact with it using `systemctl` and `journalctl`.

*   **`systemctl status m3tal-api`**: Check the status of the `m3tal-api` service.

    ```bash
    systemctl status m3tal-api
    ```

*   **`journalctl -u m3tal-api -f`**: View and follow the logs for the `m3tal-api` service in real-time.

    ```bash
    journalctl -u m3tal-api -f
    ```

---

## Direct Docker Compose Commands (Fallback)

In scenarios where direct control is needed, you can use `docker compose` commands on the individual Compose files located in `/docker/`.

*   **Starting a specific stack (e.g., `my-stack-compose.yml`)**:

    ```bash
    docker compose -f /docker/my-stack-compose.yml up -d
    ```

*   **Stopping a specific stack**:

    ```bash
    docker compose -f /docker/my-stack-compose.yml down
    ```

*   **Viewing logs for a specific service within a stack (e.g., `my-service` in `my-stack-compose.yml`)**:

    ```bash
    docker compose -f /docker/my-stack-compose.yml logs -f my-service
    ```

---

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