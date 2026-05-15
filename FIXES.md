**DocCritic Audit Report: M3TAL Media Server**

**Auditor:** Senior DevOps Auditor, M3TAL Platform
**Status:** **FAILED**
**Verdict:** The documentation is currently a "developer's souvenir." It describes the architecture well but fails to provide a functional deployment path for a new operator. A new user will experience immediate runtime failures due to missing environment configuration and directory assumptions.

---

### 🚨 Issue List

#### **BLOCKER**
*   **Missing `.env` Template:** There is no documentation on how to generate the initial environment file. A user running `./m3tal init` will likely fail if the application expects pre-existing environment variables that aren't defined.
*   **Undefined `BASE_STORAGE_PATH` Requirement:** The "Path Consistency Rule" mentions `/mnt` but does not explicitly state that the user *must* create this directory or mount a drive there before running the orchestrator.
*   **No Dependency Pre-check:** The documentation fails to state that the host requires Docker and Docker Compose (v2+) to be installed.

#### **WARNING**
*   **Network Port Blindness:** There is zero information regarding which ports (80, 443, 8080) must be open on the host firewall. Traefik handles routing, but if the user has a local web server (e.g., Nginx/Apache) running, the `m3tal up` command will silently fail or conflict.
*   **Ambiguous `init` behavior:** It is unclear if `./m3tal init` creates the `.env` file automatically or if the user is expected to copy an example file (which doesn't exist in the repo).

#### **SUGGESTION**
*   **"Quick Start" Logic Gap:** The transition from `go build` to `./m3tal init` assumes the binary has permission to manipulate the file system. Add a note about `chmod +x` and required user permissions (docker socket access).
*   **Architecture Diagram Clarity:** The diagram shows `m3tal-goback` as an external dependency, but doesn't explain how to link it (e.g., does it need to be in the same Docker network, or is it a URL-based API?).

---

### 🛠 Suggested Fixes

#### 1. Add a Pre-flight Checklist (Documentation)
Add a section before "Quick Start":
> **System Prerequisites:**
> - Linux host with Docker Engine and Docker Compose (v2.x) installed.
> - User must be in the `docker` group.
> - Ensure port 80 and 443 are available.
> - Mount your media storage to `/mnt` on the host, or update `BASE_STORAGE_PATH` via `./m3tal config`.

#### 2. Create an `.env` initialization step
Modify the "Quick Start" to:
```bash
# 1. Compile
go build -o m3tal main.go 

# 2. Setup Configuration
./m3tal setup  # (New command: prompts for API keys/storage paths)
# OR 
cp .env.example .env && ./m3tal config
```

#### 3. Formalize the Path Constraint
Update the **Path Consistency Rule** section:
> **Strict Path Requirement:** 
> The orchestrator mandates that the host path defined in `BASE_STORAGE_PATH` is accessible. If `/mnt` does not exist on your host, run: `sudo mkdir -p /mnt && sudo chown $USER:$USER /mnt`.

#### 4. Document Firewall/Network Exposure
Add a "Network Requirements" section:
> **Networking:** 
> M3TAL uses Traefik as an ingress controller. Ensure your host firewall allows traffic on TCP ports 80 and 443. Access the dashboard at `http://<host-ip>`.

#### 5. Clarify External Dependencies
In the "Troubleshooting" section, explicitly define:
> **m3tal-goback:** This must be reachable via the network defined in your `.env`. If running on the same host, use `http://host.docker.internal:PORT`.

---
**Auditor Note:** *Fix these items, and the project might actually be deployable by a human. Re-submit for audit once the `setup` flow is non-destructive and explicit.*