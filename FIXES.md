## DocCritic Audit Report: M3TAL README.md

**Verdict:** FAILED (Multiple BLOCKER issues identified)

---

**Issue List:**

1.  **BLOCKER:** APT installation commands are present but the README states "M3TAL is installed via APT — not built from source." however, the ground truth provides the exact three command block required. The README is missing the 3-command keyring+repo+install block.
    *   **Required Fix:** Add the following 3-command block to the "Installation" section:
        ```bash
        curl -fsSL https://jakej985-rgb.github.io/m3tal-apt-key/public.key | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg
        echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-apt-key stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list
        sudo apt update && sudo apt install -y m3tal
        ```

2.  **BLOCKER:** Docker dependency is missing a clear statement that Docker Engine + Docker Compose V2 are required. While it mentions Docker internally, it does not explicitly state the prerequisites.
    *   **Required Fix:** Add a clear statement in the "Prerequisites" section: "Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation. M3TAL internally orchestrates Docker containers using Docker Engine and Docker Compose V2."

3.  **BLOCKER:** Deployment lifecycle explanation is incomplete. The README mentions `m3tal up` and the `/docker` directory, but it fails to explicitly explain how stacks work, the role of `/docker` as a symlink to `/opt/m3tal/stack/`, and the mechanism for adding new compose files to the deployment.
    *   **Required Fix:** Enhance the "Deployment Lifecycle" section to include:
        *   A clear statement that `/docker` is a symlink to `/opt/m3tal/stack/`.
        *   An explicit explanation that `m3tal up` runs `docker compose` across all `*-compose.yml` files in `/docker/`.
        *   A more detailed explanation on how to add new compose files (e.g., placing them directly in `/docker/` and `m3tal up` will pick them up).

4.  **BLOCKER:** Traefik routing explanation is insufficient. The README mentions Traefik as the HTTP gateway and its use of labels and dynamic config for routing, but it lacks a clear, consolidated explanation of how services get exposed via Traefik, including the role of the `proxy` network and the `traefik.enable=true` label.
    *   **Required Fix:** In the "Traefik Gateway" section, add a clear summary: "Traefik acts as the HTTP gateway and is configured to route traffic to Docker services based on `traefik` labels defined in their respective Compose files. For a service to be exposed by Traefik, it must have `traefik.enable=true` in its labels and be connected to the `proxy` Docker network. Traefik also utilizes dynamic configuration files from `/etc/traefik/dynamic/` for advanced routing."

5.  **WARNING:** The port table is missing port 8081.
    *   **Required Fix:** Add port 8081 to the "Port Map" table with its corresponding service and description.

6.  **WARNING:** The README mentions `m3tal-api.service` but does not explicitly mention `systemctl` for managing it.
    *   **Required Fix:** In the "Service Management" section, explicitly state: "The M3TAL API daemon is managed as a systemd service named `m3tal-api.service`. Use standard `systemctl` commands for its operation:".

7.  **WARNING:** The firewall note is present but could be more specific about `ufw/iptables`.
    *   **Required Fix:** In the "Firewall Considerations" section, amend the sentence to read: "If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall (e.g., `ufw allow 80/tcp` or equivalent iptables rules)."

8.  **WARNING:** Some sections, particularly the "Core Components" and the introductory paragraphs of "Deployment Lifecycle" and "Dashboard Access", lean towards descriptive language that could be perceived as marketing copy rather than purely technical documentation.
    *   **Required Fix:** Review these sections and rephrase to be more direct and technical, focusing on *what* the component does and *how* it functions rather than its benefits. For example, instead of "A unified Go binary serving as the single entrypoint for all M3TAL operations," use "The CLI binary (`/usr/bin/m3tal`) is a Go executable used for all M3TAL operations."

9.  **SUGGESTION:** The "Quick Demo" section is present but could be improved by being more explicit about the commands and their immediate impact.
    *   **Required Fix:** Enhance the "Quick Demo" section to be more direct about the commands: "To quickly get the M3TAL dashboard running in `local` mode (default), execute `m3tal dash up`. This command specifically manages the `m3tal-dashboard` container. To deploy all other M3TAL managed stacks, including Traefik and any user-defined Docker Compose files placed in `/docker/`, execute `m3tal up`."