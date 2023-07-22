// Content managed by Project Forge, see [projectforge.md] for details.
import {opt} from "./dom";

export function audit(msg: string, code: boolean, ...args: any) { // eslint-disable-line @typescript-eslint/no-explicit-any
  const out = msg.replace(/{(\d+)}/g, function (match, number) {
    return typeof args[number] != 'undefined' ? args[number] : match;
  });
  const el = opt("#audit-log");
  if (!el) {
    console.log("### Audit ###\n" + out);
  } else {
    const li = document.createElement("li");
    if (code) {
      const pre = document.createElement("pre");
      pre.innerText = out;
      li.appendChild(pre);
    } else {
      li.innerText = out;
    }
    el.appendChild(li);
  }
}
