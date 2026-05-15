### **DocCritic Audit Report: M3TAL Core**

**Verdict: FAILED (BLOCKER)**
As a new user, I cannot successfully deploy this system. The documentation assumes a "magic" environment where the user knows how to configure environment variables, where the system expects data to live, and how the network layer is exposed. The current README is a high-level architectural overview, not an operator's manual.

---

### **Issue List**

#### **BLOCKER**
1. **Missing `.env` / Configuration Instructions:** The documentation mentions `BASE_STORAGE_PATH` and `API_TOKEN`, but there is no guide on *where* to create an environment file or how the system reads configuration. If I run `m3tal init`, it will likely fail or generate invalid configs because it lacks these inputs.
2. **Missing Port/Access Documentation:** I am told the `Traefik Gateway` is deployed, but I have no idea what port to hit in my browser. Is it 80? 443? 8080? If I don't know the port, I cannot access the dashboard.
3. **Hard-coded Assumption of `/mnt`:** The "Path Consistency Rule" mandates `/mnt`. If a user is on a system where `/mnt` is restricted or unused (e.g., a NAS mount at `/data`), the user has no instructions on how to alias or satisfy this constraint without breaking the orchestrator.

#### **WARNING**
4. **Ambiguous `m3tal up` dependencies:** Does `m3tal up` require `docker-compose` to be pre-installed on the host? Does it use the Docker Python SDK or call the `docker-compose` binary? The distinction between the Go orchestrator and the `deploy/stack` manifests is unclear regarding dependencies.
5. **Lack of "Day 0" Setup:** The "Quick Start" is too fast. It ignores the required initial setup of the `m3tal` configuration file. Where does `m3tal.yaml` (or similar) live? `/etc/m3tal/`? `~/.config/m3tal/`?

#### **SUGGESTION**
6. **"Doctor" Command Transparency:** Since you provide a `m3tal doctor` command, promote it in the "Quick Start" section. It’s the user's best friend for debugging the missing environment variables mentioned above.

---

### **Suggested Fixes**

*   **Fix for #1:** Add a section: **"Initial Configuration"**.
    > Create a file at `/etc/m3tal/config.env`:
    > ```bash
    > BASE_STORAGE_PATH=/home/user/media
    > API_TOKEN=generate_a_secure_random_string
    > ```
    > Mention that `m3tal init` pulls from this location.

*   **Fix for #2:** Add a **"Networking & Access"** section:
    > "By default, M3TAL exposes the Dashboard on port `80` (HTTP) and `443` (HTTPS). Ensure these ports are open on your host firewall."

*   **Fix for #3:** Clarify the "Path Consistency Rule":
    > "The orchestrator mandates that your `BASE_STORAGE_PATH` is mapped to `/mnt` inside the container. If your media is stored at `/data`, the orchestrator will automatically bind `/data` to `/mnt` within the container context."

*   **Fix for #4 & #5:** Update "Quick Start" with a logical flow:
    1. `sudo apt install m3tal-core`
    2. `sudo m3tal config set path /path/to/media`
    3. `sudo m3tal init`
    4. `m3tal doctor` (Verify all systems green)
    5. `m3tal up`

*   **Documentation Meta-Improvement:** Add a "Requirements" header specifying that `Docker Engine` and `Docker Compose V2` must be installed prior to running the `m3tal-core` package.