import { els } from "./dom";

export function menuInit() {
  for (const n of els(".menu-container .final")) {
    n.scrollIntoView({ block: "center" });
  }
}
