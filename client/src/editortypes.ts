// Content managed by Project Forge, see [projectforge.md] for details.
export function typeKey(t?: string) {
  if (t && t !== "") {
    return t;
  }
  return "string";
}

export type Column = {
  key: string
  title: string
  description?: string
  type?: string
};

export type Editor = {
  key: string
  title: string
  columns: Column[]
  textarea: HTMLTextAreaElement,
  rows: { [key: string]: unknown; }[]
  table?: Element,
};
