# Project Forge TUI Framework Plan (Crush-Style UX)

## Goal
Build a reusable MVC-style TUI framework in Project Forge that:
- uses the Charm ecosystem (Bubble Tea, Bubbles, Lip Gloss, Huh, and selected adjacent libs),
- reproduces the Crush-style layout/interaction model,
- exposes an API analogous to `controller.Act` for TUI flows,
- renders views in raw Go (no templates).

## App Definition (What We Are Building)
The framework will power a concrete Project Forge terminal app with this startup menu:

1. `Projects`
- List active projects from `st.Services.Projects.Projects()`.
- Selecting a project opens a project workspace screen with action groups.

2. `Doctor`
- Expose doctor checks from `app/doctor/checks` and action workflows from `app/project/action/doctor.go`.
- Support run all, run single check, inspect errors, and solve when available.

3. `Settings`
- Configuration and admin tasks for the TUI and app runtime.
- Includes environment/config toggles, feature flags, and service diagnostics.

4. `About`
- Build/version metadata, environment summary, links, and quick docs.

Extensibility requirement:
- Startup menu is data-driven via a registry (not hardcoded switch blocks).
- Colors/theme should be handled by lipgloss, and built from a `Theme` (`./app/lib/theme`), defaulting to `./app/lib/theme/defaults.go`
- New top-level sections should be added by implementing a `Screen` and registering menu item(s).
- Reuse existing menu types from `app/lib/menu/item.go` and `app/lib/menu/items.go` (`menu.Item`, `menu.Items`).

### Projects Workspace (Detailed)
After choosing a project, show a project-focused workspace with:

- Primary actions:
  - `preview`, `generate`, `audit`, `build` (from `project/action.ProjectTypes`)
- Secondary actions:
  - `debug`, `rules`, `svg`, and optional hidden/internal actions as an advanced mode
- Git tools:
  - status, fetch/pull/push, commit, reset, history, magic flows (matching `app/lib/git/actions.go` and `app/controller/cproject/git.go`)
- Supporting panes:
  - project metadata, recent logs/output, status/errors, and action result summaries

Execution model:
- Actions run through the same backend action pipeline (`action.Apply` / `action.ApplyAll`) used by existing controllers.
- The TUI is an orchestration layer, not a second business-logic path.

## Research Summary (Primary Sources)
- Crush architecture and UI model:
  - https://github.com/charmbracelet/crush
  - `/tmp/crush_repo/internal/ui/model/ui.go` (state machine, update routing, layout generation, draw pipeline)
  - `/tmp/crush_repo/internal/ui/model/sidebar.go` (sidebar composition, dynamic truncation limits)
  - `/tmp/crush_repo/internal/ui/model/session.go` (files list with truncation and summary)
  - `/tmp/crush_repo/internal/ui/model/lsp.go` and `/tmp/crush_repo/internal/ui/model/mcp.go` (status sections)
  - `/tmp/crush_repo/internal/ui/AGENTS.md` (UI architecture conventions used by Crush itself)
- Bubble Tea (Elm-style model/update/view and command pattern):
  - https://github.com/charmbracelet/bubbletea
- Bubbles components:
  - https://github.com/charmbracelet/bubbles
- Lip Gloss styles and layout composition:
  - https://github.com/charmbracelet/lipgloss
- Huh forms, including embedded use as a `tea.Model`:
  - https://github.com/charmbracelet/huh
- Ultraviolet (cell-based renderer used by Crush):
  - https://github.com/charmbracelet/ultraviolet
- Wish (optional remote SSH serving for TUIs):
  - https://github.com/charmbracelet/wish
- Glamour (optional markdown rendering in panes):
  - https://github.com/charmbracelet/glamour

## Crush UX/Architecture Traits To Emulate
From source analysis, the key Crush behaviors worth copying are:

1. Single top-level orchestrator model
- One root UI model routes all messages, focus, and state transitions.
- Child components are mostly “dumb”: they expose methods and return `tea.Cmd` where needed.

