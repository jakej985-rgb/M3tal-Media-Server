**To:** DocSmith  
**From:** DocCritic, Senior DevOps Auditor  
**Subject:** AUDIT REPORT: M3TAL Core Orchestrator README  

---

### **Verdict: FAILED**
The current documentation is a "developer’s brain dump" rather than an operational guide. It suffers from dangerous assumptions regarding filesystem state and lacks essential network exposure details. While the APT instructions are present, the project is currently impossible to deploy reliably in a production or clean-room environment.

---

### **Issue List**

#### **BLOCKER**
1. **Host-Level Dependency Assumption**: The `Filesystem Standard` lists `/mnt/m3tal-media` as a requirement, but provides no instructions on how the user should configure this. If the directory does not exist, the orchestrator/containers will fail to start (or worse, create root-owned directories on the host).
2. **Missing Port Exposure/Traefik**: The documentation mentions a "Dashboard" but does not define which ports it binds to, nor does it mention the Traefik configuration required to route traffic to the containerized components. Users cannot access the services they just deployed.
3. **Missing "Stop/Teardown" Workflow**: A "Quick Start" is useless without a "Quick Shutdown." Users need to know how to tear down the environment properly to avoid orphaned containers or zombie processes.

#### **WARNING**
1. **Ambiguous `m3tal up` context**: The documentation claims `m3tal up` starts the infrastructure defined in `deploy/stack`. It is unclear if this command expects the user to be in a specific directory or if it assumes a hardcoded path.
2. **Missing Docker Socket Permissions**: The `m3tal-orchestrator` service requires access to `/var/run/docker.sock`. Running this without noting the requirement for the user to be in the `docker` group or requiring `sudo` for the CLI will lead to immediate permission denied errors.

#### **SUGGESTION**
1. **Marketing Overload**: The section "Ecosystem Integration Rules" is fluff. A user reading a README cares about *how* to deploy, not the architectural philosophy of the "Go-Native Migration." Move this to a `/docs/ARCH.md` file.
2. **Environment Variables**: The `m3tal-orchestrator` section mentions `M3TAL_ROOT`, but there is no guide on how to provide this variable. Is it via a `.env` file? Shell export? CLI flag?

---

### **Required Fixes**

1. **Implement Pre-flight Checks**: Add a section before `m3tal setup` that verifies the environment:
   ```bash
   # Add this to documentation
   sudo mkdir -p /mnt/m3tal-media
   sudo chown $USER:$USER /mnt/m3tal-media
   ```
2. **Define Networking**: Explicitly define the port mappings in the Docker section. Example:
   > "The Dashboard is exposed on host port 8080. Ensure Traefik is configured to point to this container port if using a reverse proxy."
3. **Clarify CLI usage**: Explicitly state the relationship between the CLI and the Docker socket:
   > "Note: The `m3tal` binary requires Docker socket access. Ensure your user is added to the `docker` group (`sudo usermod -aG docker $USER`) before running `m3tal up`."
4. **Add Teardown**: Add:
   ```bash
   m3tal dash down
   m3tal down
   ```
5. **Clean the Copy**: Replace the "Ecosystem Integration Rules" with a **Configuration Reference** table showing exactly which Environment Variables are required for the CLI to function.

---

**Auditor Note:** *Stop treating the README like a manifesto. A README is a checklist for a sysadmin. Keep it technical, keep it actionable, and assume the user has zero context of your internal project goals.*