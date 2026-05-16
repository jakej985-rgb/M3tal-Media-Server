To: DocSmith, M3TAL Orchestration Lead
From: DocCritic, Senior DevOps Auditor
Subject: Audit Report [M3TAL-ARCH-001] – Deployment Documentation

---

### **Verdict: FAILED**
The current documentation is an architectural whitepaper, not a functional deployment guide. As a new user, I cannot deploy this system. It assumes a "perfect" environment, lacks critical configuration details, and provides no recovery path if the environment assumptions are unmet.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing Repository/Package Validation:** The Quick Start relies on a `curl` and `apt` process that assumes these assets are live and populated. Without verification of the GPG key/APT repo content, the user is left dead-in-the-water if the URL is unreachable or 404.
2.  **Missing `.env` Template:** The orchestrator requires global configuration, but there is no example of the required contents for `/etc/m3tal/.env`. What variables (API keys, DB URLs, persistence paths) must be defined before `m3tal up`?
3.  **Traefik Gateway Omission:** The documentation mentions "Traefik Ownership" but provides zero instructions on how to configure, start, or secure the Traefik entry point. Users need to know which ports (80/443) are required and how Traefik maps to the internal stack.
4.  **Hardware Path Assumptions:** The docs assume `/mnt` and `/mnt/m3tal-media` exist. If a user runs `m3tal up` without these mount points present, will the Docker containers fail to start (Volume bind-mount errors)? 

#### **WARNING**
5.  **Ambiguous `m3tal init`:** It is unclear what `m3tal init` performs. Does it generate secrets? Does it require sudo? Does it create the directories specified in the "Filesystem" table?
6.  **"Source/m3tal-stack" usage:** The table mentions `/opt/m3tal/stack` but doesn't explain if the user is responsible for populating this or if `m3tal init` generates the manifests.

#### **SUGGESTION**
7.  **Diagnostic Feedback:** The documentation mentions `m3tal doctor`. It should be listed as a mandatory post-install step to verify the environment before the first `m3tal up`.
8.  **Port Mapping Table:** Add a table defining exposed ports (e.g., Dashboard, API, Traefik Dashboard) so firewall/security group rules can be configured in advance.

---

### **Required Remediation Plan**

1.  **Environment Setup:** Add a command `m3tal setup-dirs` or clarify that `m3tal init` performs `mkdir -p` on all required paths (`/etc/m3tal`, `/var/lib/m3tal`, `/mnt/m3tal-media`).
2.  **Configuration Reference:** Include a `SAMPLE_ENV` block in the README. Users cannot "guess" the keys required for a Go-based backend.
3.  **Traefik Bootstrapping:** Add a "Network Prerequisites" section. Explain that the user must open port 80/443 and that the Orchestrator will provision the Traefik container automatically upon first run.
4.  **Failure States:** Explicitly state what happens if the APT install fails. Is there a manual binary download path?
5.  **Cleanup the "Ghost" Modules:** Your "Current Repository Status" note admits there is no code. The README should include a "Development" section with instructions for building from source for early adopters, as the APT repo might not be production-ready yet.

**DocCritic’s Final Word:** *Stop describing the architecture and start describing the machine. If a user has to ask "what goes in my .env?", the documentation has failed.*