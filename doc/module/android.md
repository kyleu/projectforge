<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# Android

This is a module for [Project Forge](https://projectforge.dev). It provides a webview-based application and Android build

https://github.com/kyleu/projectforge/tree/master/module/android

### License 

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage
- An Android application template is provided in `./tools/android`
- To function, the `android` build option must be enabled
- The application is based on a webview, no changes should be needed
- A script is provided at `./bin/android.sh` that will copy and build the app
- Icons and settings can be configured in the Project Forge UI
- Once built, you can find the Android Studio project at `./build/dist/mobile_android_app_arm64`
