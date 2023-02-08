export const appKey = "{{{ .Key }}}";
export const appName = "{{{ .Name }}}";

export function svgRef(key: string, size?: number, cls?: string): string {
  if (!size) {
    size = 18;
  }
  if (cls) {
    return `<svg class="${cls}" style="width: ${size}px; height: ${size + "px"};"><use xlink:href="#svg-${key}"></use></svg>`;
  }
  return `<svg class="icon" style="width: ${size}px; height: ${size + "px"};"><use xlink:href="#svg-${key}"></use></svg>`;
}

export function svg(key: string, size?: number, cls?: string) {
  return {"__html": svgRef(key, size, cls)};
}

export function expandCollapse(extra?: string) {
  if (!extra) {
    extra = "";
  }
  let e = svgRef("right", 15, "expand");
  let c = svgRef("down", 15, "collapse");
  return {"__html": e + c + extra};
}

export function focusDelay(el: HTMLInputElement | HTMLTextAreaElement) {
  setTimeout(() => {
    el.setSelectionRange(el.value.length, el.value.length);
    el.focus();
  }, 100);
  return true;
}
