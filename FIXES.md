### **Audit Report: M3TAL Media Server Documentation**

**Auditor:** DocCritic, Senior DevOps Auditor  
**Status:** **FAILED**  
**Verdict:** **BLOCKER**. As a new user, I cannot successfully deploy this stack. The documentation assumes deep internal knowledge of the project's state machine and fails to provide the critical "Day 0" configuration steps required to bridge the gap between cloning the repo and running `./m3tal up`.

---

### **Detailed Issue List**

#### **BLOCKER**
1.  **Missing Initial Configuration Step:** The `Quick Start` commands `./m3tal init` and `./m3tal up` will fail immediately because there is no instruction to create or populate the `.env` file. The software cannot "Initialize" if the required environment variables are unknown to the system.
2.  **Implicit Dependency on `m3tal-goback`:** The documentation states the dashboard communicates with `m3tal-goback` (a separate repository), but provides no guidance on how to obtain, configure, or reach this external service. The user will be left with a broken dashboard.
3.  **Host Path Blindness:** The `Path Consistency Rule` demands `/mnt` exist on the host, but `./m3tal init` does not explain if it creates this directory or if the user must pre-create it with specific permissions.

#### **WARNING**
4.  **Traefik / Port Visibility:** No mention of which ports must be opened on the host firewall. A user doesn't know if they need to expose 80, 443, 8080, or custom ports to access the dashboard.
5.  **Docker Prerequisites:** The documentation assumes the user has `docker` and `docker-compose` installed and the current user is in the `docker` group. If I run this as a non-root user, it will likely fail with a permission error.

#### **SUGGESTION**
6.  **"Doctor" Command Timing:** The documentation should explicitly state: "If `up` fails, run `./m3tal doctor` to verify your environment."
7.  **Absolute Path Clarity:** "Absolute path" is mentioned, but a concrete example (e.g., `/home/user/media` vs `~/media`) would prevent common configuration errors.

---

### **Suggested Fixes**

**1. Create a `Pre-flight` section in README:**
*   Add a step: `cp .env.example .env`.
*   Explain the requirement to populate `BASE_STORAGE_PATH` and `API_TOKEN` before running `init`.

**2. Clarify `m3tal-goback` integration:**
*   Add a warning: *"Prerequisite: This orchestrator requires an active `m3tal-goback` instance. Please ensure the API_URL in your `.env` points to your backend instance before launching."*

**3. Explicit Port Documentation:**
*   Add a `Networking` sub-section in `README.md` listing the ports:
    *   `80/443` (Traefik Gateway)
    *   `8080` (Dashboard)
    *   Add command: `ufw allow 80/tcp && ufw allow 443/tcp`.

**4. Path/Permission Setup:**
*   Update the `Path Consistency Rule` section:
    *   "Ensure the host directory defined in `BASE_STORAGE_PATH` exists and is writable by the user executing the `./m3tal` binary: `mkdir -p /mnt/m3tal && chown $USER:$USER /mnt/m3tal`."

**5. Docker Group Check:**
*   Add a "System Requirements" block:
    *   "Ensure your user is added to the docker group: `sudo usermod -aG docker $USER` (requires re-login)."

---

**Auditor Note:** *You are building a high-performance orchestration tool. Do not assume your users are mind-readers. If the CLI manages lifecycle and storage, the documentation must explicitly guide the user through the initial environment setup, or the "sub-millisecond control plane" is useless because the system never reaches a `Running` state.*