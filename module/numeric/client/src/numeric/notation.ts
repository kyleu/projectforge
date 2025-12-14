import { Numeric, NumericSource } from "./numeric";
import { formatWithCommas, noSpecialFormatting, showCommas } from "./utils";

export const Settings = {
  exponentCommas: {
    show: false,
    min: 100000,
    max: 1000000000
  },
  exponentDefaultPlaces: 0,
  isInfinite: (x: Numeric) => (Number.isFinite(x.m) ? x.gte(Numeric.maxValue) : x.m === Infinity)
};

export abstract class Notation {
  public abstract get name(): string;

  public abstract formatNumeric(v: Numeric, places: number, placesExponent: number): string;

  public format(v: NumericSource, places = 0, placesUnder1000 = 0, placesExponent = places): string {
    if (typeof v === "number" && !Number.isFinite(v)) {
      return v === -Infinity ? this.negativeInfinite : this.infinite;
    }
    const n = new Numeric(v);
    if (isNaN(n.m)) {
      return this.nan;
    }
    if (Settings.isInfinite(n.abs())) {
      return n.sign() < 0 ? this.negativeInfinite : this.infinite;
    }
    if (n.e < -300) {
      return n.sign() < 0
        ? this.formatVerySmallNegativeNumeric(n.abs(), placesUnder1000)
        : this.formatVerySmallNumeric(n, placesUnder1000);
    }
    if (n.e < 3) {
      const number = n.toNumber();
      return number < 0
        ? this.formatNegativeUnder1000(Math.abs(number), placesUnder1000)
        : this.formatUnder1000(number, placesUnder1000);
    }
    return n.sign() < 0
      ? this.formatNegativeNumeric(n.abs(), places, placesExponent)
      : this.formatNumeric(n, places, placesExponent);
  }

  public get negativeInfinite(): string {
    return `-${this.infinite}`;
  }

  public infinite = "Infinite";
  public nan = "NaN";

  public formatVerySmallNegativeNumeric(v: Numeric, places: number): string {
    return `-${this.formatVerySmallNumeric(v, places)}`;
  }

  public formatVerySmallNumeric(v: Numeric, places: number): string {
    // We switch to very small formatting as soon as 1e-300 due to precision loss, so value.toNumber() might not be zero.
    return this.formatUnder1000(v.toNumber(), places);
  }

  public formatNegativeUnder1000(v: number, places: number): string {
    return `-${this.formatUnder1000(v, places)}`;
  }

  public formatUnder1000(v: number, places: number): string {
    return v.toFixed(places);
  }

  public formatNegativeNumeric(v: Numeric, places: number, placesExponent: number): string {
    return `-${this.formatNumeric(v, places, placesExponent)}`;
  }

  protected formatExponent(
    exponent: number,
    precision: number = Settings.exponentDefaultPlaces,
    specialFormat: (n: number, p: number) => string = (n) => n.toString(),
    largeExponentPrecision: number = Math.max(2, precision)
  ): string {
    if (noSpecialFormatting(exponent)) {
      return specialFormat(exponent, Math.max(precision, 1));
    }
    if (showCommas(exponent)) {
      return formatWithCommas(specialFormat(exponent, 0));
    }
    specialFormat(0, 0);
    return this.formatNumeric(new Numeric(exponent), largeExponentPrecision, largeExponentPrecision);
  }
}
