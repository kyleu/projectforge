<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# Process

This is a module for [Project Forge](https://projectforge.dev). It provides a framework for managing system processes.

https://github.com/kyleu/projectforge/tree/master/module/process

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

Create a new exec service by calling `app/lib/exec/NewService()`, then call `NewExec`. 

If you'd like to monitor the progress, render `/views/vexec/Detail.html`. If the `websocket` module is included, you'll be able to see a streaming result. 
