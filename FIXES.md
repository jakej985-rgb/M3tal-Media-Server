### **DocCritic Audit Report: M3TAL Core Orchestrator**

**Verdict:** 🔴 **NON-DEPLOYABLE**
This documentation reads like a marketing brochure rather than a technical manual. It fails to answer the "How" for a new user. There is no installation guide, no configuration schema, and the project currently assumes the user is an expert who already knows the internal secrets of the M3TAL architecture.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing Installation Steps:** There is no `go build` instruction, no binary release download info, and no guidance on how to install the `m3tal` CLI to `/usr/bin/`.
2.  **Missing Environment Configuration:** `/etc/m3tal/.env` is defined as the "Source of Truth," but there is no template or list of required environment variables (e.g., API keys, database URLs, port mappings).
3.  **Docker Orchestration Gaps:** The provided Docker snippet is orphaned. There is no instruction on how to initialize the `/opt/m3tal/stack` directory, nor a `docker-compose.yml` that links the Core, API, and Dashboard.
4.  **Assumed Directory Structure:** The guide assumes `/mnt/m3tal-media` and `/var/lib/m3tal/` exist. If a user runs the Docker command without creating these, it will trigger permission errors or mount host filesystem noise.

#### **WARNING**
5.  **Lack of Port/Gateway Details:** The documentation mentions "Traefik" in the prompt context but provides zero info on how to expose the services. How does the user access the Dashboard? What ports must be open?
6.  **"Core-First" Ambiguity:** It is unclear if the `m3tal` binary *runs* the stack or if the stack *runs* the binary. The architecture description contradicts the Docker snippet.

#### **SUGGESTION**
7.  **Missing Quickstart/Bootstrapping:** Add a `m3tal init` or similar command explanation to automate the creation of the required directory structure.

---

### **Suggested Fixes**

**1. Add a "Getting Started" Section:**
*   Add a shell script block for installation:
    ```bash
    git clone ...
    go build -o m3tal ./cmd/m3tal
    sudo cp m3tal /usr/bin/
    sudo mkdir -p /etc/m3tal /var/lib/m3tal
    ```

**2. Provide an `.env.example`:**
Create a section detailing mandatory variables:
*   `M3TAL_API_KEY`: Secret for inter-module communication.
*   `M3TAL_DASH_PORT`: Port for the web interface.
*   `M3TAL_MEDIA_PATH`: Host path for media.

**3. Provide a `docker-compose.yml` Template:**
Don't just show the `m3tal-core` service. Provide a full composition that links `m3tal-goback` and `m3tal-godash` using Docker networking:
```yaml
services:
  core:
    image: m3tal/core
    volumes:
      - /var/lib/m3tal:/var/lib/m3tal
  api:
    image: m3tal/goback
    depends_on: [core]
  dashboard:
    image: m3tal/godash
    ports:
      - "80:80" # Map to Traefik/Host
```

**4. Define "First Run" Procedures:**
Add a section: *"Before starting the containers, run `m3tal setup` to generate the default configuration manifest in `/opt/m3tal/stack`."* (If this doesn't exist, create it—it is vital for automation).

**5. Clear Path Assumptions:**
Explicitly tell the user: "Ensure your host machine has read/write permissions for `/mnt/m3tal-media`. If this directory does not exist, the orchestrator will fail to initialize the media library."