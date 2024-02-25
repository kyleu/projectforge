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

export function modalNew(key: string, title: string) {
  const el = document.createElement("div");
  el.classList.add("modal");
  el.id = "modal-" + key;
  el.style.display = "none";

  const backdrop = document.createElement("a");
  backdrop.classList.add("backdrop");
  backdrop.href = "#";
  el.appendChild(backdrop);

  const content = document.createElement("div");
  content.classList.add("modal-content");
  el.appendChild(content);

  const hd = document.createElement("div");
  hd.classList.add("modal-header");
  content.appendChild(hd);

  const hdClose = document.createElement("a");
  hdClose.classList.add("modal-close");
  hdClose.href = "#";
  hdClose.innerText = "×";
  hd.appendChild(hdClose);

  const hdTitle = document.createElement("h2");
  hdTitle.innerText = title;
  hd.appendChild(hdTitle);

  const bd = document.createElement("div");
  bd.classList.add("modal-body");
  content.appendChild(bd);

  document.body.appendChild(el);
  return el;
}

export function modalGetOrCreate(key: string, title: string): HTMLElement {
  const el = document.getElementById("modal-" + key);
  if (el) {
    const bodies = el.getElementsByClassName("modal-body");
    if (bodies.length !== 1) {
      throw new Error("unable to find modal body");
    }
    return el;
  }
  return modalNew(key, title);
}

export function modalGetBody(m: HTMLElement) {
  const bodies = m.getElementsByClassName("modal-body");
  if (bodies.length !== 1) {
    throw new Error("unable to find modal body");
  }
  return bodies[0];
}

export function modalGetHeader(m: HTMLElement) {
  const bodies = m.getElementsByClassName("modal-header");
  if (bodies.length !== 1) {
    throw new Error("unable to find modal header");
  }
  return bodies[0];
}

export function modalSetTitle(m: HTMLElement, title: string) {
  const hd = modalGetHeader(m);
  const t = hd.getElementsByTagName("h2");
  if (t.length !== 1) {
    throw new Error("unable to find modal title");
  }
  t[0].innerText = title;
}
