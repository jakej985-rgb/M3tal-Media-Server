**DOC-CRITIC AUDIT REPORT v2024.1**
**Auditor:** DocCritic (Senior DevOps Auditor)
**Project:** M3TAL Media Server (v1.4)
**Status:** **REJECTED**

---

### ⚖️ VERDICT
**The documentation is currently UNUSABLE for a new user.** While the architectural overview is high-quality, the installation path contains a "hallucinated" prerequisite and critical gaps regarding repository structure and environment initialization. A user following these steps verbatim will experience a build failure within 60 seconds.

---

### 🚩 DETAILED ISSUE LIST

#### 1. [BLOCKER] Impossible Prerequisite: Go 1.26+
The documentation requires **Go 1.26+**.
*   **The Issue:** As of today, the current stable version of Go is 1.22. Go 1.26 does not exist. A user attempting to install "the latest" will find 1.22, assume their environment is "too old," and stop.
*   **Suggested Fix:** Update the prerequisite to a realistic version (e.g., `Go 1.21+` or `Go 1.22+`).

#### 2. [BLOCKER] The "Empty Shell" Problem (Source Management)
The README mentions building binaries in `source/` and using `source/m3tal-stack`.
*   **The Issue:** It is unclear if these directories are included in the main repository, managed via Git Submodules, or if the user needs to clone the "Related Projects" manually into the `source/` folder. If I clone just this repo, `build.sh` will likely fail because `source/` is empty or missing.
*   **Suggested Fix:** Explicitly state: "Clone with submodules: `git clone --recursive ...`" OR provide a setup script that fetches the dependencies.

#### 3. [BLOCKER] Binary Location Ambiguity
The Quick Start says to run `./m3tal init`.
*   **The Issue:** The `build.sh` script "compiles Go-native binaries found in `source/`." It does not state that it moves the resulting orchestrator binary to the root directory. If `go build` outputs to `./source/m3tal/m3tal`, the Quick Start command `./m3tal` will return `command not found`.
*   **Suggested Fix:** Add a step: `mv source/orchestrator/m3tal .` or ensure `build.sh` explicitly handles the pathing.

#### 4. [WARNING] The `.env` Paradox
The docs list a table of `.env` variables and state that `m3tal init` generates tokens.
*   **The Issue:** Does `m3tal init` create the `.env` file from scratch? Does it append to one? Or does the user need to `cp .env.example .env` first? If the user runs `init` without a pre-existing `.env`, the system might crash or use hardcoded defaults that contradict the "Configuration" table.
*   **Suggested Fix:** Add a step: `cp .env.example .env` before running `./m3tal init`, or clarify that `init` bootstraps the file.

#### 5. [WARNING] Python/Flask Runtime Omission
The Dashboard is described as a "Python/Flask-based interface."
*   **The Issue:** Prerequisites only list Docker and Go. If the Dashboard runs inside Docker, this is fine. However, if the `m3tal` orchestrator attempts to run the dashboard locally for dev purposes, the user is missing Python 3.x and `pip install -r requirements.txt`.
*   **Suggested Fix:** Clarify if the Dashboard is **exclusively** containerized. If not, add Python to Prerequisites.

#### 6. [WARNING] Host Path Assumptions (`/mnt`)
The "Relationship Mapping" mentions "consistent path mapping at `/mnt` for storage operations."
*   **The Issue:** On many Linux distros (and specifically WSL), `/mnt` requires root permissions. If the Docker container maps to a host path in `/mnt` that hasn't been `chown`ed to the user, the media server will have permission denied errors.
*   **Suggested Fix:** Add a note: `sudo mkdir -p /mnt/media && sudo chown $USER:$USER /mnt/media`.

#### 7. [SUGGESTION] Traefik Port Conflict
`HTTP_PORT` is set to `8080`.
*   **The Issue:** Port `8080` is a very common default for development tools (Jenkins, Tomcat, other proxy dashboards). 
*   **Suggested Fix:** Recommend checking port availability or suggest `80` for a "Media Server" experience, with `8080` kept for the Traefik internal dashboard only.

#### 8. [SUGGESTION] Architecture vs. Reality
The README states the CLI communicates with the API on port `5050` to perform changes.
*   **The Issue:** If the API is not yet running (because the user hasn't run `./m3tal up`), will `./m3tal config set` fail?
*   **Suggested Fix:** Clarify which CLI commands are "Offline" (modifying files) and which are "Online" (interacting with the API).

---

### 🛠️ SUMMARY OF REQUIRED CHANGES
1.  Change **Go 1.26** to **Go 1.21**.
2.  Add **Git Submodule** instructions.
3.  Clarify the **Dashboard's** runtime (Docker vs. Host).
4.  Add a **Permission Warning** for `/mnt`.
5.  Standardize the **Binary Output Path** in the `build.sh` description.