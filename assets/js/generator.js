import {
  PHOENIX_START,
  VERSION,
  archiveCurrentDaily,
  calcGateScore,
  defaultDaily,
  evaluateGlobalWarnings,
  exportStateJSON,
  filterLogsByRange,
  phoenixDayInfo,
  saveState,
  seedState,
  setState,
  state,
} from "./state.js";
import { clamp, el, escapeHTML, id, todayISO } from "./utils.js";

let viewEl;
let tabsEls = [];
let importInput;
let setStatus = () => {};
let sprintTimer = null;

const GATE_ITEMS = [
  ["klartext", "Klarer, verständlicher Befund"],
  ["belege", "Mindestens 2 belastbare Belege"],
  ["struktur", "Saubere Struktur / Dateinamen"],
  ["prio", "P0/P1/P2 nachvollziehbar"],
  ["mobil", "Mobile-First geprüft"],
  ["intern", "Nur intern, keine echte Beauftragung"],
  ["exec", "Kurzfassung für Tag X möglich"],
];

export const setupRenderer = ({ view, tabs, importEl, onStatus }) => {
  viewEl = view;
  tabsEls = tabs ? Array.from(tabs) : [];
  importInput = importEl;
  setStatus = typeof onStatus === "function" ? onStatus : () => {};
};

export const render = () => {
  if (!viewEl) return;
  const now = new Intl.DateTimeFormat("de-DE", {
    dateStyle: "full",
    timeStyle: "short",
  }).format(new Date());
  el("#datestamp").textContent = now;

  const info = phoenixDayInfo();
  el("#phoenixRange").textContent =
    info.passed <= 0
      ? `Start ab ${PHOENIX_START}`
      : info.left <= 0
      ? `Tag ${info.total} · Fenster abgeschlossen`
      : `Tag ${info.passed}/${info.total} · noch ${info.left}`;

  tabsEls.forEach((tab) => {
    tab.classList.toggle("active", tab.dataset.tab === state.activeTab);
  });

  viewEl.innerHTML = "";

  const renderMap = {
    home: renderHome,
    fokus: renderFokus,
    prototypen: renderPrototypen,
    evidence: renderEvidence,
    outreach: renderOutreach,
    manifest: renderManifest,
    generator: renderGenerator,
  };
  (renderMap[state.activeTab] || renderHome)();
};

const renderHome = () => {
  const wrap = document.createElement("div");
  wrap.className = "grid";

  const warns = evaluateGlobalWarnings(state);
  if (warns.length) {
    const w = document.createElement("div");
    w.className = "card";
    w.innerHTML = `
      <h2>Compliance-Radar</h2>
      <div class="list">
        ${warns
          .map(
            (w) =>
              `<div class="pill ${w.lvl === "warn" ? "warnc" : "okc"}">${escapeHTML(w.text)}</div>`
          )
          .join("")}
      </div>
    `;
    wrap.append(w);
  }

  wrap.append(
    cardDailyProtection(),
    cardMIT(),
    cardLogs(),
    cardBuddha(),
    cardFinanceMini()
  );

  const exportCard = document.createElement("div");
  exportCard.className = "card";
  exportCard.innerHTML = `
    <h2>Backup & Tools</h2>
    <div class="row">
      <button class="btn" id="exp">JSON-Backup exportieren</button>
      <button class="btn" id="imp">JSON importieren</button>
    </div>
    <div class="muted" style="margin-top:6px">
      Import ersetzt den kompletten aktuellen Stand.
      Vorher Backup anlegen.
    </div>
  `;
  exportCard.querySelector("#imp").onclick = () => importInput && importInput.click();
  exportCard.querySelector("#exp").onclick = exportJSON;
  wrap.append(exportCard);

  viewEl.append(wrap);
};

const cardDailyProtection = () => {
  const d = document.createElement("div");
  const day = state.daily.date || todayISO();
  const mins = state.daily.activeMinutes || 0;
  const over = mins > 180;
  d.className = "card";
  d.innerHTML = `
    <h2>Tägliches Schutz-Protokoll</h2>
    <div class="row">
      <div class="ctl">
        <label>Datum</label>
        <input id="dpDate" type="date" value="${day}">
      </div>
      <div class="ctl">
        <label>Arbeitsähnliche Minuten (Mini-Betrieb + Aufbau)</label>
        <input id="dpMins" type="number" min="0" max="600" step="5" value="${mins}">
      </div>
    </div>
    <div class="row wrap" style="margin-top:6px">
      <div class="pill">
        <span>Limit Leitplanke</span>
        <span class="tag" style="color:${over ? "#fb7185" : "#22c55e"}">${mins}/180</span>
      </div>
      <label class="pill">
        <input type="radio" name="mode" value="A" ${state.daily.mode === "A" ? "checked" : ""}>
        <span>Modus A · Lernen/Struktur</span>
      </label>
      <label class="pill">
        <input type="radio" name="mode" value="B" ${state.daily.mode === "B" ? "checked" : ""}>
        <span>Modus B · Übungs-Audits/Templates</span>
      </label>
    </div>
    <div class="muted" style="margin-top:4px">
      Trage nur Tätigkeiten ein, die als Arbeit gelesen werden könnten.
      Regeneration, Buddha-Zeit, Spaziergang bleiben draußen.
    </div>
  `;
  const dateEl = d.querySelector("#dpDate");
  const minEl = d.querySelector("#dpMins");

  dateEl.onchange = (e) => {
    const newDate = e.target.value || todayISO();
    archiveCurrentDaily();
    const exist = state.logbook.find((l) => l.date === newDate);
    state.daily = exist ? JSON.parse(JSON.stringify(exist)) : defaultDaily(newDate);
    saveState();
    render();
  };
  minEl.onchange = (e) => {
    const v = parseInt(e.target.value || "0", 10);
    state.daily.activeMinutes = clamp(Number.isNaN(v) ? 0 : v, 0, 600);
    saveState();
    render();
  };
  d.querySelectorAll('input[name="mode"]').forEach((r) => {
    r.onchange = (ev) => {
      state.daily.mode = ev.target.value;
      saveState();
    };
  });
  return d;
};

