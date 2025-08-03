import { type Column, typeKey } from "./editortypes";

function stringInput(id: string, col: Column, x: { [p: string]: unknown }, onChange?: () => void) {
  const input = document.createElement("input");
  input.name = col.key;
  input.id = id;
  const v = x[col.key];
  if (v !== null && v) {
    input.value = v.toString();
  }
  if (!onChange) {
    onChange = () => {
      x[col.key] = input.value;
    };
  }
  input.onchange = onChange;
  return input;
}

function boolInput(col: Column, x: { [p: string]: unknown }) {
  const vx = x[col.key];
  const v = vx !== null && (vx === "true" || vx === true);

  const div = document.createElement("div");

  [true, false].forEach((b: boolean) => {
    const label = document.createElement("label");
    div.appendChild(label);
    const input = document.createElement("input");
    label.appendChild(input);
    input.name = col.key;
    input.type = "radio";
    input.value = b ? "true" : "false";
    input.checked = v === b;
    input.onclick = () => {
      x[col.key] = b;
      return true;
    };
    const span = document.createElement("span");
    label.appendChild(span);
    span.innerText = b ? " True " : " False ";
  });
  return div;
}

function intInput(id: string, col: Column, x: { [p: string]: unknown }) {
  const input = stringInput(id, col, x, () => {
    x[col.key] = parseInt(input.value, 10);
  });
  input.type = "number";
  return input;
}

function typeInput(id: string, col: Column, x: { [p: string]: unknown }) {
  const textarea = document.createElement("textarea");
  textarea.name = col.key;
  textarea.id = id;
  textarea.value = JSON.stringify(x[col.key], null, 2);
  return textarea;
}

export function createEditorInput(id: string, col: Column, x: { [p: string]: unknown }): HTMLElement {
  const t = typeKey(col.type);
  switch (t) {
    case "bool":
      return boolInput(col, x);
    case "int":
      return intInput(id, col, x);
    case "type":
      return typeInput(id, col, x);
    default:
      return stringInput(id, col, x);
  }
}
