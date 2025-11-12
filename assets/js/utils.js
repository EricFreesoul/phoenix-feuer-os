export const el = (sel, root = document) => root.querySelector(sel);
export const els = (sel, root = document) => Array.from(root.querySelectorAll(sel));
export const id = () => Math.random().toString(36).slice(2, 10);
export const todayISO = () => new Date().toISOString().slice(0, 10);
export const at00 = (d) => new Date(`${d}T00:00:00`);
export const clamp = (value, min, max) => (value < min ? min : value > max ? max : value);
export const escapeHTML = (input = "") =>
  input.replace(/[&<>"']/g, (m) => ({
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#039;",
  })[m]);

export const copyToClipboard = (text) => {
  if (typeof navigator === "undefined" || !navigator.clipboard) {
    return Promise.reject(new Error("Clipboard API unavailable"));
  }
  return navigator.clipboard.writeText(text);
};
