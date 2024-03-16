<!-- $PF_GENERATE_ONCE$ -->
---
toc: false
---

# {{{ .Name }}} Notebook

This is your notebook, use it for whatever you'd like!

```js echo
const go = await FileAttachment("./data/GoExample.json").json().catch(e => [["Go Error", e]]);
const js = await FileAttachment("./data/TypeScriptExample.json").json().catch(e => [["TypeScript Error", e]]);
const py = await FileAttachment("./data/PythonExample.json").json().catch(e => [["Python Error", e]]);

const data = [].concat(go).concat(js).concat(py);
```

${Inputs.table(data, {"columns": ["0", "1"], "header": {"0": "Language", "1": "Current Time"}})}

