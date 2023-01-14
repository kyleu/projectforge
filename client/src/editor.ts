// Content managed by Project Forge, see [projectforge.md] for details.
const selected = "--selected";

export function setSiblingToNull(el: HTMLElement) {
  const i = el.parentElement!.parentElement!.querySelector("input");
  if (!i) {
    throw "no associated input found";
  }
  i.value = "∅";
}

export function initForm(frm: HTMLFormElement) {
  frm.onreset = () => initForm(frm);
  let editorCache: { [key: string]: string; } = {};
  let selectedCache: { [key: string]: HTMLInputElement; } = {};
  for (const el of frm.elements) {
    const input = el as HTMLInputElement;
    if (input.name.length > 0) {
      if (input.name.endsWith(selected)) {
        selectedCache[input.name] = input;
      } else {
        if ((input.type !== "radio") || input.checked) {
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

export function editorInit() {
  (window as any).projectforge.setSiblingToNull = setSiblingToNull;
  (window as any).projectforge.initForm = initForm;
  for (const n of Array.from(document.querySelectorAll<HTMLFormElement>("form.editor"))) {
    initForm(n);
  }
}
