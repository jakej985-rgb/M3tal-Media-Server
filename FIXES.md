To: DocSmith
From: DocCritic, Senior DevOps Auditor
Subject: Audit Report - M3TAL Core Orchestrator README

### Verdict: FAILED
The current documentation is an architectural overview, not a deployment guide. It assumes the user has deep tribal knowledge of your specific environment. It fails to address critical security, networking, and filesystem pathing requirements, making a successful deployment highly improbable for a new user.

---

### Issue List

#### 1. BLOCKER: Missing Port/Ingress documentation
The documentation mentions a "Dashboard" and "Backend API" but provides zero information on how to access them.
*   **Fix**: Explicitly document that Traefik is expected as a gateway or provide the mapping for exposed ports (e.g., `8080` for dashboard). If Traefik labels are required, they must be documented in a "Network Requirements" section.

#### 2. BLOCKER: Fragile Filesystem Assumptions
The README defines `/mnt/m3tal-media` as a requirement but provides no instructions on how to set this up. If the directory is missing, does the service crash? Does it need specific permissions (`chown`) for the docker user?
*   **Fix**: Add a setup step: `sudo mkdir -p /mnt/m3tal-media && sudo chown $USER:$USER /mnt/m3tal-media`.

#### 3. WARNING: Ambiguous Docker Deployment
The "Deployment: Docker Configuration" section is a YAML snippet with no context. Where does this go? Is it a `docker-compose.yml` file? Where is the file supposed to live?
*   **Fix**: Rename to "Manual Compose Deployment" and provide a full file path (e.g., `docker-compose.yml` in the project root) and a clear command to execute it (`docker compose up -d`).

#### 4. WARNING: Hidden APT Requirements
You assume the user is using Debian/Ubuntu. If they are on RHEL, Fedora, or Arch, the installation will fail silently or explicitly.
*   **Fix**: Add a warning banner: *"Warning: Official APT repository support is currently limited to Debian/Ubuntu-based distributions. Building from source is required for other Linux distributions."*

#### 5. SUGGESTION: Remove "Buzzword" Marketing
Phrases like "Go-Native Migration Active" and "Modular Infrastructure Platform" add zero technical value.
*   **Fix**: Remove the final "Ecosystem Integration Rules" and "Related Projects" section fluff. Keep documentation strictly technical. A user needs to know how to deploy, not the project's current status in your internal roadmap.

#### 6. SUGGESTION: CLI Demo Gaps
The `m3tal up` command implies that it automatically pulls containers or expects them in `/opt/m3tal/stack`. If the user hasn't cloned that repository or copied files there, the command will fail.
*   **Fix**: Add a step: `git clone <repo> /opt/m3tal` before running `m3tal up`.

---

### Suggested README Structure (Abbreviated)

**[Installation]**
*   *(Keep existing APT block)*
*   **System Prep:** 
    ```bash
    sudo mkdir -p /mnt/m3tal-media
    sudo chown $USER:$USER /mnt/m3tal-media
    ```

**[Networking & Access]**
*   **Dashboard:** Accessible at `http://localhost:8080` (or define Traefik entrypoint).
*   **API:** Internal only. Ensure Docker bridge allows communication between `m3tal-goback` and `m3tal-orchestrator`.

**[Deployment]**
*   **Standard Usage:**
    1. Initialize: `m3tal setup`
    2. Deploy Stack: `m3tal up`
    3. Start UI: `m3tal dash up`

**[Troubleshooting]**
*   *Note: If `m3tal up` fails, ensure `/opt/m3tal/stack` contains valid `docker-compose.yml` manifests.*

---
**Audit Complete.** Correct these items before resubmission.