# Git

The **`git`** module provides comprehensive Git repository management and automation capabilities for your application.
It offers both basic Git operations and intelligent "magic" workflows that automate common development tasks.

## Overview

This module enables your application to:

- **Repository Management**: Perform standard Git operations (status, commit, push, pull, etc.)
- **Automated Workflows**: Use "magic" commands that intelligently handle common Git scenarios
- **History Analysis**: Examine commit history and repository state
- **GitHub Integration**: Work with GitHub repositories and operations
- **Repository Creation**: Initialize and set up new Git repositories

⚠️ **Note**: This module is marked as "dangerous" as it performs system-level Git operations that can modify repository state.

## Key Features

### Core Git Operations
- Repository status checking and analysis
- Commit creation with automatic staging
- Push/pull operations with conflict detection
- Branch management and switching
- Stash operations for temporary changes

### Magic Workflows
- Intelligent commit automation based on repository state
- Automatic conflict resolution strategies
- Smart push/pull decisions based on repository analysis
- Stash management during complex operations

### Advanced Features
- Commit history analysis and traversal
- Repository state validation
- GitHub API integration
- Repository creation and initialization
- Outdated branch detection

## Package Structure

### Core Service

- **`Service`** - Main Git service that manages repository operations
  - Configurable repository path
  - Context-aware operations
  - Comprehensive logging integration

### Repository Operations

- **`status.go`** - Repository status checking and file state analysis
- **`commit.go`** - Commit creation, staging, and commit counting
- **`push.go`** / **`pull.go`** - Remote synchronization operations
- **`fetch.go`** - Remote repository updates without merging
- **`branch.go`** - Branch management and switching operations

### Advanced Operations

- **`magic.go`** - Intelligent workflow automation
  - Analyzes repository state
  - Performs appropriate actions automatically
  - Handles complex scenarios with minimal user input

- **`history.go`** - Commit history analysis and traversal
  - Configurable history depth
  - Commit metadata extraction
  - Timeline analysis

- **`stash.go`** - Temporary change management
  - Save work-in-progress changes
  - Apply stashed changes with conflict resolution

### Utilities

- **`git.go`** - Core Git command execution and error handling
- **`result.go`** - Standardized operation result formatting
- **`createrepo.go`** - Repository initialization and setup
- **`reset.go`** / **`undo.go`** - Repository state restoration
- **`gh.go`** - GitHub-specific operations and API integration

## Usage Examples

### Basic Repository Operations

```go
// Create a new Git service for a repository
gitSvc := git.NewService("myrepo", "/path/to/repo")

// Check repository status
result, err := gitSvc.Status(ctx, logger)
if err != nil {
    return err
}

// Commit changes
result, err = gitSvc.Commit(ctx, "Add new feature", logger)
if err != nil {
    return err
}

// Push to remote
result, err = gitSvc.Push(ctx, logger)
```

### Magic Workflow Automation

```go
// Automatically handle common Git workflows
// This will analyze the repo state and perform appropriate actions
result, err := gitSvc.Magic(ctx, "Automated commit", false, logger)
if err != nil {
    return err
}

// Dry run to see what would be done
result, err = gitSvc.Magic(ctx, "Test commit", true, logger)
```

### History Analysis

```go
// Get commit history
args := &git.HistoryArgs{
    Limit: 50,
    // Additional filtering options
}
result, err := gitSvc.History(ctx, args, logger)
```

## Configuration

The Git module integrates with the telemetry system and uses the application's logging framework. No additional configuration is required beyond ensuring Git is available in the system PATH.

### Environment Integration

- Uses `telemetry.RunProcessSimple` for Git command execution
- Integrates with structured logging for operation tracking
- Respects application-level timeout and context cancellation

## Dependencies

### Required External Tools
- **Git**: Must be available in system PATH
- **GitHub CLI** (optional): For enhanced GitHub integration

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/git
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Git Documentation](https://git-scm.com/docs) - Official Git documentation
- [GitHub CLI Documentation](https://cli.github.com/manual/) - GitHub CLI reference
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
