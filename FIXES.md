**Verdict: FAILED**

The README is missing critical information regarding Docker installation and Traefik routing configuration, which are fundamental to deploying and operating M3TAL.

---

### Issues:

1.  **BLOCKER**: APT installation: README MUST show the 3-command keyring+repo+install block.
    *   **Issue**: The README *does* provide the correct 3-command APT installation block.
    *   **Classification**: Not Applicable (Present)

2.  **BLOCKER**: Docker dependency: README MUST state that Docker Engine + Docker Compose V2 are required.
    *   **Issue**: The README states that "Docker Engine and Docker Compose V2 are strictly REQUIRED," but it does not specify *how* to install them. Users need clear instructions for installing Docker Engine and Docker Compose V2.
    *   **Classification**: BLOCKER
    *   **Required Fix**: Add a section or bullet points detailing the installation steps for Docker Engine and Docker Compose V2.

3.  **BLOCKER**: Deployment lifecycle: README MUST explain how stacks work (/docker dir, m3tal up, adding new compose files).
    *   **Issue**: The README *does* explain that `/docker/` is a symlink to `/opt/m3tal/stack/`, that `m3tal up` runs `docker compose` on files in `/docker/`, and how to add new compose files.
    *   **Classification**: Not Applicable (Present)

4.  **BLOCKER**: Traefik routing: README MUST explain that Traefik is the HTTP gateway and how services get exposed (labels or dynamic config).
    *   **Issue**: The README mentions Traefik as the reverse proxy and that services are exposed via labels. It provides a good example for a custom service. However, it *does not* clearly explain how Traefik itself is configured (e.g., the existence and purpose of `traefik.yml` and `dynamic/` directory) or how the internal M3TAL services (like the API itself) are exposed *via* Traefik. The ground truth shows Traefik routes `api.DOMAIN` to `host.docker.internal:8080` and the dashboard via `dash.DOMAIN`. This internal routing logic is not fully conveyed.
    *   **Classification**: BLOCKER
    *   **Required Fix**: Explicitly state that Traefik uses `traefik.yml` for static configuration and the `dynamic/` directory for dynamic configuration. Clarify how Traefik routes internal M3TAL services (e.g., API and dashboard in Traefik mode) by referencing the relevant configuration or labels.

5.  **WARNING**: Port table: README SHOULD list ports 80, 8080, 8081, 8082.
    *   **Issue**: The README *does* list ports 80, 8080, 8081, and 8082 in its "Port Map" section.
    *   **Classification**: Not Applicable (Present)

6.  **WARNING**: Service management: README SHOULD mention systemctl for managing m3tal-api.service.
    *   **Issue**: The README *does* mention `systemctl` commands for managing `m3tal-api.service`.
    *   **Classification**: Not Applicable (Present)

7.  **WARNING**: Firewall note: README SHOULD remind users to allow port 80 in ufw/iptables.
    *   **Issue**: The README *does* provide a firewall reminder for port 80.
    *   **Classification**: Not Applicable (Present)

8.  **WARNING**: Tone: Flag as WARNING if the writing is marketing copy rather than technical documentation.
    *   **Issue**: The tone is generally technical and appropriate for documentation. There are no overt marketing phrases.
    *   **Classification**: Not Applicable (Appropriate Tone)

9.  **SUGGESTION**: Quick demo: README SHOULD have a working Quick Start section.
    *   **Issue**: The README *does* have a "Quick Demo" section with `m3tal dash up`. While this demonstrates starting a single component, a true "Quick Start" would ideally guide a user through a minimal functional deployment (e.g., setting up initial config and running `m3tal up` to see the dashboard accessible).
    *   **Classification**: SUGGESTION
    *   **Required Fix**: Enhance the "Quick Demo" to include steps for initial configuration (e.g., `m3tal config wizard`) and then a full `m3tal up` to demonstrate a basic working system accessible via its default mode.

---