/* eslint-disable max-lines, max-lines-per-function */
import { describe, expect, test } from "vitest";
import { Numeric } from "./numeric";
import { ScientificNotation } from "./formatter/scientific";
import { EngineeringNotation } from "./formatter/engineering";
import { StandardNotation } from "./formatter/standard";
import { EnglishNotation } from "./formatter/english";

interface example {
  name: string;
  mantissa: number;
  exponent: number;
  std: string;
  sci: string;
  eng: string;
  en: string;
}

describe("NumberFormat", () => {
  const tests: example[] = [
    { name: "NaN", mantissa: NaN, exponent: 10, std: "NaN", sci: "NaN", eng: "NaN", en: "not a number" },
    {
      name: "Infinity",
      mantissa: Infinity,
      exponent: 10,
      std: "Infinite",
      sci: "Infinite",
      eng: "Infinite",
      en: "an infinitely large positive number"
    },
    {
      name: "Negative Infinity",
      mantissa: -Infinity,
      exponent: 10,
      std: "-Infinite",
      sci: "-Infinite",
      eng: "-Infinite",
      en: "an infinitely large negative number"
    },
    { name: "Zero", mantissa: 0, exponent: 1, std: "0", sci: "0", eng: "0", en: "zero" },
    { name: "Zero scaled", mantissa: 0, exponent: 100, std: "0", sci: "0", eng: "0", en: "zero" },
    { name: "Power zero", mantissa: 1.2, exponent: 0, std: "1", sci: "1", eng: "1", en: "one point two" },
    { name: "One", mantissa: 0.1, exponent: 1, std: "1", sci: "1", eng: "1", en: "one" },
    {
      name: "Five remainder",
      mantissa: 0.54321,
      exponent: 1,
      std: "5",
      sci: "5",
      eng: "5",
      en: "five point four three two"
    },

    { name: "Twelve", mantissa: 1.2, exponent: 1, std: "12", sci: "12", eng: "12", en: "twelve" },
    {
      name: "Twelve one decimal place",
      mantissa: 1.23,
      exponent: 1,
      std: "12",
      sci: "12",
      eng: "12",
      en: "twelve point three"
    },
    {
      name: "Twelve two decimal places",
      mantissa: 1.234,
      exponent: 1,
      std: "12",
      sci: "12",
      eng: "12",
      en: "twelve point three four"
    },
    {
      name: "Twelve three decimal places",
      mantissa: 1.2345,
      exponent: 1,
      std: "12",
      sci: "12",
      eng: "12",
      en: "twelve point three four five"
    },
    {
      name: "Twelve remainder",
      mantissa: 1.2345678,
      exponent: 1,
      std: "12",
      sci: "12",
      eng: "12",
      en: "twelve point three four six"
    },

    {
      name: "Hundred round",
      mantissa: 1.2,
      exponent: 2,
      std: "120",
      sci: "120",
      eng: "120",
      en: "one hundred and twenty"
    },
    {
      name: "Hundred remainder",
      mantissa: 1.234567,
      exponent: 2,
      std: "123",
      sci: "123",
      eng: "123",
      en: "one hundred and twenty three point four five seven"
    },

    {
      name: "Small number < 1000",
      mantissa: 4.567891,
      exponent: 2,
      std: "457",
      sci: "457",
      eng: "457",
      en: "four hundred and fifty six point seven eight nine"
    },
    {
      name: "Negative small number",
      mantissa: -1.2345,
      exponent: 2,
      std: "-123",
      sci: "-123",
      eng: "-123",
      en: "negative one hundred and twenty three point four five"
    },

    {
      name: "Thousand round",
      mantissa: 1,
      exponent: 3,
      std: "1.000 K",
      sci: "1.000e3",
      eng: "1.000e3",
      en: "one thousand"
    },
    {
      name: "Thousand remainder",
      mantissa: 1.23456,
      exponent: 3,
      std: "1.235 K",
      sci: "1.235e3",
      eng: "1.235e3",
      en: "one point two three five thousand"
    },

    {
      name: "Ten thousand round",
      mantissa: 1.2,
      exponent: 4,
      std: "12.000 K",
      sci: "1.200e4",
      eng: "12.000e3",
      en: "twelve thousand"
    },
    {
      name: "Ten thousand remainder",
      mantissa: 1.23456,
      exponent: 4,
      std: "12.346 K",
      sci: "1.235e4",
      eng: "12.346e3",
      en: "twelve point three four six thousand"
    },

    {
      name: "Million round",
      mantissa: 1,
      exponent: 6,
      std: "1.000 M",
      sci: "1.000e6",
      eng: "1.000e6",
      en: "one million"
    },
    {
      name: "Million remainder",
      mantissa: 1.23456,
      exponent: 6,
      std: "1.235 M",
      sci: "1.235e6",
      eng: "1.235e6",
      en: "one point two three five million"
    },
    {
      name: "Negative million",
      mantissa: -9.87654321,
      exponent: 6,
      std: "-9.877 M",
      sci: "-9.877e6",
      eng: "-9.877e6",
      en: "negative nine point eight seven seven million"
    },

    {
      name: "Billion with 1 decimal",
      mantissa: 1.5,
      exponent: 9,
      std: "1.500 B",
      sci: "1.500e9",
      eng: "1.500e9",
      en: "one point five billion"
    },
    {
      name: "Billion remainder",
      mantissa: 9.87654321,
      exponent: 9,
      std: "9.877 B",
      sci: "9.877e9",
      eng: "9.877e9",
      en: "nine point eight seven seven billion"
    },
    {
      name: "Billion large remainder",
      mantissa: -1.56789987654321,
      exponent: 10,
      std: "-15.679 B",
      sci: "-1.568e10",
      eng: "-15.679e9",
      en: "negative fifteen point six seven nine billion"
    },

    {
      name: "Trillion with no decimals",
      mantissa: 2,
      exponent: 12,
      std: "2.000 T",
      sci: "2.000e12",
      eng: "2.000e12",
      en: "two trillion"
    },

    {
      name: "Quadrillion",
      mantissa: 1.23,
      exponent: 15,
      std: "1.230 Qa",
      sci: "1.230e15",
      eng: "1.230e15",
      en: "one point two three quadrillion"
    },

    {
      name: "Decillion",
      mantissa: 1,
      exponent: 33,
      std: "1.000 Dc",
      sci: "1.000e33",
      eng: "1.000e33",
      en: "one decillion"
    },

    {
      name: "Undecillion",
      mantissa: 1.23,
      exponent: 37,
      std: "12.300 UDc",
      sci: "1.230e37",
      eng: "12.300e36",
      en: "twelve point three undecillion"
    },

    {
      name: "Quindecillion round",
      mantissa: 1,
      exponent: 50,
      std: "100.000 QtDc",
      sci: "1.000e50",
      eng: "100.000e48",
      en: "one hundred quindecillion"
    },
    {
      name: "Quindecillion remainder",
      mantissa: 1.8888888,
      exponent: 50,
      std: "188.889 QtDc",
      sci: "1.889e50",
      eng: "188.889e48",
      en: "one hundred and eighty eight point eight eight nine quindecillion"
    },

    {
      name: "Septentrigintillion round",
      mantissa: 1,
      exponent: 115,
      std: "10.000 SpTg",
      sci: "1.000e115",
      eng: "10.000e114",
      en: "ten septentrigintillion"
    },
    {
      name: "Septentrigintillion remainder",
      mantissa: 1.8888888,
      exponent: 116,
      std: "188.889 SpTg",
      sci: "1.889e116",
      eng: "188.889e114",
      en: "one hundred and eighty eight point eight eight nine septentrigintillion"
    },

    {
      name: "Novenonagintillion",
      mantissa: 1.8888888,
      exponent: 300,
      std: "1.889 NNn",
      sci: "1.889e300",
      eng: "1.889e300",
      en: "one point eight eight nine novenonagintillion"
    },

    {
      name: "Centillion round",
      mantissa: 3,
      exponent: 303,
      std: "3.000 Ce",
      sci: "3.000e303",
      eng: "3.000e303",
      en: "three centillion"
    },
    {
      name: "Uncentillion remainder",
      mantissa: 3.14159,
      exponent: 307,
      std: "31.416 UCe",
      sci: "3.142e307",
      eng: "31.416e306",
      en: "thirty one point four one six uncentillion"
    },
    {
      name: "So big",
      mantissa: 3.14159,
      exponent: 8000000,
      std: "314.159 DMC-SxSeScMI-QtSeSc",
      sci: "3e8.000e6",
      eng: "314e8.000e6",
      en: "three hundred and fourteen point one five nine duomicro-sesexagintasescentimilli-quinsexagintasescentillion"
    }
  ];

  tests.forEach(({ name, mantissa, exponent, std, sci, eng, en }) => {
    const stdFmt = new StandardNotation();
    const sciFmt = new ScientificNotation();
    const engFmt = new EngineeringNotation();
    const enFmt = new EnglishNotation();
    test(name, () => {
      const n = Numeric.from(mantissa, exponent);

      const stdRes = stdFmt.format(n, 3);
      expect(stdRes).toBe(std);

      const sciRes = sciFmt.format(n, 3);
      expect(sciRes).toBe(sci);

      const engRes = engFmt.format(n, 3);
      expect(engRes).toBe(eng);

      const enRes = enFmt.format(n, 3, 3);
      expect(enRes).toBe(en);
    });
  });
});
