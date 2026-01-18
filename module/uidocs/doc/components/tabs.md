# Tabs

A tabbed navigation component that allows users to switch between different content panels. This component works entirely without JavaScript, using CSS and HTML radio buttons to manage the active tab state.

## How It Works

The tab system uses HTML radio buttons with the same `name` attribute to ensure only one tab can be selected at a time. CSS `:checked` pseudo-selectors show the corresponding content panel when a tab is selected.

## Basic Usage

```html
<div class="tabs">
  <input name="tab-group" type="radio" id="tab-alpha" class="input" checked/>
  <label for="tab-alpha" class="label">Alpha</label>
  <div class="panel">
    <p>This is the Alpha tab content. You can put any HTML content here.</p>
  </div>

  <input name="tab-group" type="radio" id="tab-beta" class="input"/>
  <label for="tab-beta" class="label">Beta</label>
  <div class="panel">
    <p>This is the Beta tab content. Each tab can have completely different content.</p>
  </div>
</div>
```

## Dynamic Tab Generation

For dynamic content, you can use template loops to generate tabs:

```html
<div class="tabs">
  {%- for i, option := range []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"} -%}
  <input name="tab-group" type="radio" id="tab-{%s option %}" class="input" {%- if i == 0 %} checked{%- endif %}/>
  <label for="tab-{%s option %}" class="label">{%s option %}</label>
  <div class="panel">
    <h3>{%s option %} Content</h3>
    <p>This is the content for the {%s option %} tab. Each tab can contain unique content tailored to its purpose.</p>
  </div>
  {%- endfor -%}
</div>
```

### Container
- `<div class="tabs">`: The main container that holds all tab elements

### Tab Pattern
Each tab consists of three elements that must appear in this order:

1. **Radio Input**: `<input name="tab-group" type="radio" id="tab-id" class="input"/>`
   - Controls which tab is active
   - All tabs in a group must share the same `name` attribute
   - Each tab needs a unique `id`

2. **Tab Label**: `<label for="tab-id" class="label">Tab Title</label>`
   - The clickable tab header
   - `for` attribute must match the radio input's `id`

3. **Content Panel**: `<div class="panel">Content here</div>`
   - Contains the tab's content
   - Shown when the corresponding radio button is checked

### Important Notes

- **Consistent Naming**: All radio inputs in a tab group must have the same `name` attribute
- **Unique IDs**: Each radio input must have a unique `id`
- **Matching Labels**: Each label's `for` attribute must match its radio input's `id`
- **Order Matters**: The input, label, and panel must appear in that exact order for CSS to work correctly
