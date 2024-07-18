/* eslint-disable no-undef */
import esbuild from "esbuild";

esbuild.build({
  entryPoints: ["src/client.ts"],
  bundle: true,
  minify: true,
  sourcemap: true,
  outfile: "../assets/client.js",
  logLevel: "info"
}).catch((e) => console.error(e.message));
