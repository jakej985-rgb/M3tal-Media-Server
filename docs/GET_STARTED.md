# GET STARTED: The M3TAL Journey

Welcome to the M3TAL Ecosystem. You are 10 minutes away from a fully orchestrated, Go-native infrastructure stack. 

---

## 1. Installation (The Foundation)
Deploy the M3TAL orchestrator via our APT repository to ensure seamless updates and system integration.

```bash
# Add M3TAL Repository
curl -sL https://apt.m3tal.io/setup.sh | sudo -E bash -

# Install Core Orchestrator
sudo apt update && sudo apt install m3tal -y
```

---

## 2. Configuration Wizard
Once installed, initialize your environment. The wizard will prompt you for your preferred domain, local networking, and security keys.

```bash
sudo m3tal config wizard
```

**What this does:**
*   Generates your `/etc/m3tal/.env` file.
*   Sets up the systemd daemon for `m3tal-api.service`.
*   Prepares the directory structure for your services.

---

## 3. Filesystem Contract
M3TAL respects a strict filesystem hierarchy. Understanding this is key to maintaining your stack:

| Location | Purpose |
| :--- | :--- |
| `/etc/m3tal/.env` | **Primary Configuration:** Global variables for your stack. |
| `/var/lib/m3tal/state.db` | **State Database:** The brain of the M3TAL orchestrator. |
| `/opt/m3tal/stack` | **Deployment Assets:** Compose files and Traefik configs. |
| `/docker` | **Symlink:** Convenient shortcut to `/opt/m3tal/stack`. |

---

## 4. Starting the Ecosystem
With configuration complete, fire up the orchestrator:

```bash
sudo systemctl enable --now m3tal-api
sudo m3tal up
```

Your stack is now live and managed by the Go-native orchestrator.

---

## 5. Access & Management
M3TAL exposes two primary gateways for your infrastructure:

*   **Dashboard Web Interface:** `http://localhost:8082`
*   **Traefik Gateway:** `http://localhost:8080`

### Authentication
*   **Default Credentials:** `admin / admin`
*   **Security Tip:** If you ever lose access or need to rotate credentials, run the following command from your terminal:
    ```bash
    sudo m3tal dashpass
    ```

---

## 6. Pro-Tips
*   **Logs:** Need to see what the orchestrator is doing? Run `sudo journalctl -u m3tal-api -f`.
*   **Updates:** To refresh your stack to the latest version, simply run `sudo apt update && sudo apt install m3tal` followed by `sudo m3tal up`.
*   **Safety:** Never manually edit the database in `/var/lib/m3tal/`. Always use the CLI orchestrator to ensure state consistency.

*You are now operational. Welcome to the M3TAL era.*