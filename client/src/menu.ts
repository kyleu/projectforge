// Content managed by Project Forge, see [projectforge.md] for details.
export function menuInit() {
  for (const n of Array.from(document.querySelectorAll(".menu-container .final"))) {
    n.scrollIntoView({block: "nearest"});
  }
}
