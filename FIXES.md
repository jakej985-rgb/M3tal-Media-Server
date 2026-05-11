### **DocCritic Audit Report: M3TAL Control Plane**

**Verdict: FAILED (BLOCKER)**
The current documentation is an architectural overview, not a deployment guide. A new user will fail at Step 1 because the project lacks critical build, dependency installation, and environment initialization procedures. You have described *what* the project is, but not *how* to make it run.

---

### **Issue List**

#### **BLOCKER**
*   **Missing Dependency Installation:** The README assumes a `venv` exists, but there is no instruction to create one or install `requirements.txt`.
*   **Missing Build Steps (Go-Native):** The project relies on a `source/go-backend`, but there are no instructions to compile the Go binary. A user trying to run `python m3tal.py up` will likely encounter "file not found" errors because the backend is not built.
*   **Missing `.env` Template:** The docs mention a `.env` file but provide no command to create it from a template (e.g., `cp .env.example .env`).
*   **Implicit Docker Compose Dependency:** The CLI tool `m3tal.py` is expected to manage the stack, but it is unclear if the user needs to manually run `docker-compose up` or if the Python script handles it via subprocess.

#### **WARNING**
*   **Hardcoded `/mnt` Path:** Requiring absolute paths on the root filesystem (`/mnt`) is dangerous for users on macOS, Windows (WSL), or systems where `/mnt` is restricted or reserved. It assumes the user has root/sudo access to the host partition.
*   **Traefik/Gateway Omission:** You mention a "Core Orchestrator" and "Dashboard," but there is zero information on how to access these services once deployed. What is the URL? What ports are exposed?
*   **Unclear Python Version Management:** You specify Python 3.9+, but provide no `requirements.txt` location or automated setup script.

#### **SUGGESTION**
*   **Clarify CLI Capabilities:** Does `m3tal.py` detect if dependencies are missing? It should perform a "Pre-flight Check" to verify Go, Docker, and environment variables before attempting to launch.
*   **Versioning Logic:** The README states v1.4.0.3, but the repo layout seems to be in flux. Add a "Health Check" command to the CLI to verify the stack status.

---

### **Suggested Fixes**

#### **1. Immediate Deployment Procedure Update**
Replace your "Installation" section with:
```bash
# Initialize Environment
git clone ...
cd M3tal-Media-Server

# Setup Python
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Setup Backend (Go)
cd source/go-backend
go build -o m3tal-backend .
cd ../..

# Setup Configuration
cp .env.example .env
# Edit .env and set DASHBOARD_SECRET
```

#### **2. Add a Pre-flight Diagnostic Table**
Add this to the "Infrastructure Requirements" section:
| Component | Status Check |
| :--- | :--- |
| Go Compiler | `go version` |
| Docker | `docker compose version` |
| Python | `python3 --version` |
| Storage | `ls -ld /mnt/media` |

#### **3. Explicit Access Instructions**
Add a "Post-Deployment" section:
> **Accessing M3TAL**
> Once `python m3tal.py up` completes successfully:
> * **Dashboard:** `http://localhost:8080`
> * **API Endpoint:** `http://localhost:9000` (Defined in `source/go-backend`)
> * **Logs:** Monitor via `python m3tal.py logs`

#### **4. Path Configuration Flexibility**
Modify the documentation to allow for local path overrides:
* "By default, M3TAL uses `/mnt`. For non-Linux environments, please update `BASE_STORAGE_PATH` in your `.env` to a directory where you have read/write permissions."

**Final Note:** You are treating the README as a technical design document. Treat it as a **product manual**. If I have to guess how to build the Go backend, your documentation is incomplete.