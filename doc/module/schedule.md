# Scheduled Jobs

The **`schedule`** module provides a scheduled job engine and admin UI for your application, built on [gocron](https://github.com/go-co-op/gocron) v2.

## Overview

The schedule module enables applications to run jobs on cron expressions, fixed intervals, or calendar-based schedules. It provides:

- **Job scheduling** with cron, duration, daily/weekly/monthly, and one-time schedules
- **Web UI** for monitoring job status, execution history, and results
- **Telemetry integration** with distributed tracing support
- **Error handling** with panic recovery and structured logging
- **Execution tracking** with run counts and last-result storage

## Features

### Job Management
- Create jobs with custom functions, schedules, and metadata
- Singleton mode to prevent overlapping executions (reschedules when busy)
- Job tagging for organization and filtering
- UUID identifiers with optional friendly names

### Monitoring and Debugging
- Web interface showing all scheduled jobs
- Individual job details with last and next execution times
- Execution count tracking and last result capture
- Optional debug logging control

### Integration
- Telemetry spans for each job execution
- Structured logging with configurable verbosity
- Panic recovery to prevent job failures from crashing the scheduler

### Storage and Retention
- Execution counts and results are kept in memory only
- Only the most recent result is stored per job
- Restarting the app resets counts and results

## Usage

### Basic Setup

The schedule service is created and started automatically when the module is included:

```go
// Service is available in app.State
jobs := as.Services.Schedule.ListJobs()
```

### Creating Jobs

```go
// Define a job function
jobFunc := func(ctx context.Context, logger util.Logger) (any, error) {
    logger.Info("executing my job")
    // Perform work here
    return "job completed", nil
}

// Schedule the job
job, err := as.Services.Schedule.NewJob(
    ctx,
    "daily-cleanup",
    gocron.DailyJob(
        1,
        gocron.NewAtTimes(gocron.NewAtTime(2, 0, 0)),
    ),
    jobFunc,
    true,
    logger,
    "maintenance", "cleanup",
)
```

### Schedule Types

The module supports various schedule types through gocron:

```go
// Cron expressions (set withSeconds true for 6-field cron)
gocron.CronJob("0 2 * * *", false)

// Fixed interval
gocron.DurationJob(time.Hour)

// Daily at 3:30 PM
gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(15, 30, 0)))

// Weekly on Monday at 9 AM
gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(gocron.NewAtTime(9, 0, 0)))

// One-time run 10 minutes from now
gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(10 * time.Minute)))
```

### Job Function Pattern

Job functions should follow this signature:

```go
type JobFunc func(ctx context.Context, logger util.Logger) (any, error)

func myJob(ctx context.Context, logger util.Logger) (any, error) {
    // Use the provided context for cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Perform work
        result := doWork()
        return result, nil
    }
}
```

Return values are stored as the most recent result and rendered in the admin UI. Prefer JSON-friendly types for easy display.

### Accessing Job Information

```go
// List all jobs
jobs := as.Services.Schedule.ListJobs()

// Get specific job
job := as.Services.Schedule.GetJob(jobID)

// Check execution counts (in memory)
count := as.Services.Schedule.ExecCounts[jobID]

// Get last result (in memory)
result := as.Services.Schedule.Results[jobID]
```

## Web Interface and URLs

The module provides admin pages for job monitoring:

- **`/admin/schedule`** - Lists all scheduled jobs with status
- **`/admin/schedule/{id}`** - Detailed view for a specific job

These endpoints respect the standard Project Forge rendering formats via `?format=` (json, yaml, xml, csv) or the `Accept` header.

## CLI

There are no CLI commands for this module. Manage jobs in code or through the admin UI.

## Configuration

No environment variables are required. The scheduler is started automatically.

### Logging Control

Debug logging is enabled by default. Disable it if you want quieter logs:

```go
as.Services.Schedule.EnableLogging(false)
```

### Advanced Scheduling

For advanced gocron options (updates, removal, or global job options), use the underlying scheduler directly:

```go
engine := as.Services.Schedule.Engine
```

## Dependencies

### Required Modules
- **core** - Provides admin routing and layout (always included)

### Recommended Modules
- **user** - Access control for admin routes

### External Libraries
- **[gocron v2](https://pkg.go.dev/github.com/go-co-op/gocron/v2)** - Scheduling engine

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/schedule
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [User Module](user.md) - Access control for admin routes
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
