# DocCritic Audit Report: M3TAL Core Orchestrator

**Auditor:** DocCritic, Senior DevOps Auditor  
**Verdict:** **REJECTED - REQUIRES REWORK**

While the documentation provides a baseline for installation, it suffers from "Developer Tunnel Vision." It assumes the user has a pre-configured filesystem, lacks critical networking information, and relies on marketing fluff rather than technical specification. 

---

### Issue List

#### 1. BLOCKER: Missing Network/Access Requirements
The documentation mentions a Dashboard and API, but fails to define **ports**. Users cannot access the services if they do not know which host ports to open or which Traefik/Reverse Proxy configurations are expected.
*   **Fix:** Add a "Network Requirements" section listing standard ports (e.g., 8080, 9000) and indicate if a reverse proxy is required for the dashboard.

#### 2. BLOCKER: Dev-Only Assumption (`/mnt` existence)
The documentation mandates `/mnt/m3tal-media` but provides no instruction on how to provision this. A fresh Debian install does not automatically mount external drives to this path. 
*   **Fix:** Add a step for directory initialization (`sudo mkdir -p /mnt/m3tal-media && sudo chown $USER:$USER /mnt/m3tal-media`) or explain that this is a placeholder for a mount point.

#### 3. WARNING: Ambiguous "Docker Deployment"
The provided `docker-compose.yaml` snippet is an orphan. It is unclear if the user is supposed to create this file, where to save it, or if `m3tal up` generates it automatically.
*   **Fix:** Explicitly state: "The `m3tal up` command automatically triggers a `docker compose -f /opt/m3tal/stack/docker-compose.yml up -d` operation. Ensure your stack file is located at [path]."

#### 4. WARNING: Missing Lifecycle Management
There is no information on how to update, stop, or remove the services. "Up" is only half the lifecycle.
*   **Fix:** Add a section for `m3tal down` and `m3tal status` or `m3tal logs`.

#### 5. SUGGESTION: "Marketing vs. Manual"
The footer ("Modular Infrastructure Platform. Status: Go-Native Migration Active") is fluff. Documentation should be terse and functional.
*   **Fix:** Remove marketing slogans. Replace with a "System Requirements" table (CPU/RAM/Kernel version).

#### 6. SUGGESTION: Clarification on `m3tal setup`
The command `m3tal setup` is mentioned, but the user is not warned about what it does (e.g., does it overwrite existing configs? Does it require root?).
*   **Fix:** Add a disclaimer: "Warning: `m3tal setup` creates `/etc/m3tal/config.yaml`. If this file exists, it will not be overwritten by default."

---

### Suggested README Refinements

**Replace the "Deployment" section with this:**

#### Deployment & Networking
The Orchestrator manages the container lifecycle. Ensure your system meets the following prerequisites:

1.  **Storage:** Create the media mount point:
    ```bash
    sudo mkdir -p /mnt/m3tal-media
    sudo chown -R $USER:$USER /mnt/m3tal-media
    ```
2.  **Ports:** Ensure the following ports are available on the host:
    *   `8080`: Dashboard UI
    *   `9000`: Backend API
    *   `9001`: Metrics/Health-check

**Replace "Ecosystem Integration" with a "System Lifecycle" section:**
*   `m3tal up`: Pulls images and initializes the stack defined in `/opt/m3tal/stack/`.
*   `m3tal down`: Safely shuts down containers.
*   `m3tal status`: Returns a JSON object containing container health and API connectivity status.

---

**DocCritic Final Note:** *Stop treating your users like they already know how you built the tool. Define your expectations for the host environment explicitly. If the system crashes because `/mnt` doesn't exist, that's a failure of the documentation, not the user.*