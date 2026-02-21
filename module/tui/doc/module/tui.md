# Terminal UI/UX

The **`tui`** module adds a full-screen terminal interface built on Bubble Tea/Lip Gloss and wired into the same app services as the web UI.
It is launched from the CLI and includes an HTTP server status indicator, optional documentation browsing, and any other screens you add.

## Overview

- `tui` CLI command that starts the terminal UI and the HTTP server lifecycle (`app/cmd/tui.go`)
- Stack-based screen router and shared shell (`app/controller/tui/framework/root.go`)
- Typed screen contract with per-page state and telemetry spans (`app/controller/tui/screens/screen.go`, `app/controller/tui/mvc/pagestate.go`)
- Built-in screens for projects, docs, doctor, settings, and about (`app/controller/tui/screens/`)
- Log drawer support for recent runtime logs (up to 200 retained entries)

## CLI Usage

Run the TUI from your application binary:

```bash
./build/debug/{{{ .Exec }}} tui
```

You can use the standard root flags (for example):

```bash
./build/debug/{{{ .Exec }}} tui --addr 127.0.0.1 --port 9000 --working_dir .
```

Behavior notes from `app/cmd/tui.go`:

- Initializes quiet logging and streams log messages into the TUI log drawer.
- Initializes app state, loads server routes, and attempts to start the HTTP server.
- Passes server URL (or startup error) into the status bar.
- Runs the Bubble Tea program in the alternate screen.
- Shuts down HTTP server and app state cleanly on exit.

On WebAssembly builds, this command is unavailable (`app/cmd/tui_stub.go`).

## Screen Model

The TUI bootstraps screens in `app/controller/tui/screens/bootstrap.go` and starts at `mainmenu`.

Primary menu screens:

- `settings` - runtime diagnostics
- `about` - build and metadata details

Navigation uses typed transitions in `app/controller/tui/mvc/transition.go`:

- `Stay`, `Push`, `Pop`, `Replace`, `Route`, `Quit`

## Keybindings

Global bindings from `app/controller/tui/framework/root.go`:

- `/` toggles the log drawer
- `ctrl+c` quits immediately

Common screen bindings:

- `up/down` (and often `j/k`) move selection
- `enter` opens or executes
- `b`/`esc`/`backspace` goes back

Main menu also supports `q` to quit.

## Layout, Style, and Status

- Responsive layout solver in `app/controller/tui/layout/layout.go`
- Compact mode when terminal is smaller than `100x24`
- Theme-derived styles in `app/controller/tui/style/style.go`
- Shared components for menu lists and two-line status bar in `app/controller/tui/components/`
- Status bar displays action/help text and server URL or server startup error

## State and Telemetry

Shared app dependencies live in `app/controller/tui/mvc/state.go`.
Each pushed screen gets its own `PageState` with:

- scoped context and logger
- mutable UI fields (`Title`, `Cursor`, `Status`, `Error`, `Data`)
- telemetry span `tui:<action>` closed on pop/replace/quit

`mvc.Act` (`app/controller/tui/mvc/act.go`) provides optional panic-safe action wrappers with timing logs.

## Development Workflow

For rapid local iteration, use:

```bash
./bin/tui.sh
```

This helper rebuilds and reruns the TUI on file changes (requires `watchexec`).

## Configuration

This module does not introduce module-specific environment variables.
Runtime behavior is primarily controlled via existing CLI flags (`--addr`, `--port`, `--working_dir`, `--config_dir`, `--verbose`).

## Dependencies

- Requires a non-WASM target (native terminal runtime)
- Uses Charmbracelet libraries: Bubble Tea, Lip Gloss, Glamour
- Integrates with existing app services (projects, modules, doctor checks, git helpers)

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/tui
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Help Module](help.md)
- [charmbracelet](https://github.com/charmbracelet)
- [Project Forge Documentation](https://projectforge.dev)
