export const appKey = "{{{ .Key }}}";
export const appName = "{{{ .Name }}}";

export function svgRef(key: string, size?: number, cls?: string): string {
  size ??= 18;
  cls ??= "icon";
  return `<svg class="${cls}" style="width: ${size.toString()}px; height: ${size.toString()}px;"><use xlink:href="#svg-${key}"></use></svg>`;
}

export function svg(key: string, size?: number, cls?: string) {
  return { __html: svgRef(key, size, cls) };
}

export function expandCollapse(extra?: string) {
  extra ??= "";
  const e = svgRef("right", 15, "expand-collapse");
  return { __html: e + extra };
}

export function focusDelay(el: HTMLInputElement | HTMLTextAreaElement) {
  setTimeout(() => {
    el.setSelectionRange(el.value.length, el.value.length);
    el.focus();
  }, 100);
  return true;
}

export function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function unknownToString(value: unknown): string {
  if (value === undefined || value === null) {
    return "-";
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "symbol") {
    return value.toString();
  }
  if (value instanceof Date) {
    return value.toISOString();
  }
  if (Array.isArray(value)) {
    return `[${value.map((item) => unknownToString(item)).join(", ")}]`;
  }
  if (isRecord(value)) {
    return JSON.stringify(value, null, 2);
  }
  if (typeof value === "function") {
    return value.name || value.toString();
  }
  if (value && typeof value === "object") {
    try {
      return JSON.stringify(value);
    } catch {
      return Object.prototype.toString.call(value);
    }
  }
  return "";
}

// eslint-disable-next-line @typescript-eslint/no-unnecessary-type-parameters
export function jsonParse<T>(json: string): T | null {
  try {
    return JSON.parse(json) as T;
  } catch {
    return null;
  }
}
