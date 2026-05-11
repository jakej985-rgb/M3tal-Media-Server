To: Project Maintainers
From: DocCritic, Senior DevOps Auditor
Subject: Audit Report - M3TAL Control Plane (v1.4.0.3)

---

### **Verdict: FAILED**
**Reasoning:** The documentation is currently a "developer’s journal" rather than a deployment guide. It fails to account for critical path dependencies, lacks internal verification logic, and obscures the network topology. A new user will encounter "Permission Denied" or "Connection Refused" errors within minutes of deployment.

---

### **Issue List**

*   **[BLOCKER] Missing Initialization of `m3tal.py` requirements:** The `install.py` script exists, but the document does not explain that `m3tal.py` requires its own dependency resolution or system-level binary compilation (Go modules).
*   **[BLOCKER] Hardcoded `/mnt` Assumption:** You are forcing a root-level directory creation (`/mnt/`) on the host. This will fail on systems with restricted permissions (macOS/non-root Linux users) or where `/mnt` is a protected mount point.
*   **[BLOCKER] Traefik/Gateway Port Omission:** The documentation mentions a "Dashboard" and "Backend," but fails to mention that `m3tal-stack` likely deploys a Traefik gateway. Users need to know which port to hit to actually access the UI.
*   **[WARNING] Unclear `.env` lifecycle:** The guide says "Ensure your .env file is configured," but does not provide a `.env.example` file or instructions on *how* to generate it via the installer.
*   **[WARNING] Build Step Omission:** The project uses a "Go-native backend." There is no instruction to `go build` the binary in `source/go-backend` before running `m3tal.py up`.
*   **[SUGGESTION] Ambiguous CLI context:** `python m3tal.py` is called at the root, but the backend lives in `source/`. It is unclear if `m3tal.py` handles the compilation of Go binaries or expects them pre-compiled.

---

### **Suggested Fixes**

#### 1. Decouple from `/mnt`
Do not hardcode `/mnt`. Update `source/m3tal-stack/docker-compose.yml` to use environment variables for paths.
*   **Fix:** Add `STORAGE_ROOT=${STORAGE_ROOT:-./data}` to your `.env` template so the project runs out-of-the-box in the local folder.

#### 2. Explicit Compilation Steps
Add a section specifically for the Go backend.
*   **Fix:**
    ```bash
    cd source/go-backend
    go mod download
    go build -o m3tal-backend .
    ```

#### 3. Network/Access Table
Provide a clear "Service Mapping" table so the user knows what to open in their firewall.
*   **Example Table:**
    | Service | Port | Access |
    | :--- | :--- | :--- |
    | Dashboard | 8080 | http://host-ip:8080 |
    | Backend API | 9000 | Internal only |
    | Traefik UI | 8081 | http://host-ip:8081 |

#### 4. Automated Environment Validation
Update `install.py` to perform a pre-flight check.
*   **Fix:** Ensure `install.py` validates the existence of Docker, Go, and the presence of a generated `.env` file based on a provided `template.env`.

#### 5. Clarify `m3tal.py` Dependency
Document whether `m3tal.py` is a wrapper that executes the compiled binary.
*   **Fix:** Add a section under **CLI Commands**: *"Note: Ensure you have compiled the Go backend (see step X) before calling `m3tal.py up`, as the Python orchestrator expects the binary to be present in `source/go-backend/bin/`."*

---
**Auditor Note:** *Stop treating the deployment process as common knowledge. If it's not documented, it doesn't exist to the user.*