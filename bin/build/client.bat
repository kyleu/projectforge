@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Uses `esbuild` to compile the scripts in `client`
rem Requires node, tsc, and esbuild available on the path

cd %~dp0\..\..\client

@ECHO ON
node build.js
