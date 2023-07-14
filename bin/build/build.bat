@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Builds the app (or just use make build)

cd %~dp0\..\..

os=${1:-darwin}
arch=${2:-amd64}
fn=${3:-projectforge}

echo "Building [$os $arch]..."
@ECHO ON
env GOOS=$os GOARCH=$arch make build-release
md build\$os\$arch
move "build\release\$fn" "build\$os\$arch\$fn"
