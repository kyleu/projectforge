# Process

This is a module for [Project Forge](https://projectforge.dev). It allows the execution of JavaScript files using a built-in interpreter.

https://github.com/kyleu/projectforge/tree/main/module/scripting

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

Create a new `Script` service by calling `scripting.NewService(filesystem.NewService("./data"), "scripts")`.

A UI is provided for ad-hoc scripting and filesystem management.

Your scripts can expose test case examples that will automatically be run:

```javascript
function test(name, t) {
  return `Hello [${name}] from [${t}] script`;
}

const examples = {
  "test": [["a", "x"], ["b", "y"], ["c", "z"]]
};
```
