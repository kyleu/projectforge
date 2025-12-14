import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import eslintConfigPrettier from "eslint-config-prettier";
import path from "node:path";
import { fileURLToPath } from "node:url";

const tsconfigRootDir = path.dirname(fileURLToPath(import.meta.url));

const tsFiles = ["**/*.ts", "**/*.tsx"];
const ignores = ["build.js", "eslint.config.js", "jest.config.js"];

export default tseslint.config(
  { ignores },
  eslint.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  ...tseslint.configs.stylisticTypeChecked,
  {
    files: tsFiles,
    languageOptions: {
      ecmaVersion: 2021,
      sourceType: "module",
      parserOptions: {
        project: "./tsconfig.json",
        tsconfigRootDir: tsconfigRootDir
      }
    },
    rules: {
      "max-lines": ["error"],
      "max-lines-per-function": ["error", { max: 100 }],
      "@typescript-eslint/no-unnecessary-condition": "off"
    }
  },
  eslintConfigPrettier
);
