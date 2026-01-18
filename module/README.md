# Project Forge Modules

This directory contains the template modules used by Project Forge to generate applications. Each module is a self-contained feature bundle that contributes files to a generated project when enabled.

## How modules work

- Each module is a template tree that is merged into the target project.
- Files are rendered with Go templates using custom delimiters: `{{{` and `}}}`.
- The template context is the Project Forge model (project.json plus computed metadata).

## Documentation

Module docs index: [`doc/module/`](../doc/module/)

## Module layout

```
module/<name>/
  .module.json          # Metadata (name, description, config vars, dependencies)
  README.md             # Module-specific overview
  doc/module/<name>.md  # Full module documentation
  app/                  # Go source templates
  assets/               # Static asset templates
  client/               # TypeScript/JS templates
  views/                # quicktemplate HTML templates
  tools/                # Build/tooling templates
```

## Editing notes

- Files under `module/` are templates, not the running Project Forge app.
- Template directives can appear in any file type and are processed before files are written to a generated project.
- See the module's `doc/module/<name>.md` for module-specific behavior and configuration.
- After template changes, regenerate a target project and run its build/tests to validate the output.
