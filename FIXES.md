## M3TAL README Audit Report

**Verdict:** FAILED - Critical information is missing, making successful deployment and operation impossible without external knowledge.

---

### Issue List:

1.  **BLOCKER: APT installation command block missing.**
    *   **Description:** The README fails to provide the complete 3-command sequence (keyring, repo, install) required for APT installation. While the commands are present, they are not presented as a single, contiguous block, and the `sudo apt update` is separated from the `sudo apt install`.
    *   **Required Fix:** Present the three commands as a single, executable block as per the GROUND TRUTH.

2.  **BLOCKER: Docker dependency not explicitly stated.**
    *   **Description:** While the README mentions Docker Engine and Docker Compose V2 in the prerequisites, it does not explicitly state that *both* are required. The GROUND TRUTH clarifies that M3TAL uses both internally.
    *   **Required Fix:** Explicitly state "Docker Engine **and** Docker Compose V2 are required."

3.  **BLOCKER: Deployment lifecycle explanation is incomplete.**
    *   **Description:** The README explains that `m3tal up` operates on `*-compose.yml` files in `/docker/`, but it omits crucial details:
        *   It doesn't explicitly state that `/docker` is a symlink to `/opt/m3tal/stack/`.
        *   It doesn't explain that `m3tal up` *orchestrates* Docker Compose V2 internally for these files.
        *   It doesn't mention the `m3tal dash up` command for managing the dashboard specifically.
        *   The explanation of adding new compose files is present, but it relies on the missing context of how `m3tal up` works.
    *   **Required Fix:** Incorporate the details about the symlink, the direct use of `docker compose`, and the `m3tal dash up` command.

4.  **BLOCKER: Traefik routing explanation is insufficient.**
    *   **Description:** The README states Traefik is the HTTP gateway and mentions it uses labels and dynamic config. However, it lacks specific examples of *how* services get exposed through Traefik:
        *   It doesn't show the Traefik labels needed for a service to be discovered by Traefik.
        *   It doesn't clearly link the `api.DOMAIN` routing to the `host.docker.internal:8080` target as described in the GROUND TRUTH.
        *   The example for exposing a custom service is provided, which is good, but the core explanation of how Traefik *works* is abstract.
    *   **Required Fix:** Clarify that Traefik uses `traefik.enable=true` and associated labels for service discovery. Explicitly mention that Traefik routes to internal container ports or host-local ports like `host.docker.internal:8080`.

5.  **WARNING: Port table missing required ports.**
    *   **Description:** The README's port table lists ports 80, 8080, 8081, and 8082. However, it's missing explicit clarification on how each port is accessed (public vs. host-local).
    *   **Required Fix:** Add an "Access" column to the port table to specify whether a port is "Public", "Host-local", or both.

6.  **WARNING: Service management details are incomplete.**
    *   **Description:** While the README mentions `systemctl` for `m3tal-api.service`, it doesn't provide the specific commands for managing it (start, stop, status, logs).
    *   **Required Fix:** Include the standard `systemctl` commands for managing the `m3tal-api.service` as found in the GROUND TRUTH.

7.  **WARNING: Firewall note is absent.**
    *   **Description:** The README does not include a reminder for users to allow port 80 (and potentially 443) in their firewall, which is crucial for Traefik to function publicly.
    *   **Required Fix:** Add a section reminding users to configure their firewall (ufw/iptables) to allow public access to port 80.

8.  **WARNING: Tone contains marketing copy.**
    *   **Description:** Phrases like "M3TAL is a system composed of a Go-based CLI, a Go API daemon, a Python/Flask dashboard..." and "serving as the primary command-line interface for all M3TAL operations" lean towards descriptive marketing rather than direct technical documentation.
    *   **Required Fix:** Rephrase descriptive sentences to be more direct and technical. For example, instead of "serving as the primary command-line interface for all M3TAL operations," simply state "The CLI binary...".

9.  **SUGGESTION: Quick Start section is not ideal.**
    *   **Description:** The "Quick Demo" section outlines two commands (`m3tal dash up` and `m3tal up`). However, it doesn't provide a step-by-step, actionable "quick start" that a new user could follow from installation to seeing a functional dashboard. It's more of a summary of commands.
    *   **Required Fix:** Create a concise "Quick Start" guide that assumes the user has completed the installation and guides them through setting up and accessing the dashboard, covering the `DASHBOARD_EXPOSE_MODE` in a practical way.

---