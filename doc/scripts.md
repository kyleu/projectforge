<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# Scripts

There's a variety of shell scripts available in `./bin`. Here's a few of them:

- `bootstrap.sh`: Downloads and installs the Go libraries and tools needed in other scripts
- `build/android.sh`: Builds the Android library and application
- `build/build.sh`: Builds the app (or just use make build)
- `build/client.sh`: Uses `esbuild` to compile the scripts in `client`
- `build/client-watch.sh`: Builds the TypeScript resources, then watches for changes via `watchexec`
- `build/desktop.sh`: Uses `tools/desktop` to build a desktop application
- `build/desktop-release.sh`: Meant to be run as part of the release process, builds desktop apps
- `build/ios.sh`: Builds the iOS framework and application
- `build/release.sh`: Runs `goreleaser`
- `build/release-test.sh`: Runs `goreleaser` in test mode
- `check.sh`: Runs code statistics, checks for outdated dependencies, then runs linters
- `dev.sh`: Starts the app, reloading on changes
- `format.sh`: Formatting code from all projects
- `tag.sh`: Tags the git repo using the first argument or the incremented minor version
- `templates.sh`: Builds all the templates using quicktemplate, skipping if unchanged
- `test.sh`: Runs all the tests
- `util/view-binary-size.sh`: Visualizes space usage for the release binary
- `util/view-cpu-profile.sh`: Starts a `pprof` server using the (previously-recorded) CPU profile at `./tmp/cpu.pprof`
- `util/view-go-deps.sh`: Uses `gomod` to visualize the module graph
- `util/view-mem-profile.sh`: Starts a `pprof` server using the (previously-recorded) heap dump at `./tmp/mem.pprof`
- `workspace.sh`: Loads `dev.sh` and `client-watch.sh` in a split window. Requires iTerm
