import { describe, expect, it } from "vitest";
import { expandCollapse, isRecord, jsonParse, svgRef, unknownToString } from "../util";

describe("util", () => {
  it("unknownToString handles core primitives", () => {
    expect(unknownToString(undefined)).toBe("-");
    expect(unknownToString(null)).toBe("-");
    expect(unknownToString("hi")).toBe("hi");
    expect(unknownToString(123)).toBe("123");
    expect(unknownToString(true)).toBe("true");

    const sym = Symbol("sym");
    expect(unknownToString(sym)).toBe(sym.toString());
  });

  it("unknownToString handles dates, arrays, objects, and functions", () => {
    const date = new Date("2024-01-01T00:00:00Z");
    expect(unknownToString(date)).toBe(date.toISOString());

    expect(unknownToString(["a", 1])).toBe("[a, 1]");
    expect(unknownToString({ a: 1 })).toBe('{\n  "a": 1\n}');

    function namedFn() {
      return true;
    }
    expect(unknownToString(namedFn)).toBe("namedFn");
  });

  it("isRecord returns true for objects and false for arrays", () => {
    expect(isRecord({ ok: true })).toBe(true);
    expect(isRecord([])).toBe(false);
    expect(isRecord(null)).toBe(false);
    expect(isRecord(new Date())).toBe(true);
  });

  it("jsonParse returns values or null", () => {
    expect(jsonParse<{ a: number }>('{"a": 1}')).toEqual({ a: 1 });
    expect(jsonParse("{not json}")).toBeNull();
  });

  it("svgRef and expandCollapse output expected markup", () => {
    const svg = svgRef("check");
    expect(svg).toContain('class="icon"');
    expect(svg).toContain("width: 18px");
    expect(svg).toContain("height: 18px");
    expect(svg).toContain("#svg-check");

    const expanded = expandCollapse("!");
    expect(expanded.__html).toContain("#svg-right");
    expect(expanded.__html).toContain("expand-collapse");
    expect(expanded.__html).toContain("!");
  });
});
