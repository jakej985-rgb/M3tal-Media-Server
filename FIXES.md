**Audit Report: M3TAL Core Command Center (v1.5)**
**Auditor:** DocCritic, Senior DevOps Auditor
**Verdict:** **UNUSABLE / BLOCKER**

The current documentation suffers from "Institutional Knowledge Syndrome." It assumes the user understands the internal state-management of the Go orchestrator without providing the necessary validation steps. The instructions are riddled with port conflicts, missing directory creation steps, and ambiguous state-management requirements.

---

### 🚨 BLOCKER Issues

1.  **Port Conflict (System-Breaking):** The docs state `8080` is the default for both the `HTTP_PORT` (Traefik) and the `DASHBOARD_PORT`. You cannot bind Traefik and an internal service to the same port on the host. 
    *   *Fix:* Separate the internal service ports from the Host-exposed Traefik ports. Ensure the README explicitly states which ports must be open on the host vs. which are internal.
2.  **Missing `init` context:** You instruct users to run `./m3tal init`, but never define what this script does to the host environment (e.g., creating folders, generating SSL certs, writing files to `/usr/share/m3tal/`). 
    *   *Fix:* Include a section on what `./m3tal init` modifies on the host file system.
3.  **Dependency Black Box:** The `build.sh` script is referenced but its requirements are not clearly defined. Does it require `go` installed on the host? Yes. Does it require `git`? Yes. 
    *   *Fix:* Explicitly list build-time dependencies (e.g., `build-essential`, `golang-go`).

---

### ⚠️ WARNING Issues

1.  **Storage Assumption:** You state the orchestrator maps `BASE_STORAGE_PATH` to `/mnt`. If the user has not created the directory defined in `BASE_STORAGE_PATH`, will the orchestrator create it, or will Docker create it as `root` (preventing user read/write access later)?
    *   *Fix:* Add a step: `mkdir -p ./data && chown $USER:$USER ./data`.
2.  **Dashboard/API URL Confusion:** You provide URLs like `http://m3tal.localhost:8080`. Users will get "Connection Refused" unless they have a local DNS resolver or edit their `/etc/hosts` file.
    *   *Fix:* Add a mandatory step to edit `/etc/hosts` to map these domains to `127.0.0.1`.
3.  **Inconsistent Cleanup:** You mention `./m3tal down` in one place and `make down` in another.
    *   *Fix:* Standardize all commands to the orchestrator binary. Do not suggest `make` unless a `Makefile` is provided in the repo.

---

### 💡 SUGGESTION Issues

1.  **Ambiguous CLI Syntax:** You list `./m3tal dashpass admin yourpassword`.
    *   *Improvement:* Provide a flag-based or prompt-based alternative to avoid exposing passwords in the bash history.
2.  **Traefik Admin Panel:** You mention `http://traefik.localhost:8080`. Providing a raw, unprotected link to the Traefik dashboard is a security risk.
    *   *Improvement:* Add a note about protecting the Traefik dashboard with basic auth in the configuration.
3.  **Host Path Clarification:** The table lists `BASE_STORAGE_PATH` as `./data`. If a user moves the binary or changes working directories, this breaks. 
    *   *Improvement:* Recommend the use of absolute paths (e.g., `/home/username/m3tal/data`) in the `.env` file to prevent path resolution issues.

---

### Summary of Recommended Actions
1.  **Revise Ports:** Define non-conflicting default ports (e.g., Traefik 80/443; API/Dashboard as internal Docker-only ports).
2.  **Pre-flight Script:** Add a `check.sh` script that verifies `docker`, `go`, and port availability before the user attempts to run the orchestrator.
3.  **Explicit `/etc/hosts` Guidance:** Include the exact lines the user needs to add to their host machine for the `.localhost` domains to resolve.
4.  **Ownership Check:** In the `init` protocol, verify if the `BASE_STORAGE_PATH` has the correct UID/GID permissions for the user running the orchestrator.