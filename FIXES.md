## M3TAL README Audit

**Verdict: FAILED**

The README is missing critical information regarding APT installation and deployment lifecycle, which are blockers for successful deployment. Additionally, there are missing details about Traefik routing and service management that warrant warnings.

---

### Issue List:

1.  **BLOCKER: APT installation missing 3-command block**
    *   **Description:** The README does not present the full 3-command sequence (keyring, repo, install) required for APT installation. While the individual commands are present, they are not grouped together as a single, actionable block.
    *   **Required Fix:** Consolidate the three APT installation commands into a single, contiguous code block.

2.  **BLOCKER: Docker dependency missing explicit statement**
    *   **Description:** While the README mentions Docker Engine and Docker Compose V2 are required, it does not explicitly state that both *must* be installed. The phrasing "are strictly REQUIRED and must be installed prior to M3TAL installation and operation" is good, but the explicit mention of *both* components being installed is missing in a clear, single statement.
    *   **Required Fix:** Ensure the README clearly states that both "Docker Engine" and "Docker Compose V2" are required and must be installed.

3.  **BLOCKER: Deployment lifecycle missing explanation of stacks and `/docker` directory usage**
    *   **Description:** The README mentions `m3tal up` and the `/docker` directory but fails to explain that `m3tal up` operates on all `*-compose.yml` files within `/docker/` and that `/docker` is a symlink to `/opt/m3tal/stack/`. This fundamental understanding of how M3TAL deploys services is missing.
    *   **Required Fix:** Add a clear explanation of how M3TAL uses Docker Compose within the `/docker` directory (which is a symlink to `/opt/m3tal/stack/`), and how `m3tal up` orchestrates these files.

4.  **BLOCKER: Traefik routing explanation incomplete**
    *   **Description:** The README states Traefik is the HTTP gateway and mentions it interprets labels. However, it fails to explicitly mention *how* services get exposed (i.e., via labels or dynamic config) and doesn't clearly link the dynamic configuration files to the `/docker/dynamic/` path. The example for exposing a custom service is a good start but needs to be integrated into a clearer overall explanation.
    *   **Required Fix:** Provide a clearer explanation of Traefik's role as an HTTP gateway, emphasizing that services are exposed via Traefik *labels* in their compose files or through *dynamic configuration files*. Explicitly mention that dynamic configuration files are typically placed in `/docker/dynamic/` (or within `/opt/m3tal/stack/dynamic/`).

5.  **WARNING: Port table missing port 8080**
    *   **Description:** The "Port Map" table in the README is missing the entry for port `8080`, which is the host-local port for the M3TAL Go API daemon.
    *   **Required Fix:** Add an entry for port `8080` to the "Port Map" table, detailing its service and access method.

6.  **WARNING: Service management section incomplete**
    *   **Description:** The README mentions `m3tal-api.service` and `systemctl` commands but does not explicitly state that this is the *primary* way to manage the API daemon.
    *   **Required Fix:** Clarify that `systemctl` is the standard and recommended method for managing the `m3tal-api.service`.

7.  **WARNING: Firewall note missing specific port 80 reminder**
    *   **Description:** The "Firewall Considerations" section mentions allowing port `80` if Traefik is used for public access but doesn't explicitly remind users to allow port `80` in their firewall configuration (e.g., `ufw` or `iptables`).
    *   **Required Fix:** Add a specific reminder within the "Firewall Considerations" section to ensure port `80` is allowed in the firewall.

8.  **WARNING: Tone is not purely technical documentation**
    *   **Description:** The "Overview" section uses phrases like "technical details and operational procedures" and "describes the system's architecture, component interactions, deployment mechanisms, and operational guidelines." While informative, it borders on marketing copy rather than a direct, technical introduction.
    *   **Required Fix:** Rephrase the "Overview" to be more direct and technical, focusing on what the document *is* rather than what it *describes*. For example, "This document details the M3TAL system's installation, configuration, and operational procedures."

9.  **SUGGESTION: Quick demo section is not a "quick demo"**
    *   **Description:** The "Quick Demo" section is essentially a scaled-down installation and deployment guide. A true "quick demo" would involve a very simple, single-command or minimal set of commands to show a basic functionality, not a multi-step process.
    *   **Required Fix:** Re-evaluate the "Quick Demo" section. If it's intended to be a quick start, rename it accordingly. If a true "quick demo" (e.g., spinning up a minimal service with one command to showcase functionality) is desired, create a separate, simpler section for it.

---

### Required Fixes Per Issue:

1.  **APT installation missing 3-command block:**
    *   **Fix:** Combine the three APT commands into a single code block:
        ```bash
        curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg
        echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list
        sudo apt update && sudo apt install -y m3tal
        ```

2.  **Docker dependency missing explicit statement:**
    *   **Fix:** Modify the "Prerequisites" section to clearly state: "Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation."

3.  **Deployment lifecycle missing explanation of stacks and `/docker` directory usage:**
    *   **Fix:** Enhance the "Deployment Lifecycle" section with the following:
        "M3TAL orchestrates Docker containers using Docker Compose V2. The `m3tal up` command is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory. The `/docker` directory is a user-facing symlink alias to the canonical stack directory located at `/opt/m3tal/stack/`. Therefore, placing any `*-compose.yml` file into `/docker/` will cause it to be deployed by `m3tal up`."

4.  **Traefik routing explanation incomplete:**
    *   **Fix:** Update the "Traefik Gateway" section. For example:
        "Traefik acts as the primary HTTP gateway. Services are exposed and routed to Traefik in two primary ways:
        *   **Traefik Labels:** By adding `traefik.enable=true` and other relevant `traefik.http.*` labels to a service definition in a Docker Compose file.
        *   **Dynamic Configuration:** Via static configuration files in the `/docker/dynamic/` directory (which maps to `/opt/m3tal/stack/dynamic/`). These files define routers and services, such as routing `api.DOMAIN` to the host-local Go API on port `8080`."
        The example for exposing a custom user service should be presented as an illustration of using labels.

5.  **Port table missing port 8080:**
    *   **Fix:** Add the following row to the "Port Map" table:
        | 8080 | M3TAL API daemon (Go) | Host-local | The internal port the M3TAL API daemon listens on. |

6.  **Service management section incomplete:**
    *   **Fix:** Modify the "Service Management" section to clarify:
        "The M3TAL API daemon operates as a systemd service, `m3tal-api.service`. You can manage its lifecycle using standard systemctl commands, which is the recommended method."

7.  **Firewall note missing specific port 80 reminder:**
    *   **Fix:** Update the "Firewall Considerations" section:
        "If you are using Traefik for public-facing access (i.e., `DASHBOARD_EXPOSE_MODE=traefik` or exposing other services via Traefik), ensure that host port `80` (and `443` if HTTPS is configured) is open in your firewall. For example, with `ufw`: `sudo ufw allow 80/tcp`."

8.  **Tone is not purely technical documentation:**
    *   **Fix:** Rephrase the "Overview" to be more direct and technical. For example:
        "This document provides technical details and operational procedures for the M3TAL system, covering installation, configuration, component interactions, and deployment."

9.  **Quick demo section is not a "quick demo":**
    *   **Fix:** Rename the "Quick Demo" section to "Quick Start". If a true, minimal "quick demo" is desired, create a separate section for it that shows a single, impactful command for basic functionality.