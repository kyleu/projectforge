# Upgrade

The **`upgrade`** module provides automated in-place version upgrades for your application using GitHub Releases.
It enables applications to update themselves to the latest version without manual intervention.

## Overview

This module adds self-upgrade capabilities to your application by:

- **GitHub Integration**: Automatically checks GitHub Releases for newer versions
- **In-place Updates**: Downloads and replaces the current binary seamlessly
- **Cross-platform Support**: Works across all supported operating systems and architectures
- **Safe Upgrades**: Validates downloads and provides rollback mechanisms

## Key Features

### Automated Updates
- Checks GitHub Releases API for latest versions
- Compares current version with available releases
- Downloads appropriate binary for current platform/architecture

### Safety & Reliability
- Validates downloaded binaries before replacement
- Creates backup of current binary
- Provides rollback capability if upgrade fails
- Graceful error handling and status reporting

### User Experience
- Simple command-line interface
- Progress indicators during download
- Clear success/failure messaging
- Non-disruptive upgrade process

## Usage

### Basic Upgrade Command

```bash
# Upgrade to the latest release
./your-app upgrade

# Check for available upgrades without installing
./your-app upgrade --check

# Upgrade to a specific version
./your-app upgrade --version v1.2.3
```

### Command Options

- `--check` - Check for updates without installing
- `--version` - Upgrade to a specific version tag
- `--force` - Force upgrade even if current version is newer
- `--dry-run` - Show what would be upgraded without making changes

## Package Structure

### Core Components

- **`cmd/upgrade.go`** - CLI command implementation
  - Version checking and comparison
  - GitHub API integration
  - Download and installation logic

- **`controller/upgrade.go`** - Web interface for upgrades
  - Admin panel upgrade interface
  - Progress tracking and status display
  - Browser-based upgrade management

### Libraries

- **`lib/upgrade/`** - Core upgrade functionality
  - GitHub Release API client
  - Binary download and validation
  - Platform detection and asset selection

## Troubleshooting

### Common Issues

**Network/GitHub API Issues:**
- Set `GITHUB_TOKEN` for higher rate limits
- Check firewall settings for GitHub API access
- Verify repository configuration

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/upgrade
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
- [Customizing Guide](../customizing.md) - Advanced customization options
- [GitHub Releases API](https://docs.github.com/en/rest/releases) - GitHub API documentation
