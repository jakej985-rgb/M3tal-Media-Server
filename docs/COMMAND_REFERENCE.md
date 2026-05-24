```markdown
# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for all available `m3tal` CLI commands and their usage.

## Getting Started

### Installation

To install M3TAL, follow these steps:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

### Systemd Service Management

M3TAL's API daemon is managed by systemd. You can use `systemctl` to control and monitor it:

- **Check status:**
  ```bash
  sudo systemctl status m3tal-api
  ```

- **Restart the service:**
  ```bash
  sudo systemctl restart m3tal-api
  ```

- **View logs:**
  ```bash
  sudo journalctl -u m3tal-api -f
  ```

## Core M3TAL Commands

### Interactive TUI

- **`sudo m3tal`**:
  Opens the interactive TUI Control Center, presenting a numbered menu for common operations.
  ```bash
  sudo m3tal
  ```

### Initialization and Configuration

- **`m3tal init`**:
  Generates the default `/etc/m3tal/.env` configuration file. Use this on your first installation.
  ```bash
  m3tal init
  ```

- **`m3tal config wizard`**:
  Launches an interactive wizard to guide you through configuring your `/etc/m3tal/.env` file.
  ```bash
  m3tal config wizard
  ```

- **`m3tal config set KEY VALUE`**:
  Sets a single environment variable in your `/etc/m3tal/.env` file.
  ```bash
  m3tal config set DOMAIN mydomain.com
  ```

- **`m3tal config get KEY`**:
  Reads and displays the value of a single environment variable from your `/etc/m3tal/.env` file.
  ```bash
  m3tal config get DOMAIN
  ```

- **`m3tal config scan`**:
  Lists all environment variables across all M3TAL stacks, including their current and default values.
  ```bash
  m3tal config scan
  ```

- **`m3tal config list`**:
  Displays the current contents of your `/etc/m3tal/.env` file.
  ```bash
  m3tal config list
  ```

### Dashboard Management

- **`m3tal dashpass [username] [password]`**:
  Updates the password for a dashboard user. If `username` and `password` are omitted, it will prompt interactively.
  ```bash
  # Interactive usage
  m3tal dashpass

  # With username and password
  m3tal dashpass admin new_secure_password
  ```

- **`m3tal dash up`**:
  Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.
  ```bash
  m3tal dash up
  ```

- **`m3tal dash down`**:
  Stops the dashboard container.
  ```bash
  m3tal dash down
  ```

- **`m3tal dash restart`**:
  Restarts the dashboard container.
  ```bash
  m3tal dash restart
  ```

- **`m3tal dash logs`**:
  Streams the logs from the dashboard container in real-time.
  ```bash
  m3tal dash logs
  ```

- **`m3tal dash status`**:
  Shows the current status of the dashboard container.
  ```bash
  m3tal dash status
  ```

### Stack Management

- **`m3tal up`**:
  Runs `docker compose up` across all `*-compose.yml` files found in the `/docker/` directory, starting all your defined stacks.
  ```bash
  m3tal up
  ```

- **`m3tal down`**:
  Runs `docker compose down` across all your defined stacks.
  ```bash
  m3tal down
  ```

- **`m3tal logs`**:
  Streams aggregated logs from all currently running M3TAL stacks.
  ```bash
  m3tal logs
  ```

## Docker Compose Fallback

M3TAL uses Docker Compose under the hood. In situations where direct control is needed, you can use `docker compose` commands on the individual compose files within `/docker/`.

- **Starting a specific stack (e.g., `my-app-compose.yml`):**
  ```bash
  docker compose -f /docker/my-app-compose.yml up -d
  ```

- **Stopping a specific stack:**
  ```bash
  docker compose -f /docker/my-app-compose.yml down
  ```

- **Viewing logs for a specific stack:**
  ```bash
  docker compose -f /docker/my-app-compose.yml logs -f
  ```

- **Bringing up all stacks directly (equivalent to `m3tal up`):**
  ```bash
  docker compose $(find /docker/ -name '*-compose.yml' | paste -sd ' -f ') up -d
  ```

- **Bringing down all stacks directly (equivalent to `m3tal down`):**
  ```bash
  docker compose $(find /docker/ -name '*-compose.yml' | paste -sd ' -f ') down
  ```

## Important Directories and Files

- **`/etc/m3tal/.env`**: The primary configuration file for M3TAL.
- **`/var/lib/m3tal/state.db`**: SQLite database for M3TAL state.
- **`/docker/`**: Symlink to `/opt/m3tal/stack/`. This is where all your stack compose files reside.
- **`/docker/users.json`**: Stores dashboard credentials.

## Dashboard Access Modes

The dashboard can be accessed in two modes, controlled by `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`:

### 1. `local` (Default)

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
-   **Access:** `http://HOST_IP:8082` or `http://localhost:8082`
-   **Usage:** Ideal for LAN-only setups and initial testing.

### 2. `traefik`

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
-   **Access:** `http://dash.DOMAIN` (Requires Traefik to be running via `m3tal up`)
-   **Usage:** Suitable for domain-based setups with Traefik as a reverse proxy.

---

*This document was generated by DocSmith, the M3TAL Ecosystem Documentation Architect.*
```