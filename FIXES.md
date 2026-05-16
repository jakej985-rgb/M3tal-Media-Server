**To:** M3TAL Development Team  
**From:** DocCritic, Senior DevOps Auditor  
**Subject:** AUDIT REPORT: M3TAL Core Orchestrator README  

### **Verdict: FAILED**
**Status:** The current documentation is an architectural manifesto, not a deployment manual. It describes *what* the system is, but provides zero actionable intelligence on *how* to deploy it. A new user cannot get this running based on the provided text.

---

### **Issue List**

1. **BLOCKER: Missing Initialization Workflow**  
   There is no "Quick Start" or "Installation" section. How do I build the binary? How do I initialize `/var/lib/m3tal`? Does the orchestrator auto-generate its own config, or do I need to create it manually?

2. **BLOCKER: `.env` Schema Undefined**  
   The documentation mentions `/etc/m3tal/.env` as the "Source of Truth" but fails to provide a template or required variables. The system will fail to start without a defined schema.

3. **BLOCKER: Missing Orchestrator Usage Commands**  
   The README mentions `m3tal` CLI, but provides zero usage examples. How do I initiate a stack? How do I trigger the `m3tal-goback` integration?

4. **WARNING: Implicit Infrastructure Requirements**  
   The project assumes a specific host directory structure (`/mnt/m3tal-media`, `/opt/m3tal/stack`). If these paths don't exist, will the container crash or auto-create them? This is a "Dev-only" assumption that breaks in production environments.

5. **WARNING: Networking/Gateway Blindness**  
   The documentation mentions a Dashboard and API, but fails to document port mapping, Traefik labels, or ingress requirements. Users will not be able to access the UI.

6. **SUGGESTION: Docker Deployment Gaps**  
   The `docker-compose` snippet provided is bare-bones. It lacks `network_mode`, `labels` (for Traefik), and proper `restart` policies. Using `latest` tags is a security risk; pin to specific versions or hashes.

---

### **Suggested Fixes**

*   **Add an "Initialization" section:**
    *   `mkdir -p /etc/m3tal /var/lib/m3tal /opt/m3tal/stack`
    *   Provide a `sample.env` file block.
    *   Document the binary build process: `go build -o m3tal main.go` or "Download the release from [link]".
*   **Add a "Gateway" section:**
    *   Define necessary ports (e.g., 8080 for API, 3000 for Dash).
    *   Include a Traefik label snippet:
        ```yaml
        labels:
          - "traefik.http.routers.m3tal.rule=Host(`m3tal.example.com`)"
          - "traefik.http.services.m3tal.loadbalancer.server.port=8080"
        ```
*   **Clarify Paths:**
    *   Add a warning: *"Ensure host directories /mnt/m3tal-media exist before deployment to prevent volume bind-mount errors."*
*   **Add "Basic Commands":**
    *   `m3tal --help`
    *   `m3tal stack up [project-name]`
    *   `m3tal status`
*   **Define Architecture:**
    *   Provide a minimal `docker-compose.yml` that pulls the full stack (Core + API + Dash) rather than just the Core, so the user has a functioning UI upon first run.

**DocCritic's Final Note:** You have built a engine, but you haven't provided the keys to the ignition. Fix the installation flow immediately or expect 100% of new users to abandon this project.