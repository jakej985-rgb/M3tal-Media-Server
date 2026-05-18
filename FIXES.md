**Verdict:** FAILED

**Reason:** The README is missing critical information regarding Traefik's service discovery mechanism and has several missing or incorrectly documented port mappings and service management details.

---

**Numbered Issue List:**

1.  **BLOCKER: Traefik routing explanation is incomplete.**
    *   **Issue:** The README states Traefik discovers services by interpreting labels and uses a file provider for dynamic configuration. However, it fails to explicitly state *how* services get exposed. While an example is provided for a custom service, it doesn't clearly articulate that services *must* have `traefik.enable=true` and associated labels to be routed by Traefik, or that dynamic configuration files are another method. The ground truth indicates Traefik uses `providers.docker.exposedByDefault: false` meaning explicit labeling is required.
    *   **Required Fix:** Clarify that services must have `traefik.enable=true` and relevant labels in their compose files to be exposed by Traefik. Explicitly mention that Traefik does *not* expose services by default.

2.  **BLOCKER: Missing explanation of Docker Compose V2 usage.**
    *   **Issue:** The README states "M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose`". While this implies V2, it doesn't explicitly state that Docker Engine *and* Docker Compose V2 are *required* prerequisites as per the ground truth.
    *   **Required Fix:** In the "Prerequisites" section, explicitly state that "Docker Engine + Docker Compose V2 are required."

3.  **BLOCKER: Missing explanation of deployment lifecycle for stacks.**
    *   **Issue:** The README explains how `m3tal up` works with `*-compose.yml` files in `/docker/`. It also mentions `/opt/m3tal/stack/` as the canonical source. However, it does not fully explain how adding *new* compose files to `/docker/` is the mechanism to deploy additional stacks. The section "Adding a New Stack" correctly describes this, but it's under "Deployment Lifecycle" and could be more integrated into the core explanation of how stacks are managed. The ground truth emphasizes that `m3tal up` runs `docker compose` across ALL `*-compose.yml` files in `/docker/`.
    *   **Required Fix:** Integrate the concept that `m3tal up` deploys *all* compose files found in `/docker/` as distinct stacks. Ensure the explanation clearly links new compose files in `/docker/` to new deployments managed by `m3tal up`.

4.  **WARNING: Port map table is incomplete and has inaccuracies.**
    *   **Issue:** The README's "Port Map" section lists ports 80, 8080, 8081, and 8082.
        *   Port 8080 is described as "M3TAL API daemon (Go) - Host-local". Ground truth confirms this.
        *   Port 8082 is described as "M3TAL Dashboard - Direct port (local mode) or via Traefik (traefik mode)". Ground truth shows it is directly bound in `m3tal-compose.local.yml` and the container internally listens on 8082.
        *   Port 80 is described as "Traefik HTTP entry point - Public (when Traefik mode is active)". Ground truth confirms this.
        *   Port 8081 is described as "Traefik dashboard - Host-local only". Ground truth confirms this (`127.0.0.1:8081:8080`).
    *   However, the ground truth also indicates that Traefik itself uses port 80 for its `web` entrypoint, and the `routing-compose.yml` shows Traefik listening on `:80`. The README's port map *should* list port 80 as Traefik's HTTP entry point for general traffic, not just "public (when Traefik mode is active)". Furthermore, the internal Traefik dashboard port is 8080, which is then mapped to `127.0.0.1:8081` on the host. The table implies the host-local port is 8081, which is correct, but it should be explicitly stated that this is Traefik's *dashboard* port.
    *   **Required Fix:**
        *   Clarify that port 80 is the public HTTP entry point for Traefik.
        *   Ensure all listed ports (80, 8080, 8081, 8082) are explicitly mentioned.
        *   The table is mostly correct according to the ground truth provided, but a slight rephrasing for clarity around Traefik's ports is warranted.

5.  **WARNING: Service management details are incomplete.**
    *   **Issue:** The README mentions `systemctl status m3tal-api`, `systemctl restart m3tal-api`, and `journalctl -u m3tal-api -f`. While these are correct, the ground truth provides slightly more detail on how the API daemon is managed by systemd. The README implies that `m3tal-api.service` is the *only* service managed by systemd, which seems to be the case based on the provided ground truth.
    *   **Required Fix:** No immediate fix required if the ground truth confirms only `m3tal-api.service` is systemd managed and the provided commands are sufficient. The current documentation is adequate for this item.

