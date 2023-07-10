@ECHO OFF

rem Builds the app (or just use make build)

cd %~dpnx0\..\..

os=${1:-darwin}
arch=${2:-amd64}
fn=${3:-{{{ .Exec }}}}

echo "Building [$os $arch]..."
env GOOS=$os GOARCH=$arch make build-release
md build\$os\$arch
move "build\release\$fn" "build\$os\$arch\$fn"
