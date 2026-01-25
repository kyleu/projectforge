/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it, vi } from "vitest";
import { flashCreate, flashInit } from "../flash";

afterEach(() => {
  vi.useRealTimers();
  vi.clearAllTimers();
  document.body.innerHTML = "";
});

describe("flash", () => {
  it("flashCreate builds a container and message", () => {
    flashCreate("k1", "success", "Hello");

    const container = document.getElementById("flash-container");
    expect(container).not.toBeNull();

    const flash = container?.querySelector(".flash");
    expect(flash).not.toBeNull();

    const content = container?.querySelector(".content") as HTMLElement;
    expect(content.className).toContain("flash-success");
    expect(content.innerText).toBe("Hello");
  });

  it("flashCreate fades and removes messages", () => {
    vi.useFakeTimers();
    flashCreate("k2", "error", "Oops");

    const flash = document.querySelector(".flash") as HTMLElement;
    expect(flash).not.toBeNull();

    vi.advanceTimersByTime(5000);
    expect(flash.style.opacity).toBe("0");

    vi.advanceTimersByTime(500);
    expect(document.querySelector(".flash")).toBeNull();
  });

  it("flashInit returns creator and fades existing flashes", () => {
    vi.useFakeTimers();
    const container = document.createElement("div");
    container.id = "flash-container";
    const flash = document.createElement("div");
    flash.className = "flash";
    container.appendChild(flash);
    document.body.appendChild(container);

    const create = flashInit();
    expect(typeof create).toBe("function");

    vi.advanceTimersByTime(5000);
    vi.advanceTimersByTime(500);
    expect(document.querySelector(".flash")).toBeNull();
  });
});
