import { Socket, SocketMessage } from "./socket";
import { req } from "./dom";

export function socketLog(
  debug: boolean,
  parentEl: HTMLElement,
  terminal: boolean,
  url: string,
  extraHandlers: ((m: SocketMessage) => void)[]
) {
  const o = () => {
    if (debug) {
      console.log("[socket]: open");
    }
  };

  const newRow = () => {
    if (terminal) {
      const row = document.createElement("tr");
      const numTH = document.createElement("th");
      numTH.innerText = parentEl.children.length.toString();
      const textTD = document.createElement("td");
      row.append(numTH, textTD);
      parentEl.append(row);
      return textTD;
    }
    const div = document.createElement("div");
    parentEl.append(div);
    return div;
  };

  let currRow: HTMLElement | null = null;
  const appendText = (row: HTMLElement, text: string) => {
    row.appendChild(document.createTextNode(text));
  };
  const r = (m: SocketMessage) => {
    if (m.cmd !== "output") {
      if (extraHandlers.length === 0) {
        console.log("unknown command [" + m.cmd + "] received");
      } else {
        for (const h of extraHandlers) {
          h(m);
        }
      }
      return;
    }
    const htmlParam = m.param.html;
    if (htmlParam === undefined) {
      console.log("no [html] key in message param: " + JSON.stringify(m, null, 2));
    }
    let html: string[] = [];
    if (Array.isArray(htmlParam)) {
      html = htmlParam.filter((item): item is string => typeof item === "string");
    } else if (typeof htmlParam === "string") {
      html = [htmlParam];
    }

    let content = "";
    for (const x of html) {
      currRow ??= newRow();
      if (x === "\n") {
        if (content === "") {
          content = "\u00a0";
        }
        appendText(currRow, content);
        content = "";
        currRow = null;
      } else {
        content += x;
      }
    }
    if (currRow) {
      appendText(currRow, content);
    }
    const c = req("#content");
    c.scrollTo(0, c.scrollHeight);
  };

  const e = (svc: string, err: string) => {
    console.error("socket error", svc, err);
  };

  return new Socket(false, o, r, e, url);
}
