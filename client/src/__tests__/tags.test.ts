/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it } from "vitest";
import { tagsWire } from "../tags";

function buildEditor(value = "alpha, beta") {
  const wrapper = document.createElement("div");
  wrapper.className = "tag-editor";
  wrapper.innerHTML = `
    <input class="result" value="${value}">
    <div class="tags"></div>
    <div class="clear"></div>
  `;
  document.body.appendChild(wrapper);
  return wrapper;
}

afterEach(() => {
  document.body.innerHTML = "";
});

describe("tags", () => {
  it("tagsWire renders tags and add button", () => {
    const wrapper = buildEditor("alpha, beta, , ");
    tagsWire(wrapper);

    const input = wrapper.querySelector("input.result") as HTMLInputElement;
    expect(input.style.display).toBe("none");

    const items = wrapper.querySelectorAll(".tags .item");
    expect(items).toHaveLength(2);
    expect((items[0]?.querySelector(".value") as HTMLElement).innerText).toBe("alpha");
    expect((items[1]?.querySelector(".value") as HTMLElement).innerText).toBe("beta");

    const add = wrapper.querySelector(".add-item") as HTMLElement;
    const clear = wrapper.querySelector(".clear") as HTMLElement;
    expect(add).not.toBeNull();
    expect(clear).not.toBeNull();
    expect(add.nextElementSibling).toBe(clear);
  });

  it("removes tags and updates hidden input", () => {
    const wrapper = buildEditor();
    tagsWire(wrapper);

    const close = wrapper.querySelector(".item .close") as HTMLElement;
    close.click();

    const input = wrapper.querySelector("input.result") as HTMLInputElement;
    expect(input.value).toBe("beta");
    expect(wrapper.querySelectorAll(".tags .item")).toHaveLength(1);
  });

  it("edits tags and persists changes", () => {
    const wrapper = buildEditor();
    tagsWire(wrapper);

    const firstItem = wrapper.querySelector(".tags .item") as HTMLElement;
    const value = firstItem.querySelector(".value") as HTMLElement;
    const editor = firstItem.querySelector(".editor") as HTMLInputElement;

    value.click();
    expect(editor.value).toBe("alpha");
    expect(value.style.display).toBe("none");
    expect(editor.style.display).toBe("block");

    editor.value = "gamma";
    editor.dispatchEvent(new Event("blur"));

    expect(value.innerText).toBe("gamma");
    const input = wrapper.querySelector("input.result") as HTMLInputElement;
    expect(input.value).toBe("gamma, beta");
  });

  it("adds new tags via the add button", () => {
    const wrapper = buildEditor();
    tagsWire(wrapper);

    const before = wrapper.querySelectorAll(".tags .item").length;
    const add = wrapper.querySelector(".add-item") as HTMLElement;
    add.click();

    const after = wrapper.querySelectorAll(".tags .item").length;
    expect(after).toBe(before + 1);

    const newItem = wrapper.querySelectorAll(".tags .item")[after - 1] as HTMLElement;
    const value = newItem.querySelector(".value") as HTMLElement;
    const editor = newItem.querySelector(".editor") as HTMLInputElement;
    expect(value.style.display).toBe("none");
    expect(editor.style.display).toBe("block");
  });
});
