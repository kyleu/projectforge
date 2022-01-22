// Content managed by Project Forge, see [projectforge.md] for details.
export function linkInit() {
  for (const l of Array.from(document.getElementsByClassName("link-confirm"))) {
    const el = (l as HTMLAnchorElement);
    el.onclick = function() {
      let msg = el.dataset.message as string;
      if (msg && msg.length === 0) {
        msg = "Are you sure?";
      }
      return confirm(msg);
    }
  }
}
