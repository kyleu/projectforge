# Core

This is a module for [Project Forge](https://projectforge.dev). It provides common utilities for a Go application.

https://github.com/kyleu/projectforge/tree/master/module/core

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Packages

See `customizing.md` for details

- `cmd` contains the main CLI actions
- `controller` contains HTTP actions for the server UI, see [faq.md](../faq.md) for details
- `lib/filter` is used by the UI for sorting and filtering
- `lib/log` contains custom zap loggers and appenders
- `lib/menu` is used by the UI to draw the left nav and breadcrumbs
- `lib/telemetry` allows tracing via OpenTelemetry and metrics via Prometheus, used everywhere
- `lib/theme` contains UI themes for controlling the UI look and feel
- `lib/user` defines user, accounts, and permissions
- `util` contains dozens of useful helper functions