const cardMIT = () => {
  const d = document.createElement("div");
  d.className = "card";
  d.innerHTML = `<h2>Primär-Direktiven für heute (max. 3)</h2>`;
  const list = document.createElement("div");
  list.className = "list";
  state.home.mit.forEach((item) => {
    const it = document.createElement("div");
    it.className = "item";
    const chk = document.createElement("div");
    chk.className = "check" + (item.done ? " done" : "");
    chk.textContent = item.done ? "✓" : "";
    chk.onclick = () => {
      item.done = !item.done;
      saveState();
      render();
    };
    const txt = document.createElement("div");
    txt.innerHTML = `<strong>${escapeHTML(item.text)}</strong>`;
    it.append(chk, txt);
    list.append(it);
  });
  const add = document.createElement("div");
  add.className = "row";
  add.innerHTML = `
    <div class="ctl" style="flex:1">
      <label>Neue Direktive</label>
      <input id="mitNew" type="text" maxlength="90" placeholder="Kurz, klar, ohne Pathos.">
    </div>
    <button class="btn ok" id="mitAdd">+</button>
  `;
  add.querySelector("#mitAdd").onclick = () => {
    const v = add.querySelector("#mitNew").value.trim();
    if (!v) return;
    if (state.home.mit.length >= 3) state.home.mit.shift();
    state.home.mit.push({ id: id(), text: v, done: false });
    saveState();
    render();
  };
  d.append(list, add);
  return d;
};

const cardLogs = () => {
  const d = document.createElement("div");
  d.className = "card";
  d.innerHTML = `
    <h2>Logbuch · Drei Perspektiven</h2>
    <div class="ctl">
      <label>Gesundheits-/Belastungslog (für DRV · sachlich)</label>
      <textarea id="lgHealth" placeholder="z. B. 120 Min Belastung, danach Kopfschmerz, 30 Min Ruhe.">${escapeHTML(
        state.daily.healthLog || ""
      )}</textarea>
    </div>
    <div class="ctl">
      <label>Executive Summary (für IV · Vorbereitung, keine verdeckten Umsätze)</label>
      <textarea id="lgExec" placeholder="z. B. Interne Prototypenarbeit, kein Kundenkontakt, keine Rechnungen.">${escapeHTML(
        state.daily.execLog || ""
      )}</textarea>
    </div>
    <div class="ctl">
      <label>Fakten-Log (Was genau erledigt? Neutral.)</label>
      <textarea id="lgFacts" placeholder="z. B. 10:00–11:30: Test-Audit Struktur überarbeitet; 1 Mini-Textauftrag (15 €) dokumentiert.">${escapeHTML(
        state.daily.factsLog || ""
      )}</textarea>
    </div>
    <div class="muted" style="margin-top:4px">
      Diese drei Ebenen müssen zusammenpassen.
      Widersprüche sind Gift in jedem Verfahren.
    </div>
  `;
  d.querySelector("#lgHealth").oninput = (e) => {
    state.daily.healthLog = e.target.value;
    saveState();
  };
  d.querySelector("#lgExec").oninput = (e) => {
    state.daily.execLog = e.target.value;
    saveState();
  };
  d.querySelector("#lgFacts").oninput = (e) => {
    state.daily.factsLog = e.target.value;
    saveState();
  };
  return d;
};

const cardBuddha = () => {
  const d = document.createElement("div");
  d.className = "card";
  d.innerHTML = `
    <h2>Buddha-Checkpoint 🐾</h2>
    <div class="muted">
      Buddha ist dein Kater und dein Taktgeber.
      Wenn du überziehst, Pause mit ihm (ca. 10 Min), dann ein klar gesetzter 90-Minuten-Fokusblock.
      Dieser Bereich ist rein privat und wird in keinem externen Report ausgegeben.
    </div>
    <div class="ctl" style="margin-top:6px">
      <label>Notiz (optional)</label>
      <textarea id="bdNote" placeholder="z. B. 2x Pause mit Buddha, Stress runter, dann fokussierter Sprint.">${escapeHTML(
        state.daily.buddhaNote || ""
      )}</textarea>
    </div>
  `;
  d.querySelector("#bdNote").oninput = (e) => {
    state.daily.buddhaNote = e.target.value;
    saveState();
  };
  return d;
};

