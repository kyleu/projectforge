# [audit]

This is a module for [Project Forge](https://projectforge.dev). It provides an audit framework for tracking changes.

https://github.com/kyleu/projectforge/tree/master/module/audit

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

The `audit` module depends on `database`, and provides tables for tracking changes to other objects. 

Best paired with the `export` module, you can use `audit` to save audit logs. See `./app/lib/audit` for details
