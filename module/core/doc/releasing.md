# Releasing

{{{ .Name }}} uses `goreleaser` to create build artifacts.

You can release to GitHub using `./bin/build/release.sh`, or test the release by running `./bin/build/release-test.sh`.

Your releases are available at {{{ .Info.Sourcecode }}}/releases

### Checksums

All release binaries are checksummed, available in `checksums.txt` in the root of the release
{{{ if .Build.Changelog }}}
### Changelog

A changelog will be created based on the commit history, including all authors and messages
{{{ end }}}
### Docker Images

Multiple Docker images will be created. The main image is `ghcr.io/{{{ .Info.Org }}}/{{{ .Key }}}/x.x.x`, and a debug image is provided at `ghcr.io/{{{ .Info.Org }}}/{{{ .Key }}}/x.x.x-debug` that includes `delve` for debugging
{{{ if .Build.Homebrew }}}
### Homebrew

Packages for macOS and Linux will be pushed to Homebrew at `{{{ .Info.Org }}}/homebrew-{{{ .Info.Org }}}`
{{{ end }}}{{{ if .Build.NFPMS }}}
### NFPMS

The build will produce `apk`, `deb`, and `rpm` packages for each supported Linux architecture
{{{ end }}}{{{ if .Build.BOM }}}
### BOM

The build will create a bill of materials for each binary
{{{ end }}}{{{ if .BuildNotarize }}}
### Notarization

Release binaries for macOS and iOS are notarized using Apple Notarization services
{{{ end }}}{{{ if .Build.Signing }}}
### Signing

Release binaries and the checksum file are signed using `gpg`
{{{ end }}}{{{ if .Build.Snapcraft }}}
### Snapcraft

The build will produce `snap` packages for each supported Linux architecture
{{{ end }}}
### Source Code

The source code will be bundled in the release, available as `{{{ .Key }}}_x.x.x_source.zip`

### Universal Binaries

A universal macOS app will be created, containing the complete app for x86-64 and ARM
{{{ if .BuildDesktop }}}
### Desktop Build

A standalone desktop application, bundling the server and a web view, will be built for Linux, macOS, and Windows
{{{ end }}}{{{ if .BuildMobile }}}
### Mobile Build

A standalone mobile app, bundling the server and a web view, will be built for Android and iOS
{{{ end }}}{{{ if .BuildWASM }}}
### Web Build

The server is compiled to WASM and available as a WebAssembly module{{{ end }}}
