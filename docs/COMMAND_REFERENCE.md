```markdown
# M3TAL CLI Command Reference

As DocSmith, the M3TAL Ecosystem Documentation Architect, my purpose is to provide you with a comprehensive and precise guide to managing your M3TAL environment. This document serves as your definitive cheat-sheet for the `m3tal` command-line interface, detailing every command, its purpose, and practical usage examples.

The M3TAL ecosystem is built upon robust open-source technologies, primarily leveraging Docker Engine and Docker Compose V2 for container orchestration, a Go-based API daemon for core logic, and systemd for service management. This document also covers the essential interactions with these underlying systems.

---

## M3TAL System Architecture Overview

Before diving into commands, let's establish the foundational architecture of M3TAL. Understanding these components and their interactions is crucial for effective management.

### Core Components

*   **CLI Binary (`/usr/bin/m3tal`)**: Your primary interface. A unified Go binary installed via APT, serving as the single entrypoint for all M3TAL operations.
*   **API Daemon (`m3tal-api.service`)**: The brain of M3TAL. A Go binary running as a systemd service on `http://localhost:8080`. It manages Docker interactions, maintains the state database (`/var/lib/m3tal/state.db`), and exposes API routes.
*   **M3TAL Dashboard (`m3tal-dashboard` container)**: Your graphical control panel. A Python/Flask container running on port 8082. It communicates with the API daemon at `http://host.docker.internal:8080` (from within the Docker network).
*   **Traefik Gateway (`routing-compose.yml`)**: The intelligent front door. A reverse proxy container exposing services by domain name on host port 80 (HTTP) and 443 (HTTPS, if configured). It uses Docker labels and a file provider for dynamic routing.
*   **Cloudflared (Optional, via `routing-compose.yml`)**: For secure, zero-config internet access to your services via Cloudflare Tunnels, bypassing port forwarding.

### Filesystem Contract

M3TAL adheres to a strict filesystem layout for reliability and predictability:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. **Managed by `m3tal config wizard`.**    |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.   |
| `/opt/m3tal/stack/`         | Canonical stack directory. Contains M3TAL's core compose files and Traefik config. |
| `/docker`                   | **User-facing symlink to `/opt/m3tal/stack/`. All user stack operations should reference this path.** |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |
| `/docker/dynamic/`          | Traefik file provider directory for dynamic configuration (e.g., API routing). |

### Docker / Compose Runtime

M3TAL is deeply integrated with **Docker Engine** and **Docker Compose V2**. These are hard dependencies.

*   The `m3tal up` command orchestrates all `*-compose.yml` files found within the `/docker/` directory.
*   The `m3tal dash up` command specifically manages the M3TAL dashboard container, which is defined in `/docker/m3tal-compose.yml`.
*   Your custom application stacks should reside in `/docker/`. Simply place a `my-app-compose.yml` file there, and M3TAL will manage it.
*   All compose files share a unified environment configuration provided by `/etc/m3tal/.env`, which is injected via the `--env-file` flag during `docker compose` operations.

### Deployment Lifecycle — Day 2 Operations

To install and manage a new application stack within M3TAL:

1.  **Place your compose file:** Create or copy your Docker Compose file (e.g., `my-stack-compose.yml`) into the `/docker/` directory.
2.  **Configure environment variables:** Ensure any required variables for your stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` for an interactive setup or `m3tal config set KEY value` for specific variables.
3.  **Start your stack:** Run `m3tal up` to start all services defined in `/docker/`, including your new stack. Alternatively, for fine-grained control over a single stack, use `docker compose -f /docker/my-stack-compose.yml --env-file /etc/m3tal/.env up -d`.

### Traefik Routing Architecture

Traefik is deployed via `routing-compose.yml` and acts as the central ingress point for all web traffic.

*   **Entry Point**: It binds to port 80 (HTTP) and typically port 443 (HTTPS) on the host system.
*   **Service Discovery**: Traefik automatically discovers and routes traffic to containers that have appropriate Docker labels (e.g., `traefik.http.routers.myapp.rule=Host(\`app.domain.com\`)`).
*   **Dynamic Configuration**: It loads additional routing rules from `/docker/dynamic/` via its file provider, enabling hot-reloads of configuration changes.
*   **M3TAL API & Dashboard Routing**:
    *   `api.YOUR_DOMAIN` is routed directly to the M3TAL API daemon running on the host at `http://host.docker.internal:8080` (defined in `dynamic/api.yml`).
    *   `dash.YOUR_DOMAIN` is routed to the M3TAL dashboard container (via labels in `m3tal-compose.traefik.yml`).
