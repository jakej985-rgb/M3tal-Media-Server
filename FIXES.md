### **DocCritic Audit Report: M3TAL Core**

**Verdict: BLOCKER**
The documentation is fundamentally broken for any user attempting to deploy this system. While the "native package" flow exists, the documentation fails to bridge the gap between installation and operational runtime. It makes dangerous assumptions about host directory structures, neglects security configuration, and provides zero guidance on critical post-installation failures.

---

### **Issue List**

#### **BLOCKER**
*   **[BLOCKER] Missing `.env` lifecycle:** The `m3tal up` command implies a containerized stack, but there is no mention of how the `API_TOKEN` defined in `config.yaml` is injected into the Docker environment. Does `m3tal up` auto-generate a `.env` file for the Compose stack? If not, the containers will fail to authenticate.
*   **[BLOCKER] Assumption of `/etc/m3tal/` permissions:** Installing via `apt` places config in `/etc/m3tal/`. As a standard user, I cannot write to this directory. The docs do not mention running as `sudo` for `config` commands or provide a path for user-space configuration overrides.
*   **[BLOCKER] Hidden Host Requirements:** The "Path Consistency Rule" states the system maps host paths to `/mnt` inside the container. If I run `m3tal config set path /home/user/media`, does the orchestrator *create* the host path? Does it require specific permissions? If I don't have a `/mnt` directory on the host (e.g., standard Debian/Ubuntu), will the container fail?

#### **WARNING**
*   **[WARNING] Traefik Access/Gateway:** The docs state it exposes port 80/443. It does *not* mention that it requires Traefik labels in the `m3tal-stack` or how to configure TLS (SSL certificates). A user starting this will likely get a "Connection Refused" or an untrusted certificate error without a `traefik.yml` or dynamic configuration guide.
*   **[WARNING] Missing `m3tal-api` definition:** The Architecture diagram shows `m3tal-api`, but there is no documentation on whether this is a standalone container I need to pull, a separate repo, or if it is bundled within the `m3tal-stack`.

#### **SUGGESTION**
*   **[SUGGESTION] `m3tal doctor` ambiguity:** The command exists but the docs don't explain what "health" looks like. Add an example output of a "Successful" vs "Failed" `doctor` run so users can identify common misconfigurations.
*   **[SUGGESTION] Missing `docker-compose.yml` location:** If I build from source, where is the template stack located? The docs mention `deploy/stack/` but do not explain if the binary expects these files to be in a specific relative path to the binary, or if it embeds them.

---

### **Suggested Fixes**

1.  **Clarify Configuration Management:** Add a section explicitly defining the relationship between `/etc/m3tal/config.yaml` and the Docker environment variables. 
    *   *Fix:* "The `m3tal` binary automatically generates a `.env` file in the stack directory based on your `config.yaml`. Ensure you run configuration commands with `sudo`."
2.  **Explicit Path/Volume Guidance:**
    *   *Fix:* Add a note: "Ensure the user running `m3tal` has read/write permissions to your `BASE_STORAGE_PATH`. If using external mounts, ensure they are mounted to the host prior to running `m3tal up`."
3.  **Networking/Traefik Documentation:**
    *   *Fix:* Add a `docs/NETWORKING.md` reference or a quick-start note: "By default, Traefik is configured for HTTP. To enable HTTPS, create `acme.json` and configure `traefik.yml` in `/etc/m3tal/stack/`."
4.  **CLI Permissions:**
    *   *Fix:* Clearly state: "All `m3tal` commands managing system state must be executed with root privileges."
5.  **Fix the "Missing API" Gap:**
    *   *Fix:* Clarify if `m3tal-api` is deployed by the orchestrator. If it is part of the `m3tal-stack`, label it as such in the architecture section.