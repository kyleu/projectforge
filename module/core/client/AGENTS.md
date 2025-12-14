# Client (TypeScript) Agent Notes

This directory contains {{{ .Name }}}'s progressive-enhancement TypeScript and CSS.
Keep it lightweight, dependency-minimal, and resilient when expected DOM elements are absent.

## Commands

- Install: `npm ci`
- Dev (watch): `npm run dev`
- Build (typecheck + bundle): `npm run build`
- Lint: `npm run lint`
- Format: `npm run format`

## Build Outputs

- Bundles to `assets/client.js` and `assets/client.css` (and sourcemaps) via `client/build.js`.
- `tsc` is used for typechecking only (`noEmit`); do not emit JS/CSS via `tsc`.

## Code Expectations

- Progressive enhancement: everything should be safe to run on pages that don’t include a given feature.
- Avoid large frameworks and runtime dependencies; prefer small utilities and DOM APIs.
- Keep initialization centralized in `client/src/client.ts` and feature-specific logic in separate modules.
