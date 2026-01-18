# Theme Catalog

The **`themecatalog`** module provides comprehensive theme management and customization capabilities for your application.
It includes 24 built-in themes plus palette-driven generation and custom theme creation tools.

## Overview

This module extends the core theming system by providing:

- **Built-in Theme Collection**: 24 professionally designed themes with light/dark variants
- **Custom Theme Creation**: Generate themes from any hex color using advanced palette algorithms
- **Theme Management UI**: Intuitive web interface for theme selection and customization
- **Real-time Preview**: Live theme preview and editing capabilities

## Key Features

### Built-in Themes
- **24 pre-designed themes** including all your favorite colors
- **Dual color schemes** - each theme includes both light and dark variants
- **Professionally balanced palettes** optimized for readability and accessibility

### Theme Generation
- **Palette-based generation** using the `github.com/muesli/gamut` palettes (Crayola, CSS, RAL, Resene, Wikipedia)
- **Custom color creation** from any hex color input
- **Intelligent color harmony** automatically generates complementary colors
- **Accessibility considerations** maintains proper contrast ratios

### Management Interface
- **Interactive theme browser** at `/theme` endpoint
- **Real-time preview** see changes instantly
- **Theme editing tools** modify and customize existing themes
- **Export capabilities** save custom themes for reuse

## Package Structure

### Core Components

- **`app/lib/theme/`** - Theme generation and management logic
  - Color palette algorithms
  - Theme validation and processing
  - Gamut integration for advanced color theory

- **`app/controller/clib/themecatalog.go`** - HTTP handlers for theme operations
  - Theme selection and switching
  - Custom theme creation endpoints
  - Theme preview and editing

- **`views/vtheme/`** - Theme management UI templates
  - Theme selection interface
  - Color picker and customization forms
  - Preview and comparison views

## Usage

### Accessing the Theme Catalog

Navigate to `/theme` in your application to access the theme management interface.
Use `/theme/palette/{palette}` to browse palette themes (`crayola`, `css`, `ral`, `resene`, `wikipedia`).

### Creating Custom Themes

1. **From Color**: Enter any hex color to generate a full theme palette (via the color picker on `/theme`)
2. **From Existing**: Modify any built-in theme as a starting point
3. **Manual Creation**: Use the color picker to create completely custom palettes

### Applying Themes

Themes can be applied:
- **Per-user**: Individual theme preferences
- **Application-wide**: Default theme for all users
- **Dynamic switching**: Runtime theme changes without restart

## Configuration

The module integrates with the core theme system configuration:

- `app_nav_color_light` - Navigation color for light themes
- `app_nav_color_dark` - Navigation color for dark themes

## Dependencies

- **[Gamut](https://github.com/muesli/gamut)** - Color palette generation and named palettes
- **Core Module** - Required for base theming infrastructure

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/themecatalog
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
