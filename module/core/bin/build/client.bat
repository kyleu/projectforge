@ECHO OFF

rem Uses `esbuild` to compile the scripts in `client`
rem Requires node, tsc, and esbuild available on the path

cd %~dp0\..\..\client

@ECHO ON
node build.js
