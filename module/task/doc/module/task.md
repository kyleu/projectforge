# Task Engine

The `task` module provides a task registry and execution flow with an admin UI. Tasks are registered in code, accept query/field inputs, and return a structured Result with logs and data. Runs can be synchronous or streamed asynchronously over WebSocket.

## Overview

- Register tasks with metadata (title, category, icon, tags, warnings)
- Define input fields using `util.FieldDesc` for UI forms and query parsing
- Execute tasks synchronously or via WebSocket-backed async runs
- Capture structured results, logs, tags, and timing details

## Package Structure

### Core Components

- `app/lib/task/` - task types and execution helpers
  - `Task`, `Result`, and `Service`
  - Task registration and execution
  - Result logging helpers

- `app/controller/clib/task.go` - HTTP handlers for the admin UI
- `views/vtask/` - task list, detail, and result views

## Usage

### Registering a Task

```go
tf := func(ctx context.Context, res *task.Result, logger util.Logger) *task.Result {
    res.Log("starting backup")
    format := res.Args.GetStringOpt("format")
    if format == "" {
        format = "sql"
    }
    // Perform backup logic...
    return res.Complete(util.ValueMap{"file": "backup." + format})
}

t := task.NewTask("backup", "Database Backup", "maintenance", "database", "Creates a database backup", tf)
t.Fields = util.FieldDescs{
    {Key: "format", Title: "Format", Type: "string", Default: "sql", Choices: []string{"sql", "tar"}},
}
t.Expensive = true
t.Dangerous = "May lock tables during backup"
t.Tags = []string{"db", "maintenance"}

if err := as.Services.Task.RegisterTask(t); err != nil {
    return err
}
```

### Running Tasks in Code

```go
args := util.ValueMap{"format": "tar"}
res := t.Run(ctx, "ad-hoc", args, logger)
if !res.IsOK() {
    logger.Warnf("task failed: %s", res.Error)
}
```

For batch execution with rate limiting, use the service helper:

```go
t.MaxConcurrent = 4
results, err := as.Services.Task.RunAll(ctx, t, "batch", []util.ValueMap{
    {"format": "sql"},
    {"format": "tar"},
}, logger)
```

### Running Tasks via HTTP/UI

A full UI is provided at `/admin/task`.

Task routes accept query parameters for task fields. The `category` query parameter is reserved for the run label.

- `GET /admin/task` - list registered tasks
- `GET /admin/task/{key}` - task detail + form
- `GET /admin/task/{key}/run` - execute synchronously and render result
- `GET /admin/task/{key}/run?async=true` - render the async UI, which opens a WebSocket to start the task
- `GET /admin/task/{key}/start` - WebSocket upgrade endpoint used by the async UI
- `GET /admin/task/{key}/remove` - remove a dynamically registered task

All task pages support the standard `?format=` query param (json, csv, xml, yaml, debug).

## Task Fields and Metadata

`Task.Fields` uses `util.FieldDesc` to describe input fields:

- `Type` supports `string`, `int`, `float`, `bool`, `[]string`, and `time`
- `Default` provides fallback values when missing
- `Choices` populates select lists in the UI

Additional metadata:

- `Dangerous` (string) and `Expensive` (bool) for warnings/labels
- `Tags` for categorization
- `WebURL` to override the default `/admin/task/{key}` URL

## Concurrency and Rate Limiting

`Service.RunAll` uses `Task.MaxConcurrent` to limit parallelism:

- `0` uses `runtime.NumCPU()`
- `-1` uses a fixed limit of `128`

Single `Run` calls are not rate limited by default.

## Configuration

This module does not define environment variables. Configure behavior per task using `Task.MaxConcurrent`, `Task.Fields`, and your task function implementation.

## Dependencies

- **[process](process.md)**
- **[exec](exec.md)** (optional, for OS process-backed execution)
- **[websocket](websocket.md)** (optional, for live monitoring UI)

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/task
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Process Module Documentation](process.md)
- [Project Forge Documentation](https://projectforge.dev)
