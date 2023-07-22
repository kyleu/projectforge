// Content managed by Project Forge, see [projectforge.md] for details.
import {opt} from "./dom";

export function audit(msg: string, code: boolean, ...args: any) { // eslint-disable-line @typescript-eslint/no-explicit-any
  const out = msg.replace(/\{(\d+)\}/gu, (match, number) => {
    return typeof args[number] === "undefined" ? match : args[number];
  });
  const el = opt("#audit-log");
  if (el) {
    const li = document.createElement("li");
    if (code) {
      const pre = document.createElement("pre");
      pre.innerText = out;
      li.appendChild(pre);
    } else {
      li.innerText = out;
    }
    el.appendChild(li);
  } else {
    console.log("### Audit ###\n" + out);
  }
}
