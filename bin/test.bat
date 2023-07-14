@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Runs all the tests. Pass "-c" to clear the cache first, "-w" to watch for changes.

cd %~dp0\..

@ECHO ON
gotestsum -- -race app\...
