# AGENT.md

This file provides guidance for AI agents and humans working with the Project Forge codebase.

## Table of Contents

- [Project Overview](#project-overview)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Build Commands](#build-commands)
- [Testing](#testing)
- [Code Style & Conventions](#code-style--conventions)
- [Module System](#module-system)
- [Key Technologies](#key-technologies)
- [Common Patterns](#common-patterns)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Additional Resources](#additional-resources)

## Project Overview

**Project Forge** is a web application generator and management tool for Go applications. It provides:

- **Module-based architecture**: Mix and match features via self-contained modules
- **Multi-platform builds**: Desktop, mobile, WASM, and traditional server applications
- **Zero-JS functionality**: Full functionality without JavaScript, with progressive enhancement
- **High performance**: Sub-second builds, <20KB total payload, <1ms response times
- **Developer experience**: Instant compiles, comprehensive tooling, extensive documentation

### What Project Forge Does

1. **Generates** complete Go web applications with your chosen feature set
2. **Manages** application lifecycle through updates and module changes  
3. **Builds** applications for 60+ platform/architecture combinations
4. **Provides** a rich MVC framework with templating, routing, and utilities

Project Forge is itself an application managed by Project Forge

## Quick Start

```bash
# Build the application
make build

# Run in development mode with live reload
./bin/dev.sh # or `make dev`

# Run tests
./bin/test.sh # or `make test`

# Lint code
./bin/check.sh

# Format code
./bin/format.sh

# Build for production
./bin/build/build.sh # or `make build-release`
```

## Project Structure

```
app/                   # Application code
├── cmd/               # CLI commands and the main entrypoint of the app
├── controller/        # HTTP request handlers (MVC controllers, usually)
├── lib/               # Common logic and services, usually provided by modules
└── util/              # Utility functions and helpers
assets/               # Static files (CSS, JS, images)
bin/                  # Build and development scripts
client/               # TypeScript/JavaScript client code
doc/                  # Project documentation
help/                 # Embedded help files
module/               # Available Project Forge modules
test/                 # Integration and E2E tests
tools/                # Platform-specific build tools
views/                # HTML templates (quicktemplate)
```

### Key Directories

- **`/app/controller/`**: HTTP handlers organized by feature area
- **`/app/lib/`**: Module-provided logic, organized by module
- **`/views/`**: HTML templates using quicktemplate syntax
- **`/module/`**: Self-contained feature modules
- **`/bin/`**: Development and build automation scripts

## Development Workflow

### 1. Making Changes

```bash
# Start development server with live reload
./bin/dev.sh

# The server will automatically:
# - Rebuild Go code on changes
# - Recompile templates on changes  
# - Rebuild client assets on changes
# - Restart the server as needed
```

### 2. Before Committing

```bash
# Format code
./bin/format.sh

# Run linting and validation
./bin/check.sh

# Run tests
./bin/test.sh

# Ensure everything builds
make build
```

### 3. Template Development

```bash
# Compile templates manually
./bin/templates.sh

# Templates use quicktemplate syntax (.html files -> .html.go files)
# Located in /views/ directory
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Build debug binary |
| `make build-release` | Build optimized release binary |
| `make dev` | Start development server with live reload |
| `make lint` | Run linters and code quality checks |
| `make clean` | Remove build artifacts and compiled templates |
| `make templates` | Compile quicktemplate files |
| `make help` | Show available make targets |

### Script Commands

| Script | Description |
|--------|-------------|
| `./bin/dev.sh` | Development server with live reload |
| `./bin/check.sh` | Lint, format check, and validation |
| `./bin/format.sh` | Format Go code with gofumpt |
| `./bin/test.sh` | Run all tests |
| `./bin/test.sh -w` | Run tests in watch mode |
| `./bin/test.sh -c` | Run tests with clean cache |
| `./bin/templates.sh` | Compile quicktemplate files |
| `./bin/coverage.sh` | Generate test coverage report |

## Testing

### Running Tests

```bash
# All tests
./bin/test.sh # or gotestsum -- ./app/...

# Single package
go test ./app/util -v

# Single test  
go test ./app/util -run TestPlural

# With coverage
./bin/coverage.sh

# Watch mode
./bin/test.sh -w

# Clean cache and run
./bin/test.sh -c
```

### Test Organization

- **Unit tests**: Alongside source code with `_test.go` suffix
- **Integration tests**: In `/test/` directory
- **E2E tests**: Playwright tests in `/test/playwright/`

## Code Style & Conventions

### Go Code Standards

- **Formatting**: Use `bin/format.sh` or `gofumpt` (enforced by linters)
- **Max line length**: 160 characters
- **Function length**: Max 100 lines recommended
- **Cyclomatic complexity**: Max 30 recommended
- **Error handling**: Always check errors, provide context
- **Naming**: PascalCase for exported, camelCase for internal

### Import Organization

```go
import (
    // Standard library
    "context"
    "fmt"
    
    // Third-party
    "github.com/gorilla/mux"
    "go.uber.org/zap"
    
    // Project imports
    "projectforge.dev/projectforge/app/util"
)
```

### Template Conventions

- Use (quicktemplate)[https://github.com/vayala/quicktemplate] syntax (not html/template)
- Templates in `/views/` compile to `.html.go` files
- Follow existing naming patterns (PascalCase for template names)
- Leverage existing component templates in `/views/components/`

## Module System

### Understanding Modules

Modules are self-contained feature packages that can be mixed and matched:

- **Location**: `/module/{module-name}/`
- **Configuration**: Each module has a `.module.json` file
- **Documentation**: Module docs in `/module/{name}/doc/module/{name}.md`
- **Dependencies**: Modules can depend on other modules


### Available Modules

See [README.md](README.md#available-modules) for a complete list.

### Module Structure

The files in the `module` directory are templates, using Go templating syntax, using `{{{` and `}}}` as tokens

The files are run through the template engine using the project's model as the template data, producing files in the generated applications

## Key Technologies

### Core Stack

See [technology.md](doc/technology.md) for a complete list.

### Frontend

- **TypeScript**: Progressive enhancement client code
- **ESBuild**: Fast asset bundling
- **CSS**: Modern CSS with CSS custom properties
- **SVG**: Embedded icon system

### Development Tools

- **Air**: Live reload for development
- **gofumpt**: Opinionated Go formatting
- **golangci-lint**: Comprehensive Go linting
- **gotestsum**: Enhanced test runner

## Common Patterns

### Controllers (MVC)

Controllers handle HTTP requests and are organized by feature:

```go
func Home(w http.ResponseWriter, r *http.Request) {
	Act("home", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := "any variable" // will be used for JSON and other data-oriented content types; it can be any type
		ps.SetTitleAndData("My Home Page", ret) // alternately, set properties of `ps` directly
        page := &views.Debug{} // HTML quicktemplate file; this one just displays the data set in the line above
		return Render(r, as, page, ps) // this will render the page or return the content type's data; return `"/foo", nil` to redirect
	})
}
```

### Services

Business logic lives in service packages:

```go
type Service struct {
    logger *zap.Logger
}

func (s *Service) SomeFeature(ctx context.Context) error {
    // Business logic here
}
```

### Error Handling

```go
import "github.com/pkg/errors"

func process() error {
    data, err := fetch()
    if err != nil {
        return errors.Wrap(err, "failed to fetch data")
    }
    return nil
}
```

### Configuration

Use environment variables with sensible defaults. Document the allowed envvars in [running.md](doc/running.md):

```go
count := util.GetEnvInt("my_variable", 42000)
```

## Configuration

### Project Configuration

Project settings for both Project Forge and the projects it manages live in `.projectforge/project.json`. See [project.schema.md](assets/schema/project.schema.md) for the definition

### Environment Variables

Common environment variables:

- `PORT`: HTTP server port (default: 42000)
- `DEBUG`: Enable debug logging (default: false)
- `CONFIG_DIR`: Configuration directory override
- `ADDR`: Server bind address (default: 0.0.0.0)

## Troubleshooting

### Common Issues

**Build fails with template errors:**
```bash
# Recompile templates
./bin/templates.sh
make build
```

**Tests failing:**
```bash
# Clean test cache
./bin/test.sh -c
```

**Linting errors:**
```bash
# Auto-fix formatting
./bin/format.sh

# Check what's wrong
./bin/check.sh
```

**Development server not reloading:**
- Check that `air` is installed and working
- Verify `.air.toml` configuration exists
- Restart with `./bin/dev.sh`

### Performance Issues

- Use `go tool pprof` for profiling (see `/bin/util/` scripts)
- Check telemetry data for bottlenecks
- Review database query patterns
- Verify proper HTTP caching headers

## Additional Resources

### Documentation

- [Installation Guide](doc/installation.md)
- [Running Guide](doc/running.md) 
- [Contributing Guide](doc/contributing.md)
- [Technology Overview](doc/technology.md)
- [Module Documentation](README.md#available-modules)

### Build Scripts

Explore `/bin/` directory for additional utilities:
- `/bin/util/`: Profiling and analysis tools
- `/bin/workspace.sh`: Helper script for macOS iTerm2 users
- `/bin/tag.sh`: Version tagging
- `/bin/export.sh`: Project export

### Development Files

- `Makefile`: Primary build targets
- `.air.toml`: Live reload configuration (if present)  
- `go.mod`: Go dependency management
- `client/package.json`: Frontend dependencies

### Getting Help

- Check existing documentation in `/doc/`
- Review module-specific docs in `/module/{name}/doc/`
- Examine similar patterns in existing code
- Look at example applications in README.md

---

*This file is maintained for both AI agents and human developers. When making changes, ensure information remains accurate and examples stay current with the codebase.*