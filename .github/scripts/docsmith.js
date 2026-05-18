import fs from "fs";
import path from "path";
import { GoogleGenerativeAI } from "@google/generative-ai";

const genAI = new GoogleGenerativeAI(process.env.GOOGLE_AI_KEY);

// Ensure docs directory exists
fs.mkdirSync("docs", { recursive: true });

// --- Detect M3TAL repo structure ---
function detectStructure() {
  const structure = {
    hasDocker: fs.existsSync("deploy/stack/m3tal-compose.yml"),
    hasGo: fs.existsSync("go.mod"),
    hasCLI: fs.existsSync("cmd/m3tal/main.go"),
    hasTraefik: fs.existsSync("deploy/stack/traefik.yml"),
    hasRouting: fs.existsSync("deploy/stack/routing-compose.yml"),
    hasCloudflared: fs.existsSync("deploy/stack/cloudflared-config.yml"),
    services: [],
  };

  if (fs.existsSync("deploy")) {
    const sourceDirs = fs.readdirSync("deploy", { withFileTypes: true })
      .filter(d => d.isDirectory())
      .map(d => d.name);
    structure.services = sourceDirs;
  }

  return structure;
}

const structure = detectStructure();

// Load dynamic parsed context from previous parser steps
const dockerState = fs.existsSync("docs/docker.json")
  ? fs.readFileSync("docs/docker.json", "utf-8")
  : "[]";

const envState = fs.existsSync("docs/env.json")
  ? fs.readFileSync("docs/env.json", "utf-8")
  : "[]";

// Read real Traefik config files for DocCritic context grounding
const traefikStatic = fs.existsSync("deploy/stack/traefik.yml")
  ? fs.readFileSync("deploy/stack/traefik.yml", "utf-8")
  : "(not found)";

const traefikDynamic = fs.existsSync("deploy/stack/dynamic/api.yml")
  ? fs.readFileSync("deploy/stack/dynamic/api.yml", "utf-8")
  : "(not found)";

// ─────────────────────────────────────────────────────────────────────────────
// SYSTEM CONTEXT — complete picture of the REAL M3TAL architecture
// This context is injected into every prompt so Gemini writes accurate docs.
// ─────────────────────────────────────────────────────────────────────────────
const context = `
## M3TAL System Architecture (Ground Truth)

### Components
- **CLI binary** (\`/usr/bin/m3tal\`): Unified Go binary installed via APT. Single entrypoint for all operations.
- **API daemon** (\`m3tal-api.service\`): Go binary running as a systemd service on port 8080. Manages Docker, state DB, and API routes.
- **Dashboard container** (\`m3tal-dashboard\`): Python/Flask container running on port 8082. Communicates with the API daemon at \`http://host.docker.internal:8080\`.
- **Traefik gateway** (\`routing-compose.yml\`): Reverse proxy container exposing services by domain name on port 80. Uses file provider for dynamic routing.
- **Cloudflared** (\`routing-compose.yml\`): Optional Cloudflare tunnel container for zero-config internet access.

### Filesystem Contract (MUST document explicitly)
| Path | Purpose |
|------|---------|
| \`/etc/m3tal/.env\` | Primary configuration file. Managed by \`m3tal config wizard\`. |
| \`/var/lib/m3tal/state.db\` | SQLite state database. Auto-created by the API daemon. |
| \`/opt/m3tal/stack/\` | Canonical stack directory. Contains compose files and Traefik config. |
| \`/docker\` | Symlink → \`/opt/m3tal/stack/\`. This is the user-facing path for all stack operations. |
| \`/docker/users.json\` | Dashboard credential store. Managed by \`m3tal dashpass\`. |

### Docker / Compose Runtime (BLOCKER FIX — must explain clearly)
- M3TAL uses **Docker Engine + Docker Compose V2** under the hood. These are hard dependencies.
- The \`m3tal up\` command runs \`docker compose\` across all \`*-compose.yml\` files found in \`/docker/\`.
- The \`m3tal dash up\` command specifically manages the dashboard container via \`/docker/m3tal-compose.yml\`.
- User stacks live in \`/docker/\`. Adding a new stack means placing a \`*-compose.yml\` file there.
- The compose files use a shared env file at \`/etc/m3tal/.env\` via the \`--env-file\` flag.

### Deployment Lifecycle — Day 2 Operations (BLOCKER FIX)
Installing a new stack:
1. Place your compose file in \`/docker/my-stack-compose.yml\`
2. Ensure required variables are set in \`/etc/m3tal/.env\` (use \`m3tal config wizard\` or \`m3tal config set KEY value\`)
3. Run \`m3tal up\` to start all stacks, or \`docker compose -f /docker/my-stack-compose.yml up -d\` for a single stack.

### Traefik Routing Architecture (BLOCKER FIX — explain the gateway)
Traefik is deployed as a container via \`routing-compose.yml\`. It:
- Binds port 80 on the host as the HTTP entry point.
- Discovers services automatically via Docker labels (e.g. \`traefik.http.routers.myapp.rule=Host(\\\`app.domain.com\\\`)\`).
- Loads dynamic config from \`/docker/dynamic/\` (file provider, hot-reload).
- Routes \`api.DOMAIN\` → \`http://host.docker.internal:8080\` (the Go API daemon) via \`dynamic/api.yml\`.
- Routes \`dash.DOMAIN\` → the dashboard container on port 8082 via labels in \`m3tal-compose.traefik.yml\`.
- The Traefik dashboard itself is accessible at \`http://localhost:8081\` (local only).

Traefik static config (traefik.yml):
\`\`\`yaml
${traefikStatic}
\`\`\`

Dynamic routing example (dynamic/api.yml):
\`\`\`yaml
${traefikDynamic}
\`\`\`

### Service Management — systemd (WARNING FIX)
- The API daemon is managed by systemd as \`m3tal-api.service\`.
- Commands: \`systemctl status m3tal-api\`, \`systemctl restart m3tal-api\`, \`journalctl -u m3tal-api -f\`
- The CLI itself (\`m3tal\`) also manages the daemon: \`m3tal dash up\` triggers \`systemctl start m3tal-api\`.

### Port Map
| Port | Service | Access |
|------|---------|--------|
| 80 | Traefik HTTP entry point | Public |
| 8080 | M3TAL API daemon (Go) | Host-local (via Traefik or direct) |
| 8081 | Traefik dashboard | Host-local only |
| 8082 | M3TAL Dashboard (Python/Flask) | Via Traefik or direct |

### APT Installation (ALWAYS include this exact block)
\`\`\`bash
# 1. Add the GPG signing key
curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg

# 2. Add the APT repository
echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list

# 3. Install
sudo apt update && sudo apt install -y m3tal
\`\`\`

### Docker Services State JSON:
${dockerState}

### Environment Variables State JSON:
${envState}
`;

