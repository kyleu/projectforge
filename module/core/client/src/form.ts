export function setSiblingToNull(el: HTMLElement) {
  const i = el.parentElement!.parentElement!.querySelector("input");
  if (!i) {
    throw "no associated input found";
  }
  i.value = "∅";
}
