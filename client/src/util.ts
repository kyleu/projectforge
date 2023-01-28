// Content managed by Project Forge, see [projectforge.md] for details.
export const appKey = "projectforge";
export const appName = "Project Forge";

export function svgRef(key: string, size: number): string {
  return `<svg class="icon" style="width: ${size}px; height: ${size}px;"><use xlink:href="#svg-${key}"></use></svg>`;
}

export function svg(key: string, cls?: string) {
  return {
    "__html": `<svg class="${cls || ""}" style="width: 18px; height: 18px;"><use href="#svg-${key}"></use></svg>`
  }
}
