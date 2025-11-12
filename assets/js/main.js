import { initTabs } from "./tabs.js";
import { registerStatusHandler, saveState, state } from "./state.js";
import { copyToClipboard, el, els } from "./utils.js";
import { getImportInput, handleImportData, render, setupRenderer } from "./generator.js";

const view = el("#view");
const tabs = els(".tab");
const importEl = el("#importFile");
const themeBtn = el("#themeBtn");

const setStatus = (text) => {
  const statusEl = el("#status");
  if (statusEl) statusEl.textContent = text;
};

registerStatusHandler(setStatus);

setupRenderer({ view, tabs, importEl, onStatus: setStatus });
initTabs(tabs);
initTheme();
bindImport();
bindGlobal();
render();

function initTheme() {
  document.documentElement.dataset.theme = state.theme === "light" ? "light" : "dark";
  if (themeBtn) {
    themeBtn.onclick = () => {
      state.theme = state.theme === "dark" ? "light" : "dark";
      document.documentElement.dataset.theme = state.theme;
      saveState();
    };
  }
}

function bindImport() {
  const input = getImportInput();
  if (!input) return;
  input.addEventListener("change", (e) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
      try {
        const data = JSON.parse(ev.target.result);
        handleImportData(data);
      } catch (err) {
        console.error(err);
        setStatus("Import ungültig");
      }
    };
    reader.readAsText(file);
  });
}

function bindGlobal() {
  document.addEventListener("click", (e) => {
    const copyBtn = e.target.closest(".copy-btn");
    if (copyBtn) {
      const targetId = copyBtn.dataset.copyTarget;
      const source = targetId && el(`#${targetId}`);
      if (source) {
        const text = source.textContent || source.value || "";
        copyToClipboard(text)
          .then(() => {
            setStatus("kopiert");
            setTimeout(() => setStatus("bereit"), 600);
          })
          .catch(() => setStatus("Kopieren fehlgeschlagen"));
      }
    }
  });
}

// expose helpers for potential future modules
export { setStatus };
