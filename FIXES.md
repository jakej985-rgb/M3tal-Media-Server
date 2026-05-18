**Audit Report: M3TAL Core Orchestrator README**
**Auditor:** DocCritic, Senior DevOps Auditor
**Date:** 2024-05-22

---

### **Verdict: FAILED**
The documentation is insufficient for a production-grade deployment. While the APT installation instructions are present, the project suffers from "assume-it-works" syndrome, particularly regarding storage prerequisites, network exposure, and the operational lifecycle. The documentation prioritizes marketing-speak over actionable engineering instructions.

---

### **Issue List**

#### **BLOCKER**
1. **Unmet Dependency Assumptions (`/mnt`):** The documentation assumes `/mnt/m3tal-media` exists. If a user runs `m3tal up` without this path, the Docker container will likely crash or create a directory owned by `root`.
2. **Missing Network/Port Documentation:** No mention of the Traefik/Gateway ports. A user will be unable to access the dashboard after running `m3tal dash up` because they don't know which port to hit or how to configure ingress.

#### **WARNING**
1. **Vague Docker Stack Instructions:** You mention `deploy/stack` and `deploy/dashboard` in the CLI commands but provide no guide on how to deploy these if the user isn't using the CLI, or what environment variables are required for the `docker-compose.yml` files.
2. **"Go-Native" Buzzword Overload:** The README spends more time highlighting the rewrite than explaining how to troubleshoot the orchestration logic.

#### **SUGGESTION**
1. **Marketing Fluff:** Phrases like "Modular Infrastructure Platform" and "Go-Native Migration" add no technical value. Remove these to keep the README concise.
2. **Binary Verification:** The APT instructions don't mention checking the binary version or the system service status (e.g., `systemctl` if applicable).

---

### **Suggested Fixes**

#### **1. Addressing Storage (BLOCKER)**
Add a "System Preparation" section before "Installation":
> **System Preparation**
> M3TAL requires a persistent media mount. Ensure this directory exists and is accessible:
> ```bash
> sudo mkdir -p /mnt/m3tal-media
> # Replace 'user' with your actual non-root user
> sudo chown -R $USER:$USER /mnt/m3tal-media
> ```

#### **2. Adding Network/Access Info (BLOCKER)**
Add an "Accessing the Platform" section:
> **Access & Ports**
> By default, the M3TAL dashboard and API are exposed via the following ports:
> - **Dashboard:** `http://localhost:8080` (or configured Traefik port)
> - **API Gateway:** `http://localhost:9000`
> Ensure your firewall allows incoming traffic on these ports if hosting externally:
> `sudo ufw allow 8080/tcp`

#### **3. Improving Docker Clarity (WARNING)**
Clarify the relationship between the CLI and the Compose files:
> **Deployment Notes**
> The `m3tal up` command executes `docker-compose -f /opt/m3tal/stack/docker-compose.yml up -d`. If you need to troubleshoot the stack manually, use:
> `docker compose -f /opt/m3tal/stack/docker-compose.yml ps`

#### **4. Removing Marketing Fluff (SUGGESTION)**
*   **Remove:** "The ecosystem is transitioning away from legacy wrappers." (Keep documentation forward-looking).
*   **Remove:** "M3TAL — Modular Infrastructure Platform." (The title at the top is sufficient).
*   **Add:** A "Troubleshooting" section that lists `m3tal status` or log locations (e.g., `/var/log/m3tal.log`).

---

**Auditor Note:** *DocSmith, clean up the technical gaps. Users don't care that it's "Go-native"; they care that it doesn't crash when they mount the storage. Fix the `/mnt` assumption immediately.*