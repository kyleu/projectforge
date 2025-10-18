/* eslint-disable max-lines */
import { Numeric } from "../numeric";
import { toFixedEngineering } from "../utils";
import { EngineeringNotation } from "./engineering";

const units = [
  "zero",
  "one",
  "two",
  "three",
  "four",
  "five",
  "six",
  "seven",
  "eight",
  "nine",
  "ten",
  "eleven",
  "twelve",
  "thirteen",
  "fourteen",
  "fifteen",
  "sixteen",
  "seventeen",
  "eighteen",
  "nineteen"
];

const tens = ["", "ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"];

const prefixes = [
  ["", "un", "duo", "tre", "quattuor", "quin", "se", "septe", "octo", "nove"],
  [
    "",
    "deci",
    "viginti",
    "triginta",
    "quadraginta",
    "quinquaginta",
    "sexaginta",
    "septuaginta",
    "octoginta",
    "nonaginta"
  ],
  ["", "centi", "ducenti", "trecenti", "quadringenti", "quingenti", "sescenti", "septingenti", "octingenti", "nongenti"]
];

const prefixes2 = [
  "",
  "milli-",
  "micro-",
  "nano-",
  "pico-",
  "femto-",
  "atto-",
  "zepto-",
  "yocto-",
  "xono-",
  "veco-",
  "meco-",
  "dueco-",
  "treco-",
  "tetreco-",
  "penteco-",
  "hexeco-",
  "hepteco-",
  "octeco-",
  "enneco-",
  "icoso-",
  "meicoso-",
  "dueicoso-",
  "trioicoso-",
  "tetreicoso-",
  "penteicoso-",
  "hexeicoso-",
  "hepteicoso-",
  "octeicoso-",
  "enneicoso-",
  "triaconto-",
  "metriaconto-",
  "duetriaconto-",
  "triotriaconto-",
  "tetretriaconto-",
  "pentetriaconto-",
  "hexetriaconto-",
  "heptetriaconto-",
  "octtriaconto-",
  "ennetriaconto-",
  "tetraconto-",
  "metetraconto-",
  "duetetraconto-",
  "triotetraconto-",
  "tetretetraconto-",
  "pentetetraconto-",
  "hexetetraconto-",
  "heptetetraconto-",
  "octetetraconto-",
  "ennetetraconto-",
  "pentaconto-",
  "mepentaconto-",
  "duepentaconto-",
  "triopentaconto-",
  "tetrepentaconto-",
  "pentepentaconto-",
  "hexepentaconto-",
  "heptepentaconto-",
  "octepentaconto-",
  "ennepentaconto-",
  "hexaconto-",
  "mehexaconto-",
  "duehexaconto-",
  "triohexaconto-",
  "tetrehexaconto-",
  "pentehexaconto-",
  "hexehexaconto-",
  "heptehexaconto-",
  "octehexaconto-",
  "ennehexaconto-",
  "heptaconto-",
  "meheptaconto-",
  "dueheptaconto-",
  "trioheptaconto-",
  "tetreheptaconto-",
  "penteheptaconto-",
  "hexeheptaconto-",
  "hepteheptaconto-",
  "octeheptaconto-",
  "enneheptaconto-",
  "octaconto-",
  "meoctaconto-",
  "dueoctaconto-",
  "triooctaconto-",
  "tetreoctaconto-",
  "penteoctaconto-",
  "hexeoctaconto-",
  "hepteoctaconto-",
  "octeoctaconto-",
  "enneoctaconto-",
  "ennaconto-",
  "meennaconto-",
  "dueeennaconto-",
  "trioennaconto-",
  "tetreennaconto-",
  "penteennaconto-",
  "hexeennaconto-",
  "hepteennaconto-",
  "octeennaconto-",
  "enneennaconto-",
  "hecto-",
  "mehecto-",
  "duehecto-"
];

const prefixCO = ["c", "o"];
const prefixDCTQS = ["d", "c", "t", "q", "s"];
const prefixSepteNove = ["septe", "nove"];
const prefixTreOrSe = ["tre", "se"];
const prefixVO = ["v", "o"];
const prefixVTQ = ["v", "t", "q"];
const smallPrefixes = ["", "thousand", "million", "billion", "trillion"];

export class EnglishNotation extends EngineeringNotation {
  public override get name(): string {
    return "English";
  }

  public override get negativeInfinite(): string {
    return "an infinitely large negative number";
  }

  public override infinite = "an infinitely large positive number";
  public override nan = "not a number";

  public formatNegativeVerySmallNumeric(v: Numeric, places: number): string {
    return `negative one ${this.formatNumeric(v.recip(), places).replace(/ /gu, "-").replace("--", "-")}th`;
  }

  public override formatVerySmallNumeric(v: Numeric, places: number): string {
    return `one ${this.formatNumeric(v.recip(), places).replace(/ /gu, "-").replace("--", "-")}th`;
  }

  public override formatNegativeUnder1000(v: number, places: number): string {
    return `negative ${this.formatNumeric(new Numeric(v), places)}`;
  }

  public override formatUnder1000(v: number, places: number): string {
    const n = Numeric.raw(v, 0);
    return this.formatNumeric(n, places);
  }

