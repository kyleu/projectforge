import { rowEditHandler } from "./editorobject";
import { unknownToString } from "./util";
import { type Column, type Editor, type Type, typeToString } from "./editortypes";

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
      c.innerText = v;
    } else {
      c.innerText = typeToString(v as Type);
    }
  } else {
    c.innerText = unknownToString(v);
  }
  return c;
}

function createTableRow(e: Editor, idx: number, x: Record<string, unknown>): HTMLElement {
  const r = document.createElement("tr");
  r.dataset.index = idx.toString();
  e.columns.forEach((col) => {
    const c = createTableCell(col, x[col.key]);
    r.appendChild(c);
  });
  r.onclick = rowEditHandler(idx, e, x);
  return r;
}

export function createTable(e: Editor): HTMLElement {
  const tbl = document.createElement("table");
  e.table = tbl;
  tbl.classList.add("min-200");
  tbl.appendChild(createTableHead(e.columns));
  const tbody = document.createElement("tbody");
  e.rows.forEach((row, idx) => {
    const tr = createTableRow(e, idx, row);
    tbody.appendChild(tr);
  });
  tbl.appendChild(tbody);

  const div = document.createElement("div");
  div.classList.add("overflow");
  div.classList.add("full-width");
  div.appendChild(tbl);
  return div;
}
