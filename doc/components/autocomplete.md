# Autocomplete

Server-driven autocomplete using the browser's native `datalist` functionality.

To create an autocomplete, create a normal form input, and call the `autocomplete` method. 
It takes several arguments, including the URL, the querystring field name, and transformations methods for creating the `datalist`.

The URL must return a JSON array of objects.

The result of the `title` function must include the search content or the browser will hide it.

```html
<input id="input-user" />

<script>
  document.addEventListener("DOMContentLoaded", function() {
    const input = document.getElementById("input-user");
    
    // returns the text description of the object
    const title = function(o) {
      return o["name"] + " (" + o["id"] + ")";
    }
    
    // returns the actual value that will be submitted
    const val = function(o) {
      return o["id"];
    }

    // autocomplete(el: HTMLInputElement, url: string, field: string, title: (x: any) => string, val: (x: any) => string)
    {your project}.autocomplete(input, "/user", "q", title, val);
  });
</script>
```

If you use the `export` module in your app, autocompletes will be generated automatically for model relationships.  
