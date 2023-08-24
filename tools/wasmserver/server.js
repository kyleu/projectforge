// Content managed by Project Forge, see [projectforge.md] for details.
function startApp() {
  const getLog = (s) => "### Client: " + s;
  if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("sw.js").then(reg => {
      reg.addEventListener('statechange', event => {
        console.log(getLog("received `statechange` event"), { reg, event });
      });
      console.log(getLog("service worker registered"), reg);
    }).catch(err => {
      console.error(getLog("service worker registration failed"), err);
    });
    navigator.serviceWorker.addEventListener('controllerchange', event => {
      console.log(getLog("received `controllerchange` event"), event);
    });
    navigator.serviceWorker.ready.then((reg) => reg.active.postMessage({ type: 'clientattached' }));
  } else {
    console.error(getLog("serviceWorker is missing from `navigator`"));
  }
}
