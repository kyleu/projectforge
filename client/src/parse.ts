export const Parse = {
  date: (x: unknown, dflt?: () => Date): Date => {
    const d = Parse.dateOpt(x);
    if (d !== undefined) {
      return d;
    }
    if (dflt === undefined) {
      throw new Error(`invalid date input [${x}] of type [${typeof x}]`);
    }
    return dflt();
  },

  dateOpt: (x: unknown): Date | undefined => {
    if (x instanceof Date) {
      return x;
    }
    if (typeof x === "string") {
      return new Date(x);
    }
    return undefined;
  },

  float: (x: unknown, dflt?: () => number): number => {
    const s = Parse.floatOpt(x);
    if (s !== undefined) {
      return s;
    }
    if (dflt === undefined) {
      throw new Error(`invalid float input [${x}] of type [${typeof x}]`);
    }
    return dflt();
  },

  floatOpt: (x: unknown): number | undefined => {
    if (typeof x === "number") {
      return x;
    }
    if (typeof x === "string") {
      return parseFloat(x);
    }
    return undefined;
  },

  int: (x: unknown, dflt?: () => number): number => {
    const s = Parse.intOpt(x);
    if (s !== undefined) {
      return s;
    }
    if (dflt === undefined) {
      throw new Error(`invalid integer input [${x}] of type [${typeof x}]`);
    }
    return dflt();
  },

  intOpt: (x: unknown): number | undefined => {
    if (typeof x === "number") {
      return x;
    }
    if (typeof x === "string") {
      return parseInt(x, 10);
    }
    return undefined;
  },

  obj: (x: unknown, dflt?: () => { [key: string]: unknown }): { [key: string]: unknown } => {
    const o = Parse.objOpt(x);
    if (o !== undefined) {
      return o;
    }
    if (dflt === undefined) {
      throw new Error(`invalid object input [${x}] of type [${typeof x}]`);
    }
    return dflt();
  },

  objOpt: (x: unknown): { [key: string]: unknown } | undefined => {
    if (typeof x === "object" && x !== null) {
      return x as { [key: string]: unknown };
    }
    return undefined;
  },

  string: (x: unknown, dflt?: () => string): string => {
    const s = Parse.stringOpt(x);
    if (s !== undefined) {
      return s;
    }
    if (dflt === undefined) {
      throw new Error(`invalid string input [${x}] of type [${typeof x}]`);
    }
    return dflt();
  },

  stringOpt: (x: unknown): string | undefined => {
    if (typeof x === "string") {
      return x;
    }
    return undefined;
  }
};
