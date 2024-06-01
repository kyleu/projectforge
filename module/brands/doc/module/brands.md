# Brand Icons

This is a module for [Project Forge](https://projectforge.dev). It provides thousands of SVG icons from [simple-icons](https://github.com/simple-icons/simple-icons) for representing common logos

https://github.com/kyleu/projectforge/tree/master/module/brands

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- Brand icons are available in `lib/icons`, in a very large source file
- There's a gallery available in the admin settings section
- To use a brand icon, simply reference an SVG starting with the prefix `brand-` in your template

```html
{%= components.SVGRef("brand-apple", 64, 64, "", ps) %}
```
