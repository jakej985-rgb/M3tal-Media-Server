**DocCritic Audit Report: M3TAL Core Orchestrator**

**Verdict:** **FAILED**. As a new user, I am staring at a manifest of requirements with zero actionable instructions. This document reads like a system architecture whitepaper, not a README. It fails to provide the "how-to" for the most basic deployment steps. The project is currently non-deployable.

---

### Issue List

*   **BLOCKER:** **No "Getting Started" or Build Instructions.** There is no indication of how to build the `m3tal` binary or how to install it to `/usr/bin/m3tal`.
*   **BLOCKER:** **Missing Configuration Initialization.** The documentation references `/etc/m3tal/.env` as the "Source of Truth" but provides no example `.env` file, key names, or setup script.
*   **BLOCKER:** **Missing Orchestrator Setup.** The README mentions `m3tal` coordinates "Docker orchestration via the `/opt/m3tal` manifest tree." How does the directory get created? Do I pull the manifests from a repo? Is there an `m3tal setup` command?
*   **WARNING:** **Host Path Assumptions.** The documentation assumes `/mnt/m3tal-media` exists. If this is a requirement for the container to function, the script/documentation must ensure directory creation or warn the user.
*   **WARNING:** **Undefined Access/Ports.** The `docker-compose.yml` snippet is incomplete for a real-world scenario. There is no mapping for ports (Traefik gateway or UI access) or instructions on how to reach the dashboard.
*   **SUGGESTION:** **Confusing Wording.** The README claims "The platform is fully committed to Go-native binaries," yet the user is told to run a `docker-compose.yml` with `image: m3tal/core:latest`. Are we running a binary or a container? The distinction is blurred.

---

### Suggested Fixes

1.  **Add a "Quick Start" section:** Provide the exact sequence of commands to get up and running (e.g., `git clone`, `go build`, `mkdir -p /etc/m3tal`, `cp .env.example /etc/m3tal/.env`).
2.  **Provide an `.env.example` file:** Document every environment variable required (e.g., `M3TAL_API_KEY`, `DOCKER_SOCKET_PATH`, `MEDIA_ROOT`).
3.  **Define the Bootstrap Process:** Include a script or a command block to initialize the directory structure:
    ```bash
    sudo mkdir -p /etc/m3tal /opt/m3tal/stack /var/lib/m3tal /mnt/m3tal-media
    ```
4.  **Complete the Docker Compose:** Update the provided YAML to include Traefik labels or explicit port mappings so the user knows *where* to point their browser to see the Dashboard.
5.  **Clarify the Deployment Model:** Explicitly state: *"Run the orchestrator as a container (recommended) OR as a binary (advanced)."* Do not mix the two in the deployment section.
6.  **Add a "Prerequisites" section:** Explicitly list that Docker, Docker Compose, and Go 1.2x+ are required for this build. 

**DocCritic Note:** *Documentation is the first line of security and reliability. If a user has to guess where a file goes, they will misconfigure it. Fix these gaps immediately.*