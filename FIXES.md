**To:** M3TAL Development Team
**From:** DocCritic, Senior DevOps Auditor
**Subject:** AUDIT REPORT: M3TAL Media Server (v1.7) Documentation
**Date:** 2023-10-27

---

### **Verdict: FAILED**
The current documentation is an architectural "napkin sketch." It assumes the user is already a project contributor who knows the codebase. As a new user, I cannot deploy this; I am left with a directory of files and no actionable instructions on environment configuration, prerequisites, or external dependencies.

---

### **Detailed Issue List**

#### **BLOCKER**
*   **[BLOCKER] Missing `.env` template:** The documentation mentions `BASE_STORAGE_PATH` and `API_TOKEN` but provides no template or example `.env` file. A user running `./m3tal init` will likely encounter a crash or silent failure without these keys.
*   **[BLOCKER] Missing `m3tal-goback` setup:** The `Operational Flow` states the Dashboard *exclusively* communicates with `m3tal-goback`. However, there is no instruction on how to deploy or link this required external dependency. Without it, the system is fundamentally broken.
*   **[BLOCKER] Host Preparation:** The project mandates `/mnt` for storage. Does the user need to mount a drive to `/mnt` manually? Does the orchestrator create it? If it requires root access, the documentation is silent.

#### **WARNING**
*   **[WARNING] Traefik Access Info:** The documentation mentions the Traefik gateway but does not define which ports must be opened on the host firewall (e.g., 80, 443, 8080).
*   **[WARNING] Ambiguous CLI Usage:** `./m3tal init` is described as "Syncs configuration," but it’s unclear if this step *creates* the missing `.env` or just reads it.
*   **[WARNING] Missing Build Dependencies:** Is `Docker` and `Docker Compose` (v2) installed? The documentation implies they are required but fails to list them as prerequisites.

#### **SUGGESTION**
*   **[SUGGESTION] Dashboard Port Mapping:** Clearly define the internal vs. external ports for the Dashboard so users know where to point their browsers.
*   **[SUGGESTION] Log/Output expectation:** Add a "What to expect" section. After `./m3tal up`, the user has no idea if the process is blocking or if they should check container health.

---

### **Suggested Fixes**

1.  **Add a `Getting Started` Prerequisites section:**
    *   Explicitly list: Docker Engine (v20+), Docker Compose (v2+), and Go 1.21+.
    *   Provide a command to verify these: `docker compose version && go version`.

2.  **Provide an `.env.example` file:**
    *   Create a hidden file in the repo `example.env` and instruct the user: 
        `cp example.env .env && ./m3tal config` (to allow the binary to populate the values).

3.  **Clarify the `m3tal-goback` dependency:**
    *   Add a warning: "This system requires an instance of `m3tal-goback` to be reachable. If running locally, you must clone and deploy it separately." 
    *   Provide the expected variable name for the backend URL in the `.env`.

4.  **Networking/Firewall Section:**
    *   Update `NETWORKING.md` to list required open ports: 
        *   `80/TCP` (HTTP Gateway)
        *   `443/TCP` (HTTPS Gateway)
        *   `8080/TCP` (Traefik Dashboard - Optional)

5.  **Path Consistency Warning:**
    *   Add a bolded note: *"Ensure your media directory is mounted at `/mnt` on the host machine before running `init`. If using a custom path, ensure it is symlinked to `/mnt`."*

6.  **Update `Quick Start`:**
    *   Change the order: 
        1. Install Prerequisites.
        2. Clone `m3tal-goback`.
        3. Copy `.env`.
        4. Build binary.
        5. `./m3tal init`.
        6. `./m3tal up`.

**DocCritic Note:** *Do not release this to the public until a user can go from `git clone` to `dashboard login` without reading the source code. Fix these blockers immediately.*