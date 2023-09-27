# Help

This is a module for [Project Forge](https://projectforge.dev). It provides Markdown help files that integrate into the UI

https://github.com/kyleu/projectforge/tree/master/module/help

### License 

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

Each web action provides a `key` as the first argument to `Act`. 
If that key matches a Markdown file in `./doc/help`, a help link will be rendered in the top navigation, linking to a modal containing an HTML version of the files contents. 
