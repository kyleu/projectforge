// Content managed by Project Forge, see [projectforge.md] for details.
export function modalInit() {
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      if (document.location.hash.startsWith("#modal-")) {
        document.location.hash = "";
      }
    }
  });
}
