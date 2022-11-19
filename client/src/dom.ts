// Content managed by Project Forge, see [projectforge.md] for details.
export function els<T extends HTMLElement>(selector: string, context?: HTMLElement): readonly T[] {
  let result: NodeListOf<Element>;
  if (context) {
    result = context.querySelectorAll(selector);
  } else {
    result = document.querySelectorAll(selector);
  }
  const ret: T[] = []
  result.forEach(v => {
    ret.push(v as T);
  });
  return ret;
}

export function opt<T extends HTMLElement>(selector: string, context?: HTMLElement): T | undefined {
  const e = els<T>(selector, context);
  switch (e.length) {
    case 0:
      return undefined;
    case 1:
      return e[0];
    default:
      console.warn(`found [${e.length}] elements with selector [${selector}], wanted zero or one`);
  }
}

export function req<T extends HTMLElement>(selector: string, context?: HTMLElement): T {
  const res = opt<T>(selector, context);
  if (!res) {
    console.warn(`no element found for selector [${selector}]`);
  }
  return res!;
}

export function setHTML(el: string | HTMLElement, html: string) {
  if (typeof el === "string") {
    el = req(el);
  }
  el.innerHTML = html;
  return el;
}

export function setDisplay(el: string | HTMLElement, condition: boolean, v: string = "block") {
  if (typeof el === "string") {
    el = req(el);
  }

  el.style.display = condition ? v : "none";
  return el;
}

export function setText(el: string | HTMLElement, text: string): HTMLElement {
  if (typeof el === "string") {
    el = req(el);
  }
  el.innerText = text;
  return el;
}

export function clear(el: string | HTMLElement) {
  return setHTML(el, "");
}
