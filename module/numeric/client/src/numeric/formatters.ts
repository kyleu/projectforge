import {EnglishNotation} from "./formatter/english";
import {StandardNotation} from "./formatter/standard";
import {ScientificNotation} from "./formatter/scientific";
import {EngineeringNotation} from "./formatter/engineering";
import {RomanNotation} from "./formatter/roman";

export const formatters = {
  "standard": new StandardNotation(),
  "scientific": new ScientificNotation(),
  "engineering": new EngineeringNotation(),
  "english": new EnglishNotation(),
  "roman": new RomanNotation()
};
