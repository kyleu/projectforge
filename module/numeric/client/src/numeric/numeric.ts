/* eslint-disable max-lines */
import { methods } from "./numeric-impl";
import { Pool } from "./utils";

const EXP_LIMIT = 9e15;
const NUMBER_EXP_MAX = 308;
const NUMBER_EXP_MIN = -324;
const ROUND_TOLERANCE = 1e-10;

const powersOf10: number[] = [];
for (let i = NUMBER_EXP_MIN + 1; i <= NUMBER_EXP_MAX; i++) {
  powersOf10.push(Number("1e" + i.toString()));
}
const indexOf0InPowersOf10 = 323;

export function powerOf10(power: number): number {
  return powersOf10[power + indexOf0InPowersOf10];
}

export interface MantissaExponent {
  m: number;
  e: number;
}

export type NumericSource = MantissaExponent | number | string;

export class Numeric {
  public static pool = new Pool<Numeric>();

  public static zero = Numeric.pool.getOrAddPooled(0, () => new Numeric(0));
  public static one = Numeric.pool.getOrAddPooled(1, () => new Numeric(1));
  public static negOne = Numeric.pool.getOrAddPooled(-1, () => new Numeric(-1));
  public static numericNaN = new Numeric(Number.NaN);
  public static positiveInfinity = new Numeric(Number.POSITIVE_INFINITY);
  public static negativeInfinity = new Numeric(Number.NEGATIVE_INFINITY);
  public static maxValue = Numeric.raw(1, EXP_LIMIT);
  public static minValue = Numeric.raw(1, -EXP_LIMIT);

  public static raw(mantissa: number, exponent: number) {
    if (!isFinite(mantissa) || !isFinite(exponent)) {
      if (mantissa === Infinity || exponent === Infinity) {
        return Numeric.positiveInfinity;
      }
      if (mantissa === -Infinity || exponent === -Infinity) {
        return Numeric.negativeInfinity;
      }
      return Numeric.numericNaN;
    }
    const n = new Numeric();
    n.m = mantissa;
    n.e = exponent;
    return n;
  }

  public static from(mantissa: number, exponent: number) {
    return Numeric.raw(mantissa, exponent).normalize();
  }

  public static fromNum(v: number) {
    return new Numeric().setFromNumber(v);
  }

  public static pow10(v: number): Numeric {
    if (Number.isInteger(v)) {
      return Numeric.from(1, v);
    }
    return Numeric.from(Math.pow(10, v % 1), Math.trunc(v));
  }

  public m = NaN;
  public e = NaN;

  constructor(x?: NumericSource) {
    if (x === undefined) {
      this.m = 0;
      this.e = 0;
    } else if (typeof x === "number") {
      this.setFromNumber(x);
    } else if (typeof x === "string") {
      this.setFromString(x);
    } else if (x instanceof Numeric) {
      this.copyFrom(x);
    } else if (x.m === undefined) {
      throw new Error(`unsupported Numeric source type [${JSON.stringify(x)}]`);
    } else {
      this.copyFrom(x);
    }
  }

  public isNaN(): boolean {
    return this.m !== this.m;
  }

  private copyFrom(v: MantissaExponent): void {
    this.m = v.m;
    this.e = v.e ?? 0;
  }

  public clone(): Numeric {
    return Numeric.raw(this.m, this.e);
  }

  public floor(): Numeric {
    return methods.floor(this);
  }

  public ceil(): Numeric {
    return methods.ceil(this);
  }

  public log(base: Numeric): Numeric {
    return methods.log(this, base);
  }

  public log10(): number {
    return this.e + Math.log10(this.m);
  }

  public absLog10(): number {
    return this.e + Math.log10(Math.abs(this.m));
  }

  public pow(v: number): Numeric {
    return methods.pow(this, v);
  }

  public neg(): Numeric {
    return Numeric.raw(-this.m, this.e);
  }

  public sign(): number {
    return Math.sign(this.m);
  }

  public abs(): Numeric {
    return Numeric.raw(Math.abs(this.m), this.e);
  }

  public ln(): Numeric {
    const m = this.m;
    if (m !== m) {
      return Numeric.numericNaN;
    }
    if (!isFinite(m)) {
      return m === Infinity ? Numeric.positiveInfinity : Numeric.numericNaN;
    }
    if (m <= 0) {
      return m === 0 ? Numeric.negativeInfinity : Numeric.numericNaN;
    }
    const value = Math.log(m) + this.e * Math.LN10;
    if (!isFinite(value)) {
      if (value === Infinity) {
        return Numeric.positiveInfinity;
      }
      if (value === -Infinity) {
        return Numeric.negativeInfinity;
      }
      return Numeric.numericNaN;
    }
    return Numeric.fromNum(value);
  }

