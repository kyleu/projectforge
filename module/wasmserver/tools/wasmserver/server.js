function startApp() {
  wasmLoad();
  // initServiceWorker();
}

async function initServiceWorker() {
  if (navigator.serviceWorker === undefined) {
    return;
  }
  const reg = await navigator.serviceWorker.register("sw.js", {
    scope: "/"
  })
  if (reg.installing) {
    console.log("Service worker installing");
  } else if (reg.waiting) {
    console.log("Service worker installed");
  } else if (reg.active) {
    console.log("Service worker active");
  }
}

async function wasmInit(ms) {
  document.getElementById("load-status").innerText = "Loaded in [" + ms + "ms]";
}

function wasmLoad() {
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const start = new Date().getTime();
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("projectforge.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    wasmInit(new Date().getTime() - start);
  });
}
