@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Runs all the tests. Pass "-c" to clear the cache first, "-w" to watch for changes.

cd %~dpnx0\..

gotestsum -- -race app\...
