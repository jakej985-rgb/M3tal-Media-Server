**To:** M3TAL Development Team  
**From:** DocCritic, Senior DevOps Auditor  
**Subject:** AUDIT REPORT: M3TAL Media Server Documentation  
**Status:** **FAILED**

---

### **Verdict**
**BLOCKER.** The provided README is essentially a "Developer’s Cheat Sheet," not deployment documentation. A new user following these instructions will encounter immediate runtime failures, silent configuration errors, and environment-level permission lockouts. The documentation assumes a "perfect" pre-configured host environment without providing the necessary boot-strapping steps.

---

### **Issue List**

#### **BLOCKER**
1. **Missing Initialization Context:** `./m3tal init` is cited as the entry point, but there is no documentation on what it creates. Does it generate a `.env` file? Does it require root/sudo access for Docker socket management?
2. **Missing Environment Configuration:** There is no instruction on how to generate the mandatory `.env` file (e.g., `cp .env.example .env`). The system requires `BASE_STORAGE_PATH` and `API_TOKEN`, but provides no template or guidance on valid input.
3. **Implicit Dependency Failure:** The "Path Consistency Rule" assumes `/mnt` exists and is writable. If the user runs this on a clean machine (e.g., Ubuntu server), the stack will crash immediately because `/mnt` is often restricted or absent.

#### **WARNING**
4. **Networking Blind Spots:** The documentation references Traefik but fails to state which ports must be open on the host firewall. A user will be unable to access the Dashboard because `80/443` are not mentioned as requirements.
5. **Dashboard Connectivity:** The `Troubleshooting` section mentions `m3tal-goback` must be "reachable," but provides no configuration step to set the API endpoint URL in the environment variables.

#### **SUGGESTION**
6. **Binary Lifecycle:** The Quick Start suggests `go build`. It should also suggest a `make` target or a release binary download link for users who do not have Go toolchains installed on their production servers.
7. **Pre-flight Checks:** The `./m3tal doctor` command is excellent, but it is listed in the reference table rather than the "Quick Start." It should be the *first* command a user runs after building.

---

### **Suggested Fixes**

**1. Create an "Environment Setup" Section:**
Before running `init`, the user must be instructed:
```bash
cp .env.example .env
nano .env # Edit BASE_STORAGE_PATH and API_TOKEN
```

**2. Clarify Host Prerequisites:**
Explicitly define the host expectation:
> "Ensure the host directory defined in `BASE_STORAGE_PATH` exists and is owned by the current user. The M3TAL orchestrator requires Docker Socket access; ensure your user is in the `docker` group."

**3. Update Quick Start Flow:**
Revise the steps to ensure safety:
```bash
# 1. Compile
go build -o m3tal main.go

# 2. Pre-flight check (New Step)
./m3tal doctor

# 3. Configure (New Step)
./m3tal init # This should generate the .env template if missing

# 4. Launch
./m3tal up
```

**4. Networking/Firewall Disclosure:**
Add a "Connectivity" section:
> "M3TAL exposes services via Traefik. Ensure the following host ports are open:
> * **80 (HTTP)**: Entry point for web traffic.
> * **443 (HTTPS)**: Default for secure dashboard access."

**5. Architecture Clarity:**
Clarify that `m3tal-goback` is a *separate repo* that must be running *before* `m3tal up`. The current wording makes it sound like a magic internal dependency that might (or might not) exist. Provide the expected URL format (e.g., `BACKEND_URL=http://<host-ip>:port`).