2. Explicit UI state machine
- Distinct states (`onboarding`, `initialize`, `landing`, `chat`) and explicit focus states.
- State changes trigger layout recomputation and component resizing.

3. Deterministic layout function
- A single `generateLayout(width,height)` returns rectangles for all regions.
- Non-compact mode: main + editor on the left, fixed-width right sidebar.
- Compact mode: top header + main + editor, with details overlay.

4. Layered render pass
- Clear screen.
- Draw base panes (header/sidebar/main/editor).
- Draw secondary overlays (pills, completions).
- Draw dialogs last so they always win z-order.
- Final cursor placement computed from layout + focused component.

5. Event ingestion pattern
- Program input (`tea.Msg`) + async app/service events (pubsub) pushed into the same update loop.
- No expensive or blocking work in `Update`; work deferred to `tea.Cmd`.

6. Sidebar with adaptive truncation
- Sections for session metadata, files, LSPs, MCPs.
- Dynamic height budgeting across sections and explicit “...and N more” rows.
- ANSI-safe truncation and width calculations.

7. Strong key/help system
- Key maps grouped by context.
- Short help and full help generated from bindings.
- Focus-sensitive behavior (`tab` toggles editor/chat focus).

## Target Framework Design For Project Forge

### 1) High-Level Packages
Proposed under `app/controller/tui`:

- `framework/`
  - Runtime harness around `tea.Program`.
  - Message bus bridge for app/service events.
  - Panic recovery, lifecycle, shutdown orchestration.
- `mvc/`
  - Controller contracts (`Act` equivalent), route/state context, middleware hooks.
- `layout/`
  - Pane tree + rectangle solver.
  - Compact/non-compact breakpoints and adaptive sidebar sizing.
- `view/`
  - Raw-Go render helpers for pane composition.
  - Optional UV-backed renderer abstraction (default string view first).
- `components/`
  - Reusable editor, list, status, sidebar sections, dialogs, command palette.
- `style/`
  - Semantic tokens and style registry (Lip Gloss-based).
- `events/`
  - Typed app events and adapters from Project Forge services/logs.
- `screens/`
  - Concrete flows (`mainmenu`, `projects`, `project`, `doctor`, `settings`, `about`).

### 2) MVC Surface: TUI Version of `controller.Act`
Introduce an API with similar ergonomics to web `controller.Act`:

```go
type ActFn func(ts *TUIState, ps *PageState) (Transition, error)

func Act(key string, m *RootModel, f ActFn)
```

Recommended semantics:
- `key`: telemetry/logging/action key (`tui.home`, `tui.project.list`, etc).
- `TUIState`: shared app dependencies (`*app.State`, logger, settings, event services).
- `PageState`: per-screen ephemeral state (focus, selected ids, flash/status messages, pane metadata).
- `Transition`: next action (`Stay`, `Push(screen)`, `Replace(screen)`, `Pop`, `Quit`, `OpenDialog`, `Redirect(screenID)`).

Act-like guarantees to provide:
- panic-safe execution wrapper,
- timing/metrics capture,
- error interception to status/diagnostics pane,
- optional middleware chain before/after handler.

### 3) Root Model Contract

```go
type Screen interface {
    Key() string
    Init(*TUIState, *PageState) tea.Cmd
    Update(*TUIState, *PageState, tea.Msg) (Transition, tea.Cmd, error)
    View(*TUIState, *PageState, layout.Rects) string
    Help(*TUIState, *PageState) HelpBindings
}
```

Why this shape:
- keeps screen logic MVC-ish,
- keeps all state explicit and serializable for tests,
- lets root manage navigation stack + global overlays.

