// Content managed by Project Forge, see [projectforge.md] for details.
import {setHTML} from "./dom";

declare global {
  namespace JSX { // eslint-disable-line @typescript-eslint/no-namespace, no-shadow
    type IntrinsicElements = {
      [elemName: string]: unknown;
    }
  }
}

// noinspection JSUnusedGlobalSymbols
export function JSX(tag: string, attrs: any[], ...args: Node[]) { // eslint-disable-line @typescript-eslint/no-explicit-any
  const e = document.createElement(tag);
  for (let name in attrs) {
    if (name === "for") {
      name = "for";
    }
    if (name === "className") {
      name = "class";
    }
    if (name && attrs.hasOwnProperty(name)) { // eslint-disable-line no-prototype-builtins
      const v = attrs[name];
      if (name === "dangerouslySetInnerHTML") {
        setHTML(e, v.__html); // eslint-disable-line no-underscore-dangle
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
