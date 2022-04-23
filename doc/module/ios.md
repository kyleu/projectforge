<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# iOS

This is a module for [Project Forge](https://projectforge.dev). It provides a webview-based application and iOS build

https://github.com/kyleu/projectforge/tree/master/module/ios

### License 

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- An iOS application template is provided in `./tools/ios`
- To function, the `iOS` build option must be enabled
- The application is based on a webview, no changes should be needed
- A script is provided at `./bin/ios.sh` that will copy and build the app
- Icons and settings can be configured in the Project Forge UI
- Once built, you can find the Xcode project at `./build/dist/mobile_ios_app_arm64`
