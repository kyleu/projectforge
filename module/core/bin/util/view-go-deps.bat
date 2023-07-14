@ECHO OFF

rem Uses gomod to visualize the module graph
rem Requires gomod available on the path

cd %~dp0\..\..

echo "building dependency SVG..."
@ECHO ON
gomod graph | dot -Tsvg -o tmp\deps.svg
