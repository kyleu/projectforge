# DOM Utilities

Several utilities are provided to make it more convenient to work with the dom.

```typescript
import {els, opt, req, setDisplay} from "./dom";

// els<T extends HTMLElement>(selector: string, context?: HTMLElement): readonly T[]
const list1 = els("button"); // find all buttons
const list2 = els("button", someElement); // find all buttons within provided parent

// opt<T extends HTMLElement>(selector: string, context?: HTMLElement): T | undefined
const elOption = opt<HTMLButtonElement>("#my-button"); // returns undefined if not present 

// req<T extends HTMLElement>(selector: string, context?: HTMLElement): T
const el = req<HTMLDivElement>("#my-div");

setHTML(el, "<em>Hello!</em>");
setDisplay(el, false);
setText(el, "Hello!");
clear(el);
```
