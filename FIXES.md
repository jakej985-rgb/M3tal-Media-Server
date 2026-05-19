## M3TAL System Documentation Audit Report

**Verdict: FAILED - Critical Gaps Identified**

The README.md is a good starting point, but it contains several critical omissions and areas that could lead to confusion for users attempting to deploy and operate the M3TAL system.

---

### Issue List:

1.  **BLOCKER: APT installation commands missing.**
    *   **Description:** The README mentions installation via APT but fails to provide the complete, three-command sequence (keyring, repo, install) required for users to successfully add the M3TAL repository and install the package.
    *   **Required Fix:** Provide the exact `curl | gpg`, `echo | tee`, and `apt update && apt install` commands as shown in the Ground Truth.

2.  **BLOCKER: Docker dependency not explicitly stated.**
    *   **Description:** While the README mentions Docker Engine and Docker Compose V2 in the prerequisites, it does not explicitly state that *both* are required, and it fails to specify that Docker Compose V2 is the *required* version.
    *   **Required Fix:** Clearly state that "Docker Engine + Docker Compose V2 are required" in the prerequisites section.

3.  **BLOCKER: Deployment lifecycle is incomplete.**
    *   **Description:** The README describes `m3tal up` but does not fully explain how stacks work, specifically mentioning the `/docker` directory but not its relationship to `/opt/m3tal/stack/` as the canonical source of truth, nor the `m3tal up` command's behavior of orchestrating *all* `*-compose.yml` files. The explanation of adding new compose files is also a bit vague regarding how they are picked up.
    *   **Required Fix:** Explicitly state that `/docker` is a symlink to `/opt/m3tal/stack/`, that `m3tal up` processes *all* `*-compose.yml` files in `/docker/`, and clarify that placing a new file in `/docker/` will make it part of the `m3tal up` deployment.

4.  **BLOCKER: Traefik routing explanation is insufficient.**
    *   **Description:** The README mentions Traefik as the HTTP gateway and that services are exposed via labels or dynamic config. However, it lacks a clear, step-by-step explanation of *how* services are exposed, relying on abstract statements and a single dynamic config example without fully contextualizing how it applies to user services. The critical `traefik.enable=true` label is not explicitly called out as a requirement for services to be discovered.
    *   **Required Fix:** Clearly state that services are exposed via Traefik labels (e.g., `traefik.enable=true`) and mention dynamic configuration files for more complex routing. Provide a clear example of how to expose a custom user service using Traefik labels, as demonstrated in the Ground Truth.

5.  **WARNING: Missing port table.**
    *   **Description:** The README lacks a dedicated table listing the essential ports (80, 8080, 8081, 8082) and their associated services/access methods.
    *   **Required Fix:** Add a port table that clearly lists ports 80, 8080, 8081, and 8082 with their corresponding services and access types (public, host-local).

6.  **WARNING: Service management information is incomplete.**
    *   **Description:** While the README mentions `m3tal-api.service` and systemctl, it doesn't provide the specific command to *start* the service, which is a common first step after installation.
    *   **Required Fix:** Include the `systemctl start m3tal-api.service` command in the Service Management section.

7.  **WARNING: Firewall note is vague.**
    *   **Description:** The firewall note mentions allowing port 80 for Traefik but doesn't explicitly mention port 443, which is also commonly used for HTTPS and is present in the `.env.example` for Traefik.
    *   **Required Fix:** Update the firewall note to include mentioning port 443 if HTTPS is configured.

8.  **SUGGESTION: Marketing copy tone detected.**
    *   **Description:** Phrases like "unified Go binary installed via APT, serving as the single entrypoint for all system operations" and "This document provides technical details and operational procedures for the M3TAL system" lean towards marketing rather than strictly technical documentation.
    *   **Required Fix:** Rephrase sentences to be more direct and technical. For example, "The M3TAL CLI is a Go binary installed via APT, used for system operations."

9.  **SUGGESTION: Quick demo section is missing.**
    *   **Description:** A "Quick Start" or "Quick Demo" section that guides a user through a minimal deployment and verification is absent, which is a common and highly beneficial element for technical documentation.
    *   **Required Fix:** Add a "Quick Start" section that provides a minimal set of commands to get a basic deployment running (e.g., installing, configuring, and starting the dashboard in local mode).

---

### Required Fixes Summary:

1.  **APT Installation:** Add the complete 3-command APT installation block.
2.  **Docker Dependency:** Explicitly state "Docker Engine + Docker Compose V2 are required."
3.  **Deployment Lifecycle:** Clarify the relationship between `/docker` and `/opt/m3tal/stack/`, and how `m3tal up` processes all `*-compose.yml` files in `/docker/`.
4.  **Traefik Routing:** Provide a clear explanation of exposing services via Traefik labels (`traefik.enable=true`) and include an example of exposing a custom user service.
5.  **Port Table:** Add a table listing ports 80, 8080, 8081, and 8082 with service and access details.
6.  **Service Management:** Add the `systemctl start m3tal-api.service` command.
7.  **Firewall Note:** Include a mention of port 443 for HTTPS if configured.
8.  **Tone:** Rephrase marketing-oriented sentences to be more technical.
9.  **Quick Demo:** Add a "Quick Start" section with basic deployment commands.