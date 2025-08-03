import { Notation } from "../notation";
import { Numeric } from "../numeric";
import { formatMantissaBaseTen, formatMantissaBaseTenZero, formatMantissaWithExponent } from "../utils";

export class ScientificNotation extends Notation {
  public override get name(): string {
    return "Scientific";
  }

  public override formatNumeric(v: Numeric, places: number, placesExponent: number): string {
    const f = formatMantissaWithExponent(
      formatMantissaBaseTen,
      this.formatExponent.bind(this),
      10,
      1,
      formatMantissaBaseTenZero
    );
    return f(v, places, placesExponent);
  }
}