const cardFinanceMini = () => {
  const d = document.createElement("div");
  d.className = "card";
  const year = new Date().getFullYear().toString();
  const ein = state.finances.entries
    .filter((e) => e.type === "EIN" && (e.year === year || (e.date || "").startsWith(year)))
    .reduce((n, e) => n + (Number(e.amount) || 0), 0);
  const aus = state.finances.entries
    .filter((e) => e.type === "AUS" && (e.year === year || (e.date || "").startsWith(year)))
    .reduce((n, e) => n + (Number(e.amount) || 0), 0);
  d.innerHTML = `
    <h2>Mini-Betrieb · Finanz-Snapshot</h2>
    <div class="row wrap">
      <div class="pill"><span>Einnahmen ${year}</span><span class="tag">${ein.toFixed(2)} €</span></div>
      <div class="pill"><span>Ausgaben ${year}</span><span class="tag">${aus.toFixed(2)} €</span></div>
      <div class="pill"><span>Ergebnis</span><span class="tag">${(ein - aus).toFixed(2)} €</span></div>
    </div>
    <div class="row" style="margin-top:6px">
      <div class="ctl">
        <label>Datum</label>
        <input id="fiDate" type="date" value="${todayISO()}">
      </div>
      <div class="ctl">
        <label>Typ</label>
        <select id="fiType">
          <option value="EIN">Einnahme (Mini-Auftrag)</option>
          <option value="AUS">Ausgabe (Tool, Technik)</option>
        </select>
      </div>
    </div>
    <div class="row" style="margin-top:4px">
      <div class="ctl">
        <label>Betrag (€)</label>
        <input id="fiAmount" type="number" min="0" step="0.01">
      </div>
      <div class="ctl">
        <label>Kategorie / Notiz</label>
        <input id="fiNote" type="text" placeholder="z. B. Plattform X · Kurzauftrag">
      </div>
    </div>
    <button class="btn ok block" id="fiAdd" style="margin-top:6px">Eintrag speichern</button>
    <div class="muted" style="margin-top:4px">
      Nur echte Buchungen erfassen.
      Diese Liste bildet deine Kurz-EÜR und Hinzuverdienstbasis ab.
    </div>
  `;
  d.querySelector("#fiAdd").onclick = () => {
    const date = d.querySelector("#fiDate").value || todayISO();
    const type = d.querySelector("#fiType").value;
    const amount = parseFloat(d.querySelector("#fiAmount").value || "0");
    const note = d.querySelector("#fiNote").value.trim();
    if (!amount) return;
    state.finances.entries.unshift({
      id: id(),
      date,
      type,
      amount,
      note,
      year: date.slice(0, 4),
    });
    saveState();
    render();
  };
  return d;
};

const renderFokus = () => {
  const wrap = document.createElement("div");
  wrap.className = "grid";
  wrap.append(cardSprint(), cardFocusLog());
  viewEl.append(wrap);
};

const cardSprint = () => {
  const card = document.createElement("div");
  card.className = "card";
  card.innerHTML = `
    <h2>Fokus-Sprint</h2>
    <div class="row">
      <div class="ctl">
        <label>Dauer (Minuten)</label>
        <input id="spDur" type="number" min="30" max="180" step="15" value="${state.sprint.durationMin}">
      </div>
      <div class="ctl">
        <label>Timer</label>
        <div class="badge" id="spTimer">${fmtTimer(state.sprint.remainingSec)}</div>
      </div>
    </div>
    <div class="ctl" style="margin-top:6px">
      <label>Notiz</label>
      <textarea id="spNote" placeholder="Was ist das Ziel dieses Fokusblocks?">${escapeHTML(state.sprint.note || "")}</textarea>
    </div>
    <div class="row" style="margin-top:6px">
      <button class="btn primary" id="spStart">${state.sprint.running ? "Pause" : "Start"}</button>
      <button class="btn" id="spReset">Zurücksetzen</button>
    </div>
  `;
  card.querySelector("#spDur").onchange = (e) => {
    const v = parseInt(e.target.value || "0", 10);
    state.sprint.durationMin = clamp(Number.isNaN(v) ? 30 : v, 30, 180);
    if (!state.sprint.startedAt) {
      state.sprint.remainingSec = state.sprint.durationMin * 60;
    }
    saveState();
    render();
  };
  card.querySelector("#spNote").oninput = (e) => {
    state.sprint.note = e.target.value;
    saveState();
  };
  card.querySelector("#spStart").onclick = () => {
    if (state.sprint.running) {
      stopSprint(false);
      saveState();
      render();
    } else {
      const next = (state.daily.activeMinutes || 0) + state.sprint.durationMin;
      if (next > 180 && typeof window !== "undefined") {
        const ok = window.confirm(
          "Dieser Sprint überschreitet die 3h-Leitplanke. Nur starten, wenn du die Folgen im Gesundheitslog dokumentierst."
        );
        if (!ok) return;
      }
      startSprint();
    }
  };
  card.querySelector("#spReset").onclick = () => {
    stopSprint(true);
    saveState();
    render();
  };
  return card;
};

const cardFocusLog = () => {
  const card = document.createElement("div");
  card.className = "card";
  card.innerHTML = `
    <h2>Ergebnis-Fokus</h2>
    <div class="muted">Nutze diesen Raum für Outcomes, nicht für Aufgabenlisten.</div>
    <div class="ctl" style="margin-top:6px">
      <label>Executive Kurzfassung</label>
      <textarea id="spExec" placeholder="Welche Ergebnisse liefert der Sprint für Tag X?">${escapeHTML(
        state.daily.execLog || ""
      )}</textarea>
    </div>
  `;
  card.querySelector("#spExec").oninput = (e) => {
    state.daily.execLog = e.target.value;
    saveState();
  };
  return card;
};

