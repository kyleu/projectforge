# Modal Windows

This opens a "window" over the current document, blocking access to the page until the modal is dismissed by either clicking outside it or clicking the close button.

Full functionality is available without JavaScript, and care has been taken to preserve the navigation history (back/forward works). 

To create a modal, use markup like so:

```html
<div id="modal-foo" class="modal" style="display: none;">
  <a class="backdrop" href="#"></a>
  <div class="modal-content">
    <div class="modal-header">
      <a href="#" class="modal-close">×</a>
      <h2>The Modal's Title</h2>
    </div>
    <div class="modal-body">
      The Modal's Body
    </div>
  </div>
</div>
```

Then, to activate it, simply make an anchor:

```html
<a href="#modal-foo">Open modal</a>
```