6.  **WARNING: Firewall note is not prominent enough.**
    *   **Issue:** The firewall note about allowing port 80 is placed under the "Installation" section. While it's present, it might be missed by users focusing on post-installation operations.
    *   **Required Fix:** Consider reiterating the firewall note in a more prominent location, perhaps in the "Traefik Gateway" section, or as a distinct "Important Notes" section.

7.  **WARNING: Tone is not strictly marketing copy, but could be more direct.**
    *   **Issue:** Phrases like "M3TAL Ecosystem Documentation" and "M3TAL Ecosystem" and "unified entrypoint for all M3TAL operations" lean slightly towards marketing language. The document should focus on purely technical and operational descriptions.
    *   **Required Fix:** Rephrase introductory and descriptive sentences to be more direct and less promotional. For example, "M3TAL System Documentation" instead of "M3TAL Ecosystem Documentation."

8.  **SUGGESTION: Quick demo section is present and functional.**
    *   **Issue:** The "Quick Demo" section provides functional `m3tal dash up` and `m3tal up` commands. This meets the requirement.
    *   **Required Fix:** No action needed.

---

**Detailed Fixes per Issue:**

1.  **Traefik routing explanation:**
    *   **Location:** "Traefik Gateway" section.
    *   **Fix:** Add the following clarification:
        > "Traefik automatically discovers and routes traffic to Docker services by interpreting Traefik labels defined within their Docker Compose service definitions. **Crucially, services are not exposed by Traefik by default and require `traefik.enable=true` along with other relevant labels to be discoverable and routable.**"

2.  **Docker Compose V2 requirement:**
    *   **Location:** "Prerequisites" section.
    *   **Fix:** Modify the existing sentence to be more explicit:
        > "Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL relies on Docker Compose V2 for its internal orchestration and uses Docker Engine + Docker Compose V2 internally."
        to:
        > "Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2."

3.  **Deployment lifecycle for stacks:**
    *   **Location:** "Deployment Lifecycle" section.
    *   **Fix:** Enhance the explanation of `m3tal up` to emphasize its role in deploying all compose files:
        > "M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on **all `*-compose.yml` files located within the `/docker/` directory, effectively deploying each as an independent stack.**"
        And ensure "Adding a New Stack" clearly states this:
        > "To deploy a new Docker Compose stack within the M3TAL ecosystem:
        >
        > 1. Place your Docker Compose file (e.g., `my-stack-compose.yml`) directly into the `/docker/` directory. **This file will be automatically included by `m3tal up`.**"

4.  **Port map table:**
    *   **Location:** "Port Map" section.
    *   **Fix:** Modify the table and surrounding text to be more precise:
        > **Port Map**
        >
        > The following table lists the primary network ports utilized by the M3TAL system:
        >
        > | Port | Service | Access | Description |
        > |------|---------|--------|-------------|
        > | 80 | Traefik HTTP entry point | Public | The public-facing HTTP port for services exposed via Traefik. |
        > | 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |
        > | 8081 | Traefik dashboard | Host-local only | The internal Traefik dashboard port, accessible only from the host machine. |
        > | 8082 | M3TAL Dashboard | Direct port (local mode) or via Traefik (traefik mode) | The port the M3TAL Dashboard container listens on internally. Access method depends on `DASHBOARD_EXPOSE_MODE`. |
        >
        > **Note:** These are the primary M3TAL-managed ports. User-added Docker Compose stacks may expose additional ports on their own, depending on their configurations.

5.  **Service management:**
    *   **Location:** "Service Management" section.
    *   **Fix:** No changes required if the ground truth confirms only `m3tal-api.service` is systemd managed. The current documentation is sufficient.

6.  **Firewall note prominence:**
    *   **Location:** Add a "Firewall Considerations" section or add a note in the "Traefik Gateway" section.
    *   **Fix:** Consider adding a section like this:
        > **Firewall Considerations**
        >
        > If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp`).

7.  **Tone:**
    *   **Location:** Throughout the document, especially the introduction.
    *   **Fix:** Replace marketing-oriented phrasing with direct technical descriptions.
        *   Change "M3TAL Ecosystem Documentation" to "M3TAL System Documentation."
        *   Change "unified entrypoint for all M3TAL operations" to "single command-line interface for M3TAL operations."

8.  **Quick demo:**
    *   **Location:** "Quick Demo" section.
    *   **Fix:** No action needed. The section is present and functional.