### **DocCritic Audit Report: M3TAL Core Orchestrator**

**Verdict: BLOCKER.**
As a new user, I cannot deploy this. The documentation describes an architecture but provides no "Getting Started" path. It assumes the user is already a developer familiar with the internal manifest tree and binary compilation. The documentation is currently an architectural manifesto, not an installation manual.

---

### **Issue List**

#### **BLOCKERS**
1.  **[BLOCKER] Missing Installation/Compilation Instructions:** The README mentions a `m3tal` binary but provides no instructions on how to build it from source (`go build`) or where to download it.
2.  **[BLOCKER] No Configuration Lifecycle:** There is no documentation on how to initialize `/etc/m3tal/.env`. What variables are required? Is there a `.env.example`? Without these, the core container will fail to start or connect to the API.
3.  **[BLOCKER] Missing Network/Port Mapping:** The Docker Compose snippet provided is missing the `ports` declaration. Even if it relies on Traefik, the internal port exposure is undocumented. How does `m3tal-godash` talk to `m3tal-goback`?

#### **WARNINGS**
4.  **[WARNING] Path Assumption Error:** The documentation assumes `/mnt/m3tal-media` exists on the host. If a user runs the container without creating this directory, Docker will create it as `root`, causing permission issues for the media user.
5.  **[WARNING] Docker Socket Security:** The instructions mount `/var/run/docker.sock` with read-write access. This is a massive security risk. There is no mention of security considerations or best practices (e.g., using a proxy like `socket-proxy`).
6.  **[WARNING] Missing Orchestration Workflow:** The document mentions `/opt/m3tal/stack` but doesn't explain how to trigger the deployment of these stacks. Does the binary watch this folder? Do I need to run an `m3tal init` command?

#### **SUGGESTIONS**
7.  **[SUGGESTION] Prerequisites Section:** A "System Requirements" section (e.g., Docker, Docker Compose, Go 1.x+) is missing.
8.  **[SUGGESTION] Missing Startup Command:** There is no instruction on how to run the binary once compiled. Does it require systemd? Is it intended to run in the foreground?

---

### **Suggested Fixes**

*   **Implement a "Quick Start" section:**
    ```bash
    # 1. Clone & Build
    git clone ...
    go build -o m3tal ./cmd/m3tal/main.go
    
    # 2. Setup Config
    sudo mkdir -p /etc/m3tal /opt/m3tal/stack /mnt/m3tal-media
    cp .env.example /etc/m3tal/.env
    ```

*   **Provide a mandatory `.env` template:** Create a file listing required variables (e.g., `M3TAL_API_KEY`, `M3TAL_CORE_PORT`, `DATABASE_URL`).

*   **Update Docker Compose with networking:**
    ```yaml
    services:
      m3tal-core:
        ...
        ports:
          - "8080:8080" # Map API port
        networks:
          - m3tal-net
    networks:
      m3tal-net:
        driver: bridge
    ```

*   **Add "First-Run" instructions:**
    *   Explicitly define that the user must run `m3tal setup` to generate the `/var/lib/m3tal` directory structure.
    *   Provide a warning about `mkdir /mnt/m3tal-media` and `chown` permissions.

*   **Security Disclaimer:** Add a note regarding the Docker socket: *"Warning: Mounting the Docker socket grants the container root-level access to the host. Ensure the environment is secured."*