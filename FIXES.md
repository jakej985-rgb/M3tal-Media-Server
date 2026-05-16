# AUDIT REPORT: M3TAL Core Orchestrator Documentation
**Auditor:** DocCritic, Senior DevOps Auditor  
**Status:** **FAILED**

---

### **Verdict**
**CRITICAL FAILURE.** The current documentation is an architectural whitepaper, not an installation guide. As a new user, I have no idea how to go from "cloned repo" to "running service." You assume global system paths (`/opt`, `/etc`, `/mnt`) exist and provide no mechanism to create them. You mention a `m3tal` binary but provide no `Makefile` or build instructions.

---

### **Issue List**

#### **BLOCKER**
*   **Missing Build/Install Instructions:** The README mentions a `m3tal` binary, but there is no command to build it (e.g., `go build -o m3tal main.go`).
*   **Missing Setup Script:** The document assumes directories like `/etc/m3tal` and `/opt/m3tal` exist. There is no `setup.sh` or installation step to bootstrap these required paths/permissions.
*   **Environment Variables:** You mention `/etc/m3tal/.env` as the "Source of Truth," but provide no template or documentation for what variables are required (e.g., DB_URL, API_KEY, PORTS).
*   **Docker Orchestration Logic:** You provide a snippet for `m3tal-core`, but you do not explain how the user is supposed to deploy the *rest* of the stack (the `m3tal-stack`).

#### **WARNING**
*   **Traefik Gateway Omission:** You claim "Traefik Ownership," but provide zero guidance on how to configure the Traefik entry point, port mapping (80/443), or internal network names. A new user will have a conflict with their host port 80 immediately.
*   **Assuming /mnt exists:** You treat `/mnt/m3tal-media` as a default. If the user doesn't have a drive mounted at `/mnt`, the container mount will create an empty root-owned directory, causing permission errors.

#### **SUGGESTION**
*   **Confusing Hierarchy:** The README describes `/docker` as the "Primary User Entry Point" but then provides a Docker Compose block that relies on `/opt/m3tal`. These instructions are contradictory.

---

### **Suggested Fixes**

1.  **Add a `Quickstart` Section:**
    ```bash
    # Example addition
    git clone ...
    sudo ./scripts/install.sh # Needs to create /etc/m3tal, /opt/m3tal, etc.
    cp .env.example /etc/m3tal/.env
    go build -o m3tal .
    sudo mv m3tal /usr/bin/
    ```

2.  **Provide an `.env.example`:** Create a file in the repo detailing required variables:
    *   `M3TAL_API_PORT`
    *   `DOCKER_NETWORK_NAME`
    *   `TRAEFIK_ENTRYPOINT`

3.  **Clarify the Docker Lifecycle:** Provide a `docker-compose.yml` that includes the full stack (Traefik, Goback, Godash, and Core), rather than just the Core container. If the Core is meant to spin up other containers, provide an example of the command used to trigger this via the binary (e.g., `m3tal start --stack`).

4.  **Add a Directory Validator:** Add a check in your Go code (at startup) to verify that the required paths (`/etc/m3tal`, etc.) exist. If they don't, have the program output a clear error: `CRITICAL: Directory /etc/m3tal not found. Run ./setup.sh`.

5.  **Standardize Pathing:** Decide on one "Entry Point." If `/opt/m3tal` is the source of truth, stop referencing `/docker` as the user entry point. It creates cognitive load for no reason.

**DocCritic Note:** *Fix these before I return for the Stage 2 Compliance Audit. Documentation is code; if it doesn't run, it's a bug.*