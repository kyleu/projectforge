# Project Forge

![app logo](./assets/favicon.png)

[Project Forge](https://projectforge.dev) is an application that allows you to generate, manage, and grow web applications built using the Go programming
language.
You control your application's features, provided via "modules" that enable everything from databases to OAuth.
When creating a new application with Project Forge, a standard Golang project is created, which includes utilities and APIs to help you build the application
you want, without compromise.

All projects managed by Project Forge provide an HTTP server using [quicktemplate](https://github.com/valyala/quicktemplate) for HTML templates (and SQL, if
enabled).
An MVC framework is provided (but not required) that handles content negotiation, hierarchical menus, breadcrumbs, OAuth to dozens of providers, stateless user
profiles, dark mode support, SVG management, syntax highlighting, form components, and embedded assets.

Project Forge applications can support any UI framework, but the included UI renders a JS-dependency-free page, heavily optimized for speed and accessibility,
with a modern UX that works surprisingly well without JavaScript.
The about page is animated, themed, and responsive, and only creates three requests (HTML, CSS, JS) totaling less than 20KB zipped.
It serves in less than a millisecond, and renders in Chrome in less than 20ms.
Progressive enhancement is provided by an included ESBuild TypeScript project, though all functionality is supported with JavaScript disabled.

Your application can (optionally) build for _every_ platform; desktop and mobile webview apps, WASM, and notarized universal macOS binaries.
If you enable all the build options, it will produce almost 60 builds for various platforms. They all produce a ~20MB native binary.
The binaries produced can be configured to auto-upgrade from GitHub Releases, or be upgraded by the user using a CLI or UI (module "upgrade" must be in your
project).
CI/CD workflows based on GitHub Actions are provided, handling building, testing, linting, and publishing to GitHub Releases (and any configured Docker repos).

## Download

https://projectforge.dev/download

## Source code

https://github.com/kyleu/projectforge

## App Features

A project managed by Project Forge...

- Has a beautiful and fast HTML UI, no JavaScript required
- Optimizes for speed and developer experience, sub-second turnaround times
- Builds in seconds, live-reloads in dev mode when code or templates change
- Builds native apps for dozens of platforms; mobile, desktop, weird architectures, macOS universal apps
- Produces a small self-contained binary, includes full REST server and command line interface
- Optionally supports OAuth with dozens of providers, full theme and stateless session support
- Provides a type-safe API for working with PostgreSQL, MySQL, and SQLite databases (or no database at all)
- Uses ESBuild for a TypeScript client with progressive enhancement, embedded SVGs, and modern CSS

## Available Modules

|                                                                  |                                                                                                |
|------------------------------------------------------------------|------------------------------------------------------------------------------------------------|
| [Core](module/core/doc/module/core.md)                           | Provides common utilities for a Go application                                                 |
| [Android App](module/android/doc/module/android.md)              | Webview-based application and Android build                                                    |
| [Audit](module/audit/doc/module/audit.md)                        | Using the database module, provides an audit framework for tracking changes                    |
| [Brand Icons](module/brands/doc/module/brands.md)                | Provides thousands of SVG icons from [simple-icons] for representing common logos              |
| [Database](module/database/doc/module/database.md)               | Provides an API for accessing relational databases                                             |
| [Database UI](module/databaseui/doc/module/databaseui.md)        | Provides a UI for registered databases                                                         |
| [Desktop](module/desktop/doc/module/desktop.md)                  | Provides a desktop application using the system's webview                                      |
| [Doc Browser](module/docbrowse/doc/module/docbrowse.md)          | Provides a UI for browsing the documentation                                                   |
| [Export](module/export/doc/module/export.md)                     | Generates code based on the project's schema                                                   |
| [Expression](module/expression/doc/module/expression.md)         | Exposes CEL engine for evaluating arbitrary expressions                                        |
| [Filesystem](module/filesystem/doc/module/filesystem.md)         | Provides an abstraction around local and remote filesystems                                    |
| [Git](module/graphql/doc/module/git.md)                          | Helper classes for performing operations on git repositories                                   |
| [GraphQL](module/graphql/doc/module/graphql.md)                  | Supports GraphQL APIs within your application                                                  |
| [Help](module/help/doc/module/help.md)                           | Provides Markdown help files that integrate into the UI                                        |
| [HTTP Archive](module/har/doc/module/har.md)                     | Provides classes for parsing HTTP Archive (*.har) files                                        |
| [iOS App](module/ios/doc/module/ios.md)                          | Webview-based application and iOS build                                                        |
| [JSX](module/jsx/doc/module/jsx.md)                              | Provides a slim JSX implementation for scripting                                               |
| [Marketing Site](module/marketing/doc/module/marketing.md)       | Provides a website for downloads, tutorials, and marketing                                     |
| [Migration](module/migration/doc/module/migration.md)            | Database migrations and a common database pool                                                 |
| [MySQL](module/mysql/doc/module/mysql.md)                        | Provides an API for accessing MySQL databases                                                  |
| [Notarize](module/notarize/doc/module/notarize.md)               | Sends files to Apple for notarization                                                          |
| [Notebook](module/notebook/doc/module/notebook.md)               | Provides an Observable Framework notebook                                                      |
| [Numeric](module/numeric/doc/module/numeric.md)                  | It provides TypeScript and Golang implementations for managing large numbers.                  |
| [OAuth](module/oauth/doc/module/oauth.md)                        | Provides logins and session management for many OAuth providers                                |
| [OpenAPI](module/openapi/doc/module/openapi.md)                  | Embeds the Swagger UI, using your OpenAPI specification                                        |
| [Playwright](module/playwright/doc/module/playwright.md)         | Adds a project for testing the UI using playwright.dev                                         |
| [Plot](module/plot/doc/module/plot.md)                           | Library for visualizing data using Observable Plot                                             |
| [PostgreSQL](module/postgres/doc/module/postgres.md)             | Provides an API for accessing PostgreSQL databases                                             |
| [Process](module/process/doc/module/process.md)                  | Provides a framework for managing system processes                                             |
| [Proxy](module/proxy/doc/module/proxy.md)                        | Provides an HTTP proxy while still enforcing this app's security                               |
| [Queue](module/queue/doc/module/queue.md)                        | Provides a simple message queue based on SQLite                                                |
| [Read-only DB](module/readonlydb/doc/module/readonlydb.md)       | Adds a read-only database connection                                                           |
| [Rich Editor](module/richedit/doc/module/richedit.md)            | It provides a rich editing experience with a decent fallback when scripting is disabled        |
| [Sandbox](module/sandbox/doc/module/sandbox.md)                  | Useful playgrounds for testing custom functions                                                |
| [Scheduled Jobs](module/schedule/doc/module/schedule.md)         | Provides a scheduled job engine and UI based on gocron                                         |
| [Scripting](module/scripting/doc/module/scripting.md)            | Allows the execution of JavaScript files using a built-in interpreter                          |
| [Search](module/search/doc/module/search.md)                     | Adds search facilities to the top-level navigation bar                                         |
| [SQL Server](module/sqlserver/doc/module/sqlserver.md)           | Provides an API for accessing MSSQL databases                                                  |
| [SQLite](module/sqlite/doc/module/sqlite.md)                     | Provides an API for accessing SQLite databases                                                 |
| [Task](module/themecatalog/doc/module/task.md)                   | Provides an engine for executing and monitoring tasks                                          |
| [Theme Catalog](module/themecatalog/doc/module/themecatalog.md)  | Includes a dozen default themes, and facilities to create additional                           |
| [Types](module/types/doc/module/types.md)                        | Classes for representing common data types                                                     |
| [Upgrade](module/upgrade/doc/module/upgrade.md)                  | Provides in-place version upgrades using GitHub Releases                                       |
| [User](module/user/doc/module/user.md)                           | Classes for representing a user                                                                |
| [WebAssembly Client](module/wasmclient/doc/module/wasmclient.md) | Provides a WASM library and HTML host for an HTTP client                                       |
| [WebAssembly Server](module/wasmserver/doc/module/wasmserver.md) | Build your normal app as an http server, but load it as a WebAssembly module or Service Worker |
| [WebSocket](module/websocket/doc/module/websocket.md)            | Provides an API for hosting WebSocket connections                                              |

## Example Applications

- [Rituals.dev](https://rituals.dev) ([GitHub](https://github.com/kyleu/rituals)):
  Work with your team to estimate work, track your progress, and gather feedback.
  - It's a full websocket-driven rich client application, but also works fine without JavaScript

- [TODO Forge](https://todo.kyleu.dev) ([GitHub](https://github.com/kyleu/todoforge)):
  Manages collections of todo items.
  - Almost entirely generated using Project Forge, this is a "stock" application

- [Load Toad](https://loadtoad.kyleu.dev) ([GitHub](https://github.com/kyleu/loadtoad)):
  A tool for uploading HTTP Archives (`.har` files) and running load tests.
  - Also supports client-defined JavaScript, executed in-process on the server

- [Admini](https://admini.dev) ([GitHub](https://github.com/kyleu/admini)):
  A database management application, basically. It does other stuff too.
  - This one is weird, it tried to build a user-defined admin app, but it just ended up looking like a 1990's web portal

- [NPN](https://npn.dev) ([GitHub](https://github.com/kyleu/npn)):
  Basically Postman, it helps you explore and test HTTP services with a focus on speed and correctness.
  - This uses a Vue.js-based rich client application, and a websocket to handle communication

- [Solitaire](https://solitaire.kyleu.dev) ([GitHub](https://github.com/kyleu/solitaire)):
  An example game, not really anything right now.
  - It mainly exists as a testbed for me, and to show the features of Project Forge projects

_More examples coming soon..._

## Documentation

- [Installation](doc/installation.md)
- [Contributing](doc/contributing.md)
- [Customizing](doc/customizing.md)
- [Releasing](doc/releasing.md)
- [Running](doc/running.md)
- [Scripts](doc/scripts.md)
- [Technology](doc/technology.md)

## Licensing

The Project Forge application is released under [MIT](LICENSE.md) license, and all modules are [CC0](https://creativecommons.org/publicdomain/zero/1.0/).
