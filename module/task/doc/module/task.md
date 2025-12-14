# Task Engine

The **`task`** module provides a comprehensive task execution engine for your application.
It enables registration, execution, and monitoring of background tasks with real-time progress tracking and a rich web interface.

## Overview

This module provides:

- **Task Registration**: Define and register executable tasks with metadata and parameters
- **Execution Engine**: Synchronous and asynchronous task execution with concurrency control
- **Progress Monitoring**: Real-time progress streaming via WebSocket connections
- **Web Interface**: Rich UI for browsing, configuring, and executing tasks
- **Result Management**: Detailed task results with logging, timing, and structured output

## Key Features

### Task Management
- Register tasks with titles, descriptions, categories, and icons
- Support for task parameters with validation and form generation
- Task categorization and tagging for organization
- Optional danger warnings and expense flags for critical operations

### Execution Engine
- **Synchronous execution**: Immediate task execution with direct results
- **Asynchronous execution**: Background task processing with progress streaming
- **Concurrency control**: Configurable maximum concurrent tasks (defaults to CPU count)
- **Rate limiting**: Prevents system overload during intensive operations

### Real-time Monitoring
- WebSocket-based progress streaming with terminal-style output
- Live status updates during task execution
- Detailed timing information and performance metrics
- Error handling and failure reporting

### Web Interface
- Task browser with search and filtering capabilities
- Dynamic form generation based on task parameters
- Real-time progress display with streaming logs
- Result viewer with structured data presentation

## Package Structure

### Core Components

- **`app/lib/task/`** - Task engine implementation
  - Task registration and management
  - Execution orchestration and scheduling
  - Progress tracking and result collection
  - WebSocket streaming infrastructure

- **`app/controller/task.go`** - HTTP request handlers
  - Task listing and detail pages
  - Task execution endpoints (sync/async)
  - WebSocket connection management
  - Result display and export

### UI Components

- **`views/vtask/`** - Task management interface
  - Task browser and search functionality
  - Parameter form generation and validation
  - Real-time progress display
  - Result presentation and formatting

## Usage Examples

### Registering a Task

```go
// Register a simple task
taskService.Register("backup", &task.Task{
    Title:       "Database Backup",
    Description: "Creates a backup of the application database",
    Category:    "maintenance",
    Icon:        "database",
    Fields: field.Fields{
        field.NewField("format", "", "backup format", field.Select, field.Options{
            "sql": "SQL Dump",
            "tar": "Compressed Archive",
        }),
    },
    Fn: func(ctx context.Context, logger util.Logger, ps *cutil.PageState, args map[string]any) (any, error) {
        format := args["format"].(string)
        // Perform backup logic
        return map[string]any{"file": "backup.sql", "size": 1024000}, nil
    },
})
```

### Executing Tasks

**Synchronous execution:**
```go
result, err := taskService.Run(ctx, "backup", args, logger, ps)
```

**Asynchronous execution with progress:**
```go
run, err := taskService.Start(ctx, "backup", args, logger, ps)
// Monitor progress via run.Status() and run.Progress()
```

### Task Parameters

Tasks support various parameter types:
- **Text fields**: String input with validation
- **Select fields**: Dropdown choices with predefined options
- **Boolean flags**: Checkbox inputs for feature toggles
- **File uploads**: File selection for data processing tasks

## Configuration

The task module supports configuration through environment variables:

### Execution Control
- `TASK_MAX_CONCURRENT` - Maximum concurrent task executions (default: CPU count)
- `TASK_TIMEOUT` - Default timeout for task execution (default: 5 minutes)
- `TASK_WEBSOCKET_ENABLED` - Enable real-time progress streaming (default: true)

### Performance Tuning
- `TASK_BUFFER_SIZE` - WebSocket message buffer size (default: 1000)
- `TASK_LOG_LEVEL` - Task-specific logging level
- `TASK_RETENTION_DAYS` - Days to retain completed task results

## Dependencies

This module requires:

- **[process](../process/process.md)** - Process management and execution infrastructure
- **Core WebSocket support** - For real-time progress streaming
- **Telemetry integration** - Task execution metrics and tracing

## Integration Points

### Service Registration
The task service is automatically registered with the application's service container and available via dependency injection.

### Web Routes
- `/admin/task` - Task browser and management interface
- `/admin/task/{key}` - Individual task detail and execution
- `/admin/task/{key}/start` - Asynchronous task execution endpoint
- `/admin/task/{key}/run` - Synchronous task execution endpoint

### WebSocket Endpoints
- `/ws/task/{runID}` - Real-time progress streaming for running tasks

## Advanced Features

### Custom Task Categories
Organize tasks by category for better navigation:
- `maintenance` - System maintenance and cleanup tasks
- `data` - Data processing and migration tasks
- `admin` - Administrative and configuration tasks
- `report` - Report generation and analytics tasks

### Task Metadata
Enhance tasks with rich metadata:
- **Icons**: Visual indicators using the application's icon system
- **Danger flags**: Warnings for potentially destructive operations
- **Expense flags**: Indicators for resource-intensive tasks
- **Tags**: Additional categorization and filtering options

### Error Handling
Comprehensive error management:
- Detailed error messages with context
- Stack trace collection for debugging
- Automatic retry mechanisms for transient failures
- Graceful degradation for partial failures

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/task
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Process Module Documentation](process.md) - Required process management functionality
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
