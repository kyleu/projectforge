// Content managed by Project Forge, see [projectforge.md] for details.
import {els, req} from "./dom";

type Column = {
  key: string
  title: string
  description?: string
  type?: string
};

function createTableHead(cols: Column[]): HTMLElement {
  const thead = document.createElement("thead");
  const r = document.createElement("tr");
  let first = true;
  cols.forEach((col) => {
    const c = document.createElement("th");
    if (first) {
      c.classList.add("shrink");
      first = false;
    }
    c.innerText = col.title;
    r.appendChild(c);
  });
  thead.appendChild(r);
  return thead;
}

function createTableRow(cols: Column[], x: { [key: string]: unknown; }): [HTMLElement, HTMLElement[]] {
  const r = document.createElement("tr");
  const tds: HTMLElement[] = [];
  cols.forEach((col) => {
    const c = document.createElement("td");
    tds.push(c);
    const v = x[col.key];
    if (v === undefined || v === null) {
      const em = document.createElement("em");
      em.innerText = "-";
      c.appendChild(em);
    } else if (col.type === "code") {
      const pre = document.createElement("pre");
      pre.innerText = JSON.stringify(v, null, 2);
      c.appendChild(pre);
    } else {
      c.innerText = v.toString();
    }
    r.appendChild(c);
  });
  return [r, tds];
}

function createTable(cols: Column[], rows: { [key: string]: unknown; }[]): HTMLElement {
  const tbl = document.createElement("table");
  const allTds: HTMLElement[][] = [];
  tbl.classList.add("min-200");
  tbl.appendChild(createTableHead(cols));
  const tbody = document.createElement("tbody");
  rows.forEach((row) => {
    const [tr, tds] = createTableRow(cols, row);
    tbody.appendChild(tr);
    allTds.push(tds);
  });
  tbl.appendChild(tbody);
  return tbl;
}

export function createEditor(el: HTMLElement): void {
  const key = el.dataset.key ?? "editor";
  const columnsStr = el.dataset.columns ?? "[]";
  const columns: Column[] = JSON.parse(columnsStr.replace(/\\"/gu, "\""));

  const inp: HTMLTextAreaElement = req<HTMLTextAreaElement>("textarea", el);
  let curr: Column[] = JSON.parse(inp.value);
  if (curr === undefined || curr === null) {
    curr = [];
  }
  if (!Array.isArray(curr)) {
    throw new Error("input value for element [" + key + "] of type [" + typeof curr + "] must be an array");
  }

  const tbl = createTable(columns, curr);

  els(".toggle-editor-" + key).forEach((toggle) => {
    toggle.innerText = "Edit";
    toggle.onclick = () => {
      if (toggle.innerText === "Edit") {
        toggle.innerText = "View";
        tbl.hidden = true;
        inp.hidden = false;
      } else {
        toggle.innerText = "Edit";
        tbl.hidden = false;
        inp.hidden = true;
      }
    };
  });

  el.appendChild(tbl);
  inp.hidden = true;
}

export function editorInit() {
  els(".rich-editor").forEach((x) => {
    createEditor(x);
  });
}
