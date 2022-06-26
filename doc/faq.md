# FAQ

## What is Project Forge?

[Project Forge](https://projectforge.dev) helps you build and maintain feature-rich applications written in the Go programming language.

It's a code-generation platform that uses modules, which are self-contained features for your app. 
All managed projects expose a web and CLI interface, and additional modules are available for everything from OAuth to database access.

## What CLI commands are available?

- `all`: (default action) Starts the main http server on port 40000 and the marketing site on port 40001
- `audit`: Audits the project files, detecting invalid files and empty folders
- `build`: Builds the project, many options available
- `completion`: Generate the autocompletion script for the specified shell
- `create`: Creates a new project
- `debug`: Dumps a ton of information about the project
- `doctor`: Makes sure your machine has the required dependencies
- `help`: Help about any command
- `merge`: Merges changed files as required
- `preview`: Show what would happen if you did merge
- `server`: Starts the http server on port 40000 (by default)
- `site`: Starts the marketing site on port 40000 (by default)
- `slam`: Slams all files to the target, ignoring changes
- `svg`: Builds the project's SVG files
- `update`: Refreshes downloaded assets such as modules
- `upgrade`: Upgrades projectforge to the latest published version
- `validate`: Validates the project config
- `version`: Displays the version and exits


## What build options are available?

- `build`: Runs [make build]
- `clean`: Runs [make clean]
- `cleanup`: Cleans up file permissions
- `clientBuild`: Runs [bin/build/client.sh]
- `clientInstall`: Installs dependencies for the TypeScript client
- `deployments`: Manages deployments
- `deps`: Manages Go dependencies
- `format`: Runs [bin/format.sh]
- `imports`: Reorders the imports
- `lint`: Runs [bin/check.sh]
- `packages`: Visualize your application's packages
- `test`: Does a test
- `tidy`: Runs [go mod tidy]


## What modules are available?

- `android`: Webview-based application and Android build
- `audit`: Using the database module, provides an audit framework for tracking changes
- `core`: Provides common utilities for a Go application
- `database`: Provides an API for accessing relational databases
- `databaseui`: Provides a UI for registered databases
- `desktop`: Provides a desktop application using the system's webview
- `docbrowse`: Provides a UI for browsing the documentation
- `export`: Generates code based on the project's schema
- `expression`: Exposes CEL engine for evaluating arbitrary expressions
- `filesystem`: Provides an abstraction around local and remote filesystems
- `graphql`: Supports GraphQL APIs within your application
- `ios`: Webview-based application and iOS build
- `marketing`: Provides a website for downloads, tutorials, and marketing
- `migration`: Database migrations and a common PostgreSQL or SQLite database
- `mysql`: Provides an API for accessing MySQL databases
- `notarize`: Sends files to Apple for notarization
- `oauth`: Provides logins and session management for many OAuth providers
- `postgreSQL`: Provides an API for accessing PostgreSQL databases
- `readonlydb`: Adds a read-only database connection
- `sandbox`: Useful playgrounds for testing custom functions
- `schema`: Classes for representing a collection of models; depends on the types module
- `search`: Adds search facilities to the top-level navigation bar
- `sqlite`: Provides an API for accessing SQLite databases
- `types`: Classes for representing common data types
- `upgrade`: Provides in-place version upgrades using Github Releases
- `websocket`: Provides an API for hosting WebSocket connections
