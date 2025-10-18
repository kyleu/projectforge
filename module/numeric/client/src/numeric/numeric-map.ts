import { Numeric, NumericSource } from "./numeric";

export type NumericMap = Record<string, Numeric>;

export function parseNumericMap(obj: Record<string, unknown>): NumericMap {
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [key, new Numeric(value as NumericSource)])
  );
}
