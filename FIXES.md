### **Audit Report: M3TAL Core Documentation**
**Auditor:** Senior DevOps Auditor (DocCritic)  
**Status:** **NON-COMPLIANT**

The current documentation assumes a "happy path" that does not exist for a clean-slate user. The instructions are dangerously vague regarding permissions, networking, and the actual state of the filesystem.

---

### **Verdict**
**REJECTED.** The project is currently non-deployable for a new user. The lack of environment variable documentation, directory permission requirements, and Traefik configuration details will lead to immediate runtime failures.

---

### **Issue List**

#### **BLOCKER**
*   **Missing `.env` Specification:** `m3tal init` generates secrets, but there is no documentation on *required* variables (e.g., API keys, database credentials, or domain names) that the user must set *before* running `m3tal up`.
*   **Directory Permissions:** The documentation references `/etc/m3tal/`, `/var/lib/m3tal/`, and `/opt/m3tal/`. As a standard user, `m3tal` will fail to create these or write to them without `sudo`. The installation instructions are inconsistent about user vs. root privilege.
*   **Network/Port Mapping:** There is no mention of which ports must be opened (e.g., 80/443 for Traefik). Users running this on a cloud VPS will find themselves blocked by firewalls with no guidance.

#### **WARNING**
*   **"Black Box" `m3tal up`:** The documentation hides the orchestration logic. If the `docker-compose.yml` resides in `/opt/m3tal/stack`, what does `m3tal up` actually do? Does it run `docker compose -f ... up -d`? If it fails, the user has no way to debug the underlying Docker stack.
*   **Traefik Configuration:** You state "Traefik Ownership" but provide no entrypoint or config example. A user will spin this up and have no idea how to route traffic to the dashboard.

#### **SUGGESTION**
*   **Missing "Doctor" context:** The `m3tal doctor` command is mentioned in the CLI reference but not the Quick Start. It should be the final verification step in the setup flow.
*   **Filesystem Ambiguity:** Why is there a `/docker` symlink? This is confusing for Linux admins who expect standard paths. Clarify if `/docker` is intended for user-mounted media or if it is strictly internal.

---

### **Suggested Fixes**

1.  **Add a "Configuration" Section:**
    Create an `Environment Configuration` table detailing:
    *   `M3TAL_DOMAIN`: External URL for Traefik.
    *   `M3TAL_API_KEY`: Secret for internal Auth.
    *   `M3TAL_DATA_PATH`: Override for persistent storage.

2.  **Explicit Permission Steps:**
    Add a setup block:
    ```bash
    # Ensure system directories exist and are writable
    sudo mkdir -p /etc/m3tal /var/lib/m3tal /opt/m3tal
    sudo chown $USER:$USER /etc/m3tal /var/lib/m3tal /opt/m3tal
    ```

3.  **Traefik & Ports:**
    Add a networking section:
    *   *Required Ports:* `80 (HTTP)`, `443 (HTTPS)`.
    *   Provide a sample `traefik.yml` or clarify if the `m3tal` binary autogenerates it.

4.  **Refine Quick Start:**
    Update the Quick Start sequence:
    1. Install.
    2. `m3tal init` (Explain that this creates `/etc/m3tal/.env`).
    3. **"Edit your `/etc/m3tal/.env` file before proceeding."**
    4. `m3tal up`.
    5. `m3tal doctor` (Verify status).

5.  **Troubleshooting/Logs:**
    Add a "Debugging" section explaining how to view the underlying stack status:
    *   `docker compose -f /opt/m3tal/stack/docker-compose.yml ps`

**DocCritic Note:** *Do not treat your CLI as a magical abstraction. In DevOps, if the user cannot see the `docker-compose.yml` or the environment variables, the system is impossible to audit or repair.*