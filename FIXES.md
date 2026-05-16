**To:** M3TAL Development Team  
**From:** DocCritic, Senior DevOps Auditor  
**Subject:** AUDIT REPORT: M3TAL-CORE-README-001  

---

### **Verdict: FAILED**
The current documentation is an architectural vision statement, not a deployment guide. It lacks the functional "Last Mile" instructions required for a user to actually boot the system. It assumes the user is already an expert in the internal M3TAL ecosystem and ignores common failure points during fresh installation.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing "Getting Started" Workflow:** The documentation mentions a `m3tal` binary but provides no instructions on how to acquire it (Build from source? Binary release? `go install`?).
2.  **Zero Configuration Guide:** There is no documentation for `/etc/m3tal/.env`. A new user has no idea what environment variables are required for the containers to actually start (e.g., API keys, port mappings, database credentials).
3.  **Missing Traefik/Gateway Specs:** The README mentions an "API-only communication" protocol but provides no guidance on how to route traffic to the `goback` or `godash` containers. Without Traefik or an exposed port mapping, the ecosystem is a black box.

#### **WARNING**
4.  **Implicit Host Requirements:** The guide assumes the existence of `/mnt/m3tal-media` and `/opt/m3tal`. If a user blindly runs the Docker Compose, the volume mounting will either fail or create root-owned directories on the host, leading to permissions hell.
5.  **Ambiguous Orchestration:** The instructions say the `m3tal` binary "coordinates Docker orchestration," but it does not specify if the user needs to manually run a `docker-compose up` or if the `m3tal` binary triggers the stack itself.

#### **SUGGESTION**
6.  **"Runbook" Format:** The documentation is too abstract. It needs a "Quick Start" section with copy-pasteable commands.
7.  **Container Tagging:** Using `latest` for the `m3tal-core` image is a DevOps anti-pattern. Use versioned tags to prevent unexpected production breaks.

---

### **Required Fixes**

1.  **Add Installation Section:**
    *   Provide a `git clone` command.
    *   Add a `make build` or `go build -o m3tal ./cmd/main.go` instruction.
    *   Explain how to move the binary to `/usr/bin/m3tal`.

2.  **Add Config Template:**
    *   Include a `m3tal.env.example` block in the README so users know what to put in `/etc/m3tal/.env`.
    *   **Mandatory keys:** `DB_URL`, `API_KEY`, `DOCKER_NETWORK`, `PORT`.

3.  **Define Infrastructure Pre-reqs:**
    *   Add a script snippet: 
        ```bash
        sudo mkdir -p /etc/m3tal /opt/m3tal/stack /mnt/m3tal-media
        sudo chown -R $USER:$USER /opt/m3tal
        ```

4.  **Refine Docker Section:**
    *   Provide a full `docker-compose.yml` example that includes `m3tal-goback` and `m3tal-godash`. 
    *   Define the network configuration. A user cannot guess how these containers "talk" to each other via the API.

5.  **Clarify CLI Role:** 
    *   Explicitly state: "Run `m3tal setup` to initialize the `/opt/m3tal` manifest tree before starting the Docker stack."

---
*DocCritic Note: Do not push this to production README until a user who has never seen this repo can successfully reach the dashboard by following the steps sequentially.*