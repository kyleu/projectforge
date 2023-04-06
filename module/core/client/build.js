/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
const esbuild = require("esbuild");

esbuild.build({
  entryPoints: ["src/client.ts"],
  bundle: true,
  minify: true,
  sourcemap: true,
  outfile: "../assets/client.js",
  logLevel: "info"
}).catch((e) => console.error(e.message));
