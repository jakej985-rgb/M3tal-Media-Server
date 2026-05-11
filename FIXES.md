### **DocCritic Audit Report: M3TAL Control Plane**

**Verdict: BLOCKED.**
As a new user, I cannot deploy this software. The documentation is currently an "architectural overview" rather than an actionable deployment guide. It assumes I am a developer familiar with the repository structure rather than a user trying to install a platform.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing `.env` Template/Guide**: You state a `.env` is required, but provide no example file or documentation on which variables are mandatory (API keys, ports, paths, etc.). Without this, the Go-backend/Docker stack will fail immediately.
2.  **Infrastructure Initialization Failure**: The "Infrastructure Requirements" section requires `/mnt` paths, but `install.py` does not appear to automate this. Expecting a user to manually `mkdir` and `chown` specific paths is a common point of failure.
3.  **Ambiguous Deployment Path**: The README mentions `source/m3tal-stack`, but it is never explained how to link that folder to the `./m3tal` CLI. Does `m3tal up` look for a specific file? Do I need to copy files?
4.  **No Gateway/Port documentation**: There is no mention of the Traefik gateway or the ports required to access the Dashboard. If a user runs the stack, they are blind as to which URL/Port to visit.

#### **WARNING**
5.  **Confusing CLI/Installer Overlap**: The README suggests running `install.py` *and* `go build` to launch the system. It is unclear if the installer sets up the environment or if it is a legacy file. 
6.  **Dependency Hell**: You list `Python` as a prerequisite for a "Go-Native" project, yet the `README` doesn't explain if the Python requirements (Flask, etc.) need a `venv` or `pip install -r requirements.txt`.

#### **SUGGESTION**
7.  **Pre-flight Check**: Add a command (e.g., `./m3tal check`) that verifies if the Docker daemon is running and if the required `/mnt` directories exist.
8.  **Formatting**: The README relies on "Architectural Blueprint" but lacks a "Quickstart" section. Move the "Installation" steps to the top.

---

### **Suggested Fixes**

*   **For Issue 1 (.env):** Create a file named `.env.example` in the root. Update the README to say: `cp .env.example .env && nano .env`.
*   **For Issue 2 (Paths):** Update `install.py` to check for these paths automatically. If missing, the script should prompt: *"M3TAL requires /mnt storage. Should I create these for you? [Y/n]"*.
*   **For Issue 3 (Deployment):** Add a section under "Launch": 
    > "The `m3tal` CLI automatically references the stack definitions in `source/m3tal-stack`. Ensure your `.env` contains the `COMPOSE_PROJECT_NAME=m3tal` variable."
*   **For Issue 4 (Access):** Explicitly state: 
    > "Access the M3TAL Dashboard at `http://<host-ip>:8080`. Ensure port 8080 is open on your firewall."
*   **For Issue 6 (Dependencies):** If the project is Go-native, deprecate the need for Python entirely if possible. If not, include a `pip install -r requirements.txt` step in the installation section. 
*   **General Structure:**
    1.  Prerequisites.
    2.  Quickstart (Clone -> Env -> Install -> Run).
    3.  Detailed Config (Storage/Networking).
    4.  Architecture (Move the "Blueprints" here).

**DocCritic Note:** *An orchestrator is only as good as the first 5 minutes of a user's experience. Currently, this platform results in an "Exit 1" before the user even types the first command.*