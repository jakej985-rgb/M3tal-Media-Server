**Verdict: FAILED**

The README is missing critical information regarding the deployment lifecycle, specifically how `m3tal up` interacts with Docker Compose files and Traefik's configuration. Additionally, it lacks clear instructions on Docker dependency, making it a BLOCKER for successful deployment.

---

## Audit Issues:

1.  **BLOCKER: Missing Docker dependency statement.**
    *   **Description:** The README does not explicitly state that Docker Engine and Docker Compose V2 are required. While it mentions Docker orchestration, the specific versions and necessity of both components are not clearly called out as prerequisites.
    *   **Required Fix:** Add a clear statement in the "Prerequisites" section that "Docker Engine and Docker Compose V2 are strictly REQUIRED and must be installed prior to M3TAL installation and operation."

2.  **BLOCKER: Incomplete Deployment Lifecycle explanation.**
    *   **Description:** The README describes `m3tal up` as a wrapper around `docker compose` that operates on all `*-compose.yml` files in `/docker/`. However, it fails to explain that `/docker` is a symlink to `/opt/m3tal/stack/` and that user-added compose files are placed here. It also does not explicitly mention how new compose files are incorporated into the `m3tal up` process beyond placing them in `/docker`. The interaction with Traefik for routing new services is also not fully clarified here.
    *   **Required Fix:**
        *   Clarify that `/docker` is a symlink to `/opt/m3tal/stack/`.
        *   Explicitly state that `m3tal up` processes all `*-compose.yml` files within `/opt/m3tal/stack/` (via the `/docker` symlink).
        *   Elaborate on how to add new services: "To deploy a new service, place its `your-service-compose.yml` file directly into the `/docker/` directory. This file, along with all others ending in `-compose.yml` in this directory, will be processed by `m3tal up`."
        *   Reinforce that Traefik's routing for these new services relies on labels defined within their compose files, as outlined in the "Traefik Gateway" section.

3.  **BLOCKER: Incomplete Traefik routing explanation.**
    *   **Description:** While the README mentions Traefik as a reverse proxy and its use of labels, it doesn't fully explain how Traefik is configured and how it routes traffic. Specifically, it doesn't mention the `traefik.yml` static configuration or the `dynamic` configuration directory for routing, nor does it clearly link Traefik's operation to the `routing-compose.yml`. The dynamic routing to host-local ports (like the Go API via `host.docker.internal:8080`) is not explicitly detailed.
    *   **Required Fix:**
        *   Explain that Traefik uses both static configuration (`traefik.yml`) and dynamic configuration (files in `/etc/traefik/dynamic/`) for routing.
        *   Mention that `routing-compose.yml` is responsible for launching Traefik and Cloudflared.
        *   Clarify that Traefik discovers services via Docker labels and can also route to host-local services using dynamic configuration files that reference `host.docker.internal`.

4.  **WARNING: Missing Port Table detail (80, 8080, 8081, 8082).**
    *   **Description:** The "Port Map" section lists ports 80, 8080, 8081, and 8082. However, the descriptions are slightly ambiguous. For instance, "M3TAL API daemon (Go)" is described as "Host-local" access, which is correct, but it doesn't explicitly mention it's the internal port for API communication from other services or the dashboard in local mode. Similarly, the description for port 8082 could be clearer about its access method depending on `DASHBOARD_EXPOSE_MODE`.
    *   **Required Fix:**
        *   Refine the description for port 8080: "The internal port the M3TAL API daemon listens on. Used for communication with the dashboard (in local mode) and other internal M3TAL services."
        *   Refine the description for port 8082: "The internal port the M3TAL Dashboard container listens on. Accessed directly at `HOST_IP:8082` in local mode, or via Traefik (`dash.DOMAIN`) in traefik mode."

5.  **WARNING: Service management details are minimal.**
    *   **Description:** The README mentions `systemctl` for `m3tal-api.service` but only provides status, restart, and logs. It doesn't include common commands like `enable` or `disable` which are crucial for managing a systemd service.
    *   **Required Fix:** Add `enable` and `disable` to the list of `systemctl` commands:
        *   **Enable service on boot:** `systemctl enable m3tal-api`
        *   **Disable service on boot:** `systemctl disable m3tal-api`

6.  **WARNING: Firewall note could be more prominent.**
    *   **Description:** The firewall note is present but could be more clearly associated with initial setup and general networking requirements, especially if Traefik is involved.
    *   **Required Fix:** Consider moving the "Firewall Considerations" section to be immediately after the "Installation" section, as it's a common post-installation step that can affect service accessibility.

7.  **SUGGESTION: Quick demo could be improved.**
    *   **Description:** The "Quick Demo" section is a good start, but the explanation of `m3tal dash up` could be more descriptive. It mentions "retrieving the latest compose files," which is not explicitly how it works; it primarily applies overrides based on `.env` variables. The dependency on `DASHBOARD_EXPOSE_MODE` for accessing the dashboard could also be reinforced here.
    *   **Required Fix:**
        *   For `m3tal dash up`: "This command specifically starts the `m3tal-dashboard` container. It ensures the correct compose overrides (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) are applied based on your `DASHBOARD_EXPOSE_MODE` setting in `/etc/m3tal/.env`. If using `local` mode, the dashboard will be accessible at `http://HOST_IP:8082`."
        *   For `m3tal up`: "This command orchestrates and deploys all `*-compose.yml` files found in the `/docker/` directory (which links to `/opt/m3tal/stack/`). This includes core M3TAL components like Traefik and any user-defined compose files you have placed there."