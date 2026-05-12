# DocCritic Audit Report: M3TAL Platform (v1.4)

**Verdict: FAILED**
The documentation provides a high-level overview but lacks the technical specificity required to actually deploy the software. As a new user, I am left guessing about the state of my filesystem, the validity of generated secrets, and how the underlying orchestration stack (Traefik/Docker) interacts with the host.

---

### Issue List

#### 1. BLOCKER: Missing Filesystem Initialization
The docs mention `./data` as a default, but there is no instruction to create this directory. Docker volume mounts will fail or create root-owned directories if the directory doesn't exist before `up`.
*   **Fix:** Add `mkdir -p ./data` to the `init` or `Quick Start` section.

#### 2. BLOCKER: Ambiguous `.env` Generation
The documentation says `./m3tal init` generates secrets, but it does not specify if it creates the `.env` file from a template (e.g., `.env.example`). If the file is missing, the orchestration will fail.
*   **Fix:** Explicitly state: "The `init` command will generate a `.env` file from the `example.env` template. If it does not exist, run `cp example.env .env` before proceeding."

#### 3. WARNING: Confusing "Source of Truth" vs. Manual Override
The docs forbid manual `docker compose` usage but mention `make up`. If `make up` triggers `docker compose` internally, the user needs to know where the compose file is located (e.g., `source/m3tal-stack/docker-compose.yml`) to troubleshoot mounting issues.
*   **Fix:** Explicitly define the path of the docker-compose file being orchestrated.

#### 4. WARNING: Traefik Configuration Void
You define Traefik as the gateway on port 8080/443, but you provide no instructions on how to configure SSL (e.g., Let's Encrypt/CertResolver) or how to map services to domains.
*   **Fix:** Add a "Networking & SSL" section explaining that custom routing requires modifying `source/m3tal-stack/traefik.yml` or labels in individual services.

#### 5. SUGGESTION: Port Collision Risks
The documentation lists ports 8080, 8082, 5050, 443. These are extremely common defaults.
*   **Fix:** Add a warning: "Ensure ports 80, 443, 5050, and 8080 are not already in use by other services (like Nginx or local web servers) before starting."

#### 6. SUGGESTION: Inconsistent "Build" instructions
Step 1 lists `make build` and `build.sh`. Step 3 lists `make up`. Does `make` look for a local `Makefile`? The document does not confirm the existence of a `Makefile` in the root directory.
*   **Fix:** Confirm the existence of a `Makefile` and include its contents or clarify that `make` is an optional wrapper for the CLI binary.

---

### Summary of Required Documentation Updates

| Action | Priority |
| :--- | :--- |
| Create a `Makefile` or document specific shell command requirements. | High |
| Add a "Pre-flight Check" section to ensure `./data` exists. | High |
| Document the `.env.example` file usage workflow. | High |
| Clarify how Traefik handles SSL (or if it's currently insecure/HTTP-only). | Medium |
| Add a note on `sudo` usage if the user is not in the `docker` group. | Medium |

**DocCritic Note:** *Stop assuming the environment is perfect. A user will break this within 30 seconds if the `/data` directory isn't managed or the `.env` file is missing. Fix these, and I'll reconsider the deployment readiness.*