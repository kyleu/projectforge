# AGENTS.md: Working with Project Forge Modules

This guide is for agents (and humans) editing files under `module/`.

The short version: files in `module/<name>/` are **templates**, not runtime app files. Project Forge renders them with a `*template.Context`, compares the output to a target project, then applies diffs.

## Scope and intent

Use this doc when you need to:
- Understand how module templates become generated project files
- Predict why a file is included, skipped, or preserved
- Prepare for future "weave generated changes back into templates" work

For module catalog/overview docs, start with `module/README.md`.

## Module anatomy

Each module directory usually contains:
- `.module.json`: module metadata (name, description, requires, config vars, port offsets, etc.)
- `README.md`: short module readme
- `doc/module/<key>.md`: module usage docs
- template trees such as `app/`, `views/`, `client/`, `assets/`, `tools/`, `bin/`

Loaded metadata is represented by `app/module/module.go` (`type Module`).

Important `.module.json` fields:
- `requires`: dependency modules that must also be enabled
- `priority`: module ordering precedence when files collide (later load wins by path key)
- `configVars`: module config var declarations merged into template context
- `portOffsets`: named service offsets merged into template context
- `dangerous`: flagged by project validation when safe mode is enabled

Dependency validation is enforced in `app/project/validate.go`.

## Template render pipeline

High-level flow (preview/generate both use this):
1. Collect files from all enabled modules (`app/module/files.go`)
2. Exclude special files like `.module.json`, module `README.md`, binary icon assets
3. Add generated extras (`README.md`, export outputs, VS Code settings/launch)
4. Build template context (`app/project/template/template.go`)
5. Render templated filenames first (`{{{ ... }}}` in paths)
6. Preserve marked sections from the current target file (`file.ReplaceSections`)
7. Render file content with Go templates and custom delimiters
8. Diff rendered source vs target project (`app/file/diff`)
9. Apply generation action based on diff status (`app/project/action/generate.go`)

Main entrypoints:
- `app/project/action/diffs.go`
- `app/project/action/util.go`
- `app/project/action/generate.go`
- `app/file/diff/load.go`

## Template syntax and context

Template engine:
- Go `text/template`
- Delimiters are `{{{` and `}}}` (not `{{` and `}}`)

Context type:
- `app/project/template/template.go` (`type Context`)

Core context fields:
- Project identity: `Key`, `Name`, `Exec`, `Version`, `Package`, `Args`, `Port`
- Module state: `Modules`, `Tags`
- Config and metadata: `ConfigVars`, `PortOffsets`, `Config`, `Info`, `Build`, `Theme`
- Derived: `DatabaseEngine`, `Ignore`, `IgnoreGrep`, `Linebreak`

Helper methods are on `*Context` across:
- `app/project/template/templatehelper.go`
- `app/project/template/templatecontent.go`
- `app/project/template/templatetypescript.go`
- `app/project/template/templateservices.go`

Common usage:
```gotemplate
{{{ .Package }}}
{{{ if .HasModule "search" }}} ... {{{ end }}}
{{{ if .HasModules "database" "migration" }}} ... {{{ end }}}
```

## Filename templating

Paths can be templated too. If a source path contains `{{{`, Project Forge renders the path before content rendering.

This is done in `runTemplates(...)` in `app/project/action/diffs.go`.

## Special `$PF_*` control tags

These tags are parsed by `app/file/sections.go` and `app/file/diff/load.go`.

### `$PF_SECTION_START(name)$` / `$PF_SECTION_END(name)$`

Purpose:
- Preserve user-edited regions from the target file during regeneration.

Behavior:
- If the incoming template contains section tags and target file exists, section content is copied from target into rendered content before diffing.
- Section keys must match in source and target.
- Malformed or nested section markers can fail generation.

### `$PF_HAS_MODULE(modA,modB,...)$`

Purpose:
- File-level module gate.

Behavior:
- Must be on the first meaningful line.
- If project modules do not include all listed modules, the file is skipped.
- If included, this marker line is removed from content before diffing/writing.

### `<!--- $PF_GENERATE_ONCE$ --->` (or equivalent inline marker)

Purpose:
- Protect files from subsequent generate updates after initial creation.

Behavior:
- During diff load, files containing this marker are treated as skipped.
- During write-to-project generation, the marker line is stripped before writing output.

### `$PF_INJECT_*`

Purpose:
- Indexed content injection hooks (`file.Inject`).

Status:
- Supported by file utilities; less common in module templates than section/module/generate-once markers.

## Diff statuses and generation actions

Diff status model (`app/file/diff/status.go`):
- `different`: both files exist, content differs
- `new`: exists in source render, missing in target
- `missing`: exists in target, not in source render
- `identical`
- `skipped`

Generate action (`app/project/action/generate.go`) applies only actionable diffs:
- `to=project` (default): writes rendered source into target file
- `to=ignore`: adds path to `project.Info.IgnoredFiles`
- `to=rebuild`: placeholder, currently not implemented

## Ignore behavior (important)

There are two different ignore mechanisms:
- `Project.Ignore`: path patterns used in various project checks and helper output
- `Project.Info.IgnoredFiles`: concrete file paths skipped by generate diff/apply

For generation decisions, `IgnoredFiles` is the key list passed into diffing.

## File collisions across modules

When multiple modules provide the same output path:
- Later-loaded module file wins for that path key in `Service.GetFiles(...)`
- Module ordering is by module `priority` (`Modules.Sort()`)

Practical implication:
- Overriding a file via another module depends on module priority/load order.

## Practical guidance for agent edits

When modifying module templates:
1. Identify whether the target file is plain output, section-preserved, module-gated, or generate-once.
2. Check whether filename templating affects the output path.
3. Check context usage (`.HasModule`, `.Build*`, `.Database*`, `.Config`, `.Info`).
4. Validate module dependencies in `.module.json` when introducing new cross-module references.
5. Run preview/generate on a representative project and inspect diff statuses.

Suggested verification loop:
- `./bin/dev.sh` (or `make dev`) for rapid iteration
- `./bin/test.sh` and `./bin/check.sh` before finalizing

## Notes for future "weave back" implementation

Current implementation is source-template -> generated-project only.

For reverse weaving, this doc captures the minimum facts needed to build robust mapping logic:
- Template path rendering rules
- File-level inclusion gates (`$PF_HAS_MODULE`)
- Protected content boundaries (`$PF_SECTION_*`)
- One-time ownership markers (`$PF_GENERATE_ONCE$`)
- Module precedence on path collisions
- Context-derived conditional branches

If you are building that feature, start at:
- `app/project/action/diffs.go`
- `app/file/diff/load.go`
- `app/file/sections.go`
- `app/project/template/*.go`
