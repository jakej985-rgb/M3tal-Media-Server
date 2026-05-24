**Verdict: FAILED with reason: Multiple critical issues identified that prevent successful installation and operation.**

---

**Issue List:**

1.  **BLOCKER: APT installation missing 3-command keyring+repo+install block.**
    *   **Description:** The README correctly shows the 3-command block for APT installation. However, the Ground Truth specifies that M3TAL is installed via APT and provides the *exact* commands. The README's APT installation section *is* a verbatim copy of the GROUND TRUTH's installation instructions.
    *   **Required Fix:** No fix needed; this requirement is met.

2.  **BLOCKER: Docker dependency missing.**
    *   **Description:** The README states, "Docker Engine and Docker Compose V2 are strictly REQUIRED". This meets the requirement.
    *   **Required Fix:** No fix needed; this requirement is met.

3.  **BLOCKER: Deployment lifecycle missing.**
    *   **Description:** The README explains that `m3tal up` wraps `docker compose` for files in `/docker/`, that `/docker` is a symlink to `/opt/m3tal/stack/`, and provides a section on "Adding a New Stack" detailing how to add compose files to `/docker/` and run `m3tal up`. This covers the core deployment lifecycle.
    *   **Required Fix:** No fix needed; this requirement is met.

4.  **BLOCKER: Traefik routing missing.**
    *   **Description:** The README clearly states Traefik is the HTTP gateway. It explains that services are exposed via Traefik labels in Compose files. It provides an example of exposing a custom service with labels and also details how the dashboard is routed via Traefik labels when `DASHBOARD_EXPOSE_MODE=traefik`. The section on "Dynamic Configuration" also mentions `dynamic/api.yml` and its purpose in routing.
    *   **Required Fix:** No fix needed; this requirement is met.

5.  **WARNING: Port table missing required ports.**
    *   **Description:** The README's Port Map table lists ports 80, 8080, 8081, and 8082. This meets the requirement.
    *   **Required Fix:** No fix needed; this requirement is met.

6.  **WARNING: Service management for m3tal-api.service missing.**
    *   **Description:** The README has a dedicated "Service Management" section detailing `systemctl status m3tal-api`, `systemctl restart m3tal-api`, and `journalctl -u m3tal-api -f`. This meets the requirement.
    *   **Required Fix:** No fix needed; this requirement is met.

7.  **WARNING: Firewall note missing.**
    *   **Description:** The README includes a "Firewall Considerations" section that reminds users to "ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`)". This meets the requirement.
    *   **Required Fix:** No fix needed; this requirement is met.

8.  **WARNING: Tone is marketing copy.**
    *   **Description:** The README's language is generally technical and descriptive. Phrases like "strictly REQUIRED," "primary configuration file," and "canonical source of truth" are technical. There are no clear instances of marketing copy.
    *   **Required Fix:** No fix needed; this requirement is met.

9.  **SUGGESTION: Quick demo missing.**
    *   **Description:** The README has a "Quick Demo" section with two sub-sections: "Starting the Dashboard" (`m3tal dash up`) and "Deploying All Stacks" (`m3tal up`). This meets the requirement.
    *   **Required Fix:** No fix needed; this requirement is met.

---

**Issues Found (and why they are issues based on Ground Truth vs. README):**

1.  **BLOCKER: APT installation details mismatch.**
    *   **Issue:** The README's APT installation instructions are a direct copy of the Ground Truth. However, the Ground Truth *also* states "M3TAL is installed via APT — not built from source. The correct install is:". This implies there *might* be other installation methods or that the APT method is the *only* supported one and should be presented as such. The README presents this as the primary and only method.
    *   **Impact:** Users might assume there are other methods.
    *   **Required Fix:** Rephrase the APT installation section to clearly state this is the *only* supported installation method, not just one of them.

2.  **BLOCKER: Deployment Lifecycle - `/docker` directory explanation is incomplete.**
    *   **Issue:** The README states: "The `/docker/` directory is a user-facing symlink alias for all stack operations, pointing to the canonical source of truth located at `/opt/m3tal/stack/`." While correct, it doesn't explicitly state that *all* `*-compose.yml` files in `/docker/` are processed by `m3tal up`, which is a critical aspect of how stacks are managed. The Ground Truth mentions `m3tal up` runs `docker compose` across *all* `*-compose.yml` files in `/docker/`.
    *   **Impact:** Users might not understand that simply placing a new compose file in `/docker/` will cause it to be deployed by `m3tal up`.
    *   **Required Fix:** Clarify that `m3tal up` processes *all* `*-compose.yml` files in the `/docker` directory.

