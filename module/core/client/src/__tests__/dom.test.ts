/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it, vi } from "vitest";
import { clear, els, opt, req, setDisplay, setHTML, setText } from "../dom";
import { requireNonNull } from "./testUtils";

afterEach(() => {
  vi.restoreAllMocks();
  document.body.innerHTML = "";
});

describe("dom helpers", () => {
  it("els returns matching elements within a context", () => {
    document.body.innerHTML = '<div class="root"><span class="item"></span><span class="item"></span></div>';
    const root = requireNonNull(document.querySelector<HTMLElement>(".root"), "root");
    const items = els<HTMLElement>(".item", root);
    expect(items).toHaveLength(2);
  });

  it("opt returns undefined or a single element and warns on multiples", () => {
    document.body.innerHTML = '<div class="one"></div><div class="two"></div><div class="two"></div>';
    const warn = vi.spyOn(console, "warn").mockImplementation(() => undefined);

    expect(opt<HTMLElement>(".missing")).toBeUndefined();
    expect(opt<HTMLElement>(".one")).toBeInstanceOf(HTMLElement);
    expect(opt<HTMLElement>(".two")).toBeUndefined();
    expect(warn).toHaveBeenCalledOnce();
  });

  it("req throws when missing", () => {
    expect(() => req(".missing")).toThrow(/no element found/iu);
  });

  it("setHTML supports selector and element inputs", () => {
    document.body.innerHTML = '<div id="slot"></div>';
    const el = setHTML("#slot", "<span>hi</span>");
    expect(el.querySelector("span")?.textContent).toBe("hi");

    const direct = document.createElement("div");
    setHTML(direct, "yo");
    expect(direct.innerHTML).toBe("yo");
  });

  it("setDisplay toggles display values", () => {
    document.body.innerHTML = '<div id="slot"></div>';
    const el = requireNonNull(document.querySelector<HTMLElement>("#slot"), "slot");

    setDisplay(el, false);
    expect(el.style.display).toBe("none");

    setDisplay("#slot", true, "flex");
    expect(el.style.display).toBe("flex");
  });

  it("setText and clear update text and HTML", () => {
    document.body.innerHTML = '<div id="slot"><span>keep</span></div>';
    const el = requireNonNull(document.querySelector<HTMLElement>("#slot"), "slot");
    setText(el, "hello");
    expect(el.innerText).toBe("hello");

    clear(el);
    expect(el.innerHTML).toBe("");
  });
});
