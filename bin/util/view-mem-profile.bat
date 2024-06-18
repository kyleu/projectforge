@ECHO OFF

rem Starts a pprof server using the (previously-recorded) heap dump at ./tmp/mem.pprof

cd %~dp0\..\..

echo "=== launching profiler for mem.pprof ==="
@ECHO ON
go tool pprof -http=":8000" build\debug\projectforge tmp\mem.pprof
