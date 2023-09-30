# Accordion

It's an accordion. No JavaScript required. You'll need to import `views/components` at the top of your template.

```html
<ul class="accordion">
  <li>
    <input id="accordion-a" type="checkbox" hidden />
    <label for="accordion-a">{%= components.ExpandCollapse(3, ps) %} Option A</label>
    <div class="bd"><div><div>
      Option A!
    </div></div></div>
  </li>
  <li>
    <input id="accordion-b" type="checkbox" hidden />
    <label for="accordion-b">{%= components.ExpandCollapse(3, ps) %} Option B</label>
    <div class="bd"><div><div>
      Option B!
    </div></div></div>
  </li>
  <li>
    <input id="accordion-b" type="checkbox" hidden />
    <label for="accordion-b">{%= components.ExpandCollapse(3, ps) %} Option C (not animated)</label>
    <div class="bd-no-animation">Option C!</div>
  </li>
</ul>
```