const startSprint = () => {
  if (state.sprint.running) return;
  state.sprint.running = true;
  if (!state.sprint.startedAt) {
    state.sprint.startedAt = Date.now();
    state.sprint.remainingSec = state.sprint.durationMin * 60;
  }
  if (sprintTimer) clearInterval(sprintTimer);
  sprintTimer = setInterval(() => {
    if (state.sprint.remainingSec > 0) {
      state.sprint.remainingSec -= 1;
      const badge = el("#spTimer");
      if (badge) badge.textContent = fmtTimer(state.sprint.remainingSec);
    } else {
      stopSprint(false);
      if (typeof window !== "undefined") {
        window.alert("Fokus-Sprint beendet. Ergebnis im Fakten-/Exec-Log sichern.");
      }
      state.home.kpis.sprints = (state.home.kpis.sprints || 0) + 1;
      state.daily.activeMinutes = (state.daily.activeMinutes || 0) + state.sprint.durationMin;
      saveState();
      render();
    }
  }, 1000);
  saveState();
  render();
};

const stopSprint = (reset) => {
  state.sprint.running = false;
  if (sprintTimer) {
    clearInterval(sprintTimer);
    sprintTimer = null;
  }
  if (reset) {
    state.sprint.startedAt = null;
    state.sprint.remainingSec = state.sprint.durationMin * 60;
    state.sprint.note = "";
  }
};

const fmtTimer = (sec) => {
  const s = Math.max(0, sec | 0);
  const m = Math.floor(s / 60).toString().padStart(2, "0");
  const r = (s % 60).toString().padStart(2, "0");
  return `${m}:${r}`;
};

const renderPrototypen = () => {
  const wrap = document.createElement("div");
  wrap.className = "grid";

  const mk = document.createElement("div");
  mk.className = "card";
  mk.innerHTML = `
    <h2>Prototyp (Übungs-Audit) anlegen</h2>
    <div class="muted">
      Nur Test-/Wettbewerbsdomains oder eigene Projekte.
      Kein echter Kunde, keine Rechnung, kein Lockangebot.
    </div>
    <div class="row" style="margin-top:6px">
      <div class="ctl">
        <label>Domain/Label</label>
        <input id="aDom" type="text" placeholder="z. B. demo-immo.test">
      </div>
      <div class="ctl">
        <label>Status</label>
        <select id="aStage">
          <option>NEU</option>
          <option>IN_ARBEIT</option>
          <option>ENTWURF</option>
          <option>FERTIG_INTERNAL</option>
        </select>
      </div>
    </div>
    <button class="btn ok block" id="aAdd" style="margin-top:6px">Prototyp speichern</button>
  `;
  mk.querySelector("#aAdd").onclick = () => {
    const dom = mk.querySelector("#aDom").value.trim();
    const st = mk.querySelector("#aStage").value;
    if (!dom) return;
    state.audits.unshift({
      id: id(),
      domain: dom,
      stage: st,
      created: todayISO(),
      gate: {},
      gateScore: 0,
    });
    saveState();
    render();
  };
  wrap.append(mk);

  const filter = document.createElement("div");
  filter.className = "card";
  filter.innerHTML = `
    <h2>Prototypen-Übersicht (${state.audits.length})</h2>
    <div class="row">
      <div class="ctl">
        <label>Suche</label>
        <input id="fText" type="text" placeholder="Domain / Label">
      </div>
      <div class="ctl">
        <label>Status</label>
        <select id="fStage">
          <option value="">(alle)</option>
          <option>NEU</option>
          <option>IN_ARBEIT</option>
          <option>ENTWURF</option>
          <option>FERTIG_INTERNAL</option>
        </select>
      </div>
    </div>
  `;
  wrap.append(filter);

  const listCard = document.createElement("div");
  listCard.className = "card";
  const ul = document.createElement("div");
  ul.className = "list";

  const applyFilter = () => {
    ul.innerHTML = "";
    const q = (filter.querySelector("#fText").value || "").toLowerCase();
    const st = filter.querySelector("#fStage").value;
    state.audits
      .filter((a) => (!q || a.domain.toLowerCase().includes(q)) && (!st || a.stage === st))
      .forEach((a) => {
        const linked = state.evidence.filter((ev) => ev.auditId === a.id);
        const it = document.createElement("div");
        it.className = "item";
        it.innerHTML = `
          <div style="flex:1">
            <strong>${escapeHTML(a.domain)}</strong>
            <div class="muted">
              Status: ${a.stage} · seit ${a.created}
              ${a.gateScore >= 90 ? ` · <span class="okc">Gate ${a.gateScore}</span>` : ""}
            </div>
            <details style="margin-top:4px">
              <summary>Gate-Check (Ziel ≥ 90)</summary>
              ${renderGateChecklistHTML(a)}
            </details>
            <details style="margin-top:4px">
              <summary>Evidence (${linked.length})</summary>
              <div class="muted">
                ${
                  linked.length
                    ? linked
                        .map((ev) => `- [${ev.type}] ${escapeHTML(ev.title || "(ohne Titel)")}`)
                        .join("<br>")
                    : "Noch keine Evidence verknüpft."
                }
              </div>
            </details>
          </div>
          <div class="ctl" style="width:100px">
            <label>Stage</label>
            <select class="auditStage" data-id="${a.id}">
              ${["NEU", "IN_ARBEIT", "ENTWURF", "FERTIG_INTERNAL"]
                .map((s) => `<option ${a.stage === s ? "selected" : ""}>${s}</option>`)
                .join("")}
            </select>
          </div>
          <button class="btn ghost danger" data-del="${a.id}">✕</button>
        `;
        ul.append(it);
      });
  };

  filter.querySelectorAll("input,select").forEach((n) => (n.oninput = applyFilter));
  applyFilter();

  ul.addEventListener("change", (e) => {
    const sel = e.target.closest(".auditStage");
    if (sel) {
      const a = state.audits.find((x) => x.id === sel.dataset.id);
      if (a) {
        a.stage = sel.value;
        saveState();
      }
    }
    const gateChk = e.target.closest('input[type="checkbox"][data-gate]');
    if (gateChk) {
      const aid = gateChk.dataset.aid;
      const key = gateChk.dataset.gate;
      const a = state.audits.find((x) => x.id === aid);
      if (a) {
        a.gate = a.gate || {};
        a.gate[key] = !!gateChk.checked;
        a.gateScore = calcGateScore(a.gate);
        saveState();
        render();
      }
    }
  });

  ul.addEventListener("click", (e) => {
    const del = e.target.closest("button[data-del]");
    if (del) {
      const idd = del.dataset.del;
      const i = state.audits.findIndex((x) => x.id === idd);
      if (i > -1) {
        state.audits.splice(i, 1);
        saveState();
        render();
      }
    }
  });

  listCard.append(ul);
  wrap.append(listCard);
  viewEl.append(wrap);
};

