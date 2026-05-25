# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for all available `m3tal` CLI commands, covering system management, configuration, and container operations.

## Table of Contents

*   [Interactive TUI](#interactive-tui)
*   [System Initialization & Health](#system-initialization--health)
*   [Configuration Management](#configuration-management)
*   [Dashboard Management](#dashboard-management)
*   [System-wide Operations](#system-wide-operations)
*   [Systemd Service Management](#systemd-service-management)
*   [Direct Docker Compose Fallback](#direct-docker-compose-fallback)

---

## Interactive TUI

This command launches the M3TAL TUI Control Center, offering a numbered menu for various operations.

*   **Command:** `sudo m3tal`
*   **Description:** Opens the interactive TUI Control Center.
*   **Usage Example:**
    ```bash
    sudo m3tal
    ```

---

## System Initialization & Health

Commands for setting up M3TAL and verifying its operational status.

### Initialize M3TAL Environment

Generates the default `/etc/m3tal/.env` file. This should be run on first installation.

*   **Command:** `m3tal init`
*   **Description:** Generates `/etc/m3tal/.env` from defaults. Use on first install.
*   **Usage Example:**
    ```bash
    m3tal init
    ```

### Pre-flight Health Check

Performs a series of checks to ensure M3TAL is ready to operate, including Docker connectivity, `.env` file validity, and port availability.

*   **Command:** `m3tal doctor`
*   **Description:** Pre-flight health check: Docker connectivity, .env validity, port availability.
*   **Usage Example:**
    ```bash
    m3tal doctor
    ```

---

## Configuration Management

Commands for managing M3TAL's environment variables and configuration files.

### Configuration Wizard

An interactive wizard to guide you through configuring the `/etc/m3tal/.env` file.

*   **Command:** `m3tal config wizard`
*   **Description:** Interactive wizard to configure `/etc/m3tal/.env`.
*   **Usage Example:**
    ```bash
    m3tal config wizard
    ```

### Set Single Environment Variable

Sets a specific environment variable in the `.env` file.

*   **Command:** `m3tal config set KEY VALUE`
*   **Description:** Set a single env var.
*   **Usage Example:**
    ```bash
    m3tal config set DOMAIN mydomain.com
    ```

### Get Single Environment Variable

Reads and displays the value of a specific environment variable from the `.env` file.

*   **Command:** `m3tal config get KEY`
*   **Description:** Read a single env var.
*   **Usage Example:**
    ```bash
    m3tal config get DOMAIN
    ```

### Scan All Environment Variables

Lists all environment variables across all configured stacks, showing their current values.

*   **Command:** `m3tal config scan`
*   **Description:** List all env vars across all stacks.
*   **Usage Example:**
    ```bash
    m3tal config scan
    ```

### List Current .env File Contents

Displays the entire content of the current `/etc/m3tal/.env` file.

*   **Command:** `m3tal config list`
*   **Description:** List current .env file contents.
*   **Usage Example:**
    ```bash
    m3tal config list
    ```

---

## Dashboard Management

Commands for managing the M3TAL dashboard container.

### Update Dashboard Password

Updates the password for a dashboard user. If no username or password is provided, it will prompt interactively.

*   **Command:** `m3tal dashpass [username] [password]`
*   **Description:** Update dashboard user password. Interactive if args omitted.
*   **Usage Example (interactive):**
    ```bash
    sudo m3tal dashpass
    ```
*   **Usage Example (with arguments):**
    ```bash
    sudo m3tal dashpass admin new_secret_password
    ```

### Update and Start Dashboard

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

*   **Command:** `m3tal dash up`
*   **Description:** Pull latest dashboard compose config from GitHub, then start the dashboard container.
*   **Usage Example:**
    ```bash
    sudo m3tal dash up
    ```

### Stop Dashboard Container

Stops the M3TAL dashboard container.

*   **Command:** `m3tal dash down`
*   **Description:** Stop the dashboard container.
*   **Usage Example:**
    ```bash
    sudo m3tal dash down
    ```

### Restart Dashboard Container

Restarts the M3TAL dashboard container.

*   **Command:** `m3tal dash restart`
*   **Description:** Restart the dashboard container.
*   **Usage Example:**
    ```bash
    sudo m3tal dash restart
    ```

### Stream Dashboard Logs

Streams the logs from the M3TAL dashboard container in real-time.

*   **Command:** `m3tal dash logs`
*   **Description:** Stream dashboard container logs.
*   **Usage Example:**
    ```bash
    sudo m3tal dash logs
    ```

### Show Dashboard Container Status

Displays the current status of the M3TAL dashboard container.

*   **Command:** `m3tal dash status`
*   **Description:** Show dashboard container status.
*   **Usage Example:**
    ```bash
    sudo m3tal dash status
    ```

---

## System-wide Operations

Commands that affect all running M3TAL stacks.

### Start All Stacks

Runs `docker compose up` across all `*-compose.yml` files located in `/docker/`. This starts all your deployed services.

*   **Command:** `m3tal up`
*   **Description:** Run docker compose up across all *-compose.yml files in /docker/.
*   **Usage Example:**
    ```bash
    sudo m3tal up -d
    ```
    *(Note: `-d` is commonly used to run in detached mode)*

### Stop All Stacks

Runs `docker compose down` for all services managed by M3TAL.

*   **Command:** `m3tal down`
*   **Description:** Run docker compose down across all stacks.
*   **Usage Example:**
    ```bash
    sudo m3tal down
    ```

### Stream All Logs

Aggregates and streams logs from all currently running M3TAL stacks.

*   **Command:** `m3tal logs`
*   **Description:** Stream aggregated logs from all running stacks.
*   **Usage Example:**
    ```bash
    sudo m3tal logs -f
    ```
    *(Note: `-f` follows the log output)*

---

## Systemd Service Management

The M3TAL API daemon runs as a systemd service. You can manage it using `systemctl` and view its logs with `journalctl`.

### Check API Service Status

Shows the current status of the `m3tal-api.service`.

*   **Command:** `systemctl status m3tal-api`
*   **Description:** Checks the status of the M3TAL API daemon.
*   **Usage Example:**
    ```bash
    systemctl status m3tal-api
    ```

### Restart API Service

Restarts the `m3tal-api.service`.

*   **Command:** `systemctl restart m3tal-api`
*   **Description:** Restarts the M3TAL API daemon.
*   **Usage Example:**
    ```bash
    sudo systemctl restart m3tal-api
    ```

### View API Service Logs

Streams the logs from the `m3tal-api.service` in real-time.

*   **Command:** `journalctl -u m3tal-api -f`
*   **Description:** Streams logs from the M3TAL API daemon.
*   **Usage Example:**
    ```bash
    sudo journalctl -u m3tal-api -f
    ```

---

## Direct Docker Compose Fallback

In scenarios where the `m3tal` CLI commands might not be sufficient or for direct control, you can interact with Docker Compose directly. M3TAL orchestrates Docker Compose files located in `/docker/`.

### Example: Starting a Specific Stack

To start a stack defined in `/docker/my-stack-compose.yml`:

*   **Command:** `docker compose -f /docker/my-stack-compose.yml up -d`
*   **Usage Example:**
    ```bash
    sudo docker compose -f /docker/my-stack-compose.yml up -d
    ```

### Example: Stopping a Specific Stack

To stop a stack defined in `/docker/my-stack-compose.yml`:

*   **Command:** `docker compose -f /docker/my-stack-compose.yml down`
*   **Usage Example:**
    ```bash
    sudo docker compose -f /docker/my-stack-compose.yml down
    ```

### Example: Viewing Logs for a Specific Stack Service

To view logs for the `my-service` container within a stack defined in `/docker/my-stack-compose.yml`:

*   **Command:** `docker compose -f /docker/my-stack-compose.yml logs -f my-service`
*   **Usage Example:**
    ```bash
    sudo docker compose -f /docker/my-stack-compose.yml logs -f my-service
    ```

### Example: Listing All Docker Compose Projects

To see all Docker Compose projects M3TAL manages (which are typically the individual compose files in `/docker/`):

*   **Command:** `docker compose ls`
*   **Usage Example:**
    ```bash
    sudo docker compose ls
    ```

---

## APT Installation

This section details the steps required to install M3TAL via APT.

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```