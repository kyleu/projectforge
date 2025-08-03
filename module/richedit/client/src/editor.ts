import { els, req } from "./dom";
import { createTable } from "./editortable";
import type { Column, Editor } from "./editortypes";

const rawLabel = "Raw JSON";

function extractEditor(el: HTMLElement) {
  const key = el.dataset.key ?? "editor";
  const title = el.dataset.title ?? "Object";
  const columnsStr = el.dataset.columns ?? "[]";
  const columns: Column[] = JSON.parse(columnsStr.replace(/\\"/gu, '"'));

  const inp: HTMLTextAreaElement = req<HTMLTextAreaElement>("textarea", el);
  let curr: { [key: string]: unknown }[] = JSON.parse(inp.value);
  if (curr === undefined || curr === null) {
    curr = [];
  }
  if (!Array.isArray(curr)) {
    throw new Error("input value for element [" + key + "] of type [" + typeof curr + "] must be an array");
  }
  inp.hidden = true;

  const e: Editor = { key: key, title: title, columns: columns, textarea: inp, rows: curr };

  const tbl = createTable(e);
  el.appendChild(tbl);

  return e;
}

function editorShow(el: HTMLElement, e: Editor) {
  els(".toggle-editor-" + e.key).forEach((toggle) => {
    toggle.innerText = rawLabel;
  });
  e.rows = JSON.parse(e.textarea.value);
  if (e.rows === undefined || e.rows === null) {
    e.rows = [];
  }
  e.table?.remove();
  createTable(e);
  if (e.table) {
    el.appendChild(e.table);
  }
  e.textarea.hidden = true;
}

function editorHide(e: Editor, editorLabel: string) {
  els(".toggle-editor-" + e.key).forEach((toggle) => {
    toggle.innerText = editorLabel;
  });
  e.table?.remove();
  e.textarea.value = JSON.stringify(e.rows, null, 2);
  e.textarea.hidden = false;
}

function wireToggles(el: HTMLElement, e: Editor) {
  els(".toggle-editor-" + e.key).forEach((toggle) => {
    if (toggle.innerText === "") {
      toggle.innerText = "Editor";
    }
    const editorLabel = toggle.innerText;
    toggle.innerText = rawLabel;
    toggle.style.display = "inline";
    toggle.onclick = () => {
      if (toggle.innerText === rawLabel) {
        editorHide(e, editorLabel);
      } else {
        editorShow(el, e);
      }
    };
  });
}

export function createEditor(el: HTMLElement): void {
  const e = extractEditor(el);
  wireToggles(el, e);
}

export function editorInit() {
  els(".rich-editor").forEach((x) => {
    createEditor(x);
  });
}
