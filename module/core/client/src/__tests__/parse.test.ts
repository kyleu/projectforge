import { describe, expect, it } from "vitest";
import { Parse } from "../parse";

describe("Parse", () => {
  it("parses basic types", () => {
    expect(Parse.int("42")).toBe(42);
    expect(Parse.float("3.14")).toBeCloseTo(3.14);
    expect(Parse.string("hello")).toBe("hello");

    const obj = { ok: true };
    expect(Parse.obj(obj)).toBe(obj);

    const date = Parse.date("2020-01-01T00:00:00Z");
    expect(date).toBeInstanceOf(Date);
  });

  it("throws on invalid input when no default is provided", () => {
    expect(() => Parse.obj(123)).toThrow(/invalid object input/iu);
    expect(() => Parse.string(123)).toThrow(/invalid string input/iu);
    expect(() => Parse.date(123)).toThrow(/invalid date input/iu);
  });

  it("uses defaults when provided", () => {
    const dflt = new Date("2022-02-02T00:00:00Z");
    expect(Parse.date(123, () => dflt)).toBe(dflt);
    expect(Parse.string(123, () => "fallback")).toBe("fallback");
  });
});
