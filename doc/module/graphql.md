# GraphQL

The **`graphql`** module provides comprehensive GraphQL API support for your application.
It enables you to build modern, type-safe GraphQL APIs with an intuitive development interface.

## Overview

This module provides:

- **GraphQL Service**: Schema registration and query execution with OpenTelemetry tracing
- **Multi-Schema Support**: Register and manage multiple GraphQL schemas in a single application
- **GraphiQL Interface**: Interactive query builder and schema explorer
- **Built-in Scalars**: Custom scalars for Time, UUID, and other common types
- **Panic Recovery**: Robust error handling with telemetry integration

## Key Features

### GraphQL Service
- Schema composition with include system
- Execution tracking and telemetry
- Field resolver support using [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)
- Built-in scalars for common Go types

### Development Tools
- GraphiQL web interface for query development
- Schema introspection and documentation
- Real-time query execution with error handling
- Multiple schema management

### Performance & Reliability
- OpenTelemetry distributed tracing integration
- Panic recovery with detailed error reporting
- Execution count tracking per schema
- Efficient schema registration and caching

## Package Structure

### Core Services

- **`lib/graphql/`** - GraphQL service implementation
  - Schema registration and management
  - Query execution with tracing
  - Multi-schema support
  - Built-in scalar types (Time, UUID)

### Controllers

- **`controller/clib/graphql.go`** - HTTP handlers for GraphQL endpoints
  - Schema listing and selection
  - GraphiQL interface serving
  - Query execution via HTTP POST
  - Content negotiation for API responses

### Schema Definition

- **`gql/schema.graphql`** - Base GraphQL schema with includes
  - Built-in scalar definitions
  - Example types and queries
  - Schema composition system

## Usage

### Basic Setup

The GraphQL service is automatically wired into `app.State`. Register schemas and execute queries:

```go
// Register a schema
as.Services.GraphQL.RegisterSchema("my-api", mySchema)

// Execute a query
result := as.Services.GraphQL.Execute(ctx, "my-api", query, variables)
```

### Schema Development

Define your GraphQL schema using the include system in `app/gql/schema.graphql`:

```graphql
# scalar Time
# scalar UUID
# include "example.graphql"

type Query {
  # Your queries here
}
```

### GraphiQL Interface

Access the interactive GraphiQL interface at `/admin/graphql` to:
- Browse schema documentation
- Write and test queries
- Explore available types and fields
- View query execution results

## Configuration

The module works out-of-the-box with no additional configuration required. Schemas are registered programmatically through the service.

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/graphql
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [GraphQL Specification](https://graphql.org) - Official GraphQL documentation
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
