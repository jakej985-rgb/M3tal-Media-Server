## **DocCritic Audit Report: M3TAL Core Orchestrator**

**To:** DocSmith / M3TAL Architectural Team
**From:** DocCritic, Senior DevOps Auditor
**Subject:** Documentation Audit – M3TAL Media Server

---

### **Verdict: REJECTED**
The current documentation is an architectural vision document, not an operational manual. As a new user, I am left with a pile of assumptions, missing environment requirements, and no clear path to verify the deployment. You are assuming the user has a perfectly configured system environment that matches your internal lab; real-world users will fail at step 3.

---

### **Detailed Issue List**

#### **BLOCKER**
1. **Missing `.env` Configuration:** The documentation references `/etc/m3tal/.env` as the "Source of Truth" but provides no template, no variable list (e.g., API keys, database credentials, volume paths), and no instructions on how to generate it.
2. **Missing Filesystem Pre-requisites:** You reference `/mnt/m3tal-media` and `/opt/m3tal/stack`. If these directories do not exist on a clean install, does the `m3tal init` command create them with the correct permissions (chown/chmod)? If not, the deployment will crash on permission errors.
3. **Traefik Configuration Silence:** You claim "Traefik ownership" but provide zero instructions on how to expose the services. Which ports need to be open on the host firewall? Where is the Traefik entry point configuration?

#### **WARNING**
4. **"Magic" APT Repository:** The `curl` command to add the repository assumes the user is root or has global sudo privileges. It fails to mention that `m3tal` binary might need `CAP_NET_BIND_SERVICE` or Docker socket access.
5. **Lack of "Doctor" Logic:** If `m3tal up` fails, the documentation provides no troubleshooting steps. 
6. **Docker-Compose Ambiguity:** You refer to `/docker` as a "User Entry Point" (symlink to `/opt/m3tal/stack`), but you never explain if the user needs to manually create this symlink or if the CLI does it.

#### **SUGGESTION**
7. **Developer vs. Operator Clarity:** The README mentions Go 1.21+ requirements. Clarify if the user *needs* to compile the code or if the APT install is sufficient. A binary-consumer doesn't need a Go toolchain.
8. **Logging/Observability:** Where do the logs live? If the dashboard fails to connect to the backend, there is no mention of `m3tal logs` or log file paths.

---

### **Suggested Fixes**

1.  **Environment Setup:** 
    *   Add a section: `Configuration Setup`. Provide a sample `.env` file structure.
    *   `m3tal init` must explicitly state that it checks for `/mnt/m3tal-media` and logs a warning/error if it is missing or unmounted.
2.  **Port Mapping Table:** Add a table to the README:
    | Service | Port | Access |
    | :--- | :--- | :--- |
    | Traefik Web | 80/443 | External |
    | M3TAL API | 8080 | Internal (Traefik) |
    | M3TAL Dash | 3000 | Internal (Traefik) |
3.  **Deployment Verification:** Add a "Verification" section after `m3tal up`:
    *   *Check container status:* `docker ps | grep m3tal`
    *   *Check logs:* `m3tal logs`
    *   *Check connectivity:* `curl -v http://localhost:8080/health`
4.  **Permission Safety:** Ensure the `post-install` script for the APT package creates the directories:
    ```bash
    sudo mkdir -p /etc/m3tal /opt/m3tal/stack /var/lib/m3tal
    sudo chown $USER:$USER /etc/m3tal /var/lib/m3tal
    ```
5.  **Standardize Pathing:** Remove the ambiguity about the `/docker` symlink. Either state "Run `ln -s /opt/m3tal/stack /docker`" or have the CLI `init` process handle it automatically.

---
*DocCritic Note: Do not push this to production until a "Troubleshooting" section is added. Currently, this documentation forces a support ticket on every user.*