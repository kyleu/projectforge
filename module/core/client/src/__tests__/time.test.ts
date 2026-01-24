/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it, vi } from "vitest";
import { relativeTime } from "../time";

function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  return (
    year.toString() +
    "-" +
    (month < 10 ? "0" + month.toString() : month.toString()) +
    "-" +
    (day < 10 ? "0" + day.toString() : day.toString())
  );
}

afterEach(() => {
  vi.useRealTimers();
  vi.clearAllTimers();
  document.body.innerHTML = "";
});

describe("relativeTime", () => {
  it("renders fallback for invalid timestamps", () => {
    const el = document.createElement("span");
    el.dataset.timestamp = "not-a-date";
    const result = relativeTime(el);
    expect(result).toBe("-");
    expect(el.textContent).toBe("-");
  });

  it("renders absolute dates for old timestamps", () => {
    const now = Date.now();
    const oldDate = new Date(now - 40 * 24 * 60 * 60 * 1000);
    const el = document.createElement("span");
    el.dataset.timestamp = oldDate.toISOString();
    const result = relativeTime(el);
    const expected = formatDate(new Date(el.dataset.timestamp));
    expect(result).toBe(expected);
    expect(el.textContent).toBe(expected);
  });

  it("renders relative time for recent timestamps", () => {
    vi.useFakeTimers();
    const now = new Date("2024-01-01T00:00:00Z");
    vi.setSystemTime(now);

    const el = document.createElement("span");
    document.body.appendChild(el);
    const recent = new Date(now.getTime() - 2 * 60 * 1000);
    el.dataset.timestamp = recent.toISOString();

    const result = relativeTime(el);
    expect(result).toBe("2 minutes ago");
    expect(el.textContent).toBe("2 minutes ago");
  });
});
