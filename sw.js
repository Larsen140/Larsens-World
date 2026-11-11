/* Larsens World: immer zuerst die neueste Version aus dem Netz holen; ohne Internet die zuletzt geladene Version zeigen. */
const C='lw-v1';
self.addEventListener('install',()=>self.skipWaiting());
self.addEventListener('activate',e=>e.waitUntil(self.clients.claim()));
self.addEventListener('fetch',e=>{
  const r=e.request;if(r.method!=='GET'||new URL(r.url).origin!==location.origin)return;
  e.respondWith(fetch(r).then(res=>{if(res.ok){const cp=res.clone();caches.open(C).then(c=>c.put(r,cp)).catch(()=>{})}return res})
    .catch(()=>caches.match(r,{ignoreSearch:true}).then(m=>m||caches.match('./',{ignoreSearch:true}))));
});
