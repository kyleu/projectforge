import {els} from "./dom";

export function linkInit() {
  for (const el of els(".link-confirm")) {
    el.onclick = () => {
      let msg = el.dataset.message as string;
      if (msg && msg.length === 0) {
        msg = "Are you sure?";
      }
      return confirm(msg);
    };
  }
}
