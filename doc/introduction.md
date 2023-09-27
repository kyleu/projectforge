# Introduction

[Project Forge](https://projectforge.dev) is an application that allows you to generate, manage, and grow web applications built using the Go programming language.
You control your application's features, provided via "modules" that enable everything from databases to OAuth.
When creating a new application with Project Forge, a standard Golang project is created, which includes utilities and APIs to help you build the application you want, without compromise.

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

