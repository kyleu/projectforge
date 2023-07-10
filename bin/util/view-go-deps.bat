@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Uses gomod to visualize the module graph
rem Requires gomod available on the path

cd %~dpnx0\..\..

@ECHO ON
echo "building dependency SVG..."
gomod graph | dot -Tsvg -o tmp\deps.svg

open tmp\deps.svg
