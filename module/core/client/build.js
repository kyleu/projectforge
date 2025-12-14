import esbuild from "esbuild";

const args = new Set(process.argv.slice(2));
const watch = args.has("--watch") || args.has("-w");
const dev = args.has("--dev");
const minify = !dev && !args.has("--no-minify");

const options = {
  entryPoints: ["src/client.ts"],
  bundle: true,
  minify: minify,
  sourcemap: true,
  outfile: "../assets/client.js",
  logLevel: "info",
  platform: "browser",
  target: ["es2021"],
  define: {
    "process.env.NODE_ENV": JSON.stringify(minify ? "production" : "development")
  }
};

if (watch) {
  {{{ .TypeScriptProjectWarning }}}const ctx = await esbuild.context(options);
  await ctx.watch();
  console.log("esbuild: watching for changes...");
} else {
  await esbuild.build(options);{{{ .TypeScriptProjectContent }}}
}