*   **Traefik Dashboard**: The Traefik dashboard itself is accessible locally at `http://localhost:8081`.

**Traefik Static Configuration (`traefik.yml`):**
(Located in `/opt/m3tal/stack/traefik/config/`)
```yaml
entryPoints:
  web:
    address: ":80"
  websecure: # Assumes HTTPS is enabled
    address: ":443"

providers:
  docker:
    exposedByDefault: false
    network: proxy # M3TAL's shared Docker network for services
  file:
    directory: /etc/traefik/dynamic # Maps to /docker/dynamic/ on host
    watch: true
```

**Dynamic Routing Example (`/docker/dynamic/api.yml`):**
```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)" # Uses the DOMAIN env var from .env
      service: api
      entryPoints:
        - web # Or websecure for HTTPS

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080" # Routes to the host API daemon
```

### Port Map

| Port | Service                                | Access                 |
| :--- | :------------------------------------- | :--------------------- |
| 80   | Traefik HTTP entry point               | Public (via DNS)       |
| 443  | Traefik HTTPS entry point              | Public (via DNS)       |
| 8080 | M3TAL API daemon (Go)                  | Host-local (via Traefik or direct) |
| 8081 | Traefik dashboard                      | Host-local only        |
| 8082 | M3TAL Dashboard (Python/Flask)         | Via Traefik or direct  |

### APT Installation

If you haven't already, install M3TAL using the official APT repository:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
```

---

## M3TAL CLI Command Reference

### Core M3TAL Commands

#### `sudo m3tal`

Opens the interactive Text User Interface (TUI) Control Center. This provides a user-friendly, menu-driven interface for common M3TAL operations without needing to remember specific CLI commands.

```bash
# Access the interactive TUI
sudo m3tal
```

#### `m3tal init`

Generates the initial `/etc/m3tal/.env` configuration file from M3TAL's default values. This command should be run upon first installation or if your `.env` file is missing. It will *not* overwrite an existing file unless forced (not recommended).

```bash
# Initialize the M3TAL environment file
m3tal init
```

#### `m3tal doctor`

Performs a pre-flight health check of your M3TAL system. This command verifies crucial dependencies and configurations, including Docker connectivity, the validity of `/etc/m3tal/.env`, and the availability of essential network ports.

```bash
# Run a diagnostic check on the M3TAL system
m3tal doctor
```

### Configuration Management (`m3tal config`)

The `m3tal config` subcommand group allows for precise control over your `/etc/m3tal/.env` configuration.

#### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. This is the recommended method for initial setup and major configuration changes, ensuring all required variables are set correctly.

```bash
# Start the interactive configuration wizard
m3tal config wizard
```

#### `m3tal config set KEY VALUE`

Sets a single environment variable within `/etc/m3tal/.env` to the specified value. This command is useful for quick, targeted adjustments.

```bash
# Set the primary domain for M3TAL services
m3tal config set DOMAIN "your.domain.com"

# Update the timezone
m3tal config set TZ "America/New_York"
```

#### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

```bash
# Get the currently configured domain
m3tal config get DOMAIN

# Retrieve the configured API token
m3tal config get API_TOKEN
```

#### `m3tal config scan`

Lists all environment variables that M3TAL understands, along with their current values (if set in `.env`) and default values. This helps in understanding available configuration options across all potential stacks.

```bash
# List all known environment variables and their states
m3tal config scan
```

#### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file, showing all explicitly set configuration variables.

```bash
# Display the full contents of the M3TAL environment file
m3tal config list
```

### Dashboard Management (`m3tal dash`)

The `m3tal dash` subcommand group is dedicated to managing the M3TAL Dashboard container and its credentials.

#### `m3tal dashpass [username] [password]`

Updates the password for a specified dashboard user. If `username` and `password` are omitted, the command will prompt you interactively for the necessary details. This manages the `/docker/users.json` file.

```bash
# Set the password for the 'admin' user interactively
m3tal dashpass

