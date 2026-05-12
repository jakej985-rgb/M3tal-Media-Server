# Audit Report: M3TAL Control Plane (v1.4.0.3)

**Auditor:** DocCritic, Senior DevOps Auditor  
**Verdict:** **NON-COMPLIANT / BLOCKER**

The current documentation assumes a "happy path" environment and fails to address critical security, networking, and environment state requirements. A user attempting to deploy this today will encounter permission errors, container networking failures, and absolute uncertainty regarding reverse proxy configuration.

---

### 🔴 BLOCKER Issues

1.  **Hardcoded `/mnt` Path Assumptions:**
    *   **Issue:** The `scripts/setup.sh` script assumes it has permission to modify `/mnt` (often root-owned or mount points). It also fails to account for existing data in those paths.
    *   **Fix:** Implement a check in `setup.sh` to verify `BASE_STORAGE_PATH` exists, is writable, and create a nested structure (e.g., `$BASE_STORAGE_PATH/m3tal/data`) instead of assuming the root of a partition.
2.  **Missing Dashboard/API Python Dependencies:**
    *   **Issue:** There is zero instruction on how to run the `source/dashboard` (the Python app). Is it inside a container? Is it local? If local, there is no `requirements.txt` execution step.
    *   **Fix:** Explicitly state if the Dashboard is containerized in the compose stack. If so, document the `Dockerfile` build context. If local, add a step to create a venv and run `pip install -r requirements.txt`.
3.  **Traefik / Reverse Proxy Gap:**
    *   **Issue:** The doc mentions Traefik in the notes, but provides no `docker-compose.yaml` labels or network configuration to actually *use* it.
    *   **Fix:** Provide a "Traefik Example" section in the `m3tal-stack/docker-compose.yml` showing labels for automatic discovery (e.g., `traefik.enable=true`, `traefik.http.routers...`).

---

### 🟡 WARNING Issues

1.  **Incomplete `.env` instructions:**
    *   **Issue:** You mention generating tokens with `openssl`, but you do not show the user *where* in the `.env` file to put them.
    *   **Fix:** Provide a template snippet or a `sed` command to inject the output of `openssl` directly into the `.env` file.
2.  **Service Access Confusion:**
    *   **Issue:** The table lists `8082` for the dashboard. If the user is on a remote VPS or a specific LAN IP, `localhost` will fail.
    *   **Fix:** Advise users to replace `localhost` with their host's static LAN IP and warn that browsers may block insecure connections if they attempt to access via HTTP instead of HTTPS/Traefik.

---

### 🔵 SUGGESTION Issues

1.  **"M3TAL Setup" Command:**
    *   **Issue:** The project name implies a "Control Plane," yet the setup is fragmented (run script, then cd here, then run docker).
    *   **Fix:** Create a top-level `Makefile` or `m3tal.sh` wrapper script. One command (`make install` or `./m3tal.sh up`) should handle setup, `.env` generation, and deployment.
2.  **Missing "Cleanup" documentation:**
    *   **Issue:** No mention of how to safely stop or purge the stack.
    *   **Fix:** Add a section for "Maintenance & Cleanup" covering `docker compose down -v` and safe backup of the `BASE_STORAGE_PATH`.
3.  **Volume Mapping Transparency:**
    *   **Issue:** Users have no idea what data is stored where.
    *   **Fix:** Add a small table explaining the bind-mounts: e.g., `/mnt/m3tal/db` -> `/var/lib/postgresql/data`.

---

### Auditor's Final Word
*The documentation is currently written for the original developer. It assumes the user already knows the relationship between the folders and the container orchestration. **Re-factor to prioritize a single-entry-point setup process.** Without clear instructions on how the Dashboard and API interact via Docker networks, users will report "Connection Refused" errors within 10 minutes of operation.*