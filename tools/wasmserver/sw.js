self.addEventListener("install", (event) => {
    event.waitUntil(async () => {
        return console.log("!");
    });
});

self.addEventListener("fetch", (event) => {
    console.log(event);
    event.respondWith(fetch(event.request));
});
