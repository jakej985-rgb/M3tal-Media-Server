**DocCritic Audit Report: M3TAL Platform**

**To:** M3TAL Engineering Team  
**From:** DocCritic, Senior DevOps Auditor  
**Subject:** Documentation Audit - Repository "M3TAL Media Server"

---

### **Verdict: BLOCKER**
The current documentation is an "insider’s manual." It assumes the user already possesses the environment, knows where to put secrets, and understands how the proprietary `m3tal` binary interacts with host-level dependencies. A new user will fail at the `./m3tal init` stage because the prerequisites for the orchestrator are not documented.

---

### **Detailed Issue List**

#### **BLOCKER**
1. **Missing `.env` Template:** The orchestrator relies on `BASE_STORAGE_PATH` and `API_TOKEN`. There is no `template.env` or instruction on where to create the configuration file. `m3tal init` will likely crash or error out if the environment is missing.
2. **Missing Dependency Check:** The README mentions `m3tal init` validates paths, but fails to state that the host *must* have the Docker daemon running, a specific Go version, and potentially specific system permissions (Docker socket access) to run the `m3tal` binary.
3. **Hardcoded Path Dependency:** The instruction "the host `BASE_STORAGE_PATH` is always mounted to `/mnt`" is a critical architectural requirement. If a user tries to run this on a standard Linux distro without an existing `/mnt` mapping, the stack will fail.

#### **WARNING**
4. **Networking Blind Spot:** The documentation mentions Traefik but ignores ports. Which ports must be open on the host firewall? (80/443? 8080?) A user behind a restrictive firewall will be unable to reach the dashboard.
5. **Orchestrator Bootstrapping:** The "Quick Start" assumes `./m3tal` can be run immediately after `go build`. It doesn't mention if the binary needs root/sudo privileges to manipulate Docker, or if it requires being added to the `docker` group.

#### **SUGGESTION**
6. **Binary Location:** The build command produces `./m3tal`. It should suggest moving this to `/usr/local/bin` or creating an alias so the user isn't stuck running it from the repo root forever.

---

### **Suggested Fixes**

**1. Create a `bootstrap.sh` or `setup` section:**
Add a step:
```bash
cp .env.example .env
# Edit the .env file with your specific paths and tokens
nano .env
```

**2. Add a `Prerequisites` section to the README:**
*   Docker Engine installed and current user added to `docker` group.
*   Go 1.21+ installed.
*   Required Ports: `80`, `443` (Traefik), and `8080` (Internal).

**3. Clarify the `/mnt` requirement:**
Explicitly warn the user:
> *Warning: M3TAL requires that the host directory defined in `BASE_STORAGE_PATH` be available. Ensure your system user has read/write permissions to this directory before running `./m3tal init`.*

**4. Add a "Port Exposure" table:**
| Port | Protocol | Purpose |
| :--- | :--- | :--- |
| 80 | TCP | Traefik HTTP Gateway |
| 443 | TCP | Traefik HTTPS Gateway |
| 8080 | TCP | Internal Orchestration API |

**5. Fix the `./m3tal init` gap:**
Explicitly explain what the orchestrator does during `init`. Does it write to `/etc/m3tal`? Does it pull images? Does it generate the `docker-compose.yml`? A user needs to know *what* they are running.

**DocCritic’s Closing Note:** *Documentation that assumes the reader knows the system architecture is useless for onboarding. Fix the .env workflow and the port mapping, or you will spend your entire lifecycle providing manual support via Slack/Email.*