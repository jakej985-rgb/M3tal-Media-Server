# DocCritic Audit Report: M3TAL Core Orchestrator

**Auditor:** DocCritic, Senior DevOps Auditor  
**Verdict:** **REJECTED - NEEDS IMMEDIATE REVISION**

The documentation currently suffers from significant "developer-bias." It assumes the user has existing infrastructure knowledge, fails to address basic environment requirements (like mounting), and lacks critical networking details required for a successful deployment.

---

### Detailed Issue List

#### 1. BLOCKER: Missing Filesystem Initialization
The `Filesystem Standard` section mandates `/mnt/m3tal-media`, but the installation guide does not mention creating this directory or mounting a drive. If the user runs `m3tal up` without this path existing, the Docker volume mapping will fail (or create a root-owned directory), breaking the system.
*   **Fix:** Add a section "System Preparation" requiring the user to create the mount point and set permissions: `sudo mkdir -p /mnt/m3tal-media && sudo chown $USER:$USER /mnt/m3tal-media`.

#### 2. BLOCKER: Missing Port/Access Information
There is zero information regarding how to access the Dashboard or API after running `m3tal dash up`. I have no idea if I should visit `http://localhost:8080` or if a Traefik gateway is expecting specific host headers.
*   **Fix:** Add a "Connectivity" section defining default ports (e.g., Dashboard 80, API 8080) and any Traefik labels required for local access.

#### 3. WARNING: Ambiguous Docker Interaction
The "Deployment: Docker Configuration" snippet is presented as if the user should paste it into a file, but it doesn't specify *where* that file goes or if it's meant to be managed by `m3tal up`.
*   **Fix:** Clarify that `m3tal up` automatically triggers the stack found in `/opt/m3tal/stack/docker-compose.yml`. Explain how to modify that file if the user needs to adjust network configurations.

#### 4. WARNING: Hidden Assumptions regarding `m3tal setup`
`m3tal setup` assumes the user has write access to `/etc/m3tal`. On most Linux distros, a standard user cannot write to `/etc/`.
*   **Fix:** Explicitly state if the CLI needs `sudo` or if the user needs to create the directory with specific permissions first.

#### 5. SUGGESTION: Overly Dramatic Language
The README uses phrases like "Go-native architectural requirements" and "Ecosystem Integration Rules." This sounds like marketing fluff.
*   **Fix:** Remove buzzwords. Replace with "Requirements" and "Deployment Standards." The focus should be on *what the system does*, not *why the architecture is superior*.

#### 6. SUGGESTION: Missing "Troubleshooting/Log" Section
The documentation provides no way to debug a failed `m3tal up` command.
*   **Fix:** Add a section on how to check logs (`m3tal logs` or standard `docker compose logs`).

---

### Suggested README Structure Improvements

1.  **System Prep:** Command to create and chown `/mnt/m3tal-media`.
2.  **Installation:** (Keep your APT section, it is technically correct, but clarify the `sudo` requirement for `m3tal setup`).
3.  **Quickstart:**
    *   `m3tal setup`
    *   `m3tal up`
    *   *Add: Verify status via `m3tal dash up`*
4.  **Networking & Access:**
    *   Default Dashboard URL: `http://localhost`
    *   API Endpoint: `http://localhost:8080`
    *   Explain if Traefik is handled automatically.
5.  **Troubleshooting:**
    *   `m3tal --help`
    *   Location of logs.

**DocCritic Note:** *Stop treating the README like a technical manifesto and start treating it like a support manual for a busy SysAdmin. Make it dry, precise, and actionable.*