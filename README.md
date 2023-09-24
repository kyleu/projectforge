# Project Forge
![app logo](./assets/favicon.png)

[Project Forge](https://projectforge.dev) is an application that allows you to generate, manage, and grow web applications built using the Go language.
You control your application's features, provided via "modules" that enable everything from databases to OAuth.
When creating a new application with Project Forge, a standard Golang project is created, including a ton of utilities and APIs to help you build the application you want, without compromise.

All projects managed by Project Forge provide an HTTP server based on [fasthttp](https://github.com/valyala/fasthttp), and use [quicktemplate](https://github.com/valyala/quicktemplate) for HTML templates (and SQL, if enabled).
An MVC framework is provided (but not required) that handles content negotiation, hierarchical menus, breadcrumbs, OAuth to dozens of providers, stateless user profiles, dark mode support, SVG management, syntax highlighting, form components, and embedded assets.

Project Forge applications can support any UI framework, but the included UI renders a JS-dependency-free page, heavily optimized for speed and modern UX.
The about page is animated, themed, and responsive, and only creates three requests (HTML, CSS, JS) totaling less than 20KB zipped.
It serves in less than a millisecond, and renders in Chrome in less than 20ms.
Progressive enhancement is provided by an included ESBuild TypeScript project, though all functionality is supported with JavaScript disabled.

Your application can (optionally) build for _every_ platform; desktop and mobile webview apps, WASM, and notarized universal macOS binaries.
If you enable all the build options, it will produce almost 60 builds for various platforms. They all produce a ~20MB native binary.
The binaries produced can be configured to auto-upgrade from GitHub Releases, or be upgraded by the user using a CLI or UI (module "upgrade" must be in your project).
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


## Example Applications

- [Rituals.dev](https://rituals.dev) ([GitHub](https://github.com/kyleu/rituals)):
Work with your team to estimate work, track your progress, and gather feedback.

- [TODO Forge](https://todo.kyleu.dev) ([GitHub](https://github.com/kyleu/todoforge)):
  A tool for uploading HTTP Archives (`.har` files) and running load tests.

- [Load Toad](https://loadtoad.kyleu.dev) ([GitHub](https://github.com/kyleu/loadtoad)):
  A tool for uploading HTTP Archives (`.har` files) and running load tests.

- [Admini](https://admini.dev) ([GitHub](https://github.com/kyleu/admini)):
  A database management application, basically. It does other stuff too.

- [NPN](https://npn.dev) ([GitHub](https://github.com/kyleu/npn)):
  Basically Postman, it helps you explore and test HTTP services with a focus on speed and correctness.

- [Solitaire](https://solitaire.kyleu.dev) ([GitHub](https://github.com/kyleu/solitaire)):
  An example game, mainly exists to show the features of Project Forge projects.

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
