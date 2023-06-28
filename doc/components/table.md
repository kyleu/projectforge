# Tables

Helpers are provided for emitting HTML tables. 

You can use the methods in `TableHeader.html` to create `<th>` entries that are styled, resizable, and work without JavaScript.

```go
TableHeader(section, key, title, params, icon, uri, tooltip, sortable, cls, resizable, ps)
```

You can use the methods in `Table.html` to create table rows containing a title, help text, and a form input, all wired together. Look through the file, there are helpers for almost all input types.

```go
TableInput(key, id, title, value, indent, "help text")
```