const renderGateChecklistHTML = (audit) => {
  const gate = audit.gate || {};
  return `
    <div class="list" style="margin-top:4px">
      ${GATE_ITEMS.map(
        ([k, label]) => `
          <label class="row">
            <input type="checkbox" data-gate="${k}" data-aid="${audit.id}" ${gate[k] ? "checked" : ""}>
            <span>${label}</span>
          </label>
        `
      ).join("")}
      <div class="pill" style="margin-top:4px">
        <span>Score</span><span class="tag">${audit.gateScore || 0}</span>
      </div>
    </div>
  `;
};

const renderEvidence = () => {
  const wrap = document.createElement("div");
  wrap.className = "grid";

  const form = document.createElement("div");
  form.className = "card";
  form.innerHTML = `
    <h2>Evidence hinzufügen</h2>
    <div class="row">
      <div class="ctl" style="flex:1">
        <label>Titel</label>
        <input id="evTitle" type="text" placeholder="z. B. P1 · Mobile Performance">
      </div>
      <div class="ctl">
        <label>Typ</label>
        <select id="evType">
          <option>SEO</option>
          <option>UX</option>
          <option>Tech</option>
          <option>P0</option>
          <option>P1</option>
          <option>P2</option>
        </select>
      </div>
    </div>
    <div class="ctl" style="margin-top:6px">
      <label>Audit</label>
      <select id="evAudit">
        <option value="">(optional) Prototyp verknüpfen</option>
        ${state.audits
          .map((a) => `<option value="${a.id}">${escapeHTML(a.domain)}</option>`)
          .join("")}
      </select>
    </div>
    <div class="ctl" style="margin-top:6px">
      <label>Notiz</label>
      <textarea id="evNote" placeholder="Kurzbeschreibung, Quelle, Kontext"> </textarea>
    </div>
    <button class="btn ok block" id="evAdd" style="margin-top:6px">Evidence speichern</button>
  `;
  form.querySelector("#evAdd").onclick = () => {
    const title = form.querySelector("#evTitle").value.trim();
    const type = form.querySelector("#evType").value;
    const auditId = form.querySelector("#evAudit").value || null;
    const note = form.querySelector("#evNote").value.trim();
    if (!title && !note) return;
    state.evidence.unshift({
      id: id(),
      title,
      type,
      auditId,
      note,
      created: todayISO(),
    });
    saveState();
    render();
  };
  wrap.append(form);

  const list = document.createElement("div");
  list.className = "card";
  list.innerHTML = `<h2>Evidence (${state.evidence.length})</h2>`;
  const ul = document.createElement("div");
  ul.className = "list";
  state.evidence.forEach((ev) => {
    const it = document.createElement("div");
    it.className = "item";
    const audit = ev.auditId && state.audits.find((a) => a.id === ev.auditId);
    it.innerHTML = `
      <div style="flex:1">
        <strong>[${ev.type}] ${escapeHTML(ev.title || "(ohne Titel)")}</strong>
        <div class="muted">${ev.created}${audit ? ` · ${escapeHTML(audit.domain)}` : ""}</div>
        ${ev.note ? `<div class="muted" style="margin-top:4px">${escapeHTML(ev.note)}</div>` : ""}
      </div>
      <button class="btn ghost danger" data-del="${ev.id}">✕</button>
    `;
    ul.append(it);
  });
  ul.addEventListener("click", (e) => {
    const del = e.target.closest("button[data-del]");
    if (!del) return;
    const idd = del.dataset.del;
    const idx = state.evidence.findIndex((x) => x.id === idd);
    if (idx > -1) {
      state.evidence.splice(idx, 1);
      saveState();
      render();
    }
  });
  list.append(ul);
  wrap.append(list);

  viewEl.append(wrap);
};

