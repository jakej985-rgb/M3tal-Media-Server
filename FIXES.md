To: M3TAL Engineering Team
From: DocCritic, Senior DevOps Auditor
Subject: Audit Report - M3TAL Media Server Documentation

### Verdict: **FAIL**
The current documentation is a "Happy Path" abstraction that fails to account for the actual requirements of a production-grade infrastructure tool. A new user following these instructions will encounter immediate environment failures, path conflicts, and connectivity dead-ends.

---

### Detailed Issue List

#### **BLOCKER**
1.  **Missing Environment Bootstrap**: The `Quick Start` implies `./m3tal init` works out-of-the-box, but there is zero instruction on how to create the initial `.env` file or where the template resides. Users will face "File Not Found" errors immediately.
2.  **Hardcoded Path Assumption**: The documentation demands `BASE_STORAGE_PATH` be mounted to `/mnt` inside containers, but provides no instruction on how to configure this on the host. If the user does not have a `/mnt` directory or the necessary permissions, the containers will fail to start.
3.  **Dependency Ambiguity**: You mention `m3tal-goback` is an "external" requirement. The docs do not define how to link this. Without a clear configuration step for the API endpoint and token, the dashboard will be a broken UI upon first load.

#### **WARNING**
4.  **Networking Opacity**: You reference a "Traefik Gateway," but there are no instructions on which ports must be open on the host firewall. A user doesn't know if they should point their browser to `localhost:80`, `localhost:8080`, or something else.
5.  **Lack of Docker Requirements**: The document assumes Docker and Docker Compose (v2) are installed. It should explicitly state the engine requirements.

#### **SUGGESTION**
6.  **"Quick Start" vs "Setup"**: The current flow mixes compilation and initialization. Separate "Prerequisites" from "Deployment."
7.  **Diagnostic Feedback**: Since you have a `m3tal doctor` command, tell the user to run this *before* they try to go live.

---

### Suggested Fixes

**1. Add a "Prerequisites" section:**
> *   **Docker Engine** (v20.10+) with Compose V2 plugin enabled.
> *   **Go 1.21+** installed.
> *   **Host Directory**: Ensure you have a target directory for media (e.g., `/mnt/media`).

**2. Add "Configuration" step before "Quick Start":**
> 1. Copy the template: `cp .env.example .env`
> 2. Edit `.env` to set your `BASE_STORAGE_PATH` (e.g., `/home/user/media`).
> 3. Run `./m3tal config` to validate your environment variables.

**3. Clarify Networking:**
> Add a note: "Traefik is configured to bind to port 80/443 on the host. Ensure these ports are available. Access your dashboard at `http://localhost` after running `./m3tal up`."

**4. Update Quick Start sequence:**
```bash
# 1. Prerequisites (Check if Docker/Go are installed)
# 2. Configure Environment
cp .env.example .env
# 3. Compile
go build -o m3tal main.go
# 4. Initialize and Validate
./m3tal init
./m3tal doctor # Verify host paths and Docker status
# 5. Launch
./m3tal up
```

**5. Backend Connectivity:**
> Add a dedicated sub-section: **Connecting the Backend**
> "The Dashboard requires `m3tal-goback`. Ensure `GOBACK_URL` and `API_TOKEN` are set in your `.env` file before executing `up`."