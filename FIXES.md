# Audit Report: M3TAL Media Server Documentation
**Auditor:** DocCritic, Senior DevOps Auditor  
**Status:** **FAILED**

---

### Verdict
The documentation is currently a "developer-to-developer" shorthand that fails to provide a viable path for a new user. As a Senior Auditor, I cannot recommend this for production use or even local evaluation. The project relies on "hidden" state and assumes a perfectly configured host environment without documenting the prerequisites.

---

### Issue List

#### 1. BLOCKER: Missing Environment Initialization
The `Quick Start` assumes `./m3tal init` will magically work, but there is no documentation on where the initial configuration files or `.env` templates originate.
*   **Fix:** Add a step to generate a `.env.example` or run a command that prompts for necessary environment variables before `init`.

#### 2. BLOCKER: Absolute Path Pre-requisite
The documentation states: "The M3TAL ecosystem mandates that the host `BASE_STORAGE_PATH` is always mounted to `/mnt`". 
*   **Issue:** If I am on macOS or a custom Linux partition where `/mnt` is restricted or non-existent, the documentation provides no guidance on how to satisfy this constraint.
*   **Fix:** Document exactly how to configure `BASE_STORAGE_PATH` in the `.env` file and warn users that this directory *must* exist on the host prior to execution.

#### 3. WARNING: Undocumented Port Dependencies
`Traefik` is mentioned as the gateway, but no ports are specified. Is it 80/443? Does it conflict with existing local services (like an existing local web server)?
*   **Fix:** Add a "Networking Requirements" section explicitly listing the required host ports (e.g., 80, 443, 8080) and how to handle collisions.

#### 4. WARNING: Hidden `m3tal-goback` Dependency
The `Operational Flow` mentions the Dashboard communicates with an *external* `m3tal-goback` service. 
*   **Issue:** A new user will run `./m3tal up` and be greeted with a broken dashboard because they didn't know they needed to deploy or configure a separate backend.
*   **Fix:** Provide a clear "Prerequisites" section or a `docker-compose.override.yml` example to link the backend, or define the environment variable where this URL is set.

#### 5. SUGGESTION: Docker/Compose Version Requirements
There is no mention of `docker-compose` version compatibility or if `docker compose` (v2) is required.
*   **Fix:** Add a "System Requirements" section: `Docker Engine >= 20.10`, `Docker Compose >= 2.0`.

---

### Suggested Documentation Remediation (Draft)

**Revised Quick Start Section:**

1. **Install Prerequisites**: Ensure `Docker` and `Go 1.21+` are installed.
2. **Environment Setup**:
   ```bash
   cp .env.example .env
   # Edit .env and set BASE_STORAGE_PATH to an absolute path on your host (e.g., /home/user/media)
   ./m3tal config set BASE_STORAGE_PATH=/home/user/media
   ```
3. **Initialize**: 
   ```bash
   ./m3tal init 
   # This will validate your /mnt mount and verify Docker connectivity.
   ```
4. **Launch**:
   ```bash
   ./m3tal up
   # Access your dashboard at http://localhost:<DASHBOARD_PORT>
   ```

**Revised Networking Note:**
> **Note:** The stack requires ports 80 and 443. Ensure no local processes (e.g., Nginx, Apache) are bound to these ports. If conflicts occur, update the `HTTP_PORT` and `HTTPS_PORT` variables in your `.env` file and re-run `./m3tal init`.

---
**Auditor Final Note:** *The "Path Consistency Rule" is a brittle design pattern. Consider abstracting the mount path via the `m3tal-stack` docker-compose file instead of forcing the host to adapt to a hardcoded `/mnt` requirement.*