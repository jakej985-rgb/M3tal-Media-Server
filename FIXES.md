**DocCritic Audit Report: M3TAL Core Documentation**
**Auditor:** Senior DevOps Auditor
**Severity Level:** CRITICAL
**Status:** **FAILED**

As a new user attempting to deploy this environment, I am currently staring at a broken shell. Your documentation assumes significant "tribal knowledge" regarding environment variables, absolute pathing, and container network dependencies.

---

### 🚨 Detailed Issue List

*   **BLOCKER: Missing `.env` Schema.**
    You tell the user to `cp template.env .env`, but provide no documentation on what variables are required. Is `BASE_STORAGE_PATH` mandatory? What about database credentials, API keys, or Docker registry settings?
*   **BLOCKER: Missing `/mnt` Pre-flight Check.**
    You mandate a strict `/mnt` mapping. If I am on macOS or a system where `/mnt` is not writable or does not exist, the orchestrator will likely panic. The docs don't mention if the orchestrator creates this or if the user must manually create it.
*   **BLOCKER: Undefined `build.sh` Requirements.**
    You command the user to run `chmod +x build.sh` and `./build.sh`. Does this require a specific Go version environment, or does it download dependencies? If it fails, where are the build logs?
*   **WARNING: Traefik Ingress Connectivity.**
    You list ports in a table, but you don't explain how the host maps to these. If Traefik is running on 8080, but I have another service there, the install will fail. There is no instruction on how to handle port collisions or verifying Traefik connectivity.
*   **WARNING: "API-Only Communication" Ambiguity.**
    You claim the dashboard talks to the backend via internal networks, but provide no `docker-compose.yml` insight. If I need to debug network issues, I don't know the network name.
*   **SUGGESTION: Binary Execution Context.**
    The documentation doesn't state if `./m3tal` needs root/sudo privileges to manipulate Docker sockets or bind to host ports.

---

### ✅ Suggested Fixes

1.  **Environment Documentation:** Add a table below the `Quick Start` section defining every key in `template.env`.
    *   *Example:* `BASE_STORAGE_PATH`: Absolute path to media directory (e.g., `/home/user/media`).
2.  **Explicit Path Validation:** Update the `m3tal init` process to perform a pre-flight check.
    *   *Docs update:* "Ensure your `BASE_STORAGE_PATH` exists on the host. If using Linux, ensure your user has ownership via `sudo chown $USER:$USER /mnt`."
3.  **Build Documentation:** Clarify `build.sh`.
    *   *Suggestion:* Add an output snippet of what a successful build looks like. Mention that it requires `go` to be in the `$PATH`.
4.  **Network/Port Awareness:** Explicitly state the dependency on port 80/443 for Traefik.
    *   *Add:* "Warning: Ensure ports 80, 443, and 8080 are free before running `m3tal up`."
5.  **Troubleshooting the Host File:** Your `/etc/hosts` instructions are manual. Provide a one-liner for Linux/macOS users to automate this:
    *   *Add:* `echo "127.0.0.1 m3tal.localhost api.localhost traefik.localhost" | sudo tee -a /etc/hosts`
6.  **Permission Clarity:** Explicitly mention if the orchestrator requires `sudo`. If it doesn't, ensure the documentation mentions the user must be in the `docker` group.

---

**Verdict:** **UNUSABLE**. The project is currently a "Black Box." A new user cannot successfully deploy this without guessing variable names or debugging the Go binary's source code. **Documentation must be updated before this is considered production-ready.**