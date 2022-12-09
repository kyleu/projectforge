# Tag Editor

A drag-and-drop tag editor is provided. 
It allows you to manipulate an array of strings, represented as a comma-separated list when submitted.
Users with scripting disabled will see a normal text input.

```html
<!-- FormInputTags(key string, id string, values []string, ps *cutil.PageState, placeholder ...string) -->
{%= components.FormInputTags("choices", "choices", []string{"a", "b", "c"}, ps) %}
```

Or, create one manually:
```html
<div class="tag-editor">
  <input class="result" name="choices" value="a,b,c" />
  <div class="tags"></div>
  <div class="clear"></div>
</div>
```
