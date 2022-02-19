<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# [core]

This is a module for [Project Forge](https://projectforge.dev). It provides common utilities for a Go application.

https://github.com/kyleu/projectforge/tree/master/module/core

### License

Licensed under [CC0](https://creativecommons.org/share-your-work/public-domain/cc0)

### Packages
- `cmd` contains the main CLI actions; start at `run.go`
- `controller` contains HTTP actions for the server UI
- `lib/filesystem` contains an interface for manipulating files
- `lib/filter` is used by the UI for sorting and filtering
- `lib/log` contains custom zap loggers and appenders
- `lib/menu` is used by the UI to draw the leftnav and breadcrumbs
- `lib/telemetry` allows tracing and metrics, used everywhere
- `lib/theme` contains UI themes for controlling the UI look and feel
- `lib/user` defines user, accounts, and permissions
- `util` contains dozens of useful helper functions
