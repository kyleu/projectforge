import { isRecord, jsonParse } from "./util";
import { els, req } from "./dom";
import { createTable } from "./editortable";
import type { Column, Editor } from "./editortypes";

const rawLabel = "Raw JSON";

function parseRows(value: string, context?: string): Record<string, unknown>[] {
  const parsed: unknown = JSON.parse(value);
  if (parsed == null) {
    return [];
  }
  if (!Array.isArray(parsed)) {
    throw new Error("input value" + (context ? " for " + context : "") + " must be an array");
  }
  return parsed.map((row) => {
    if (!isRecord(row)) {
      throw new Error("row values" + (context ? " for " + context : "") + " must be objects");
    }
    return row;
  });
}

function extractEditor(el: HTMLElement) {
  const key = el.dataset.key ?? "editor";
  const title = el.dataset.title ?? "Object";
  const columnsStr = el.dataset.columns ?? "[]";
  const columns = jsonParse<Column[]>(columnsStr) ?? [];

  const inp: HTMLTextAreaElement = req<HTMLTextAreaElement>("textarea", el);
  const curr = parseRows(inp.value, "editor [" + key + "]");
  inp.hidden = true;

  const e: Editor = { key: key, title: title, columns: columns, textarea: inp, rows: curr };

  const wrapper = createTable(e);
  el.appendChild(wrapper);

  return e;
}

function editorShow(el: HTMLElement, e: Editor) {
  els(".toggle-editor-" + e.key).forEach((toggle) => {
    toggle.innerText = rawLabel;
  });
  e.rows = parseRows(e.textarea.value, "editor [" + e.key + "]");
  e.tableWrapper?.remove();
  const wrapper = createTable(e);
  el.appendChild(wrapper);
  e.textarea.hidden = true;
}

function editorHide(e: Editor, editorLabel: string) {
  els(".toggle-editor-" + e.key).forEach((toggle) => {
    toggle.innerText = editorLabel;
  });
  e.tableWrapper?.remove();
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
