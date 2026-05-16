### **DocCritic Audit Report: M3TAL Core Orchestrator**

**Verdict: FAILED**
The current documentation is an architectural whitepaper, not an installation manual. It describes what the system *is* but provides zero actionable steps on how to make it *run*. A new user would be unable to deploy this project without significant trial-and-error, manual path creation, and guesswork regarding environment variables.

---

### **Issue List**

#### **BLOCKER**
*   **[BLOCKER] Missing Initialization/Setup:** There is no mention of `m3tal.py setup` or any equivalent binary initialization. The user has no idea how to generate the `/etc/m3tal/.env` file or the required folder structure.
*   **[BLOCKER] Missing Docker Compose Instructions:** The provided YAML snippet is a fragment. There is no `docker-compose.yml` file structure, no networking instructions, and no mention of how to actually trigger the deployment.
*   **[BLOCKER] Filesystem Assumptions:** The docs assume `/opt/m3tal` and `/mnt/m3tal-media` exist. If a user runs the container without these directories pre-created, Docker will often create them as `root`-owned directories, leading to immediate permission failures.
*   **[BLOCKER] Missing Port Mapping/Traefik:** The Core Orchestrator is meant to be a gateway. There are no port mappings defined in the YAML, nor is there documentation on how Traefik or any ingress controller should point to the system.

#### **WARNING**
*   **[WARNING] Missing Configuration Schema:** The `.env` file is mentioned as the "Source of Truth," but there is no example of what keys are required (e.g., API keys, database URLs, log levels).
*   **[WARNING] Dependency Orchestration:** It is unclear if the user should run `m3tal-goback` and `m3tal-godash` manually, or if the Core Orchestrator is meant to spawn them.

#### **SUGGESTION**
*   **[SUGGESTION] Prerequisites Section:** Explicitly state the need for Docker, Docker Compose, and Go (if compiling from source).
*   **[SUGGESTION] Command Examples:** Add "Quick Start" code blocks for common tasks (e.g., `m3tal init`, `m3tal up`).

---

### **Suggested Fixes**

1.  **Create a `setup.sh` or `m3tal.py` script:** Include a command that automatically runs:
    ```bash
    mkdir -p /etc/m3tal /opt/m3tal/stack /var/lib/m3tal /mnt/m3tal-media
    # Generate initial .env with default template
    ```
2.  **Provide a complete `docker-compose.yml`:** Include a template that defines the network and labels for Traefik.
    ```yaml
    services:
      m3tal-core:
        image: m3tal/core:latest
        labels:
          - "traefik.enable=true"
          - "traefik.http.routers.m3tal.rule=Host(`m3tal.local`)"
    ```
3.  **Document the `.env` schema:** Provide a `.env.example` file:
    ```text
    M3TAL_API_KEY=your_secret_key
    STORAGE_PATH=/mnt/m3tal-media
    LOG_LEVEL=info
    ```
4.  **Add an "Installation Flow" section:**
    *   **Step 1:** Clone repo.
    *   **Step 2:** Run `sudo ./scripts/setup.sh` (to create paths and permissions).
    *   **Step 3:** Configure `/etc/m3tal/.env`.
    *   **Step 4:** Deploy via `docker-compose up -d`.

**DocCritic's Final Note:** *Stop treating this repo like a research project and start treating it like a product. Users need commands they can copy-paste, not a manifesto on Go-native architecture.*