/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it } from "vitest";
import { initForm, setSiblingToNull } from "../form";
import { requireNonNull } from "./testUtils";

afterEach(() => {
  document.body.innerHTML = "";
});

describe("form", () => {
  it("setSiblingToNull sets the associated input to null marker", () => {
    document.body.innerHTML = `
      <div id="wrap">
        <div><input id="target" value="ok"></div>
        <div><span id="btn"></span></div>
      </div>
    `;
    const btn = requireNonNull(document.querySelector<HTMLElement>("#btn"), "button");
    setSiblingToNull(btn);

    const input = requireNonNull(document.querySelector<HTMLInputElement>("#target"), "input");
    expect(input.value).toBe("∅");
  });

  it("setSiblingToNull throws when no associated input exists", () => {
    const el = document.createElement("div");
    expect(() => {
      setSiblingToNull(el);
    }).toThrow(/no associated input found/iu);
  });

  it("initForm updates selected cache when value changes", () => {
    document.body.innerHTML = `
      <form id="frm" class="editor">
        <input name="foo" value="a">
        <input name="foo--selected" type="checkbox">
      </form>
    `;
    const frm = requireNonNull(document.querySelector<HTMLFormElement>("#frm"), "form");
    initForm(frm);

    const input = requireNonNull(frm.querySelector<HTMLInputElement>('input[name="foo"]'), "foo input");
    const selected = requireNonNull(frm.querySelector<HTMLInputElement>('input[name="foo--selected"]'), "selected input");
    expect(selected.checked).toBe(false);

    input.value = "b";
    input.dispatchEvent(new Event("change"));
    expect(selected.checked).toBe(true);

    input.value = "a";
    input.dispatchEvent(new Event("keyup"));
    expect(selected.checked).toBe(false);
  });
});
