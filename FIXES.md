## M3TAL README Audit Report

**Verdict: FAILED - Several critical pieces of information are missing or unclear, preventing a user from successfully deploying and operating M3TAL.**

---

### Issue List:

1.  **BLOCKER: Missing Docker Compose V2 Clarification for `m3tal up`**
    *   **Description:** The README states that `m3tal up` operates on all `*-compose.yml` files in `/docker/`. However, it does not explicitly mention that `m3tal up` internally uses `docker compose` (V2) and that this command is the underlying mechanism for orchestrating these files. This could lead users to believe it's a proprietary orchestration system without understanding its foundation.
    *   **Required Fix:** Clearly state that `m3tal up` is a wrapper around `docker compose` V2 and that it processes all `.yml` files within the `/docker` directory.

2.  **BLOCKER: Incomplete Deployment Lifecycle Explanation**
    *   **Description:** While the README mentions `m3tal up` and the `/docker` directory, it doesn't fully explain the stack deployment mechanism as it relates to how Docker Compose itself works. Specifically, it doesn't clearly link the user-facing `/docker` directory to the underlying `/opt/m3tal/stack/` directory where the actual compose files and Traefik configurations reside. The statement "all stack files reside" is also a bit ambiguous – it should clarify that this is where the *M3TAL-managed* stack files are.
    *   **Required Fix:** Refine the "Deployment Lifecycle" section to explicitly state that `/docker` is a symlink to `/opt/m3tal/stack/`, and that `m3tal up` essentially runs `docker compose up` across all `.yml` files found within this stack directory.

3.  **BLOCKER: Unclear Traefik Service Exposure Mechanism**
    *   **Description:** The README mentions Traefik as the HTTP gateway and explains it uses "Traefik labels defined within their Docker Compose service definitions" for discovery. However, it doesn't explicitly state that `traefik.enable=true` is a prerequisite for a service to be discovered by Traefik. It also doesn't clarify that Traefik uses the `proxy` network to communicate with services.
    *   **Required Fix:** Add a note that `traefik.enable=true` is required for Traefik to pick up a service. Also, mention that Traefik and its proxied services typically communicate over the `proxy` Docker network.

4.  **WARNING: Missing Port 8081 in Port Table**
    *   **Description:** The "Port Map" section lists ports 80, 8080, and 8082 but omits port 8081, which is used for the Traefik dashboard (host-local only).
    *   **Required Fix:** Add port 8081 to the "Port Map" table with its corresponding service (Traefik dashboard) and access details (host-local only).

5.  **WARNING: Inconsistent Dashboard Port Definition**
    *   **Description:** The "Port Map" table lists `8081` for the Traefik dashboard, but the ground truth indicates `127.0.0.1:8081:8080` in the `routing-compose.yml`. This implies that the *internal* port Traefik listens on is 8080, and it's exposed locally on the host as 8081. The table should reflect this accurately.
    *   **Required Fix:** Clarify in the "Port Map" that port 8081 on the host maps to port 8080 *inside* the Traefik container.

6.  **WARNING: Ambiguous `m3tal-api.service` Management**
    *   **Description:** The "Service Management — systemd" section correctly identifies `m3tal-api.service` but doesn't explicitly state that this service is what runs the *M3TAL API daemon*. It's implied, but a direct statement would improve clarity.
    *   **Required Fix:** Add a sentence to explicitly state that `m3tal-api.service` manages the M3TAL API daemon.

7.  **WARNING: Incomplete Firewall Reminder**
    *   **Description:** The firewall reminder only mentions allowing port 80 for Traefik. It does not mention port 443, which is also listed in the `.env.example` for Traefik and is crucial for HTTPS configurations.
    *   **Required Fix:** Update the "Firewall Considerations" to include a reminder for port 443 if HTTPS is configured.

8.  **SUGGESTION: Marketing Tone in "Deployment Lifecycle"**
    *   **Description:** The phrase "effectively deploying each as an independent stack" and "user-facing symlink alias for all stack operations" leans slightly towards marketing copy rather than purely technical documentation. While not a blocker, it could be more direct.
    *   **Required Fix:** Rephrase sentences to be more concise and technically focused. For example, "The `/docker` directory is a user-facing symlink to the canonical stack directory `/opt/m3tal/stack/`, containing all M3TAL-managed compose files and configurations. `m3tal up` orchestrates all `.yml` files found here using `docker compose`."

9.  **SUGGESTION: Missing Quick Start Detail for Dashboard**
    *   **Description:** The "Quick Demo" section mentions `m3tal dash up` but doesn't explain *why* a user might need to run this command specifically for the dashboard, or what override it applies based on `DASHBOARD_EXPOSE_MODE`. It also implies this is a separate step from `m3tal up`.
    *   **Required Fix:** Clarify that `m3tal dash up` is an optional command to manage the dashboard specifically, and its behavior is dictated by `DASHBOARD_EXPOSE_MODE`. Explain that `m3tal up` will also deploy the dashboard (and any override configurations), making `m3tal dash up` potentially redundant for a full deployment unless used for isolated dashboard management.

---