  public eq(v: NumericSource): boolean {
    return methods.eq(this, new Numeric(v));
  }

  public neq(v: NumericSource): boolean {
    return !methods.eq(this, new Numeric(v));
  }

  public lt(v: NumericSource): boolean {
    return methods.lt(this, new Numeric(v));
  }
  public lte(v: NumericSource): boolean {
    return methods.lte(this, new Numeric(v));
  }

  public gt(v: NumericSource): boolean {
    return methods.gt(this, new Numeric(v));
  }

  public gte(v: NumericSource): boolean {
    return methods.gte(this, new Numeric(v));
  }

  public recip(): Numeric {
    return Numeric.raw(1 / this.m, -this.e);
  }

  public normalize(): this {
    if (this.m >= 1 && this.m < 10) {
      return this;
    }
    if (!isFinite(this.m)) {
      return this;
    }
    if (this.m === 0) {
      this.m = 0;
      this.e = 0;
      return this;
    }
    const tempExponent = Math.floor(Math.log10(Math.abs(this.m)));
    this.m = tempExponent === NUMBER_EXP_MIN ? (this.m * 10) / 1e-323 : this.m / powerOf10(tempExponent);
    this.e += tempExponent;
    return this;
  }

  public isFinite(): boolean {
    return isFinite(this.m);
  }

  public isZero(): boolean {
    return this.m === 0;
  }

  public toNumber(): number {
    if (!this.isFinite()) {
      return this.m;
    }
    if (this.e > NUMBER_EXP_MAX) {
      return this.m > 0 ? Number.POSITIVE_INFINITY : Number.NEGATIVE_INFINITY;
    }
    if (this.e < NUMBER_EXP_MIN) {
      return 0;
    }
    if (this.e === NUMBER_EXP_MIN) {
      return this.m > 0 ? 5e-324 : -5e-324;
    }
    const result = this.m * powerOf10(this.e);
    if (!isFinite(result) || this.e < 0) {
      return result;
    }
    const resultRounded = Math.round(result);
    if (Math.abs(resultRounded - result) < ROUND_TOLERANCE) {
      return resultRounded;
    }
    return result;
  }

  public toInt(): number {
    return 0 // TODO;
  }

  public add(v: NumericSource): Numeric {
    return methods.add(this, new Numeric(v));
  }

  public sub(v: NumericSource): Numeric {
    return methods.add(this, new Numeric(v).neg());
  }

  public mult(v: NumericSource): Numeric {
    return methods.mult(this, new Numeric(v));
  }

  public div(v: NumericSource): Numeric {
    return methods.div(this, new Numeric(v));
  }

  public toString(): string {
    return `n(${this.m.toString()}e${this.e.toString()})`;
  }

  private setFromNumber(v: number): this {
    if (!isFinite(v)) {
      this.m = v;
      this.e = 0;
      return this;
    }

    if (v === 0) {
      this.m = 0;
      this.e = 0;
      return this;
    }

    this.e = Math.floor(Math.log10(Math.abs(v)));
    this.m = this.e === NUMBER_EXP_MIN ? (v * 10) / 1e-323 : v / powerOf10(this.e);
    this.normalize();
    return this;
  }

  private setFromString(v: string): void {
    v = v.toLowerCase();
    if (v.includes("e")) {
      const parts = v.split("e");
      if (parts.length !== 2) {
        throw new Error("invalid [Numeric] argument: " + v);
      }
      this.m = parseFloat(parts[0]);
      if (isNaN(this.m)) {
        this.m = 1;
      }
      this.e = parseFloat(parts[1]);
      this.normalize();
      return;
    }
    if (v === "nan" || v === "NaN") {
      this.copyFrom(Numeric.numericNaN);
      return;
    }
    this.setFromNumber(parseFloat(v));
    this.normalize();
    if (this.isNaN()) {
      throw new Error("invalid [Numeric] argument: " + v);
    }
  }
}

export const maxNumeric = Numeric.raw(1, EXP_LIMIT);
export const minNumeric = Numeric.raw(1, -EXP_LIMIT);
