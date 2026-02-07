# Rich Terminal UI/UX Plan

Goal: build a polished, interactive terminal UI using charmbracelet libraries. Entry point is `RunTUI` in `app/controller/tui/tui.go` (already wired). This plan is implementation-focused and should be followed in a later session.

## Scope and initial experience

- Show a splash screen with branding and a short prompt ("press any key" or "enter to continue").
- Transition to a project picker (data from `st.Services.Project.List()`).
- When a project is selected, print its name and exit (initial milestone).
- Ensure graceful exit on Ctrl+C and Esc.

## Suggested libraries

- `bubbletea` for the core event loop and model.
- `bubbles` for list, textinput, help, spinner.
- `lipgloss` for styling and layout.
- (Optional) `glamour` for rich text help screens and markdown rendering later.

## Architecture and model design

- Use a small, explicit state machine: `splash -> loading -> projects -> done/exit`.
- Central `model` struct should include:
  - `state` enum
  - `width`, `height` for layout and reflow
  - `projects` slice and a `list.Model`
  - `err` for display + exit
  - `ready` boolean to avoid rendering before `WindowSizeMsg`
  - `styles` (lipgloss styles grouped in a struct)
  - `keys` (keymap struct for help)
- Separate rendering for each state with a `viewSplash`, `viewLoading`, `viewList`, `viewError`.
- Keep business logic (loading projects) separate from view logic where possible.

## Data loading and async flow

- Load projects via a `tea.Cmd` that calls `st.Services.Project.List()` in a goroutine.
- Use a `loading` state with a spinner until data arrives.
- Handle error: show a friendly error screen with instructions to quit.

## Interaction and navigation

- Splash: any key advances, Ctrl+C exits.
- Project list:
  - Up/Down or j/k to navigate.
  - Enter selects.
  - / to start a filter (if using `bubbles/list` filter).
  - Esc to go back to splash (optional) or quit.
- Help footer showing keybinds (use `bubbles/help`).

## Styling and layout

- Define a style palette once (primary, muted, danger, highlight).
- Use `lipgloss.Place` to center splash and loading screens.
- For list view, create a header (title + subtitle), body (list), footer (help).
- Respect terminal size: use `WindowSizeMsg` and recompute widths.
- Consider a minimal theme that works in both light/dark terminals (avoid hard-coded colors).

## Terminal UX considerations

- Avoid flicker: update list size on resize only when `ready`.
- Show a spinner during loading.
- Provide empty-state messaging when no projects exist.
- Provide clear exit hint ("q to quit").

## Error handling

- Wrap list retrieval errors with context.
- Render an error view with a short message + details (truncated to fit).
- Exit with non-zero status when errors are fatal.

## File structure suggestions (if needed)

- `app/controller/tui/tui.go`: `RunTUI`, initial model setup, `tea.NewProgram`.
- `app/controller/tui/model.go`: model struct + update routing.
- `app/controller/tui/view.go`: view helpers + styles.
- `app/controller/tui/keys.go`: keymap and help view.
- `app/controller/tui/data.go`: commands for loading projects.

## Incremental implementation plan

1. Create model + state machine with splash and basic update loop.
2. Add async load command and loading screen.
3. Add list model and render list view.
4. Wire keybindings and selection behavior.
5. Add styles and footer help.
6. Handle errors and empty states.
7. Add polish (transitions, consistent padding, exit codes).

## Suggested improvements (post-milestone)

- Project details panel on the right (metadata, description).
- Recent projects section or pinned favorites.
- Search across projects and filter by tags.
- Multi-step flow: select project -> actions (open, run, build).
- Persist last selection (user config).
- Themed icons via unicode (optional, keep ASCII fallback).
