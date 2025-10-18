import { els } from "./dom";

export function linkInit() {
  for (const el of els(".link-confirm")) {
    el.onclick = () => {
      return confirm(el.dataset.message ?? "Are you sure?");
    };
  }
}
