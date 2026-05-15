## **DocCritic Audit Report: M3TAL Core Documentation**

**Verdict:** **FAILED**. As a new user, I cannot successfully deploy this system. The documentation suffers from "developer myopia," assuming local environment configurations (like specific paths, service availability, and network routing) that are not explicitly documented. It fails to bridge the gap between "installing a binary" and "making the system actually functional."

---

### **Issue List**

#### **BLOCKER**
*   **[BLOCKER] Missing `.env` lifecycle:** The `m3tal up` command implies a Docker Compose deployment, but there is no instruction on how to generate or populate the `.env` file required by Docker Compose stacks. If `m3tal init` generates it, this is not stated.
*   **[BLOCKER] Traefik Configuration:** The docs mention Traefik exposes 80/443, but provide zero instructions on how to configure Traefik labels, domain names, or SSL/TLS certificates. The system will likely fail to start or remain inaccessible behind an unconfigured proxy.
*   **[BLOCKER] Dependency Orchestration:** The `m3tal-api` is referenced in the mermaid diagram and the troubleshooting section, but it is **not** defined in the "System Components" or "Quick Start" as a manual install requirement. Is it part of the stack? Does the user have to run it separately?

#### **WARNING**
*   **[WARNING] Path Assumption:** The documentation states the orchestrator maps storage to `/mnt` inside containers, but it doesn't clarify if the user needs to create `/mnt` on their host or if the `m3tal init` command handles the directory creation/permissioning.
*   **[WARNING] Binary vs. Source Conflict:** The "Quick Start" uses a Debian package, but the "Development" section implies building from source. There is no warning that mixing these (e.g., global binary vs. local build) will cause version skew or configuration path conflicts.

#### **SUGGESTION**
*   **[SUGGESTION] Initial Setup Workflow:** Add a "First Run Checklist" section. Currently, the user jumps from `apt install` to `m3tal up` without being told they need to define an `API_TOKEN` or service credentials in a specific file.
*   **[SUGGESTION] Port Conflict Warning:** If the user has a local Nginx, Apache, or another service on 80/443, `m3tal up` will crash. Add a warning about checking host port availability.

---

### **Suggested Fixes**

1.  **Deployment Pre-flight:**
    *   Add a step: `m3tal config setup`. Clearly document that this creates `/etc/m3tal/config.yaml` AND the necessary `.env` files for the Docker stack.
2.  **Explicit Proxy Instructions:**
    *   Add a section: `🌐 Configuring Traefik`. Explain that users must define their domain in the config or provide an example `docker-compose.override.yml` for users who want to use a specific local domain.
3.  **Component Clarification:**
    *   Update "System Components" to clarify: *"The `m3tal-api` is automatically managed by the `m3tal up` command via the underlying Docker stack; no manual installation is required."*
4.  **Path Resolution Logic:**
    *   Explicitly state: *"Run `m3tal init` with sudo permissions. It will ensure the host directory specified in `BASE_STORAGE_PATH` exists and is accessible to the Docker user."*
5.  **Environment Isolation:**
    *   Add a bold note: **"Important: Choose one installation method. Do not mix the Debian-packaged binary with custom-compiled binaries from the source root to avoid configuration drift."**
6.  **Configuration Example:**
    *   Provide a full `config.yaml` example in the documentation, including the `API_TOKEN` field, to prevent users from guessing the required keys.