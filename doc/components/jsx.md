# JSX Templating

A very tiny custom JSX engine is provided. To use it, create a `.tsx` file in `./client/src`, like so:

```typescript
import * as JSX from "./jsx";

export function myComponent() {
  return <div style="color: red;">Check out this JSX!</div>;
}
```
