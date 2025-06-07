# Scheduled Jobs

A Project Forge module that provides a scheduled job engine and web UI based on [gocron](https://github.com/go-co-op/gocron).

## Overview

The schedule module enables applications to run jobs on a schedule using cron expressions or intervals. It provides:

- **Job scheduling** with cron expressions and interval-based execution
- **Web UI** for monitoring job status, execution history, and results  
- **Telemetry integration** with distributed tracing support
- **Error handling** with panic recovery and structured logging
- **Execution tracking** with run counts and result storage

## Features

### Job Management
- Create jobs with custom functions, schedules, and metadata
- Singleton mode to prevent overlapping executions
- Job tagging for organization and filtering
- Execution count tracking per job

### Monitoring & Debugging
- Web interface showing all scheduled jobs
- Individual job details with execution history
- Last and next execution times
- Result storage with return values and errors
- Optional debug logging control

### Integration
- Telemetry spans for each job execution
- Structured logging with configurable verbosity
- Panic recovery to prevent job failures from crashing the scheduler

## Usage

### Basic Setup

The schedule service is automatically initialized and started when the module is included:

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
    "daily-cleanup",           // job name
    gocron.DailyAt("02:00"),  // schedule (cron or interval)
    jobFunc,                  // function to execute
    true,                     // singleton mode
    logger,                   // logger instance
    "maintenance", "cleanup", // tags
)
```

### Schedule Types

The module supports various schedule types through gocron:

```go
// Cron expressions
gocron.CronJob("0 2 * * *", false) // Daily at 2 AM

// Intervals
gocron.DurationJob(time.Hour)       // Every hour
gocron.DailyAt("15:30")            // Daily at 3:30 PM
gocron.WeeklyOn(time.Monday, "09:00") // Weekly on Monday at 9 AM
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

### Accessing Job Information

```go
// List all jobs
jobs := as.Services.Schedule.ListJobs()

// Get specific job
job := as.Services.Schedule.GetJob(jobID)

// Check execution counts
count := as.Services.Schedule.ExecCounts[jobID]

// Get last result
result := as.Services.Schedule.Results[jobID]
```

## Web Interface

The module provides admin pages for job monitoring:

- **`/admin/schedule`** - Lists all scheduled jobs with status
- **`/admin/schedule/{id}`** - Detailed view for a specific job

The interface shows:
- Job ID, name, and tags
- Last and next execution times  
- Total execution count
- Latest execution result and any errors

## Configuration

### Environment Variables

- Jobs inherit the application's logging configuration
- Debug logging for job execution can be controlled programmatically:

```go
as.Services.Schedule.EnableLogging(false) // Disable debug logs
```

### Service Options

When creating jobs, you can specify:
- **Singleton mode**: Prevents overlapping executions
- **Tags**: For categorization and filtering
- **Custom names**: For easier identification in logs and UI

## Implementation Details

### Components

- **`Service`**: Main scheduler service wrapping gocron
- **`Job`**: Job metadata and execution information  
- **`Result`**: Execution results with timing and error data
- **Controllers**: Web interface for job monitoring

### Thread Safety

The service uses mutexes to ensure thread-safe access to:
- Execution counts map
- Results storage
- Job state information

### Error Handling

- Panics are recovered and logged as errors
- Job errors are captured and stored in results
- Failed jobs don't affect other scheduled jobs
- Scheduler continues running even after job failures

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/schedule
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)
