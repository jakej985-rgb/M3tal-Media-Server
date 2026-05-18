import fs from "fs";

fs.mkdirSync("docs", { recursive: true });

// Scan the actual Go source layout: cmd/ and internal/ packages
const modules = [];

const scanDir = (dir, type) => {
  if (!fs.existsSync(dir)) return;
  fs.readdirSync(dir, { withFileTypes: true })
    .filter(d => d.isDirectory() && !d.name.startsWith("."))
    .forEach(d => {
      modules.push({
        name: d.name,
        type,
        path: `${dir}/${d.name}`
      });
    });
};

// Core Go binary entrypoints
scanDir("cmd", "CLI binary");

// Internal Go packages
scanDir("internal", "Go package");

// Deploy artifacts
if (fs.existsSync("deploy")) {
  fs.readdirSync("deploy", { withFileTypes: true })
    .filter(d => d.isDirectory() && !d.name.startsWith("."))
    .forEach(d => {
      modules.push({
        name: d.name,
        type: "Deploy artifact",
        path: `deploy/${d.name}`
      });
    });
}

// Top-level markers
if (fs.existsSync("go.mod")) {
  const modContent = fs.readFileSync("go.mod", "utf-8");
  const modLine = modContent.split("\n").find(l => l.startsWith("module "));
  modules.push({
    name: modLine ? modLine.replace("module ", "").trim() : "go-module",
    type: "Go module",
    path: "go.mod"
  });
}

fs.writeFileSync("docs/m3tal-services.json", JSON.stringify(modules, null, 2));
console.log(`Discovered ${modules.length} M3TAL source modules.`);
