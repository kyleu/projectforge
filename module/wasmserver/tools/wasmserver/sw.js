const skipCache = false;

self.addEventListener("install", (event) => {
  event.waitUntil(async () => {
    return console.log("!");
  });
});

self.addEventListener("fetch", (event) => {
  event.respondWith(handleFetch(event.request));
});

async function handleFetch(request) {
  if (request.url.indexOf("/app") > -1) {
    log(json(fetchToJSON(request)));
    return fetch(request);
  } else {
    return cacheFirst(request);
  }
}

async function cacheFirst(request) {
  if (skipCache) {
    return fetch(request);
  }
  const responseFromCache = await caches.match(request);
  if (responseFromCache) {
    return responseFromCache;
  }
  const responseFromNetwork = await fetch(request);
  putInCache(request, responseFromNetwork.clone());
  return responseFromNetwork;
}

async function putInCache(request, response) {
  const cache = await caches.open("v1");
  await cache.put(request, response);
}

function fetchToJSON(r) {
  const headers = Object.fromEntries(r.headers.entries());

  var urlParams = new URLSearchParams(r.url.queryString);
  const qs = [];
  for (const [key, value] of urlParams) {
    qs.push({key, value});
  }

  const postData = {
    "mimeType": "", "params": [], "text": ""
  };

  return {"method": r.method, "url": r.url, "headers": headers, "queryString": qs, "postData": postData};
}

function log() {
  console.log(...arguments);
}

function json(x) {
  return JSON.stringify(x, null, 2)
}
