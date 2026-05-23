Greetings, intrepid M3TAL operator! DocSmith here, your M3TAL Ecosystem Documentation Architect. You've landed on the definitive guide to wrangling the M3TAL CLI, the unified entry point to your self-hosted infrastructure. This cheat-sheet will equip you with every command you need, from initial setup to day-to-day operations and advanced troubleshooting.

## M3TAL Ecosystem Overview

The M3TAL ecosystem is designed for streamlined, self-hosted deployment of your services. It leverages **Docker Engine** and **Docker Compose V2** as its foundation, orchestrated by a powerful Go-based **API daemon** running as a `systemd` service (`m3tal-api.service`). All commands are channeled through the unified `/usr/bin/m3tal` CLI binary.

Your services, or "stacks," are defined by `*-compose.yml` files located in the `/docker/` directory (which is a symlink to `/opt/m3tal/stack/`). M3TAL automatically discovers and manages these stacks.

### Core Components
-   **M3TAL CLI (`/usr/bin/m3tal`):** Your primary interaction point.
-   **M3TAL API Daemon (`m3tal-api.service`):** The brains of the operation, running on port 8080. It manages Docker, maintains state in `/var/lib/m3tal/state.db`, and exposes API routes.
-   **M3TAL Dashboard (`m3tal-dashboard` container):** A Python/Flask container providing a user-friendly web interface, communicating with the API daemon.
-   **Traefik Gateway (`routing-compose.yml`):** A reverse proxy container (on port 80) for domain-based routing, automatically discovering services via Docker labels.
-   **Cloudflared (`routing-compose.yml`):** An optional Cloudflare tunnel container for secure, zero-config internet access to your services.

### Filesystem Contract

M3TAL relies on a strict filesystem layout for configuration and data persistence:

| Path                        | Purpose                                                                                                                                                                                                                                                                                                                                     |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `/etc/m3tal/.env`           | **Primary configuration file.** Contains all environment variables for M3TAL and your Docker stacks. **Managed by `m3tal config wizard`**. Do not edit manually unless you know precisely what you're doing.                                                                                                                                    |
| `/var/lib/m3tal/state.db`   | SQLite state database. Automatically created and managed by the M3TAL API daemon. Stores internal state, service information, and other operational data.                                                                                                                                                                                     |
| `/opt/m3tal/stack/`         | The canonical directory for M3TAL's core Docker Compose files (e.g., `routing-compose.yml`, `m3tal-compose.yml`) and dynamic Traefik configurations.                                                                                                                                                                                          |
| `/docker`                   | **User-facing symlink to `/opt/m3tal/stack/`.** This is where you place *your* custom `*-compose.yml` files for new services. M3TAL scans this directory for all `*-compose.yml` files when running `m3tal up`.                                                                                                                            |
| `/docker/users.json`        | Dashboard credential store. Contains hashed user credentials for the M3TAL Dashboard. **Managed by `m3tal dashpass`**.                                                                                                                                                                                                                      |
| `/docker/dynamic/api.yml`   | Traefik dynamic configuration file, routing `api.${DOMAIN}` to the M3TAL API daemon.                                                                                                                                                                                                                                                        |
| `/docker/m3tal-compose.yml` | Base Docker Compose file for the M3TAL Dashboard. Updated by `m3tal dash up`.                                                                                                                                                                                                                                                               |
| `/docker/routing-compose.yml` | Base Docker Compose file for Traefik and Cloudflared.                                                                                                                                                                                                                                                                                       |
| `/docker/m3tal-compose.local.yml` | Override for the M3TAL Dashboard when `DASHBOARD_EXPOSE_MODE=local`. Adds a direct port binding.                                                                                                                                                                                                                                             |
| `/docker/m3tal-compose.traefik.yml` | Override for the M3TAL Dashboard when `DASHBOARD_EXPOSE_MODE=traefik`. Adds Traefik labels for routing.                                                                                                                                                                                                                                   |

### Dashboard Access Modes (CRITICAL!)

The M3TAL Dashboard provides a web-based interface for managing your ecosystem. Its accessibility is controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

#### Mode 1: `local` (Default)

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=local`
-   **Mechanism:** M3TAL uses the `m3tal-compose.local.yml` override, which directly binds the dashboard container's internal port 8082 to the host's `DASHBOARD_PORT` (defaulting to 8082). Traefik is **not** involved in routing to the dashboard in this mode.
-   **Access Via:** `http://YOUR_HOST_IP:8082` or `http://localhost:8082` (if accessing from the host machine itself).
-   **Best For:** LAN-only setups, first-time installations, local development, or when you don't need domain-based access for the dashboard.

