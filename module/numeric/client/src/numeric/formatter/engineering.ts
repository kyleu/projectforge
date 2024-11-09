import {Notation} from "../notation";
import {Numeric} from "../numeric";
import {formatMantissaBaseTen, formatMantissaBaseTenZero, formatMantissaWithExponent} from "../utils";

export class EngineeringNotation extends Notation {
  public override get name(): string {
    return "Engineering";
  }

  public override formatNumeric(v: Numeric, places: number, placesExponent: number): string {
    const f = formatMantissaWithExponent(formatMantissaBaseTen, this.formatExponent.bind(this), 10, 3, formatMantissaBaseTenZero);
    return f(v, places, placesExponent);
  }
}
