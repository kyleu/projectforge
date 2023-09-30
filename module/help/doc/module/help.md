# Help

This is a module for [Project Forge](https://projectforge.dev). It provides Markdown help files that integrate into the UI

https://github.com/kyleu/projectforge/tree/master/module/help

### License 

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

Each web action provides a `key` as the first argument to `Act`. 
If that key matches a Markdown file in `./doc/help`, a help link will be rendered in the top navigation, linking to a modal containing an HTML version of the files contents.

As an example, the home page uses an key of `home`, so the app will render the help content from `help/home.md`.
If the Markdown file contains a header as its first line, the header content will be used as the title. 
Otherwise, the page's title will be used.

Files named `_header.md` and `_footer.md` will be included in every help page, in the places you'd expect.
