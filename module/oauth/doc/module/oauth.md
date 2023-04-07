# OAuth

This is a module for [Project Forge](https://projectforge.dev). It provides logins and session management for many OAuth providers.

https://github.com/kyleu/projectforge/tree/master/module/oauth

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- Provides OAuth sessions for the web UI
- By default, your project doesn't require any permissions or OAuth
- To enable OAuth, set environment variables as described in the setup page at `/admin/settings`
- To enable permissions, add the following code to your `appInit` function:

```go
// The first parameter indicates if all actions should be allowed by default
user.SetPermissions(false,
    user.Perm("/admin", "github:projectforge.dev", true), // grant admin access to users signed into GitHub with an email domain ending in projectforge.dev
    user.Perm("/admin", "*", false), // deny other users access to admin
    user.Perm("/", "*", true), // allow all other traffic
)
```