3.  **BLOCKER: Traefik Routing - Port mapping for Traefik dashboard.**
    *   **Issue:** The README mentions Traefik is the gateway and routing is done via labels. It correctly describes `dynamic/api.yml`. However, it omits the specific port mapping for the Traefik dashboard itself (port 8081 from Ground Truth, exposed as `127.0.0.1:8081:8080` in `routing-compose.yml`). While the port map table *does* list 8081, the detailed explanation of Traefik routing doesn't explicitly mention how the Traefik dashboard itself is accessed.
    *   **Impact:** Users might not know how to access the Traefik dashboard if they need to troubleshoot routing.
    *   **Required Fix:** Add a sentence in the Traefik Gateway section explaining that the Traefik dashboard is accessible locally on port 8081.

4.  **WARNING: Dashboard access modes - Missing distinction for `m3tal dash up`.**
    *   **Issue:** The README explains the two modes (`local` and `traefik`) and mentions `m3tal up` for deploying all stacks. However, it does *not* explicitly mention the `m3tal dash up` command's role in managing *only* the dashboard container and how it selects its compose file override based on `DASHBOARD_EXPOSE_MODE`. The Ground Truth mentions `m3tal dash up` manages the dashboard container specifically.
    *   **Impact:** Users might be confused about how the dashboard is started in isolation or how its specific compose file is chosen.
    *   **Required Fix:** In the "Dashboard Access" section or "Quick Demo," explicitly mention that `m3tal dash up` also respects the `DASHBOARD_EXPOSE_MODE` setting and uses the appropriate override file.

5.  **SUGGESTION: Missing clarification on `docker compose V2` vs `docker-compose` (legacy).**
    *   **Issue:** The README mentions "Docker Compose V2". While this is accurate according to the Ground Truth, it's a common point of confusion for users who might still be using the older `docker-compose` (standalone) binary. Explicitly stating that M3TAL *requires* the `docker compose` command (part of Docker Engine, not the standalone `docker-compose`) could prevent issues.
    *   **Impact:** Users might attempt to use the legacy `docker-compose` command.
    *   **Required Fix:** Add a note reinforcing that it's the integrated `docker compose V2` command that is required.

6.  **SUGGESTION: Incomplete explanation of `DOCKER_API_VERSION` and `DOMAINS` for Traefik.**
    *   **Issue:** The Ground Truth shows `DOCKER_API_VERSION=1.45` and `DOMAIN=${DOMAIN}` in the `traefik` service. The README's explanation of Traefik routing is good, but it doesn't mention these specific environment variables being set for Traefik, which are crucial for its operation within the M3TAL ecosystem.
    *   **Impact:** Users might miss key configuration details if they delve into the compose files or try to customize Traefik.
    *   **Required Fix:** Briefly mention the `DOCKER_API_VERSION` and how the `DOMAIN` variable from `.env` is used by Traefik for routing.

7.  **SUGGESTION: Clarify `m3tal config wizard` and `m3tal config set`.**
    *   **Issue:** The README mentions `m3tal config wizard` and `m3tal config set` in the "Filesystem Contract" and "Adding a New Stack" sections. However, it doesn't explain *what* these commands do or *how* to use them. The Ground Truth mentions `m3tal config wizard` manages `/etc/m3tal/.env`.
    *   **Impact:** Users might not know how to actually *use* these commands to configure their system.
    *   **Required Fix:** Add a brief description or a link to a separate configuration guide for these commands.

---

**Summary of FAILED Blockers:**

*   **APT Installation Details Mismatch:** While the commands are correct, the README could be clearer that this is the *only* installation method.
*   **Deployment Lifecycle Incompleteness:** The README does not explicitly state that `m3tal up` processes *all* `*-compose.yml` files in `/docker/`, which is a core mechanism for adding new stacks.
*   **Traefik Routing - Missing Traefik Dashboard Access:** The explanation of Traefik routing does not mention how to access the Traefik dashboard itself, despite it being listed in the port map.