/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it } from "vitest";
import { initForm, setSiblingToNull } from "../form";

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
    const btn = document.getElementById("btn") as HTMLElement;
    setSiblingToNull(btn);

    const input = document.getElementById("target") as HTMLInputElement;
    expect(input.value).toBe("∅");
  });

  it("setSiblingToNull throws when no associated input exists", () => {
    const el = document.createElement("div");
    expect(() => setSiblingToNull(el)).toThrow(/no associated input found/iu);
  });

  it("initForm updates selected cache when value changes", () => {
    document.body.innerHTML = `
      <form id="frm" class="editor">
        <input name="foo" value="a">
        <input name="foo--selected" type="checkbox">
      </form>
    `;
    const frm = document.getElementById("frm") as HTMLFormElement;
    initForm(frm);

    const input = frm.querySelector('input[name="foo"]') as HTMLInputElement;
    const selected = frm.querySelector('input[name="foo--selected"]') as HTMLInputElement;
    expect(selected.checked).toBe(false);

    input.value = "b";
    input.dispatchEvent(new Event("change"));
    expect(selected.checked).toBe(true);

    input.value = "a";
    input.dispatchEvent(new Event("keyup"));
    expect(selected.checked).toBe(false);
  });
});
