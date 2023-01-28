export const appKey = "{{{ .Key }}}";
export const appName = "{{{ .Name }}}";

export function svgRef(key: string, size: number): string {
  return `<svg class="icon" style="width: ${size}px; height: ${size}px;"><use xlink:href="#svg-${key}"></use></svg>`;
}

export function svg(key: string, cls?: string) {
  return {
    "__html": `<svg class="${cls || ""}" style="width: 18px; height: 18px;"><use href="#svg-${key}"></use></svg>`
  }
}
