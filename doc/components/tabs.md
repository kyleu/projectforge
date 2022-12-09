# Tabs

It's tabbed navigation. No JavaScript required.

```html
<div class="tabs">
  {%- for _, o := range []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"} -%}
  <input name="type" type="radio" id="tab-{%s o %}" class="input"/>
  <label for="tab-{%s o %}" class="label">{%s o %}</label>
  <div class="panel">
    <p>This is a tab named {%s o %}</p>
  </div>
  {%- endfor -%}
</div>
```

