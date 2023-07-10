@ECHO OFF

rem Starts a pprof server using the (previously-recorded) heap dump at ./tmp/mem.pprof

cd %~dpnx0\..\..

echo "=== launching profiler for mem.pprof ==="
go tool pprof -http=":8000" build\debug\{{{ .Exec }}} tmp\mem.pprof
