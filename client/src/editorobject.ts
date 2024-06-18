import {createEditorInput} from "./editorfield";
import type {Editor} from "./editortypes";
import {modalGetBody, modalGetOrCreate, modalSetTitle} from "./modal";
import {createTable} from "./editortable";

function createEditorButtons() {
  const btns = document.createElement("div");
  btns.classList.add("mt");

  const applyBtn = document.createElement("button");
  btns.appendChild(applyBtn);
  applyBtn.innerText = "Apply";
  applyBtn.type = "submit";

  const spacer = document.createElement("span");
  btns.appendChild(spacer);
  spacer.innerText = " ";

  const cancelBtn = document.createElement("button");
  btns.appendChild(cancelBtn);
  cancelBtn.innerText = "Cancel";
  cancelBtn.type = "button";
  // cancelBtn.onclick = () => window.history.back();
  cancelBtn.onclick = () => {
    window.location.href = "#";
  };

  return btns;
}

function createEditor(e: Editor, x: { [p: string]: unknown }, onComplete: (row: { [p: string]: unknown }) => void) {
  const editCopy: {[p: string]: unknown} = JSON.parse(JSON.stringify(x));

  const div = document.createElement("div");
  div.classList.add("overflow", "full-width");

  const form = document.createElement("form");
  div.appendChild(form);
  form.classList.add("expanded");
  form.onsubmit = () => {
    onComplete(editCopy);
    return false;
  };

  const t = document.createElement("table");
  form.appendChild(t);
  t.classList.add("min-200", "full-width");

  const tbody = document.createElement("tbody");
  t.appendChild(tbody);

  e.columns.forEach((col) => {
    const id = `richedit-${col.key}`;
    const tr = document.createElement("tr");
    tbody.appendChild(tr);

    const th = document.createElement("th");
    tr.appendChild(th);
    th.classList.add("shrink");

    const label = document.createElement("label");
    th.appendChild(label);
    label.htmlFor = id;
    label.innerText = col.title;

    const td = document.createElement("td");
    tr.appendChild(td);

    const input = createEditorInput(id, col, editCopy);
    td.appendChild(input);
  });

  form.appendChild(createEditorButtons());

  return div;
}

function onEditComplete(e: Editor, idx: number, row: { [p: string]: unknown }) {
  e.rows[idx] = row;
  e.textarea.value = JSON.stringify(e.rows, null, 2);
  if (e.table) {
    e.table.replaceWith(createTable(e));
  }
  window.location.href = "#";
}

export function rowEditHandler(idx: number, e: Editor, x: { [p: string]: unknown }) {
  return () => {
    const modal = modalGetOrCreate("rich-editor", "Rich Editor");
    modalSetTitle(modal, "Editing Row");
    const onComplete = (row: { [p: string]: unknown }) => {
      onEditComplete(e, idx, row);
    };
    const objectEditor = createEditor(e, x, onComplete);
    modalGetBody(modal).replaceChildren(objectEditor);
    window.location.href = "#modal-rich-editor";
  };
}
