**DocCritic Audit Report: M3TAL Core Orchestrator**

**Verdict: FAILED**
The documentation exhibits dangerous assumptions regarding host system state and fails to provide the basic networking and deployment configuration required for a "Core Orchestrator." It reads more like a project announcement than a technical manual.

---

### Issue List

#### 1. BLOCKER: Missing Mount Point Validation
The README assumes the host environment already has `/mnt/m3tal-media` configured. A new user will experience immediate runtime failures if this directory does not exist or lacks correct permissions.
*   **Fix:** Add a pre-flight check section or a `m3tal setup` step that explicitly validates/creates required directories (e.g., `mkdir -p /mnt/m3tal-media && chown $USER:$USER /mnt/m3tal-media`).

#### 2. BLOCKER: Missing Port/Access Documentation
The guide mentions a Dashboard and API but provides zero information on how to access them. Users don't know which ports to open on their firewall or where to point their browsers.
*   **Fix:** Include a "Network Access" section explicitly stating that the dashboard runs on port `[X]` (e.g., Traefik/Dashboard gateway).

#### 3. WARNING: Docker Deployment Ambiguity
The "Deployment: Docker Configuration" section provides a snippet of a YAML file, but does not explain *where* this file lives or how to execute it. Is it a `docker-compose.yml` file? Do I run `docker compose up`?
*   **Fix:** Provide the full `docker-compose.yml` boilerplate and explicit start commands. 

#### 4. WARNING: Ecosystem "Quick Demo" Inconsistency
You mention `m3tal dash up` in the Quick Demo, but the "Architecture Overview" says the dashboard is a separate project (`m3tal-godash`). The user has no instruction on how to link these two repositories.
*   **Fix:** Clarify if `m3tal dash up` triggers a git clone of the other repo, or if the user must install the dashboard independently.

#### 5. SUGGESTION: Marketing Fluff Removal
Phrases like "Go-Native Architectural Requirements" and "Modular Infrastructure Platform" add zero value to an auditor or a developer.
*   **Fix:** Strip all marketing adjectives. Documentation should be functional, not a sales pitch.

#### 6. SUGGESTION: APT Repository Safety
The installation instructions rely on `tee` with `sudo` directly. This is fine, but you should explicitly note the distro requirements (e.g., "Tested on Debian 12 / Ubuntu 22.04+").
*   **Fix:** Add a "Supported Platforms" section.

---

### Suggested README Refactor (Critical Sections)

**[REPLACE "Deployment" Section with the following:]**

#### Deployment
The Orchestrator requires a standardized Docker Compose environment to manage services. 

1. Create `/opt/m3tal/stack/docker-compose.yml`:
```yaml
services:
  orchestrator:
    image: m3tal/core:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /mnt/m3tal-media:/mnt/m3tal-media
    ports:
      - "8080:8080" # Dashboard Access
```
2. Initialize and start:
```bash
sudo mkdir -p /mnt/m3tal-media
m3tal setup
m3tal up
```
*Access the dashboard at `http://<host-ip>:8080`.*

---

**Auditor Note:** *Do not push to main until the `/mnt` directory creation is automated or explicitly documented as a manual prerequisite.*