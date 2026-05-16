import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";
import { GoogleGenerativeAI } from "@google/generative-ai";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, "../../");

const genAI = new GoogleGenerativeAI(process.env.GOOGLE_AI_KEY);

const readmePath = path.join(rootDir, "README.md");
const readme = fs.readFileSync(readmePath, "utf-8");

const prompt = `
You are DocCritic, a Senior DevOps Auditor for the M3TAL platform.

Act like a NEW USER trying to run this project. Be extremely strict.

Check for:
- Missing install steps (make build, m3tal setup, .env configuration).
- Missing Docker instructions (deploy/stack usage).
- Missing ports / access info (Traefik gateway).
- Confusing wording or technical gaps.
- Dev-only assumptions (e.g., assuming /mnt already exists).
- Overly dramatic marketing copy or buzzwords (documentation should be a technical guide, not a sales pitch).

Classify issues:
- BLOCKER (Project cannot be deployed with this documentation)
- WARNING (Documentation is confusing or incomplete)
- SUGGESTION (Improvements for clarity)

Return:
- A clear Verdict.
- Detailed issue list with "BLOCKER", "WARNING", or "SUGGESTION" prefixes.
- Suggested fixes for every issue.

README:
${readme}
`;

async function run() {
  const modelsToTry = [
    "gemini-3.1-flash-lite",   // Tier 1: High-Volume Primary (15 RPM / 500 RPD)
    "gemini-3-flash-preview",  // Tier 1: Performance Primary (5 RPM / 20 RPD)
    "gemini-2.5-flash",        // Tier 2: Secondary Fallback
    "gemini-2.0-flash",        // Tier 3: Stable Fallback
    "gemini-1.5-flash"         // Tier 4: Legacy Fallback
  ];

  for (const modelName of modelsToTry) {
    try {
      console.log(`Attempting to use model: ${modelName}`);
      const model = genAI.getGenerativeModel({ model: modelName });
      const result = await model.generateContent(prompt);
      const output = result.response.text();
      const fixesPath = path.join(rootDir, "FIXES.md");
      fs.writeFileSync(fixesPath, output);

      console.log("--- DocCritic Audit Report ---");
      console.log(output);
      console.log("------------------------------");

      // Log results but do NOT fail CI (per user request: "shouldnt fail if qa has change remendations")
      const failedPath = path.join(rootDir, ".doc-failed");
      if (output.includes("BLOCKER")) {
        console.warn("DocCritic found BLOCKER issues. Please review FIXES.md and address them when possible.");
        fs.writeFileSync(failedPath, "true");
      } else {
        console.log("DocCritic audit passed (No Blockers).");
        if (fs.existsSync(failedPath)) fs.unlinkSync(failedPath);
      }

      console.log(`DocCritic complete using ${modelName}`);
      return;
    } catch (error) {
      console.warn(`Model ${modelName} failed: ${error.message}`);
    }
  }

  console.error("All Gemini models failed. Please check your API key and region permissions.");
  process.exit(1);
}

run();