#### Mode 2: `traefik`

-   **Configuration:** `DASHBOARD_EXPOSE_MODE=traefik`
-   **Mechanism:** M3TAL uses the `m3tal-compose.traefik.yml` override, which adds specific Docker labels to the dashboard container. Traefik, if running (via `m3tal up`), discovers these labels and automatically routes `dash.${DOMAIN}` to the dashboard on port 8082.
-   **Access Via:** `http://dash.YOUR_DOMAIN` (e.g., `http://dash.example.com`). This requires Traefik to be running via `m3tal up` and your `DOMAIN` variable to be correctly configured in `/etc/m3tal/.env`.
-   **Best For:** Domain-based access, integrating the dashboard with other services behind a Traefik reverse proxy, and external access when `Cloudflared` is also used.

### Deployment Lifecycle: Adding a New Stack

1.  **Create your Compose File:** Develop your `docker-compose.yml` file for your new service.
2.  **Place in `/docker/`:** Save your compose file in the `/docker/` directory, e.g., `/docker/my-app-compose.yml`.
3.  **Configure Environment Variables:** Ensure any necessary environment variables for your stack are defined in `/etc/m3tal/.env`. Use `m3tal config wizard` or `m3tal config set KEY value` to manage these.
4.  **Start all Stacks:** Run `m3tal up` to deploy your new stack alongside all existing M3TAL services.

### Traefik Routing Architecture

Traefik, deployed via `routing-compose.yml`, acts as the reverse proxy for your M3TAL ecosystem.

-   **Port Binding:** Binds host port 80 (HTTP) as its main entry point.
-   **Service Discovery:** Automatically discovers services via Docker labels (e.g., those added by `m3tal-compose.traefik.yml` for the dashboard).
-   **Dynamic Configuration:** Loads additional dynamic routing rules from `/docker/dynamic/` (e.g., `api.yml` for the M3TAL API daemon). These configurations hot-reload without restarting Traefik.
-   **M3TAL API Routing:** Routes `api.${DOMAIN}` to the M3TAL API daemon (`http://host.docker.internal:8080`) via `/docker/dynamic/api.yml`.

### Port Map

