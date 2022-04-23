# Database

This is a module for [Project Forge](https://projectforge.dev). It provides an API for accessing relational databases.

https://github.com/kyleu/projectforge/tree/master/module/database

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage
- SQL files in `queries` will be compiled with quicktemplate
- Package `app/lib/database` provides many utility classes and services
- To use this module, include one of the engine-specific modules like `postgres`, `mysql`, or `sqlite`
