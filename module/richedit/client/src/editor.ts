import {els, req} from "./dom";
import {modalGetBody, modalGetOrCreate, modalSetTitle} from "./modal";

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

function rowEdit(idx: number, cols: Column[], x: { [p: string]: unknown }) {
  return () => {
    const objectEditor = document.createElement("pre");
    objectEditor.innerText = JSON.stringify(x, null, 2);

    const modal = modalGetOrCreate("rich-editor", "Rich Editor");
    modalSetTitle(modal, "Editing row [" + idx + "]");
    modalGetBody(modal).replaceChildren(objectEditor);
    window.location.href = "#modal-rich-editor";

    console.log("!", idx, cols, x);
  }
}

function createTableCell(col: Column, v: unknown): HTMLElement {
  const c = document.createElement("td");
  if (v === undefined || v === null) {
    const em = document.createElement("em");
    em.innerText = "-";
    c.appendChild(em);
  } else if (col.type === "code" || col.type === "json") {
    const pre = document.createElement("pre");
    pre.innerText = JSON.stringify(v, null, 2);
    c.appendChild(pre);
  } else if (col.type === "type") {
    if (typeof v === "string") {
      c.innerText = v.toString();
    } else {
      if (typeof v === "object") {
        c.innerText = (v as any)["k"].toString();
      } else {
        const pre = document.createElement("pre");
        pre.innerText = JSON.stringify(v, null, 2);
        c.appendChild(pre);
      }
    }
  } else {
    c.innerText = v.toString();
  }
  return c;
}

function createTableRow(idx: number, cols: Column[], x: { [key: string]: unknown; }): HTMLElement {
  const r = document.createElement("tr");
  r.dataset.index = idx.toString();
  cols.forEach((col) => {
    const c = createTableCell(col, x[col.key]);
    r.appendChild(c);
  });
  r.onclick = rowEdit(idx, cols, x);
  return r;
}

function createTable(cols: Column[], rows: { [key: string]: unknown; }[]): HTMLElement {
  const tbl = document.createElement("table");
  tbl.classList.add("min-200");
  tbl.appendChild(createTableHead(cols));
  const tbody = document.createElement("tbody");
  rows.forEach((row, idx) => {
    const tr = createTableRow(idx, cols, row);
    tbody.appendChild(tr);
  });
  tbl.appendChild(tbody);

  const div = document.createElement("div");
  div.classList.add("overflow");
  div.classList.add("full-width");
  div.appendChild(tbl);
  return div;
}

const rawLabel = "Raw JSON";

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

  let tbl = createTable(columns, curr);

  els(".toggle-editor-" + key).forEach((toggle) => {
    if (toggle.innerText === "") {
      toggle.innerText = "Editor";
    }
    const editorLabel = toggle.innerText;
    toggle.innerText = rawLabel;
    toggle.onclick = () => {
      if (toggle.innerText === rawLabel) {
        toggle.innerText = editorLabel;
        tbl.remove();
        inp.hidden = false;
      } else {
        toggle.innerText = rawLabel;
        tbl.remove();
        curr = JSON.parse(inp.value);
        if (curr === undefined || curr === null) {
          curr = [];
        }
        tbl = createTable(columns, curr);
        el.appendChild(tbl);
        inp.hidden = true;
      }
    };
    toggle.style.display = "inline";
  });

  el.appendChild(tbl);
  inp.hidden = true;
}

export function editorInit() {
  els(".rich-editor").forEach((x) => {
    createEditor(x);
  });
}
