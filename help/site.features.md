# Features

Project Forge is modular. You enable the features you need via modules in `.projectforge/project.json`.

Common modules include:

- `core`: base runtime, controllers, templates, and UI
- `database` + `postgres`/`mysql`/`sqlite`: database access and migrations
- `oauth` + `user`: authentication and user management
- `search`: global search UI and providers
- `export`: schema-driven code generation (Go + TypeScript)
- `websocket`: real-time UI and server events

Use the modules page in the UI to browse available modules and see what each one adds.
