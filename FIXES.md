### **DocCritic Audit Report: M3TAL Core**

**Verdict:** **FAILED**. As a new user, I cannot successfully deploy this stack. The documentation assumes a "perfect" environment and fails to provide critical implementation details regarding file system dependencies and binary initialization.

---

## Issue List

#### **BLOCKER**
*   **[BLOCKER] `/mnt` Dependency:** The docs state `/mnt` is a "Required" directory. On macOS or Windows (Docker Desktop), the `/mnt` root path does not exist and cannot be created by the user. If the Go orchestrator blindly expects `/mnt` to exist on the host, the `m3tal up` command will crash or fail to mount volumes.
*   **[BLOCKER] Missing `m3tal.py` initialization:** The architecture section mentions a Python/Flask dashboard, but there is no instruction on how to install Python dependencies (`pip install -r requirements.txt`) or if the orchestrator handles this. I don't know if I need to run a `venv` or if it's containerized.
*   **[BLOCKER] `.env` Variable Blindness:** The `cp template.env .env` command is provided, but there is no instruction to actually *edit* the file. The orchestrator will likely fail if `BASE_STORAGE_PATH` points to a non-existent or default location.

#### **WARNING**
*   **[WARNING] Port Conflict Potential:** The "Pre-flight Checklist" notes ports 80/443 must be free. On most Linux distributions, these are protected ports requiring `sudo` or have existing services (Apache/Nginx). The documentation ignores the `sudo` requirement or capability mapping.
*   **[WARNING] Binary Compilation:** `build.sh` is invoked, but there is no mention of Go dependency management (`go mod download`). If the environment isn't pre-warmed, the build will fail.

#### **SUGGESTION**
*   **[SUGGESTION] Traefik Gateway Config:** There is no mention of how the Traefik configuration in `source/m3tal-stack/` expects the network to be defined. If a user tries to run this on a host with an existing Traefik instance, it will collide.
*   **[SUGGESTION] CLI "Help" command:** The CLI reference is good, but `m3tal` is a new tool. Please add `./m3tal --help` or `--version` verification steps for initial smoke testing.

---

### **Suggested Fixes**

1.  **Resolve `/mnt` Issue:** 
    *   **Fix:** Update the documentation to allow a custom `BASE_STORAGE_PATH` that defaults to a local subdirectory (e.g., `./data`) rather than the system root `/mnt`. 
    *   *Documentation:* "Ensure `BASE_STORAGE_PATH` is set to an absolute path. If on macOS/Windows, avoid root-level paths like `/mnt` and use your user home directory."

2.  **Explicit `.env` Configuration:**
    *   **Fix:** Change the quickstart to:
        ```bash
        cp template.env .env
        nano .env  # Update BASE_STORAGE_PATH and any API keys
        ```

3.  **Dependency Initialization:**
    *   **Fix:** Add a dedicated "System Preparation" step:
        ```bash
        # Initialize Go dependencies
        go mod tidy
        # If the dashboard is not containerized, instructions for the Python environment:
        cd source/dashboard && pip install -r requirements.txt
        ```

4.  **Networking/Privilege Warning:**
    *   **Fix:** Add a "Pro-tip" to the Networking section: "If you cannot bind to ports 80/443, update `docker-compose.yml` to map these to high-range ports (e.g., 8080:80)."

5.  **Build Safety:**
    *   **Fix:** Ensure `build.sh` performs a `go mod vendor` or `go mod download` check to prevent confusing "module not found" errors during compilation.