# Export

The **`export`** module is a code generation module for [Project Forge](https://projectforge.dev) applications. It generates comprehensive Go code based on your project's schema definition, creating standardized service interfaces, model contracts, and data export functionality.

## Overview

This module is **optional** and marked as **dangerous** due to its code generation capabilities. When enabled, it provides:

- **Service Interface Generation**: Generic CRUD service interfaces for all schema models
- **Model Contracts**: Standardized interfaces and implementations for data models
- **Data Export**: CSV and other format export capabilities for all models
- **Search Functionality**: Integrated search interfaces for model queries
- **Soft Delete Support**: Configurable soft delete patterns for models

## Key Features

### Code Generation
- Generates Go service interfaces from schema definitions
- Creates standardized CRUD operations for all models
- Implements consistent error handling patterns
- Provides type-safe model interfaces

### Data Management
- Generic service layer with configurable implementations
- Soft delete support with automatic filtering
- Search and filtering capabilities
- Batch operations and transactions

### Export Capabilities
- CSV export for all generated models
- Configurable export formats
- Data transformation pipelines
- Streaming export for large datasets

## Package Structure

### Generated Services

- **Service Interfaces** - Generic CRUD operations for each model
  - Create, Read, Update, Delete operations
  - List and search functionality
  - Bulk operations support
  - Transaction-aware implementations

- **Model Contracts** - Standardized model interfaces
  - Type-safe property access
  - Validation and serialization
  - Export format implementations
  - Search index integration

### Export Infrastructure

- **Export Services** - Data export functionality
  - Multiple format support (CSV, JSON, etc.)
  - Streaming for large datasets
  - Configurable field selection
  - Data transformation hooks

## Configuration

This module generates code based on your project's schema configuration. The generated services respect:

- Model-specific configuration options
- Soft delete settings per model
- Custom validation rules
- Export field configurations

## Important Notes

The generated code follows Project Forge conventions and integrates seamlessly with other modules.

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/export
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
