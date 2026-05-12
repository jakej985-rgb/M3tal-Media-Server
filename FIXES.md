### **DocCritic Audit Report: M3TAL Control Plane (v1.4.0.3)**

**Verdict:** ❌ **REJECTED**
The current documentation is an "insider-only" document. It assumes the user has deep knowledge of the project's internal directory structure and hidden dependencies. A new user will fail at Step 4 due to missing runtime context and environment variables.

---

### **Detailed Issue List**

#### **BLOCKER**
1.  **Missing `m3tal.py` / Orchestrator Context:** Step 4 runs `./m3tal`. Is this the Go binary, or is it a wrapper? The `m3tal.py` mentioned in my instructions is absent, and the Go build instructions for `m3tal-api` are provided but never utilized in a startup command.
2.  **Environment Variable Omission:** The `.env.example` file structure is not documented. If `BASE_STORAGE_PATH` is required, what are the *other* required keys? The deployment will crash if a user leaves the default `.env.example` as-is.
3.  **Missing Docker-Compose Entry:** Step 4 calls `./m3tal up`, but nowhere does it explain how `m3tal` interacts with `source/m3tal-stack/docker-compose.yml`. Does the binary auto-generate a compose file? Run it directly? This is a "black box" that will result in a deployment failure.
4.  **Missing Port/Gateway Documentation:** There is zero mention of ports (e.g., 80, 443, 8080) or Traefik/Reverse Proxy integration. A user will have no idea which URL to hit to verify the deployment.

#### **WARNING**
5.  **Host Dependency Assumption:** The documentation assumes `/mnt` exists or is the correct target. On many systems (e.g., standard Debian/Ubuntu), `/mnt` is empty and requires root permissions or specific mounting.
6.  **"Go-Native" Confusion:** The stack involves Go binaries *and* a Python/Flask dashboard. The instructions do not explain how the Python dashboard dependencies (pip install) are handled or if they are containerized within the stack.

#### **SUGGESTION**
7.  **Service Connectivity Table:** The current table is useless for a user; it lists file paths, not endpoints (e.g., `http://localhost:8080`).
8.  **Setup Script Transparency:** `scripts/setup.sh` is a "blind execution." The README should document what this script modifies (e.g., does it install Docker? Does it edit fstab?).

---

### **Suggested Fixes**

*   **For Blocker 1 & 3:** Explicitly state: "The `m3tal` binary acts as a wrapper for `docker compose`. Ensure your Docker daemon is running, as the binary calls `docker compose -f source/m3tal-stack/docker-compose.yml up -d` internally."
*   **For Blocker 2:** Provide a table of required `.env` variables:
    *   `BASE_STORAGE_PATH`: (e.g., `/mnt/m3tal_data`)
    *   `API_PORT`: (e.g., `8080`)
    *   `DASHBOARD_PORT`: (e.g., `5000`)
*   **For Blocker 4:** Add a "Verification" section: "Once deployed, access your dashboard at `http://<YOUR_HOST_IP>:5000`."
*   **For Warning 5:** Add a pre-flight check: "Ensure `/mnt` is writable by your user or adjust `BASE_STORAGE_PATH` to a directory you own (e.g., `/home/user/m3tal`). Run `mkdir -p $BASE_STORAGE_PATH` before proceeding."
*   **For Warning 6:** Clarify if the Python dashboard is built via Dockerfile. If so, add `docker compose build` to the deployment steps. If not, add a `pip install -r requirements.txt` step.
*   **General:** Add a "Quick Start" diagram or a command-line tree showing exactly where `m3tal` and `m3tal-api` binaries should be placed after compilation.