### 3.1) Menu and Screen Registry
Add a first-class registry for top-level navigation and future modules using existing menu models:
- `menu.Item` for nodes (`Key`, `Title`, `Description`, `Icon`, `Route`, `Hidden`, `Children`, etc.)
- `menu.Items` for collections and utility methods (`Get`, `Visible`, `GetByPath`, etc.)
- Each item's `Route` points to a screen key resolved by a `ScreenRegistry`.

Rules:
- `MainMenuScreen` renders `menu.Items.Visible()`.
- Each item points to a screen key (for example via `Key` or `Route`) resolved by a `ScreenRegistry`.
- Features can register additional items at bootstrap time.
- Menu rendering is generic list UI; business behavior lives in the target screen.

### 4) Layout Blueprint (Crush-Compatible)
Use a canonical pane layout that mirrors Crush:

- Non-compact (default desktop terminal):
  - `MainPane` (chat/content)
  - `EditorPane` bottom input strip
  - `SidebarPane` right, fixed baseline width (start with 30 cols)
  - `StatusPane` bottom line(s)
- Compact (small terminal):
  - `HeaderPane`
  - `MainPane`
  - `EditorPane`
  - optional `DetailsOverlay` on toggle
  - `StatusPane`

Layout solver requirements:
- one pure function from `(state, width, height)` to rectangles,
- min size constraints with graceful degrade,
- height budget allocator for sidebar sections,
- predictable z-order for overlays.

### 5) Components To Build First

1. `StatusBar`
- short/full help, transient notices, errors.

2. `Editor`
- Bubbles `textarea` wrapper with prompt styles and history hooks.

3. `VirtualList`
- viewport-like list for chat/log/event streams, with focus + selection.

4. `Sidebar`
- sections: context/project, modified files, LSP, MCP.
- ellipsis behavior and “...and N more”.

5. `DialogOverlay`
- modal stack; always rendered last.

6. `CommandPalette`
- slash/ctrl+p style launcher (possible integration with Huh selectors).

7. `ActionRunner`
- common component for running project/doctor actions and streaming status/output.

8. `MenuList`
- reusable selectable list for top-level menu and project/action selection screens.

### 5.1) App Screens To Build First

1. `MainMenuScreen`
- Startup landing menu with `Projects`, `Doctor`, `Settings`, `About`.

2. `ProjectsScreen`
- Lists projects from `st.Services.Projects.Projects()`.
- Includes search/filter/sort and quick metadata.

3. `ProjectScreen`
- Shows selected project details and action catalog.
- Runs actions (`preview`, `generate`, `audit`, `build`) and git operations.

4. `DoctorScreen`
- Shows checks relevant to current modules.
- Supports run all/run one/solve and readable result detail.

5. `SettingsScreen`
- Runtime config and admin operations.

6. `AboutScreen`
- version/build metadata and environment summary.

### 6) Charm Libraries: Adoption Matrix
Required baseline:
- `bubbletea`: main event loop and command runtime.
- `bubbles`: textarea/help/key/spinner/viewport/list primitives.
- `lipgloss`: styles and joins.

Recommended:
- `huh`: modal flows and structured input forms inside TUI (config/auth/wizards).
- `x/ansi` or ANSI-safe helpers: truncation, width-safe slicing.

Evaluate after baseline:
- `ultraviolet`: if we need cell-level rendering/perf parity closer to Crush.
- `glamour`: markdown-rich content panes.
- `wish`: remote/SSH mode for multi-user or headless server TUI access.

## Implementation Phases

### Phase 0: Scaffolding (in current tree)
- Create framework skeleton under `app/controller/tui/` (packages above).
- Keep existing `TUI.Run` entrypoint; wire it to new root runtime.
- Add `plan` ADR-style notes for architecture decisions.
- Add `MainMenuScreen` scaffold and screen/menu registries.

Acceptance:
- app boots into a visible main menu with 4 entries:
  - `Projects`
  - `Doctor`
  - `Settings`
  - `About`

### Phase 1: MVC Core
- Implement `Act` analog and screen stack transitions.
- Add panic/error/timing wrappers comparable to web controller flow.
- Add screen registry and typed route IDs.
- Build navigation flow: main menu -> each section -> back.