| Port | Service                     | Access                                      |
| :--- | :-------------------------- | :------------------------------------------ |
| 80   | Traefik HTTP Entry Point    | Public (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API Daemon (Go)       | Host-local only                             |
| 8081 | Traefik Dashboard           | Host-local only                             |
| 8082 | M3TAL Dashboard (container) | Direct (local mode) or via Traefik (traefik mode) |

### APT Installation

To install the M3TAL CLI and API daemon:

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

This section details every command available in the M3TAL CLI, complete with real-world usage examples.

### Core M3TAL Commands

#### `sudo m3tal`

Opens the interactive M3TAL TUI Control Center. This provides a numbered menu for common operations and real-time status updates.
**Note:** `sudo` is required as many operations interact with Docker and system-level files.

**Description:** Launches the main interactive Text User Interface (TUI) for M3TAL. From here, you can navigate menus to perform configuration, stack management, dashboard operations, and more, all without remembering specific CLI commands. It's the central hub for most day-to-day tasks.

**Usage Example:**

```bash
sudo m3tal
```

#### `m3tal init`

Generates the primary configuration file, `/etc/m3tal/.env`, from built-in defaults. This command is crucial for a first-time setup or to reset the configuration.

**Description:** Initializes the M3TAL configuration. If `/etc/m3tal/.env` does not exist, it creates it with default values for all essential environment variables. If it already exists, it will prompt to overwrite, which can be useful for restoring default settings or if the file becomes corrupted.

**Usage Example:**

```bash
sudo m3tal init
```

#### `m3tal doctor`

Performs a pre-flight health check of the M3TAL system.

**Description:** Runs a series of diagnostic checks to ensure the M3TAL environment is ready for operation. This includes verifying Docker daemon connectivity, checking the validity and permissions of `/etc/m3tal/.env`, and confirming that essential ports (like 80, 8080, 8082) are not in use by other applications, preventing conflicts before starting services.

**Usage Example:**

```bash
sudo m3tal doctor
```

### Configuration Commands (`m3tal config`)

These commands are used to manage the `/etc/m3tal/.env` configuration file.

#### `m3tal config wizard`

Launches an interactive wizard to guide you through configuring `/etc/m3tal/.env`.

**Description:** Provides a step-by-step interactive prompt to configure all essential M3TAL environment variables. This is the recommended way to set up your system after `m3tal init`, allowing you to define domain names, expose modes, user IDs, and other critical parameters safely and easily.

**Usage Example:**

```bash
sudo m3tal config wizard
```

#### `m3tal config set KEY VALUE`

Sets a single environment variable in `/etc/m3tal/.env` to a specified value.

**Description:** Allows for direct modification of individual configuration keys. This is useful for making quick, atomic changes without going through the full wizard. The command will validate the key and attempt to update the value in the `.env` file.

**Usage Example:**

```bash
sudo m3tal config set DASHBOARD_EXPOSE_MODE traefik
sudo m3tal config set DOMAIN example.com
```

#### `m3tal config get KEY`

Retrieves and displays the current value of a specific environment variable from `/etc/m3tal/.env`.

**Description:** Fetches and prints the value associated with a given configuration key. Useful for quickly checking the current setting of a specific variable without opening the entire `.env` file.

**Usage Example:**

```bash
sudo m3tal config get DOMAIN
```

#### `m3tal config scan`

Lists all detected environment variables across all active Docker Compose stacks, including their sources.

**Description:** Scans all `*-compose.yml` files in `/docker/` and the `/etc/m3tal/.env` file to identify all environment variables used by the M3TAL ecosystem. It helps in understanding which variables are expected by which services and what their default or configured values are.

**Usage Example:**

```bash
sudo m3tal config scan
```

#### `m3tal config list`

Displays the entire contents of the current `/etc/m3tal/.env` file.

**Description:** Prints the complete list of environment variables and their values as currently configured in `/etc/m3tal/.env`. This provides a full overview of your M3TAL system's primary configuration.

**Usage Example:**

```bash
sudo m3tal config list
```

### Dashboard Management Commands (`m3tal dash`)

These commands specifically manage the M3TAL Dashboard container and its access credentials.

#### `m3tal dashpass [username] [password]`

Updates or sets the password for a dashboard user in `/docker/users.json`. If no arguments are provided, it runs interactively.

**Description:** This command manages dashboard user credentials.
-   If run without arguments, it prompts for a username and a new password, then confirms the password.
-   If `username` and `password` are provided, it directly sets the password for that user, creating the user if they don't exist. This overwrites any existing password for the specified user.
The passwords are hashed and stored in `/docker/users.json`.

**Usage Examples:**

```bash
# Interactive mode (recommended for security)
sudo m3tal dashpass

# Non-interactive mode (use with caution, avoid exposing passwords in history)
sudo m3tal dashpass admin MySecureP@ssw0rd!
```

#### `m3tal dash up`

Pulls the latest dashboard compose configuration from GitHub and then starts the dashboard container.

**Description:** Ensures your M3TAL Dashboard is running with the most up-to-date configuration. It downloads `m3tal-compose.yml` and its overrides (`local.yml`, `traefik.yml`) to `/docker/`, then starts the `m3tal-dashboard` service using the appropriate override based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`.

**Usage Example:**

```bash
sudo m3tal dash up
```

#### `m3tal dash down`

Stops the M3TAL Dashboard container.

**Description:** Gracefully stops and removes the `m3tal-dashboard` container. This does not remove its volumes or configuration.

**Usage Example:**

```bash
sudo m3tal dash down
```

#### `m3tal dash restart`

Restarts the M3TAL Dashboard container.

**Description:** Stops the currently running `m3tal-dashboard` container and then starts it again. This is useful for applying minor configuration changes that affect only the dashboard or for troubleshooting.

**Usage Example:**

```bash
sudo m3tal dash restart
```

#### `m3tal dash logs`

Streams aggregated logs from the M3TAL Dashboard container.

**Description:** Displays the real-time output (stdout/stderr) from the `m3tal-dashboard` container. This is invaluable for debugging issues related to the dashboard, checking its startup process, or monitoring its activity.

**Usage Example:**

```bash
sudo m3tal dash logs
```

#### `m3tal dash status`

Shows the current status of the M3TAL Dashboard container.

**Description:** Provides a quick overview of whether the `m3tal-dashboard` container is running, stopped, restarting, or in another state, along with its container ID and uptime.

**Usage Example:**

```bash
sudo m3tal dash status
```

### Stack Management Commands (`m3tal`)

These commands manage all Docker Compose stacks defined in `/docker/`.

#### `m3tal up`

Runs `docker compose up -d` across all `*-compose.yml` files in `/docker/`.

**Description:** This is the primary command to deploy or update all your M3TAL services. It iterates through every `*-compose.yml` file found in `/docker/` (including `routing-compose.yml`, `m3tal-compose.yml`, and your custom stack files) and brings them up in detached mode (`-d`). It will create containers, networks, and volumes as defined in your compose files.

**Usage Example:**

```bash
sudo m3tal up
```

#### `m3tal down`

Runs `docker compose down` across all `*-compose.yml` files in `/docker/`.

**Description:** Gracefully stops and removes all containers, networks, and volumes (if not explicitly marked as external) defined by all `*-compose.yml` files in `/docker/`. Use this to shut down your entire M3TAL ecosystem.

**Usage Example:**

```bash
sudo m3tal down
```

#### `m3tal logs`

Streams aggregated logs from all running Docker Compose stacks.

**Description:** Collects and displays real-time log output from all containers managed by M3TAL. This provides a holistic view of the system's activity and is essential for diagnosing issues affecting multiple services.

**Usage Example:**

```bash
sudo m3tal logs
```

---

## Systemd Service Management

The M3TAL API daemon runs as a `systemd` service, `m3tal-api.service`. You can manage and monitor it using standard `systemctl` and `journalctl` commands.

#### `systemctl status m3tal-api`

Checks the current status of the M3TAL API daemon.

**Description:** Displays whether the `m3tal-api.service` is active (running), inactive (stopped), or in a failed state. It also shows recent log entries and other process information.

**Usage Example:**

```bash
systemctl status m3tal-api
```

#### `systemctl restart m3tal-api`

Restarts the M3TAL API daemon.

**Description:** Stops and then restarts the `m3tal-api.service`. This is useful after making manual changes to the `/etc/m3tal/.env` file that the API daemon needs to pick up, or for general troubleshooting.

**Usage Example:**

```bash
sudo systemctl restart m3tal-api
```

#### `journalctl -u m3tal-api -f`

Streams real-time logs from the M3TAL API daemon.

**Description:** Provides a continuous, real-time output of the logs generated by the `m3tal-api.service`. This is the most effective way to monitor the API daemon's operations, diagnose startup issues, or see how it interacts with Docker and other components.

**Usage Example:**

```bash
sudo journalctl -u m3tal-api -f
```

---

## Docker Direct Commands (Fallback)

While M3TAL provides a high-level abstraction, it's built on Docker Engine and Docker Compose V2. Knowing the direct Docker commands can be invaluable for advanced troubleshooting or specific operations that M3TAL's CLI might not directly expose.

**Important:** When using direct Docker Compose commands, always navigate to the `/docker/` directory (or `/opt/m3tal/stack/`) first, as M3TAL expects all compose files to be relative to this location.

#### Navigating to the Stacks Directory

```bash
cd /docker/
# or
cd /opt/m3tal/stack/
```

#### `docker ps`

Lists all running Docker containers.

**Description:** Shows a summary of all containers currently running on your system, including their ID, image, command, creation time, status, ports, and name.

**Usage Example:**

```bash
sudo docker ps
```

#### `docker compose -f <compose_file.yml> up -d`

Starts a specific Docker Compose stack in detached mode.

**Description:** This command will bring up only the services defined in the specified compose file. This is useful if you only want to manage a single stack without affecting others. Replace `<compose_file.yml>` with the actual filename.

**Usage Example:**

```bash
# Start only the routing stack
sudo docker compose -f routing-compose.yml up -d

# Start a custom stack you added
sudo docker compose -f my-media-stack-compose.yml up -d
```

#### `docker compose -f <compose_file.yml> down`

Stops and removes containers for a specific Docker Compose stack.

**Description:** Stops and removes containers, networks, and volumes associated with the specified compose file.

**Usage Example:**

```bash
# Stop only the routing stack
sudo docker compose -f routing-compose.yml down

# Stop a custom stack
sudo docker compose -f my-media-stack-compose.yml down
```

#### `docker compose -f <compose_file.yml> logs -f`

Streams logs for a specific Docker Compose stack.

**Description:** Displays real-time log output for all services defined within a particular compose file.

**Usage Example:**

```bash
# Stream logs for the routing stack
sudo docker compose -f routing-compose.yml logs -f

# Stream logs for the M3TAL dashboard
sudo docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml logs -f
```

#### `docker exec -it <container_name_or_id> /bin/bash`

Executes a command inside a running container, often to get a shell.

**Description:** Allows you to run commands directly within a running container. The `-it` flags provide an interactive TTY. `/bin/bash` or `/bin/sh` is commonly used to open a shell.

**Usage Example:**

```bash
# Get a shell inside the M3TAL dashboard container
sudo docker exec -it m3tal-dashboard /bin/bash

# Get a shell inside the Traefik container
sudo docker exec -it traefik /bin/sh
```

---

That's your complete reference, operator. With this knowledge, you are fully equipped to navigate and control your M3TAL ecosystem. Remember, when in doubt, consult the TUI with `sudo m3tal` or leverage `m3tal doctor` for diagnostics. Happy self-hosting!