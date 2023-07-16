# Notarize

This is a module for [Project Forge](https://projectforge.dev). It sends files to Apple for notarization.

https://github.com/kyleu/projectforge/tree/master/module/notarize

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- Notarizes built artifacts using Apple's services
- To function, your project's `SigningIdentity` must be set and the `Notarize` build option must be enabled
- You'll need to install [gon](https://github.com/mitchellh/gon) and, if you're on Apple Silicon, patch the binary
- It takes forever and sends you multiple emails, so use cautiously
