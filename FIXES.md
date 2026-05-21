## DocCritic Audit Report - M3TAL README.md

**Verdict: FAILED - Several critical pieces of information are missing or inaccurate, hindering successful deployment and operation.**

---

### Issue List:

1.  **BLOCKER: Missing Docker Dependency Statement**
    *   **Description:** The README does not explicitly state that Docker Engine and Docker Compose V2 are required *before* installing M3TAL. While it mentions them, it does not present them as a prerequisite for installation.
    *   **Classification:** BLOCKER
    *   **Required Fix:** Add a clear prerequisite statement for Docker Engine and Docker Compose V2 installation before detailing the APT installation steps.

2.  **WARNING: Incomplete Traefik Routing Explanation**
    *   **Description:** The README mentions Traefik as the gateway and discusses routing via labels and dynamic configuration. However, it does not explicitly state *how* Traefik is made aware of services in the first place. The Ground Truth indicates Traefik uses `providers.docker.exposedByDefault: false` and requires labels or dynamic config. The README implies labels are the primary mechanism but could be clearer. It also doesn't explicitly mention Traefik's static configuration file (`traefik.yml`) or its role in defining entrypoints and providers.
    *   **Classification:** WARNING
    *   **Required Fix:** Clarify that Traefik discovers services via Docker labels or dynamic configuration. Explicitly mention that Traefik itself is configured to use these mechanisms (e.g., through `traefik.yml` which defines `providers.docker` and `providers.file`).

3.  **WARNING: Port Table Incomplete**
    *   **Description:** The "Port Map" section in the README lists ports 80, 8080, 8081, and 8082. The Ground Truth indicates these are indeed the primary ports. However, the description for port 8080 could be more precise, stating it's the M3TAL Go API daemon listening on the *host-local* port, not just "host-local". Also, the Traefik dashboard port (8081) is described as "Host-local only," which is correct, but it's important to reinforce this from a security perspective and its mapping.
    *   **Classification:** WARNING
    *   **Required Fix:** Refine the descriptions for port 8080 to clearly indicate it's the Go API daemon listening on a host-local port. Explicitly mention Traefik's mapping of 8081 to `127.0.0.1:8081` for clarity.

4.  **WARNING: Marketing Copy Detected**
    *   **Description:** The "Overview" section uses marketing-oriented language ("This document provides technical details and operational procedures") rather than purely technical documentation. While minor, it detracts from a strictly technical tone.
    *   **Classification:** WARNING
    *   **Required Fix:** Rephrase the "Overview" section to be more direct and technical, e.g., "This document details the technical architecture and operational procedures for the M3TAL system."

5.  **SUGGESTION: Quick Demo Section Lacks a Concise "First Run" Example**
    *   **Description:** The "Quick Demo" section has two items: "Start the Dashboard only" and "Deploy all M3TAL stacks." While useful, it doesn't provide a single, cohesive "first-time user" quick start that walks through a minimal deployment from installation to accessing the dashboard. This would significantly improve user onboarding.
    *   **Classification:** SUGGESTION
    *   **Required Fix:** Add a "Quick Start" section that combines installation, configuration (e.g., `m3tal config wizard`), and accessing the dashboard in its default `local` mode.

6.  **WARNING: Deployment Lifecycle - `m3tal up` Explanation Not Fully Grounded**
    *   **Description:** The README states "`m3tal up` is a wrapper around `docker compose` that operates on all `*-compose.yml` files located within the `/docker/` directory". The Ground Truth clarifies that this means `docker compose` is run *across all* `*-compose.yml` files. The README's wording could be interpreted as `m3tal up` individually processing each file, rather than composing them together. The Ground Truth also mentions `/docker` is a symlink to `/opt/m3tal/stack/`, which is relevant for understanding the source of truth.
    *   **Classification:** WARNING
    *   **Required Fix:** Clarify that `m3tal up` executes `docker compose` commands that consider all `.yml` files in `/docker/` as part of a single composition. Explicitly mention that `/docker` is a symlink to `/opt/m3tal/stack/`.

7.  **WARNING: Dashboard Access Mode Explanation Could Be More Precise on Overrides**
    *   **Description:** The README explains the two modes (`local` and `traefik`) and their access methods. It mentions override files (`m3tal-compose.local.yml`, `m3tal-compose.traefik.yml`). However, it doesn't explicitly state that these override files are *applied by M3TAL* to the base `m3tal-compose.yml` based on the `DASHBOARD_EXPOSE_MODE` setting.
    *   **Classification:** WARNING
    *   **Required Fix:** Add a sentence to each mode's explanation to clarify that M3TAL dynamically applies the relevant override file (`m3tal-compose.local.yml` or `m3tal-compose.traefik.yml`) based on the `DASHBOARD_EXPOSE_MODE` setting.

8.  **WARNING: Firewall Note Could Be More Prominent**
    *   **Description:** The "Firewall Considerations" section is present, reminding users to allow port 80. However, it doesn't explicitly link this to *all* public-facing services managed by Traefik, not just the dashboard in `traefik` mode.
    *   **Classification:** WARNING
    *   **Required Fix:** Broaden the firewall note to emphasize that port 80 (and 443 for HTTPS) must be open for any service exposed publicly via Traefik, not solely tied to the dashboard's mode.

---

### Required Fixes Summary:

1.  **Prerequisite Statement:** Add a clear prerequisite for Docker Engine and Docker Compose V2 *before* the APT installation steps.
2.  **Traefik Explanation:** Elaborate on how Traefik discovers services (labels/dynamic config) and mention its static configuration (`traefik.yml`) and its role.
3.  **Port Table Refinement:** Improve descriptions for ports 8080 and 8081 for clarity and security context.
4.  **Tone Adjustment:** Rephrase the "Overview" section to be purely technical.
5.  **Quick Start Enhancement:** Create a concise "first-time user" quick start example.
6.  **Deployment Lifecycle Clarity:** Detail how `m3tal up` composes `.yml` files and mention the symlink nature of `/docker`.
7.  **Dashboard Mode Precision:** Explain that M3TAL applies override files based on `DASHBOARD_EXPOSE_MODE`.
8.  **Firewall Note Prominence:** Broaden the firewall reminder to cover all public Traefik-exposed services.