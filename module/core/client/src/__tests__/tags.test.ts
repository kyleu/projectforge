/* @vitest-environment happy-dom */
import { afterEach, describe, expect, it } from "vitest";
import { tagsWire } from "../tags";
import { requireNonNull } from "./testUtils";

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

    const input = requireNonNull(wrapper.querySelector<HTMLInputElement>("input.result"), "result input");
    expect(input.style.display).toBe("none");

    const items = wrapper.querySelectorAll<HTMLElement>(".tags .item");
    expect(items).toHaveLength(2);
    const firstItem = requireNonNull(items.item(0), "first tag");
    const secondItem = requireNonNull(items.item(1), "second tag");
    const firstValue = requireNonNull(firstItem.querySelector<HTMLElement>(".value"), "first tag value");
    const secondValue = requireNonNull(secondItem.querySelector<HTMLElement>(".value"), "second tag value");
    expect(firstValue.innerText).toBe("alpha");
    expect(secondValue.innerText).toBe("beta");

    const add = requireNonNull(wrapper.querySelector<HTMLElement>(".add-item"), "add button");
    const clear = requireNonNull(wrapper.querySelector<HTMLElement>(".clear"), "clear");
    expect(add).not.toBeNull();
    expect(clear).not.toBeNull();
    expect(add.nextElementSibling).toBe(clear);
  });

  it("removes tags and updates hidden input", () => {
    const wrapper = buildEditor();
    tagsWire(wrapper);

    const close = requireNonNull(wrapper.querySelector<HTMLElement>(".item .close"), "close button");
    close.click();

    const input = requireNonNull(wrapper.querySelector<HTMLInputElement>("input.result"), "result input");
    expect(input.value).toBe("beta");
    expect(wrapper.querySelectorAll(".tags .item")).toHaveLength(1);
  });

  it("edits tags and persists changes", () => {
    const wrapper = buildEditor();
    tagsWire(wrapper);

    const firstItem = requireNonNull(wrapper.querySelector<HTMLElement>(".tags .item"), "first tag");
    const value = requireNonNull(firstItem.querySelector<HTMLElement>(".value"), "tag value");
    const editor = requireNonNull(firstItem.querySelector<HTMLInputElement>(".editor"), "tag editor");

    value.click();
    expect(editor.value).toBe("alpha");
    expect(value.style.display).toBe("none");
    expect(editor.style.display).toBe("block");

    editor.value = "gamma";
    editor.dispatchEvent(new Event("blur"));

    expect(value.innerText).toBe("gamma");
    const input = requireNonNull(wrapper.querySelector<HTMLInputElement>("input.result"), "result input");
    expect(input.value).toBe("gamma, beta");
  });

  it("adds new tags via the add button", () => {
    const wrapper = buildEditor();
    tagsWire(wrapper);

    const before = wrapper.querySelectorAll(".tags .item").length;
    const add = requireNonNull(wrapper.querySelector<HTMLElement>(".add-item"), "add button");
    add.click();

    const after = wrapper.querySelectorAll(".tags .item").length;
    expect(after).toBe(before + 1);

    const newItem = requireNonNull(wrapper.querySelectorAll<HTMLElement>(".tags .item").item(after - 1), "new tag");
    const value = requireNonNull(newItem.querySelector<HTMLElement>(".value"), "new tag value");
    const editor = requireNonNull(newItem.querySelector<HTMLInputElement>(".editor"), "new tag editor");
    expect(value.style.display).toBe("none");
    expect(editor.style.display).toBe("block");
  });
});
