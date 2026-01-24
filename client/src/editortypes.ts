export function typeKey(t?: string) {
  if (t && t !== "") {
    return t;
  }
  return "string";
}

export type Type =
  | string
  | {
      k: string;
      t: Record<string, unknown>;
    };

export function typeToString(t: Type): string {
  if (typeof t === "string") {
    return t;
  }
  switch (t.k) {
    case "enum":
      return "enum(" + String(t.t.ref) + ")";
    case "list":
      return "[]" + typeToString(t.t.v as Type);
    default:
      return t.k;
  }
}

export interface Column {
  key: string;
  title: string;
  description?: string;
  type?: string;
}

export interface Editor {
  key: string;
  title: string;
  columns: Column[];
  textarea: HTMLTextAreaElement;
  rows: Record<string, unknown>[];
  table?: HTMLTableElement;
  tableWrapper?: HTMLDivElement;
}
