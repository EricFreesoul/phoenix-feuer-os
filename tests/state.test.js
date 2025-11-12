import { describe, expect, it } from "vitest";
import {
  PHOENIX_END,
  PHOENIX_START,
  calcGateScore,
  filterLogsByRange,
  phoenixDayInfo,
} from "../assets/js/state.js";

const isoDaysAgo = (days) => {
  const date = new Date();
  date.setDate(date.getDate() - days);
  return date.toISOString().slice(0, 10);
};

describe("phoenixDayInfo", () => {
  it("returns a consistent day range", () => {
    const info = phoenixDayInfo();
    const start = new Date(`${PHOENIX_START}T00:00:00`);
    const end = new Date(`${PHOENIX_END}T00:00:00`);
    const total = Math.round((end - start) / (1000 * 60 * 60 * 24)) + 1;
    expect(info.total).toBe(total);
    expect(info.passed + info.left).toBe(info.total);
  });
});

describe("filterLogsByRange", () => {
  const logs = [
    { date: isoDaysAgo(0), activeMinutes: 10 },
    { date: isoDaysAgo(10), activeMinutes: 20 },
    { date: isoDaysAgo(40), activeMinutes: 30 },
  ];

  it("keeps only entries within the last 30 days", () => {
    const filtered = filterLogsByRange(logs, "30");
    expect(filtered.map((l) => l.date)).toEqual([logs[1].date, logs[0].date]);
  });

  it("sorts entries chronologically when requesting all logs", () => {
    const filtered = filterLogsByRange(logs, "all");
    const expected = [...logs].sort((a, b) => (a.date > b.date ? 1 : -1)).map((l) => l.date);
    expect(filtered.map((l) => l.date)).toEqual(expected);
  });
});

describe("calcGateScore", () => {
  it("calculates percentage of completed gate items", () => {
    const score = calcGateScore({ klartext: true, belege: true, struktur: true });
    expect(score).toBe(43);
  });

  it("returns 100 when every gate is fulfilled", () => {
    const allTrue = {
      klartext: true,
      belege: true,
      struktur: true,
      prio: true,
      mobil: true,
      intern: true,
      exec: true,
    };
    expect(calcGateScore(allTrue)).toBe(100);
  });
});
