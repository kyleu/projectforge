# Tag Editor

An interactive drag-and-drop tag editor component that allows users to manage collections of tags or keywords. The component provides an intuitive interface for adding, removing, and reordering tags while maintaining full functionality for users with JavaScript disabled.

## Key Features

- **Drag-and-Drop Interface**: Intuitive tag management with visual feedback
- **Progressive Enhancement**: Works as a text input without JavaScript
- **Form Integration**: Seamlessly integrates with standard HTML forms
- **Comma-Separated Output**: Automatically formats tags as comma-separated values when submitted

## How It Works

The tag editor enhances a standard text input with JavaScript to provide an interactive interface. Tags are stored as a comma-separated string in the underlying input field, ensuring compatibility with standard form processing.

### Without JavaScript
Users see a regular text input where they can enter comma-separated values:
```
apple, banana, cherry, date
```

### With JavaScript
The same data is presented as individual, interactive tag elements that can be:
- Added by typing and pressing Enter
- Removed by clicking the × button
- Reordered by dragging and dropping
- Edited by double-clicking

## Basic Usage

### Using the Form Helper

The simplest way to create a tag editor is using the `FormInputTags` helper:

```html
{%= components.FormInputTags("choices", "choices", []string{"apple", "banana", "cherry"}, ps) %}
```

**Parameters:**
- `key`: string - The form field name
- `id`: string - The DOM element ID
- `values`: []string - Initial tag values
- `ps`: *cutil.PageState - Current page state
- `placeholder`: ...string - Optional placeholder text

### Manual HTML Structure

You can also create the tag editor manually for more control:

```html
<div class="tag-editor">
  <input class="result" name="choices" value="apple,banana,cherry" />
  <div class="tags"></div>
  <div class="clear"></div>
</div>
```

**Structure Elements:**
- `tag-editor`: Container div with the main CSS class
- `result`: Hidden input that stores the comma-separated tag values
- `tags`: Container where individual tag elements are displayed
- `clear`: Clearfix element for proper layout