// ─────────────────────────────────────────────────────────────────────────────
// DOCUMENTATION TARGETS
// ─────────────────────────────────────────────────────────────────────────────
const targets = [
  {
    path: "README.md",
    name: "README",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate the primary README.md based on the REAL system architecture documented below.

STRICT RULES:
- NO marketing copy. No "robust", "cohesive", "autonomous" buzzwords. Write like a DevOps engineer, not a salesperson.
- NO placeholders. Every section must contain real, actionable content.
- Explain the actual runtime: M3TAL uses Docker Engine + Docker Compose V2. State this explicitly.
- MUST include the exact APT installation block from the context below. Do not invent a different install method.
- MUST include a "Deployment Lifecycle" section explaining:
  - How stacks work: compose files in /docker/, brought up with \`m3tal up\`
  - How to add a new stack (place compose file in /docker/, run \`m3tal up\`)
- MUST include a "Traefik Gateway" section explaining:
  - Traefik runs as a container, binds port 80
  - Services are exposed by adding Traefik labels to their compose service definition
  - \`api.DOMAIN\` routes to the Go API daemon via dynamic config
  - \`dash.DOMAIN\` routes to the dashboard container
- MUST include a "Service Management" section showing systemctl commands for m3tal-api.service
- MUST include a Quick Demo section with: install → config wizard → m3tal dash up → open browser
- MUST include a port table (80, 8080, 8081, 8082)
- Firewall note: remind users to allow port 80 in ufw/iptables for Traefik

${context}
`
  },
  {
    path: "docs/GET_STARTED.md",
    name: "Getting Started Guide",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate docs/GET_STARTED.md — a complete, newbie-friendly setup guide for first-time users.

STRICT RULES:
- NO marketing copy. Write clear operational steps only.
- Step 1: Prerequisites — explicitly state: "Docker Engine and Docker Compose V2 must be installed."
  Show: \`docker --version && docker compose version\`
- Step 2: Install M3TAL via APT (include the exact 3-command APT block from context).
- Step 3: Run the configuration wizard: \`sudo m3tal config wizard\` — explain what each prompt means.
- Step 4: Start the routing stack (Traefik): \`m3tal up\` — explain this starts all compose files in /docker/.
- Step 5: Start the dashboard: \`m3tal dash up\` — explain it pulls the image and starts the container.
- Step 6: Open browser at \`http://YOUR_IP:8082\` (or \`http://dash.DOMAIN\` if Traefik is configured).
- Step 7: Log in — default credentials, how to change with \`sudo m3tal dashpass\`.
- Filesystem Contract section: document /etc/m3tal/.env, /var/lib/m3tal/state.db, /docker symlink.
- Port table: 80 (Traefik), 8080 (API), 8082 (Dashboard).
- Firewall note: \`sudo ufw allow 80\` if Traefik is exposed.
- Service management: show \`systemctl status m3tal-api\`, \`journalctl -u m3tal-api -f\`.

${context}
`
  },
  {
    path: "docs/COMMAND_REFERENCE.md",
    name: "Command Line Reference",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate docs/COMMAND_REFERENCE.md — a complete CLI cheat-sheet.

STRICT RULES:
- Document every command with a real usage example. No placeholders.
- Include these commands and their subcommands:
  - \`sudo m3tal\` → Opens the interactive TUI Control Center (numbered menu).
  - \`m3tal init\` → Generates /etc/m3tal/.env from defaults. Use on first install.
  - \`m3tal doctor\` → Pre-flight health check: Docker connectivity, .env validity, port availability.
  - \`m3tal config wizard\` → Interactive wizard to configure /etc/m3tal/.env.
  - \`m3tal config set KEY VALUE\` → Set a single env var.
  - \`m3tal config get KEY\` → Read a single env var.
  - \`m3tal config scan\` → List all env vars across all stacks.
  - \`m3tal config list\` → List current .env file contents.
  - \`m3tal dashpass [username] [password]\` → Update dashboard user password. Interactive if args omitted.
  - \`m3tal dash up\` → Pull latest dashboard compose config from GitHub, then start the dashboard container.
  - \`m3tal dash down\` → Stop the dashboard container.
  - \`m3tal dash restart\` → Restart the dashboard container.
  - \`m3tal dash logs\` → Stream dashboard container logs.
  - \`m3tal dash status\` → Show dashboard container status.
  - \`m3tal up\` → Run docker compose up across all *-compose.yml files in /docker/.
  - \`m3tal down\` → Run docker compose down across all stacks.
  - \`m3tal logs\` → Stream aggregated logs from all running stacks.
- Include a section on systemd service management: \`systemctl status m3tal-api\`, \`journalctl -u m3tal-api -f\`.
- Include a Docker section showing direct compose commands as fallback.

${context}
`
  },
  {
    path: "docs/ENVIRONMENT_VARIABLES.md",
    name: "Environment Variables Reference",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate docs/ENVIRONMENT_VARIABLES.md — a complete environment variable reference.

STRICT RULES:
- Use the Environment Variables State JSON from the context as your source of truth.
- For EVERY variable in the JSON, document:
  - Name, description, default value, example value, which component uses it.
- Note that DASHBOARD_SECRET and API_TOKEN are auto-generated on first \`m3tal init\` — users should NOT set them manually unless rotating.
- Note that BASE_STORAGE_PATH controls where media data is stored and defaults to /mnt in production deployments (not ./data as in the template).
- Note that DOMAIN controls Traefik routing rules — setting it enables \`dash.DOMAIN\` and \`api.DOMAIN\` routes.
- Explain that all variables are read from \`/etc/m3tal/.env\` by both the CLI and all compose stacks via --env-file.
- Group variables logically: Core, Auth, Network, Storage, Traefik, VPN, System.
- Include a quick-reference table at the top.

${context}
`
  }
];

// ─────────────────────────────────────────────────────────────────────────────
// MODEL FALLBACK CHAIN
// ─────────────────────────────────────────────────────────────────────────────
const modelsToTry = [
  "gemini-2.5-flash",        // Primary
  "gemini-2.5-flash-lite",   // Fallback
  "gemini-1.5-flash",        // Legacy
  "gemini-pro"               // Universal
];

async function generateWithFallback(target) {
  for (const modelName of modelsToTry) {
    try {
      console.log(`[${target.name}] Attempting to use model: ${modelName}`);
      const model = genAI.getGenerativeModel({ model: modelName });
      const result = await model.generateContent(target.prompt);
      const output = result.response.text();
      fs.writeFileSync(target.path, output);
      console.log(`[${target.name}] ✅ ${target.path} updated using ${modelName}.`);
      return;
    } catch (error) {
      console.warn(`[${target.name}] Model ${modelName} failed: ${error.message}`);
    }
  }
  console.error(`[${target.name}] ❌ All Gemini models failed to generate content.`);
}

async function run() {
  for (const target of targets) {
    await generateWithFallback(target);
  }
  console.log("✅ All documentation targets updated.");
}

run();
