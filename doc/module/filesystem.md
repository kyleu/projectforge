<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# Filesystem

This is a module for [Project Forge](https://projectforge.dev). It provides an abstraction around local and remote filesystems

https://github.com/kyleu/projectforge/tree/master/module/filesystem

### License 

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- The interface `FileLoader` is provided, with a single implementation `FileSystem`
- The `Services` interface contains a default config filesystem
- You can create new filesystems with `NewFileSystem`
