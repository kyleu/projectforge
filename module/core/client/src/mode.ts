import {els} from "./dom";

const l = "mode-light";
const d = "mode-dark";

export function modeInit() {
  for (const el of els<HTMLInputElement>(".mode-input")) {
    el.onclick = () => {
      switch (el.value) {
        case "":
          document.body.classList.remove(l);
          document.body.classList.remove(d);
          break;
        case "light":
          document.body.classList.add(l);
          document.body.classList.remove(d);
          break;
        case "dark":
          document.body.classList.remove(l);
          document.body.classList.add(d);
          break;
      }
    };
  }
}
