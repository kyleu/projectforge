@ECHO OFF

rem Starts a pprof server using the (previously-recorded) CPU profile at ./tmp/cpu.pprof

cd %~dpnx0\..\..

echo "=== launching profiler for cpu.pprof ==="
go tool pprof -http=":8000" build\debug\{{{ .Exec }}} tmp\cpu.pprof

