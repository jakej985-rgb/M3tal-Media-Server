## DocCritic Audit Report: M3TAL Platform (v1.7)

**Verdict:** **REJECTED.**
The current documentation assumes an "expert-in-the-room" level of knowledge and leaves critical operational gaps that will cause a new user's first deployment to fail immediately. The orchestration logic (Go vs. Bash vs. Docker Compose) is opaque, and the path/permission requirements are brittle.

---

### 🚨 Issue List

#### 1. BLOCKER: Missing `.env` Variable Definition
The `template.env` is mentioned but not documented. A user has no idea what `BASE_STORAGE_PATH` or other required keys are supposed to look like.
*   **Fix:** Include a table in `README.md` or a link to a completed `example.env` explaining every required key (e.g., `BASE_STORAGE_PATH`, `API_KEY`, `DB_PASSWORD`).

#### 2. BLOCKER: Ambiguous `m3tal.py` vs Binary usage
The intro mentions a "Python 3.10 Legacy Dashboard" but the Quick Start uses a Go binary `./m3tal`. Is the Python service managed by the Go binary? Does the user need to run `pip install` for the legacy service?
*   **Fix:** Clarify the dependency chain. If the Go binary handles the Python service, state it clearly. If manual intervention is needed, provide the `pip install -r requirements.txt` steps.

#### 3. BLOCKER: The `/mnt` Directory Assumption
The documentation mandates an absolute path for `BASE_STORAGE_PATH` but warns that an empty `/mnt` inside the container causes failure. It does not explain how the `m3tal` binary validates this path or what error message the user should expect.
*   **Fix:** Add a validation step in the CLI. The `m3tal init` command should perform a pre-flight check: "Does `BASE_STORAGE_PATH` exist? Is it writable? If not, throw a clean error and exit."

#### 4. WARNING: Traefik Network Configuration
You specify `traefik.localhost` and `m3tal.localhost`, but you don't explain *how* the Traefik container actually binds to the host ports. Does the user need to open ports 80/443 on the host? What if they are already in use (e.g., by Apache or Nginx)?
*   **Fix:** Document the required Host Ports (80/443/8080) and advise on potential port conflicts.

#### 5. WARNING: Build Process Ambiguity
`./build.sh` is a "black box." Does it compile the Go binary only? Does it pull Docker images? Does it build the Python dashboard?
*   **Fix:** Document the output of `build.sh`. What should the user see if it succeeds? Where is the resulting binary placed?

#### 6. SUGGESTION: "Clean Room" Initialization
The "Quick Start" assumes the user can just run `./m3tal init`. It does not explicitly mention that this generates the underlying Docker Compose files in `source/m3tal-stack/`.
*   **Fix:** Add a "What to expect" section. Explain that `init` generates `docker-compose.yml` from a template.

---

### 📝 Recommended Immediate Actions for Authors:

1.  **Refine the `template.env`**: Add comments to every line in the file. A user shouldn't have to guess what `API_KEY` is for.
2.  **Explicit Path Instructions**: Change the instruction:
    *   *From:* `mkdir -p /path/to/your/media`
    *   *To:* "Create a folder (e.g., `/home/user/media`), then set `BASE_STORAGE_PATH=/home/user/media` in your `.env` file."
3.  **Add a Pre-Flight Command**: Add a `./m3tal doctor` command. This is standard in modern CLI tools. It should print:
    *   Docker: [OK]
    *   .env: [OK]
    *   /mnt access: [OK]
    *   Traefik Ports: [80: Available]
4.  **Clarify Ingress**: Explicitly state: "Ensure ports 80 and 443 are not currently occupied by another web server."

**DocCritic Note:** *You have a powerful project, but your current README expects the user to be a mind-reader. Fix the installation path and environment configuration documentation, or you will spend all your time answering "Why is my container empty?" issues.*