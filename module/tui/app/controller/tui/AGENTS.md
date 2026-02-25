# TUI Agent Notes

This directory contains the Bubble Tea terminal UI for {{{ .Name }}}.
Use this document when adding new screens, actions, components, or framework behavior.

## Scope

`/app/controller/tui` owns:

- TUI startup and integration with app state (`tui.go`, `tui_init.go`)
- Root navigation/event loop (`framework/`)
- Per-screen MVC contracts and state (`mvc/`, `screens/`)
- Terminal layout and style primitives (`layout/`, `style/`, `components/`)

The TUI is intentionally thin and orchestration-focused. Reuse existing app services (`ts.App.Services`) for business logic.

## Architecture

### Runtime Flow

1. `tui.NewTUI(...)` builds the TUI context and refreshes projects in `InitTUI`.
2. `runTUI(...)` in `app/cmd/tui.go` bootstraps the screen registry with `registry.Bootstrap()`.
3. `(*TUI).Run(...)` creates `mvc.State` and starts Bubble Tea with `framework.RootModel`.
4. `RootModel` is responsible for:
   - Terminal window sizing
   - Key-level global controls (`/` log drawer, `ctrl+c` quit)
   - Navigation stack (`push`, `pop`, `replace`, `quit`)
   - Rendering a common shell (main area, sidebar/header, editor hint, status/help)
5. Active screen handles screen-local `Init`, `Update`, and `View`.

### Navigation Model

- Navigation is stack-based (`framework.RootModel.stack`).
- Screen code returns `mvc.Transition` values (`Stay`, `Push`, `Pop`, `Replace`, `Quit`, `Route`).
- Each stack entry has its own `mvc.PageState` with context + telemetry span.
- `PageState.Close()` is called on pop/replace/quit to complete spans.

### Screen Contract

Screens implement `screens.Screen`:

- `Key() string`: stable route key
- `Init(*mvc.State, *mvc.PageState) tea.Cmd`: screen initialization
- `Update(..., tea.Msg) (mvc.Transition, tea.Cmd, error)`: message handling + optional async command
- `View(..., layout.Rects) string`: fully rendered body string
- `Help(...) screens.HelpBindings`: short/full keybinding help strings

Prefer deterministic screen behavior:

- Keep selection state in `PageState` (`Cursor`, `Data`).
- Set `ps.Status` for normal progress, `ps.SetError(err)` for failures.
- Return `tea.Cmd` for long-running work and post a typed result message back into `Update`.

## Directory Guide

- `tui.go`: process-level lifecycle, log tail capture hook
- `framework/root.go`: root model, global event handling, shell rendering, transition application
- `mvc/state.go`: immutable shared dependencies for screens
- `mvc/pagestate.go`: per-screen mutable UI state + telemetry span lifecycle
- `mvc/transition.go`: navigation command model
- `registry/bootstrap.go`: full registry wiring, including settings/admin screens
- `screens/*.go`: screen implementations
- `components/`: reusable view helpers (menu list, status bar)
- `layout/layout.go`: responsive terminal region solver
- `style/style.go`: theme-to-lipgloss style mapping

## Extending the TUI

### Add a New Screen

1. Add a key constant in `screens/keys.go`.
2. Implement a new file in `screens/` satisfying `Screen`.
3. Register core routes in `registry/bootstrap.go`:
   - Use `reg.Register(item, screen)` if it should appear in the main menu.
   - Use `reg.AddScreen(screen)` for hidden/detail routes.
4. Return navigation transitions from existing screens where appropriate.
5. Add/adjust tests for any layout/state helper logic.

### Preferred Screen Pattern

Use this shape for action-driven screens (see `project.go`, `doctor.go`):

- In `Init`: set `ps.Title`, set initial `ps.Status`, preload any data needed for `View`.
- In `Update`:
   - Handle movement keys (`up/down`, `j/k`)
   - Handle execution keys (`enter`, `r`, `a`, etc.)
   - Handle back keys (`esc`, `backspace`, `b`) with `mvc.Pop()`
   - Handle typed async result message and update `ps.Data` / status / error
- In `View`: derive UI only from `ts`, `ps`, and `rects`; avoid side effects.

### Async Work

For non-trivial tasks:

- Build a `tea.Cmd` closure that uses `ps.Context` and `ps.Logger`.
- Do service calls inside the command (not directly in key handler path).
- Return a typed message struct (`fooResultMsg`) containing `title`, `lines`, `err`.
- In `Update`, apply result to `ps.Data["result"]` and `ps.SetStatusText(...)`.

This keeps the UI responsive and consistent with Bubble Tea update semantics.

## Design Rules

- Keep TUI code as controller/orchestration logic, not domain logic.
- Do not bypass app services when functionality already exists in `app/lib` or `app/project` packages.
- Keep help bindings accurate with actual key handling.
- Always support keyboard-only navigation.
- Respect compact vs non-compact layouts via `layout.Solve(...)` and `rects.Compact`.
- Avoid rendering strings wider than available width; use truncation helpers where needed.
- Keep state local:
   - shared app dependencies in `mvc.State`
   - per-screen mutable values in `mvc.PageState`

## Styling and Layout

- Instantiate styles from `style.New(ts.Theme)` in each view.
- Render into the widths/heights provided by `layout.Rects`; do not hardcode terminal assumptions.
- In compact mode, stack panels vertically.
- In wide mode, use split-pane layouts with constrained sidebars.
- Reuse `components.RenderMenuList` and `components.RenderStatus` for consistency.

## Error Handling and Telemetry

- Use `mvc.Act` when wrapping complex transactional state updates to get panic protection + timing logs.
- `PageState` starts a telemetry span per page (`tui:<action>`); ensure pages are closed via transitions.
- Surface user-visible failures with `ps.SetError(err)` and keep `ps.Status` meaningful on success.

## Testing Guidance

Current tests cover layout solving (`layout/layout_test.go`).
When extending the TUI, prioritize tests for pure logic:

- layout calculations
- truncation/format helpers
- registry/menu wiring where practical

For interaction-heavy changes, do manual verification with `make build` and by running the TUI path in a terminal.

## Workflow Checklist

Before submitting TUI changes:

1. `./bin/format.sh`
2. `go test ./app/controller/tui/...`
3. `make build`
4. Manually verify keybindings and navigation stack behavior for changed screens

## Notes

- Keep `AGENTS.md` and code as the source of truth for current behavior.
- Keep this document updated when screen contracts, keybindings, root layout, or extension patterns change.
