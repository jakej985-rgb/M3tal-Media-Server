**DocCritic Audit Report: M3TAL-CORE-001**
**Auditor:** Senior DevOps Auditor, M3TAL Platform
**Status:** **REJECTED**

As a new user attempting to deploy the M3TAL stack, I am currently stuck. The documentation reads like a whitepaper rather than a deployment guide. It assumes the infrastructure is already provisioned, compiled, and configured, offering zero "Getting Started" path for a clean install.

---

### 🚨 Issue List

| ID | Severity | Description |
| :--- | :--- | :--- |
| 1 | **BLOCKER** | **Missing Compilation Steps:** The doc mentions "Go 1.21+" but provides no `go build` command to generate the `m3tal` binary. |
| 2 | **BLOCKER** | **Configuration Void:** No template or instructions provided for `/etc/m3tal/.env`. The system will crash on boot without required variables (API keys, DB paths, etc.). |
| 3 | **BLOCKER** | **Missing Orchestration Setup:** How is `/opt/m3tal` populated? Does the binary create these directories? Does the user have to clone them? |
| 4 | **WARNING** | **Traefik Ambiguity:** The doc claims "The Orchestrator maintains the base Traefik proxy," but provides no `docker-compose.yml` that actually configures the Traefik entry points or network requirements. |
| 5 | **WARNING** | **Mount Point Assumptions:** The documentation assumes `/mnt/m3tal-media` exists. A new user's system will fail if this directory is not pre-created or owned by the correct user/UID. |
| 6 | **SUGGESTION** | **Binary Installation:** No instructions on how to install the compiled binary to `/usr/bin/m3tal`. |

---

### ✅ Suggested Fixes

**1. Add a "Quick Start" Installation Section:**
Provide the exact commands to get the binary ready:
```bash
# Compilation
go mod tidy
go build -o m3tal ./cmd/m3tal/main.go
sudo cp m3tal /usr/bin/m3tal

# Initialization
sudo mkdir -p /etc/m3tal /opt/m3tal/stack /var/lib/m3tal
# Provide a template for .env
cp .env.example /etc/m3tal/.env 
```

**2. Provide a Reference Docker Compose:**
The "Deployment" section currently shows a single service snippet. Users need a full `docker-compose.yml` that includes `m3tal-goback`, `m3tal-godash`, and `traefik` to see how labels and networking connect.

**3. Explicit Directory Setup:**
Add a "Pre-requisite Setup" section:
```bash
sudo mkdir -p /mnt/m3tal-media
sudo chown -R $USER:$USER /mnt/m3tal-media
```

**4. Document Traefik Ports:**
Explicitly list required host ports (e.g., 80, 443, 8080) so the user knows what to open in their firewall (UFW/Cloud provider).

**5. Define the `.env` schema:**
Add a table defining what keys are mandatory:
* `M3TAL_API_KEY`: Required for communication between back/dash.
* `M3TAL_DB_PATH`: Location for internal state.

---

**Verdict:** 
**BLOCKER.** The documentation is currently an architectural overview, not a functional deployment guide. A new user cannot deploy this system without guessing commands or reverse-engineering the Go code. **Documentation must be updated to include a "Quick Start" and a valid `docker-compose.yml` reference.**