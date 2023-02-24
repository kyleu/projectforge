// Content managed by Project Forge, see [projectforge.md] for details.
/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
const esbuild = require("esbuild");

esbuild.build({
  entryPoints: ["src/client.ts"],
  bundle: true,
  minify: true,
  sourcemap: true,
  outfile: "../assets/client.js",
  watch: process.argv[2] === "watch" ? {
    onRebuild(error, result) {
      if (error) {
        console.error("watch build failed:", error);
      } else {
        console.log("watch build succeeded:", result);
      }
    }
  } : null,
  logLevel: "info"
}).catch((e) => console.error(e.message));
