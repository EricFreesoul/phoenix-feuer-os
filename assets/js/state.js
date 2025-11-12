import { at00, clamp, id, todayISO } from "./utils.js";

export const VERSION = "12.2-omega-final";
export const STORAGE_KEY = "PHOENIX_V12_STATE_OMEGA";
export const PHOENIX_START = "2025-11-07";
export const PHOENIX_END = "2026-01-31";

const createMemoryStorage = () => {
  const store = new Map();
  return {
    getItem: (key) => (store.has(key) ? store.get(key) : null),
    setItem: (key, value) => store.set(key, String(value)),
    removeItem: (key) => store.delete(key),
  };
};

const storage = (() => {
  try {
    if (typeof window !== "undefined" && window.localStorage) {
      return window.localStorage;
    }
  } catch (err) {
    // ignore and fall back to memory storage
  }
  return createMemoryStorage();
})();

let statusHandler = () => {};

export const registerStatusHandler = (fn) => {
  statusHandler = typeof fn === "function" ? fn : () => {};
};

export const defaultDaily = (date = todayISO()) => ({
  date,
  mode: "A",
  activeMinutes: 0,
  healthLog: "",
  execLog: "",
  factsLog: "",
  buddhaNote: "",
  outreachToday: 0,
});

export const seedState = () => ({
  version: VERSION,
  theme: "dark",
  activeTab: "generator",
  lastSaved: Date.now(),
  phoenix: {
    start: PHOENIX_START,
    end: PHOENIX_END,
    targetRoute: "A",
    pKonto: false,
    businessAccountCold: false,
    piPrepared: false,
    drvNoticeSent: false,
    ltaPrepared: false,
    zeitjournalStable: false,
    tagXPackageDraft: false,
    routeDecisionNote: "",
  },
  daily: defaultDaily(),
  logbook: [],
  finances: {
    entries: [],
  },
  home: {
    mit: [
      { id: id(), text: "Heute Mini-Betrieb transparent halten (keine Expansion, keine Grauzone).", done: false },
      { id: id(), text: "Zeitjournal & Logs sauber pflegen.", done: false },
      { id: id(), text: "80-Tage-Route kurz spiegeln.", done: false },
    ],
    kpis: { sprints: 0 },
  },
  sprint: {
    running: false,
    startedAt: null,
    durationMin: 90,
    remainingSec: 90 * 60,
    note: "",
    buddhaPauseBefore: false,
  },
  audits: [],
  evidence: [],
  outreach: [],
});

const normalizeState = (raw) => {
  if (!raw || typeof raw !== "object") {
    return seedState();
  }
  const base = seedState();
  const next = {
    ...base,
    ...raw,
    phoenix: { ...base.phoenix, ...(raw.phoenix || {}) },
    daily: raw.daily || defaultDaily(),
    logbook: Array.isArray(raw.logbook) ? raw.logbook : [],
    finances: {
      entries: Array.isArray(raw?.finances?.entries) ? raw.finances.entries : [],
    },
    home: {
      mit: Array.isArray(raw?.home?.mit) ? raw.home.mit : [],
      kpis: { sprints: 0, ...(raw?.home?.kpis || {}) },
    },
    sprint: { ...base.sprint, ...(raw.sprint || {}) },
    audits: Array.isArray(raw.audits) ? raw.audits : [],
    evidence: Array.isArray(raw.evidence) ? raw.evidence : [],
    outreach: Array.isArray(raw.outreach) ? raw.outreach : [],
  };
  next.version = raw.version || VERSION;
  next.theme = raw.theme === "light" ? "light" : "dark";
  next.activeTab = raw.activeTab || "generator";
  return next;
};

const loadStateInternal = () => {
  try {
    const raw = storage.getItem(STORAGE_KEY);
    if (!raw) return null;
    return normalizeState(JSON.parse(raw));
  } catch (err) {
    console.warn("Load failed", err);
    return null;
  }
};

export let state = loadStateInternal() || seedState();

export const setState = (next) => {
  state = normalizeState(next);
  saveState();
};

export const phoenixDayInfo = () => {
  const start = at00(PHOENIX_START);
  const end = at00(PHOENIX_END);
  const now = at00(todayISO());
  const total = Math.round((end - start) / (1000 * 60 * 60 * 24)) + 1;
  const passed = clamp(Math.round((now - start) / (1000 * 60 * 60 * 24)) + 1, 0, total);
  const left = clamp(total - passed, 0, total);
  return { total, passed, left };
};

