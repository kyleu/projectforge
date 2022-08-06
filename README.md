<!--- $PF_IGNORE$ -->
# Project Forge
![app logo](./assets/favicon.png)

[Project Forge](https://projectforge.dev) is an application that allows you to generate, manage, and grow web applications built using the Go language. 
You control your application's features, provided via "modules" that enable everything from databases to OAuth. 
A standard Golang project is created, with a ton of utilities and APIs to help you build the application you want, with no compromises.

All projects managed by Project Forge provide an HTTP server based on fasthttp, and use quicktemplate for HTML templates (and SQL, if enabled). 
An MVC framework is provided (but not required) that handles content negotiation, hierarchical menus, breadcrumbs, OAuth to dozens of providers, stateless user profiles, dark mode support, SVG management, syntax highlighting, form components, and embedded assets.

You're welcome to use any UI framework you'd like, but the included UI renders a JS-dependency-free page, heavily optimized for speed and modern UX. 
The about page is animated, themed, and responsive, but is only three requests (HTML, CSS, JS) totaling less than 40KB zipped. 
It serves in less than a millisecond, and renders in Chrome in less than 20ms. 
Progressive enhancement is provided by an included ESBuild TypeScript project, though all functionality is supported with JavaScript disabled.

Your project can (optionally) build for _every_ platform; desktop and mobile webview apps, WASM, universal macOS binaries, frickin' Plan9 and Solaris. 
If you check all the boxes it'll produce like 60 builds. They all produce a ~20MB native binary. 
The binaries produced can be configured to auto-upgrade from GitHub Releases, or upgraded by the user using a CLI or UI (module "upgrade" must be in your project). 
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

## Documentation

- [Installation](doc/installation.md)
- [Contributing](doc/contributing.md)
- [Customizing](doc/customizing.md)
- [Contributing](doc/contributing.md)
- [Technology](doc/technology.md)

## Licensing

The Project Forge application is released under [MIT](LICENSE.md) license, and all modules are [CC0](https://creativecommons.org/publicdomain/zero/1.0/).
