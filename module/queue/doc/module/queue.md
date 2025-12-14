# Queue

The **`queue`** module provides a database-backed message queue system for your application.
It enables reliable, persistent messaging between application components with topic-based routing and automatic retry mechanisms.

## Overview

This module provides a simple yet robust message queue built on SQLite, offering:

- **Database-backed persistence**: Messages survive application restarts
- **Topic-based messaging**: Organize messages by logical channels
- **Retry mechanisms**: Automatic retry handling with configurable limits
- **Admin interface**: Web-based queue monitoring and management
- **Transactional support**: Database transaction integration

## Key Features

### Reliability
- Persistent message storage using SQLite database
- Automatic retry mechanisms with configurable limits
- Transactional message handling
- Dead letter handling for failed messages

### Performance
- Blocking and non-blocking message retrieval
- Configurable timeout and limit settings
- Efficient database queries with proper indexing
- Minimal overhead for high-throughput scenarios

### Observability
- Message status tracking (sent/received counts)
- Retry count monitoring
- Web-based admin interface for queue inspection
- Integration with core telemetry systems

### Developer Experience
- Simple send/receive API
- Topic-based message organization
- JSON parameter support for structured data
- Built-in admin controllers for debugging

## Package Structure

### Core Queue Components

- **`lib/queue/queue.go`** - Main queue implementation
  - Message sending and receiving operations
  - Topic-based message routing
  - Configurable timeout and retry handling
  - Database transaction integration

- **`lib/queue/message.go`** - Message structure and operations
  - Message ID generation and tracking
  - Topic and parameter handling
  - Retry count management
  - JSON serialization support

- **`lib/queue/status.go`** - Queue status and metrics
  - Sent/received message counting
  - Performance metrics collection
  - Status reporting for monitoring

### Admin Interface

- **`controller/clib/queue.go`** - HTTP controllers for queue management
  - Queue status viewing and monitoring
  - Message inspection and debugging
  - Administrative operations
  - JSON API endpoints for external tools

### Database Schema

- **`queries/`** - SQL DDL and operations
  - Message table definitions with proper indexing
  - Queue-specific database operations
  - Migration scripts for schema updates

## Usage Examples

### Basic Message Sending

```go
import "{{{ .Package }}}/app/lib/queue"

// Create a new queue instance
q := queue.New(db, logger)

// Send a message to a topic
params := map[string]any{
    "userId": 123,
    "action": "send_email",
    "email": "user@example.com",
}

err := q.Send(ctx, "email_notifications", params)
if err != nil {
    return errors.Wrap(err, "failed to send message")
}
```

### Message Processing

```go
// Receive messages from a topic (blocking)
message, err := q.Receive(ctx, "email_notifications", 30*time.Second)
if err != nil {
    return errors.Wrap(err, "failed to receive message")
}

if message != nil {
    // Process the message
    err = processEmailNotification(message.Params)
    if err != nil {
        // Message will be retried automatically
        return errors.Wrap(err, "failed to process message")
    }
}
```

### Non-blocking Operations

```go
// Non-blocking receive (returns immediately)
message, err := q.ReceiveNonBlocking(ctx, "background_tasks")
if err != nil {
    return err
}

if message == nil {
    // No messages available
    return nil
}

// Process message...
```

## Configuration

The queue module supports configuration through environment variables and database settings:

### Queue Behavior
- `queue_timeout` - Default timeout for blocking operations (default: 30s)
- `queue_retry_limit` - Maximum retry attempts per message (default: 3)
- `queue_cleanup_interval` - Interval for cleaning up old messages

### Performance Tuning
- `queue_batch_size` - Number of messages to process in batch operations
- `queue_max_pending` - Maximum number of pending messages per topic
- `queue_worker_count` - Number of concurrent message processors

## Dependencies

This module requires the following:

### Required Modules
- **`sqlite`** - SQLite database support for message persistence
- **`core`** - Foundation framework and utilities

### External Dependencies
- **SQLite database** - For persistent message storage
- **Database migrations** - For schema management

## Admin Interface

The module provides a web-based admin interface accessible at `/admin/queue`:

- **Queue Status**: Overview of all queues and message counts
- **Message Browser**: View pending and processed messages
- **Topic Management**: Monitor topics and their message flow
- **Retry Monitoring**: Track failed messages and retry attempts

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/queue
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [SQLite Module Documentation](sqlite.md) - Database backend
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
