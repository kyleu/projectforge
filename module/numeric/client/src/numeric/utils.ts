import {Numeric} from "./numeric";
import {Settings} from "./notation";

export class Pool<T> {
  private pool = new Map<number | string, T>();

  public tryGetPooled = (v: number | string): T | undefined => {
    return this.pool.get(v);
  };

  public getOrAddPooled = (v: number | string, fn: () => T): T => {
    let pooled = this.pool.get(v);
    if (pooled === undefined) {
      pooled = Object.freeze(fn());
      this.pool.set(v, pooled);
    }

    return pooled;
  };
}

function commaSection(v: string, index: number): string {
  if (index === 0) {
    return v.slice(-3);
  }
  return v.slice(-3 * (index + 1), -3 * index);
}

function addCommas(v: string): string {
  return Array.from(Array(Math.ceil(v.length / 3))).map((_, i) => commaSection(v, i)).reverse().join(",");
}

export function formatWithCommas(v: number | string): string {
  const decimalPointSplit = v.toString().split(".");
  decimalPointSplit[0] = decimalPointSplit[0].replace(/\w+$/gu, addCommas);
  return decimalPointSplit.join(".");
}

export function fixMantissaOverflow(v: Numeric, places: number, threshold: number, powerOffset: number): Numeric {
  const pow10 = 10 ** places;
  const isOverflowing = Math.round(v.m * pow10) >= threshold * pow10;
  if (isOverflowing) {
    return Numeric.from(1, v.e + powerOffset);
  }
  return v;
}

export function toEngineering(v: Numeric): Numeric {
  const exponentOffset = v.e % 3;
  return Numeric.raw(v.m * 10 ** exponentOffset, v.e - exponentOffset);
}

export function toLongScale(v: Numeric): Numeric {
  const mod = v.e < 6 ? 3 : 6;
  const exponentOffset = v.e % mod;
  return Numeric.raw(v.m * 10 ** exponentOffset, v.e - exponentOffset);
}

export function toFixedEngineering(v: Numeric, places: number): Numeric {
  return fixMantissaOverflow(toEngineering(v), places, 1000, 3);
}

export function toFixedLongScale(v: Numeric, places: number): Numeric {
  const overflowPlaces = v.e < 6 ? 3 : 6;
  return fixMantissaOverflow(toLongScale(v), places, 10 ** overflowPlaces, overflowPlaces);
}

const SUBSCRIPT_NUMBERS = ["₀", "₁", "₂", "₃", "₄", "₅", "₆", "₇", "₈", "₉"];

export function toSubscript(v: number): string {
  return v.toFixed(0).split("").map((x) => {
    return x === "-" ? "₋" : SUBSCRIPT_NUMBERS[parseInt(x, 10)];
  }).join("");
}

const SUPERSCRIPT_NUMBERS = ["⁰", "¹", "²", "³", "⁴", "⁵", "⁶", "⁷", "⁸", "⁹"];

export function toSuperscript(v: number): string {
  return v.toFixed(0).split("").map((x) => {
    return x === "-" ? "⁻" : SUPERSCRIPT_NUMBERS[parseInt(x, 10)];
  }).join("");
}

const STANDARD_ABBREVIATIONS = [
  "K", "M", "B", "T", "Qa", "Qt", "Sx", "Sp", "Oc", "No"
];

const STANDARD_PREFIXES = [
  ["", "U", "D", "T", "Qa", "Qt", "Sx", "Sp", "O", "N"],
  ["", "Dc", "Vg", "Tg", "Qd", "Qi", "Se", "St", "Og", "Nn"],
  ["", "Ce", "Dn", "Tc", "Qe", "Qu", "Sc", "Si", "Oe", "Ne"]
];

const STANDARD_PREFIXES_2 = ["", "MI-", "MC-", "NA-", "PC-", "FM-", "AT-", "ZP-"];

export function abbreviateStandard(rawExp: number): string {
  const exp = rawExp - 1;
  if (exp === -1) {
    return "";
  }
  if (exp < STANDARD_ABBREVIATIONS.length) {
    return STANDARD_ABBREVIATIONS[exp];
  }
  const prefix = [];
  let e = exp;
  while (e > 0) {
    prefix.push(STANDARD_PREFIXES[prefix.length % 3][e % 10]);
    e = Math.floor(e / 10);
  }
  while (prefix.length % 3 !== 0) {
    prefix.push("");
  }
  let abbreviation = "";
  for (let i = prefix.length / 3 - 1; i >= 0; i--) {
    abbreviation += prefix.slice(i * 3, i * 3 + 3).join("") + STANDARD_PREFIXES_2[i];
  }
  return abbreviation.replace(/-[A-Z]{2}-/gu, "-").replace(/U([A-Z]{2}-)/gu, "$1").replace(/-$/u, "");
}

export function noSpecialFormatting(exponent: number): boolean {
  return exponent < Settings.exponentCommas.min;
}

export function showCommas(exponent: number): boolean {
  return Settings.exponentCommas.show && exponent < Settings.exponentCommas.max;
}

export function isExponentFullyShown(exponent: number): boolean {
  return noSpecialFormatting(exponent) || showCommas(exponent);
}

// eslint-disable-next-line max-params
export function formatMantissaWithExponent(
  mantissaFormatting: (n: number, precision: number) => string,
  exponentFormatting: (n: number, precision: number) => string, base: number, steps: number,
  mantissaFormattingIfExponentIsFormatted?: (n: number, precision: number) => string,
  separator: string = "e", forcePositiveExponent: boolean = false
):
    ((n: Numeric, precision: number, precisionExponent: number) => string) {
  return function(n: Numeric, precision: number, precisionExponent: number): string {
    const realBase = base ** steps;
    let exponent = Math.floor(n.log(realBase)) * steps;
    if (forcePositiveExponent) {
      exponent = Math.max(exponent, 0);
    }
    let mantissa = n.div(new Numeric(base).pow(exponent)).toNumber();
    if (!(mantissa > 1 && mantissa < realBase)) {
      const adjust = Math.floor(Math.log(mantissa) / Math.log(realBase));
      mantissa /= Math.pow(realBase, adjust);
      exponent += steps * adjust;
    }
    let m = mantissaFormatting(mantissa, precision);
    if (m === mantissaFormatting(realBase, precision)) {
      m = mantissaFormatting(1, precision);
      exponent += steps;
    }
    if (exponent === 0) {
      return m;
    }
    const e = exponentFormatting(exponent, precisionExponent);
    if (typeof mantissaFormattingIfExponentIsFormatted !== "undefined" && !isExponentFullyShown(exponent)) {
      m = mantissaFormattingIfExponentIsFormatted(mantissa, precision);
    }
    return `${m}${separator}${e}`;
  };
}

export function formatMantissaBaseTen(n: number, precision: number): string {
  return n.toFixed(Math.max(0, precision));
}

export function formatMantissaBaseTenZero(n: number): string {
  return formatMantissaBaseTen(n, 0);
}

export function formatMantissa(base: number, digits: string): ((n: number, precision: number) => string) {
  return function(n: number, precision: number): string {
    let value = Math.round(n * base ** Math.max(0, precision));
    const d = [];
    while (value > 0 || d.length === 0) {
      d.push(digits[value % base]);
      value = Math.floor(value / base);
    }
    let result = d.reverse().join("");
    if (precision > 0) {
      result = result.padStart(precision + 1, "0");
      result = `${result.slice(0, -precision)}.${result.slice(-precision)}`;
    }
    return result;
  };
}
