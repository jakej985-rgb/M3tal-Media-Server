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
    hasDashboard: fs.existsSync("deploy/stack/m3tal-compose.yml"), // Standardized path
    hasCLI: fs.existsSync("cmd/m3tal/main.go"),
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

const servicesState = fs.existsSync("docs/m3tal-services.json")
  ? fs.readFileSync("docs/m3tal-services.json", "utf-8")
  : "[]";

const context = `
Repo Architectural Map:
- Orchestrator (CLI): ${structure.hasCLI ? "cmd/m3tal (present)" : "missing"}
- Infrastructure Stacks: ${structure.hasDocker ? "deploy/stack (compose and traefik configurations)" : "missing"}
- Backend API: ${structure.hasGo ? "Go native (root go.mod)" : "missing"}
- Detected Services/Modules: ${structure.services.join(", ") || "none"}

Docker Services State JSON:
${dockerState}

Environment Variables State JSON:
${envState}

M3TAL Core Architectural Rules:
- M3tal-Core is a Go-native unified binary (\`m3tal\`) acting as the orchestrator control plane.
- It exposes a systemd backend API daemon (\`m3tal-api.service\`) that communicates with the frontend dashboard container.
- Persistent files are isolated strictly per system contract:
  - Configuration: \`/etc/m3tal/.env\`
  - Database & state: \`/var/lib/m3tal/state.db\`
  - Stack files: \`/opt/m3tal/stack\` (user-facing symlink \`/docker\`)
- Dashboard authentication uses abcrypt token store in \`users.json\` (healed automatically on startup).
`;

const targets = [
  {
    path: "README.md",
    name: "README",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate or update the primary README.md for this repository based on its REAL architectural layout.

STRICT RULES:
- Use real structure: Do NOT invent features or directories.
- Relationship Mapping: Explain how the CLI, systemd Backend API service, and Python Dashboard container interact.
- Installation: Emphasize that building from source is NOT needed. Show the APT repository keyring setup and install commands:
  \`\`\`bash
  # Add GPG Key
  curl -fsSL https://jakej985-rgb.github.io/m3tal-core/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/m3tal-archive-keyring.gpg
  # Add Repository
  echo "deb [signed-by=/usr/share/keyrings/m3tal-archive-keyring.gpg] https://jakej985-rgb.github.io/m3tal-core stable main" | sudo tee /etc/apt/sources.list.d/m3tal.list
  # Install M3TAL
  sudo apt update && sudo apt install -y m3tal
  \`\`\`
- Quick Demo: Show how to open the interactive selection menu via \`sudo m3tal\` and configure with \`m3tal config wizard\`.
- Style: Premium, professional, step-by-step newbie start guide. Focus on actionable installation and usage instructions. Speak directly to a user trying to install this for the first time.

${context}
`
  },
  {
    path: "docs/GET_STARTED.md",
    name: "Getting Started Guide",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate a comprehensive, newbie-friendly docs/GET_STARTED.md.

STRICT RULES:
- 10-Minute Setup Journey: Walk the user through APT repository installation, running the configuration wizard (\`sudo m3tal init\` or \`sudo m3tal config wizard\`), and starting the services.
- Ports: Document port \`8082\` (Dashboard web interface) and port \`8080\` (Traefik gateway).
- Filesystem Contract: Explicitly explain the locations:
  - Configuration: \`/etc/m3tal/.env\`
  - Database: \`/var/lib/m3tal/state.db\`
  - Stack directory: \`/opt/m3tal/stack\` (and user-facing symlink \`/docker\`).
- First Login: Emphasize that default credentials are \`admin/admin\` (or custom set during wizard) and how users can reset credentials using \`sudo m3tal dashpass\`.
- Keep it clean, direct, and incredibly premium.

${context}
`
  },
  {
    path: "docs/COMMAND_REFERENCE.md",
    name: "Command Line Reference",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate a professional, structured cheat-sheet docs/COMMAND_REFERENCE.md.

STRICT RULES:
- Document every CLI command in the \`m3tal\` binary:
  - \`m3tal\`: Launches the TUI-based interactive Control Center selection menu.
  - \`m3tal init\`: Generates default environment configurations.
  - \`m3tal doctor\`: Runs comprehensive pre-flight health checks.
  - \`m3tal dashpass [username] [password]\`: Manages dashboard auth credentials (interactive fallback if args omitted).
  - \`m3tal config\`: Subcommands (\`wizard\`, \`set\`, \`get\`, \`scan\`, \`list\`).
  - \`m3tal dash\`: Manages dashboard container state (\`up\`, \`down\`, \`start\`, \`stop\`, \`restart\`, \`logs\`, \`status\`).
  - \`m3tal up\` / \`m3tal down\`: Controls the environment stacks mapped in \`/docker\`.
  - \`m3tal logs\`: Streams aggregated container logs.
- Provide direct command line examples for each option, including their descriptions and flags.

${context}
`
  },
  {
    path: "docs/ENVIRONMENT_VARIABLES.md",
    name: "Environment Variables Reference",
    prompt: `You are DocSmith, the M3TAL Ecosystem Documentation Architect.
Generate a clean, tabular, and highly structured reference docs/ENVIRONMENT_VARIABLES.md.

STRICT RULES:
- Read the environment variable state JSON.
- For EVERY variable in the JSON, create a premium documentation block detailing:
  - Name (e.g. \`BASE_STORAGE_PATH\`, \`DASHBOARD_PORT\`, \`ADMIN_PASSWORD\`, \`CF_TUNNEL_TOKEN\`)
  - Description of what it controls
  - Default value
  - Common configuration values
- Emphasize which values are automatically generated by the system on first boot (\`DASHBOARD_SECRET\`, \`API_TOKEN\`).

${context}
`
  }
];

const modelsToTry = [
  "gemini-3-flash",          // Tier 1: Primary (5 RPM / 20 RPD)
  "gemini-3.1-flash-lite",   // Tier 2: High-Volume Fallback (15 RPM / 500 RPD)
  "gemini-2.5-flash",        // Tier 3: Secondary Fallback (5 RPM / 20 RPD)
  "gemini-2.5-flash-lite",   // Tier 3: Tertiary Fallback (10 RPM / 20 RPD)
  "gemini-1.5-flash",        // Legacy Fallback
  "gemini-pro"               // Universal Fallback
];

async function generateWithFallback(target) {
  for (const modelName of modelsToTry) {
    try {
      console.log(`[${target.name}] Attempting to use model: ${modelName}`);
      const model = genAI.getGenerativeModel({ model: modelName });
      const result = await model.generateContent(target.prompt);
      const output = result.response.text();
      fs.writeFileSync(target.path, output);
      console.log(`[${target.name}] File ${target.path} updated successfully using ${modelName}.`);
      return;
    } catch (error) {
      console.warn(`[${target.name}] Model ${modelName} failed: ${error.message}`);
    }
  }
  console.error(`[${target.name}] All Gemini models failed to generate content.`);
}

async function run() {
  for (const target of targets) {
    await generateWithFallback(target);
  }
  console.log("All documentation targets updated successfully!");
}

run();
