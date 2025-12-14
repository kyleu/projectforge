# Brand Icons

The **`brands`** module provides comprehensive brand icon support for your application.
It integrates thousands of high-quality SVG brand icons from [simple-icons](https://github.com/simple-icons/simple-icons), making it easy to display logos and branding elements in your applications.

## Overview

This module offers:

- **Thousands of Brand Icons**: Access to 2,500+ SVG icons representing popular brands, services, and technologies
- **Optimized Performance**: Lightweight SVG icons with minimal payload impact
- **Easy Integration**: Simple template helpers for rendering brand icons
- **Consistent Styling**: Standardized 24x24 viewBox with customizable styling
- **Admin Gallery**: Built-in administration interface for browsing available icons

## Key Features

### Comprehensive Icon Library
- 2,500+ brand icons from simple-icons project
- Regular updates with new brand additions
- Consistent SVG format with 24x24 viewBox
- Includes popular technology, social media, and service brands

### Performance Optimized
- Lightweight SVG format
- Minimal impact on bundle size
- Efficient symbol-based rendering
- CSS class integration for styling

### Developer Experience
- Simple template integration
- Administrative gallery for icon discovery
- Type-safe icon references
- Customizable styling options

## Package Structure

### Core Components

- **`lib/icons/`** - Brand icon management and rendering
  - `Icon` struct for individual brand icons
  - `Icons` collection type for bulk operations
  - `Library` for icon discovery and management
  - SVG rendering and symbol generation utilities

### Icon Data
- **`brands.go`** - Complete collection of brand icon definitions
  - Icon metadata (title, color, aliases, guidelines, license)
  - SVG path data for each brand
  - License and usage information

### Admin Interface
- Gallery view for browsing all available brand icons
- Search and filtering capabilities
- Icon preview and metadata display

## Usage

### Basic Icon Rendering

Use the `SVGRef` component helper (and its friends, defined in `SVG.html`) to render brand icons in templates:

```html
{%= components.SVGRef("brand-apple", 64, 64, "logo-class", ps) %}
{%= components.SVGRef("brand-github", 32, 32, "social-icon", ps) %}
{%= components.SVGRef("brand-react", 48, 48, "tech-badge", ps) %}
```

### Icon Parameters
- **Icon Key**: Brand icon identifier (e.g., "brand-apple", "brand-google")
- **Width/Height**: Pixel dimensions for the rendered icon
- **CSS Class**: Custom styling class for the icon
- **Page State**: `*cutil.PageState` object

### Available Icons

Icons are prefixed with `brand-` and use the simple-icons naming convention:

- **Technology**: `brand-react`, `brand-vue`, `brand-golang`, `brand-python`
- **Social Media**: `brand-twitter`, `brand-linkedin`, `brand-instagram`
- **Services**: `brand-github`, `brand-gitlab`, `brand-aws`, `brand-docker`
- **Companies**: `brand-apple`, `brand-google`, `brand-microsoft`

### Admin Gallery

Access the brand icon gallery through the admin settings:
1. Navigate to Admin Settings
2. Select "Brand Icons" section
3. Browse, search, and preview available icons
4. Copy icon keys for use in templates

## Configuration

No additional configuration is required. The module automatically:
- Loads all brand icon definitions
- Registers admin gallery routes
- Provides template helpers for icon rendering using theme colors

## Integration

### Template Integration

```html
<!-- Header with company logo -->
<header class="header">
  {%= components.SVGRef("brand-" + ps.Data.CompanyBrand, 40, 40, "company-logo", ps) %}
  <h1>Company Name</h1>
</header>

<!-- Technology stack display -->
<div class="tech-stack">
  {%- for _, tech := range technologies -%}
    {%= components.SVGRef("brand-" + tech.Key, 32, 32, "tech-icon", ps) %}
  {%- endfor -%}
</div>
```

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/brands
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [simple-icons](https://github.com/simple-icons/simple-icons) - Brand icon source
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
