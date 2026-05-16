### **Audit Report: M3TAL Core Orchestrator Documentation**
**Auditor:** DocCritic, Senior DevOps Auditor  
**Status:** **REJECTED**

---

### **Verdict**
The current documentation is an architectural whitepaper, not an installation guide. As a new user, I have no idea how to actually run this software. The instructions assume I am already part of the M3TAL development team. The lack of a clear "Quick Start" guide and missing prerequisite configuration steps makes this repository unusable for a general user.

---

### **Detailed Issue List**

#### **BLOCKER**
1. **Missing Initialization Procedure:** The documentation mentions a `m3tal` binary, but fails to provide a `build` or `install` command. Do I run `go build`? Is there a Makefile?
2. **Missing `.env` Template:** The docs state `/etc/m3tal/.env` is the "Source of Truth," but provide no template or required variables (e.g., API keys, database URLs, mount paths). A blank install will crash immediately.
3. **Missing "First-Run" logic:** There is no command provided to initialize the state (e.g., `m3tal setup` or `m3tal init`).
4. **Volume Dependency Crash:** The Docker configuration maps `/mnt:/mnt`. If I don't have `/mnt/m3tal-media` pre-created with correct permissions (UID/GID), the container will likely throw permission errors or boot-loop.

#### **WARNING**
5. **Ambiguous Entry Point:** The README mentions both `/opt/m3tal/stack` and `/docker`. Which one is the intended directory for users to run `docker compose up`? 
6. **Network/Port Siloing:** The documentation mentions Traefik but provides zero information on required ports (e.g., 80, 443, 8080) or how to link the API and Dashboard containers to the Orchestrator network.
7. **No "Source/Usage" workflow:** I see references to `m3tal-goback` and `m3tal-godash`, but no instruction on how to deploy them *together* with the Core. Do I need one `docker-compose.yml` that handles all three? 

#### **SUGGESTION**
8. **Lack of "Prerequisites" Check:** Add a shell script or command to verify that Docker, Go, and the necessary directory structures exist before attempting installation.
9. **Role Confusion:** The "System Architecture" section describes the Orchestrator, API, and Dashboard as separate components, but the deployment section only shows the Core. I need a unified `docker-compose.yml` example.

---

### **Required Remediation Steps**

1.  **Add a `Quick Start` section:**
    *   Step 1: Clone repo.
    *   Step 2: `make build` (or equivalent).
    *   Step 3: `cp .env.example /etc/m3tal/.env`.
    *   Step 4: `m3tal init` (ensure this creates the `/opt/m3tal` tree).
2.  **Provide a Full Stack `docker-compose.yml`:** Include the Orchestrator, Backend, and Dashboard in a single YAML file so users aren't guessing how to network them.
3.  **Define the Environment:** Create a `.env.example` file in the root of the repo with every required variable clearly commented.
4.  **Add a "Ports & Networking" table:** Clearly list which ports the user must expose to the host for the Dashboard and API.
5.  **Explicit Path Creation:** Add a note or a small script snippet: `mkdir -p /mnt/m3tal-media && chown -R 1000:1000 /mnt/m3tal-media` to prevent permission-denied errors on startup.

**DocCritic's Final Note:** *Stop treating the README as a marketing document and start treating it as a technical manual. If a user can't be at a `docker compose up` command within 5 minutes of reading, the documentation has failed.*