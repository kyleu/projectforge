/* eslint-disable max-lines */
import { Numeric } from "./numeric";

export const methods = {
  add: (left: Numeric, right: Numeric): Numeric => {
    if (right.isZero()) {
      return left.clone();
    }
    if (left.isZero()) {
      return right.clone();
    }
    const ret = left.clone();
    if (!isFinite(ret.m) || !isFinite(right.m)) {
      ret.m += right.m; //Infinity + -Infinity === NaN
      ret.e = isNaN(ret.m) ? NaN : Infinity;
      return ret;
    }

    const difference = ret.e - right.e;
    if (Math.abs(difference) >= 16) {
      if (difference < 0) {
        ret.m = right.m;
        ret.e = right.e;
      }
      return ret;
    }

    if (difference >= 0) {
      ret.m += right.m / 10 ** difference;
    } else {
      ret.m = right.m + ret.m * 10 ** difference;
      ret.e = right.e;
    }

    const after = Math.abs(ret.m);
    if (after === 0) {
      right.e = 0;
      return ret;
    } else if (after >= 10) {
      ret.m /= 10;
      ret.e++;
    } else if (after < 1) {
      const digits = -Math.floor(Math.log10(after));
      ret.m *= 10 ** digits;
      ret.e -= digits;
    }

    ret.m = Math.round(ret.m * 1e14) / 1e14;
    if (Math.abs(ret.m) === 10) {
      ret.m /= 10;
      ret.e++;
    }
    return ret;
  },

  mult: (left: Numeric, right: Numeric): Numeric => {
    if (left.m === 0 || right.m === 0) {
      return Numeric.zero.clone();
    }

    const ret = left.clone();
    ret.e += right.e;
    ret.m *= right.m;

    if (!isFinite(ret.e)) {
      if (ret.e === -Infinity) {
        ret.m = 0;
        ret.e = 0;
      } else {
        ret.m = ret.e === Infinity ? Infinity : NaN;
      }
      return ret;
    }
    if (Math.abs(ret.m) >= 10) {
      ret.m /= 10;
      ret.e++;
    }
    ret.m = Math.round(ret.m * 1e14) / 1e14;
    if (Math.abs(ret.m) === 10) {
      ret.m /= 10;
      ret.e++;
    }
    return ret;
  },

  div: (left: Numeric, right: Numeric): Numeric => {
    if (right.m === 0) {
      return Numeric.numericNaN.clone();
    } else if (left.m === 0) {
      return Numeric.zero.clone();
    }
    const ret = left.clone();
    ret.e -= right.e;
    ret.m /= right.m;

    if (!isFinite(ret.e)) {
      if (ret.e === -Infinity) {
        ret.m = 0;
        ret.e = 0;
      } else {
        ret.m = ret.e === Infinity ? Infinity : NaN;
      }
      return ret;
    }

    if (Math.abs(ret.m) < 1) {
      ret.m *= 10;
      ret.e--;
    }

    ret.m = Math.round(ret.m * 1e14) / 1e14;
    if (Math.abs(ret.m) === 10) {
      ret.m /= 10;
      ret.e++;
    }
    return ret;
  },

  pow: (left: Numeric, power: number): Numeric => {
    if (power === 0) {
      if (left.m === 0 || isNaN(left.m)) {
        return Numeric.numericNaN.clone();
      }
      return Numeric.one.clone();
    }
    if (left.m === 0) {
      if (power < 0) {
        return Numeric.numericNaN.clone();
      }
      return Numeric.zero.clone();
    }
    if (!isFinite(power)) {
      if (left.e === 0 && left.m === 1) {
        return left;
      }
      if (left.m < 0 || isNaN(power) || isNaN(left.m)) {
        return Numeric.numericNaN.clone();
      } else if ((power === -Infinity && left.e >= 0) || (power === Infinity && left.e < 0)) {
        return Numeric.zero.clone();
      } else {
        return Numeric.positiveInfinity.clone();
      }
    }

    const ret = left.clone();
    const negative = ret.m < 0 ? Math.abs(power) % 2 : null;
    if (negative !== null) {
      if (Math.floor(power) !== power) {
        return Numeric.numericNaN.clone();
      }
      ret.m *= -1;
    }

    const base10 = power * (Math.log10(ret.m) + ret.e);
    if (!isFinite(base10)) {
      if (base10 === -Infinity) {
        return Numeric.zero.clone();
      } else if (isNaN(ret.m)) {
        return Numeric.numericNaN.clone();
      } else {
        return negative === 1 ? Numeric.negativeInfinity.clone() : Numeric.positiveInfinity.clone();
      }
    }

    const target = Math.floor(base10);
    ret.m = 10 ** (base10 - target);
    ret.e = target;

    ret.m = Math.round(ret.m * 1e14) / 1e14;
    if (ret.m === 10) {
      ret.m = 1;
      ret.e++;
    }

    if (negative === 1) {
      ret.m *= -1;
    }
    return ret;
  },

  log: (left: Numeric, base: Numeric): Numeric => {
    if (base.m === 0 || (base.e === 0 && Math.abs(base.m) === 1)) {
      return Numeric.numericNaN.clone();
    }
    if (left.e === 0 && Math.abs(left.m) === 1) {
      if (left.m === 1) {
        return Numeric.zero.clone();
      }
      return Numeric.numericNaN.clone();
    }
    if (left.m === 0) {
      if (isNaN(base.m)) {
        return Numeric.numericNaN.clone();
      }
      return base.e < 0 ? Numeric.positiveInfinity.clone() : Numeric.negativeInfinity.clone();
    }
    if (!isFinite(base.m)) {
      return Numeric.numericNaN.clone();
    }
    if (!isFinite(left.m)) {
      if (isNaN(left.m) || left.m === -Infinity) {
        return Numeric.numericNaN.clone();
      }
      return Math.abs(base.m) < 1 ? Numeric.negativeInfinity.clone() : Numeric.positiveInfinity.clone();
    }

    const ret = left.clone();
    const negative = ret.m < 0;
    if (negative) {
      if (base.m > 0) {
        return Numeric.numericNaN.clone();
      }
      ret.m *= -1;
    }

    const tooSmall = ret.e < 0; //Minor issue with negative power
    const base10 = Math.log10(Math.abs(Math.log10(ret.m) + ret.e));
    const target = Math.floor(base10);
    ret.m = 10 ** (base10 - target);
    ret.e = target;

    if (tooSmall) {
      ret.m *= -1;
    }
    if (base.e !== 1 || base.m !== 1) {
      ret.m /= Math.log10(Math.abs(base.m)) + base.e;

      const after = Math.abs(ret.m);
      if (after < 1 || after >= 10) {
        const digits = Math.floor(Math.log10(after));
        ret.m /= 10 ** digits;
        ret.e += digits;
      }
    }

    if (base.m < 0 || negative) {
      if (ret.e < 0) {
        return Numeric.numericNaN.clone();
      }
      //Due to floats (1.1 * 100 !== 110), test is done in this way (also we assume that big numbers are never uneven)
      const test = ret.e < 16 ? Math.abs(Math.round(ret.m * 1e14) / 10 ** (14 - ret.e)) % 2 : 0;
      if (base.m < 0 && (negative ? test !== 1 : test !== 0)) {
        return Numeric.numericNaN.clone();
      }
    }

    ret.m = Math.round(ret.m * 1e14) / 1e14;
    if (Math.abs(ret.m) === 10) {
      ret.m /= 10;
      ret.e++;
    }
    return ret;
  },

  lt: (left: Numeric, right: Numeric): boolean => {
    if (left.m === 0 || right.m === 0 || left.e === right.e) {
      return left.m < right.m;
    }
    if (right.m > 0) {
      return left.m < 0 ? true : right.e > left.e;
    }
    return left.m < 0 && right.e < left.e;
  },

  lte: (left: Numeric, right: Numeric): boolean => {
    if (left.m === 0 || right.m === 0 || left.e === right.e) {
      return left.m <= right.m;
    }
    if (right.m > 0) {
      return left.m < 0 ? true : right.e > left.e;
    } //NaN safety
    return left.m < 0 && right.e < left.e;
  },

  gt: (left: Numeric, right: Numeric): boolean => {
    if (left.m === 0 || right.m === 0 || left.e === right.e) {
      return left.m > right.m;
    }
    if (left.m > 0) {
      return right.m < 0 ? true : left.e > right.e;
    } //NaN safety
    return right.m < 0 && left.e < right.e;
  },

  gte: (left: Numeric, right: Numeric): boolean => {
    if (left.m === 0 || right.m === 0 || left.e === right.e) {
      return left.m >= right.m;
    }
    if (left.m > 0) {
      return right.m < 0 ? true : left.e > right.e;
    } //NaN safety
    return right.m < 0 && left.e < right.e;
  },

  eq: (left: Numeric, right: Numeric): boolean => {
    return left.e === right.e && left.m === right.m;
  },

  neq: (left: Numeric, right: Numeric): boolean => {
    return left.e !== right.e || left.m !== right.m;
  },

  trunc: (left: Numeric): Numeric => {
    if (left.e < 14) {
      if (left.e < 0) {
        left.m = 0;
        left.e = 0;
      } else {
        left.m = Math.trunc(left.m * 10 ** left.e) / 10 ** left.e;
      }
    }
    return left;
  },

  floor: (left: Numeric): Numeric => {
    if (left.e < 14) {
      if (left.e < 0) {
        left.m = left.m < 0 ? -1 : 0;
        left.e = 0;
      } else {
        left.m = Math.floor(left.m * 10 ** left.e) / 10 ** left.e;

        if (left.m === -10) {
          left.m = -1;
          left.e++;
        }
      }
    }
    return left;
  },

  ceil: (left: Numeric): Numeric => {
    if (left.e < 14) {
      if (left.e < 0) {
        left.m = left.m < 0 ? 0 : 1;
        left.e = 0;
      } else {
        left.m = Math.ceil(left.m * 10 ** left.e) / 10 ** left.e;

        if (left.m === 10) {
          left.m = 1;
          left.e++;
        }
      }
    }
    return left;
  },

  round: (left: Numeric): Numeric => {
    if (left.e < 14) {
      if (left.e < 0) {
        left.m = left.e === -1 ? Math.round(left.m / 10) : 0;
        left.e = 0;
      } else {
        left.m = Math.round(left.m * 10 ** left.e) / 10 ** left.e;

        if (Math.abs(left.m) === 10) {
          left.m /= 10;
          left.e++;
        }
      }
    }
    return left;
  }
};
