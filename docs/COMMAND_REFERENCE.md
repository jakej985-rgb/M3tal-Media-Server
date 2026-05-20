As DocSmith, the M3TAL Ecosystem Documentation Architect, I present the definitive CLI cheat-sheet for mastering your M3TAL deployment. This guide covers every essential command, demystifies the system's architecture, and provides crucial insights into its operation.

---

# docs/COMMAND_REFERENCE.md — M3TAL CLI Cheat-Sheet

Welcome to the M3TAL Ecosystem. This document serves as your comprehensive guide to interacting with your M3TAL deployment via the command-line interface. From initial setup to day-to-day operations and advanced troubleshooting, you'll find everything you need to manage your services efficiently.

## M3TAL CLI Commands

The `m3tal` binary is your primary interface for managing the M3TAL ecosystem. Below is a complete reference of its commands.

### Core System Management

#### `sudo m3tal`
Opens the interactive Text User Interface (TUI) Control Center. This provides a numbered menu for common operations and system status.
**Usage Example:**
```bash
sudo m3tal
```

#### `m3tal init`
Generates the primary configuration file, `/etc/m3tal/.env`, from default templates. This command should be used on first installation or to reset the configuration.
**Usage Example:**
```bash
m3tal init
```

#### `m3tal doctor`
Performs a pre-flight health check of your M3TAL system. This includes verifying Docker connectivity, validating the `/etc/m3tal/.env` file, and checking for essential port availability. Crucial for troubleshooting.
**Usage Example:**
```bash
m3tal doctor
```

### Configuration Management

M3TAL uses `/etc/m3tal/.env` as its central configuration file. These commands facilitate its management.

#### `m3tal config wizard`
Launches an interactive wizard to guide you through configuring or updating the `/etc/m3tal/.env` file. Highly recommended for initial setup and major configuration changes.
**Usage Example:**
```bash
m3tal config wizard
```

#### `m3tal config set KEY VALUE`
Sets a specific environment variable (`KEY`) to a desired `VALUE` in `/etc/m3tal/.env`. This is useful for precise, single-variable adjustments.
**Usage Example:**
```bash
m3tal config set DOMAIN myhome.com
```

#### `m3tal config get KEY`
Retrieves and displays the current value of a specified environment variable (`KEY`) from `/etc/m3tal/.env`.
**Usage Example:**
```bash
m3tal config get DASHBOARD_PORT
```

#### `m3tal config scan`
Scans and lists all recognized environment variables across all M3TAL stacks, showing their current values. This helps identify all configurable parameters relevant to your deployment.
**Usage Example:**
```bash
m3tal config scan
```

#### `m3tal config list`
Displays the entire contents of the current `/etc/m3tal/.env` file.
**Usage Example:**
```bash
m3tal config list
```

### M3TAL Dashboard Management

These commands are dedicated to controlling and interacting with the M3TAL Dashboard container (`m3tal-dashboard`).

#### `m3tal dashpass [username] [password]`
Updates the password for a specified dashboard user. If `username` and `password` are omitted, the command will prompt you interactively. Dashboard credentials are stored in `/docker/users.json`.
**Usage Examples:**
```bash
# Interactive mode
m3tal dashpass

# Direct mode for user 'admin'
m3tal dashpass admin newstrongpassword123
```

