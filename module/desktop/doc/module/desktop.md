# Desktop

This is a module for [Project Forge](https://projectforge.dev). It provides a desktop application using the system's webview

https://github.com/kyleu/projectforge/tree/master/module/desktop

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage
- Desktop applications are provided for Linux, macOS, and Windows
- To function, the `Desktop` build option must be enabled
- The application is based on a webview, no changes should be needed
- The server application is bundled as a library and automatically started
- Because of cross-compilation issues, building is done in a Docker image
- To build the desktop apps, run `./bin/build/desktop.sh`
