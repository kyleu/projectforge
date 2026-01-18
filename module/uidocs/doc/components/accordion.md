# Accordion

An accordion component that allows users to expand and collapse content sections. This component works entirely without JavaScript, using CSS and HTML form elements to provide smooth animations and interactions.

## Overview

The accordion component is perfect for organizing content into collapsible sections, making it ideal for FAQs, documentation sections, or any interface where you want to save vertical space while keeping content accessible. Each section can be independently expanded or collapsed by clicking on its header.

## Key Features

- **No JavaScript Required**: Uses pure CSS and HTML form elements for functionality
- **Smooth Animations**: Built-in CSS transitions for expand/collapse actions
- **Accessible**: Proper labeling and keyboard navigation support
- **Flexible Content**: Each section can contain any HTML content
- **Optional Animation**: Choose between animated and non-animated sections

## Prerequisites

You'll need to import `views/components` at the top of your template to use the `ExpandCollapse` helper component.

## Basic Usage

```html
<ul class="accordion">
  <li>
    <input id="accordion-a" type="checkbox" hidden />
    <label for="accordion-a">{%= components.ExpandCollapse(3, ps) %} Option A</label>
    <div class="bd"><div><div>
      Option A content goes here. This can include any HTML elements,
      text, images, or other components.
    </div></div></div>
  </li>
  <li>
    <input id="accordion-b" type="checkbox" hidden />
    <label for="accordion-b">{%= components.ExpandCollapse(3, ps) %} Option B</label>
    <div class="bd"><div><div>
      Option B content goes here. Each section is independent
      and can be expanded or collapsed separately.
    </div></div></div>
  </li>
  <li>
    <input id="accordion-c" type="checkbox" hidden />
    <label for="accordion-c">{%= components.ExpandCollapse(3, ps) %} Option C (not animated)</label>
    <div class="bd-no-animation">
      Option C content without animation. Use this class when you want
      instant show/hide without the smooth transition effect.
    </div>
  </li>
</ul>
```

## Structure Breakdown

### Container
- `<ul class="accordion">`: The main container that holds all accordion sections

### Individual Sections
Each accordion section consists of three parts:

1. **Hidden Checkbox**: `<input type="checkbox" hidden />` - Controls the expand/collapse state
2. **Clickable Label**: `<label>` - The header that users click to toggle the section
3. **Content Area**: `<div class="bd">` or `<div class="bd-no-animation">` - Contains the collapsible content

### Important Notes

- **Unique IDs**: Each checkbox must have a unique `id` attribute
- **Matching Labels**: The `for` attribute on each label must match its corresponding checkbox `id`
- **Triple Div Structure**: For animated sections, use the nested `<div><div><div>` structure for proper animation
- **ExpandCollapse Helper**: The `{%= components.ExpandCollapse(3, ps) %}` adds the expand/collapse icon

## Animation Options

### Animated Sections (Default)
Use `class="bd"` with the triple-div structure for smooth expand/collapse animations:

```html
<div class="bd"><div><div>
  Your content here
</div></div></div>
```

### Non-Animated Sections
Use `class="bd-no-animation"` for instant show/hide without transitions:

```html
<div class="bd-no-animation">
  Your content here
</div>
```
