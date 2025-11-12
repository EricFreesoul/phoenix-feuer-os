import { saveState, state } from "./state.js";
import { render } from "./generator.js";

export const initTabs = (tabButtons = []) => {
  tabButtons.forEach((tab) => {
    tab.addEventListener("click", () => {
      state.activeTab = tab.dataset.tab;
      saveState();
      render();
    });
  });
};