#### `m3tal dash up`
Pulls the latest dashboard Docker Compose configuration files (`m3tal-compose.yml`, `m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from GitHub, then starts the M3TAL dashboard container using the appropriate configuration based on `DASHBOARD_EXPOSE_MODE` in `/etc/m3tal/.env`.
**Usage Example:**
```bash
m3tal dash up
```

#### `m3tal dash down`
Stops and removes the M3TAL dashboard container.
**Usage Example:**
```bash
m3tal dash down
```

#### `m3tal dash restart`
Restarts the M3TAL dashboard container.
**Usage Example:**
```bash
m3tal dash restart
```

#### `m3tal dash logs`
Streams the real-time logs from the M3TAL dashboard container, useful for debugging.
**Usage Example:**
```bash
m3tal dash logs
```

#### `m3tal dash status`
Shows the current status of the M3TAL dashboard container (e.g., running, stopped, unhealthy).
**Usage Example:**
```bash
m3tal dash status
```

### Stack Management (Docker Compose)

These commands manage all Docker Compose stacks defined by `*-compose.yml` files in the `/docker/` directory.

#### `m3tal up`
Runs `docker compose up -d` across all `*-compose.yml` files found in `/docker/`. This command starts or recreates all defined services in the background.
**Usage Example:**
```bash
m3tal up
```

#### `m3tal down`
Runs `docker compose down` across all `*-compose.yml` files in `/docker/`. This stops and removes all containers, networks, and volumes defined by these compose files.
**Usage Example:**
```bash
m3tal down
```

#### `m3tal logs`
Streams aggregated real-time logs from all running Docker containers managed by M3TAL. This provides a unified view for monitoring all services.
**Usage Example:**
```bash
m3tal logs
```

---

## M3TAL System Architecture Overview

M3TAL is a robust, modular system designed for self-hosted infrastructure. Its core components work in concert:

*   **CLI binary** (`/usr/bin/m3tal`): Your primary command-line interface, a unified Go binary for all operations.
*   **API daemon** (`m3tal-api.service`): A Go binary running as a systemd service on port `8080`. It manages Docker interactions, maintains a local SQLite state database (`/var/lib/m3tal/state.db`), and exposes API routes.
*   **Dashboard container** (`m3tal-dashboard`): A Python/Flask container (internal port `8082`) that provides a user-friendly web interface. It communicates with the API daemon at `http://host.docker.internal:8080`.
*   **Traefik gateway** (`routing-compose.yml`): A reverse proxy container that exposes services (including the API and Dashboard, if configured) by domain name on host port `80`. It uses a file provider for dynamic routing configuration.
*   **Cloudflared** (`routing-compose.yml`): An optional Cloudflare tunnel container, offering zero-configuration internet access for your services via Cloudflare Tunnels.

---

## Filesystem Contract

M3TAL adheres to a strict filesystem contract for its operations. Understanding these paths is critical for configuration and troubleshooting:

| Path                        | Purpose                                                              |
| :-------------------------- | :------------------------------------------------------------------- |
| `/etc/m3tal/.env`           | Primary configuration file. Managed by `m3tal config wizard`.        |
| `/var/lib/m3tal/state.db`   | SQLite state database. Auto-created and managed by the API daemon.   |
| `/opt/m3tal/stack/`         | Canonical directory for all Docker Compose files and Traefik config. |
| `/docker`                   | Symlink to `/opt/m3tal/stack/`. This is the user-facing path.        |
| `/docker/users.json`        | Dashboard credential store. Managed by `m3tal dashpass`.             |
| `/docker/dynamic/`          | Directory for Traefik's dynamic file provider configuration.         |

---

## Dashboard Access Modes

The M3TAL Dashboard can be accessed in two distinct modes, controlled by the `DASHBOARD_EXPOSE_MODE` variable in `/etc/m3tal/.env`.

### Mode 1: `local` (Default)
*   **`DASHBOARD_EXPOSE_MODE=local`**
*   Uses the `m3tal-compose.local.yml` Docker Compose override file.
*   **Direct Port Binding:** Adds a direct port binding, typically `${DASHBOARD_PORT:-8082}:8082`, exposing the dashboard directly on the host's IP address.
*   **Access via:** `http://HOST_IP:8082` or `http://localhost:8082` (where `HOST_IP` is the IP address of your M3TAL server).
*   **Requirements:** No Traefik required. Works out-of-the-box.
*   **Best for:** LAN-only setups, first-time users, local testing, or scenarios where a reverse proxy isn't desired.

### Mode 2: `traefik`
*   **`DASHBOARD_EXPOSE_MODE=traefik`**
*   Uses the `m3tal-compose.traefik.yml` Docker Compose override file.
*   **Traefik Integration:** Adds specific Traefik labels to the dashboard container, allowing Traefik to automatically route `dash.${DOMAIN}` (e.g., `dash.myhome.com`) to the dashboard's internal port `8082`.
*   **Access via:** `http://dash.DOMAIN` (e.g., `http://dash.myhome.com`).
*   **Requirements:** Traefik must be running via `m3tal up` for this mode to function.
*   **Best for:** Domain-based setups, deployments with multiple services behind a central reverse proxy, and internet-facing access (especially when combined with Cloudflared).

**To switch modes:**
1.  Run `m3tal config set DASHBOARD_EXPOSE_MODE [local|traefik]`
2.  Run `m3tal dash up` to apply the changes and restart the dashboard.

---

## Docker / Compose Runtime Explained

M3TAL leverages **Docker Engine** and **Docker Compose V2** as fundamental dependencies for container orchestration.

*   **`m3tal up` Mechanism:** When you execute `m3tal up`, the CLI iterates through all files ending in `-compose.yml` within the `/docker/` directory and executes `docker compose up -d` for each of them. This ensures all your configured services are started or updated.
*   **`m3tal dash up` Specifics:** This command has a specialized workflow:
    1.  It first downloads the latest `m3tal-compose.yml` (base configuration) and its override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`) from the official GitHub repository.
    2.  It then reads the `DASHBOARD_EXPOSE_MODE` variable from `/etc/m3tal/.env`.
    3.  Finally, it starts the dashboard container using `docker compose` with the appropriate override file applied (`-f m3tal-compose.yml -f m3tal-compose.local.yml` or `-f m3tal-compose.yml -f m3tal-compose.traefik.yml`).
*   **User Stacks:** To add new services or applications to your M3TAL ecosystem, simply place your Docker Compose file (e.g., `my-app-compose.yml`) into the `/docker/` directory. Ensure any required environment variables are set in `/etc/m3tal/.env` (use `m3tal config wizard` or `m3tal config set`). After placing the file, run `m3tal up` to deploy your new stack.

---

## Traefik Routing Architecture

Traefik acts as the central ingress for your M3TAL services, routing incoming HTTP requests to the correct containers.

*   **Deployment:** Traefik is deployed as a Docker container via `routing-compose.yml` (part of the core M3TAL stack).
*   **Entry Point:** It binds to port `80` on the host machine, serving as the HTTP entry point for all incoming web traffic.
*   **Service Discovery:** Traefik automatically discovers services by reading Docker labels applied to containers. This enables dynamic routing without manual configuration.
*   **File Provider:** It also uses a file provider to load dynamic routing configurations from `/docker/dynamic/`. This allows for highly flexible and hot-reloading routing rules.
*   **API Routing:** An example of file-based routing is how Traefik routes `api.DOMAIN` to the M3TAL API daemon (running directly on the host at `http://host.docker.internal:8080`) via `/docker/dynamic/api.yml`.
*   **Dashboard Routing:** If `DASHBOARD_EXPOSE_MODE=traefik` is set, Traefik routes `dash.DOMAIN` to the M3TAL Dashboard container using labels defined in `m3tal-compose.traefik.yml`.

**Traefik Static Configuration (`traefik.yml`):**
```yaml
entryPoints:
  web:
    address: ":80"

providers:
  docker:
    exposedByDefault: false
    network: proxy # M3TAL uses a common 'proxy' network for Traefik to discover services
  file:
    directory: /etc/traefik/dynamic # Traefik's internal path for dynamic config, maps to /docker/dynamic on host
    watch: true
```

**Dynamic Routing Example (`/docker/dynamic/api.yml`):**
```yaml
http:
  routers:
    api:
      rule: "Host(`api.${DOMAIN}`)" # Routes traffic for api.yourdomain.com
      service: api
      entryPoints:
        - web

  services:
    api:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:8080" # Routes to the M3TAL API daemon
```

---

## Systemd Service Management

The M3TAL API daemon, a crucial backend component, runs as a `systemd` service named `m3tal-api.service`. This ensures it starts on boot and is managed consistently by the operating system.

**Common `systemctl` Commands:**

*   **Check Status:** View the current status of the API daemon.
    ```bash
    systemctl status m3tal-api
    ```
*   **Restart Service:** Restart the API daemon.
    ```bash
    sudo systemctl restart m3tal-api
    ```
*   **Enable/Disable on Boot:**
    ```bash
    sudo systemctl enable m3tal-api   # Enable auto-start on boot
    sudo systemctl disable m3tal-api  # Disable auto-start on boot
    ```
*   **Stream Logs:** View real-time logs from the API daemon for debugging.
    ```bash
    journalctl -u m3tal-api -f
    ```

---

## Port Map

Understanding the port assignments is crucial for network configuration and firewall rules.

| Port | Service                               | Access                                      |
| :--- | :------------------------------------ | :------------------------------------------ |
| 80   | Traefik HTTP entry point              | Public (if `DASHBOARD_EXPOSE_MODE=traefik`) |
| 8080 | M3TAL API daemon (Go)                 | Host-local only                             |
| 8081 | Traefik dashboard (admin interface)   | Host-local only                             |
| 8082 | M3TAL Dashboard (Python/Flask)        | Direct port (if `local` mode) or via Traefik (if `traefik` mode) |

---

## APT Installation

To install M3TAL on Debian/Ubuntu-based systems, follow these steps. This process ensures you receive automatic updates.

```bash
# 1. Add the GPG signing key for package verification
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the M3TAL APT repository to your sources list
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update your package lists and install M3TAL
sudo apt update && sudo apt install -y m3tal
```

---

## Direct Docker Compose Fallback

While `m3tal` CLI commands are the recommended way to manage your stacks, understanding direct `docker compose` commands can be invaluable for advanced debugging or in scenarios where the `m3tal` binary might not be fully functional.

**M3TAL's `m3tal up` command essentially runs:**

```bash
cd /docker/
docker compose -f routing-compose.yml up -d
docker compose -f m3tal-compose.yml -f m3tal-compose.<mode>.yml up -d # where <mode> is local or traefik
# ... and so on for any other *-compose.yml files
```

You can replicate this behavior manually:

1.  **Navigate to the stack directory:**
    ```bash
    cd /docker/
    ```

2.  **Start individual stacks:**
    *   **Routing (Traefik/Cloudflared):**
        ```bash
        docker compose -f routing-compose.yml up -d
        ```
    *   **M3TAL Dashboard (Local Mode):**
        ```bash
        docker compose -f m3tal-compose.yml -f m3tal-compose.local.yml up -d
        ```
    *   **M3TAL Dashboard (Traefik Mode):**
        ```bash
        docker compose -f m3tal-compose.yml -f m3tal-compose.traefik.yml up -d
        ```
    *   **Your Custom Stack (e.g., `my-stack-compose.yml`):**
        ```bash
        docker compose -f my-stack-compose.yml up -d
        ```

3.  **Stop individual stacks:** Replace `up -d` with `down`.
    *   **Example: Stop Routing Stack**
        ```bash
        docker compose -f routing-compose.yml down
        ```

4.  **View logs for specific containers:**
    ```bash
    docker logs -f m3tal-dashboard
    docker logs -f traefik
    ```