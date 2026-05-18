# M3TAL CLI Command Reference

This document provides a comprehensive cheat-sheet for the M3TAL Command Line Interface (CLI). It covers all available commands and their usage with real-world examples.

---

## Table of Contents

1.  [Interactive Control Center](#interactive-control-center)
2.  [M3TAL Initialization and Configuration](#m3tal-initialization-and-configuration)
    *   [Initialization](#initialization)
    *   [Configuration Wizard](#configuration-wizard)
    *   [Setting Environment Variables](#setting-environment-variables)
    *   [Getting Environment Variables](#getting-environment-variables)
    *   [Scanning Environment Variables](#scanning-environment-variables)
    *   [Listing Current Environment Variables](#listing-current-environment-variables)
3.  [Dashboard Management](#dashboard-management)
    *   [Updating Dashboard Password](#updating-dashboard-password)
    *   [Starting the Dashboard](#starting-the-dashboard)
    *   [Stopping the Dashboard](#stopping-the-dashboard)
    *   [Restarting the Dashboard](#restarting-the-dashboard)
    *   [Viewing Dashboard Logs](#viewing-dashboard-logs)
    *   [Checking Dashboard Status](#checking-dashboard-status)
4.  [System-Wide Operations](#system-wide-operations)
    *   [Starting All Stacks](#starting-all-stacks)
    *   [Stopping All Stacks](#stopping-all-stacks)
    *   [Streaming Aggregated Logs](#streaming-aggregated-logs)
5.  [Systemd Service Management](#systemd-service-management)
6.  [Direct Docker Compose Usage (Fallback)](#direct-docker-compose-usage-fallback)

---

## 1. Interactive Control Center

The primary entry point for interacting with M3TAL is its interactive TUI Control Center.

| Command           | Description                                    | Example Usage                                   |
| :---------------- | :--------------------------------------------- | :---------------------------------------------- |
| `sudo m3tal`      | Opens the interactive TUI Control Center with a numbered menu for various operations. | `sudo m3tal`                                    |

---

## 2. M3TAL Initialization and Configuration

This section covers commands related to setting up and managing M3TAL's configuration.

### Initialization

| Command     | Description                                                              | Example Usage                               |
| :---------- | :----------------------------------------------------------------------- | :------------------------------------------ |
| `m3tal init`  | Generates the default `/etc/m3tal/.env` configuration file. Use on first install. | `m3tal init`                                |

### Configuration Wizard

| Command              | Description                                     | Example Usage                   |
| :------------------- | :---------------------------------------------- | :------------------------------ |
| `m3tal config wizard`  | Initiates an interactive wizard to configure `/etc/m3tal/.env`. | `m3tal config wizard`           |

### Setting Environment Variables

| Command                | Description                               | Example Usage                         |
| :--------------------- | :---------------------------------------- | :------------------------------------ |
| `m3tal config set KEY VALUE` | Sets a single environment variable in `.env`. | `m3tal config set DASHBOARD_PORT 8083` |

### Getting Environment Variables

| Command              | Description                             | Example Usage                    |
| :------------------- | :-------------------------------------- | :------------------------------- |
| `m3tal config get KEY` | Reads and displays a single env var.    | `m3tal config get DASHBOARD_PORT` |

### Scanning Environment Variables

| Command            | Description                                 | Example Usage             |
| :----------------- | :------------------------------------------ | :------------------------ |
| `m3tal config scan` | Lists all environment variables across all stacks. | `m3tal config scan`       |

### Listing Current Environment Variables

| Command            | Description                           | Example Usage             |
| :----------------- | :------------------------------------ | :------------------------ |
| `m3tal config list` | Lists the contents of the current `.env` file. | `m3tal config list`       |

---

## 3. Dashboard Management

Commands specific to managing the M3TAL dashboard container.

### Updating Dashboard Password

| Command                       | Description                                                              | Example Usage                                 |
| :---------------------------- | :----------------------------------------------------------------------- | :-------------------------------------------- |
| `m3tal dashpass [username] [password]` | Updates the dashboard user password. If arguments are omitted, it becomes interactive. | `m3tal dashpass admin mynewsecurepassword`    |
|                               |                                                                          | `m3tal dashpass` (interactive prompt)         |

### Starting the Dashboard

| Command         | Description                                                                                                                                                                                | Example Usage   |
| :-------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------- |
| `m3tal dash up` | Pulls the latest dashboard compose configuration from GitHub, reads `DASHBOARD_EXPOSE_MODE` from `/etc/m3tal/.env`, and starts the dashboard container with the appropriate override. | `m3tal dash up` |

### Stopping the Dashboard

| Command          | Description                 | Example Usage       |
| :--------------- | :-------------------------- | :------------------ |
| `m3tal dash down` | Stops the dashboard container. | `m3tal dash down`   |

### Restarting the Dashboard

| Command             | Description                   | Example Usage         |
| :------------------ | :---------------------------- | :-------------------- |
| `m3tal dash restart` | Restarts the dashboard container. | `m3tal dash restart` |

### Viewing Dashboard Logs

| Command         | Description                     | Example Usage       |
| :-------------- | :------------------------------ | :------------------ |
| `m3tal dash logs` | Streams the dashboard container logs. | `m3tal dash logs`   |

### Checking Dashboard Status

| Command           | Description                    | Example Usage         |
| :------------------ | :----------------------------- | :-------------------- |
| `m3tal dash status` | Shows the dashboard container status. | `m3tal dash status` |

---

## 4. System-Wide Operations

Commands that affect all running M3TAL stacks and services.

### Starting All Stacks

| Command   | Description                                                      | Example Usage |
| :-------- | :--------------------------------------------------------------- | :------------ |
| `m3tal up`  | Runs `docker compose up -d` across all `*-compose.yml` files in `/docker/`. | `m3tal up`    |

### Stopping All Stacks

| Command     | Description                                            | Example Usage   |
| :---------- | :----------------------------------------------------- | :-------------- |
| `m3tal down`  | Runs `docker compose down` across all stacks in `/docker/`. | `m3tal down`    |

### Streaming Aggregated Logs

| Command    | Description                               | Example Usage   |
| :--------- | :---------------------------------------- | :-------------- |
| `m3tal logs` | Streams aggregated logs from all running stacks. | `m3tal logs`    |

---

## 5. Systemd Service Management

The M3TAL API daemon is managed by `systemd`.

| Command                            | Description                                 | Example Usage                      |
| :--------------------------------- | :------------------------------------------ | :--------------------------------- |
| `systemctl status m3tal-api`       | Checks the status of the `m3tal-api` service. | `systemctl status m3tal-api`       |
| `systemctl restart m3tal-api`      | Restarts the `m3tal-api` service.           | `systemctl restart m3tal-api`      |
| `journalctl -u m3tal-api -f`       | Streams logs for the `m3tal-api` service.   | `journalctl -u m3tal-api -f`       |

---

## 6. Direct Docker Compose Usage (Fallback)

In cases where direct control is needed, you can use Docker Compose commands manually. M3TAL orchestrates these commands across files located in `/docker/`.

**Note:** These commands assume you are in the `/docker/` directory or have correctly set up your Docker Compose context.

| Command                                     | Description                                                                                                  | Example Usage                                                                   |
| :------------------------------------------ | :----------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------ |
| `docker compose -f /docker/my-stack.yml up` | Starts a specific stack defined by `/docker/my-stack.yml`.                                                   | `docker compose -f /docker/plex-compose.yml up -d`                              |
| `docker compose -f /docker/my-stack.yml down` | Stops a specific stack defined by `/docker/my-stack.yml`.                                                    | `docker compose -f /docker/plex-compose.yml down`                               |
| `docker compose -f /docker/routing-compose.yml ps` | Lists the containers for the routing stack (Traefik, Cloudflared).                                           | `docker compose -f /docker/routing-compose.yml ps`                              |
| `docker compose -f /docker/m3tal-compose.yml logs m3tal-dashboard` | Displays logs for the `m3tal-dashboard` container specifically.                                              | `docker compose -f /docker/m3tal-compose.yml logs m3tal-dashboard -f`           |
| `docker compose -f /docker/routing-compose.yml pull` | Pulls the latest images for services defined in `/docker/routing-compose.yml`.                               | `docker compose -f /docker/routing-compose.yml pull`                            |

---
```