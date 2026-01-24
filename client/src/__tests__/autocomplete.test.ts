/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it, vi } from "vitest";
import { autocompleteInit } from "../autocomplete";

afterEach(() => {
  vi.useRealTimers();
  vi.restoreAllMocks();
  document.body.innerHTML = "";
});

describe("autocomplete", () => {
  it("ignores stale responses", async () => {
    vi.useFakeTimers();

    const input = document.createElement("input");
    input.id = "test-input";
    document.body.appendChild(input);

    const pending: { resolve: (data: unknown) => void }[] = [];
    const fetchMock = vi.fn(() => {
      let resolveJson!: (data: unknown) => void;
      const jsonPromise = new Promise<unknown>((resolve) => {
        resolveJson = resolve;
      });
      pending.push({ resolve: resolveJson });
      return Promise.resolve({ json: () => jsonPromise });
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const autocomplete = autocompleteInit();
    autocomplete(
      input,
      "/search",
      "q",
      (x) => (x as { t: string }).t,
      (x) => (x as { v: string }).v
    );

    input.value = "a";
    input.dispatchEvent(new Event("input"));
    await vi.runAllTimersAsync();

    input.value = "ab";
    input.dispatchEvent(new Event("input"));
    await vi.runAllTimersAsync();

    pending[1].resolve([{ v: "ab", t: "AB" }]);
    await Promise.resolve();
    await Promise.resolve();

    pending[0].resolve([{ v: "a", t: "A" }]);
    await Promise.resolve();
    await Promise.resolve();

    const list = document.getElementById("test-input-list") as HTMLDataListElement;
    const values = Array.from(list.querySelectorAll("option")).map((option) => option.value);
    expect(values).toContain("ab");
    expect(values).not.toContain("a");
  });
});
