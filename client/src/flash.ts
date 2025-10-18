function fade(el: HTMLElement) {
  setTimeout(() => {
    el.style.opacity = "0";
    setTimeout(() => {
      el.remove();
    }, 500);
  }, 5000);
}

export function flashCreate(key: string, level: "success" | "error", msg: string) {
  let container = document.getElementById("flash-container");
  if (container === null) {
    container = document.createElement("div");
    container.id = "flash-container";
    document.body.insertAdjacentElement("afterbegin", container);
  }
  const fl = document.createElement("div");
  fl.className = "flash";

  const radio = document.createElement("input");
  radio.type = "radio";
  radio.style.display = "none";
  radio.id = "hide-flash-" + key;
  fl.appendChild(radio);

  const label = document.createElement("label");
  label.htmlFor = "hide-flash-" + key;
  const close = document.createElement("span");
  close.innerHTML = "×";
  label.appendChild(close);
  fl.appendChild(label);

  const content = document.createElement("div");
  content.className = "content flash-" + level;
  content.innerText = msg;
  fl.appendChild(content);

  container.appendChild(fl);
  fade(fl);
}

export function flashInit() {
  const container = document.getElementById("flash-container");
  if (container === null) {
    return flashCreate;
  }
  const x = container.querySelectorAll<HTMLElement>(".flash");
  if (x.length > 0) {
    x.forEach(fade);
  }
  return flashCreate;
}
