/* eslint-disable max-lines */
import { Pool } from "./utils";

const EXP_LIMIT = 9e15;
const NUMBER_EXP_MAX = 308;
const NUMBER_EXP_MIN = -324;
const ROUND_TOLERANCE = 1e-10;

const powersOf10: number[] = [];
for (let i = NUMBER_EXP_MIN + 1; i <= NUMBER_EXP_MAX; i++) {
  powersOf10.push(Number("1e" + i));
}
const indexOf0InPowersOf10 = 323;

export function powerOf10(power: number): number {
  return powersOf10[power + indexOf0InPowersOf10];
}

export type MantissaExponent = {
  m: number;
  e: number;
};

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
    if (x === undefined || x === null) {
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
    // eslint-disable-next-line no-self-compare
    return this.m !== this.m;
  }

  private copyFrom(v: MantissaExponent): void {
    this.m = v.m;
    this.e = v.e ?? 0;
  }

  public clone(): Numeric {
    return Numeric.raw(this.m, this.e);
  }

  public log(base: number): number {
    return (Math.LN10 / Math.log(base)) * this.log10();
  }

  public log10(): number {
    return this.e + Math.log10(this.m);
  }

  public absLog10(): number {
    return this.e + Math.log10(Math.abs(this.m));
  }

  public pow(v: number | Numeric): Numeric {
    if (this.m === 0) {
      return this;
    }
    const numberValue = v instanceof Numeric ? v.toNumber() : v;
    const temp = this.e * numberValue;
    let newMantissa;
    if (Number.isSafeInteger(temp)) {
      newMantissa = Math.pow(this.m, numberValue);
      if (isFinite(newMantissa) && newMantissa !== 0) {
        return Numeric.from(newMantissa, temp);
      }
    }
    const newExponent = Math.trunc(temp);
    const residue = temp - newExponent;
    newMantissa = Math.pow(10, numberValue * Math.log10(this.m) + residue);
    if (isFinite(newMantissa) && newMantissa !== 0) {
      return Numeric.from(newMantissa, newExponent);
    }
    const result = Numeric.pow10(numberValue * this.absLog10()); // this is 2x faster and gives same values AFAIK
    if (this.sign() === -1) {
      if (Math.abs(numberValue % 2) === 1) {
        return result.neg();
      } else if (Math.abs(numberValue % 2) === 0) {
        return result;
      }
      return Numeric.numericNaN;
    }
    return result;
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

  public eq(v: NumericSource): boolean {
    const n = new Numeric(v);
    return this.e === n.e && this.m === n.m;
  }

  public neq(v: NumericSource): boolean {
    return !this.eq(v);
  }

  public lt(v: NumericSource): boolean {
    const n = new Numeric(v);
    if (!this.isFinite()) {
      return !n.isFinite();
    }
    if (this.m === 0) {
      return n.m > 0;
    }
    if (n.m === 0) {
      return this.m <= 0;
    }
    if (this.e === n.e) {
      return this.m < n.m;
    }
    if (this.m > 0) {
      return n.m > 0 && this.e < n.e;
    }
    return n.m > 0 || this.e > n.e;
  }
  public lte(v: NumericSource): boolean {
    return !this.gt(v);
  }

  public gt(v: NumericSource): boolean {
    const n = new Numeric(v);
    if (this.m === 0) {
      return n.m < 0;
    }
    if (n.m === 0) {
      return this.m > 0;
    }
    if (this.e === n.e) {
      return this.m > n.m;
    }
    if (this.m > 0) {
      return n.m < 0 || this.e > n.e;
    }
    return n.m < 0 && this.e < n.e;
  }

  public gte(v: NumericSource): boolean {
    return !this.lt(v);
  }

  public recip(): Numeric {
    return Numeric.raw(1 / this.m, -this.e);
  }

  public normalize(): Numeric {
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

  public mul(v: NumericSource): Numeric {
    if (typeof v === "number") {
      if (v < 1e307 && v > -1e307) {
        return Numeric.from(this.m * v, this.e);
      }
      return Numeric.from(this.m * 1e-307 * v, this.e + 307);
    }
    const n = new Numeric(v);
    return Numeric.from(this.m * n.m, this.e + n.e);
  }

  public div(v: NumericSource): Numeric {
    return this.mul(new Numeric(v).recip());
  }

  public toString(): string {
    return `n(${this.m}e${this.e})`;
  }

  private setFromNumber(v: number): Numeric {
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
    if (v.indexOf("e") !== -1) {
      const parts = v.split("e");
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
