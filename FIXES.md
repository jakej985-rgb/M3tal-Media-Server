### **DocCritic Audit Report: M3TAL Control Plane**

**Verdict:** **UNUSABLE**. As a new user, I cannot deploy this project. The documentation suffers from "expert bias," assuming the user understands the internal orchestration logic without explicit instructions. It currently functions more like a high-level architecture overview than a deployment manual.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing `m3tal.py`/Installation logic:** The docs mention `m3tal.py` in the prompt requirements, but it is absent from the `Step-by-Step` guide. Does the orchestrator rely on Python scripts or Go binaries?
2.  **`BASE_STORAGE_PATH` Ambiguity:** The docs state it "Must map to /mnt." If I am on a fresh server, `/mnt` does not exist or is managed by the OS. Will the system fail if I don't manually create these mount points?
3.  **Traefik Gateway Missing:** The architectural blueprint mentions a "unified deployment stack" and "Docker Proxy," but there is no mention of Traefik or an Ingress configuration. If I expose port 8082, is it publicly accessible? How do I secure it?
4.  **No `m3tal-stack` Execution context:** The guide says "Spin up the infrastructure: `./m3tal up`." What exactly does this binary *do*? Does it perform `docker compose up`? If the Go binary fails (e.g., missing permissions to the Docker socket), there is no recovery documentation.

#### **WARNING**
5.  **Undefined Environment Dependencies:** `cp .env.example .env` is a start, but there is no mention of whether the user needs to generate the secrets (e.g., `openssl rand -hex 32`). A user pasting "example" tokens into production is a security risk.
6.  **Dependency Hell:** It lists Go 1.21+ *and* Python 3.10+ as prerequisites, but it is unclear if the user needs to install these on the *host* or if the `m3tal-stack` handles them inside containers.

#### **SUGGESTION**
7.  **Service Access Gap:** Accessing `localhost:8082` assumes a bare-metal GUI environment. Most DevOps users will be deploying this on a remote headless server. Documentation should mention SSH tunneling or proxy configuration for remote access.
8.  **Structure:** The documentation lacks a "Verification" step that confirms the stack is healthy *before* the user tries to load the web interface.

---

### **Suggested Fixes**

*   **For Blocker 1 & 4 (Orchestration):** Clarify the `m3tal` binary's role. If it is a wrapper, document the requirement: *'Ensure your user is in the `docker` group, as the `m3tal` binary interacts directly with the Docker socket.'*
*   **For Blocker 2 (Storage):** Provide a `mkdir` command explicitly for the required folders. 
    *   *Fix:* `mkdir -p /mnt/m3tal/config /mnt/m3tal/data`. 
*   **For Blocker 3 (Traefik/Security):** Add a section on "Network Security." If Traefik is used, provide a sample snippet of the `docker-compose.yml` labels required.
*   **For Warning 5 (Secrets):** Update the `Configure Settings` step:
    *   *Fix:* "Generate secure keys using `openssl rand -hex 16` and update the following fields in `.env`."
*   **For Warning 6 (Env):** Clarify the runtime environment:
    *   *Fix:* "The Go binary and Python dashboard are designed to run [on the host / inside containers]. You do [not] need to install Go/Python on the host system."
*   **For Suggestion 7 (Remote Access):** Add a note: "If running on a remote headless server, access the Dashboard via SSH tunnel: `ssh -L 8082:localhost:8082 user@your-server-ip`."