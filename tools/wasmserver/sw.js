importScripts("wasm_exec.js");

const wasmURI = "projectforge.wasm";
let WasmApp, WasmAppStatus;

function getLog(s) {
  return "### ServiceWorker: " + s;
}

function explodedPromise() {
  var status = "pending", value = null;
  const get = () => ({ status, value });
  var resolve, reject;
  var promise = new Promise((_resolve, _reject) => {
    resolve = (_value) => {
      status = "resolved";
      value = _value;
      _resolve(_value);
    };
    reject = (_value) => {
      status = "rejected"
      value = _value;
      _reject(_value)
    };
  });
  return [promise, get, resolve, reject];
}

const loadWasmApp = (() => {
  var locked = false;
  var currentEtag;
  var resolveApp, rejectApp;
  [WasmApp, WasmAppStatus, resolveApp, rejectApp] = explodedPromise();
  return async (trigger = "unknown") => {
    if (locked) {
      console.log(getLog("skipped redundant checking for new wasm app"), { trigger });
      return;
    }
    try {
      locked = true;
      let response = await fetch(wasmURI, { cache: "no-cache" });
      let newEtag = response.headers.get('etag');
      if (newEtag && newEtag === currentEtag) {
        console.log(getLog("skipped reinstalling wasm app with matching etag"), { trigger, etag: newEtag });
        locked = false;
        return;
      }

      {
        var { status } = WasmAppStatus();
        if (status === "resolved") {
          console.log(getLog("stopping old wasm app"), { trigger, etag: currentEtag });
          [WasmApp, WasmAppStatus, resolveApp, rejectApp] = explodedPromise();
        }
      }
      try {
        console.log("installing new App", { trigger, newEtag });
        const go = new Go();
        go.argv = [wasmURI, "wasm"];
        const newApp = await WebAssembly.instantiateStreaming(response, go.importObject);
        go.run(newApp.instance);
        resolveApp(newApp.instance);
        currentEtag = newEtag;
      } catch (error) {
        console.error(getLog("failed to install wasm app"), { error })
        rejectApp(error);
      }
    } catch (error) {
      console.error(getLog("error thrown while updating wasm app"), { error })
    }
    finally {
      locked = false;
    }
  };
})();

setInterval(() => loadWasmApp("interval"), 15 * 60 * 1000);

self.addEventListener("install", (event) => {
  console.log(getLog("install event"));
  event.waitUntil(loadWasmApp("install"));
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  console.log(getLog("activate event"));
  event.waitUntil(loadWasmApp("activate"));
  event.waitUntil(clients.claim());
});

self.addEventListener('message', (event) => {
  console.log(getLog("message event"), { type: event.data.type, event });
  if (event.data.type === 'clientattached') {
    event.waitUntil(loadWasmApp("clientattached"));
  }
});

function isLocal(url) {
  if (url.toString().startsWith("https://projectforge.dev/")) {
    return true;
  }
  for (const s of ["favicon.ico", "logo.svg", "projectforge.wasm", "server.js", "sw.js", "wasm_exec.js"]) {
    if (url.pathname.endsWith(s)) {
      return true;
    }
  }
  return false;
}

self.addEventListener("fetch", (event) => {
  let url = new URL(event.request.url);
  const sameOrigin = url.origin === event.target.location.origin;
  let shouldOverride = sameOrigin && !isLocal(url) && WasmAppStatus().status === "resolved";
  console.log(getLog("fetch event received"), { overriding: shouldOverride, method: event.request.method, url, event });
  if (!shouldOverride) {
    return;
  }
  event.respondWith((async () => {
    try {
      const headers = [];
      event.request.headers.forEach((value, key) => {
        headers.push([key, value]);
      });
      const reqBody = await event.request.text();
      return getResponse(event.request, headers, reqBody);
    } catch (error) {
      console.error(getLog("error calling wasm app"), { error, event });
    }
  })());
});

async function getResponse(req, headers, reqBody) {
  await WasmApp;
  console.log(getLog("request to wasm app"), req);
  const rsp = goFetch(req, headers, reqBody);
  console.log(getLog("response received from wasm app"), rsp);
  return rsp
}
