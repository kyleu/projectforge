// Content managed by Project Forge, see [projectforge.md] for details.
function debounce(callback: (...args: unknown[]) => void, wait: number) {
  let timeoutId = 0;
  return (...args: unknown[]) => {
    if (timeoutId !== 0) {
      window.clearTimeout(timeoutId);
    }
    timeoutId = window.setTimeout(() => {
      callback(null, ...args);
    }, wait);
  };
}

function autocomplete(el: HTMLInputElement, url: string, field: string, title: (x: unknown) => string, val: (x: unknown) => string) {
  if (!el) {
    return;
  }
  const listId = el.id + "-list";
  const list = document.createElement("datalist");

  const loadingOpt = document.createElement("option");
  loadingOpt.value = "";
  loadingOpt.innerText = "Loading...";
  list.appendChild(loadingOpt);

  list.id = listId;
  el.parentElement?.prepend(list);

  el.setAttribute("autocomplete", "off");
  el.setAttribute("list", listId);

  const cache: {
    [_: string]: {
      url: string;
      complete: boolean;
      data: unknown[];
      frag: DocumentFragment;
    }
  } = {};
  let lastQuery = "";

  function getURL(q: string): string {
    const dest = url;
    if (dest.includes("?")) {
      return dest + "&t=json&" + field + "=" + encodeURIComponent(q);
    }
    return dest + "?t=json&" + field + "=" + encodeURIComponent(q);
  }

  function datalistUpdate(q: string) {
    const c = cache[q];
    if (!c || !c.frag) {
      return;
    }
    lastQuery = q;
    list.replaceChildren(c.frag.cloneNode(true));
  }

  function f() {
    const q = el.value;
    if (q.length === 0) {
      return;
    }
    const dest = getURL(q);
    let proceed: boolean = !q || !lastQuery;
    if (!proceed) {
      const l = cache[lastQuery];
      if (l) {
        proceed = !l.data.find((d) => q === val(d));
      }
    }
    if (!proceed) {
      return;
    }
    if (cache[q] && cache[q].url === dest) {
      datalistUpdate(q);
      return;
    }

    fetch(dest, {credentials: "include"}).then((res) => res.json()).then((data) => {
      if (!data) {
        return;
      }
      const arr = Array.isArray(data) ? data : [data];
      const frag = document.createDocumentFragment();
      let optMax = 10;
      for (let d = 0; d < arr.length && optMax > 0; d++) {
        const v = val(arr[d]);
        const t = title(arr[d]);
        if (v) {
          const option = document.createElement("option");
          option.value = v;
          option.innerText = t;
          frag.appendChild(option);
          optMax--;
        }
      }
      cache[q] = {url: dest, data: arr, frag: frag, complete: false};
      datalistUpdate(q);
    });
  }

  el.oninput = debounce(f, 250);
  console.log("managing [" + el.id + "] autocomplete");
}

export function autocompleteInit() {
  return autocomplete;
}
