/* Minimal service worker so Android offers "App installieren". It caches nothing: every request goes to the network, so updates arrive immediately. */
self.addEventListener('install',()=>self.skipWaiting());
self.addEventListener('activate',e=>e.waitUntil(self.clients.claim()));
self.addEventListener('fetch',e=>e.respondWith(fetch(e.request)));
