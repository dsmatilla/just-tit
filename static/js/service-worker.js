const version = "0.5";
const cacheName = `just-tit-${version}`;
self.addEventListener('install', e => {
    e.waitUntil(caches.open(cacheName).then(cache => {
        return cache.addAll([`/`]).then(() => self.skipWaiting())
    }))
});
self.addEventListener('activate', event => {
    event.waitUntil(self.clients.claim())
});
self.addEventListener('fetch', event => {
    event.respondWith(caches.open(cacheName).then(cache => cache.match(event.request, {ignoreSearch: 0})).then(response => {
        return response || fetch(event.request)
    }))
})