  public override formatNegativeNumeric(v: Numeric, places: number): string {
    return `negative ${this.formatNumeric(v, places)}`;
  }

  public override formatNumeric(v: Numeric, places: number): string {
    if (v.eq(0)) {
      return "zero";
    }
    // Format in the form of "one xth" when number is less than or equal 0.001.
    if (v.lte(0.001)) {
      return this.formatVerySmallNumeric(v, places);
    }

    const engineering = toFixedEngineering(v, places);
    const precision = 10 ** -places;

    // Prevent 0.002 from being formatted as "two undefined" and alike.
    if (v.lte(0.01)) {
      return this.formatUnits(v.toNumber() + precision / 2, places);
    }

    // Calculate the actual mantissa and exponent.
    const ceiled = engineering.m + precision / 2 >= 1000;
    const mantissa = ceiled ? 1 : engineering.m + precision / 2;
    const exponent = engineering.e + (ceiled ? 1 : 0);

    const unit = this.formatUnits(mantissa, places);
    const abbreviation = this.formatPrefixes(exponent);
    return abbreviation ? `${unit} ${abbreviation}` : unit;
  }

  private formatUnits(e: number, p: number): string {
    const ans = [];
    const origin = e;
    let precision = 10 ** -p;
    // The hundred place.
    if (e >= 100) {
      const a = Math.floor(e / 100);
      ans.push(`${units[a]} hundred`);
      e -= a * 100;
    }
    // The tens and units place. Because 11-19 in English is only one word this has to be separated.
    if (e < 20) {
      if (e >= 1 && ans.length > 0) {
        ans.push("and");
      }
      const a = Math.floor(e);
      ans.push(e < 1 && origin > 1 ? "" : units[a]);
      e -= a;
    } else {
      if (ans.length > 0) {
        ans.push("and");
      }
      let a = Math.floor(e / 10);
      ans.push(tens[a]);
      e -= a * 10;
      a = Math.floor(e);
      if (a !== 0) {
        ans.push(units[a]);
        e -= a;
      }
    }
    // Places after the decimal point.
    if (e >= 10 ** -p && p > 0) {
      ans.push("point");
      let a = 0;
      while (e >= precision && a < p) {
        ans.push(units[Math.floor(e * 10)]);
        e = e * 10 - Math.floor(e * 10);
        precision *= 10;
        a++;
      }
    }
    return ans.filter((i) => i !== "").join(" ");
  }

  private formatPrefixes(e: number): string {
    e = Math.floor(e / 3) - 1;
    // Quick returns.
    if (e <= 3) {
      return smallPrefixes[e + 1];
    }
    // I don't know how to clean this please send help
    let index2 = 0;
    const prefix = [prefixes[0][e % 10]];
    while (e >= 10) {
      e = Math.floor(e / 10);
      prefix.push(prefixes[++index2 % 3][e % 10]);
    }
    index2 = Math.floor(index2 / 3);
    while (prefix.length % 3 !== 0) {
      prefix.push("");
    }
    let abbreviation = "";
    while (index2 >= 0) {
      if (
        prefix[index2 * 3] !== "un" ||
        prefix[index2 * 3 + 1] !== "" ||
        prefix[index2 * 3 + 2] !== "" ||
        index2 === 0
      ) {
        let abb2 = prefix[index2 * 3 + 1] + prefix[index2 * 3 + 2];
        // Special cases.
        if (prefixTreOrSe.includes(prefix[index2 * 3]) && prefixVTQ.includes(abb2.substring(0, 1))) {
          abb2 = `s${abb2}`;
        }
        if (prefix[index2 * 3] === "se" && prefixCO.includes(abb2.substring(0, 1))) {
          abb2 = `x${abb2}`;
        }
        if (prefixSepteNove.includes(prefix[index2 * 3]) && prefixVO.includes(abb2.substring(0, 1))) {
          abb2 = `m${abb2}`;
        }
        if (prefixSepteNove.includes(prefix[index2 * 3]) && prefixDCTQS.includes(abb2.substring(0, 1))) {
          abb2 = `n${abb2}`;
        }
        abbreviation += prefix[index2 * 3] + abb2;
      }
      if (prefix[index2 * 3] !== "" || prefix[index2 * 3 + 1] !== "" || prefix[index2 * 3 + 2] !== "") {
        abbreviation += prefixes2[index2];
      }
      index2--;
    }
    abbreviation = abbreviation.replace(/-$/u, "");
    let ret = `${abbreviation}illion`;
    ret = ret.replace("i-illion", "illion");
    ret = ret.replace("iillion", "illion");
    ret = ret.replace("aillion", "illion");
    ret = ret.replace("oillion", "illion");
    ret = ret.replace("eillion", "illion");
    ret = ret.replace("unillion", "untillion");
    ret = ret.replace("duillion", "duotillion");
    ret = ret.replace("trillion", "tretillion");
    ret = ret.replace("quattuorillion", "quadrillion");
    ret = ret.replace("quinillion", "quintillion");
    ret = ret.replace("sillion", "sextillion");
    ret = ret.replace("novillion", "nonillion");
    return ret;
  }
}