const renderOutreach = () => {
  const wrap = document.createElement("div");
  wrap.className = "grid";

  const card = document.createElement("div");
  card.className = "card";
  card.innerHTML = `
    <h2>Explorationskontakte</h2>
    <div class="muted">Nur mit dokumentierter Startklausel weiterziehen.</div>
    <div class="row" style="margin-top:6px">
      <div class="ctl">
        <label>Name</label>
        <input id="ocName" type="text" placeholder="z. B. Lisa · Netzwerk">
      </div>
      <div class="ctl">
        <label>Kanal</label>
        <input id="ocChannel" type="text" placeholder="LinkedIn / Forum / Meetup">
      </div>
    </div>
    <div class="row" style="margin-top:6px">
      <div class="ctl" style="flex:1">
        <label>Handle / Referenz</label>
        <input id="ocHandle" type="text" placeholder="@handle oder URL">
      </div>
      <div class="ctl" style="flex:1">
        <label>Notiz</label>
        <input id="ocNote" type="text" placeholder="z. B. Startklausel zugesagt, Q1 2025">
      </div>
    </div>
    <button class="btn primary block" id="ocAdd" style="margin-top:6px">Kontakt sichern</button>
  `;
  card.querySelector("#ocAdd").onclick = () => {
    const contact = card.querySelector("#ocName").value.trim();
    const channel = card.querySelector("#ocChannel").value.trim();
    if (!contact || !channel) return;
    const handle = card.querySelector("#ocHandle").value.trim();
    const notes = card.querySelector("#ocNote").value.trim();
    state.outreach.unshift({
      id: id(),
      contact,
      channel,
      handle,
      notes,
      status: "EXPLORATION",
      created: todayISO(),
    });
    state.daily.outreachToday = (state.daily.outreachToday || 0) + 1;
    saveState();
    render();
  };
  wrap.append(card);

  const list = document.createElement("div");
  list.className = "card";
  list.innerHTML = `<h2>Kontakte (${state.outreach.length})</h2>`;
  const ul = document.createElement("div");
  ul.className = "list";
  state.outreach.forEach((o) => {
    const col = o.status === "EXPLORATION (Klausel OK)" ? "okc" : "warnc";
    const it = document.createElement("div");
    it.className = "item";
    it.innerHTML = `
      <div style="flex:1">
        <strong>${escapeHTML(o.contact)}</strong>
        <div class="muted">
          ${escapeHTML(o.channel)}${o.handle ? " · " + escapeHTML(o.handle) : ""}
          · ${o.created}
          · <span class="${col}">${o.status}</span>
        </div>
        ${o.notes ? `<div class="muted" style="margin-top:3px">${escapeHTML(o.notes)}</div>` : ""}
      </div>
      <div class="row">
        <button class="btn ghost ok" data-act="OK" data-id="${o.id}">Klausel OK</button>
        <button class="btn ghost danger" data-act="DEL" data-id="${o.id}">✕</button>
      </div>
    `;
    ul.append(it);
  });
  ul.addEventListener("click", (e) => {
    const b = e.target.closest("button[data-act]");
    if (!b) return;
    const idd = b.dataset.id;
    const o = state.outreach.find((x) => x.id === idd);
    if (!o) return;
    if (b.dataset.act === "OK") o.status = "EXPLORATION (Klausel OK)";
    if (b.dataset.act === "DEL") {
      const i = state.outreach.findIndex((x) => x.id === idd);
      if (i > -1) state.outreach.splice(i, 1);
    }
    saveState();
    render();
  });
  list.append(ul);
  wrap.append(list);

  viewEl.append(wrap);
};

const renderManifest = () => {
  const d = document.createElement("div");
  d.className = "card roadmap-content";
  const { total, passed } = phoenixDayInfo();
  const p = state.phoenix;
  d.innerHTML = `
    <h2>PHOENIX 80-Tage-Manifest</h2>
    <p class="muted">
      Dieses Manifest verknüpft deinen real laufenden Mini-Betrieb mit dem 80-Tage-Fenster:
      Ziel: prüffähige Entschuldung, stabiles Leistungsbild, startklares Business-Modell.
    </p>
    <h3>Phase 1 · Fundament (Tag 1–14)</h3>
    <ul>
      <li>Zeitjournal täglich führen (aktive Minuten + 3-Log-System).</li>
      <li>Kurz-EÜR Rückblick 2023–heute erstellen (Mini-Betrieb).</li>
      <li>DRV informieren & Unterlagen in Arbeit anzeigen.</li>
    </ul>
    <h3>Phase 2 · Kalte Werkstatt & Mini-Betrieb (Tag 15–40)</h3>
    <ul>
      <li>2–5 Prototypen-Audits intern erstellen (Fokus Immobilien/SEO).</li>
      <li>Evidence-Datenbank füllen (P0/P1/UX/Tech-Befunde).</li>
      <li>Mini-Betrieb klein, transparent, lückenlos verbucht.</li>
    </ul>
    <h3>Phase 3 · Schnitt & Routenentscheidung (Tag 41–65)</h3>
    <ul>
      <li>Route A: Vorbereitung des formellen Cuts der Selbstständigkeit, um Verbraucherinsolvenz zu ermöglichen.</li>
      <li>Route B: Vorbereitung einer Regelinsolvenz mit sauberem Freigabeantrag nach § 35 Abs. 2 InsO.</li>
      <li>Tag-X-Paket-Entwurf mit Prototypen & Evidence.</li>
    </ul>
    <h3>Phase 4 · Versiegelung (Tag 66–80)</h3>
    <ul>
      <li>Vollständige Antragsunterlagen fertigstellen.</li>
      <li>Tag-X-Paket finalisieren (noch nicht versenden ohne Beratung).</li>
      <li>JSON-Backup exportieren und sicher ablegen.</li>
    </ul>
    <h3>Live-Status (autonom aus deinen Eingaben)</h3>
    <ul>
      <li>Zeitjournal-Tage: <strong>${new Set(state.logbook.map((l) => l.date)).size}</strong> / 80</li>
      <li>Fokus-Sprints: <strong>${state.home.kpis.sprints || 0}</strong></li>
      <li>FERTIG_INTERNAL Prototypen: <strong>${state.audits.filter((a) => a.stage === "FERTIG_INTERNAL").length}</strong> (Ziel 3–5)</li>
      <li>Explorationskontakte mit Klausel OK: <strong>${state.outreach.filter((o) => o.status === "EXPLORATION (Klausel OK)").length}</strong></li>
    </ul>
    <p class="muted">
      Tag ${passed}/${total}: Dieses Dashboard ersetzt kein Gericht, aber es liefert deiner anwaltlichen Vertretung, DRV und dem Insolvenzverwalter eine saubere, konsistente Geschichte.
    </p>
  `;
  viewEl.append(d);
};

