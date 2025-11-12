const APP_VERSION = "12.2-omega-final";
const CACHE = `feuer-${APP_VERSION}`;
const ASSET_PATHS = [
  "./",
  "./index.html",
  "./manifest.webmanifest",
  "./assets/icons/icon-192.png",
  "./assets/icons/icon-512.png"
];
const VERSIONED_ASSETS = ASSET_PATHS.map(versionedUrl);
const FALLBACK_KEY = versionedUrl("./");

let lastActivationHadUpdate = false;

self.addEventListener("install", (e) => {
  self.skipWaiting();
  e.waitUntil(caches.open(CACHE).then((c) => c.addAll(VERSIONED_ASSETS)));
});

self.addEventListener("activate", (e) => {
  e.waitUntil((async () => {
    const keys = await caches.keys();
    const oldKeys = keys.filter((k) => k !== CACHE);
    lastActivationHadUpdate = oldKeys.length > 0;
    await Promise.all(oldKeys.map((k) => caches.delete(k)));
    await self.clients.claim();
    const clients = await self.clients.matchAll({ type: "window" });
    clients.forEach((client) => client.postMessage({
      type: "SW_VERSION",
      version: APP_VERSION,
      updated: lastActivationHadUpdate
    }));
    lastActivationHadUpdate = false;
  })());
});

self.addEventListener("fetch", (e) => {
  const req = e.request;
  if (req.method !== "GET") return;

  const url = new URL(req.url);
  if (url.origin !== self.location.origin) return;

  const cacheKey = versionedUrl(url);

  e.respondWith((async () => {
    const cache = await caches.open(CACHE);
    const cached = await cache.match(cacheKey);
    if (cached) return cached;

    try {
      const res = await fetch(req);
      if (res && res.ok) {
        const copy = res.clone();
        cache.put(cacheKey, copy);
      }
      return res;
    } catch (err) {
      const fallback = await cache.match(cacheKey) || await cache.match(FALLBACK_KEY);
      if (fallback) return fallback;
      throw err;
    }
  })());
});

self.addEventListener('message', (event) => {
  if (!event.data) return;
  if (event.data.type === 'SKIP_WAITING') self.skipWaiting();
  if (event.data.type === 'REQUEST_VERSION') {
    if (event.source && typeof event.source.postMessage === 'function') {
      event.source.postMessage({
        type: 'SW_VERSION',
        version: APP_VERSION,
        updated: false
      });
    }
  }
});

function versionedUrl(input) {
  const base = typeof input === "string" ? input : input.toString();
  const url = new URL(base, self.location.origin);
  url.searchParams.set("v", APP_VERSION);
  return url.toString();
}
