# Accordion

You'll need to import `views/components` at the top of your template

```html
<ul class="accordion">
  <li>
    <input id="accordion-a" type="checkbox" hidden />
    <label for="accordion-a">{%= components.ExpandCollapse(3, ps) %} Option A</label>
    <div class="bd">Option A!</div>
  </li>
  <li>
    <input id="accordion-b" type="checkbox" hidden />
    <label for="accordion-b">{%= components.ExpandCollapse(3, ps) %} Option B</label>
    <div class="bd">Option B!</div>
  </li>
  <li>
    <input id="accordion-c" type="checkbox" hidden />
    <label for="accordion-c">{%= components.ExpandCollapse(3, ps) %} Option C</label>
    <div class="bd">Option C!</div>
  </li>
</ul>
```
