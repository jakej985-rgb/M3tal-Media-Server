# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for all available M3TAL CLI commands.

## Table of Contents

1.  [Interactive Control Center](#interactive-control-center)
2.  [System Initialization and Configuration](#system-initialization-and-configuration)
3.  [Configuration Management](#configuration-management)
4.  [Dashboard Management](#dashboard-management)
5.  [System-Wide Operations](#system-wide-operations)
6.  [Systemd Service Management](#systemd-service-management)
7.  [Direct Docker Compose Fallback](#direct-docker-compose-fallback)

---

## Interactive Control Center

The primary interface for managing your M3TAL ecosystem.

| Command              | Description                                           | Example Usage                                 |
| :------------------- | :---------------------------------------------------- | :-------------------------------------------- |
| `sudo m3tal`         | Opens the interactive TUI Control Center (numbered menu). | `sudo m3tal`                                  |

---

## System Initialization and Configuration

Commands for setting up and preparing your M3TAL environment.

| Command                | Description                                                    | Example Usage                                                                                                              |
| :--------------------- | :------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------- |
| `m3tal init`           | Generates `/etc/m3tal/.env` from defaults. Use on first install. | `m3tal init`                                                                                                               |
| `m3tal doctor`         | Pre-flight health check: Docker connectivity, .env validity, port availability. | `m3tal doctor`                                                                                                             |
| `m3tal config wizard`  | Interactive wizard to configure `/etc/m3tal/.env`.             | `m3tal config wizard`                                                                                                      |

---

## Configuration Management

Commands for directly manipulating M3TAL's environment variables.

| Command                  | Description                                | Example Usage                                     |
| :----------------------- | :----------------------------------------- | :------------------------------------------------ |
| `m3tal config set KEY VALUE` | Set a single env var.                      | `m3tal config set DOMAIN mydomain.com`            |
| `m3tal config get KEY`     | Read a single env var.                     | `m3tal config get DASHBOARD_PORT`                 |
| `m3tal config scan`        | List all env vars across all stacks.       | `m3tal config scan`                               |
| `m3tal config list`        | List current `.env` file contents.         | `m3tal config list`                               |

---

## Dashboard Management

Commands specifically for the M3TAL dashboard container.

| Command                 | Description                                                                        | Example Usage                          |
| :---------------------- | :--------------------------------------------------------------------------------- | :------------------------------------- |
| `m3tal dashpass [username] [password]` | Update dashboard user password. Interactive if args omitted.       | `m3tal dashpass myuser mysecurepassword` |
| `m3tal dash up`         | Pull latest dashboard compose config from GitHub, then start the dashboard container. | `m3tal dash up`                        |
| `m3tal dash down`       | Stop the dashboard container.                                                      | `m3tal dash down`                      |
| `m3tal dash restart`    | Restart the dashboard container.                                                   | `m3tal dash restart`                   |
| `m3tal dash logs`       | Stream dashboard container logs.                                                   | `m3tal dash logs`                      |
| `m3tal dash status`     | Show dashboard container status.                                                   | `m3tal dash status`                    |

---

## System-Wide Operations

Commands that affect all running M3TAL stacks.

| Command        | Description                                                                   | Example Usage |
| :------------- | :---------------------------------------------------------------------------- | :------------ |
| `m3tal up`     | Run `docker compose up` across all `*-compose.yml` files in `/docker/`.       | `m3tal up`    |
| `m3tal down`   | Run `docker compose down` across all stacks.                                  | `m3tal down`  |
| `m3tal logs`   | Stream aggregated logs from all running stacks.                               | `m3tal logs`  |

---

## Systemd Service Management

The M3TAL API daemon is managed by `systemd`.

| Command                               | Description                                         | Example Usage                                    |
| :------------------------------------ | :-------------------------------------------------- | :----------------------------------------------- |
| `systemctl status m3tal-api`          | Check the status of the M3TAL API service.          | `systemctl status m3tal-api`                     |
| `systemctl restart m3tal-api`         | Restart the M3TAL API service.                      | `systemctl restart m3tal-api`                    |
| `journalctl -u m3tal-api -f`          | Stream logs from the M3TAL API service in real-time. | `journalctl -u m3tal-api -f`                     |

---

## Direct Docker Compose Fallback

In situations where the `m3tal` CLI might not suffice, you can directly use `docker compose` commands. Ensure you are in the `/docker/` directory.

| Command                                                                             | Description                                                                                                       | Example Usage                                                                                                                                                                                                                          |
| :---------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `docker compose up -d`                                                              | Start all services defined in `*.yml` files in the current directory (typically `/docker/`) in detached mode.     | `cd /docker && docker compose up -d`                                                                                                                                                                                                   |
| `docker compose down`                                                               | Stop and remove all services defined in `*.yml` files in the current directory.                                   | `cd /docker && docker compose down`                                                                                                                                                                                                    |
| `docker compose ps`                                                                 | List the running containers managed by `docker compose`.                                                          | `cd /docker && docker compose ps`                                                                                                                                                                                                      |
| `docker compose logs [service_name]`                                                | View logs for a specific service or all services if none is specified.                                            | `cd /docker && docker compose logs m3tal-dashboard`                                                                                                                                                                                    |
| `docker compose pull [service_name]`                                                | Download newer versions of a service's images.                                                                    | `cd /docker && docker compose pull m3tal-dashboard`                                                                                                                                                                                    |
| `docker compose build [service_name]`                                               | Build or rebuild services.                                                                                        | `cd /docker && docker compose build m3tal-dashboard`                                                                                                                                                                                   |
| `docker compose exec [service_name] [command]`                                      | Execute a command inside a running container.                                                                     | `cd /docker && docker compose exec m3tal-dashboard bash` (to get a shell inside the dashboard container)                                                                                                                                  |
| `docker compose config`                                                             | Validate and view the Compose configuration.                                                                      | `cd /docker && docker compose config`                                                                                                                                                                                                  |
| `docker compose --file m3tal-compose.yml --file m3tal-compose.local.yml up -d`      | Start services using specific compose files (e.g., for dashboard local mode).                                     | `cd /docker && docker compose --file m3tal-compose.yml --file m3tal-compose.local.yml up -d`                                                                                                                                           |
| `docker compose --file m3tal-compose.yml --file m3tal-compose.traefik.yml up -d`    | Start services using specific compose files (e.g., for dashboard traefik mode).                                   | `cd /docker && docker compose --file m3tal-compose.yml --file m3tal-compose.traefik.yml up -d`                                                                                                                                         |
| `docker compose -f /opt/m3tal/stack/routing-compose.yml up -d`                      | Start the Traefik gateway service using its specific compose file.                                                | `docker compose -f /opt/m3tal/stack/routing-compose.yml up -d`                                                                                                                                                                         |
| `docker compose -f /opt/m3tal/stack/routing-compose.yml down`                       | Stop the Traefik gateway service.                                                                                 | `docker compose -f /opt/m3tal/stack/routing-compose.yml down`                                                                                                                                                                          |

---

## APT Installation

To install or update M3TAL, use the following APT commands:

```bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Update package list and install M3TAL
sudo apt update && sudo apt install -y m3tal
```