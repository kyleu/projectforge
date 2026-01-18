# Scripting

This module provides server-side JavaScript execution capabilities for your application using the [Goja](https://github.com/dop251/goja) JavaScript runtime.

## Features

- **JavaScript Execution**: Run ES5+ JavaScript files using the Goja interpreter
- **Web-based Script Management**: Create, edit, and view script execution and example results in a web UI at `/admin/scripting`
- **Automatic Testing**: Scripts can include example test cases that are automatically executed
- **Filesystem Integration**: Scripts are stored as `.js` files and managed through the filesystem module
- **Built-in Helpers**: Console logging plus a small set of utility helpers
- **Search Integration**: Scripts are indexed and searchable within your application

WARNING: This module is marked as dangerous because it executes arbitrary JavaScript on the server. Restrict access to trusted admins.

## Installation

Add this module to your project configuration, by using the UI or editing `project.json` directly:

```json
{
  "modules": ["filesystem", "scripting"]
}
```

## Usage

### Service Setup

Create a new scripting service in your application:

```go
import (
    "your-app/app/lib/filesystem"
    "your-app/app/lib/scripting"
)

// Create filesystem service pointing to script directory
fs := filesystem.NewService("./data")

// Create scripting service with filesystem and script path
scriptService := scripting.NewService(fs, "scripts")
```

### Web Interface

Access the scripting interface at `/admin/scripting` to:
- List all available scripts
- View script content, load results, and example output
- Create and edit scripts
- Delete scripts
- Run example functions automatically

### JavaScript Environment

Each script runs in an isolated JavaScript VM. The runtime exposes a small set of helpers:

- `console.debug`, `console.log`, `console.info`, `console.warn`, `console.error` (server logging)
- `randomString(length)` - Generate random strings
- `microsToMillis(value)` - Convert microseconds to millisecond strings

### Example Testing

Scripts can include automated test examples that will be executed when viewing the script:

```javascript
// Your script functions
function greet(name, title) {
  return `Hello ${title} ${name}!`;
}

function calculate(a, b, operation) {
  switch(operation) {
    case 'add': return a + b;
    case 'multiply': return a * b;
    default: return 0;
  }
}

// Test examples - will be automatically executed
const examples = {
  "greet": [
    ["Alice", "Dr."],
    ["Bob", "Mr."],
    ["Carol", "Prof."]
  ],
  "calculate": [
    [5, 3, "add"],
    [4, 7, "multiply"],
    [10, 2, "divide"]
  ]
};
```

The `examples` object maps function names to arrays of parameter sets. Each parameter set will be passed to the function and the results displayed in the web interface.

### Script Management

Scripts are stored as `.js` files in the configured directory (by default under the filesystem root at `scripts/`). The service provides methods to:

- `ListScripts()` - Get all JavaScript files
- `LoadScript(path)` - Load script content
- `SaveScript(path, content)` - Save script to filesystem
- `DeleteScript(path)` - Remove script file

### Error Handling

Scripts are executed in isolated VMs with comprehensive error handling:
- Syntax errors are caught during script loading
- Runtime errors are captured and displayed
- VM isolation prevents scripts from affecting each other

## API Endpoints

The module provides these routes:

- `GET /admin/scripting` - List all scripts
- `GET /admin/scripting/new` - New script form
- `POST /admin/scripting/new` - Create a new script
- `GET /admin/scripting/{key}` - View a script (loads it and runs examples)
- `GET /admin/scripting/{key}/edit` - Edit script form
- `POST /admin/scripting/{key}/edit` - Save script updates
- `GET /admin/scripting/{key}/delete` - Delete script

## Security Considerations

- Scripts can execute arbitrary JavaScript code
- No sandboxing beyond VM isolation is provided
- Scripts have access to all functions exposed to the JavaScript environment
- Only deploy in trusted environments where script execution is acceptable
- Restrict access to admin-only routes and consider network restrictions if scripts might make external calls

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/scripting
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
