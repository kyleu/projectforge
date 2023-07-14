@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Starts a pprof server using the (previously-recorded) CPU profile at ./tmp/cpu.pprof

cd %~dp0\..\..

echo "=== launching profiler for cpu.pprof ==="
@ECHO ON
go tool pprof -http=":8000" build\debug\projectforge tmp\cpu.pprof

