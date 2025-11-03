import { Numeric, NumericSource } from "./numeric";

type NumericInput = Numeric | NumericSource;

export class NumericMap implements Iterable<[string, Numeric]> {
  private readonly data: Record<string, Numeric>;

  constructor(entries: NumericMap | Record<string, NumericInput> = {}) {
    this.data = {};

    if (entries instanceof NumericMap) {
      entries.forEach((value, key) => {
        this.data[key] = value;
      });
      return;
    }

    Object.entries(entries).forEach(([key, value]) => {
      this.data[key] = NumericMap.coerce(value);
    });
  }

  static parse(obj: Record<string, unknown>): NumericMap {
    const entries: Record<string, NumericInput> = {};
    Object.entries(obj).forEach(([key, value]) => {
      entries[key] = value as NumericInput;
    });
    return new NumericMap(entries);
  }

  static fromRecord(record: Record<string, NumericInput>): NumericMap {
    return new NumericMap(record);
  }

  private static coerce(value: NumericInput): Numeric {
    return value instanceof Numeric ? value : new Numeric(value);
  }

  get(key: string): Numeric | undefined {
    return this.data[key];
  }

  getOrZero(key: string): Numeric {
    return this.get(key) ?? Numeric.zero.clone();
  }

  set(key: string, value: NumericInput): this {
    this.data[key] = NumericMap.coerce(value);
    return this;
  }

  has(key: string): boolean {
    return Object.prototype.hasOwnProperty.call(this.data, key);
  }

  keys(): string[] {
    return Object.keys(this.data);
  }

  values(): Numeric[] {
    return Object.values(this.data);
  }

  entries(): [string, Numeric][] {
    return Object.entries(this.data).map(([key, value]) => [key, value] as [string, Numeric]);
  }

  forEach(callback: (value: Numeric, key: string) => void): void {
    this.entries().forEach(([key, value]) => {
      callback(value, key);
    });
  }

  clone(): NumericMap {
    return new NumericMap(this);
  }

  lt(other: NumericMap): boolean {
    return this.entries().every(([key, value]) => value.lt(other.getOrZero(key)));
  }

  gt(other: NumericMap): boolean {
    return this.entries().every(([key, value]) => value.gt(other.getOrZero(key)));
  }

  eq(other: NumericMap): boolean {
    return this === other || this.entries().every(([key, value]) => value.eq(other.getOrZero(key)));
  }

  add(other: Numeric): NumericMap {
    return this.mapWithNumeric(other, (a, b) => a.add(b));
  }

  addMap(other: NumericMap): NumericMap {
    return this.merge(other, (a, b) => a.add(b));
  }

  sub(other: Numeric): NumericMap {
    return this.mapWithNumeric(other, (a, b) => a.sub(b));
  }

  subMap(other: NumericMap): NumericMap {
    return this.merge(other, (a, b) => a.sub(b));
  }

  mul(other: Numeric): NumericMap {
    return this.mapWithNumeric(other, (a, b) => a.mult(b));
  }

  mulMap(other: NumericMap): NumericMap {
    return this.merge(other, (a, b) => a.mult(b));
  }

  div(other: Numeric): NumericMap {
    return this.mapWithNumeric(other, (a, b) => a.div(b));
  }

  divMap(other: NumericMap): NumericMap {
    return this.merge(other, (a, b) => a.div(b));
  }

  floor(): NumericMap {
    return this.map((value) => value.floor());
  }

  ceil(): NumericMap {
    return this.map((value) => value.ceil());
  }

  min(): Numeric {
    return this.entries().reduce((min, [_, value]) => (value.lt(min) ? value : min), Numeric.positiveInfinity);
  }

  max(): Numeric {
    return this.entries().reduce((max, [_, value]) => (value.gt(max) ? value : max), Numeric.negativeInfinity);
  }

  ln(): NumericMap {
    return this.map((value) => value.ln());
  }

  map(fn: (value: Numeric) => Numeric): NumericMap {
    const entries: Record<string, Numeric> = {};
    this.entries().forEach(([key, value]) => {
      entries[key] = fn(value);
    });
    return new NumericMap(entries);
  }

  merge(other: NumericMap, fn: (a: Numeric, b: Numeric) => Numeric): NumericMap {
    const entries: Record<string, Numeric> = {};
    const keys = new Set([...this.keys(), ...other.keys()]);
    keys.forEach((key) => {
      entries[key] = fn(this.getOrZero(key), other.getOrZero(key));
    });
    return new NumericMap(entries);
  }

  private mapWithNumeric(other: Numeric, fn: (a: Numeric, b: Numeric) => Numeric): NumericMap {
    const entries: Record<string, Numeric> = {};
    this.keys().forEach((key) => {
      entries[key] = fn(this.getOrZero(key), other);
    });
    return new NumericMap(entries);
  }

  [Symbol.iterator](): IterableIterator<[string, Numeric]> {
    return this.entries()[Symbol.iterator]();
  }
}