Acceptance:
- all top-level screens navigate push/pop/replace cleanly.

### Phase 2: Crush-Style Layout + Sidebar
- Implement deterministic layout solver with compact mode breakpoints.
- Implement right sidebar sections and truncation behavior.
- Add bottom editor strip + main viewport.

Acceptance:
- visual structure matches target screenshot class:
  - large main pane,
  - right context sidebar,
  - bottom input/status/help rows,
  - graceful compact fallback.

### Phase 3: Interaction Fidelity
- Keymaps by focus/context.
- Focus switching (editor/main/sidebar/dialog).
- Mouse wheel and selection support where feasible.
- Command palette and slash routing.
- Project action picker and execution UX parity with existing web/CLI actions.

Acceptance:
- full keyboard-only workflow across panes and overlays.

### Phase 4: App Integration
- Adapter layer from Project Forge services/logs/events to tea messages.
- Reuse `app.State` and existing logger in `TUIState`.
- Integrate:
  - projects list from `st.Services.Projects.Projects()`
  - project actions via `project/action`
  - doctor checks via `app/doctor/checks`
  - git actions via `app/lib/git`

Acceptance:
- fully functional first app slice:
  - main menu,
  - project selection,
  - project action execution,
  - doctor execution,
  - settings/about views.

### Phase 5: Hardening + Extensibility
- snapshot tests for layout math and key view strings.
- smoke tests for screen transitions and command routing.
- docs for adding new screens/components.

Acceptance:
- adding a new screen requires only registry entry + screen struct.

## Testing Strategy
- Unit:
  - layout rectangle math,
  - truncation/ellipsis logic,
  - transition reducer.
- Component:
  - sidebar rendering with constrained widths/heights,
  - editor + command palette interactions.
- Integration:
  - root model update loop with synthetic tea messages,
  - app event adapter to UI events.
- Golden/snapshot (recommended):
  - deterministic view output at fixed terminal sizes.

## Risks and Mitigations
- Risk: over-coupling to Crush internals.
  - Mitigation: copy behavioral patterns, not file-by-file architecture.
- Risk: premature complexity in custom renderer.
  - Mitigation: start with Bubble Tea string `View`, add UV path only if needed.
- Risk: async event races.
  - Mitigation: single-threaded state mutation in `Update`, all IO via `tea.Cmd`.
- Risk: style drift over time.
  - Mitigation: semantic token palette in `style/` and snapshot tests.

## Suggested Initial API Sketch

```go
// app/controller/tui/mvc/act.go
package mvc

type ActFn func(ts *State, ps *PageState) (Transition, error)

func Act(key string, rt Runtime, ps *PageState, fn ActFn) (Transition, tea.Cmd)
```

```go
// app/controller/tui/framework/root.go
package framework

type RootModel struct {
    state      *mvc.State
    nav        *mvc.Navigator
    overlays   *OverlayStack
    layout     layout.Rects
    keymap     KeyMap
}
```

```go
// app/controller/tui/screens/mainmenu.go
func (s *MainMenuScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error)
```

## Concrete Next Steps
1. Add package skeleton (`framework`, `mvc`, `layout`, `style`, `components`, `screens`).
2. Implement `mvc.Act` + `Transition` enum + `Navigator`.
3. Add `menu.Items`-backed menu registry and `ScreenRegistry`; implement `MainMenuScreen`.
4. Implement layout solver with non-compact + compact variants.
5. Implement `StatusBar`, `Sidebar`, `Editor`, `MenuList`, `ActionRunner` components.
6. Wire root model into `app/controller/tui/tui.go` `Run` method.
7. Ship first functional app slice:
  - `ProjectsScreen` from `st.Services.Projects.Projects()`
  - `ProjectScreen` with `preview/generate/audit/build` + git operations
  - `DoctorScreen` with run/solve flows
