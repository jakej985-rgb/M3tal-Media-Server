import fs from "fs";

if (!fs.existsSync("docs")) {
  fs.mkdirSync("docs");
}

if (!fs.existsSync("source")) {
  console.log("No source directory found.");
  process.exit(0);
}

const services = fs.readdirSync("source", { withFileTypes: true })
  .filter(d => d.isDirectory() && !d.name.startsWith('.'))
  .map(d => {
    const hasGo = fs.existsSync(`source/${d.name}/go.mod`);
    const hasPy = fs.existsSync(`source/${d.name}/requirements.txt`);
    return {
      name: d.name,
      type: hasGo ? "Go" : (hasPy ? "Python" : "Standard"),
      path: `source/${d.name}`
    };
  });

fs.writeFileSync("docs/m3tal-services.json", JSON.stringify(services, null, 2));
console.log(`Discovered ${services.length} M3TAL source services.`);
