/* @vitest-environment happy-dom */
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { requireNonNull } from "./testUtils";

let socketLog: typeof import("../socketlog").socketLog;
let lastSocketInstance: {
  open: () => void;
  recv: (m: { channel: string; cmd: string; param: Record<string, unknown> }) => void;
  err: (svc: string, err: string) => void;
} | null = null;

vi.mock("../socket", () => {
  class MockSocket {
    open: () => void;
    recv: (m: { channel: string; cmd: string; param: Record<string, unknown> }) => void;
    err: (svc: string, err: string) => void;

    constructor(
      _debug: boolean,
      o: () => void,
      r: (m: { channel: string; cmd: string; param: Record<string, unknown> }) => void,
      e: (svc: string, err: string) => void
    ) {
      this.open = o;
      this.recv = r;
      this.err = e;
      lastSocketInstance = { open: o, recv: r, err: e };
    }
  }
  return { Socket: MockSocket };
});

beforeEach(async () => {
  lastSocketInstance = null;
  ({ socketLog } = await import("../socketlog"));
});

afterEach(() => {
  document.body.innerHTML = "";
});

describe("socketlog", () => {
  it("renders output rows and scrolls", () => {
    const content = document.createElement("div");
    content.id = "content";
    const scrollSpy = vi.fn();
    content.scrollTo = scrollSpy;
    document.body.appendChild(content);

    const parent = document.createElement("div");
    document.body.appendChild(parent);

    socketLog(false, parent, false, "/connect", []);
    expect(lastSocketInstance).not.toBeNull();

    lastSocketInstance?.recv({
      channel: "c",
      cmd: "output",
      param: { html: ["hello", "\n", "world"] }
    });

    expect(parent.children).toHaveLength(2);
    expect(parent.children[0]?.textContent).toBe("hello");
    expect(parent.children[1]?.textContent).toBe("world");
    expect(scrollSpy).toHaveBeenCalled();
  });

  it("renders terminal output into table rows", () => {
    const content = document.createElement("div");
    content.id = "content";
    const scrollSpy = vi.fn();
    content.scrollTo = scrollSpy;
    document.body.appendChild(content);

    const parent = document.createElement("tbody");
    document.body.appendChild(parent);

    socketLog(false, parent, true, "/connect", []);
    lastSocketInstance?.recv({
      channel: "c",
      cmd: "output",
      param: { html: ["x"] }
    });

    const row = requireNonNull(parent.querySelector<HTMLTableRowElement>("tr"), "table row");
    const th = requireNonNull(row.querySelector<HTMLTableCellElement>("th"), "table header");
    const td = requireNonNull(row.querySelector<HTMLTableCellElement>("td"), "table cell");

    expect(th.textContent).toBe("0");
    expect(td.textContent).toBe("x");
  });

  it("routes unknown commands to extra handlers", () => {
    const content = document.createElement("div");
    content.id = "content";
    content.scrollTo = vi.fn();
    document.body.appendChild(content);

    const parent = document.createElement("div");
    const handler = vi.fn();
    socketLog(false, parent, false, "/connect", [handler]);

    lastSocketInstance?.recv({
      channel: "c",
      cmd: "unknown",
      param: {}
    });

    expect(handler).toHaveBeenCalledTimes(1);
  });
});
