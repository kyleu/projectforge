# Process

The **`process`** module provides comprehensive system process management capabilities for [Project Forge](https://projectforge.dev) applications. It enables spawning, monitoring, and controlling system processes with real-time output streaming and web-based administration.

## Overview

This module provides:

- **Process Execution**: Spawn and manage system processes with custom commands and environments
- **Real-time Monitoring**: Live output streaming via WebSockets with process state tracking
- **Web Administration**: Browser-based interface for creating, monitoring, and controlling processes
- **Process History**: Execution tracking with PIDs, exit codes, and timing information

⚠️ **Security Note**: This module is marked as **dangerous** because it allows executing arbitrary system commands. Use with caution and appropriate access controls.

## Key Features

### Process Management
- Start processes with custom commands, arguments, and environment variables
- Monitor process state (running, completed, failed)
- Terminate running processes
- Track process IDs, exit codes, and execution duration

### Real-time Streaming
- Live terminal output via WebSocket connections (requires `websocket` module)
- Buffered output for processes that complete before connection
- Thread-safe concurrent process management

### Web Interface
- Admin panel at `/admin/exec` for process management
- Create new processes with custom parameters
- View running processes and execution history
- Real-time output display with WebSocket integration

### Performance & Safety
- Non-blocking process execution
- Proper resource cleanup and process termination
- Configurable output buffering
- Thread-safe operations for concurrent access

## Usage

### Basic Process Execution

```go
// Create a new exec service
execSvc := exec.NewService()

// Create and start a process
proc := execSvc.NewExec("echo", []string{"Hello, World!"}, "", nil)
err := proc.Start()
if err != nil {
    return err
}

// Wait for completion
err = proc.Wait()
```

### Web Interface

Navigate to `/admin/exec` to:
1. Create new processes with custom commands
2. Monitor running processes
3. View real-time output (with WebSocket support)
4. Terminate processes as needed

### Real-time Monitoring

Include the `websocket` module to enable live output streaming:

```html
<!-- Process detail view with WebSocket streaming -->
{%= components.RenderExecDetail(exec, ps) %}
```

## Configuration

The process module uses standard Project Forge configuration patterns:

- Process timeout limits
- Output buffer sizes
- WebSocket connection settings
- Security restrictions for allowed commands

## Dependencies

### Required Modules
- **`websocket`** - Required for real-time output streaming

### Recommended Modules
- **`user`** - For access control to process management
- **`audit`** - For logging process execution events

## Security Considerations

⚠️ **Important**: This module allows execution of arbitrary system commands. Implement appropriate security measures:

- Restrict access to trusted users only
- Validate and sanitize process commands
- Consider running in sandboxed environments
- Monitor and audit process execution
- Implement command whitelisting if needed

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/process
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [WebSocket Module](websocket.md) - Required for real-time features
- [User Module](user.md) - For access control integration
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
