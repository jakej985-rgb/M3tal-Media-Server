**Audit Report: M3TAL Core Orchestrator Documentation**
**Auditor:** DocCritic, Senior DevOps Auditor
**Date:** 2023-10-27
**Subject:** Documentation Audit - M3TAL Core

---

### **Verdict: FAILED**
The provided documentation is high-level architectural marketing material, not a functional installation guide. As a new user, I have no idea how to compile, install, or initialize the state required to make this platform operational. The current README is a "roadmap," not a "manual."

---

### **Issue List**

1.  **[BLOCKER] Missing Compilation/Installation Steps:** The README mentions a "Go-native binary," but provides no `go build` instructions, no binary release source, and no instructions on how to move the artifact to `/usr/bin/m3tal`.
2.  **[BLOCKER] Missing Initialization Command:** There is no documentation for a `m3tal setup` or `m3tal init` process. How are the `/etc/m3tal/.env` and `/var/lib/m3tal/` directories generated? 
3.  **[BLOCKER] Environment Variable Definition:** The `.env` file is referenced as the "Global Configuration Source of Truth," but no schema or template is provided. What keys are required (e.g., API keys, port bindings, database credentials)?
4.  **[WARNING] Infrastructure Assumptions:** The documentation assumes `/mnt/m3tal-media` exists and is mounted. If this directory is missing, will the service crash? The docs should mandate a check or provide a script to ensure directory integrity.
5.  **[WARNING] Missing Network/Gateway Exposure:** The documentation mentions a Dashboard and API, but fails to mention how to route traffic to them (e.g., Traefik/Nginx configurations, exposed port numbers).
6.  **[SUGGESTION] Docker Compose Incompleteness:** The provided snippet is a raw block. It lacks a `docker-compose.yml` file structure, versioning, and networking context needed for a multi-container ecosystem (Core/Back/Dash).

---

### **Suggested Fixes**

*   **For [1 & 2]:** Add a "Getting Started" section.
    *   *Add:* `go build -o m3tal ./cmd/m3tal && sudo mv m3tal /usr/bin/`
    *   *Add:* A `m3tal setup` command documentation that creates the directory tree: `mkdir -p /etc/m3tal /var/lib/m3tal /opt/m3tal/stack`.
*   **For [3]:** Include an `.env.example` block in the README detailing:
    ```text
    M3TAL_API_PORT=8080
    M3TAL_DASH_PORT=3000
    M3TAL_STORAGE_ROOT=/mnt/m3tal-media
    M3TAL_SECRET_KEY=...
    ```
*   **For [4]:** Add a "Pre-requisites" section advising the user to configure their mount points *before* starting the stack, and provide a shell one-liner: `[ -d /mnt/m3tal-media ] || sudo mkdir -p /mnt/m3tal-media`.
*   **For [5]:** Create a `docker-compose.yml` that defines a shared internal network (`m3tal-net`) so the dashboard and backend can communicate internally. Clearly state which ports must be opened on the host firewall.
*   **For [6]:** Provide a full `docker-compose.yml` example. The current snippet is isolated and useless for a user trying to link the `goback` and `godash` services together.

---

**Auditor's Closing Note:** *Fix these blockers immediately. A system is only as robust as the user's ability to deploy it without trial and error. Re-submit once a "Quick Start" section is added.*