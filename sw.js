/* Larsens World: eigene Dateien immer zuerst frisch aus dem Netz (no-store), nur ohne Netz aus dem Speicher. */
const C='lw-v2';
self.addEventListener('install',()=>self.skipWaiting());
self.addEventListener('activate',e=>e.waitUntil(caches.keys().then(ks=>Promise.all(ks.filter(k=>k!==C).map(k=>caches.delete(k)))).then(()=>self.clients.claim())));
self.addEventListener('fetch',e=>{
  const r=e.request;if(r.method!=='GET')return;const u=new URL(r.url);if(u.origin!==location.origin)return;
  e.respondWith(fetch(r,{cache:'no-store'}).then(res=>{if(res.ok&&!u.pathname.endsWith('version.txt')){const cp=res.clone();caches.open(C).then(c=>c.put(u.pathname,cp)).catch(()=>{})}return res})
    .catch(()=>caches.match(u.pathname).then(m=>m||caches.match(new URL('./',location).pathname))));
});
