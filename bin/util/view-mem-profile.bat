@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Starts a pprof server using the (previously-recorded) heap dump at ./tmp/mem.pprof

cd %~dpnx0\..\..

echo "=== launching profiler for mem.pprof ==="
go tool pprof -http=":8000" build\debug\projectforge tmp\mem.pprof