export const filterLogsByRange = (allLogs, range) => {
  const today = at00(todayISO());
  const sorted = Array.isArray(allLogs) ? allLogs.slice() : [];
  if (range === "all") {
    return sorted.sort((a, b) => at00(a.date) - at00(b.date));
  }
  const cut = new Date(today);
  if (range === "30") cut.setDate(today.getDate() - 30);
  if (range === "60") cut.setDate(today.getDate() - 60);
  return sorted
    .filter((entry) => at00(entry.date) >= cut)
    .sort((a, b) => at00(a.date) - at00(b.date));
};

export const calcGateScore = (gate = {}) => {
  const items = [
    "klartext",
    "belege",
    "struktur",
    "prio",
    "mobil",
    "intern",
    "exec",
  ];
  const ok = items.reduce((count, key) => count + (gate[key] ? 1 : 0), 0);
  return Math.round((100 * ok) / items.length);
};

export const applyGuards = (daily) => {
  if (!daily) return;
  const health = (daily.healthLog || "").toLowerCase();
  const exec = (daily.execLog || "").toLowerCase();
  const mins = daily.activeMinutes || 0;

  if (mins > 180 && !health.includes("erschöpf")) {
    statusHandler("Warnung: >3h ohne Erschöpfungsprotokoll (DRV-Risiko prüfen).");
  }
  if (exec && !exec.match(/intern|prototype|prototyp|vorbereitung|ohne einnahmen|minibetrieb|mini-betrieb|keine rechnung/)) {
    statusHandler("Hinweis: Executive-Log schärfen (kein Eindruck verdeckter Vollselbstständigkeit).");
  }
};

export const evaluateGlobalWarnings = (current) => {
  const warnings = [];
  const data = current || state;
  const daysSet = new Set(data.logbook.map((l) => l.date));
  if (daysSet.size < 30) {
    warnings.push({ lvl: "warn", text: "Zeitjournal unter 30 Tagen – für DRV & IV noch dünn. Kontinuierlich weiterführen." });
  } else if (daysSet.size < 60) {
    warnings.push({ lvl: "info", text: `Zeitjournal ${daysSet.size} Tage – solide Basis im Aufbau.` });
  } else {
    data.phoenix.zeitjournalStable = true;
  }
  const year = new Date().getFullYear().toString();
  const sumYear = data.finances.entries
    .filter((e) => e.type === "EIN" && (e.year === year || (e.date || "").startsWith(year)))
    .reduce((n, e) => n + (Number(e.amount) || 0), 0);
  if (sumYear > 15000) {
    warnings.push({
      lvl: "warn",
      text: `Hinzuverdienst ${sumYear.toFixed(2)} € – aktuelle Grenzen prüfen (individuelles Rentenkonto, Bescheid).`,
    });
  }
  const riskOut = data.outreach.filter((o) => o.status !== "EXPLORATION (Klausel OK)");
  if (riskOut.length) {
    warnings.push({ lvl: "warn", text: "Explorationskontakte ohne bestätigte Startklausel vorhanden – prüfen und nachziehen." });
  }
  const fert = data.audits.filter((a) => a.stage === "FERTIG_INTERNAL").length;
  if (fert < 3) {
    warnings.push({ lvl: "info", text: `Bisher ${fert} fertige Prototypen. Ziel Tag X: mind. 3–5 belastbare Übungs-Audits.` });
  }
  return warnings;
};

export const saveState = () => {
  applyGuards(state.daily);
  state.lastSaved = Date.now();
  try {
    storage.setItem(STORAGE_KEY, JSON.stringify(state));
    statusHandler("gespeichert");
    if (typeof window !== "undefined") {
      window.setTimeout(() => statusHandler("bereit"), 600);
    }
  } catch (err) {
    console.warn("Save failed", err);
    statusHandler("Speicherfehler");
  }
};

export const archiveCurrentDaily = () => {
  const cur = state.daily;
  if (!cur || !cur.date) return;
  const idx = state.logbook.findIndex((l) => l.date === cur.date);
  const clone = JSON.parse(JSON.stringify(cur));
  if (idx > -1) {
    state.logbook[idx] = clone;
  } else {
    state.logbook.push(clone);
  }
};

export const exportStateJSON = () => JSON.stringify(state, null, 2);