const renderGenerator = () => {
  const d = document.createElement("div");
  d.className = "grid";

  const card = document.createElement("div");
  card.className = "card";
  card.innerHTML = `
    <h2>⚡ Beweismittel-Generator</h2>
    <p class="muted">
      Aus deinen Daten entstehen drei Texte:
      Tag-X-Paket für den Insolvenzverwalter,
      Leistungsbild für die DRV,
      Kurz-EÜR-Report für Steuer/Anwalt.
      Alles offline, du kopierst nur Inhalte in deine finalen PDFs.
    </p>
    <h3>1 · Tag-X-Paket (IV / § 35 Abs. 2 InsO)</h3>
    <button class="btn primary block" id="gIV">Entwurf erzeugen</button>
    <h3 style="margin-top:10px;">2 · DRV-Leistungsbild</h3>
    <div class="row">
      <div class="ctl">
        <label>Zeitraum</label>
        <select id="gDRVRange">
          <option value="30">Letzte 30 Tage</option>
          <option value="60">Letzte 60 Tage</option>
          <option value="all">Komplette PHOENIX-Phase</option>
        </select>
      </div>
      <button class="btn ok block" id="gDRV">Report erzeugen</button>
    </div>
    <h3 style="margin-top:10px;">3 · Kurz-EÜR / Hinzuverdienst</h3>
    <button class="btn block" id="gEUE">Kurz-EÜR erzeugen</button>
    <div class="template-box" id="genOutBox" style="display:none;">
      <strong id="genTitle">Ergebnis</strong>
      <button class="btn block copy-btn" data-copy-target="genOut">Inhalt kopieren</button>
      <pre id="genOut"></pre>
    </div>
  `;
  d.append(card);
  viewEl.append(d);

  const outBox = el("#genOutBox");
  const outTit = el("#genTitle");
  const out = el("#genOut");

  el("#gIV").onclick = () => {
    archiveCurrentDaily();
    const audits = state.audits.filter((a) => a.stage === "FERTIG_INTERNAL");
    const outreach = state.outreach.filter((o) => o.status === "EXPLORATION (Klausel OK)");
    const execLogs = filterLogsByRange([...state.logbook, state.daily], "30").filter((l) => (l.execLog || "").trim());
    const route = state.phoenix.targetRoute;

    let s = "";
    s += "ENTWURF · TAG-X-PAKET / ANTRAG AUF FREIGABE SELBSTSTÄNDIGER TÄTIGKEIT (§ 35 ABS. 2 INSO)\n";
    s += "-------------------------------------------------------------------------------\n\n";
    s += "1. Ausgangslage\n";
    s += "- Laufender Mini-Betrieb (freiberuflich, geringes Volumen, vollständig erklärt).\n";
    s += "- Geplante Entschuldung im Rahmen " +
      (route === "A"
        ? "einer Verbraucherinsolvenz nach Einstellung der selbstständigen Tätigkeit.\n"
        : "einer Regelinsolvenz mit beantragter Freigabe meiner selbstständigen Tätigkeit.\n");
    s += "- Parallel: Aufbau eines klar umrissenen Angebots 'Mobile-First SEO Audits (Immobilien)'.\n\n";

    s += "2. Prototypen (interne Referenzen)\n";
    if (!audits.length) {
      s += "- Noch keine FERTIG_INTERNAL-Prototypen dokumentiert.\n";
    } else {
      audits.forEach((a, i) => {
        s += `- Prototyp ${i + 1}: ${a.domain} (FERTIG_INTERNAL, Gate-Score ${a.gateScore || 0})\n`;
        const evs = state.evidence.filter((ev) => ev.auditId === a.id);
        if (evs.length) {
          s += `  Evidence-Auszug:\n`;
          evs.slice(0, 5).forEach((ev) => {
            s += `  · [${ev.type}] ${ev.title}: ${ev.note.slice(0, 120)}\n`;
          });
        }
      });
    }
    s += "\n3. Markt-Pipeline (unverbindlich)\n";
    if (!outreach.length) {
      s += "- Keine verbindlichen Voranfragen. Nur Konzeptniveau.\n";
    } else {
      outreach.forEach((o) => {
        s += `- ${o.contact} (${o.channel}): Explorationskontakt mit dokumentierter Startklausel.\n`;
      });
    }

    s += "\n4. Compliance-Profil\n";
    s += "- Keine Rechnungsstellung für Audits innerhalb des PHOENIX-Zeitraums ohne Freigabe.\n";
    s += "- Mini-Betrieb über EÜR und Logbuch dokumentiert.\n";
    s += "- Executive-Logs (Auszug letzte 30 Tage) belegen Vorbereitung statt faktischer Vollausübung.\n\n";
    execLogs.forEach((l) => {
      s += `  ${l.date}: ${l.execLog}\n`;
    });

    s += "\n5. Vorschlag Abführungsmodell (Platzhalter)\n";
    s += "- Nach Erreichen stabiler Einkünfte aus freigegebener Tätigkeit: feste monatliche Quote,\n";
    s += "  orientiert an vergleichbaren Teilzeitnettoeinkommen.\n";

    outTit.textContent = "Tag-X-Paket · Rohentwurf für Rechtsberatung / Insolvenzverwalter";
    out.textContent = s;
    outBox.style.display = "block";
  };

  el("#gDRV").onclick = () => {
    archiveCurrentDaily();
    const range = el("#gDRVRange").value;
    const logsRaw = [...state.logbook];
    const idx = logsRaw.findIndex((l) => l.date === state.daily.date);
    if (idx > -1) logsRaw[idx] = JSON.parse(JSON.stringify(state.daily));
    else logsRaw.push(JSON.parse(JSON.stringify(state.daily)));

    const logs = filterLogsByRange(logsRaw, range);
    let s = "";
    s += "PROTOKOLL · LEISTUNGSBILD (INTERNER ENTWURF FÜR DRV-KOMMUNIKATION)\n";
    s += "-----------------------------------------------------------------\n";
    s += `Zeitraum: ${range === "all" ? "gesamtes PHOENIX-Fenster" : `letzte ${range} Tage`}\n\n`;
    let total = 0;
    if (!logs.length) {
      s += "(Keine Einträge im gewählten Zeitraum.)\n";
    } else {
      logs.forEach((l) => {
        total += l.activeMinutes || 0;
        s += `Datum: ${l.date}\n`;
        s += `Arbeitsähnliche Zeit: ${l.activeMinutes || 0} Min\n`;
        s += `Belastungsnotiz: ${l.healthLog || "(keine notiert)"}\n`;
        s += "--------------------------------------------------\n";
      });
    }
    const avg = logs.length ? Math.round(total / logs.length) : 0;
    s += `\nDurchschnittliche arbeitsähnliche Zeit: ${avg} Min/Tag.\n`;
    s += "Hinweis: Entwurf. Konkrete Bewertung und Grenzen immer mit der DRV abstimmen.\n";

    outTit.textContent = "Leistungsbild-Protokoll · Entwurf für DRV";
    out.textContent = s;
    outBox.style.display = "block";
  };

  el("#gEUE").onclick = () => {
    archiveCurrentDaily();
    const byYear = {};
    state.finances.entries.forEach((e) => {
      const y = e.year || (e.date || "").slice(0, 4);
      if (!byYear[y]) byYear[y] = { ein: 0, aus: 0 };
      if (e.type === "EIN") byYear[y].ein += Number(e.amount) || 0;
      if (e.type === "AUS") byYear[y].aus += Number(e.amount) || 0;
    });
    let s = "";
    s += "KURZ-EÜR (INTERNER AUSZUG)\n";
    s += "--------------------------\n\n";
    const years = Object.keys(byYear).sort();
    if (!years.length) {
      s += "(Noch keine Einträge vorhanden.)\n";
    } else {
      years.forEach((y) => {
        const ein = byYear[y].ein;
        const aus = byYear[y].aus;
        s += `Jahr ${y}:\n`;
        s += `  Einnahmen: ${ein.toFixed(2)} €\n`;
        s += `  Ausgaben:  ${aus.toFixed(2)} €\n`;
        s += `  Gewinn:    ${(ein - aus).toFixed(2)} €\n\n`;
      });
    }
    s += "Hinweis: Werte dienen als Basis für die offizielle Anlage S/EÜR in ELSTER.\n";

    outTit.textContent = "Kurz-EÜR · Grundlage für Steuer / Beratung";
    out.textContent = s;
    outBox.style.display = "block";
  };
};

const exportJSON = () => {
  archiveCurrentDaily();
  saveState();
  const blob = new Blob([exportStateJSON()], { type: "application/json" });
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  a.download = `${todayISO()}_PHOENIX_v12_backup.json`;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(a.href);
  setStatus("Export erstellt");
};

export const handleImportData = (data) => {
  const merged = { ...seedState(), ...data, version: VERSION };
  setState(merged);
  setStatus("Import ok");
  if (typeof window !== "undefined") {
    window.location.reload();
  }
};

export const getImportInput = () => importInput;
