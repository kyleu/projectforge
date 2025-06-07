# Scripting

This module provides server-side JavaScript execution capabilities for [Project Forge](https://projectforge.dev) applications using the [Goja](https://github.com/dop251/goja) JavaScript runtime.

## Features

- **JavaScript Execution**: Run ES5+ JavaScript files using the Goja interpreter
- **Web-based Script Management**: Create, edit, and execute scripts through a web UI at `/admin/scripting`
- **Automatic Testing**: Scripts can include example test cases that are automatically executed
- **Filesystem Integration**: Scripts are stored as `.js` files and managed through the filesystem module
- **Search Integration**: Scripts are indexed and searchable within Project Forge applications

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
- View script content and execution results
- Create and edit scripts
- Run scripts and view output
- Test example functions automatically

### JavaScript Environment

Each script runs in an isolated JavaScript VM

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

Scripts are stored as `.js` files in the configured directory. The service provides methods to:

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

When included in a Project Forge application, the module provides these routes:

- `GET /admin/scripting` - List all scripts
- `GET /admin/scripting/{path}` - View specific script and run examples
- `POST /admin/scripting/{path}` - Execute script with custom parameters
- `PUT /admin/scripting/{path}` - Save script content
- `DELETE /admin/scripting/{path}` - Delete script

## Security Considerations

- Scripts can execute arbitrary JavaScript code
- No sandboxing beyond VM isolation is provided
- Scripts have access to all functions exposed to the JavaScript environment
- Only deploy in trusted environments where script execution is acceptable
- Consider network restrictions if scripts might make external calls

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/scripting
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)
