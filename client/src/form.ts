import { els } from "./dom";

const selected = "--selected";

export function setSiblingToNull(el: HTMLElement) {
  const i = el.parentElement?.parentElement?.querySelector("input");
  if (!i) {
    throw new Error("no associated input found");
  }
  i.value = "∅";
}

export function initForm(frm: HTMLFormElement) {
  frm.onreset = () => initForm(frm);
  const editorCache: { [key: string]: string } = {};
  const selectedCache: { [key: string]: HTMLInputElement } = {};
  for (const el of frm.elements) {
    const input = el as HTMLInputElement;
    if (input.name.length > 0) {
      if (input.name.endsWith(selected)) {
        selectedCache[input.name] = input;
      } else {
        if (input.type !== "radio" || input.checked) {
          editorCache[input.name] = input.value;
        }
        const evt = () => {
          const cv = selectedCache[input.name + selected];
          if (cv) {
            cv.checked = editorCache[input.name] !== input.value;
          }
        };
        input.onchange = evt;
        input.onkeyup = evt;
      }
    }
  }
}

export function formInit(): [(el: HTMLElement) => void, (frm: HTMLFormElement) => void] {
  for (const n of els<HTMLFormElement>("form.editor")) {
    initForm(n);
  }
  return [setSiblingToNull, initForm];
}
