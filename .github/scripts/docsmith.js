import fs from "fs";
import { GoogleGenerativeAI } from "@google/generative-ai";

const genAI = new GoogleGenerativeAI(process.env.GOOGLE_AI_KEY);
const model = genAI.getGenerativeModel({ model: "gemini-1.5-pro" });

const generatedReadme = fs.existsSync("README.generated.md")
  ? fs.readFileSync("README.generated.md", "utf-8")
  : "";

if (!generatedReadme) {
  console.error("README.generated.md not found.");
  process.exit(1);
}

const prompt = `
You are DocSmith, the M3TAL Ecosystem Aesthetic Architect.

Task:
Rewrite the following RAW generated README to be professional, high-density, and aligned with the "Mission Control" aesthetic.

Rules:
- DO NOT invent features or services.
- PRESERVE all technical details from the raw README (ports, variable names, service paths).
- STRUCTURE the information using modern markdown (callouts, clean tables, dividers).
- IMPROVE the language to sound authoritative and professional.
- ADD a standard M3TAL header and footer.

RAW GENERATED README:
${generatedReadme}
`;

async function run() {
  try {
    const result = await model.generateContent(prompt);
    const output = result.response.text();
    fs.writeFileSync("README.md", output);
    console.log("README polished by DocSmith");
  } catch (error) {
    console.error("DocSmith failed:", error);
    process.exit(1);
  }
}

run();
