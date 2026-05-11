**To:** M3TAL Development Team  
**From:** DocCritic, Senior DevOps Auditor  
**Date:** 2023-10-27  
**Subject:** Audit Report – M3TAL Control Plane Documentation

---

### **Verdict: FAILED**
The current documentation is a "developer’s assumption" trap. It fails to provide a path to a successful deployment for an external user. The instructions rely on tribal knowledge (hidden requirements in scripts) rather than explicit documentation. It is currently impossible to deploy this stack reliably using only the provided README.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing `m3tal.py` Build Steps:** The README mentions `source/go-backend` and `source/dashboard`, but never explains if/how to build these Go binaries or compile the dashboard before running `python m3tal.py up`.
2.  **Missing Traefik/Gateway Configuration:** The documentation mentions "Traefik" implicitly, but provides zero information on how to route traffic, which ports are exposed on the host, or how to access the dashboard after deployment.
3.  **Dependency Hell:** `install.py` is mentioned, but `requirements.txt` or `go.mod` handling is ignored. Does `install.py` handle binary compilation, or does the user need to manually `go build`?

#### **WARNING**
4.  **Implicit Path Assumptions:** The documentation forces `/mnt` (root directory) usage. This requires `sudo` access and assumes the user has a secondary partition or is comfortable writing to the system root. This is dangerous and non-standard for containerized deployments.
5.  **`.env` Lifecycle:** The user is told to create a `.env` file, but there is no `env.example` provided. A user has no way of knowing all required variables (e.g., `DB_PASSWORD`, `TRAEFIK_TAGS`, etc.).
6.  **Venv Ambiguity:** You suggest `source venv/bin/activate` but never provide the command to *create* the virtual environment.

#### **SUGGESTION**
7.  **Directory Structure:** The project structure is complex (`source/m3tal-stack`, `source/go-backend`, etc.). Include a tree visualization to help users locate the components.
8.  **Port Mapping:** A table of default ports (80, 443, 8080, etc.) is mandatory for a system that acts as an "Orchestrator."

---

### **Suggested Fixes**

1.  **Refine "Installation" Section:**
    *   Add: `pip install -r requirements.txt`
    *   Add: `cd source/go-backend && go build -o bin/backend .` (and explain if `m3tal.py` expects this binary to exist).
2.  **Provide an Environment Template:** 
    *   Create `cp .env.example .env`. Document every variable in the example file.
3.  **Decouple Storage:** 
    *   Allow an environment variable `BASE_STORAGE_PATH` to default to a user-local directory (e.g., `./data`) rather than forcing `sudo` creation of `/mnt`.
4.  **Clarify `m3tal-stack` Usage:**
    *   Explicitly state that `m3tal.py` calls `docker compose -f source/m3tal-stack/docker-compose.yml`.
5.  **Expose Port Info:**
    *   Add a "Connectivity" section: 
        *   `Dashboard`: `http://localhost:8080`
        *   `Traefik Dashboard`: `http://localhost:8081`
6.  **Provide Setup Commands:** 
    *   Explicitly list: 
        ```bash
        python3 -m venv venv
        source venv/bin/activate
        pip install -r requirements.txt
        ```

---

**Auditor Note:** *You are currently building a tool for system automation. If the "Installation" experience is as chaotic as the current README suggests, users will abandon the project immediately. Fix these blockers before the next release.*