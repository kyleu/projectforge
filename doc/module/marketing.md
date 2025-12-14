# Marketing Site

The **`marketing`** module provides a comprehensive marketing website infrastructure for your application.
It creates a separate marketing site that runs alongside your main application, perfect for showcasing features, providing downloads, and engaging users.

## Overview

This module provides:

- **Separate Marketing Site**: Runs on a dedicated port (main port + 1)
- **Download Management**: GitHub release integration with detailed download pages
- **Documentation System**: Markdown-based content management
- **Marketing Pages**: Professional landing pages and feature showcases
- **SEO Optimization**: Search engine friendly structure and metadata

## Key Features

### Multi-Port Architecture
- Marketing site runs independently on port offset +1
- Clean separation between application and marketing concerns
- Can be deployed separately or together
- Zero impact on main application performance

### Content Management
- Markdown-based documentation system
- Dynamic page generation from `/doc` directory
- File-based content that's version controlled
- Easy content updates without code changes

### Download Infrastructure
- GitHub releases integration
- Automatic asset discovery and listing
- Platform-specific download links
- Release notes and changelog display
- Download statistics and tracking

### Professional Presentation
- Clean, modern marketing page layouts
- Responsive design for all devices
- Fast loading times and optimized assets
- Progressive enhancement compatibility

## Configuration

### Startup Options
Start the marketing site alongside your application:

```bash
# Start with marketing site
./app all
./app site

# Marketing site only
./app marketing
```

### Content Management
- Place marketing content in `./doc/` directory
- Supports standard Markdown with frontmatter
- Automatic menu generation from directory structure
- Images and assets served from `./assets/marketing/`

### GitHub Integration
Configure GitHub repository for downloads:
- Automatic detection from git remote
- Release asset discovery and organization
- Platform-specific download recommendations
- Release notes integration

## Key Use Cases

### Software Distribution
- Professional download pages for releases
- Platform-specific installation instructions
- Release notes and changelog display
- Binary verification and checksums

### Documentation Hub
- User guides and tutorials
- API documentation
- Getting started guides
- Feature demonstrations

### Marketing Website
- Product feature showcases
- Comparison pages
- Pricing information
- Contact and support pages

### Developer Resources
- Code examples and snippets
- Integration guides
- Best practices documentation
- Community resources

## Dependencies

**Required Modules:**
- `filesystem` - File operations and content management

**Compatible With:**
- `database` - For download analytics
- `user` - For user-specific download tracking
- `admin` - For content management interface

## Performance

### Optimizations
- Static asset caching and compression
- Markdown compilation caching
- Minimal JavaScript requirements
- Optimized image delivery

### Metrics
- Fast page load times (<100ms for cached content)
- Lightweight asset footprint
- SEO-optimized markup and structure
- Mobile-first responsive design

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/marketing
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Filesystem Module](filesystem.md) - Required dependency for content management
- [Configuration Guide](../running.md) - Environment variables and startup options
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
