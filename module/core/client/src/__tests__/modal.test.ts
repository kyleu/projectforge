/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it, vi } from "vitest";
import { modalGetBody, modalGetHeader, modalGetOrCreate, modalInit, modalNew, modalSetTitle } from "../modal";

afterEach(() => {
  vi.restoreAllMocks();
  document.body.innerHTML = "";
  window.location.hash = "";
});

describe("modal", () => {
  it("modalNew builds the modal structure", () => {
    const modal = modalNew("test", "Hello");
    expect(modal.id).toBe("modal-test");
    expect(modal.classList.contains("modal")).toBe(true);

    const backdrop = modal.querySelector(".backdrop");
    const content = modal.querySelector(".modal-content");
    const header = modal.querySelector(".modal-header");
    const body = modal.querySelector(".modal-body");
    const title = modal.querySelector("h2");

    expect(backdrop).not.toBeNull();
    expect(content).not.toBeNull();
    expect(header).not.toBeNull();
    expect(body).not.toBeNull();
    expect(title?.textContent).toBe("Hello");
  });

  it("modalGetOrCreate returns existing modal and throws on invalid body", () => {
    const modal = modalNew("existing", "Title");
    const existing = modalGetOrCreate("existing", "Ignored");
    expect(existing).toBe(modal);
    expect(modalGetBody(existing)).toBeInstanceOf(HTMLElement);
    expect(modalGetHeader(existing)).toBeInstanceOf(HTMLElement);

    const bad = document.createElement("div");
    bad.id = "modal-bad";
    document.body.appendChild(bad);
    expect(() => modalGetOrCreate("bad", "Nope")).toThrow(/unable to find modal body/iu);
  });

  it("modalSetTitle updates the header title", () => {
    const modal = modalNew("title", "Old");
    modalSetTitle(modal, "New");
    expect(modal.querySelector("h2")?.textContent).toBe("New");
  });

  it("modalInit handles escape and click close", () => {
    const back = vi.spyOn(window.history, "back").mockImplementation(() => undefined);

    const backdrop = document.createElement("a");
    backdrop.className = "backdrop";
    document.body.appendChild(backdrop);

    const close = document.createElement("a");
    close.className = "modal-close";
    document.body.appendChild(close);

    modalInit();

    window.location.hash = "#modal-test";
    document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));
    expect(back).toHaveBeenCalledTimes(1);

    const clickEvent = new MouseEvent("click", { bubbles: true, cancelable: true });
    backdrop.dispatchEvent(clickEvent);
    close.dispatchEvent(new MouseEvent("click", { bubbles: true, cancelable: true }));

    expect(back).toHaveBeenCalledTimes(3);
  });
});
