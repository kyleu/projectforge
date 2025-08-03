import { Notation } from "../notation";
import { Numeric } from "../numeric";
import { abbreviateStandard, formatMantissaBaseTen, formatMantissaWithExponent } from "../utils";

export class StandardNotation extends Notation {
  public override readonly name = "Standard";

  public override formatNumeric(v: Numeric, places: number, placesExponent: number): string {
    const f = formatMantissaWithExponent(formatMantissaBaseTen, abbreviateStandard, 1000, 1, undefined, " ", true);
    return f(v, places, placesExponent);
  }
}
