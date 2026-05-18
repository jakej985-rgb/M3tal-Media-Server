## M3TAL README Audit

**Verdict: FAILED**

The README fails to provide critical information required for installation and operation, leading to a BLOCKER classification. It also contains a significant amount of marketing copy and lacks essential operational details.

---

### Issues:

1.  **BLOCKER**: APT installation instructions are missing.
    *   **Reason**: The README does not provide the 3-command keyring+repo+install block necessary for APT-based installation as described in the Ground Truth.
    *   **Required Fix**: Add the following block to the README, likely in an "Installation" section:
        ```bash
        curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg
        echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list
        sudo apt update && sudo apt install -y m3tal
        ```

2.  **BLOCKER**: Docker dependency is not explicitly stated.
    *   **Reason**: The README does not explicitly mention that Docker Engine and Docker Compose V2 are required.
    *   **Required Fix**: Add a clear statement in the README, perhaps near the deployment section, such as: "M3TAL requires Docker Engine and Docker Compose V2 to be installed and running on your system."

3.  **BLOCKER**: Deployment lifecycle and stack management are not explained.
    *   **Reason**: The README lacks explanation on how Docker stacks work, the purpose of the `/docker` directory (or its symlink `/opt/m3tal/stack/`), how `m3tal up` orchestrates services using `*-compose.yml` files, and how to add new compose files.
    *   **Required Fix**: Add a section detailing the deployment lifecycle:
        *   Explain that M3TAL uses Docker Compose to manage services.
        *   Clarify that user-managed compose files should reside in `/docker/` (which symlinks to `/opt/m3tal/stack/`).
        *   Describe that `m3tal up` will read all `*-compose.yml` files in `/docker/` and pass them to `docker compose`.
        *   Mention `m3tal dash up` for managing the dashboard specifically.

4.  **BLOCKER**: Traefik routing explanation is missing.
    *   **Reason**: The README does not explain that Traefik is the HTTP gateway or how services are exposed (e.g., via labels in compose files).
    *   **Required Fix**: Add a section explaining:
        *   Traefik acts as the primary HTTP gateway.
        *   Services are exposed to Traefik using Docker labels (as seen in `m3tal-compose.traefik.yml` and `routing-compose.yml`).
        *   Mention how Traefik's static and dynamic configuration is handled (e.g., `traefik.yml`, `dynamic/` directory).

5.  **WARNING**: Port table is incomplete.
    *   **Reason**: The README does not list the standard ports expected by users: 80, 8080, 8081, 8082. It only mentions 8081 for Traefik and 8082 as a default for the dashboard (via a variable).
    *   **Required Fix**: Add a comprehensive port table, including the expected ports and their purposes, matching the Ground Truth:
        *   **Port 80**: Traefik HTTP (public access)
        *   **Port 8080**: M3TAL Go API (internal, potentially exposed via Traefik)
        *   **Port 8081**: Traefik Dashboard (host-local access)
        *   **Port 8082**: M3TAL Dashboard container (host-local access if `DASHBOARD_EXPOSE_MODE=local`)

6.  **WARNING**: Service management details are missing.
    *   **Reason**: The README does not mention `systemctl` for managing the `m3tal-api.service`.
    *   **Required Fix**: Add a brief note about managing the M3TAL API service using `systemctl`, providing common commands like `systemctl status m3tal-api`, `systemctl restart m3tal-api`, and `journalctl -u m3tal-api -f`.

7.  **WARNING**: Firewall reminder is absent.
    *   **Reason**: The README does not remind users to configure their firewall (e.g., ufw, iptables) to allow incoming traffic on port 80.
    *   **Required Fix**: Add a note advising users to ensure port 80 is open in their firewall configuration.

8.  **WARNING**: Tone includes marketing copy.
    *   **Reason**: Phrases like "Autonomous, self-healing media automation platform" are marketing-oriented and not technical documentation. The "Architecture (Source Services)" section also lists internal Go packages which is not typical for end-user documentation.
    *   **Required Fix**: Refine the language to be purely technical and focus on what a user needs to know to install, configure, and operate the system. Remove the detailed Go package breakdown and focus on functional components.

9.  **SUGGESTION**: Quick Start section is missing.
    *   **Reason**: The README lacks a "Quick Start" or "Getting Started" section that guides a user through a minimal, functional deployment.
    *   **Required Fix**: Add a "Quick Start" section. This should ideally guide a user through:
        *   Prerequisites (Docker).
        *   The APT installation commands.
        *   Initial configuration (e.g., `m3tal config wizard` if it exists, or setting up the `.env` file).
        *   Running `m3tal up`.
        *   Accessing the dashboard.

---

**Note**: The current deployment command `python m3tal.py init` is unclear and does not align with the APT installation described in the Ground Truth. This needs clarification or correction. The README also lists services and paths that are implementation details rather than user-facing operational information.