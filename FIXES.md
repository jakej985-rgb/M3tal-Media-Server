**Verdict: FAILED**

**Reason:** The README is missing critical information regarding Docker Compose V2, the specific commands for `m3tal up`, and incomplete documentation on Traefik routing for user-defined services.

---

**Issues:**

1.  **BLOCKER: APT installation:** The README does not explicitly show the 3-command keyring+repo+install block as required by the audit criteria.
    *   **Required Fix:** The README should clearly present the three commands for APT installation: `curl ... | sudo gpg ...`, `echo ... | sudo tee ...`, and `sudo apt update && sudo apt install -y m3tal`.
    *   **Status:** **Resolved** (The provided README already has this correct.)

2.  **BLOCKER: Docker dependency:** The README fails to state that Docker Engine + Docker Compose V2 are required.
    *   **Required Fix:** Add a clear statement that both Docker Engine and Docker Compose V2 are mandatory dependencies.
    *   **Status:** **Resolved** (The provided README already has this correct.)

3.  **BLOCKER: Deployment lifecycle:** The README's explanation of how stacks work is incomplete. It mentions `m3tal up` and the `/docker` directory but doesn't fully explain that `/docker` is a symlink to `/opt/m3tal/stack/` or how adding new compose files into this directory is handled by `m3tal up`.
    *   **Required Fix:** Clarify that `/docker` is a symlink to `/opt/m3tal/stack/` and that `m3tal up` processes all `*-compose.yml` files within this directory. Explicitly state that placing a new `my-stack-compose.yml` in `/docker/` (which points to `/opt/m3tal/stack/`) will cause `m3tal up` to deploy it.
    *   **Status:** **Partially Resolved** (The README states `/docker` is a symlink to `/opt/m3tal/stack/` and that `m3tal up` operates on files in `/docker`. However, it doesn't explicitly link placing a file in `/docker` to `m3tal up` processing it, which is a key part of the lifecycle.)

4.  **BLOCKER: Traefik routing:** The README's explanation of Traefik routing is insufficient. It mentions Traefik as the HTTP gateway and that services are exposed via labels or dynamic config, but it doesn't clearly explain how Traefik discovers services managed by Docker Compose (specifically through the `proxy` network and `traefik.enable=true` labels). It also needs to clarify the role of Traefik in the `DASHBOARD_EXPOSE_MODE=traefik` scenario.
    *   **Required Fix:**
        *   Explain that Traefik, by default, only discovers services with `traefik.enable=true` labels and that services must be on the `proxy` network.
        *   Clarify how Traefik routes to the dashboard when `DASHBOARD_EXPOSE_MODE=traefik` by referencing the labels applied to the dashboard service.
        *   Ensure the example for exposing a custom user service is complete and demonstrates the necessary labels and network configuration.
    *   **Status:** **Partially Resolved** (The README mentions Traefik labels for exposing services, and the custom service example is provided. However, it could be more explicit about Traefik's discovery mechanism via Docker labels and the `proxy` network for services it manages directly.)

5.  **WARNING: Port table:** The README lists ports 80, 8080, 8081, and 8082, but the description for port 80 is incomplete, and the description for port 8082 is also lacking detail about its access method based on `DASHBOARD_EXPOSE_MODE`.
    *   **Required Fix:**
        *   For port 80, specify "Traefik HTTP (public)" as the service.
        *   For port 8082, clarify that access depends on `DASHBOARD_EXPOSE_MODE` (direct port or via Traefik).
    *   **Status:** **Partially Resolved** (The port table is present, but the descriptions need enhancement as noted above.)

6.  **WARNING: Service management:** The README mentions systemctl for managing m3tal-api.service.
    *   **Required Fix:** None needed, this criterion is met.
    *   **Status:** **Resolved** (The provided README already has this correct.)

7.  **WARNING: Firewall note:** The README reminds users to allow port 80 in ufw/iptables.
    *   **Required Fix:** None needed, this criterion is met.
    *   **Status:** **Resolved** (The provided README already has this correct.)

8.  **WARNING: Tone:** The writing style is generally technical and appropriate.
    *   **Required Fix:** None needed.
    *   **Status:** **Resolved** (The provided README already has this correct.)

9.  **SUGGESTION: Quick demo:** The README has a working Quick Start section.
    *   **Required Fix:** None needed, this criterion is met.
    *   **Status:** **Resolved** (The provided README already has this correct.)

---

**Summary of BLOCKER Issues and Required Fixes:**

*   **Deployment Lifecycle (Issue 3):**
    *   **Explanation:** The README states that `/docker` is a symlink to `/opt/m3tal/stack/` and that `m3tal up` operates on files in `/docker`. However, it needs to explicitly state that adding a new `*-compose.yml` file into `/docker/` will lead to `m3tal up` deploying it as part of the overall stack. This is a fundamental aspect of how users extend the system.
    *   **Fix:** In the "Adding a New Stack" section, explicitly state: "Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. This file will be automatically discovered and included by `m3tal up` when it processes all `*-compose.yml` files in the stack directory."

*   **Traefik Routing (Issue 4):**
    *   **Explanation:** The README needs to be more explicit about Traefik's service discovery and routing mechanisms for Docker services. This includes:
        *   How Traefik finds services (Docker labels, `traefik.enable=true`, `proxy` network).
        *   The specific mechanism for `DASHBOARD_EXPOSE_MODE=traefik` (i.e., Traefik labels on the dashboard service).
        *   Ensuring the user service example is clear and functional.
    *   **Fix:**
        *   In the "Traefik Gateway" section, add: "Traefik, configured with the Docker provider, automatically discovers services with the `traefik.enable=true` label within the `proxy` Docker network. To expose a service, it must be connected to this network and have appropriate routing labels defined."
        *   In the "Dashboard Access" section, under "2. Traefik Mode," add: "When `DASHBOARD_EXPOSE_MODE=traefik`, Traefik is configured via labels on the `m3tal-dashboard` service definition (within its compose file) to route traffic from `dash.DOMAIN` to the dashboard's internal port `8082`."
        *   In the "Exposing a Custom User Service via Traefik" example, ensure the `networks` section clearly defines the `proxy` network as external and explicitly states the service must be connected to it.

*   **Port Table (Issue 5):**
    *   **Explanation:** The descriptions for ports 80 and 8082 are vague and could be more informative.
    *   **Fix:** Update the "Port Map" table:
        *   For Port 80: Change "Traefik HTTP (public)" to "Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik."
        *   For Port 8082: Change "M3TAL Dashboard" to "M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`."