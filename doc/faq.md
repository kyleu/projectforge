# FAQ

1. [What is Project Forge?](#faq-1)
2. [What modules are available?](#faq-2)
3. [What CLI commands are available?](#faq-3)
4. [What actions are available for my project from the UI?](#faq-4)
5. [What build options are available?](#faq-5)
6. [What git actions are available?](#faq-6)
7. [How do I make a new HTTP action?](#faq-7)
8. [What is provided in the TypeScript application?](#faq-8)
9. [How do I use custom modules?](#faq-9)
10. [How do I manage more than one Project Forge project at once?](#faq-10)
11. [How are static assets served?](#faq-11)
12. [How do I work with SVG icons?](#faq-12)
13. [What's this PF_SECTION crap?](#faq-13)
14. [What administrative functions are available?](#faq-14)
15. [How can I secure my application?](#faq-15)
16. [What options are available for the HTML menu and breadcrumbs?](#faq-16)
17. [What search facilities are available?](#faq-17)
18. [What utility methods are available for my app?](#faq-18)
19. [How do I add a new model to the "export" facilities?](#faq-19)

## What is Project Forge?<a name="faq-1"></a>

[Project Forge](https://projectforge.dev) helps you build and maintain feature-rich applications written in the Go programming language.

It's a code-generation platform that uses modules, which are self-contained features for your app.
All managed projects expose a web and CLI interface, and additional modules are available for everything from OAuth to database access.


## What modules are available?<a name="faq-2"></a>

- `android`: Webview-based application and Android build
- `audit`: Audit framework for tracking changes in the system
- `brands`: Provides thousands of SVG icons from `simple-icons` for representing common logos
- `core`: Common utilities for a Go application
- `database`: API for accessing relational databases
- `databaseui`: UI for registered databases, allows auditing and tracing
- `desktop`: Desktop application for all major operating systems, using the system's webview
- `docbrowse`: UI for browsing the Markdown documentation as HTML
- `export`: Generates code based on the project's schema
- `expression`: CEL engine for evaluating arbitrary expressions
- `filesystem`: An abstraction around local and remote filesystems
- `git`: Helper classes for performing operations on git repositories
- `graphql`: Supports GraphQL APIs within your application
- `har`: Utilities for working with HTTP Archive files
- `help`: Markdown-backed help files that integrate into the UI
- `ios`: Webview-based application and iOS build
- `jsx`: A custom JSX engine for reactive UIs without React
- `marketing`: Serves a website for downloads, tutorials, and marketing
- `migration`: Database migrations and a common database pool
- `mysql`: API for accessing MySQL databases
- `notarize`: Sends files to Apple for notarization
- `notebook`: Provides an Observable Framework notebook
- `oauth`: Login and session management for many OAuth providers
- `openapi`: Embeds the Swagger UI, using your OpenAPI specification
- `playwright`: Adds a project for testing the UI using playwright.dev
- `postgres`: API for accessing PostgreSQL databases
- `process`: Framework and UI for managing system processes.
- `proxy`: Provides an HTTP proxy while still enforcing this app's security
- `queue`: Provides a simple message queue based on SQLite
- `readonlydb`: Read-only database connection
- `richedit`: Provides a rich editing experience with a decent fallback when scripting is disabled
- `sandbox`: Useful playgrounds for testing custom functions
- `schedule`: Provides a scheduled job engine and UI based on gocron
- `scripting`: Allows the execution of JavaScript files using a built-in interpreter
- `search`: Adds search facilities to the top-level navigation bar
- `sqlite`: Provides an API for accessing SQLite databases
- `system`: Provides logic and a UI for getting the status of the system the app is running on
- `themecatalog`: A dozen default themes, and facilities to create additional
- `types`: Classes for representing common data types
- `upgrade`: In-place version upgrades using GitHub Releases
- `user`: Classes for representing persistent user records, application usage
- `wasmclient`: WASM library and HTML host for a custom WASM application
- `wasmserver`: Runs your unmodified HTTP server as a Service Worker
- `websocket`: API for hosting WebSocket connections


## What CLI commands are available?<a name="faq-3"></a>

- `all`: (default action) Starts the main http server on port 40000 and the marketing site on port 40001
- `audit`: Audits the project files, detecting invalid files and configuration
- `build`: Builds the project, many options available, see below
- `completion`: Generate the autocompletion script for the specified shell
- `create`: Creates a new project
- `debug`: Dumps information about the project
- `doctor`: Makes sure your machine has the required dependencies
- `generate`: Applies pending changes to files as required
- `git`: Handles git operations, many options available, see below
- `help`: Help about any command
- `preview`: Shows what would happen if you generate
- `server`: Starts the http server on port 40000 (by default)
- `site`: Starts the marketing site on port 40000 (by default)
- `svg`: Builds the project's SVG files
- `update`: Refreshes downloaded assets such as modules
- `upgrade`: Upgrades projectforge to the latest published version
- `validate`: Validates the project config
- `version`: Displays the version and exits


## What actions are available for my project from  the UI?<a name="faq-4"></a>

- `preview`: Shows what would happen if you generate
- `generate`: Applies pending changes to files as required
- `audit`: Audits the project files, detecting invalid files and configuration
- `build`: Builds the project, many options available, see below
- `files`: Browse your project's filesystem
- `svg`: Options for managing the SVG icons in the system
- `git`: Handles git operations, many options available, see below


## What build options are available?<a name="faq-5"></a>

- `build`: Runs [make build]
- `clean`: Runs [make clean]
- `cleanup`: Cleans up file permissions
- `clientBuild`: Runs [bin/build/client.sh]
- `clientInstall`: Installs dependencies for the TypeScript client
- `deployments`: Manages deployments
- `deps`: Manages Go dependencies
- `format`: Runs [bin/format.sh]
- `imports`: Organizes the imports in source files and templates
- `lint`: Runs [bin/check.sh]
- `lint-client`: Runs [bin/check-client.sh]
- `packages`: Visualize your application's packages
- `test`: Does a test
- `tidy`: Runs [go mod tidy]


## What git actions are available?<a name="faq-6"></a>

Common git actions like `status`, `fetch`, `pull` and `commit` do exactly what you'd expect.
There's also a `magic` action, which is a little crazy:
- If there are uncommitted changes, it will stash them
- If there are pending upstream commits, it will pull them
- If a stash was previously created, it will pop and commit it
- If commits exist that haven't been pushed, it will push them


## How do I make a new HTTP action?<a name="faq-7"></a>

Create a new file in `app/controller` or a child directory, then add your method.
The usual form is:

```go
// All controller actions require an http request/response
func CurrentTime(w http.ResponseWriter, r *http.Request) {
	// The Act method wires up the page state for your logic
	Act("current.time", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		// The PageState::Title is used as the HTML title
		ps.Title = "The current server time for " + util.AppName
		// For example purposes
		t := util.TimeCurrent()
		// PageState::Data will be rendered as JSON or XML if the Content-Type of the request matches
		ps.Data = t
		// The Render method will send the template contents if HTML is requested. The final argument refers to the active menu key
		return Render(r, as, &views.CurrentTime{Time: t}, ps, "time")
	})
}
```

Add your action to `./app/controller/routes/routes.go` like so:

```go
makeRoute(r, http.MethodGet, "/time", controller.CurrentTime)
```

Then create a menu item in `./app/controller/cmenu/menu.go` like so:

```go
i := &menu.Item{Key: "time", Title: "Current Time", Description: "...", Icon: "star", Route: "/time"}
ret = append(ret, i)
```

Templates are written in [quicktemplate](https://github.com/valyala/quicktemplate), and generally take this form:

```html
{% code type CurrentTime struct {
  layout.Basic
  Time time.Time
} %}

{% func (p *CurrentTime) Body(as *app.State, ps *cutil.PageState) %}
  <div class="card">
    <h3>{%= components.SVGIcon(`calendar`, ps) %} The Current Time!</h3>
    <p>{%s fmt.Sprint(p.Time) %}</p>
  </div>
{% endfunc %}
```

The `layout.Basic` renders menu and navigation, and is swappable for custom response types.


## What is provided in the TypeScript application?<a name="faq-8"></a>

The TypeScript project (available at `./client`) is dependency-free, lightweight, and is built with ESBuild.
Most code is wired in automatically when the build's `.js` output is requested. 
See the code in `./client/src` for details.

CSS is also built by ESBuild, see `./client/src/client.css` and the files in `./client/src/style`.


## How do I use custom modules?<a name="faq-9"></a>

If you want to make your own module, you can load it in Project Forge by adding the following section to the "info" object in `./.projectforge/project.json`:

```json
{
  "moduleDefs": [
    {
      "key": "foo",
      "path": "./module_foo",
      "url": "https://github.com/me/foo"
    },
    {
      "key": "*",
      "path": "../modules",
      "url": "https://github.com/me/modules"
    }
  ]
}
```


## How do I manage more than one Project Forge project at once?<a name="faq-10"></a>

Create the file `./projectforge/additional-projects.json`, containing a JSON array of string paths to other projects.


## How are static assets served?<a name="faq-11"></a>

Assets are embedded into the servier binary at build time using Go's `embed` package.
This brings a lot of advantages, but you'll need to rebuild before you see changes to your assets.
The script `./bin/dev.sh` handles the rebuild/restart automatically.


## How do I work with SVG icons?<a name="faq-12"></a>

To see available icons, click the `SVG` button on your project's page.
To add a new icon, enter the name of the icon from [Line Awesome](https://icons8.com/line-awesome) or paste a URL.
When added, the icon is rewritten to support themes/styling and references.

To reference an icon, add `{%= components.SVG("star") %}` in your template.
Other helper methods available, see `./views/compnents/SVG.html`


## What's this PF_SECTION crap?<a name="faq-13"></a>

Most files are either fully-generated by Project Forge, or completely custom to your application.
Rarely, files are managed but also support custom code. The text you place inside a PF_SECTION will be preserved when your project is updated.

To exclude a file from code generation entirely, add it to the `ignoredFiles` array in your project's configuration.


## What administrative functions are available?<a name="faq-14"></a>

Assuming you've got access, going to `/admin` will show you the available actions:

- CPU/memory profiling and visualization, runtime heap dumps
- Garbage collection management and metrics
- View your projects Go modules, internal and external
- Theme management, with theme catalogs and live preview
- Lots of HTTP tools, for investigating sessions or request


## How can I secure my application?<a name="faq-15"></a>

TODO


## What options are available for the HTML menu and breadcrumbs?<a name="faq-16"></a>

TODO


## What search facilities are available?<a name="faq-17"></a>

TODO


## What utility methods are available for my app?<a name="faq-18"></a>

_So many_. See the files in `./app/util`, there's a ton of juicy stuff.


## How do I add a new model to the "export" facilities?<a name="faq-19"></a>

Export configuration is stored in `./.projectforge/export`.
A rudimentary UI is provided in the "Export" section of your project's page. 
To be honest, the easiest way to do this is to copy/paste and modify an existing model definition.
There's some good examples at [rituals.dev](https://github.com/kyleu/rituals/tree/master/.projectforge/export/models).
