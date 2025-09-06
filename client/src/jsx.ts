import { setHTML } from "./dom";

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace JSX {
    type IntrinsicElements = {
      [elemName: string]: unknown;
    };
  }
}

// noinspection JSUnusedGlobalSymbols
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function JSX(tag: string, attrs: any[], ...args: Node[]) {
  const e = document.createElement(tag);
  for (let name in attrs) {
    if (name === "for") {
      name = "for";
    }
    if (name === "className") {
      name = "class";
    }
    // eslint-disable-next-line no-prototype-builtins
    if (name && attrs.hasOwnProperty(name)) {
      const v = attrs[name];
      if (name === "dangerouslySetInnerHTML") {
        // eslint-disable-next-line no-underscore-dangle
        setHTML(e, v.__html);
      } else if (v === true) {
        e.setAttribute(name, name);
      } else if (v !== false && v !== null && v !== undefined) {
        e.setAttribute(name, v.toString());
      }
    }
  }
  for (let child of args) {
    if (Array.isArray(child)) {
      child.forEach((c) => {
        if (child === undefined || child === null) {
          throw new Error(`child array for tag [${tag}] is ${child}\n${e.outerHTML}`);
        }
        if (c === undefined || c === null) {
          throw new Error(`child for tag [${tag}] is ${c}\n${e.outerHTML}`);
        }
        if (typeof c === "string") {
          c = document.createTextNode(c);
        }
        e.appendChild(c);
      });
    } else if (child === undefined || child === null) {
      throw new Error(`child for tag [${tag}] is ${child}\n${e.outerHTML}`);
    } else {
      if (!child.nodeType) {
        child = document.createTextNode(child.toString());
      }
      e.appendChild(child);
    }
  }
  return e;
}
