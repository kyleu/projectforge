import {els, opt, req, setDisplay} from "./dom";
import {svgRef} from "./util";

function compareOrder(elem1: HTMLElement, elem2: HTMLElement) {
  if (elem1.parentElement !== elem2.parentElement) {
    return null;
  }
  if (elem1 === elem2) {
    return 0;
  }
  if (elem1.compareDocumentPosition(elem2) & Node.DOCUMENT_POSITION_FOLLOWING) {
    return -1;
  }
  return 1;
}

let draggedElement: HTMLElement;

function tagsUpdate(editorEl: HTMLElement) {
  const values: string[] = [];
  const elements = els(".item .value", editorEl);
  for (const el of elements) {
    values.push(el.innerText);
  }
  const inputEl = req<HTMLInputElement>("input.result", editorEl);
  inputEl.value = values.join(", ");
}

function tagsRemove(itemEl: HTMLElement) {
  const editorEl = itemEl.parentElement?.parentElement;
  itemEl.remove();
  if (editorEl) {
    tagsUpdate(editorEl);
  }
}

function tagsEdit(itemEl: HTMLElement) {
  const value = req(".value", itemEl);
  const edit = req<HTMLInputElement>(".editor", itemEl);
  edit.value = value.innerText;
  const apply = () => {
    if (edit.value === "") {
      itemEl.remove();
      return;
    }
    value.innerText = edit.value;
    setDisplay(value, true);
    setDisplay(edit, false);
    const editorEl = itemEl.parentElement?.parentElement;
    if (editorEl) {
      tagsUpdate(editorEl);
    }
  };
  edit.onblur = apply;
  edit.onkeydown = (event) => {
    if (event.code === "Enter") {
      event.preventDefault();
      apply();
      return false;
    }
  };
  setDisplay(value, false);
  setDisplay(edit, true);
  edit.focus();
}

function tagsRender(v: string, editorEl: HTMLElement): HTMLDivElement {
  const item = document.createElement("div");
  item.className = "item";
  item.draggable = true;
  item.ondragstart = (e) => {
    e.dataTransfer?.setDragImage(document.createElement("div"), 0, 0);
    item.classList.add("dragging");
    draggedElement = item;
  };
  item.ondragover = () => {
    const order = compareOrder(item, draggedElement);
    if (!order) {
      return;
    }
    const baseElement = order === -1 ? item : item.nextSibling;
    draggedElement.parentElement?.insertBefore(draggedElement, baseElement);
    tagsUpdate(editorEl);
  };
  item.ondrop = (e) => {
    e.preventDefault();
  };
  item.ondragend = (e) => {
    item.classList.remove("dragging");
    e.preventDefault();
  };

  const value = document.createElement("div");
  value.innerText = v;
  value.className = "value";
  value.onclick = () => {
    tagsEdit(item);
  };
  item.appendChild(value);

  const editor = document.createElement("input");
  editor.className = "editor";
  item.appendChild(editor);

  const close = document.createElement("div");
  close.innerHTML = svgRef("times", 13);
  close.className = "close";
  close.onclick = () => {
    tagsRemove(item);
  };
  item.appendChild(close);

  return item;
}

function tagsAdd(tagContainerEl: HTMLElement, editorEl: HTMLElement) {
  const itemEl = tagsRender("", editorEl);
  tagContainerEl.appendChild(itemEl);
  tagsEdit(itemEl);
}

export function tagsWire(el: HTMLElement) {
  const input = opt<HTMLInputElement>("input.result", el);
  if (!input) {
    return;
  }
  const tagContainer = req<HTMLDivElement>(".tags", el);
  const vals = input.value.split(",").map((x) => x.trim()).filter((k) => k !== "");

  setDisplay(input, false);
  tagContainer.innerHTML = "";
  for (const v of vals) {
    tagContainer.appendChild(tagsRender(v, el));
  }

  opt(".add-item", el)?.remove();

  const add = document.createElement("div");
  add.className = "add-item";
  add.innerHTML = svgRef("plus", 22);
  add.onclick = () => {
    tagsAdd(tagContainer, el);
  };
  el.insertBefore(add, req(".clear", el));
}

export function tagsInit() {
  for (const el of els(".tag-editor")) {
    tagsWire(el);
  }
  return tagsWire;
}
