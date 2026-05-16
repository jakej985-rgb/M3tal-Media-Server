**DocCritic Audit Report: M3TAL Core Orchestrator**

**Verdict: FAILED**
As a new user, I am currently unable to deploy this project. The documentation describes *what* the system is, but provides no instruction on *how* to build, configure, or initialize it. It assumes the user already possesses a pre-configured server environment.

---

### Issue List

#### 1. BLOCKER: Missing Build/Installation Instructions
The documentation mentions a `m3tal` binary but provides no instructions on how to compile it from the source (Go 1.21+ is listed, but no `go build` command is provided). There is no "Getting Started" guide to fetch the repository and install the binary to `/usr/bin/`.

*   **Fix:** Add a "Quick Start" section with `go build -o m3tal .` and binary installation steps.

#### 2. BLOCKER: Missing `.env` Configuration Schema
The docs reference `/etc/m3tal/.env` as the "Global Configuration Source of Truth," but provide no template, example, or explanation of required variables (e.g., API keys, database URLs, storage paths).

*   **Fix:** Provide a `.env.example` file in the repo and document the mandatory environment variables required for the orchestrator to boot.

#### 3. BLOCKER: Missing Initialization Workflow
The docs reference a `/opt/m3tal` directory and `/var/lib/m3tal/` state, but fail to mention if the CLI initializes these directories. A new user will encounter a "Permission Denied" or "File Not Found" error immediately.

*   **Fix:** Include an `m3tal setup` command or documentation on creating these directories and setting correct ownership (`chown`).

#### 4. WARNING: Ambiguous Docker Usage
The "Deployment (Docker)" section provides a raw YAML snippet but fails to specify how to *run* it. Should I save this as `docker-compose.yml`? Do I need to clone the `m3tal-goback` and `m3tal-godash` repos manually, or does the orchestrator pull them?

*   **Fix:** Provide a complete `docker-compose.yml` example that integrates the backend, dashboard, and Traefik, showing how the Orchestrator links them.

#### 5. WARNING: Missing Gateway/Access Information
The documentation mentions Traefik but doesn't explain how to access the services. What are the default ports? What is the expected entry point for the dashboard?

*   **Fix:** Explicitly state the Traefik entry points (e.g., `http://localhost:8080`) and any required labels for auto-discovery.

#### 6. SUGGESTION: Hardcoded Path Assumptions
The reliance on `/mnt` and `/opt/m3tal` is strict. On some systems (or rootless Docker setups), these paths may not exist or require elevated privileges.

*   **Fix:** Add a section on "Prerequisites/Environment Prep" that explains how to create the directory structure or how to override paths via the `.env` file.

---

### Required Actions for Remediation:
1.  **Add `m3tal setup` logic:** Create a command to auto-generate the directory hierarchy (`/etc/m3tal`, `/opt/m3tal`) with appropriate permissions.
2.  **Provide a "First Run" Guide:**
    *   Step 1: Clone repo.
    *   Step 2: `go build`.
    *   Step 3: `sudo ./m3tal setup` (creates folders/configs).
    *   Step 4: Configure `/etc/m3tal/.env`.
    *   Step 5: `docker compose up -d`.
3.  **Document Environment Variables:** Clearly define every key required in the `.env` file for the stack to reach "Ready" state.