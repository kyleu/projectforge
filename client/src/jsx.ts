import { unknownToString } from "./util";
import { setHTML } from "./dom";

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace JSX {
    type Element = HTMLElement;
    type IntrinsicElements = Record<string, unknown>;
  }
}

// noinspection JSUnusedGlobalSymbols
export function JSX(tag: string, attrs: Record<string, unknown> | null | undefined, ...args: Node[]) {
  const e = document.createElement(tag);
  const attrSource = attrs ?? {};
  for (const key of Object.keys(attrSource)) {
    let name = key;
    if (name === "className") {
      name = "class";
    }
    if (name && Object.prototype.hasOwnProperty.call(attrSource, key)) {
      const v: unknown = attrSource[key];
      if (name === "dangerouslySetInnerHTML") {
        setHTML(e, unknownToString((v as { __html: string }).__html));
      } else if (v === true) {
        e.setAttribute(name, name);
      } else if (v !== false && v !== null && v !== undefined) {
        e.setAttribute(name, unknownToString(v));
      }
    }
  }
  for (let child of args) {
    if (Array.isArray(child)) {
      child.forEach((c) => {
        if (child === undefined || child === null) {
          throw new Error(`child array for tag [${tag}] is [never]\n${e.outerHTML}`);
        }
        if (c === undefined || c === null) {
          throw new Error(`child for tag [${tag}] is not defined\n${e.outerHTML}`);
        }
        if (typeof c === "string") {
          c = document.createTextNode(c);
        }
        e.appendChild(c);
      });
    } else if (child === undefined || child === null) {
      throw new Error(`child for tag [${tag}] is empty\n${e.outerHTML}`);
    } else {
      if (!child.nodeType) {
        child = document.createTextNode(unknownToString(child));
      }
      e.appendChild(child);
    }
  }
  return e;
}
