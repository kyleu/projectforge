# Project Forge

Project Forge is a Go-first application generator and project manager. It creates a full Go codebase (no runtime dependency on Project Forge), then lets you evolve it by regenerating modules, targets, and schema-driven code as your project changes.

## Quick Start

1. Create a project:
   - `projectforge create`
2. Start the dev server:
   - `./bin/dev.sh` (or `make dev`)
3. Open the UI:
   - `http://localhost:40000`
4. Change features:
   - edit `.projectforge/project.json`, then run `projectforge generate`

## Useful Pages

- Install Project Forge (this help section)
- Features and modules
- Doctor (dependency checks and diagnostics)
