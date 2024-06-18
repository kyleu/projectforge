import {opt} from "./dom";

function renderAudit(msg: string, ...codes: unknown[]) { // eslint-disable-line @typescript-eslint/no-explicit-any
  const li = document.createElement("li");
  li.innerText = msg;
  for (const code of codes) {
    const pre = document.createElement("pre");
    if (typeof code === "string") {
      pre.innerText = code;
    } else {
      pre.innerText = JSON.stringify(code, null, 2);
    }
    li.appendChild(pre);
  }
  return li;
}

export function audit(msg: string, ...codes: any) { // eslint-disable-line @typescript-eslint/no-explicit-any
  const el = opt("#audit-log");
  if (el) {
    el.appendChild(renderAudit(msg, ...codes));
  } else {
    console.log("### Audit ###\n" + msg, ...codes);
  }
}
