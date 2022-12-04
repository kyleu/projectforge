<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# Process

This is a module for [Project Forge](https://projectforge.dev). It provides a framework for managing system processes.

https://github.com/kyleu/projectforge/tree/master/module/process

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

The `audit` module depends on `database`, and provides tables for tracking changes to other objects. 

Best paired with the `export` module, you can use `audit` to save audit logs. See `./app/lib/audit` for details
