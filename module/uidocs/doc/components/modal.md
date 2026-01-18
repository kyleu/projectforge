# Modal Windows

Modal windows create overlay dialogs that appear on top of the main content, requiring user interaction before they can continue. This component provides full functionality without JavaScript while preserving browser navigation history.

## Key Features

- **No JavaScript Required**: Pure CSS and HTML implementation
- **Navigation History Preserved**: Back/forward browser buttons work correctly
- **Accessible**: Proper focus management and keyboard navigation
- **Flexible Content**: Support for any HTML content including forms and media
- **Multiple Dismiss Options**: Close via backdrop click, close button, pressing escape, or browser back button
- **Responsive**: Adapts to different screen sizes

## How It Works

The modal system uses CSS `:target` pseudo-selector to show/hide modals based on URL fragments. When a user clicks a link with a hash (`#modal-id`), the browser navigates to that fragment, making the modal visible. Clicking the backdrop or close button navigates back to `#`, hiding the modal.

## Basic Usage

### Step 1: Create the Modal HTML

```html
<div id="modal-example" class="modal" style="display: none;">
  <a class="backdrop" href="#"></a>
  <div class="modal-content">
    <div class="modal-header">
      <a href="#" class="modal-close">×</a>
      <h2>Modal Title</h2>
    </div>
    <div class="modal-body">
      <p>Your modal content goes here. This can include text, forms, images, or any other HTML elements.</p>
    </div>
    <!-- optional footer -->
    <div class="modal-footer">
      <a href="#" class="btn btn-secondary">Cancel</a>
      <a href="/delete/123" class="btn btn-danger">Delete</a>
    </div>
  </div>
</div>
```

### Step 2: Create a Trigger Link

```html
<a href="#modal-example">Open Modal</a>
```

## Structure Breakdown

### Modal Container
- `<div id="modal-example" class="modal" style="display: none;">`: The main modal container
  - **ID**: Must be unique and match the trigger link's hash
  - **Class**: `modal` applies the modal styling
  - **Style**: `display: none;` hides the modal by default

### Backdrop
- `<a class="backdrop" href="#"></a>`: Clickable area behind the modal content
  - Clicking this closes the modal by navigating to `#`
  - Covers the entire viewport when modal is open

### Modal Content
- `<div class="modal-content">`: Contains the actual modal content
  - This is the white/styled box that appears in the center

### Modal Header
- `<div class="modal-header">`: Contains the title and close button
  - `<a href="#" class="modal-close">×</a>`: Close button (typically an X)
  - `<h2>`: Modal title
