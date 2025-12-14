# Plot

The **`plot`** module provides data visualization capabilities for your application using [Observable Plot](https://observablehq.com/plot).
It enables developers to create interactive, responsive charts and graphs with minimal code.

## Overview

This module provides:

- **Observable Plot Integration**: Pre-configured Observable Plot and D3.js libraries
- **Template Components**: Ready-to-use chart templates and helper functions
- **Responsive Design**: Charts that automatically resize with the browser window
- **TypeScript Support**: Full TypeScript integration for type-safe chart development

## Key Features

### Performance
- Minified D3.js and Observable Plot libraries
- Efficient chart rendering with automatic cleanup
- Responsive charts that adapt to container size
- Progressive enhancement (graceful degradation without JavaScript)

### Developer Experience
- Pre-built chart templates for common use cases
- Simple API for creating custom visualizations
- Automatic asset loading and script management
- Integration with your application's templating system

### Chart Types
- Horizontal bar charts with customizable styling
- Extensible framework for additional chart types
- Built-in tooltips and interactions
- Responsive legends and axes

## Package Structure

### Assets
- **`assets/plot/d3.min.js`** - Minified D3.js library for data manipulation
- **`assets/plot/plot.min.js`** - Minified Observable Plot library for visualization

### Components
- **`views/components/Plot.html`** - Template functions for chart creation
  - `PlotAssets()` - Loads required JavaScript libraries
  - `PlotCall()` - Creates responsive chart containers
  - `PlotHorizontalBar()` - Pre-built horizontal bar chart template

## Usage

### Basic Setup

First, include the plot assets in your template:

```html
{%= PlotAssets() %}
```

### Creating a Horizontal Bar Chart

```html
{% code
  data := []util.ValueMap{
    {"name": "Category A", "value": 10},
    {"name": "Category B", "value": 20},
    {"name": "Category C", "value": 15},
  }
%}

<div id="chart-container"></div>
{%= PlotHorizontalBar("chart-container", data, "value", "name", "${d.name}: ${d.value}", 80) %}
```

### Custom Charts

Use the `PlotCall()` function to create custom visualizations:

```html
<div id="custom-chart"></div>
<script type="module">
  const customChart = (div) => {
    // Your Observable Plot code here
    return Plot.plot({
      marks: [
        // Your marks
      ]
    });
  };

  {%= PlotCall("custom-chart", "customChart(div)") %}
</script>
```

## Configuration

The module automatically handles:
- Asset loading and caching
- Responsive behavior on window resize
- Memory cleanup when charts are destroyed
- Integration with your application's asset pipeline

## Dependencies

- **D3.js**: Data manipulation and DOM interaction
- **Observable Plot**: High-level charting grammar
- **Core Module**: Template system and utilities

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/plot
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Observable Plot Documentation](https://observablehq.com/plot) - Complete Observable Plot reference
- [D3.js Documentation](https://d3js.org) - D3.js API reference
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