# Set the password for 'admin' directly (non-interactive)
m3tal dashpass admin MySecurePa$$w0rd!
```

#### `m3tal dash up`

Pulls the latest dashboard Docker Compose configuration from GitHub, then starts or updates the M3TAL dashboard container using `/docker/m3tal-compose.yml`. This also ensures the M3TAL API daemon (`m3tal-api.service`) is running.

```bash
# Pull and start the M3TAL Dashboard
m3tal dash up
```

#### `m3tal dash down`

Stops and removes the M3TAL dashboard container.

```bash
# Stop the M3TAL Dashboard
m3tal dash down
```

#### `m3tal dash restart`

Restarts the M3TAL dashboard container. This is useful for applying configuration changes or resolving minor issues.

```bash
# Restart the M3TAL Dashboard
m3tal dash restart
```

#### `m3tal dash logs`

Streams the aggregated logs from the M3TAL dashboard container, useful for debugging and monitoring its activity.

```bash
# Stream M3TAL Dashboard logs
m3tal dash logs
```

#### `m3tal dash status`

Shows the current operational status of the M3TAL dashboard container.

```bash
# Check the status of the M3TAL Dashboard
m3tal dash status
```

### Stack Management (`m3tal`)

These commands interact with all Docker Compose stacks defined in the `/docker/` directory.

#### `m3tal up`

Executes `docker compose up -d` across all `*-compose.yml` files found in the `/docker/` directory. This command starts or recreates all services defined in your M3TAL ecosystem, ensuring they are running in detached mode.

```bash
# Start or update all M3TAL and user-defined Docker Compose stacks
m3tal up
```

#### `m3tal down`

Executes `docker compose down` across all `*-compose.yml` files in the `/docker/` directory. This command stops and removes all services and networks associated with your M3TAL ecosystem.

```bash
# Stop and remove all M3TAL and user-defined Docker Compose stacks
m3tal down
```

#### `m3tal logs`

Streams aggregated logs from all currently running Docker Compose stacks managed by M3TAL. This provides a consolidated view of activity across your entire containerized environment.

```bash
# Stream aggregated logs from all running M3TAL stacks
m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon is managed as a `systemd` service. Understanding these commands is crucial for debugging and direct control over the daemon's lifecycle.

#### `systemctl status m3tal-api`

Checks the current status of the M3TAL API daemon service. This command shows if the service is active, its uptime, and recent log entries.

```bash
# Check the status of the M3TAL API daemon
systemctl status m3tal-api
```

#### `systemctl restart m3tal-api`

Restarts the M3TAL API daemon service. This is useful if you make changes that require the API to reload, or if the API is experiencing issues.

```bash
# Restart the M3TAL API daemon
sudo systemctl restart m3tal-api
```

#### `journalctl -u m3tal-api -f`

Streams the real-time logs for the M3TAL API daemon service directly from `journald`. This is invaluable for live debugging and observing the daemon's behavior.

```bash
# Stream live logs from the M3TAL API daemon
sudo journalctl -u m3tal-api -f
```

---

## Direct Docker Compose Fallback

While the `m3tal` CLI provides a convenient abstraction, you might occasionally need to interact directly with Docker Compose for specific, advanced, or troubleshooting scenarios. M3TAL uses Docker Compose V2, which means the command is `docker compose` (without the hyphen).

Remember that M3TAL's compose files are located in `/docker/` (symlinked from `/opt/m3tal/stack/`) and they all rely on the shared environment file at `/etc/m3tal/.env`.

#### Starting a specific stack

To bring up only a single stack (e.g., your custom `my-app-compose.yml`) without affecting others:

```bash
# Start a specific application stack in detached mode
sudo docker compose -f /docker/my-app-compose.yml --env-file /etc/m3tal/.env up -d
```

#### Stopping a specific stack

To stop and remove services for a particular stack:

```bash
# Stop and remove services for a specific stack
sudo docker compose -f /docker/my-app-compose.yml --env-file /etc/m3tal/.env down
```

#### Viewing logs for specific services

To stream logs from services within a particular stack:

```bash
# Stream logs from the 'dashboard' service within the M3TAL dashboard stack
sudo docker compose -f /docker/m3tal-compose.yml --env-file /etc/m3tal/.env logs -f dashboard
```

#### Listing all containers for a stack

To view the status of containers defined in a specific compose file:

```bash
# List containers for the M3TAL routing stack
sudo docker compose -f /docker/routing-compose.yml --env-file /etc/m3tal/.env ps
```

---

This concludes your M3TAL CLI cheat-sheet. With these commands and architectural insights, you are well-equipped to manage your M3TAL ecosystem effectively. Happy orchestrating!
```